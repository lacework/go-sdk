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
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"golang.org/x/crypto/ssh"
)

type AWSRunner struct {
	Runner           Runner
	Region           string
	AvailabilityZone string
	InstanceID       string
}

func NewAWSRunner(amiImageId, userFromCLIArg, host, region, availabilityZone, instanceID string, callback ssh.HostKeyCallback) (*AWSRunner, error) {
	// Look up the AMI name of the runner
	imageName, err := getAMIName(amiImageId, region)
	if err != nil {
		return nil, err
	}

	// Heuristically assign SSH username based on AMI name
	detectedUsername, err := getSSHUsername(userFromCLIArg, imageName)
	if err != nil {
		return nil, err
	}

	defaultCallback, err := DefaultKnownHosts()
	if err == nil && callback == nil {
		callback = defaultCallback
	}

	runner := New(detectedUsername, host, callback)

	return &AWSRunner{
		*runner,
		region,
		availabilityZone,
		instanceID,
	}, nil
}

func (run AWSRunner) RunSession() error {
	c := ssm.New(ssm.Options{
		Region: run.Region,
	})

	input := &ssm.StartSessionInput{
		Target:       aws.String(run.InstanceID),
		DocumentName: aws.String("SSM-SessionManagerRunShell"), // default, but some IAM roles require us to specify
		Reason:       aws.String("Lacework CLI agent install"),
	}

	_, err := c.StartSession(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}

func (run AWSRunner) SendAndUseIdentityFile() error {
	pubBytes, privBytes, err := GetKeyBytes()
	if err != nil {
		return err
	}

	err = run.SendPublicKey(pubBytes)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(privBytes)
	if err != nil {
		return err
	}
	run.Runner.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}

	return nil
}

// Helper function to send a public key to a test instance. Uses
// EC2InstanceConnect. The AWS account used to run the tests must
// have EC2InstanceConnect permissions attached to its IAM role.
// First checks to make sure the instance is still running.
func (run AWSRunner) SendPublicKey(pubBytes []byte) error {
	// Send public key
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}
	cfg.Region = run.Region
	svc := ec2instanceconnect.NewFromConfig(cfg)

	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: &run.AvailabilityZone,
		InstanceId:       &run.InstanceID,
		InstanceOSUser:   aws.String(run.Runner.User),
		SSHPublicKey:     aws.String(string(pubBytes)),
	}

	_, err = svc.SendSSHPublicKey(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}

// AssociateInstanceProfileWithRunner associates a given instance profile with the
// receiving runner. First checks if there are any instance profiles already associated
// with the runner, and returns an error if so (since a runner can only have one instance
// profile associated with it). Then associates the instance profile with the runner.
func (run AWSRunner) AssociateInstanceProfileWithRunner(cfg aws.Config, instanceProfile types.InstanceProfile) error {
	c := ec2.New(ec2.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	// Check to see if there are any instance profiles already associated with the runner
	describeInput := &ec2.DescribeIamInstanceProfileAssociationsInput{
		Filters: []ec2types.Filter{
			{
				Name: aws.String("instance-id"),
				Values: []string{
					run.InstanceID,
				},
			},
		},
	}
	describeOutput, err := c.DescribeIamInstanceProfileAssociations(context.Background(), describeInput)
	if err != nil {
		return err
	}

	alreadyAssociated, err := run.isCorrectInstanceProfileAlreadyAssociated(cfg, describeOutput.IamInstanceProfileAssociations)
	if err != nil {
		return err
	}

	if alreadyAssociated { // use the existing, correctly configured instance profile
		return nil
	} else { // associate our own instance profile
		associateInput := &ec2.AssociateIamInstanceProfileInput{
			IamInstanceProfile: &ec2types.IamInstanceProfileSpecification{
				Arn: instanceProfile.Arn,
			},
			InstanceId: aws.String(run.InstanceID),
		}
		_, err = c.AssociateIamInstanceProfile(context.Background(), associateInput)
		if err != nil {
			return err
		}

		return nil
	}
}

// isCorrectInstanceProfileAlreadyAssociated takes a list of instance profile associations
// and checks if there is an instance profile associated and if this instance
// profile has the correct policy for SSM access. Returns `true, nil` if so. Returns
// `false, nil` if there is no instance profile associated. Returns `false, <error>` if
// there is an incorrect instance profile associated, or if there was an error in
// executing this function.
func (run AWSRunner) isCorrectInstanceProfileAlreadyAssociated(cfg aws.Config, associations []ec2types.IamInstanceProfileAssociation) (bool, error) {
	if len(associations) <= 0 { // no instance profile associated
		return false, nil
	}
	instanceProfileName := strings.Split(*associations[0].IamInstanceProfile.Arn, "instance-profile/")[1]

	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	getInstanceProfileInput := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(instanceProfileName),
	}
	getInstanceProfileOutput, err := c.GetInstanceProfile(context.Background(), getInstanceProfileInput)
	if err != nil {
		return false, err
	}

	// Check to see if the instance profile associated with the runner has the correct policy

	// if foundArn == *instanceProfile.Arn { // found our instance profile associated, use it
	// 	// should already have the role attached if it was associated
	// 	return nil

	if len(getInstanceProfileOutput.InstanceProfile.Roles) <= 0 { // can only have max one role
		return false, fmt.Errorf(
			"runner %v already has an instance profile (%v) attached, does not have a role",
			run,
			getInstanceProfileOutput.InstanceProfile,
		)
	}

	// Check which policies are associated with this instance profile's role
	listAttachedRolePoliciesInput := iam.ListAttachedRolePoliciesInput{
		RoleName: getInstanceProfileOutput.InstanceProfile.Roles[0].RoleName,
	}
	listAttachedRolePoliciesOutput, err := c.ListAttachedRolePolicies(context.Background(), &listAttachedRolePoliciesInput)
	if err != nil {
		return false, err
	}

	ssmPolicy := "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
	for _, policy := range listAttachedRolePoliciesOutput.AttachedPolicies {
		if *policy.PolicyArn == ssmPolicy {
			return true, nil // everything is configured correctly, we can return now
		}
	}

	// The runner has an instance profile attached, the instance profile has a role,
	// and the role does not have the policy we need for SSM. We can't install on
	// this instance, return an error
	return false, fmt.Errorf(
		"runner %v already has an instance profile (%v) attached, does not have policy %s",
		run,
		getInstanceProfileOutput.InstanceProfile,
		ssmPolicy,
	)
}

