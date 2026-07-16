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
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/v2/api"
)

// repairMethodAwsCfgOrg repairs an account onboarded via the AWS Config organization template.
const repairMethodAwsCfgOrg = "aws-cfg-org"

// aws-cfg-org flags.
var (
	repairAccountID  string // --account-id
	repairAll        bool   // --all
	repairStackName  string // --stack-name (the AWS Config org management stack)
	repairAwsProfile string // --aws-profile
	repairAwsRegion  string // --aws-region
	repairDryRun     bool   // --dry-run
)

func init() {
	f := cloudAccountRepairCmd.Flags()
	f.StringVar(&repairAccountID,
		"account-id", "", "AWS account id whose AwsCfg integration is missing (required unless --all)")
	f.BoolVar(&repairAll,
		"all", false, "find and repair every account in the targeted OUs that is missing its member stack or integration")
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
	if repairAll && repairAccountID != "" {
		return errors.New("--account-id and --all are mutually exclusive")
	}
	if !repairAll && repairAccountID == "" {
		return errors.New("--account-id is required (or use --all)")
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

	if repairAll {
		return repairAwsCfgOrgAll(ctx, cfg, cfn, laceworkAccount, stackSetName)
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
		err = createMemberInstances(ctx, cfn, stackSetName, cfg.Region, targetOUs, []string{repairAccountID})
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

	name, roleArn, externalID, err := awsCfgRegistration(laceworkAccount, repairAccountID, stackID)
	if err != nil {
		return err
	}

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

	cli.StartProgress(" Creating integration...")
	intgGuid, alreadyOnboarded, err := createAwsCfgIntegration(name, roleArn, externalID, repairAccountID)
	cli.StopProgress()
	if err != nil {
		return err
	}
	if alreadyOnboarded {
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

	if cli.JSONOutput() {
		return cli.OutputJSON(
			repairRegisterResult("registered", repairAccountID, name, roleArn, externalID, intgGuid))
	}
	cli.OutputHuman("\nThe cloud account integration was created with GUID %s.\n", intgGuid)
	return nil
}

// awsCfgRegistration derives the integration name, role ARN and external id for one member account
// from its stack instance id, exactly as the config_org member template builds them
// (role: lacework-config-role-<UID>; externalId: lweid:aws:v2:<tenant>:<acct>:LW<UID>).
func awsCfgRegistration(laceworkAccount, accountID, stackID string) (name, roleArn, externalID string, err error) {
	uid := stackUID(stackID)
	if uid == "" {
		return "", "", "", errors.Errorf("could not derive the role UID from stack id %q", stackID)
	}
	name = fmt.Sprintf("%s-Config", laceworkAccount)
	roleArn = fmt.Sprintf("arn:aws:iam::%s:role/lacework-config-role-%s", accountID, uid)
	externalID = fmt.Sprintf("lweid:aws:v2:%s:%s:LW%s", laceworkAccount, accountID, uid)
	return name, roleArn, externalID, nil
}

// createAwsCfgIntegration registers the AwsCfg integration with the Lacework API. Mirroring the
// setup Lambda, an already-registered account is success (alreadyOnboarded), not an error.
func createAwsCfgIntegration(name, roleArn, externalID, accountID string) (
	intgGuid string, alreadyOnboarded bool, err error,
) {
	awsCfg := api.NewCloudAccount(name, api.AwsCfgCloudAccount, api.AwsCfgData{
		AwsAccountID: accountID,
		Credentials: api.AwsCfgCredentials{
			RoleArn:    roleArn,
			ExternalID: externalID,
		},
	})
	resp, err := cli.LwApi.V2.CloudAccounts.Create(awsCfg)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "aws account is already used") {
			return "", true, nil
		}
		return "", false, errors.Wrap(err, "unable to create cloud account integration")
	}
	return resp.Data.IntgGuid, false, nil
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

// createMemberInstances re-deploys the member stack into the given accounts by creating stack
// instances scoped to those accounts within the targeted OUs, and waits for the operation to finish.
// INTERSECTION means nothing is deployed to an account that is not under a targeted OU.
// ponytail: FailureToleranceCount 0 / MaxConcurrentCount 1 aborts the whole batch on the first bad
// account (re-run friendly); raise both if fleets get large enough to hurt.
func createMemberInstances(
	ctx context.Context, cfn *cloudformation.Client, stackSetName, region string, ous, accountIDs []string,
) error {
	out, err := cfn.CreateStackInstances(ctx, &cloudformation.CreateStackInstancesInput{
		StackSetName: aws.String(stackSetName),
		Regions:      []string{region},
		CallAs:       cfntypes.CallAsSelf,
		DeploymentTargets: &cfntypes.DeploymentTargets{
			OrganizationalUnitIds: ous,
			Accounts:              accountIDs,
			AccountFilterType:     cfntypes.AccountFilterTypeIntersection,
		},
		OperationPreferences: &cfntypes.StackSetOperationPreferences{
			FailureToleranceCount: aws.Int32(0),
			MaxConcurrentCount:    aws.Int32(1),
		},
	})
	if err != nil {
		return errors.Wrapf(err, "unable to create stack instances for %s", strings.Join(accountIDs, ", "))
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

// repairAwsCfgOrgAll (--all) sweeps every ACTIVE account under the StackSet's targeted OUs, finds
// the ones missing their member stack instance and/or AwsCfg integration, rebuilds the missing
// instances in one StackSet operation, and registers every unregistered account.
func repairAwsCfgOrgAll(
	ctx context.Context, cfg aws.Config, cfn *cloudformation.Client, laceworkAccount, stackSetName string,
) error {
	targetOUs := targetOUsParam(ctx, cfn)
	if len(targetOUs) == 0 {
		return errors.Errorf(
			"stack %s has no OrganizationalUnits parameter; there are no targeted OUs to sweep", repairStackName)
	}

	org := organizations.NewFromConfig(cfg)

	cli.StartProgress(" Enumerating accounts in the targeted OUs...")
	expected, err := accountsInOUs(ctx, org, targetOUs)
	cli.StopProgress()
	if err != nil {
		return err
	}

	cli.StartProgress(" Listing member stack instances and existing integrations...")
	instances, err := stackInstanceIDs(ctx, cfn, stackSetName)
	var registered map[string]bool
	if err == nil {
		registered, err = registeredAwsCfgAccounts()
	}
	cli.StopProgress()
	if err != nil {
		return err
	}

	missingInstance, missingIntegration := diffRepairTargets(expected, instances, registered)

	if len(missingInstance) == 0 && len(missingIntegration) == 0 {
		if cli.JSONOutput() {
			return cli.OutputJSON(map[string]interface{}{
				"targetOUs":        targetOUs,
				"expectedAccounts": len(expected),
				"repaired":         []interface{}{},
			})
		}
		cli.OutputHuman("All %d accounts in the targeted OUs are onboarded; nothing to do.\n", len(expected))
		return nil
	}

	if repairDryRun {
		if cli.JSONOutput() {
			return cli.OutputJSON(map[string]interface{}{
				"targetOUs":            targetOUs,
				"expectedAccounts":     len(expected),
				"missingStackInstance": missingInstance,
				"missingIntegration":   missingIntegration,
				"dryRun":               true,
			})
		}
		cli.OutputHuman("Found %d accounts in the targeted OUs (%s):\n", len(expected), strings.Join(targetOUs, ","))
		cli.OutputHuman("  missing stack instance (role gone): %s\n", orNone(missingInstance))
		cli.OutputHuman("  missing integration only:           %s\n", orNone(missingIntegration))
		cli.OutputHuman("\nDry-run: nothing was created. Re-run without --dry-run to repair.\n")
		return nil
	}

	if len(missingInstance) > 0 {
		cli.OutputHuman("Rebuilding %d member stack instance(s): %s\n",
			len(missingInstance), strings.Join(missingInstance, ", "))
		cli.StartProgress(" Waiting for the stack instances (rebuilding IAM roles)...")
		err = createMemberInstances(ctx, cfn, stackSetName, cfg.Region, targetOUs, missingInstance)
		cli.StopProgress()
		if err != nil {
			return err
		}
		instances, err = stackInstanceIDs(ctx, cfn, stackSetName)
		if err != nil {
			return err
		}
	}

	// Register everything unregistered. For just-rebuilt accounts the template's setup Lambda
	// usually wins the race, which surfaces here as already-onboarded - still a full recovery.
	var results []map[string]interface{}
	failed := 0
	for _, accountID := range append(missingInstance, missingIntegration...) {
		stackID, ok := instances[accountID]
		if !ok {
			results = append(results, map[string]interface{}{
				"accountId": accountID, "action": "error",
				"error": "no stack instance was created; is the account still in a targeted OU?",
			})
			failed++
			continue
		}
		result := repairOneRegistration(laceworkAccount, accountID, stackID)
		if result["action"] == "error" {
			failed++
		}
		results = append(results, result)
	}

	if cli.JSONOutput() {
		err = cli.OutputJSON(map[string]interface{}{
			"targetOUs":        targetOUs,
			"expectedAccounts": len(expected),
			"repaired":         results,
		})
		if err != nil {
			return err
		}
	} else {
		cli.OutputHuman("\n")
		for _, r := range results {
			switch r["action"] {
			case "registered":
				cli.OutputHuman("  %s: registered (GUID %s)\n", r["accountId"], r["intgGuid"])
			case "already-onboarded":
				cli.OutputHuman("  %s: already onboarded\n", r["accountId"])
			default:
				cli.OutputHuman("  %s: FAILED - %s\n", r["accountId"], r["error"])
			}
		}
		cli.OutputHuman("\nRepaired %d of %d account(s).\n", len(results)-failed, len(results))
	}
	if failed > 0 {
		return errors.Errorf("%d of %d repairs failed", failed, len(results))
	}
	return nil
}

// repairOneRegistration derives the registration for one account and creates the integration,
// returning a result map for the summary (never an error - bulk repair reports and moves on).
func repairOneRegistration(laceworkAccount, accountID, stackID string) map[string]interface{} {
	name, roleArn, externalID, err := awsCfgRegistration(laceworkAccount, accountID, stackID)
	if err != nil {
		return map[string]interface{}{"accountId": accountID, "action": "error", "error": err.Error()}
	}
	intgGuid, alreadyOnboarded, err := createAwsCfgIntegration(name, roleArn, externalID, accountID)
	if err != nil {
		return map[string]interface{}{"accountId": accountID, "action": "error", "error": err.Error()}
	}
	if alreadyOnboarded {
		return map[string]interface{}{"accountId": accountID, "action": "already-onboarded"}
	}
	return repairRegisterResult("registered", accountID, name, roleArn, externalID, intgGuid)
}

// diffRepairTargets classifies the expected accounts against the current stack instances and
// Lacework integrations: no instance means the IAM role is gone (rebuild + register); an instance
// without an integration means register only.
func diffRepairTargets(
	expected []string, instances map[string]string, registered map[string]bool,
) (missingInstance, missingIntegration []string) {
	for _, accountID := range expected {
		if _, ok := instances[accountID]; !ok {
			missingInstance = append(missingInstance, accountID)
		} else if !registered[accountID] {
			missingIntegration = append(missingIntegration, accountID)
		}
	}
	return missingInstance, missingIntegration
}

// accountsInOUs returns every ACTIVE account id under the given OUs, recursing into child OUs to
// mirror StackSets' recursive OU targeting. The organization management account is excluded:
// service-managed StackSets never deploy to it, so it can never have a member instance.
func accountsInOUs(ctx context.Context, org *organizations.Client, ous []string) ([]string, error) {
	mgmtOut, err := org.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		return nil, errors.Wrap(err, "unable to describe the organization")
	}
	mgmtAccountID := aws.ToString(mgmtOut.Organization.MasterAccountId)

	seen := map[string]bool{}
	var accounts []string
	queue := append([]string{}, ous...)
	for len(queue) > 0 {
		parent := queue[0]
		queue = queue[1:]

		accPag := organizations.NewListAccountsForParentPaginator(org,
			&organizations.ListAccountsForParentInput{ParentId: aws.String(parent)})
		for accPag.HasMorePages() {
			page, err := accPag.NextPage(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to list accounts under %s", parent)
			}
			for _, a := range page.Accounts {
				id := aws.ToString(a.Id)
				if a.Status == orgtypes.AccountStatusActive && id != mgmtAccountID && !seen[id] {
					seen[id] = true
					accounts = append(accounts, id)
				}
			}
		}

		ouPag := organizations.NewListOrganizationalUnitsForParentPaginator(org,
			&organizations.ListOrganizationalUnitsForParentInput{ParentId: aws.String(parent)})
		for ouPag.HasMorePages() {
			page, err := ouPag.NextPage(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to list child OUs under %s", parent)
			}
			for _, child := range page.OrganizationalUnits {
				queue = append(queue, aws.ToString(child.Id))
			}
		}
	}
	sort.Strings(accounts)
	return accounts, nil
}

// stackInstanceIDs maps member account id -> stack id for every instance in the StackSet.
func stackInstanceIDs(
	ctx context.Context, cfn *cloudformation.Client, stackSetName string,
) (map[string]string, error) {
	instances := map[string]string{}
	pag := cloudformation.NewListStackInstancesPaginator(cfn, &cloudformation.ListStackInstancesInput{
		StackSetName: aws.String(stackSetName),
		CallAs:       cfntypes.CallAsSelf,
	})
	for pag.HasMorePages() {
		page, err := pag.NextPage(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to list stack instances for %s", stackSetName)
		}
		for _, s := range page.Summaries {
			if id := aws.ToString(s.StackId); id != "" {
				instances[aws.ToString(s.Account)] = id
			}
		}
	}
	return instances, nil
}

// registeredAwsCfgAccounts returns the AWS account ids that already have an AwsCfg integration,
// derived from each integration's cross-account role ARN.
func registeredAwsCfgAccounts() (map[string]bool, error) {
	res, err := cli.LwApi.V2.CloudAccounts.ListByType(api.AwsCfgCloudAccount)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list AwsCfg cloud accounts")
	}
	registered := map[string]bool{}
	for _, raw := range res.Data {
		if accountID, _, ok := deriveAwsAccountAndRole(raw); ok {
			registered[accountID] = true
		}
	}
	return registered, nil
}

// orNone renders an account list for human output, with an explicit marker when empty.
func orNone(accounts []string) string {
	if len(accounts) == 0 {
		return "(none)"
	}
	return strings.Join(accounts, ", ")
}
