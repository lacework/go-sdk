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
func SetupSSMAccess(cfg aws.Config, roleName string) (types.Role, error) {
	cli.Log.Debugw("setting up role", "passed roleName", roleName)
	role, err := setupSSMRole(cfg, roleName)
	if err != nil {
		return role, err
	}

	err = attachSSMPoliciesToRole(cfg, role)
	if err != nil {
		return role, err
	}

	// TODO do we need to give the user permission to switch to the role?

	return role, nil
}

func TeardownSSMAccess(cfg aws.Config, role types.Role) error {
	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

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
		cli.Log.Debugw("detaching policy", "policy", attachedPolicy, "role", role)
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
	cli.Log.Debugw("deleting role", "role", role)
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
func setupSSMRole(cfg aws.Config, roleName string) (types.Role, error) {
	if roleName != "" {
		role, err := getRoleFromName(cfg, roleName)
		return role, err
	} else { // user did not provide a role, creating one now
		role, err := createSSMRole(cfg)
		return role, err
	}
}

func getRoleFromName(cfg aws.Config, roleName string) (types.Role, error) {
	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

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
func createSSMRole(cfg aws.Config) (types.Role, error) {
	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	currentUserARN, err := getCurrentUserARN(cfg)
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

func getCurrentUserARN(cfg aws.Config) (string, error) {
	c := sts.New(sts.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

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
			"Action": "sts:AssumeRole"
		}
	]
}`

// attachSSMPoliciesToRole takes a role, calls the IAM API to attach
// policies required for SSM to the role, and returns the role along
// with any errors.
func attachSSMPoliciesToRole(cfg aws.Config, role types.Role) error {
	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	cli.Log.Debug("attaching policy to role")
	input := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonSSMFullAccess"),
		RoleName:  role.RoleName,
	}
	_, err := c.AttachRolePolicy(context.Background(), input)

	return err
}
