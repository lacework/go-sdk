package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/lwgenerate/aws"
	"github.com/pkg/errors"
)

// survey.Validator for aws ARNs
//
// This isn't service/type specific but rather just validates that an ARN was entered that matches valid ARN formats
func validAwsArnFormat(val interface{}) error {
	// the reflect value of the result
	value := reflect.ValueOf(val)

	// if the value passed is not a string
	if value.Kind() != reflect.String {
		return errors.New("value must be a string")
	}

	// if value doesn't match regex, return invalid arn
	// original source: https://regex101.com/r/pOfxYN/1
	matchRegEx := `^arn:(?P<Partition>[^:\n]*):(?P<Service>[^:\n]*):(?P<Region>[^:\n]*):(?P<AccountID>[^:\n]*):(?P<Ignore>(?P<ResourceType>[^:\/\n]*)[:\/])?(?P<Resource>.*)$`
	ok, _ := regexp.MatchString(matchRegEx, value.String())
	if !ok {
		return errors.New("invalid arn supplied")
	}

	return nil
}

func promptAwsCtQuestions(config *aws.GenerateAwsTfConfigurationArgs) error {
	// Only ask these questions if configure cloudtrail is true
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: "Use consolidated Cloudtrail?", Default: config.ConsolidatedCloudtrail},
			Response: &config.ConsolidatedCloudtrail,
		},
		{
			Prompt: &survey.Input{
				Message: "Specify an existing bucket ARN used for Cloudtrail logs:",
				Default: config.ExistingCloudtrailBucketArn,
			},
			Checks:   []*bool{&config.ConsolidatedCloudtrail},
			Required: true,
			Opts:     []survey.AskOpt{survey.WithValidator(validAwsArnFormat)},
			Response: &config.ExistingCloudtrailBucketArn,
		},
	}, config.Cloudtrail); err != nil {
		return err
	}

	// If a new bucket is to be created; should the force destroy bit be set?
	newBucket := config.ExistingCloudtrailBucketArn == ""
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Confirm{
			Message: "Should the new S3 bucket have force destroy enabled?",
			Default: config.ForceDestroyS3Bucket},
		Response: &config.ForceDestroyS3Bucket,
		Checks:   []*bool{&config.Cloudtrail, &newBucket}}); err != nil {
		return err
	}

	return nil
}

func promptAwsExistingIamQuestions(config *aws.GenerateAwsTfConfigurationArgs) error {
	// ensure struct is initialized
	if config.ExistingIamRole == nil {
		config.ExistingIamRole = &aws.ExistingIamRoleDetails{}
	}

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt: &survey.Input{
				Message: "Specify an existing IAM role name for Cloudtrail access",
				Default: config.ExistingIamRole.Name},
			Response: &config.ExistingIamRole.Name,
			Opts:     []survey.AskOpt{survey.WithValidator(survey.Required)},
		},
		{
			Prompt: &survey.Input{
				Message: "Specify an existing IAM role ARN for Cloudtrail access",
				Default: config.ExistingIamRole.Arn,
			},
			Response: &config.ExistingIamRole.Arn,
			Opts:     []survey.AskOpt{survey.WithValidator(survey.Required), survey.WithValidator(validAwsArnFormat)},
		},
		{
			Prompt: &survey.Input{
				Message: "Specify the external ID to be used with the existing IAM role",
				Default: config.ExistingIamRole.ExternalId,
			},
			Response: &config.ExistingIamRole.ExternalId,
			Opts:     []survey.AskOpt{survey.WithValidator(survey.Required)},
		}}); err != nil {
		return err
	}

	return nil
}

func promptAwsAdditionalAccountQuestions(config *aws.GenerateAwsTfConfigurationArgs) error {
	// For each added account, collect it's profile name and the region that should be used
	accountDetails := []aws.AwsSubAccount{}
	askAgain := true

	// Determine the profile for the main account
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{
			Message: "Before adding subaccounts, your primary AWS account profile name must be set; which profile should the main account use?",
			Default: config.AwsProfile,
			Help:    "This is the main account where your cloudtrail resources are created",
		},
		Response: &config.AwsProfile,
		Required: true,
	}); err != nil {
		return nil
	}

	// For each account to add, collect the aws profile and region to use
	for askAgain {
		var accountProfileName string
		var accountProfileRegion string
		accountQuestions := []SurveyQuestionWithValidationArgs{
			{
				Prompt:   &survey.Input{Message: "Supply the profile name for this additional AWS account:"},
				Required: true,
				Response: &accountProfileName,
			},
			{
				Prompt:   &survey.Input{Message: "What region should be used for this account?"},
				Required: true,
				Response: &accountProfileRegion,
			},
		}

		if err := SurveyMultipleQuestionWithValidation(accountQuestions); err != nil {
			return err
		}

		accountDetails = append(
			accountDetails,
			aws.AwsSubAccount{AwsProfile: accountProfileName, AwsRegion: accountProfileRegion})

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Prompt:   &survey.Confirm{Message: "Add another AWS account?"},
			Response: &askAgain}); err != nil {
			return err
		}
	}
	config.SubAccounts = accountDetails

	return nil
}

