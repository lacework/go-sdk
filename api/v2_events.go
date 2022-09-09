//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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
	"time"
)

// EventService is the service that interacts with
// the Events schema from the Lacework APIv2 Server
type EventService struct {
	client *Client
}

// Search returns a list of Events matching the specified criteria
// e.g.
//
//   var (
//       now      = time.Now().UTC()
//       before   = now.AddDate(0, 0, -7) // 7 days from ago
//       filters  = api.SearchFilter{
//           TimeFilter: &api.TimeFilter{
//               StartTime: &before,
//               EndTime:   &now,
//           },
//       }
//   )
//   lacework.V2.Events.Search(response, filters)
//
func (svc *EventService) Search(filters SearchFilter) (response EventResponse, err error) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2EventsSearch, filters, &response)
	return
}

type EventResponse struct {
	Data   []V2Event    `json:"data"`
	Paging V2Pagination `json:"paging"`
}

// Fulfill Pageable interface (look at api/v2.go)
func (r EventResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *EventResponse) ResetPaging() {
	r.Paging = V2Pagination{}
}

type V2Event struct {
	EndTime    time.Time `json:"endTime"`
	EventCount int       `json:"eventCount"`
	EventType  string    `json:"eventType"`
	Id         int       `json:"id"`
	SrcEvent   any       `json:"srcEvent"`
	SrcType    string    `json:"srcType"`
	StartTime  time.Time `json:"startTime"`
}

type V2EventFile struct {
	FilePath     string `json:"file_path"`
	FileDataHash string `json:"filedata_hash"`
	Hostname     string `json:"hostname"`
	Mid          int    `json:"mid"`
	Username     string `json:"username"`
}
