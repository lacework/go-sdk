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
	Data    []Policy `json:"data"`
	Ok      bool     `json:"ok"`
	Message string   `json:"message"`
}

type Policy struct {
	ID           string `json:"policy_id"`
	Title        string `json:"title"`
	Enabled      bool   `json:"enabled"`
	AlertEnabled bool   `json:"alert_enabled"`
	Frequency    string `json:"eval_frequency"`
	Severity     string `json:"severity"`
	QueryID      string `json:"lql_id"`
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
	err = svc.client.RequestEncoderDecoder("POST", apiPolicy, p, &response)
	return
}

func (svc *PolicyService) GetAll() (PolicyResponse, error) {
	return svc.GetByID("")
}

func (svc *PolicyService) GetByID(policyID string) (
	response PolicyResponse,
	err error,
) {
	uri := apiPolicy

	if policyID != "" {
		uri += "?POLICY_ID=" + url.QueryEscape(policyID)
	}

	err = svc.client.RequestDecoder("GET", uri, nil, &response)
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

	uri := fmt.Sprintf("%s?POLICY_ID=%s", apiPolicy, url.QueryEscape(policyID))
	err = svc.client.RequestEncoderDecoder("PATCH", uri, p, &response)
	return
}

func (svc *PolicyService) Delete(policyID string) (
	response map[string]interface{}, // endpoint currently 204's so no response content
	err error,
) {
	if policyID == "" {
		err = errors.New("policy ID must be provided")
		return
	}

	err = svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf("%s?POLICY_ID=%s", apiPolicy, url.QueryEscape(policyID)),
		nil,
		&response,
	)
	return
}
