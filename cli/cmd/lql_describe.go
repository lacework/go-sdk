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

package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

const (
	lqlDescribeBadInputMsg string = "Please specify a valid data source"
	lqlDescribeDebugMsg    string = "describing LQL data source"
	lqlDescribeNotFoundMsg string = "There is nothing to describe.\n"
	lqlDescribeUnableMsg   string = "unable to describe an LQL data source"
)

var (
	// lqlDescribeCmd represents the lql describe command
	lqlDescribeCmd = &cobra.Command{
		Use:   "describe <data source>",
		Short: "describe an LQL data source",
		Long:  `Describe an LQL data source.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  describeQuerySource,
	}
)

func init() {
	lqlCmd.AddCommand(lqlDescribeCmd)
}

func describeToTable(describeData []api.LQLDescribeData) (out [][]string) {
	if len(describeData) == 0 {
		return
	}

	// "Field Name", Placement", "Type", "Required", "Default"
	for _, param := range describeData[0].Parameters {
		out = append(out, []string{
			param.Name,
			"Parameters",
			param.Type,
			fmt.Sprintf("%v", param.Required),
			param.Default,
		})
	}
	for _, schema := range describeData[0].Schema {
		out = append(out, []string{
			schema.Name,
			"Schema",
			schema.Type,
			"",
			"",
		})
	}
	return
}

func describeQuerySource(_ *cobra.Command, args []string) error {
	var dataSource string

	if len(args) != 0 && args[0] != "" {
		dataSource = args[0]
	} else {
		return errors.Wrap(
			errors.New(lqlDescribeBadInputMsg),
			lqlDescribeUnableMsg,
		)
	}

	cli.Log.Debugw(lqlDescribeDebugMsg, "data source", dataSource)

	describe, err := cli.LwApi.LQL.Describe(dataSource)

	if err != nil {
		return errors.Wrap(err, lqlDescribeUnableMsg)
	} else if cli.JSONOutput() {
		return cli.OutputJSON(describe.Data)
	} else if len(describe.Data) == 0 {
		cli.OutputHuman(lqlDescribeNotFoundMsg)
	} else {
		cli.OutputHuman(
			renderSimpleTable(
				[]string{"Field Name", "Placement", "Type", "Required", "Default"},
				describeToTable(describe.Data),
			),
		)
	}

	return nil
}
