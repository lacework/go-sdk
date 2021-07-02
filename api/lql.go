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
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/lacework/go-sdk/lwtime"
)

type Query struct {
	ID             string                   `json:"queryId,omitempty" yaml:"queryId,omitempty"`
	QueryText      string                   `json:"queryText,omitempty" yaml:"queryText,omitempty"`
	EvaluatorID    string                   `json:"evaluatorId,omitempty" yaml:"evaluatorId,omitempty"`
	Owner          string                   `json:"owner,omitempty" yaml:"owner,omitempty"`
	LastUpdateTime string                   `json:"lastUpdateTime,omitempty" yaml:"lastUpdateTime,omitempty"`
	LastUpdateUser string                   `json:"lastUpdateUser,omitempty" yaml:"lastUpdateUser,omitempty"`
	ResultSchema   []map[string]interface{} `json:"resultSchema,omitempty" yaml:"resultSchema,omitempty"`
}

func ParseQuery(s string) (Query, error) {
	var query Query
	var err error

	// valid json
	if err = json.Unmarshal([]byte(s), &query); err == nil {
		return query, err
	}
	// valid yaml
	query = Query{}
	err = yaml.Unmarshal([]byte(s), &query)
	if err == nil && !reflect.DeepEqual(query, Query{}) { // empty string unmarshals w/o error
		return query, nil
	}
	// invalid policy
	return query, errors.New("query must be valid JSON or YAML")
}

func ParseQueryTime(inTime string) (time.Time, error) {
	// empty
	if inTime == "" {
		return time.Time{}, errors.New(fmt.Sprintf("unable to parse time (%s)", inTime))
	}
	// parse time as relative
	if t, err := lwtime.ParseRelative(inTime); err == nil {
		return t, err
	}
	// parse time as RFC3339
	if t, err := time.Parse(time.RFC3339, inTime); err == nil {
		return t, err
	}
	// parse time as millis
	if i, err := strconv.ParseInt(inTime, 10, 64); err == nil {
		return time.Unix(0, i*int64(time.Millisecond)), err
	}
	return time.Time{}, errors.New(fmt.Sprintf("unable to parse time (%s)", inTime))
}

func ValidateQueryRange(start, end time.Time) (err error) {
	// validate range
	if start.After(end) {
		err = errors.New("date range should have a start time before the end time")
		return
	}
	return nil
}

type QueryResponse struct {
	Data    Query  `json:"data"`
	Message string `json:"message"`
}

type LQLQueriesResponse struct {
	Data    []Query `json:"data"`
	Message string  `json:"message"`
}

// LQLService is a service that interacts with the LQL
// endpoints from the Lacework Server
type LQLService struct {
	client *Client
}

func (svc *LQLService) Create(q string) (
	response QueryResponse,
	err error,
) {
	var query Query
	if query, err = ParseQuery(q); err != nil {
		return
	}
	err = svc.client.RequestEncoderDecoder("POST", apiV2LQL, query, &response)
	return
}

func (svc *LQLService) Update(q string) (
	response QueryResponse,
	err error,
) {
	var query Query
	if query, err = ParseQuery(q); err != nil {
		return
	}
	err = svc.client.RequestEncoderDecoder(
		"PATCH",
		fmt.Sprintf("%s/%s", apiV2LQL, url.QueryEscape(query.ID)),
		query,
		&response,
	)
	return
}

func (svc *LQLService) GetQueries() (
	response LQLQueriesResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET", apiV2LQL, nil, &response)
	return
}

func (svc *LQLService) GetByID(queryID string) (
	response QueryResponse,
	err error,
) {
	if queryID == "" {
		err = errors.New("query ID must be provided")
		return
	}
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf("%s/%s", apiV2LQL, url.QueryEscape(queryID)),
		nil,
		&response,
	)
	return
}

func (svc *LQLService) Run(q, start, end string) (
	response map[string]interface{},
	err error,
) {
	var v2Query Query
	if v2Query, err = ParseQuery(q); err != nil {
		return
	}

	var startTimeRange, endTimeRange time.Time
	if startTimeRange, err = ParseQueryTime(start); err != nil {
		return
	}
	if endTimeRange, err = ParseQueryTime(end); err != nil {
		return
	}
	if err = ValidateQueryRange(startTimeRange, endTimeRange); err != nil {
		return
	}

	v1Query := map[string]string{
		"query_text":       v2Query.QueryText,
		"start_time_range": startTimeRange.UTC().Format(time.RFC3339),
		"end_time_range":   endTimeRange.UTC().Format(time.RFC3339),
	}
	err = svc.client.RequestEncoderDecoder("POST", apiLQLQuery, v1Query, &response)
	return
}
