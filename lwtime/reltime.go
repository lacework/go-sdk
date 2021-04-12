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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type RelDate struct {
	Year  int
	Month time.Month
	Day   int
}

func Mondays(year int) (mondays []RelDate) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())
	d, _ := time.ParseDuration("24h")
	startYear := start.Year()
	for startYear == start.Year() {
		if start.Weekday() == time.Monday {
			year, month, day := start.Date()
			mondays = append(mondays, RelDate{year, month, day})
		}
		start = start.Add(d)
	}
	return mondays
}

type RelTimeUnit string

const (
	relTimeRE string      = `^([+-])?(?:(\d+)(\w+))?(?:@(\w+))?$`
	Year      RelTimeUnit = "y"
	Month     RelTimeUnit = "mon"
	Week      RelTimeUnit = "w"
	Day       RelTimeUnit = "d"
	Hour      RelTimeUnit = "h"
	Minute    RelTimeUnit = "m"
	Second    RelTimeUnit = "s"
)

func (rtu RelTimeUnit) IsValid() bool {
	switch RelTimeUnit(strings.ToLower(string(rtu))) { // inline lowercase conversion
	case Year, Month, Week, Day, Hour, Minute, Second, RelTimeUnit(""):
		return true
	}
	return false
}

func (rtu RelTimeUnit) SnapTime(inTime time.Time) (outTime time.Time, err error) {
	// immediately short circuit if snap is invalid
	if !rtu.IsValid() {
		err = errors.New(fmt.Sprintf(
			"snap (%s) is not a valid relative time unit", rtu))
		return
	}

	year, month, day := inTime.Date()
	hour := inTime.Hour()
	minute := inTime.Minute()
	second := inTime.Second()
	nano := inTime.Nanosecond()

	switch RelTimeUnit(strings.ToLower(string(rtu))) {
	case Week:
		year, week := inTime.ISOWeek()
		relDate := Mondays(year)[week-1]
		outTime = time.Date(
			relDate.Year,
			relDate.Month,
			relDate.Day,
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

type RelTime struct {
	Num  string
	INum int
	Unit RelTimeUnit
	Snap RelTimeUnit
}

func (rt *RelTime) Parse(s string) error {
	var rt_parts []string
	// now is equivelant to +0s
	if s == "now" {
		s = "+0s"
	}
	// regex
	re := regexp.MustCompile(relTimeRE)
	if rt_parts = re.FindStringSubmatch(s); s == "" || rt_parts == nil {
		return errors.New(fmt.Sprintf("relative time specifier (%s) is invalid", s))
	}
	// Num
	if rt_parts[1] == "-" {
		rt.Num = rt_parts[1] + rt_parts[2]
	} else {
		rt.Num = rt_parts[2]
	}
	var err error
	rt.INum, err = strconv.Atoi(rt.Num)
	if err != nil {
		rt.Num = "0"
		rt.INum = 0
	}
	// Unit
	rt.Unit = RelTimeUnit(strings.ToLower(string(rt_parts[3])))
	if !rt.Unit.IsValid() {
		return errors.New(fmt.Sprintf("invalid unit for relative time specifier (%s)", s))
	}
	if rt.Unit == RelTimeUnit("") {
		rt.Unit = Second
	}
	// Weeeeeeek
	if rt.Unit == Week {
		rt.INum = rt.INum * 7
		rt.Num = strconv.Itoa(rt.INum)
		rt.Unit = Day
	}
	// Snap
	rt.Snap = RelTimeUnit(strings.ToLower(string(rt_parts[4])))
	if !rt.Snap.IsValid() {
		return errors.New(fmt.Sprintf("invalid snap for relative time specifier (%s)", s))
	}
	return nil
}

func (rt RelTime) Time() (outTime time.Time, err error) {
	outTime, err = rt.Time_(time.Now())
	return
}

func (rt RelTime) Time_(inTime time.Time) (outTime time.Time, err error) {
	baseErr := "unable to construct time object"
	switch rt.Unit {
	case Year:
		outTime = inTime.AddDate(rt.INum, 0, 0)
	case Month:
		outTime = inTime.AddDate(0, rt.INum, 0)
	case Day:
		outTime = inTime.AddDate(0, 0, rt.INum)
	case Hour, Minute, Second:
		var d time.Duration
		d, err = time.ParseDuration(fmt.Sprintf("%s%s", rt.Num, rt.Unit))
		if err != nil {
			return
		}
		outTime = inTime.Add(d)
	default:
		err = errors.Wrap(
			errors.New(fmt.Sprintf("relative time unit (%s) is invalid", rt.Unit)),
			baseErr,
		)
		return
	}
	if rt.Snap != "" {
		outTime, err = rt.Snap.SnapTime(outTime)
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
