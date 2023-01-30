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

func TestConvertAzureSuppressions(t *testing.T) {
	var (
		fakeServer = lacework.MockServer()
	)

	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("suppressions/azure/allExceptions",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "AzureList() should be a GET method")
			Suppressions := rawAzureSuppressions()
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

	response, err := c.V2.Suppressions.Azure.List()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	var recsWithSuppressions []string
	for key, sup := range response {
		if len(sup.SuppressionConditions) >= 1 {
			recsWithSuppressions = append(recsWithSuppressions, key)
		}
	}
	assert.Equal(t, 13, len(recsWithSuppressions))
	expectedRecsWithSup := []string{"Azure_CIS_131_3_10", "Azure_CIS_131_3_11", "Azure_CIS_131_3_5",
		"Azure_CIS_131_3_6", "Azure_CIS_131_3_8", "Azure_CIS_131_5_1_4", "Azure_CIS_131_5_1_5",
		"Azure_CIS_131_7_2", "Azure_CIS_131_7_3", "Azure_CIS_131_9_1", "Azure_CIS_131_9_10",
		"Azure_CIS_131_9_5", "Azure_LW_Monitoring_1"}
	sort.Strings(expectedRecsWithSup)
	sort.Strings(recsWithSuppressions)
	assert.Equal(t, expectedRecsWithSup, recsWithSuppressions)

	convertedPolicyExceptions, payloadsText, discardedSuppressions := convertAzureSuppressions(
		response, genAzurePoliciesExceptionConstraintsMap())
	assert.Equal(t, 5, len(discardedSuppressions))
	assert.Equal(t, 71, len(payloadsText))
	var actualConvertedPolicyIds []string
	expectedConvertedPolicyIds := []string{"lacework-global-532", "lacework-global-532",
		"lacework-global-532", "lacework-global-532", "lacework-global-532", "lacework-global-532",
		"lacework-global-532", "lacework-global-532", "lacework-global-533", "lacework-global-533",
		"lacework-global-533", "lacework-global-533", "lacework-global-533", "lacework-global-533",
		"lacework-global-533", "lacework-global-533", "lacework-global-533", "lacework-global-533",
		"lacework-global-533", "lacework-global-533", "lacework-global-533", "lacework-global-533",
		"lacework-global-533", "lacework-global-533", "lacework-global-533", "lacework-global-533",
		"lacework-global-557", "lacework-global-557", "lacework-global-557", "lacework-global-557",
		"lacework-global-557", "lacework-global-557", "lacework-global-557", "lacework-global-557",
		"lacework-global-557", "lacework-global-557", "lacework-global-557", "lacework-global-557",
		"lacework-global-557", "lacework-global-557", "lacework-global-557", "lacework-global-582",
		"lacework-global-582", "lacework-global-582", "lacework-global-582", "lacework-global-582",
		"lacework-global-582", "lacework-global-582", "lacework-global-582", "lacework-global-582",
		"lacework-global-587", "lacework-global-635", "lacework-global-635", "lacework-global-635",
		"lacework-global-635", "lacework-global-635", "lacework-global-635", "lacework-global-635",
		"lacework-global-635", "lacework-global-635", "lacework-global-636", "lacework-global-636",
		"lacework-global-642", "lacework-global-642", "lacework-global-642", "lacework-global-642",
		"lacework-global-642", "lacework-global-642", "lacework-global-642", "lacework-global-642",
		"lacework-global-642"}

	for _, entry := range convertedPolicyExceptions {
		for key, _ := range entry {
			actualConvertedPolicyIds = append(actualConvertedPolicyIds, key)
		}
	}
	sort.Strings(actualConvertedPolicyIds)
	sort.Strings(expectedConvertedPolicyIds)
	assert.Equal(t, expectedConvertedPolicyIds, actualConvertedPolicyIds)
}

