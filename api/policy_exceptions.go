//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"errors"
	"fmt"
)

// policyExceptionsService is the service that interacts with
// the Exceptions schema from the Lacework APIv2 Server
type policyExceptionsService struct {
	client *Client
}

// List returns a list of the Policy Exceptions for a policy ID.
func (svc policyExceptionsService) List(policyID string) (response PolicyExceptionsResponse, err error) {
	if policyID == "" {
		return response, errors.New("specify a policy ID")
	}
	err = svc.client.RequestDecoder("GET", fmt.Sprintf(apiV2PolicyExceptions, policyID), nil, &response)
	return
}

// Get returns a raw response of the Policy Exception with the matching policy ID and exception ID.
func (svc policyExceptionsService) Get(policyID string, exceptionID string, response interface{}) error {
	if exceptionID == "" || policyID == "" {
		return errors.New("specify exception and policy IDs")
	}
	apiPath := fmt.Sprintf(apiV2PolicyExceptionsFromExceptionID, exceptionID, policyID)
	return svc.client.RequestDecoder("GET", apiPath, nil, &response)
}

func (svc policyExceptionsService) Delete(policyID string, exceptionID string) error {
	if exceptionID == "" || policyID == "" {
		return errors.New("specify exception and policy IDs")
	}

	return svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf(apiV2PolicyExceptionsFromExceptionID, exceptionID, policyID),
		nil,
		nil,
	)
}

// Create creates a single Policy Exception
func (svc *policyExceptionsService) Create(policyID string, policy PolicyException) (
	response PolicyExceptionResponse,
	err error,
) {
	if policyID == "" {
		return response, errors.New("specify a policy ID")
	}
	err = svc.client.RequestEncoderDecoder("POST", fmt.Sprintf(apiV2PolicyExceptions, policyID),
		policy, &response)
	return
}

// Update updates a single Policy Exception
func (svc policyExceptionsService) Update(policyID string, exception PolicyException) (response PolicyExceptionResponse, err error) {
	if exception.ExceptionID == "" || policyID == "" {
		return response, errors.New("specify exception and policy IDs")
	}
	apiPath := fmt.Sprintf(apiV2PolicyExceptionsFromExceptionID, exception.ExceptionID, policyID)
	// Request is invalid if it contains the ExceptionID, LastUpdatedTime or LastUpdatedUser fields.
	exception.ExceptionID = ""
	exception.LastUpdateUser = ""
	exception.LastUpdateTime = ""

	err = svc.client.RequestEncoderDecoder("PATCH", apiPath, exception, &response)
	return
}

type PolicyExceptionResponse struct {
	Data PolicyException `json:"data"`
}

type PolicyExceptionsResponse struct {
	Data []PolicyException `json:"data"`
}
type PolicyException struct {
	ExceptionID    string                      `json:"exceptionId,omitempty"`
	Description    string                      `json:"description"`
	Constraints    []PolicyExceptionConstraint `json:"constraints"`
	LastUpdateTime string                      `json:"lastUpdateTime,omitempty"`
	LastUpdateUser string                      `json:"lastUpdateUser,omitempty"`
}

type PolicyExceptionConstraint struct {
	FieldKey    string   `json:"fieldKey"`
	FieldValues []string `json:"fieldValues"`
}
