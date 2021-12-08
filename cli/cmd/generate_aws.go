package cmd

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/lwgenerate/aws"
	"github.com/pkg/errors"
)

var (
	// Define question text here so they can be reused in testing
	QuestionEnableConfig                = "Enable Config Integration?"
	QuestionEnableCloudtrail            = "Enable Cloudtrail Integration?"
	QuestionAwsRegion                   = "Specify the AWS region to be used by Cloudtrail, SNS, and S3:"
	QuestionConsolidatedCloudtrail      = "Use consolidated Cloudtrail?"
	QuestionUseExistingCloudtrail       = "Use an existing Cloudtrail?"
	QuestionCloudtrailExistingBucketArn = "Specify an existing bucket ARN used for Cloudtrail logs:"
	QuestionForceDestroyS3Bucket        = "Should the new S3 bucket have force destroy enabled?"
	QuestionExistingIamRoleName         = "Specify an existing IAM role name for Cloudtrail access"
	QuestionExistingIamRoleArn          = "Specify an existing IAM role ARN for Cloudtrail access"
	QuestionExistingIamRoleExtId        = "Specify the external ID to be used with the existing IAM role"
	QuestionPrimaryAwsAccountProfile    = "Before adding subaccounts, your primary AWS account profile name must be set; which profile should the main account use?"
	QuestionSubAccountProfileName       = "Supply the profile name for this additional AWS account:"
	QuestionSubAccountRegion            = "What region should be used for this account?"
	QuestionSubAccountAddMore           = "Add another AWS account?"
	QuestionSubAccountReplace           = "Currently configured AWS subaccounts: %s, replace?"
	QuestionConfigAdvanced              = "Configure advanced integration options?"
	QuestionAnotherAdvancedOpt          = "Configure another advanced integration option"
	QuestionCustomizeOutputLocation     = "Provide the location for the output to be written:"

	// select options
	AdvancedOptDone        = "Done"
	AdvancedOptCloudTrail  = "Additional Cloudtrail options"
	AdvancedOptIamRole     = "Configure Lacework integration with an existing IAM role"
	AdvancedOptAwsAccounts = "Add additional AWS Accounts to Lacework"
	AdvancedOptLocation    = "Customize output location"

	// original source: https://regex101.com/r/pOfxYN/1
	AwsArnRegex = `^arn:(?P<Partition>[^:\n]*):(?P<Service>[^:\n]*):(?P<Region>[^:\n]*):(?P<AccountID>[^:\n]*):(?P<Ignore>(?P<ResourceType>[^:\/\n]*)[:\/])?(?P<Resource>.*)$`
	// regex used for validating region input; note intentionally does not match gov cloud
	AwsRegionRegex  = `(us|ap|ca|cn|eu|sa)-(central|(north|south)?(east|west)?)-\d`
	AwsProfileRegex = `([A-Za-z_0-9-]+)`
)

// create survey.Validator for string with regex
func validateStringWithRegex(val interface{}, regex string, errorString string) error {
	// the reflect value of the result
	value := reflect.ValueOf(val)

	// if the value passed is not a string
	if value.Kind() != reflect.String {
		return errors.New("value must be a string")
	}

	// if value doesn't match regex, return invalid arn
	ok, err := regexp.MatchString(regex, value.String())
	if err != nil {
		return errors.Wrap(err, "failed to validate input")
	}

	if !ok {
		return errors.New(errorString)
	}

	return nil
}

// survey.Validator for aws ARNs
//
// This isn't service/type specific but rather just validates that an ARN was entered that matches valid ARN formats
func validateAwsArnFormat(val interface{}) error {
	return validateStringWithRegex(val, AwsArnRegex, "invalid arn supplied")
}

// survey.Validator for aws region
func validateAwsRegion(val interface{}) error {
	return validateStringWithRegex(val, AwsRegionRegex, "invalid region name supplied")
}

// survey.Validator for aws profile
func validateAwsProfile(val interface{}) error {
	return validateStringWithRegex(val, fmt.Sprintf(`^%s$`, AwsProfileRegex), "invalid profile name supplied")
}

