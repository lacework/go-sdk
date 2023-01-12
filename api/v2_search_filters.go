//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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
	"errors"
	"math"
	"time"
)

// SearchFilter is the representation of an advanced search payload
// for retrieving information out of the Lacework APIv2 Server
//
// An advanced example of a SearchFilter to search for an Agent
// Access Token that matches the provider token alias and return
// only the token found:
//
//		SearchFilter{
//			Filters: []Filter{
//				Filter{
//					Field:      "tokenAlias",
//					Expression: "eq",
//					Value:      "k8s-deployment,
//				},
//			},
//			Returns: []string{"accessToken"},
//		}
type SearchFilter struct {
	*TimeFilter `json:"timeFilter,omitempty"`
	Filters     []Filter `json:"filters,omitempty"`
	Returns     []string `json:"returns,omitempty"`
}

type Filter struct {
	Expression string   `json:"expression,omitempty"`
	Field      string   `json:"field,omitempty"`
	Value      string   `json:"value,omitempty"`
	Values     []string `json:"values,omitempty"`
}

type TimeFilter struct {
	StartTime *time.Time `json:"startTime,omitempty"`
	EndTime   *time.Time `json:"endTime,omitempty"`
}

type SearchResponse interface {
	GetDataLength() int
}

type SearchableFilter interface {
	GetTimeFilter() *TimeFilter
	SetStartTime(*time.Time)
	SetEndTime(*time.Time)
}

// V2ApiMaxSearchHistoryDays defines the maximum number of days in the past api v2 allows to be searched
const V2ApiMaxSearchHistoryDays = 92

// V2ApiMaxSearchWindowDays defines the maximum number of days in a single request api v2 allows to be searched
const V2ApiMaxSearchWindowDays = 7

type search func(response interface{}, filters SearchableFilter) error

// WindowedSearchFirst performs a new search of a specific time frame size,
// until response data is found or the max searchable days is reached
func WindowedSearchFirst(fn search, size int, max int, response SearchResponse, filter SearchableFilter) error {
	if size > max {
		return errors.New("window size cannot be greater than max history")
	}

	// if start and end time are the same, adjust the windows
	timeDifference := int(math.RoundToEven(filter.GetTimeFilter().EndTime.Sub(*filter.GetTimeFilter().StartTime).Hours() / 24))

	if timeDifference == 0 {
		newStart := filter.GetTimeFilter().StartTime.AddDate(0, 0, -size)
		filter.SetStartTime(&newStart)
	}

	for i := timeDifference; i < max; i += size {
		err := fn(&response, filter)
		if err != nil {
			return err
		}
		if response.GetDataLength() != 0 {
			return nil
		}

		// adjust window
		newStart := filter.GetTimeFilter().StartTime.AddDate(0, 0, -size)
		newEnd := filter.GetTimeFilter().EndTime.AddDate(0, 0, -size)

		// ensure we do not go over the max allowed searchable days
		searchableDays := time.Since(newStart).Hours() / 24
		if int(searchableDays) > max {
			newStart = time.Now().AddDate(0, 0, -max)
		}
		filter.SetStartTime(&newStart)
		filter.SetEndTime(&newEnd)
	}
	return nil
}
