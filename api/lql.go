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
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type NewQuery struct {
	QueryID     string `json:"queryId" yaml:"queryId"`
	QueryText   string `json:"queryText" yaml:"queryText"`
	EvaluatorID string `json:"evaluatorId" yaml:"evaluatorId"`
}

type UpdateQuery struct {
	QueryText string `json:"queryText" yaml:"queryText"`
}

type Query struct {
	QueryID        string                   `json:"queryId" yaml:"queryId"`
	QueryText      string                   `json:"queryText" yaml:"queryText"`
	EvaluatorID    string                   `json:"evaluatorId" yaml:"evaluatorId"`
	Owner          string                   `json:"owner" yaml:"owner"`
	LastUpdateTime string                   `json:"lastUpdateTime" yaml:"lastUpdateTime"`
	LastUpdateUser string                   `json:"lastUpdateUser" yaml:"lastUpdateUser"`
	ResultSchema   []map[string]interface{} `json:"resultSchema" yaml:"resultSchema"`
}

type QueryResponse struct {
	Data    Query  `json:"data"`
	Message string `json:"message"`
}

type QueriesResponse struct {
	Data    []Query `json:"data"`
	Message string  `json:"message"`
}

func validateQueryRange(start, end time.Time) (err error) {
	// validate range
	if start.After(end) {
		err = errors.New("date range should have a start time before the end time")
		return
	}
	return nil
}

// QueryService is a service that interacts with the Queries
// endpoints from the Lacework Server
type QueryService struct {
	client *Client
}

func (svc *QueryService) Create(nq NewQuery) (
	response QueryResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2Queries, nq, &response)
	return
}

func (svc *QueryService) Update(queryID string, uq UpdateQuery) (
	response QueryResponse,
	err error,
) {
	if queryID == "" {
		err = errors.New("query ID must be provided")
		return
	}
	err = svc.client.RequestEncoderDecoder(
		"PATCH",
		fmt.Sprintf("%s/%s", apiV2Queries, url.QueryEscape(queryID)),
		uq,
		&response,
	)
	return
}

func (svc *QueryService) GetQueries() (
	response QueriesResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET", apiV2Queries, nil, &response)
	return
}

func (svc *QueryService) GetByID(queryID string) (
	response QueryResponse,
	err error,
) {
	if queryID == "" {
		err = errors.New("query ID must be provided")
		return
	}
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf("%s/%s", apiV2Queries, url.QueryEscape(queryID)),
		nil,
		&response,
	)
	return
}

func (svc *QueryService) Execute(queryText string, start time.Time, end time.Time) (
	response map[string]interface{},
	err error,
) {
	if queryText == "" {
		err = errors.New("query text must be provided")
		return
	}
	if err = validateQueryRange(start, end); err != nil {
		return
	}
	query := map[string]string{
		"query_text":       queryText,
		"start_time_range": start.UTC().Format(time.RFC3339),
		"end_time_range":   end.UTC().Format(time.RFC3339),
	}
	err = svc.client.RequestEncoderDecoder("POST", apiLQLQuery, query, &response)
	return
}
