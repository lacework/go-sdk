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
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNaturalAdjectiveIsValid(t *testing.T) {
	assert.True(t, naturalAdjective("today").isValid())
	assert.True(t, naturalAdjective("YESTERDAY").isValid()) // case-sensitivity test
	assert.True(t, naturalAdjective("this").isValid())
	assert.True(t, naturalAdjective("current").isValid())
	assert.True(t, naturalAdjective("previous").isValid())
	assert.True(t, naturalAdjective("last").isValid())
	assert.False(t, naturalAdjective("").isValid())
	assert.False(t, naturalAdjective("invalid").isValid())
}

type NaturalLoadRTUTest struct {
	Name     string
	Input    string
	Expected natural
	OK       bool
}

var NaturalLoadRTUTests = []NaturalLoadRTUTest{
	NaturalLoadRTUTest{
		"year",
		"year",
		natural{unit: Year},
		true,
	},
	NaturalLoadRTUTest{
		"month",
		"MONTH",
		natural{unit: Month}, // case-sensitivity test
		true,
	},
	NaturalLoadRTUTest{
		"week",
		"weeks",
		natural{unit: Week}, // plural test
		true,
	},
	NaturalLoadRTUTest{
		"day",
		"day",
		natural{unit: Day},
		true,
	},
	NaturalLoadRTUTest{
		"hour",
		"hour",
		natural{unit: Hour},
		true,
	},
	NaturalLoadRTUTest{
		"minute",
		"Minute",
		natural{unit: Minute},
		true,
	},
	NaturalLoadRTUTest{
		"second",
		"second",
		natural{unit: Second},
		true,
	},
	NaturalLoadRTUTest{
		"empty",
		"",
		natural{},
		false,
	},
	NaturalLoadRTUTest{
		"invalid",
		"invalid",
		natural{},
		false,
	},
}

func TestNaturalLoadRTU(t *testing.T) {
	for _, nlrtut := range NaturalLoadRTUTests {
		t.Run(nlrtut.Name, func(t *testing.T) {
			var nt natural
			ok := nt.loadRelativeUnit(nlrtut.Input)
			assert.Equal(t, nlrtut.OK, ok)
			assert.Equal(t, nlrtut.Expected, nt)
		})
	}
}

type NewNaturalTest struct {
	Name     string
	Input    string
	Expected natural
	Error    error
}

var NewNaturalTests = []NewNaturalTest{
	NewNaturalTest{
		"today",
		"TODAY",
		natural{This, "1", 1, Day},
		nil,
	},
	NewNaturalTest{
		"yesterday",
		"yesterday",
		natural{Previous, "1", 1, Day},
		nil,
	},
	NewNaturalTest{
		"this",
		"this year",
		natural{This, "1", 1, Year},
		nil,
	},
	NewNaturalTest{
		"current",
		"current week",
		natural{Current, "1", 1, Week},
		nil,
	},
	NewNaturalTest{
		"previous",
		"previous month",
		natural{Previous, "1", 1, Month},
		nil,
	},
	NewNaturalTest{
		"previous-invalid",
		"previous 7 months",
		natural{},
		errors.New("natural time (previous 7 months) is invalid"),
	},
	NewNaturalTest{
		"last-singular",
		"last month",
		natural{Last, "1", 1, Month},
		nil,
	},
	NewNaturalTest{
		"last-plural",
		"last 8 hours",
		natural{Last, "8", 8, Hour},
		nil,
	},
	NewNaturalTest{
		"empty",
		"",
		natural{},
		errors.New("natural time () is invalid"),
	},
	NewNaturalTest{
		"invalid",
		"random garbage",
		natural{},
		errors.New("natural time (random garbage) is invalid"),
	},
}

func TestNewNatural(t *testing.T) {
	for _, pnt := range NewNaturalTests {
		t.Run(pnt.Name, func(t *testing.T) {
			nt, err := newNatural(pnt.Input)

			if pnt.Error == nil {
				assert.Equal(t, pnt.Expected, nt)
				return
			}
			assert.Equal(t, pnt.Error.Error(), err.Error())
		})
	}
}

