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
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestEventsSeverity(t *testing.T) {
	var (
		critical = api.Event{Severity: "1"}
		high     = api.Event{Severity: "2"}
		medium   = api.Event{Severity: "3"}
		low      = api.Event{Severity: "4"}
		info     = api.Event{Severity: "5"}
		unknown  = api.Event{Severity: "9"}
	)

	assert.Equal(t, "Critical", critical.SeverityString())
	assert.Equal(t, "High", high.SeverityString())
	assert.Equal(t, "Medium", medium.SeverityString())
	assert.Equal(t, "Low", low.SeverityString())
	assert.Equal(t, "Info", info.SeverityString())
	assert.Equal(t, "Unknown", unknown.SeverityString())
}

func TestEventsListRangeError(t *testing.T) {
	var (
		now    = time.Now().UTC()
		from   = now.AddDate(0, 0, -7) // 7 days from now
		c, err = api.NewClient("test", api.WithToken("TOKEN"))
	)
	assert.Nil(t, err)

	// a tipical user input error could be that they provide the
	// date range the other way around, from should be the start
	// time, and now should be the end time
	response, err := c.Events.ListDateRange(now, from)
	assert.Empty(t, response)
	if assert.NotNil(t, err) {
		assert.Equal(t,
			"data range should have a start time before the end time",
			err.Error(), "error message mismatch",
		)
	}
}

func TestEventsList(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/events/GetEventsForDateRange",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List or ListDateRange should be a GET method")

			start, ok := r.URL.Query()["START_TIME"]
			if assert.True(t, ok,
				"START_TIME parameter missing") {

				end, ok := r.URL.Query()["END_TIME"]
				if assert.True(t, ok,
					"END_TIME parameter missing") {

					// verify that start and end times are 7 days apart
					// and that the start time is before the end time
					startTime, err := time.Parse(time.RFC3339, start[0])
					assert.Nil(t, err)
					endTime, err := time.Parse(time.RFC3339, end[0])
					assert.Nil(t, err)

					assert.True(t,
						startTime.Before(endTime),
						"the start time should not be after the end time",
					)
					assert.True(t,
						startTime.AddDate(0, 0, 7).Equal(endTime),
						"the data range is not 7 days apart",
					)
					fmt.Fprintf(w, arrayOfEventsResponse(start[0]))
				}
			}
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Events.List()
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		assert.Equal(t,
			api.EventsCount{
				High:  2,
				Low:   1,
				Total: 3,
			},
			response.GetEventsCount(),
		)

		if assert.Equal(t, 3, len(response.Events)) {
			assert.Equal(t, "EventTypeGoesHere", response.Events[0].EventType)
		}
	}
}

func TestEventsDetailsErrorEmptyID(t *testing.T) {
	c, err := api.NewClient("test", api.WithToken("TOKEN"))
	assert.Nil(t, err)

	// a tipical user input error could be that they provide an empty event_id
	response, err := c.Events.Details("")
	assert.Empty(t, response)
	if assert.NotNil(t, err) {
		assert.Equal(t,
			"event_id cannot be empty",
			err.Error(), "error message mismatch",
		)
	}
}

func TestEventsDetails(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.MockAPI(
		"external/events/GetEventDetails",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Details should be a GET method")

			eventID, ok := r.URL.Query()["EVENT_ID"]
			if assert.True(t, ok,
				"EVENT_ID parameter missing") {

				fmt.Fprintf(w, eventDetailsResponse(eventID[0]))
			}
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Events.Details("123")
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		if assert.Equal(t, 1, len(response.Events)) {
			eventDetails := response.Events[0]
			assert.Equal(t, "123", eventDetails.EventID)
			assert.Equal(t, "User", eventDetails.EventActor)
			assert.Equal(t, "UserTracking", eventDetails.EventModel)
			assert.Equal(t, "NewExternalServerDNSConn", eventDetails.EventType)

			// TODO @afiune assert EntityMap
			assert.Equal(t,
				"ruby chef-client",
				eventDetails.EntityMap.Application[0].Application,
			)
		}
	}
}

