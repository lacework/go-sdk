package azure

import (
	"fmt"
	"slices"
)

// CheckDirectoryRoles validates the Entra ID directory roles that deployment
// needs but that ARM permission checks cannot see. Deployment creates an
// Entra ID application (all integration types when a new AD application is
// created; agentless always creates its own), and config/activity log also
// assign the Directory Readers role to it, which requires Privileged Role
// Administrator. Runs unconditionally: subscription Owner/Contributor
// (IsAdmin) is orthogonal to directory roles.
func CheckDirectoryRoles(p *Preflight) error {
	p.verboseWriter.Write("Checking Entra ID directory roles")

	hasAny := func(roleIDs ...string) bool {
		for _, id := range roleIDs {
			if slices.Contains(p.caller.DirectoryRoles, id) {
				return true
			}
		}
		return false
	}

	canCreateApp := hasAny(
		ApplicationAdministratorRoleID,
		CloudApplicationAdministratorRoleID,
		GlobalAdministratorRoleID,
	)
	canAssignDirectoryRole := hasAny(
		PrivilegedRoleAdministratorRoleID,
		GlobalAdministratorRoleID,
	)

	for _, integrationType := range p.integrationTypes {
		// The agentless module always creates its own Entra ID application;
		// config/activity log only do so when not reusing an existing one.
		createsApp := integrationType == Agentless || !p.useExistingAdApplication
		// Only the ad-application module (config/activity log, new app only)
		// assigns the Directory Readers role.
		assignsDirectoryRole := integrationType != Agentless && !p.useExistingAdApplication

		if createsApp && !canCreateApp {
			p.errors[integrationType] = append(p.errors[integrationType], fmt.Sprintf(
				"Required Entra ID directory role missing: Application Administrator "+
					"(or Cloud Application Administrator) to create the Lacework Entra ID "+
					"application for %s. Activate the role in PIM or re-authenticate if it "+
					"was just assigned", integrationType))
		}
		if assignsDirectoryRole && !canAssignDirectoryRole {
			p.errors[integrationType] = append(p.errors[integrationType], fmt.Sprintf(
				"Required Entra ID directory role missing: Privileged Role Administrator "+
					"to assign the Directory Readers role to the Lacework Entra ID "+
					"application for %s. Activate the role in PIM or re-authenticate if it "+
					"was just assigned", integrationType))
		}
	}

	return nil
}
