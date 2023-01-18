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

package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/lacework/go-sdk/internal/lacework"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestRecommendationsAwsList(t *testing.T) {
	var (
		expectedLen = 161
		fakeServer  = lacework.MockServer()
	)

	fakeServer.MockToken("TOKEN")
	fakeServer.UseApiV2()
	fakeServer.MockAPI("recommendations/aws",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "List() should be a GET method")
			Recommendations := generateRecommendations()
			fmt.Fprintf(w, Recommendations)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Recommendations.Aws.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response))
	for _, rec := range response {
		assert.NotEmpty(t, rec.ID)
	}
}

func TestRecommendationsAwsPatch(t *testing.T) {
	var (
		fakeServer = lacework.MockServer()
	)

	fakeServer.MockToken("TOKEN")
	fakeServer.UseApiV2()
	fakeServer.MockAPI("recommendations/aws",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PATCH", r.Method, "AwsPatch() should be a PATCH method")
			PatchResponse := generateRecommendationsPatchResponse()
			fmt.Fprintf(w, PatchResponse)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	recommendationPatch := api.RecommendationStateV2{"LW_S3_1": "disable"}

	response, err := c.V2.Recommendations.Aws.Patch(recommendationPatch)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	for k, v := range response.Data[0] {
		assert.Equal(t, k, "LW_S3_1")
		assert.False(t, v.Enabled)
	}
	recList := response.RecommendationList()
	assert.Equal(t, 1, len(recList))
	assert.Equal(t, "LW_S3_1", recList[0].ID)
	assert.False(t, recList[0].State)
}

func generateRecommendations() string {
	return listRecommendations()
}

func generateRecommendationsPatchResponse() string {
	return `{
     "data": [
         {
             "LW_S3_1": {
                 "enabled": false
             }
         }
     ],
     "ok": true,
     "message": "SUCCESS"
 }`
}

func listRecommendations() string {
	return `{
     "data": [
         {
             "LW_AWS_NETWORKING_15": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_14": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_17": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_16": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_19": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_18": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_11": {
                 "enabled": true
             },
             "LW_S3_18": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_10": {
                 "enabled": true
             },
             "LW_S3_19": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_13": {
                 "enabled": true
             },
             "LW_S3_16": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_12": {
                 "enabled": true
             },
             "LW_S3_17": {
                 "enabled": false
             },
             "AWS_CIS_1_22": {
                 "enabled": true
             },
             "LW_S3_14": {
                 "enabled": true
             },
             "AWS_CIS_1_23": {
                 "enabled": true
             },
             "LW_S3_15": {
                 "enabled": true
             },
             "AWS_CIS_1_20": {
                 "enabled": true
             },
             "LW_S3_12": {
                 "enabled": false
             },
             "AWS_CIS_1_21": {
                 "enabled": true
             },
             "LW_S3_13": {
                 "enabled": false
             },
             "LW_S3_10": {
                 "enabled": true
             },
             "LW_S3_11": {
                 "enabled": true
             },
             "AWS_CIS_1_24": {
                 "enabled": true
             },
             "LW_S3_1": {
                 "enabled": true
             },
             "LW_S3_2": {
                 "enabled": true
             },
             "LW_S3_3": {
                 "enabled": true
             },
             "LW_S3_4": {
                 "enabled": true
             },
             "LW_S3_5": {
                 "enabled": true
             },
             "LW_S3_6": {
                 "enabled": false
             },
             "LW_S3_7": {
                 "enabled": true
             },
             "LW_S3_8": {
                 "enabled": true
             },
             "LW_S3_9": {
                 "enabled": false
             },
             "AWS_CIS_2_8": {
                 "enabled": true
             },
             "AWS_CIS_2_7": {
                 "enabled": true
             },
             "AWS_CIS_2_6": {
                 "enabled": true
             },
             "AWS_CIS_2_5": {
                 "enabled": true
             },
             "AWS_CIS_2_4": {
                 "enabled": false
             },
             "AWS_CIS_2_3": {
                 "enabled": true
             },
             "AWS_CIS_2_2": {
                 "enabled": true
             },
             "AWS_CIS_2_1": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_1": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_2": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_3": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_4": {
                 "enabled": true
             },
             "AWS_CIS_2_9": {
                 "enabled": true
             },
             "LW_S3_21": {
                 "enabled": false
             },
             "LW_S3_20": {
                 "enabled": true
             },
             "LW_AWS_RDS_1": {
                 "enabled": false
             },
             "LW_AWS_NETWORKING_37": {
                 "enabled": true
             },
             "LW_AWS_SERVERLESS_5": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_36": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_39": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_38": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_31": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_30": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_33": {
                 "enabled": true
             },
             "LW_AWS_SERVERLESS_1": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_32": {
                 "enabled": true
             },
             "LW_AWS_SERVERLESS_2": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_35": {
                 "enabled": true
             },
             "LW_AWS_SERVERLESS_3": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_34": {
                 "enabled": true
             },
             "LW_AWS_SERVERLESS_4": {
                 "enabled": false
             },
             "AWS_CIS_3_7": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_26": {
                 "enabled": true
             },
             "AWS_CIS_3_6": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_25": {
                 "enabled": true
             },
             "AWS_CIS_3_5": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_28": {
                 "enabled": true
             },
             "AWS_CIS_3_4": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_27": {
                 "enabled": true
             },
             "AWS_CIS_3_3": {
                 "enabled": true
             },
             "AWS_CIS_3_2": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_29": {
                 "enabled": true
             },
             "AWS_CIS_3_1": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_20": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_22": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_21": {
                 "enabled": true
             },
             "AWS_CIS_3_9": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_24": {
                 "enabled": true
             },
             "AWS_CIS_3_8": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_23": {
                 "enabled": true
             },
             "LW_AWS_IAM_13": {
                 "enabled": true
             },
             "LW_AWS_IAM_14": {
                 "enabled": false
             },
             "LW_AWS_IAM_11": {
                 "enabled": true
             },
             "LW_AWS_IAM_12": {
                 "enabled": false
             },
             "LW_AWS_NETWORKING_51": {
                 "enabled": false
             },
             "LW_AWS_NETWORKING_50": {
                 "enabled": false
             },
             "AWS_CIS_3_10": {
                 "enabled": true
             },
             "LW_AWS_IAM_2": {
                 "enabled": true
             },
             "LW_AWS_IAM_3": {
                 "enabled": true
             },
             "LW_AWS_IAM_1": {
                 "enabled": true
             },
             "LW_AWS_IAM_6": {
                 "enabled": true
             },
             "LW_AWS_IAM_7": {
                 "enabled": true
             },
             "LW_AWS_IAM_4": {
                 "enabled": true
             },
             "LW_AWS_IAM_5": {
                 "enabled": true
             },
             "LW_AWS_IAM_8": {
                 "enabled": true
             },
             "LW_AWS_IAM_9": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_48": {
                 "enabled": true
             },
             "AWS_CIS_4_5": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_47": {
                 "enabled": true
             },
             "AWS_CIS_4_4": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_49": {
                 "enabled": false
             },
             "AWS_CIS_4_2": {
                 "enabled": true
             },
             "AWS_CIS_4_1": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_40": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_42": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_41": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_44": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_43": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_46": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_45": {
                 "enabled": true
             },
             "AWS_CIS_3_14": {
                 "enabled": true
             },
             "AWS_CIS_3_13": {
                 "enabled": true
             },
             "AWS_CIS_3_12": {
                 "enabled": true
             },
             "AWS_CIS_3_11": {
                 "enabled": true
             },
             "AWS_CIS_3_15": {
                 "enabled": true
             },
             "LW_AWS_GENERAL_SECURITY_4": {
                 "enabled": true
             },
             "LW_AWS_GENERAL_SECURITY_5": {
                 "enabled": true
             },
             "LW_AWS_GENERAL_SECURITY_2": {
                 "enabled": false
             },
             "LW_AWS_GENERAL_SECURITY_3": {
                 "enabled": true
             },
             "LW_AWS_GENERAL_SECURITY_8": {
                 "enabled": true
             },
             "LW_AWS_GENERAL_SECURITY_6": {
                 "enabled": true
             },
             "LW_AWS_GENERAL_SECURITY_7": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_9": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_5": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_6": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_7": {
                 "enabled": true
             },
             "LW_AWS_NETWORKING_8": {
                 "enabled": true
             },
             "LW_AWS_ELASTICSEARCH_3": {
                 "enabled": true
             },
             "LW_AWS_ELASTICSEARCH_2": {
                 "enabled": true
             },
             "LW_AWS_ELASTICSEARCH_1": {
                 "enabled": true
             },
             "LW_AWS_ELASTICSEARCH_4": {
                 "enabled": true
             },
             "AWS_CIS_1_9": {
                 "enabled": false
             },
             "LW_AWS_MONGODB_4": {
                 "enabled": true
             },
             "AWS_CIS_1_8": {
                 "enabled": true
             },
             "LW_AWS_MONGODB_3": {
                 "enabled": true
             },
             "AWS_CIS_1_7": {
                 "enabled": true
             },
             "LW_AWS_MONGODB_6": {
                 "enabled": true
             },
             "AWS_CIS_1_6": {
                 "enabled": true
             },
             "LW_AWS_MONGODB_5": {
                 "enabled": true
             },
             "AWS_CIS_1_5": {
                 "enabled": true
             },
             "AWS_CIS_1_4": {
                 "enabled": true
             },
             "AWS_CIS_1_3": {
                 "enabled": true
             },
             "AWS_CIS_1_2": {
                 "enabled": true
             },
             "AWS_CIS_1_11": {
                 "enabled": true
             },
             "AWS_CIS_1_12": {
                 "enabled": true
             },
             "AWS_CIS_1_10": {
                 "enabled": true
             },
             "AWS_CIS_1_15": {
                 "enabled": true
             },
             "AWS_CIS_1_16": {
                 "enabled": true
             },
             "AWS_CIS_1_13": {
                 "enabled": true
             },
             "AWS_CIS_1_14": {
                 "enabled": true
             },
             "AWS_CIS_1_1": {
                 "enabled": true
             },
             "AWS_CIS_1_19": {
                 "enabled": true
             },
             "LW_AWS_IAM_10": {
                 "enabled": true
             },
             "AWS_CIS_1_17": {
                 "enabled": true
             },
             "LW_AWS_GENERAL_SECURITY_1": {
                 "enabled": true
             },
             "LW_AWS_MONGODB_2": {
                 "enabled": true
             },
             "LW_AWS_MONGODB_1": {
                 "enabled": true
             }
         }
     ],
     "ok": true,
     "message": "SUCCESS"
 }`
}
