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
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	reLQL                  string = `(?ms)^(\w+)\([^)]+\)\s*{`
	LQLQueryTranslateError string = "unable to translate query blob"
)

type LQLQuery struct {
	ID             string `json:"LQL_ID,omitempty"`
	StartTimeRange string `json:"START_TIME_RANGE,omitempty"`
	EndTimeRange   string `json:"END_TIME_RANGE,omitempty"`
	QueryText      string `json:"QUERY_TEXT"`
	// QueryBlob is a special string that supports type conversion
	// back and forth from LQL to JSON
	QueryBlob string `json:"-"`
}

func (q *LQLQuery) Validate(allowEmptyTimes bool) error {
	// translate
	if err := q.Translate(); err != nil {
		return err
	}
	// validate range
	if err := q.ValidateRange(allowEmptyTimes); err != nil {
		return err
	}
	// validate query
	if q.QueryText == "" {
		return errors.New("query should not be empty")
	}
	return nil
}

func (q *LQLQuery) Translate() error {
	// query
	if err := q.TranslateQuery(); err != nil {
		return err
	}
	// start
	if start, err := q.TranslateTime(q.StartTimeRange); err != nil {
		return err
	} else {
		q.StartTimeRange = start
	}
	// end
	if end, err := q.TranslateTime(q.EndTimeRange); err != nil {
		return nil
	} else {
		q.EndTimeRange = end
	}
	return nil
}

func (q *LQLQuery) TranslateQuery() error {
	// empty
	if q.QueryText != "" {
		return nil
	}
	// json
	var t LQLQuery

	if err := json.Unmarshal([]byte(q.QueryBlob), &t); err == nil {
		if q.StartTimeRange == "" {
			q.StartTimeRange = t.StartTimeRange
		}
		if q.EndTimeRange == "" {
			q.EndTimeRange = t.EndTimeRange
		}
		q.QueryText = t.QueryText
		return err
	}
	// lql
	if matched, _ := regexp.MatchString(reLQL, q.QueryBlob); matched {
		q.QueryText = q.QueryBlob
		return nil
	}
	return errors.New(LQLQueryTranslateError)
}

func (q LQLQuery) TranslateTime(inTime string) (outTime string, err error) {
	// empty
	if inTime == "" {
		return
	}
	// parse time as RFC3339
	var t time.Time
	if t, err = time.Parse(time.RFC3339, inTime); err == nil {
		outTime = t.UTC().Format(time.RFC3339)
		return
	}
	// parse time as millis
	var i int64
	if i, err = strconv.ParseInt(inTime, 10, 64); err == nil {
		outTime = time.Unix(0, i*int64(time.Millisecond)).UTC().Format(time.RFC3339)
		return
	}
	err = errors.New("unable to parse time (" + inTime + ")")
	return
}

func (q LQLQuery) ValidateRange(allowEmptyTimes bool) (err error) {
	// validate start
	var start time.Time
	if q.StartTimeRange != "" {
		if start, err = time.Parse(time.RFC3339, q.StartTimeRange); err != nil {
			return
		}
	} else if allowEmptyTimes {
		start = time.Unix(0, 0)
	} else {
		err = errors.New("start time must not be empty")
		return
	}
	// validate end
	var end time.Time
	if q.EndTimeRange != "" {
		if end, err = time.Parse(time.RFC3339, q.EndTimeRange); err != nil {
			return
		}
	} else if allowEmptyTimes {
		end = time.Now()
	} else {
		err = errors.New("end time must not be empty")
		return
	}
	// validate range
	if start.After(end) {
		err = errors.New("date range should have a start time before the end time")
		return
	}
	return nil
}

type LQLQueryResponse struct {
	Data    []LQLQuery `json:"data"`
	Ok      bool       `json:"ok"`
	Message string     `json:"message"`
}

// LQLService is a service that interacts with the LQL
// endpoints from the Lacework Server
type LQLService struct {
	client *Client
}

func (svc *LQLService) CreateQuery(query string) (
	response LQLQueryResponse,
	err error,
) {
	lqlQuery := LQLQuery{QueryBlob: query}
	if err = lqlQuery.Validate(true); err != nil {
		return
	}

	err = svc.client.RequestEncoderDecoder("POST", ApiLQL, lqlQuery, &response)
	return
}

func (svc *LQLService) GetQueries() (LQLQueryResponse, error) {
	return svc.GetQueryByID("")
}

func (svc *LQLService) GetQueryByID(queryID string) (
	response LQLQueryResponse,
	err error,
) {
	uri := ApiLQL

	if queryID != "" {
		uri += "?LQL_ID=" + url.QueryEscape(queryID)
	}

	err = svc.client.RequestDecoder("GET", uri, nil, &response)
	return
}

func (svc *LQLService) RunQuery(query, start, end string) (
	response map[string]interface{},
	err error,
) {
	lqlQuery := LQLQuery{
		StartTimeRange: start,
		EndTimeRange:   end,
		QueryBlob:      query,
	}
	if err = lqlQuery.Validate(false); err != nil {
		return
	}

	err = svc.client.RequestEncoderDecoder("POST", ApiLQLQuery, lqlQuery, &response)
	return
}
