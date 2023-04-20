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

package lacework

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// Mock is a quick HTTP server that can be used to mock a Lacework API
// server, you can use it to avoid using a real server in our unit tests
//
// A simple usage:
//
//     func TestSomethingNew(t *testing.T) {
//         fakeServer := lacework.NewServer()
//         fakeServer.MockToken("TOKEN")
//         defer fakeServer.Close()
//
//         // Make sure to pass the fake API server URL
//         c, err := api.NewClient("test", api.WithURL(fakeServer.URL()))
//         if assert.Nil(t, err) {
//         	// The client c is ready to be used
//         }
//     }
type Mock struct {
	Mux        *http.ServeMux
	Server     *httptest.Server
	ApiVersion string
}

// MockServer returns a new mocked http server with a mutex
func MockServer() *Mock {
	mux := http.NewServeMux()
	return &Mock{
		Mux:        mux,
		Server:     httptest.NewServer(mux),
		ApiVersion: "v1",
	}
}

func MockUnstartedServer() *Mock {
	mux := http.NewServeMux()
	return &Mock{
		Mux:        mux,
		Server:     httptest.NewUnstartedServer(mux),
		ApiVersion: "v2",
	}
}

func (m *Mock) UseApiV2() {
	m.ApiVersion = "v2"
}

// MockAPI will mock the api path inside the server mutex with the provided handler function
func (m *Mock) MockAPI(p string, handler func(http.ResponseWriter, *http.Request)) {
	m.Mux.HandleFunc(fmt.Sprintf("/api/%s/%s", m.ApiVersion, p), handler)
}

func (s *Mock) MockToken(token string) {
	s.MockAPI("access/tokens", func(w http.ResponseWriter, r *http.Request) {
		expiration := time.Now().AddDate(0, 0, 1)
		fmt.Fprintf(w, `
      {
        "data": [{
          "expiresAt": "`+expiration.Format("Jan 02 2006 15:04")+`",
          "token": "`+token+`"
        }],
        "ok": true,
        "message": "SUCCESS"
      }
    `)
	})
}

func (s *Mock) MockTokenV2(token string) {
	s.UseApiV2()
	s.MockAPI("access/tokens", func(w http.ResponseWriter, r *http.Request) {
		expiration := time.Now().AddDate(0, 0, 1)
		fmt.Fprintf(w, `
      {
        "expiresAt": "`+expiration.Format(time.RFC3339)+`",
        "token": "`+token+`"
      }
    `)
	})
}

func (m Mock) URL() string {
	if m.Server != nil {
		return m.Server.URL
	}
	return ""
}

func (m Mock) Close() {
	if m.Server != nil {
		m.Server.Close()
	}
}
