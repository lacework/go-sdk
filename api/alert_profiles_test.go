//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestAlertProfilesGet(t *testing.T) {
	var (
		guid         = "LW_PROFILE_EXAMPLE"
		apiPath      = fmt.Sprintf("AlertProfiles/%s", guid)
		alertProfile = singleMockAlertProfile(guid)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateAlertProfileResponse(alertProfile))
			}
		},
	)

	fakeServer.MockAPI("AlertProfiles/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when alert profile exists", func(t *testing.T) {
		var response api.AlertProfileResponse
		err := c.V2.Alert.Profiles.Get(guid, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, "LW_PROFILE_EXAMPLE", response.Data.Guid)
			assert.Equal(t, "MID", response.Data.Fields[0].Name)
			assert.Equal(t, "MID", response.Data.DescriptionKeys[0].Name)
			assert.Equal(t, "{{MID}}", response.Data.DescriptionKeys[0].Spec)
			assert.Equal(t, "HE_File_Violation", response.Data.Alerts[0].Name)
		}
	})

	t.Run("when alert profile does NOT exist", func(t *testing.T) {
		var response api.AlertProfileResponse
		err := c.V2.Alert.Profiles.Get("UNKNOWN_INTG_GUID", response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/AlertProfiles/UNKNOWN_INTG_GUID")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestAlertProfilesDelete(t *testing.T) {
	var (
		guid         = "LW_PROFILE_EXAMPLE"
		apiPath      = fmt.Sprintf("AlertProfiles/%s", guid)
		alertProfile = singleMockAlertProfile(guid)
		getResponse  = generateAlertProfileResponse(alertProfile)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if getResponse != "" {
				switch r.Method {
				case "GET":
					fmt.Fprintf(w, getResponse)
				case "DELETE":
					// once deleted, empty the getResponse so that
					// further GET requests return 404s
					getResponse = ""
				}
			} else {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	fakeServer.MockAPI("AlertProfiles/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("verify alert profile exists", func(t *testing.T) {
		var response api.AlertProfileResponse
		err := c.V2.Alert.Profiles.Get(guid, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, "LW_PROFILE_EXAMPLE", response.Data.Guid)
			assert.Equal(t, "MID", response.Data.Fields[0].Name)
			assert.Equal(t, "MID", response.Data.DescriptionKeys[0].Name)
			assert.Equal(t, "{{MID}}", response.Data.DescriptionKeys[0].Spec)
			assert.Equal(t, "HE_File_Violation", response.Data.Alerts[0].Name)
		}
	})

	t.Run("when alert profile has been deleted", func(t *testing.T) {
		err := c.V2.Alert.Profiles.Delete(guid)
		assert.Nil(t, err)

		var response api.AlertProfileResponse
		err = c.V2.Alert.Profiles.Get(guid, &response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/AlertProfiles/LW_")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestAlertProfilesList(t *testing.T) {
	var (
		allGUIDs      []string
		alertProfiles = generateGuids(&allGUIDs, 3)
		expectedLen   = len(allGUIDs)
		fakeServer    = lacework.MockServer()
	)

	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("AlertProfiles",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			alertProfiles := []string{
				generateAlertProfiles(alertProfiles),
			}
			fmt.Fprintf(w,
				generateAlertProfilesResponse(
					strings.Join(alertProfiles, ", "),
				),
			)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Alert.Profiles.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGUIDs, d.Guid)
	}
}

func TestAlertProfileUpdate(t *testing.T) {
	var (
		guid       = "LW_PROFILE_EXAMPLE"
		apiPath    = fmt.Sprintf("AlertProfiles/%s", guid)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "Update() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "HE_File_Violation", "alert profile alerts are missing")
		}

		fmt.Fprintf(w, generateAlertProfileResponse(singleMockAlertProfile(guid)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	alertProfile := api.AlertProfile{
		Alerts: []api.AlertTemplate{{Name: "HE_File_Violation",
			EventName:   "LW Host Entity File Violation Alert",
			Description: "{{_OCCURRENCE}} Violation for file {{PATH}} on machine {{MID}}",
			Subject:     "{{_OCCURRENCE}} violation detected for file {{PATH}} on machine {{MID}}"},
		},
	}
	assert.Equal(t, "HE_File_Violation", alertProfile.Alerts[0].Name, "an alert profile alerts name should match")
	assert.Equal(t, "LW Host Entity File Violation Alert", alertProfile.Alerts[0].EventName, "an alert profile event name should match")
	assert.Equal(t, "{{_OCCURRENCE}} Violation for file {{PATH}} on machine {{MID}}", alertProfile.Alerts[0].Description, "an alert profile description should match")
	assert.Equal(t, "{{_OCCURRENCE}} violation detected for file {{PATH}} on machine {{MID}}", alertProfile.Alerts[0].Subject, "an alert profile subject should match")

	response, err := c.V2.Alert.Profiles.Update("LW_PROFILE_EXAMPLE", alertProfile.Alerts)
	if assert.NoError(t, err) {
		assert.NotNil(t, response)
		assert.Equal(t, "LW_PROFILE_EXAMPLE", response.Data.Guid)
		assert.Equal(t, "MID", response.Data.Fields[0].Name)
		assert.Equal(t, "MID", response.Data.DescriptionKeys[0].Name)
		assert.Equal(t, "{{MID}}", response.Data.DescriptionKeys[0].Spec)
		assert.Equal(t, "HE_File_Violation", response.Data.Alerts[0].Name)
	}
}

func generateAlertProfiles(guids []string) string {
	alertProfiles := make([]string, len(guids))
	for i, guid := range guids {
		alertProfiles[i] = singleMockAlertProfile(guid)
	}
	return strings.Join(alertProfiles, ", ")
}

func generateAlertProfilesResponse(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func generateAlertProfileResponse(data string) string {
	return `
		{
			"data": ` + data + `
		}
	`
}

func singleMockAlertProfile(id string) string {
	return fmt.Sprintf(`
{
        "alertProfileId": %q,
        "extends": "LW_TEST",
        "fields": [
            {
                "name": "MID"
            },
            {
                "name": "PATH"
            },
            {
                "name": "_OCCURRENCE"
            }
        ],
        "descriptionKeys": [
            {
                "name": "MID",
                "spec": "{{MID}}"
            },
            {
                "name": "_OCCURRENCE",
                "spec": "{{_OCCURRENCE}}"
            },
            {
                "name": "PATH",
                "spec": "{{PATH}}"
            }
        ],
        "alerts": [
            {
                "name": "HE_File_Violation",
                "eventName": "LW Host Entity File Violation Alert",
                "description": "{{_OCCURRENCE}} Violation for file {{PATH}} on machine {{MID}}",
                "subject": "{{_OCCURRENCE}} violation detected for file {{PATH}} on machine {{MID}}"
            }
        ]
}
	`, id)
}
