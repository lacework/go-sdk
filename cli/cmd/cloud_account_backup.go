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
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/v2/api"
)

// backup command flags
var (
	backupType           string   // -t/--type
	backupAccountIDs     []string // --account-ids
	backupAccountIDsFile string   // --account-ids-file
	backupFile           string   // -f/--file (output)
)

func init() {
	cloudAccountCommand.AddCommand(cloudAccountBackupCmd)

	cloudAccountBackupCmd.Flags().StringVarP(&backupType,
		"type", "t", "", "cloud account type to back up (e.g. AwsCfg)")
	cloudAccountBackupCmd.Flags().StringSliceVar(&backupAccountIDs,
		"account-ids", nil, "comma-separated AWS account ids to back up (AWS types only)")
	cloudAccountBackupCmd.Flags().StringVar(&backupAccountIDsFile,
		"account-ids-file", "", "file with AWS account ids, one per line (AWS types only)")
	cloudAccountBackupCmd.Flags().StringVarP(&backupFile,
		"file", "f", "", "output backup file (default lacework-cloud-account-backup-<ts>.json)")
}

// cloudAccountBackup is the on-disk format written by `backup` and consumed by the other bulk
// commands. It stores full raw integration records so any cloud-account type can be restored
// faithfully by `restore`, and AWS role ARNs remain available to `cleanup`.
type cloudAccountBackup struct {
	CreatedAt    time.Time             `json:"createdAt"`
	LWAccount    string                `json:"lwAccount"`
	LWSubaccount string                `json:"lwSubaccount,omitempty"`
	Type         string                `json:"type"`
	Count        int                   `json:"count"`
	Integrations []api.CloudAccountRaw `json:"integrations"`
}

var cloudAccountBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Back up cloud account integrations of a given type to a file",
	Long: `Back up cloud account integrations of a given type into a local JSON file.

The backup file is the input for the ` + "`delete --bulk`" + `, ` + "`restore`" + ` and ` + "`cleanup`" + ` commands,
which together let you retire integrations in bulk and restore or clean them up afterwards.

By default every integration of the given type is backed up. For AWS types you can narrow the
scope to specific AWS accounts with --account-ids or --account-ids-file (the account id is parsed
from each integration's role ARN).

    lacework cloud-account backup --type AwsCfg
    lacework cloud-account backup --type AwsCfg --account-ids 111111111111,222222222222
    lacework cloud-account backup --type AwsCfg --account-ids-file accounts.txt --file backup.json`,
	Args: cobra.NoArgs,
	RunE: runCloudAccountBackup,
}

func runCloudAccountBackup(_ *cobra.Command, _ []string) error {
	caType, found := api.FindCloudAccountType(backupType)
	if !found {
		return errors.Errorf("unknown cloud account type '%s'", backupType)
	}

	accountScope, err := loadAccountScope()
	if err != nil {
		return err
	}
	if len(accountScope) > 0 && !isAwsCrossAccountType(caType.String()) {
		return errors.Errorf(
			"account-id filtering is not supported for type '%s'; omit --account-ids to back up all",
			caType.String(),
		)
	}

	cli.StartProgress("Fetching cloud accounts...")
	res, err := cli.LwApi.V2.CloudAccounts.ListByType(caType)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to list cloud accounts")
	}

	var (
		selected   []api.CloudAccountRaw
		skipped    int
		maskedExtI int
	)
	for _, raw := range res.Data {
		if len(accountScope) > 0 {
			accountID, _, ok := deriveAwsAccountAndRole(raw)
			if !ok {
				cli.Log.Debugw("skipping integration with unparseable role ARN",
					"guid", raw.IntgGuid, "name", raw.Name)
				skipped++
				continue
			}
			if _, want := accountScope[accountID]; !want {
				continue
			}
		}
		if hasCreds, masked := integrationExternalIDMasked(raw); hasCreds && masked {
			maskedExtI++
		}
		selected = append(selected, raw)
	}

	if skipped > 0 {
		cli.OutputHuman("Skipped %d integration(s) with an unparseable role ARN.\n", skipped)
	}
	if len(selected) == 0 {
		cli.OutputHuman("No matching cloud accounts found. Nothing was backed up.\n")
		return nil
	}

	backup := cloudAccountBackup{
		CreatedAt:    time.Now().UTC(),
		LWAccount:    cli.Account,
		LWSubaccount: cli.Subaccount,
		Type:         caType.String(),
		Count:        len(selected),
		Integrations: selected,
	}

	path, err := resolveBackupPath()
	if err != nil {
		return err
	}
	if err := writeCloudAccountBackup(path, backup); err != nil {
		return err
	}

	if maskedExtI > 0 {
		cli.OutputHuman(
			"WARNING: %d integration(s) returned a masked/empty externalId. Restoring them with "+
				"`restore` will need the externalId you already hold.\n", maskedExtI)
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(map[string]interface{}{"file": path, "count": len(selected)})
	}
	cli.OutputHuman("Backed up %d %s integration(s) to %s\n", len(selected), caType.String(), path)
	return nil
}

