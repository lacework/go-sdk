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
	// https://docs.lacework.com/console/aws-compliance-policy-exceptions-criteria#lacework-custom-policies-for-aws-iam
	// https://docs.lacework.com/console/cis-aws-140-benchmark-report#identity-and-access-management
	// old ID to new ID mapping, using the old Constraints with the hope they match the new Constraints
	awsEquivalencesMap = map[string]string{
		"AWS_CIS_1_2":  "lacework-global-39",
		"AWS_CIS_1_3":  "lacework-global-41",
		"AWS_CIS_1_4":  "lacework-global-43",
		"AWS_CIS_1_9":  "lacework-global-37",
		"AWS_CIS_1_10": "lacework-global-38",
		"AWS_CIS_1_12": "lacework-global-34",
		"AWS_CIS_1_13": "lacework-global-35",
		"AWS_CIS_1_14": "lacework-global-69",
		"AWS_CIS_1_15": "lacework-global-33",
		"AWS_CIS_1_16": "lacework-global-44", // no iam policies to users
		//"AWS_CIS_1_19":  "lacework-global-31", manual
		//"AWS_CIS_1_20":  "lacework-global-32", manual
		//"AWS_CIS_1_21":  "lacework-global-70", manual
		"AWS_CIS_1_22":  "lacework-global-46",
		"AWS_CIS_1_23":  "lacework-global-40",
		"AWS_CIS_1_24":  "lacework-global-45",
		"AWS_CIS_2_1":   "lacework-global-53",
		"AWS_CIS_2_2":   "lacework-global-75",
		"AWS_CIS_2_3":   "lacework-global-54", // s3 bucket cloudtrail log
		"AWS_CIS_2_4":   "lacework-global-55",
		"AWS_CIS_2_5":   "lacework-global-76",
		"AWS_CIS_2_6":   "lacework-global-56", // s3 bucket cloudtrail log
		"AWS_CIS_2_7":   "lacework-global-77",
		"AWS_CIS_2_8":   "lacework-global-78",
		"AWS_CIS_2_9":   "lacework-global-79",
		"AWS_CIS_3_1":   "lacework-global-57",
		"AWS_CIS_3_2":   "lacework-global-58",
		"AWS_CIS_3_3":   "lacework-global-59",
		"AWS_CIS_3_4":   "lacework-global-60",
		"AWS_CIS_3_5":   "lacework-global-61",
		"AWS_CIS_3_6":   "lacework-global-82",
		"AWS_CIS_3_7":   "lacework-global-83",
		"AWS_CIS_3_8":   "lacework-global-62",
		"AWS_CIS_3_9":   "lacework-global-84",
		"AWS_CIS_3_10":  "lacework-global-85",
		"AWS_CIS_3_11":  "lacework-global-86",
		"AWS_CIS_3_12":  "lacework-global-63",
		"AWS_CIS_3_13":  "lacework-global-64",
		"AWS_CIS_3_14":  "lacework-global-65",
		"AWS_CIS_4_1":   "lacework-global-68",
		"AWS_CIS_4_2":   "lacework-global-68",
		"AWS_CIS_4_3":   "lacework-global-79",
		"AWS_CIS_4_4":   "lacework-global-87",
		"LW_S3_1":       "lacework-global-130",
		"LW_S3_2":       "lacework-global-131",
		"LW_S3_3":       "lacework-global-132",
		"LW_S3_4":       "lacework-global-133",
		"LW_S3_5":       "lacework-global-134",
		"LW_S3_6":       "lacework-global-135",
		"LW_S3_7":       "lacework-global-136",
		"LW_S3_8":       "lacework-global-137",
		"LW_S3_9":       "lacework-global-138",
		"LW_S3_10":      "lacework-global-139",
		"LW_S3_11":      "lacework-global-140",
		"LW_S3_12":      "lacework-global-94",
		"LW_S3_13":      "lacework-global-95",
		"LW_S3_14":      "lacework-global-217",
		"LW_S3_15":      "lacework-global-96",
		"LW_S3_16":      "lacework-global-97",
		"LW_S3_18":      "lacework-global-98",
		"LW_S3_19":      "lacework-global-99",
		"LW_S3_20":      "lacework-global-100",
		"LW_S3_21":      "lacework-global-101",
		"LW_AWS_IAM_1":  "lacework-global-115",
		"LW_AWS_IAM_2":  "lacework-global-116",
		"LW_AWS_IAM_3":  "lacework-global-117",
		"LW_AWS_IAM_4":  "lacework-global-118",
		"LW_AWS_IAM_5":  "lacework-global-119",
		"LW_AWS_IAM_6":  "lacework-global-120",
		"LW_AWS_IAM_7":  "lacework-global-121",
		"LW_AWS_IAM_11": "lacework-global-181", // non-root user
		"LW_AWS_IAM_12": "lacework-global-142",
		"LW_AWS_IAM_13": "lacework-global-141",
		"LW_AWS_IAM_14": "lacework-global-105",
		// "AWS_CIS_4_5" : "88 (Manual)",
		"LW_AWS_NETWORKING_1":       "lacework-global-227", // sec-group
		"LW_AWS_NETWORKING_2":       "lacework-global-145", // network acl
		"LW_AWS_NETWORKING_3":       "lacework-global-146", // network acl
		"LW_AWS_NETWORKING_4":       "lacework-global-147",
		"LW_AWS_NETWORKING_5":       "lacework-global-148",
		"LW_AWS_NETWORKING_6":       "lacework-global-149",
		"LW_AWS_NETWORKING_7":       "lacework-global-228",
		"LW_AWS_NETWORKING_8":       "lacework-global-229",
		"LW_AWS_NETWORKING_9":       "lacework-global-230",
		"LW_AWS_NETWORKING_10":      "lacework-global-231",
		"LW_AWS_NETWORKING_11":      "lacework-global-199",
		"LW_AWS_NETWORKING_12":      "lacework-global-150",
		"LW_AWS_NETWORKING_13":      "lacework-global-151",
		"LW_AWS_NETWORKING_14":      "lacework-global-152",
		"LW_AWS_NETWORKING_15":      "lacework-global-153",
		"LW_AWS_NETWORKING_16":      "lacework-global-225",
		"LW_AWS_NETWORKING_17":      "lacework-global-226",
		"LW_AWS_NETWORKING_18":      "lacework-global-154",
		"LW_AWS_NETWORKING_19":      "lacework-global-155",
		"LW_AWS_NETWORKING_20":      "lacework-global-156",
		"LW_AWS_NETWORKING_21":      "lacework-global-104",
		"LW_AWS_NETWORKING_22":      "lacework-global-106",
		"LW_AWS_NETWORKING_23":      "lacework-global-107",
		"LW_AWS_NETWORKING_24":      "lacework-global-108",
		"LW_AWS_NETWORKING_25":      "lacework-global-109",
		"LW_AWS_NETWORKING_26":      "lacework-global-110",
		"LW_AWS_NETWORKING_27":      "lacework-global-111",
		"LW_AWS_NETWORKING_28":      "lacework-global-112",
		"LW_AWS_NETWORKING_29":      "lacework-global-113",
		"LW_AWS_NETWORKING_30":      "lacework-global-114",
		"LW_AWS_NETWORKING_31":      "lacework-global-218",
		"LW_AWS_NETWORKING_32":      "lacework-global-219",
		"LW_AWS_NETWORKING_33":      "lacework-global-220",
		"LW_AWS_NETWORKING_34":      "lacework-global-221",
		"LW_AWS_NETWORKING_35":      "lacework-global-222",
		"LW_AWS_NETWORKING_36":      "lacework-global-148",
		"LW_AWS_NETWORKING_37":      "lacework-global-102",
		"LW_AWS_NETWORKING_38":      "lacework-global-223",
		"LW_AWS_NETWORKING_39":      "lacework-global-184",
		"LW_AWS_NETWORKING_40":      "lacework-global-103",
		"LW_AWS_NETWORKING_41":      "lacework-global-125", // cloudfront
		"LW_AWS_NETWORKING_42":      "lacework-global-126", // cloudfront
		"LW_AWS_NETWORKING_43":      "lacework-global-127",
		"LW_AWS_NETWORKING_44":      "lacework-global-231",
		"LW_AWS_NETWORKING_45":      "lacework-global-482",
		"LW_AWS_NETWORKING_46":      "lacework-global-157",
		"LW_AWS_NETWORKING_47":      "lacework-global-128",
		"LW_AWS_NETWORKING_49":      "lacework-global-159",
		"LW_AWS_NETWORKING_50":      "lacework-global-129", // cloudfront
		"LW_AWS_NETWORKING_51":      "lacework-global-483",
		"LW_AWS_MONGODB_1":          "lacework-global-196", // not documented
		"LW_AWS_MONGODB_2":          "lacework-global-196",
		"LW_AWS_MONGODB_3":          "lacework-global-197",
		"LW_AWS_MONGODB_4":          "lacework-global-197",
		"LW_AWS_MONGODB_5":          "lacework-global-198",
		"LW_AWS_MONGODB_6":          "lacework-global-198",
		"LW_AWS_GENERAL_SECURITY_1": "lacework-global-89", // ec2 tags
		"LW_AWS_GENERAL_SECURITY_2": "lacework-global-90",
		"LW_AWS_GENERAL_SECURITY_3": "lacework-global-160",
		"LW_AWS_GENERAL_SECURITY_4": "lacework-global-171",
		"LW_AWS_GENERAL_SECURITY_5": "lacework-global-91",
		"LW_AWS_GENERAL_SECURITY_6": "lacework-global-92",
		"LW_AWS_GENERAL_SECURITY_7": "lacework-global-182",
		"LW_AWS_GENERAL_SECURITY_8": "lacework-global-183",
		"LW_AWS_SERVERLESS_1":       "lacework-global-179",
		"LW_AWS_SERVERLESS_2":       "lacework-global-180",
		"LW_AWS_SERVERLESS_4":       "lacework-global-143",
		"LW_AWS_SERVERLESS_5":       "lacework-global-144",
		"LW_AWS_RDS_1":              "lacework-global-93",
		"LW_AWS_ELASTICSEARCH_1":    "lacework-global-122",
		"LW_AWS_ELASTICSEARCH_2":    "lacework-global-123",
		"LW_AWS_ELASTICSEARCH_3":    "lacework-global-124",
		"LW_AWS_ELASTICSEARCH_4":    "lacework-global-161",
	}

	// suppressionsMigrateAwsCmd represents the aws sub-command inside the suppressions migrate command
	suppressionsMigrateAwsCmd = &cobra.Command{
		Use:     "migrate",
		Aliases: []string{"mig"},
		Short:   "Migrate legacy suppressions for AWS to mapped policy exceptions",
		RunE:    suppressionsAwsMigrate,
	}

	// suppressionsListAwsCmd represents the aws sub-command inside the suppressions list command
	suppressionsListAwsCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List legacy suppressions for AWS",
		RunE:    suppressionsAwsList,
	}
)

