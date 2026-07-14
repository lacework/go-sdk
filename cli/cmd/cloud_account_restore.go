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
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/v2/api"
)

// restore command flags
var (
	restoreFile   string // -f/--file (input)
	restoreDryRun bool   // --dry-run
)

func init() {
	cloudAccountCommand.AddCommand(cloudAccountRestoreCmd)

	cloudAccountRestoreCmd.Flags().StringVarP(&restoreFile,
		"file", "f", "", "backup file to restore integrations from")
	cloudAccountRestoreCmd.Flags().BoolVar(&restoreDryRun,
		"dry-run", false, "show what would be re-created without creating")
}

var cloudAccountRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Re-create cloud account integrations from a backup file",
	Long: `Re-create the cloud account integrations captured in a backup file, undoing a
previous 'cloud-account delete --bulk'. Each integration is re-created from its stored
name, type and data.

    lacework cloud-account restore --file backup.json

Note 1: re-creating an integration produces a NEW integration GUID; a table of old -> new GUIDs
is printed at the end.

Note 2: some integration secrets (such as an AWS externalId) are returned masked by the API and
therefore stored masked in the backup. Re-creating those needs the original secret value.`,
	Args: cobra.NoArgs,
	RunE: cloudAccountRestore,
}

func cloudAccountRestore(_ *cobra.Command, _ []string) error {
	if restoreFile == "" {
		return errors.New("--file is required")
	}

	backup, err := readCloudAccountBackup(restoreFile)
	if err != nil {
		return err
	}

	cli.OutputHuman("Found %d integration(s) to re-create in %s:\n", backup.Count, restoreFile)
	var maskedWarn int
	for _, intg := range backup.Integrations {
		if hasCreds, masked := integrationExternalIDMasked(intg); hasCreds && masked {
			maskedWarn++
		}
		cli.OutputHuman("  %s  %s  (%s)\n", intg.IntgGuid, intg.Name, intg.Type)
	}
	if maskedWarn > 0 {
		cli.OutputHuman(
			"WARNING: %d record(s) have a masked/empty secret (e.g. externalId); re-creating them "+
				"may fail or produce a broken integration unless the original secret is supplied.\n", maskedWarn)
	}

	if restoreDryRun {
		cli.OutputHuman("\nDry-run: no integrations were re-created. Re-run without --dry-run to restore.\n")
		return nil
	}

	if !confirmBulkOperation(fmt.Sprintf("Re-create %d cloud account integration(s)?", backup.Count)) {
		cli.OutputHuman("Aborted. No integrations were re-created.\n")
		return nil
	}

	var created, failed int
	var mappings [][]string // old guid -> new guid (+ name)
	for _, intg := range backup.Integrations {
		caType, found := api.FindCloudAccountType(intg.Type)
		if !found {
			cli.OutputHuman("Failed to re-create %s: unknown type %q\n", intg.Name, intg.Type)
			failed++
			continue
		}

		account := api.NewCloudAccount(intg.Name, caType, intg.Data)
		account.Enabled = intg.Enabled

		cli.StartProgress(fmt.Sprintf(" Re-creating %s...", intg.Name))
		resp, err := cli.LwApi.V2.CloudAccounts.Create(account)
		cli.StopProgress()
		if err != nil {
			cli.OutputHuman("Failed to re-create %s (%s): %s\n", intg.Name, intg.Type, err.Error())
			failed++
			continue
		}
		mappings = append(mappings, []string{intg.IntgGuid, resp.Data.IntgGuid, intg.Name})
		created++
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(restoreMappingsJSON(mappings, created, failed))
	}

	if len(mappings) > 0 {
		cli.OutputHuman("\nRe-created integrations (the integration GUID changes on re-create):\n")
		cli.OutputHuman(renderSimpleTable(
			[]string{"Old GUID", "New GUID", "Name"}, mappings))
	}
	cli.OutputHuman("\nRe-created %d integration(s); %d failure(s).\n", created, failed)
	if failed > 0 {
		return errors.Errorf("%d integration(s) could not be re-created", failed)
	}
	return nil
}

func restoreMappingsJSON(mappings [][]string, created, failed int) map[string]interface{} {
	remapped := make([]map[string]string, 0, len(mappings))
	for _, m := range mappings {
		remapped = append(remapped, map[string]string{"oldGuid": m[0], "newGuid": m[1], "name": m[2]})
	}
	return map[string]interface{}{
		"created":  created,
		"failed":   failed,
		"mappings": remapped,
	}
}
