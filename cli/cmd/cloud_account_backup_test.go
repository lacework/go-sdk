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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lacework/go-sdk/v2/api"
)

func awsCfgRaw(guid, name, roleArn, externalID string, isOrg int) api.CloudAccountRaw {
	raw := api.NewCloudAccount(name, api.AwsCfgCloudAccount, map[string]interface{}{
		"crossAccountCredentials": map[string]interface{}{
			"roleArn":    roleArn,
			"externalId": externalID,
		},
	})
	raw.IntgGuid = guid
	raw.IsOrg = isOrg
	return raw
}

func TestParseRoleArn(t *testing.T) {
	cases := []struct {
		name        string
		arn         string
		wantAccount string
		wantRole    string
		wantErr     bool
	}{
		{"standard", "arn:aws:iam::123456789012:role/acme-laceworkcwsrole-sa", "123456789012", "acme-laceworkcwsrole-sa", false},
		{"role with path", "arn:aws:iam::210987654321:role/lw/config-role", "210987654321", "lw/config-role", false},
		{"gov partition", "arn:aws-us-gov:iam::111122223333:role/r", "111122223333", "r", false},
		{"not a role", "arn:aws:iam::123456789012:user/bob", "", "", true},
		{"garbage", "not-an-arn", "", "", true},
		{"empty", "", "", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			account, role, err := parseRoleArn(c.arn)
			if c.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, c.wantAccount, account)
			assert.Equal(t, c.wantRole, role)
		})
	}
}

func TestDeriveAwsAccountAndRole(t *testing.T) {
	account, role, ok := deriveAwsAccountAndRole(
		awsCfgRaw("g1", "n1", "arn:aws:iam::123456789012:role/some-role", "ext", 0))
	assert.True(t, ok)
	assert.Equal(t, "123456789012", account)
	assert.Equal(t, "some-role", role)

	_, _, ok = deriveAwsAccountAndRole(awsCfgRaw("g2", "n2", "bad-arn", "ext", 0))
	assert.False(t, ok)

	// A record without cross-account credentials (e.g. a non-AWS type) yields no target.
	raw := api.NewCloudAccount("gcp", api.GcpCfgCloudAccount, map[string]interface{}{"id": "proj"})
	_, _, ok = deriveAwsAccountAndRole(raw)
	assert.False(t, ok)
}

func TestIntegrationExternalIDMasked(t *testing.T) {
	hasCreds, masked := integrationExternalIDMasked(
		awsCfgRaw("g", "n", "arn:aws:iam::123456789012:role/r", "lweid:aws:v2:acct:123:uid", 0))
	assert.True(t, hasCreds)
	assert.False(t, masked)

	hasCreds, masked = integrationExternalIDMasked(
		awsCfgRaw("g", "n", "arn:aws:iam::123456789012:role/r", "****", 0))
	assert.True(t, hasCreds)
	assert.True(t, masked)

	hasCreds, _ = integrationExternalIDMasked(
		api.NewCloudAccount("gcp", api.GcpCfgCloudAccount, map[string]interface{}{"id": "proj"}))
	assert.False(t, hasCreds)
}

func TestIsAwsCrossAccountType(t *testing.T) {
	assert.True(t, isAwsCrossAccountType("AwsCfg"))
	assert.True(t, isAwsCrossAccountType("AwsCtSqs"))
	assert.True(t, isAwsCrossAccountType("AwsEksAudit"))
	assert.False(t, isAwsCrossAccountType("GcpCfg"))
	assert.False(t, isAwsCrossAccountType("AzureCfg"))
	assert.False(t, isAwsCrossAccountType(""))
}

func TestIsCustomerManagedPolicy(t *testing.T) {
	assert.True(t, isCustomerManagedPolicy("arn:aws:iam::123456789012:policy/LaceworkCWSAuditPolicy-123456789012"))
	assert.False(t, isCustomerManagedPolicy("arn:aws:iam::aws:policy/SecurityAudit"))
}

func TestBackupRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "backup.json")
	original := cloudAccountBackup{
		LWAccount: "acme",
		Type:      "AwsCfg",
		Count:     2,
		Integrations: []api.CloudAccountRaw{
			awsCfgRaw("g1", "int-1", "arn:aws:iam::111111111111:role/r1", "ext1", 0),
			awsCfgRaw("g2", "int-2", "arn:aws:iam::222222222222:role/r2", "ext2", 0),
		},
	}
	require.NoError(t, writeCloudAccountBackup(path, original))

	got, err := readCloudAccountBackup(path)
	require.NoError(t, err)
	assert.Equal(t, "acme", got.LWAccount)
	assert.Equal(t, "AwsCfg", got.Type)
	require.Len(t, got.Integrations, 2)
	assert.Equal(t, "g1", got.Integrations[0].IntgGuid)

	account, role, ok := deriveAwsAccountAndRole(got.Integrations[1])
	assert.True(t, ok)
	assert.Equal(t, "222222222222", account)
	assert.Equal(t, "r2", role)
}

func TestReadCloudAccountBackupErrors(t *testing.T) {
	_, err := readCloudAccountBackup(filepath.Join(t.TempDir(), "missing.json"))
	assert.Error(t, err)

	empty := filepath.Join(t.TempDir(), "empty.json")
	require.NoError(t, writeCloudAccountBackup(empty, cloudAccountBackup{Type: "AwsCfg"}))
	_, err = readCloudAccountBackup(empty)
	assert.Error(t, err, "backup with no integrations should error")
}

func TestCollectCleanupTargets(t *testing.T) {
	backup := cloudAccountBackup{
		Integrations: []api.CloudAccountRaw{
			awsCfgRaw("g1", "a", "arn:aws:iam::111111111111:role/r1", "e", 0),
			// duplicate (account, role) -> deduped
			awsCfgRaw("g2", "b", "arn:aws:iam::111111111111:role/r1", "e", 0),
			awsCfgRaw("g3", "c", "arn:aws:iam::222222222222:role/r2", "e", 0),
			// unparseable ARN
			awsCfgRaw("g4", "d", "bad", "e", 0),
			// non-AWS type -> unsupported
			api.NewCloudAccount("gcp", api.GcpCfgCloudAccount, map[string]interface{}{"id": "p"}),
		},
	}

	targets, unsupported, unparseable := collectCleanupTargets(backup, "")
	assert.Len(t, targets, 2, "duplicate (account, role) should be deduped")
	assert.Equal(t, 1, unparseable)
	assert.Equal(t, 1, unsupported["GcpCfg"])

	// Type filter keeps only AwsCfg records.
	targets, _, _ = collectCleanupTargets(backup, "AwsCfg")
	assert.Len(t, targets, 2)
}
