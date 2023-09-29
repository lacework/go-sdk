package api

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

type ComponentDataService struct {
	client *Client
}

const URL_TYPE_DEFAULT = "Default"
const URL_TYPE_SAST_TABLES = "SastTables"

var URL_TYPES = []string{URL_TYPE_DEFAULT, URL_TYPE_SAST_TABLES}

type ComponentDataInitialRequest struct {
	Name             string          `json:"name"`
	Tags             []string        `json:"tags"`
	SupportedMethods []string        `json:"supportedMethods"`
	Documents        []*DocumentSpec `json:"documents"`
	UrlType          string          `json:"urlType"`
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
	UploadGuid string `json:"uploadGuid"`
	UrlType    string `json:"urlType"`
}

type ComponentDataCompleteResponseRaw struct {
	Data *ComponentDataCompleteResponse `json:"data,omitempty"`
}

type ComponentDataCompleteResponse struct {
	Guid string `json:"guid,omitempty"`
}

func (svc *ComponentDataService) UploadFiles(
	name string, tags []string, paths []string) (string, error) {
	return svc.doUploadFiles(name, tags, paths, URL_TYPE_DEFAULT)
}

func (svc *ComponentDataService) UploadSastTables(
	name string, paths []string) (string, error) {
	return svc.doUploadFiles(name, []string{"sast"}, paths, URL_TYPE_SAST_TABLES)
}

func (svc *ComponentDataService) doUploadFiles(
	name string, tags []string, paths []string, urlType string) (string, error) {
	var hasValidType = false
	for _, validType := range URL_TYPES {
		if urlType == validType {
			hasValidType = true
			break
		}
	}
	if !hasValidType {
		return "", errors.Errorf("Invalid URL type: (%s)", urlType)
	}
	initialRequest, err := buildComponentDataInitialRequest(name, tags, paths, urlType)
	if err != nil {
		return "", err
	}
	var initialResponse ComponentDataInitialResponseRaw
	err = doWithExponentialBackoffWaiting(func() error {
		return svc.client.RequestEncoderDecoder(http.MethodPost,
			apiV2ComponentDataRequest,
			initialRequest,
			&initialResponse,
		)
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
		UploadGuid: initialResponse.Data.Guid,
		UrlType:    urlType,
	}
	var completeResponse ComponentDataCompleteResponseRaw
	err = doWithExponentialBackoffWaiting(func() error {
		return svc.client.RequestEncoderDecoder(http.MethodPost,
			apiV2ComponentDataComplete,
			completeRequest,
			&completeResponse,
		)
	})
	if err != nil {
		return "", err
	}
	if initialResponse.Data.Guid != completeResponse.Data.Guid {
		return "", errors.New("expected the initial GUID and the one returned on completion to match")
	}
	return initialResponse.Data.Guid, nil
}

func buildComponentDataInitialRequest(
	name string, tags []string, paths []string, urlType string,
) (*ComponentDataInitialRequest, error) {
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
		UrlType:          urlType,
	}, nil
}

func (svc *ComponentDataService) putFileToS3(path string, uploadUrls map[string]string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, uploadUrls[filepath.Base(path)], bytes.NewReader(contents))
	if err != nil {
		return err
	}
	resp, err := svc.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			return errors.Errorf("Upload to S3 failed (%s): %s", resp.Status, body)
		}
		return errors.Errorf("Upload to S3 failed (%s)", resp.Status)
	}
	return nil
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
