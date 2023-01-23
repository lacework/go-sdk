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

package cmd

import (
	"fmt"
	"net/http"
	"sort"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

func TestConvertAwsSuppressions(t *testing.T) {
	var (
		fakeServer = lacework.MockServer()
	)

	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("suppressions/aws/allExceptions",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "AwsList() should be a GET method")
			Suppressions := rawAwsSuppressions()
			fmt.Fprintf(w, Suppressions)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Suppressions.Aws.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	var recsWithSuppressions []string
	suppCondCount := 0
	for key, sup := range response {
		if len(sup.SuppressionConditions) >= 1 {
			recsWithSuppressions = append(recsWithSuppressions, key)
			suppCondCount += len(sup.SuppressionConditions)
		}
	}
	assert.Equal(t, 5, len(recsWithSuppressions))
	expectedRecsWithSup := []string{"AWS_CIS_1_14", "AWS_CIS_2_7", "AWS_CIS_4_4",
		"LW_AWS_NETWORKING_50",
		"LW_S3_1"}
	sort.Strings(expectedRecsWithSup)
	sort.Strings(recsWithSuppressions)
	assert.Equal(t, expectedRecsWithSup, recsWithSuppressions)

	convertedPolicyExceptions, payloadsText, discardedSuppressions := convertAwsSuppressions(
		response, genAwsPoliciesExceptionConstraintsMap())
	assert.Equal(t, 1, len(discardedSuppressions))
	assert.Equal(t, suppCondCount, len(payloadsText))
	var actualConvertedPolicyIds []string
	expectedConvertedPolicyIds := []string{"lacework-global-129", "lacework-global-87",
		"lacework-global-69", "lacework-global-130", "lacework-global-130", "lacework-global-130",
		"lacework-global-77"}

	for _, entry := range convertedPolicyExceptions {
		for key, _ := range entry {
			actualConvertedPolicyIds = append(actualConvertedPolicyIds, key)
		}
	}
	sort.Strings(actualConvertedPolicyIds)
	sort.Strings(expectedConvertedPolicyIds)
	assert.Equal(t, expectedConvertedPolicyIds, actualConvertedPolicyIds)
}

func TestConvertAwsResourceNamesSupCondition(t *testing.T) {
	awsPoliciesExceptionConstraintsMap := genAwsPoliciesExceptionConstraintsMap()
	resourceNamesConstraint := convertSupCondition([]string{"foobar",
		"arn:partition:service:region:account-id:resource-id"},
		"resourceNames",
		awsPoliciesExceptionConstraintsMap["lacework-global-100"])
	expectedResourceNamesConstraint := api.PolicyExceptionConstraint{
		FieldKey:    "resourceNames",
		FieldValues: []any{"foobar", "resource-id"},
	}
	assert.Equal(t, expectedResourceNamesConstraint, resourceNamesConstraint)
}

func rawAwsSuppressions() string {
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

func genAwsPoliciesExceptionConstraintsMap() map[string][]string {
	return map[string][]string{
		"lacework-global-100": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-101": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-102": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-103": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-104": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-105": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-106": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-107": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-108": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-109": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-110": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-111": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-112": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-113": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-114": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-115": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-116": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-117": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-118": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-119": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-120": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-121": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-122": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-123": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-124": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-125": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-126": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-127": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-128": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-129": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-130": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-131": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-132": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-133": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-134": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-135": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-136": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-137": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-138": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-139": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-140": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-141": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-142": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-143": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-144": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-145": {
			"accountIds",
			"regionNames",
			"resourceNames",
		},
		"lacework-global-146": {
			"accountIds",
			"regionNames",
			"resourceNames",
		},
		"lacework-global-147": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-148": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-149": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-150": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-151": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-152": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-153": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-154": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-155": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-156": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-157": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-159": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-160": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-161": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-171": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-179": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-180": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-181": {
			"accountIds",
		},
		"lacework-global-182": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-183": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-184": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-196": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-197": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-198": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-199": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-217": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-218": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-219": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-220": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-221": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-222": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-223": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-225": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-226": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-227": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-228": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-229": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-230": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-231": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-31": nil,
		"lacework-global-32": nil,
		"lacework-global-33": nil,
		"lacework-global-34": {
			"accountIds",
		},
		"lacework-global-35": {
			"accountIds",
		},
		"lacework-global-37": {
			"accountIds",
		},
		"lacework-global-38": {
			"accountIds",
		},
		"lacework-global-39": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-40": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-41": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-43": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-44": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-45": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-46": {
			"accountIds",
		},
		"lacework-global-482": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-483": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-53": {
			"accountIds",
		},
		"lacework-global-54": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-55": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-56": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-58": {
			"accountIds",
		},
		"lacework-global-59": {
			"accountIds",
		},
		"lacework-global-60": {
			"accountIds",
		},
		"lacework-global-61": {
			"accountIds",
		},
		"lacework-global-62": {
			"accountIds",
		},
		"lacework-global-63": {
			"accountIds",
		},
		"lacework-global-64": {
			"accountIds",
		},
		"lacework-global-65": {
			"accountIds",
		},
		"lacework-global-68": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-69": {
			"accountIds",
		},
		"lacework-global-70": nil,
		"lacework-global-75": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-76": {
			"accountIds",
			"regionNames",
		},
		"lacework-global-77": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-78": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-79": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-82": {
			"accountIds",
		},
		"lacework-global-83": {
			"accountIds",
		},
		"lacework-global-84": {
			"accountIds",
		},
		"lacework-global-85": {
			"accountIds",
		},
		"lacework-global-86": {
			"accountIds",
		},
		"lacework-global-87": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-89": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-90": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-91": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-92": {
			"accountIds",
		},
		"lacework-global-93": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-94": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-95": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-96": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-97": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-98": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-99": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
	}
}
