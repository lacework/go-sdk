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
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/gammazero/workerpool"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	agentInstallAWSSSMCmd = &cobra.Command{
		Use:   "ec2ssm",
		Args:  cobra.NoArgs,
		Short: "Use SSM to securely install on EC2 instances",
		RunE:  installAWSSSM,
		Long: `This command installs the agent on all EC2 instances in an AWS account using SSM.

To filter by one or more regions:

    lacework agent aws-install ec2ssm --include_regions us-west-2,us-east-2

To filter by instance tag:

    lacework agent aws-install ec2ssm --tag TagName,TagValue

To filter by instance tag key:

    lacework agent aws-install ec2ssm --tag_key TagName

To provide an agent access token of your choice, use the command 'lacework agent token list',
select a token and pass it to the '--token' flag. This flag must be selected if the
'--noninteractive' flag is set.

    lacework agent aws-install ec2ssm --token <token>

AWS credentials are read from the following environment variables:
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_SESSION_TOKEN (optional)
- AWS_REGION (optional)`,
	}
)

func init() {
	// 'agent aws-install ec2ssm' flags
	agentInstallAWSSSMCmd.Flags().StringVar(&agentCmdState.InstallTagKey,
		"tag_key", "", "only install agents on infra with this tag key set",
	)
	agentInstallAWSSSMCmd.Flags().StringSliceVar(&agentCmdState.InstallTag,
		"tag", []string{}, "only install agents on infra with this tag",
	)
	agentInstallAWSSSMCmd.Flags().StringSliceVarP(&agentCmdState.InstallIncludeRegions,
		"include_regions", "r", []string{}, "list of regions to filter on",
	)
	agentInstallAWSSSMCmd.Flags().StringVar(&agentCmdState.InstallAgentToken,
		"token", "", "agent access token",
	)
	agentInstallAWSSSMCmd.Flags().IntVarP(
		&agentCmdState.InstallMaxParallelism,
		"max_parallelism",
		"n",
		50,
		"maximum number of workers executing AWS API calls, set if rate limits are lower or higher than normal",
	)
	agentInstallAWSSSMCmd.Flags().StringVar(
		&agentCmdState.InstallBYORole,
		"iam_role_name",
		"",
		"IAM role name (not ARN) with SSM policy, if not provided then an ephemeral role will be created",
	)
}

