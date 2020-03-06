package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

type clientOption func(*client)

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

// NewRequest generates a new http request
func (c *client) NewRequest(method string, apiURL string, body io.Reader) (*http.Request, error) {
	apiPath, err := url.Parse(c.apiPath(apiURL))
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(apiPath)
	request, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// set all necessary headers
	headers := map[string]string{
		"Method": request.Method,
		"Accept": "application/json",
	}

	// handle the special case that we are requesting an access token
	if apiURL == apiTokens {
		headers["X-LW-UAKS"] = c.auth.secret
	} else {
		headers["Authorization"] = c.auth.token
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	if body != nil {
		// @afiune we should detect the content-type from the body
		// instead of hard-coding it here
		request.Header.Set("Content-Type", "application/json")
	}

	// parse and encode query string values
	values := request.URL.Query()
	request.URL.RawQuery = values.Encode()

	return request, nil
}

// Do calls request.Do() directly
func (c *client) Do(req *http.Request) (*http.Response, error) {
	return c.c.Do(req)
}

// DoDecoder is used to Do() the request and decode it into the provided interface, all at once
func (c *client) DoDecoder(req *http.Request, v interface{}) (*http.Response, error) {
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(res)
	if err != nil {
		return res, err
	}

	var resBuf bytes.Buffer
	resTee := io.TeeReader(res.Body, &resBuf)
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resTee)
		} else {
			err = json.NewDecoder(resTee).Decode(v)
			if err != nil {
				return res, err
			}
		}
	}
	return res, nil
}

// requestDecoder performs a request on an endpoint, and decodes the response into the passed in Type
func (c *client) requestDecoder(method, path string, body io.Reader, v interface{}) error {
	request, err := c.NewRequest(method, path, body)
	if err != nil {
		return err
	}

	res, err := c.DoDecoder(request, v)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return err
}

type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode)
}

// CheckResponse checks the provided response and generates an Error
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}
