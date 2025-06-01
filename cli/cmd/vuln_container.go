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
	"github.com/lacework/go-sdk/v2/api"
	"github.com/lacework/go-sdk/v2/internal/array"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
	"sort"
	"strings"
	"time"
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

	// DEPRECATED
	vulContainerShowAssessmentCmd.Flags().BoolVar(
		&vulCmdState.ImageID, "image_id", false,
		"tread the provided sha256 hash as image id",
	)
	errcheckWARN(vulContainerShowAssessmentCmd.Flags().MarkDeprecated(
		"image_id", "by default we now look up both, image_id and image_digest at once.",
	))
}

func setPollFlag(cmds ...*flag.FlagSet) {
	for _, cmd := range cmds {
		if cmd != nil {
			cmd.BoolVar(&vulCmdState.Poll, "poll", false, "poll until the vulnerability scan completes")
		}
	}
}

// Returns registry domains from the customer's container registry integrations
// This will not return registries scanned by inline and proxy scanner integrations
func getPlatformScannerIntegrationContainerRegistries() ([]string, error) {
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

func getContainerRegistries() ([]string, error) {
	var (
		registries = make([]string, 0)
	)
	// Get registry domains from container registry integrations
	registries, err := getPlatformScannerIntegrationContainerRegistries()
	if err != nil {
		cli.Log.Debugw("error trying to retrieve configured registries", "error", err)
	}

	// Build filter to fetch all container evaluations in the last 7 days. 7 days is an api limitation
	// This is required to find registries that are only scanned by proxy and inline scanner since these
	// integrations don't include the registry domain.
	end := time.Now()
	start := end.Add(-24 * 7 * time.Hour)
	filter := api.SearchFilter{
		TimeFilter: &api.TimeFilter{
			StartTime: &start,
			EndTime:   &end,
		},
		Returns: []string{"evalCtx"},
		Filters: []api.Filter{
			{
				Expression: "not_in",
				Field:      "evalCtx.image_info.registry",
				Values:     registries,
			},
		},
	}
	registryMap := map[string]bool{}
	for {
		response, err := cli.LwApi.V2.Vulnerabilities.Containers.SearchAllPages(filter)
		if err != nil {
			return registries, errors.Wrap(err, "unable to search for container registries")
		}
		// Use a map to get distinct registries from response
		for _, ctr := range response.Data {
			registryMap[strings.TrimSpace(ctr.EvalCtx.ImageInfo.Registry)] = true
		}
		// Convert map to slice
		for reg := range registryMap {
			registries = append(registries, reg)
		}
		// If we hit the 500,000 total row limit from the api, create a filter for registries we haven't seen yet
		if len(response.Data) == 500000 {
			filter.Filters = []api.Filter{
				{
					Expression: "not_in",
					Field:      "evalCtx.image_info.registry",
					Values:     registries,
				},
			}
			continue
		}
		break
	}
	// Sort registries alphabetically
	sort.Strings(registries)
	return registries, nil
}
