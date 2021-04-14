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

var (
	// lqlDescribeCmd represents the lql describe command
	lqlDescribeCmd = &cobra.Command{
		Use:   "describe <data_source>",
		Short: "describe an LQL data source",
		Long:  `Describe an LQL data source.`,
		Args:  cobra.ExactArgs(1),
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
			fmt.Sprintf("%t", param.Required),
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
	cli.Log.Debugw("describing LQL data source", "data source", args[0])

	describe, err := cli.LwApi.LQL.Describe(args[0])

	if err != nil {
		return errors.Wrap(err, "unable to describe LQL data source")
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(describe.Data)
	}
	if len(describe.Data) == 0 {
		return yikes("unable to describe LQL data source")
	}
	cli.OutputHuman(
		renderSimpleTable(
			[]string{"Field Name", "Placement", "Type", "Required", "Default"},
			describeToTable(describe.Data),
		),
	)
	return nil
}
