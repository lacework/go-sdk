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

	"github.com/gammazero/workerpool"
	"github.com/spf13/cobra"
)

var (
	agentInstallAWSSSHCmd = &cobra.Command{
		Use:   "ec2ssh <token>",
		Args:  cobra.ExactArgs(1),
		Short: "Use SSH to securely connect to EC2 instances",
		Long: `This command installs the agent on all EC2 instances in an AWS account
using SSH.

To filter by one or more regions:

    lacework agent aws-install ec2ssh <token> --include_regions us-west-2,us-east-2

To filter by instance tag:

    lacework agent aws-install ec2ssh <token> --tag TagName,TagValue

To filter by instance tag key:

    lacework agent aws-install ec2ssh <token> --tag_key TagName

You will need to provide an SSH authentication method. This authentication method
should work for all instances that your tag or region filters select. Instances must
be routable from your local host.

To authenticate using username and password:

    lacework agent aws-install ec2ssh <token> --ssh_username <your-user> --ssh_password <secret>

To authenticate using an identity file:

    lacework agent aws-install ec2ssh <token> -i /path/to/your/key

The environment should contain AWS credentials in the following variables:
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_SESSION_TOKEN (optional),
- AWS_REGION (optional)

This command will automatically add hosts with successful connections to
'~/.ssh/known_hosts' unless specified with '--trust_host_key=false'.`,
		RunE: installAWSSSH,
	}
)

func init() {
	// 'agent install ec2ssh' flags
	agentInstallAWSSSHCmd.Flags().StringVar(&agentCmdState.InstallTagKey,
		"tag_key", "", "only install agents on infra with this tag key",
	)
	agentInstallAWSSSHCmd.Flags().StringSliceVar(&agentCmdState.InstallTag,
		"tag", []string{}, "only select instances with this tag",
	)
	agentInstallAWSSSHCmd.Flags().StringVarP(&agentCmdState.InstallIdentityFile,
		"identity_file", "i", defaultSshIdentityKey,
		"identity (private key) for public key authentication",
	)
	agentInstallAWSSSHCmd.Flags().BoolVar(&agentCmdState.InstallTrustHostKey,
		"trust_host_key", true, "automatically add host keys to the ~/.ssh/known_hosts file",
	)
	agentInstallAWSSSHCmd.Flags().StringSliceVarP(&agentCmdState.InstallIncludeRegions,
		"include_regions", "r", []string{}, "list of regions to filter on",
	)
	agentInstallAWSSSHCmd.Flags().StringVar(&agentCmdState.InstallPassword,
		"ssh_password", "", "password for authentication",
	)
	agentInstallAWSSSHCmd.Flags().StringVar(&agentCmdState.InstallSshUser,
		"ssh_username", "", "username to login with",
	)
	agentInstallAWSSSHCmd.Flags().IntVar(&agentCmdState.InstallSshPort,
		"ssh_port", 22, "port to connect to on the remote host",
	)
	agentInstallAWSSSHCmd.Flags().IntVarP(
		&agentCmdState.InstallMaxParallelism,
		"max_parallelism",
		"n",
		50,
		"maximum number of workers executing AWS API calls, set if rate limits are lower or higher than normal",
	)
}

func installAWSSSH(_ *cobra.Command, args []string) error {
	runners, err := awsDescribeInstances()
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	wp := workerpool.New(agentCmdState.InstallMaxParallelism)
	for _, runner := range runners {
		wg.Add(1)

		// In order to use `cl.Execute()`, the input func() must not take any arguments.
		// Copy the runner info to dedicated variable in the goroutine to prevent race overwrite
		runnerCopyWg := new(sync.WaitGroup)
		runnerCopyWg.Add(1)

		wp.Submit(func() {
			threadRunner := *runner
			runnerCopyWg.Done()
			cli.Log.Debugw("threadRunner info: ",
				"user", threadRunner.Runner.User,
				"region", threadRunner.Region,
				"az", threadRunner.AvailabilityZone,
				"instance_id", threadRunner.InstanceID,
				"hostname", threadRunner.Runner.Hostname,
			)

			err := threadRunner.Runner.UseIdentityFile(agentCmdState.InstallIdentityFile)
			if err != nil {
				cli.Log.Warnw("unable to use provided identity file", "err", err, "thread_runner", threadRunner.InstanceID)
			}

			if err := verifyAccessToRemoteHost(&threadRunner.Runner); err != nil {
				cli.Log.Debugw("verifyAccessToRemoteHost failed", "err", err, "thread_runner", threadRunner.InstanceID)
			}

			if alreadyInstalled := isAgentInstalledOnRemoteHost(&threadRunner.Runner); alreadyInstalled != nil {
				cli.Log.Debugw("agent already installed on host, skipping", "thread_runner", threadRunner.InstanceID)
			}

			var token string
			if len(args) <= 0 || args[0] == "" {
				// user didn't provide an agent token
				cli.Log.Warnw("agent token not provided", "thread_runner", threadRunner.InstanceID)
			} else {
				token = args[0]
			}
			cmd := fmt.Sprintf("sudo sh -c \"curl -sSL %s | sh -s -- %s\"", agentInstallDownloadURL, token)
			err = runInstallCommandOnRemoteHost(&threadRunner.Runner, cmd)
			if err != nil {
				cli.Log.Debugw("runInstallCommandOnRemoteHost failed", "thread_runner", threadRunner.InstanceID)
			}
			wg.Done()
		})
		runnerCopyWg.Wait()
	}

	wg.Wait()
	wp.StopWait()

	return nil
}
