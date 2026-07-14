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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
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
		"aws-region", "", "AWS region for the IAM/STS clients")
	cloudAccountCleanupCmd.Flags().BoolVar(&cleanupDryRun,
		"dry-run", false, "show what would be deleted without deleting")
}

var cloudAccountCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Delete leftover cloud resources for the integrations in a backup file",
	Long: `Delete the cloud-side resources that a deleted integration left behind. For AWS cross-account
integrations this removes the IAM role and the customer-managed/inline policies created for the
integration; AWS-managed policies (such as SecurityAudit) are only detached, never deleted.

The role to remove is taken from each backed-up integration's role ARN, so no resource names are
guessed. The command assumes a role (default OrganizationAccountAccessRole) into each account, so run
it with credentials for your AWS Organizations management account.

    lacework cloud-account cleanup --file backup.json --aws-profile management
    lacework cloud-account cleanup --file backup.json --type AwsCfg --dry-run

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
		if !confirmBulkOperation(fmt.Sprintf("Delete IAM role(s) and their policies in %d target(s)?", len(targets))) {
			cli.OutputHuman("Aborted. No resources were deleted.\n")
			return nil
		}
	}

	ctx := context.Background()
	base, err := buildAwsBaseConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to load AWS credentials")
	}

	var cleaned, missing, skipped int
	for _, t := range targets {
		cli.OutputHuman("\naccount=%s role=%s\n", t.accountID, t.roleName)
		iamClient, err := assumeIntoAccount(ctx, base, t.accountID)
		if err != nil {
			cli.OutputHuman("  skipped: unable to assume %s in %s: %s\n",
				cleanupAssumeRoleName, t.accountID, err.Error())
			skipped++
			continue
		}
		found, err := cleanupIamRole(ctx, iamClient, t.roleName, !cleanupDryRun)
		if err != nil {
			cli.OutputHuman("  error: %s\n", err.Error())
			skipped++
			continue
		}
		if !found {
			cli.OutputHuman("  role not found - already cleaned up.\n")
			missing++
			continue
		}
		cleaned++
	}

	verb := "cleaned"
	if cleanupDryRun {
		verb = "would clean"
	}
	cli.OutputHuman("\n%s %d role(s); %d already gone; %d skipped.\n", verb, cleaned, missing, skipped)
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

// assumeIntoAccount assumes <assume-role-name> in the target account and returns an IAM client.
func assumeIntoAccount(ctx context.Context, base aws.Config, accountID string) (*iam.Client, error) {
	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, cleanupAssumeRoleName)
	provider := stscreds.NewAssumeRoleProvider(sts.NewFromConfig(base), roleArn)
	cfg := base.Copy()
	cfg.Credentials = aws.NewCredentialsCache(provider)
	// Fail fast if the assume-role can't be resolved.
	if _, err := cfg.Credentials.Retrieve(ctx); err != nil {
		return nil, err
	}
	return iam.NewFromConfig(cfg), nil
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

	// Attached managed policies: detach all; delete the customer-managed ones.
	attachPager := iam.NewListAttachedRolePoliciesPaginator(c, &iam.ListAttachedRolePoliciesInput{RoleName: &roleName})
	for attachPager.HasMorePages() {
		page, err := attachPager.NextPage(ctx)
		if err != nil {
			return true, err
		}
		for _, p := range page.AttachedPolicies {
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
	}

	// Inline policies.
	inlinePager := iam.NewListRolePoliciesPaginator(c, &iam.ListRolePoliciesInput{RoleName: &roleName})
	for inlinePager.HasMorePages() {
		page, err := inlinePager.NextPage(ctx)
		if err != nil {
			return true, err
		}
		for _, name := range page.PolicyNames {
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
