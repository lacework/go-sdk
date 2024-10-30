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
	assert.Equal(t, "OCI", api.OciResourceGroup.String(), "wrong resource group type")
	assert.Equal(t, "KUBERNETES", api.KubernetesResourceGroup.String(), "wrong resource group type")
}

func TestFindResourceGroupType(t *testing.T) {
	groupFound, found := api.FindResourceGroupType("NONE_EXISTING_GROUP")
	assert.False(t, found, "resource group type should not be found")
	assert.Equal(t, 0, int(groupFound), "wrong resource group type")
	assert.Equal(t, "None", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("AWS")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "AWS", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("GCP")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "GCP", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("AZURE")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "AZURE", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("CONTAINER")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "CONTAINER", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("MACHINE")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "MACHINE", groupFound.String(), "wrong resource group type")

	groupFound, found = api.FindResourceGroupType("OCI")
	assert.True(t, found, "resource group type should exist")
	assert.Equal(t, "OCI", groupFound.String(), "wrong resource group type")

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
		awsResourceGUIDs        = []string{intgguid.New(), intgguid.New()}
		azureResourceGUIDs      = []string{intgguid.New(), intgguid.New()}
		containerResourceGUIDs  = []string{intgguid.New()}
		gcpResourceGUIDs        = []string{intgguid.New()}
		machineResourceGUIDs    = []string{intgguid.New()}
		ociResourceGUIDs        = []string{intgguid.New()}
		kubernetesResourceGUIDs = []string{intgguid.New()}
		allGroups               = [][]string{awsResourceGUIDs, azureResourceGUIDs, containerResourceGUIDs,
			gcpResourceGUIDs, machineResourceGUIDs, ociResourceGUIDs, kubernetesResourceGUIDs}
		allGuids   []string
		fakeServer = lacework.MockServer()
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
				generateResourceGroups(ociResourceGUIDs, "OCI"),
				generateResourceGroups(kubernetesResourceGUIDs, "KUBERNETES"),
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
		case api.OciResourceGroup.String():
			resourceGroups[i] = singleOciResourceGroup(guid)
		case api.KubernetesResourceGroup.String():
			resourceGroups[i] = singleKubernetesResourceGroup(guid)

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
					  "field": "Account",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter2": {
					  "field": "Organization ID",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter3": {
					  "field": "Resource Tag",
					  "operation": "EQUALS",
					  "key": "*",
					  "values": [
						  "*"
					  ]
				  },
				  "filter4": {
					  "field": "Region",
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
					  "field": "Tenant ID",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter2": {
					  "field": "Subscription ID",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter3": {
					  "field": "Resource Tag",
					  "operation": "EQUALS",
					  "key": "*",
					  "values": [
						  "*"
					  ]
				  },
				  "filter4": {
					  "field": "Region",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter5": {
					  "field": "Tenant Name",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter6": {
					  "field": "Subscription Name",
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
					  "field": "Container Tag",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ],
					  "key": "*"
				  },
				  "filter2": {
					  "field": "Container Label",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ],
					  "key": "*"
				  },
				  "filter3": {
					  "field": "Image Repo",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter4": {
					  "field": "Image Registry",
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
					  "field": "Organization ID",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter2": {
					  "field": "Folder",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter3": {
					  "field": "Project ID",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter4": {
					  "field": "Resource Label",
					  "operation": "EQUALS",
					  "key": "*",
					  "values": [
						  "*"
					  ]
				  },
				  "filter5": {
					  "field": "Region",
					  "operation": "EQUALS",
					  "values": [
						  "*"
					  ]
				  },
				  "filter6": {
					  "field": "Organization Name",
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
					  "field": "Machine Tag",
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

func singleOciResourceGroup(id string) string {
	return `
	{
        "guid": "` + id + `",
		"description": "All OCI Resources",
        "isDefaultBoolean": true,
        "resourceGroupGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "OCI",
        "enabled": 1,
		"query": {
			 "filters": {
				 "filter1": {
					 "field": "Compartment ID",
					 "operation": "EQUALS",
					 "values": [
						 "*"
					 ]
				 },
				 "filter2": {
					 "field": "Compartment Name",
					 "operation": "EQUALS",
					 "values": [
						 "*"
					 ]
				 },
				 "filter3": {
					 "field": "Resource Tag",
					 "operation": "EQUALS",
					 "key": "*",
					 "values": [
						 "*"
					 ]
				 },
				 "filter4": {
					 "field": "Region",
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

func singleKubernetesResourceGroup(id string) string {

	return `
	{
        "guid": "` + id + `",
		"description": "All Kubernetes Resources",
        "isDefaultBoolean": true,
        "resourceGroupGuid": "` + id + `",
        "resourceName": "group_name",
        "resourceType": "KUBERNETES",
        "enabled": 1,
		"query": {
			 "filters": {
				 "filter1": {
					 "field": "AWS Account",
					 "operation": "EQUALS",
					 "values": [
						 "*"
					 ]
				 },
				 "filter2": {
					 "field": "AWS Region",
					 "operation": "EQUALS",
					 "values": [
						 "*"
					 ]
				 },
				 "filter3": {
					 "field": "Cluster Name",
					 "operation": "EQUALS",
					 "values": [
						 "*"
					 ]
				 },
				 "filter4": {
					 "field": "Namespace",
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
