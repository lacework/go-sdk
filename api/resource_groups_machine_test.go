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

func TestResourceGroupMachineGet(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetMachineResourceGroup() should be a GET method")
		fmt.Fprintf(w, generateResourceGroupResponse(singleMachineResourceGroup(resourceGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ResourceGroups.GetMachineResourceGroup(resourceGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
	assert.Equal(t, "group_name", response.Data.Name)
	assert.Equal(t, "All Machine Tags", response.Data.GetProps().Description)
	assert.Equal(t, []map[string]string{{"*": "*"}}, response.Data.GetProps().MachineTags)
}

func TestResourceGroupsMachineUpdate(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateMachineResourceGroup() should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, resourceGUID, "ResourceGUID missing")
			assert.Contains(t, body, "group_name", "Resource Group  name is missing")
			assert.Contains(t, body, "MACHINE", "wrong Resource Group  type")
			assert.Contains(t, body, "Updated", "wrong description")
			assert.Contains(t, body, "[{\"tag\":\"machineTag\"}]", "wrong machine tags")
		}

		fmt.Fprintf(w, generateResourceGroupResponse(singleMachineResourceGroup(resourceGUID)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	resourceGroup := api.NewResourceGroup("group_name",
		api.MachineResourceGroup,
		api.MachineResourceGroupData{
			Props: api.MachineResourceGroupProps{
				Description: "Updated",
				MachineTags: []map[string]string{{"tag": "machineTag"}},
			},
		},
	)
	assert.Equal(t, "group_name", resourceGroup.Name, "Machine Resource Group name mismatch")
	assert.Equal(t, "MACHINE", resourceGroup.Type, "a new Machine Resource Group should match its type")
	assert.Equal(t, 1, resourceGroup.Enabled, "a new Machine Resource Group should be enabled")
	resourceGroup.ResourceGuid = resourceGUID

	response, err := c.V2.ResourceGroups.UpdateMachineResourceGroup(resourceGroup)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
}

func singleMachineResourceGroup(id string) string {
	return `
	{
        "guid": "` + id + `",
        "isDefault": "1",
        "props": "{\"DESCRIPTION\":\"All Machine Tags\",\"MACHINE_TAGS\":[{\"*\":\"*\"}],\"UPDATED_BY\":null,\"LAST_UPDATED\":1586453993565}",
        "resourceGuid": "` + id + `",
        "resourceName": "group_name",
				"resourceType": "MACHINE",
        "enabled": 1
	}
	`
}
