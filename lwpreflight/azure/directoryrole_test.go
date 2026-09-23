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
	// the shape azidentity produces when the token endpoint answers: the first
	// line says only that authentication failed, the reason is in the body
	p := preflightWithRoles(nil, nil, Agentless, Config)
	p.graphPermissionsErr = errors.New("failed to get token: ClientSecretCredential authentication failed. \n" +
		"POST https://login.microsoftonline.com/tenant/oauth2/v2.0/token\n" +
		"--------------------------------------------------------------------------------\n" +
		"RESPONSE 401: 401 Unauthorized\n" +
		"--------------------------------------------------------------------------------\n" +
		"{\n" +
		"  \"error\": \"invalid_client\",\n" +
		"  \"error_description\": \"AADSTS7000215: Invalid client secret provided. Ensure the secret being " +
		"sent in the request is the client secret value, not the client secret ID, for a secret added to app " +
		"'00000000-0000-0000-0000-000000000000'. Trace ID: abc\\r\\nCorrelation ID: def\\r\\nTimestamp: 2026-09-23 19:00:00Z\"\n" +
		"}\n" +
		"--------------------------------------------------------------------------------\n" +
		"To troubleshoot, visit https://aka.ms/azsdk/go/identity/troubleshoot#client-secret")
	assert.NoError(t, CheckDirectoryRoles(p))

	require.Len(t, p.errors[Agentless], 1)
	require.Len(t, p.errors[Config], 2)
	// the caller learns the second path was never looked at, and why, on both
	// kinds of message, without the request dump
	for _, msg := range append(p.errors[Agentless], p.errors[Config]...) {
		assert.Contains(t, msg, "could not be read: AADSTS7000215: Invalid client secret provided.)")
		assert.NotContains(t, msg, "RESPONSE 401")
		assert.NotContains(t, msg, "Trace ID")
		assert.NotContains(t, msg, "\n")
	}

	assert.True(t, p.caller.GraphPermissionsUnread)

	// nothing appended when the permissions were read fine
	p = preflightWithRoles(nil, nil, Agentless)
	assert.NoError(t, CheckDirectoryRoles(p))
	assert.False(t, p.caller.GraphPermissionsUnread)
	require.Len(t, p.errors[Agentless], 1)
	assert.NotContains(t, p.errors[Agentless][0], "could not be read")
}

func TestGraphErrorReason(t *testing.T) {
	// no response from the token endpoint: the first line is the whole story
	assert.Equal(t, "failed to get token: dial tcp: lookup login.microsoftonline.com: no such host",
		graphErrorReason(errors.New("failed to get token: dial tcp: lookup login.microsoftonline.com: no such host")))
	// a multi-line error without an AADSTS code keeps its first line, trimmed
	assert.Equal(t, "failed to get token: DefaultAzureCredential: failed to acquire a token.",
		graphErrorReason(errors.New("failed to get token: DefaultAzureCredential: failed to acquire a token. \nAttempted credentials:")))
	// dotted names inside the first sentence do not cut it short
	assert.Equal(t, "AADSTS90002: Tenant 'contoso.onmicrosoft.com' not found.",
		graphErrorReason(errors.New(`"AADSTS90002: Tenant 'contoso.onmicrosoft.com' not found. Check to make sure you have the correct tenant ID."`)))
}

func TestCheckDirectoryRolesReportsCapabilities(t *testing.T) {
	cases := []struct {
		name                 string
		roles, graph         []string
		canCreate, canAssign bool
	}{
		{"none", nil, nil, false, false},
		{"application administrator only", []string{ApplicationAdministratorRoleID}, nil, true, false},
		{"privileged role administrator only", []string{PrivilegedRoleAdministratorRoleID}, nil, false, true},
		{"global administrator", []string{GlobalAdministratorRoleID}, nil, true, true},
		{"graph owned-by only", nil, []string{GraphApplicationReadWriteOwnedByPermission}, true, false},
		{"graph both", nil, []string{GraphApplicationReadWriteAllPermission, GraphRoleManagementReadWriteDirectoryPermission}, true, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// reported whatever the existing-application answers, so a caller
			// that waived the checks can still tell whether it had a choice
			p := preflightWithRoles(c.roles, map[IntegrationType]bool{Config: true, ActivityLog: true}, Config, ActivityLog)
			p.caller.GraphPermissions = c.graph
			require.NoError(t, CheckDirectoryRoles(p))
			assert.Equal(t, c.canCreate, p.caller.CanCreateApplication)
			assert.Equal(t, c.canAssign, p.caller.CanAssignDirectoryRole)
			assert.Empty(t, p.errors)
		})
	}
}
