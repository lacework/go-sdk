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

func TestResourceGroupAzureGet(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetAzureResourceGroup() should be a GET method")
		fmt.Fprintf(w, generateResourceGroupResponse(singleAzureResourceGroup(resourceGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ResourceGroups.GetAzure(resourceGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
	assert.Equal(t, "group_name", response.Data.Name)
	assert.Equal(t, "All Tenants and Subscriptions", response.Data.Props.Description)
	assert.Equal(t, []string{"*"}, response.Data.Props.Subscriptions)
	assert.Equal(t, "*", response.Data.Props.Tenant)
}

func TestResourceGroupsAzureUpdate(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateAzureResourceGroup() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "group_name", "Resource Group  name is missing")
			assert.Contains(t, body, "AZURE", "wrong Resource Group  type")
			assert.Contains(t, body, "Updated", "wrong description")
			assert.Contains(t, body, "[\"abc123\",\"cba321\"]", "wrong subscriptions")
		}

		fmt.Fprintf(w, generateResourceGroupResponse(singleAzureResourceGroupUpdateResponse(resourceGUID)))
	})

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	resourceGroup := api.NewResourceGroup("group_name",
		api.AzureResourceGroup,
		api.AzureResourceGroupProps{
			Description:   "Updated",
			Subscriptions: []string{"abc123", "cba321"},
			Tenant:        "tenant",
		},
	)
	assert.Equal(t, "group_name", resourceGroup.Name, "Azure Resource Group name mismatch")
	assert.Equal(t, "AZURE", resourceGroup.Type, "a new Azure Resource Group should match its type")
	assert.Equal(t, 1, resourceGroup.Enabled, "a new Azure Resource Group should be enabled")
	resourceGroup.ResourceGuid = resourceGUID

	response, err := c.V2.ResourceGroups.UpdateAzure(&resourceGroup)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
}

func singleAzureResourceGroup(id string) string {
	return `
	{
        "guid": "` + id + `",
        "isDefault": "1",
        "props": "{\"DESCRIPTION\":\"All Tenants and Subscriptions\",\"SUBSCRIPTIONS\":[\"*\"],\"TENANT\":\"*\",\"UPDATED_BY\":null,\"LAST_UPDATED\":1586453993500}",
        "resourceGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "AZURE",
        "enabled": 1
	}
	`
}

func singleAzureResourceGroupUpdateResponse(id string) string {
	return `
	{
        "guid": "` + id + `",
        "isDefault": 1,
        "props": {
			"description":"All Tenants and Subscriptions",
			"subscriptions":["*"],
			"tenant":"*",
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
