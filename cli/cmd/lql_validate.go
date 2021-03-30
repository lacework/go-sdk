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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	lqlValidateDebugMsg   string = "validating LQL query"
	lqlValidateSuccessMsg string = "LQL query validated successfully.\n"
	lqlValidateUnableMsg  string = "unable to validate LQL query"
)

var (
	// lqlValidateCmd represents the lql validate command
	lqlValidateCmd = &cobra.Command{
		Use:   "validate [query_id]",
		Short: "validate an LQL query",
		Long:  `Validate an LQL query.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  validateQuery,
	}
)

func init() {
	lqlCmd.AddCommand(lqlValidateCmd)

	setQueryFlags(lqlValidateCmd.Flags())
}

func validateQuery(cmd *cobra.Command, args []string) error {
	query, err := inputQuery(cmd, args)
	if err != nil {
		return errors.Wrap(err, lqlValidateUnableMsg)
	}
	return CompileQueryAndOutput(query)
}

func CompileQueryAndOutput(query string) error {
	cli.Log.Debugw(lqlValidateDebugMsg, "query", query)

	compile, err := cli.LwApi.LQL.CompileQuery(query)

	if err != nil {
		return errors.Wrap(err, lqlValidateUnableMsg)
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(compile.Data)
	}
	cli.OutputHuman(lqlValidateSuccessMsg)
	return nil
}
