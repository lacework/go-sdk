//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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
	"encoding/base64"
	"fmt"
	"strings"
)

// GetServiceNowRest gets a single ServiceNowRest alert channel matching the
// provided integration guid
func (svc *AlertChannelsService) GetServiceNowRest(guid string) (
	response ServiceNowRestAlertChannelResponseV2,
	err error,
) {
	err = svc.get(guid, &response)
	return
}

// UpdateServiceNowRest updates a single ServiceNowRest integration on the Lacework Server
func (svc *AlertChannelsService) UpdateServiceNowRest(data AlertChannel) (
	response ServiceNowRestAlertChannelResponseV2,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

func (snow *ServiceNowRestDataV2) EncodeCustomTemplateFile(template string) {
	encodedTemplate := base64.StdEncoding.EncodeToString([]byte(template))
	snow.CustomTemplateFile = fmt.Sprintf("data:application/json;name=i.json;base64,%s", encodedTemplate)
}

func (snow *ServiceNowRestDataV2) DecodeCustomTemplateFile() (string, error) {
	if len(snow.CustomTemplateFile) == 0 {
		return "", nil
	}

	var (
		b64      = strings.Split(snow.CustomTemplateFile, ",")
		raw, err = base64.StdEncoding.DecodeString(b64[1])
	)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

type ServiceNowRestAlertChannelResponseV2 struct {
	Data ServiceNowRestAlertChannelV2 `json:"data"`
}

type ServiceNowRestAlertChannelV2 struct {
	v2CommonIntegrationData
	Data ServiceNowRestDataV2 `json:"data"`
}

type ServiceNowRestDataV2 struct {
	Username           string `json:"userName"`
	Password           string `json:"password"`
	InstanceURL        string `json:"instanceUrl"`
	CustomTemplateFile string `json:"CustomTemplateFile,omitempty"`
	IssueGrouping      string `json:"issueGrouping,omitempty"`
}
