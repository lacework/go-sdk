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

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestAgentInfoSearch(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("AgentInfo/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockAgentInfoResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response := api.AgentInfoResponse{}
	err = c.V2.AgentInfo.Search(&response, api.SearchFilter{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, response.Paging.Rows)
	if assert.Equal(t, 2, len(response.Data)) {
		// Machine 1
		assert.Equal(t, "5.4.2-8f652c57", response.Data[0].AgentVersion)
		assert.Equal(t, "circle-ci-test-node", response.Data[0].Hostname)
		assert.Equal(t, "ACTIVE", response.Data[0].Status)
		assert.Equal(t, 30, response.Data[0].Mid)
		// Machine 2
		assert.Equal(t, "5.8.0-6b21f03f", response.Data[1].AgentVersion)
		assert.Equal(t, "gke-hipstershop-default-pool-abc12343-0ivv", response.Data[1].Hostname)
		assert.Equal(t, "ACTIVE", response.Data[1].Status)
		assert.Equal(t, 168, response.Data[1].Mid)
	}
}

func mockAgentInfoResponse() string {
	return `
{
  "data": [
    {
      "agentVersion": "5.4.2-8f652c57",
      "createdTime": "2020-10-23T17:33:35.646Z",
      "hostname": "circle-ci-test-node",
      "ipAddr": "1.0.1.1",
      "lastUpdate": "2022-08-04T19:34:57.480Z",
      "mid": 30,
      "mode": "normal",
      "os": "Linux",
      "status": "ACTIVE",
      "tags": {
        "Account": "123456789012",
        "AmiId": "ami-1234abcd1234abcd1",
        "ExternalIp": "2.2.2.2",
        "Hostname": "ip-1-0-1-1.us-west-2.compute.internal",
        "InstanceId": "i-abcde12345",
        "InternalIp": "1.0.1.1",
        "LwTokenShort": "abc1234e545abc12349eababc12342",
        "SubnetId": "subnet-abc1234fabc12349c",
        "VmInstanceType": "t2.micro",
        "VmProvider": "AWS",
        "VpcId": "vpc-09abc1234abc1234d",
        "Zone": "us-west-2a",
        "arch": "amd64",
        "os": "linux"
      }
    },
    {
      "agentVersion": "5.8.0-6b21f03f",
      "createdTime": "2022-08-04T18:40:47.681Z",
      "hostname": "gke-hipstershop-default-pool-abc12343-0ivv",
      "ipAddr": "1.2.0.1",
      "lastUpdate": "2022-08-04T19:33:47.733Z",
      "mid": 168,
      "mode": "ebpf",
      "os": "Linux",
      "status": "ACTIVE",
      "tags": {
        "Cluster": "hipstershop",
        "Env": "k8s",
        "ExternalIp": "3.2.1.1",
        "lw_KubernetesCluster": "hipstershop",
        "os": "linux"
      }
    }
  ],
  "paging": {
    "rows": 2,
    "totalRows": 2,
    "urls": {
      "nextPage": null
    }
  }
}
	`
}
