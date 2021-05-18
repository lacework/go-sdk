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
)

type LQLDescribeResponse struct {
	Data    []LQLDescribeData `json:"data"`
	Ok      bool              `json:"ok"`
	Message string            `json:"message"`
}

type LQLDescribeData struct {
	Complexity  int                     `json:"complexity"`
	MaxDuration int                     `json:"maxDuration"`
	Parameters  []LQLDescribeParameters `json:"parameters"`
	PrimaryKey  []interface{}           `json:"primaryKey"`
	Props       map[string]interface{}  `json:"props"`
	Schema      []LQLDescribeSchema     `json:"schema"`
	Type        string                  `json:"type"`
}

type LQLDescribeParameters struct {
	Default  string                 `json:"default"`
	Name     string                 `json:"name"`
	Props    map[string]interface{} `json:"props"`
	Required bool                   `json:"required"`
	Type     string                 `json:"type"`
}

type LQLDescribeSchema struct {
	Name  string                 `json:"name"`
	Props map[string]interface{} `json:"props"`
	Type  string                 `json:"type"`
}

func (svc *LQLService) Describe(dataSource string) (
	response LQLDescribeResponse,
	err error,
) {
	uri := fmt.Sprintf("%s/%s", apiLQLDescribe, url.QueryEscape(dataSource))

	err = svc.client.RequestDecoder("GET", uri, nil, &response)
	return
}
