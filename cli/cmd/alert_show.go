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

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// alertShowCmd represents the alert list command
	alertShowCmd = &cobra.Command{
		Use:   "show <alert_id>",
		Short: "Show details about a specific alert",
		Long: `Show details about a specific alert.

There are different types of alert details that can be shown to assist
with alert investigation. These types are referred to as alert detail scopes.

The following alert detail scopes are available:

  * Details (default)
  * Investigation
  * Events
  * RelatedAlerts
  * Integrations
  * Timeline

View an alert's timeline details:

  lacework alert show <alert_id> --scope Timeline
`,
		Args: cobra.ExactArgs(1),
		RunE: showAlert,
	}
)

func init() {
	// show command
	alertCmd.AddCommand(alertShowCmd)

	// scope flag
	alertShowCmd.Flags().StringVar(
		&alertCmdState.Scope,
		"scope", api.AlertDetailsScope.String(),
		"type of alert details to show",
	)
}

func showAlert(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("showing alert", "alert", args[0])

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("alert ID must be a number")
	}

	switch alertCmdState.Scope {
	case api.AlertDetailsScope.String():
		err = showAlertDetails(id)
	case api.AlertInvestigationScope.String():
		err = showAlertInvestigation(id)
	case api.AlertEventsScope.String():
		err = showAlertEvents(id)
	case api.AlertRelatedAlertsScope.String():
		err = showRelatedAlerts(id)
	case api.AlertIntegrationsScope.String():
		err = showAlertIntegrations(id)
	case api.AlertTimelineScope.String():
		err = showAlertTimeline(id)
	default:
		err = errors.New(fmt.Sprintf("scope (%s) is not recognized", alertCmdState.Scope))
	}
	if err != nil {
		return err
	}

	// breadcrumb
	if !cli.JSONOutput() {
		url := alertLinkBuilder(id)
		cli.OutputHuman(
			fmt.Sprintf(
				"\nFor further investigation of this alert navigate to %s\n",
				url,
			),
		)
	}
	return nil
}
