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
	"github.com/spf13/cobra"
)

var (
	alertCmdState = struct {
		Comment string
		Reason  int
	}{}

	// alertCmd represents the lql parent command
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
)

func init() {
	// add the lql command
	rootCmd.AddCommand(alertCmd)
}
