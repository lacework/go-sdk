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

	instances "github.com/lacework/go-sdk/lwcloud/gcp/resources/instances"
	"github.com/lacework/go-sdk/lwrunner"
)

func gcpDescribeInstances(orgID string) ([]*lwrunner.GCPRunner, error) {
	discoveredInstances, err := instances.EnumerateInstancesInOrg(context.Background(), nil, "", "", nil, nil)
	if err != nil {
		return nil, err
	}

	runners := []*lwrunner.GCPRunner{}

	for projectName, projectInstances := range discoveredInstances {
		for _, instance := range projectInstances {
			runner, err := lwrunner.NewGCPRunner(
				instance.PublicIP,
				instance.Zone,
				instance.InstanceID,
				projectName,
				verifyHostCallback,
			)
			if err != nil {
				return nil, err
			}
			runners = append(runners, runner)
		}
	}

	return runners, nil
}
