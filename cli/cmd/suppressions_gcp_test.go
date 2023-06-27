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

func TestConvertGcpSuppressions(t *testing.T) {
	var (
		fakeServer = lacework.MockServer()
	)

	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("suppressions/gcp/allExceptions",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "GcpList() should be a GET method")
			Suppressions := rawGcpSuppressions()
			fmt.Fprintf(w, Suppressions)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.V2.Suppressions.Gcp.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	var recsWithSuppressions []string
	for key, sup := range response {
		if len(sup.SuppressionConditions) >= 1 {
			recsWithSuppressions = append(recsWithSuppressions, key)
		}
	}
	assert.Equal(t, 5, len(recsWithSuppressions))
	expectedRecsWithSup := []string{"GCP_CIS12_1_8", "GCP_CIS12_4_5", "GCP_CIS12_6_2_5",
		"GCP_CIS12_6_3_3", "GCP_CIS12_6_3_6"}
	sort.Strings(expectedRecsWithSup)
	sort.Strings(recsWithSuppressions)
	assert.Equal(t, expectedRecsWithSup, recsWithSuppressions)

	convertedPolicyExceptions, payloadsText, discardedSuppressions := convertGcpSuppressions(
		response, genGcpPoliciesExceptionConstraintsMap())
	assert.Equal(t, 2, len(discardedSuppressions))
	assert.Equal(t, 3, len(payloadsText))
	var actualConvertedPolicyIds []string
	expectedConvertedPolicyIds := []string{"lacework-global-290", "lacework-global-268",
		"lacework-global-287"}

	for _, entry := range convertedPolicyExceptions {
		for key, _ := range entry {
			actualConvertedPolicyIds = append(actualConvertedPolicyIds, key)
		}
	}
	sort.Strings(actualConvertedPolicyIds)
	sort.Strings(expectedConvertedPolicyIds)
	assert.Equal(t, expectedConvertedPolicyIds, actualConvertedPolicyIds)
}

func TestConvertGcpSupCondition(t *testing.T) {
	gcpPoliciesExceptionConstraintsMap := genGcpPoliciesExceptionConstraintsMap()
	resourceNamesConstraint1 := convertGcpResourceNameSupConditions([]string{"foobar",
		"buzz"},
		"resourceName",
		gcpPoliciesExceptionConstraintsMap["lacework-global-234"])
	expectedResourceNamesConstraint1 := api.PolicyExceptionConstraint{
		FieldKey:    "resourceName",
		FieldValues: []any{"*/foobar", "*/buzz"},
	}

	resourceNamesConstraint2 := convertGcpResourceNameSupConditions([]string{"*"},
		"resourceName",
		gcpPoliciesExceptionConstraintsMap["lacework-global-234"])
	expectedResourceNamesConstraint2 := api.PolicyExceptionConstraint{
		FieldKey:    "resourceName",
		FieldValues: []any{"*"},
	}

	projectConstraint1 := convertSupCondition([]string{"ALL_PROJECTS"},
		"projects",
		gcpPoliciesExceptionConstraintsMap["lacework-global-234"])
	expectedProjectsConstraint1 := api.PolicyExceptionConstraint{
		FieldKey:    "projects",
		FieldValues: []any{"*"},
	}

	projectConstraint2 := convertSupCondition([]string{"foobar"},
		"projects",
		gcpPoliciesExceptionConstraintsMap["lacework-global-234"])
	expectedProjectsConstraint2 := api.PolicyExceptionConstraint{
		FieldKey:    "projects",
		FieldValues: []any{"foobar"},
	}

	organizationConstraint1 := convertSupCondition([]string{"ALL_ORGANIZATIONS"},
		"organizations",
		gcpPoliciesExceptionConstraintsMap["lacework-global-234"])
	expectedOrganizationConstraint1 := api.PolicyExceptionConstraint{
		FieldKey:    "organizations",
		FieldValues: []any{"*"},
	}

	organizationConstraint2 := convertSupCondition([]string{"foobar"},
		"organizations",
		gcpPoliciesExceptionConstraintsMap["lacework-global-234"])
	expectedOrganizationConstraint2 := api.PolicyExceptionConstraint{
		FieldKey:    "organizations",
		FieldValues: []any{"foobar"},
	}

	assert.Equal(t, expectedResourceNamesConstraint1, resourceNamesConstraint1)
	assert.Equal(t, expectedResourceNamesConstraint2, resourceNamesConstraint2)
	assert.Equal(t, expectedProjectsConstraint1, projectConstraint1)
	assert.Equal(t, expectedProjectsConstraint2, projectConstraint2)
	assert.Equal(t, expectedOrganizationConstraint1, organizationConstraint1)
	assert.Equal(t, expectedOrganizationConstraint2, organizationConstraint2)
}

