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
	"net/url"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type alertCmdStateType struct {
	Comment  string
	End      string
	Fixable  bool
	Range    string
	Reason   int
	Scope    string
	Severity string
	Status   string
	Start    string
	Type     string
}

// hasFilters returns true if certain filters are present
// in the command state.  excludes time filters (start, end, range).
func (s alertCmdStateType) hasFilters() bool {
	// severity / status / type filters
	if s.Severity != "" || s.Status != "" || s.Type != "" {
		return true
	}
	return s.Fixable
}

var (
	alertCmdState = alertCmdStateType{}

	// alertCmd represents the alert parent command
	alertCmd = &cobra.Command{
		Use:     "alert",
		Aliases: []string{"alerts"},
		Short:   "Inspect and manage alerts",
		Long: `Inspect and manage alerts.

Lacework provides real-time alerts that are interactive and manageable.
Each alert contains various metadata information, such as severity level, type, status, alert category, and associated tags.

You can also post a comment to an alert's timeline; or change an alert status from Open to Closed.

For more information about alerts, visit:

https://docs.lacework.com/console/alerts-overview

To view all alerts in your Lacework account.

    lacework alert ls

To show an alert.

    lacework alert show <alert_id>

To close an alert.

    lacework alert close <alert_id>
`,
	}

	alertOpenCmd = &cobra.Command{
		Use:   "open <alert_id>",
		Short: "Open a specified alert in a web browser",
		Long:  `Open a specified alert in a web browser.`,
		Args:  cobra.ExactArgs(1),
		RunE:  openAlert,
	}
)

func init() {
	// add the alert command
	rootCmd.AddCommand(alertCmd)

	// add the alert open command
	alertCmd.AddCommand(alertOpenCmd)
}

// Generates a URL similar to:
//
//	=> https://account.lacework.net/ui/investigation/monitor/AlertInbox/123/details?accountName=subaccount
func alertLinkBuilder(id int) string {
	u, err := url.Parse(
		fmt.Sprintf(
			"https://%s.lacework.net/ui/investigation/monitor/AlertInbox/%d/details",
			cli.Account,
			id,
		),
	)
	if err != nil {
		return ""
	}

	q := u.Query()
	q.Set("accountName", cli.Account)
	if cli.Subaccount != "" {
		q.Set("accountName", cli.Subaccount)
	}
	if r := q.Encode(); r != "" {
		u.RawQuery = r
	}
	return u.String()
}

func openAlert(_ *cobra.Command, args []string) error {
	cli.Log.Debugw("opening alert", "alert", args[0])

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("alert ID must be a number")
	}

	// ALLY-1233: Need to switch to alertLinkBuilder when new Alerting UI becomes generally available
	url := alertLinkBuilder(id)

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform\n\nNavigate to %s", url)
	}
	if err != nil {
		return errors.Wrap(err, "unable to open web browser")
	}

	return nil
}