func arrayOfEventsResponse(t string) string {
	return `
{
  "data": [
    {
      "end_time": "` + t + `",
      "event_id": "1",
      "event_type": "EventTypeGoesHere",
      "severity": "4",
      "start_time": "` + t + `"
    },
    {
      "END_TIME": "` + t + `",
      "EVENT_ID": "2",
      "EVENT_TYPE": "EventTypeGoesHere",
      "SEVERITY": "2",
      "START_TIME": "` + t + `"
    },
    {
      "end_time": "` + t + `",
      "event_id": "3",
      "event_type": "EventTypeGoesHere",
      "severity": "2",
      "start_time": "` + t + `"
    }
  ]
}
`
}

// @afiune real response from a demo environment
func eventDetailsResponse(id string) string {
	return `
{
  "data": [
    {
      "END_TIME": "2020-04-20T20:00:00Z",
      "ENTITY_MAP": {
        "Application": [
          {
            "APPLICATION": "ruby chef-client",
            "EARLIEST_KNOWN_TIME": "2020-04-20T19:00:00Z",
            "HAS_EXTERNAL_CONNS": 1,
            "IS_CLIENT": 0,
            "IS_SERVER": 1
          }
        ],
        "FileExePath": [
          { "EXE_PATH": "/opt/chef/embedded/bin/ruby" },
          { "EXE_PATH": "/opt/chef/embedded/bin/ruby" },
          { "EXE_PATH": "/opt/chef/embedded/bin/ruby" }
        ],
        "Machine": [
          {
            "CPU_PERCENTAGE": 0.39,
            "EXTERNAL_IP": "10.0.2.15",
            "HOSTNAME": "default-centos-8.vagrantup.com",
            "INSTANCE_ID": "",
            "INTERNAL_IP_ADDR": "10.0.2.15",
            "IS_EXTERNAL": 1
          },
          {
            "CPU_PERCENTAGE": 0.31,
            "EXTERNAL_IP": "10.0.2.15",
            "HOSTNAME": "default-centos-8.vagrantup.com",
            "INSTANCE_ID": "",
            "INTERNAL_IP_ADDR": "10.0.2.15",
            "IS_EXTERNAL": 1
          },
          {
            "CPU_PERCENTAGE": 0.42,
            "EXTERNAL_IP": "10.0.2.15",
            "HOSTNAME": "default-centos-8.vagrantup.com",
            "INSTANCE_ID": "",
            "INTERNAL_IP_ADDR": "10.0.2.15",
            "IS_EXTERNAL": 1
          }
        ],
        "Process": [
          {
            "CMDLINE": "chef-client worker: ppid=3803;start=19:47:41;",
            "CPU_PERCENTAGE": 0,
            "HOSTNAME": "default-centos-8.vagrantup.com",
            "PROCESS_ID": 3810,
            "PROCESS_START_TIME": "2020-04-20T19:47:40Z"
          },
          {
            "CMDLINE": "chef-client worker: ppid=5328;start=19:54:08;",
            "CPU_PERCENTAGE": 0.1,
            "HOSTNAME": "default-centos-8.vagrantup.com",
            "PROCESS_ID": 5346,
            "PROCESS_START_TIME": "2020-04-20T19:54:08Z"
          },
          {
            "CMDLINE": "chef-client worker: ppid=12057;start=19:47:54;",
            "CPU_PERCENTAGE": 0.87,
            "HOSTNAME": "default-centos-8.vagrantup.com",
            "PROCESS_ID": 12062,
            "PROCESS_START_TIME": "2020-04-20T19:47:54Z"
          }
        ],
        "User": [
          {
            "MACHINE_HOSTNAME": "default-centos-8.vagrantup.com",
            "USERNAME": "root"
          },
          {
            "MACHINE_HOSTNAME": "default-centos-8.vagrantup.com",
            "USERNAME": "root"
          },
          {
            "MACHINE_HOSTNAME": "default-centos-8.vagrantup.com",
            "USERNAME": "root"
          }
        ]
      },
      "EVENT_ACTOR": "User",
      "EVENT_ID": "` + id + `",
      "EVENT_MODEL": "UserTracking",
      "EVENT_TYPE": "NewExternalServerDNSConn",
      "START_TIME": "2020-04-20T19:00:00Z"
    }
  ]
}
`
}
