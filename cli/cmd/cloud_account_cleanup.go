//
// Copyright:: Copyright 2026, Lacework Inc.
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
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/v2/api"
)

// cleanupTarget is a unique IAM role to remove, derived from a backed-up integration's role ARN.
type cleanupTarget struct {
	accountID string
	roleName  string
}

// cleanup command flags
var (
	cleanupFile           string // -f/--file (input)
	cleanupType           string // -t/--type (optional filter)
	cleanupAssumeRoleName string // --assume-role-name
	cleanupAwsProfile     string // --aws-profile
	cleanupAwsRegion      string // --aws-region
	cleanupNoAssumeRole   bool   // --no-assume-role
	cleanupDryRun         bool   // --dry-run
)

func init() {
	cloudAccountCommand.AddCommand(cloudAccountCleanupCmd)

	cloudAccountCleanupCmd.Flags().StringVarP(&cleanupFile,
		"file", "f", "", "backup file identifying the resources to clean up")
	cloudAccountCleanupCmd.Flags().StringVarP(&cleanupType,
		"type", "t", "", "only clean up records of this type (selects the cleanup strategy)")
	cloudAccountCleanupCmd.Flags().StringVar(&cleanupAssumeRoleName,
		"assume-role-name", "OrganizationAccountAccessRole",
		"IAM role name to assume in each account for cleanup")
	cloudAccountCleanupCmd.Flags().StringVar(&cleanupAwsProfile,
		"aws-profile", "", "AWS profile to use as the base credentials for cleanup")
	cloudAccountCleanupCmd.Flags().StringVar(&cleanupAwsRegion,
		"aws-region", "", "AWS region for the base STS client and the first region searched for stacks")
	cloudAccountCleanupCmd.Flags().BoolVar(&cleanupNoAssumeRole,
		"no-assume-role", false,
		"use the current AWS credentials directly instead of assuming a role in each account")
	cloudAccountCleanupCmd.Flags().BoolVar(&cleanupDryRun,
		"dry-run", false, "show what would be deleted without deleting")
}

var cloudAccountCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Delete leftover cloud resources for the integrations in a backup file",
	Long: `Delete the cloud-side resources that a deleted integration left behind. For AWS cross-account
integrations this removes the IAM role and the customer-managed/inline policies created for the
integration; AWS-managed policies (such as SecurityAudit) are only detached, never deleted.

If the integration was onboarded with the Lacework CloudFormation template, the owning
CloudFormation stack is detected (via the role) and deleted as well, so no orphaned/drifted stack is
left behind. Integrations not created by CloudFormation (e.g. Terraform, console, or API) have no
stack and get IAM cleanup only. The IAM cleanup and stack deletion are independent - a failure of one
does not block the other.

The role to remove is taken from each backed-up integration's role ARN, so no resource names are
guessed. By default the command assumes a role (default OrganizationAccountAccessRole) into each
target account, so run it with credentials for your AWS Organizations management account. If a target
account is the same account your credentials are already in (for example the management account's own
integration), the assume-role is skipped and the current credentials are used directly. Pass
--no-assume-role to always use the current credentials without assuming a role.

    lacework cloud-account cleanup --file backup.json --aws-profile management
    lacework cloud-account cleanup --file backup.json --type AwsCfg --dry-run
    lacework cloud-account cleanup --file backup.json --no-assume-role   # run directly in one account

Use --type to restrict cleanup to a single integration type (it also selects the cleanup strategy).
Types without a cleanup implementation are reported and skipped.`,
	Args: cobra.NoArgs,
	RunE: cloudAccountCleanup,
}

