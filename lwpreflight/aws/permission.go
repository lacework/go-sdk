package aws

import (
	"fmt"
	"regexp"
)

func CheckPermissions(p *Preflight) error {
	if p.caller.IsAdmin {
		return nil
	}

	for _, integrationType := range p.integrationTypes {
		p.verboseWriter.Write(fmt.Sprintf("Checking permissions for %s", integrationType))

		requiredPermissions := RequiredPermissions[integrationType]
		for _, permission := range requiredPermissions {
			// First check plain permission strings
			matched := p.permissions[permission]
			if !matched {
				// Then check permission strings that contain wildcard(*)
				for _, pattern := range p.permissionsWithWildcard {
					matched, _ = regexp.MatchString(pattern, permission)
					if matched {
						break
					}
				}
			}
			if !matched {
				p.errors[integrationType] = append(
					p.errors[integrationType],
					fmt.Sprintf("Required permission missing: %s", permission),
				)
			}
		}
	}
	return nil
}
