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
	ImageName        string
}

func NewAWSRunner(amiImageId, userFromCLIArg, host, region, availabilityZone, instanceID string, filterSSH bool, callback ssh.HostKeyCallback) (*AWSRunner, error) {
	// Look up the AMI name of the runner
	imageName, err := getAMIName(amiImageId, region)
	if err != nil {
		return nil, err
	}

	// Heuristically assign SSH username based on AMI name
	var detectedUsername string
	if filterSSH {
		detectedUsername, err = getSSHUsername(userFromCLIArg, imageName)
		if err != nil {
			return nil, err
		}
	} else {
		detectedUsername = "no_ssh_username_provided"
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
		imageName,
	}, nil
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
// Returns the association ID or an error.
func (run AWSRunner) AssociateInstanceProfileWithRunner(cfg aws.Config, instanceProfile types.InstanceProfile) (string, error) {
	c := ec2.New(ec2.Options{
		Credentials: cfg.Credentials,
		Region:      run.Region,
	})

	// Check to see if there are any instance profiles already associated with the runner
	describeOutput, err := c.DescribeIamInstanceProfileAssociations(
		context.Background(),
		&ec2.DescribeIamInstanceProfileAssociationsInput{
			Filters: []ec2types.Filter{
				{
					Name: aws.String("instance-id"),
					Values: []string{
						run.InstanceID,
					},
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	associationID, err := run.isCorrectInstanceProfileAlreadyAssociated(cfg, describeOutput.IamInstanceProfileAssociations)
	if err != nil {
		return "", err
	}

	if associationID != "" { // use the existing, correctly configured instance profile
		return associationID, nil
	} else { // associate our own instance profile
		associateOutput, err := c.AssociateIamInstanceProfile(
			context.Background(),
			&ec2.AssociateIamInstanceProfileInput{
				IamInstanceProfile: &ec2types.IamInstanceProfileSpecification{
					Arn: instanceProfile.Arn,
				},
				InstanceId: aws.String(run.InstanceID),
			},
		)
		if err != nil {
			return "", err
		}

		return *associateOutput.IamInstanceProfileAssociation.AssociationId, nil
	}
}

// isCorrectInstanceProfileAlreadyAssociated takes a list of instance profile associations
// and checks if there is an instance profile associated and if this instance
// profile has the correct policy for SSM access. Returns `<assoc. id>, nil` if so. Returns
// `"", nil` if there is no instance profile associated. Returns `"", <error>` if
// there is an incorrect instance profile associated, or if there was an error in
// executing this function.
func (run AWSRunner) isCorrectInstanceProfileAlreadyAssociated(cfg aws.Config, associations []ec2types.IamInstanceProfileAssociation) (string, error) {
	if len(associations) <= 0 { // no instance profile associated
		return "", nil
	}
	instanceProfileName := strings.Split(*associations[0].IamInstanceProfile.Arn, "instance-profile/")[1]

	c := iam.New(iam.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	getInstanceProfileOutput, err := c.GetInstanceProfile(
		context.Background(),
		&iam.GetInstanceProfileInput{
			InstanceProfileName: aws.String(instanceProfileName),
		},
	)
	if err != nil {
		return "", err
	}

	// Check to see if the instance profile associated with the runner has the correct policy

	if len(getInstanceProfileOutput.InstanceProfile.Roles) <= 0 { // can only have max one role
		return "", fmt.Errorf(
			"runner %v already has an instance profile (%v) attached, does not have a role",
			run,
			getInstanceProfileOutput.InstanceProfile,
		)
	}

	// Check which policies are associated with this instance profile's role
	listAttachedRolePoliciesOutput, err := c.ListAttachedRolePolicies(
		context.Background(),
		&iam.ListAttachedRolePoliciesInput{
			RoleName: getInstanceProfileOutput.InstanceProfile.Roles[0].RoleName,
		},
	)
	if err != nil {
		return "", err
	}

	for _, policy := range listAttachedRolePoliciesOutput.AttachedPolicies {
		if *policy.PolicyArn == SSMInstancePolicy {
			return *associations[0].AssociationId, nil // everything is configured correctly, we can return now
		}
	}

	// The runner has an instance profile attached, the instance profile has a role,
	// and the role does not have the policy we need for SSM. We can't install on
	// this instance, return an error
	return "", fmt.Errorf(
		"runner %v already has an instance profile (%v) attached, does not have policy %s",
		run,
		getInstanceProfileOutput.InstanceProfile,
		SSMInstancePolicy,
	)
}

func (run AWSRunner) DisassociateInstanceProfileFromRunner(cfg aws.Config, associationID string) error {
	c := ec2.New(ec2.Options{
		Credentials: cfg.Credentials,
		Region:      run.Region,
	})

	_, err := c.DisassociateIamInstanceProfile(
		context.Background(),
		&ec2.DisassociateIamInstanceProfileInput{
			AssociationId: aws.String(associationID),
		},
	)

	return err
}

const SSMInstancePolicy string = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"

// RunSSMCommandOnRemoteHost takes a shell command to install the agent on the runner
// the runner and executes it using SSM. `operation` must be one of the commands allowed
// by the SSM document. This function will not return until the command is in a terminal
// state, or until 2min have passed.
func (run AWSRunner) RunSSMCommandOnRemoteHost(cfg aws.Config, operation string) (ssm.GetCommandInvocationOutput, error) {
	c := ssm.New(ssm.Options{
		Credentials: cfg.Credentials,
		Region:      run.Region,
	})

	sendCommandOutput, err := c.SendCommand(
		context.Background(),
		&ssm.SendCommandInput{
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
		},
	)
	if err != nil {
		return ssm.GetCommandInvocationOutput{}, err
	}

	var getCommandInvocationOutput *ssm.GetCommandInvocationOutput

	// Sleep while waiting for the command to execute
	const durationTensOfSeconds = 12
	for i := 0; i < durationTensOfSeconds; i++ {
		time.Sleep(10 * time.Second)

		getCommandInvocationOutput, err = c.GetCommandInvocation(
			context.Background(),
			&ssm.GetCommandInvocationInput{
				CommandId:  sendCommandOutput.Command.CommandId,
				InstanceId: aws.String(run.InstanceID),
			},
		)
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

	return *getCommandInvocationOutput, fmt.Errorf("command %s did not finish in %dmin, final state %v, stdout %s, stderr %s",
		*sendCommandOutput.Command.CommandId,
		durationTensOfSeconds/6,
		*getCommandInvocationOutput,
		*getCommandInvocationOutput.StandardOutputContent,
		*getCommandInvocationOutput.StandardErrorContent,
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
			return strings.Contains(imageName, "amazon_linux"), "ec2-user"
		},
		func(imageName string) (bool, string) { return strings.Contains(imageName, "amzn2-ami"), "ec2-user" },
	}
}
