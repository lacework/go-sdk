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
	"reflect"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
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

type PoliciesYAML struct {
	Policies []Policy `yaml:policies`
}

// Use of pointers intentional here for optimal json.Marshal() behavior
// This permits distinction between bool/int zero values vs. missing/unspecified
type Policy struct {
	ID           string                 `json:"policy_id,omitempty" yaml:"policy_id"`
	Title        string                 `json:"title,omitempty" yaml:"title"`
	Enabled      *bool                  `json:"enabled,omitempty" yaml:"enabled"`
	AlertEnabled *bool                  `json:"alert_enabled,omitempty" yaml:"alert_enabled"`
	Frequency    string                 `json:"eval_frequency,omitempty" yaml:"eval_frequency"`
	Severity     string                 `json:"severity,omitempty" yaml:"severity"`
	QueryID      string                 `json:"lql_id,omitempty" yaml:"lql_id"`
	AlertProfile string                 `json:"alert_profile,omitempty" yaml:"alert_profile"`
	Limit        *int                   `json:"limit,omitempty" yaml:"limit"`
	Description  string                 `json:"description,omitempty" yaml:"description"`
	Remediation  string                 `json:"remediation,omitempty" yaml:"remediation"`
	Properties   map[string]interface{} `json:"properties,omitempty" yaml:"properties"`
}

func TranslatePolicy(s string) (Policy, error) {
	var policy Policy
	var err error

	// valid json
	if err = json.Unmarshal([]byte(s), &policy); err == nil {
		return policy, err
	}
	// nested yaml
	var policies PoliciesYAML

	if err = yaml.Unmarshal([]byte(s), &policies); err == nil {
		if len(policies.Policies) > 0 {
			return policies.Policies[0], err
		}
	}
	// straight yaml
	policy = Policy{}
	err = yaml.Unmarshal([]byte(s), &policy)
	if err == nil && !reflect.DeepEqual(policy, Policy{}) { // empty string unmarshals w/o error
		return policy, nil
	}
	// invalid policy
	return policy, errors.New("policy must be valid JSON or YAML")
}

func (svc *PolicyService) Create(policy string) (
	response PolicyResponse,
	err error,
) {
	var p Policy
	if p, err = TranslatePolicy(policy); err != nil {
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
	var p Policy
	if p, err = TranslatePolicy(policy); err != nil {
		return
	}

	// retreive policyID from payload and delete it
	if p.ID != "" {
		policyID = p.ID
		p.ID = ""
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
