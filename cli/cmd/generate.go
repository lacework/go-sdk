package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type AwsGenerateCommandExtraState struct {
	Output                string
	UseExistingCloudtrail bool
	AwsSubAccounts        []string
	TerraformApply        bool
}

func (a *AwsGenerateCommandExtraState) isEmpty() bool {
	return a.Output == "" && !a.UseExistingCloudtrail && len(a.AwsSubAccounts) == 0 && !a.TerraformApply
}

// Flush current state of the struct to disk, provided it's not empty
func (a *AwsGenerateCommandExtraState) writeCache() {
	if !a.isEmpty() {
		cli.WriteAssetToCache(CachedAssetAwsExtraState, time.Now().Add(time.Hour*1), a)
	}
}

var (
	QuestionRunTfPlan        = "Run Terraform plan now?"
	QuestionUsePreviousCache = "Previous IaC generation detected, load cached values?"

	// iac-generate command is used to create IaC code for various environments
	generateTfCommand = &cobra.Command{
		Use:     "iac-generate",
		Aliases: []string{"iac"},
		Short:   "Create IaC code",
		Long:    "Create IaC content for various different cloud environments and configurations",
	}
)

func init() {
	// add the iac-generate command
	cloudAccountCommand.AddCommand(generateTfCommand)

	// Add global flags for iac generation
	generateTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandExtraState.Output, "output", "", "location to write generated content")

	// add flags to sub commands
	// TODO Share the help with the interactive generation
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.Cloudtrail, "cloudtrail", false, "enable cloudtrail integration")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.Config, "config", false, "enable config integration")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.AwsRegion, "aws_region", "", "specify aws region")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.AwsProfile, "aws_profile", "default", "specify aws profile")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ExistingCloudtrailBucketArn,
		"existing_bucket_arn",
		"",
		"specify existing cloudtrail s3 bucket ARN")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsExistingRoleState.Arn,
		"existing_iam_role_arn",
		"",
		"specify existing iam role arn to use")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsExistingRoleState.Name,
		"existing_iam_role_name",
		"",
		"specify existing iam role name to use")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsExistingRoleState.ExternalId,
		"existing_iam_role_externalid",
		"",
		"specify existing iam role external_id to use")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ExistingSnsTopicArn,
		"existing_sns_topic_arn",
		"",
		"specify existing sns topic arn")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.ConsolidatedCloudtrail,
		"consolidated_cloudtrail",
		false,
		"use consolidated trail")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.ForceDestroyS3Bucket,
		"force_destroy_s3",
		false,
		"enable force destroy s3 bucket")
	generateAwsTfCommand.PersistentFlags().StringSliceVar(
		&GenerateAwsCommandExtraState.AwsSubAccounts,
		"aws_subaccount",
		[]string{},
		"configure an additional aws account; value format must be <aws profile>:<region>")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandExtraState.TerraformApply,
		"apply",
		false,
		"run terraform apply without executing plan or prompting",
	)

	// add sub-commands to the iac-generate command
	generateTfCommand.AddCommand(generateAwsTfCommand)
}

type SurveyQuestionWithValidationArgs struct {
	Prompt survey.Prompt
	// Supplied checks can be used to validate IF the question should be asked
	Checks   []*bool
	Response interface{}
	Opts     []survey.AskOpt
	Required bool
}

// SurveyQuestionInteractiveOnly Prompt use for question, only if the CLI is in interactive mode
func SurveyQuestionInteractiveOnly(question SurveyQuestionWithValidationArgs) error {
	// Do validations pass?
	ok := true
	for _, v := range question.Checks {
		if !*v {
			ok = false
		}
	}

	// If the optional check doesn't pass, skip
	if !ok {
		return nil
	}

	// If required is set, add that question opt
	if question.Required {
		question.Opts = append(question.Opts, survey.WithValidator(survey.Required))
	}

	// Add custom icon
	question.Opts = append(question.Opts, survey.WithIcons(promptIconsFunc))

	// If noninteractive is not set, ask the question
	if !cli.nonInteractive {
		err := survey.AskOne(question.Prompt, question.Response, question.Opts...)
		if err != nil {
			return err
		}
	}

	return nil
}

// SurveyMultipleQuestionWithValidation Prompt for many values at once
//
// checks: If supplied check(s) are true, questions will be asked
func SurveyMultipleQuestionWithValidation(questions []SurveyQuestionWithValidationArgs, checks ...bool) error {
	// Do validations pass?
	ok := true
	for _, v := range checks {
		if !v {
			ok = false
		}
	}

	// Ask questions
	if ok {
		for _, qs := range questions {
			if err := SurveyQuestionInteractiveOnly(qs); err != nil {
				return err
			}
		}
	}
	return nil
}

// determineOutputDirPath get output directory location based on how the output location was set
func determineOutputDirPath(location string) (string, error) {
	// determine code output path
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// If location was passed, return that location
	if location != "" {
		return filepath.FromSlash(location), nil
	}

	// If location was not passed, assemble it with lacework from os homedir
	return filepath.FromSlash(fmt.Sprintf("%s/%s", dirname, "lacework")), nil
}

