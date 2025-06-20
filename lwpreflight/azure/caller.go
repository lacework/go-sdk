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
	IsAdmin     bool // true if the caller has Owner or Contributor role
}

func FetchCaller(p *Preflight) error {
	p.verboseWriter.Write("Discovering caller information")

	// Get caller identity
	token, err := p.cred.GetToken(context.Background(), policy.TokenRequestOptions{
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
	isAdmin, err := checkAdminRole(p.cred, claims.ObjectID, p.subscriptionID)
	if err != nil {
		return err
	}

	p.caller = Caller{
		ObjectID:    claims.ObjectID,
		DisplayName: claims.DisplayName,
		PrincipalID: claims.PrincipalID,
		TenantID:    claims.TenantID,
		IsAdmin:     isAdmin,
	}

	return nil
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
				if strings.Contains(strings.ToLower(roleID), "/providers/microsoft.authorization/roledefinitions/owner") ||
					strings.Contains(strings.ToLower(roleID), "/providers/microsoft.authorization/roledefinitions/contributor") {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

type JWTClaims struct {
	ObjectID    string `json:"oid"`
	DisplayName string `json:"name"`
	PrincipalID string `json:"sub"`
	TenantID    string `json:"tid"`
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