func suppressionsAwsList(_ *cobra.Command, _ []string) error {
	var (
		suppressions map[string]api.SuppressionV2
		err          error
	)

	suppressions, err = cli.LwApi.V2.Suppressions.Aws.List()
	if err != nil {
		if strings.Contains(err.Error(), "No active AWS accounts") {
			cli.OutputHuman("No active AWS accounts found. " +
				"Unable to get legacy aws suppressions\n")
			return nil
		}
		return errors.Wrap(err, "Unable to get legacy aws suppressions")
	}

	if len(suppressions) == 0 {
		cli.OutputHuman("No legacy AWS suppressions found.\n")
		return nil
	}
	return cli.OutputJSON(suppressions)
}

func suppressionsAwsMigrate(_ *cobra.Command, _ []string) error {
	var (
		suppressionsMap map[string]api.SuppressionV2
		err             error

		convertedPolicyExceptions []map[string]api.PolicyException
		payloadsText              []string
		discardedSuppressions     []map[string]api.SuppressionV2
	)
	suppressionsMap, err = cli.LwApi.V2.Suppressions.Aws.List()
	if err != nil {
		if strings.Contains(err.Error(), "No active AWS accounts") {
			cli.OutputHuman("No active AWS accounts found. " +
				"Unable to get legacy aws suppressions")
			return nil
		}
		return errors.Wrap(err, "Unable to get legacy aws suppressions")
	}

	if len(suppressionsMap) == 0 {
		cli.OutputHuman("No legacy AWS suppressions found.\n")
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
		_, payloadsText, discardedSuppressions = convertAwsSuppressions(
			suppressionsMap,
			policyExceptionsConstraintsMap,
		)
		printPayloadsText(payloadsText)
		printDiscardedSuppressions(discardedSuppressions)
	case autoMigration:
		convertedPolicyExceptions, _, discardedSuppressions = convertAwsSuppressions(
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

func convertAwsSuppressions(
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
		mappedPolicyId, ok := awsEquivalencesMap[id]
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

				accountIdsConstraint := convertSupCondition(suppression.AccountIds,
					"accountIds",
					policyIdExceptionsTemplate)
				if accountIdsConstraint.FieldKey != "" {
					convertedConstraints = append(convertedConstraints, accountIdsConstraint)
				}

				regionNamesConstraint := convertSupCondition(suppression.RegionNames,
					"regionNames",
					policyIdExceptionsTemplate)
				if regionNamesConstraint.FieldKey != "" {
					convertedConstraints = append(convertedConstraints, regionNamesConstraint)
				}

				resourceNamesConstraint := convertSupCondition(suppression.ResourceNames,
					"resourceNames",
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
