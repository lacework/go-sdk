package client_test

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
//         c, err := client.New("test", client.WithURL(fakeAPI.URL()))
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
