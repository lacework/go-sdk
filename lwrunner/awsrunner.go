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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
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

	// Logic for when there is already an instance profile associated with the runner
	// If the CLI fails or is interrupted, this may be our instance profile
	if len(describeOutput.IamInstanceProfileAssociations) > 0 {
		foundArn := *describeOutput.IamInstanceProfileAssociations[0].IamInstanceProfile.Arn
		if foundArn == *instanceProfile.Arn { // found our instance profile associated, use it
			// should already have the role attached if it was associated
			return nil
		} else {
			return fmt.Errorf(
				"runner %v already has an instance profile (%v) attached",
				run,
				describeOutput.IamInstanceProfileAssociations[0],
			)
		}
	}

	// Associate our own instance profile
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
	fmt.Println("successfully associated with instance ID", run.InstanceID)

	return nil
}

func (run AWSRunner) IsAgentInstalledOnRemoteHostSSM(cfg aws.Config, documentName string, agentVersionCmd string) error {
	c := ssm.New(ssm.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})

	input := &ssm.SendCommandInput{
		DocumentName: aws.String(documentName),
		Comment:      aws.String("checks if Lacework Agent is installed on host"),
		InstanceIds: []string{
			run.InstanceID,
		},
		Parameters: map[string][]string{
			"commands": {
				agentVersionCmd,
			},
		},
	}

	_, err := c.SendCommand(context.Background(), input)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
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
