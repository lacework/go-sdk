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
	"time"

	"github.com/pkg/errors"
)

// AgentAccessTokensService is the service that interacts with
// the AgentAccessTokens schema from the Lacework APIv2 Server
type AgentAccessTokensService struct {
	client *Client
}

// List returns a list of Agent Access Tokens
func (svc *AgentAccessTokensService) List() (response AgentAccessTokensResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2AgentAccessTokens, nil, &response)
	return
}

// Create creates a single Agent Access Token
func (svc *AgentAccessTokensService) Create(alias, desc string) (
	response AgentAccessTokenResponse,
	err error,
) {
	if alias == "" {
		err = errors.New("token alias is required")
		return
	}

	err = svc.client.RequestEncoderDecoder("POST",
		apiV2AgentAccessTokens,
		AgentAccessTokenRequest{
			TokenAlias: alias,
			Enabled:    1,
			Props: &AgentAccessTokenProps{
				Description: desc,
			},
		},
		&response,
	)
	return
}

// Get returns an Agent Access Token with the matching ID (token)
func (svc *AgentAccessTokensService) Get(token string) (
	response AgentAccessTokenResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET",
		fmt.Sprintf(apiV2AgentAccessTokenFromID, token),
		nil,
		&response,
	)
	return
}

// Update updates an Agent Access Token with the provided request data
func (svc *AgentAccessTokensService) Update(token string, data AgentAccessTokenRequest) (
	response AgentAccessTokenResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("PATCH",
		fmt.Sprintf(apiV2AgentAccessTokenFromID, token),
		data,
		&response,
	)
	return
}

// UpdateState updates only the state of an Agent Access Token (enable or disable)
func (svc *AgentAccessTokensService) UpdateState(token string, enable bool) (
	response AgentAccessTokenResponse,
	err error,
) {

	request := AgentAccessTokenRequest{Enabled: 0}
	if enable {
		request.Enabled = 1
	}
	err = svc.client.RequestEncoderDecoder("PATCH",
		fmt.Sprintf(apiV2AgentAccessTokenFromID, token),
		request,
		&response,
	)
	return
}

// SearchAlias will search for an Agent Access Token that matches the provider token alias
func (svc *AgentAccessTokensService) SearchAlias(alias string) (
	response AgentAccessTokensResponse,
	err error,
) {

	if alias == "" {
		err = errors.New("specify a token alias to search")
		return
	}
	err = svc.client.RequestEncoderDecoder("POST",
		apiV2AgentAccessTokensSearch,
		SearchFilter{
			Filters: []Filter{
				Filter{
					Field:      "tokenAlias",
					Expression: "eq",
					Value:      alias,
				},
			},
		},
		&response,
	)
	return
}

type AgentAccessToken struct {
	AccessToken string                `json:"accessToken"`
	CreatedTime time.Time             `json:"createdTime"`
	Props       AgentAccessTokenProps `json:"props,omitempty"`
	TokenAlias  string                `json:"tokenAlias"`
	Enabled     int                   `json:"tokenEnabled"`
	Version     string                `json:"version"`
}

func (t AgentAccessToken) State() bool {
	return t.Enabled == 1
}

func (t AgentAccessToken) PrettyState() string {
	if t.State() {
		return "Enabled"
	}
	return "Disabled"
}

type AgentAccessTokenProps struct {
	CreatedTime time.Time `json:"createdTime,omitempty"`
	Description string    `json:"description,omitempty"`
}

type AgentAccessTokenResponse struct {
	Data AgentAccessToken `json:"data"`
}

type AgentAccessTokensResponse struct {
	Data []AgentAccessToken `json:"data"`
}

type AgentAccessTokenRequest struct {
	Enabled    int                    `json:"tokenEnabled"`
	TokenAlias string                 `json:"tokenAlias,omitempty"`
	Props      *AgentAccessTokenProps `json:"props,omitempty"`
}
