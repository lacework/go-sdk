package client

import "fmt"

const (
	//apiIntegrationType = "external/integrations/type/%s/"
	apiIntegrations = "external/integrations"
	apiTokens       = "access/tokens"
)

// WithApiV2 configures the client to use the API version 2 (/api/v2)
func WithApiV2() Option {
	return clientFunc(func(c *Client) error {
		c.apiVersion = "v2"
		return nil
	})
}

// ApiVersion returns the API client version
func (c *Client) ApiVersion() string {
	return c.apiVersion
}

// apiPath builds a path by using the current API version
func (c *Client) apiPath(p string) string {
	return fmt.Sprintf("/api/%s/%s", c.apiVersion, p)
}
