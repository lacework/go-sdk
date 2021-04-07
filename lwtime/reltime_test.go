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

package lwtime_test

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/lwtime"
)

func TestMondays(t *testing.T) {
	for y := 1970; y <= time.Now().Year(); y++ {
		t.Run(strconv.Itoa(y), func(t *testing.T) {
			mondays := lwtime.Mondays(y)
			assert.GreaterOrEqual(t, 53, len(mondays))
			assert.LessOrEqual(t, 52, len(mondays))
		})
	}
}

func TestRelTimeUnitIsValid(t *testing.T) {
	assert.True(t, lwtime.RelTimeUnit("y").IsValid())
	assert.True(t, lwtime.RelTimeUnit("MON").IsValid()) // case-sensitivity test
	assert.True(t, lwtime.RelTimeUnit("w").IsValid())
	assert.True(t, lwtime.RelTimeUnit("d").IsValid())
	assert.True(t, lwtime.RelTimeUnit("h").IsValid())
	assert.True(t, lwtime.RelTimeUnit("m").IsValid())
	assert.True(t, lwtime.RelTimeUnit("s").IsValid())
	assert.True(t, lwtime.RelTimeUnit("").IsValid())
	assert.False(t, lwtime.RelTimeUnit("no such thing").IsValid())
}

type RTUSnapTest struct {
	Name                string
	ReferenceTime       string
	ReferenceTimeFormat string
	Snap                lwtime.RelTimeUnit
	Output              string
	Error               error
}

var RTUSnapTests []RTUSnapTest = []RTUSnapTest{
	RTUSnapTest{
		"year",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		lwtime.RelTimeUnit("y"),
		"2006-01-01T00:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"month",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		lwtime.RelTimeUnit("MON"), // case-sensitivity test
		"2006-02-01T00:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"week",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		lwtime.RelTimeUnit("w"),
		"2006-01-30T00:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"leap-week",
		"2016-03-03T15:04:05-07:00",
		time.RFC3339,
		lwtime.RelTimeUnit("w"),
		"2016-02-29T00:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"day",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		lwtime.RelTimeUnit("d"),
		"2006-02-02T00:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"hour",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		lwtime.RelTimeUnit("h"),
		"2006-02-02T15:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"minute",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		lwtime.RelTimeUnit("m"),
		"2006-02-02T15:04:00-07:00",
		nil,
	},
	RTUSnapTest{
		"second",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		lwtime.RelTimeUnit("s"),
		"2006-02-02T15:04:05-07:00",
		nil,
	},
	RTUSnapTest{
		"bad-snap",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		lwtime.RelTimeUnit("no such thing"),
		"",
		errors.New("snap (no such thing) is not a valid relative time unit"),
	},
}

func TestRelTimeUnitSnap(t *testing.T) {
	for _, rtust := range RTUSnapTests {
		t.Run(rtust.Name, func(t *testing.T) {
			refTime, err := time.Parse(rtust.ReferenceTimeFormat, rtust.ReferenceTime)
			if err != nil {
				assert.FailNow(t, "unable to parse reference time")
			}
			outTime, err := lwtime.RelTimeUnit(rtust.Snap).SnapTime(refTime)
			if rtust.Error == nil {
				assert.Equal(t, rtust.Output, outTime.Format(time.RFC3339))
				return
			}
			assert.Equal(t, rtust.Error.Error(), err.Error())
		})
	}
}

type RelTimeParseTest struct {
	Name   string
	Input  string
	Output lwtime.RelTime
	Error  error
}

var RelTimeParseTests []RelTimeParseTest = []RelTimeParseTest{
	RelTimeParseTest{
		"now",
		"now",
		lwtime.RelTime{
			"0",
			0,
			lwtime.RelTimeUnit("s"),
			lwtime.RelTimeUnit(""),
		},
		nil,
	},
	RelTimeParseTest{
		"minus-10-seconds",
		"-10s",
		lwtime.RelTime{
			"-10",
			-10,
			lwtime.RelTimeUnit("s"),
			lwtime.RelTimeUnit(""),
		},
		nil,
	},
	RelTimeParseTest{
		"plus-5-minutes",
		"+5m",
		lwtime.RelTime{
			"5",
			5,
			lwtime.RelTimeUnit("m"),
			lwtime.RelTimeUnit(""),
		},
		nil,
	},
	RelTimeParseTest{
		"at-hour",
		"@h",
		lwtime.RelTime{
			"0",
			0,
			lwtime.RelTimeUnit("s"),
			lwtime.RelTimeUnit("h"),
		},
		nil,
	},
	RelTimeParseTest{
		"minus-1-day-at-hour",
		"-1d@h",
		lwtime.RelTime{
			"-1",
			-1,
			lwtime.RelTimeUnit("d"),
			lwtime.RelTimeUnit("h"),
		},
		nil,
	},
	RelTimeParseTest{
		"minus-1-week-at-week",
		"-1w@w",
		lwtime.RelTime{
			"-7",
			-7,
			lwtime.RelTimeUnit("d"),
			lwtime.RelTimeUnit("w"),
		},
		nil,
	},
	RelTimeParseTest{
		"minus-1-mon",
		"-1MON", // case-sensitivity test
		lwtime.RelTime{
			"-1",
			-1,
			lwtime.RelTimeUnit("mon"),
			lwtime.RelTimeUnit(""),
		},
		nil,
	},
	RelTimeParseTest{
		"plus-3-years",
		"3y",
		lwtime.RelTime{
			"3",
			3,
			lwtime.RelTimeUnit("y"),
			lwtime.RelTimeUnit(""),
		},
		nil,
	},
	RelTimeParseTest{
		"empty",
		"",
		lwtime.RelTime{},
		errors.New("relative time specifier () is invalid"),
	},
	RelTimeParseTest{
		"completely-bad",
		"completely bad",
		lwtime.RelTime{},
		errors.New("relative time specifier (completely bad) is invalid"),
	},
	RelTimeParseTest{
		"minus-1-invalid",
		"-1x",
		lwtime.RelTime{},
		errors.New("invalid unit for relative time specifier (-1x)"),
	},
	RelTimeParseTest{
		"at-x",
		"@x",
		lwtime.RelTime{},
		errors.New("invalid snap for relative time specifier (@x)"),
	},
}