func rawGcpSuppressions() string {
	return `{
	"data": [{
		"recommendationExceptions": {
			"GCP_CIS12_1_1": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_10": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_11": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_12": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_13": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_14": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_15": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_2": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_3": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_4": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_5": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_6": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_7": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_1_8": {
				"enabled": true,
				"suppressionConditions": [{
					"organizationIds": [
						"ALL_ORGANIZATIONS"
					],
					"projectIds": [
						"ALL_PROJECTS"
					],
					"resourceNames": [
						"test"
					]
				}]
			},
			"GCP_CIS12_1_9": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_1": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_10": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_11": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_12": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_2": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_3": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_4": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_5": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_6": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_7": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_8": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_2_9": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_3_1": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_3_10": {
				"enabled": false,
				"suppressionConditions": null
			},
			"GCP_CIS12_3_3": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_3_4": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_3_5": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_3_6": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_3_7": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_3_8": {
				"enabled": true,
				"suppressionConditions": []
			},
			"GCP_CIS12_3_9": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_4_1": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_4_10": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_4_11": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_4_2": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_4_3": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_4_4": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_4_5": {
				"enabled": true,
				"suppressionConditions": [{
					"organizationIds": [
						"ALL_ORGANIZATIONS"
					],
					"projectIds": [
						"ALL_PROJECTS"
					],
					"resourceNames": [
						"*"
					]
				}]
			},
			"GCP_CIS12_4_6": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_4_7": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_4_8": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_4_9": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_5_1": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_5_2": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_1_1": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_1_2": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_1_3": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_1": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_10": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_11": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_12": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_13": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_14": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_15": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_16": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_2": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_3": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_4": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_5": {
				"enabled": true,
				"suppressionConditions": [{
					"organizationIds": [
						"ALL_ORGANIZATIONS"
					],
					"projectIds": [
						"ALL_PROJECTS"
					],
					"resourceNames": [
						"*"
					]
				}]
			},
			"GCP_CIS12_6_2_6": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_7": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_8": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_2_9": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_3_1": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_3_2": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_3_3": {
				"enabled": true,
				"suppressionConditions": [{
					"organizationIds": [
						"ALL_ORGANIZATIONS"
					],
					"projectIds": [
						"ALL_PROJECTS"
					],
					"resourceNames": [
						"*"
					]
				}]
			},
			"GCP_CIS12_6_3_4": {
				"enabled": false,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_3_5": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_3_6": {
				"enabled": true,
				"suppressionConditions": [{
					"organizationIds": [
						"ALL_ORGANIZATIONS"
					],
					"projectIds": [
						"ALL_PROJECTS"
					],
					"resourceLabels": [{
						"key": "key",
						"value": "value"
					}],
					"resourceNames": [
						"*"
					]
				}]
			},
			"GCP_CIS12_6_3_7": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_4": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_5": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_6": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_6_7": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_7_1": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_7_2": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_CIS12_7_3": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_1": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_10": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_11": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_12": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_13": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_14": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_15": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_16": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_17": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_18": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_2": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_3": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_4": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_5": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_6": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_7": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_8": {
				"enabled": true,
				"suppressionConditions": null
			},
			"GCP_K8S_1_9": {
				"enabled": true,
				"suppressionConditions": null
			}
		}
	}],
	"message": "SUCCESS",
	"ok": true
}`
}

