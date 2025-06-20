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
	roles := make(map[string]bool)

	iamSvc, err := iam.NewService(context.Background(), p.gcpClientOption)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		for _, b := range policy.Bindings {
			// Continue if already processed
			if roles[b.Role] {
				continue
			}
			for _, m := range b.Members {
				if strings.Contains(strings.ToLower(m), strings.ToLower(p.caller.Email)) {
					role, err := iamSvc.Roles.Get(b.Role).Do()
					if err == nil {
						permissions = append(permissions, role.IncludedPermissions...)
					}
					roles[b.Role] = true
					break
				}
			}
		}
	}

	for _, permission := range permissions {
		p.permissions[permission] = true
	}

	return nil
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
