//
// Author:: Fortinet
// Copyright:: Copyright 2026, Fortinet
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
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// repairMethodFunc recovers a missing cloud-account integration for one onboarding method.
type repairMethodFunc func() error

// repairMethods maps a --method value to its handler. To add a new onboarding method, drop a
// cloud_account_repair_<method>.go file that implements the handler and registers its own flags,
// then add one entry here. Method values are named <cloud>-<template>.
var repairMethods = map[string]repairMethodFunc{
	repairMethodAwsCfgOrg: repairAwsCfgOrg,
}

func supportedRepairMethods() []string {
	names := make([]string, 0, len(repairMethods))
	for n := range repairMethods {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// repairMethod is the --method flag; the method-specific flags live in each method's own file.
var repairMethod string

func init() {
	cloudAccountCommand.AddCommand(cloudAccountRepairCmd)

	cloudAccountRepairCmd.Flags().StringVar(&repairMethod, "method", repairMethodAwsCfgOrg,
		"onboarding template to repair (supported: "+strings.Join(supportedRepairMethods(), ", ")+")")
}

var cloudAccountRepairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Re-register a missing cloud-account integration from its onboarding template",
	Long: `Re-register a missing Lacework cloud-account integration for one account onboarded through a
Lacework onboarding template, when the integration is missing - whether the underlying IAM role is
still healthy or was dropped (in which case it is rebuilt first).

The template is selected with --method, named <cloud>-<template>. Only "aws-cfg-org" (AWS Config
organization onboarding) is implemented today; more AWS templates and other clouds (gcp, azure) are
planned and will be added as additional --method values. The rest of this help describes the
aws-cfg-org method.

The org-config setup Lambda registers each member account with Lacework but does NOT reconcile: if a
registration (or the member stack) is dropped, replaying is the only recovery. For aws-cfg-org this
command is state-aware and does the minimum needed for the account:

  - member stack instance present, integration missing -> re-register the integration only
  - member stack instance missing (role gone) but the account is in a targeted OU -> re-create the
    stack instance (which rebuilds the IAM role), wait for it, then register the integration
  - integration already present -> nothing to do
  - account not in any targeted OU, or the management StackSet is gone -> report and do nothing

It registers directly against the Lacework API (deriving the member role ARN and external id from
the stack instance), so you get the result synchronously instead of the setup Lambda's 5-minute
fire-and-forget. When it registers directly the integration is named "<LaceworkAccount>-Config".

Naming note: in the rebuild path (member stack instance missing), re-creating the instance fires the
template's own setup Lambda, which registers the integration first. In that case this command reports
"already onboarded" and the integration keeps the Lambda's name - the account is still fully
recovered, just registered by the Lambda rather than directly.

Run it with credentials for your AWS Organizations management account (the account that owns the
management stack):

    lacework cloud-account repair --account-id 123456789012 --stack-name lacework-aws-cfg-org
    lacework cloud-account repair --account-id 123456789012 --stack-name lacework-aws-cfg-org --dry-run

Creating a stack instance only works when the account already belongs to an OU the StackSet targets.
If it does not, add the account to a targeted OU (or add its OU to the stack's OrganizationalUnits
parameter) first - this command reports that case and does nothing.`,
	Args: cobra.NoArgs,
	RunE: cloudAccountRepair,
}

// cloudAccountRepair dispatches to the handler registered for the selected --method.
func cloudAccountRepair(_ *cobra.Command, _ []string) error {
	run, ok := repairMethods[repairMethod]
	if !ok {
		return errors.Errorf("unsupported --method %q; supported: %s",
			repairMethod, strings.Join(supportedRepairMethods(), ", "))
	}
	return run()
}
