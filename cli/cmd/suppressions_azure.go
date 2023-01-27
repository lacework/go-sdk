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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/fatih/color"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	azureEquivalencesMap = map[string]string{
		"Azure_CIS_131_1_1":  "lacework-global-514",
		"Azure_CIS_131_1_2":  "lacework-global-597",
		"Azure_CIS_131_1_3":  "lacework-global-499",
		"Azure_CIS_131_1_4":  "lacework-global-500",
		"Azure_CIS_131_1_5":  "lacework-global-501",
		"Azure_CIS_131_1_6":  "lacework-global-503",
		"Azure_CIS_131_1_7":  "lacework-global-504",
		"Azure_CIS_131_1_8":  "lacework-global-505",
		"Azure_CIS_131_1_9":  "lacework-global-506",
		"Azure_CIS_131_1_10": "lacework-global-507",
		"Azure_CIS_131_1_11": "lacework-global-508",
		"Azure_CIS_131_1_12": "lacework-global-509",
		//"Azure_CIS_131_1_13": "N/A",
		"Azure_CIS_131_1_14": "lacework-global-590",
		"Azure_CIS_131_1_15": "lacework-global-510",
		"Azure_CIS_131_1_16": "lacework-global-591",
		"Azure_CIS_131_1_17": "lacework-global-592",
		"Azure_CIS_131_1_18": "lacework-global-593",
		"Azure_CIS_131_1_19": "lacework-global-594",
		"Azure_CIS_131_1_20": "lacework-global-511",
		"Azure_CIS_131_1_21": "lacework-global-512",
		"Azure_CIS_131_1_22": "lacework-global-513",
		"Azure_CIS_131_1_23": "lacework-global-595",
		"Azure_CIS_131_2_1":  "lacework-global-598",
		"Azure_CIS_131_2_2":  "lacework-global-599",
		"Azure_CIS_131_2_3":  "lacework-global-601",
		"Azure_CIS_131_2_4":  "lacework-global-602",
		"Azure_CIS_131_2_5":  "lacework-global-604",
		//"Azure_CIS_131_2_6": "N/A",
		"Azure_CIS_131_2_7":   "lacework-global-605",
		"Azure_CIS_131_2_8":   "lacework-global-607",
		"Azure_CIS_131_2_9":   "lacework-global-614",
		"Azure_CIS_131_2_10":  "lacework-global-613",
		"Azure_CIS_131_2_11":  "lacework-global-524",
		"Azure_CIS_131_2_12":  "lacework-global-523",
		"Azure_CIS_131_2_13":  "lacework-global-526",
		"Azure_CIS_131_2_14":  "lacework-global-527",
		"Azure_CIS_131_2_15":  "lacework-global-525",
		"Azure_CIS_131_3_1":   "lacework-global-528",
		"Azure_CIS_131_3_2":   "lacework-global-530",
		"Azure_CIS_131_3_3":   "lacework-global-616",
		"Azure_CIS_131_3_4":   "lacework-global-531",
		"Azure_CIS_131_3_5":   "lacework-global-532",
		"Azure_CIS_131_3_6":   "lacework-global-533",
		"Azure_CIS_131_3_7":   "lacework-global-617",
		"Azure_CIS_131_3_8":   "lacework-global-535",
		"Azure_CIS_131_3_9":   "lacework-global-618",
		"Azure_CIS_131_3_10":  "lacework-global-619",
		"Azure_CIS_131_3_11":  "lacework-global-620",
		"Azure_CIS_131_4_1_1": "lacework-global-537",
		"Azure_CIS_131_4_1_2": "lacework-global-540",
		"Azure_CIS_131_4_1_3": "lacework-global-541",
		"Azure_CIS_131_4_2_1": "lacework-global-622",
		"Azure_CIS_131_4_2_2": "lacework-global-623",
		"Azure_CIS_131_4_2_3": "lacework-global-624",
		"Azure_CIS_131_4_2_4": "lacework-global-625",
		"Azure_CIS_131_4_2_5": "lacework-global-542",
		"Azure_CIS_131_4_3_1": "lacework-global-543",
		"Azure_CIS_131_4_3_2": "lacework-global-551",
		"Azure_CIS_131_4_3_3": "lacework-global-544",
		"Azure_CIS_131_4_3_4": "lacework-global-545",
		"Azure_CIS_131_4_3_5": "lacework-global-546",
		"Azure_CIS_131_4_3_6": "lacework-global-547",
		"Azure_CIS_131_4_3_7": "lacework-global-548",
		"Azure_CIS_131_4_3_8": "lacework-global-549",
		"Azure_CIS_131_4_4":   "lacework-global-539",
		"Azure_CIS_131_4_5":   "lacework-global-621",
		"Azure_CIS_131_5_1_1": "lacework-global-554",
		"Azure_CIS_131_5_1_2": "lacework-global-555",
		"Azure_CIS_131_5_1_3": "lacework-global-556",
		"Azure_CIS_131_5_1_4": "lacework-global-630",
		"Azure_CIS_131_5_1_5": "lacework-global-557",
		"Azure_CIS_131_5_2_1": "lacework-global-558",
		"Azure_CIS_131_5_2_2": "lacework-global-559",
		"Azure_CIS_131_5_2_3": "lacework-global-560",
		"Azure_CIS_131_5_2_4": "lacework-global-561",
		//"Azure_CIS_131_5_2_5": "N/A",
		//"Azure_CIS_131_5_2_6": "N/A",
		"Azure_CIS_131_5_2_7": "lacework-global-562",
		"Azure_CIS_131_5_2_8": "lacework-global-563",
		"Azure_CIS_131_5_2_9": "lacework-global-564",
		"Azure_CIS_131_5_3":   "lacework-global-553",
		"Azure_CIS_131_6_1":   "lacework-global-568",
		"Azure_CIS_131_6_2":   "lacework-global-569",
		"Azure_CIS_131_6_3":   "lacework-global-538",
		"Azure_CIS_131_6_4":   "lacework-global-633",
		"Azure_CIS_131_6_5":   "lacework-global-634",
		"Azure_CIS_131_6_6":   "lacework-global-570",
		"Azure_CIS_131_7_1":   "lacework-global-573",
		"Azure_CIS_131_7_2":   "lacework-global-635",
		"Azure_CIS_131_7_3":   "lacework-global-636",
		"Azure_CIS_131_7_4":   "lacework-global-574",
		"Azure_CIS_131_7_5":   "lacework-global-522",
		"Azure_CIS_131_7_6":   "lacework-global-637",
		"Azure_CIS_131_7_7":   "lacework-global-638",
		"Azure_CIS_131_8_1":   "lacework-global-575",
		"Azure_CIS_131_8_2":   "lacework-global-577",
		"Azure_CIS_131_8_3":   "lacework-global-645",
		"Azure_CIS_131_8_4":   "lacework-global-579",
		//"Azure_CIS_131_8_5": "N/A",
		"Azure_CIS_131_9_1":  "lacework-global-642",
		"Azure_CIS_131_9_2":  "lacework-global-580",
		"Azure_CIS_131_9_3":  "lacework-global-581",
		"Azure_CIS_131_9_4":  "lacework-global-643",
		"Azure_CIS_131_9_5":  "lacework-global-582",
		"Azure_CIS_131_9_6":  "lacework-global-583",
		"Azure_CIS_131_9_7":  "lacework-global-584",
		"Azure_CIS_131_9_8":  "lacework-global-585",
		"Azure_CIS_131_9_9":  "lacework-global-586",
		"Azure_CIS_131_9_10": "lacework-global-587",
		"Azure_CIS_131_9_11": "lacework-global-644",
	}

	// suppressionsMigrateAzureCmd represents the azure sub-command inside the suppressions migrate
	//command
	suppressionsMigrateAzureCmd = &cobra.Command{
		Use:     "migrate",
		Aliases: []string{"mig"},
		Short:   "Migrate legacy suppressions for Azure to mapped policy exceptions",
		RunE:    suppressionsAzureMigrate,
	}

	// suppressionsListAzureCmd represents the azure sub-command inside the suppressions list
	//command
	suppressionsListAzureCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List legacy suppressions for Azure",
		RunE:    suppressionsAzureList,
	}
)

