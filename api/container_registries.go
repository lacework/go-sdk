//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package api

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

// ContainerRegistriesService is the service that interacts with
// the ContainerRegistries schema from the Lacework APIv2 Server
type ContainerRegistriesService struct {
	client *Client
}

// NewContainerRegistry returns an instance of the ContainerRegistryRaw struct with the
// provided Container Registry integration type, name and raw data as an interface{}.
//
// NOTE: This function must be used by any Container Registry type.
//
// Basic usage: Initialize a new GhcrContainerRegistry integration struct, then
//              use the new instance to do CRUD operations
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   ghcr := api.NewContainerRegistry("foo",
//     api.GhcrContainerRegistry,
//     api.GhcrData{
//       Credentials: api.GhcrCredentials {
//         Username: "bubu",
//         Password: "supers3cret",
//         Ssl: true,
//       },
//     },
//   )
//
//   client.V2.ContainerRegistries.Create(ghcr)
//
func NewContainerRegistry(name string, regType containerRegistryType, data interface{}) ContainerRegistryRaw {
	reg := ContainerRegistryRaw{
		v2CommonIntegrationData: v2CommonIntegrationData{
			Name:    name,
			Type:    "ContVulnCfg",
			Enabled: 1,
		},
	}

	switch regType {
	case GcpGarContainerRegistry:
		reg.Data = verifyGcpGarContainerRegistry(data)
	case GhcrContainerRegistry:
		reg.Data = verifyGhcrContainerRegistry(data)
	case InlineScannerContainerRegistry:
		reg.Data = verifyInlineScannerContainerRegistry(data)
	case ProxyScannerContainerRegistry:
		reg.Data = verifyProxyScannerContainerRegistry(data)
	case AwsEcrContainerRegistry:
		reg.Data = verifyAwsEcrContainerRegistry(data)
	case DockerhubContainerRegistry:
		reg.Data = verifyDockerhubContainerRegistry(data)
	case DockerhubV2ContainerRegistry:
		reg.Data = verifyDockerhubV2ContainerRegistry(data)
	case GcpGcrContainerRegistry:
		reg.Data = verifyGcpGcrContainerRegistry(data)
	default:
		reg.Data = data
	}

	return reg
}

// ContainerRegistry is an interface that helps us implement a few functions
// that any Container Registry might use, there are some cases, like during
// Update, where we need to get the ID of the Container Registry and its type,
// this will allow users to pass any Container Registry that implements these
// methods
type ContainerRegistry interface {
	ID() string
	ContainerRegistryType() containerRegistryType
}

type containerRegistryType int

const (
	// type that defines a non-existing Container Registry integration
	NoneContainerRegistry containerRegistryType = iota
	GcpGarContainerRegistry
	GhcrContainerRegistry
	InlineScannerContainerRegistry
	ProxyScannerContainerRegistry
	AwsEcrContainerRegistry
	DockerhubContainerRegistry
	DockerhubV2ContainerRegistry
	GcpGcrContainerRegistry
)

// ContainerRegistryTypes is the list of available Container Registry integration types
var ContainerRegistryTypes = map[containerRegistryType]string{
	NoneContainerRegistry:          "None",
	GcpGarContainerRegistry:        "GCP_GAR",
	GhcrContainerRegistry:          "GHCR",
	InlineScannerContainerRegistry: "INLINE_SCANNER",
	ProxyScannerContainerRegistry:  "PROXY_SCANNER",
	AwsEcrContainerRegistry:        "AWS_ECR",
	DockerhubContainerRegistry:     "DOCKERHUB",
	DockerhubV2ContainerRegistry:   "V2_REGISTRY",
	GcpGcrContainerRegistry:        "GCP_GCR",
}

// String returns the string representation of a Container Registry integration type
func (i containerRegistryType) String() string {
	return ContainerRegistryTypes[i]
}

// FindContainerRegistryType looks up inside the list of available container registry types
// the matching type from the provided string, if none, returns NoneContainerRegistry
func FindContainerRegistryType(containerRegistry string) (containerRegistryType, bool) {
	for cType, cStr := range ContainerRegistryTypes {
		if cStr == containerRegistry {
			return cType, true
		}
	}
	return NoneContainerRegistry, false
}

// List returns a list of Container Registry integrations
func (svc *ContainerRegistriesService) List() (response ContainerRegistriesResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2ContainerRegistries, nil, &response)
	return
}

