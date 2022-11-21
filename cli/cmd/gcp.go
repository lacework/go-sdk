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
	"github.com/lacework/go-sdk/lwrunner"
)

func gcpDescribeInstancesInProject(projectID string) ([]*lwrunner.GCPRunner, error) {
	discoveredInstances, err := instances.EnumerateInstancesInProject(context.Background(), nil, "", projectID)
	if err != nil {
		return nil, err
	}

	runners := []*lwrunner.GCPRunner{}

	for _, instance := range discoveredInstances {
		runner, err := lwrunner.NewGCPRunner(
			instance.PublicIP,
			instance.Zone,
			instance.InstanceID,
			projectID,
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
