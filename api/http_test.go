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

package api_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestNewRequest(t *testing.T) {
	// TODO @afiune to-be-implemented!
}

func TestDoDecoder(t *testing.T) {
	// TODO @afiune to-be-implemented!
}

func TestDoDecoder204(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI(
		"foo",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusNoContent)
		},
	)
	defer fakeServer.Close()

	c, _ := api.NewClient("foo", api.WithURL(fakeServer.URL()), api.WithTokenFromKeys("KEY", "SECRET"), api.WithExpirationTime(-60))
	request, _ := c.NewRequest("GET", "foo", nil)

	var v interface{}
	_, err := c.DoDecoder(request, v)
	assert.Nil(t, err)
}

func TestRequestDecoder(t *testing.T) {
	// TODO @afiune to-be-implemented!
}

func TestDo(t *testing.T) {
	// TODO @afiune to-be-implemented!
}

// When a new request is made, an expired token should be refreshed
func TestNewRequestWithExpiredToken(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	c, _ := api.NewClient("foo", api.WithURL(fakeServer.URL()), api.WithTokenFromKeys("KEY", "SECRET"), api.WithExpirationTime(-60))
	expiredToken := c.TokenExpired()
	_, err := c.NewRequest("GET", "foo", nil)

	assert.Equal(t, expiredToken, true, "Token starting state should be expired")
	assert.Nilf(t, err, "Error refreshing token %v", err)
	assert.Equal(t, c.TokenExpired(), false, "Expired token should be refreshed")
}

// httpBodySniffer is like a request sniffer, it reads the body
// from the provided request without closing it
func httpBodySniffer(r *http.Request) string {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return ""
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}