func cloudAccountCleanup(_ *cobra.Command, _ []string) error {
	if cleanupFile == "" {
		return errors.New("--file is required")
	}

	// Optional type filter, which also validates that the type supports cleanup.
	typeFilter := ""
	if cleanupType != "" {
		caType, found := api.FindCloudAccountType(cleanupType)
		if !found {
			return errors.Errorf("unknown cloud account type '%s'", cleanupType)
		}
		if !isAwsCrossAccountType(caType.String()) {
			return errors.Errorf("cleanup is not supported for type '%s'", caType.String())
		}
		typeFilter = caType.String()
	}

	backup, err := readCloudAccountBackup(cleanupFile)
	if err != nil {
		return err
	}

	targets, unsupported, unparseable := collectCleanupTargets(backup, typeFilter)

	for t, n := range unsupported {
		cli.OutputHuman("Skipping %d record(s) of type %s: cleanup not supported for this type.\n", n, t)
	}
	if unparseable > 0 {
		cli.OutputHuman("Skipping %d record(s) with an unparseable role ARN.\n", unparseable)
	}
	if len(targets) == 0 {
		cli.OutputHuman("No cloud resources to clean up.\n")
		return nil
	}

	cli.OutputHuman("%d IAM role(s) targeted for cleanup:\n", len(targets))
	for _, t := range targets {
		cli.OutputHuman("  account=%s role=%s\n", t.accountID, t.roleName)
	}

	if !cleanupDryRun {
		if !confirmBulkOperation(fmt.Sprintf(
			"Delete IAM role(s)/policies and any owning CloudFormation stack in %d target(s)?", len(targets))) {
			cli.OutputHuman("Aborted. No resources were deleted.\n")
			return nil
		}
	}

	ctx := context.Background()
	base, err := buildAwsBaseConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to load AWS credentials")
	}

	// Learn which account the base credentials belong to, so we can skip the assume-role for any
	// target that is already that account (e.g. the management account's own integration). If this
	// lookup fails we fall back to assuming for every cross-account target.
	baseAccountID := ""
	if !cleanupNoAssumeRole {
		if baseAccountID, err = callerAccountID(ctx, base); err != nil {
			cli.OutputHuman("Warning: could not determine the current AWS account (%s); "+
				"will assume %s for every target.\n", err.Error(), cleanupAssumeRoleName)
		}
	}

	finder := newStackFinder(ctx, base)

	var cleaned, missing, skipped int
	var stacksDeleted, stacksFailed int
	for _, t := range targets {
		cli.OutputHuman("\naccount=%s role=%s\n", t.accountID, t.roleName)
		cfg, err := accountConfigForTarget(ctx, base, baseAccountID, t.accountID)
		if err != nil {
			cli.OutputHuman("  skipped: unable to assume %s in %s: %s\n",
				cleanupAssumeRoleName, t.accountID, err.Error())
			skipped++
			continue
		}
		iamClient := iam.NewFromConfig(cfg)

		// Pass 0 (read-only, before deleting the role): find the owning CloudFormation stack, if any,
		// by querying DescribeStackResources for the role across regions. Must run while the role
		// still exists; any detection error is a warning and never blocks the IAM cleanup below.
		cli.StartProgress(fmt.Sprintf(" Searching for the CloudFormation stack owning %s...", t.roleName))
		stackName, stackRegion, derr := finder.find(ctx, cfg, t.roleName)
		cli.StopProgress()
		if derr != nil {
			cli.OutputHuman("  warning: could not check for a CloudFormation stack: %s\n", derr.Error())
		}

		// Pass 1: delete the role and its policies (unchanged behavior).
		found, err := cleanupIamRole(ctx, iamClient, t.roleName, !cleanupDryRun)
		switch {
		case err != nil:
			cli.OutputHuman("  error: %s\n", err.Error())
			skipped++
		case !found:
			cli.OutputHuman("  role not found - already cleaned up.\n")
			missing++
		default:
			cleaned++
		}

		// Pass 2: delete the owning stack (best-effort, independent of the IAM result above), using a
		// CloudFormation client in the stack's own region.
		if stackName != "" {
			cfnCfg := cfg.Copy()
			if stackRegion != "" {
				cfnCfg.Region = stackRegion
			}
			cfnClient := cloudformation.NewFromConfig(cfnCfg)
			if err := deleteCfnStack(ctx, cfnClient, stackName, !cleanupDryRun); err != nil {
				cli.OutputHuman("  stack cleanup failed for %s: %s\n", stackName, err.Error())
				stacksFailed++
			} else {
				stacksDeleted++
			}
		}
	}

	roleVerb, stackVerb := "cleaned", "deleted"
	if cleanupDryRun {
		roleVerb, stackVerb = "would clean", "would delete"
	}
	cli.OutputHuman("\n%s %d role(s); %d already gone; %d skipped.\n", roleVerb, cleaned, missing, skipped)
	cli.OutputHuman("%s %d CloudFormation stack(s); %d failed.\n", stackVerb, stacksDeleted, stacksFailed)
	if cleanupDryRun {
		cli.OutputHuman("Dry-run only. Re-run without --dry-run to delete.\n")
	}
	return nil
}

