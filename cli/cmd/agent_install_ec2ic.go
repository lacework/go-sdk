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
	"sync"

	"github.com/korovkin/limiter"
	"github.com/spf13/cobra"
)

var (
	agentInstallAWSEC2ICCmd = &cobra.Command{
		Use:   "ec2ic",
		Args:  cobra.NoArgs,
		Short: "Use EC2InstanceConnect to securely connect to EC2 instances",
		RunE:  installAWSEC2IC,
		Long: `This command installs the agent on all EC2 instances in an AWS account using EC2InstanceConnect.

To filter by one or more regions:

    lacework agent aws-install ec2ic --include_regions us-west-2,us-east-2

To filter by instance tag:

    lacework agent aws-install ec2ic --tag TagName,TagValue

To filter by instance tag key:

    lacework agent aws-install ec2ic --tag_key TagName

To explicitly specify the username for all SSH logins:

    lacework agent aws-install ec2ic --ssh_username <your-user>

AWS credentials are read from the following environment variables:
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_SESSION_TOKEN (optional)
- AWS_REGION (optional)

This command will automatically add hosts with successful connections to
'~/.ssh/known_hosts' unless specified with '--trust_host_key=false'.`,
	}
)

func init() {
	// 'agent install ec2ic' flags
	agentInstallAWSEC2ICCmd.Flags().StringVar(&agentCmdState.InstallTagKey,
		"tag_key", "", "only install agents on infra with this tag key set",
	)
	agentInstallAWSEC2ICCmd.Flags().StringSliceVar(&agentCmdState.InstallTag,
		"tag", []string{}, "only install agents on infra with this tag",
	)
	agentInstallAWSEC2ICCmd.Flags().StringVar(&agentCmdState.InstallAgentToken,
		"token", "", "agent access token",
	)
	agentInstallAWSEC2ICCmd.Flags().BoolVar(&agentCmdState.InstallTrustHostKey,
		"trust_host_key", true, "automatically add host keys to the ~/.ssh/known_hosts file",
	)
	agentInstallAWSEC2ICCmd.Flags().StringSliceVarP(&agentCmdState.InstallIncludeRegions,
		"include_regions", "r", []string{}, "list of regions to filter on",
	)
	agentInstallAWSEC2ICCmd.Flags().StringVar(&agentCmdState.InstallSshUser,
		"ssh_username", "", "username to login with",
	)
	agentInstallAWSEC2ICCmd.Flags().IntVarP(&agentCmdState.InstallMaxParallelism, "max_parallelism", "p", 5, "maximum number of workers executing AWS API calls, set if rate limits are lower or higher than normal")
}

func installAWSEC2IC(_ *cobra.Command, _ []string) error {
	runners, err := awsDescribeInstances()
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	cl := limiter.NewConcurrencyLimiter(agentCmdState.InstallMaxParallelism)
	for _, runner := range runners {
		wg.Add(1)

		// In order to use `cl.Execute()`, the input func() must not take any arguments.
		// Copy the runner info to dedicated variable in the goroutine to prevent race overwrite
		runnerCopyWg := new(sync.WaitGroup)
		runnerCopyWg.Add(1)

		cl.Execute(func() {
			threadRunner := *runner
			runnerCopyWg.Done()
			cli.Log.Debugw("runner info: ",
				"user", threadRunner.Runner.User,
				"region", threadRunner.Region,
				"az", threadRunner.AvailabilityZone,
				"instance ID", threadRunner.InstanceID,
				"hostname", threadRunner.Runner.Hostname,
			)
			err := threadRunner.SendAndUseIdentityFile()
			if err != nil {
				cli.Log.Debugw("ec2ic key send failed", "err", err, "runner", threadRunner.InstanceID)
			}

			if err := verifyAccessToRemoteHost(&threadRunner.Runner); err != nil {
				cli.Log.Debugw("verifyAccessToRemoteHost failed", "err", err, "runner", threadRunner.InstanceID)
			}

			if alreadyInstalled := isAgentInstalledOnRemoteHost(&threadRunner.Runner); alreadyInstalled != nil {
				cli.Log.Debugw("agent already installed on host, skipping", "runner", threadRunner.InstanceID)
			}

			token := agentCmdState.InstallAgentToken
			if token == "" {
				cli.Log.Warnw("agent token not provided", "runner", threadRunner.InstanceID)
			}
			cmd := fmt.Sprintf("sudo sh -c \"curl -sSL %s | sh -s -- %s\"", agentInstallDownloadURL, token)
			err = runInstallCommandOnRemoteHost(&threadRunner.Runner, cmd)
			if err != nil {
				cli.Log.Debugw("runInstallCommandOnRemoteHost failed", "err", err, "runner", threadRunner.InstanceID)
			}
			if threadRunner != *runner {
				cli.Log.Debugw("mutated runner", "threadRunner", threadRunner, "runner", runner)
			}
			wg.Done()
		})
		runnerCopyWg.Wait()
	}
	wg.Wait()
	cl.WaitAndClose()

	return nil
}
