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

package lwrunner

import (
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/stretchr/testify/assert"
)

const GOOD_SSH_USERNAME = "customer_ssh_username"

func TestSSHUsernameLookupChecksForCLIArgUsernameFirst(t *testing.T) {
	user, err := getSSHUsername(GOOD_SSH_USERNAME, "some_ubuntu_ami")
	assert.NoError(t, err)
	assert.Equal(t, user, GOOD_SSH_USERNAME)
	assert.NotEqual(t, user, "ubuntu")
}

func TestSSHUsernameLookupChecksEnvBeforeAMI(t *testing.T) {
	t.Setenv("LW_SSH_USER", GOOD_SSH_USERNAME)

	user, err := getSSHUsername("", "some_ubuntu_ami")
	assert.NoError(t, err)
	assert.Equal(t, user, GOOD_SSH_USERNAME)
	assert.NotEqual(t, user, "ubuntu")
}

func TestSSHUsernameLookupFailsOnBadImageName(t *testing.T) {
	user, err := getSSHUsername("", "ami_bad_image_name")
	assert.Error(t, err)
	assert.Empty(t, user)
}

func TestSSHUsernameFromAmazonLinuxIsEC2User(t *testing.T) {
	user, err := getSSHUsername("", "amzn2-ami-hvm-x86_64-gp2")
	assert.NoError(t, err)
	assert.Equal(t, "ec2-user", user)
}

func TestGetSSMCommandInvocationStdOut(t *testing.T) {
	cases := []struct {
		mockCommandOutput *ssm.GetCommandInvocationOutput
		expectedResult    string
	}{
		{
			mockCommandOutput: &ssm.GetCommandInvocationOutput{
				StandardOutputContent: aws.String("foobarbaz"),
			},
			expectedResult: "foobarbaz",
		},
		{
			mockCommandOutput: &ssm.GetCommandInvocationOutput{
				StandardOutputContent: aws.String(""),
			},
			expectedResult: "",
		},
		{
			mockCommandOutput: &ssm.GetCommandInvocationOutput{
				StandardOutputContent: nil,
			},
			expectedResult: "",
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			stdoutContent := GetSSMCommandInvocationStdOut(*tt.mockCommandOutput)
			assert.Equal(t, tt.expectedResult, stdoutContent)
		})
	}
}

func TestGetSSMCommandInvocationStdErr(t *testing.T) {
	cases := []struct {
		mockCommandOutput *ssm.GetCommandInvocationOutput
		expectedResult    string
	}{
		{
			mockCommandOutput: &ssm.GetCommandInvocationOutput{
				StandardErrorContent: aws.String("foobarbaz"),
			},
			expectedResult: "foobarbaz",
		},
		{
			mockCommandOutput: &ssm.GetCommandInvocationOutput{
				StandardErrorContent: aws.String(""),
			},
			expectedResult: "",
		},
		{
			mockCommandOutput: &ssm.GetCommandInvocationOutput{
				StandardErrorContent: nil,
			},
			expectedResult: "",
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			stdoutContent := GetSSMCommandInvocationStdErr(*tt.mockCommandOutput)
			assert.Equal(t, tt.expectedResult, stdoutContent)
		})
	}
}
