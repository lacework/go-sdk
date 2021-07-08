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
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwtime"
)

var (
	queryEvaluator = "Cloudtrail"
	queryID        = "my_lql"
	newQueryText   = `my_lql { source { CloudTrailRawEvents } return { INSERT_ID } }`
	newQuery       = api.NewQuery{
		EvaluatorID: queryEvaluator,
		QueryID:     queryID,
		QueryText:   newQueryText,
	}
	newQueryJSON = fmt.Sprintf(`{
	"evaluatorId": "%s",
	"queryId": "%s",
	"queryText": "%s"
	}`, queryEvaluator, queryID, newQueryText)
	newQueryYAML = fmt.Sprintf(`---
	evaluatorId: %s
	queryId: %s
	queryText: %s
	`, queryEvaluator, queryID, newQueryText)
)

type parseQueryTest struct {
	Name     string
	Input    string
	Return   error
	Expected api.NewQuery
}

var parseQueryTests = []parseQueryTest{
	parseQueryTest{
		Name:     "empty-blob",
		Input:    "",
		Return:   errors.New("query must be valid JSON or YAML"),
		Expected: api.NewQuery{},
	},
	parseQueryTest{
		Name:     "junk-blob",
		Input:    "this is junk",
		Return:   errors.New("query must be valid JSON or YAML"),
		Expected: api.NewQuery{},
	},
	parseQueryTest{
		Name:     "partial-blob",
		Input:    "{",
		Return:   errors.New("query must be valid JSON or YAML"),
		Expected: api.NewQuery{},
	},
	parseQueryTest{
		Name:   "json-blob",
		Input:  newQueryJSON,
		Return: nil,
		Expected: api.NewQuery{
			QueryID:     queryID,
			QueryText:   newQueryText,
			EvaluatorID: queryEvaluator,
		},
	},
	parseQueryTest{
		Name:   "yaml-blob",
		Input:  newQueryYAML,
		Return: nil,
		Expected: api.NewQuery{
			QueryID:     queryID,
			QueryText:   newQueryText,
			EvaluatorID: queryEvaluator,
		},
	},
}

func TestparseQuery(t *testing.T) {
	for _, pqt := range parseQueryTests {
		t.Run(pqt.Name, func(t *testing.T) {
			actual, err := parseQuery(pqt.Input)
			if err == nil {
				assert.Equal(t, pqt.Return, err)
			} else {
				assert.Equal(t, pqt.Return.Error(), err.Error())
			}
			assert.Equal(t, pqt.Expected, actual)
		})
	}
}

type parseQueryTimeTest struct {
	Name       string
	Input      string
	ReturnTime string
	ReturnErr  error
}

var (
	atDay, _            = lwtime.ParseRelative("@d")
	parseQueryTimeTests = []parseQueryTimeTest{
		parseQueryTimeTest{
			Name:       "valid-rfc-utc",
			Input:      "2021-03-31T00:00:00Z",
			ReturnTime: "2021-03-31T00:00:00Z",
			ReturnErr:  nil,
		},
		parseQueryTimeTest{
			Name:       "valid-rfc-central",
			Input:      "2021-03-31T00:00:00-05:00",
			ReturnTime: "2021-03-31T05:00:00Z",
			ReturnErr:  nil,
		},
		parseQueryTimeTest{
			Name:       "valid-milli",
			Input:      "1617230464000",
			ReturnTime: "2021-03-31T22:41:04Z",
			ReturnErr:  nil,
		},
		parseQueryTimeTest{
			Name:       "valid-relative",
			Input:      "@d",
			ReturnTime: atDay.UTC().Format(time.RFC3339),
			ReturnErr:  nil,
		},
		parseQueryTimeTest{
			Name:       "empty",
			Input:      "",
			ReturnTime: "0001-01-01T00:00:00Z",
			ReturnErr:  errors.New("unable to parse time ()"),
		},
		parseQueryTimeTest{
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
