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

func TestvalidateQueryRange(t *testing.T) {
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