func suppressionsAzureList(_ *cobra.Command, _ []string) error {
	var (
		suppressions map[string]api.SuppressionV2
		err          error
	)

	suppressions, err = cli.LwApi.V2.Suppressions.Azure.List()
	if err != nil {
		if strings.Contains(err.Error(), "No active Azure accounts") {
			cli.OutputHuman("No active Azure accounts found. " +
				"Unable to get legacy Azure suppressions\n")
			return nil
		}
		return errors.Wrap(err, "Unable to get legacy Azure suppressions")
	}

	if len(suppressions) == 0 {
		cli.OutputHuman("No legacy Azure suppressions found.\n")
		return nil
	}
	return cli.OutputJSON(suppressions)
}

func suppressionsAzureMigrate(_ *cobra.Command, _ []string) error {
	var (
		suppressionsMap map[string]api.SuppressionV2
		err             error

		convertedPolicyExceptions []map[string]api.PolicyException
		payloadsText              []string
		discardedSuppressions     []map[string]api.SuppressionV2
	)
	suppressionsMap, err = cli.LwApi.V2.Suppressions.Azure.List()
	if err != nil {
		if strings.Contains(err.Error(), "No active Azure accounts") {
			cli.OutputHuman("No active Azure accounts found. " +
				"Unable to get legacy Azure suppressions\n")
			return nil
		}
		return errors.Wrap(err, "Unable to get legacy Azure suppressions")
	}

	if len(suppressionsMap) == 0 {
		cli.OutputHuman("No legacy Azure suppressions found.\n")
		return nil
	}

	answer := ""
	manualMigration := "Output translated legacy suppressions as policy exception commands to be" +
		" run manually (Recommended)"
	autoMigration := "Auto migrate legacy suppressions.\nDISCLAIMER: " +
		"By selecting this option, you accept liability for the migration and " +
		"any compliance violations missed as a result of the added exceptions"
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Select{
			Message: "Select your legacy suppression migration approach?",
			Options: []string{
				manualMigration,
				autoMigration,
			},
		},
		Response: &answer,
	}); err != nil {
		return err
	}

	// get a list of all policies and parse the valid exception constraints and create a map of
	// {"<policyId>": [<validPolicyConstraints>]}
	policyExceptionsConstraintsMap := getPoliciesExceptionConstraintsMap()

	switch answer {
	case manualMigration:
		_, payloadsText, discardedSuppressions = convertAzureSuppressions(
			suppressionsMap,
			policyExceptionsConstraintsMap,
		)
		printPayloadsText(payloadsText)
		printDiscardedSuppressions(discardedSuppressions)
	case autoMigration:
		convertedPolicyExceptions, _, discardedSuppressions = convertAzureSuppressions(
			suppressionsMap,
			policyExceptionsConstraintsMap,
		)
		printConvertedSuppressions(convertedPolicyExceptions)
		confirm := false
		err := survey.AskOne(&survey.Confirm{
			Message: "Confirm the above exceptions have been reviewed and you wish to continue" +
				" with the auto migration.",
		}, &confirm)
		if err != nil {
			return err
		}
		if confirm {
			autoConvertSuppressions(convertedPolicyExceptions)
			printDiscardedSuppressions(discardedSuppressions)
			cli.OutputHuman(color.GreenString("To view the newly created Exceptions, " +
				"try running `lacework policy-exceptions list <policyId>"))
		} else {
			cli.OutputHuman("Cancelled Legacy Suppression to Exception migration!")
		}
	}

	return nil
}

