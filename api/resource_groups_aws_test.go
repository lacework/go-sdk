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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestResourceGroupAwsGet(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAwsResourceGroup() should be a GET method")
		fmt.Fprintf(w, generateResourceGroupResponse(singleAwsResourceGroup(resourceGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ResourceGroups.GetAws(resourceGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
	assert.Equal(t, "group_name", response.Data.Name)
	assert.Equal(t, "All Aws Accounts", response.Data.Props.Description)
	assert.Equal(t, []string{"*"}, response.Data.Props.AccountIDs)
}

func TestResourceGroupsAwsUpdate(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAwsResourceGroup() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "group_name", "Resource Group  name is missing")
			assert.Contains(t, body, "AWS", "wrong Resource Group  type")
			assert.Contains(t, body, "Updated", "wrong description")
			assert.Contains(t, body, "[\"abc123\",\"cba321\"]", "wrong account ids")
		}

		fmt.Fprintf(w, generateResourceGroupResponse(singleAwsResourceGroupUpdateResponse(resourceGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	resourceGroup := api.NewResourceGroup("group_name",
		api.AwsResourceGroup,
		api.AwsResourceGroupProps{
			Description: "Updated",
			AccountIDs:  []string{"abc123", "cba321"},
		},
	)
	assert.Equal(t, "group_name", resourceGroup.Name, "Aws Resource Group name mismatch")
	assert.Equal(t, "AWS", resourceGroup.Type, "a new Aws Resource Group should match its type")
	assert.Equal(t, 1, resourceGroup.Enabled, "a new Aws Resource Group should be enabled")
	resourceGroup.ResourceGuid = resourceGUID

	response, err := c.V2.ResourceGroups.UpdateAws(&resourceGroup)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
}

func TestMarshallResourceGroupLastUpdatedTime(t *testing.T) {
	var res api.AwsResourceGroupResponse
	err := json.Unmarshal([]byte(generateResourceGroupResponse(singleAwsResourceGroupUpdateResponse("test"))), &res)
	if err != nil {
		log.Fatal("Unable to unmarshall aws resource group string")
	}
	jsonString, err := json.Marshal(res)
	if err != nil {
		log.Fatal("Unable to marshall aws resource group string")
	}

	assert.Equal(t, res.Data.Props.LastUpdated.ToTime().UnixNano()/int64(time.Millisecond), int64(1586453993470))
	assert.Contains(t, string(jsonString), "2020-04-09T17:39:53Z")
}

func singleAwsResourceGroup(id string) string {
	return `
	{
        "guid": "` + id + `",
        "isDefault": "1",
        "props": "{\"DESCRIPTION\":\"All Aws Accounts\",\"ACCOUNT_IDS\":[\"*\"],\"UPDATED_BY\":null,\"LAST_UPDATED\":1586453993470}",
        "resourceGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "AWS",
        "enabled": 1
	}
	`
}

func singleAwsResourceGroupUpdateResponse(id string) string {
	return `
	{
        "guid": "` + id + `",
        "isDefault": 1,
        "props": {
			"description":"All Aws Accounts",
			"accountIds":["*"],
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
