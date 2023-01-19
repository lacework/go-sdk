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
	"github.com/fatih/color"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/lacework/go-sdk/api"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var (
	// top-level suppressions command
	suppressionsCommand = &cobra.Command{
		Use:     "suppressions",
		Hidden:  true,
		Aliases: []string{"suppression", "sup", "sups"},
		Short:   "Manage legacy suppressions",
		Long:    "Manage legacy suppressions",
	}

	// suppressionsAwsCmd represents the aws sub-command inside the suppressions command
	suppressionsAwsCmd = &cobra.Command{
		Use:   "aws",
		Short: "Manage legacy suppressions for aws",
	}

	// suppressionsAzureCmd represents the aws sub-command inside the suppressions command
	suppressionsAzureCmd = &cobra.Command{
		Use:   "azure",
		Short: "Manage legacy suppressions for azure",
	}

	// suppressionsGcpCmd represents the aws sub-command inside the suppressions command
	suppressionsGcpCmd = &cobra.Command{
		Use:   "gcp",
		Short: "Manage legacy suppressions for gcp",
	}
)

func init() {
	rootCmd.AddCommand(suppressionsCommand)
	// aws
	suppressionsCommand.AddCommand(suppressionsAwsCmd)
	suppressionsAwsCmd.AddCommand(suppressionsListAwsCmd)
	suppressionsAwsCmd.AddCommand(suppressionsMigrateAwsCmd)
	// azure
	suppressionsCommand.AddCommand(suppressionsAzureCmd)
	suppressionsAzureCmd.AddCommand(suppressionsListAzureCmd)
	// gcp
	suppressionsCommand.AddCommand(suppressionsGcpCmd)
	suppressionsGcpCmd.AddCommand(suppressionsListGcpCmd)
	suppressionsGcpCmd.AddCommand(suppressionsMigrateGcpCmd)
}

func autoConvertSuppressions(convertedPolicyExceptions []map[string]api.PolicyException) {
	cli.StartProgress("Creating policy exceptions ...")
	for _, exceptionMap := range convertedPolicyExceptions {
		for policyId, exception := range exceptionMap {
			response, err := cli.LwApi.V2.Policy.Exceptions.Create(policyId, exception)
			if err != nil {
				cli.Log.Debug(err, "unable to create exception")
				cli.OutputHuman(color.RedString(
					"Error creating policy exception to create exception. %e"),
					err)
				continue
			}
			cli.OutputHuman("Exception created for PolicyId: %s - ExceptionId: %s\n\n",
				color.GreenString(policyId), color.BlueString(response.Data.ExceptionID))
		}
	}

	cli.StopProgress()
}

func printPayloadsText(payloadsText []string) {
	if len(payloadsText) >= 1 {
		cli.OutputHuman(color.YellowString("#### Legacy Suppressions --> Exceptions payloads\n\n"))
		for _, payload := range payloadsText {
			cli.OutputHuman(color.GreenString("%s \n\n", payload))
		}
	} else {
		cli.OutputHuman("No legacy suppressions found that could be migrated\n")
	}
}

func printConvertedSuppressions(convertedSuppressions []map[string]api.PolicyException) {
	if len(convertedSuppressions) >= 1 {
		cli.OutputHuman(color.YellowString("#### Converted legacy suppressions in Policy Exception" +
			" format" +
			"\n"))
		for _, exception := range convertedSuppressions {
			err := cli.OutputJSON(exception)
			if err != nil {
				return
			}
		}
		colorizeR := color.New(color.FgRed, color.Bold)
		cli.OutputHuman(colorizeR.Sprintf("WARNING: Before continuing, " +
			"please thoroughly inspect the above exceptions to ensure they are valid and" +
			" required. By continuing, you accept liability for any compliance violations" +
			" missed as a result of the above exceptions!\n\n"))

	}
}

