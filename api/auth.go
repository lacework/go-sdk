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

const defaultTokenExpiryTime = 3600

// authConfig representing information like key_id, secret and token
// used for authenticating requests
type authConfig struct {
	keyID      string
	secret     string
	token      string
	expiration int
}

// WithApiKeys sets the key_id and secret used to generate API access tokens
func WithApiKeys(id, secret string) Option {
	return clientFunc(func(c *Client) error {
		if c.auth == nil {
			c.auth = &authConfig{}
		}
		c.auth.keyID = id
		c.auth.secret = secret
		return nil
	})
}

// WithTokenFromKeys sets the API access keys and triggers a new token generation
// NOTE: Order matters when using this option, use it at the end of a NewClient() func
func WithTokenFromKeys(id, secret string) Option {
	return clientFunc(func(c *Client) error {
		if c.auth == nil {
			c.auth = &authConfig{}
		}

		_, err := c.GenerateTokenWithKeys(id, secret)
		return err
	})
}

// WithToken sets the token used to authenticate the API requests
func WithToken(token string) Option {
	return clientFunc(func(c *Client) error {
		c.auth.token = token
		return nil
	})
}

// WithExpirationTime configures the token expiration time
func WithExpirationTime(t int) Option {
	return clientFunc(func(c *Client) error {
		c.auth.expiration = t
		return nil
	})
}

// GenerateToken generates a new access token
func (c *Client) GenerateToken() (response tokenResponse, err error) {
	if c.auth.keyID == "" || c.auth.secret == "" {
		err = fmt.Errorf("unable to generate access token: auth keys missing")
		return
	}

	body, err := jsonReader(tokenRequest{c.auth.keyID, c.auth.expiration})
	if err != nil {
		return
	}

	err = c.RequestDecoder("POST", apiTokens, body, &response)
	if err != nil {
		return
	}

	if len(response.Data) > 0 {
		// @afiune how do we handle cases where there is more than one token
		c.auth.token = response.Data[0].Token
	}

	return
}

// GenerateTokenWithKeys generates a new access token with the provided keys
func (c *Client) GenerateTokenWithKeys(keyID, secretKey string) (tokenResponse, error) {
	c.auth.keyID = keyID
	c.auth.secret = secretKey
	return c.GenerateToken()
}

type tokenResponse struct {
	Data    []tokenData `json:"data"`
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
}

func (tr tokenResponse) Token() string {
	if len(tr.Data) > 0 {
		// @afiune how do we handle cases where there is more than one token
		return tr.Data[0].Token
	}

	return ""
}

type tokenData struct {
	ExpiresAt string `json:"expiresAt"`
	Token     string `json:"token"`
}

type tokenRequest struct {
	KeyId      string `json:"keyId"`
	ExpiryTime int    `json:"expiryTime"`
}
