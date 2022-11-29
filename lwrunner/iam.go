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
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetUserArn() (string, error) {
	c := sts.New(sts.Options{})
	output, err := c.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return *output.Arn, nil
}

type TrustedUser struct {
	ARN string
}

const trustPolicyDocumentTemplate = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": { "AWS": "%s" },
			"Action": "sts:AssumeRole",
		}
	]
}`

// CreateSSMRole makes a call to the AWS API to create an IAM role.
// Returns the newly created role and any errors.
func CreateSSMRole() (types.Role, error) {
	currentUserArn, err := GetUserArn()
	if err != nil {
		return types.Role{}, err
	}

	trustPolicyDocument := fmt.Sprintf(trustPolicyDocumentTemplate, currentUserArn)

	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: &trustPolicyDocument,
		RoleName:                 aws.String("Lacework-Agent-SSM-Install-Role"),
		Description: aws.String(
			`Ephemeral role to install Lacework agents using SSM; created by the Lacework CLI.
			Safe to delete if found`,
		),
	}

	c := iam.New(iam.Options{})
	output, err := c.CreateRole(context.Background(), input)
	if err != nil {
		return types.Role{}, err
	}

	return *output.Role, nil
}

// SetupSSMAccess sets up an IAM role for SSM and attaches it to
// the machine's instance profile. If `roleArn` is not the empty
// string, then use that role instead of creating a new one.
func (run AWSRunner) SetupSSMAccess(roleArn string) error {
	if roleArn != "" {
		// do something
		return nil
	}

	// User did not provide a role, creating one now
	role, err := CreateSSMRole()
	if err != nil {
		return err
	}
	run.SSMIAMRole = role

	return nil
}