func promptAwsCtQuestions(config *aws.GenerateAwsTfConfigurationArgs, extraState *AwsGenerateCommandExtraState) error {
	// Only ask these questions if configure cloudtrail is true
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionConsolidatedCloudtrail, Default: config.ConsolidatedCloudtrail},
			Response: &config.ConsolidatedCloudtrail,
		},
		{
			Prompt:   &survey.Confirm{Message: QuestionUseExistingCloudtrail, Default: extraState.UseExistingCloudtrail},
			Response: &extraState.UseExistingCloudtrail,
		},
		{
			Prompt:   &survey.Input{Message: QuestionCloudtrailExistingBucketArn, Default: config.ExistingCloudtrailBucketArn},
			Checks:   []*bool{&extraState.UseExistingCloudtrail},
			Required: true,
			Opts:     []survey.AskOpt{survey.WithValidator(validateAwsArnFormat)},
			Response: &config.ExistingCloudtrailBucketArn,
		},
	}, config.Cloudtrail); err != nil {
		return err
	}

	// If a new bucket is to be created; should the force destroy bit be set?
	newBucket := config.ExistingCloudtrailBucketArn == ""
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Message: QuestionForceDestroyS3Bucket, Default: config.ForceDestroyS3Bucket},
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
			Prompt:   &survey.Input{Message: QuestionExistingIamRoleName, Default: config.ExistingIamRole.Name},
			Response: &config.ExistingIamRole.Name,
			Opts:     []survey.AskOpt{survey.WithValidator(survey.Required)},
		},
		{
			Prompt:   &survey.Input{Message: QuestionExistingIamRoleArn, Default: config.ExistingIamRole.Arn},
			Response: &config.ExistingIamRole.Arn,
			Opts:     []survey.AskOpt{survey.WithValidator(survey.Required), survey.WithValidator(validateAwsArnFormat)},
		},
		{
			Prompt:   &survey.Input{Message: QuestionExistingIamRoleExtId, Default: config.ExistingIamRole.ExternalId},
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
			Message: QuestionPrimaryAwsAccountProfile,
			Default: config.AwsProfile,
		},
		Opts:     []survey.AskOpt{survey.WithValidator(validateAwsProfile)},
		Response: &config.AwsProfile,
		Required: true,
	}); err != nil {
		return nil
	}

	// If there are existing sub accounts configured (i.e., from the CLI) display them and ask if they want to add more
	if len(config.SubAccounts) > 0 {
		subAccountListing := []string{}
		for _, account := range config.SubAccounts {
			subAccountListing = append(subAccountListing, fmt.Sprintf("%s:%s", account.AwsProfile, account.AwsRegion))
		}

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Prompt: &survey.Confirm{
				Message: fmt.Sprintf(
					QuestionSubAccountReplace,
					strings.Trim(strings.Join(strings.Fields(fmt.Sprint(subAccountListing)), ", "), "[]"),
				),
			},
			Response: &askAgain}); err != nil {
			return err
		}
	}

	// For each account to add, collect the aws profile and region to use
	for askAgain {
		var accountProfileName string
		var accountProfileRegion string
		accountQuestions := []SurveyQuestionWithValidationArgs{
			{
				Prompt:   &survey.Input{Message: QuestionSubAccountProfileName},
				Opts:     []survey.AskOpt{survey.WithValidator(validateAwsProfile)},
				Required: true,
				Response: &accountProfileName,
			},
			{
				Prompt:   &survey.Input{Message: QuestionSubAccountRegion},
				Opts:     []survey.AskOpt{survey.WithValidator(validateAwsRegion)},
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
			Prompt:   &survey.Confirm{Message: QuestionSubAccountAddMore},
			Response: &askAgain}); err != nil {
			return err
		}
	}

	// If we created new accounts, re-write config
	if len(accountDetails) > 0 {
		config.SubAccounts = accountDetails
	}

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
	if err := validateOutputLocation(value.String()); err != nil {
		return err
	}

	return nil
}

