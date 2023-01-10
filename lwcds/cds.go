package lwcds

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type CdsInitialRequest struct {
	Version   string
	Name      string
	Tags      []string
	Documents []*DocumentSpec
}

type DocumentSpec struct {
	Name string
	Size int64
}

type CdsInitialResponse struct {
	Version       string
	Status        string
	Guid          string
	Failure       *CdsFailure
	UploadMethods []*CdsUploadMethod
}

type CdsFailure struct {
	Issue   string
	Message string
}

type CdsUploadMethod struct {
	Method string
	Info   map[string]string
}

type CdsUploadResponse struct {
	Version string
	Status  string
	Failure *CdsFailure
}

type CdsCompleteRequest struct {
	Version string
	Guid    string
}

type CdsCompleteResponse struct {
	Version string
	Status  string
	Failure *CdsFailure
}

func UploadFiles(name string, tags []string, paths []string) (string, error) {
	laceworkApi, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithToken(os.Getenv("LW_API_TOKEN")),
		api.WithApiV2(),
	)
	if err != nil {
		return "", err
	}
	initialRequest, err := BuildInitialRequest(name, tags, paths)
	if err != nil {
		return "", err
	}
	initialBytes, err := json.Marshal(initialRequest)
	if err != nil {
		return "", err
	}
	initialReq, err := laceworkApi.NewRequest("POST", "cds/requestUpload", bytes.NewBuffer(initialBytes))
	if err != nil {
		return "", err
	}
	initialResp, err := http.DefaultClient.Do(initialReq)
	if err != nil {
		return "", err
	}
	if initialResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code: %v (%s)", initialResp.StatusCode, initialResp.Status)
	}
	var initialResponse CdsInitialResponse
	err = json.NewDecoder(initialResp.Body).Decode(&initialResponse)
	if err != nil {
		return "", err
	}
	if initialResponse.Failure != nil {
		return "", errorFromFailure(initialResponse.Failure)
	}
	if initialResponse.Status != "SUCCESS" {
		return "", errors.New("upload failed for unknown reason")
	}
	var chosenMethod *CdsUploadMethod
	for _, method := range initialResponse.UploadMethods {
		if method.Method == "POST_TARBALL" {
			chosenMethod = method
		}
	}
	if chosenMethod == nil {
		return "", errors.New("couldn't find a supported upload method in the initial response")
	}
	fileReader := ReaderForPaths(paths)
	req, err := laceworkApi.NewRequest("POST", chosenMethod.Info["upload"], &fileReader)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code: %v (%s)", resp.StatusCode, resp.Status)
	}
	var uploadResponse CdsUploadResponse
	err = json.NewDecoder(resp.Body).Decode(&uploadResponse)
	if err != nil {
		return "", err
	}
	if uploadResponse.Failure != nil {
		return "", errorFromFailure(uploadResponse.Failure)
	}
	if uploadResponse.Status != "SUCCESS" {
		return "", errors.New("upload failed for unknown reason")
	}
	completeRequest := CdsCompleteRequest{
		Version: "0.1.0",
		Guid:    initialResponse.Guid,
	}
	completeBytes, err := json.Marshal(completeRequest)
	if err != nil {
		return "", err
	}
	completeReq, err := laceworkApi.NewRequest("POST", chosenMethod.Info["complete"], bytes.NewBuffer(completeBytes))
	if err != nil {
		return "", err
	}
	completeResp, err := http.DefaultClient.Do(completeReq)
	if err != nil {
		return "", err
	}
	if completeResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code: %v (%s)", completeResp.StatusCode, completeResp.Status)
	}
	var completeResponse CdsUploadResponse
	err = json.NewDecoder(resp.Body).Decode(&completeResp)
	if err != nil {
		return "", err
	}
	if completeResponse.Failure != nil {
		return "", errorFromFailure(completeResponse.Failure)
	}
	if completeResponse.Status != "SUCCESS" {
		return "", errors.New("upload failed for unknown reason")
	}
	return "", nil
}

func BuildInitialRequest(name string, tags []string, paths []string) (*CdsInitialRequest, error) {
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
	return &CdsInitialRequest{
		Version:   "0.1.0",
		Name:      name,
		Tags:      tags,
		Documents: documents,
	}, nil
}

func errorFromFailure(failure *CdsFailure) error {
	return fmt.Errorf("failed upload due to %s (%s)", failure.Issue, failure.Message)
}

func ReaderForPaths(paths []string) TarballReader {
	var outBuffer = new(bytes.Buffer)
	return TarballReader{
		paths:                paths,
		indexOfCurrentFile:   -1,
		currentFileHandle:    nil,
		readPosInCurrentFile: 0,
		lengthOfCurrentFile:  0,
		outputBuffer:         outBuffer,
		fileReadBuffer:       make([]byte, 1024), // TODO: Tune this to something bigger?
		tarballWriter:        tar.NewWriter(outBuffer),
		finishedWriting:      false,
	}
}

type TarballReader struct {
	paths              []string
	indexOfCurrentFile int
	currentFileHandle  fs.File
	readPosInCurrentFile int64
	lengthOfCurrentFile int64
	outputBuffer   *bytes.Buffer
	fileReadBuffer []byte
	tarballWriter  *tar.Writer
	finishedWriting bool
}

func (r *TarballReader) Read(p []byte) (n int, err error) {
	if r.outputBuffer.Len() > 0 {
		return r.outputBuffer.Read(p)
	}
	if r.finishedWriting {
		return 0, io.EOF
	}
	if r.readPosInCurrentFile == r.lengthOfCurrentFile {
		if r.currentFileHandle != nil {
			err := r.currentFileHandle.Close()
			if err != nil {
				return 0, err
			}
		}
		r.indexOfCurrentFile += 1
		if r.indexOfCurrentFile == len(r.paths) {
			err := r.tarballWriter.Close()
			if err != nil {
				return 0, err
			}
			r.finishedWriting = true
		} else {
			path := r.paths[r.indexOfCurrentFile]
			file, err := os.Open(path)
			if err != nil {
				return 0, err
			}
			r.currentFileHandle = file
			fileInfo, err := r.currentFileHandle.Stat()
			if err != nil {
				return 0, err
			}
			r.lengthOfCurrentFile = fileInfo.Size()
			err = r.tarballWriter.WriteHeader(&tar.Header{
				Name: filepath.Base(path),
				Mode: 0600,
				Size: fileInfo.Size(),
			})
			if err != nil {
				return 0, err
			}
		}
	} else {
		n, err = r.currentFileHandle.Read(r.fileReadBuffer)
		if err != nil {
			return 0, err
		}
		r.readPosInCurrentFile += int64(n)
		_, err = r.tarballWriter.Write(r.fileReadBuffer[0:n])
		if err != nil {
			return 0, err
		}
	}
	return r.outputBuffer.Read(p)
}