func installAWSSSM(_ *cobra.Command, _ []string) error {
	token := agentCmdState.InstallAgentToken
	if token == "" {
		if cli.InteractiveMode() {
			// user didn't provide an agent token
			cli.Log.Debugw("agent token not provided, asking user to select one now")
			var err error
			token, err = selectAgentAccessToken()
			if err != nil {
				return err
			}
		} else {
			return errors.New("user did not provide or interactively select an agent token")
		}
	}

	runners, err := awsDescribeInstances()
	if err != nil {
		return err
	}

	cfg, err := GetConfig()
	if err != nil {
		return err
	}
	role, instanceProfile, err := SetupSSMAccess(cfg, agentCmdState.InstallBYORole, token)
	defer func() {
		err := TeardownSSMAccess(cfg, role, instanceProfile, agentCmdState.InstallBYORole) // clean up after ourselves
		cli.Log.Warnw("got an error while tearing down IAM infra", "error", err)
	}()
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	wp := workerpool.New(agentCmdState.InstallMaxParallelism)
	for _, runner := range runners {
		wg.Add(1)

		// In order to use `wp.Submit()`, the input func() must not take any arguments.
		// Copy the runner info to dedicated variable in the goroutine to prevent race overwrite
		runnerCopyWg := new(sync.WaitGroup)
		runnerCopyWg.Add(1)

		wp.Submit(func() {
			defer wg.Done()

			threadRunner := *runner
			runnerCopyWg.Done()

			cli.Log.Debugw("runner info: ",
				"user", threadRunner.Runner.User,
				"region", threadRunner.Region,
				"az", threadRunner.AvailabilityZone,
				"instance_id", threadRunner.InstanceID,
				"hostname", threadRunner.Runner.Hostname,
			)

			// Attach an instance profile with our new role to the runner
			err = threadRunner.AssociateInstanceProfileWithRunner(cfg, instanceProfile)
			if err != nil {
				cli.Log.Debugw("failed to attach instance profile to runner",
					"error", err,
					"instance ID", threadRunner.InstanceID,
					"role", role,
					"instance profile", instanceProfile,
				)
				return
			} else {
				cli.OutputHuman(fmt.Sprintf("successfully associated with instance ID %s\n", threadRunner.InstanceID))
			}

			// Establish SSM Command connection to the runner

			// Check if agent is already installed on the host, skip if yes
			// Sleep for up to 5min to wait for instance profile to associate with instance
			var ssmError error
			var commandOutput ssm.GetCommandInvocationOutput
			for i := 0; i < 5; i++ {
				cli.Log.Debugw("waiting for instance profile to associate with instance, sleeping 1min",
					"iteration number (time slept in minutes)", i,
					"instance ID", threadRunner.InstanceID,
				)
				time.Sleep(1 * time.Minute)

				const agentVersionCmd = "sudo sh -c '/var/lib/lacework/datacollector -v'"
				commandOutput, ssmError = threadRunner.RunSSMCommandOnRemoteHost(cfg, agentVersionCmd)
				if ssmError != nil {
					cli.Log.Debugw("error when checking if agent already installed on host, retrying",
						"ssmError", ssmError,
						"runner", threadRunner.InstanceID,
					)
				} else if commandOutput.Status == ssmtypes.CommandInvocationStatusCancelled ||
					commandOutput.Status == ssmtypes.CommandInvocationStatusTimedOut {
					cli.Log.Debugw("command did not complete successfully, retrying",
						"command output", commandOutput,
						"runner", threadRunner.InstanceID,
					)
				} else if commandOutput.Status == ssmtypes.CommandInvocationStatusSuccess {
					cli.Log.Debugw("agent already installed on host, skipping",
						"runner", threadRunner.InstanceID,
					)
					return
				} else if commandOutput.Status == ssmtypes.CommandInvocationStatusFailed {
					cli.Log.Debugw("no agent found on host, proceeding to install",
						"command output", commandOutput,
						"time slept in minutes", i,
						"runner", threadRunner.InstanceID,
					)
					break
				} else {
					cli.Log.Debugw("unexpected command exit, skipping this runner",
						"command output", commandOutput,
						"runner", threadRunner.InstanceID,
					)
					return
				}
			}
			if ssmError != nil { // SSM still erroring after 5min of sleep, skip this host
				cli.Log.Debugw("error when checking if agent already installed on host, skipping runner",
					"ssmError", ssmError,
					"command output", commandOutput,
					"runner", threadRunner.InstanceID,
				)
				return
			}

			// Install the agent on the host
			// No need to sleep because instance profile already associated
			const runInstallCmdTmpl = "sudo sh -c 'curl -sSL %s | sh -s -- %s'"
			runInstallCmd := fmt.Sprintf(runInstallCmdTmpl, agentInstallDownloadURL, token)
			commandOutput, err := threadRunner.RunSSMCommandOnRemoteHost(cfg, runInstallCmd)
			if err != nil {
				cli.Log.Debugw("runInstallCommandOnRemoteHost failed",
					"error", err,
					"runner", threadRunner.InstanceID,
				)
			} else if commandOutput.Status == ssmtypes.CommandInvocationStatusSuccess {
				cli.OutputHuman("Lacework agent installed successfully on host %s\n\n", threadRunner.InstanceID)
				cli.OutputHuman(fmtSuccessfulAgentInstallString(*commandOutput.StandardOutputContent))
			} else {
				cli.Log.Debugw("Install command did not return `Success` exit status on host",
					"runner", threadRunner.InstanceID,
					"status", commandOutput,
				)
			}
		})
		runnerCopyWg.Wait()
	}
	wg.Wait()
	wp.StopWait()

	return nil
}
