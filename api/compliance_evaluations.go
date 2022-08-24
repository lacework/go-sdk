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

type ComplianceEvaluationService struct {
	client *Client
}

type complianceEvaluationDataset string

const AwsComplianceEvaluationDataset complianceEvaluationDataset = "AwsCompliance"

// Search expects the response and the search filters
//
// e.g.
//
//  var (
//	  awsComplianceEvaluationSearchResponse api.ComplianceEvaluationAwsResponse
//	  filter = api.ComplianceEvaluationSearch{
//		  SearchFilter: api.SearchFilter{
//			  Filters: []api.Filter{{
//				  Expression: "eq",
//				  Field:      "resource",
//				  Value:      arn:aws:s3:::my-bucket,
//			  }},
//		  },
//		  Dataset: api.AwsComplianceEvaluationDataset,
//	  }
//  )
//   lacework.V2.ComplianceEvaluation.Search(&awsComplianceEvaluationSearchResponse, filters)
//
func (svc *ComplianceEvaluationService) Search(response interface{}, filters ComplianceEvaluationSearch) error {
	return svc.client.RequestEncoderDecoder("POST", apiV2ComplianceEvaluationsSearch, filters, response)
}

type ComplianceEvaluationSearch struct {
	SearchFilter
	Dataset complianceEvaluationDataset `json:"dataset"`
}
