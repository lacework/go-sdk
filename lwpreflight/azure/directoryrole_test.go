package azure

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lacework/go-sdk/v2/lwpreflight/verbosewriter"
)

func preflightWithRoles(
	roles []string,
	useExistingAdApplication map[IntegrationType]bool,
	types ...IntegrationType,
) *Preflight {
	return &Preflight{
		integrationTypes:         types,
		useExistingAdApplication: useExistingAdApplication,
		caller:                   Caller{DirectoryRoles: roles},
		errors:                   map[IntegrationType][]string{},
		verboseWriter:            verbosewriter.New(),
	}
}

func preflightWithGraphPermissions(permissions []string, types ...IntegrationType) *Preflight {
	p := preflightWithRoles(nil, nil, types...)
	p.caller.GraphPermissions = permissions
	return p
}

func TestCheckDirectoryRolesMissingAll(t *testing.T) {
	p := preflightWithRoles(nil, nil, Config, ActivityLog, Agentless)
	assert.NoError(t, CheckDirectoryRoles(p))

	// config/activity log need app creation + directory role assignment
	assert.Len(t, p.errors[Config], 2)
	assert.Len(t, p.errors[ActivityLog], 2)
	// agentless creates an app but assigns no directory role
	assert.Len(t, p.errors[Agentless], 1)
	assert.Contains(t, p.errors[Agentless][0], "Application Administrator")
}

func TestCheckDirectoryRolesMissingPrivilegedRoleAdmin(t *testing.T) {
	p := preflightWithRoles(
		[]string{ApplicationAdministratorRoleID}, nil, Config, ActivityLog, Agentless)
	assert.NoError(t, CheckDirectoryRoles(p))

	assert.Len(t, p.errors[Config], 1)
	assert.Contains(t, p.errors[Config][0], "Privileged Role Administrator")
	assert.Len(t, p.errors[ActivityLog], 1)
	// agentless does not need Privileged Role Administrator
	assert.Empty(t, p.errors[Agentless])
}

func TestCheckDirectoryRolesGlobalAdminSatisfiesAll(t *testing.T) {
	p := preflightWithRoles(
		[]string{GlobalAdministratorRoleID}, nil, Config, ActivityLog, Agentless)
	assert.NoError(t, CheckDirectoryRoles(p))
	assert.Empty(t, p.errors)
}

func TestCheckDirectoryRolesExistingAdApplication(t *testing.T) {
	p := preflightWithRoles(nil, map[IntegrationType]bool{
		Config: true, ActivityLog: true,
	}, Config, ActivityLog, Agentless)
	assert.NoError(t, CheckDirectoryRoles(p))

	// existing AD app: config/activity log neither create an app nor assign roles
	assert.Empty(t, p.errors[Config])
	assert.Empty(t, p.errors[ActivityLog])
	// agentless always creates its own app
	assert.Len(t, p.errors[Agentless], 1)
}

func TestCheckDirectoryRolesMixedExistingAdApplication(t *testing.T) {
	tests := []struct {
		name                string
		existingIntegration IntegrationType
		newApplicationType  IntegrationType
	}{
		{"config reuses an application", Config, ActivityLog},
		{"activity log reuses an application", ActivityLog, Config},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := preflightWithRoles(nil, map[IntegrationType]bool{
				test.existingIntegration: true,
			}, Config, ActivityLog)
			assert.NoError(t, CheckDirectoryRoles(p))

			assert.Empty(t, p.errors[test.existingIntegration])
			assert.Len(t, p.errors[test.newApplicationType], 2)
		})
	}
}

func TestCheckDirectoryRolesGraphPermissionsSatisfyAll(t *testing.T) {
	// no directory role at all, both capabilities held as Graph app permissions
	p := preflightWithGraphPermissions([]string{
		GraphApplicationReadWriteAllPermission,
		GraphRoleManagementReadWriteDirectoryPermission,
	}, Config, ActivityLog, Agentless)
	assert.NoError(t, CheckDirectoryRoles(p))
	assert.Empty(t, p.errors)
}

func TestCheckDirectoryRolesGraphPermissionsPartial(t *testing.T) {
	tests := []struct {
		name          string
		permissions   []string
		configErrors  int
		agentlessErrs int
		wantError     string
	}{
		{
			name:          "app creation only, cannot assign the directory role",
			permissions:   []string{GraphApplicationReadWriteAllPermission},
			configErrors:  1,
			agentlessErrs: 0,
			wantError:     "Privileged Role Administrator",
		},
		{
			name:          "owned-by variant also creates applications",
			permissions:   []string{GraphApplicationReadWriteOwnedByPermission},
			configErrors:  1,
			agentlessErrs: 0,
			wantError:     "Privileged Role Administrator",
		},
		{
			name:          "role assignment only, cannot create the application",
			permissions:   []string{GraphRoleManagementReadWriteDirectoryPermission},
			configErrors:  1,
			agentlessErrs: 1,
			wantError:     "Application Administrator",
		},
		{
			name:          "an unrelated permission satisfies nothing",
			permissions:   []string{"Directory.Read.All"},
			configErrors:  2,
			agentlessErrs: 1,
			wantError:     "Application Administrator",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := preflightWithGraphPermissions(test.permissions, Config, Agentless)
			assert.NoError(t, CheckDirectoryRoles(p))

			assert.Len(t, p.errors[Config], test.configErrors)
			assert.Len(t, p.errors[Agentless], test.agentlessErrs)
			assert.Contains(t, p.errors[Config][0], test.wantError)
		})
	}
}

func TestCheckDirectoryRolesMixedRoleAndGraphPermission(t *testing.T) {
	// the case the change exists for: one capability from a directory role, the
	// other from a Graph application permission
	p := preflightWithRoles([]string{ApplicationAdministratorRoleID}, nil, Config, Agentless)
	p.caller.GraphPermissions = []string{GraphRoleManagementReadWriteDirectoryPermission}
	assert.NoError(t, CheckDirectoryRoles(p))
	assert.Empty(t, p.errors)

	// and the reverse direction
	p = preflightWithRoles([]string{PrivilegedRoleAdministratorRoleID}, nil, Config, Agentless)
	p.caller.GraphPermissions = []string{GraphApplicationReadWriteOwnedByPermission}
	assert.NoError(t, CheckDirectoryRoles(p))
	assert.Empty(t, p.errors)
}

func TestCheckDirectoryRolesUnreadableGraphPermissions(t *testing.T) {
	p := preflightWithRoles(nil, nil, Agentless)
	p.graphPermissionsErr = errors.New("AADSTS900023: tenant not found")
	assert.NoError(t, CheckDirectoryRoles(p))

	require.Len(t, p.errors[Agentless], 1)
	// the caller learns the second path was never looked at, so the missing
	// directory role is not reported as the whole story
	assert.Contains(t, p.errors[Agentless][0], "could not be read")
	assert.Contains(t, p.errors[Agentless][0], "AADSTS900023")

	// nothing appended when the permissions were read fine
	p = preflightWithRoles(nil, nil, Agentless)
	assert.NoError(t, CheckDirectoryRoles(p))
	require.Len(t, p.errors[Agentless], 1)
	assert.NotContains(t, p.errors[Agentless][0], "could not be read")
}
