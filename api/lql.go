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

// LQLQuery is a special string that supports type conversion
// back and forth from LQL to JSON
type LQLQuery string

func (q LQLQuery) json() (
	qMap map[string]interface{},
	err error,
) {
	matched, _ := regexp.MatchString(reLQL, string(q))
	if matched {
		qMap = make(map[string]interface{})
		qMap["QUERY_TEXT"] = string(q)
		return
	}
	qBytes := []byte(q)
	err = json.Unmarshal(qBytes, &qMap)
	return
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
	var qMap map[string]interface{}
	qMap, err = LQLQuery(query).json()
	if err != nil {
		return
	}
	err = svc.client.RequestEncoderDecoder(
		"POST",
		apiLQL,
		qMap,
		&response,
	)
	return
}

func (svc *LQLService) CompileQuery(query string) (
	response map[string]interface{},
	err error,
) {
	var qMap map[string]interface{}
	qMap, err = LQLQuery(query).json()
	if err != nil {
		return
	}
	err = svc.client.RequestEncoderDecoder(
		"POST",
		apiLQLCompile,
		qMap,
		&response,
	)
	return
}

func (svc *LQLService) DataSources() (
	response map[string]interface{},
	err error,
) {
	err = svc.client.RequestDecoder(
		"GET",
		apiLQLDataSources,
		nil,
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

func (svc *LQLService) Describe(dataSource string) (
	response map[string]interface{},
	err error,
) {
	uri := "%v/%v"
	uri = fmt.Sprintf(uri, apiLQLDescribe, url.QueryEscape(dataSource))

	err = svc.client.RequestDecoder(
		"GET",
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
	var qMap map[string]interface{}
	qMap, err = LQLQuery(query).json()
	if err != nil {
		return
	}
	if start != "" {
		qMap["START_TIME_RANGE"] = start
	}
	if end != "" {
		qMap["END_TIME_RANGE"] = end
	}
	err = svc.client.RequestEncoderDecoder(
		"POST",
		apiLQLQuery,
		qMap,
		&response,
	)
	return
}

func (svc *LQLService) UpdateQuery(query string) (
	response map[string]interface{},
	err error,
) {
	var qMap map[string]interface{}
	qMap, err = LQLQuery(query).json()
	if err != nil {
		return
	}
	err = svc.client.RequestEncoderDecoder(
		"PATCH",
		apiLQL,
		qMap,
		&response,
	)
	return
}
