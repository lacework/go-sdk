//go:build report_definitions && !windows

//
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
	"os"
	"strings"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// Unstable test disabled as part of GROW-1396
func _TestReportDefinitionsPromptCreate(t *testing.T) {
	os.Setenv("LW_NOCACHE", "true")
	defer os.Setenv("LW_NOCACHE", "")
	var final string

	dir := createTOMLConfigFromCIvars()
	defer os.RemoveAll(dir)

	runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectsCliOutput(t, c, []MsgRspHandler{
				MsgRsp{cmd.CreateReportDefinitionQuestion, "n"},
				MsgRsp{cmd.CreateReportDefinitionReportNameQuestion, "TestCLIPromptReportName"},
				MsgRsp{cmd.CreateReportDefinitionDisplayNameQuestion, "Test CLI Prompt Display Name"},
				MsgRsp{cmd.CreateReportDefinitionReportSubTypeQuestion, "AWS"},
				MsgRsp{cmd.CreateReportDefinitionSectionTitleQuestion, "CLI Prompt Test"},
				Select{cmd.CreateReportDefinitionPoliciesQuestion},
				MsgRsp{cmd.CreateReportDefinitionAddSectionQuestion, "n"},
			})
			final, _ = c.ExpectEOF()
		},
		"report-definition",
		"create",
	)

	assert.Contains(t, final, "New report definition created")

	assert.NoError(t, removeReportDefinitionResponse(final))
}

func removeReportDefinitionResponse(final string) error {
	res := strings.Split(final, "lacework report-definition show ")
	fmt.Println(res[1])
	id := strings.TrimSuffix(res[1], " \r\n")

	_, _, exitcode := LaceworkCLIWithTOMLConfig(
		"rd", "delete", id,
	)

	if exitcode != 0 {
		return errors.Errorf("unable to remove report definition %s", id)
	}
	return nil
}
