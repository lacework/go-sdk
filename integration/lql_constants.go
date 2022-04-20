//go:build query || policy

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

package integration

const (
	queryID           string = "CLI_AWS_CTA_IntegrationTest"
	queryHostID       string = "CLI_Host_Files_IntegrationTest"
	queryText         string = "{ source { CloudTrailRawEvents } return { INSERT_ID } }"
	queryUpdateText   string = "{ source { CloudTrailRawEvents } return { INSERT_ID, INSERT_TIME } }"
	queryJSONTemplate string = `{
	"queryID": "%s",
	"queryText": "%s"
}`
	queryURL     string = "https://raw.githubusercontent.com/lacework/go-sdk/main/integration/test_resources/lql/CLI_AWS_CTA_IntegrationTest.yaml"
	queryHostURL string = "https://raw.githubusercontent.com/lacework/go-sdk/main/integration/test_resources/lql/CLI_Host_Files_IntegrationTest.yaml"
)