// RunSSMCommandOnRemoteHost takes a shell command to install the agent on
// the runner and executes it using SSM. The user must pass in a valid
// `documentName`. `operation` must be one of the commands allowed by the SSM
// document. This function will not return until the command is in a terminal
// state.
func (run AWSRunner) RunSSMCommandOnRemoteHost(cfg aws.Config, operation string) (ssm.GetCommandInvocationOutput, error) {
	c := ssm.New(ssm.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	input := &ssm.SendCommandInput{
		DocumentName: aws.String("AWS-RunShellScript"),
		Comment:      aws.String("this command is for installing the Lacework Agent"),
		InstanceIds: []string{
			run.InstanceID,
		},
		Parameters: map[string][]string{
			"commands": {
				operation,
			},
		},
	}

	sendCommandOutput, err := c.SendCommand(context.Background(), input)
	if err != nil {
		// return ssmtypes.CommandInvocation{}, err
		return ssm.GetCommandInvocationOutput{}, err
	}

	var getCommandInvocationOutput *ssm.GetCommandInvocationOutput
	getCommandInvocationInput := &ssm.GetCommandInvocationInput{
		CommandId:  sendCommandOutput.Command.CommandId,
		InstanceId: aws.String(run.InstanceID),
	}

	// Wait for up to a minute for the command to execute
	for i := 0; i < 6; i++ {
		time.Sleep(10 * time.Second)

		getCommandInvocationOutput, err = c.GetCommandInvocation(context.Background(), getCommandInvocationInput)
		if err != nil {
			return ssm.GetCommandInvocationOutput{}, err
		}

		// Check if the command has reached a "terminal state"
		if getCommandInvocationOutput.Status == ssmtypes.CommandInvocationStatusSuccess ||
			getCommandInvocationOutput.Status == ssmtypes.CommandInvocationStatusCancelled ||
			getCommandInvocationOutput.Status == ssmtypes.CommandInvocationStatusTimedOut ||
			getCommandInvocationOutput.Status == ssmtypes.CommandInvocationStatusFailed {
			return *getCommandInvocationOutput, nil
		}
	}

	return *getCommandInvocationOutput, fmt.Errorf("command %s did not finish in 1min, final state %v",
		*sendCommandOutput.Command.CommandId,
		*getCommandInvocationOutput,
	)
}

// getAMIName takes an AMI image ID and an AWS region name as input
// and calls the AWS API to get the name of the AMI. Returns the AMI
// name or an error if unsuccessful.
func getAMIName(amiImageId, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}
	cfg.Region = region
	svc := ec2.NewFromConfig(cfg)
	input := ec2.DescribeImagesInput{
		ImageIds: []string{
			amiImageId,
		},
	}
	result, err := svc.DescribeImages(context.Background(), &input)
	if err != nil {
		return "", err
	}
	if len(result.Images) != 1 {
		return "", fmt.Errorf("expected to find only one AMI, instead found %v", result.Images)
	}

	return *result.Images[0].Name, nil
}

// getSSHUsername takes any username passed as a CLI arg,
// an AMI image name, a shell environment, and returns
// the username for SSHing into the AWS runner or the empty
// string and an error if the AMI is not supported.
// It first checks if `LW_SSH_USER` is set and returns it if so.
// Then it checks the AMI image name to heuristically determine the
// SSH username.
func getSSHUsername(userFromCLIArg, imageName string) (string, error) {
	if userFromCLIArg != "" { // from CLI arg
		return userFromCLIArg, nil
	}
	usernameLUT := getSSHUsernameLookupTable()
	for _, matchFn := range usernameLUT {
		if match, foundName := matchFn(imageName); match {
			return foundName, nil
		}
	}
	// No matching AMI found, return an error
	return "", fmt.Errorf("no SSH username found for AMI %s, set as arg or shell env", imageName)
}

// getSSHUsernameLookupTable returns a lookup table for heuristically
// determining SSH username based on AMI.
// The first row of the table it returns is a function that checks
// `LW_SSH_USER` in the shell environment.
func getSSHUsernameLookupTable() []func(string) (bool, string) {
	return []func(string) (bool, string){
		func(_ string) (bool, string) { return os.Getenv("LW_SSH_USER") != "", os.Getenv("LW_SSH_USER") }, // THIS ROW MUST BE FIRST IN THE TABLE
		func(imageName string) (bool, string) { return strings.Contains(imageName, "ubuntu"), "ubuntu" },
		func(imageName string) (bool, string) {
			return strings.Contains(imageName, "amazon_linux"), "amazon_linux"
		},
		func(imageName string) (bool, string) { return strings.Contains(imageName, "amzn2-ami"), "amzn2-ami" },
	}
}
