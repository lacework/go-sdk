//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestRecommendationsAzureCISGetReport(t *testing.T) {
	var (
		expectedLen = 89
		fakeServer  = lacework.MockServer()
	)

	fakeServer.MockToken("TOKEN")
	fakeServer.UseApiV2()
	fakeServer.MockAPI("recommendations/azure",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetReport() should be a GET method")
			Recommendations := listAzureRecommendations()
			fmt.Fprintf(w, Recommendations)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Recommendations.Azure.GetReport("CIS_1_0")
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response))
	for _, rec := range response {
		assert.NotEmpty(t, rec.ID)
	}
}

func TestRecommendationsAzureCIS131GetReport(t *testing.T) {
	var (
		expectedLen = 114
		fakeServer  = lacework.MockServer()
	)

	fakeServer.MockToken("TOKEN")
	fakeServer.UseApiV2()
	fakeServer.MockAPI("recommendations/azure",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetReport() should be a GET method")
			Recommendations := listAzureRecommendations()
			fmt.Fprintf(w, Recommendations)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Recommendations.Azure.GetReport("CIS_1_3_1")
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response))
	for _, rec := range response {
		assert.NotEmpty(t, rec.ID)
	}
}

func listAzureRecommendations() string {
	return `{
    "data": [
        {
            "Azure_CIS_5_3": {
                "enabled": true
            },
            "Azure_CIS_5_2": {
                "enabled": true
            },
            "Azure_CIS_5_1": {
                "enabled": true
            },
            "Azure_CIS_131_7_1": {
                "enabled": true
            },
            "Azure_CIS_131_7_3": {
                "enabled": true
            },
            "Azure_CIS_131_7_2": {
                "enabled": true
            },
            "Azure_CIS_131_7_5": {
                "enabled": true
            },
            "Azure_CIS_5_9": {
                "enabled": true
            },
            "Azure_CIS_131_7_4": {
                "enabled": true
            },
            "Azure_CIS_5_8": {
                "enabled": true
            },
            "Azure_CIS_131_7_7": {
                "enabled": true
            },
            "Azure_CIS_5_7": {
                "enabled": true
            },
            "Azure_CIS_131_7_6": {
                "enabled": true
            },
            "Azure_CIS_5_6": {
                "enabled": true
            },
            "Azure_CIS_5_5": {
                "enabled": true
            },
            "Azure_CIS_5_4": {
                "enabled": true
            },
            "Azure_CIS_5_11": {
                "enabled": true
            },
            "Azure_CIS_5_12": {
                "enabled": true
            },
            "Azure_CIS_5_13": {
                "enabled": true
            },
            "Azure_CIS_5_10": {
                "enabled": true
            },
            "Azure_CIS_6_2": {
                "enabled": true
            },
            "Azure_CIS_6_1": {
                "enabled": true
            },
            "Azure_CIS_131_8_2": {
                "enabled": false
            },
            "Azure_CIS_131_8_1": {
                "enabled": true
            },
            "Azure_CIS_131_8_4": {
                "enabled": true
            },
            "Azure_CIS_131_8_3": {
                "enabled": true
            },
            "Azure_CIS_131_8_5": {
                "enabled": true
            },
            "Azure_CIS_6_5": {
                "enabled": true
            },
            "Azure_CIS_6_4": {
                "enabled": true
            },
            "Azure_CIS_6_3": {
                "enabled": true
            },
            "Azure_CIS_131_1_1": {
                "enabled": true
            },
            "Azure_CIS_3_5": {
                "enabled": true
            },
            "Azure_CIS_3_4": {
                "enabled": true
            },
            "Azure_CIS_131_1_3": {
                "enabled": true
            },
            "Azure_CIS_131_4_1_1": {
                "enabled": true
            },
            "Azure_CIS_3_3": {
                "enabled": true
            },
            "Azure_CIS_131_1_2": {
                "enabled": true
            },
            "Azure_CIS_131_4_1_2": {
                "enabled": true
            },
            "Azure_CIS_3_2": {
                "enabled": true
            },
            "Azure_CIS_131_1_5": {
                "enabled": true
            },
            "Azure_CIS_3_1": {
                "enabled": true
            },
            "Azure_CIS_131_1_4": {
                "enabled": true
            },
            "Azure_CIS_131_1_7": {
                "enabled": true
            },
            "Azure_CIS_131_1_6": {
                "enabled": true
            },
            "Azure_CIS_131_1_9": {
                "enabled": true
            },
            "Azure_CIS_131_9_1": {
                "enabled": true
            },
            "Azure_CIS_131_1_8": {
                "enabled": true
            },
            "Azure_CIS_131_9_3": {
                "enabled": true
            },
            "Azure_CIS_131_9_2": {
                "enabled": true
            },
            "Azure_CIS_131_4_1_3": {
                "enabled": true
            },
            "Azure_CIS_131_9_5": {
                "enabled": true
            },
            "Azure_CIS_131_9_4": {
                "enabled": true
            },
            "Azure_CIS_131_9_7": {
                "enabled": true
            },
            "Azure_CIS_3_7": {
                "enabled": true
            },
            "Azure_CIS_131_9_6": {
                "enabled": true
            },
            "Azure_CIS_3_6": {
                "enabled": true
            },
            "Azure_CIS_131_1_10": {
                "enabled": true
            },
            "Azure_CIS_131_9_9": {
                "enabled": true
            },
            "Azure_CIS_131_9_8": {
                "enabled": true
            },
            "Azure_CIS_131_2_2": {
                "enabled": true
            },
            "Azure_CIS_131_2_1": {
                "enabled": true
            },
            "Azure_CIS_131_4_2_1": {
                "enabled": true
            },
            "Azure_CIS_131_2_4": {
                "enabled": true
            },
            "Azure_CIS_131_2_3": {
                "enabled": true
            },
            "Azure_CIS_131_2_6": {
                "enabled": true
            },
            "Azure_CIS_131_2_5": {
                "enabled": true
            },
            "Azure_CIS_131_2_8": {
                "enabled": true
            },
            "Azure_CIS_131_2_7": {
                "enabled": true
            },
            "Azure_CIS_131_2_9": {
                "enabled": true
            },
            "Azure_CIS_131_4_2_2": {
                "enabled": true
            },
            "Azure_CIS_131_4_2_3": {
                "enabled": true
            },
            "Azure_CIS_131_4_2_4": {
                "enabled": true
            },
            "Azure_CIS_131_4_2_5": {
                "enabled": true
            },
            "Azure_CIS_4_1_5": {
                "enabled": true
            },
            "Azure_CIS_4_1_6": {
                "enabled": true
            },
            "Azure_CIS_4_1_3": {
                "enabled": true
            },
            "Azure_CIS_4_1_4": {
                "enabled": true
            },
            "Azure_CIS_4_1_7": {
                "enabled": true
            },
            "Azure_CIS_4_1_8": {
                "enabled": true
            },
            "Azure_CIS_4_1_1": {
                "enabled": true
            },
            "Azure_CIS_4_1_2": {
                "enabled": true
            },
            "Azure_CIS_1_7": {
                "enabled": true
            },
            "Azure_CIS_1_6": {
                "enabled": true
            },
            "Azure_CIS_131_3_1": {
                "enabled": true
            },
            "Azure_CIS_1_5": {
                "enabled": true
            },
            "Azure_CIS_1_4": {
                "enabled": true
            },
            "Azure_CIS_131_3_3": {
                "enabled": true
            },
            "Azure_CIS_1_3": {
                "enabled": true
            },
            "Azure_CIS_131_3_2": {
                "enabled": true
            },
            "Azure_CIS_1_2": {
                "enabled": true
            },
            "Azure_CIS_131_3_5": {
                "enabled": true
            },
            "Azure_CIS_1_1": {
                "enabled": true
            },
            "Azure_CIS_131_3_4": {
                "enabled": true
            },
            "Azure_CIS_131_3_7": {
                "enabled": true
            },
            "Azure_CIS_131_4_3_5": {
                "enabled": true
            },
            "Azure_CIS_131_3_6": {
                "enabled": true
            },
            "Azure_CIS_131_4_3_6": {
                "enabled": true
            },
            "Azure_CIS_131_3_9": {
                "enabled": true
            },
            "Azure_CIS_131_4_3_7": {
                "enabled": true
            },
            "Azure_CIS_131_3_8": {
                "enabled": true
            },
            "Azure_CIS_131_4_3_8": {
                "enabled": true
            },
            "Azure_CIS_131_4_3_1": {
                "enabled": true
            },
            "Azure_CIS_131_4_3_2": {
                "enabled": true
            },
            "Azure_CIS_131_4_3_3": {
                "enabled": true
            },
            "Azure_CIS_1_9": {
                "enabled": true
            },
            "Azure_CIS_131_4_3_4": {
                "enabled": true
            },
            "Azure_CIS_1_8": {
                "enabled": true
            },
            "Azure_CIS_131_2_11": {
                "enabled": true
            },
            "Azure_CIS_131_2_10": {
                "enabled": true
            },
            "Azure_CIS_131_5_1_1": {
                "enabled": true
            },
            "Azure_CIS_131_5_1_4": {
                "enabled": true
            },
            "Azure_CIS_131_5_1_5": {
                "enabled": true
            },
            "Azure_CIS_131_5_1_2": {
                "enabled": true
            },
            "Azure_CIS_131_5_1_3": {
                "enabled": true
            },
            "Azure_CIS_131_2_15": {
                "enabled": true
            },
            "Azure_CIS_131_2_14": {
                "enabled": true
            },
            "Azure_CIS_131_2_13": {
                "enabled": true
            },
            "Azure_CIS_131_2_12": {
                "enabled": true
            },
            "Azure_CIS_4_2_4": {
                "enabled": true
            },
            "Azure_CIS_4_2_5": {
                "enabled": true
            },
            "Azure_CIS_4_2_2": {
                "enabled": true
            },
            "Azure_CIS_4_2_3": {
                "enabled": true
            },
            "Azure_CIS_4_2_8": {
                "enabled": true
            },
            "Azure_CIS_4_2_6": {
                "enabled": true
            },
            "Azure_CIS_4_2_7": {
                "enabled": true
            },
            "Azure_CIS_4_2_1": {
                "enabled": true
            },
            "Azure_CIS_2_6": {
                "enabled": true
            },
            "Azure_CIS_2_5": {
                "enabled": true
            },
            "Azure_CIS_131_5_2_9": {
                "enabled": true
            },
            "Azure_CIS_2_3": {
                "enabled": true
            },
            "Azure_CIS_2_1": {
                "enabled": true
            },
            "Azure_CIS_131_4_4": {
                "enabled": true
            },
            "Azure_CIS_131_1_19": {
                "enabled": true
            },
            "Azure_CIS_131_1_18": {
                "enabled": true
            },
            "Azure_CIS_131_1_17": {
                "enabled": true
            },
            "Azure_CIS_131_4_5": {
                "enabled": true
            },
            "Azure_CIS_131_1_16": {
                "enabled": true
            },
            "Azure_CIS_131_1_15": {
                "enabled": true
            },
            "Azure_CIS_131_1_14": {
                "enabled": true
            },
            "Azure_CIS_131_1_13": {
                "enabled": true
            },
            "Azure_CIS_2_9": {
                "enabled": true
            },
            "Azure_CIS_131_1_12": {
                "enabled": true
            },
            "Azure_CIS_2_8": {
                "enabled": true
            },
            "Azure_CIS_131_1_11": {
                "enabled": true
            },
            "Azure_CIS_2_7": {
                "enabled": true
            },
            "Azure_CIS_131_1_21": {
                "enabled": true
            },
            "Azure_CIS_131_1_20": {
                "enabled": true
            },
            "Azure_CIS_131_5_2_3": {
                "enabled": true
            },
            "Azure_CIS_131_5_2_4": {
                "enabled": true
            },
            "Azure_CIS_131_5_2_1": {
                "enabled": true
            },
            "Azure_CIS_131_5_2_2": {
                "enabled": true
            },
            "Azure_CIS_131_5_2_7": {
                "enabled": true
            },
            "Azure_CIS_131_5_2_8": {
                "enabled": true
            },
            "Azure_CIS_131_5_2_5": {
                "enabled": true
            },
            "Azure_CIS_131_5_2_6": {
                "enabled": true
            },
            "Azure_CIS_131_1_23": {
                "enabled": true
            },
            "Azure_CIS_131_1_22": {
                "enabled": true
            },
            "Azure_CIS_7_1": {
                "enabled": true
            },
            "Azure_CIS_1_22": {
                "enabled": true
            },
            "Azure_CIS_1_23": {
                "enabled": true
            },
            "Azure_CIS_131_5_3": {
                "enabled": true
            },
            "Azure_CIS_7_6": {
                "enabled": true
            },
            "Azure_CIS_7_5": {
                "enabled": true
            },
            "Azure_CIS_7_4": {
                "enabled": true
            },
            "Azure_CIS_7_3": {
                "enabled": true
            },
            "Azure_CIS_7_2": {
                "enabled": true
            },
            "Azure_CIS_1_20": {
                "enabled": true
            },
            "Azure_CIS_1_21": {
                "enabled": true
            },
            "Azure_CIS_1_13": {
                "enabled": true
            },
            "Azure_CIS_1_14": {
                "enabled": true
            },
            "Azure_CIS_1_11": {
                "enabled": true
            },
            "Azure_CIS_1_12": {
                "enabled": true
            },
            "Azure_CIS_1_17": {
                "enabled": true
            },
            "Azure_CIS_131_9_10": {
                "enabled": true
            },
            "Azure_CIS_1_15": {
                "enabled": true
            },
            "Azure_CIS_1_16": {
                "enabled": true
            },
            "Azure_CIS_131_9_11": {
                "enabled": true
            },
            "Azure_CIS_1_19": {
                "enabled": true
            },
            "Azure_LW_IAM_3": {
                "enabled": true
            },
            "Azure_LW_IAM_2": {
                "enabled": true
            },
            "Azure_LW_IAM_1": {
                "enabled": true
            },
            "Azure_CIS_1_10": {
                "enabled": true
            },
            "Azure_CIS_2_15": {
                "enabled": true
            },
            "Azure_CIS_2_14": {
                "enabled": true
            },
            "Azure_CIS_2_13": {
                "enabled": true
            },
            "Azure_CIS_2_12": {
                "enabled": true
            },
            "Azure_CIS_2_19": {
                "enabled": true
            },
            "Azure_CIS_2_18": {
                "enabled": true
            },
            "Azure_CIS_131_6_2": {
                "enabled": true
            },
            "Azure_CIS_2_17": {
                "enabled": true
            },
            "Azure_CIS_131_6_1": {
                "enabled": true
            },
            "Azure_CIS_2_16": {
                "enabled": true
            },
            "Azure_CIS_131_6_4": {
                "enabled": true
            },
            "Azure_CIS_131_6_3": {
                "enabled": true
            },
            "Azure_CIS_131_6_6": {
                "enabled": true
            },
            "Azure_CIS_131_6_5": {
                "enabled": true
            },
            "Azure_CIS_8_3": {
                "enabled": true
            },
            "Azure_CIS_8_2": {
                "enabled": true
            },
            "Azure_CIS_8_1": {
                "enabled": true
            },
            "Azure_LW_Monitoring_1": {
                "enabled": false
            },
            "Azure_CIS_2_11": {
                "enabled": true
            },
            "Azure_CIS_2_10": {
                "enabled": true
            },
            "Azure_CIS_131_3_10": {
                "enabled": true
            },
            "Azure_CIS_131_3_11": {
                "enabled": true
            }
        }
    ],
    "ok": true,
    "message": "SUCCESS"
}`
}
