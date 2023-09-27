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
	Name           string `json:"name"`
	Description    string `json:"description"`
	Version        string `json:"version"`
	Size           int64  `json:"size"`
	Component_type string `json:"type"`
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
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Component_type string   `json:"type"`
	Versions       []string `json:"versions"`
}

func (svc *ComponentsService) ListComponentVersions(component string, os string, arch string) (response ListComponentVersionsResponse, err error) {
	apiPath := fmt.Sprintf(apiV2ComponentsVersions, component, os, arch)

	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)

	return
}

type FetchComponentResponse struct {
	Data []Artifact `json:"data"`
}

type Artifact struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	Size           int64  `json:"size"`
	InstallMessage string `json:"installMessage"`
	UpdateMessage  string `json:"updateMessage"`
	ArtifactUrl    string `json:"artifact_url"`
}

func (svc *ComponentsService) FetchComponentArtifact(component string, os string, arch string, version string) (response FetchComponentResponse, err error) {
	apiPath := fmt.Sprintf(apiV2ComponentsFetch, component, os, arch, version)

	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)

	return
}
