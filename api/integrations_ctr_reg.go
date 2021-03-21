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
	EcrRegistry
	GcrRegistry
)

// RegistryTypes is the list of available registry types
var RegistryTypes = map[registryType]string{
	NoneRegistry:      "NONE",
	DockerHubRegistry: "DOCKERHUB",
	DockerV2Registry:  "V2_REGISTRY",
	EcrRegistry:       "AWS_ECR",
	GcrRegistry:       "GCP_GCR",
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

func NewDockerHubRegistryIntegration(name string, data ContainerRegData) ContainerRegIntegration {
	data.RegistryType = DockerHubRegistry.String()
	data.RegistryDomain = "index.docker.io"
	return NewContainerRegIntegration(name, data)
}

func NewDockerV2RegistryIntegration(name string, data ContainerRegData) ContainerRegIntegration {
	data.RegistryType = DockerV2Registry.String()
	return NewContainerRegIntegration(name, data)
}

func NewGcrRegistryIntegration(name string, data ContainerRegData) ContainerRegIntegration {
	data.RegistryType = GcrRegistry.String()
	return NewContainerRegIntegration(name, data)
}

// CreateContainerRegistry creates a container registry integration on the Lacework Server
func (svc *IntegrationsService) CreateContainerRegistry(integration ContainerRegIntegration) (
	response ContainerRegIntResponse,
	err error,
) {
	err = svc.create(integration, &response)
	return
}

// GetContainerRegistry gets a container registry integration that matches with
// the provided integration guid on the Lacework Server
func (svc *IntegrationsService) GetContainerRegistry(guid string) (
	response ContainerRegIntResponse,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateContainerRegistry updates a single container registry integration
func (svc *IntegrationsService) UpdateContainerRegistry(integration ContainerRegIntegration) (
	response ContainerRegIntResponse,
	err error,
) {
	err = svc.update(integration.IntgGuid, integration, &response)
	return
}

// ListContainerRegistries lists the CONT_VULN_CFG external integrations available on the Lacework Server
func (svc *IntegrationsService) ListContainerRegistries() (response ContainerRegIntResponse, err error) {
	err = svc.listByType(ContainerRegistryIntegration, &response)
	return
}

type ContainerRegIntResponse struct {
	Data    []ContainerRegIntegration `json:"data"`
	Ok      bool                      `json:"ok"`
	Message string                    `json:"message"`
}

type ContainerRegIntegration struct {
	commonIntegrationData
	Data ContainerRegData `json:"DATA"`
}

type ContainerRegData struct {
	// @afiune the container registry schema contains a few different DATA types,
	// and because of that we are adding ALL fields that we could possibly have
	// for ALL container registry types (look at the variable RegistryTypes) with
	// the exception of AWS_ECR, this integration has a different credentials field
	// and because of that we have to define it separately
	Credentials  ContainerRegCreds `json:"CREDENTIALS" mapstructure:"CREDENTIALS"`
	RegistryType string            `json:"REGISTRY_TYPE" mapstructure:"REGISTRY_TYPE"`

	// for GCP_GCR integrations, the registry domain has to be one of:
	// => [ "gcr.io", "us.gcr.io", "eu.gcr.io", "asia.gcr.io" ]
	RegistryDomain string `json:"REGISTRY_DOMAIN" mapstructure:"REGISTRY_DOMAIN"`
	LimitByTag     string `json:"LIMIT_BY_TAG" mapstructure:"LIMIT_BY_TAG"`
	LimitByLabel   string `json:"LIMIT_BY_LABEL" mapstructure:"LIMIT_BY_LABEL"`
	LimitByRep     string `json:"LIMIT_BY_REP,omitempty" mapstructure:"LIMIT_BY_REP"`
	LimitNumImg    int    `json:"LIMIT_NUM_IMG,omitempty" mapstructure:"LIMIT_NUM_IMG"`
}

type ContainerRegCreds struct {
	// for docker hub registry (DOCKERHUB)
	Username string `json:"USERNAME,omitempty" mapstructure:"USERNAME"`
	Password string `json:"PASSWORD,omitempty" mapstructure:"PASSWORD"`

	// for docker V2 registry (V2_REGISTRY)
	SSL bool `json:"SSL,omitempty" mapstructure:"SSL"`

	// for GCR registry (GCP_GCR)
	ClientEmail  string `json:"CLIENT_EMAIL,omitempty" mapstructure:"CLIENT_EMAIL"`
	ClientID     string `json:"CLIENT_ID,omitempty" mapstructure:"CLIENT_ID"`
	PrivateKey   string `json:"PRIVATE_KEY,omitempty" mapstructure:"PRIVATE_KEY"`
	PrivateKeyID string `json:"PRIVATE_KEY_ID,omitempty" mapstructure:"PRIVATE_KEY_ID"`
}
