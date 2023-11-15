//
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

import "github.com/spf13/cobra"

var (
	cdkInitCmd = &cobra.Command{
		Use:               "cdk-init",
		Short:             "internal use only",
		Long:              `Internal command that initializes [[.Component]] to be used by the Lacework CLI as a component.`,
		Hidden:            true,
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		RunE: func(_ *cobra.Command, _ []string) (err error) {
			// This event is only executed after this component is installed
			//
			// Use this event to deploy any necessary file, cache initialization, to install additional
			// libraries packaged within the component. Anything that you need to do before the user
			// starts executing commands.

			//
			// More information at:
			// => https://lacework.atlassian.net/l/cp/Bm0ZrmRG
			return nil
		},
	}

	cdkReconfigureCmd = &cobra.Command{
		Use:               "cdk-reconfigure <current_version> <new_or_older_version>",
		Short:             "internal use only",
		Long:              `Internal command that reconfigures [[.Component]] when it is updated or downgraded.`,
		Hidden:            true,
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		RunE: func(_ *cobra.Command, _ []string) (err error) {
			// There are two places where this event is used;
			//
			//  1) When a component is updated and,
			//  2) When a component is downgraded to an older version (rollback)
			//
			// Rollbacks are needed mostly for when components are updated to a version with a bug or
			// something that prevents the user to continue using the component productively.
			//
			// Note that when this lifecycle event is executed, it receives two arguments, the current
			// version and the version to update or rollback the component.
			//
			// More information at:
			// => https://lacework.atlassian.net/l/cp/Bm0ZrmRG
			return nil
		},
	}

	cdkCleanupCmd = &cobra.Command{
		Use:               "cdk-cleanup",
		Short:             "internal use only",
		Long:              `Internal command that cleans anything that [[.Component]] deployed during cdk-init or cdk-reconfigure events.`,
		Hidden:            true,
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		RunE: func(_ *cobra.Command, _ []string) (err error) {
			// Before this component is removed (uninstalled), this lifecycle event is executed.
			//
			// Here is where you should perform a cleanup of anything that was deployed during the
			// initialization (cdk-init) and/or the reconfiguration (cdk-reconfigure) events.
			//
			// More information at:
			// => https://lacework.atlassian.net/l/cp/Bm0ZrmRG
			return nil
		},
	}
)

func init() {
	// add the commands for the CDK lifecycle events
	rootCmd.AddCommand(cdkInitCmd)
	rootCmd.AddCommand(cdkReconfigureCmd)
	rootCmd.AddCommand(cdkCleanupCmd)
}
