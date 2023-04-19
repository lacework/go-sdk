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
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lacework/go-sdk/lwdomain"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	defaultTimeout    = 60 * time.Second
	defaultTLSTimeout = 63 * time.Second
)

type Client struct {
	id         string
	account    string
	subaccount string
	apiVersion string
	baseURL    *url.URL
	auth       *authConfig
	c          *http.Client
	log        *zap.Logger
	headers    map[string]string
	callbacks  LifecycleCallbacks

	Policy *PolicyService

	V2 *V2Endpoints
}

type Option interface {
	apply(c *Client) error
}

type clientFunc func(c *Client) error

func (fn clientFunc) apply(c *Client) error {
	return fn(c)
}

// New generates a new Lacework API client
//
// Example of basic usage
//
//	lacework, err := api.NewClient("demo")
//	if err == nil {
//	    lacework.Integrations.List()
//	}
func NewClient(account string, opts ...Option) (*Client, error) {
	if account == "" {
		return nil, errors.New("account cannot be empty")
	}

	// verify if the user provided the full qualified domain name
	if strings.Contains(account, ".lacework.net") {
		domain, err := lwdomain.New(account)
		if err != nil {
			return nil, err
		}
		account = domain.String()
	}

	baseURL, err := url.Parse(fmt.Sprintf("https://%s.lacework.net", account))
	if err != nil {
		return nil, err
	}

	c := &Client{
		id:         newID(),
		account:    account,
		baseURL:    baseURL,
		apiVersion: "v1",
		headers: map[string]string{
			"User-Agent": fmt.Sprintf("Go Client/%s", Version),
		},
		auth: &authConfig{
			expiration: DefaultTokenExpiryTime,
		},
		c: &http.Client{Timeout: defaultTimeout,
			Transport: &http.Transport{TLSHandshakeTimeout: defaultTLSTimeout}},
	}

	c.V2 = NewV2Endpoints(c)

	// init logger, this could change if a user calls api.WithLogLevel()
	c.initLogger("")

	for _, opt := range opts {
		if err := opt.apply(c); err != nil {
			return c, err
		}
	}

	c.log.Info("api client created",
		zap.String("url", c.baseURL.String()),
		zap.String("version", c.apiVersion),
		zap.String("log_level", c.log.Level().CapitalString()),
		zap.Int("timeout", c.auth.expiration),
	)
	return c, nil
}

// CopyClient generates a copy of the provider Lacework API Go client
//
// Example of basic usage
//
//	client, err := api.NewClient("demo")
//	if err == nil {
//	    client.Integrations.List()
//	}
//
//	clientCopy, err := api.CopyClient(client, api.WithOrgAccess())
//	if err == nil {
//	    clientCopy.Integrations.List()
//	}
func CopyClient(origin *Client, opts ...Option) (*Client, error) {
	dest := new(Client)
	*dest = *origin

	// no client should have the same ID
	dest.id = newID()

	for _, opt := range opts {
		if err := opt.apply(dest); err != nil {
			return dest, err
		}
	}
	return dest, nil
}

// WithSubaccount sets a subaccount into an API client
func WithSubaccount(subaccount string) Option {
	return clientFunc(func(c *Client) error {
		if subaccount != "" {
			c.log.Debug("setting up client", zap.String("subaccount", subaccount))
			c.subaccount = subaccount
			c.log.Debug("setting up header", zap.String("Account-Name", subaccount))
			c.headers["Account-Name"] = subaccount
		}
		return nil
	})
}

// WithTimeout changes the default client timeout
func WithTimeout(timeout time.Duration) Option {
	return clientFunc(func(c *Client) error {
		c.log.Debug("setting up client", zap.Reflect("timeout", timeout))
		c.c.Timeout = timeout
		return nil
	})
}

// WithTransport changes the default transport to increase TLSHandshakeTimeout
func WithTransport(transport *http.Transport) Option {
	return clientFunc(func(c *Client) error {
		c.c.Transport = transport
		return nil
	})
}

// WithURL sets the base URL, this options is only available for test purposes
func WithURL(baseURL string) Option {
	return clientFunc(func(c *Client) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		c.log.Debug("setting up client", zap.String("url", baseURL))
		c.baseURL = u
		return nil
	})
}

// WithHeader configures a HTTP Header to pass to every request
func WithHeader(header, value string) Option {
	return clientFunc(func(c *Client) error {
		if header != "" && value != "" {
			c.log.Debug("setting up header", zap.String(header, value))
			c.headers[header] = value
		}
		return nil
	})
}

// WithOrgAccess sets the Org-Access Header to access the organization level data sets
func WithOrgAccess() Option {
	return clientFunc(func(c *Client) error {
		c.log.Debug("setting up header", zap.String("Org-Access", "true"))
		c.headers["Org-Access"] = "true"
		return nil
	})
}

// URL returns the base url configured
func (c *Client) URL() string {
	return c.baseURL.String()
}

// ValidAuth verifies that the client has valid authentication
func (c *Client) ValidAuth() bool {
	return c.auth.token != ""
}

// OrgAccess check if the Org-Access header is set to 'true', if so,
// the client is configured to manage org level dataset
func (c *Client) OrgAccess() bool {
	return c.headers["Org-Access"] == "true"
}

// newID generates a new client id, this id is useful for logging purposes
// when there are more than one client running on the same machine
func newID() string {
	now := time.Now().UTC().UnixNano()
	seed := rand.New(rand.NewSource(now))
	return strconv.FormatInt(seed.Int63(), 16)
}
