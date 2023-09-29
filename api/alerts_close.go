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

type alertCloseReason int

const (
	AlertCloseReasonOther alertCloseReason = iota
	AlertCloseReasonFalsePositive
	AlertCloseReasonNotEnoughInfo
	AlertCloseReasonMalicious
	AlertCloseReasonExpected
	AlertCloseReasonExpectedBehavior
)

// String returns the string representation of an Alert closure reason
func (i alertCloseReason) String() string {
	return AlertCloseReasons[i]
}

type alertCloseReasons map[alertCloseReason]string

// AlertCloseReasons is the list of available Alert closure reasons
var AlertCloseReasons = alertCloseReasons{
	AlertCloseReasonOther:            "Other",
	AlertCloseReasonFalsePositive:    "False positive",
	AlertCloseReasonNotEnoughInfo:    "Not enough information",
	AlertCloseReasonMalicious:        "Malicious and have resolution in place",
	AlertCloseReasonExpected:         "Expected because of routine testing",
	AlertCloseReasonExpectedBehavior: "Expected Behavior",
}

func (acr alertCloseReasons) GetOrderedReasonStrings() []string {
	reasons := []string{}
	for i := 0; i < len(acr); i++ {
		reasons = append(reasons, acr[alertCloseReason(i)])
	}
	return reasons
}

type AlertCloseRequest struct {
	AlertID int    `json:"-"`
	Reason  int    `json:"reason"`
	Comment string `json:"comment,omitempty"`
}

type AlertCloseResponse struct {
	Message string `json:"message"`
}

func (svc *AlertsService) Close(request AlertCloseRequest) (
	response AlertCloseResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder(
		"POST",
		fmt.Sprintf(apiV2AlertsClose, request.AlertID),
		request,
		&response,
	)
	return
}
