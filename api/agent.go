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

import (
	"fmt"
)

// AgentsService is a service that interacts with the Agent Access Tokens
// endpoints from the Lacework Server
type AgentsService struct {
	client *Client
}

// ListTokens returns a list of agent access tokens in a Lacework account
func (svc *AgentsService) ListTokens() (response AgentTokensResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiAgentTokens, nil, &response)
	return
}

// GetToken returns details about an agent access token
func (svc *AgentsService) GetToken(token string) (response AgentTokensResponse, err error) {
	err = svc.client.RequestDecoder("GET", fmt.Sprintf(apiAgentTokenFromID, token), nil, &response)
	return
}

// CreateToken creates a new agent access token
func (svc *AgentsService) CreateToken(name, desc string) (response AgentTokensResponse, err error) {
	err = svc.client.RequestEncoderDecoder("POST",
		apiAgentTokens,
		AgentTokenRequest{
			TokenAlias: name,
			Enabled:    1,
			Props: &AgentTokenProps{
				Description: desc,
			},
		},
		&response,
	)
	return
}

// UpdateToken updates an agent access token with the provided request data
func (svc *AgentsService) UpdateToken(token string, data AgentTokenRequest) (
	response AgentTokensResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("PUT",
		fmt.Sprintf(apiAgentTokenFromID, token),
		data,
		&response,
	)
	return
}

// UpdateTokenStatus updates only the status of an agent access token (enable or disable)
func (svc *AgentsService) UpdateTokenStatus(token string, enable bool) (
	response AgentTokensResponse,
	err error,
) {

	request := AgentTokenRequest{Enabled: 0}
	if enable {
		request.Enabled = 1
	}
	err = svc.client.RequestEncoderDecoder("PUT",
		fmt.Sprintf(apiAgentTokenFromID, token),
		request,
		&response,
	)
	return
}

type AgentTokensResponse struct {
	Data    []AgentToken `json:"data"`
	Ok      bool         `json:"ok"`
	Message string       `json:"message"`
}

type AgentToken struct {
	AccessToken     string           `json:"ACCESS_TOKEN"`
	Account         string           `json:"ACCOUNT"`
	LastUpdatedTime *Json16DigitTime `json:"LAST_UPDATED_TIME"`
	Props           *AgentTokenProps `json:"PROPS,omitempty"`
	TokenAlias      string           `json:"TOKEN_ALIAS"`
	Enabled         string           `json:"TOKEN_ENABLED"`
	Version         string           `json:"VERSION"`
}

// @afiune this API returns a string as a boolean, so we have to do this mokeypatch
func (t AgentToken) PrettyStatus() string {
	if t.Enabled == "true" {
		return "Enabled"
	}
	return "Disabled"
}
func (t AgentToken) Status() bool {
	return t.Enabled == "true"
}

func (t AgentToken) EnabledInt() int {
	if t.Enabled == "true" {
		return 1
	}
	return 0
}

type AgentTokenRequest struct {
	TokenAlias string           `json:"TOKEN_ALIAS,omitempty"`
	Enabled    int              `json:"TOKEN_ENABLED"`
	Props      *AgentTokenProps `json:"PROPS,omitempty"`
}

type AgentTokenProps struct {
	CreatedTime *Json16DigitTime `json:"CREATED_TIME,omitempty"`
	Description string           `json:"DESCRIPTION,omitempty"`
}
