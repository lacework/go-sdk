//
// Author:: Fortinet
// Copyright:: Copyright 2026, Fortinet
// License:: Apache License, Version 2.0
//

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/v2/lwpreflight/gcp"
)

var (
	preflightGcpState struct {
		agentless       bool
		auditLog        bool
		config          bool
		gkeAuditLog     bool
		region          string
		orgID           string
		projectID       string
		accessToken     string
		credentialsFile string
	}

	preflightGcpCmd = &cobra.Command{
		Use:   "gcp",
		Short: "Run preflight checks against a GCP project",
		Long: `Run preflight checks against a GCP project to verify the caller has the IAM
permissions required by the selected Lacework integrations. Credentials are
resolved from --access-token, --credentials-file, or the
GOOGLE_APPLICATION_CREDENTIALS environment variable.

At least one integration flag must be set: --agentless, --audit-log,
--config, or --gke-audit-log.`,
		Args: cobra.NoArgs,
		RunE: runPreflightGcp,
	}
)

func init() {
	flags := preflightGcpCmd.Flags()
	flags.BoolVar(&preflightGcpState.agentless, "agentless", false,
		"check permissions for the Agentless integration")
	flags.BoolVar(&preflightGcpState.auditLog, "audit-log", false,
		"check permissions for the Audit Log integration")
	flags.BoolVar(&preflightGcpState.config, "config", false,
		"check permissions for the Config integration")
	flags.BoolVar(&preflightGcpState.gkeAuditLog, "gke-audit-log", false,
		"check permissions for the GKE Audit Log integration")
	flags.StringVar(&preflightGcpState.region, "region", "",
		"GCP region to use for region-scoped checks")
	flags.StringVar(&preflightGcpState.orgID, "org-id", "",
		"GCP organization ID; sets the integration to org-level when non-empty")
	flags.StringVar(&preflightGcpState.projectID, "project-id", "",
		"GCP project ID (required)")
	flags.StringVar(&preflightGcpState.accessToken, "access-token", "",
		"GCP OAuth2 access token")
	flags.StringVar(&preflightGcpState.credentialsFile, "credentials-file", "",
		"Path to a GCP service account credentials JSON file")
}

func runPreflightGcp(_ *cobra.Command, _ []string) error {
	s := preflightGcpState
	if !s.agentless && !s.auditLog && !s.config && !s.gkeAuditLog {
		return errors.New(
			"at least one of --agentless, --audit-log, --config, --gke-audit-log must be set",
		)
	}
	if s.projectID == "" {
		return errors.New("--project-id is required")
	}

	params := gcp.Params{
		Agentless:       s.agentless,
		AuditLog:        s.auditLog,
		Config:          s.config,
		GkeAuditLog:     s.gkeAuditLog,
		Region:          s.region,
		OrgID:           s.orgID,
		ProjectID:       s.projectID,
		AccessToken:     s.accessToken,
		CredentialsFile: s.credentialsFile,
	}

	pf, err := gcp.New(params)
	if err != nil {
		return errors.Wrap(err, "unable to initialize GCP preflight")
	}
	if !cli.HumanOutput() {
		pf.SetVerboseWriter(silentVerboseWriter())
	}

	cli.StartProgress("Running GCP preflight checks...")
	result, err := pf.Run()
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "GCP preflight failed")
	}

	if cli.JSONOutput() {
		if err := cli.OutputJSON(result); err != nil {
			return err
		}
		return preflightExitError(toStringErrorMap(result.Errors))
	}

	renderGcpHumanResult(result, integrationsRequestedGcp(s))
	return preflightExitError(toStringErrorMap(result.Errors))
}

func renderGcpHumanResult(result *gcp.Result, integrations []string) {
	cli.OutputHuman("Preflight check: GCP\n\n")
	cli.OutputHuman("Caller\n")
	cli.OutputHuman("  Email:   %s\n", result.Caller.Email)
	cli.OutputHuman("  User ID: %s\n", result.Caller.UserID)

	if len(integrations) > 0 {
		cli.OutputHuman("\nIntegrations checked: %s\n", joinIntegrations(integrations))
	}

	renderIntegrationErrors(integrations, toStringErrorMap(result.Errors))

	if len(result.Details.SchedulerRegions) > 0 {
		cli.OutputHuman("\nDetails\n")
		cli.OutputHuman("  Cloud Scheduler regions: %d\n", len(result.Details.SchedulerRegions))
	}
}

func integrationsRequestedGcp(s struct {
	agentless       bool
	auditLog        bool
	config          bool
	gkeAuditLog     bool
	region          string
	orgID           string
	projectID       string
	accessToken     string
	credentialsFile string
}) []string {
	out := []string{}
	if s.agentless {
		out = append(out, string(gcp.Agentless))
	}
	if s.auditLog {
		out = append(out, string(gcp.AuditLog))
	}
	if s.config {
		out = append(out, string(gcp.Config))
	}
	if s.gkeAuditLog {
		out = append(out, string(gcp.GkeAuditLog))
	}
	return out
}
