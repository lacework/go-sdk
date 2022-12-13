package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var (
	// constraint types
	// ResourceNames, in the old CIS1.1, works for EC2 instance ID,
	// VPC ID, Group ID, ARN, ELB ID, Lambda name, IAM policy, etc)
	allAwsConditions           = []string{"accountIds", "regionNames", "resourceNames", "resourceTags"}
	noRegionNames              = []string{"accountIds", "resourceNames", "resourceTags"}
	accountIdsOnly             = []string{"accountIds"}
	accountIdsAndResourceNames = []string{"accountIds", "resourceNames"}
	noResourceTags             = []string{"accountIds", "regionNames", "resourceNames"}

	// https://docs.lacework.com/console/aws-compliance-policy-exceptions-criteria#lacework-custom-policies-for-aws-iam
	// https://docs.lacework.com/console/cis-aws-140-benchmark-report#identity-and-access-management
	// old ID to new ID mapping, using the old Constraints with the hope they match the new Constraints
	awsEquivalencesMap = map[string]string{
		"AWS_CIS_1_2":   "lacework-global-39",
		"AWS_CIS_1_3":   "lacework-global-41",
		"AWS_CIS_1_4":   "lacework-global-43",
		"AWS_CIS_1_9":   "lacework-global-37",
		"AWS_CIS_1_10":  "lacework-global-38",
		"AWS_CIS_1_12":  "lacework-global-34",
		"AWS_CIS_1_13":  "lacework-global-35",
		"AWS_CIS_1_14":  "lacework-global-69",
		"AWS_CIS_1_15":  "lacework-global-33",
		"AWS_CIS_1_16":  "lacework-global-44", // no iam policies to users
		"AWS_CIS_1_19":  "lacework-global-31", // manual?
		"AWS_CIS_1_20":  "lacework-global-32", // manual?
		"AWS_CIS_1_21":  "lacework-global-70", // manual?
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

	policyExceptionTemplatesMap = map[string][]string{
		"lacework-global-39":  accountIdsAndResourceNames,
		"lacework-global-41":  allAwsConditions,
		"lacework-global-43":  accountIdsAndResourceNames,
		"lacework-global-37":  accountIdsOnly,
		"lacework-global-38":  accountIdsOnly,
		"lacework-global-34":  accountIdsOnly,
		"lacework-global-35":  accountIdsOnly,
		"lacework-global-69":  accountIdsOnly,
		"lacework-global-33":  accountIdsOnly,
		"lacework-global-44":  accountIdsAndResourceNames, // no iam policies to users
		"lacework-global-31":  accountIdsOnly,             // manual?
		"lacework-global-32":  accountIdsOnly,             // manual?
		"lacework-global-70":  noRegionNames,              // manual?
		"lacework-global-46":  accountIdsOnly,
		"lacework-global-40":  accountIdsAndResourceNames,
		"lacework-global-45":  accountIdsAndResourceNames,
		"lacework-global-53":  accountIdsOnly,
		"lacework-global-75":  noResourceTags,
		"lacework-global-54":  noRegionNames, // s3 bucket cloudtrail log
		"lacework-global-55":  noResourceTags,
		"lacework-global-76":  noResourceTags,
		"lacework-global-56":  noRegionNames, // s3 bucket cloudtrail log
		"lacework-global-77":  allAwsConditions,
		"lacework-global-78":  allAwsConditions,
		"lacework-global-79":  allAwsConditions,
		"lacework-global-57":  accountIdsOnly,
		"lacework-global-58":  accountIdsOnly,
		"lacework-global-59":  accountIdsOnly,
		"lacework-global-60":  accountIdsOnly,
		"lacework-global-61":  accountIdsOnly,
		"lacework-global-82":  accountIdsOnly,
		"lacework-global-83":  accountIdsOnly,
		"lacework-global-62":  accountIdsOnly,
		"lacework-global-84":  accountIdsOnly,
		"lacework-global-85":  accountIdsOnly,
		"lacework-global-86":  accountIdsOnly,
		"lacework-global-63":  accountIdsOnly,
		"lacework-global-64":  accountIdsOnly,
		"lacework-global-65":  accountIdsOnly,
		"lacework-global-68":  allAwsConditions,
		"lacework-global-87":  allAwsConditions,
		"lacework-global-130": allAwsConditions,
		"lacework-global-131": allAwsConditions,
		"lacework-global-132": allAwsConditions,
		"lacework-global-133": allAwsConditions,
		"lacework-global-134": allAwsConditions,
		"lacework-global-135": allAwsConditions,
		"lacework-global-136": allAwsConditions,
		"lacework-global-137": allAwsConditions,
		"lacework-global-138": allAwsConditions,
		"lacework-global-139": allAwsConditions,
		"lacework-global-140": allAwsConditions,
		"lacework-global-94":  allAwsConditions,
		"lacework-global-95":  allAwsConditions,
		"lacework-global-217": allAwsConditions,
		"lacework-global-96":  allAwsConditions,
		"lacework-global-97":  allAwsConditions,
		"lacework-global-98":  allAwsConditions,
		"lacework-global-99":  allAwsConditions,
		"lacework-global-100": allAwsConditions,
		"lacework-global-101": allAwsConditions,
		"lacework-global-115": accountIdsAndResourceNames,
		"lacework-global-116": accountIdsAndResourceNames,
		"lacework-global-117": accountIdsAndResourceNames,
		"lacework-global-118": accountIdsAndResourceNames,
		"lacework-global-119": accountIdsAndResourceNames,
		"lacework-global-120": accountIdsAndResourceNames,
		"lacework-global-121": accountIdsAndResourceNames,
		"lacework-global-181": accountIdsOnly, // non-root user
		"lacework-global-142": accountIdsAndResourceNames,
		"lacework-global-141": accountIdsAndResourceNames,
		"lacework-global-105": accountIdsAndResourceNames,
		"lacework-global-227": noRegionNames, // sec-group
		"lacework-global-145": noRegionNames, // network acl
		"lacework-global-146": noRegionNames, // network acl
		"lacework-global-147": accountIdsAndResourceNames,
		"lacework-global-148": allAwsConditions,
		"lacework-global-149": allAwsConditions,
		"lacework-global-228": allAwsConditions,
		"lacework-global-229": allAwsConditions,
		"lacework-global-230": allAwsConditions,
		"lacework-global-231": allAwsConditions,
		"lacework-global-199": allAwsConditions,
		"lacework-global-150": allAwsConditions,
		"lacework-global-151": allAwsConditions,
		"lacework-global-152": allAwsConditions,
		"lacework-global-153": allAwsConditions,
		"lacework-global-225": allAwsConditions,
		"lacework-global-226": allAwsConditions,
		"lacework-global-154": allAwsConditions,
		"lacework-global-155": allAwsConditions,
		"lacework-global-156": allAwsConditions,
		"lacework-global-104": allAwsConditions,
		"lacework-global-106": allAwsConditions,
		"lacework-global-107": allAwsConditions,
		"lacework-global-108": allAwsConditions,
		"lacework-global-109": allAwsConditions,
		"lacework-global-110": allAwsConditions,
		"lacework-global-111": allAwsConditions,
		"lacework-global-112": allAwsConditions,
		"lacework-global-113": allAwsConditions,
		"lacework-global-114": allAwsConditions,
		"lacework-global-218": allAwsConditions,
		"lacework-global-219": allAwsConditions,
		"lacework-global-220": allAwsConditions,
		"lacework-global-221": allAwsConditions,
		"lacework-global-222": allAwsConditions,
		"lacework-global-102": allAwsConditions,
		"lacework-global-223": allAwsConditions,
		"lacework-global-184": allAwsConditions,
		"lacework-global-103": allAwsConditions,
		"lacework-global-125": noRegionNames, // cloudfront
		"lacework-global-126": noRegionNames, // cloudfront
		"lacework-global-127": allAwsConditions,
		"lacework-global-482": allAwsConditions,
		"lacework-global-157": allAwsConditions,
		"lacework-global-128": allAwsConditions,
		"lacework-global-159": allAwsConditions,
		"lacework-global-129": noRegionNames, // cloudfront
		"lacework-global-483": allAwsConditions,
		"lacework-global-196": allAwsConditions,
		"lacework-global-197": allAwsConditions,
		"lacework-global-198": allAwsConditions,
		"lacework-global-89":  noRegionNames, // ec2 tags
		"lacework-global-90":  allAwsConditions,
		"lacework-global-160": noRegionNames,
		"lacework-global-171": noRegionNames,
		"lacework-global-91":  noRegionNames,
		"lacework-global-92":  noResourceTags,
		"lacework-global-182": noRegionNames,
		"lacework-global-183": noRegionNames,
		"lacework-global-179": allAwsConditions,
		"lacework-global-180": allAwsConditions,
		"lacework-global-143": allAwsConditions,
		"lacework-global-144": allAwsConditions,
		"lacework-global-93":  allAwsConditions,
		"lacework-global-122": noResourceTags,
		"lacework-global-123": noResourceTags,
		"lacework-global-124": noResourceTags,
		"lacework-global-161": noResourceTags,
	}

	// suppressionsMigrateAwsCmd represents the aws sub-command inside the suppressions migrate command
	suppressionsMigrateAwsCmd = &cobra.Command{
		Use:     "aws",
		Aliases: []string{"aws"},
		Short:   "Migrate legacy suppressions for AWS to mapped policy exceptions",
		RunE:    suppressionsAwsMigrate,
	}

	// suppressionsListAwsCmd represents the aws sub-command inside the suppressions list command
	suppressionsListAwsCmd = &cobra.Command{
		Use:     "aws",
		Aliases: []string{"aws"},
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
		return errors.Wrap(err, "unable to get legacy aws suppressions")
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
	answer := ""
	manualMigration := "Output translated legacy suppressions as policy exception commands to be run manually"
	autoMigration := "Auto migrate legacy suppressions. DISCLAIMER: " +
		"By selecting this option, you accept responsibility for " +
		"the migration and any compliance violations missed as a result of the added exceptions"
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Select{
			Message: "Which migration approach would you like to take?",
			Options: []string{
				manualMigration,
				autoMigration,
			},
		},
		Response: &answer,
	}); err != nil {
		return err
	}

	suppressionsMap, err = cli.LwApi.V2.Suppressions.Aws.List()
	if err != nil {
		return errors.Wrap(err, "unable to get legacy aws suppressions")
	}

	switch answer {
	case manualMigration:
		_, payloadsText, discardedSuppressions = convertAwsSuppressions(suppressionsMap)
		printPayloadsText(payloadsText)
		printDiscardedSuppressions(discardedSuppressions)
	case autoMigration:
		convertedPolicyExceptions, _, discardedSuppressions = convertAwsSuppressions(
			suppressionsMap)
		printConvertedSuppressions(convertedPolicyExceptions)
		confirm := false
		err := survey.AskOne(&survey.Confirm{
			Message: "Confirm that you have reviewed the above exceptions and wish to continue" +
				" with the auto migration.",
		}, &confirm)
		if err != nil {
			return err
		}
		if confirm {
			autoConvertAwsSuppressions(convertedPolicyExceptions)
			printDiscardedSuppressions(discardedSuppressions)
			cli.OutputHuman("To view the newly created Exceptions, " +
				"try running `lacework policy-exceptions list <policyId>")
		} else {
			cli.OutputHuman("Cancelled Legacy Suppression to Exception migration!")
		}
	}

	return nil
}

