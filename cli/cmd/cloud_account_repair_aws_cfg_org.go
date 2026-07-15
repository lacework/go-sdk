//
// Author:: Fortinet
// Copyright:: Copyright 2026, Fortinet
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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/v2/api"
)

// repairMethodAwsCfgOrg repairs an account onboarded via the AWS Config organization template.
const repairMethodAwsCfgOrg = "aws-cfg-org"

// aws-cfg-org flags.
var (
	repairAccountID  string // --account-id
	repairStackName  string // --stack-name (the AWS Config org management stack)
	repairAwsProfile string // --aws-profile
	repairAwsRegion  string // --aws-region
	repairDryRun     bool   // --dry-run
)

func init() {
	f := cloudAccountRepairCmd.Flags()
	f.StringVar(&repairAccountID,
		"account-id", "", "AWS account id whose AwsCfg integration is missing (required)")
	f.StringVar(&repairStackName,
		"stack-name", "", "name of the Lacework AWS Config org management CloudFormation stack (required)")
	f.StringVar(&repairAwsProfile,
		"aws-profile", "", "AWS profile for the org management account credentials")
	f.StringVar(&repairAwsRegion,
		"aws-region", "", "AWS region the management stack/stackset live in (default us-east-1)")
	f.BoolVar(&repairDryRun,
		"dry-run", false, "show what would be re-registered without calling the Lacework API")
}

// repairAwsCfgOrg re-registers a missing AwsCfg integration for one AWS account onboarded through the
// AWS Config organization CloudFormation template. It runs with AWS Organizations management-account
// credentials and reads the member StackSet instance to derive the role ARN and external id.
func repairAwsCfgOrg() error {
	if repairAccountID == "" {
		return errors.New("--account-id is required")
	}
	if repairStackName == "" {
		return errors.New("--stack-name is required")
	}

	ctx := context.Background()
	region := repairAwsRegion
	if region == "" {
		region = "us-east-1"
	}
	opts := []func(*config.LoadOptions) error{config.WithRegion(region)}
	if repairAwsProfile != "" {
		opts = append(opts, config.WithSharedConfigProfile(repairAwsProfile))
	}
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return errors.Wrap(err, "unable to load AWS credentials")
	}
	cfn := cloudformation.NewFromConfig(cfg)

	// 1. LaceworkAccount (tenant) from the management stack parameters.
	cli.StartProgress(" Reading the management stack...")
	laceworkAccount, err := stackParameter(ctx, cfn, repairStackName, "LaceworkAccount")
	if err != nil {
		cli.StopProgress()
		return err
	}

	// 2. The StackSet name is the physical id of the management stack's LaceworkStackset resource.
	stackSetName, err := stackResourcePhysicalID(ctx, cfn, repairStackName, "LaceworkStackset")
	cli.StopProgress()
	if err != nil {
		return err
	}

	// 3. The member account's stack instance carries the StackId we derive the role UID from.
	cli.StartProgress(fmt.Sprintf(" Locating the stack instance for %s...", repairAccountID))
	stackID, found, err := memberStackID(ctx, cfn, stackSetName, repairAccountID)
	cli.StopProgress()
	if err != nil {
		return err
	}

	// No stack instance means the IAM role does not exist; re-create the instance (which rebuilds the
	// role) before registering. This only lands the account if it is under an OU the StackSet targets.
	if !found {
		targetOUs := targetOUsParam(ctx, cfn)
		if repairDryRun {
			if cli.JSONOutput() {
				return cli.OutputJSON(map[string]interface{}{
					"accountId": repairAccountID,
					"action":    "create-stack-instance-and-register",
					"stackSet":  stackSetName,
					"targetOUs": targetOUs,
					"dryRun":    true,
				})
			}
			cli.OutputHuman("Account %s has no stack instance in stackset %s.\n", repairAccountID, stackSetName)
			cli.OutputHuman("Dry-run: would create-stack-instances (targeting OUs %s) to rebuild the "+
				"role, then register the integration. Re-run without --dry-run.\n", strings.Join(targetOUs, ","))
			return nil
		}
		if len(targetOUs) == 0 {
			return errors.Errorf(
				"stack %s has no OrganizationalUnits parameter; cannot create a stack instance", repairStackName)
		}
		cli.StartProgress(fmt.Sprintf(" Creating the stack instance for %s (rebuilding the IAM role)...", repairAccountID))
		err = createMemberInstance(ctx, cfn, stackSetName, cfg.Region, targetOUs, repairAccountID)
		cli.StopProgress()
		if err != nil {
			return err
		}
		stackID, found, err = memberStackID(ctx, cfn, stackSetName, repairAccountID)
		if err != nil {
			return err
		}
		if !found {
			return errors.Errorf(
				"no stack instance was created for %s: the account is not under any OU the stackset "+
					"targets (%s). Add the account to a targeted OU (or add its OU to the stack's "+
					"OrganizationalUnits parameter), then retry", repairAccountID, strings.Join(targetOUs, ","))
		}
		cli.OutputHuman("Stack instance created for %s.\n", repairAccountID)
	}

	uid := stackUID(stackID)
	if uid == "" {
		return errors.Errorf("could not derive the role UID from stack id %q", stackID)
	}

	// Role name and external id are both keyed off the same UID by the config_org member template
	// (role: lacework-config-role-<UID>; externalId: lweid:aws:v2:<tenant>:<acct>:LW<UID>).
	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/lacework-config-role-%s", repairAccountID, uid)
	externalID := fmt.Sprintf("lweid:aws:v2:%s:%s:LW%s", laceworkAccount, repairAccountID, uid)
	name := fmt.Sprintf("%s-Config", laceworkAccount)

	cli.OutputHuman("Re-registering AwsCfg integration:\n")
	cli.OutputHuman("  account-id:  %s\n", repairAccountID)
	cli.OutputHuman("  name:        %s\n", name)
	cli.OutputHuman("  role-arn:    %s\n", roleArn)
	cli.OutputHuman("  external-id: %s\n", externalID)

	if repairDryRun {
		if cli.JSONOutput() {
			return cli.OutputJSON(repairRegisterResult("register", repairAccountID, name, roleArn, externalID, ""))
		}
		cli.OutputHuman("\nDry-run: nothing was created. Re-run without --dry-run to register.\n")
		return nil
	}

	awsCfg := api.NewCloudAccount(name, api.AwsCfgCloudAccount, api.AwsCfgData{
		AwsAccountID: repairAccountID,
		Credentials: api.AwsCfgCredentials{
			RoleArn:    roleArn,
			ExternalID: externalID,
		},
	})

	cli.StartProgress(" Creating integration...")
	resp, err := cli.LwApi.V2.CloudAccounts.Create(awsCfg)
	cli.StopProgress()
	if err != nil {
		// Mirror the setup Lambda: an already-registered account is success, not an error.
		if strings.Contains(strings.ToLower(err.Error()), "aws account is already used") {
			if cli.JSONOutput() {
				return cli.OutputJSON(map[string]interface{}{
					"accountId": repairAccountID,
					"action":    "none",
					"status":    "already-onboarded",
				})
			}
			cli.OutputHuman("\nAccount %s is already onboarded; nothing to do.\n", repairAccountID)
			return nil
		}
		return errors.Wrap(err, "unable to create cloud account integration")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(
			repairRegisterResult("registered", repairAccountID, name, roleArn, externalID, resp.Data.IntgGuid))
	}
	cli.OutputHuman("\nThe cloud account integration was created with GUID %s.\n", resp.Data.IntgGuid)
	return nil
}

