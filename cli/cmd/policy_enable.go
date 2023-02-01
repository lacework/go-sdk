//
// Author:: Darren Murray(<darren.murray@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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
	"github.com/lacework/go-sdk/internal/pointer"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// policyEnableTagCmd represents the policy enable command
	policyEnableTagCmd = &cobra.Command{
		Use:   "enable [policy_id]",
		Short: "Enable policies",
		Long: `Enable policies by ID or all policies matching a tag.
To enter the policy enable prompt:

	lacework policy enable

To enable a single policy by its ID:

	lacework policy enable lacework-policy-id

To enable many policies by ID provide a list of policy ids:

	lacework policy enable lacework-policy-id-one lacework-policy-id-two

To enable all policies for AWS CIS 1.4.0:

	lacework policy enable --tag framework:cis-aws-1-4-0

To enable all policies for GCP CIS 1.3.0:

	lacework policy enable --tag framework:cis-gcp-1-3-0

`,
		PreRunE: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 && policyCmdState.Tag != "" {
				return errors.New("'--tag' flag may not be use in conjunction with 'policy_id' arg")
			}
			policyCmdState.State = pointer.BoolPtr(true)
			return nil
		},
		RunE: setPoliciesState,
	}
)

func init() {
	// add sub-commands to the policy command
	policyCmd.AddCommand(policyEnableTagCmd)

	// policy enable specific flags
	policyEnableTagCmd.Flags().StringVar(
		&policyCmdState.Tag,
		"tag", "", "enable all policies with the specified tag",
	)
}
