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
	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// This is an experimental command.
	// policyDisableTagCmd represents the policy disable command
	policyDisableTagCmd = &cobra.Command{
		Use:   "disable [policy_id]",
		Short: "Disable policies",
		Long: `Disable policies by ID or all policies matching a tag.

To disable a single policy by its ID:

	lacework policy disable lacework-policy-id

To disable many policies by ID provide a list of policy ids:

	lacework policy disable lacework-policy-id-one lacework-policy-id-two

To disable all policies for AWS CIS 1.4.0:

	lacework policy disable --tag framework:cis-aws-1-4-0

To disable all policies for GCP CIS 1.3.0:

	lacework policy disable --tag framework:cis-gcp-1-3-0
`,
		PreRunE: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 && policyCmdState.Tag != "" {
				return errors.New("'--tag' flag may not be use in conjunction with 'policy_id' arg")
			}
			policyCmdState.State = ptr.Bool(false)
			return nil
		},
		RunE: setPoliciesState,
	}
)

func init() {
	// add sub-commands to the policy command
	policyCmd.AddCommand(policyDisableTagCmd)

	// policy disable specific flags
	policyDisableTagCmd.Flags().StringVar(
		&policyCmdState.Tag,
		"tag", "", "disable all policies with the specified tag",
	)
}
