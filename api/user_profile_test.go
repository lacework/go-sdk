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
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestV2UserProfile(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI("UserProfile",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "UserProfile should be a GET method")

			fmt.Fprintf(w, userProfileResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.UserProfile.Get()
	assert.Nil(t, err)

	if assert.NotNil(t, response) {
		if assert.Equal(t, 1, len(response.Data)) {
			profile := response.Data[0]
			assert.Equal(t, "salim.afiunemaya@lacework.net", profile.Username)

			if assert.Equal(t, 1, len(profile.Accounts)) {
				assert.Equal(t, "CUSTOMERDEMO", profile.Accounts[0].AccountName)
			}
		}
	}
}

// @afiune real response from a demo environment
func userProfileResponse() string {
	return `
{
  "data": [
    {
      "username": "salim.afiunemaya@lacework.net",
      "orgAccount": true,
      "url": "customerdemo.lacework.net",
      "orgAdmin": true,
      "orgUser": false,
      "accounts": [
        {
          "admin": true,
          "accountName": "CUSTOMERDEMO",
          "custGuid": "CUSTOMER_123455854C4272A5AC58B9FAA369C0ABB05564A91DA0ED9",
          "userGuid": "CUSTOMER_12345E3EDECD89F30125BCCFCFD308CCABE8DF908A08DD3",
          "userEnabled": 1
        }
      ]
    }
  ]
}
`
}
