//
// Author:: Nicholas Schmeller (<nick.schmeller@lacework.net>)
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
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/gammazero/workerpool"
	"github.com/lacework/go-sdk/lwrunner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	agentInstallAWSSSMCmd = &cobra.Command{
		Use:   "ec2ssm",
		Args:  cobra.NoArgs,
		Short: "Use SSM to securely install the Lacework agent on EC2 instances",
		RunE:  installAWSSSM,
		Long: `This command installs the Lacework agent on all EC2 instances in an AWS account using SSM.

This command will create a role and instance profile with 'SSMManagedInstanceCore'
attached and associate that instance profile with the target instances. If the target
instances already have associated instance profiles, this command will not change
their state. This command will teardown the IAM role and instance profile before exiting.

This command authenticates with AWS credentials from well-known locations on the user's
machine. The principal associated with these credentials should have the
'AmazonEC2FullAccess', 'IAMFullAccess' and 'AmazonSSMFullAccess' policies attached.

Target instances must have the SSM agent installed and running for successful
installation.

To skip IAM role / instance profile creation and instance profile association:

    lacework agent aws-install ec2ssm --skip_iam_role_creation

To provide a preexisting IAM role with the 'SSMManagedInstanceCore' policy

    lacework agent aws-install ec2ssm --iam_role_name IAMRoleName

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

To explicitly specify the server URL that the agent will connect to:

    lacework agent aws-install ec2ssm --server_url https://your.server.url.lacework.net

To specify an AWS credential profile other than 'default':

    lacework agent aws-install ec2ssm --credential_profile aws-profile-name

AWS credentials are read from the following environment variables:
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_SESSION_TOKEN (optional)
- AWS_REGION`,
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
	agentInstallAWSSSMCmd.Flags().BoolVar(
		&agentCmdState.InstallSkipCreatInfra,
		"skip_iam_role_creation",
		false,
		"set this flag to skip creating an IAM role and instance profile and associating the instance profile."+
			" Assumes all instances are already setup for SSM",
	)
	agentInstallAWSSSMCmd.Flags().BoolVarP(
		&agentCmdState.InstallDryRun,
		"dry_run",
		"d",
		false,
		"set this flag to print out the target instances and exit",
	)
	agentInstallAWSSSMCmd.Flags().BoolVarP(
		&agentCmdState.InstallForceReinstall,
		"force_reinstall",
		"f",
		false,
		"set this flag to force-reinstall the agent, even if already running on the target instance",
	)
	agentInstallAWSSSMCmd.Flags().StringVar(&agentCmdState.InstallServerURL,
		"server_url", "https://agent.lacework.net", "server URL that agents will talk to, prefixed with `https://`",
	)
	agentInstallAWSSSMCmd.Flags().StringVar(&agentCmdState.InstallAWSProfile,
		"credential_profile", "default", "AWS credential profile to use",
	)
}

