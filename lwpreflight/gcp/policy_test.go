package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	cloudresourcemanagerV3 "google.golang.org/api/cloudresourcemanager/v3"
)

func policy(bindings ...*cloudresourcemanagerV3.Binding) []*cloudresourcemanagerV3.Policy {
	return []*cloudresourcemanagerV3.Policy{{Bindings: bindings}}
}

func TestRolesForCaller(t *testing.T) {
	cases := []struct {
		name     string
		policies []*cloudresourcemanagerV3.Policy
		email    string
		expected []string
	}{
		{
			// CAD-2290: roles/owner does not include storage.buckets.get,
			// storage.buckets.{get,set}IamPolicy, storage.objects.{delete,list}.
			// Those come from the legacy bucket role bound to projectOwner.
			name: "owner also gets the legacy bucket role",
			policies: policy(&cloudresourcemanagerV3.Binding{
				Role: "roles/owner", Members: []string{"user:MPan@fortinet.com"},
			}),
			email:    "mpan@fortinet.com",
			expected: []string{"roles/owner", "roles/storage.legacyBucketOwner"},
		},
		{
			name: "editor also gets the legacy bucket role",
			policies: policy(&cloudresourcemanagerV3.Binding{
				Role: "roles/editor", Members: []string{"user:mpan@fortinet.com"},
			}),
			email:    "mpan@fortinet.com",
			expected: []string{"roles/editor", "roles/storage.legacyBucketOwner"},
		},
		{
			name: "a role bound in more than one policy is only returned once",
			policies: []*cloudresourcemanagerV3.Policy{
				{Bindings: []*cloudresourcemanagerV3.Binding{{
					Role: "roles/owner", Members: []string{"user:mpan@fortinet.com"},
				}}},
				{Bindings: []*cloudresourcemanagerV3.Binding{{
					Role: "roles/owner", Members: []string{"user:mpan@fortinet.com"},
				}}},
			},
			email:    "mpan@fortinet.com",
			expected: []string{"roles/owner", "roles/storage.legacyBucketOwner"},
		},
		{
			name: "viewer gets the legacy reader role",
			policies: policy(&cloudresourcemanagerV3.Binding{
				Role: "roles/viewer", Members: []string{"user:mpan@fortinet.com"},
			}),
			email:    "mpan@fortinet.com",
			expected: []string{"roles/viewer", "roles/storage.legacyBucketReader"},
		},
		{
			name: "custom roles are returned as-is",
			policies: policy(&cloudresourcemanagerV3.Binding{
				Role: "projects/p/roles/lw_agentless", Members: []string{"user:mpan@fortinet.com"},
			}),
			email:    "mpan@fortinet.com",
			expected: []string{"projects/p/roles/lw_agentless"},
		},
		{
			name: "bindings for other principals are ignored",
			policies: policy(&cloudresourcemanagerV3.Binding{
				Role: "roles/owner", Members: []string{"user:someone-else@fortinet.com"},
			}),
			email:    "mpan@fortinet.com",
			expected: []string{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			roles := rolesForCaller(c.policies, c.email)
			got := make([]string, 0, len(roles))
			for r := range roles {
				got = append(got, r)
			}
			assert.ElementsMatch(t, c.expected, got)
		})
	}
}