func printDiscardedSuppressions(discardedSuppressions []map[string]api.SuppressionV2) {
	if len(discardedSuppressions) >= 1 {
		cli.OutputHuman(color.YellowString("#### Discarded legacy suppressions\n"))
		for _, suppression := range discardedSuppressions {
			err := cli.OutputJSON(suppression)
			if err != nil {
				return
			}
		}
	}
}

func convertSupCondition(supConditions []string, fieldKey string,
	policyIdExceptionsTemplate []string) api.PolicyExceptionConstraint {
	if len(supConditions) >= 1 && slices.Contains(
		policyIdExceptionsTemplate, fieldKey) {

		var condition []any
		// verify for aws:
		// if "ALL_ACCOUNTS" OR "ALL_REGIONS" is in the suppression condition slice
		// verify for gcp:
		// if "ALL_ORGANIZATIONS" OR "ALL_PROJECTS" is in the suppression condition slice
		// if so we should ignore the supplied conditions and replace with a wildcard *
		if (slices.Contains(supConditions, "ALL_ACCOUNTS") && fieldKey == "accountIds") ||
			(slices.Contains(supConditions, "ALL_REGIONS") && fieldKey == "regionNames") {
			condition = append(condition, "*")
		} else if (slices.Contains(supConditions, "ALL_ORGANIZATIONS") && fieldKey == "organizations") ||
			(slices.Contains(supConditions, "ALL_PROJECTS") && fieldKey == "projects") {
			condition = append(condition, "*")
		} else if fieldKey == "resourceNames" {
			condition = convertResourceNamesSupConditions(supConditions)
		} else if fieldKey == "resourceName" {
			// resourceName singular is specific to GCP
			condition = convertGcpResourceNameSupConditions(supConditions)
		} else {
			condition = convertToAnySlice(supConditions)
		}

		return api.PolicyExceptionConstraint{
			FieldKey:    fieldKey,
			FieldValues: condition,
		}
	}
	return api.PolicyExceptionConstraint{}
}

func convertResourceNamesSupConditions(supConditions []string) []any {
	var conditions []any
	for _, condition := range supConditions {
		ok := arn.IsARN(condition)
		if ok {
			parsedEntry, _ := arn.Parse(condition)
			condition = parsedEntry.Resource
		}
		conditions = append(conditions, condition)
	}
	return conditions
}

func convertGcpResourceNameSupConditions(supConditions []string) []any {
	var conditions []any
	for _, condition := range supConditions {
		// skip this logic if we already have a wildcard
		if condition != "*" {
			// It appears that for GCP, the resourceName field for policy exceptions is in fact expecting
			// users to provider the full GCP resource_id.
			// Example resourceId: //compute.googleapis.com/projects/gke-project-01-c8403ba1/zones/us-central1-a/instances/squid-proxy
			// This was not the case for legacy suppressions and in most cases it's unlikely that the
			// users will have provided this. Instead, we are more likely to have
			// the resource name provided. To cover this scenario we prepend the resource name
			// from the legacy suppression with "*/" to make it match the resource name while
			// wildcarding the rest of the resourceId
			condition = "*/" + condition
		}
		conditions = append(conditions, condition)
	}
	return conditions
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

func getPoliciesExceptionConstraintsMap() map[string][]string {
	// get a list of all policies and parse the valid exception constraints and return a map of
	// {"<policyId>": [<validPolicyConstraints>]}
	policies, err := cli.LwApi.V2.Policy.List()
	if err != nil {
		return nil
	}

	policiesSupportedConstraints := make(map[string][]string)
	for _, policy := range policies.Data {
		exceptionConstraints := getPolicyExceptionConstraintsSlice(policy.ExceptionConfiguration)
		policiesSupportedConstraints[policy.PolicyID] = exceptionConstraints
	}

	return policiesSupportedConstraints
}

func convertToAnySlice(slice []string) []any {
	s := make([]interface{}, len(slice))
	for i, v := range slice {
		s[i] = v
	}
	return s
}
