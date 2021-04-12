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

type NatTimeAdjective string

const (
	// both plural and singular regular expressions should contain 3 capture groups
	natTimePluralRE   string = `^(last)\s(\d+)\s(years|months|weeks|days|hours|minutes|seconds)$`
	natTimeSingularRE string = `^(this|current|previous|last)(\s)(year|month|week|day|hour|minute|second)$`

	Today     NatTimeAdjective = "today"
	Yesterday NatTimeAdjective = "yesterday"
	This      NatTimeAdjective = "this"
	Current   NatTimeAdjective = "current"
	Previous  NatTimeAdjective = "previous"
	Last      NatTimeAdjective = "last"
)

func (nta NatTimeAdjective) IsValid() bool {
	switch NatTimeAdjective(strings.ToLower(string(nta))) { // inline lowercase conversion
	case Today, Yesterday, This, Current, Previous, Last:
		return true
	}
	return false
}

type NatTime struct {
	Adjective NatTimeAdjective
	Num       string
	INum      int
	Unit      RelTimeUnit
}

func (nt *NatTime) LoadRelTimeUnit(u string) bool {
	switch strings.ToLower(u) {
	case "year", "years":
		nt.Unit = Year
	case "month", "months":
		nt.Unit = Month
	case "week", "weeks":
		nt.Unit = Week
	case "day", "days":
		nt.Unit = Day
	case "hour", "hours":
		nt.Unit = Hour
	case "minute", "minutes":
		nt.Unit = Minute
	case "second", "seconds":
		nt.Unit = Second
	default:
		return false
	}
	return true
}

func (nt *NatTime) Parse(s string) error {
	s = strings.ToLower(s)
	// Today
	if NatTimeAdjective(s) == Today {
		nt.Adjective = This
		nt.Num = "1"
		nt.INum = 1
		nt.Unit = RelTimeUnit("d")
		return nil
	}
	// Yesterday
	if NatTimeAdjective(s) == Yesterday {
		nt.Adjective = Previous
		nt.Num = "1"
		nt.INum = 1
		nt.Unit = RelTimeUnit("d")
		return nil
	}
	// Singular
	var nt_parts []string
	singularRE := regexp.MustCompile(natTimeSingularRE)
	if nt_parts = singularRE.FindStringSubmatch(s); s == "" || nt_parts == nil {
		// Plural
		pluralRE := regexp.MustCompile(natTimePluralRE)
		if nt_parts = pluralRE.FindStringSubmatch(s); s == "" || nt_parts == nil {
			return errors.New(fmt.Sprintf("natural time (%s) is invalid", s))
		}
	}
	// Adjective
	nt.Adjective = NatTimeAdjective(nt_parts[1])
	if !nt.Adjective.IsValid() {
		// this would indicate a code mismatch between enumerated Adjectives and Regex
		return errors.New(fmt.Sprintf("invalid adjective for natural time (%s)", s))
	}
	// Num
	nt.Num = nt_parts[2]
	var err error
	nt.INum, err = strconv.Atoi(nt.Num)
	if err != nil {
		nt.Num = "1"
		nt.INum = 1
	}
	// Unit
	if ok := nt.LoadRelTimeUnit(nt_parts[3]); !ok {
		return errors.New(fmt.Sprintf("invalid unit for natural time (%s)", s))
	}
	return nil
}

type NatTimeRange struct {
	Start string
	End   string
}

func (ntr *NatTimeRange) Parse(s string) error {
	nt := NatTime{}
	if err := nt.Parse(s); err != nil {
		return err
	}

	switch nt.Adjective {
	case This, Current:
		ntr.Start = fmt.Sprintf("@%s", nt.Unit)
		ntr.End = "now"
	case Previous:
		ntr.Start = fmt.Sprintf("-1%s@%s", nt.Unit, nt.Unit)
		ntr.End = fmt.Sprintf("@%s", nt.Unit)
	case Last:
		ntr.Start = fmt.Sprintf("-%s%s", nt.Num, nt.Unit)
		ntr.End = "now"
	default:
		// This would indicate a code mismatch between NatTime.Parse() and NatTimeRange.Parse()
		return errors.New(
			fmt.Sprintf("invalid adjective for natural time (%s)", s))
	}
	return nil
}

func (ntr NatTimeRange) Range() (start time.Time, end time.Time, err error) {
	start, end, err = ntr.Range_(time.Now())
	return
}

func (ntr NatTimeRange) Range_(t time.Time) (start time.Time, end time.Time, err error) {
	baseErr := "unable to compute natural time range"
	rt := RelTime{}

	// start
	if err = rt.Parse(ntr.Start); err != nil {
		err = errors.Wrap(
			errors.New(fmt.Sprintf(
				"invalid relative time (%s) for range start", ntr.Start)),
			baseErr,
		)
		return
	}
	if start, err = rt.Time_(t); err != nil {
		err = errors.Wrap(err, baseErr)
		return
	}
	// end
	if err = rt.Parse(ntr.End); err != nil {
		err = errors.Wrap(
			errors.New(fmt.Sprintf(
				"invalid relative time (%s) for range end", ntr.End)),
			baseErr,
		)
		return
	}
	if end, err = rt.Time_(t); err != nil {
		return
	}
	return
}
