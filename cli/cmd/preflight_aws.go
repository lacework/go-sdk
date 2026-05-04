//
// Author:: Fortinet
// Copyright:: Copyright 2026, Fortinet
// License:: Apache License, Version 2.0
//

package cmd

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/v2/lwpreflight/aws"
)

var (
	preflightAwsState struct {
		agentless       bool
		config          bool
		cloudtrail      bool
		eksAuditLog     bool
		isOrg           bool
		simulate        bool
		region          string
		profile         string
		accessKeyID     string
		secretAccessKey string
		sessionToken    string
	}

	preflightAwsCmd = &cobra.Command{
		Use:   "aws",
		Short: "Run preflight checks against an AWS account",
		Long: `Run preflight checks against an AWS account to verify the caller has the
permissions required by the selected Lacework integrations. Credentials are
resolved using the standard AWS SDK chain (environment variables, shared
config files, EC2 instance profile) unless explicit --profile or
--access-key-id/--secret-access-key flags are provided.

At least one integration flag must be set: --agentless, --config,
--cloudtrail, or --eks-audit-log.

By default, the caller's identity-based policies are inspected locally. Pass
--simulate to evaluate each required action through the IAM policy simulator
instead — this also accounts for permissions boundaries and unconditional
Organizations service control policies, which a local policy walk cannot see.
Note: the simulator skips SCPs that have any conditions, and does not
evaluate resource control policies (RCPs). Condition keys (e.g. aws:SourceIp,
aws:MultiFactorAuthPresent, aws:PrincipalTag/*) are not supplied, so policies
that grant access only when such conditions are met may be reported as denied
even though the call would succeed in production.`,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE:         runPreflightAws,
	}
)

func init() {
	flags := preflightAwsCmd.Flags()
	flags.BoolVar(&preflightAwsState.agentless, "agentless", false,
		"check permissions for the Agentless integration")
	flags.BoolVar(&preflightAwsState.config, "config", false,
		"check permissions for the Config integration")
	flags.BoolVar(&preflightAwsState.cloudtrail, "cloudtrail", false,
		"check permissions for the CloudTrail integration")
	flags.BoolVar(&preflightAwsState.eksAuditLog, "eks-audit-log", false,
		"check permissions for the EKS Audit Log integration")
	flags.BoolVar(&preflightAwsState.isOrg, "is-org", false,
		"treat the account as an AWS Organizations management account")
	flags.BoolVar(&preflightAwsState.simulate, "simulate", false,
		"use IAM SimulatePrincipalPolicy (covers permissions boundaries and unconditional SCPs)")
	flags.StringVar(&preflightAwsState.region, "region", "",
		"AWS region to use for API calls")
	flags.StringVar(&preflightAwsState.profile, "profile", "",
		"AWS shared config profile to load credentials from")
	flags.StringVar(&preflightAwsState.accessKeyID, "access-key-id", "",
		"AWS access key ID (paired with --secret-access-key)")
	flags.StringVar(&preflightAwsState.secretAccessKey, "secret-access-key", "",
		"AWS secret access key (paired with --access-key-id)")
	flags.StringVar(&preflightAwsState.sessionToken, "session-token", "",
		"AWS session token for temporary credentials")
}

func runPreflightAws(_ *cobra.Command, _ []string) error {
	s := preflightAwsState
	if !s.agentless && !s.config && !s.cloudtrail && !s.eksAuditLog {
		return errors.New(
			"at least one of --agentless, --config, --cloudtrail, --eks-audit-log must be set",
		)
	}

	params := aws.Params{
		Agentless:       s.agentless,
		Config:          s.config,
		CloudTrail:      s.cloudtrail,
		EksAuditLog:     s.eksAuditLog,
		IsOrg:           s.isOrg,
		Simulate:        s.simulate,
		Region:          s.region,
		Profile:         s.profile,
		AccessKeyID:     s.accessKeyID,
		SecretAccessKey: s.secretAccessKey,
		SessionToken:    s.sessionToken,
	}

	pf, err := aws.New(params)
	if err != nil {
		return errors.Wrap(err, "unable to initialize AWS preflight")
	}
	if !cli.HumanOutput() || !cli.InteractiveMode() {
		pf.SetVerboseWriter(silentVerboseWriter())
	}

	cli.StartProgress("Running AWS preflight checks...")
	result, err := pf.Run()
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "AWS preflight failed")
	}

	if cli.JSONOutput() {
		if err := cli.OutputJSON(result); err != nil {
			return err
		}
		return preflightExitError(toStringErrorMap(result.Errors))
	}

	renderAwsHumanResult(result, integrationsRequestedAws(s))
	return preflightExitError(toStringErrorMap(result.Errors))
}

func renderAwsHumanResult(result *aws.Result, integrations []string) {
	cli.OutputHuman("Preflight check: AWS\n\n")
	cli.OutputHuman("Caller\n")
	cli.OutputHuman("  ARN:     %s\n", result.Caller.ARN)
	cli.OutputHuman("  Account: %s\n", result.Caller.AccountID)
	cli.OutputHuman("  User ID: %s\n", result.Caller.UserID)
	cli.OutputHuman("  Name:    %s\n", result.Caller.Name)
	cli.OutputHuman("  Admin:   %t\n", result.Caller.IsAdmin)

	if len(integrations) > 0 {
		cli.OutputHuman("\nIntegrations checked: %s\n", strings.Join(integrations, ", "))
	}

	renderIntegrationErrors(integrations, toStringErrorMap(result.Errors))

	cli.OutputHuman("\nDetails\n")
	if len(result.Details.Regions) > 0 {
		cli.OutputHuman("  Enabled regions: %s\n", strings.Join(result.Details.Regions, ", "))
	}
	if result.Details.OrgID != "" {
		cli.OutputHuman("  Organization ID:        %s\n", result.Details.OrgID)
		cli.OutputHuman("  Management account ID:  %s\n", result.Details.ManagementAccountID)
		cli.OutputHuman("  Caller is mgmt account: %t\n", result.Details.IsManagementAccount)
	}
	if result.Details.ExistingTrail.Name != "" {
		cli.OutputHuman("  Eligible CloudTrail:    %s\n", result.Details.ExistingTrail.Name)
	}
	if len(result.Details.EKSClusters) > 0 {
		cli.OutputHuman("  EKS clusters: %d\n", len(result.Details.EKSClusters))
	}
}

func integrationsRequestedAws(s struct {
	agentless       bool
	config          bool
	cloudtrail      bool
	eksAuditLog     bool
	isOrg           bool
	simulate        bool
	region          string
	profile         string
	accessKeyID     string
	secretAccessKey string
	sessionToken    string
}) []string {
	out := []string{}
	if s.agentless {
		out = append(out, string(aws.Agentless))
	}
	if s.config {
		out = append(out, string(aws.Config))
	}
	if s.cloudtrail {
		out = append(out, string(aws.CloudTrail))
	}
	if s.eksAuditLog {
		out = append(out, string(aws.EksAuditLog))
	}
	return out
}

// toStringErrorMap converts a provider's map[IntegrationType][]string into the
// generic map[string][]string the shared renderer expects. The IntegrationType
// type alias differs per provider, so we normalise once here.
func toStringErrorMap[K ~string](in map[K][]string) map[string][]string {
	out := make(map[string][]string, len(in))
	for k, v := range in {
		out[string(k)] = v
	}
	return out
}
