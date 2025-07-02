package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
)

func FetchPolicies(p *Preflight) error {
	if p.caller.IsAdmin {
		return nil
	}

	p.verboseWriter.Write(fmt.Sprintf("Discovering role assigments for subscription %s", p.azureConfig.subscriptionID))

	// Get role assignments for the caller
	client, err := armauthorization.NewRoleAssignmentsClient(p.azureConfig.subscriptionID, p.azureConfig.cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create role assignments client: %v", err)
	}

	pager := client.NewListPager(&armauthorization.RoleAssignmentsClientListOptions{
		Filter: to.Ptr(fmt.Sprintf("principalId eq '%s'", p.caller.ObjectID)),
	})

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get next page: %v", err)
		}

		for _, assignment := range page.Value {
			if assignment.Properties != nil && assignment.Properties.RoleDefinitionID != nil {
				// Get role definition to check permissions
				roleDefClient, err := armauthorization.NewRoleDefinitionsClient(p.azureConfig.cred, nil)
				if err != nil {
					return err
				}

				roleDef, err := roleDefClient.GetByID(context.Background(), *assignment.Properties.RoleDefinitionID, nil)
				if err != nil {
					return err
				}

				if roleDef.Properties != nil && roleDef.Properties.Permissions != nil {
					for _, permission := range roleDef.Properties.Permissions {
						for _, action := range permission.Actions {
							if action != nil && *action == "*" {
								p.caller.IsAdmin = true
								return nil
							}
							if strings.Contains(*action, "*") {
								p.permissionsWithWildcard = append(p.permissionsWithWildcard, *action)
							} else if action != nil {
								p.permissions[*action] = true
							}
						}
					}
				}
			}
		}
	}

	return nil
}
