//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestEventsV2Search(t *testing.T) {
	var (
		id         = rand.Intn(10000)
		apiPath    = "Events/search"
		fakeServer = lacework.MockServer()
	)
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	defer fakeServer.Close()

	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
		fmt.Fprintf(w, generateEventsV2Response(mockEventsV2List(id)))
	})

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	now := time.Now().UTC()
	before := now.AddDate(0, 0, -7) // last 7 days
	filter := api.SearchFilter{
		TimeFilter: &api.TimeFilter{
			StartTime: &before,
			EndTime:   &now,
		},
	}

	response, err := c.V2.Events.Search(filter)
	assert.Nil(t, err)
	assert.NotNil(t, response)

	event := response.Data[0]
	type srcEventStruct struct {
		ExePath  string `json:"exe_path"`
		Hostname string `json:"hostname"`
		Mid      int    `json:"mid"`
		Pid      int    `json:"pid_hash"`
		Username string `json:"username"`
	}
	expectedSrcEvent := srcEventStruct{
		ExePath:  "/usr/bin/python3.9",
		Hostname: "ip-11-111-1-111.us-east-1.compute.internal",
		Mid:      123456,
		Pid:      123456789,
		Username: "root",
	}

	var srcEvent srcEventStruct
	srcEventJson, _ := json.Marshal(event.SrcEvent)
	err = json.Unmarshal(srcEventJson, &srcEvent)
	assert.NoError(t, err)
	assert.Equal(t, expectedSrcEvent, srcEvent)
	assert.Equal(t, 1, event.EventCount)
	assert.Equal(t, id, event.Id)
	start, _ := time.Parse(time.RFC3339, "2022-08-13T14:00:00.000Z")
	assert.Equal(t, start, event.StartTime)
	end, _ := time.Parse(time.RFC3339, "2022-08-13T15:00:00.000Z")
	assert.Equal(t, end, event.EndTime)
	assert.Equal(t, "SuspiciousApplicationLaunched", event.EventType)
}

func generateEventsV2Response(data string) string {
	return `
		{
			"data": [` + data + `]
		}
	`
}

func mockEventsV2List(id int) string {
	return fmt.Sprintf(`        {
            "endTime": "2022-08-13T15:00:00.000Z",
            "eventCount": 1,
            "eventType": "SuspiciousApplicationLaunched",
            "id": %d,
            "srcEvent": {
                "exe_path": "/usr/bin/python3.9",
                "hostname": "ip-11-111-1-111.us-east-1.compute.internal",
                "mid": 123456,
                "pid_hash": 123456789,
                "username": "root"
            },
            "srcType": "Process",
            "startTime": "2022-08-13T14:00:00.000Z"
        },
        {
            "endTime": "2022-08-13T15:00:00.000Z",
            "eventCount": 1,
            "eventType": "PolicyViolationChanged",
            "id": 955156,
            "srcEvent": {
                "activity": "New",
                "activity_end_time": "2022-08-13 08:00:00.000 -0700",
                "activity_start_time": "2022-08-13 07:00:00.000 -0700",
                "filedata_hash": "abcderfg1234567",
                "last_modified_time": "2022-03-10 14:20:55.123 -0800",
                "mid": 123345,
                "path": "/etc/eks/containerd/containerd-config.toml",
                "size": 469
            },
            "srcType": "File",
            "startTime": "2022-08-13T14:00:00.000Z"
        }`, id)
}
