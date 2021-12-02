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
	"encoding/json"
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
	assert.NoError(t, err)

	response, err := c.V2.TeamMembers.List()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.UserGuid)
	}
}

func TestTeamMembers_List_WithTimeFieldsAsInts(t *testing.T) {
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
				generateTeamMembersWithTimeFieldsAsInts(teamMemberGuids),
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
	assert.NoError(t, err)

	response, err := c.V2.TeamMembers.List()
	assert.NoError(t, err)
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
		teamMember = singleMockTeamMember(intgGUID, "vatasha.white@lacework.net")
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
	assert.NoError(t, err)

	t.Run("when the team member exists", func(t *testing.T) {
		var response api.TeamMemberResponse
		err := c.V2.TeamMembers.Get(intgGUID, &response)
		assert.NoError(t, err)
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
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "api/v2/TeamMembers/FAKE_GUID")
			assert.Contains(t, err.Error(), "[404] 404 page not found")
		}
	})
}

func TestTeamMembers_Get_WithTimeFieldsAsInts(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("TeamMembers/%s", intgGUID)
		teamMember = singleMockTeamMemberWithTimeFieldsAsInts(intgGUID, "vatasha.white@lacework.net")
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
	assert.NoError(t, err)

	t.Run("when the team member exists", func(t *testing.T) {
		var response api.TeamMemberResponse
		err := c.V2.TeamMembers.Get(intgGUID, &response)
		assert.NoError(t, err)
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
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "api/v2/TeamMembers/FAKE_GUID")
			assert.Contains(t, err.Error(), "[404] 404 page not found")
		}
	})
}

func TestTeamMembers_Create(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		teamMember = singleMockTeamMember(intgGUID, "vatasha.white@lacework.net")
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
	assert.NoError(t, err)

	t.Run("when the team member is successfully created", func(t *testing.T) {
		props := api.TeamMemberProps{
			AccountAdmin:           false,
			Company:                "Lacework",
			CreatedTime:            "2021-11-16T16:33:17.573Z",
			FirstName:              "Vatasha",
			JitCreated:             false,
			LastLoginTime:          "2021-11-17T16:33:17.573Z",
			LastName:               "White",
			LastSessionCreatedTime: "0",
			OrgAdmin:               false,
			OrgUser:                false,
			UpdatedTime:            "2021-11-18T16:33:17.573Z",
		}
		tm := api.NewTeamMember("vatasha.white@lacework.net", props)
		response, err := c.V2.TeamMembers.Create(tm)
		assert.NoError(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.UserGuid)
			assert.Equal(t, "vatasha.white@lacework.net", response.Data.UserName)
			assert.Equal(t, "Lacework", response.Data.Props.Company)
		}
	})
}

func TestTeamMembers_Create_WithTimeFieldsAsInts(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		teamMember = singleMockTeamMemberWithTimeFieldsAsInts(intgGUID, "vatasha.white@lacework.net")
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
	assert.NoError(t, err)

	t.Run("when the team member is successfully created", func(t *testing.T) {
		props := api.TeamMemberProps{
			AccountAdmin:           false,
			Company:                "Lacework",
			CreatedTime:            "2021-11-16T16:33:17.573Z",
			FirstName:              "Vatasha",
			JitCreated:             false,
			LastLoginTime:          0,
			LastName:               "White",
			LastSessionCreatedTime: 0,
			OrgAdmin:               false,
			OrgUser:                false,
			UpdatedTime:            0,
		}
		tm := api.NewTeamMember("vatasha.white@lacework.net", props)
		response, err := c.V2.TeamMembers.Create(tm)
		assert.NoError(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.UserGuid)
			assert.Equal(t, "vatasha.white@lacework.net", response.Data.UserName)
			assert.Equal(t, "Lacework", response.Data.Props.Company)
		}
	})
}

