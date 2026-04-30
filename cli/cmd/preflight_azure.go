//
// Author:: Lacework Inc.
// Copyright:: Copyright 2026, Lacework Inc.
// License:: Apache License, Version 2.0
//

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/v2/lwpreflight/azure"
)

var (
	preflightAzureState struct {
		agentless      bool
		config         bool
		activityLog    bool
		subscriptionID string
		tenantID       string
		clientID       string
		clientSecret   string
		region         string
	}

	preflightAzureCmd = &cobra.Command{
		Use:   "azure",
		Short: "Run preflight checks against an Azure subscription",
		Long: `Run preflight checks against an Azure subscription to verify the caller has
the role assignments required by the selected Lacework integrations.
Credentials are resolved using DefaultAzureCredential unless --client-id and
--client-secret are provided for a service principal.

At least one integration flag must be set: --agentless, --config, or
--activity-log.`,
		Args: cobra.NoArgs,
		RunE: runPreflightAzure,
	}
)

func init() {
	flags := preflightAzureCmd.Flags()
	flags.BoolVar(&preflightAzureState.agentless, "agentless", false,
		"check permissions for the Agentless integration")
	flags.BoolVar(&preflightAzureState.config, "config", false,
		"check permissions for the Config integration")
	flags.BoolVar(&preflightAzureState.activityLog, "activity-log", false,
		"check permissions for the Activity Log integration")
	flags.StringVar(&preflightAzureState.subscriptionID, "subscription-id", "",
		"Azure subscription ID (required)")
	flags.StringVar(&preflightAzureState.tenantID, "tenant-id", "",
		"Azure tenant ID (required when using --client-id/--client-secret)")
	flags.StringVar(&preflightAzureState.clientID, "client-id", "",
		"Azure service principal client ID")
	flags.StringVar(&preflightAzureState.clientSecret, "client-secret", "",
		"Azure service principal client secret")
	flags.StringVar(&preflightAzureState.region, "region", "",
		"Azure region to use for region-scoped checks")
}

func runPreflightAzure(_ *cobra.Command, _ []string) error {
	s := preflightAzureState
	if !s.agentless && !s.config && !s.activityLog {
		return errors.New(
			"at least one of --agentless, --config, --activity-log must be set",
		)
	}
	if s.subscriptionID == "" {
		return errors.New("--subscription-id is required")
	}

	params := azure.Params{
		Agentless:      s.agentless,
		Config:         s.config,
		ActivityLog:    s.activityLog,
		SubscriptionID: s.subscriptionID,
		TenantID:       s.tenantID,
		ClientID:       s.clientID,
		ClientSecret:   s.clientSecret,
		Region:         s.region,
	}

	pf, err := azure.New(params)
	if err != nil {
		return errors.Wrap(err, "unable to initialize Azure preflight")
	}
	if !cli.HumanOutput() {
		pf.SetVerboseWriter(silentVerboseWriter())
	}

	cli.StartProgress("Running Azure preflight checks...")
	result, err := pf.Run()
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "Azure preflight failed")
	}

	if cli.JSONOutput() {
		if err := cli.OutputJSON(result); err != nil {
			return err
		}
		return preflightExitError(toStringErrorMap(result.Errors))
	}

	renderAzureHumanResult(result, integrationsRequestedAzure(s))
	return preflightExitError(toStringErrorMap(result.Errors))
}

func renderAzureHumanResult(result *azure.Result, integrations []string) {
	cli.OutputHuman("Preflight check: Azure\n\n")
	cli.OutputHuman("Caller\n")
	cli.OutputHuman("  Object ID:    %s\n", result.Caller.ObjectID)
	cli.OutputHuman("  Display name: %s\n", result.Caller.DisplayName)
	cli.OutputHuman("  Tenant ID:    %s\n", result.Caller.TenantID)
	cli.OutputHuman("  Admin:        %t\n", result.Caller.IsAdmin)

	if len(integrations) > 0 {
		cli.OutputHuman("\nIntegrations checked: %s\n", joinIntegrations(integrations))
	}

	renderIntegrationErrors(integrations, toStringErrorMap(result.Errors))

	if len(result.Details.Regions) > 0 {
		cli.OutputHuman("\nDetails\n")
		cli.OutputHuman("  Available regions: %d\n", len(result.Details.Regions))
	}
}

func integrationsRequestedAzure(s struct {
	agentless      bool
	config         bool
	activityLog    bool
	subscriptionID string
	tenantID       string
	clientID       string
	clientSecret   string
	region         string
}) []string {
	out := []string{}
	if s.agentless {
		out = append(out, string(azure.Agentless))
	}
	if s.config {
		out = append(out, string(azure.Config))
	}
	if s.activityLog {
		out = append(out, string(azure.ActivityLog))
	}
	return out
}
