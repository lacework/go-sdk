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
	Queries  map[string]LCLQuery  `json:"queries"`
	Policies map[string]LCLPolicy `json:"policies"`
}

func IsLCLInstalled(state lwcomponent.State) bool {
	component := state.GetComponent(lclComponentName)

	if component == nil || component.Status() != lwcomponent.Installed {
		return false
	}
	return true
}

func LoadLCL(state lwcomponent.State) (*LaceworkContentLibrary, error) {
	index := new(LaceworkContentLibrary)
	component := state.GetComponent(lclComponentName)

	if component == nil || component.Status() != lwcomponent.Installed {
		return index, errors.New("Lacework Content Library is not available")
	}

	stdout, _, err := component.RunAndReturn([]string{lclIndexPath}, nil)
	if err != nil || stdout == "" {
		return index, errors.Wrap(err, "unable to retrieve index from Lacework Content Library")
	}

	if err := json.Unmarshal([]byte(stdout), index); err != nil {
		return index, errors.Wrap(err, "unable to parse Lacework Content Library index")
	}
	return index, nil
}

func (lcl LaceworkContentLibrary) ListQueries() api.QueriesResponse {
	var queries []api.Query

	for id := range lcl.Queries {
		queries = append(queries, api.Query{QueryID: id})
	}
	return api.QueriesResponse{Data: queries}
}
