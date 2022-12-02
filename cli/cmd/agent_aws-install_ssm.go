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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	agentInstallAWSSSMCmd = &cobra.Command{
		Use:   "ssm",
		Args:  cobra.NoArgs,
		Short: "Use SSM to securely install on EC2 instances",
		RunE:  installAWSSSM,
		Long: `This command installs the agent on all EC2 instances in an AWS account using SSM.

To filter by one or more regions:

    lacework agent aws-install ssm --include_regions us-west-2,us-east-2

To filter by instance tag:

    lacework agent aws-install ssm --tag TagName,TagValue

To filter by instance tag key:

    lacework agent aws-install ssm --tag_key TagName

To provide an agent access token of your choice, use the command 'lacework agent token list',
select a token and pass it to the '--token' flag. This flag must be selected if the
'--noninteractive' flag is set.

    lacework agent aws-install ssm --token <token>

AWS credentials are read from the following environment variables:
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_SESSION_TOKEN (optional)
- AWS_REGION (optional)`,
	}
)

func init() {
	// 'agent aws-install ssm' flags
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
		"IAM role name (not ARN) to use for SSM, if not provided then an ephemeral role will be created",
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
	role, instanceProfile, err := SetupSSMAccess(cfg, agentCmdState.InstallBYORole)
	// defer TeardownSSMAccess(cfg, role, instanceProfile, agentCmdState.InstallBYORole) // clean up after ourselves
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

			// TODO Remove me, debug code
			// BEGIN DEBUG CODE
			// c := ec2.New(ec2.Options{
			// 	Credentials: cfg.Credentials,
			// 	Region:      cfg.Region,
			// })
			// // test that the instance profile association fails
			// associateInput := &ec2.AssociateIamInstanceProfileInput{ // PLOT TWIST IT WORKED
			// 	IamInstanceProfile: &ec2types.IamInstanceProfileSpecification{
			// 		Arn: instanceProfile.Arn,
			// 		// Arn: aws.String("arn:aws:iam::561021084946:instance-profile/Lacework-Agent-SSM-Install-Instance-Profile"),
			// 		// Arn: aws.String("arn:aws:iam::561021084946:instance-profile/Test-Debugging-Instance-Profile-Additional-Chars-Lacework-Agent-SSM-Install-Instance-Profile"),
			// 	},
			// 	InstanceId: aws.String(threadRunner.InstanceID),
			// }
			// associateOutput, err := c.AssociateIamInstanceProfile(context.Background(), associateInput)
			// // the association might have failed because the instance profile didn't exist
			// // look it up again to see if this is true
			// getInstanceProfileInput := &iam.GetInstanceProfileInput{
			// 	InstanceProfileName: instanceProfile.InstanceProfileName,
			// }
			// iamClient := iam.New(iam.Options{
			// 	Credentials: cfg.Credentials,
			// 	Region:      cfg.Region,
			// })
			// instanceProfOut, _ := iamClient.GetInstanceProfile(context.Background(), getInstanceProfileInput)
			// cli.Log.Debugw("DEBUG instance profile associations", "output", associateOutput, "error", err, "instance profile", instanceProfile, "fresh lookup of instance profile", instanceProfOut)
			// END DEBUG CODE

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
			}

			// TODO establish SSH access / SSM Command connection to the runner

			if err := verifyAccessToRemoteHost(&threadRunner.Runner); err != nil {
				cli.Log.Debugw("verifyAccessToRemoteHost failed", "err", err, "runner", threadRunner.InstanceID)
				return
			}

			if alreadyInstalled := isAgentInstalledOnRemoteHost(&threadRunner.Runner); alreadyInstalled != nil {
				cli.Log.Debugw("agent already installed on host, skipping", "runner", threadRunner.InstanceID)
				return
			}

			cmd := fmt.Sprintf("sudo sh -c \"curl -sSL %s | sh -s -- %s\"", agentInstallDownloadURL, token)
			err = runInstallCommandOnRemoteHost(&threadRunner.Runner, cmd)
			if err != nil {
				cli.Log.Debugw("runInstallCommandOnRemoteHost failed", "err", err, "runner", threadRunner.InstanceID)
			}
			if threadRunner != *runner {
				cli.Log.Debugw("mutated runner", "thread_runner", threadRunner, "runner", runner)
			}
		})
		runnerCopyWg.Wait()
	}
	wg.Wait()
	wp.StopWait()

	return nil
}
