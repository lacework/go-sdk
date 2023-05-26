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
	agentInstallGCPOSLCmd = &cobra.Command{
		Use:   "osl",
		Args:  cobra.ExactArgs(1),
		Short: "Use OSLogin to securely connect to GCE instances",
		RunE:  installGCPOSL,
		Long: `This command installs the agent on all GCE instances in a GCP organization using OSLogin.

The username of the GCP user or service account, in the form ` + "`users/<username>`" + `, is a
required argument.

This command will attempt to query the GCE metadata server for the current project. If this
command is not run on a GCE instance, pass the project ID as:

    lacework agent gcp-install osl <gcp_username> --project_id my-project-id

To filter by one or more regions:

    lacework agent gcp-install osl <gcp_username> --include_regions us-west1,europe-west2

To filter by instance metadata:

    lacework agent gcp-install osl <gcp_username> --metadata MetadataKey,MetadataValue

To filter by instance metadata key:

    lacework agent gcp-install osl <gcp_username> --metadata_key MetadataKey

To provide an agent access token of your choice, use the command 'lacework agent token list',
select a token and pass it to the '--token' flag. This flag must be selected if the
'--noninteractive' flag is set.

    lacework agent gcp-install osl <gcp_username> --token <token>

To explicitly specify the server URL that the agent will connect to:

    lacework agent gcp-install osl --server_url https://your.server.url.lacework.net

GCP credentials are read using the following environment variables:
- GOOGLE_APPLICATION_CREDENTIALS

This command will automatically add hosts with successful connections to
'~/.ssh/known_hosts' unless specified with '--trust_host_key=false'.`,
	}
)

func init() {
	// 'agent gcp-install osl' flags
	agentInstallGCPOSLCmd.Flags().StringVar(
		&agentCmdState.InstallTagKey,
		"metadata_key",
		"",
		"only install agents on infra with this metadata key set",
	)
	agentInstallGCPOSLCmd.Flags().StringSliceVar(
		&agentCmdState.InstallTag,
		"metadata",
		[]string{},
		"only install agents on infra with this metadata",
	)
	agentInstallGCPOSLCmd.Flags().StringSliceVarP(
		&agentCmdState.InstallIncludeRegions,
		"include_regions",
		"r",
		[]string{},
		"list of regions to filter on",
	)
	agentInstallGCPOSLCmd.Flags().BoolVar(
		&agentCmdState.InstallTrustHostKey,
		"trust_host_key",
		true,
		"automatically add host keys to the ~/.ssh/known_hosts file",
	)
	agentInstallGCPOSLCmd.Flags().IntVarP(
		&agentCmdState.InstallMaxParallelism,
		"max_parallelism",
		"n",
		50,
		"maximum number of workers executing GCP API calls, set if rate limits are lower or higher than normal",
	)
	agentInstallGCPOSLCmd.Flags().StringVar(
		&agentCmdState.InstallProjectId,
		"project_id",
		"",
		"ID of the GCP project, set if metadata server does not provide",
	)
	agentInstallGCPOSLCmd.Flags().StringVar(
		&agentCmdState.InstallAgentToken,
		"token",
		"",
		"agent access token",
	)
	agentInstallGCPOSLCmd.Flags().StringVar(
		&agentCmdState.InstallServerURL,
		"server_url",
		"https://agent.lacework.net",
		"server URL that agents will talk to, prefixed with `https://`",
	)
}

func installGCPOSL(_ *cobra.Command, args []string) error {
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

	var projectID string
	if agentCmdState.InstallProjectId != "" {
		projectID = agentCmdState.InstallProjectId // prioritize CLI flag
	} else if mdProjID, err := gcpGetProjectIDFromMetadataServer(); mdProjID != "" && err == nil {
		projectID = mdProjID // if flag not passed, check the metadata server
	} else {
		return fmt.Errorf("could not find project ID, no metadata server (%v) and ID not passed as flag", err)
	}

	runners, err := gcpDescribeInstancesInProject(args[0], projectID)
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
				"zone", threadRunner.AvailabilityZone,
				"instance_id", threadRunner.InstanceID,
				"hostname", threadRunner.Runner.Hostname,
			)
			err := threadRunner.SendAndUseIdentityFile()
			if err != nil {
				cli.Log.Debugw("osl key send failed", "err", err, "runner", threadRunner.InstanceID)
				return
			}

			if err := verifyAccessToRemoteHost(&threadRunner.Runner); err != nil {
				cli.Log.Debugw("verifyAccessToRemoteHost failed", "err", err, "runner", threadRunner.InstanceID)
				return
			}

			if alreadyInstalled := isAgentInstalledOnRemoteHost(&threadRunner.Runner); alreadyInstalled != nil {
				cli.Log.Debugw("agent already installed on host, skipping", "runner", threadRunner.InstanceID)
				return
			}

			cmd := fmt.Sprintf("sudo sh -c \"curl -sSL %s | sh -s -- %s -U %s\"", agentInstallDownloadURL, token, agentCmdState.InstallServerURL)
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
