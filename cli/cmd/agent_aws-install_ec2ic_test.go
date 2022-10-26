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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAwsInstallEC2IC(t *testing.T) {
	if _, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY"); !ok {
		t.Skip("aws credentials not found in environment, skipping test")
	}

	cli.LogLevel = "DEBUG"
	agentCmdState.InstallTrustHostKey = true
	agentCmdState.InstallTagKey = "CaptureTheFlagPlayer"
	cli.NonInteractive()

	err := installAWSEC2IC(nil, nil)
	assert.NoError(t, err)
}
