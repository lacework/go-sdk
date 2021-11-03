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

	"github.com/pkg/errors"
)

type NewQuery struct {
	QueryID     string `json:"queryId" yaml:"queryId"`
	QueryText   string `json:"queryText" yaml:"queryText"`
	EvaluatorID string `json:"evaluatorId,omitempty" yaml:"evaluatorId"`
}

type UpdateQuery struct {
	QueryText string `json:"queryText"`
}

type Query struct {
	QueryID        string                   `json:"queryId"`
	QueryText      string                   `json:"queryText"`
	EvaluatorID    string                   `json:"evaluatorId"`
	Owner          string                   `json:"owner"`
	LastUpdateTime string                   `json:"lastUpdateTime"`
	LastUpdateUser string                   `json:"lastUpdateUser"`
	ResultSchema   []map[string]interface{} `json:"resultSchema"`
}

type QueryResponse struct {
	Data    Query  `json:"data"`
	Message string `json:"message"`
}

type QueriesResponse struct {
	Data    []Query `json:"data"`
	Message string  `json:"message"`
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

func (svc *QueryService) Update(id string, uq UpdateQuery) (
	response QueryResponse,
	err error,
) {
	if id == "" {
		err = errors.New("query ID must be provided")
		return
	}
	err = svc.client.RequestEncoderDecoder(
		"PATCH",
		fmt.Sprintf("%s/%s", apiV2Queries, url.QueryEscape(id)),
		uq,
		&response,
	)
	return
}

func (svc *QueryService) List() (
	response QueriesResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET", apiV2Queries, nil, &response)
	return
}

func (svc *QueryService) Get(id string) (
	response QueryResponse,
	err error,
) {
	if id == "" {
		err = errors.New("query ID must be provided")
		return
	}
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf("%s/%s", apiV2Queries, url.QueryEscape(id)),
		nil,
		&response,
	)
	return
}
