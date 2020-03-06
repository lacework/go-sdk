package client

import "fmt"

const (
	apiIntegrations    = "external/integrations"
	apiIntegrationType = "external/integrations/type/%s/"
	apiTokens          = "access/tokens"
)

// WithApiV2 configures the client to use the API version 2 (/api/v2)
func WithApiV2() Option {
	return clientFunc(func(c *client) {
		c.apiVersion = "v2"
	})
}

// ApiVersion returns the API client version
func (c *client) ApiVersion() string {
	return c.apiVersion
}

// apiPath builds a path by using the current API version
func (c *client) apiPath(p string) string {
	return fmt.Sprintf("/api/%s/%s", c.apiVersion, p)
}
