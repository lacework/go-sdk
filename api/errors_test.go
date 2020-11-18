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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestErrorWithDataMessage(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/any/endpoint",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `
{
  "data": {
    "statusMessage": "We may receive error messages from here. Catch it!"
  },
  "ok": true,
  "message": "SUCCESS"
}
`, http.StatusInternalServerError)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	var v interface{}
	err = c.RequestDecoder("GET", "external/any/endpoint", nil, v)
	assert.Nil(t, v)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "[GET] http://")
		assert.Contains(t, err.Error(), "/external/any/endpoint")
		assert.Contains(t,
			err.Error(),
			"[500] We may receive error messages from here. Catch it!",
		)
	}
}

func TestErrorWithDataStatusMessage(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/any/endpoint",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `
{
  "data": {
    "Status": "NotFound",
    "Message": "We should catch this error message!"
  },
  "ok": true,
  "message": "SUCCESS"
}
`, http.StatusBadRequest)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	var v interface{}
	err = c.RequestDecoder("GET", "external/any/endpoint", nil, v)
	assert.Nil(t, v)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "[GET] http://")
		assert.Contains(t, err.Error(), "/external/any/endpoint")
		assert.Contains(t,
			err.Error(), "[400] We should catch this error message!",
		)
	}
}

func TestErrorWithDataErrorMsg(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/any/endpoint",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `
{
  "data": {
    "ErrorMsg": "This is an error message we should display!"
  },
  "ok": true,
  "message": "SUCCESS"
}
`, http.StatusTooManyRequests)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	var v interface{}
	err = c.RequestDecoder("GET", "external/any/endpoint", nil, v)
	assert.Nil(t, v)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "[GET] http://")
		assert.Contains(t, err.Error(), "/external/any/endpoint")
		assert.Contains(t,
			err.Error(),
			"[429] This is an error message we should display!",
		)
	}
}

func TestErrorWithoutValidErrorMessageInDataField(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/any/endpoint",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `
{
  "data": {
    "MessageError": "We would never catch this error!"
  },
  "ok": true,
  "message": "SUCCESS"
}
`, http.StatusTooManyRequests)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	var v interface{}
	err = c.RequestDecoder("GET", "external/any/endpoint", nil, v)
	assert.Nil(t, v)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "[GET] http://")
		assert.Contains(t, err.Error(), "/external/any/endpoint")
		assert.Contains(t,
			err.Error(),
			"[429] Too Many Requests", // This is better than nothing
		)
		assert.NotContains(t,
			err.Error(),
			"We would never catch this error!",
		)
	}
}

func TestErrorInvalidStatusCode(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/any/endpoint",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `
{
  "data": {
    "YetAnotherErrorField": "We would never catch this error!"
  },
  "ok": true,
  "message": "SUCCESS"
}
`, 432) // Note that this status code is not valid
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	var v interface{}
	err = c.RequestDecoder("GET", "external/any/endpoint", nil, v)
	assert.Nil(t, v)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "[GET] http://")
		assert.Contains(t, err.Error(), "/external/any/endpoint")
		assert.Contains(t,
			err.Error(),
			"[432] Unknown", // This is better than nothing, I guess
		)
		assert.NotContains(t,
			err.Error(),
			"We would never catch this error!",
		)
	}
}
