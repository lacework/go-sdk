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
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/lacework/go-sdk/lwrunner"
)

// setupSSMAccess sets up an IAM role for SSM and attaches it to
// the machine's instance profile. Takes role name as argument;
// pass the empty string to create a new role.
// Then creates SSM document.
func setupSSMAccess(cfg aws.Config, roleName string, token string) (types.Role, types.InstanceProfile, error) {
	cli.Log.Debugw("setting up role", "passed roleName", roleName)
	role, err := setupSSMRole(cfg, roleName)
	if err != nil {
		return role, types.InstanceProfile{}, err
	}

	err = attachSSMPoliciesToRole(cfg, role)
	if err != nil {
		return role, types.InstanceProfile{}, err
	}

	// Create instance profile and add the role to it
	instanceProfile, err := setupInstanceProfile(cfg, role)
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
func teardownSSMAccess(cfg aws.Config, role types.Role, instanceProfile types.InstanceProfile, byoRoleName string) error {
	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	// Only destroy the instance profile if it's ours
	if *instanceProfile.InstanceProfileName == instanceProfileName {
		// Remove role from instance profile
		cli.Log.Debugw("removing role from instance profile", "role", role, "instance profile", instanceProfile)
		removeInput := &iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: instanceProfile.InstanceProfileName,
			RoleName:            role.RoleName,
		}
		_, err := c.RemoveRoleFromInstanceProfile(context.Background(), removeInput)
		if err != nil {
			return err
		}

		// Delete instance profile
		cli.Log.Debugw("deleting instance profile", "instance profile", instanceProfile)
		deleteProfileInput := &iam.DeleteInstanceProfileInput{
			InstanceProfileName: instanceProfile.InstanceProfileName,
		}
		_, err = c.DeleteInstanceProfile(context.Background(), deleteProfileInput)
		if err != nil {
			return err
		}
	}

	if byoRoleName != "" || *role.RoleName != roleName {
		return nil // we didn't create this role, we should not delete it
	}

	// List managed policies attached to this role (assume there are no inline policies)
	listInput := &iam.ListAttachedRolePoliciesInput{
		RoleName: role.RoleName,
	}
	listOutput, err := c.ListAttachedRolePolicies(context.Background(), listInput)
	if err != nil {
		return err
	}

	// Detach managed policies
	for _, attachedPolicy := range listOutput.AttachedPolicies {
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
	deleteRoleInput := &iam.DeleteRoleInput{
		RoleName: role.RoleName,
	}
	_, err = c.DeleteRole(context.Background(), deleteRoleInput)
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
		return getRoleFromName(cfg, roleName)
	} else { // user did not provide a role, creating one now
		return createSSMRole(cfg)
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
// Returns information about the newly created role and any errors.
func createSSMRole(cfg aws.Config) (types.Role, error) {
	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	// Check if role already exists (in case CLI is run after interrupt or error)
	getInput := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}
	getOutput, err := c.GetRole(context.Background(), getInput)
	if err == nil && getOutput.Role != nil {
		return *getOutput.Role, err // we previously created the role, use it
	}

	trustPolicyDocument := `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": { "Service": "ec2.amazonaws.com" },
			"Action": "sts:AssumeRole"
		}
	]
}`
	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: &trustPolicyDocument,
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
				Key:   aws.String("LaceworkAutomation"),
				Value: aws.String("agent-ssm-install"),
			},
		},
	}
	output, err := c.CreateRole(context.Background(), input)
	if err != nil {
		return types.Role{}, err
	}

	return *output.Role, nil
}

const roleName string = "Lacework-Agent-SSM-Install-Role"

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
		PolicyArn: aws.String(lwrunner.SSMInstancePolicy),
		RoleName:  role.RoleName,
	}
	_, err := c.AttachRolePolicy(context.Background(), input)

	return err
}

func setupInstanceProfile(cfg aws.Config, role types.Role) (types.InstanceProfile, error) {
	instanceProfile, err := createInstanceProfile(cfg)
	if err != nil {
		return instanceProfile, err
	}

	err = addRoleToInstanceProfile(cfg, role, instanceProfile)
	if err != nil {
		return types.InstanceProfile{}, err
	}

	return instanceProfile, nil
}

func createInstanceProfile(cfg aws.Config) (types.InstanceProfile, error) {
	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	// Check if instance profile already exists (in case CLI is run after interrupt or error)
	getInput := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(instanceProfileName),
	}
	getOutput, err := c.GetInstanceProfile(context.Background(), getInput)
	if err == nil && getOutput.InstanceProfile != nil {
		return *getOutput.InstanceProfile, err // we previously created the profile, use it
	}

	cli.Log.Debug("creating instance profile")
	createInput := &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(instanceProfileName),
		Tags: []types.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(roleName),
			},
			{
				Key:   aws.String("LaceworkAutomation"),
				Value: aws.String("agent-ssm-install"),
			},
		},
	}
	createOutput, err := c.CreateInstanceProfile(context.Background(), createInput)
	if err != nil {
		return types.InstanceProfile{}, err
	}

	// Sleep for 15sec to wait for the instance profile to "settle in"
	time.Sleep(15 * time.Second)

	return *createOutput.InstanceProfile, err
}

const instanceProfileName string = "Lacework-Agent-SSM-Install-Instance-Profile"

func addRoleToInstanceProfile(cfg aws.Config, role types.Role, instanceProfile types.InstanceProfile) error {
	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	cli.Log.Debugw("checking if the role is already associated with the instance profile")
	if len(instanceProfile.Roles) > 0 {
		cli.Log.Debugw(
			"found a role already associated with the instance profile",
			"found role", instanceProfile.Roles[0],
			"our role", role,
		)
		if *instanceProfile.Roles[0].Arn == *role.Arn {
			return nil // the correct role is already attached to the instance profile
		} else { // someone else modified this instance profile. Fail now
			return fmt.Errorf(
				"tried to use role %s but pre-existing instance profile %s already has role %s",
				*role.Arn,
				instanceProfileName,
				*instanceProfile.Roles[0].Arn,
			)
		}
	}

	cli.Log.Debugw("adding role to instance profile", "role", role, "instance profile", instanceProfile)
	addInput := &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: instanceProfile.InstanceProfileName,
		RoleName:            role.RoleName,
	}
	_, err := c.AddRoleToInstanceProfile(context.Background(), addInput)

	return err
}
