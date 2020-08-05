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

var (
	// lqlCmd represents the lql command
	lqlCmd = &cobra.Command{
		Use:    "lql <query>",
		Hidden: true,
		Short:  "run an LQL query",
		Long: `Run an LQL query.

A simple example of an LQL query:

  $ lacework lql 'SimpleLQL_3(RecentComplianceReports Data) {SELECT Data.*}'

NOTE: This feature is not yet available!`,
		Args: cobra.ExactArgs(1),
		RunE: runLQLQuery,
	}
)

func init() {
	// add the lql command
	rootCmd.AddCommand(lqlCmd)
}

func runLQLQuery(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("running LQL query", "query", args[0])
	response, err := cli.LwApi.LQL.Query(args[0])
	if err != nil {
		return errors.Wrap(err, "unable to run LQL query")
	}

	if data, ok := response["data"]; ok {
		err := cli.OutputJSON(data)
		return err
	}

	if err := cli.OutputJSON(response); err != nil {
		return errors.Wrap(err, "unable to format json response")
	}
	return nil
}
