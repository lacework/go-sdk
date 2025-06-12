package azure

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
)

func FetchPolicies(p *Preflight) error {
	if p.caller.IsAdmin {
		return nil
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to create credential: %v", err)
	}

	// Get role assignments for the caller
	client, err := armauthorization.NewRoleAssignmentsClient(p.subscriptionID, cred, nil)
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
				log.Printf("Fetching role definition for ID: %s", *assignment.Properties.RoleDefinitionID)
				// Get role definition to check permissions
				roleDefClient, err := armauthorization.NewRoleDefinitionsClient(cred, nil)
				if err != nil {
					return fmt.Errorf("failed to create role definitions client: %v", err)
				}

				roleDef, err := roleDefClient.GetByID(context.Background(), *assignment.Properties.RoleDefinitionID, nil)
				if err != nil {
					return fmt.Errorf("failed to get role definition: %v", err)
				}

				if roleDef.Properties != nil && roleDef.Properties.Permissions != nil {
					log.Printf("Permissions for role %s:", *roleDef.Properties.RoleName)
					for _, permission := range roleDef.Properties.Permissions {
						for _, action := range permission.Actions {
							if action != nil {
								log.Printf("  - %s", *action)
								p.permissions[*action] = true
							}
						}
					}
				}
			}
		}
	}

	log.Printf("Complete list of permissions stored: %v", p.permissions)
	return nil
}