func TestRelTimeParse(t *testing.T) {
	for _, rtpt := range RelTimeParseTests {
		t.Run(rtpt.Name, func(t *testing.T) {
			input := lwtime.RelTime{}
			err := input.Parse(rtpt.Input)
			if rtpt.Error == nil {
				assert.Nil(t, err)
				assert.Equal(t, rtpt.Output, input)
				return
			}
			assert.Equal(t, rtpt.Error.Error(), err.Error())
		})
	}
}

type RelTimeTime_Test struct {
	Name                string
	ReferenceTime       string
	ReferenceTimeFormat string
	Input               string
	Output              string
	Error               error
}

var RelTimeTime_Tests []RelTimeTime_Test = []RelTimeTime_Test{
	RelTimeTime_Test{
		"now",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		"now",
		"2006-02-02T15:04:05-07:00",
		nil,
	},
	RelTimeTime_Test{
		"minus-10-seconds",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		"-10s",
		"2006-02-02T15:03:55-07:00",
		nil,
	},
	RelTimeTime_Test{
		"plus-5-minutes",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		"5m",
		"2006-02-02T15:09:05-07:00",
		nil,
	},
	RelTimeTime_Test{
		"minus-12-hours-at-hour",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		"-12h@h",
		"2006-02-02T03:00:00-07:00",
		nil,
	},
	RelTimeTime_Test{
		"minus-1-week-at-week",
		"2006-01-01T15:04:05-07:00",
		time.RFC3339,
		"-1w@w",
		"2005-12-19T00:00:00-07:00",
		nil,
	},
	RelTimeTime_Test{
		"plus-1-month-at-day",
		"2007-01-29T15:04:05-07:00",
		time.RFC3339,
		"1MON@d", // case-sensitivity test
		"2007-03-01T00:00:00-07:00",
		nil,
	},
	RelTimeTime_Test{
		"minus-2-years",
		"2009-02-02T15:04:05-07:00",
		time.RFC3339,
		"-2y",
		"2007-02-02T15:04:05-07:00",
		nil,
	},
	RelTimeTime_Test{
		"minus-730-days", // 2 years in days
		"2009-02-02T15:04:05-07:00",
		time.RFC3339,
		"-730d",
		"2007-02-03T15:04:05-07:00",
		nil,
	},
	RelTimeTime_Test{
		"minus-3000-years",
		"2009-02-02T15:04:05-07:00",
		time.RFC3339,
		"-3000y",
		"",
		errors.New("unable to construct time object: time predates epoch"),
	},
	RelTimeTime_Test{
		"minus-3-invalid",
		"2009-02-02T15:04:05-07:00",
		time.RFC3339,
		"-3x",
		"",
		errors.New("unable to construct time object: relative time unit (x) is invalid"),
	},
	RelTimeTime_Test{
		"minus-3-years-at-invalid",
		"2009-02-02T15:04:05-07:00",
		time.RFC3339,
		"-3y@x",
		"",
		errors.New("unable to construct time object: snap (x) is not a valid relative time unit"),
	},
}

func TestRelTimeTime_(t *testing.T) {
	for _, rt_t := range RelTimeTime_Tests {
		t.Run(rt_t.Name, func(t *testing.T) {
			refTime, err := time.Parse(rt_t.ReferenceTimeFormat, rt_t.ReferenceTime)
			if err != nil {
				assert.FailNow(t, "unable to parse reference time")
			}
			rt := lwtime.RelTime{}
			err = rt.Parse(rt_t.Input)
			// if we're expecting an error don't worry about this
			if rt_t.Error == nil {
				assert.Nil(t, err)
			}
			outTime, err := rt.Time_(refTime)

			if rt_t.Error == nil {
				assert.Nil(t, err)
				assert.Equal(t, rt_t.Output, outTime.Format(time.RFC3339))
				return
			}
			assert.Equal(t, rt_t.Error.Error(), err.Error())
		})
	}
}
