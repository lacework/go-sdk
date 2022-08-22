//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

import "time"

type ComplianceEvaluationAwsResponse struct {
	Data   []ComplianceEvaluationAws `json:"data"`
	Paging V2Pagination              `json:"paging"`
}

func (r ComplianceEvaluationAwsResponse) PageInfo() *V2Pagination {
	return &r.Paging
}
func (r *ComplianceEvaluationAwsResponse) ResetPaging() {
	r.Paging = V2Pagination{}
}

type ComplianceEvaluationAws struct {
	Account struct {
		AccountId    string `json:"AccountId"`
		AccountAlias string `json:"Account_Alias"`
	} `json:"account"`
	EvalType       string    `json:"evalType"`
	Id             string    `json:"id"`
	Reason         string    `json:"reason"`
	Recommendation string    `json:"recommendation"`
	ReportTime     time.Time `json:"reportTime"`
	Resource       string    `json:"resource"`
	Section        string    `json:"section"`
	Severity       string    `json:"severity"`
	Status         string    `json:"status"`
}
