//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestEntities_MachineDetails_Search(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/MachineDetails/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockMachineDetailsResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response := api.MachineDetailsEntityResponse{}
	err = c.V2.Entities.Search(&response, api.SearchFilter{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "mock-hostname", response.Data[0].Hostname)
	}
}

func TestEntities_MachineDetails_List(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/MachineDetails/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockPaginationResponsePage1())
		},
	)
	fakeServer.MockAPI("NextPage/vuln/containers/abc123",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "NextPage() should be a GET method")
			fmt.Fprintf(w, mockPaginationResponsePage2())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Entities.ListMachineDetails()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	// only one since all pages would be two
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "mock-1-hostname", response.Data[0].Hostname)
		assert.Equal(t,
			"https://account.lacework.net/api/v2/NextPage/vuln/containers/abc123",
			response.Paging.Urls.NextPage)
	}
}

func TestEntities_MachineDetails_List_WithFilters(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/MachineDetails/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockPaginationResponsePage1())
		},
	)
	fakeServer.MockAPI("NextPage/vuln/containers/abc123",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "NextPage() should be a GET method")
			fmt.Fprintf(w, mockPaginationResponsePage2())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	var (
		now    = time.Now().UTC()
		before = now.AddDate(0, 0, -7) // last 7 days
	)

	entityFilters := api.SearchFilter{
		TimeFilter: &api.TimeFilter{
			StartTime: &before,
			EndTime:   &now,
		},
	}

	response, err := c.V2.Entities.ListMachineDetailsWithFilters(entityFilters)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	// only one since all pages would be two
	if assert.Equal(t, 1, len(response.Data)) {
		assert.Equal(t, "mock-1-hostname", response.Data[0].Hostname)
		assert.Equal(t,
			"https://account.lacework.net/api/v2/NextPage/vuln/containers/abc123",
			response.Paging.Urls.NextPage)
	}
}

func TestEntities_MachineDetails_List_All(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/MachineDetails/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockPaginationResponsePage1())
		},
	)
	fakeServer.MockAPI("NextPage/vuln/containers/abc123",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "NextPage() should be a GET method")
			fmt.Fprintf(w, mockPaginationResponsePage2())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Entities.ListAllMachineDetails()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 2, len(response.Data)) {
		assert.Equal(t, "mock-1-hostname", response.Data[0].Hostname)
		assert.Equal(t, "mock-2-hostname", response.Data[1].Hostname)
		assert.Empty(t, response.Paging, "paging should be empty")
	}
}

func TestEntities_MachineDetails_List_All_WithFilters(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/MachineDetails/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockPaginationResponsePage1())
		},
	)
	fakeServer.MockAPI("NextPage/vuln/containers/abc123",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "NextPage() should be a GET method")
			fmt.Fprintf(w, mockPaginationResponsePage2())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	var (
		now    = time.Now().UTC()
		before = now.AddDate(0, 0, -7) // last 7 days
	)

	entityFilters := api.SearchFilter{
		TimeFilter: &api.TimeFilter{
			StartTime: &before,
			EndTime:   &now,
		},
	}

	response, err := c.V2.Entities.ListAllMachineDetailsWithFilters(entityFilters)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 2, len(response.Data)) {
		assert.Equal(t, "mock-1-hostname", response.Data[0].Hostname)
		assert.Equal(t, "mock-2-hostname", response.Data[1].Hostname)
		assert.Empty(t, response.Paging, "paging should be empty")
	}
}

func TestEntities_MachineDetails_ListAll_EmptyData(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/MachineDetails/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockPaginationEmptyResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Entities.ListAllMachineDetails()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, len(response.Data))
}

func TestEntities_MachineDetails_ListAll_WithFilters_EmptyData(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/MachineDetails/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockPaginationEmptyResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	var (
		now    = time.Now().UTC()
		before = now.AddDate(0, 0, -7) // last 7 days
	)

	entityFilters := api.SearchFilter{
		TimeFilter: &api.TimeFilter{
			StartTime: &before,
			EndTime:   &now,
		},
	}

	response, err := c.V2.Entities.ListAllMachineDetailsWithFilters(entityFilters)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, len(response.Data))
}

func mockMachineDetailsResponse() string {
	return `
{
  "data": [
    {
      "awsInstanceId": "i-abc12345678901234",
      "awsZone": "us-west-2",
      "createdTime": "2022-01-20T20:42:03.912Z",
      "domain": "(none)",
      "hostname": "mock-hostname",
      "kernel": "Linux",
      "kernelRelease": "5.3.0-1035-aws",
      "kernelVersion": "#37-Ubuntu SMP Sun Sep 6 01:17:41 UTC 2020",
      "mid": 51,
      "os": "Ubuntu",
      "osVersion": "18.04",
      "tags": {
        "Account": "123456789012",
        "AmiId": "ami-abc12345678901234",
        "ExternalIp": "1.2.3.4",
        "Hostname": "ip-10-0-1-1.us-west-2.compute.internal",
        "InstanceId": "i-abc12345678901234",
        "InternalIp": "10.0.1.1",
        "LwTokenShort": "abcdefghijklm12345678901234567",
        "SubnetId": "subnet-abc1234xyz567890c",
        "VmInstanceType": "t4g.nano",
        "VmProvider": "AWS",
        "VpcId": "vpc-abc12345678901234",
        "Zone": "us-west-2",
        "arch": "arm64",
        "os": "linux"
      }
    }
  ],
  "paging": {
    "rows": 1,
    "totalRows": 1,
    "urls": {
      "nextPage": null
    }
  }
}
	`
}