func promptCustomizeOutputLocation(config *aws.GenerateAwsTfConfigurationArgs, extraState *AwsGenerateCommandExtraState) error {
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionCustomizeOutputLocation, Default: extraState.Output},
		Response: &extraState.Output,
		Opts:     []survey.AskOpt{survey.WithValidator(validPathExists)},
		Required: true,
	}); err != nil {
		return err
	}

	return nil
}

func askAdvancedOptions(config *aws.GenerateAwsTfConfigurationArgs, extraState *AwsGenerateCommandExtraState) error {
	answer := ""

	// Prompt for options
	for answer != AdvancedOptDone {
		// Construction of this slice is a bit strange at first look, but the reason for that is because we have to do string
		// validation to know which option was selected due to how survey works; and doing it by index (also supported) is
		// difficult when the options are dynamic (which they are)
		//
		// Only ask about more accounts if consolidated cloudtrail is setup (matching scenarios doc)
		options := []string{AdvancedOptCloudTrail, AdvancedOptIamRole}
		if config.ConsolidatedCloudtrail {
			options = append(options, AdvancedOptAwsAccounts)
		}
		options = append(options, AdvancedOptLocation, AdvancedOptDone)
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
		case AdvancedOptCloudTrail:
			if err := promptAwsCtQuestions(config, extraState); err != nil {
				return err
			}
		case AdvancedOptIamRole:
			if err := promptAwsExistingIamQuestions(config); err != nil {
				return err
			}
		case AdvancedOptAwsAccounts:
			if err := promptAwsAdditionalAccountQuestions(config); err != nil {
				return err
			}
		case AdvancedOptLocation:
			if err := promptCustomizeOutputLocation(config, extraState); err != nil {
				return err
			}
		}

		// Re-prompt if not done
		innerAskAgain := true
		if answer == AdvancedOptDone {
			innerAskAgain = false
		}

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Checks:   []*bool{&innerAskAgain},
			Prompt:   &survey.Confirm{Message: QuestionAnotherAdvancedOpt, Default: false},
			Response: &innerAskAgain,
		}); err != nil {
			return err
		}

		if !innerAskAgain {
			answer = AdvancedOptDone
		}
	}

	return nil
}

func configOrCloudtrailEnabled(config *aws.GenerateAwsTfConfigurationArgs) *bool {
	cloudtrailOrConfigEnabled := config.Cloudtrail || config.Config
	return &cloudtrailOrConfigEnabled
}

// entry point for launching a survey to build out the required generation parameters
func promptAwsGenerate(
	config *aws.GenerateAwsTfConfigurationArgs,
	existingIam *aws.ExistingIamRoleDetails,
	extraState *AwsGenerateCommandExtraState,
) error {
	// Cache for later use if generation is abandon and in interactive mode
	if cli.InteractiveMode() {
		defer cli.WriteAssetToCache(CachedAssetIacParams, time.Now().Add(time.Hour*1), config)
		defer cli.WriteAssetToCache(CachedAssetAwsExtraState, time.Now().Add(time.Hour*1), extraState)
	}

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
				Prompt:   &survey.Confirm{Message: QuestionEnableConfig, Default: config.Config},
				Response: &config.Config,
			},
			{
				Prompt:   &survey.Confirm{Message: QuestionEnableCloudtrail, Default: config.Cloudtrail},
				Response: &config.Cloudtrail,
			},
		}); err != nil {
		return err
	}

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionAwsRegion, Default: config.AwsRegion},
		Response: &config.AwsRegion,
		Opts:     []survey.AskOpt{survey.WithValidator(survey.Required), survey.WithValidator(validateAwsRegion)},
		Checks:   []*bool{configOrCloudtrailEnabled(config)},
	}); err != nil {
		return err
	}

	// Validate one of config or cloudtrail was enabled; otherwise error out
	if !config.Config && !config.Cloudtrail {
		return errors.New("must enable cloudtrail or config")
	}

	// Find out if the customer wants to specify more advanced features
	askAdvanced := false
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Message: QuestionConfigAdvanced, Default: askAdvanced},
		Response: &askAdvanced,
	}); err != nil {
		return err
	}

	// Keep prompting for advanced options until the say done
	if askAdvanced {
		if err := askAdvancedOptions(config, extraState); err != nil {
			return err
		}
	}

	return nil
}
