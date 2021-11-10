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
package integration

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestAlertRuleRead(t *testing.T) {
	alertRule, createErr := createAlertRuleWithSlackAlertChannel()
	if createErr != nil {
		log.Fatal(createErr)
	}
	defer LaceworkCLIWithTOMLConfig("alert-rules", "delete", alertRule.Data.Guid)

	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert-rules", "show", alertRule.Data.Guid)

	expectedHeaders := []string{
		"GUID",
		"NAME",
		"DESCRIPTION",
		"ENABLED",
		"CHANNELS",
		"SEVERITIES",
		"EVENT CATEGORIES",
	}

	expectedFields := []string{
		"Alert Rule Test",
		"This is a test Alert Rule",
		"Critical, High",
		"Compliance",
	}

	t.Run("verify table headers", func(t *testing.T) {
		for _, header := range expectedHeaders {
			assert.Contains(t, out.String(), header, "STDOUT table headers changed, please check")
		}
	})

	t.Run("verify table fields", func(t *testing.T) {
		for _, field := range expectedFields {
			assert.Contains(t, out.String(), field, "STDOUT table fields changed, please check")
		}
	})
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestAlertRulesJsonOutput(t *testing.T) {
	alertRule, createErr := createAlertRuleWithSlackAlertChannel()
	if createErr != nil {
		log.Fatal(createErr)
	}
	defer LaceworkCLIWithTOMLConfig("alert-rules", "delete", alertRule.Data.Guid)

	out, err, exitcode := LaceworkCLIWithTOMLConfig("alert-rules", "list", "--json")

	var rule api.AlertRulesResponse
	jsonErr := json.Unmarshal(out.Bytes(), &rule)

	if assert.NoError(t, jsonErr) {
		assert.NotNil(t, rule.Data[0].Filter.Name)
		assert.Empty(t, err.String(), "STDERR should be empty")
	}
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func createAlertRuleWithSlackAlertChannel() (alertRule api.AlertRuleResponse, err error) {
	var slackChannel string
	lacework, err := api.NewClient(os.Getenv("CI_V2_ACCOUNT"),
		api.WithSubaccount(os.Getenv("CI_ACCOUNT")),
		api.WithApiKeys(os.Getenv("CI_API_KEY"), os.Getenv("CI_API_SECRET")),
		api.WithApiV2(),
	)
	if err != nil {
		return
	}
	channels, err := lacework.V2.AlertChannels.List()

	if err != nil {
		return
	}

	for _, channel := range channels.Data {
		if channel.Name == "#tech-ally-notify" {
			slackChannel = channel.ID()
		}
	}

	rule := api.NewAlertRule("Alert Rule Test", api.AlertRuleConfig{
		Channels:        []string{slackChannel},
		Description:     "This is a test Alert Rule",
		Severities:      api.NewAlertRuleSeverities([]string{"Critical", "High"}),
		EventCategories: []string{"Compliance"},
	})

	return lacework.V2.AlertRules.Create(rule)
}
