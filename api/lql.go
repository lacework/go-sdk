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
	"regexp"

	"github.com/pkg/errors"
)

const (
	reLQL string = `(?ms)^(\w+)\(\w+\s\w+\)\s*{`
)

type LQLQuery struct {
	StartTimeRange string `json:"START_TIME_RANGE,omitempty"`
	EndTimeRange   string `json:"END_TIME_RANGE,omitempty"`
	QueryText      string `json:"QUERY_TEXT"`
	// QueryBlob is a special string that supports type conversion
	// back and forth from LQL to JSON
	QueryBlob string `json:"-"`
}

func (q *LQLQuery) translate() {
	// if QueryText is not already populated
	if q.QueryText == "" {
		// if QueryBlob is JSON
		var t LQLQuery
		err := json.Unmarshal([]byte(q.QueryBlob), &t)

		if err == nil {
			q.QueryText = t.QueryText
			// if QueryBlob is LQL
		} else if matched, _ := regexp.MatchString(reLQL, q.QueryBlob); matched {
			q.QueryText = q.QueryBlob
		}
	}
}

// LQLService is a service that interacts with the LQL
// endpoints from the Lacework Server
type LQLService struct {
	client *Client
}

func (svc *LQLService) CreateQuery(query string) (
	response map[string]interface{},
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

func (svc *LQLService) DeleteQuery(queryID string) (
	response map[string]interface{},
	err error,
) {
	var uri string

	if queryID != "" {
		uri = apiLQL + "?LQL_ID=" + url.QueryEscape(queryID)
	} else {
		uri = apiLQL
	}

	err = svc.client.RequestDecoder(
		"DELETE",
		uri,
		nil,
		&response,
	)
	return
}

func (svc *LQLService) GetQueries() (
	response map[string]interface{},
	err error,
) {
	return svc.GetQueryByID("")
}

func (svc *LQLService) GetQueryTextByID(queryID string) (
	query string,
	err error,
) {
	msg := "unable to retrieve query"
	var response map[string]interface{}

	response, err = svc.GetQueryByID(queryID)

	if err != nil {
		err = errors.Wrap(err, msg)
	} else if data, ok := response["data"]; ok {
		queries, ok := data.([]interface{})
		if ok && len(queries) != 0 {
			if in, ok := queries[0].(map[string]interface{}); ok {
				for k, v := range in {
					if k == "QUERY_TEXT" || k == "query_text" {
						query = fmt.Sprintf("%v", v)
						err = nil
						return
					}
				}
			}
		}
		err = errors.New(msg)
	}
	return
}

func (svc *LQLService) GetQueryByID(queryID string) (
	response map[string]interface{},
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

func (svc *LQLService) UpdateQuery(query string) (
	response map[string]interface{},
	err error,
) {
	lqlQuery := LQLQuery{QueryBlob: query}
	lqlQuery.translate()

	err = svc.client.RequestEncoderDecoder(
		"PATCH",
		apiLQL,
		lqlQuery,
		&response,
	)
	return
}
