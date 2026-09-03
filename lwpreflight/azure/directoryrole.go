package azure

import (
	"fmt"
	"slices"
)

// CheckDirectoryRoles validates the Entra ID privileges that deployment needs
// but that ARM permission checks cannot see. Deployment creates an Entra ID
// application (all integration types when a new AD application is created;
// agentless always creates its own), and config/activity log also assign the
// Directory Readers role to it. Either a directory role or the equivalent
// Microsoft Graph application permission satisfies each requirement. Runs
// unconditionally: subscription Owner/Contributor (IsAdmin) is orthogonal to
// both.
func CheckDirectoryRoles(p *Preflight) error {
	p.verboseWriter.Write("Checking Entra ID directory roles")

	// Directory roles are GUIDs and Graph application permissions are dotted
	// names, so the two claim spaces cannot collide in one list.
	held := slices.Concat(p.caller.DirectoryRoles, p.caller.GraphPermissions)
	hasAny := func(grants ...string) bool {
		for _, grant := range grants {
			if slices.Contains(held, grant) {
				return true
			}
		}
		return false
	}

	canCreateApp := hasAny(
		ApplicationAdministratorRoleID,
		CloudApplicationAdministratorRoleID,
		GlobalAdministratorRoleID,
		GraphApplicationReadWriteOwnedByPermission,
		GraphApplicationReadWriteAllPermission,
	)
	canAssignDirectoryRole := hasAny(
		PrivilegedRoleAdministratorRoleID,
		GlobalAdministratorRoleID,
		GraphRoleManagementReadWriteDirectoryPermission,
	)

	// A caller can hold the Graph permission rather than the directory role, so
	// a failure to read those permissions makes either message a guess.
	unread := ""
	if p.graphPermissionsErr != nil {
		unread = fmt.Sprintf(
			" (Microsoft Graph application permissions could not be read: %v)",
			p.graphPermissionsErr)
	}

	for _, integrationType := range p.integrationTypes {
		usesExistingAdApplication := p.useExistingAdApplication[integrationType]

		// The agentless module always creates its own Entra ID application;
		// config/activity log only do so when not reusing an existing one.
		createsApp := integrationType == Agentless || !usesExistingAdApplication
		// Only the ad-application module (config/activity log, new app only)
		// assigns the Directory Readers role.
		assignsDirectoryRole := integrationType != Agentless && !usesExistingAdApplication

		if createsApp && !canCreateApp {
			p.errors[integrationType] = append(p.errors[integrationType], fmt.Sprintf(
				"Required Entra ID directory role missing: Application Administrator "+
					"(or Cloud Application Administrator, or the "+
					GraphApplicationReadWriteOwnedByPermission+" Microsoft Graph "+
					"permission) to create the Lacework Entra ID application for %s%s",
				integrationType, unread))
		}
		if assignsDirectoryRole && !canAssignDirectoryRole {
			p.errors[integrationType] = append(p.errors[integrationType], fmt.Sprintf(
				"Required Entra ID directory role missing: Privileged Role Administrator "+
					"(or the "+GraphRoleManagementReadWriteDirectoryPermission+
					" Microsoft Graph permission) to assign the Directory Readers role "+
					"to the Lacework Entra ID application for %s%s",
				integrationType, unread))
		}
	}

	return nil
}
