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

type InitialResponseRaw struct {
	Data *InitialResponse `json:"data,omitempty"`
}

type InitialResponse struct {
	Guid          string          `json:"guid,omitempty"`
	UploadMethods []*UploadMethod `json:"uploadMethods,omitempty"`
}

type UploadMethod struct {
	Method string            `json:"method,omitempty"`
	Info   map[string]string `json:"info,omitempty"`
}

type CompleteRequest struct {
	Guid string `json:"guid"`
}

type CompleteResponseRaw struct {
	Data *CompleteResponse `json:"data,omitempty"`
}

type CompleteResponse struct {
	Guid string `json:"guid,omitempty"`
}

func (svc *ComponentDataService) UploadFiles(name string, tags []string, paths []string) (string, error) {
	initialRequest, err := buildInitialRequest(name, tags, paths)
	if err != nil {
		return "", err
	}
	var initialResponse InitialResponseRaw
	err = doWithExponentialBackoff(func() error {
		return svc.client.RequestEncoderDecoder(http.MethodPost, "v2/ComponentData/requestUpload", initialRequest, &initialResponse)
	})
	if err != nil {
		return "", err
	}
	var chosenMethod *UploadMethod
	for _, method := range initialResponse.Data.UploadMethods {
		if method.Method == "AwsS3" {
			chosenMethod = method
		}
	}
	if chosenMethod == nil {
		return "", errors.New("couldn't find a supported upload method in the upload request response")
	}
	for _, path := range paths {
		err = doWithExponentialBackoff(func() error {
			return putFileToS3(svc, path, chosenMethod.Info)
		})
		if err != nil {
			return "", err
		}
	}
	completeRequest := CompleteRequest{
		Guid: initialResponse.Data.Guid,
	}
	var completeResponse CompleteResponseRaw
	err = doWithExponentialBackoff(func() error {
		return svc.client.RequestEncoderDecoder(http.MethodPost, "v2/ComponentData/completeUpload", completeRequest, &completeResponse)
	})
	if err != nil {
		return "", err
	}
	if initialResponse.Data.Guid != completeResponse.Data.Guid {
		return "", errors.New("expected the initial GUID and the one returned on completion to match")
	}
	return initialResponse.Data.Guid, nil
}

func buildInitialRequest(name string, tags []string, paths []string) (*ComponentDataInitialRequest, error) {
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

func putFileToS3(svc *ComponentDataService, path string, uploadUrls map[string]string) error {
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
	if err != nil {
		return err
	}
	return nil
}

func doWithExponentialBackoff(f func() error) error {
	err := f()
	if err == nil {
		return nil
	}
	for wait := 2; wait < 60; wait *= 2 {
		time.Sleep(time.Duration(wait) * time.Second)
		err := f()
		if err == nil {
			return nil
		}
	}
	return err
}