// repairRegisterResult builds the --json payload for the register path. The dry-run ("register") and
// applied ("registered") outcomes share the derived role/external id; the applied one also carries
// the server-assigned integration guid.
func repairRegisterResult(action, accountID, name, roleArn, externalID, intgGuid string) map[string]interface{} {
	r := map[string]interface{}{
		"accountId":  accountID,
		"action":     action,
		"name":       name,
		"roleArn":    roleArn,
		"externalId": externalID,
	}
	if intgGuid != "" {
		r["intgGuid"] = intgGuid
	} else {
		r["dryRun"] = true
	}
	return r
}

// stackParameter returns the value of a named parameter on a CloudFormation stack.
func stackParameter(ctx context.Context, cfn *cloudformation.Client, stackName, key string) (string, error) {
	out, err := cfn.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)})
	if err != nil {
		return "", errors.Wrapf(err, "unable to describe stack %s", stackName)
	}
	if len(out.Stacks) == 0 {
		return "", errors.Errorf("stack %s not found", stackName)
	}
	for _, p := range out.Stacks[0].Parameters {
		if aws.ToString(p.ParameterKey) == key {
			return aws.ToString(p.ParameterValue), nil
		}
	}
	return "", errors.Errorf("stack %s has no %s parameter", stackName, key)
}

// stackResourcePhysicalID returns the physical resource id for a logical resource on a stack.
func stackResourcePhysicalID(
	ctx context.Context, cfn *cloudformation.Client, stackName, logicalID string,
) (string, error) {
	out, err := cfn.DescribeStackResources(ctx, &cloudformation.DescribeStackResourcesInput{
		StackName:         aws.String(stackName),
		LogicalResourceId: aws.String(logicalID),
	})
	if err != nil {
		return "", errors.Wrapf(err, "unable to read %s from stack %s", logicalID, stackName)
	}
	if len(out.StackResources) == 0 {
		return "", errors.Errorf("stack %s has no %s resource", stackName, logicalID)
	}
	return aws.ToString(out.StackResources[0].PhysicalResourceId), nil
}

