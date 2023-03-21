//go:build report_definitions && !windows

// Author:: Darren Murray (<dmurray-lacework@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestReportDefinitionsPromptUpdate(t *testing.T) {
	report := fetchCustomReportDefinition()

	if report == nil {
		t.Skip("skipping report definitions update. No custom reports found.")
	}

	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	dir := createTOMLConfigFromCIvars()
	defer os.RemoveAll(dir)

	suffix := time.Now().UnixMilli()
	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.UpdateReportDefinitionQuestion, "n"},
				MsgRsp{cmd.UpdateReportDefinitionReportNameQuestion, "CLI Test Updated"},
				MsgRsp{cmd.UpdateReportDefinitionDisplayNameQuestion, fmt.Sprintf("Test CLI Update-%d", suffix)},
				MsgRsp{cmd.UpdateReportDefinitionEditSectionQuestion, "n"},
				MsgRsp{cmd.UpdateReportDefinitionAddSectionQuestion, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"report-definition",
		"update",
		report.ReportDefinitionGuid,
	)

	assert.Contains(t, final, "Report definition updated")
}

func TestReportDefinitionsUpdateFromFile(t *testing.T) {
	report := fetchCustomReportDefinition()

	if report == nil {
		t.Skip("skipping report definitions update. No custom reports found.")
	}

	customReportYaml, err := yaml.Marshal(report.Config())
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	file, err := createTemporaryFile("TestReportDefinitionUpdateFile", string(customReportYaml))
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	out, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"rd", "update", report.ReportDefinitionGuid, "--file", file.Name())

	assert.Empty(t, stderr.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	assert.Contains(t, out.String(), "Report definition updated")
}

func fetchCustomReportDefinition() *api.ReportDefinition {
	lacework, err := api.NewClient(os.Getenv("CI_ACCOUNT"),
		api.WithSubaccount(os.Getenv("CI_SUBACCOUNT")),
		api.WithApiKeys(os.Getenv("CI_API_KEY"), os.Getenv("CI_API_SECRET")),
		api.WithApiV2(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// List all report definitions
	reports, err := lacework.V2.ReportDefinitions.List()
	if err != nil {
		log.Fatal(err)
	}

	// Select custom report
	for _, r := range reports.Data {
		if r.IsCustom() && strings.Contains(r.ReportName, "CLI Test") {
			return &r
		}
	}

	// if test custom report doesn't exist create one
	rptCfg := api.ReportDefinitionConfig{
		ReportName:    "CLI Test",
		DisplayName:   "CLI Test",
		ReportType:    api.ReportDefinitionTypeCompliance.String(),
		SubReportType: api.ReportDefinitionSubTypeAws.String(),
		Sections: []api.ReportDefinitionSection{{
			Title:    "CLI Test Report",
			Policies: []string{"lacework-global-1"},
		}},
	}

	customReport := api.NewReportDefinition(rptCfg)

	report, err := lacework.V2.ReportDefinitions.Create(customReport)
	if err != nil {
		log.Fatal(err)
	}

	return &report.Data
}
