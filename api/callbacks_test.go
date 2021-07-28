//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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
	"net/http"
	"testing"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestCallbacks(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI(
		"foo",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusBadRequest)
		},
	)
	defer fakeServer.Close()

	c, _ := api.NewClient("foo",
		api.WithApiV2(),
		api.WithURL(fakeServer.URL()),
		api.WithToken("TOKEN"),
		api.WithLifecycleCallbacks(api.LifecycleCallbacks{
			RequestCallback: func(httpCode int, _ http.Header) error {
				assert.Equal(t, http.StatusBadRequest, httpCode)
				return nil
			},
		}),
	)
	request, _ := c.NewRequest("GET", "foo", nil)

	_, err := c.Do(request)
	assert.Nil(t, err)
}

func TestCallbacksOnRequestError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI(
		"foo",
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Duration(10))
		},
	)
	defer fakeServer.Close()

	c, _ := api.NewClient("foo",
		api.WithApiV2(),
		api.WithURL(fakeServer.URL()),
		api.WithTimeout(time.Duration(1)),
		api.WithToken("TOKEN"),
		api.WithLifecycleCallbacks(api.LifecycleCallbacks{
			RequestCallback: func(httpCode int, _ http.Header) error {
				return nil
			},
		}),
	)
	request, _ := c.NewRequest("GET", "foo", nil)
	_, err := c.Do(request)
	assert.NotNil(t, err)
}
