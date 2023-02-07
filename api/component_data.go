package api

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

type ComponentDataService struct {
	client *Client
}

type ComponentDataInitialRequest struct {
	Name             string          `json:"name"`
	Tags             []string        `json:"tags"`
	SupportedMethods []string        `json:"supportedMethods"`
	Documents        []*DocumentSpec `json:"documents"`
}

type DocumentSpec struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type ComponentDataInitialResponseRaw struct {
	Data *ComponentDataInitialResponse `json:"data,omitempty"`
}

type ComponentDataInitialResponse struct {
	Guid          string                       `json:"guid,omitempty"`
	UploadMethods []*ComponentDataUploadMethod `json:"uploadMethods,omitempty"`
}

type ComponentDataUploadMethod struct {
	Method string            `json:"method,omitempty"`
	Info   map[string]string `json:"info,omitempty"`
}

type ComponentDataCompleteRequest struct {
	Guid string `json:"guid"`
}

type ComponentDataCompleteResponseRaw struct {
	Data *ComponentDataCompleteResponse `json:"data,omitempty"`
}

type ComponentDataCompleteResponse struct {
	Guid string `json:"guid,omitempty"`
}

func (svc *ComponentDataService) UploadFiles(name string, tags []string, paths []string) (string, error) {
	initialRequest, err := buildComponentDataInitialRequest(name, tags, paths)
	if err != nil {
		return "", err
	}
	var initialResponse ComponentDataInitialResponseRaw
	err = doWithExponentialBackoffWaiting(func() error {
		return svc.client.RequestEncoderDecoder(http.MethodPost, apiV2ComponentDataRequest, initialRequest, &initialResponse)
	})
	if err != nil {
		return "", err
	}
	var chosenMethod *ComponentDataUploadMethod
	for _, method := range initialResponse.Data.UploadMethods {
		if method.Method == "AwsS3" {
			chosenMethod = method
		}
	}
	if chosenMethod == nil {
		return "", errors.New("couldn't find a supported upload method in the upload request response")
	}
	for _, path := range paths {
		err = doWithExponentialBackoffWaiting(func() error {
			return svc.putFileToS3(path, chosenMethod.Info)
		})
		if err != nil {
			return "", err
		}
	}
	completeRequest := ComponentDataCompleteRequest{
		Guid: initialResponse.Data.Guid,
	}
	var completeResponse ComponentDataCompleteResponseRaw
	err = doWithExponentialBackoffWaiting(func() error {
		return svc.client.RequestEncoderDecoder(http.MethodPost, apiV2ComponentDataComplete, completeRequest, &completeResponse)
	})
	if err != nil {
		return "", err
	}
	if initialResponse.Data.Guid != completeResponse.Data.Guid {
		return "", errors.New("expected the initial GUID and the one returned on completion to match")
	}
	return initialResponse.Data.Guid, nil
}

func buildComponentDataInitialRequest(name string, tags []string, paths []string) (*ComponentDataInitialRequest, error) {
	documents := make([]*DocumentSpec, 0, len(paths))
	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			return nil, err
		}
		documents = append(documents, &DocumentSpec{
			Name: filepath.Base(path),
			Size: info.Size(),
		})
	}
	return &ComponentDataInitialRequest{
		Name:             name,
		Tags:             tags,
		SupportedMethods: []string{"AwsS3"},
		Documents:        documents,
	}, nil
}

func (svc *ComponentDataService) putFileToS3(path string, uploadUrls map[string]string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	req, err := http.NewRequest(http.MethodPut, uploadUrls[filepath.Base(path)], file)
	if err != nil {
		return err
	}
	_, err = svc.client.Do(req)
	return err
}

func doWithExponentialBackoffWaiting(f func() error) error {
	return DoWithExponentialBackoff(f, func(x int) {
		time.Sleep(time.Duration(x) * time.Second)
	})
}

func DoWithExponentialBackoff(f func() error, wait func(x int)) error {
	err := f()
	if err == nil {
		return nil
	}
	for waitTime := 2; waitTime < 60; waitTime *= 2 {
		wait(waitTime)
		err := f()
		if err == nil {
			return nil
		}
	}
	return err
}
