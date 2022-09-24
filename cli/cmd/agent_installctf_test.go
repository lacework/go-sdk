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
	"os"
	"sync"
	"testing"

	"github.com/lacework/go-sdk/lwrunner"
	"github.com/stretchr/testify/assert"
)

// Requires AWS credentials in the shell environment
// Lists runners, sends keys, attempts to connect
// Example command to run:
// `aws-vault exec default -- go test -run TestAwsEC2ICFindRunnersToCapture`
// If AWS credentials are already present in the shell environment, only use:
// `go test -run TestAwsEC2ICFindRunnersToCapture`
func TestAwsEC2ICFindRunnersToCapture(t *testing.T) {
	if _, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY"); !ok {
		t.Skip("aws credentials not found in environment, skipping test")
	}

	cli.LogLevel = "DEBUG"
	agentCmdState.InstallTrustHostKey = true
	agentCmdState.CTFInfraTagKey = "CaptureTheFlagPlayer"
	cli.NonInteractive()

	runners, err := awsFindRunnersToCapture()
	assert.NoError(t, err)

	wg := new(sync.WaitGroup)
	for _, runner := range runners {
		wg.Add(1)
		go func(runner *lwrunner.AWSRunner) {
			out := fmt.Sprintf("--------- Runner ---------\nRegion: %v\nInstance ID: %v\n", runner.Region, runner.InstanceID)

			err = runner.SendAndUseIdentityFile()
			assert.NoError(t, err)

			err = verifyAccessToRemoteHost(&runner.Runner)
			assert.NoError(t, err)

			if alreadyInstalled := isAgentInstalledOnRemoteHost(&runner.Runner); alreadyInstalled != nil {
				out += alreadyInstalled.Error() + "\n"
			}

			fmt.Println(out)
			wg.Done()
		}(runner)
	}
	wg.Wait()
}
