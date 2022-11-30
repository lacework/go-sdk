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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// SetupSSMRole sets up an IAM role for SSM and attaches it to
// the machine's instance profile. Takes role name as argument;
// pass the empty string to create a new role.
func SetupSSMAccess(roleName string) (types.Role, error) {
	role, err := setupSSMRole(roleName)
	if err != nil {
		return role, err
	}

	err = attachSSMPoliciesToRole(role)
	if err != nil {
		return role, err
	}

	// TODO do we need to give the user permission to switch to the role?

	return role, nil
}

func TeardownSSMAccess(role types.Role) error {
	c := iam.New(iam.Options{})

	// List managed policies attached to this role (assume there are no inline policies)
	listInput := &iam.ListAttachedRolePoliciesInput{
		RoleName: role.RoleName,
	}
	output, err := c.ListAttachedRolePolicies(context.Background(), listInput)
	if err != nil {
		return err
	}

	// Detach managed policies
	for _, attachedPolicy := range output.AttachedPolicies {
		detachInput := &iam.DetachRolePolicyInput{
			PolicyArn: attachedPolicy.PolicyArn,
			RoleName:  role.RoleName,
		}
		_, err := c.DetachRolePolicy(context.Background(), detachInput)
		if err != nil {
			return err
		}
	}

	// Delete the role
	deleteInput := &iam.DeleteRoleInput{
		RoleName: role.RoleName,
	}
	_, err = c.DeleteRole(context.Background(), deleteInput)
	if err != nil {
		return err
	}

	return nil
}

// setupSSMRole sets up an IAM role for SSM to assume.
// If `roleName` is not the empty string, then use that role
// instead of creating a new one.
func setupSSMRole(roleName string) (types.Role, error) {
	if roleName != "" {
		role, err := getRoleFromName(roleName)
		return role, err
	} else { // user did not provide a role, creating one now
		role, err := createSSMRole()
		return role, err
	}
}

func getRoleFromName(roleName string) (types.Role, error) {
	c := iam.New(iam.Options{})

	input := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}
	output, err := c.GetRole(context.Background(), input)
	if err != nil {
		return types.Role{}, err
	}

	return *output.Role, nil
}

// createSSMRole makes a call to the AWS API to create an IAM role.
// This role allows the current user to assume it.
// Returns information about the newly created role and any errors.
func createSSMRole() (types.Role, error) {
	c := iam.New(iam.Options{})

	currentUserARN, err := getCurrentUserARN()
	if err != nil {
		return types.Role{}, err
	}

	trustPolicyDocument := fmt.Sprintf(trustPolicyDocumentTemplate, currentUserARN)

	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: &trustPolicyDocument,
		RoleName:                 aws.String("Lacework-Agent-SSM-Install-Role"),
		Description: aws.String(
			`Ephemeral role to install Lacework agents using SSM; created by the Lacework CLI.
			Safe to delete if found`,
		),
	}
	output, err := c.CreateRole(context.Background(), input)
	if err != nil {
		return types.Role{}, err
	}

	return *output.Role, nil
}

func getCurrentUserARN() (string, error) {
	c := sts.New(sts.Options{})

	output, err := c.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return *output.Arn, nil
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

// attachSSMPoliciesToRole takes a role, calls the IAM API to attach
// policies required for SSM to the role, and returns the role along
// with any errors.
func attachSSMPoliciesToRole(role types.Role) error {
	c := iam.New(iam.Options{})

	input := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonSSMFullAccess"),
		RoleName:  aws.String("AmazonSSMFullAccess"),
	}
	_, err := c.AttachRolePolicy(context.Background(), input)

	return err
}