// collectCleanupTargets dedups (account, role) targets from the backup, honoring an optional type
// filter, and returns counts of records skipped as unsupported-type or unparseable-ARN.
func collectCleanupTargets(backup cloudAccountBackup, typeFilter string) (
	[]cleanupTarget, map[string]int, int,
) {
	var targets []cleanupTarget
	seen := map[string]struct{}{}
	unsupported := map[string]int{}
	unparseable := 0

	for _, intg := range backup.Integrations {
		if typeFilter != "" && intg.Type != typeFilter {
			continue
		}
		if !isAwsCrossAccountType(intg.Type) {
			unsupported[intg.Type]++
			continue
		}
		accountID, roleName, ok := deriveAwsAccountAndRole(intg)
		if !ok {
			unparseable++
			continue
		}
		key := accountID + "/" + roleName
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		targets = append(targets, cleanupTarget{accountID: accountID, roleName: roleName})
	}
	return targets, unsupported, unparseable
}

// buildAwsBaseConfig loads the base AWS config from the shared credential chain, honoring the
// --aws-profile/--aws-region flags. IAM/STS need a region even though IAM is global.
func buildAwsBaseConfig(ctx context.Context) (aws.Config, error) {
	region := cleanupAwsRegion
	if region == "" {
		region = "us-east-1"
	}
	opts := []func(*config.LoadOptions) error{config.WithRegion(region)}
	if cleanupAwsProfile != "" {
		opts = append(opts, config.WithSharedConfigProfile(cleanupAwsProfile))
	}
	return config.LoadDefaultConfig(ctx, opts...)
}

// callerAccountID returns the AWS account id that the given credentials belong to.
func callerAccountID(ctx context.Context, cfg aws.Config) (string, error) {
	out, err := sts.NewFromConfig(cfg).GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return aws.ToString(out.Account), nil
}

// accountConfigForTarget returns the aws.Config to use for a target account. It uses the base
// credentials directly when --no-assume-role is set or when the target is already the base
// credentials' own account (e.g. the management account's own integration); otherwise it assumes
// <assume-role-name> in the target account. The returned config lets the caller build both IAM and
// CloudFormation clients.
func accountConfigForTarget(ctx context.Context, base aws.Config, baseAccountID, accountID string) (aws.Config, error) {
	if cleanupNoAssumeRole || (baseAccountID != "" && accountID == baseAccountID) {
		return base, nil
	}
	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, cleanupAssumeRoleName)
	provider := stscreds.NewAssumeRoleProvider(sts.NewFromConfig(base), roleArn)
	cfg := base.Copy()
	cfg.Credentials = aws.NewCredentialsCache(provider)
	// Fail fast if the assume-role can't be resolved.
	if _, err := cfg.Credentials.Retrieve(ctx); err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}

// standardRegions is the fallback set of regions to search when the account's enabled regions can't
// be listed (e.g. missing ec2:DescribeRegions permission).
var standardRegions = []string{
	"us-east-1", "us-east-2", "us-west-1", "us-west-2",
	"ca-central-1", "sa-east-1",
	"eu-west-1", "eu-west-2", "eu-west-3", "eu-central-1", "eu-north-1", "eu-south-1",
	"ap-south-1", "ap-southeast-1", "ap-southeast-2",
	"ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-east-1",
	"me-south-1", "af-south-1",
}

// stackFinder locates the CloudFormation stack that owns an IAM role. CloudFormation is regional and
// IAM roles do not carry the aws:cloudformation:* tags, so we can't read the region off the role;
// instead we query DescribeStackResources by physical resource id (the role name) across candidate
// regions. The region of the first hit is promoted to the front, so once one stack is found the rest
// of the run usually resolves in a single call.
type stackFinder struct {
	regions []string
}

