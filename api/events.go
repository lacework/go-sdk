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
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// EventsService is a service that interacts with the Events endpoints
// from the Lacework Server
type EventsService struct {
	client *Client
}

// List leverages ListRange and returns a list of events from the last 7 days
func (svc *EventsService) List() (EventsResponse, error) {
	var (
		now  = time.Now().UTC()
		from = now.AddDate(0, 0, -7) // 7 days from now
	)

	return svc.ListRange(from, now)
}

// ListRange returns a list of Lacework events during the specified date range.
//
// Requirements and specifications:
// * The dates format should be: yyyy-MM-ddTHH:mm:ssZ (example 2019-07-11T21:11:00Z)
// * The START_TIME and END_TIME must be specified in UTC
// * The difference between the START_TIME and END_TIME must not be greater than 7 days
// * The START_TIME must be less than or equal to three months from current date
// * The number of records produced is limited to 5000
func (svc *EventsService) ListRange(start, end time.Time) (
	response EventsResponse,
	err error,
) {
	if start.After(end) {
		err = errors.New("data range should have a start time before the end time")
		return
	}

	apiPath := fmt.Sprintf(
		"%s?START_TIME=%s&END_TIME=%s",
		apiEventsDateRange,
		start.UTC().Format(time.RFC3339),
		end.UTC().Format(time.RFC3339),
	)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

type EventsResponse struct {
	Events []Event `json:"data"`
}

type Event struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Severity  string    `json:"severity"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (e *Event) SeverityString() string {
	switch e.Severity {
	case "1":
		return "Critical"
	case "2":
		return "High"
	case "3":
		return "Medium"
	case "4":
		return "Low"
	case "5":
		return "Info"
	default:
		return "Unknown"
	}
}

type EventsCount struct {
	Critical int
	High     int
	Medium   int
	Low      int
	Info     int
	Total    int
}

func (er *EventsResponse) GetEventsCount() EventsCount {
	counts := EventsCount{}
	for _, e := range er.Events {
		switch e.Severity {
		case "1":
			counts.Critical++
		case "2":
			counts.High++
		case "3":
			counts.Medium++
		case "4":
			counts.Low++
		case "5":
			counts.Info++
		}
		counts.Total++
	}
	return counts
}
