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

package lwtime

import (
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestMondays(t *testing.T) {
	for y := 1970; y <= time.Now().Year(); y++ {
		t.Run(strconv.Itoa(y), func(t *testing.T) {
			mondays := mondays(y)
			assert.GreaterOrEqual(t, 53, len(mondays))
			assert.LessOrEqual(t, 52, len(mondays))
		})
	}
}

func TestRelativeUnitIsValid(t *testing.T) {
	assert.True(t, relativeUnit("y").isValid())
	assert.True(t, relativeUnit("MON").isValid()) // case-sensitivity test
	assert.True(t, relativeUnit("w").isValid())
	assert.True(t, relativeUnit("d").isValid())
	assert.True(t, relativeUnit("h").isValid())
	assert.True(t, relativeUnit("m").isValid())
	assert.True(t, relativeUnit("s").isValid())
	assert.True(t, relativeUnit("").isValid())
	assert.False(t, relativeUnit("no such thing").isValid())
}

type RTUSnapTest struct {
	Name                string
	ReferenceTime       string
	ReferenceTimeFormat string
	Snap                relativeUnit
	Output              string
	Error               error
}

var RTUSnapTests = []RTUSnapTest{
	RTUSnapTest{
		"year",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		relativeUnit("y"),
		"2006-01-01T00:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"month",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		relativeUnit("MON"), // case-sensitivity test
		"2006-02-01T00:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"week",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		relativeUnit("w"),
		"2006-01-30T00:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"leap-week",
		"2016-03-03T15:04:05-07:00",
		time.RFC3339,
		relativeUnit("w"),
		"2016-02-29T00:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"day",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		relativeUnit("d"),
		"2006-02-02T00:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"hour",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		relativeUnit("h"),
		"2006-02-02T15:00:00-07:00",
		nil,
	},
	RTUSnapTest{
		"minute",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		relativeUnit("m"),
		"2006-02-02T15:04:00-07:00",
		nil,
	},
	RTUSnapTest{
		"second",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		relativeUnit("s"),
		"2006-02-02T15:04:05-07:00",
		nil,
	},
	RTUSnapTest{
		"bad-snap",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		relativeUnit("no such thing"),
		"",
		errors.New("snap (no such thing) is not a valid relative time unit"),
	},
}

func TestRelativeUnitSnap(t *testing.T) {
	for _, rtust := range RTUSnapTests {
		t.Run(rtust.Name, func(t *testing.T) {
			refTime, err := time.Parse(rtust.ReferenceTimeFormat, rtust.ReferenceTime)
			if err != nil {
				assert.FailNow(t, "unable to parse reference time")
			}
			outTime, err := relativeUnit(rtust.Snap).snapTime(refTime)
			if rtust.Error == nil {
				assert.Equal(t, rtust.Output, outTime.Format(time.RFC3339))
				return
			}
			assert.Equal(t, rtust.Error.Error(), err.Error())
		})
	}
}

type newRelativeTest struct {
	Name   string
	Input  string
	Output relative
	Error  error
}

var newRelativeTests = []newRelativeTest{
	newRelativeTest{
		"now",
		"now",
		relative{
			"0",
			0,
			relativeUnit("s"),
			relativeUnit(""),
		},
		nil,
	},
	newRelativeTest{
		"minus-10-seconds",
		"-10s",
		relative{
			"-10",
			-10,
			relativeUnit("s"),
			relativeUnit(""),
		},
		nil,
	},
	newRelativeTest{
		"plus-5-minutes",
		"+5m",
		relative{
			"5",
			5,
			relativeUnit("m"),
			relativeUnit(""),
		},
		nil,
	},
	newRelativeTest{
		"at-hour",
		"@h",
		relative{
			"0",
			0,
			relativeUnit("s"),
			relativeUnit("h"),
		},
		nil,
	},
	newRelativeTest{
		"minus-1-day-at-hour",
		"-1d@h",
		relative{
			"-1",
			-1,
			relativeUnit("d"),
			relativeUnit("h"),
		},
		nil,
	},
	newRelativeTest{
		"minus-1-week-at-week",
		"-1w@w",
		relative{
			"-7",
			-7,
			relativeUnit("d"),
			relativeUnit("w"),
		},
		nil,
	},
	newRelativeTest{
		"minus-1-mon",
		"-1MON", // case-sensitivity test
		relative{
			"-1",
			-1,
			relativeUnit("mon"),
			relativeUnit(""),
		},
		nil,
	},
	newRelativeTest{
		"plus-3-years",
		"3y",
		relative{
			"3",
			3,
			relativeUnit("y"),
			relativeUnit(""),
		},
		nil,
	},
	newRelativeTest{
		"empty",
		"",
		relative{},
		errors.New("relative time specifier () is invalid"),
	},
	newRelativeTest{
		"completely-bad",
		"completely bad",
		relative{},
		errors.New("relative time specifier (completely bad) is invalid"),
	},
	newRelativeTest{
		"minus-1-invalid",
		"-1x",
		relative{},
		errors.New("invalid unit for relative time specifier (-1x)"),
	},
	newRelativeTest{
		"at-x",
		"@x",
		relative{},
		errors.New("invalid snap for relative time specifier (@x)"),
	},
}

func TestNewRelative(t *testing.T) {
	for _, nrt := range newRelativeTests {
		t.Run(nrt.Name, func(t *testing.T) {
			input, err := newRelative(nrt.Input)

			if nrt.Error == nil {
				assert.Nil(t, err)
				assert.Equal(t, nrt.Output, input)
				return
			}
			assert.Equal(t, nrt.Error.Error(), err.Error())
		})
	}
}

type RelativeTimeTest struct {
	Name                string
	ReferenceTime       string
	ReferenceTimeFormat string
	Input               string
	Output              string
	Error               error
}

var RelativeTimeTests = []RelativeTimeTest{
	RelativeTimeTest{
		"now",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		"now",
		"2006-02-02T15:04:05-07:00",
		nil,
	},
	RelativeTimeTest{
		"minus-10-seconds",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		"-10s",
		"2006-02-02T15:03:55-07:00",
		nil,
	},
	RelativeTimeTest{
		"plus-5-minutes",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		"5m",
		"2006-02-02T15:09:05-07:00",
		nil,
	},
	RelativeTimeTest{
		"minus-12-hours-at-hour",
		"2006-02-02T15:04:05-07:00",
		time.RFC3339,
		"-12h@h",
		"2006-02-02T03:00:00-07:00",
		nil,
	},
	RelativeTimeTest{
		"minus-1-week-at-week",
		"2006-01-01T15:04:05-07:00",
		time.RFC3339,
		"-1w@w",
		"2005-12-19T00:00:00-07:00",
		nil,
	},
	RelativeTimeTest{
		"plus-1-month-at-day",
		"2007-01-29T15:04:05-07:00",
		time.RFC3339,
		"1MON@d", // case-sensitivity test
		"2007-03-01T00:00:00-07:00",
		nil,
	},
	RelativeTimeTest{
		"minus-2-years",
		"2009-02-02T15:04:05-07:00",
		time.RFC3339,
		"-2y",
		"2007-02-02T15:04:05-07:00",
		nil,
	},
	RelativeTimeTest{
		"minus-730-days", // 2 years in days
		"2009-02-02T15:04:05-07:00",
		time.RFC3339,
		"-730d",
		"2007-02-03T15:04:05-07:00",
		nil,
	},
	RelativeTimeTest{
		"minus-3000-years",
		"2009-02-02T15:04:05-07:00",
		time.RFC3339,
		"-3000y",
		"",
		errors.New("unable to construct time object: time predates epoch"),
	},
	RelativeTimeTest{
		"minus-3-invalid",
		"2009-02-02T15:04:05-07:00",
		time.RFC3339,
		"-3x",
		"",
		errors.New("unable to construct time object: relative time unit (x) is invalid"),
	},
	RelativeTimeTest{
		"minus-3-years-at-invalid",
		"2009-02-02T15:04:05-07:00",
		time.RFC3339,
		"-3y@x",
		"",
		errors.New("unable to construct time object: snap (x) is not a valid relative time unit"),
	},
}

func TestRelativeTime(t *testing.T) {
	for _, rtt := range RelativeTimeTests {
		t.Run(rtt.Name, func(t *testing.T) {
			refTime, err := time.Parse(rtt.ReferenceTimeFormat, rtt.ReferenceTime)
			if err != nil {
				assert.FailNow(t, "unable to parse reference time")
			}

			rel, err := newRelative(rtt.Input)
			// if we're expecting an error don't worry about this
			if rtt.Error == nil {
				assert.Nil(t, err)
			}
			outTime, err := rel.time(refTime)

			if rtt.Error == nil {
				assert.Nil(t, err)
				assert.Equal(t, rtt.Output, outTime.Format(time.RFC3339))
				return
			}
			assert.Equal(t, rtt.Error.Error(), err.Error())
		})
	}
}