func newStackFinder(ctx context.Context, base aws.Config) *stackFinder {
	var regions []string
	if out, err := ec2.NewFromConfig(base).DescribeRegions(ctx, &ec2.DescribeRegionsInput{}); err == nil {
		for _, r := range out.Regions {
			if n := aws.ToString(r.RegionName); n != "" {
				regions = append(regions, n)
			}
		}
	}
	if len(regions) == 0 {
		regions = append(regions, standardRegions...)
	}
	// Search the explicitly requested / base region first.
	preferred := cleanupAwsRegion
	if preferred == "" {
		preferred = base.Region
	}
	return &stackFinder{regions: moveRegionToFront(regions, preferred)}
}

// find returns the name and region of the stack that owns roleName in the given account config, or
// empty strings if no stack manages it (e.g. a non-CloudFormation role). Must run while the role
// still exists.
func (f *stackFinder) find(ctx context.Context, cfg aws.Config, roleName string) (string, string, error) {
	var lastErr error
	for _, region := range f.regions {
		rc := cfg.Copy()
		rc.Region = region
		out, err := cloudformation.NewFromConfig(rc).DescribeStackResources(ctx,
			&cloudformation.DescribeStackResourcesInput{PhysicalResourceId: &roleName})
		if err != nil {
			// A "does not exist" ValidationError just means no stack in this region; keep looking.
			if !strings.Contains(err.Error(), "does not exist") {
				lastErr = err
			}
			continue
		}
		for _, r := range out.StackResources {
			if r.StackName != nil && *r.StackName != "" {
				f.regions = moveRegionToFront(f.regions, region)
				return *r.StackName, region, nil
			}
		}
	}
	// Only surface an error if nothing was found AND a non-"not found" error occurred.
	return "", "", lastErr
}

func moveRegionToFront(regions []string, region string) []string {
	if region == "" {
		return regions
	}
	out := make([]string, 0, len(regions)+1)
	out = append(out, region)
	for _, r := range regions {
		if r != region {
			out = append(out, r)
		}
	}
	return out
}

// deleteCfnStack deletes a CloudFormation stack and waits for the deletion to complete. When apply is
// false it only logs what it would do. Deleting a stack whose resources were already removed
// out-of-band is idempotent - CloudFormation treats the missing resources as already deleted.
func deleteCfnStack(ctx context.Context, c *cloudformation.Client, stackName string, apply bool) error {
	logCleanup(apply, "delete CloudFormation stack", stackName)
	if !apply {
		return nil
	}
	if _, err := c.DeleteStack(ctx, &cloudformation.DeleteStackInput{StackName: &stackName}); err != nil {
		return err
	}
	waiter := cloudformation.NewStackDeleteCompleteWaiter(c)
	return waiter.Wait(ctx, &cloudformation.DescribeStacksInput{StackName: &stackName}, 10*time.Minute)
}

// cleanupIamRole detaches/deletes the role's policies and the role itself. When apply is false it
// only inspects and logs what it would do. Returns false if the role does not exist.
func cleanupIamRole(ctx context.Context, c *iam.Client, roleName string, apply bool) (bool, error) {
	if _, err := c.GetRole(ctx, &iam.GetRoleInput{RoleName: &roleName}); err != nil {
		if isNoSuchEntity(err) {
			return false, nil
		}
		return false, err
	}

	// Gather the role's policies up front so we can report a count before touching anything.
	attached, err := collectAttachedPolicies(ctx, c, roleName)
	if err != nil {
		return true, err
	}
	inline, err := collectInlinePolicies(ctx, c, roleName)
	if err != nil {
		return true, err
	}

	customerManaged := 0
	for _, p := range attached {
		if isCustomerManagedPolicy(aws.ToString(p.PolicyArn)) {
			customerManaged++
		}
	}
	cli.OutputHuman("  %d policy(ies) to delete (%d customer-managed, %d inline); %d AWS-managed to detach only\n",
		customerManaged+len(inline), customerManaged, len(inline), len(attached)-customerManaged)

	// Attached managed policies: detach all; delete the customer-managed ones.
	for _, p := range attached {
		arnStr := aws.ToString(p.PolicyArn)
		name := aws.ToString(p.PolicyName)
		if apply {
			if _, err := c.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
				RoleName: &roleName, PolicyArn: p.PolicyArn,
			}); err != nil && !isNoSuchEntity(err) {
				return true, err
			}
		}
		logCleanup(apply, "detach", name)
		if isCustomerManagedPolicy(arnStr) {
			if apply {
				if err := deleteCustomerManagedPolicy(ctx, c, arnStr); err != nil {
					return true, err
				}
			}
			logCleanup(apply, "delete customer-managed policy", name)
		} else {
			cli.OutputHuman("  keep AWS-managed policy %s (detached only)\n", name)
		}
	}

	// Inline policies.
	for _, name := range inline {
		if apply {
			n := name
			if _, err := c.DeleteRolePolicy(ctx, &iam.DeleteRolePolicyInput{
				RoleName: &roleName, PolicyName: &n,
			}); err != nil && !isNoSuchEntity(err) {
				return true, err
			}
		}
		logCleanup(apply, "delete inline policy", name)
	}

	// The role itself.
	if apply {
		if _, err := c.DeleteRole(ctx, &iam.DeleteRoleInput{RoleName: &roleName}); err != nil && !isNoSuchEntity(err) {
			return true, err
		}
	}
	logCleanup(apply, "delete role", roleName)
	return true, nil
}

