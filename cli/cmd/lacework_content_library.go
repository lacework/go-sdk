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

package cmd

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwcomponent"
)

const (
	lclComponentName string = "lacework-content-library"
	lclIndexPath     string = "content.index"
)

type LCLReference struct {
	ID   string `json:"id"`
	Type string `json:"content_type"`
	Path string `json:"path"`
	URI  string `json:"uri"`
}

type LCLQuery struct {
	References []LCLReference `json:"references"`
}

type LCLPolicy struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Tags        []string       `json:"tags"`
	References  []LCLReference `json:"references"`
}

type LaceworkContentLibrary struct {
	Component *lwcomponent.Component
	Queries   map[string]LCLQuery  `json:"queries"`
	Policies  map[string]LCLPolicy `json:"policies"`
}

func IsLCLInstalled(state lwcomponent.State) bool {
	component := state.GetComponent(lclComponentName)

	if component == nil || component.Status() != lwcomponent.Installed {
		return false
	}
	return true
}

func LoadLCL(state lwcomponent.State) (*LaceworkContentLibrary, error) {
	lcl := new(LaceworkContentLibrary)
	lcl.Component = state.GetComponent(lclComponentName)

	index, err := lcl.run(lclIndexPath)
	if err != nil {
		return new(LaceworkContentLibrary), errors.Wrap(
			err, "unable to load Lacework Content Library")
	}

	if err := json.Unmarshal([]byte(index), lcl); err != nil {
		return new(LaceworkContentLibrary), errors.Wrap(
			err, "unable to load Lacework Content Library")
	}
	return lcl, nil
}

func (lcl LaceworkContentLibrary) run(path string) (string, error) {
	if lcl.Component == nil || lcl.Component.Status() != lwcomponent.Installed {
		return "", errors.New("Lacework Content Library is not installed")
	}
	stdout, _, err := lcl.Component.RunAndReturn([]string{path}, nil)
	return stdout, err
}

func (lcl LaceworkContentLibrary) getReferenceForQuery(id string) (LCLReference, error) {
	var ref LCLReference

	if id == "" {
		return ref, errors.New("query ID must be provided")
	}
	if _, ok := lcl.Queries[id]; !ok {
		return ref, errors.New("query does not exist in library")
	}
	if len(lcl.Queries[id].References) < 1 {
		return ref, errors.New("query exists but is malformed")
	}
	ref = lcl.Queries[id].References[0]
	return ref, nil
}

func (lcl LaceworkContentLibrary) ListQueries() api.QueriesResponse {
	var queries []api.Query

	for id := range lcl.Queries {
		queries = append(queries, api.Query{QueryID: id})
	}
	return api.QueriesResponse{Data: queries}
}

func (lcl LaceworkContentLibrary) GetQuery(id string) (api.QueryResponse, error) {
	var response api.QueryResponse

	// get query reference
	ref, err := lcl.getReferenceForQuery(id)
	if err != nil {
		return response, err
	}
	// check query path
	if ref.Path == "" {
		return response, errors.New("query exists but is malformed")
	}
	// get query string
	queryString, err := lcl.run(ref.Path)
	if err != nil {
		return response, err
	}
	// parse query string
	newQuery, err := api.ParseNewQuery(queryString)
	if err != nil {
		return response, queryErrorCrumbs(queryString)
	}
	response.Data = api.Query{
		QueryID:     newQuery.QueryID,
		QueryText:   newQuery.QueryText,
		EvaluatorID: newQuery.EvaluatorID,
	}
	return response, nil
}