// Create creates a single Container Registry integration
func (svc *ContainerRegistriesService) Create(integration ContainerRegistryRaw) (
	response ContainerRegistryResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// Delete deletes a Container Registry integration that matches the provided guid
func (svc *ContainerRegistriesService) Delete(guid string) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}

	return svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf(apiV2ContainerRegistryFromGUID, guid),
		nil,
		nil,
	)
}

// Get returns a raw response of the Container Registry with the matching integration guid.
//
// To return a more specific Go struct of a Container Registry integration, use the proper
// method such as GetGhcr() where the function name is composed by:
//
//  Get<Type>(guid)
//
//    Where <Type> is the Container Registry integration type.
func (svc *ContainerRegistriesService) Get(guid string, response interface{}) error {
	return svc.get(guid, &response)
}

type ContainerRegistryRaw struct {
	v2CommonIntegrationData
	Data        interface{}    `json:"data,omitempty"`
	ServerToken *V2ServerToken `json:"serverToken,omitempty"`
}

func (reg ContainerRegistryRaw) StateString() string {
	switch reg.ContainerRegistryType() {
	case InlineScannerContainerRegistry, ProxyScannerContainerRegistry:
		return "Ok"
	default:
		return reg.v2CommonIntegrationData.StateString()
	}
}

type V2ServerToken struct {
	ServerToken string `json:"serverToken"`
	Uri         string `json:"uri"`
}

func (reg ContainerRegistryRaw) GetData() any {
	return reg.Data
}

func (reg ContainerRegistryRaw) GetCommon() v2CommonIntegrationData {
	return reg.v2CommonIntegrationData
}

func (reg ContainerRegistryRaw) ContainerRegistryType() containerRegistryType {
	if casting, ok := reg.Data.(map[string]interface{}); ok {
		if regType, exist := casting["registryType"]; exist {
			t, _ := FindContainerRegistryType(regType.(string))
			return t
		}
	}

	m := structs.Map(reg.Data)
	if regType, exist := m["RegistryType"]; exist {
		t, _ := FindContainerRegistryType(regType.(string))
		return t
	}

	return NoneContainerRegistry
}

func (reg ContainerRegistryRaw) ContainerRegistryDomain() string {
	if casting, ok := reg.Data.(map[string]interface{}); ok {
		if domain, exist := casting["registryDomain"]; exist {
			return domain.(string)
		}
	}

	if structs.IsStruct(reg.Data) {
		m := structs.Map(reg.Data)
		if domain, exist := m["RegistryDomain"]; exist {
			return domain.(string)
		}
	}
	return ""
}

type ContainerRegistryResponse struct {
	Data ContainerRegistryRaw `json:"data"`
}

type ContainerRegistriesResponse struct {
	Data []ContainerRegistryRaw `json:"data"`
}

func (svc *ContainerRegistriesService) create(data interface{}, response interface{}) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2ContainerRegistries, data, response)
}

func (svc *ContainerRegistriesService) get(guid string, response interface{}) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}
	apiPath := fmt.Sprintf(apiV2ContainerRegistryFromGUID, guid)
	return svc.client.RequestDecoder("GET", apiPath, nil, response)
}

func (svc *ContainerRegistriesService) update(guid string, data interface{}, response interface{}) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}
	apiPath := fmt.Sprintf(apiV2ContainerRegistryFromGUID, guid)
	return svc.client.RequestEncoderDecoder("PATCH", apiPath, data, response)
}
