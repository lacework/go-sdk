//
// Author:: Fortinet
// Copyright:: Copyright 2026, Fortinet
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackUID(t *testing.T) {
	cases := []struct {
		name    string
		stackID string
		want    string
	}{
		{
			name:    "standard stack arn",
			stackID: "arn:aws:cloudformation:us-west-2:123456789012:stack/lacework-cfg/abcd1234-ef56-7890-abcd-ef1234567890",
			want:    "abcd1234",
		},
		{
			name: "eu region member stack",
			stackID: "arn:aws:cloudformation:eu-west-1:210987654321:stack/StackSet-lacework-x/" +
				"deadbeef-0000-1111-2222-333344445555",
			want: "deadbeef",
		},
		{
			name:    "empty",
			stackID: "",
			want:    "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, stackUID(c.stackID))
		})
	}
}

// The role name and external id must key off the same UID, exactly as the config_org member
// template builds them, or the backend AssumeRole validation rejects the re-registration.
func TestRepairDerivationMatchesTemplate(t *testing.T) {
	stackID := "arn:aws:cloudformation:us-east-1:123456789012:stack/StackSet-lw/abcd1234-ef56-7890-abcd-ef1234567890"
	uid := stackUID(stackID)

	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/lacework-config-role-%s", "123456789012", uid)
	externalID := fmt.Sprintf("lweid:aws:v2:%s:%s:LW%s", "TENANT", "123456789012", uid)

	assert.Equal(t, "arn:aws:iam::123456789012:role/lacework-config-role-abcd1234", roleArn)
	assert.Equal(t, "lweid:aws:v2:TENANT:123456789012:LWabcd1234", externalID)
}

// The --json payload for the register path must carry the derived role/external id, and only the
// applied outcome should include the server-assigned guid (dry-run has none).
func TestRepairRegisterResult(t *testing.T) {
	roleArn := "arn:aws:iam::123456789012:role/lacework-config-role-abcd1234"
	externalID := "lweid:aws:v2:TENANT:123456789012:LWabcd1234"

	t.Run("dry-run has no guid and dryRun true", func(t *testing.T) {
		r := repairRegisterResult("register", "123456789012", "TENANT-Config", roleArn, externalID, "")
		assert.Equal(t, "123456789012", r["accountId"])
		assert.Equal(t, "register", r["action"])
		assert.Equal(t, "TENANT-Config", r["name"])
		assert.Equal(t, roleArn, r["roleArn"])
		assert.Equal(t, externalID, r["externalId"])
		assert.Equal(t, true, r["dryRun"])
		_, hasGuid := r["intgGuid"]
		assert.False(t, hasGuid, "dry-run must not report a guid")
	})

	t.Run("applied carries guid and dryRun false", func(t *testing.T) {
		r := repairRegisterResult("registered", "123456789012", "TENANT-Config", roleArn, externalID, "LWINTG-GUID")
		assert.Equal(t, "registered", r["action"])
		assert.Equal(t, "LWINTG-GUID", r["intgGuid"])
		_, hasDryRun := r["dryRun"]
		assert.False(t, hasDryRun, "applied outcome must not report dryRun")
	})
}

// --all classifies each expected account into exactly one bucket: no stack instance means rebuild
// then register; an instance without an integration means register only; fully onboarded means skip.
func TestDiffRepairTargets(t *testing.T) {
	expected := []string{"111111111111", "222222222222", "333333333333", "444444444444"}
	instances := map[string]string{
		"111111111111": "arn:aws:cloudformation:us-east-1:111111111111:stack/s/aaaa1111-x",
		"222222222222": "arn:aws:cloudformation:us-east-1:222222222222:stack/s/bbbb2222-x",
		// 333333333333 and 444444444444 have no instance
	}
	registered := map[string]bool{
		"111111111111": true, // healthy: instance + integration
		// 222222222222 has an instance but no integration
	}

	missingInstance, missingIntegration := diffRepairTargets(expected, instances, registered)

	assert.Equal(t, []string{"333333333333", "444444444444"}, missingInstance)
	assert.Equal(t, []string{"222222222222"}, missingIntegration)

	// healthy fleet: both buckets empty but non-nil, so the --json envelope renders [] not null
	missingInstance, missingIntegration = diffRepairTargets([]string{"111111111111"}, instances, registered)
	assert.NotNil(t, missingInstance)
	assert.NotNil(t, missingIntegration)
	assert.Empty(t, missingInstance)
	assert.Empty(t, missingIntegration)
}
