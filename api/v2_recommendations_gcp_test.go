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

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestRecommendationsGcpCISGetReport(t *testing.T) {
	var (
		expectedLen = 63
		fakeServer  = lacework.MockServer()
	)

	fakeServer.MockToken("TOKEN")
	fakeServer.UseApiV2()
	fakeServer.MockAPI("recommendations/gcp",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetReport() should be a GET method")
			Recommendations := listGcpRecommendations()
			fmt.Fprintf(w, Recommendations)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Recommendations.Gcp.GetReport("CIS_1_0")
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response))
	for _, rec := range response {
		assert.NotEmpty(t, rec.ID)
	}
}

func TestRecommendationsGcpCIS12GetReport(t *testing.T) {
	var (
		expectedLen = 83
		fakeServer  = lacework.MockServer()
	)

	fakeServer.MockToken("TOKEN")
	fakeServer.UseApiV2()
	fakeServer.MockAPI("recommendations/gcp",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GetReport() should be a GET method")
			Recommendations := listGcpRecommendations()
			fmt.Fprintf(w, Recommendations)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Recommendations.Gcp.GetReport("CIS_1_2")
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLen, len(response))
	for _, rec := range response {
		assert.NotEmpty(t, rec.ID)
	}
}

func listGcpRecommendations() string {
	return `{
    "data": [
        {
            "GCP_CIS12_6_2_1": {
                "enabled": true
            },
            "GCP_K8S_1_10": {
                "enabled": false
            },
            "GCP_K8S_1_12": {
                "enabled": false
            },
            "GCP_K8S_1_11": {
                "enabled": false
            },
            "GCP_K8S_1_14": {
                "enabled": false
            },
            "GCP_K8S_1_13": {
                "enabled": false
            },
            "GCP_K8S_1_16": {
                "enabled": false
            },
            "GCP_K8S_1_15": {
                "enabled": false
            },
            "GCP_K8S_1_18": {
                "enabled": false
            },
            "GCP_K8S_1_2": {
                "enabled": false
            },
            "GCP_K8S_1_17": {
                "enabled": true
            },
            "GCP_K8S_1_3": {
                "enabled": false
            },
            "GCP_K8S_1_4": {
                "enabled": true
            },
            "GCP_K8S_1_5": {
                "enabled": true
            },
            "GCP_K8S_1_6": {
                "enabled": false
            },
            "GCP_K8S_1_7": {
                "enabled": false
            },
            "GCP_K8S_1_8": {
                "enabled": false
            },
            "GCP_K8S_1_9": {
                "enabled": false
            },
            "GCP_CIS_4_1": {
                "enabled": true
            },
            "GCP_CIS_4_3": {
                "enabled": true
            },
            "GCP_CIS_4_2": {
                "enabled": true
            },
            "GCP_CIS_4_5": {
                "enabled": true
            },
            "GCP_CIS_4_4": {
                "enabled": true
            },
            "GCP_CIS_4_6": {
                "enabled": true
            },
            "GCP_CIS_1_10": {
                "enabled": true
            },
            "GCP_CIS12_2_12": {
                "enabled": true
            },
            "GCP_CIS_1_11": {
                "enabled": true
            },
            "GCP_CIS12_2_11": {
                "enabled": true
            },
            "GCP_CIS12_2_10": {
                "enabled": true
            },
            "GCP_CIS_1_12": {
                "enabled": true
            },
            "GCP_CIS_1_13": {
                "enabled": true
            },
            "GCP_CIS12_1_5": {
                "enabled": true
            },
            "GCP_CIS12_5_2": {
                "enabled": true
            },
            "GCP_CIS12_6_2_5": {
                "enabled": true
            },
            "GCP_CIS12_1_4": {
                "enabled": true
            },
            "GCP_CIS12_6_2_4": {
                "enabled": true
            },
            "GCP_CIS12_1_3": {
                "enabled": false
            },
            "GCP_CIS12_6_2_3": {
                "enabled": true
            },
            "GCP_CIS12_1_2": {
                "enabled": true
            },
            "GCP_CIS12_6_2_2": {
                "enabled": true
            },
            "GCP_CIS12_1_9": {
                "enabled": true
            },
            "GCP_CIS12_6_2_9": {
                "enabled": true
            },
            "GCP_CIS12_1_8": {
                "enabled": true
            },
            "GCP_CIS12_6_2_8": {
                "enabled": true
            },
            "GCP_CIS12_1_7": {
                "enabled": true
            },
            "GCP_CIS12_6_2_7": {
                "enabled": true
            },
            "GCP_CIS12_1_6": {
                "enabled": true
            },
            "GCP_CIS12_5_1": {
                "enabled": false
            },
            "GCP_CIS12_6_2_6": {
                "enabled": true
            },
            "GCP_CIS12_1_1": {
                "enabled": false
            },
            "GCP_CIS12_3_10": {
                "enabled": false
            },
            "GCP_CIS12_6_1_2": {
                "enabled": true
            },
            "GCP_CIS12_6_1_1": {
                "enabled": false
            },
            "GCP_CIS_3_2": {
                "enabled": true
            },
            "GCP_CIS_7_11": {
                "enabled": true
            },
            "GCP_CIS_3_1": {
                "enabled": true
            },
            "GCP_CIS_7_10": {
                "enabled": true
            },
            "GCP_CIS_3_4": {
                "enabled": true
            },
            "GCP_CIS_7_13": {
                "enabled": true
            },
            "GCP_CIS_3_3": {
                "enabled": true
            },
            "GCP_CIS_7_12": {
                "enabled": true
            },
            "GCP_CIS_3_6": {
                "enabled": true
            },
            "GCP_CIS_7_2": {
                "enabled": true
            },
            "GCP_CIS_3_5": {
                "enabled": true
            },
            "GCP_CIS_7_1": {
                "enabled": true
            },
            "GCP_CIS_3_8": {
                "enabled": true
            },
            "GCP_CIS_7_4": {
                "enabled": true
            },
            "GCP_CIS_3_7": {
                "enabled": true
            },
            "GCP_CIS_7_3": {
                "enabled": true
            },
            "GCP_K8S_1_1": {
                "enabled": false
            },
            "GCP_CIS_7_18": {
                "enabled": true
            },
            "GCP_CIS12_2_9": {
                "enabled": true
            },
            "GCP_CIS12_6_5": {
                "enabled": false
            },
            "GCP_CIS_7_15": {
                "enabled": true
            },
            "GCP_CIS12_6_6": {
                "enabled": false
            },
            "GCP_CIS_7_14": {
                "enabled": true
            },
            "GCP_CIS12_6_7": {
                "enabled": true
            },
            "GCP_CIS_7_17": {
                "enabled": true
            },
            "GCP_CIS_7_16": {
                "enabled": true
            },
            "GCP_CIS12_2_4": {
                "enabled": false
            },
            "GCP_CIS12_2_3": {
                "enabled": true
            },
            "GCP_CIS12_2_2": {
                "enabled": true
            },
            "GCP_CIS12_2_1": {
                "enabled": true
            },
            "GCP_CIS12_6_1_3": {
                "enabled": true
            },
            "GCP_CIS12_6_4": {
                "enabled": false
            },
            "GCP_CIS12_2_8": {
                "enabled": true
            },
            "GCP_CIS12_2_7": {
                "enabled": true
            },
            "GCP_CIS12_2_6": {
                "enabled": true
            },
            "GCP_CIS12_2_5": {
                "enabled": true
            },
            "GCP_CIS_7_6": {
                "enabled": true
            },
            "GCP_CIS_7_5": {
                "enabled": true
            },
            "GCP_CIS_7_8": {
                "enabled": true
            },
            "GCP_CIS_7_7": {
                "enabled": true
            },
            "GCP_CIS_7_9": {
                "enabled": true
            },
            "GCP_CIS12_1_12": {
                "enabled": true
            },
            "GCP_CIS_2_11": {
                "enabled": true
            },
            "GCP_CIS12_1_11": {
                "enabled": true
            },
            "GCP_CIS12_1_10": {
                "enabled": true
            },
            "GCP_CIS_2_10": {
                "enabled": true
            },
            "GCP_CIS12_1_15": {
                "enabled": true
            },
            "GCP_CIS12_1_14": {
                "enabled": true
            },
            "GCP_CIS12_1_13": {
                "enabled": true
            },
            "GCP_CIS12_6_2_10": {
                "enabled": true
            },
            "GCP_CIS_2_3": {
                "enabled": true
            },
            "GCP_CIS_2_2": {
                "enabled": true
            },
            "GCP_CIS12_6_2_12": {
                "enabled": true
            },
            "GCP_CIS_2_5": {
                "enabled": true
            },
            "GCP_CIS_6_1": {
                "enabled": true
            },
            "GCP_CIS12_6_2_11": {
                "enabled": true
            },
            "GCP_CIS_2_4": {
                "enabled": true
            },
            "GCP_CIS12_6_2_14": {
                "enabled": true
            },
            "GCP_CIS_2_7": {
                "enabled": true
            },
            "GCP_CIS_6_3": {
                "enabled": true
            },
            "GCP_CIS12_6_2_13": {
                "enabled": true
            },
            "GCP_CIS_2_6": {
                "enabled": true
            },
            "GCP_CIS_6_2": {
                "enabled": true
            },
            "GCP_CIS12_6_2_16": {
                "enabled": true
            },
            "GCP_CIS_2_9": {
                "enabled": true
            },
            "GCP_CIS12_6_2_15": {
                "enabled": true
            },
            "GCP_CIS_2_8": {
                "enabled": true
            },
            "GCP_CIS_6_4": {
                "enabled": true
            },
            "GCP_CIS12_3_8": {
                "enabled": true
            },
            "GCP_CIS12_3_9": {
                "enabled": false
            },
            "GCP_CIS_2_1": {
                "enabled": true
            },
            "GCP_CIS12_3_4": {
                "enabled": true
            },
            "GCP_CIS12_3_5": {
                "enabled": true
            },
            "GCP_CIS12_7_1": {
                "enabled": true
            },
            "GCP_CIS12_3_6": {
                "enabled": false
            },
            "GCP_CIS12_7_2": {
                "enabled": true
            },
            "GCP_CIS12_3_7": {
                "enabled": false
            },
            "GCP_CIS12_7_3": {
                "enabled": true
            },
            "GCP_CIS12_3_1": {
                "enabled": true
            },
            "GCP_CIS12_3_2": {
                "enabled": true
            },
            "GCP_CIS12_3_3": {
                "enabled": false
            },
            "GCP_CIS_1_4": {
                "enabled": true
            },
            "GCP_CIS_1_3": {
                "enabled": true
            },
            "GCP_CIS12_4_10": {
                "enabled": false
            },
            "GCP_CIS_1_6": {
                "enabled": true
            },
            "GCP_CIS_5_2": {
                "enabled": true
            },
            "GCP_CIS12_4_11": {
                "enabled": false
            },
            "GCP_CIS_1_5": {
                "enabled": true
            },
            "GCP_CIS_5_1": {
                "enabled": true
            },
            "GCP_CIS_1_8": {
                "enabled": true
            },
            "GCP_CIS_1_7": {
                "enabled": true
            },
            "GCP_CIS_5_3": {
                "enabled": true
            },
            "GCP_CIS_1_9": {
                "enabled": true
            },
            "GCP_CIS12_4_7": {
                "enabled": false
            },
            "GCP_CIS12_4_8": {
                "enabled": false
            },
            "GCP_CIS12_4_9": {
                "enabled": true
            },
            "GCP_CIS_1_2": {
                "enabled": true
            },
            "GCP_CIS_1_1": {
                "enabled": true
            },
            "GCP_CIS12_4_3": {
                "enabled": false
            },
            "GCP_CIS12_6_3_4": {
                "enabled": true
            },
            "GCP_CIS12_4_4": {
                "enabled": true
            },
            "GCP_CIS12_6_3_3": {
                "enabled": true
            },
            "GCP_CIS12_4_5": {
                "enabled": false
            },
            "GCP_CIS12_6_3_2": {
                "enabled": true
            },
            "GCP_CIS12_4_6": {
                "enabled": true
            },
            "GCP_CIS12_6_3_1": {
                "enabled": true
            },
            "GCP_CIS12_6_3_7": {
                "enabled": true
            },
            "GCP_CIS12_4_1": {
                "enabled": false
            },
            "GCP_CIS12_6_3_6": {
                "enabled": true
            },
            "GCP_CIS12_4_2": {
                "enabled": false
            },
            "GCP_CIS12_6_3_5": {
                "enabled": false
            }
        }
    ],
    "ok": true,
    "message": "SUCCESS"
}`
}
