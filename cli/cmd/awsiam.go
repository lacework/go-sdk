//
// Author:: Nicholas Schmeller (<nick.schmeller@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/google/uuid"
	"github.com/lacework/go-sdk/lwrunner"
)

// setupSSMAccess sets up an IAM role for SSM and attaches it to
// the machine's instance profile. Takes role name as argument;
// pass the empty string to create a new role.
// Then creates SSM document.
func setupSSMAccess(cfg aws.Config, roleName string, token string) (types.Role, types.InstanceProfile, error) {
	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	cli.Log.Debugw("setting up role", "passed roleName", roleName)
	role, err := setupSSMRole(c, roleName)
	if err != nil {
		return role, types.InstanceProfile{}, err
	}

	err = attachSSMPoliciesToRole(c, role)
	if err != nil {
		return role, types.InstanceProfile{}, err
	}

	// Create instance profile and add the role to it
	instanceProfile, err := setupInstanceProfile(c, role)
	if err != nil {
		return role, instanceProfile, err
	}

	return role, instanceProfile, nil
}

// teardownSSMAccess destroys all the infra created during the execution of this program.
// Specifically, this function:
// - Removes the role from the instance profile
// - Deletes the instance profile
// - Detaches all managed policies from the role (assumes no inline policies attached)
// - Deletes the role
func teardownSSMAccess(
	cfg aws.Config, role types.Role, instanceProfile types.InstanceProfile, byoRoleName string,
) error {
	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	if taggedLaceworkResource(instanceProfile.Tags) {
		cli.Log.Debugw("removing role from instance profile", "role", role, "instance profile", instanceProfile)
		_, err := c.RemoveRoleFromInstanceProfile(
			context.Background(),
			&iam.RemoveRoleFromInstanceProfileInput{
				InstanceProfileName: instanceProfile.InstanceProfileName,
				RoleName:            role.RoleName,
			},
		)
		if err != nil {
			return err
		}

		cli.Log.Debugw("deleting instance profile", "instance profile", instanceProfile)
		_, err = c.DeleteInstanceProfile(
			context.Background(),
			&iam.DeleteInstanceProfileInput{
				InstanceProfileName: instanceProfile.InstanceProfileName,
			},
		)
		if err != nil {
			return err
		}
	}

	if byoRoleName != "" || !taggedLaceworkResource(role.Tags) {
		cli.Log.Debugw("Lacework didn't create this role, will not delete it",
			"byoRoleName", byoRoleName,
			"role", role,
		)
		return nil
	}

	cli.Log.Debug("listing managed policies attached to this role (assuming no inline policies")
	listOutput, err := c.ListAttachedRolePolicies(
		context.Background(),
		&iam.ListAttachedRolePoliciesInput{
			RoleName: role.RoleName,
		},
	)
	if err != nil {
		return err
	}

	// Detach managed policies
	for _, attachedPolicy := range listOutput.AttachedPolicies {
		cli.Log.Debugw("detaching policy", "policy", attachedPolicy, "role", role)
		_, err := c.DetachRolePolicy(
			context.Background(),
			&iam.DetachRolePolicyInput{
				PolicyArn: attachedPolicy.PolicyArn,
				RoleName:  role.RoleName,
			},
		)
		if err != nil {
			return err
		}
	}

	cli.Log.Debugw("deleting role", "role", role)
	_, err = c.DeleteRole(
		context.Background(),
		&iam.DeleteRoleInput{
			RoleName: role.RoleName,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// taggedLaceworkResource is a helper function that takes the tag set of
// a suspected Lacework resource, iterates through the tags, and determines
// if the resource belongs to Lacework. Returns `true, nil` if the resource
// belongs to Lacework and `false, nil` if the resource does not belong to
// Lacework.
func taggedLaceworkResource(tags []types.Tag) bool {
	for _, tag := range tags {
		if *tag.Key == laceworkAutomationTagKey {
			return true
		}
	}

	return false
}

const laceworkAutomationTagKey = "LaceworkAutomation"

// setupSSMRole sets up an IAM role for SSM to assume.
// If `roleName` is not the empty string, then use that role
// instead of creating a new one.
func setupSSMRole(c *iam.Client, roleName string) (types.Role, error) {
	if roleName != "" {
		return getRoleFromName(c, roleName)
	} else {
		cli.Log.Debug("user did not provide a role, creating one now")
		return createSSMRole(c)
	}
}

type iamGetRoleFromNameAPI interface {
	GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)
}

func getRoleFromName(c iamGetRoleFromNameAPI, roleName string) (types.Role, error) {
	cli.Log.Debug("fetching info about role", roleName)
	output, err := c.GetRole(
		context.Background(),
		&iam.GetRoleInput{
			RoleName: aws.String(roleName),
		},
	)
	if err != nil {
		return types.Role{}, err
	}

	return *output.Role, nil
}

type iamCreateSSMRoleAPI interface {
	CreateRole(
		ctx context.Context, params *iam.CreateRoleInput, optFns ...func(*iam.Options),
	) (*iam.CreateRoleOutput, error)
}

// createSSMRole makes a call to the AWS API to create an IAM role.
// Returns information about the newly created role and any errors.
func createSSMRole(c iamCreateSSMRoleAPI) (types.Role, error) {
	const roleNameBase string = "Lacework-Agent-SSM-Install-Role-"
	roleName := roleNameBase + uuid.New().String()[:5]

	const trustPolicyDocument = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": { "Service": "ec2.amazonaws.com" },
			"Action": "sts:AssumeRole"
		}
	]
}`

	output, err := c.CreateRole(
		context.Background(),
		&iam.CreateRoleInput{
			AssumeRolePolicyDocument: aws.String(trustPolicyDocument),
			RoleName:                 aws.String(roleName),
			Description: aws.String(
				`Ephemeral role to install Lacework agents using SSM; created by the Lacework CLI.
