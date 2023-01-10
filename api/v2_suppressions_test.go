//
// Author:: Ross Moles (<ross.moles@lacework.net>)
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
	"sort"
	"testing"

	"github.com/lacework/go-sdk/internal/lacework"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
)

func TestSuppressionsAwsList(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("suppressions/aws/allExceptions",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "AwsList() should be a GET method")
			Suppressions := generateSuppressions()
			fmt.Fprintf(w, Suppressions)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Suppressions.Aws.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	var recsWithSuppressions []string
	for key, sup := range response {
		if len(sup.SuppressionConditions) >= 1 {
			recsWithSuppressions = append(recsWithSuppressions, key)
		}
	}
	assert.Equal(t, 5, len(recsWithSuppressions))
	expectedRecsWithSup := []string{"AWS_CIS_1_14", "AWS_CIS_2_7", "AWS_CIS_4_4",
		"LW_AWS_NETWORKING_50",
		"LW_S3_1"}
	sort.Strings(expectedRecsWithSup)
	sort.Strings(recsWithSuppressions)
	assert.Equal(t, expectedRecsWithSup, recsWithSuppressions)
}

func generateSuppressions() string {
	return `{
  "data": [
    {
      "recommendationExceptions": {
        "AWS_CIS_1_1": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_10": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_11": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_12": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_13": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_14": {
          "enabled": true,
          "suppressionConditions": [
            {
              "accountIds": [
                "ALL_ACCOUNTS"
              ],
              "comments": "",
              "regionNames": [
                "ALL_REGIONS"
              ],
              "resourceNames": [],
              "resourceTags": []
            }
          ]
        },
        "AWS_CIS_1_15": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_16": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_17": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_19": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_2": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_20": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_21": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_22": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_23": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_24": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_3": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_4": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_5": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_6": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_7": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_8": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_1_9": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_2_1": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_2_2": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_2_3": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_2_4": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_2_5": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_2_6": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_2_7": {
          "enabled": true,
          "suppressionConditions": [
            {
              "accountIds": [
                "ALL_ACCOUNTS"
              ],
              "comments": "",
              "regionNames": [
                "eu-central-1",
                "eu-north-1"
              ],
              "resourceNames": [
                "*"
              ],
              "resourceTags": []
            }
          ]
        },
        "AWS_CIS_2_8": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_2_9": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_10": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_11": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_12": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_13": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_14": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_15": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_2": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_3": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_4": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_5": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_6": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_7": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_8": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_3_9": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_4_1": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_4_2": {
          "enabled": true,
          "suppressionConditions": null
        },
        "AWS_CIS_4_4": {
          "enabled": true,
          "suppressionConditions": [
            {
              "accountIds": [
                "ALL_ACCOUNTS"
              ],
              "comments": "this is a dev exception",
              "regionNames": [
                "ALL_REGIONS"
              ],
              "resourceNames": [
                "*"
              ],
              "resourceTags": [
                {
                  "key": "owner",
                  "value": "dev"
                },
                {
                  "key": "env",
                  "value": "dev1"
                }
              ]
            }
          ]
        },
        "AWS_CIS_4_5": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_ELASTICSEARCH_1": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_ELASTICSEARCH_2": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_ELASTICSEARCH_3": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_ELASTICSEARCH_4": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_GENERAL_SECURITY_1": {
          "enabled": true,
          "suppressionConditions": []
        },
        "LW_AWS_GENERAL_SECURITY_2": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_GENERAL_SECURITY_3": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_GENERAL_SECURITY_4": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_GENERAL_SECURITY_5": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_GENERAL_SECURITY_6": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_GENERAL_SECURITY_7": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_GENERAL_SECURITY_8": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_1": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_10": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_11": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_12": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_13": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_14": {
          "enabled": true,
          "suppressionConditions": []
        },
        "LW_AWS_IAM_2": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_3": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_4": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_5": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_6": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_7": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_8": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_IAM_9": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_MONGODB_1": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_MONGODB_2": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_MONGODB_3": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_MONGODB_4": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_MONGODB_5": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_MONGODB_6": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_1": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_10": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_11": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_12": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_13": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_14": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_15": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_16": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_17": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_18": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_19": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_2": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_20": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_21": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_22": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_23": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_24": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_25": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_26": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_27": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_28": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_29": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_3": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_30": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_31": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_32": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_33": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_34": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_35": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_36": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_37": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_38": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_39": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_4": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_40": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_41": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_42": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_43": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_44": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_45": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_46": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_47": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_48": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_49": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_5": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_50": {
          "enabled": true,
          "suppressionConditions": [
            {
              "accountIds": [
                "287105300711"
              ],
              "comments": "Test Suppressions for TA - LINK-819",
              "regionNames": [
                "ALL_REGIONS"
              ],
              "resourceNames": [],
              "resourceTags": [
                {
                  "key": "http_https_redirect",
                  "value": "true"
                }
              ]
            }
          ]
        },
        "LW_AWS_NETWORKING_51": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_6": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_7": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_8": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_NETWORKING_9": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_RDS_1": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_SERVERLESS_1": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_SERVERLESS_2": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_AWS_SERVERLESS_3": {
          "enabled": true,
          "suppressionConditions": []
        },
        "LW_AWS_SERVERLESS_4": {
          "enabled": true,
          "suppressionConditions": []
        },
        "LW_AWS_SERVERLESS_5": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_1": {
          "enabled": true,
          "suppressionConditions": [
            {
              "accountIds": [
                "287105300711"
              ],
              "comments": "",
              "regionNames": [
                "ALL_REGIONS"
              ],
              "resourceNames": [
                "eco-eng-cw-delivery-channel-bucket"
              ],
              "resourceTags": []
            },
            {
              "accountIds": [
                "ALL_ACCOUNTS"
              ],
              "comments": "A suppression comment",
              "regionNames": [
                "ALL_REGIONS"
              ],
              "resourceNames": [
                "Test"
              ],
              "resourceTags": []
            },
            {
              "accountIds": [
                "ALL_ACCOUNTS"
              ],
              "comments": "",
              "regionNames": [
                "ALL_REGIONS"
              ],
              "resourceNames": [
                "*"
              ],
              "resourceTags": [
                {
                  "key": "FOO",
                  "value": "BAR"
                }
              ]
            }
          ]
        },
        "LW_S3_10": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_11": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_12": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_13": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_14": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_15": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_16": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_17": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_18": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_19": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_2": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_20": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_21": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_3": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_4": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_5": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_6": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_7": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_8": {
          "enabled": true,
          "suppressionConditions": null
        },
        "LW_S3_9": {
          "enabled": true,
          "suppressionConditions": null
        }
      }
    }
  ],
  "message": "SUCCESS",
  "ok": true
}`
}
