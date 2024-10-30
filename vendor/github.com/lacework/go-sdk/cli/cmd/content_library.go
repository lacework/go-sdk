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

	"github.com/lacework/go-sdk/lwcomponent"
)

const (
	lclComponentName string = "content-library"
	lclIndexPath     string = "content.index"
)

type LCLContentType string

const (
	LCLQueryType  LCLContentType = "query"
	LCLPolicyType LCLContentType = "policy"
)

type LCLReference struct {
	ID   string         `json:"id"`
	Type LCLContentType `json:"content_type"`
	Path string         `json:"path"`
	URI  string         `json:"uri"`
}

func getPolicyReference(refs []LCLReference) (LCLReference, error) {
	for i := range refs {
		if refs[i].Type == LCLPolicyType {
			return refs[i], nil
		}
	}
	return LCLReference{}, errors.New("policy reference not found")
}

type LCLQuery struct {
	References []LCLReference `json:"references"`
}

type LCLPolicy struct {
	PolicyID    string         `json:"policyId"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Tags        []string       `json:"tags"`
	QueryID     string         `json:"queryId"`
	References  []LCLReference `json:"references"`
}

type LaceworkContentLibrary struct {
	Component  *lwcomponent.Component
	Queries    map[string]LCLQuery  `json:"queries"`
	Policies   map[string]LCLPolicy `json:"policies"`
	PolicyTags map[string][]string  `json:"policy_tags"`
}

func (c *cliState) isLCLInstalled() bool {
	return c.IsComponentInstalled(lclComponentName)
}

func (c *cliState) LoadLCL() (*LaceworkContentLibrary, error) {
	var (
		baseErr = "unable to load Lacework Content Library"
		lcl     = new(LaceworkContentLibrary)
		found   bool
	)

	if c.LwComponents == nil {
		return lcl, errors.New(baseErr)
	}

	lcl.Component, found = c.LwComponents.GetComponent(lclComponentName)
	if !found {
		return lcl, errors.Wrap(errors.New("component not installed"), baseErr)
	}

	index, err := lcl.run(lclIndexPath)
	if err != nil {
		return new(LaceworkContentLibrary), errors.Wrap(err, baseErr)
	}

	if err := json.Unmarshal([]byte(index), lcl); err != nil {
		return new(LaceworkContentLibrary), errors.Wrap(err, baseErr)
	}

	for policyID, policy := range lcl.Policies {
		for i := range policy.References {
			if policy.References[i].Type == LCLQueryType {
				policy.QueryID = policy.References[i].ID
			}
			if policy.References[i].Type == LCLPolicyType {
				policy.PolicyID = policy.References[i].ID
			}
		}
		lcl.Policies[policyID] = policy
	}

	return lcl, nil
}

func (lcl *LaceworkContentLibrary) run(path string) (string, error) {
	if lcl.Component == nil || !lcl.Component.IsInstalled() {
		return "", errors.New("Lacework Content Library is not installed")
	}
	stdout, _, err := lcl.Component.RunAndReturn([]string{path}, nil)
	return stdout, err
}

func (lcl *LaceworkContentLibrary) getReferenceForQuery(id string) (LCLReference, error) {
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

func (lcl *LaceworkContentLibrary) getReferencesForPolicy(id string) ([]LCLReference, error) {
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

func (lcl *LaceworkContentLibrary) GetQuery(id string) (string, error) {
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

func (lcl *LaceworkContentLibrary) GetPolicy(id string) (string, error) {
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

func (lcl *LaceworkContentLibrary) GetPoliciesByTag(t string) map[string]LCLPolicy {
	var (
		policies  map[string]LCLPolicy = map[string]LCLPolicy{}
		policyIDs []string
		ok        bool
	)

	if policyIDs, ok = lcl.PolicyTags[t]; !ok {
		return policies
	}

	for _, policyID := range policyIDs {
		if lclPolicy, ok := lcl.Policies[policyID]; ok {
			policies[policyID] = lclPolicy
		}
	}

	return policies
}
