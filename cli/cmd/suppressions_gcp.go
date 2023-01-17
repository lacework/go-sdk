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

	"github.com/fatih/color"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	gcpEquivalencesMap = map[string]string{
		"GCP_CIS12_1_2":   "lacework-global-233",
		"GCP_CIS12_1_3":   "lacework-global-293",
		"GCP_CIS12_1_4":   "lacework-global-234",
		"GCP_CIS12_1_5":   "lacework-global-235",
		"GCP_CIS12_1_6":   "lacework-global-236",
		"GCP_CIS12_1_7":   "lacework-global-237",
		"GCP_CIS12_1_8":   "lacework-global-294",
		"GCP_CIS12_1_9":   "lacework-global-238",
		"GCP_CIS12_1_10":  "lacework-global-239",
		"GCP_CIS12_1_11":  "lacework-global-295",
		"GCP_CIS12_1_12":  "lacework-global-296",
		"GCP_CIS12_1_13":  "lacework-global-240",
		"GCP_CIS12_1_14":  "lacework-global-241",
		"GCP_CIS12_1_15":  "lacework-global-242",
		"GCP_CIS12_2_1":   "lacework-global-245",
		"GCP_CIS12_2_2":   "lacework-global-246",
		"GCP_CIS12_2_3":   "lacework-global-298",
		"GCP_CIS12_2_4":   "lacework-global-247",
		"GCP_CIS12_2_5":   "lacework-global-248",
		"GCP_CIS12_2_6":   "lacework-global-249",
		"GCP_CIS12_2_7":   "lacework-global-250",
		"GCP_CIS12_2_8":   "lacework-global-251",
		"GCP_CIS12_2_9":   "lacework-global-252",
		"GCP_CIS12_2_10":  "lacework-global-253",
		"GCP_CIS12_2_11":  "lacework-global-254",
		"GCP_CIS12_2_12":  "lacework-global-255",
		"GCP_CIS12_3_1":   "lacework-global-300",
		"GCP_CIS12_3_2":   "lacework-global-258",
		"GCP_CIS12_3_3":   "lacework-global-259",
		"GCP_CIS12_3_4":   "lacework-global-260",
		"GCP_CIS12_3_5":   "lacework-global-261",
		"GCP_CIS12_3_6":   "lacework-global-301",
		"GCP_CIS12_3_7":   "lacework-global-302",
		"GCP_CIS12_3_8":   "lacework-global-262",
		"GCP_CIS12_3_9":   "lacework-global-263",
		"GCP_CIS12_3_10":  "lacework-global-303",
		"GCP_CIS12_4_1":   "lacework-global-264",
		"GCP_CIS12_4_2":   "lacework-global-265",
		"GCP_CIS12_4_3":   "lacework-global-266",
		"GCP_CIS12_4_4":   "lacework-global-267",
		"GCP_CIS12_4_5":   "lacework-global-268",
		"GCP_CIS12_4_6":   "lacework-global-269",
		"GCP_CIS12_4_7":   "lacework-global-304",
		"GCP_CIS12_4_8":   "lacework-global-305",
		"GCP_CIS12_4_9":   "lacework-global-306",
		"GCP_CIS12_4_10":  "lacework-global-307",
		"GCP_CIS12_4_11":  "lacework-global-308",
		"GCP_CIS12_5_1":   "lacework-global-270",
		"GCP_CIS12_5_2":   "lacework-global-310",
		"GCP_CIS12_6_1_1": "lacework-global-274",
		"GCP_CIS12_6_1_2": "lacework-global-275",
		"GCP_CIS12_6_1_3": "lacework-global-276",
		//"GCP_CIS12_6_2_1": "N/A",
		"GCP_CIS12_6_2_2": "lacework-global-312",
		"GCP_CIS12_6_2_3": "lacework-global-277",
		"GCP_CIS12_6_2_4": "lacework-global-278",
		//"GCP_CIS12_6_2_5": "N/A",
		//"GCP_CIS12_6_2_6": "N/A",
		"GCP_CIS12_6_2_7": "lacework-global-279",
		"GCP_CIS12_6_2_8": "lacework-global-280",
		//"GCP_CIS12_6_2_9": "N/A",
		//"GCP_CIS12_6_2_10": "N/A",
		//"GCP_CIS12_6_2_11": "N/A",
		//"GCP_CIS12_6_2_12": "N/A",
		"GCP_CIS12_6_2_13": "lacework-global-281",
		"GCP_CIS12_6_2_14": "lacework-global-282",
		//"GCP_CIS12_6_2_15": "N/A",
		"GCP_CIS12_6_2_16": "lacework-global-283",
		"GCP_CIS12_6_3_1":  "lacework-global-285",
		"GCP_CIS12_6_3_2":  "lacework-global-286",
		"GCP_CIS12_6_3_3":  "lacework-global-287",
		"GCP_CIS12_6_3_4":  "lacework-global-288",
		"GCP_CIS12_6_3_5":  "lacework-global-289",
		"GCP_CIS12_6_3_6":  "lacework-global-290",
		"GCP_CIS12_6_3_7":  "lacework-global-291",
		"GCP_CIS12_6_4":    "lacework-global-271",
		"GCP_CIS12_6_5":    "lacework-global-272",
		"GCP_CIS12_6_6":    "lacework-global-311",
		"GCP_CIS12_6_7":    "lacework-global-273",
		"GCP_CIS12_7_1":    "lacework-global-292",
		"GCP_CIS12_7_2":    "lacework-global-313",
		"GCP_CIS12_7_3":    "lacework-global-314",
	}

	// suppressionsMigrateGcpCmd represents the azure sub-command inside the suppressions migrate command
	suppressionsMigrateGcpCmd = &cobra.Command{
		Use:     "migrate",
		Aliases: []string{"mig"},
		Short:   "Migrate legacy suppressions for Gcp to mapped policy exceptions",
		RunE:    suppressionsGcpMigrate,
	}

	// suppressionsListGcpCmd represents the gcp sub-command inside the suppressions list command
	suppressionsListGcpCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List legacy suppressions for GCP",
		RunE:    suppressionsGcpList,
	}
)

