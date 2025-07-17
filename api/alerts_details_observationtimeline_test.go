//
// Author:: Lokesh Vadlamudi (<lvadlamudi@fortinet.com>)
// Copyright:: Copyright 2025, Fortinet Inc.
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
	"net/http"
	"testing"

	"github.com/lacework/go-sdk/v2/api"
	"github.com/lacework/go-sdk/v2/internal/lacework"
	"github.com/stretchr/testify/assert"
)

var alertObservationTimelineJSON = `{
  "data": [
    {
      "description": "Remote system discovery using Ping or ARP observed",
      "endEpoch": 1750692691,
      "entities": [
        {
          "entity_key": {
            "repo": "docker.io/example/ecommerce-website"
          },
          "entity_text": "docker.io/example/ecommerce-website",
          "entity_type": "container_repo",
          "entity_uuid": "gt6LFSoHsasqSBfQpJncoj",
          "is_detail": false,
          "is_subject": false
        },
        {
          "entity_key": {
            "mid": "4325415338702892297",
            "pid_hash": "-5106375999317139109"
          },
          "entity_props": "{\"cmdline\": \"ping -c 1 -t 1 10.0.3.6\"}",
          "entity_text": "ping -c 1 -t 1 10.0.3.6",
          "entity_type": "process",
          "entity_uuid": "TL4dYNpFrFxZRhrT2ZqNRJ",
          "is_detail": true,
          "is_subject": false
        },
        {
          "entity_key": {
            "account": "123456789012",
            "principal_id": "AROAEXAMPLE:example-instance",
            "username": "AssumedRole/123456789012:ExampleRole"
          },
          "entity_text": "AssumedRole/123456789012:ExampleRole",
          "entity_type": "ct_user",
          "entity_uuid": "ZXXdzywqQCCRz97n9Q57Mw",
          "is_detail": false,
          "is_subject": true
        },
        {
          "entity_key": {
            "hostname": "host-1.example.com",
            "mid": "4325415338702892297"
          },
          "entity_props": "{\"internal_ip_addr\": \"10.0.3.226\", \"os_type\": \"linux\"}",
          "entity_text": "host-1.example.com",
          "entity_type": "machine",
          "entity_uuid": "muupdPVN2nZWokr8HkooLK",
          "is_detail": false,
          "is_subject": true
        },
        {
          "entity_key": {
            "cluster_id": "example-cluster",
            "pod_name": "ecommerce-website-58dc7b54f9-ctp2n",
            "pod_namespace": "default"
          },
          "entity_text": "ecommerce-website-58dc7b54f9-ctp2n",
          "entity_type": "k8_pod",
          "entity_uuid": "WWJksiQgq2zBnronZquHzN",
          "is_detail": true,
          "is_subject": false
        }
      ],
      "formalTags": [
        {
          "customer_facing_id": "T1018",
          "id": "T1018",
          "name": "Remote System Discovery",
          "order_index": 25,
          "parent_id": "TA0007",
          "url": "https://attack.mitre.org/techniques/T1018"
        },
        {
          "customer_facing_id": "TA0007",
          "id": "TA0007",
          "name": "Discovery",
          "order_index": 6,
          "url": "https://attack.mitre.org/tactics/TA0007"
        }
      ],
      "observationPivotEntityUuids": [
        "muupdPVN2nZWokr8HkooLK"
      ],
      "observationType": "linux_discovery_remote_system_discovery",
      "recordUuid": "2e034bf4-683f-56f0-a91b-bc8c4510d2b3",
      "relationships": [
        {
          "dstEntityKey": {
            "mid": "4325415338702892297",
            "pid_hash": "-5161409754333993982"
          },
          "dstEntityText": "sh -x ./deepce.sh",
          "dstEntityType": "process",
          "dstEntityUuid": "CExjyKEFbc3ngP5GFeNja4",
          "relationshipDescriptor": "is a child process of process",
          "relationshipId": "c8a6fb3e23fa89ab3b977dc8de815204fe738c87",
          "srcEntityKey": {
            "mid": "4325415338702892297",
            "pid_hash": "-5112041715404371368"
          },
          "srcEntityText": "ping -c 1 127.0.0.1",
          "srcEntityType": "process",
          "srcEntityUuid": "EHnMTz74eJmgfb2iLpzkcx"
        }
      ],
      "startEpoch": 1750692691
    }
  ]
}`

func TestAlertsGetObservationTimelineMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetObservationTimeline should be a GET method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.GetObservationTimeline(alertID)
	assert.Nil(t, err)
}

func TestAlertsGetObservationTimelineOK(t *testing.T) {
	mockResponse := alertObservationTimelineJSON

	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, mockResponse)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)
	// Get actual response from SDK method
	resp, err := c.V2.Alerts.GetObservationTimeline(alertID)
	assert.Nil(t, err)

	// Marshal actual output back to JSON
	actualJSON, err := json.Marshal(resp)
	assert.Nil(t, err)

	// Compare JSON using assert.JSONEq
	assert.JSONEq(t, mockResponse, string(actualJSON))

}

func TestAlertsGetObservationTimelineError(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		fmt.Sprintf("Alerts/%d", alertID),
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, lqlErrorReponse, http.StatusInternalServerError)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Alerts.GetObservationTimeline(alertID)
	assert.NotNil(t, err)
}