func TestTeamMembers_CreateOrg(t *testing.T) {
	var fakeServer = lacework.MockServer()

	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI("TeamMembers",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "POST", r.Method, "CreateOrg() should be a POST method") {
				fmt.Fprintf(w, singleMockOrgTeamMemberCreateResponse("vatasha.white@lacework.net"))
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
		api.WithOrgAccess(),
	)
	assert.NoError(t, err)

	t.Run("when the team member is successfully created", func(t *testing.T) {
		props := api.TeamMemberProps{
			Company:   "Lacework",
			FirstName: "Vatasha",
			LastName:  "White",
		}
		tm := api.NewTeamMemberOrg("vatasha.white@lacework.net", props)
		response, err := c.V2.TeamMembers.CreateOrg(tm)
		assert.NoError(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, "customerdemo.lacework.net", response.Data.Url)
			assert.Equal(t, "vatasha.white@lacework.net", response.Data.UserName)
		}
	})
}

func TestTeamMember_Update(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		teamMember = singleMockTeamMember(intgGUID, "vatasha.white+updated@lacework.net")
		fakeServer = lacework.MockServer()
		apiPath    = fmt.Sprintf("TeamMembers/%s", intgGUID)
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PATCH", r.Method, "Update() should be a PATCH method")

			if assert.NotNil(t, r.Body) {
				body := httpBodySniffer(r)
				assert.Contains(t, body, intgGUID, "IntgGUID missing")
				assert.Contains(t, body, "company\":\"Lacework", "missing company")
				assert.Contains(t, body, "firstName\":\"Vatasha", "missing first name")
				assert.Contains(t, body, "lastName\":\"White", "missing last name")
				assert.Contains(t, body, "userEnabled\":1", "missing user enabled")
				assert.Contains(t, body, "userName\":\"vatasha.white+updated@lacework.net", "missing username")
			}
			fmt.Fprintf(w, generateTeamMemberResponse(teamMember))
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	t.Run("when the team member is successfully updated", func(t *testing.T) {
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
		tm := api.NewTeamMember("vatasha.white+updated@lacework.net", props)
		tm.UserGuid = intgGUID
		response, err := c.V2.TeamMembers.Update(tm)
		assert.NoError(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, "vatasha.white+updated@lacework.net", response.Data.UserName)
		}
	})
}

func TestTeamMember_Update_WithTimeFieldsAsInts(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		teamMember = singleMockTeamMemberWithTimeFieldsAsInts(intgGUID, "vatasha.white+updated@lacework.net")
		fakeServer = lacework.MockServer()
		apiPath    = fmt.Sprintf("TeamMembers/%s", intgGUID)
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PATCH", r.Method, "Update() should be a PATCH method")

			if assert.NotNil(t, r.Body) {
				body := httpBodySniffer(r)
				assert.Contains(t, body, intgGUID, "IntgGUID missing")
				assert.Contains(t, body, "company\":\"Lacework", "missing company")
				assert.Contains(t, body, "firstName\":\"Vatasha", "missing first name")
				assert.Contains(t, body, "lastName\":\"White", "missing last name")
				assert.Contains(t, body, "userEnabled\":1", "missing user enabled")
				assert.Contains(t, body, "userName\":\"vatasha.white+updated@lacework.net", "missing username")
			}
			fmt.Fprintf(w, generateTeamMemberResponse(teamMember))
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	t.Run("when the team member is successfully updated", func(t *testing.T) {
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
		tm := api.NewTeamMember("vatasha.white+updated@lacework.net", props)
		tm.UserGuid = intgGUID
		response, err := c.V2.TeamMembers.Update(tm)
		assert.NoError(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, "vatasha.white+updated@lacework.net", response.Data.UserName)
		}
	})
}

