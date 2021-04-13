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
	"fmt"
	"testing"
	"time"

	"github.com/lacework/go-sdk/lwtime"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNatTimeAdjectiveIsValid(t *testing.T) {
	assert.True(t, lwtime.NatTimeAdjective("today").IsValid())
	assert.True(t, lwtime.NatTimeAdjective("YESTERDAY").IsValid()) // case-sensitivity test
	assert.True(t, lwtime.NatTimeAdjective("this").IsValid())
	assert.True(t, lwtime.NatTimeAdjective("current").IsValid())
	assert.True(t, lwtime.NatTimeAdjective("previous").IsValid())
	assert.True(t, lwtime.NatTimeAdjective("last").IsValid())
	assert.False(t, lwtime.NatTimeAdjective("").IsValid())
	assert.False(t, lwtime.NatTimeAdjective("invalid").IsValid())
}

type NatTimeLoadRTUTest struct {
	Name     string
	Input    string
	Expected lwtime.NatTime
	OK       bool
}

var NatTimeLoadRTUTests = []NatTimeLoadRTUTest{
	NatTimeLoadRTUTest{
		"year",
		"year",
		lwtime.NatTime{Unit: lwtime.Year},
		true,
	},
	NatTimeLoadRTUTest{
		"month",
		"MONTH",
		lwtime.NatTime{Unit: lwtime.Month}, // case-sensitivity test
		true,
	},
	NatTimeLoadRTUTest{
		"week",
		"weeks",
		lwtime.NatTime{Unit: lwtime.Week}, // plural test
		true,
	},
	NatTimeLoadRTUTest{
		"day",
		"day",
		lwtime.NatTime{Unit: lwtime.Day},
		true,
	},
	NatTimeLoadRTUTest{
		"hour",
		"hour",
		lwtime.NatTime{Unit: lwtime.Hour},
		true,
	},
	NatTimeLoadRTUTest{
		"minute",
		"Minute",
		lwtime.NatTime{Unit: lwtime.Minute},
		true,
	},
	NatTimeLoadRTUTest{
		"second",
		"second",
		lwtime.NatTime{Unit: lwtime.Second},
		true,
	},
	NatTimeLoadRTUTest{
		"empty",
		"",
		lwtime.NatTime{},
		false,
	},
	NatTimeLoadRTUTest{
		"invalid",
		"invalid",
		lwtime.NatTime{},
		false,
	},
}

func TestNatTimeLoadRTU(t *testing.T) {
	for _, ntlrtut := range NatTimeLoadRTUTests {
		t.Run(ntlrtut.Name, func(t *testing.T) {
			nt := lwtime.NatTime{}
			ok := nt.LoadRelTimeUnit(ntlrtut.Input)
			assert.Equal(t, ntlrtut.OK, ok)
			assert.Equal(t, ntlrtut.Expected, nt)
		})
	}
}

type NatTimeParseTest struct {
	Name     string
	Input    string
	Expected lwtime.NatTime
	Error    error
}

var NatTimeParseTests = []NatTimeParseTest{
	NatTimeParseTest{
		"today",
		"TODAY",
		lwtime.NatTime{lwtime.This, "1", 1, lwtime.Day},
		nil,
	},
	NatTimeParseTest{
		"yesterday",
		"yesterday",
		lwtime.NatTime{lwtime.Previous, "1", 1, lwtime.Day},
		nil,
	},
	NatTimeParseTest{
		"this",
		"this year",
		lwtime.NatTime{lwtime.This, "1", 1, lwtime.Year},
		nil,
	},
	NatTimeParseTest{
		"current",
		"current week",
		lwtime.NatTime{lwtime.Current, "1", 1, lwtime.Week},
		nil,
	},
	NatTimeParseTest{
		"previous",
		"previous month",
		lwtime.NatTime{lwtime.Previous, "1", 1, lwtime.Month},
		nil,
	},
	NatTimeParseTest{
		"previous-invalid",
		"previous 7 months",
		lwtime.NatTime{},
		errors.New("natural time (previous 7 months) is invalid"),
	},
	NatTimeParseTest{
		"last-singular",
		"last month",
		lwtime.NatTime{lwtime.Last, "1", 1, lwtime.Month},
		nil,
	},
	NatTimeParseTest{
		"last-plural",
		"last 8 hours",
		lwtime.NatTime{lwtime.Last, "8", 8, lwtime.Hour},
		nil,
	},
	NatTimeParseTest{
		"empty",
		"",
		lwtime.NatTime{},
		errors.New("natural time () is invalid"),
	},
	NatTimeParseTest{
		"invalid",
		"random garbage",
		lwtime.NatTime{},
		errors.New("natural time (random garbage) is invalid"),
	},
}

