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

import "fmt"

type alertTemplatesService struct {
	client *Client
}

func (svc *alertTemplatesService) Create(alertProfileID string, template AlertTemplate) (
	response AlertProfileResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiV2AlertTemplates, alertProfileID)
	err = svc.client.RequestEncoderDecoder("POST", apiPath, template, &response)
	return
}

func (svc *alertTemplatesService) Update(alertProfileID string, template AlertTemplate) (
	response AlertProfileResponse,
	err error,
) {
	body := alertTemplateUpdate{template.EventName, template.Description, template.Subject}
	apiPath := fmt.Sprintf(apiV2AlertTemplatesFromGUID, alertProfileID, template.Name)
	err = svc.client.RequestEncoderDecoder("POST", apiPath, body, &response)
	return
}

func (svc *alertTemplatesService) Delete(alertProfileID string, alertTemplateID string) (
	err error,
) {
	apiPath := fmt.Sprintf(apiV2AlertTemplatesFromGUID, alertProfileID, alertTemplateID)
	err = svc.client.RequestEncoderDecoder("POST", apiPath, nil, nil)
	return
}

type AlertTemplate struct {
	Name        string `json:"name"`
	EventName   string `json:"eventName"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
}

type alertTemplatesUpdate struct {
	Alerts []AlertTemplate `json:"alerts"`
}

type alertTemplateUpdate struct {
	EventName   string `json:"eventName,omitempty"`
	Description string `json:"description,omitempty"`
	Subject     string `json:"subject,omitempty"`
}