func TestTeamMember_UpdateOrg(t *testing.T) {
	var (
		allGUIDs        []string
		teamMemberGuids = generateGuids(&allGUIDs, 2)
		fakeServer      = lacework.MockServer()
		apiPath         = fmt.Sprintf("TeamMembers/%s", teamMemberGuids[0])

		username   = fmt.Sprintf("vatasha.white+%s@lacework.net", teamMemberGuids[0])
		teamMember = singleMockTeamMembersOrgResponse(teamMemberGuids[0], teamMemberGuids[1], username)
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI("TeamMembers/search",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "POST", r.Method, "SearchUsername() should be a POST method") {
				var body api.SearchFilter
				decoder := json.NewDecoder(r.Body)
				err := decoder.Decode(&body)
				assert.NoError(t, err)
				if body.Filters[0].Value == username {
					teamMembers := []string{
						generateTeamMembers(teamMemberGuids),
					}
					fmt.Fprintf(w,
						generateTeamMembersResponse(
							strings.Join(teamMembers, ", "),
						),
					)
				} else {
					fmt.Fprint(w, "{}")
				}
			}
		},
	)

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PATCH", r.Method, "UpdateOrg() should be a PATCH method")

			if assert.NotNil(t, r.Body) {
				body := httpBodySniffer(r)
				assert.Contains(t, body, teamMemberGuids[0], "IntgGUID missing")
				assert.Contains(t, body, "company\":\"Lacework", "missing company")
				assert.Contains(t, body, "firstName\":\"Vatasha", "missing first name")
				assert.Contains(t, body, "lastName\":\"White", "missing last name")
				assert.Contains(t, body, "userEnabled\":1", "missing user enabled")
				assert.Contains(t, body, fmt.Sprintf("userName\":\"%s", username), "missing username")
			}
			fmt.Fprintf(w, teamMember)
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
		api.WithOrgAccess(),
	)
	assert.NoError(t, err)

	t.Run("when the team member is successfully updated", func(t *testing.T) {
		props := api.TeamMemberProps{
			Company:   "Lacework",
			FirstName: "Vatasha",
			LastName:  "White",
		}
		tm := api.NewTeamMemberOrg(username, props)
		tm.UserGuid = teamMemberGuids[0]
		response, err := c.V2.TeamMembers.UpdateOrg(tm)
		assert.NoError(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, username, response.Data[0].UserName)
		}
	})
}

