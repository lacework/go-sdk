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
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/compute/metadata"
	instances "github.com/lacework/go-sdk/lwcloud/gcp/resources/instances"
	resources "github.com/lacework/go-sdk/lwcloud/gcp/resources/models"
	"github.com/lacework/go-sdk/lwrunner"
)

// gcpDescribeInstancesInProject takes a GCP project ID and the username of an IAM username in the
// project associated with the credentials in use as input, and outputs a list of GCP instances
// in the project. It reads the flag value `InstallIncludeRegions` if populated to filter on regions,
// and the flag values `InstallTag` and `InstallTagKey` if populated to filter on tag.
func gcpDescribeInstancesInProject(parentUsername, projectID string) ([]*lwrunner.GCPRunner, error) {
	var discoveredInstances []resources.InstanceDetails
	var err error

	// Filter instances by region, if provided as CLI flag value
	if len(agentCmdState.InstallIncludeRegions) > 0 {
		cli.Log.Debugw("filtering on regions", "regions", agentCmdState.InstallIncludeRegions)
		for _, region := range agentCmdState.InstallIncludeRegions {
			discoveredInstances, err = instances.EnumerateInstancesInProject(context.Background(), nil, region, projectID)
			if err != nil {
				return nil, err
			}
		}
	} else {
		discoveredInstances, err = instances.EnumerateInstancesInProject(context.Background(), nil, "", projectID)
		if err != nil {
			return nil, err
		}
	}
	cli.Log.Debugw("found instances", "instances", discoveredInstances)

	runners := []*lwrunner.GCPRunner{}

	for _, instance := range discoveredInstances {
		// Filter out instances that are not in the RUNNING state
		if instance.State != "RUNNING" {
			continue
		}

		// Filter instances by tag and tag key, if provided as CLI flag values
		// NB that tags are another name for GCP "metadata"
		if len(agentCmdState.InstallTag) == 2 {
			cli.Log.Debugw("filtering on tag (metadata)", "tag", agentCmdState.InstallTag)
			if tagVal, ok := instance.Props[agentCmdState.InstallTag[0]]; ok { // is tag key present?
				if tagVal != agentCmdState.InstallTag[1] { // does tag value match?
					continue // skip this instance if filter tag key and value are not present
				}
			} else { // tag key was not present, skip
				continue
			}
		}
		if agentCmdState.InstallTagKey != "" {
			cli.Log.Debugw("filtering on tag (metadata) key", "tag key", agentCmdState.InstallTagKey)
			if _, ok := instance.Props[agentCmdState.InstallTagKey]; !ok {
				continue // skip this instance if filter tag key is not present
			}
		}
		runner, err := lwrunner.NewGCPRunner(
			instance.PublicIP,
			parentUsername,
			projectID,
			instance.Zone,
			instance.InstanceID,
			verifyHostCallback,
		)
		if err != nil {
			return nil, err
		}
		runners = append(runners, runner)
	}

	return runners, nil
}

func gcpGetProjectIDFromMetadataServer() (string, error) {
	client := metadata.NewClient(&http.Client{})

	projectID, err := client.ProjectID()
	if err != nil {
		err = fmt.Errorf("cannot get project details due to %s", err.Error())
		return "", err
	}

	return projectID, nil
}