type GetRelativeRangeTest struct {
	Name          string
	Input         string
	ExpectedStart string
	ExpectedEnd   string
	Error         error
}

var GetRelativeRangeTests = []GetRelativeRangeTest{
	GetRelativeRangeTest{
		"today",
		"today",
		"@d",
		"now",
		nil,
	},
	GetRelativeRangeTest{
		"yesterday",
		"yesterday",
		"-1d@d",
		"@d",
		nil,
	},
	GetRelativeRangeTest{
		"this",
		"this year",
		"@y",
		"now",
		nil,
	},
	GetRelativeRangeTest{
		"current",
		"current week",
		"@w",
		"now",
		nil,
	},
	GetRelativeRangeTest{
		"previous",
		"previous month",
		"-1mon@mon",
		"@mon",
		nil,
	},
	GetRelativeRangeTest{
		"last-singular",
		"last day",
		"-1d",
		"now",
		nil,
	},
	GetRelativeRangeTest{
		"last-plural",
		"last 8 hours",
		"-8h",
		"now",
		nil,
	},
	GetRelativeRangeTest{
		"invalid",
		"random garbage",
		"",
		"",
		errors.New("invalid adjective for natural time"),
	},
}

func TestGetRelativeRange(t *testing.T) {
	for _, grr := range GetRelativeRangeTests {
		t.Run(grr.Name, func(t *testing.T) {
			nt, _ := newNatural(grr.Input)
			actualStart, actualEnd, err := nt.getRelativeRange()

			if grr.Error == nil {
				assert.Equal(t, grr.ExpectedStart, actualStart)
				assert.Equal(t, grr.ExpectedEnd, actualEnd)
				return
			}
			assert.Equal(t, grr.Error.Error(), err.Error())
		})
	}
}

type NaturalFromTimeTest struct {
	Name                string
	ReferenceTime       string
	ReferenceTimeFormat string
	Input               string
	Duration            int64
	Error               error
}

var NaturalFromTimeTests = []NaturalFromTimeTest{
	NaturalFromTimeTest{
		"today",
		"2021-01-01T00:00:01Z",
		time.RFC3339,
		"today",
		-29,
		nil,
	},
	NaturalFromTimeTest{
		"yesterday",
		"2021-01-01T00:00:00Z",
		time.RFC3339,
		"yesterday",
		86400,
		nil,
	},
	NaturalFromTimeTest{
		"this",
		"2021-01-01T23:59:59Z",
		time.RFC3339,
		"this day",
		86369,
		nil,
	},
	NaturalFromTimeTest{
		"current",
		"2021-01-01T23:59:59Z",
		time.RFC3339,
		"current month",
		86369,
		nil,
	},
	NaturalFromTimeTest{
		"previous",
		"2021-01-01T23:59:59Z",
		time.RFC3339,
		"previous day",
		86400,
		nil,
	},
	NaturalFromTimeTest{
		"empty",
		"2021-01-01T23:59:59Z",
		time.RFC3339,
		"",
		0,
		errors.New("unable to compute natural time range: invalid adjective for natural time"),
	},
	NaturalFromTimeTest{
		"invalid",
		"2021-01-01T23:59:59Z",
		time.RFC3339,
		"random garbage",
		0,
		errors.New("unable to compute natural time range: invalid adjective for natural time"),
	},
}

func TestNaturalFromTime(t *testing.T) {
	for _, nft := range NaturalFromTimeTests {
		t.Run(nft.Name, func(t *testing.T) {
			refTime, err := time.Parse(nft.ReferenceTimeFormat, nft.ReferenceTime)
			if err != nil {
				assert.FailNow(t, "unable to parse reference time")
			}

			nt, _ := newNatural(nft.Input)
			start, end, err := nt.getRange(refTime)

			if nft.Error == nil {
				d := end.Unix() - start.Unix()
				assert.Equal(t, nft.Duration, d)
				return
			}
			assert.Equal(t, nft.Error.Error(), err.Error())
		})
	}
}
