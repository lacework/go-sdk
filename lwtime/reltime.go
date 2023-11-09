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
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type relativeDate struct {
	year  int
	month time.Month
	day   int
}

// The potential difference between the clocks on the client
// and the Lacework API server
const clockOffset = "+0s"

// Returns 'now' with the default or the provided clock offset
func nowClockOffset() string {
	if os.Getenv("LW_CLOCK_OFFSET") != "" {
		return os.Getenv("LW_CLOCK_OFFSET")
	}
	return clockOffset
}

func mondays(year int) (mondays []relativeDate) {
	// get the start of the year datetime
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())
	// get 24 hours of duration
	d, _ := time.ParseDuration("24h")
	// set startYear for comparison as we iterate
	startYear := start.Year()
	// iterate until startYear deviates from start.Year()
	for startYear == start.Year() {
		// if we kave a monday, add it....
		if start.Weekday() == time.Monday {
			year, month, day := start.Date()
			mondays = append(mondays, relativeDate{year, month, day})
		}
		// add our 24 hour duration
		start = start.Add(d)
	}
	return mondays
}

type relativeUnit string

const (
	relativeRE               = `^([+-])?(?:(\d+)(\w+))?(?:@(\w+))?$`
	Year        relativeUnit = "y"
	Month       relativeUnit = "mon"
	Week        relativeUnit = "w"
	Day         relativeUnit = "d"
	Hour        relativeUnit = "h"
	Minute      relativeUnit = "m"
	Second      relativeUnit = "s"
	HoursInADay              = 24
)

func (ru relativeUnit) isValid() bool {
	switch relativeUnit(strings.ToLower(string(ru))) { // inline lowercase conversion
	case Year, Month, Week, Day, Hour, Minute, Second, relativeUnit(""):
		return true
	}
	return false
}

func (ru relativeUnit) snapTime(inTime time.Time) (outTime time.Time, err error) {
	// immediately short circuit if snap is invalid
	if !ru.isValid() {
		err = errors.New(fmt.Sprintf(
			"snap (%s) is not a valid relative time unit", ru))
		return
	}

	year, month, day := inTime.Date()
	hour := inTime.Hour()
	minute := inTime.Minute()
	second := inTime.Second()
	nano := inTime.Nanosecond()

	switch relativeUnit(strings.ToLower(string(ru))) {
	case Week:
		year, week := inTime.ISOWeek()
		relDate := mondays(year)[week-1]
		outTime = time.Date(
			relDate.year,
			relDate.month,
			relDate.day,
			0, 0, 0, 0, inTime.Location(),
		)
		return
	case Year:
		month = 1
		fallthrough
	case Month:
		day = 1
		fallthrough
	case Day:
		hour = 0
		fallthrough
	case Hour:
		minute = 0
		fallthrough
	case Minute:
		second = 0
		fallthrough
	case Second:
		nano = 0
	}
	outTime = time.Date(year, month, day, hour, minute, second, nano, inTime.Location())
	return
}

type relative struct {
	num  string
	iNum int
	unit relativeUnit
	snap relativeUnit
}

func newRelative(s string) (relative, error) {
	var rel relative
	var rel_parts []string

	// now is equivelant to LW_CLOCK_OFFSET (defaults to const clockOffset(+0s))
	// prevent corner conditions with Lacework's API server
	if s == "now" {
		s = nowClockOffset()
	}
	// regex
	re := regexp.MustCompile(relativeRE)
	if rel_parts = re.FindStringSubmatch(s); s == "" || rel_parts == nil {
		return rel, errors.New(fmt.Sprintf("relative time specifier (%s) is invalid", s))
	}
	// Num
	if rel_parts[1] == "-" {
		rel.num = rel_parts[1] + rel_parts[2]
	} else {
		rel.num = rel_parts[2]
	}
	var err error
	rel.iNum, err = strconv.Atoi(rel.num)
	if err != nil {
		rel.num = "0"
		rel.iNum = 0
	}
	// Unit
	rel.unit = relativeUnit(strings.ToLower(string(rel_parts[3])))
	if !rel.unit.isValid() {
		return rel, errors.New(fmt.Sprintf("invalid unit for relative time specifier (%s)", s))
	}
	// normalize years, weeks, and days in to hours
	switch rel.unit {
	case relativeUnit(""):
		rel.unit = Second
	case Week:
		rel.iNum = rel.iNum * 7 * HoursInADay
		rel.num = strconv.Itoa(rel.iNum)
		rel.unit = Hour
	case Day:
		rel.iNum = rel.iNum * HoursInADay
		rel.num = strconv.Itoa(rel.iNum)
		rel.unit = Hour
	}
	// Snap
	rel.snap = relativeUnit(strings.ToLower(string(rel_parts[4])))
	if !rel.snap.isValid() {
		return rel, errors.New(fmt.Sprintf("invalid snap for relative time specifier (%s)", s))
	}
	return rel, nil
}

func (rel relative) time(inTime time.Time) (outTime time.Time, err error) {
	baseErr := "unable to construct time object"

	switch rel.unit {
	case Year:
		outTime = inTime.AddDate(rel.iNum, 0, 0)
	case Month:
		outTime = inTime.AddDate(0, rel.iNum, 0)
	case Day:
		outTime = inTime.AddDate(0, 0, rel.iNum)
	case Hour, Minute, Second:
		var d time.Duration
		d, err = time.ParseDuration(fmt.Sprintf("%s%s", rel.num, rel.unit))
		if err != nil {
			return
		}
		outTime = inTime.Add(d)
	default:
		err = errors.Wrap(
			errors.New(fmt.Sprintf("relative time unit (%s) is invalid", rel.unit)),
			baseErr,
		)
		return
	}
	if rel.snap != "" {
		outTime, err = rel.snap.snapTime(outTime)
	}
	if err != nil {
		err = errors.Wrap(err, baseErr)
		return
	}
	if outTime.Unix() < 0 {
		err = errors.Wrap(errors.New("time predates epoch"), baseErr)
		return
	}
	return
}

// ParseRelative parses the string representation of a Lacework relative time.
// Time object is returned in UTC.
//
//	t, err := lwtime.ParseRelative("-1y@y")
//	if err != nil {
//		...
//	}
func ParseRelative(s string) (time.Time, error) {
	// time.Now() is intentional here such that snaps work properly
	// For instance snapping to @d should snap to the start of the local day
	localTime, err := parseRelativeFromTime(s, time.Now())
	return localTime.UTC(), err
}

func parseRelativeFromTime(s string, fromTime time.Time) (time.Time, error) {
	relative, err := newRelative(s)
	if err != nil {
		return time.Time{}, err
	}

	return relative.time(fromTime)
}
