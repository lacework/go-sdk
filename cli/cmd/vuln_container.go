//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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
	"github.com/lacework/go-sdk/internal/array"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

func init() {
	// add sub-commands to the 'vulnerability container' command
	vulContainerCmd.AddCommand(vulContainerScanCmd)
	vulContainerCmd.AddCommand(vulContainerListAssessmentsCmd)
	vulContainerCmd.AddCommand(vulContainerListRegistriesCmd)
	vulContainerCmd.AddCommand(vulContainerShowAssessmentCmd)

	// add start flag to list-assessments command
	vulContainerListAssessmentsCmd.Flags().StringVar(&vulCmdState.Start,
		"start", "-24h", "start of the time range",
	)
	// add end flag to list-assessments command
	vulContainerListAssessmentsCmd.Flags().StringVar(&vulCmdState.End,
		"end", "now", "end of the time range",
	)
	// range time flag
	vulContainerListAssessmentsCmd.Flags().StringVar(&vulCmdState.Range,
		"range", "", "natural time range for query",
	)
	// add repository flag to list-assessments command
	vulContainerListAssessmentsCmd.Flags().StringSliceVarP(&vulCmdState.Repositories,
		"repository", "r", []string{}, "filter assessments for specific repositories",
	)

	// add registry flag to list-assessments command
	vulContainerListAssessmentsCmd.Flags().StringSliceVarP(&vulCmdState.Registries,
		"registry", "", []string{}, "filter assessments for specific registries",
	)

	setPollFlag(
		vulContainerScanCmd.Flags(),
	)

	setHtmlFlag(
		vulContainerScanCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	setDetailsFlag(
		vulContainerScanCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	setSeverityFlag(
		vulContainerScanCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	setFailOnSeverityFlag(
		vulContainerScanCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	setFailOnFixableFlag(
		vulContainerScanCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	setActiveFlag(
		vulContainerListAssessmentsCmd.Flags(),
	)

	setFixableFlag(
		vulContainerScanCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
		vulContainerListAssessmentsCmd.Flags(),
	)

	setPackagesFlag(
		vulContainerScanCmd.Flags(),
		vulContainerShowAssessmentCmd.Flags(),
	)

	setCsvFlag(
		vulContainerShowAssessmentCmd.Flags(),
		vulContainerListAssessmentsCmd.Flags(),
	)

	vulContainerShowAssessmentCmd.Flags().BoolVar(
		&vulCmdState.ImageID, "image_id", false,
		"tread the provided sha256 hash as image id",
	)
}

func setPollFlag(cmds ...*flag.FlagSet) {
	for _, cmd := range cmds {
		if cmd != nil {
			cmd.BoolVar(&vulCmdState.Poll, "poll", false, "poll until the vulnerability scan completes")
		}
	}
}

func getContainerRegistries() ([]string, error) {
	var (
		registries            = make([]string, 0)
		regsIntegrations, err = cli.LwApi.V2.ContainerRegistries.List()
	)
	if err != nil {
		return registries, errors.Wrap(err, "unable to get container registry integrations")
	}

	for _, i := range regsIntegrations.Data {
		// avoid adding empty registries coming from the new local_scanner and avoid adding duplicate registries
		if i.ContainerRegistryDomain() == "" || array.ContainsStr(registries, i.ContainerRegistryDomain()) {
			continue
		}

		registries = append(registries, i.ContainerRegistryDomain())
	}

	return registries, nil
}
