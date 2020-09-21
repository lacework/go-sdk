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
	"io/ioutil"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	lqlFile string

	// lqlCmd represents the lql command
	lqlCmd = &cobra.Command{
		Use:    "lql <query>",
		Hidden: true,
		Short:  "run an LQL query",
		Long: `Run an LQL query.

A simple example of an LQL query:

  $ lacework lql 'SimpleLQL_3(RecentComplianceReports Data) {SELECT Data.*}'

NOTE: This feature is not yet available!`,
		Args: cobra.MaximumNArgs(1),
		RunE: runLQLQuery,
	}
)

func init() {
	// add the lql command
	rootCmd.AddCommand(lqlCmd)

	// file flag to specify a query from disk
	lqlCmd.Flags().StringVarP(&lqlFile,
		"file", "f", "",
		"path to an LQL query to run",
	)
}

func runLQLQuery(_ *cobra.Command, args []string) error {
	var query = ""

	if len(args) != 0 && args[0] != "" {
		query = args[0]
	} else if lqlFile != "" {
		lqlQuery, err := ioutil.ReadFile(lqlFile)
		if err != nil {
			return errors.Wrap(err, "unable to read file")
		}
		query = string(lqlQuery)
	} else {
		// avoid asking for a confirmation before launching the editor
		prompt := &survey.Editor{
			Message:  "Type an LQL query to run",
			FileName: "query*.sh",
		}
		err := survey.AskOne(prompt, &query)
		if err != nil {
			return err
		}
	}

	cli.Log.Debugw("running LQL query", "query", query)
	response, err := cli.LwApi.LQL.Query(query)
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
