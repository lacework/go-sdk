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
	"fmt"
	"net/http"
	"net/http/httptest"
)

// LaceworkServer is a quick HTTP server that can be used to mock a Lacework API
// server, you can use it to avoid using a real server in our unit tests
//
// A simple usage:
//
//     func TestSomethingNew(t *testing.T) {
//         fakeAPI := NewLaceworkServer()
//         fakeAPI.MockToken("TOKEN")
//         defer fakeAPI.Close()
//
//         // Make sure to pass the fake API server URL
//         c, err := api.NewClient("test", api.WithURL(fakeAPI.URL()))
//         if assert.Nil(t, err) {
//         	// The client c is ready to be used
//         }
//     }
type LaceworkServer struct {
	Mux        *http.ServeMux
	Server     *httptest.Server
	ApiVersion string
}

// NewLaceworkServer returns a new mocked htp server with a mutex
func NewLaceworkServer() *LaceworkServer {
	mux := http.NewServeMux()
	return &LaceworkServer{
		Mux:        mux,
		Server:     httptest.NewServer(mux),
		ApiVersion: "v1",
	}
}

// MockAPI will mock the api path inside the server mutex with the provided handler function
func (api *LaceworkServer) MockAPI(p string, handler func(http.ResponseWriter, *http.Request)) {
	api.Mux.HandleFunc(fmt.Sprintf("/api/%s/%s", api.ApiVersion, p), handler)
}

func (api *LaceworkServer) MockToken(token string) {
	api.MockAPI("access/tokens", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
      {
        "data": [{
          "expiresAt": "Mar 10 2020 08:10",
          "token": "`+token+`"
        }],
        "ok": true,
        "message": "SUCCESS"
      }
    `)
	})
}

func (api LaceworkServer) URL() string {
	if api.Server != nil {
		return api.Server.URL
	}
	return ""
}

func (api LaceworkServer) Close() {
	if api.Server != nil {
		api.Server.Close()
	}
}