// collectAttachedPolicies returns all managed policies attached to the role.
func collectAttachedPolicies(ctx context.Context, c *iam.Client, roleName string) ([]iamtypes.AttachedPolicy, error) {
	var out []iamtypes.AttachedPolicy
	pager := iam.NewListAttachedRolePoliciesPaginator(c, &iam.ListAttachedRolePoliciesInput{RoleName: &roleName})
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		out = append(out, page.AttachedPolicies...)
	}
	return out, nil
}

// collectInlinePolicies returns the names of all inline policies on the role.
func collectInlinePolicies(ctx context.Context, c *iam.Client, roleName string) ([]string, error) {
	var out []string
	pager := iam.NewListRolePoliciesPaginator(c, &iam.ListRolePoliciesInput{RoleName: &roleName})
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		out = append(out, page.PolicyNames...)
	}
	return out, nil
}

// deleteCustomerManagedPolicy removes non-default versions then deletes the policy. A DeleteConflict
// (still attached elsewhere) is reported rather than treated as fatal.
func deleteCustomerManagedPolicy(ctx context.Context, c *iam.Client, policyArn string) error {
	versPager := iam.NewListPolicyVersionsPaginator(c, &iam.ListPolicyVersionsInput{PolicyArn: &policyArn})
	for versPager.HasMorePages() {
		page, err := versPager.NextPage(ctx)
		if err != nil {
			if isNoSuchEntity(err) {
				return nil
			}
			return err
		}
		for _, v := range page.Versions {
			if v.IsDefaultVersion {
				continue
			}
			if _, err := c.DeletePolicyVersion(ctx, &iam.DeletePolicyVersionInput{
				PolicyArn: &policyArn, VersionId: v.VersionId,
			}); err != nil && !isNoSuchEntity(err) {
				return err
			}
		}
	}
	if _, err := c.DeletePolicy(ctx, &iam.DeletePolicyInput{PolicyArn: &policyArn}); err != nil {
		if isNoSuchEntity(err) {
			return nil
		}
		if isDeleteConflict(err) {
			cli.OutputHuman("  WARNING: policy %s still attached elsewhere - left in place.\n", policyArn)
			return nil
		}
		return err
	}
	return nil
}

func logCleanup(apply bool, action, name string) {
	if apply {
		cli.OutputHuman("  %s %s\n", action, name)
	} else {
		cli.OutputHuman("  would %s %s\n", action, name)
	}
}

// isCustomerManagedPolicy reports whether a policy ARN is customer-managed (deletable). AWS-managed
// policy ARNs use the ":iam::aws:policy/" form and must never be deleted.
func isCustomerManagedPolicy(policyArn string) bool {
	return !strings.Contains(policyArn, ":iam::aws:policy/")
}

func isNoSuchEntity(err error) bool {
	var e *iamtypes.NoSuchEntityException
	return errors.As(err, &e)
}

func isDeleteConflict(err error) bool {
	var e *iamtypes.DeleteConflictException
	return errors.As(err, &e)
}
