//
// Author:: Darren Murray(<darren.murray@lacework.net>)
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

package cmd

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// revert command is used to rollback lacework report definition to a previous version
var reportDefinitionsRevertCommand = &cobra.Command{
	Use:   "revert <report_definition_id> <version>",
	Short: "Update a report definition",
	Long: `Update an existing custom report definition.

To revert a report definition:

    lacework report-definition revert <report_definition_id> <version>
`,
	Args: cobra.ExactArgs(2),
	RunE: revertReportDefinition,
}

func revertReportDefinition(_ *cobra.Command, args []string) error {
	var (
		err     error
		version int
	)

	if version, err = strconv.Atoi(args[1]); err != nil {
		return errors.Wrap(err, "unable to parse version")
	}

	cli.StartProgress("Reverting report definition...")
	resp, err := cli.LwApi.V2.ReportDefinitions.Revert(args[0], version)
	cli.StopProgress()

	if err != nil {
		return err
	}

	cli.OutputHuman("The report definition %s was reverted to version %d \n", resp.Data.ReportDefinitionGuid, version)
	return nil
}
