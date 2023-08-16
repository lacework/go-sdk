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

package cmd

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwtime"
)

var (
	newQuery = api.NewQuery{
		QueryID:   "my_lql",
		QueryText: `my_lql { source { CloudTrailRawEvents } return { INSERT_ID } }`,
	}
)

type parseQueryTimeTest struct {
	Name       string
	Input      string
	ReturnTime string
	ReturnErr  error
}

var (
	atDay, _            = lwtime.ParseRelative("@d")
	parseQueryTimeTests = []parseQueryTimeTest{
		{
			Name:       "valid-rfc-utc",
			Input:      "2021-03-31T00:00:00Z",
			ReturnTime: "2021-03-31T00:00:00Z",
			ReturnErr:  nil,
		},
		{
			Name:       "valid-rfc-central",
			Input:      "2021-03-31T00:00:00-05:00",
			ReturnTime: "2021-03-31T05:00:00Z",
			ReturnErr:  nil,
		},
		{
			Name:       "valid-milli",
			Input:      "1617230464000",
			ReturnTime: "2021-03-31T22:41:04Z",
			ReturnErr:  nil,
		},
		{
			Name:       "valid-relative",
			Input:      "@d",
			ReturnTime: atDay.UTC().Format(time.RFC3339),
			ReturnErr:  nil,
		},
		{
			Name:       "empty",
			Input:      "",
			ReturnTime: "0001-01-01T00:00:00Z",
			ReturnErr:  errors.New("unable to parse time ()"),
		},
		{
			Name:       "invalid",
			Input:      "jweaver",
			ReturnTime: "0001-01-01T00:00:00Z",
			ReturnErr:  errors.New("unable to parse time (jweaver)"),
		},
	}
)

func TestParseQueryTime(t *testing.T) {
	for _, pqtt := range parseQueryTimeTests {
		t.Run(pqtt.Name, func(t *testing.T) {
			outTime, err := parseQueryTime(pqtt.Input)
			if err == nil {
				assert.Equal(t, pqtt.ReturnErr, err)
			} else {
				assert.Equal(t, pqtt.ReturnErr.Error(), err.Error())
			}
			assert.Equal(t, pqtt.ReturnTime, outTime.UTC().Format(time.RFC3339))
		})
	}
}

type getRunStartProgressMessageTest struct {
	Name                string
	StartTime           string
	EndTime             string
	ReferenceTimeFormat string
	Return              string
}

var (
	getRunStartProgressMessageTests = []getRunStartProgressMessageTest{
		{
			Name:                "no-start",
			EndTime:             "2006-02-02T15:04:05-07:00",
			ReferenceTimeFormat: time.RFC3339,
			Return:              "Executing query",
		},
		{
			Name:                "no-end",
			StartTime:           "2006-02-02T15:04:05-07:00",
			ReferenceTimeFormat: time.RFC3339,
			Return:              "Executing query",
		},
		{
			Name:                "basic",
			StartTime:           "2006-02-02T15:04:05-07:00",
			EndTime:             "2006-02-03T15:04:05-07:00",
			ReferenceTimeFormat: time.RFC3339,
			Return:              "Executing query in the time range 2006-Feb-2 22:04:05 UTC - 2006-Feb-3 22:04:05 UTC",
		},
	}
)

func TestRunStartProgressMessage(t *testing.T) {
	for _, grspmt := range getRunStartProgressMessageTests {
		t.Run(grspmt.Name, func(t *testing.T) {
			args := []api.ExecuteQueryArgument{}

			if grspmt.StartTime != "" {
				startTime, startErr := time.Parse(grspmt.ReferenceTimeFormat, grspmt.StartTime)
				if startErr != nil {
					assert.FailNow(t, startErr.Error())
				}
				args = append(args, api.ExecuteQueryArgument{
					Name:  api.QueryStartTimeRange,
					Value: startTime.UTC().Format(lwtime.RFC3339Milli),
				})
			}
			if grspmt.EndTime != "" {
				endTime, endErr := time.Parse(grspmt.ReferenceTimeFormat, grspmt.EndTime)
				if endErr != nil {
					assert.FailNow(t, endErr.Error())
				}
				args = append(args, api.ExecuteQueryArgument{
					Name:  api.QueryEndTimeRange,
					Value: endTime.UTC().Format(lwtime.RFC3339Milli),
				})
			}

			assert.Equal(t, grspmt.Return, getRunStartProgressMessage(args))
		})
	}
}
