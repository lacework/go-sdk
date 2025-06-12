package azure

import (
	"fmt"
)

func CheckPermissions(p *Preflight) error {
	if p.caller.IsAdmin {
		return nil
	}

	for _, integrationType := range p.integrationTypes {
		requiredPermissions := RequiredPermissions[integrationType]
		for _, permission := range requiredPermissions {
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
