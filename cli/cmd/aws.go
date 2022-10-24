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
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gammazero/workerpool"
	"github.com/lacework/go-sdk/lwrunner"
)

func awsDescribeInstances() ([]*lwrunner.AWSRunner, error) {
	regions, err := awsDescribeRegions()
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

// awsDescribeRegions queries the AWS API to list all the regions that
// are enabled for the user's AWS account. Use the "include_regions"
// command-line flag to only get regions in this list.
func awsDescribeRegions() ([]types.Region, error) {
	// Describe all regions that are enabled for the account
	var filters []types.Filter
	if len(agentCmdState.InstallIncludeRegions) > 0 {
		filters = []types.Filter{
			{
				Name:   aws.String("region-name"),
				Values: agentCmdState.InstallIncludeRegions,
			},
		}
	}
	input := &ec2.DescribeRegionsInput{
		Filters: filters,
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
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

	output, err := svc.DescribeRegions(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return output.Regions, nil
}

func awsRegionDescribeInstances(region string) ([]*lwrunner.AWSRunner, error) {
	var (
		tagKey = agentCmdState.InstallTagKey
		tag    = agentCmdState.InstallTag
	)
	cfg, err := config.LoadDefaultConfig(context.Background())
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
	result, err := svc.DescribeInstances(context.Background(), input)
	if err != nil {
		return nil, err
	}

	runners := []*lwrunner.AWSRunner{}
	producerWg := new(sync.WaitGroup)
	// cl := limiter.NewConcurrencyLimiter(agentCmdState.InstallMaxParallelism)
	wp := workerpool.New(agentCmdState.InstallMaxParallelism)
	runnerCh := make(chan *lwrunner.AWSRunner)

	// We have multiple producers of runners and a single consumer.
	// This goroutine acts as the consumer and reads from a channel into
	// a slice. Pass a pointer to this slice and wait for this goroutine
	// to finish before returning the memory pointed to.
	consumerWg := new(sync.WaitGroup)
	consumerWg.Add(1)
	go func(runners *[]*lwrunner.AWSRunner) {
		for runner := range runnerCh {
			*runners = append(*runners, runner)
		}
		consumerWg.Done()
	}(&runners)

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if instance.PublicIpAddress != nil && instance.State.Name == "running" {
				cli.Log.Debugw("found runner",
					"public ip address", *instance.PublicIpAddress,
					"instance state name", instance.State.Name,
				)

				if err != nil {
					cli.Log.Debugw("error identifying runner", "error", err, "instance_id", *instance.InstanceId)
					continue
				}

				producerWg.Add(1)

				// In order to use `cl.Execute()`, the input func() must not take any arguments.
				// Copy the runner info to dedicated variable in the goroutine
				instanceCopyWg := new(sync.WaitGroup)
				instanceCopyWg.Add(1)

				wp.Submit(func() {
					threadInstance := instance
					instanceCopyWg.Done()
					cli.Log.Debugw("found runner",
						"public ip address", *threadInstance.PublicIpAddress,
						"instance state name", threadInstance.State.Name,
					)

					runner, err := lwrunner.NewAWSRunner(
						*threadInstance.ImageId,
						*threadInstance.PublicIpAddress,
						region,
						*threadInstance.Placement.AvailabilityZone,
						*threadInstance.InstanceId,
						verifyHostCallback,
					)
					if err != nil {
						cli.Log.Debugw("error identifying runner", "error", err, "instance_id", *threadInstance.InstanceId)
					}

					runnerCh <- runner
					producerWg.Done()
				})
				instanceCopyWg.Wait()
			}
		}
	}

	// Wait for the producers to finish, then close the producer thread pool,
	// then close the channel they're writing to, then wait for the consumer to finish
	producerWg.Wait()
	wp.StopWait()
	close(runnerCh)
	consumerWg.Wait()

	return runners, nil
}
