package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/lacework/go-sdk/v2/lwpreflight/logger"

	"google.golang.org/api/cloudresourcemanager/v1"
	cloudresourcemanagerV3 "google.golang.org/api/cloudresourcemanager/v3"

	"google.golang.org/api/iam/v1"
)

// The basic roles (owner/editor/viewer) do not list the bucket- and
// object-scoped Cloud Storage permissions in their includedPermissions, so
// expanding them via roles.get under-reports what the caller can actually do.
// Cloud Storage grants those separately: every bucket created in a project
// gets default IAM bindings for the projectOwner/projectEditor/projectViewer
// convenience values, which map to the legacy storage roles below.
var basicRoleStorageRoles = map[string]string{
	"roles/owner":  "roles/storage.legacyBucketOwner",
	"roles/editor": "roles/storage.legacyBucketOwner",
	"roles/viewer": "roles/storage.legacyBucketReader",
}

func FetchPolicies(p *Preflight) error {
	var err error
	var policies []*cloudresourcemanagerV3.Policy

	if p.orgID != "" {
		p.verboseWriter.Write(fmt.Sprintf("Discovering IAM policies for organization %s", p.orgID))
		policies, err = fetchOrgPolicies(p)
	} else {
		p.verboseWriter.Write(fmt.Sprintf("Discovering IAM policies for project %s", p.projectID))
		policies, err = fetchProjectPolicies(p)
	}
	if err != nil {
		return err
	}

	// Loop through polices and fetch permissions
	permissions := []string{}

	iamSvc, err := iam.NewService(context.Background(), p.gcpClientOption)
	if err != nil {
		return err
	}

	roles := rolesForCaller(policies, p.caller.Email)

	for role := range roles {
		r, err := iamSvc.Roles.Get(role).Do()
		if err != nil {
			return err
		}
		permissions = append(permissions, r.IncludedPermissions...)
	}

	for _, permission := range permissions {
		p.permissions[permission] = true
	}

	return nil
}

// rolesForCaller returns the set of roles bound to the caller across the given
// policies, plus the legacy Cloud Storage role implied by any basic role.
func rolesForCaller(policies []*cloudresourcemanagerV3.Policy, email string) map[string]bool {
	roles := make(map[string]bool)
	for _, policy := range policies {
		for _, b := range policy.Bindings {
			if roles[b.Role] {
				continue
			}
			for _, m := range b.Members {
				if strings.Contains(strings.ToLower(m), strings.ToLower(email)) {
					roles[b.Role] = true
					if legacyRole, ok := basicRoleStorageRoles[b.Role]; ok {
						roles[legacyRole] = true
					}
					break
				}
			}
		}
	}
	return roles
}

func fetchProjectPolicies(p *Preflight) ([]*cloudresourcemanagerV3.Policy, error) {
	ctx := context.Background()

	crmSvc, err := cloudresourcemanager.NewService(ctx, p.gcpClientOption)
	if err != nil {
		return nil, err
	}

	crmSvcV3, err := cloudresourcemanagerV3.NewService(ctx, p.gcpClientOption)
	if err != nil {
		return nil, err
	}

	response, err := crmSvc.Projects.GetAncestry(p.projectID, &cloudresourcemanager.GetAncestryRequest{}).Do()
	if err != nil {
		return nil, err
	}

	policies := []*cloudresourcemanagerV3.Policy{}

	for _, a := range response.Ancestor {
		var policy *cloudresourcemanagerV3.Policy
		var err error
		policyRequest := &cloudresourcemanagerV3.GetIamPolicyRequest{
			Options: &cloudresourcemanagerV3.GetPolicyOptions{
				RequestedPolicyVersion: 3,
			},
		}

		switch a.ResourceId.Type {
		case "organization":
			policy, err = crmSvcV3.Organizations.GetIamPolicy(
				fmt.Sprintf("organizations/%s", a.ResourceId.Id),
				policyRequest,
			).Do()
		case "project":
			policy, err = crmSvcV3.Projects.GetIamPolicy(
				fmt.Sprintf("projects/%s", a.ResourceId.Id),
				policyRequest,
			).Do()
		case "folder":
			policy, err = crmSvcV3.Folders.GetIamPolicy(
				fmt.Sprintf("folders/%s", a.ResourceId.Id),
				policyRequest,
			).Do()
		}

		if err != nil {
			logger.Log.Warnf("cannot fetch policy (continuing): %s", err.Error())
			continue
		}
		policies = append(policies, policy)
	}

	return policies, err
}

func fetchOrgPolicies(p *Preflight) ([]*cloudresourcemanagerV3.Policy, error) {
	crmSvcV3, err := cloudresourcemanagerV3.NewService(context.Background(), p.gcpClientOption)
	if err != nil {
		return nil, err
	}

	policy, err := crmSvcV3.Organizations.GetIamPolicy(
		fmt.Sprintf("organizations/%s", p.orgID),
		&cloudresourcemanagerV3.GetIamPolicyRequest{},
	).Do()
	if err != nil {
		return nil, err
	}

	return []*cloudresourcemanagerV3.Policy{policy}, nil
}