func installAWSSSM(_ *cobra.Command, _ []string) error {
	token := agentCmdState.InstallAgentToken
	if token == "" {
		if !cli.InteractiveMode() {
			return errors.New("agent token not provided. Use '--token' when running in non interactive mode")
		}
		// user didn't provide an agent token
		cli.Log.Debug("agent token not provided, asking user to select one now")
		var err error
		token, err = selectAgentAccessToken()
		if err != nil {
			return err
		}
	}

	runners, err := awsDescribeInstances(false /* filter on SSH support */)
	if err != nil {
		return err
	}

	if agentCmdState.InstallDryRun {
		cli.OutputHuman("Dry run, listing target instances for installation\n")
		for _, runner := range runners {
			cli.Log.Info(runner)
			cli.OutputHuman("target instance %v\n", *runner)
		}
		cli.OutputHuman("Dry run finished, exiting now.\n")
		return nil
	}

	cfg, err := config.LoadDefaultConfig(
		context.Background(), config.WithSharedConfigProfile(agentCmdState.InstallAWSProfile),
	)
	if err != nil {
		return err
	}

	var role types.Role
	var instanceProfile types.InstanceProfile
	if !agentCmdState.InstallSkipCreatInfra {
		cli.StartProgress("Setting up IAM role and instance profile...")

		var err error
		role, instanceProfile, err = setupSSMAccess(cfg, agentCmdState.InstallBYORole, token)
		defer func() {
			cli.StopProgress()
			err := teardownSSMAccess(cfg, role, instanceProfile, agentCmdState.InstallBYORole) // clean up after ourselves
			if err != nil {
				cli.OutputHuman("got an error %v while tearing down IAM role / infra\n", err)
				cli.Log.Debugw("IAM infra info after error",
					"role", role,
					"instance profile", instanceProfile,
					"error", err,
				)
			}
		}()
		if err != nil {
			cli.StopProgress()
			return err
		}
		cli.StopProgress()
		cli.OutputHuman(
			"Created role %s with policy %s and instance profile %s, added role to profile\n",
			*role.RoleName,
			lwrunner.SSMInstancePolicy,
			*instanceProfile.InstanceProfileName,
		)
	}

	var successfulCount int32 = 0
	totalCount := len(runners)
	cli.StartProgress(fmt.Sprintf("Installing agents on %d total instances...", totalCount))
	defer cli.StopProgress()

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

			cli.Log.Debugw("runner info",
				"user", threadRunner.Runner.User,
				"region", threadRunner.Region,
				"az", threadRunner.AvailabilityZone,
				"instance_id", threadRunner.InstanceID,
				"hostname", threadRunner.Runner.Hostname,
			)

			if !agentCmdState.InstallSkipCreatInfra {
				// Attach an instance profile with our new role to the instance
				associationID, err := threadRunner.AssociateInstanceProfileWithRunner(cfg, instanceProfile)
				if err != nil {
					cli.OutputHuman(
						"Failed to attach instance profile %s to instance %s with error %v\n",
						*instanceProfile.InstanceProfileName,
						threadRunner.InstanceID,
						err,
					)
					return
				}
				defer func(cfg aws.Config, associationID string) {
					cli.Log.Debugw("disassociating instance profile from runner",
						"association ID", associationID,
						"instance_id", threadRunner.InstanceID,
					)
					err := threadRunner.DisassociateInstanceProfileFromRunner(cfg, associationID)
					if err != nil {
						cli.Log.Debugw("failed to disassociate instance profile from runner",
							"association ID", associationID,
							"instance_id", threadRunner.InstanceID,
							"error", err,
						)
					}
				}(cfg, associationID)
			}

			// Establish SSM Command connection to the runner

			// Check if agent is already installed on the host, skip if yes
			// Sleep for up to 7min to wait for instance profile to associate with instance
			var ssmError error
			var commandOutput ssm.GetCommandInvocationOutput
			const maxSleepTime int = 8
			for i := 0; i < maxSleepTime; i++ {
				const agentVersionCmd = "sudo sh -c '/var/lib/lacework/datacollector -v'"
				commandOutput, ssmError = threadRunner.RunSSMCommandOnRemoteHost(cfg, agentVersionCmd)
				if ssmError != nil {
					cli.Log.Debugw("error when checking if agent already installed on host, retrying",
						"ssmError", ssmError,
						"instance_id", threadRunner.InstanceID,
					)
				} else if commandOutput.Status == ssmtypes.CommandInvocationStatusCancelled ||
					commandOutput.Status == ssmtypes.CommandInvocationStatusTimedOut {
					cli.Log.Debugw("command did not complete successfully, retrying",
						"command output", commandOutput,
						"instance_id", threadRunner.InstanceID,
					)
				} else if commandOutput.Status == ssmtypes.CommandInvocationStatusSuccess {
					if agentCmdState.InstallForceReinstall {
						cli.OutputHuman(
							"Lacework Agent already installed on instance %s, forcing reinstall\n",
							threadRunner.InstanceID,
						)
						break
					} else {
						cli.OutputHuman(
							"Lacework Agent already installed on instance %s, skipping\n",
							threadRunner.InstanceID,
						)
						return
					}
				} else if commandOutput.Status == ssmtypes.CommandInvocationStatusFailed {
					cli.Log.Infow("no agent found on host, proceeding to install",
						"command output", commandOutput,
						"time slept in minutes", i,
						"instance_id", threadRunner.InstanceID,
					)
					cli.OutputHuman(
						"No agent found on instance %s, proceeding to install\n",
						threadRunner.InstanceID,
					)
					break
				} else {
					cli.OutputHuman(
						"Unexpected SSM command exit %v, stderr %s, skipping instance %s\n",
						commandOutput.ResponseCode,
						lwrunner.GetSSMCommandInvocationStdErr(commandOutput),
						threadRunner.InstanceID,
					)
					return
				}

				if i < maxSleepTime-1 { // only sleep when we have a next iteration
					cli.OutputHuman(
						"Waiting for AWS to associate instance profile with instance %s, sleeping 1min, already slept %d min\n",
						threadRunner.InstanceID,
						i,
					)
					time.Sleep(1 * time.Minute)
				}
			}
			if ssmError != nil { // SSM still erroring after 7min of sleep, skip this host
				cli.Log.Warnw("error when checking if agent already installed on host, skipping runner",
					"SSM error", ssmError,
					"command output", commandOutput,
					"instance_id", threadRunner.InstanceID,
				)
				cli.OutputHuman(
					"Error %v when checking if agent already installed on instance %s, skipping\n",
					ssmError,
					threadRunner.InstanceID,
				)
				return
			}

			// Install the agent on the host
			// No need to sleep because instance profile already associated
			const runInstallCmdTmpl = "sudo sh -c 'curl -sSL %s | sh -s -- %s -U %s'"
			runInstallCmd := fmt.Sprintf(runInstallCmdTmpl, agentInstallDownloadURL, token, agentCmdState.InstallServerURL)
			commandOutput, err := threadRunner.RunSSMCommandOnRemoteHost(cfg, runInstallCmd)
			if err != nil {
				cli.OutputHuman(
					"Install failed on instance %s with error %v, stdout %s, stderr %s\n",
					threadRunner.InstanceID,
					err,
					lwrunner.GetSSMCommandInvocationStdOut(commandOutput),
					lwrunner.GetSSMCommandInvocationStdErr(commandOutput),
				)
			} else if commandOutput.Status == ssmtypes.CommandInvocationStatusSuccess {
				cli.OutputHuman("Lacework agent installed successfully on host %s\n\n", threadRunner.InstanceID)
				cli.OutputHuman(fmtSuccessfulAgentInstallString(lwrunner.GetSSMCommandInvocationStdOut(commandOutput)))
				atomic.AddInt32(&successfulCount, 1)
			} else {
				cli.Log.Warnw("Install command did not return `Success` exit status for this instance",
					"instance_id", threadRunner.InstanceID,
					"output", commandOutput,
				)
				cli.OutputHuman(
					"Install command failed with %s exit status, %s stdout, %s stderr for instance %s\n",
					commandOutput.Status,
					lwrunner.GetSSMCommandInvocationStdOut(commandOutput),
					lwrunner.GetSSMCommandInvocationStdErr(commandOutput),
					threadRunner.InstanceID,
				)
			}
		})
		runnerCopyWg.Wait()
	}
	wg.Wait()
	wp.StopWait()

	cli.OutputHuman(
		"Completed installing the Lacework Agent on %d out of %d instances.\n",
		successfulCount,
		totalCount,
	)

	return nil
}
