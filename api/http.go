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
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"go.uber.org/zap"
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

	if body != nil {
		// @afiune we should detect the content-type from the body
		// instead of hard-coding it here
		headers["Content-Type"] = "application/json"
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	// parse and encode query string values
	values := request.URL.Query()
	request.URL.RawQuery = values.Encode()

	c.log.Debug("request",
		zap.String("method", request.Method),
		zap.String("url", c.baseURL.String()),
		zap.String("endpoint", apiPath.String()),
		zap.Reflect("headers", headers),
		zap.String("body", c.httpRequestBodySniffer(request)),
	)

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

	if v != nil {
		var (
			resBuf bytes.Buffer

			// by using a TeeReader for capturing the reader’s data we avoid
			// interfering with the consumer of the reader
			resTee = io.TeeReader(res.Body, &resBuf)
		)
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resTee)
			return res, err
		}
		err = json.NewDecoder(resTee).Decode(v)
	}

	return res, err
}

// RequestDecoder performs an http request on an endpoint, and
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

// RequestEncoderDecoder leverages RequestDecoder and performs an http request that first
// encodes the provider 'data' as a JSON Reader and passes it as the body to the request
func (c *Client) RequestEncoderDecoder(method, path string, data, v interface{}) error {
	body, err := jsonReader(data)
	if err != nil {
		return err
	}
	return c.RequestDecoder(method, path, body, v)
}

// Do calls request.Do() directly
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	response, err := c.c.Do(req)
	if err == nil {
		c.log.Debug("response",
			zap.String("from_req_url", req.URL.String()),
			zap.Int("code", response.StatusCode),
			zap.String("proto", response.Proto),
			zap.Reflect("headers", response.Header),
			zap.String("body", c.httpResponseBodySniffer(response)),
		)
	}
	return response, err
}

// httpRequestBodySniffer a request sniffer, it reads the body from the
// provided request without closing it (use only for debugging purposes)
func (c *Client) httpRequestBodySniffer(r *http.Request) string {
	if !c.debugMode() {
		// prevents sniffing the request if we are not in debug mode
		return ""
	}

	if r.Body == nil || r.Body == http.NoBody {
		// No need to sniff
		return ""
	}

	var stringBody string
	r.Body, stringBody = sniffBody(r.Body)

	return stringBody
}

// httpResponseBodySniffer a response sniffer, it reads the body from the
// provided response without closing it (use only for debugging purposes)
func (c *Client) httpResponseBodySniffer(r *http.Response) string {
	if !c.debugMode() {
		// prevents sniffing the response if we are not in debug mode
		return ""
	}

	if r.Body == nil || r.ContentLength == 0 {
		// No need to sniff
		return ""
	}

	var stringBody string
	r.Body, stringBody = sniffBody(r.Body)

	return stringBody
}

// a very simple body sniffer (use only for debugging purposes)
func sniffBody(body io.ReadCloser) (io.ReadCloser, string) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, ""
	}

	if err := body.Close(); err != nil {
		return nil, ""
	}

	return ioutil.NopCloser(bytes.NewBuffer(bodyBytes)), string(bodyBytes)
}
