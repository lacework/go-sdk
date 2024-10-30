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
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// alertCommentCmd represents the alert comment command
	alertCommentCmd = &cobra.Command{
		Use:   "comment <alert_id>",
		Short: "Add a comment",
		Long: `Post a user comment on an alert's timeline .

Comments may be provided inline or via editor.
`,
		Args: cobra.ExactArgs(1),
		RunE: commentAlert,
	}
)

func init() {
	alertCmd.AddCommand(alertCommentCmd)

	// comment flag
	alertCommentCmd.Flags().StringVarP(
		&alertCmdState.Comment,
		"comment", "c", "",
		"a comment to add to the alert",
	)
}

func commentAlert(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("commenting on alert", "alert", args[0])

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("alert ID must be a number")
	}

	// if comment is not supplied inline
	// validate that the alert exists
	if alertCmdState.Comment == "" {
		exists, err := cli.LwApi.V2.Alerts.Exists(id)
		// if we are very certain the alert doesn't exist
		if !exists && err == nil {
			return errors.New(fmt.Sprintf("alert %d does not exist", id))
		}
	}

	comment, err := inputComment()
	if err != nil {
		return errors.Wrap(err, "unable to process alert comment")
	}

	cli.StartProgress(" Adding alert comment...")
	commentResponse, err := cli.LwApi.V2.Alerts.Comment(id, comment)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to add alert comment")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(commentResponse.Data)
	}

	cli.OutputHuman("Comment added successfully.")
	return nil
}