func TestNatTimeParse(t *testing.T) {
	for _, ntp := range NatTimeParseTests {
		t.Run(ntp.Name, func(t *testing.T) {
			nt := lwtime.NatTime{}
			err := nt.Parse(ntp.Input)

			if ntp.Error == nil {
				assert.Equal(t, ntp.Expected, nt)
				return
			}
			assert.Equal(t, ntp.Error.Error(), err.Error())
		})
	}
}

type NatTimeRangeParseTest struct {
	Name     string
	Input    string
	Expected lwtime.NatTimeRange
	Error    error
}

var NatTimeRangeParseTests = []NatTimeRangeParseTest{
	NatTimeRangeParseTest{
		"today",
		"today",
		lwtime.NatTimeRange{"@d", "now"},
		nil,
	},
	NatTimeRangeParseTest{
		"yesterday",
		"yesterday",
		lwtime.NatTimeRange{"-1d@d", "@d"},
		nil,
	},
	NatTimeRangeParseTest{
		"this",
		"this year",
		lwtime.NatTimeRange{"@y", "now"},
		nil,
	},
	NatTimeRangeParseTest{
		"current",
		"current week",
		lwtime.NatTimeRange{"@w", "now"},
		nil,
	},
	NatTimeRangeParseTest{
		"previous",
		"previous month",
		lwtime.NatTimeRange{"-1mon@mon", "@mon"},
		nil,
	},
	NatTimeRangeParseTest{
		"last-singular",
		"last day",
		lwtime.NatTimeRange{"-1d", "now"},
		nil,
	},
	NatTimeRangeParseTest{
		"last-plural",
		"last 8 hours",
		lwtime.NatTimeRange{"-8h", "now"},
		nil,
	},
	NatTimeRangeParseTest{
		"invalid",
		"random garbage",
		lwtime.NatTimeRange{},
		errors.New("natural time (random garbage) is invalid"),
	},
}

func TestNatTimeRangeParse(t *testing.T) {
	for _, ntrpt := range NatTimeRangeParseTests {
		t.Run(ntrpt.Name, func(t *testing.T) {
			ntr := lwtime.NatTimeRange{}
			err := ntr.Parse(ntrpt.Input)

			if ntrpt.Error == nil {
				assert.Equal(t, ntrpt.Expected, ntr)
				return
			}
			assert.Equal(t, ntrpt.Error.Error(), err.Error())
		})
	}
}

type NatTimeRangeRange_Test struct {
	Name                string
	ReferenceTime       string
	ReferenceTimeFormat string
	Input               string
	Duration            int64
	Error               error
}

var NatTimeRangeRange_Tests = []NatTimeRangeRange_Test{
	NatTimeRangeRange_Test{
		"today",
		"2021-01-01T00:00:01Z",
		time.RFC3339,
		"today",
		1,
		nil,
	},
	NatTimeRangeRange_Test{
		"yesterday",
		"2021-01-01T00:00:00Z",
		time.RFC3339,
		"yesterday",
		86400,
		nil,
	},
	NatTimeRangeRange_Test{
		"this",
		"2021-01-01T23:59:59Z",
		time.RFC3339,
		"this day",
		86399,
		nil,
	},
	NatTimeRangeRange_Test{
		"current",
		"2021-01-01T23:59:59Z",
		time.RFC3339,
		"current month",
		86399,
		nil,
	},
	NatTimeRangeRange_Test{
		"previous",
		"2021-01-01T23:59:59Z",
		time.RFC3339,
		"previous day",
		86400,
		nil,
	},
	NatTimeRangeRange_Test{
		"empty",
		"2021-01-01T23:59:59Z",
		time.RFC3339,
		"",
		0,
		errors.New("unable to compute natural time range: invalid relative time () for range start"),
	},
	NatTimeRangeRange_Test{
		"invalid",
		"2021-01-01T23:59:59Z",
		time.RFC3339,
		"random garbage",
		0,
		errors.New("unable to compute natural time range: invalid relative time () for range start"),
	},
}

func TestNatTimeRangeRange_(t *testing.T) {
	for _, ntrr_t := range NatTimeRangeRange_Tests {
		t.Run(ntrr_t.Name, func(t *testing.T) {
			refTime, err := time.Parse(ntrr_t.ReferenceTimeFormat, ntrr_t.ReferenceTime)
			if err != nil {
				assert.FailNow(t, "unable to parse reference time")
			}

			ntr := lwtime.NatTimeRange{}
			ntr.Parse(ntrr_t.Input)
			fmt.Println(ntr)
			start, end, err := ntr.Range_(refTime)

			if ntrr_t.Error == nil {
				d := end.Unix() - start.Unix()
				assert.Equal(t, ntrr_t.Duration, d)
				return
			}
			assert.Equal(t, ntrr_t.Error.Error(), err.Error())
		})
	}
}
