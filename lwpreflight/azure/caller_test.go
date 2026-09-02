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