func autoConvertAwsSuppressions(convertedPolicyExceptions []map[string]api.PolicyException) {
	cli.StartProgress("Creating policy exceptions ...")
	for _, exceptionMap := range convertedPolicyExceptions {
		for policyId, exception := range exceptionMap {
			response, err := cli.LwApi.V2.Policy.Exceptions.Create(policyId, exception)
			if err != nil {
				cli.Log.Debug(err, "unable to create exception")
				continue
			}
			cli.OutputHuman("Exception created for PolicyId: %s - ExceptionId: %s\n\n",
				policyId, response.Data.ExceptionID)
		}
	}

	cli.StopProgress()
}

func convertAwsSuppressions(suppressionsMap map[string]api.SuppressionV2) ([]map[string]api.PolicyException,
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
				discardedSuppressions = append(
					discardedSuppressions,
					map[string]api.SuppressionV2{id: suppressionInfo},
				)
			}
			continue
		}

		// get the policy exception template for the mapped policy
		// the exception template defines the exception fields that are supported for the policy
		policyIdExceptionsTemplate := policyExceptionTemplatesMap[mappedPolicyId]
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

func printPayloadsText(payloadsText []string) {
	if len(payloadsText) >= 1 {
		cli.OutputHuman("#### Legacy Suppressions --> Exceptions payloads\n")
		for _, payload := range payloadsText {
			cli.OutputHuman("%s \n\n", payload)
		}
	} else {
		cli.OutputHuman("No legacy suppressions found that could be migrated\n")
	}
}

