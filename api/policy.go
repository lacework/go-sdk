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
	"time"

	"github.com/lacework/go-sdk/internal/array"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// PolicyService is a service that interacts with the Custom Policies
// endpoints from the Lacework Server
type PolicyService struct {
	client     *Client
	Exceptions *policyExceptionsService
}

func NewV2PolicyService(c *Client) *PolicyService {
	return &PolicyService{c,
		&policyExceptionsService{c},
	}
}

type policyType int
type policyTypes map[policyType]string

const (
	PolicyTypeCompliance policyType = iota
	PolicyTypeManual
	PolicyTypeViolation
)

var ValidPolicyTypes = policyTypes{
	PolicyTypeCompliance: "Compliance",
	PolicyTypeManual:     "Manual",
	PolicyTypeViolation:  "Violation",
}

func (p policyType) String() string {
	return ValidPolicyTypes[p]
}

func (pt policyTypes) String() (types []string) {
	for _, v := range pt {
		types = append(types, v)
	}
	return
}

// ValidPolicySeverities is a list of all valid policy severities
var ValidPolicySeverities = []string{"critical", "high", "medium", "low", "info"}

type NewPolicy struct {
	PolicyID      string   `json:"policyId,omitempty" yaml:"policyId,omitempty" `
	PolicyType    string   `json:"policyType" yaml:"policyType"`
	QueryID       string   `json:"queryId" yaml:"queryId"`
	Title         string   `json:"title" yaml:"title"`
	Enabled       bool     `json:"enabled" yaml:"enabled"`
	Description   string   `json:"description" yaml:"description"`
	Remediation   string   `json:"remediation" yaml:"remediation"`
	Severity      string   `json:"severity" yaml:"severity"`
	Limit         int      `json:"limit,omitempty" yaml:"limit,omitempty"`
	EvalFrequency string   `json:"evalFrequency,omitempty" yaml:"evalFrequency,omitempty"`
	AlertEnabled  bool     `json:"alertEnabled" yaml:"alertEnabled"`
	AlertProfile  string   `json:"alertProfile,omitempty" yaml:"alertProfile,omitempty"`
	Tags          []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type newPoliciesYAML struct {
	Policies []NewPolicy `yaml:"policies"`
}

func ParseNewPolicy(s string) (NewPolicy, error) {
	var policy NewPolicy
	var err error

	// valid json
	if err = json.Unmarshal([]byte(s), &policy); err == nil {
		return policy, err
	}
	// nested yaml
	var policies newPoliciesYAML

	if err = yaml.Unmarshal([]byte(s), &policies); err == nil {
		if len(policies.Policies) > 0 {
			return policies.Policies[0], err
		}
	}
	// straight yaml
	policy = NewPolicy{}
	err = yaml.Unmarshal([]byte(s), &policy)
	if err == nil && !reflect.DeepEqual(policy, NewPolicy{}) { // empty string unmarshals w/o error
		return policy, nil
	}
	// invalid policy
	return policy, errors.New("policy must be valid JSON or YAML")
}

/* In order to properly PATCH we need to omit items that aren't specified.
For booleans and integers Golang will omit zero values false and 0 respectively.
This would prevent someone from toggling something to disabled or 0 respectively.
As such we are using pointers instead of primitives for booleans and integers in this struct
*/
type UpdatePolicy struct {
	PolicyID      string   `json:"policyId,omitempty" yaml:"policyId,omitempty"`
	PolicyType    string   `json:"policyType,omitempty" yaml:"policyType,omitempty"`
	QueryID       string   `json:"queryId,omitempty" yaml:"queryId,omitempty"`
	Title         string   `json:"title,omitempty" yaml:"title,omitempty"`
	Enabled       *bool    `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Description   string   `json:"description,omitempty" yaml:"description,omitempty"`
	Remediation   string   `json:"remediation,omitempty" yaml:"remediation,omitempty"`
	Severity      string   `json:"severity,omitempty" yaml:"severity,omitempty"`
	Limit         *int     `json:"limit,omitempty" yaml:"limit,omitempty"`
	EvalFrequency string   `json:"evalFrequency,omitempty" yaml:"evalFrequency,omitempty"`
	AlertEnabled  *bool    `json:"alertEnabled,omitempty" yaml:"alertEnabled,omitempty"`
	AlertProfile  string   `json:"alertProfile,omitempty" yaml:"alertProfile,omitempty"`
	Tags          []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type updatePoliciesYAML struct {
	Policies []UpdatePolicy `yaml:"policies"`
}

func ParseUpdatePolicy(s string) (UpdatePolicy, error) {
	var policy UpdatePolicy
	var err error

	// valid json
	if err = json.Unmarshal([]byte(s), &policy); err == nil {
		return policy, err
	}
	// nested yaml
	var policies updatePoliciesYAML

	if err = yaml.Unmarshal([]byte(s), &policies); err == nil {
		if len(policies.Policies) > 0 {
			return policies.Policies[0], err
		}
	}
	// straight yaml
	policy = UpdatePolicy{}
	err = yaml.Unmarshal([]byte(s), &policy)
	if err == nil && !reflect.DeepEqual(policy, UpdatePolicy{}) { // empty string unmarshals w/o error
		return policy, nil
	}
	// invalid policy
	return policy, errors.New("policy must be valid JSON or YAML")
}

type Policy struct {
	PolicyID               string                                               `json:"policyId" yaml:"policyId"`
	PolicyType             string                                               `json:"policyType" yaml:"-"`
	QueryID                string                                               `json:"queryId" yaml:"queryId"`
	Title                  string                                               `json:"title" yaml:"title"`
	Enabled                bool                                                 `json:"enabled" yaml:"enabled"`
	Description            string                                               `json:"description" yaml:"description"`
	Remediation            string                                               `json:"remediation" yaml:"remediation"`
	Severity               string                                               `json:"severity" yaml:"severity"`
	Limit                  int                                                  `json:"limit" yaml:"limit"`
	EvalFrequency          string                                               `json:"evalFrequency" yaml:"evalFrequency"`
	AlertEnabled           bool                                                 `json:"alertEnabled" yaml:"alertEnabled"`
	AlertProfile           string                                               `json:"alertProfile" yaml:"alertProfile"`
	Tags                   []string                                             `json:"tags" yaml:"tags"`
	Owner                  string                                               `json:"owner" yaml:"-"`
	LastUpdateTime         string                                               `json:"lastUpdateTime" yaml:"-"`
	LastUpdateUser         string                                               `json:"lastUpdateUser" yaml:"-"`
	ExceptionConfiguration map[string][]PolicyExceptionConfigurationConstraints `json:"exceptionConfiguration" yaml:"-"`
}

type PolicyExceptionConfigurationConstraints struct {
	DataType   string `json:"dataType" yaml:"dataType"`
	FieldKey   string `json:"fieldKey" yaml:"fieldKey"`
	MultiValue bool   `json:"multiValue" yaml:"multiValue"`
}

func (p *Policy) HasTag(t string) bool {
	return array.ContainsStr(p.Tags, t)
}

type PolicyResponse struct {
	Data    Policy `json:"data"`
	Message string `json:"message"`
}

type PolicyTagsResponse struct {
	Data    []string `json:"data"`
	Message string   `json:"message"`
}

type PoliciesResponse struct {
	Data    []Policy `json:"data"`
	Message string   `json:"message"`
}

func (svc *PolicyService) Create(np NewPolicy) (
	response PolicyResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2Policies, np, &response)
	return
}

func (svc *PolicyService) List() (
	response PoliciesResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET", apiV2Policies, nil, &response)
	return
}

func (svc *PolicyService) ListTags() (
	response PolicyTagsResponse,
	err error,
) {
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf("%s/Tags", apiV2Policies),
		nil,
		&response,
	)
	return
}

func (svc *PolicyService) Get(policyID string) (
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

func (svc *PolicyService) Update(up UpdatePolicy) (
	response PolicyResponse,
	err error,
) {
	if up.PolicyID == "" {
		err = errors.New("policy ID must be provided")
		return
	}
	var policyID = up.PolicyID
	up.PolicyID = "" // omit this for PATCH

	err = svc.client.RequestEncoderDecoder(
		"PATCH",
		fmt.Sprintf("%s/%s", apiV2Policies, url.QueryEscape(policyID)),
		up,
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

type BulkUpdatePolicy struct {
	PolicyID string `json:"policyId,omitempty" yaml:"policyId,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Severity string `json:"severity,omitempty" yaml:"severity,omitempty"`
}

type BulkUpdatePolicies []BulkUpdatePolicy

// UpdateMany supports updating the state(enabled/disabled) and severity of more than one
// policy using the policy bulk update api
func (svc *PolicyService) UpdateMany(policies BulkUpdatePolicies) (
	response BulkPolicyUpdateResponse,
	err error,
) {
	if len(policies) == 0 {
		err = errors.New("a list of policies must be provided")
		return
	}

	err = svc.client.RequestEncoderDecoder(
		"PATCH",
		apiV2Policies,
		policies,
		&response,
	)
	return
}

type BulkPolicyUpdateResponse struct {
	Data []BulkPolicyUpdateResponseData `json:"data"`
}

type BulkPolicyUpdateResponseData struct {
	EvaluatorId            string    `json:"evaluatorId,omitempty"`
	PolicyId               string    `json:"policyId"`
	PolicyType             string    `json:"policyType"`
	QueryId                string    `json:"queryId,omitempty"`
	QueryText              string    `json:"queryText,omitempty"`
	Title                  string    `json:"title"`
	Enabled                bool      `json:"enabled,omitempty"`
	Description            string    `json:"description"`
	Remediation            string    `json:"remediation"`
	Severity               string    `json:"severity"`
	Limit                  int       `json:"limit,omitempty"`
	EvalFrequency          string    `json:"evalFrequency,omitempty"`
	AlertEnabled           bool      `json:"alertEnabled,omitempty"`
	AlertProfile           string    `json:"alertProfile,omitempty"`
	Owner                  string    `json:"owner"`
	LastUpdateTime         time.Time `json:"lastUpdateTime"`
	LastUpdateUser         string    `json:"lastUpdateUser"`
	Tags                   []string  `json:"tags"`
	InfoLink               string    `json:"infoLink,omitempty"`
	ExceptionConfiguration struct {
		ConstraintFields []struct {
			FieldKey   string `json:"fieldKey"`
			DataType   string `json:"dataType"`
			MultiValue bool   `json:"multiValue"`
		} `json:"constraintFields"`
	} `json:"exceptionConfiguration,omitempty"`
	References            []string `json:"references,omitempty"`
	AdditionalInformation string   `json:"additionalInformation,omitempty"`
}
