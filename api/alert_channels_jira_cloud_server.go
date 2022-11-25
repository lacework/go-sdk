//
// Author:: Vatasha White (<vatasha.white@lacework.net>)
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

type jiraIssueGrouping int

const (
	NoneJiraIssueGrouping jiraIssueGrouping = iota
	EventsJiraIssueGrouping
	ResourcesJiraIssueGrouping
)

var JiraIssueGroupings = map[jiraIssueGrouping]string{
	NoneJiraIssueGrouping:      "",
	EventsJiraIssueGrouping:    "Events",
	ResourcesJiraIssueGrouping: "Resources",
}

var JiraIssueGroupingsSurvey = map[string]jiraIssueGrouping{
	"None":      NoneJiraIssueGrouping,
	"Events":    EventsJiraIssueGrouping,
	"Resources": ResourcesJiraIssueGrouping,
}

func (i jiraIssueGrouping) String() string {
	return JiraIssueGroupings[i]
}

const (
	BidirectionalJiraConfiguration = "Bidirectional"
	JiraCloudAlertType             = "JIRA_CLOUD"
	JiraServerAlertType            = "JIRA_SERVER"
)

// GetJira gets a single instance of a Jira Cloud or Jira Server alert channel with the corresponding guid
func (svc *AlertChannelsService) GetJira(guid string) (response JiraAlertChannelResponseV2, err error) {
	err = svc.get(guid, &response)
	return
}

// UpdateJira updates a single instance of a Jira Cloud or Jira Server integration on the Lacework server
func (svc *AlertChannelsService) UpdateJira(data AlertChannel) (response JiraAlertChannelResponseV2, err error) {
	err = svc.update(data.ID(), data, &response)
	return
}

type JiraDataV2 struct {
	ApiToken           string `json:"apiToken,omitempty"` // used for Jira Cloud
	CustomTemplateFile string `json:"customTemplateFile,omitempty"`
	IssueGrouping      string `json:"issueGrouping,omitempty"`
	IssueType          string `json:"issueType"`
	JiraType           string `json:"jiraType"`
	JiraUrl            string `json:"jiraUrl"`
	ProjectID          string `json:"projectId"`
	Username           string `json:"username"`
	Password           string `json:"password,omitempty"`            // used for Jira Server
	Configuration      string `json:"bidirectionalConfig,omitempty"` // used for bidirectional integration
}

type JiraAlertChannelV2 struct {
	v2CommonIntegrationData
	Data JiraDataV2 `json:"data"`
}

type JiraAlertChannelResponseV2 struct {
	Data JiraAlertChannelV2 `json:"data"`
}

func (jira *JiraDataV2) EncodeCustomTemplateFile(template string) {
	encodedTemplate := base64.StdEncoding.EncodeToString([]byte(template))
	jira.CustomTemplateFile = fmt.Sprintf("data:application/json;name=i.json;base64,%s", encodedTemplate)
}

func (jira *JiraDataV2) DecodeCustomTemplateFile() (string, error) {
	if len(jira.CustomTemplateFile) == 0 {
		return "", nil
	}

	var (
		b64      = strings.Split(jira.CustomTemplateFile, ",")
		raw, err = base64.StdEncoding.DecodeString(b64[1])
	)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}
