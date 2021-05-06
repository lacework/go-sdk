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

type naturalAdjective string

const (
	// both plural and singular regular expressions should contain 3 capture groups
	naturalPluralRE   = `^(last)\s(\d+)\s(years|months|weeks|days|hours|minutes|seconds)$`
	naturalSingularRE = `^(this|current|previous|last)(\s)(year|month|week|day|hour|minute|second)$`

	Today     naturalAdjective = "today"
	Yesterday naturalAdjective = "yesterday"
	This      naturalAdjective = "this"
	Current   naturalAdjective = "current"
	Previous  naturalAdjective = "previous"
	Last      naturalAdjective = "last"
)

func (nta naturalAdjective) isValid() bool {
	switch naturalAdjective(strings.ToLower(string(nta))) { // inline lowercase conversion
	case Today, Yesterday, This, Current, Previous, Last:
		return true
	}
	return false
}

type natural struct {
	adjective naturalAdjective
	num       string
	iNum      int
	unit      relativeUnit
}

func (nt *natural) loadRelativeUnit(u string) bool {
	switch strings.ToLower(u) {
	case "year", "years":
		nt.unit = Year
	case "month", "months":
		nt.unit = Month
	case "week", "weeks":
		nt.unit = Week
	case "day", "days":
		nt.unit = Day
	case "hour", "hours":
		nt.unit = Hour
	case "minute", "minutes":
		nt.unit = Minute
	case "second", "seconds":
		nt.unit = Second
	default:
		return false
	}
	return true
}

func newNatural(s string) (natural, error) {
	nt := natural{}
	s = strings.ToLower(s)

	// Today
	if naturalAdjective(s) == Today {
		nt.adjective = This
		nt.num = "1"
		nt.iNum = 1
		nt.unit = relativeUnit("d")
		return nt, nil
	}
	// Yesterday
	if naturalAdjective(s) == Yesterday {
		nt.adjective = Previous
		nt.num = "1"
		nt.iNum = 1
		nt.unit = relativeUnit("d")
		return nt, nil
	}
	// Singular
	var nt_parts []string
	singularRE := regexp.MustCompile(naturalSingularRE)
	if nt_parts = singularRE.FindStringSubmatch(s); s == "" || nt_parts == nil {
		// Plural
		pluralRE := regexp.MustCompile(naturalPluralRE)
		if nt_parts = pluralRE.FindStringSubmatch(s); s == "" || nt_parts == nil {
			return nt, errors.New(fmt.Sprintf("natural time (%s) is invalid", s))
		}
	}
	// Adjective
	nt.adjective = naturalAdjective(nt_parts[1])
	if !nt.adjective.isValid() {
		// this would indicate a code mismatch between enumerated Adjectives and Regex
		return nt, errors.New(fmt.Sprintf("invalid adjective for natural time (%s)", s))
	}
	// Num
	nt.num = nt_parts[2]
	var err error
	nt.iNum, err = strconv.Atoi(nt.num)
	if err != nil {
		nt.num = "1"
		nt.iNum = 1
	}
	// Unit
	if ok := nt.loadRelativeUnit(nt_parts[3]); !ok {
		return nt, errors.New(fmt.Sprintf("invalid unit for natural time (%s)", s))
	}
	return nt, nil
}

func (nt natural) range_(t time.Time) (start time.Time, end time.Time, err error) {
	var relStart, relEnd string
	baseErr := "unable to compute natural time range"

	// relatives
	if relStart, relEnd, err = nt.getRelativeRange(); err != nil {
		err = errors.Wrap(err, baseErr)
		return
	}

	// start time
	if start, err = parseRelativeFromTime(relStart, t); err != nil {
		err = errors.Wrap(err, baseErr)
		return
	}

	// end time
	if end, err = parseRelativeFromTime(relEnd, t); err != nil {
		err = errors.Wrap(err, baseErr)
		return
	}

	return
}

func (nt natural) getRelativeRange() (relStart string, relEnd string, err error) {
	// use natural adjective to determine relative start/end specifiers
	switch nt.adjective {
	case This, Current:
		relStart = fmt.Sprintf("@%s", nt.unit)
		relEnd = "now"
	case Previous:
		relStart = fmt.Sprintf("-1%s@%s", nt.unit, nt.unit)
		relEnd = fmt.Sprintf("@%s", nt.unit)
	case Last:
		relStart = fmt.Sprintf("-%s%s", nt.num, nt.unit)
		relEnd = "now"
	default:
		err = errors.New("invalid adjective for natural time")
	}
	return
}

// Parse the string representation of a Lacework natural time
//
// start, end, err := lwtime.ParseNatural("this year")
// if err != nil {
// 	...
// }
func ParseNatural(n string) (time.Time, time.Time, error) {
	return parseNaturalFromTime(n, time.Now())
}

func parseNaturalFromTime(n string, fromTime time.Time) (time.Time, time.Time, error) {
	var start, end time.Time

	natural, err := newNatural(n)
	if err != nil {
		return start, end, err
	}

	return natural.range_(fromTime)
}
