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
	"encoding/json"
	"strings"
)

// GetEmailUser gets a single EmailUser alert channel matching the
// provided integration guid
func (svc *AlertChannelsService) GetEmailUser(guid string) (
	response EmailUserAlertChannelResponse,
	err error,
) {

	// by default, expect the correct response, if not, try the workaround
	err = svc.get(guid, &response)
	if err == nil {
		return
	}

	// Workaround from APIv2
	// Bug: https://lacework.atlassian.net/browse/RAIN-20070
	//
	// This means that the response.Data.Data.ChannelProps.Recipients is a 'string'
	// instead of '[]string'. We will try to deserialize and cast to correct response
	var getResponse emailUserGetAlertChannelResponse
	err = svc.get(guid, &getResponse)
	if err != nil {
		return
	}

	// convert GET response to a consistent response
	response, err = convertGetEmailUserAlertChannelResponse(getResponse)
	return
}

// UpdateEmailUser updates a single EmailUser integration on the Lacework Server
func (svc *AlertChannelsService) UpdateEmailUser(data AlertChannel) (
	response EmailUserAlertChannelResponse,
	err error,
) {
	err = svc.update(data.ID(), data, &response)
	return
}

type EmailUserAlertChannelResponse struct {
	Data EmailUserIntegration `json:"data"`
}

type EmailUserIntegration struct {
	v2CommonIntegrationData
	Data EmailUserData `json:"data"`
}

type EmailUserData struct {
	ChannelProps      EmailUserChannelProps `json:"channelProps"`
	NotificationTypes struct {
		Properties interface{} `json:"properties,omitempty"`
	} `json:"notificationTypes"`
}

type EmailUserChannelProps struct {
	Recipients []string `json:"recipients"`
}

// Workaround from APIv2
// Bug: https://lacework.atlassian.net/browse/RAIN-20070
type emailUserGetData struct {
	ChannelProps struct {
		Recipients interface{} `json:"recipients"`
	} `json:"channelProps"`
	NotificationTypes struct {
		Properties interface{} `json:"properties,omitempty"`
	} `json:"notificationTypes"`
}
type emailUserGetIntegration struct {
	v2CommonIntegrationData
	Data emailUserGetData `json:"data"`
}
type emailUserGetAlertChannelResponse struct {
	Data emailUserGetIntegration `json:"data"`
}

func convertGetEmailUserAlertChannelResponse(
	res emailUserGetAlertChannelResponse) (EmailUserAlertChannelResponse, error) {

	recipientsString, ok := res.Data.Data.ChannelProps.Recipients.(string)
	if ok {
		// deserialize string
		res.Data.Data.ChannelProps.Recipients = strings.Split(recipientsString, ",")
	}

	return castEmailUserAlertChannelResponse(res)
}

func castEmailUserAlertChannelResponse(
	res interface{}) (r EmailUserAlertChannelResponse, err error) {
	var j []byte
	j, err = json.Marshal(res)
	if err != nil {
		return
	}
	err = json.Unmarshal(j, &r)
	return
}