func convertAzureSuppressions(
	suppressionsMap map[string]api.SuppressionV2,
	policyExceptionsConstraintsMap map[string][]string,
) ([]map[string]api.PolicyException,
	[]string, []map[string]api.SuppressionV2) {
	var (
		convertedPolicyExceptions []map[string]api.PolicyException
		payloadsText              []string
		discardedSuppressions     []map[string]api.SuppressionV2
	)

	for id, suppressionInfo := range suppressionsMap {
		// verify there is a mapped policy for this recommendation
		// if the recommendation is not a key in the map we can assume this is not mapped and
		// continue
		mappedPolicyId, ok := azureEquivalencesMap[id]
		if !ok {
			// when we don't have a mapped policy, add the legacy suppression info
			if suppressionInfo.SuppressionConditions != nil {
				suppressionInfo = updateDiscardedSupConditionsComments(suppressionInfo,
					"Legacy suppression discarded as there is no equivalent policy")
				discardedSuppressions = append(
					discardedSuppressions,
					map[string]api.SuppressionV2{id: suppressionInfo},
				)
			}
			continue
		}

		// get the supported policy exception fields for the mapped policy
		// in order to ensure we have an up-to-date list of exception constraints we need to
		// get the policy from the /api/v2/Policies/<policyId> api.
		// We then parse this into a list of constraints
		policyIdExceptionsTemplate := policyExceptionsConstraintsMap[mappedPolicyId]
		if policyIdExceptionsTemplate == nil {
			// Updating the suppression conditions comments to make it clear why these were
			// discarded
			if len(suppressionInfo.SuppressionConditions) >= 1 {
				suppressionInfo = updateDiscardedSupConditionsComments(suppressionInfo,
					fmt.Sprintf("Legacy suppression discarded as the new policy: %s does not"+
						" support exception conditions", mappedPolicyId))

				// if the list of supported constraints is empty for a policy,
				// we should let the customers know that we have discarded this suppression
				discardedSuppressions = append(
					discardedSuppressions,
					map[string]api.SuppressionV2{id: suppressionInfo},
				)
			}
			continue
		}
		if len(suppressionInfo.SuppressionConditions) >= 1 {
			for _, suppression := range suppressionInfo.SuppressionConditions {
				// used to store the converted legacy suppressions
				var convertedConstraints []api.PolicyExceptionConstraint

				resourceGroupNamesConstraint := convertSupCondition(suppression.ResourceGroupNames,
					"azureResourceGroup",
					policyIdExceptionsTemplate)
				if resourceGroupNamesConstraint.FieldKey != "" {
					convertedConstraints = append(convertedConstraints, resourceGroupNamesConstraint)
				}

				regionNamesConstraint := convertSupCondition(suppression.RegionNames,
					"regionNames",
					policyIdExceptionsTemplate)
				if regionNamesConstraint.FieldKey != "" {
					convertedConstraints = append(convertedConstraints, regionNamesConstraint)
				}

				tenantsConstraint := convertSupCondition(suppression.TenantIds,
					"tenants",
					policyIdExceptionsTemplate)
				if tenantsConstraint.FieldKey != "" {
					convertedConstraints = append(convertedConstraints, tenantsConstraint)
				}

				subscriptionsConstraint := convertSupCondition(suppression.SubscriptionIds,
					"subscriptions",
					policyIdExceptionsTemplate)
				if subscriptionsConstraint.FieldKey != "" {
					convertedConstraints = append(convertedConstraints, subscriptionsConstraint)
				}

				resourceNamesConstraint := convertSupCondition(suppression.ResourceNames,
					"resourceName",
					policyIdExceptionsTemplate)
				if resourceNamesConstraint.FieldKey != "" {
					convertedConstraints = append(convertedConstraints, resourceNamesConstraint)
				}

				resourceTagsConstraint := convertSupConditionTags(suppression.ResourceTags,
					"resourceTags",
					policyIdExceptionsTemplate)
				if resourceTagsConstraint.FieldKey != "" {
					convertedConstraints = append(convertedConstraints, resourceTagsConstraint)
				}

				description := fmt.Sprintf(
					"Migrated exception from legacy compliance policy %s. ", id,
				)
				if suppression.Comment != "" {
					description = description + fmt.Sprintf(
						"Legacy Policy comment: %s", suppression.Comment,
					)
				}
				if len(convertedConstraints) >= 1 {
					convertedPolicyExceptions = append(
						convertedPolicyExceptions,
						map[string]api.PolicyException{
							mappedPolicyId: {
								Description: description,
								Constraints: convertedConstraints,
							},
						},
					)

					exception := api.PolicyException{
						Description: description,
						Constraints: convertedConstraints,
					}
					jsonException, err := json.Marshal(exception)
					if err != nil {
						cli.Log.Error(err)
					}

					payloadsText = append(
						payloadsText,
						fmt.Sprintf(
							"lacework api post '/Exceptions?policyId=%s' -d '%s'",
							mappedPolicyId,
							jsonException,
						),
					)
				}
			}
		}
	}

	return convertedPolicyExceptions, payloadsText, discardedSuppressions
}
