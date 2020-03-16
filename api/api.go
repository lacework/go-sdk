//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package api

import "fmt"

const (
	//apiIntegrationType = "external/integrations/type/%s/"
	apiIntegrations      = "external/integrations"
	apiIntegrationByGUID = "external/integrations/%s"
	apiTokens            = "access/tokens"
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
