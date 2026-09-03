package azure

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
)

type Caller struct {
	ObjectID    string
	DisplayName string
	PrincipalID string
	TenantID    string
	// true if the caller has the Owner or Contributor RBAC role on the
	// subscription; says nothing about Entra ID directory roles (see
	// DirectoryRoles)
	IsAdmin bool
	// Entra ID directory role template IDs actively assigned to the caller,
	// from the token's wids claim
	DirectoryRoles []string
	// Microsoft Graph application permissions granted to the caller, from the
	// roles claim of a Graph-audience token. Empty for a delegated (user)
	// credential, whose effective permission is bounded by its own directory
	// roles anyway.
	GraphPermissions []string
}

func FetchCaller(p *Preflight) error {
	p.verboseWriter.Write("Discovering caller information")

	// Get caller identity
	token, err := p.azureConfig.cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	})
	if err != nil {
		return fmt.Errorf("failed to get token: %v", err)
	}

	// Parse JWT token to get caller info
	claims, err := parseJWTClaims(token.Token)
	if err != nil {
		return err
	}

	// Check if caller has Owner or Contributor role
	isAdmin, err := checkAdminRole(p.azureConfig.cred, claims.ObjectID, p.azureConfig.subscriptionID)
	if err != nil {
		return err
	}

	// Best effort: a caller can hold Graph application permissions instead of a
	// directory role. Failing to read them only costs the checks that fall back
	// to directory roles alone, so it must not fail preflight.
	graphPermissions, err := fetchGraphPermissions(p.azureConfig.cred)
	if err != nil {
		p.graphPermissionsErr = err
		p.verboseWriter.Write(fmt.Sprintf(
			"Could not read Microsoft Graph application permissions, "+
				"checking Entra ID directory roles only: %v", err))
	}

	p.caller = Caller{
		ObjectID:         claims.ObjectID,
		DisplayName:      claims.DisplayName,
		PrincipalID:      claims.PrincipalID,
		TenantID:         claims.TenantID,
		IsAdmin:          isAdmin,
		DirectoryRoles:   claims.Wids,
		GraphPermissions: graphPermissions,
	}

	return nil
}

// fetchGraphPermissions reads the caller's Microsoft Graph application
// permissions from the roles claim of a Graph-audience token. The ARM token
// cannot carry them: app roles are scoped to the resource the token is for.
// Only the token is requested, never a Graph API call, so this needs no
// permission of its own.
func fetchGraphPermissions(cred azcore.TokenCredential) ([]string, error) {
	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://graph.microsoft.com/.default"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	claims, err := parseJWTClaims(token.Token)
	if err != nil {
		return nil, err
	}

	return claims.Roles, nil
}

func checkAdminRole(cred azcore.TokenCredential, objectID, subscriptionID string) (bool, error) {
	client, err := armauthorization.NewRoleAssignmentsClient(subscriptionID, cred, nil)
	if err != nil {
		return false, fmt.Errorf("failed checkAdminRole: %v", err)
	}

	pager := client.NewListPager(&armauthorization.RoleAssignmentsClientListOptions{
		Filter: to.Ptr(fmt.Sprintf("principalId eq '%s'", objectID)),
	})

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return false, fmt.Errorf("failed checkAdminRole: %v", err)
		}

		for _, assignment := range page.Value {
			if assignment.Properties != nil && assignment.Properties.RoleDefinitionID != nil {
				roleID := *assignment.Properties.RoleDefinitionID
				roleDefClient, err := armauthorization.NewRoleDefinitionsClient(cred, nil)
				if err != nil {
					return false, fmt.Errorf("failed to create role definitions client: %v", err)
				}
				roleDef, err := roleDefClient.GetByID(context.Background(), roleID, nil)
				if err != nil {
					continue // If we can't get the role definition, skip this assignment
				}
				// check if role name is Owner or Contributor
				if roleDef.Properties != nil && roleDef.Properties.RoleName != nil {
					roleName := *roleDef.Properties.RoleName
					if roleName == "Owner" || roleName == "Contributor" {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}

type JWTClaims struct {
	ObjectID    string   `json:"oid"`
	DisplayName string   `json:"name"`
	PrincipalID string   `json:"sub"`
	TenantID    string   `json:"tid"`
	Wids        []string `json:"wids"`
	Roles       []string `json:"roles"`
}

func parseJWTClaims(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode claims: %v", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %v", err)
	}

	return &claims, nil
}