// Used to test if path supplied for output exists
func validPathExists(val interface{}) error {
	// the reflect value of the result
	value := reflect.ValueOf(val)

	// if the value passed is not a string
	if value.Kind() != reflect.String {
		return errors.New("value must be a string")
	}

	// Test if supplied path exists
	if _, err := os.Stat(filepath.FromSlash(value.String())); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("supplied path %s does not exist!", value.String()))
	}

	return nil
}

func promptCustomizeOutputLocation(config *aws.GenerateAwsTfConfigurationArgs, location *string) error {
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: "Provide the location for the output to be written:", Default: *location},
		Response: location,
		Opts:     []survey.AskOpt{survey.WithValidator(validPathExists)},
		Required: true,
	}); err != nil {
		return err
	}

	return nil
}

func askAdvancedOptions(config *aws.GenerateAwsTfConfigurationArgs, outputLocation *string) error {
	answer := ""
	done := "Done"
	askCloudTrailOptions := "Additional Cloudtrail Options"
	askIamRoleOptions := "Configure Lacework integration with an existing IAM role"
	askAdditionalAwsAccountsOptions := "Add Additional AWS Accounts to Lacework"
	askCustomizeOutputLocationOptions := "Customize Output Location"

	// Prompt for options
	for answer != done {
		// Construction of this slice is a bit strange at first look, but the reason for that is because we have to do string
		// validation to know which option was selected due to how survey works; and doing it by index (also supported) is
		// difficult when the options are dynamic (which they are)
		//
		// Only ask about more accounts if consolidated cloudtrail is setup (matching scenarios doc)
		options := []string{askCloudTrailOptions, askIamRoleOptions}
		if config.ConsolidatedCloudtrail {
			options = append(options, askAdditionalAwsAccountsOptions)
		}
		options = append(options, askCustomizeOutputLocationOptions, done)
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Prompt: &survey.Select{
				Message: "Which options would you like to enable?",
				Options: options,
			},
			Response: &answer,
		}); err != nil {
			return err
		}

		// Based on response, prompt for actions
		switch answer {
		case askCloudTrailOptions:
			if err := promptAwsCtQuestions(config); err != nil {
				return err
			}
		case askIamRoleOptions:
			if err := promptAwsExistingIamQuestions(config); err != nil {
				return err
			}
		case askAdditionalAwsAccountsOptions:
			if err := promptAwsAdditionalAccountQuestions(config); err != nil {
				return err
			}
		case askCustomizeOutputLocationOptions:
			if err := promptCustomizeOutputLocation(config, outputLocation); err != nil {
				return err
			}
		}

		// Re-prompt if not done
		innerAskAgain := true
		if answer == done {
			innerAskAgain = false
		}

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Checks:   []*bool{&innerAskAgain},
			Prompt:   &survey.Confirm{Message: "Configure another advanced integration option", Default: false},
			Response: &innerAskAgain,
		}); err != nil {
			return err
		}

		if !innerAskAgain {
			answer = done
		}
	}

	return nil
}

// entry point for launching a survey to build out the required generation parameters
func promptAwsGenerate(
	config *aws.GenerateAwsTfConfigurationArgs,
	existingIam *aws.ExistingIamRoleDetails,
	outputLocation *string,
) error {
	// Cache for later use if generation is abandon
	defer cli.WriteAssetToCache("iac-aws-generate-params", time.Now().Add(time.Hour*1), config)

	// Set ExistingIamRole details, if provided as cli flags; otherwise don't initialize
	if existingIam.Arn != "" ||
		existingIam.Name != "" ||
		existingIam.ExternalId != "" {
		config.ExistingIamRole = existingIam
	}

	// This are the core questions that should be asked.  Region required for provider block
	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Prompt:   &survey.Confirm{Message: "Enable Config Integration?", Default: config.Config},
				Response: &config.Config,
			},
			{
				Prompt:   &survey.Confirm{Message: "Enable Cloudtrail Integration?", Default: config.Cloudtrail},
				Response: &config.Cloudtrail,
			},
			{
				Checks: []*bool{&config.Config, &config.Cloudtrail},
				Prompt: &survey.Input{
					Message: "Specify the AWS region Cloudtrail, SNS, and S3 resources should use:",
					Default: config.AwsRegion,
				},
				Response: &config.AwsRegion,
				Opts:     []survey.AskOpt{survey.WithValidator(survey.Required)},
			},
		}); err != nil {
		return err
	}

	// Validate one of config or cloudtrail was enabled; otherwise error out
	if !config.Config && !config.Cloudtrail {
		return errors.New("Must enable cloudtrail or config!")
	}

	// Find out if the customer wants to specify more advanced features
	askAdvanced := false
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Message: "Configure advanced integration options?", Default: askAdvanced},
		Response: &askAdvanced,
	}); err != nil {
		return err
	}

	// Keep prompting for advanced options until the say done
	if askAdvanced {
		if err := askAdvancedOptions(config, outputLocation); err != nil {
			return err
		}
	}

	return nil
}
