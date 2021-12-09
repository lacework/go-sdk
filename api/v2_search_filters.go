//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

// SearchFilter is the representation of an advanced search payload
// for retrieving information out of the Lacework APIv2 Server
//
// An advanced example of a SearchFilter to search for an Agent
// Access Token that matches the provider token alias and return
// only the token found:
//
//		SearchFilter{
//			Filters: []Filter{
//				Filter{
//					Field:      "tokenAlias",
//					Expression: "eq",
//					Value:      "k8s-deployment,
//				},
//			},
//			Returns: []string{"accessToken"},
//		}
type SearchFilter struct {
	TimeFilter `json:"timeFilter,omitempty"`
	Filters    []Filter `json:"filters"`
	Returns    []string `json:"returns"`
}

type Filter struct {
	Expression string   `json:"expression"`
	Field      string   `json:"field"`
	Value      string   `json:"value"`
	Values     []string `json:"values"`
}

type TimeFilter struct {
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
}
