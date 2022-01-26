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

func getPolicyReference(refs []LCLReference) (LCLReference, error) {
	for i := range refs {
		if refs[i].Type == "policy" {
			return refs[i], nil
		}
	}
	return LCLReference{}, errors.New("policy reference not found")
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
	return lcl.Queries[id].References[0], nil
}

func (lcl LaceworkContentLibrary) getReferencesForPolicy(id string) ([]LCLReference, error) {
	var refs []LCLReference

	if id == "" {
		return refs, errors.New("policy ID must be provided")
	}
	if _, ok := lcl.Policies[id]; !ok {
		return refs, errors.New("policy does not exist in library")
	}
	if len(lcl.Policies[id].References) < 2 {
		return refs, errors.New("policy exists but is malformed")
	}
	return lcl.Policies[id].References, nil
}

func (lcl LaceworkContentLibrary) ListQueries() api.QueriesResponse {
	var queries []api.Query

	for id := range lcl.Queries {
		queries = append(queries, api.Query{QueryID: id})
	}
	return api.QueriesResponse{Data: queries}
}

func (lcl LaceworkContentLibrary) GetQuery(id string) (string, error) {
	// get query reference
	ref, err := lcl.getReferenceForQuery(id)
	if err != nil {
		return "", err
	}
	// check query path
	if ref.Path == "" {
		return "", errors.New("query exists but is malformed")
	}
	// get query string
	return lcl.run(ref.Path)
}

func (lcl LaceworkContentLibrary) ListPolicies() (api.PoliciesResponse, error) {
	var policies []api.Policy

	for policyID := range lcl.Policies {
		var queryRef LCLReference

		for i := range lcl.Policies[policyID].References {
			if lcl.Policies[policyID].References[i].Type == "query" {
				queryRef = lcl.Policies[policyID].References[i]
				break
			}
		}
		if queryRef.ID == "" {
			return api.PoliciesResponse{Data: policies}, errors.New(
				"unable to identify query for one or more policies")
		}
		policies = append(policies, api.Policy{
			PolicyID:    policyID,
			Title:       lcl.Policies[policyID].Title,
			Description: lcl.Policies[policyID].Description,
			QueryID:     queryRef.ID,
		})
	}
	return api.PoliciesResponse{Data: policies}, nil
}

func (lcl LaceworkContentLibrary) GetPolicy(id string) (string, error) {
	// get policy references
	refs, err := lcl.getReferencesForPolicy(id)
	if err != nil {
		return "", err
	}
	ref, err := getPolicyReference(refs)
	if err != nil || ref.Path == "" {
		return "", errors.New("policy exists but is malformed")
	}
	// get policy string
	return lcl.run(ref.Path)
}
