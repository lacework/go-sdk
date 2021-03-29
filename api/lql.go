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
)

const (
	reLQL string = `(?ms)^(\w+)\(\w+\s\w+\)\s*{`
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

func (q *LQLQuery) translate() {
	// if QueryText is populated; return
	if q.QueryText != "" {
		return
	}

	// if QueryBlob is JSON
	var t LQLQuery

	if err := json.Unmarshal([]byte(q.QueryBlob), &t); err == nil {
		q.StartTimeRange = t.StartTimeRange
		q.EndTimeRange = t.EndTimeRange
		q.QueryText = t.QueryText
		return
	}

	// if QueryBlob is LQL
	if matched, _ := regexp.MatchString(reLQL, q.QueryBlob); matched {
		q.QueryText = q.QueryBlob
	}
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
	lqlQuery.translate()

	err = svc.client.RequestEncoderDecoder(
		"POST",
		apiLQL,
		lqlQuery,
		&response,
	)
	return
}

func (svc *LQLService) GetQueries() (
	response LQLQueryResponse,
	err error,
) {
	return svc.GetQueryByID("")
}

func (svc *LQLService) GetQueryByID(queryID string) (
	response LQLQueryResponse,
	err error,
) {
	var uri string

	if queryID != "" {
		uri = apiLQL + "?LQL_ID=" + url.QueryEscape(queryID)
	} else {
		uri = apiLQL
	}

	err = svc.client.RequestDecoder(
		"GET",
		uri,
		nil,
		&response,
	)
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
	lqlQuery.translate()

	err = svc.client.RequestEncoderDecoder(
		"POST",
		apiLQLQuery,
		lqlQuery,
		&response,
	)
	return
}
