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

type DatasourcesResponse struct {
	Data    []Datasource `json:"data"`
	Message string       `json:"message"`
}

type DatasourceResponse struct {
	Data    Datasource `json:"data"`
	Message string     `json:"message"`
}

type Datasource struct {
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	ResultSchema []DatasourceSchema `json:"resultSchema"`
}

type DatasourceSchema struct {
	Name        string `json:"name"`
	DataType    string `json:"dataType"`
	Description string `json:"description"`
}

// DatasourcesService is a service that interacts with the Datasources
// endpoints from the Lacework Server
type DatasourcesService struct {
	client *Client
}

func (svc *DatasourcesService) List() (
	response DatasourcesResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET", apiV2Datasources, nil, &response)
	return
}

func (svc *DatasourcesService) Get(id string) (
	response DatasourceResponse,
	err error,
) {
	if id == "" {
		err = errors.New("datasource ID must be provided")
		return
	}
	err = svc.client.RequestDecoder(
		"GET",
		fmt.Sprintf("%s/%s", apiV2Datasources, url.QueryEscape(id)),
		nil,
		&response,
	)
	return
}
