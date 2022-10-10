//
// Author:: Nicholas Schmeller (<nick.schmeller@lacework.net>)
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
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// 'agent install ec2ic' flags
	agentInstallAWSEC2ICCmd.Flags().StringVar(&agentCmdState.CTFInfraTagKey,
		"tag_key", "", "only install agents on infra with this tag key set",
	)
	agentInstallAWSEC2ICCmd.Flags().StringSliceVar(&agentCmdState.CTFInfraTag,
		"tag", []string{}, "only install agents on infra with this tag",
	)
	agentInstallAWSEC2ICCmd.Flags().StringVar(&agentCmdState.InstallAgentToken,
		"token", "", "agent access token",
	)
	agentInstallAWSEC2ICCmd.Flags().BoolVar(&agentCmdState.InstallTrustHostKey,
		"trust_host_key", false, "automatically add host keys to the ~/.ssh/known_hosts file",
	)
	agentInstallAWSEC2ICCmd.Flags().StringSliceVarP(&agentCmdState.CTFIncludeRegions,
		"include_regions", "r", []string{}, "list of regions to filter on",
	)
}

func installAWSEC2IC(_ *cobra.Command, args []string) error {
	runners, err := awsDescribeInstances()
	if err != nil {
		return err
	}

	for _, runner := range runners {
		cli.Log.Debugw("runner: ", "runner", runner)
		cli.Log.Debugw("runner user: ", "user", runner.Runner.User)
		cli.Log.Debugw("runner region: ", "region", runner.Region)
		cli.Log.Debugw("runner az: ", "az", runner.AvailabilityZone)
		cli.Log.Debugw("runner instance ID: ", "instance ID", runner.InstanceID)
		cli.Log.Debugw("runner: ", "runner hostname", runner.Runner.Hostname)
		err := runner.SendAndUseIdentityFile()
		if err != nil {
			cli.Log.Debugw("ec2ic failed", "err", err)
			continue
		}

		if err := verifyAccessToRemoteHost(&runner.Runner); err != nil {
			cli.Log.Debugw("verifyAccessToRemoteHost failed")
			return err
		}

		if alreadyInstalled := isAgentInstalledOnRemoteHost(&runner.Runner); alreadyInstalled != nil {
			cli.Log.Debugw("agent already installed on host, skipping")
			continue
		}

		token := agentCmdState.InstallAgentToken
		if token == "" {
			// user didn't provide an agent token
			cli.Log.Debugw("agent token not provided")
			var err error
			token, err = selectAgentAccessToken()
			if err != nil {
				return err
			}
		}
		cmd := fmt.Sprintf("sudo sh -c \"curl -sSL %s | sh -s -- %s\"", agentInstallDownloadURL, token)
		err = runInstallCommandOnRemoteHost(&runner.Runner, cmd)
		if err != nil {
			cli.Log.Debugw("runInstallCommandOnRemoteHost failed")
			return err
		}
	}

	return nil
}
