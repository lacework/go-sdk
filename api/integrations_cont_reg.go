//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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

type registryType int

const (
	// type that defines a non-existing registry
	NoneRegistry registryType = iota
	DockerHubRegistry
	DockerV2Registry
)

// RegistryTypes is the list of available registry types
var RegistryTypes = map[registryType]string{
	NoneRegistry:      "NONE",
	DockerHubRegistry: "DOCKERHUB",
	DockerV2Registry:  "V2_REGISTRY",
}

// String returns the string representation of an registry type
func (i registryType) String() string {
	return RegistryTypes[i]
}

// NewContainerRegIntegration returns an instance of ContainerRegIntegration
// with the provided name and data.
//
// Basic usage: Create a Docker Hub integration
//
//   client, err := api.NewClient("account")
//   if err != nil {
//     return err
//   }
//
//   docker := api.NewContainerRegIntegration("foo",
//     api.ContainerRegData{
//       Credentials: api.ContainerRegCreds {
//         Username: "techally",
//         Password: "secret",
//       },
//       RegistryType: api.DockerHubRegistry.String(),
//       RegistryDomain: "index.docker.io",
//       LimitByTag: "*",
//       LimitByLabel: "*",
//       LimitNumImg: "5",
//     },
//   )
//
//   client.Integrations.CreateContainerRegistry(docker)
//
func NewContainerRegIntegration(name string, data ContainerRegData) ContainerRegIntegration {
	return ContainerRegIntegration{
		commonIntegrationData: commonIntegrationData{
			Name:    name,
			Type:    ContainerRegistryIntegration.String(),
			Enabled: 1,
		},
		Data: data,
	}
}

// CreateContainerRegistry creates a container registry integration on the Lacework Server
func (svc *IntegrationsService) CreateContainerRegistry(integration ContainerRegIntegration) (
	response map[string]interface{},
	//response ContainerRegIntResponse, // @afiune we can't use this :(
	err error,
) {
	err = svc.create(integration, &response)
	return
}

type ContainerRegIntegration struct {
	commonIntegrationData
	Data ContainerRegData `json:"DATA"`
}

type ContainerRegData struct {
	Credentials    ContainerRegCreds `json:"CREDENTIALS"`
	RegistryType   string            `json:"REGISTRY_TYPE"`
	RegistryDomain string            `json:"REGISTRY_DOMAIN"`
	LimitByTag     string            `json:"LIMIT_BY_TAG"`
	LimitByLabel   string            `json:"LIMIT_BY_LABEL"`
	LimitByRep     string            `json:"LIMIT_BY_REP,omitempty"`
	LimitNumImg    int               `json:"LIMIT_NUM_IMG"`
}

type ContainerRegCreds struct {
	Username string `json:"USERNAME"`
	Password string `json:"PASSWORD"`
	// @afiune this is for docker V2 registry
	SSL bool `json:"SSL,omitempty"`
}

// @afiune we can't use this response since the request sent to the
// Server is different from the one it returns as a response. :(
// If we enable this struct we will get the following error:
//
// json: cannot unmarshal string into Go struct field
//       ContainerRegData.data.DATA.LIMIT_NUM_IMG of type int
type ContainerRegIntResponse struct {
	Data    []ContainerRegIntegration `json:"data"`
	Ok      bool                      `json:"ok"`
	Message string                    `json:"message"`
}
