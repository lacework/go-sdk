//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestPagination(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
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
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response := api.MachineDetailEntityResponse{}

	t.Run("First Request", func(t *testing.T) {
		err := c.V2.Entities.Search(&response, api.SearchFilter{})
		assert.NoError(t, err)
		assert.NotNil(t, response)
		if assert.Equal(t, 1, len(response.Data)) {
			assert.Equal(t, "mock-1-hostname", response.Data[0].Hostname)
		}
	})

	t.Run("Access Page Two", func(t *testing.T) {
		pageOk, err := c.NextPage(&response)
		if assert.True(t, pageOk) && assert.NoError(t, err) {
			if assert.Equal(t, 1, len(response.Data)) {
				assert.Equal(t, "mock-2-hostname", response.Data[0].Hostname)
			}
		}
	})

	t.Run("Access Page Three - should NOT exist", func(t *testing.T) {
		pageOk, err := c.NextPage(&response)
		assert.False(t, pageOk)
		assert.NoError(t, err)
	})

	t.Run("Accessing All Pages - list all hostnames", func(t *testing.T) {
		listHosts := []string{}
		err := c.V2.Entities.Search(&response, api.SearchFilter{})
		assert.NoError(t, err)
		for {
			listHosts = append(listHosts, response.Data[0].Hostname)
			pageOk, err := c.NextPage(&response)
			if err != nil {
				t.Fail()
				t.Logf("Expected no errors, got '%s'", err.Error())
				break
			}

			if pageOk {
				continue
			}
			break
		}
		assert.Equal(t, []string{"mock-1-hostname", "mock-2-hostname"}, listHosts)
	})
}

func TestPaginationWithoutInfo(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/MachineDetails/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockPaginationResponseWithoutPaging())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response := api.MachineDetailEntityResponse{}
	err = c.V2.Entities.Search(&response, api.SearchFilter{})
	assert.NoError(t, err)
	assert.NotNil(t, response)

	pageOk, err := c.NextPage(&response)
	assert.False(t, pageOk, "a request without paging information should return pageOk false")
	assert.NoError(t, err, "a request without paging information should not error")
}

// @afiune this response should have 5000 machine details but we don't really care
// that much about it since this test is testing the information from the "paging" field
func mockPaginationResponsePage1() string {
	return `
{
  "data": [
    {
      "awsInstanceId": "i-abc12345678901234",
      "awsZone": "us-west-2",
      "createdTime": "2022-01-20T20:42:03.912Z",
      "domain": "(none)",
      "hostname": "mock-1-hostname",
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
    "rows": 5000,
    "totalRows": 5001,
    "urls": {
			"nextPage": "https://account.lacework.net/api/v2/NextPage/vuln/containers/abc123"
    }
  }
}
	`
}
func mockPaginationResponsePage2() string {
	return `
{
  "data": [
    {
      "awsInstanceId": "i-abc12345678901234",
      "awsZone": "us-west-2",
      "createdTime": "2022-01-20T20:42:03.912Z",
      "domain": "(none)",
      "hostname": "mock-2-hostname",
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
    "totalRows": 5001,
    "urls": {
      "nextPage": null 
    }
  }
}
	`
}

// @afiune purposely passing a response without paging information
func mockPaginationResponseWithoutPaging() string {
	return `
{
  "data": [
    {
      "awsInstanceId": "i-abc12345678901234",
      "awsZone": "us-west-2",
      "createdTime": "2022-01-20T20:42:03.912Z",
      "domain": "(none)",
      "hostname": "mock-1-hostname",
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
  ]
}
	`
}
