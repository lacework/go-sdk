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

package api

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type validateQueryArgumentsTest struct {
	name         string
	arguments    []ExecuteQueryArgument
	expectedTime time.Time
	retrn        error
}

var validateQueryArgumentsTests = []validateQueryArgumentsTest{
	validateQueryArgumentsTest{
		name:      "empty",
		arguments: []ExecuteQueryArgument{},
		retrn:     nil,
		//retrn: errors.New(`parsing time "" as "2006-01-02T15:04:05.000Z07:00": cannot parse "" as "2006"`),
	},
	validateQueryArgumentsTest{
		name: "start-bad",
		arguments: []ExecuteQueryArgument{
			ExecuteQueryArgument{Name: "StartTimeRange", Value: ""},
		},
		retrn: errors.New(
			`invalid StartTimeRange argument: parsing time "" as "2006-01-02T15:04:05.000Z": cannot parse "" as "2006"`),
	},
	validateQueryArgumentsTest{
		name: "start-nonutc",
		arguments: []ExecuteQueryArgument{
			ExecuteQueryArgument{Name: "StartTimeRange", Value: "2021-07-11T00:00:00.123Z07:00"},
		},
		retrn: errors.New(
			`invalid StartTimeRange argument: parsing time "2021-07-11T00:00:00.123Z07:00": extra text: "07:00"`),
	},
	validateQueryArgumentsTest{
		name: "start-good",
		arguments: []ExecuteQueryArgument{
			ExecuteQueryArgument{Name: "StartTimeRange", Value: "2021-07-12T00:00:00.000Z"},
		},
		retrn: nil,
	},
	validateQueryArgumentsTest{
		name: "end-bad",
		arguments: []ExecuteQueryArgument{
			ExecuteQueryArgument{Name: "EndTimeRange", Value: ""},
		},
		retrn: errors.New(
			`invalid EndTimeRange argument: parsing time "" as "2006-01-02T15:04:05.000Z": cannot parse "" as "2006"`),
	},
	validateQueryArgumentsTest{
		name: "end-good",
		arguments: []ExecuteQueryArgument{
			ExecuteQueryArgument{Name: "EndTimeRange", Value: "2021-07-12T00:00:00.000Z"},
		},
		retrn: nil,
	},
	validateQueryArgumentsTest{
		name: "range-bad",
		arguments: []ExecuteQueryArgument{
			ExecuteQueryArgument{Name: "StartTimeRange", Value: "2021-07-13T00:00:00.000Z"},
			ExecuteQueryArgument{Name: "EndTimeRange", Value: "2021-07-12T00:00:00.000Z"},
		},
		retrn: errors.New(
			"date range should have a start time before the end time"),
	},
	validateQueryArgumentsTest{
		name: "range-good",
		arguments: []ExecuteQueryArgument{
			ExecuteQueryArgument{Name: "StartTimeRange", Value: "2021-07-12T00:00:00.000Z"},
			ExecuteQueryArgument{Name: "EndTimeRange", Value: "2021-07-13T00:00:00.000Z"},
		},
		retrn: nil,
	},
}

func TestValidateQueryTimeString(t *testing.T) {
	for _, vqat := range validateQueryArgumentsTests {
		t.Run(vqat.name, func(t *testing.T) {
			err := validateQueryArguments(vqat.arguments)
			if err == nil {
				assert.Equal(t, vqat.retrn, err)
			} else {
				assert.Equal(t, vqat.retrn.Error(), err.Error())
			}
		})
	}
}

type validateQueryRangeTest struct {
	name           string
	startTimeRange time.Time
	endTimeRange   time.Time
	retrn          error
}

var validateQueryRangeTests = []validateQueryRangeTest{
	validateQueryRangeTest{
		name:           "ok",
		startTimeRange: time.Unix(0, 0),
		endTimeRange:   time.Unix(1, 0),
		retrn:          nil,
	},
	validateQueryRangeTest{
		name:           "empty-start",
		startTimeRange: time.Time{},
		endTimeRange:   time.Unix(1, 0),
		retrn:          nil,
	},
	validateQueryRangeTest{
		name:           "empty-end",
		startTimeRange: time.Unix(1, 0),
		endTimeRange:   time.Time{},
		retrn:          errors.New("date range should have a start time before the end time"),
	},
	validateQueryRangeTest{
		name:           "start-after-end",
		startTimeRange: time.Unix(1717333947, 0),
		endTimeRange:   time.Unix(1617333947, 0),
		retrn:          errors.New("date range should have a start time before the end time"),
	},
	validateQueryRangeTest{
		name:           "start-equal-end",
		startTimeRange: time.Unix(1617333947, 0),
		endTimeRange:   time.Unix(1617333947, 0),
		retrn:          nil,
	},
}

func TestValidateQueryRange(t *testing.T) {
	for _, vqrt := range validateQueryRangeTests {
		t.Run(vqrt.name, func(t *testing.T) {
			err := validateQueryRange(vqrt.startTimeRange, vqrt.endTimeRange)
			if err == nil {
				assert.Equal(t, vqrt.retrn, err)
			} else {
				assert.Equal(t, vqrt.retrn.Error(), err.Error())
			}
		})
	}
}
