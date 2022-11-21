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
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

var alertCommentJSON = `{
    "data": {
        "id": 38613065,
        "alertId": 966883,
        "entryType": "Comment",
        "entryAuthorType": "UserUpdate",
        "message": {
            "value": "this is fine"
        },
        "externalTime": "0",
        "user": {
            "userGuid": "QAN6DB45_FFF67EEEEEEEEEBBDACD9",
            "username": "david.hazekamp@lacework.net"
        },
        "updateContext": {}
    }
}`

func TestAlertsCommentMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d/comment", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Comment should be a POST method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.Comment(alertID, "this is a comment")
	assert.Nil(t, err)
}

func TestAlertsCommentOK(t *testing.T) {
	mockResponse := alertCommentJSON

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d/comment", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, mockResponse)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	commentExpected := api.AlertsCommentResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &commentExpected)

	var commentActual api.AlertsCommentResponse
	commentActual, err = c.V2.Alerts.Comment(alertID, "this is a comment")
	assert.Nil(t, err)

	assert.Equal(t, commentExpected, commentActual)
}

func TestAlertsCommentBad(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d/comment", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.Comment(alertID, "")
	if err == nil {
		assert.FailNow(t, "expected error not thrown")
	}
	assert.Equal(t, "alert comment must be provided", err.Error())
}

func TestAlertsCommentError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d/comment", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, lqlErrorReponse, http.StatusInternalServerError)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.Comment(alertID, "")
	assert.NotNil(t, err)
}
