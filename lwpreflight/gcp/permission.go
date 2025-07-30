package gcp

import (
	"fmt"
)

func CheckPermissions(p *Preflight) error {
	for _, integrationType := range p.integrationTypes {
		p.verboseWriter.Write(fmt.Sprintf("Checking permissions for %s", integrationType))

		requiredPermissions := RequiredPermissions
		if p.orgID != "" {
			requiredPermissions = RequiredPermissionsForOrg
		}

		for _, permission := range requiredPermissions[integrationType] {
			if !p.permissions[permission] {
				p.errors[integrationType] = append(
					p.errors[integrationType],
					fmt.Sprintf("Required permission missing: %s", permission),
				)
			}
		}
	}
	return nil
}
