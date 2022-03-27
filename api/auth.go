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

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/lacework/go-sdk/internal/format"
)

const DefaultTokenExpiryTime = 3600

// authConfig representing information like key_id, secret and token
// used for authenticating requests
type authConfig struct {
	keyID      string
	secret     string
	token      string
	expiration int
	expiresAt  time.Time
}

// WithApiKeys sets the key_id and secret used to generate API access tokens
func WithApiKeys(id, secret string) Option {
	return clientFunc(func(c *Client) error {
		if c.auth == nil {
			c.auth = &authConfig{}
		}

		c.log.Debug("setting up auth",
			zap.String("key", id),
			zap.String("secret", format.Secret(4, secret)),
		)
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
		c.log.Debug("setting up auth", zap.String("token", format.Secret(4, token)))
		c.auth.token = token
		c.auth.expiresAt = time.Now().UTC().Add(DefaultTokenExpiryTime * time.Second)
		return nil
	})
}

// WithTokenAndExpiration sets the token used to authenticate the API requests
// and additionally configures the expiration of the token
func WithTokenAndExpiration(token string, expiration time.Time) Option {
	return clientFunc(func(c *Client) error {
		c.log.Debug("setting up auth",
			zap.String("token", format.Secret(4, token)),
			zap.Time("expires_at", expiration),
		)
		c.auth.token = token
		c.auth.expiresAt = expiration.UTC()
		return nil
	})
}

// WithExpirationTime configures the token expiration time
func WithExpirationTime(t int) Option {
	return clientFunc(func(c *Client) error {
		c.log.Debug("setting up auth", zap.Int("expiration", t))
		c.auth.expiration = t
		c.auth.expiresAt = time.Now().UTC().Add(time.Duration(t) * time.Second)
		return nil
	})
}

func (c *Client) TokenExpired() bool {
	return c.auth.expiresAt.Sub(time.Now().UTC()) <= 0
}

// GenerateToken generates a new access token
func (c *Client) GenerateToken() (*TokenData, error) {
	if c.auth.keyID == "" || c.auth.secret == "" {
		return nil, fmt.Errorf("unable to generate access token: auth keys missing")
	}

	body, err := jsonReader(tokenRequest{c.auth.keyID, c.auth.expiration})
	if err != nil {
		return nil, err
	}

	request, err := c.NewRequest("POST", apiTokens, body)
	if err != nil {
		return nil, err
	}

	var tokenData TokenData
	switch c.ApiVersion() {
	case "v2":
		res, err := c.DoDecoder(request, &tokenData)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
	default:
		// we default to v1
		var tokenV1 TokenV1Response
		res, err := c.DoDecoder(request, &tokenV1)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		tokenData.Token = tokenV1.Token()
		tokenData.ExpiresAt = tokenV1.ExpiresAt()
	}

	c.log.Debug("storing token",
		zap.String("token", format.Secret(4, tokenData.Token)),
		zap.Time("expires_at", tokenData.ExpiresAt),
	)
	c.auth.token = tokenData.Token
	c.auth.expiresAt = tokenData.ExpiresAt
	if err != nil {
		c.log.Error("failed to parse token expiration response", zap.Error(err))
	}
	return &tokenData, nil
}

// GenerateTokenWithKeys generates a new access token with the provided keys
func (c *Client) GenerateTokenWithKeys(keyID, secretKey string) (*TokenData, error) {
	c.log.Debug("setting up auth",
		zap.String("key", keyID),
		zap.String("secret", format.Secret(4, secretKey)),
	)
	c.auth.keyID = keyID
	c.auth.secret = secretKey
	return c.GenerateToken()
}

type tokenRequest struct {
	KeyID      string `json:"keyId"`
	ExpiryTime int    `json:"expiryTime"`
}

// APIv2
type TokenData struct {
	ExpiresAt time.Time `json:"expiresAt"`
	Token     string    `json:"token"`
}

// APIv1
type TokenV1Data struct {
	ExpiresAt string `json:"expiresAt"`
	Token     string `json:"token"`
}

type TokenV1Response struct {
	Data    []TokenV1Data `json:"data"`
	Ok      bool          `json:"ok"`
	Message string        `json:"message"`
}

// Soon-To-Be-Deprecated
func (v1 TokenV1Response) Token() string {
	if len(v1.Data) > 0 {
		return v1.Data[0].Token
	}

	return ""
}

// Soon-To-Be-Deprecated
func (v1 TokenV1Response) ExpiresAt() time.Time {
	if len(v1.Data) > 0 {
		expiresAtTime, err := time.Parse("Jan 02 2006 15:04", v1.Data[0].ExpiresAt)
		if err == nil {
			return expiresAtTime
		}
	}

	return time.Now().UTC()
}