func TestTeamMember_Delete(t *testing.T) {
	var (
		intgGUID        = intgguid.New()
		teamMember      = singleMockTeamMember(intgGUID, "vatasha.white@lacework.net")
		fakeServer      = lacework.MockServer()
		apiPath         = fmt.Sprintf("TeamMembers/%s", intgGUID)
		responseFromGet = generateTeamMemberResponse(teamMember)
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if responseFromGet != "" {
				switch r.Method {
				case "GET":
					fmt.Fprintf(w, responseFromGet)
				case "DELETE":
					responseFromGet = ""
				}
			} else {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	t.Run("verify that the team member exists", func(t *testing.T) {
		var response api.TeamMemberResponse
		err := c.V2.TeamMembers.Get(intgGUID, &response)
		assert.NoError(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.UserGuid)
			assert.Equal(t, "vatasha.white@lacework.net", response.Data.UserName)
			assert.Equal(t, "Lacework", response.Data.Props.Company)
		}
	})

	t.Run("when the team member has been deleted", func(t *testing.T) {
		err := c.V2.TeamMembers.Delete(intgGUID)
		assert.NoError(t, err)

		var response api.TeamMemberResponse
		err = c.V2.TeamMembers.Get(intgGUID, &response)
		assert.Empty(t, response)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), apiPath)
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestTeamMember_Delete_WithTimeFieldsAsInts(t *testing.T) {
	var (
		intgGUID        = intgguid.New()
		teamMember      = singleMockTeamMemberWithTimeFieldsAsInts(intgGUID, "vatasha.white@lacework.net")
		fakeServer      = lacework.MockServer()
		apiPath         = fmt.Sprintf("TeamMembers/%s", intgGUID)
		responseFromGet = generateTeamMemberResponse(teamMember)
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if responseFromGet != "" {
				switch r.Method {
				case "GET":
					fmt.Fprintf(w, responseFromGet)
				case "DELETE":
					responseFromGet = ""
				}
			} else {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	t.Run("verify that the team member exists", func(t *testing.T) {
		var response api.TeamMemberResponse
		err := c.V2.TeamMembers.Get(intgGUID, &response)
		assert.NoError(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.UserGuid)
			assert.Equal(t, "vatasha.white@lacework.net", response.Data.UserName)
			assert.Equal(t, "Lacework", response.Data.Props.Company)
		}
	})

	t.Run("when the team member has been deleted", func(t *testing.T) {
		err := c.V2.TeamMembers.Delete(intgGUID)
		assert.NoError(t, err)

		var response api.TeamMemberResponse
		err = c.V2.TeamMembers.Get(intgGUID, &response)
		assert.Empty(t, response)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), apiPath)
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestTeamMember_DeleteOrg(t *testing.T) {
	var (
		intgGUID        = intgguid.New()
		teamMember      = singleMockTeamMember(intgGUID, "vatasha.white@lacework.net")
		fakeServer      = lacework.MockServer()
		apiPath         = fmt.Sprintf("TeamMembers/%s", intgGUID)
		responseFromGet = generateTeamMemberResponse(teamMember)
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if responseFromGet != "" {
				switch r.Method {
				case "GET":
					fmt.Fprintf(w, responseFromGet)
				case "DELETE":
					responseFromGet = ""
				}
			} else {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
		api.WithOrgAccess(),
	)
	assert.NoError(t, err)

	t.Run("verify that the team member exists", func(t *testing.T) {
		var response api.TeamMemberResponse
		err := c.V2.TeamMembers.Get(intgGUID, &response)
		assert.NoError(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, intgGUID, response.Data.UserGuid)
			assert.Equal(t, "vatasha.white@lacework.net", response.Data.UserName)
			assert.Equal(t, "Lacework", response.Data.Props.Company)
		}
	})

	t.Run("when the team member has been deleted", func(t *testing.T) {
		err := c.V2.TeamMembers.DeleteOrg(intgGUID)
		assert.NoError(t, err)

		var response api.TeamMemberResponse
		err = c.V2.TeamMembers.Get(intgGUID, &response)
		assert.Empty(t, response)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), apiPath)
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestTeamMembers_SearchUsername(t *testing.T) {
	var (
		allGUIDs        []string
		teamMemberGuids = generateGuids(&allGUIDs, 2)
		fakeServer      = lacework.MockServer()
		username        = fmt.Sprintf("vatasha.white+%s@lacework.net", teamMemberGuids[0])
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI("TeamMembers/search",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "POST", r.Method, "SearchUsername() should be a POST method") {
				var body api.SearchFilter
				decoder := json.NewDecoder(r.Body)
				err := decoder.Decode(&body)
				assert.NoError(t, err)
				if body.Filters[0].Value == username {
					teamMembers := []string{
						generateTeamMembers(teamMemberGuids),
					}
					fmt.Fprintf(w,
						generateTeamMembersResponse(
							strings.Join(teamMembers, ", "),
						),
					)
				} else {
					fmt.Fprint(w, "{}")
				}

			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	t.Run("when the team member is found successfully", func(t *testing.T) {
		response, err := c.V2.TeamMembers.SearchUsername(username)
		assert.NoError(t, err)
		if assert.NotNil(t, response.Data[0]) {
			assert.Equal(t, teamMemberGuids[0], response.Data[0].UserGuid)
			assert.Equal(t, username, response.Data[0].UserName)
			assert.Equal(t, "Lacework", response.Data[0].Props.Company)
		}
	})
}

func generateTeamMembers(guids []string) string {
	tms := make([]string, len(guids))
	for i, guid := range guids {
		username := fmt.Sprintf("vatasha.white+%s@lacework.net", guid)
		tms[i] = singleMockTeamMember(guid, username)
	}
	return strings.Join(tms, ", ")
}

func generateTeamMembersWithTimeFieldsAsInts(guids []string) string {
	tms := make([]string, len(guids))
	for i, guid := range guids {
		username := fmt.Sprintf("vatasha.white+%s@lacework.net", guid)
		tms[i] = singleMockTeamMemberWithTimeFieldsAsInts(guid, username)
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

func singleMockTeamMember(id, username string) string {
	return fmt.Sprintf(`
    {
	  "custGuid": "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA",
      "props": {
            "accountAdmin": %t,
            "company": "Lacework",
            "createdTime": "2021-11-16T16:33:17.573Z",
            "firstName": "Vatasha",
            "jitCreated": %t,
            "lastLoginTime": "2021-11-17T16:33:17.573Z",
            "lastName": "White",
			"lastSessionCreatedTime": "2021-11-17T16:33:17.573Z",
			"orgAdmin": false,
			"orgUser": false,
			"updatedTime": "2021-11-20T16:33:17.573Z"
	  },
      "userEnabled": %d,
      "userGuid": %q,
      "userName": "%s"
    }
	`, false, false, 1, id, username)
}

func singleMockTeamMemberWithTimeFieldsAsInts(id, username string) string {
	return fmt.Sprintf(`
    {
	  "custGuid": "TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA",
      "props": {
            "accountAdmin": %t,
            "company": "Lacework",
            "createdTime": "2021-11-16T16:33:17.573Z",
            "firstName": "Vatasha",
            "jitCreated": %t,
            "lastLoginTime": %d,
            "lastName": "White",
			"lastSessionCreatedTime": %d,
			"orgAdmin": false,
			"orgUser": false,
			"updatedTime": %d
	  },
      "userEnabled": %d,
      "userGuid": %q,
      "userName": "%s"
    }
	`, false, false, 0, 0, 0, 1, id, username)
}

func singleMockOrgTeamMemberCreateResponse(username string) string {
	return fmt.Sprintf(`
	{
	  "data": {
		"accounts": [
		  {
			"accountName": "Testing123",
			"admin": false,
			"custGuid": "Testing123_24737A67D8599A42DA44B1E887BE527B97684BBF0102D0C",
			"userEnabled": 1,
			"userGuid": "Testing123_DD45B1F60AD668CF479BEDB22C977E09E1D6A6897C4C0E5"
		  }
		],
		"orgAccount": true,
		"orgAdmin": false,
		"orgUser": true,
		"url": "customerdemo.lacework.net",
		"username": "%s"
	  }
	}
`, username)
}

func singleMockTeamMembersOrgResponse(userGuid, userGuid2, username string) string {
	return fmt.Sprintf(`
	{
	  "data": [
		{
		  "custGuid": "CUSTOMER_721595854C4272A5AC58B9FAA369C0ABB05564A91DA0E00",
		  "props": {
			"accountAdmin": false,
			"company": "ABC",
			"createdTime": "2021-12-02T19:08:28.757Z",
			"firstName": "vatasha updated",
			"jitCreated": false,
			"lastLoginTime": 0,
			"lastName": "white updated",
			"lastSessionCreatedTime": "2021-12-02T19:08:30.079Z",
			"orgAdmin": false,
			"orgUser": true,
			"updatedBy": "vatasha.white@lacework.net",
			"updatedTime": "2021-12-02T19:39:34.746Z"
		  },
		  "userEnabled": 1,
		  "userGuid": %q,
		  "userName": %q
		},
		{
		  "custGuid": "GITOPS_53901D0F22F80387C484022798F21EE97F4C782FBE2E39A00",
		  "props": {
			"accountAdmin": false,
			"company": "ABC",
			"createdTime": "2021-12-02T19:08:28.935Z",
			"firstName": "vatasha updated",
			"jitCreated": false,
			"lastLoginTime": 0,
			"lastName": "white updated",
			"lastSessionCreatedTime": 0,
			"orgAdmin": false,
			"orgUser": true,
			"updatedBy": "vatasha.white@lacework.net",
			"updatedTime": "2021-12-02T19:39:34.746Z"
		  },
		  "userEnabled": 1,
		  "userGuid": %q,
		  "userName": %q
		}
	  ]
	}
	`, userGuid, username, userGuid2, username)
}
