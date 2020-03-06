package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const defaultTimeout = 10 * time.Second

type client struct {
	account    string
	apiVersion string
	baseURL    *url.URL
	auth       *authConfig
	c          *http.Client
}

type Option interface {
	apply(c *client)
}

type clientFunc func(c *client)

func (fn clientFunc) apply(c *client) {
	fn(c)
}

// New generates a new Lacework API client
//
// Example of basic usage
//
//   lacework, err := client.New("demo")
//   if err == nil {
//       lacework.GetIntegrations()
//   }
func New(account string, opts ...Option) (*client, error) {
	baseUrl, err := url.Parse(fmt.Sprintf("https://%s.lacework.net", account))
	if err != nil {
		return nil, err
	}

	c := &client{
		account:    account,
		baseURL:    baseUrl,
		apiVersion: "v1",
		auth: &authConfig{
			expiration: defaultTokenExpiryTime,
		},
		c: &http.Client{Timeout: defaultTimeout},
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	return c, nil
}
