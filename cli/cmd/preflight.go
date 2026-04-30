//
// Author:: Fortinet
// Copyright:: Copyright 2026, Fortinet
// License:: Apache License, Version 2.0
//

package cmd

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/v2/lwpreflight/verbosewriter"
)

var (
	preflightCmd = &cobra.Command{
		Use:   "preflight",
		Short: "Run preflight checks against a cloud account",
		Long: `Run preflight checks against AWS, Azure, or GCP accounts to verify that the
caller identity has the IAM permissions and access required to onboard the
selected Lacework integrations. Any missing permissions are reported as
errors keyed by integration type.`,
	}
)

func init() {
	rootCmd.AddCommand(preflightCmd)
	preflightCmd.AddCommand(preflightAwsCmd)
	preflightCmd.AddCommand(preflightAzureCmd)
	preflightCmd.AddCommand(preflightGcpCmd)
}

// discardVerboseWriter satisfies verbosewriter.WriteCloser without writing
// anything, so progress chatter does not corrupt machine-readable output.
type discardVerboseWriter struct{}

func (discardVerboseWriter) Write(string) {}
func (discardVerboseWriter) Close()       {}

// silentVerboseWriter returns a verbose writer that drops all output. We use
// this whenever the user has opted out of human output (e.g. --json), so the
// preflight package does not write progress lines to the terminal.
func silentVerboseWriter() verbosewriter.WriteCloser {
	return discardVerboseWriter{}
}

// renderIntegrationErrors prints a per-integration pass/fail summary derived
// from the Errors map returned by every preflight provider. integrations is
// the list of integrations the caller actually requested, so we can surface
// "OK" for ones that ran without issues. errs is the raw map keyed by the
// provider's IntegrationType (a string alias).
func renderIntegrationErrors(integrations []string, errs map[string][]string) {
	if len(integrations) == 0 {
		cli.OutputHuman("No integrations selected.\n")
		return
	}

	sort.Strings(integrations)

	cli.OutputHuman("\nResults\n")
	for _, integration := range integrations {
		issues := errs[integration]
		if len(issues) == 0 {
			cli.OutputHuman("  %s: OK\n", integration)
			continue
		}
		cli.OutputHuman("  %s: FAIL (%d issue(s))\n", integration, len(issues))
		for _, issue := range issues {
			cli.OutputHuman("    - %s\n", issue)
		}
	}
}

// preflightExitError returns a non-nil error when any integration reported
// problems, so the CLI exits non-zero and scripts can branch on it. The
// rendered output (human or JSON) has already been emitted by the caller.
func preflightExitError(errs map[string][]string) error {
	total := 0
	for _, v := range errs {
		total += len(v)
	}
	if total == 0 {
		return nil
	}
	return errors.New(fmt.Sprintf("preflight reported %d issue(s)", total))
}