func genGcpPoliciesExceptionConstraintsMap() map[string][]string {
	return map[string][]string{
		"lacework-global-234": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-235": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-236": nil,
		"lacework-global-237": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-238": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-239": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-24": nil,
		"lacework-global-240": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-241": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-242": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-243": nil,
		"lacework-global-244": nil,
		"lacework-global-245": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-246": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-247": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-248": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-249": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-25": nil,
		"lacework-global-250": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-251": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-252": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-253": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-254": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-255": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-256": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-257": nil,
		"lacework-global-258": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-259": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-26": nil,
		"lacework-global-260": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-261": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-262": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-263": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-264": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-265": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-266": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-267": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-268": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-269": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-27": nil,
		"lacework-global-270": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-271": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-272": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-273": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-274": nil,
		"lacework-global-275": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-276": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-277": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-278": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-279": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-28": nil,
		"lacework-global-280": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-281": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-282": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-283": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-284": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-285": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-286": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-287": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-288": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-289": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-29": nil,
		"lacework-global-290": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-291": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-292": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-293": nil,
		"lacework-global-294": nil,
		"lacework-global-295": nil,
		"lacework-global-296": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-297": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-298": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-299": nil,
		"lacework-global-3":   nil,
		"lacework-global-30":  nil,
		"lacework-global-300": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-301": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-302": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-303": nil,
		"lacework-global-304": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-305": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-306": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-307": nil,
		"lacework-global-308": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-309": nil,
		"lacework-global-31":  nil,
		"lacework-global-310": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-311": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-312": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-313": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-314": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-32": nil,
		"lacework-global-33": nil,
		"lacework-global-34": {
			"accountIds",
		},
		"lacework-global-35": {
			"accountIds",
		},
		"lacework-global-36": {
			"accountIds",
			"resourceNames",
			"resourceTags",
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
		"lacework-global-4": nil,
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
		"lacework-global-42": {
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
		"lacework-global-47": {
			"accountIds",
		},
		"lacework-global-48": {
			"accountIds",
			"regionNames",
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
		"lacework-global-484": nil,
		"lacework-global-485": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-486": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-487": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-488": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-489": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-49": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-490": {
			"organizations",
			"projects",
			"resourceName",
		},
		"lacework-global-491": nil,
		"lacework-global-492": nil,
		"lacework-global-493": nil,
		"lacework-global-494": nil,
		"lacework-global-495": nil,
		"lacework-global-496": nil,
		"lacework-global-497": {
			"accountIds",
		},
		"lacework-global-498": {
			"organizations",
			"projects",
			"resourceLabel",
			"resourceName",
		},
		"lacework-global-5": nil,
		"lacework-global-50": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-51": {
			"accountIds",
			"regionNames",
			"resourceTags",
		},
		"lacework-global-52": {
			"accountIds",
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
		"lacework-global-57": {
			"accountIds",
		},
		"lacework-global-58": {
			"accountIds",
		},
		"lacework-global-59": {
			"accountIds",
		},
		"lacework-global-6": nil,
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
		"lacework-global-646": nil,
		"lacework-global-647": nil,
		"lacework-global-65": {
			"accountIds",
		},
		"lacework-global-656": nil,
		"lacework-global-657": nil,
		"lacework-global-658": nil,
		"lacework-global-659": nil,
		"lacework-global-66": {
			"accountIds",
		},
		"lacework-global-660": nil,
		"lacework-global-67": {
			"accountIds",
			"regionNames",
			"resourceNames",
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
		"lacework-global-7":  nil,
		"lacework-global-70": nil,
		"lacework-global-71": nil,
		"lacework-global-72": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-73": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-74": nil,
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
		"lacework-global-8": nil,
		"lacework-global-80": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-81": {
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
		"lacework-global-88": nil,
		"lacework-global-89": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-9": nil,
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
