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

type DescribeResponse struct {
	Data    []DescribeData `json:"data"`
	Ok      bool           `json:"ok"`
	Message string         `json:"message"`
}

type DescribeData struct {
	Complexity  int                    `json:complexity`
	MaxDuration int                    `json:maxDuration`
	Parameters  []DescribeParameters   `json:parameters`
	PrimaryKey  []interface{}          `json:primaryKey`
	Props       map[string]interface{} `json:props`
	Schema      []DescribeSchema       `json:schema`
}

type DescribeParameters struct {
	Default  string                 `json:default`
	Name     string                 `json:name`
	Props    map[string]interface{} `json:props`
	Required bool                   `json:required`
	Type     string                 `json:type`
}

type DescribeSchema struct {
	Name  string                 `json:name`
	Props map[string]interface{} `json:props`
	Type  string                 `json:type`
}

func (svc *LQLService) Describe(dataSource string) (
	response DescribeResponse,
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
