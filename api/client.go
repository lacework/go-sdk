package api

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const defaultTimeout = 10 * time.Second

type Client struct {
	account    string
	apiVersion string
	baseURL    *url.URL
	auth       *authConfig
	c          *http.Client
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
//   lacework, err := api.NewClient("demo")
//   if err == nil {
//       lacework.GetIntegrations()
//   }
func NewClient(account string, opts ...Option) (*Client, error) {
	baseURL, err := url.Parse(fmt.Sprintf("https://%s.lacework.net", account))
	if err != nil {
		return nil, err
	}

	c := &Client{
		account:    account,
		baseURL:    baseURL,
		apiVersion: "v1",
		auth: &authConfig{
			expiration: defaultTokenExpiryTime,
		},
		c: &http.Client{Timeout: defaultTimeout},
	}

	for _, opt := range opts {
		if err := opt.apply(c); err != nil {
			return c, err
		}
	}

	return c, nil
}

// WithURL sets the base URL, this options is only available for test purposes
func WithURL(baseURL string) Option {
	return clientFunc(func(c *Client) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		c.baseURL = u
		return nil
	})
}