func printConvertedSuppressions(convertedSuppressions []map[string]api.PolicyException) {
	if len(convertedSuppressions) >= 1 {
		cli.OutputHuman("#### Converted legacy suppressions in Policy Exception format\n")
		for _, exception := range convertedSuppressions {
			b, err := json.Marshal(exception)
			if err != nil {
				return
			}
			cli.OutputHuman("%s \n\n", string(b))
		}
		cli.OutputHuman("WARNING: Before continuing, " +
			"please thoroughly inspect the above exceptions to ensure they are valid and" +
			" required. Lacework is not responsible for any compliance violations missed as a " +
			"result of the above exceptions!\n\n")
	}
}

func printDiscardedSuppressions(discardedSuppressions []map[string]api.SuppressionV2) {
	if len(discardedSuppressions) >= 1 {
		cli.OutputHuman("#### Discarded legacy suppressions\n")
		for _, suppression := range discardedSuppressions {
			b, err := json.Marshal(suppression)
			if err != nil {
				return
			}
			cli.OutputHuman("%s \n\n", string(b))
		}
	}
}

func convertSupCondition(supCondition []string, fieldKey string,
	policyIdExceptionsTemplate []string) api.PolicyExceptionConstraint {
	if len(supCondition) >= 1 && slices.Contains(
		policyIdExceptionsTemplate, fieldKey) {

		var condition []any
		// verify if "ALL_ACCOUNTS" OR "ALL_REGIONS" is in the suppression condition slice
		// if so we should ignore the supplied conditions and replace with a wildcard *
		if (slices.Contains(supCondition, "ALL_ACCOUNTS") && fieldKey == "accountIds") ||
			(slices.Contains(supCondition, "ALL_REGIONS") && fieldKey == "regionNames") {
			condition = append(condition, "*")
		} else {
			condition = convertToAnySlice(supCondition)
		}

		return api.PolicyExceptionConstraint{
			FieldKey:    fieldKey,
			FieldValues: condition,
		}
	}
	return api.PolicyExceptionConstraint{}
}

func convertSupConditionTags(supCondition []map[string]string, fieldKey string,
	policyIdExceptionsTemplate []string) api.PolicyExceptionConstraint {
	if len(supCondition) >= 1 && slices.Contains(
		policyIdExceptionsTemplate, fieldKey) {

		// api.PolicyExceptionConstraint expects []any for the FieldValues
		// Therefore we need to take the supCondition []map[string]string and append each map to
		// the new convertedTags []any var
		var convertedTags []any
		for _, tagMap := range supCondition {
			convertedTags = append(convertedTags, tagMap)
		}

		return api.PolicyExceptionConstraint{
			FieldKey:    fieldKey,
			FieldValues: convertedTags,
		}
	}
	return api.PolicyExceptionConstraint{}
}

func convertToAnySlice(slice []string) []any {
	s := make([]interface{}, len(slice))
	for i, v := range slice {
		s[i] = v
	}
	return s
}
