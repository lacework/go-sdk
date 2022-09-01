//
// Author:: Darren Murray(<darren.murray@lacework.net>)
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
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// vulContainerScanCmd represents the scan sub-command inside the container vulnerability command
	vulContainerScanCmd = &cobra.Command{
		Use:   "scan <registry> <repository> <tag|digest>",
		Short: "Request an on-demand container vulnerability assessment",
		Long: `Request on-demand container vulnerability assessments and view the generated results.

To list all container registries configured in your account:

    lacework vulnerability container list-registries

**NOTE:** Scans can take up to 15 minutes to return results.

Arguments:
    <registry>    container registry where the container image has been published
    <repository>  repository name that contains the container image
    <tag|digest>  either a tag or an image digest to scan (digest format: sha256:1ee...1d3b)
    `,
		Args: cobra.ExactArgs(3),
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateSeverityFlags(); err != nil {
				return err
			}

			err := requestOnDemandContainerVulnerabilityScan(args)
			var e *vulnerabilityPolicyError
			if errors.As(err, &e) {
				c.SilenceUsage = true
			}

			return err
		},
	}
)

func requestOnDemandContainerVulnerabilityScan(args []string) error {
	cli.Log.Debugw("requesting vulnerability scan",
		"registry", args[0],
		"repository", args[1],
		"tag_or_digest", args[2],
	)
	scan, err := cli.LwApi.V2.Vulnerabilities.Containers.Scan(args[0], args[1], args[2])
	if err != nil {
		return userFriendlyErrorForOnDemandCtrVulnScan(err, args[0], args[1], args[2])
	}

	cli.Log.Debugw("vulnerability scan", "details", scan)
	if scan.Data.RequestID == "" {
		return errors.Errorf(
			"there is a problem with the vulnerability scan: %s",
			scan.Message,
		)
	}

	cli.OutputHuman(
		"A new vulnerability scan has been requested. (request_id: %s)\n\n",
		scan.Data.RequestID,
	)

	if cli.JSONOutput() {
		return cli.OutputJSON(scan.Data)
	}
	return nil
}

// Creates a user-friendly error message
func userFriendlyErrorForOnDemandCtrVulnScan(err error, registry, repo, tag string) error {
	if strings.Contains(err.Error(),
		"Could not find integration matching the registry provided",
	) || strings.Contains(err.Error(),
		"Could not find vulnerability integrations",
	) {

		registries, errReg := getContainerRegistries()
		if errReg != nil {
			cli.Log.Debugw("error trying to retrieve configured registries", "error", errReg)
			return errors.Errorf("container registry '%s' not found", registry)
		}

		if len(registries) == 0 {
			msg := `there are no container registries configured in your account.

Get started by integrating your container registry using the command:

    lacework integration create

If you prefer to configure the integration via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Settings > Integrations > Container Registry.
`
			return errors.New(fmt.Sprintf(msg, cli.Account))
		}

		msg := `container registry '%s' not found

Your account has the following container registries configured:

    > %s

To integrate a new container registry use the command:

    lacework integration create
`
		return errors.New(fmt.Sprintf(msg, registry, strings.Join(registries, "\n    > ")))
	}

	if strings.Contains(
		err.Error(),
		"Could not successfully send scan request to available integrations for given repo and label",
	) {

		msg := `container image '%s@%s' not found in registry '%s'.

This error is likely due to a problem with the container registry integration 
configured in your account. Verify that the integration was configured with 
Lacework using the correct permissions, and that the repository belongs
to the provided registry.

To view all container registries configured in your account use the command:

    lacework vulnerability container list-registries
`
		return errors.Errorf(msg, repo, tag, registry)
	}

	return errors.Wrap(err, "unable to request on-demand vulnerability scan")
}