Safe to delete if found`,
			),
			Tags: []types.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(roleName),
				},
				{
					Key:   aws.String(laceworkAutomationTagKey),
					Value: aws.String("agent-ssm-install"),
				},
			},
		},
	)
	if err != nil {
		return types.Role{}, err
	}
	cli.Log.Debugw("freshly created role", "role", *output.Role)

	return *output.Role, nil
}

// attachSSMPoliciesToRole takes a role, calls the IAM API to attach
// policies required for SSM to the role, and returns the role along
// with any errors.
func attachSSMPoliciesToRole(c *iam.Client, role types.Role) error {
	cli.Log.Debug("attaching policy to role")
	_, err := c.AttachRolePolicy(
		context.Background(),
		&iam.AttachRolePolicyInput{
			PolicyArn: aws.String(lwrunner.SSMInstancePolicy),
			RoleName:  role.RoleName,
		},
	)

	return err
}

func setupInstanceProfile(c *iam.Client, role types.Role) (types.InstanceProfile, error) {
	instanceProfile, err := createInstanceProfile(c)
	if err != nil {
		return instanceProfile, err
	}

	err = addRoleToInstanceProfile(c, role, instanceProfile)
	if err != nil {
		return types.InstanceProfile{}, err
	}

	return instanceProfile, nil
}

func createInstanceProfile(c *iam.Client) (types.InstanceProfile, error) {
	const instanceProfileNameBase string = "Lacework-Agent-SSM-Install-Instance-Profile-"
	instanceProfileName := instanceProfileNameBase + uuid.New().String()[:5]

	createOutput, err := c.CreateInstanceProfile(
		context.Background(),
		&iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(instanceProfileName),
			Tags: []types.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(instanceProfileName),
				},
				{
					Key:   aws.String(laceworkAutomationTagKey),
					Value: aws.String("agent-ssm-install"),
				},
			},
		},
	)
	if err != nil {
		return types.InstanceProfile{}, err
	}

	cli.Log.Debug("sleeping for 15sec to wait for instance profile eventual consistency")
	time.Sleep(15 * time.Second)

	return *createOutput.InstanceProfile, err
}

func addRoleToInstanceProfile(c *iam.Client, role types.Role, instanceProfile types.InstanceProfile) error {
	cli.Log.Debugw("adding role to instance profile", "role", role, "instance profile", instanceProfile)
	_, err := c.AddRoleToInstanceProfile(
		context.Background(),
		&iam.AddRoleToInstanceProfileInput{
			InstanceProfileName: instanceProfile.InstanceProfileName,
			RoleName:            role.RoleName,
		},
	)

	return err
}