func TestConvertAzureSupCondition(t *testing.T) {
	azurePoliciesExceptionConstraintsMap := genAzurePoliciesExceptionConstraintsMap()
	resourceNamesConstraint1 := convertSupCondition([]string{"isstoragecentralus"},
		"resourceName",
		azurePoliciesExceptionConstraintsMap["lacework-global-642"])
	expectedResourceNamesConstraint1 := api.PolicyExceptionConstraint{
		FieldKey:    "resourceName",
		FieldValues: []any{"isstoragecentralus"},
	}

	resourceNamesConstraint2 := convertSupCondition([]string{"*"},
		"resourceName",
		azurePoliciesExceptionConstraintsMap["lacework-global-642"])
	expectedResourceNamesConstraint2 := api.PolicyExceptionConstraint{
		FieldKey:    "resourceName",
		FieldValues: []any{"*"},
	}

	projectConstraint1 := convertSupCondition([]string{"ALL_TENANTS"},
		"tenants",
		azurePoliciesExceptionConstraintsMap["lacework-global-642"])
	expectedProjectsConstraint1 := api.PolicyExceptionConstraint{
		FieldKey:    "tenants",
		FieldValues: []any{"*"},
	}

	projectConstraint2 := convertSupCondition([]string{"foobar"},
		"tenants",
		azurePoliciesExceptionConstraintsMap["lacework-global-642"])
	expectedProjectsConstraint2 := api.PolicyExceptionConstraint{
		FieldKey:    "tenants",
		FieldValues: []any{"foobar"},
	}

	organizationConstraint1 := convertSupCondition([]string{"ALL_SUBSCRIPTIONS"},
		"subscriptions",
		azurePoliciesExceptionConstraintsMap["lacework-global-642"])
	expectedOrganizationConstraint1 := api.PolicyExceptionConstraint{
		FieldKey:    "subscriptions",
		FieldValues: []any{"*"},
	}

	organizationConstraint2 := convertSupCondition([]string{"foobar"},
		"subscriptions",
		azurePoliciesExceptionConstraintsMap["lacework-global-642"])
	expectedOrganizationConstraint2 := api.PolicyExceptionConstraint{
		FieldKey:    "subscriptions",
		FieldValues: []any{"foobar"},
	}

	assert.Equal(t, expectedResourceNamesConstraint1, resourceNamesConstraint1)
	assert.Equal(t, expectedResourceNamesConstraint2, resourceNamesConstraint2)
	assert.Equal(t, expectedProjectsConstraint1, projectConstraint1)
	assert.Equal(t, expectedProjectsConstraint2, projectConstraint2)
	assert.Equal(t, expectedOrganizationConstraint1, organizationConstraint1)
	assert.Equal(t, expectedOrganizationConstraint2, organizationConstraint2)
}

