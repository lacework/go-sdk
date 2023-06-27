//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestResourceGroupLwAccountGet(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetLwAccountResourceGroup() should be a GET method")
		fmt.Fprintf(w, generateResourceGroupResponse(singleLwAccountResourceGroup(resourceGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ResourceGroups.GetLwAccount(resourceGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
	assert.Equal(t, "group_name", response.Data.Name)
	assert.Equal(t, "All Lacework Accounts", response.Data.Props.Description)
	assert.Equal(t, []string{"*"}, response.Data.Props.LwAccounts)
}

func TestResourceGroupsLwAccountUpdate(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateLwAccountResourceGroup() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "group_name", "Resource Group name is missing")
			assert.Contains(t, body, "LW_ACCOUNT", "wrong Resource Group type")
			assert.Contains(t, body, "Updated", "wrong description")
			assert.Contains(t, body, "[\"abc123\",\"cba321\"]", "wrong lw accounts ids")
		}

		fmt.Fprintf(w, generateResourceGroupResponse(singleLwAccountResourceGroupUpdateResponse(resourceGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	resourceGroup := api.NewResourceGroup("group_name",
		api.LwAccountResourceGroup,
		api.LwAccountResourceGroupProps{
			Description: "Updated",
			LwAccounts:  []string{"abc123", "cba321"},
		},
	)
	assert.Equal(t, "group_name", resourceGroup.Name, "LwAccount Resource Group name mismatch")
	assert.Equal(t, "LW_ACCOUNT", resourceGroup.Type, "a new LwAccount Resource Group should match its type")
	assert.Equal(t, 1, resourceGroup.Enabled, "a new LwAccount Resource Group should be enabled")
	resourceGroup.ResourceGuid = resourceGUID

	response, err := c.V2.ResourceGroups.UpdateLwAccount(&resourceGroup)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
}

func singleLwAccountResourceGroup(id string) string {
	return `
	{
        "guid": "` + id + `",
        "isDefault": "1",
        "props": "{\"DESCRIPTION\":\"All Lacework Accounts\",\"LW_ACCOUNTS\":[\"*\"],\"UPDATED_BY\":null,\"LAST_UPDATED\":1582066544260}",
        "resourceGuid": "` + id + `",
        "resourceName": "group_name",
			"resourceType": "lW_ACCOUNT",
        "enabled": 1
	}
	`
}

func singleLwAccountResourceGroupUpdateResponse(id string) string {
	return `
	{
        "guid": "` + id + `",
        "isDefault": 1,
        "props": {
			"description":"All Lacework Accounts",
			"lwAccounts":["*"],
			"updatedBy":null,
			"lastUpdated":1586453993470 
			},
        "resourceGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "AWS",
        "enabled": 1
	}
	`
}
