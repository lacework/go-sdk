//
// Author:: Ross Moles (<ross.moles@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"strings"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// suppressionsListGcpCmd represents the gcp sub-command inside the suppressions list command
	suppressionsListGcpCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List legacy suppressions for GCP",
		RunE:    suppressionsGcpList,
	}
)

func suppressionsGcpList(_ *cobra.Command, _ []string) error {
	var (
		suppressions map[string]api.SuppressionV2
		err          error
	)

	suppressions, err = cli.LwApi.V2.Suppressions.Gcp.List()
	if err != nil {
		if strings.Contains(err.Error(), "No active GCP accounts") {
			cli.OutputHuman("No active GCP accounts found. " +
				"Unable to get legacy GCP suppressions\n")
			return nil
		}
		return errors.Wrap(err, "Unable to get legacy GCP suppressions")
	}

	if len(suppressions) == 0 {
		cli.OutputHuman("No legacy GCP suppressions found.\n")
		return nil
	}
	return cli.OutputJSON(suppressions)
}