func rawAzureSuppressions() string {
	return `{
	"data": [{
		"recommendationExceptions": {
			  "Azure_CIS_131_1_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_10": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_11": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_12": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_13": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_14": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_15": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_16": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_17": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_18": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_19": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_20": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_21": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_22": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_23": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_4": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_5": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_6": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_7": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_8": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_1_9": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_10": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_11": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_12": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_13": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_14": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_15": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_4": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_5": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_6": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_7": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_8": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_2_9": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_3_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_3_10": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceNames": [
					  "isstoragecentralus"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_3_11": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceNames": [
					  "sqlbackupstaxslayer"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "supportrecordings"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "tswslogs"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "isstoragecentralus"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_3_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_3_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_3_4": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_3_5": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceNames": [
					  "consumerandroid"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "consumerios"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "taxestogoandroid"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "taxestogoios"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "azurecdnlogsdevsoa"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "azurecdnlogsprodevsoa"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "azurecdnlogsprodsoa"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "azurecdnlogstsprod"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_3_6": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceGroupNames": [
					  "RG-Mobile"
					],
					"resourceNames": [
					  "consumerandroid"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Mobile"
					],
					"resourceNames": [
					  "consumerios"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Mobile"
					],
					"resourceNames": [
					  "taxestogoandroid"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Mobile"
					],
					"resourceNames": [
					  "taxestogoios"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Consumer-Dev-Microservices"
					],
					"resourceNames": [
					  "azurecdnlogsdevsoa"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Pro-Dev-Microservices"
					],
					"resourceNames": [
					  "azurecdnlogsprodevsoa"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Consumer-Prod-Microservices"
					],
					"resourceNames": [
					  "azurecdnlogsprodsoa"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-VerizonCDN"
					],
					"resourceNames": [
					  "azurecdnlogstsprod"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Persistent-QA"
					],
					"resourceNames": [
					  "persistentqastorage"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Consumer-Dev-Microservices"
					],
					"resourceNames": [
					  "stconsumerspadev2018"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Consumer-Dev-Microservices"
					],
					"resourceNames": [
					  "stconsumerspadev2019"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Consumer-Dev-Microservices"
					],
					"resourceNames": [
					  "stconsumerspadev2020"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Consumer-Dev-Microservices"
					],
					"resourceNames": [
					  "stconsumerspadev2021"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Pro-Dev-Microservices"
					],
					"resourceNames": [
					  "stprospadev2019"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Pro-Dev-Microservices"
					],
					"resourceNames": [
					  "stprospadev2020"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Pro-Dev-Microservices"
					],
					"resourceNames": [
					  "stprospadev2021"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Pro-Microservices"
					],
					"resourceNames": [
					  "stprospaprod2020"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Pro-Microservices"
					],
					"resourceNames": [
					  "stprospaprod2021"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_3_7": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_3_8": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceNames": [
					  "isstoragecentralus"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_3_9": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_1_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_1_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_1_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_2_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_2_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_2_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_2_4": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_2_5": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_3_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_3_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_3_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_3_4": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_3_5": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_3_6": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_3_7": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_3_8": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_4": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_4_5": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_1_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_1_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_1_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_1_4": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceNames": [
					  "default"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_5_1_5": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceNames": [
					  "Consumer-KV-Dev"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "Consumer-KV-PP"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "Consumer-KV-Prod"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "Consumer-KV-QA"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "IS-resourcegroup-KV"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "IT-KV"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "it-terraform-mgmt-kv"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "it-terraform-wireguard"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "Persistent-KV-QA"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "Pro-KV-Dev"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "Pro-KV-PP"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "Pro-KV-Prod"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "Pro-KV-QA"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "VITA-KV-Dev"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "VITA-KV-Prod"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_5_2_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_2_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_2_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_2_4": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_2_5": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_2_6": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_2_7": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_2_8": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_2_9": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_5_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_6_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_6_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_6_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_6_4": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_6_5": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_6_6": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_7_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_7_2": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceNames": [
					  "IT-AZ-WAF-*"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "PR-AZ-CV-*"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "EX-AZ-WEB-01"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "IT-AZ-STWAF-01"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "IS-AZ-ES-01"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "CO-AZ-BIC-01"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "IT-AZ-WUST-01"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "IT-AZ-NUST-01"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "IT-AZ-EUST-01"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_7_3": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceNames": [
					  "PR-AZ-CV-*"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceNames": [
					  "IT-AZ-WAF-*"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_7_4": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_7_5": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_7_6": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_7_7": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_8_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_8_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_8_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_8_4": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_8_5": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_9_1": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceGroupNames": [
					  "RG-TSCDN-Legacy-West"
					],
					"resourceNames": [
					  "tscdnwest"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-TSPro-Resources"
					],
					"resourceNames": [
					  "tspro-resources"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-TSCDN-Legacy-East"
					],
					"resourceNames": [
					  "tscdneast"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-IT-Resources"
					],
					"resourceNames": [
					  "updatensg"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Persistent-QA"
					],
					"resourceNames": [
					  "TS-PERSISTENT-QA"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Consumer-GTM"
					],
					"resourceNames": [
					  "ConsumerGTM-Server"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Consumer-GTM"
					],
					"resourceNames": [
					  "ConsumerGTM-Preview"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Pro-GTM"
					],
					"resourceNames": [
					  "ProGTM-Server"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Pro-GTM"
					],
					"resourceNames": [
					  "ProGTM-Preview"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_9_10": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceNames": [
					  "tspro-resources"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_9_11": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_9_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_9_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_9_4": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_9_5": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"resourceGroupNames": [
					  "RG-TSCDN-Legacy-West"
					],
					"resourceNames": [
					  "tscdnwest"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-TSPro-Resources"
					],
					"resourceNames": [
					  "tspro-resources"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-TSCDN-Legacy-East"
					],
					"resourceNames": [
					  "tscdneast"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-IT-Resources"
					],
					"resourceNames": [
					  "updatensg"
					],
					"subscriptionIds": [
					  "ALL_SUBSCRIPTIONS"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Persistent-QA"
					],
					"resourceNames": [
					  "TS-PERSISTENT-QA"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Consumer-GTM"
					],
					"resourceNames": [
					  "ConsumerGTM-Server"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Consumer-GTM"
					],
					"resourceNames": [
					  "ConsumerGTM-Preview"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Pro-GTM"
					],
					"resourceNames": [
					  "ProGTM-Server"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  },
				  {
					"resourceGroupNames": [
					  "RG-Pro-GTM"
					],
					"resourceNames": [
					  "ProGTM-Preview"
					],
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  },
			  "Azure_CIS_131_9_6": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_9_7": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_9_8": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_CIS_131_9_9": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_LW_IAM_1": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_LW_IAM_2": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_LW_IAM_3": {
				"enabled": true,
				"suppressionConditions": null
			  },
			  "Azure_LW_Monitoring_1": {
				"enabled": true,
				"suppressionConditions": [
				  {
					"subscriptionIds": [
					  "F0EAD1FB-F3B2-414D-8B85-A7DD546E8B17"
					],
					"tenantIds": [
					  "2dc69930-6a3f-4ac6-85b2-ba68071fe47c"
					]
				  }
				]
			  }
			}
	}],
	"message": "SUCCESS",
	"ok": true
}`
}

