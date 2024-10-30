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

	groupFound, found = api.FindResourceGroupType("KUBERNETES")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "KUBERNETES", groupFound.String(), "wrong resource group type")
}

func TestResourceGroupGet(t *testing.T) {
	var (
		resourceGUID = intgguid.New()
		vanillaType  = "VANILLA"
		apiPath      = fmt.Sprintf("ResourceGroups/%s", resourceGUID)
		vanillaGroup = singleVanillaResourceGroup(resourceGUID, vanillaType, "")
		fakeServer   = lacework.MockServer()
	)
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
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("when resource group exists", func(t *testing.T) {
		var response api.ResourceGroupResponse
		err := c.V2.ResourceGroups.Get(resourceGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, resourceGUID, response.Data.ResourceGroupGuid)
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
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	t.Run("verify resource group exists", func(t *testing.T) {
		var response api.ResourceGroupResponse
		err := c.V2.ResourceGroups.Get(resourceGUID, &response)
		assert.Nil(t, err)
		if assert.NotNil(t, response) {
			assert.Equal(t, resourceGUID, response.Data.ResourceGroupGuid)
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
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.ResourceGroups.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, allGuids, d.ResourceGroupGuid)
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
		case api.MachineResourceGroup.String():
			resourceGroups[i] = singleMachineResourceGroup(guid)
		}
	}
	return strings.Join(resourceGroups, ", ")
}

func singleAwsResourceGroup(id string) string {
	return `
	{
        "guid": "` + id + `",
		"description": "All Aws Resources",
        "isDefaultBoolean": true,
        "resourceGroupGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "AWS",
        "enabled": 1,
		"query": {
			  "filters": {
				  "filter1": {
					  "field": "AWS_ACCOUNT_ID",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter2": {
					  "field": "AWS_ORGANIZATION_ID",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter3": {
					  "field": "AWS_RESOURCE_TAGS",
					  "operation": "EQUALS",
					  "key": "*",
					  "values": [
						  "*"
					  ]
				  },
				  "filter4": {
					  "field": "AWS_RESOURCE_REGION",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  }
			  },
			  "expression": {
				  "operator": "OR",
				  "children": [
					  {
						  "filterName": "filter1"
					  },
					  {
						  "filterName": "filter2"
					  },
					  {
						  "filterName": "filter3"
					  },
					  {
						  "filterName": "filter4"
					  }
				  ]
			  }
		  }
	}
	`
}

func singleAzureResourceGroup(id string) string {

	return `
	{
        "guid": "` + id + `",
		"description": "All Azure Resources",
        "isDefaultBoolean": true,
        "resourceGroupGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "AZURE",
        "enabled": 1,
		"query": {
			  "filters": {
				  "filter1": {
					  "field": "AZURE_TENANT_ID",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter2": {
					  "field": "AZURE_SUBSCRIPTION_ID",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter3": {
					  "field": "AZURE_RESOURCE_TAGS",
					  "operation": "EQUALS",
					  "key": "*",
					  "values": [
						  "*"
					  ]
				  },
				  "filter4": {
					  "field": "AZURE_RESOURCE_REGION",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter5": {
					  "field": "AZURE_TENANT_NAME",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter6": {
					  "field": "AZURE_SUBSCRIPTION_NAME",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  }
			  },
			  "expression": {
				  "operator": "OR",
				  "children": [
					  {
						  "filterName": "filter1"
					  },
					  {
						  "filterName": "filter2"
					  },
					  {
						  "filterName": "filter3"
					  },
					  {
						  "filterName": "filter4"
					  },
					  {
						  "filterName": "filter5"
					  },
					  {
						  "filterName": "filter6"
					  }
				  ]
			  }
		  }
	}
	`
}

func singleContainerResourceGroup(id string) string {
	return `
	{
        "guid": "` + id + `",
		"description": "All Container Resources",
        "isDefaultBoolean": true,
        "resourceGroupGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "CONTAINER",
        "enabled": 1,
		"query": {
			  "filters": {
				  "filter1": {
					  "field": "IMAGE_TAG",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter2": {
					  "field": "CONTAINER_LABELS",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ],
					  "key": "*"
				  },
				  "filter3": {
					  "field": "IMAGE_REPO",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter4": {
					  "field": "IMAGE_REGISTRY",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				}
			  },
			  "expression": {
				  "operator": "OR",
				  "children": [
					  {
						  "filterName": "filter1"
					  },
					  {
						  "filterName": "filter2"
					  },
					  {
						  "filterName": "filter3"
					  },
					  {
						  "filterName": "filter4"
					  }
				  ]
			  }
		  }
	}
	`
}

func singleGcpResourceGroup(id string) string {
	return `
	{
        "guid": "` + id + `",
		"description": "All GCP Resources",
        "isDefaultBoolean": true,
        "resourceGroupGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "GCP",
        "enabled": 1,
		"query": {
                          "filters": {
                              "filter1": {
                                  "field": "GCP_ORGANIZATION_ID",
                                  "operation": "EQUALS",
                                  "values": [
                                      "*"
                                  ]
                              },
                              "filter2": {
                                  "field": "GCP_FOLDER_IDS",
                                  "operation": "EQUALS",
                                  "values": [
                                      "*"
                                  ]
                              },
                              "filter3": {
                                  "field": "GCP_PROJECT_ID",
                                  "operation": "EQUALS",
                                  "values": [
                                      "*"
                                  ]
                              },
                              "filter4": {
                                  "field": "GCP_RESOURCE_TAGS",
                                  "operation": "EQUALS",
                                  "key": "*",
                                  "values": [
                                      "*"
                                  ]
                              },
                              "filter5": {
                                  "field": "GCP_RESOURCE_REGION",
                                  "operation": "EQUALS",
                                  "values": [
                                      "*"
                                  ]
                              },
                              "filter6": {
                                  "field": "GCP_ORGANIZATION_NAME",
                                  "operation": "EQUALS",
                                  "values": [
                                      "*"
                                  ]
                              }
                          },
                          "expression": {
                              "operator": "OR",
                              "children": [
                                  {
                                      "filterName": "filter1"
                                  },
                                  {
                                      "filterName": "filter2"
                                  },
                                  {
                                      "filterName": "filter3"
                                  },
                                  {
                                      "filterName": "filter4"
                                  },
                                  {
                                      "filterName": "filter5"
                                  },
                                  {
                                      "filterName": "filter5"
                                  }
                              ]
                          }
                      }
	}
	`
}

func singleMachineResourceGroup(id string) string {
	return `
	{
        "guid": "` + id + `",
		"description": "All Machine Resources",
        "isDefaultBoolean": true,
        "resourceGroupGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "MACHINE",
        "enabled": 1,
		"query": {
                          "filters": {
                              "filter1": {
                                  "field": "MACHINE_TAGS",
                                  "operation": "EQUALS",
                                  "values": [
                                      "*"
                                  ],
                                  "key": "*"
                              }
                          },
                          "expression": {
                              "filterName": "filter1"
                          }
                      }
	}
	`
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

func singleVanillaResourceGroup(id string, iType string, query string) string {
	if query == "" {
		query = "{}"
	}
	return `
	{
        "isDefaultBoolean": true,
        "query": ` + query + `,
        "resourceGroupGuid": "` + id + `",
        "name": "group_name",
        "resourceType": "` + iType + `",
        "enabled": 1
	}
	`
}
