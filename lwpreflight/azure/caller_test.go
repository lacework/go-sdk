package azure

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseJWTClaimsWids(t *testing.T) {
	payload := base64.RawURLEncoding.EncodeToString([]byte(
		`{"oid":"o","tid":"t","wids":["` + PrivilegedRoleAdministratorRoleID + `"]}`))
	claims, err := parseJWTClaims("h." + payload + ".s")
	require.NoError(t, err)
	assert.Equal(t, []string{PrivilegedRoleAdministratorRoleID}, claims.Wids)

	// no wids claim: nil slice, so every directory role check fails closed
	payload = base64.RawURLEncoding.EncodeToString([]byte(`{"oid":"o"}`))
	claims, err = parseJWTClaims("h." + payload + ".s")
	require.NoError(t, err)
	assert.Nil(t, claims.Wids)
}

func TestParseJWTClaimsRoles(t *testing.T) {
	// Graph application permissions arrive in the roles claim of a
	// Graph-audience token, never in the ARM token the caller check decodes
	payload := base64.RawURLEncoding.EncodeToString([]byte(
		`{"oid":"o","roles":["` + GraphApplicationReadWriteAllPermission + `"]}`))
	claims, err := parseJWTClaims("h." + payload + ".s")
	require.NoError(t, err)
	assert.Equal(t, []string{GraphApplicationReadWriteAllPermission}, claims.Roles)

	// no roles claim: nil slice, so the check falls back to directory roles
	payload = base64.RawURLEncoding.EncodeToString([]byte(`{"oid":"o"}`))
	claims, err = parseJWTClaims("h." + payload + ".s")
	require.NoError(t, err)
	assert.Nil(t, claims.Roles)
}