func genAzurePoliciesExceptionConstraintsMap() map[string][]string {
	return map[string][]string{
		"dev7-default-14":    nil,
		"dev7-default-15":    nil,
		"dev7-ds":            nil,
		"dev7-lwcustom-1":    nil,
		"dev7-lwcustom-11":   nil,
		"dev7-lwcustom-12":   nil,
		"dev7-lwcustom-13":   nil,
		"dev7-lwcustom-14":   nil,
		"dev7-lwcustom-15":   nil,
		"dev7-lwcustom-16":   nil,
		"dev7-lwcustom-17":   nil,
		"dev7-lwcustom-18":   nil,
		"dev7-lwcustom-19":   nil,
		"dev7-lwcustom-2":    nil,
		"dev7-lwcustom-20":   nil,
		"dev7-lwcustom-200":  nil,
		"dev7-lwcustom-21":   nil,
		"dev7-lwcustom-22":   nil,
		"dev7-lwcustom-23":   nil,
		"dev7-lwcustom-24":   nil,
		"dev7-lwcustom-25":   nil,
		"dev7-lwcustom-26":   nil,
		"dev7-lwcustom-27":   nil,
		"dev7-lwcustom-28":   nil,
		"dev7-lwcustom-29":   nil,
		"dev7-lwcustom-3":    nil,
		"dev7-lwcustom-30":   nil,
		"dev7-lwcustom-31":   nil,
		"dev7-lwcustom-32":   nil,
		"dev7-lwcustom-33":   nil,
		"dev7-lwcustom-34":   nil,
		"dev7-lwcustom-35":   nil,
		"dev7-lwcustom-37":   nil,
		"dev7-lwcustom-38":   nil,
		"dev7-lwcustom-39":   nil,
		"dev7-lwcustom-4":    nil,
		"dev7-lwcustom-42":   nil,
		"dev7-lwcustom-44":   nil,
		"dev7-lwcustom-6":    nil,
		"dev7-lwcustom-9":    nil,
		"dev7-rr-1":          nil,
		"dev7-test":          nil,
		"lacework-global-1":  nil,
		"lacework-global-10": nil,
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
		"lacework-global-11": nil,
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
		"lacework-global-12": nil,
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
		"lacework-global-13": nil,
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
		"lacework-global-14": nil,
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
		"lacework-global-15": nil,
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
		"lacework-global-158": nil,
		"lacework-global-159": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-16": nil,
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
		"lacework-global-162": nil,
		"lacework-global-163": nil,
		"lacework-global-164": nil,
		"lacework-global-165": nil,
		"lacework-global-166": nil,
		"lacework-global-167": nil,
		"lacework-global-168": nil,
		"lacework-global-169": nil,
		"lacework-global-17":  nil,
		"lacework-global-170": nil,
		"lacework-global-171": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-172": nil,
		"lacework-global-173": nil,
		"lacework-global-174": nil,
		"lacework-global-175": nil,
		"lacework-global-176": nil,
		"lacework-global-177": nil,
		"lacework-global-178": nil,
		"lacework-global-179": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-18": nil,
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
		"lacework-global-185": nil,
		"lacework-global-186": nil,
		"lacework-global-187": nil,
		"lacework-global-188": nil,
		"lacework-global-189": nil,
		"lacework-global-19":  nil,
		"lacework-global-190": nil,
		"lacework-global-191": nil,
		"lacework-global-192": nil,
		"lacework-global-193": nil,
		"lacework-global-194": nil,
		"lacework-global-195": nil,
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
		"lacework-global-2":   nil,
		"lacework-global-20":  nil,
		"lacework-global-200": nil,
		"lacework-global-201": nil,
		"lacework-global-202": nil,
		"lacework-global-203": nil,
		"lacework-global-204": nil,
		"lacework-global-205": nil,
		"lacework-global-206": nil,
		"lacework-global-21":  nil,
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
		"lacework-global-22": nil,
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
		"lacework-global-224": {
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
		"lacework-global-23": nil,
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
		"lacework-global-232": nil,
		"lacework-global-233": nil,
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
		"lacework-global-315": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-316": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-317": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-318": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-319": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-32": nil,
		"lacework-global-320": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-321": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-322": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-323": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-324": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-325": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-326": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-327": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-328": nil,
		"lacework-global-329": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-33": nil,
		"lacework-global-330": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sRegionNames",
			"k8sResourceTypes",
		},
		"lacework-global-331": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-332": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-333": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-334": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-335": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-336": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-337": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-338": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-339": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-34": {
			"accountIds",
		},
		"lacework-global-340": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-341": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-342": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-343": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-344": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-345": nil,
		"lacework-global-346": nil,
		"lacework-global-347": nil,
		"lacework-global-348": nil,
		"lacework-global-349": nil,
		"lacework-global-35": {
			"accountIds",
		},
		"lacework-global-350": nil,
		"lacework-global-351": nil,
		"lacework-global-352": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-353": nil,
		"lacework-global-354": nil,
		"lacework-global-355": nil,
		"lacework-global-356": nil,
		"lacework-global-357": nil,
		"lacework-global-358": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-359": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-36": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-360": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-361": nil,
		"lacework-global-362": nil,
		"lacework-global-363": nil,
		"lacework-global-364": nil,
		"lacework-global-365": nil,
		"lacework-global-366": nil,
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
		"lacework-global-499": nil,
		"lacework-global-5":   nil,
		"lacework-global-50": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-500": nil,
		"lacework-global-501": nil,
		"lacework-global-502": nil,
		"lacework-global-503": nil,
		"lacework-global-504": nil,
		"lacework-global-505": nil,
		"lacework-global-506": nil,
		"lacework-global-507": nil,
		"lacework-global-508": nil,
		"lacework-global-509": nil,
		"lacework-global-51": {
			"accountIds",
			"regionNames",
		},
		"lacework-global-510": nil,
		"lacework-global-511": nil,
		"lacework-global-512": {
			"resourceName",
			"subscriptions",
			"tenants",
		},
		"lacework-global-513": nil,
		"lacework-global-514": nil,
		"lacework-global-515": nil,
		"lacework-global-516": nil,
		"lacework-global-517": nil,
		"lacework-global-518": nil,
		"lacework-global-519": nil,
		"lacework-global-52": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-520": nil,
		"lacework-global-521": nil,
		"lacework-global-522": nil,
		"lacework-global-523": nil,
		"lacework-global-525": nil,
		"lacework-global-526": nil,
		"lacework-global-527": nil,
		"lacework-global-528": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-529": nil,
		"lacework-global-53": {
			"accountIds",
		},
		"lacework-global-530": nil,
		"lacework-global-531": nil,
		"lacework-global-532": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-533": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-534": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-535": nil,
		"lacework-global-536": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-537": nil,
		"lacework-global-54": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-540": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-541": nil,
		"lacework-global-543": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-549": nil,
		"lacework-global-55": {
			"accountIds",
			"regionNames",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-551": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-553": nil,
		"lacework-global-554": nil,
		"lacework-global-555": {
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-556": nil,
		"lacework-global-557": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-558": {
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-559": {
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-56": {
			"accountIds",
			"resourceNames",
			"resourceTags",
		},
		"lacework-global-560": {
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-561": {
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-562": {
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-563": {
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-564": {
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-565": {
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-566": {
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-567": {
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-568": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-569": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-57": {
			"accountIds",
		},
		"lacework-global-570": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-571": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-572": nil,
		"lacework-global-573": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-574": nil,
		"lacework-global-579": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-58": {
			"accountIds",
		},
		"lacework-global-580": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-581": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-582": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-583": nil,
		"lacework-global-584": nil,
		"lacework-global-585": nil,
		"lacework-global-586": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-587": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-588": nil,
		"lacework-global-589": nil,
		"lacework-global-59": {
			"accountIds",
		},
		"lacework-global-590": nil,
		"lacework-global-591": nil,
		"lacework-global-592": nil,
		"lacework-global-593": nil,
		"lacework-global-594": nil,
		"lacework-global-595": nil,
		"lacework-global-596": nil,
		"lacework-global-597": nil,
		"lacework-global-598": nil,
		"lacework-global-599": nil,
		"lacework-global-6":   nil,
		"lacework-global-60": {
			"accountIds",
		},
		"lacework-global-600": nil,
		"lacework-global-601": nil,
		"lacework-global-602": nil,
		"lacework-global-603": nil,
		"lacework-global-604": nil,
		"lacework-global-605": nil,
		"lacework-global-606": nil,
		"lacework-global-607": nil,
		"lacework-global-608": nil,
		"lacework-global-609": nil,
		"lacework-global-61": {
			"accountIds",
		},
		"lacework-global-610": nil,
		"lacework-global-613": nil,
		"lacework-global-614": nil,
		"lacework-global-615": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-616": nil,
		"lacework-global-617": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-618": nil,
		"lacework-global-619": nil,
		"lacework-global-62": {
			"accountIds",
		},
		"lacework-global-620": nil,
		"lacework-global-621": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-626": nil,
		"lacework-global-627": nil,
		"lacework-global-628": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-629": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-63": {
			"accountIds",
		},
		"lacework-global-630": nil,
		"lacework-global-631": nil,
		"lacework-global-632": nil,
		"lacework-global-633": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-634": {
			"resourceName",
			"subscriptions",
			"tenants",
		},
		"lacework-global-635": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-636": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-637": nil,
		"lacework-global-638": nil,
		"lacework-global-639": nil,
		"lacework-global-64": {
			"accountIds",
		},
		"lacework-global-640": nil,
		"lacework-global-641": nil,
		"lacework-global-642": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-643": {
			"azureResourceGroup",
			"regionNames",
			"resourceName",
			"resourceTags",
			"subscriptions",
			"tenants",
		},
		"lacework-global-644": nil,
		"lacework-global-645": nil,
		"lacework-global-646": nil,
		"lacework-global-647": nil,
		"lacework-global-648": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-649": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-65": {
			"accountIds",
		},
		"lacework-global-650": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-651": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-652": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-653": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-654": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-655": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-656": nil,
		"lacework-global-657": nil,
		"lacework-global-658": nil,
		"lacework-global-659": nil,
		"lacework-global-66": {
			"accountIds",
		},
		"lacework-global-660": nil,
		"lacework-global-662": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-663": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-664": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-665": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-666": {
			"k8sAccountIds",
			"k8sClusterNames",
			"k8sNamespaces",
			"k8sRegionNames",
			"k8sResourceLabels",
			"k8sResourceTypes",
			"resourceNames",
		},
		"lacework-global-667": nil,
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
