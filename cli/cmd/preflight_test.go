//
// Author:: Fortinet
// Copyright:: Copyright 2026, Fortinet
// License:: Apache License, Version 2.0
//

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/v2/lwpreflight/aws"
	"github.com/lacework/go-sdk/v2/lwpreflight/azure"
	"github.com/lacework/go-sdk/v2/lwpreflight/gcp"
)

func TestPreflightExitError(t *testing.T) {
	t.Run("nil when no issues", func(t *testing.T) {
		assert.NoError(t, preflightExitError(nil))
		assert.NoError(t, preflightExitError(map[string][]string{}))
		assert.NoError(t, preflightExitError(map[string][]string{"x": {}}))
	})

	t.Run("error when any issues", func(t *testing.T) {
		err := preflightExitError(map[string][]string{
			"aws_agentless":  {"a", "b"},
			"aws_cloudtrail": {"c"},
		})
		assert.EqualError(t, err, "preflight reported 3 issue(s)")
	})
}

func TestToStringErrorMapPerProvider(t *testing.T) {
	awsErrs := map[aws.IntegrationType][]string{
		aws.Agentless:  {"missing perm 1"},
		aws.CloudTrail: {"missing perm 2", "missing perm 3"},
	}
	out := toStringErrorMap(awsErrs)
	assert.Equal(t, []string{"missing perm 1"}, out["aws_agentless"])
	assert.Equal(t, []string{"missing perm 2", "missing perm 3"}, out["aws_cloudtrail"])

	azureErrs := map[azure.IntegrationType][]string{
		azure.ActivityLog: {"missing role"},
	}
	assert.Equal(t, []string{"missing role"}, toStringErrorMap(azureErrs)["azure_activity_log"])

	gcpErrs := map[gcp.IntegrationType][]string{
		gcp.GkeAuditLog: {"missing iam"},
	}
	assert.Equal(t, []string{"missing iam"}, toStringErrorMap(gcpErrs)["gcp_gke_audit_log"])
}

func TestIntegrationsRequestedAws(t *testing.T) {
	got := integrationsRequestedAws(struct {
		agentless       bool
		config          bool
		cloudtrail      bool
		eksAuditLog     bool
		isOrg           bool
		region          string
		profile         string
		accessKeyID     string
		secretAccessKey string
		sessionToken    string
	}{agentless: true, cloudtrail: true})
	assert.Equal(t, []string{"aws_agentless", "aws_cloudtrail"}, got)
}

func TestIntegrationsRequestedAzure(t *testing.T) {
	got := integrationsRequestedAzure(struct {
		agentless      bool
		config         bool
		activityLog    bool
		subscriptionID string
		tenantID       string
		clientID       string
		clientSecret   string
		region         string
	}{config: true, activityLog: true})
	assert.Equal(t, []string{"azure_config", "azure_activity_log"}, got)
}

func TestIntegrationsRequestedGcp(t *testing.T) {
	got := integrationsRequestedGcp(struct {
		agentless       bool
		auditLog        bool
		config          bool
		gkeAuditLog     bool
		region          string
		orgID           string
		projectID       string
		accessToken     string
		credentialsFile string
	}{auditLog: true, gkeAuditLog: true})
	assert.Equal(t, []string{"gcp_audit_log", "gcp_gke_audit_log"}, got)
}

func TestSilentVerboseWriter(t *testing.T) {
	w := silentVerboseWriter()
	w.Write("anything")
	w.Close()
}
