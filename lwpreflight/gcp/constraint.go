package gcp

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func CheckOrgPolicyConstraints(p *Preflight) error {

	crmSvc, err := cloudresourcemanager.NewService(context.Background(), p.gcpClientOption)
	if err != nil {
		return err
	}

	policiesResponse, err := crmSvc.Organizations.ListOrgPolicies(
		fmt.Sprintf("organizations/%s", p.orgID),
		&cloudresourcemanager.ListOrgPoliciesRequest{},
	).Do()
	if err != nil {
		return err
	}

	for _, policy := range policiesResponse.Policies {
		if policy.BooleanPolicy == nil || !policy.BooleanPolicy.Enforced {
			continue
		}
		if strings.Contains(policy.Constraint, "disableServiceAccountKeyCreation") {
			// Populate the error for every integration
			for _, integrationType := range p.integrationTypes {
				p.errors[integrationType] = append(
					p.errors[integrationType],
					"IAM disableServiceAccountKeyCreation constraint is enabled. Please disable it to continue.",
				)
			}
			return nil
		}
	}

	return nil
}
