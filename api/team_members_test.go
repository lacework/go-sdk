//
// Author:: Vatasha White (<vatasha.white@lacework.net>)
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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestTeamMembers_List(t *testing.T) {
	var (
		allGUIDs        []string
		teamMemberGuids = generateGuids(&allGUIDs, 3)
		expectedLen     = len(allGUIDs)
		fakeServer      = lacework.MockServer()
	)

	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("TeamMembers",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			teamMembers := []string{
				generateTeamMembers(teamMemberGuids),
			}
			fmt.Fprintf(w,
				generateTeamMembersResponse(
					strings.Join(teamMembers, ", "),
				),
			)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.TeamMembers.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.UserGuid)
	}
}

func TestTeamMembers_Get(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("TeamMembers/%s", intgGUID)
		teamMember = singleMockTeamMember(intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateTeamMemberResponse(teamMember))
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when the team member exists", func(t *testing.T) {
		var response api.TeamMemberResponse
		err := c.V2.TeamMembers.Get(intgGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.UserGuid)
			assert.Equal(t, "vatasha.white@lacework.net", response.Data.UserName)
			assert.Equal(t, "Lacework", response.Data.Props.Company)
		}
	})

	t.Run("when the team member doesn't exist", func(t *testing.T) {
		var response api.TeamMemberResponse
		err := c.V2.TeamMembers.Get("FAKE_GUID", &response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/TeamMembers/FAKE_GUID")
			assert.Contains(t, err.Error(), "[404] 404 page not found")
		}
	})
}

func TestTeamMembers_Create(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		teamMember = singleMockTeamMember(intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI("TeamMembers",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "POST", r.Method, "Create() should be a POST method") {
				fmt.Fprintf(w, generateTeamMemberResponse(teamMember))
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when the team member is successfully created", func(t *testing.T) {
		props := api.TeamMemberProps{
			AccountAdmin:           false,
			Company:                "Lacework",
			CreatedTime:            "2021-11-16T16:33:17.573Z",
			FirstName:              "Vatasha",
			JitCreated:             false,
			LastLoginTime:          "0",
			LastName:               "White",
			LastSessionCreatedTime: "0",
			OrgAdmin:               false,
			OrgUser:                false,
			UpdatedTime:            "0",
		}
		tm, err := api.NewTeamMember("vatasha.white@lacework.net", "1", props)
		assert.NoError(t, err)
		response, err := c.V2.TeamMembers.Create(tm)
		assert.NoError(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.UserGuid)
			assert.Equal(t, "vatasha.white@lacework.net", response.Data.UserName)
			assert.Equal(t, "Lacework", response.Data.Props.Company)
		}
	})

	t.Run("when the team member is not successfully created because userEnabled is greater than 1", func(t *testing.T) {
		props := api.TeamMemberProps{
			AccountAdmin:           false,
			Company:                "Lacework",
			CreatedTime:            "2021-11-16T16:33:17.573Z",
			FirstName:              "Vatasha",
			JitCreated:             false,
			LastLoginTime:          "0",
			LastName:               "White",
			LastSessionCreatedTime: "0",
			OrgAdmin:               false,
			OrgUser:                false,
			UpdatedTime:            "0",
		}
		_, err := api.NewTeamMember("vatasha.white@lacework.net", "2", props)
		assert.EqualError(t, err, "userEnabled field must be 0 for disabled or 1 for enabled")
	})

	t.Run("when the team member is not successfully created because userEnabled is less than 0", func(t *testing.T) {
		props := api.TeamMemberProps{
			AccountAdmin:           false,
			Company:                "Lacework",
			CreatedTime:            "2021-11-16T16:33:17.573Z",
			FirstName:              "Vatasha",
			JitCreated:             false,
			LastLoginTime:          "0",
			LastName:               "White",
			LastSessionCreatedTime: "0",
			OrgAdmin:               false,
			OrgUser:                false,
			UpdatedTime:            "0",
		}
		_, err := api.NewTeamMember("vatasha.white@gmail.com", "-100", props)
		assert.EqualError(t, err, "userEnabled field must be 0 for disabled or 1 for enabled")
	})
}

func generateTeamMembers(guids []string) string {
	tms := make([]string, len(guids))
	for i, guid := range guids {
		tms[i] = singleMockTeamMember(guid)
	}
	return strings.Join(tms, ", ")
}

func generateTeamMembersResponse(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func generateTeamMemberResponse(data string) string {
	return `
		{
			"data": ` + data + `
		}
	`
}

func singleMockTeamMember(id string) string {
	return fmt.Sprintf(`
    {
	  "custGuid": "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA",
      "props": {
            "accountAdmin": %t,
            "company": "Lacework",
            "createdTime": "2021-11-16T16:33:17.573Z",
            "firstName": "Vatasha",
            "jitCreated": %t,
            "lastLoginTime": "0",
            "lastName": "White",
			"lastSessionCreatedTime": "0",
			"orgAdmin": false,
      		"orgUser": false,
      		"updatedTime": "0"
	  },
      "userEnabled": "1",
      "userGuid": %q,
      "userName": "vatasha.white@lacework.net" 
    }
	`, false, false, id)
}
