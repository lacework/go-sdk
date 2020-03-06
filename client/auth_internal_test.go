package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiPath(t *testing.T) {
	c1 := &client{apiVersion: "v1"}
	assert.Equal(t, "/api/v1/foo", c1.apiPath("foo"), "api path mismatch")
	assert.Equal(t,
		"/api/v1/access/tokens",
		c1.apiPath(apiTokens),
		"token api path mismatch")

	c2 := &client{apiVersion: "v2"}
	assert.Equal(t, "/api/v2/bar", c2.apiPath("bar"), "api path mismatch")
	assert.Equal(t,
		"/api/v2/external/integrations",
		c2.apiPath(apiIntegrations),
		"integrations api path mismatch")
}
