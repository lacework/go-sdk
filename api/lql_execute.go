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

	"github.com/lacework/go-sdk/lwtime"
	"github.com/pkg/errors"
)

type ExecuteQuery struct {
	QueryText   string `json:"queryText"`
	EvaluatorID string `json:"evaluatorId,omitempty"`
}

type ExecuteQueryArgument struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ExecuteQueryRequest struct {
	Query     ExecuteQuery           `json:"query"`
	Arguments []ExecuteQueryArgument `json:"arguments"`
}

type ExecuteQueryByIDRequest struct {
	QueryID   string                 `json:"queryId,omitempty"`
	Arguments []ExecuteQueryArgument `json:"arguments"`
}

func validateQueryArguments(args []ExecuteQueryArgument) (err error) {
	var hasStart, hasEnd bool
	var start, end time.Time

	for _, arg := range args {
		if arg.Name == "StartTimeRange" {
			hasStart = true
			start, err = validateQueryTimeString(arg.Value)
		}
		if err != nil {
			return errors.Wrap(err, "invalid StartTimeRange argument")
		}

		if arg.Name == "EndTimeRange" {
			hasEnd = true
			end, err = validateQueryTimeString(arg.Value)
		}
		if err != nil {
			return errors.Wrap(err, "invalid EndTimeRange argument")
		}
	}

	if hasStart && hasEnd {
		return validateQueryRange(start, end)
	}
	return nil
}

// StartTimeRange and EndTimeRange should be
func validateQueryTimeString(s string) (time.Time, error) {
	return time.Parse(lwtime.RFC3339Milli, s)
}

func validateQueryRange(start, end time.Time) (err error) {
	// validate range
	if start.After(end) {
		err = errors.New("date range should have a start time before the end time")
		return
	}
	return nil
}

func (svc *QueryService) Execute(request ExecuteQueryRequest) (
	response map[string]interface{},
	err error,
) {
	if err = validateQueryArguments(request.Arguments); err != nil {
		return
	}
	err = svc.client.RequestEncoderDecoder("POST", apiV2QueriesExecute, request, &response)
	return
}

func (svc *QueryService) ExecuteByID(request ExecuteQueryByIDRequest) (
	response map[string]interface{},
	err error,
) {
	if request.QueryID == "" {
		err = errors.New("query ID must be provided")
		return
	}
	queryID := request.QueryID
	request.QueryID = "" // omit for POST

	if err = validateQueryArguments(request.Arguments); err != nil {
		return
	}
	err = svc.client.RequestEncoderDecoder(
		"POST",
		fmt.Sprintf("%s/%s/execute", apiV2Queries, url.QueryEscape(queryID)),
		request,
		&response,
	)
	return
}
