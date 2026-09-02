package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/v2/lwpreflight/verbosewriter"
)

func preflightWithRoles(roles []string, useExistingAdApp bool, types ...IntegrationType) *Preflight {
	return &Preflight{
		integrationTypes:         types,
		useExistingAdApplication: useExistingAdApp,
		caller:                   Caller{DirectoryRoles: roles},
		errors:                   map[IntegrationType][]string{},
		verboseWriter:            verbosewriter.New(),
	}
}

func TestCheckDirectoryRolesMissingAll(t *testing.T) {
	p := preflightWithRoles(nil, false, Config, ActivityLog, Agentless)
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
		[]string{ApplicationAdministratorRoleID}, false, Config, ActivityLog, Agentless)
	assert.NoError(t, CheckDirectoryRoles(p))

	assert.Len(t, p.errors[Config], 1)
	assert.Contains(t, p.errors[Config][0], "Privileged Role Administrator")
	assert.Len(t, p.errors[ActivityLog], 1)
	// agentless does not need Privileged Role Administrator
	assert.Empty(t, p.errors[Agentless])
}

func TestCheckDirectoryRolesGlobalAdminSatisfiesAll(t *testing.T) {
	p := preflightWithRoles(
		[]string{GlobalAdministratorRoleID}, false, Config, ActivityLog, Agentless)
	assert.NoError(t, CheckDirectoryRoles(p))
	assert.Empty(t, p.errors)
}

func TestCheckDirectoryRolesExistingAdApplication(t *testing.T) {
	p := preflightWithRoles(nil, true, Config, ActivityLog, Agentless)
	assert.NoError(t, CheckDirectoryRoles(p))

	// existing AD app: config/activity log neither create an app nor assign roles
	assert.Empty(t, p.errors[Config])
	assert.Empty(t, p.errors[ActivityLog])
	// agentless always creates its own app
	assert.Len(t, p.errors[Agentless], 1)
}
