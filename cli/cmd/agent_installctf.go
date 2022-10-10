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
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/lacework/go-sdk/lwrunner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func installAWSSSH(_ *cobra.Command, args []string) error {
	runners, err := awsDescribeInstances()
	if err != nil {
		return err
	}

	for _, runner := range runners {
		cli.Log.Debugw("runner info: ", "user", runner.Runner.User, "region", runner.Region, "az", runner.AvailabilityZone, "instance ID", runner.InstanceID, "hostname", runner.Runner.Hostname)

		err := runner.Runner.UseIdentityFile(agentCmdState.InstallIdentityFile)
		if err != nil {
			return errors.Wrap(err, "unable to use provided identity file")
		}

		if err := verifyAccessToRemoteHost(&runner.Runner); err != nil {
			return errors.Wrap(err, "verifyAccessToRemoteHost failed")
		}

		if alreadyInstalled := isAgentInstalledOnRemoteHost(&runner.Runner); alreadyInstalled != nil {
			cli.Log.Debugw("agent already installed on host, skipping")
			continue
		}

		token := agentCmdState.InstallAgentToken
		if token == "" {
			// user didn't provide an agent token
			cli.Log.Debugw("agent token not provided")
			var err error
			token, err = selectAgentAccessToken()
			if err != nil {
				return err
			}
		}
		cmd := fmt.Sprintf("sudo sh -c \"curl -sSL %s | sh -s -- %s\"", agentInstallDownloadURL, token)
		err = runInstallCommandOnRemoteHost(&runner.Runner, cmd)
		if err != nil {
			cli.Log.Debugw("runInstallCommandOnRemoteHost failed")
			return err
		}
	}

	return nil
}

func installAWSEC2IC(_ *cobra.Command, args []string) error {
	runners, err := awsDescribeInstances()
	if err != nil {
		return err
	}

	for _, runner := range runners {
		cli.Log.Debugw("runner: ", "runner", runner)
		cli.Log.Debugw("runner user: ", "user", runner.Runner.User)
		cli.Log.Debugw("runner region: ", "region", runner.Region)
		cli.Log.Debugw("runner az: ", "az", runner.AvailabilityZone)
		cli.Log.Debugw("runner instance ID: ", "instance ID", runner.InstanceID)
		cli.Log.Debugw("runner: ", "runner hostname", runner.Runner.Hostname)
		err := runner.SendAndUseIdentityFile()
		if err != nil {
			cli.Log.Debugw("ec2ic failed", "err", err)
			continue
		}

		if err := verifyAccessToRemoteHost(&runner.Runner); err != nil {
			cli.Log.Debugw("verifyAccessToRemoteHost failed")
			return err
		}

		if alreadyInstalled := isAgentInstalledOnRemoteHost(&runner.Runner); alreadyInstalled != nil {
			cli.Log.Debugw("agent already installed on host, skipping")
			continue
		}

		token := agentCmdState.InstallAgentToken
		if token == "" {
			// user didn't provide an agent token
			cli.Log.Debugw("agent token not provided")
			var err error
			token, err = selectAgentAccessToken()
			if err != nil {
				return err
			}
		}
		cmd := fmt.Sprintf("sudo sh -c \"curl -sSL %s | sh -s -- %s\"", agentInstallDownloadURL, token)
		err = runInstallCommandOnRemoteHost(&runner.Runner, cmd)
		if err != nil {
			cli.Log.Debugw("runInstallCommandOnRemoteHost failed")
			return err
		}
	}

	return nil
}

func awsDescribeInstances() ([]*lwrunner.AWSRunner, error) {
	regions, err := awsFindRegions()
	if err != nil {
		return nil, err
	}

	allRunners := []*lwrunner.AWSRunner{}
	for _, region := range regions {
		regionRunners, err := awsRegionDescribeInstances(*region.RegionName)
		if err != nil {
			return nil, err
		}
		allRunners = append(allRunners, regionRunners...)
	}

	return allRunners, nil
}

// awsFindRegions queries the AWS API to list all the regions that
// are enabled for the user's AWS account. Use the "include_regions"
// command-line flag to only get regions in this list.
func awsFindRegions() ([]types.Region, error) {
	// Describe all regions that are enabled for the account
	var filters []types.Filter
	if len(agentCmdState.CTFIncludeRegions) > 0 {
		filters = []types.Filter{
			{
				Name:   aws.String("region-name"),
				Values: agentCmdState.CTFIncludeRegions,
			},
		}
	}
	input := &ec2.DescribeRegionsInput{
		Filters: filters,
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// Look for region string in shell environment first
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		region = "us-west-2" // use us-west-2 for lack of a better region
	}
	svc := ec2.New(ec2.Options{
		Credentials: cfg.Credentials,
		Region:      region,
	})

	output, err := svc.DescribeRegions(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return output.Regions, nil
}

func awsRegionDescribeInstances(region string) ([]*lwrunner.AWSRunner, error) {
	var (
		tagKey = agentCmdState.CTFInfraTagKey
		tag    = agentCmdState.CTFInfraTag
	)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	svc := ec2.New(ec2.Options{
		Credentials: cfg.Credentials,
		Region:      region,
	})

	var filters []types.Filter

	// Filter for instances that are running
	filters = append(filters, types.Filter{
		Name: aws.String("instance-state-name"),
		Values: []string{
			"running",
		},
	})

	// Filter for instances where a tag key exists
	if tagKey != "" {
		cli.Log.Debugw("found tagKey", "tagKey", tagKey)
		filters = append(filters, types.Filter{
			Name: aws.String("tag-key"),
			Values: []string{
				tagKey,
			},
		})
	}

	// Filter for instances where certain tags exist
	if len(tag) > 0 {
		cli.Log.Debugw("found tags", "tag length", len(tag), "tags", tag)
		filters = append(filters, types.Filter{
			Name:   aws.String("tag:" + tag[0]),
			Values: tag[1:],
		})
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}
	result, err := svc.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	runners := []*lwrunner.AWSRunner{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if instance.PublicIpAddress != nil && instance.State.Name == "running" {
				cli.Log.Debugw("found runner", "public ip address", *instance.PublicIpAddress, "instance state name", instance.State.Name)

				runner, err := lwrunner.NewAWSRunner(*instance.ImageId, *instance.PublicIpAddress, region, *instance.Placement.AvailabilityZone, *instance.InstanceId, verifyHostCallback)
				if err != nil {
					cli.Log.Debugw("error identifying runner", "error", err, "instance ID", *instance.InstanceId)
					continue
				}

				runners = append(runners, runner)
			}
		}
	}

	return runners, nil
}
