package api

import "fmt"

type ComponentsService struct {
	client *Client
}

type ListComponentsResponse struct {
	Data    []LatestComponent `json:"data"`
	Message string            `json:"message"`
}

type LatestComponent struct {
	Components []LatestComponentVersion `json:"components"`
}

type LatestComponentVersion struct {
	Id            int32  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Version       string `json:"version"`
	Size          int64  `json:"size"`
	ComponentType string `json:"type"`
	Deprecated    bool   `json:"deprecated"`
}

func (svc *ComponentsService) ListComponents(os string, arch string) (response ListComponentsResponse, err error) {
	apiPath := fmt.Sprintf(apiV2Components, os, arch)

	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)

	return
}

type ListComponentVersionsResponse struct {
	Data []ComponentVersions `json:"data"`
}

type ComponentVersions struct {
	Id             int32    `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Component_type string   `json:"type"`
	Deprecated     bool     `json:"deprecated"`
	Versions       []string `json:"versions"`
}

func (svc *ComponentsService) ListComponentVersions(id int32, os string, arch string) (
	response ListComponentVersionsResponse,
	err error) {
	apiPath := fmt.Sprintf(apiV2ComponentsVersions, id, os, arch)

	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)

	return
}

type FetchComponentResponse struct {
	Data []Artifact `json:"data"`
}

type Artifact struct {
	Id             int32  `json:"id"`
	Name           string `json:"name"`
	Version        string `json:"version"`
	Size           int64  `json:"size"`
	InstallMessage string `json:"installMessage"`
	UpdateMessage  string `json:"updateMessage"`
	ArtifactUrl    string `json:"artifact_url"`
}

func (svc *ComponentsService) FetchComponentArtifact(id int32, os string, arch string, version string) (
	response FetchComponentResponse,
	err error) {
	apiPath := fmt.Sprintf(apiV2ComponentsFetch, id, os, arch, version)

	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)

	return
}