// memberStackID returns the StackId of the StackSet instance for the given account. found is false
// (with no error) when the account has no instance yet - the caller decides whether to create one.
func memberStackID(
	ctx context.Context, cfn *cloudformation.Client, stackSetName, accountID string,
) (stackID string, found bool, err error) {
	out, err := cfn.ListStackInstances(ctx, &cloudformation.ListStackInstancesInput{
		StackSetName:         aws.String(stackSetName),
		StackInstanceAccount: aws.String(accountID),
		CallAs:               cfntypes.CallAsSelf,
	})
	if err != nil {
		return "", false, errors.Wrapf(err, "unable to list stack instances for %s", stackSetName)
	}
	// ponytail: a present-but-FAILED instance can carry a StackId without a live role; we still try to
	// register (backend AssumeRole validation surfaces a clear error if the role is truly absent).
	for _, s := range out.Summaries {
		if id := aws.ToString(s.StackId); id != "" {
			return id, true, nil
		}
	}
	return "", false, nil
}

// targetOUsParam returns the OU ids the StackSet targets, from the management stack's
// OrganizationalUnits parameter. Returns nil if the parameter is absent.
func targetOUsParam(ctx context.Context, cfn *cloudformation.Client) []string {
	raw, err := stackParameter(ctx, cfn, repairStackName, "OrganizationalUnits")
	if err != nil {
		return nil
	}
	var ous []string
	for _, ou := range strings.Split(raw, ",") {
		if ou = strings.TrimSpace(ou); ou != "" {
			ous = append(ous, ou)
		}
	}
	return ous
}

// createMemberInstance re-deploys the member stack into one account by creating a stack instance
// scoped to that account within the targeted OUs, and waits for the operation to finish. INTERSECTION
// with a single account means nothing is deployed if the account is not under a targeted OU.
func createMemberInstance(
	ctx context.Context, cfn *cloudformation.Client, stackSetName, region string, ous []string, accountID string,
) error {
	out, err := cfn.CreateStackInstances(ctx, &cloudformation.CreateStackInstancesInput{
		StackSetName: aws.String(stackSetName),
		Regions:      []string{region},
		CallAs:       cfntypes.CallAsSelf,
		DeploymentTargets: &cfntypes.DeploymentTargets{
			OrganizationalUnitIds: ous,
			Accounts:              []string{accountID},
			AccountFilterType:     cfntypes.AccountFilterTypeIntersection,
		},
		OperationPreferences: &cfntypes.StackSetOperationPreferences{
			FailureToleranceCount: aws.Int32(0),
			MaxConcurrentCount:    aws.Int32(1),
		},
	})
	if err != nil {
		return errors.Wrapf(err, "unable to create stack instance for %s", accountID)
	}
	return pollStackSetOperation(ctx, cfn, stackSetName, aws.ToString(out.OperationId))
}

// pollStackSetOperation blocks until a StackSet operation finishes, returning an error on any
// non-SUCCEEDED terminal status.
func pollStackSetOperation(ctx context.Context, cfn *cloudformation.Client, stackSetName, operationID string) error {
	for {
		out, err := cfn.DescribeStackSetOperation(ctx, &cloudformation.DescribeStackSetOperationInput{
			StackSetName: aws.String(stackSetName),
			OperationId:  aws.String(operationID),
			CallAs:       cfntypes.CallAsSelf,
		})
		if err != nil {
			return errors.Wrap(err, "unable to poll stack instance operation")
		}
		switch out.StackSetOperation.Status {
		case cfntypes.StackSetOperationStatusRunning,
			cfntypes.StackSetOperationStatusQueued,
			cfntypes.StackSetOperationStatusStopping:
			time.Sleep(15 * time.Second)
		case cfntypes.StackSetOperationStatusSucceeded:
			return nil
		default:
			return errors.Errorf("stack instance operation %s ended with status %s: %s",
				operationID, out.StackSetOperation.Status, aws.ToString(out.StackSetOperation.StatusReason))
		}
	}
}

// stackUID extracts the role UID the config_org template derives from a stack id: the first
// hyphen-delimited segment of the stack's GUID (the last "/"-delimited segment of the stack id).
// e.g. arn:aws:cloudformation:us-west-2:123:stack/foo/abcd1234-ef56-... -> "abcd1234".
func stackUID(stackID string) string {
	slash := strings.Split(stackID, "/")
	guid := slash[len(slash)-1]
	if guid == "" {
		return ""
	}
	return strings.Split(guid, "-")[0]
}
