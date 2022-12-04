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
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// SetupSSMRole sets up an IAM role for SSM and attaches it to
// the machine's instance profile. Takes role name as argument;
// pass the empty string to create a new role.
// Then creates SSM document.
func SetupSSMAccess(cfg aws.Config, roleName string, token string) (types.Role, types.InstanceProfile, error) {
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

	err = createSSMDocument(cfg, token)
	if err != nil {
		return role, instanceProfile, err
	}

	return role, instanceProfile, nil
}

// TeardownSSMAccess destroys all the infra created during the execution of this program.
// Specifically, this function:
// - Removes the role from the instance profile
// - Deletes the instance profile
// - Detaches all managed policies from the role
//   - This assumes there are no inline policies attached to the role
// - Deletes the role
func TeardownSSMAccess(cfg aws.Config, role types.Role, instanceProfile types.InstanceProfile, byoRoleName string) error {
	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

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

	if byoRoleName != "" { // user brought their own role, we should not delete it
		return nil
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

	// Delete the SSM document
	err = deleteSSMDocument(cfg, SSMDocumentName)
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
	const roleName string = "Lacework-Agent-SSM-Install-Role"

	// Check if role already exists (in case CLI is run after interrupt or error)
	getInput := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}
	getOutput, err := c.GetRole(context.Background(), getInput)
	if err == nil && getOutput.Role != nil {
		return *getOutput.Role, err // we previously created the role, use it
	}

	currentUserARN, err := getCurrentUserARN(cfg)
	if err != nil {
		return types.Role{}, err
	}

	trustPolicyDocument := fmt.Sprintf(trustPolicyDocumentTemplate, currentUserARN)

	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: &trustPolicyDocument,
		RoleName:                 aws.String(roleName),
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

func createSSMDocument(cfg aws.Config, token string) error {
	c := ssm.New(ssm.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})
	agentVersionCmd := "sudo sh -c \"/var/lib/lacework/datacollector -v\""
	runInstallCmd := fmt.Sprintf("sudo sh -c \"curl -sSL %s | sh -s -- %s\"", agentInstallDownloadURL, token)

	ssmDocumentContents := fmt.Sprintf(ssmDocumentTemplate, agentVersionCmd, runInstallCmd)
	cli.Log.Debugw("ssmDocumentContents", "contents", ssmDocumentContents)

	input := &ssm.CreateDocumentInput{
		Content:        aws.String(ssmDocumentContents),
		Name:           aws.String(SSMDocumentName),
		DocumentFormat: ssmtypes.DocumentFormatYaml,
		DocumentType:   ssmtypes.DocumentTypeCommand,
		TargetType:     aws.String("/AWS::EC2::Instance"),
	}
	output, err := c.CreateDocument(context.Background(), input)
	if err != nil {
		return err
	}

	cli.Log.Debugw(
		"created SSM document",
		"document", *output.DocumentDescription,
		"contents", ssmDocumentContents,
	)
	return nil
}

func deleteSSMDocument(cfg aws.Config, documentName string) error {
	c := ssm.New(ssm.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	cli.Log.Debugw("deleting SSM document", "document name", SSMDocumentName)
	input := &ssm.DeleteDocumentInput{
		Name:  aws.String(SSMDocumentName),
		Force: true,
	}
	_, err := c.DeleteDocument(context.Background(), input)

	return err
}

const ssmDocumentTemplate = `---
schemaVersion: '2.2'
description: runShellScript with command strings stored as Parameter Store parameter
parameters:
  commands:
    type: StringList
    description: (Required) The commands to run on the instance.
    allowedValues:
    - %s
    - %s
mainSteps:
- action: aws:runShellScript
  name: runShellScriptDefaultParams
  inputs:
    runCommand:
    - "{{ commands }}"`

const SSMDocumentName = "Lacework-Agent-SSM-Install-Document"
