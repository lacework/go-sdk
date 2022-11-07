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

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

const (
	ReasonUnset = -1
)

var (
	// alertCloseCmd represents the alert close command
	alertCloseCmd = &cobra.Command{
		Use:   "close <alert_id>",
		Short: "Close an alert",
		Long: `Use this command to change the status of an alert to closed.

The reason for closing the alert must be provided from the following options:

  * 0 - Other
  * 1 - False positive
  * 2 - Not enough information
  * 3 - Malicious and have resolution in place
  * 4 - Expected because of routine testing.

Reasons may be provided inline or via prompt.

If you choose Other, a comment is required and should contain a brief explanation of why the alert is closed.
Comments may be provided inline or via editor.

**Note: A closed alert cannot be reopened. You will be prompted to confirm closure of the alert.  
This prompt can be bypassed with the --noninteractive flag**
`,
		Args: cobra.ExactArgs(1),
		RunE: closeAlert,
	}
)

func init() {
	alertCmd.AddCommand(alertCloseCmd)

	// reason flag
	alertCloseCmd.Flags().IntVarP(
		&alertCmdState.Reason,
		"reason", "r", ReasonUnset,
		"the reason for closing the alert",
	)

	// comment flag
	alertCloseCmd.Flags().StringVarP(
		&alertCmdState.Comment,
		"comment", "c", "",
		"a comment to associate with the alert closure",
	)
}

func inputReason() (reason int, err error) {
	if alertCmdState.Reason != ReasonUnset {
		reason = alertCmdState.Reason
		return
	}

	prompt := &survey.Select{
		Message: "Reason:",
		Options: api.AlertCloseReasons.GetOrderedReasonStrings(),
	}
	err = survey.AskOne(prompt, &reason)

	return
}

func inputComment() (comment string, err error) {
	if alertCmdState.Comment != "" {
		comment = alertCmdState.Comment
		return
	}

	prompt := &survey.Editor{
		Message:  "Type a comment",
		FileName: "alert.comment",
	}
	err = survey.AskOne(prompt, &comment)

	return
}

func closeAlert(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("closing alert", "alert", args[0])

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("alert ID must be a number")
	}

	// if comment is not supplied inline
	// validate that the alert exists
	if alertCmdState.Comment == "" || alertCmdState.Reason == ReasonUnset {
		exists, err := cli.LwApi.V2.Alerts.Exists(id)
		// if we are very certain the alert doesn't exist
		if !exists && err == nil {
			return errors.New(fmt.Sprintf("alert %d does not exist", id))
		}
	}

	reason, err := inputReason()
	if err != nil {
		return errors.Wrap(err, "unable to process alert close reason")
	}

	comment, err := inputComment()
	if err != nil {
		return errors.Wrap(err, "unable to process alert close comment")
	}

	// ask user to confirm
	if !cli.nonInteractive {
		var confirm bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf(
				"Are you sure you want to close alert %d.  Alerts cannot be reopend.", id),
			Default: false,
		}
		err = survey.AskOne(prompt, &confirm)
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	}

	request := api.AlertCloseRequest{
		AlertID: id,
		Reason:  reason,
		Comment: comment,
	}

	cli.StartProgress(" Closing alert...")
	_, err = cli.LwApi.V2.Alerts.Close(request)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to close alert")
	}

	cli.OutputHuman("Alert %d was successfully closed.\n", id)
	return nil
}