func suppressionsGcpList(_ *cobra.Command, _ []string) error {
	var (
		suppressions map[string]api.SuppressionV2
		err          error
	)

	suppressions, err = cli.LwApi.V2.Suppressions.Gcp.List()
	if err != nil {
		if strings.Contains(err.Error(), "No active GCP accounts") {
			cli.OutputHuman("No active GCP accounts found. " +
				"Unable to get legacy GCP suppressions\n")
			return nil
		}
		return errors.Wrap(err, "Unable to get legacy GCP suppressions")
	}

	if len(suppressions) == 0 {
		cli.OutputHuman("No legacy GCP suppressions found.\n")
		return nil
	}
	return cli.OutputJSON(suppressions)
}

func suppressionsGcpMigrate(_ *cobra.Command, _ []string) error {
	var (
		suppressionsMap map[string]api.SuppressionV2
		err             error

		convertedPolicyExceptions []map[string]api.PolicyException
		payloadsText              []string
		discardedSuppressions     []map[string]api.SuppressionV2
	)
	suppressionsMap, err = cli.LwApi.V2.Suppressions.Gcp.List()
	if err != nil {
		if strings.Contains(err.Error(), "No active GCP accounts") {
			cli.OutputHuman("No active GCP accounts found. " +
				"Unable to get legacy gcp suppressions")
			return nil
		}
		return errors.Wrap(err, "Unable to get legacy gcp suppressions")
	}

	if len(suppressionsMap) == 0 {
		cli.OutputHuman("No legacy GCP suppressions found.\n")
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
		_, payloadsText, discardedSuppressions = convertGcpSuppressions(
			suppressionsMap,
			policyExceptionsConstraintsMap,
		)
		printPayloadsText(payloadsText)
		printDiscardedSuppressions(discardedSuppressions)
	case autoMigration:
		convertedPolicyExceptions, _, discardedSuppressions = convertGcpSuppressions(
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

func convertGcpSuppressions(
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
		mappedPolicyId, ok := gcpEquivalencesMap[id]
		if !ok {
			// when we don't have a mapped policy, add the legacy suppression info
			if suppressionInfo.SuppressionConditions != nil {
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
			continue
		}
		if len(suppressionInfo.SuppressionConditions) >= 1 {
			for _, suppression := range suppressionInfo.SuppressionConditions {
				// used to store the converted legacy suppressions
				var convertedConstraints []api.PolicyExceptionConstraint

				organizationIdsConstraint := convertSupCondition(suppression.OrganizationIds,
					"organizations",
					policyIdExceptionsTemplate)
				if organizationIdsConstraint.FieldKey != "" {
					convertedConstraints = append(convertedConstraints, organizationIdsConstraint)
				}

				projectIdsConstraint := convertSupCondition(suppression.ProjectIds,
					"projects",
					policyIdExceptionsTemplate)
				if projectIdsConstraint.FieldKey != "" {
					convertedConstraints = append(convertedConstraints, projectIdsConstraint)
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
