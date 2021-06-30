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
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// PolicyService is a service that interacts with the Custom Policies
// endpoints from the Lacework Server
type PolicyService struct {
	client *Client
}

// ValidPolicySeverities is a list of all valid policy severities
var ValidPolicySeverities = []string{"critical", "high", "medium", "low", "info"}

type PolicyResponse struct {
	Data    Policy `json:"data"`
	Message string `json:"message"`
}

type PoliciesResponse struct {
	Data    []Policy `json:"data"`
	Message string   `json:"message"`
}

type Policy struct {
	ID           string                 `json:"policyId"`
	Title        string                 `json:"title"`
	Enabled      bool                   `json:"enabled"`
	AlertEnabled bool                   `json:"alertEnabled"`
	Frequency    string                 `json:"evalFrequency"`
	Severity     string                 `json:"severity"`
	QueryID      string                 `json:"queryId"`
	AlertProfile string                 `json:"alertProfile"`
	Limit        int                    `json:"limit"`
	Description  string                 `json:"description"`
	Remediation  string                 `json:"remediation"`
	EvaluatorId  string                 `json:"evaluatorId"`
	Properties   map[string]interface{} `json:"properties"`
}

func (svc *PolicyService) Create(policy string) (
	response PolicyResponse,
	err error,
) {
	var p map[string]interface{}
	if err = json.Unmarshal([]byte(policy), &p); err != nil {
		err = errors.Wrap(err, "policy must be valid JSON")
		return
	}
	err = svc.client.RequestEncoderDecoder("POST", apiV2Policies, p, &response)
	return
}

func (svc *PolicyService) GetAll() (
	response PoliciesResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET", apiV2Policies, nil, &response)
	return
}

func (svc *PolicyService) GetByID(policyID string) (
	response PolicyResponse,
	err error,
) {
	if policyID == "" {
		err = errors.New("policy ID must be provided")
		return
	}
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf("%s/%s", apiV2Policies, url.QueryEscape(policyID)),
		nil,
		&response,
	)
	return
}

func (svc *PolicyService) Update(policyID, policy string) (
	response PolicyResponse,
	err error,
) {
	var p map[string]interface{}
	if err = json.Unmarshal([]byte(policy), &p); err != nil {
		err = errors.Wrap(err, "policy must be valid JSON")
		return
	}

	// retreive policyID from payload and delete it
	if payloadPolicyID, ok := p["policy_id"]; ok {
		delete(p, "policy_id")
		// if policyID is unset, take from the payload
		if policyID == "" {
			policyID = fmt.Sprintf("%v", payloadPolicyID)
		}
	}

	// if policyID is still not specified; error
	if policyID == "" {
		err = errors.New("policy ID must be provided")
		return
	}
	err = svc.client.RequestEncoderDecoder(
		"PATCH",
		fmt.Sprintf("%s/%s", apiV2Policies, url.QueryEscape(policyID)),
		p,
		&response,
	)
	return
}

func (svc *PolicyService) Delete(policyID string) (
	response PolicyResponse,
	err error,
) {
	if policyID == "" {
		err = errors.New("policy ID must be provided")
		return
	}
	err = svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf("%s/%s", apiV2Policies, url.QueryEscape(policyID)),
		nil,
		&response,
	)
	return
}
