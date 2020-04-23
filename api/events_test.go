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

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
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
	response, err := c.Events.ListRange(now, from)
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
			assert.Equal(t, "GET", r.Method, "List or ListRange should be a GET method")

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
