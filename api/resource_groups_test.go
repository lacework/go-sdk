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
	"strings"
	"testing"

	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestResourceGroupTypes(t *testing.T) {
	assert.Equal(t, "AWS", api.AwsResourceGroup.String(), "wrong resource group type")
	assert.Equal(t, "AZURE", api.AzureResourceGroup.String(), "wrong resource group type")
	assert.Equal(t, "CONTAINER", api.ContainerResourceGroup.String(), "wrong resource group type")
	assert.Equal(t, "GCP", api.GcpResourceGroup.String(), "wrong resource group type")
	assert.Equal(t, "MACHINE", api.MachineResourceGroup.String(), "wrong resource group type")
	assert.Equal(t, "LW_ACCOUNT", api.LwAccountResourceGroup.String(), "wrong resource group type")
}

func TestFindResourceGroupType(t *testing.T) {
	groupFound, found := api.FindResourceGroupType("NONE_EXISTING_GROUP")
	assert.False(t, found, "resource group type should not be found")
	assert.Equal(t, 0, int(groupFound), "wrong resource group type")
	assert.Equal(t, "None", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("AWS")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "AWS", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("AWS")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "AWS", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("CONTAINER")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "CONTAINER", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("GCP")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "GCP", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("MACHINE")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "MACHINE", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("LW_ACCOUNT")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "LW_ACCOUNT", groupFound.String(), "wrong resource group type")
}

func TestResourceGroupGet(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		vanillaType  = "VANILLA"
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		vanillaGroup = singleVanillaResourceGroup(resourceGUID, vanillaType, "")
		fakeServer   = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				fmt.Fprintf(w, generateResourceGroupResponse(vanillaGroup))
			}
		},
	)

	fakeServer.MockAPI("ResourceGroups/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when resource group exists", func(t *testing.T) {
		var response api.ResourceGroupResponse
		err := c.V2.ResourceGroups.Get(resourceGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
			assert.Equal(t, "group_name", response.Data.Name)
			assert.Equal(t, "VANILLA", response.Data.Type)
		}
	})

	t.Run("when resource group does NOT exist", func(t *testing.T) {
		var response api.ResourceGroupResponse
		err := c.V2.ResourceGroups.Get("UNKNOWN_INTG_GUID", response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/ResourceGroups/UNKNOWN_INTG_GUID")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestResourceGroupsDelete(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		vanillaType  = "VANILLA"
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		vanillaInt   = singleVanillaResourceGroup(resourceGUID, vanillaType, "")
		getResponse  = generateResourceGroupResponse(vanillaInt)
		fakeServer   = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath,
		func(w http.ResponseWriter, r *http.Request) {
			if getResponse != "" {
				switch r.Method {
				case "GET":
					fmt.Fprintf(w, getResponse)
				case "DELETE":
					getResponse = ""
				}
			} else {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	fakeServer.MockAPI("ResourceGroups/UNKNOWN_INTG_GUID",
		func(w http.ResponseWriter, r *http.Request) {
			if assert.Equal(t, "GET", r.Method, "Get() should be a GET method") {
				http.Error(w, "{ \"message\": \"Not Found\"}", 404)
			}
		},
	)

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("verify resource group exists", func(t *testing.T) {
		var response api.ResourceGroupResponse
		err := c.V2.ResourceGroups.Get(resourceGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, resourceGUID, response.Data.ResourceGuid)
			assert.Equal(t, "group_name", response.Data.Name)
			assert.Equal(t, "VANILLA", response.Data.Type)
		}
	})

	t.Run("when resource group has been deleted", func(t *testing.T) {
		err := c.V2.ResourceGroups.Delete(resourceGUID)
		assert.Nil(t, err)

		var response api.ResourceGroupResponse
		err = c.V2.ResourceGroups.Get(resourceGUID, response)
		assert.Empty(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "api/v2/ResourceGroups/MOCK_")
			assert.Contains(t, err.Error(), "[404] Not Found")
		}
	})
}

func TestResourceGroupsList(t *testing.T) {
	var (
		awsResourceGUIDs       = []string{intgguid.New(), intgguid.New()}
		azureResourceGUIDs     = []string{intgguid.New(), intgguid.New()}
		containerResourceGUIDs = []string{intgguid.New()}
		gcpResourceGUIDs       = []string{intgguid.New()}
		machineResourceGUIDs   = []string{intgguid.New()}
		allGroups              = [][]string{awsResourceGUIDs, azureResourceGUIDs, containerResourceGUIDs, gcpResourceGUIDs, machineResourceGUIDs}
		allGuids               []string
		fakeServer             = lacework.MockServer()
	)

	for _, guids := range allGroups {
		allGuids = append(allGuids, guids...)
	}
	expectedLen := len(allGuids)

	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("ResourceGroups",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			ResourceGroups := []string{
				generateResourceGroups(awsResourceGUIDs, "AWS"),
				generateResourceGroups(azureResourceGUIDs, "AZURE"),
				generateResourceGroups(containerResourceGUIDs, "CONTAINER"),
				generateResourceGroups(gcpResourceGUIDs, "GCP"),
				generateResourceGroups(machineResourceGUIDs, "MACHINE"),
			}
			fmt.Fprintf(w,
				generateResourceGroupsResponse(
					strings.Join(ResourceGroups, ", "),
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

	response, err := c.V2.ResourceGroups.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGuids, d.ResourceGuid)
	}
}

func generateResourceGroups(guids []string, iType string) string {
	resourceGroups := make([]string, len(guids))
	for i, guid := range guids {
		switch iType {
		case api.AwsResourceGroup.String():
			resourceGroups[i] = singleAwsResourceGroup(guid)
		case api.AzureResourceGroup.String():
			resourceGroups[i] = singleAzureResourceGroup(guid)
		case api.ContainerResourceGroup.String():
			resourceGroups[i] = singleContainerResourceGroup(guid)
		case api.GcpResourceGroup.String():
			resourceGroups[i] = singleGcpResourceGroup(guid)
		case api.LwAccountResourceGroup.String():
			resourceGroups[i] = singleLwAccountResourceGroup(guid)
		case api.MachineResourceGroup.String():
			resourceGroups[i] = singleMachineResourceGroup(guid)
		}
	}
	return strings.Join(resourceGroups, ", ")
}

func generateResourceGroupsResponse(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func generateResourceGroupResponse(data string) string {
	return `
		{
			"data": ` + data + `
		}
	`
}

func singleVanillaResourceGroup(id string, iType string, props string) string {
	if props == "" {
		props = "{}"
	}
	return `
	{
        "guid": "` + id + `",
        "isDefault": "1",
        "props": ` + props + `,
        "resourceGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "` + iType + `",
        "enabled": 1
	}
	`
}