// loadAccountScope merges --account-ids and --account-ids-file into a set (empty => no filter).
func loadAccountScope() (map[string]struct{}, error) {
	scope := map[string]struct{}{}
	for _, id := range backupAccountIDs {
		if id = strings.TrimSpace(id); id != "" {
			scope[id] = struct{}{}
		}
	}
	if backupAccountIDsFile != "" {
		f, err := os.Open(backupAccountIDsFile)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read account ids file")
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(strings.SplitN(scanner.Text(), "#", 2)[0])
			if line != "" {
				scope[line] = struct{}{}
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, errors.Wrap(err, "unable to read account ids file")
		}
	}
	return scope, nil
}

func defaultBackupFilename() string {
	return fmt.Sprintf("lacework-cloud-account-backup-%s.json", time.Now().UTC().Format("20060102T150405Z"))
}

// resolveBackupPath decides where to write the backup. An explicit --file always wins. Otherwise the
// filename is generated automatically and, in interactive mode, the user is prompted only for the
// directory (defaulting to the current directory); non-interactive runs use the current directory.
func resolveBackupPath() (string, error) {
	if backupFile != "" {
		return backupFile, nil
	}
	dir := "."
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: "Directory to save the backup file:", Default: dir},
		Response: &dir,
	}); err != nil {
		return "", err
	}
	return filepath.Join(dir, defaultBackupFilename()), nil
}

func writeCloudAccountBackup(path string, backup cloudAccountBackup) error {
	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return errors.Wrap(err, "unable to encode backup")
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return errors.Wrap(err, "unable to write backup file")
	}
	return nil
}

func readCloudAccountBackup(path string) (cloudAccountBackup, error) {
	var backup cloudAccountBackup
	data, err := os.ReadFile(path)
	if err != nil {
		return backup, errors.Wrap(err, "unable to read backup file")
	}
	if err := json.Unmarshal(data, &backup); err != nil {
		return backup, errors.Wrap(err, "unable to parse backup file")
	}
	if len(backup.Integrations) == 0 {
		return backup, errors.Errorf("backup file %s contains no integrations", path)
	}
	return backup, nil
}

// parseRoleArn extracts the AWS account id and role name from a role ARN.
func parseRoleArn(roleArn string) (accountID, roleName string, err error) {
	parsed, err := arn.Parse(roleArn)
	if err != nil {
		return "", "", err
	}
	if !strings.HasPrefix(parsed.Resource, "role/") {
		return "", "", errors.Errorf("ARN %q is not an IAM role", roleArn)
	}
	return parsed.AccountID, strings.TrimPrefix(parsed.Resource, "role/"), nil
}

// deriveAwsAccountAndRole pulls the role ARN out of a raw integration's data and parses it.
func deriveAwsAccountAndRole(raw api.CloudAccountRaw) (accountID, roleName string, ok bool) {
	roleArn, found := awsRoleArnFromRaw(raw)
	if !found {
		return "", "", false
	}
	accountID, roleName, err := parseRoleArn(roleArn)
	if err != nil {
		return "", "", false
	}
	return accountID, roleName, true
}

// awsRoleArnFromRaw reads data.crossAccountCredentials.roleArn from a raw integration record.
func awsRoleArnFromRaw(raw api.CloudAccountRaw) (string, bool) {
	cred, ok := crossAccountCredentials(raw)
	if !ok {
		return "", false
	}
	roleArn, ok := cred["roleArn"].(string)
	return roleArn, ok && roleArn != ""
}

// integrationExternalIDMasked reports whether the record has cross-account credentials and, if so,
// whether its externalId looks masked/empty (which matters when restoring via `restore`).
func integrationExternalIDMasked(raw api.CloudAccountRaw) (hasCreds bool, masked bool) {
	cred, ok := crossAccountCredentials(raw)
	if !ok {
		return false, false
	}
	extID, _ := cred["externalId"].(string)
	return true, strings.Trim(extID, "*Xx ") == ""
}

func crossAccountCredentials(raw api.CloudAccountRaw) (map[string]interface{}, bool) {
	m, ok := raw.Data.(map[string]interface{})
	if !ok {
		return nil, false
	}
	cred, ok := m["crossAccountCredentials"].(map[string]interface{})
	return cred, ok
}

// isAwsCrossAccountType reports whether a cloud-account type uses an AWS cross-account role ARN,
// which is what enables account-id filtering (backup) and IAM cleanup.
func isAwsCrossAccountType(typeStr string) bool {
	switch typeStr {
	case api.AwsCfgCloudAccount.String(),
		api.AwsCtSqsCloudAccount.String(),
		api.AwsEksAuditCloudAccount.String(),
		api.AwsSidekickCloudAccount.String(),
		api.AwsSidekickOrgCloudAccount.String(),
		api.AwsUsGovCfgCloudAccount.String(),
		api.AwsUsGovCtSqsCloudAccount.String():
		return true
	default:
		return false
	}
}

// confirmBulkOperation asks the operator to confirm a destructive bulk action. In non-interactive
// mode it returns true (the invocation is already explicit: a file and flags were supplied).
func confirmBulkOperation(message string) bool {
	if cli.nonInteractive {
		return true
	}
	proceed := false
	if err := survey.AskOne(&survey.Confirm{Message: message, Default: false}, &proceed); err != nil {
		return false
	}
	return proceed
}
