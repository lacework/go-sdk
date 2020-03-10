package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// NewRequest generates a new http request
func (c *Client) NewRequest(method string, apiURL string, body io.Reader) (*http.Request, error) {
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
		// verify that the client has a token, if not, try to generate one
		if c.auth.token == "" {
			if _, err = c.GenerateToken(); err != nil {
				return nil, err
			}
		}
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

// DoDecoder is used to execute (aka Do) the http request and
// decode it into the provided interface, all at once
func (c *Client) DoDecoder(req *http.Request, v interface{}) (*http.Response, error) {
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	err = checkErrorInResponse(res)
	if err != nil {
		return res, err
	}

	var (
		resBuf bytes.Buffer

		// by using a TeeReader for capturing the reader’s data we avoid
		// interfering with the consumer of the reader
		resTee = io.TeeReader(res.Body, &resBuf)
	)
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resTee)
			return res, err
		}
		err = json.NewDecoder(resTee).Decode(v)
	}

	return res, err
}

// requestDecoder performs an http request on an endpoint, and
// decodes the response into the provided interface, all at once
func (c *Client) RequestDecoder(method, path string, body io.Reader, v interface{}) error {
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

// Do calls request.Do() directly
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.c.Do(req)
}