// writeHclOutputPreCheck Prompt for confirmation if main.tf already exists; return true to continue
func writeHclOutputPreCheck(outputLocation string, filename string) (bool, error) {
	// If noninteractive, continue
	if !cli.InteractiveMode() {
		return true, nil
	}

	outputDir, err := determineOutputDirPath(outputLocation)
	if err != nil {
		return false, err
	}

	if filename == "" {
		filename = "main"
	}

	hclPath := filepath.FromSlash(fmt.Sprintf("%s/%s.tf", outputDir, filename))

	// If the file doesn't exist, carry on
	if _, err := os.Stat(hclPath); os.IsNotExist(err) {
		return true, nil
	}

	// If it does exist; confirm overwrite
	answer := false
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Message: fmt.Sprintf("%s already exists, overwrite?", hclPath)},
		Response: &answer,
	}); err != nil {
		return false, err
	}

	return answer, nil
}

// writeHclOutput Write HCL output
func writeHclOutput(hcl string, location string, filename string) (string, error) {
	// Determine write location
	dirname, err := determineOutputDirPath(location)
	if err != nil {
		return "", err
	}

	// Create directory, if needed
	if location == "" {
		directory := filepath.FromSlash(dirname)
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			err = os.Mkdir(directory, 0700)
			if err != nil {
				return "", err
			}
		}
	}

	if filename == "" {
		filename = "main"
	}

	// Create HCL file
	outputLocation := filepath.FromSlash(fmt.Sprintf("%s/%s.tf", dirname, filename))
	err = os.WriteFile(
		filepath.FromSlash(outputLocation),
		[]byte(hcl),
		0700,
	)
	if err != nil {
		return "", err
	}

	cli.StopProgress()
	return outputLocation, nil
}

// validateOutputLocation This function used to validate provided output location exists and is a directory
func validateOutputLocation(dirname string) error {
	// If output location was supplied, validate it exists
	if dirname != "" {
		outputLocation, err := os.Stat(dirname)
		if err != nil {
			return errors.Wrap(err, "could not access specified output location")
		}

		if !outputLocation.IsDir() {
			return errors.New("output location must be a directory")
		}
	}

	return nil
}

func validateAwsSubAccounts(subaccounts []string) error {
	// validate the format of supplied values is correct
	for _, account := range subaccounts {
		if ok, err := regexp.MatchString(ValidateSubAccountFlagRegex, account); !ok {
			if err != nil {
				return errors.Wrap(err, "failed to validate supplied subaccount format")
			}
			return errors.New("supplied aws subaccount in invalid format")
		}
	}

	return nil
}

// create survey.Validator for string with regex
func validateStringWithRegex(val interface{}, regex string, errorString string) error {
	switch value := val.(type) {
	case string:
		// if value doesn't match regex, return invalid arn
		ok, err := regexp.MatchString(regex, value)
		if err != nil {
			return errors.Wrap(err, "failed to validate input")
		}

		if !ok {
			return errors.New(errorString)
		}
	default:
		// if the value passed is not a string
		return errors.New("value must be a string")
	}

	return nil
}

// Used to test if path supplied for output exists
func validPathExists(val interface{}) error {
	switch value := val.(type) {
	case string:
		// Test if supplied path exists
		if err := validateOutputLocation(value); err != nil {
			return err
		}
	default:
		// if the value passed is not a string
		return errors.New("value must be a string")
	}

	return nil
}

// writeGeneratedCodeToLocation Write-out generated code to location specified
func writeGeneratedCodeToLocation(cmd *cobra.Command, hcl string, cloud string) (string, string, error) {
	//dirname, ok, location := "", false, ""
	// Write-out generated code to location specified
	dirname, err := cmd.Flags().GetString("output")
	if err != nil {
		return dirname, "", errors.Wrap(err, "failed to parse output location")
	}

	ok, err := writeHclOutputPreCheck(dirname, cloud)
	if err != nil {
		return dirname, "", errors.Wrap(err, "failed to validate output location")
	}

	if !ok {
		return dirname, "", errors.Wrap(err, "aborting to avoid overwriting existing terraform code")
	}

	location, err := writeHclOutput(hcl, dirname, cloud)
	if err != nil {
		return dirname, location, errors.Wrap(err, "failed to write terraform code to disk")
	}

	return dirname, location, nil
}

// executionPreRunChecks Execution pre-run check
func executionPreRunChecks(dirname string, locationDir string) error {
	ok, err := TerraformExecutePreRunCheck(dirname)
	if err != nil {
		return errors.Wrap(err, "failed to check for existing terraform state")
	}

	if !ok {
		cli.OutputHuman(provideGuidanceAfterExit(false, false, locationDir, "terraform"))
		return nil
	}

	if err := TerraformPlanAndExecute(locationDir); err != nil {
		return errors.Wrap(err, "failed to run terraform apply")
	}

	return nil
}
