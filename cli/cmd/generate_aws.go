package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/lwgenerate/aws"
	"github.com/pkg/errors"
)

var (
	// Define question text here so they can be reused in testing
	QuestionAwsEnableConfig             = "Enable configuration integration?"
	QuestionCustomizeConfigName         = "Customize Config integration name?"
	QuestionConfigName                  = "Specify name of config integration (optional)"
	QuestionEnableCloudtrail            = "Enable CloudTrail integration?"
	QuestionCloudtrailName              = "Specify name of cloudtrail integration (optional)"
	QuestionAwsRegion                   = "Specify the AWS region to be used by CloudTrail, SNS, and S3:"
	QuestionConsolidatedCloudtrail      = "Use consolidated CloudTrail?"
	QuestionUseExistingCloudtrail       = "Use an existing CloudTrail?"
	QuestionCloudtrailExistingBucketArn = "Specify an existing bucket ARN used for CloudTrail logs:"
	QuestionForceDestroyS3Bucket        = "Should the new S3 bucket have force destroy enabled?"
	QuestionExistingIamRoleName         = "Specify an existing IAM role name for CloudTrail access:"
	QuestionExistingIamRoleArn          = "Specify an existing IAM role ARN for CloudTrail access:"
	QuestionExistingIamRoleExtID        = "Specify the external ID to be used with the existing IAM role:"
	QuestionPrimaryAwsAccountProfile    = "Before adding sub-accounts, your primary AWS account profile name must be set;" +
		" which profile should the main account use?"
	QuestionSubAccountProfileName      = "Supply the profile name for this additional AWS account:"
	QuestionSubAccountRegion           = "What region should be used for this account?"
	QuestionSubAccountAddMore          = "Add another AWS account?"
	QuestionSubAccountReplace          = "Currently configured AWS sub-accounts: %s, replace?"
	QuestionAwsConfigAdvanced          = "Configure advanced integration options?"
	QuestionAwsAnotherAdvancedOpt      = "Configure another advanced integration option"
	QuestionAwsCustomizeOutputLocation = "Provide the location for the output to be written:"

	// S3 Bucket Questions
	QuestionBucketEnableEncryption = "Enable S3 bucket encryption when creating bucket"
	QuestionBucketSseKeyArn        = "Specify existing KMS encryption key arn for S3 bucket (optional)"
	QuestionBucketName             = "Specify name when creating S3 bucket (optional)"

	// SNS Topic Questions
	QuestionsUseExistingSNSTopic = "Use an existing SNS topic?"
	QuestionSnsTopicArn          = "Specify existing SNS topic arn"
	QuestionSnsEnableEncryption  = "Enable encryption on SNS topic when creating?"
	QuestionSnsEncryptionKeyArn  = "Specify existing KMS encryption key arn for SNS topic (optional)"
	QuestionSnsTopicName         = "Specify SNS topic name if creating new one (optional)"

	// SQS Queue Questions
	QuestionSqsEnableEncryption = "Enable encryption on SQS queue when creating"
	QuestionSqsEncryptionKeyArn = "Specify existing KMS encryption key arn for SQS queue (optional)"
	QuestionSqsQueueName        = "Specify SQS queue name when creating (optional)"

	// select options
	AwsAdvancedOptDone     = "Done"
	AdvancedOptCloudTrail  = "Additional CloudTrail options"
	AdvancedOptIamRole     = "Configure Lacework integration with an existing IAM role"
	AdvancedOptAwsAccounts = "Add additional AWS Accounts to Lacework"
	AwsAdvancedOptLocation = "Customize output location"

	// AwsArnRegex original source: https://regex101.com/r/pOfxYN/1
	AwsArnRegex = `^arn:(?P<Partition>[^:\n]*):(?P<Service>[^:\n]*):(?P<Region>[^:\n]*):(?P<AccountID>[^:\n]*):(?P<Ignore>(?P<ResourceType>[^:\/\n]*)[:\/])?(?P<Resource>.*)$`
	// AwsRegionRegex regex used for validating region input; note intentionally does not match gov cloud
	AwsRegionRegex  = `(us|ap|ca|cn|eu|sa)-(central|(north|south)?(east|west)?)-\d`
	AwsProfileRegex = `([A-Za-z_0-9-]+)`

	GenerateAwsCommandState      = &aws.GenerateAwsTfConfigurationArgs{}
	GenerateAwsExistingRoleState = &aws.ExistingIamRoleDetails{}
	GenerateAwsCommandExtraState = &AwsGenerateCommandExtraState{}
	ValidateSubAccountFlagRegex  = fmt.Sprintf(`%s:%s`, AwsProfileRegex, AwsRegionRegex)
	CachedAwsAssetIacParams      = "iac-aws-generate-params"
	CachedAssetAwsExtraState     = "iac-aws-extra-state"

	// aws command is used to generate TF code for aws
	generateAwsTfCommand = &cobra.Command{
		Use:   "aws",
		Short: "Generate and/or execute Terraform code for AWS integration",
		Long: `Use this command to generate Terraform code for deploying Lacework into an AWS environment.

By default, this command interactively prompts for the required information to setup the new cloud account.
In interactive mode, this command will:

* Prompt for the required information to setup the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version is verified as compatible for use
	* If Terraform is not installed, or the version installed is not compatible, a new version will be installed into a temporary location
	* Once Terraform is detected or installed, Terraform plan will be executed
	* The command will prompt with the outcome of the plan and allow to view more details or continue with Terraform apply
	* If confirmed, Terraform apply will be run, completing the setup of the cloud account

This command can also be run in noninteractive mode.
See help output for more details on the parameter value(s) required for Terraform code generation.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generate TF Code
			cli.StartProgress("Generating Terraform Code...")

			// Explicitly set Lacework profile if it was passed in main args
			if cli.Profile != "default" {
				GenerateAwsCommandState.LaceworkProfile = cli.Profile
			}

			// Setup modifiers for NewTerraform constructor
			mods := []aws.AwsTerraformModifier{
				aws.WithAwsProfile(GenerateAwsCommandState.AwsProfile),
				aws.WithLaceworkProfile(GenerateAwsCommandState.LaceworkProfile),
				aws.ExistingCloudtrailBucketArn(GenerateAwsCommandState.ExistingCloudtrailBucketArn),
				aws.ExistingSnsTopicArn(GenerateAwsCommandState.ExistingSnsTopicArn),
				aws.WithSubaccounts(GenerateAwsCommandState.SubAccounts...),
				aws.UseExistingIamRole(GenerateAwsCommandState.ExistingIamRole),
				aws.WithConfigName(GenerateAwsCommandState.ConfigName),
				aws.WithCloudtrailName(GenerateAwsCommandState.CloudtrailName),
				aws.WithBucketName(GenerateAwsCommandState.BucketName),
				aws.WithBucketEncryptionEnabled(GenerateAwsCommandState.BucketEncryptionEnabled),
				aws.WithBucketSSEKeyArn(GenerateAwsCommandState.BucketSseKeyArn),
				aws.WithSnsTopicName(GenerateAwsCommandState.SnsTopicName),
				aws.WithSnsEncryptionEnabled(GenerateAwsCommandState.SnsEncryptionEnabled),
				aws.WithSnsEncryptionKeyArn(GenerateAwsCommandState.SnsEncryptionKeyArn),
				aws.WithSqsQueueName(GenerateAwsCommandState.SqsQueueName),
				aws.WithSqsEncryptionEnabled(GenerateAwsCommandState.SqsEncryptionEnabled),
				aws.WithSqsEncryptionKeyArn(GenerateAwsCommandState.SqsEncryptionKeyArn),
			}

			if GenerateAwsCommandState.ForceDestroyS3Bucket {
				mods = append(mods, aws.EnableForceDestroyS3Bucket())
			}

			if GenerateAwsCommandState.ConsolidatedCloudtrail {
				mods = append(mods, aws.UseConsolidatedCloudtrail())
			}

			// Create new struct
			data := aws.NewTerraform(
				GenerateAwsCommandState.AwsRegion,
				GenerateAwsCommandState.Config,
				GenerateAwsCommandState.Cloudtrail,
				mods...)

			// Generate
			hcl, err := data.Generate()
			cli.StopProgress()

			if err != nil {
				return errors.Wrap(err, "failed to generate terraform code")
			}

			// Write-out generated code to location specified
			dirname, location, err := writeGeneratedCodeToLocation(cmd, hcl, "aws")
			if err != nil {
				return err
			}

			// Prompt to execute
			err = SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
				Prompt:   &survey.Confirm{Default: GenerateAwsCommandExtraState.TerraformApply, Message: QuestionRunTfPlan},
				Response: &GenerateAwsCommandExtraState.TerraformApply,
			})

			if err != nil {
				return errors.Wrap(err, "failed to prompt for terraform execution")
			}

			// Execute
			locationDir := filepath.Dir(location)
			if GenerateAwsCommandExtraState.TerraformApply {
				// Execution pre-run check
				err := executionPreRunChecks(dirname, locationDir, "aws")
				if err != nil {
					return err
				}
			}

			// Output where code was generated
			if !GenerateAwsCommandExtraState.TerraformApply {
				cli.OutputHuman(provideGuidanceAfterExit(false, false, locationDir, "terraform"))
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			// Validate output location is OK if supplied
			dirname, err := cmd.Flags().GetString("output")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateOutputLocation(dirname); err != nil {
				return err
			}

			// Validate aws profile, if passed
			profile, err := cmd.Flags().GetString("aws_profile")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsProfile(profile); profile != "" && err != nil {
				return err
			}

			// Validate aws region, if passed
			region, err := cmd.Flags().GetString("aws_region")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsRegion(region); region != "" && err != nil {
				return err
			}

			// Validate cloudtrail bucket arn, if passed
			arn, err := cmd.Flags().GetString("existing_bucket_arn")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsArnFormat(arn); arn != "" && err != nil {
				return err
			}

			// Validate SNS Topic Arn if passed
			arn, err = cmd.Flags().GetString("existing_sns_topic_arn")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsArnFormat(arn); arn != "" && err != nil {
				return err
			}

			// Load any cached inputs if interactive
			if cli.InteractiveMode() {
				cachedOptions := &aws.GenerateAwsTfConfigurationArgs{}
				iacParamsExpired := cli.ReadCachedAsset(CachedAwsAssetIacParams, &cachedOptions)
				if iacParamsExpired {
					cli.Log.Debug("loaded previously set values for AWS iac generation")
				}

				extraState := &AwsGenerateCommandExtraState{}
				extraStateParamsExpired := cli.ReadCachedAsset(CachedAssetAwsExtraState, &extraState)
				if extraStateParamsExpired {
					cli.Log.Debug("loaded previously set values for AWS iac generation (extra state)")
				}

				// Determine if previously cached options exists; prompt user if they'd like to continue
				answer := false
				if !iacParamsExpired || !extraStateParamsExpired {
					if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
						Prompt:   &survey.Confirm{Message: QuestionUsePreviousCache, Default: false},
						Response: &answer,
					}); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
				}

				// If the user decides NOT to use the previous values; we won't load them.  However, every time the command runs
				// we are going to write out new cached values, so if they run it - bail out - and run it again they'll get
				// re-prompted.
				if answer {
					// Merge cached inputs to current options (current options win)
					if err := mergo.Merge(GenerateAwsCommandState, cachedOptions); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
					if err := mergo.Merge(GenerateAwsCommandExtraState, extraState); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
				}
			}

			// Parse passed in subaccounts (if any)
			if len(GenerateAwsCommandExtraState.AwsSubAccounts) > 0 {
				// Validate the format of supplied values is correct
				if err := validateAwsSubAccounts(GenerateAwsCommandExtraState.AwsSubAccounts); err != nil {
					return err
				}

				awsSubAccounts := []aws.AwsSubAccount{}
				for _, account := range GenerateAwsCommandExtraState.AwsSubAccounts {
					accountDetails := strings.Split(account, ":")
					awsSubAccounts = append(awsSubAccounts, aws.NewAwsSubAccount(accountDetails[0], accountDetails[1]))
				}
				GenerateAwsCommandState.SubAccounts = awsSubAccounts
			}

			// Collect and/or confirm parameters
			err = promptAwsGenerate(GenerateAwsCommandState, GenerateAwsExistingRoleState, GenerateAwsCommandExtraState)
			if err != nil {
				return errors.Wrap(err, "collecting/confirming parameters")
			}

			return nil
		},
	}
)

type AwsGenerateCommandExtraState struct {
	Output                string
	UseExistingCloudtrail bool
	UseExistingSNSTopic   bool
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

func initGenerateAwsTfCommandFlags() {
	// add flags to sub commands
	// TODO Share the help with the interactive generation
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.Cloudtrail,
		"cloudtrail",
		false,
		"enable cloudtrail integration")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.CloudtrailName,
		"cloudtrail_name",
		"",
		"specify name of cloudtrail integration")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.Config,
		"config",
		false,
		"enable config integration")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ConfigName,
		"config_name",
		"",
		"specify name of config integration")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.AwsRegion,
		"aws_region",
		"",
		"specify aws region")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.AwsProfile,
		"aws_profile",
		"",
		"specify aws profile")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.BucketEncryptionEnabled,
		"bucket_encryption_enabled",
		true,
		"enable S3 bucket encryption when creating bucket")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.BucketName,
		"bucket_name",
		"",
		"specify bucket name when creating bucket")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.BucketSseKeyArn,
		"bucket_sse_key_arn",
		"",
		"specify existing KMS encryption key arn for bucket")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ExistingCloudtrailBucketArn,
		"existing_bucket_arn",
		"",
		"specify existing cloudtrail S3 bucket ARN")
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
		"specify existing SNS topic arn")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.ConsolidatedCloudtrail,
		"consolidated_cloudtrail",
		false,
		"use consolidated trail")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.ForceDestroyS3Bucket,
		"force_destroy_s3",
		false,
		"enable force destroy S3 bucket")
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
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandExtraState.Output,
		"output",
		"",
		"location to write generated content (default is ~/lacework/aws)",
	)
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.SnsEncryptionEnabled,
		"sns_encryption_enabled",
		true,
		"enable encryption on SNS topic when creating one")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.SnsEncryptionKeyArn,
		"sns_encryption_key_arn",
		"",
		"specify existing KMS encryption key arn for SNS topic")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.SnsTopicName,
		"sns_topic_name",
		"",
		"specify SNS topic name if creating new one")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.SqsEncryptionEnabled,
		"sqs_encryption_enabled",
		true,
		"enable encryption on SQS queue when creating")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.SqsEncryptionKeyArn,
		"sqs_encryption_key_arn",
		"",
		"specify existing KMS encryption key arn for SQS queue")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.SqsQueueName,
		"sqs_queue_name",
		"",
		"specify SQS queue name if creating new one")
}

// survey.Validator for aws ARNs
//
// This isn't service/type specific but rather just validates that an ARN was entered that matches valid ARN formats
func validateAwsArnFormat(val interface{}) error {
	return validateStringWithRegex(val, AwsArnRegex, "invalid arn supplied")
}

// Validate AWS Arn only if a value is set, this can be used for optional ARN cofiguration
func validateOptionalAwsArnFormat(val interface{}) error {
	if val.(string) != "" {
		return validateAwsArnFormat(val)
	}
	return nil
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
	}, config.Cloudtrail); err != nil {
		return err
	}
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionCloudtrailName, Default: config.CloudtrailName},
			Checks:   []*bool{existingCloudTrailDisabled(extraState)},
			Response: &config.CloudtrailName,
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
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionForceDestroyS3Bucket, Default: config.ForceDestroyS3Bucket},
			Response: &config.ForceDestroyS3Bucket,
			Checks:   []*bool{&config.Cloudtrail, &newBucket},
		},
		// If new bucket created, allow user to optionally name the bucket
		{
			Prompt:   &survey.Input{Message: QuestionBucketName, Default: config.BucketName},
			Response: &config.BucketName,
			Checks:   []*bool{&config.Cloudtrail, &newBucket},
		},
		// If new bucket created, should this have encryption enabled
		{
			Prompt:   &survey.Confirm{Message: QuestionBucketEnableEncryption, Default: config.BucketEncryptionEnabled},
			Response: &config.BucketEncryptionEnabled,
			Checks:   []*bool{&config.Cloudtrail, &newBucket},
		},
		// Allow the user to set the SSE Key ARN if required
		{
			Prompt:   &survey.Input{Message: QuestionBucketSseKeyArn, Default: config.BucketSseKeyArn},
			Response: &config.BucketSseKeyArn,
			Opts:     []survey.AskOpt{survey.WithValidator(validateOptionalAwsArnFormat)},
			Checks:   []*bool{&config.Cloudtrail, &newBucket, &config.BucketEncryptionEnabled},
		},
	}, config.Cloudtrail); err != nil {
		return err
	}
	// SNS Options
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionsUseExistingSNSTopic, Default: extraState.UseExistingSNSTopic},
			Response: &extraState.UseExistingSNSTopic,
		},
		{
			Prompt:   &survey.Input{Message: QuestionSnsTopicArn, Default: config.ExistingSnsTopicArn},
			Checks:   []*bool{&extraState.UseExistingSNSTopic},
			Required: true,
			Opts:     []survey.AskOpt{survey.WithValidator(validateAwsArnFormat)},
			Response: &config.ExistingSnsTopicArn,
		},
	}, config.Cloudtrail); err != nil {
		return err
	}
	// If a new SNS Topic is to be created
	newTopic := config.ExistingSnsTopicArn == ""
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		// If new topic created, allow user to optionally name the topic
		{
			Prompt:   &survey.Input{Message: QuestionSnsTopicName, Default: config.SnsTopicName},
			Response: &config.SnsTopicName,
			Checks:   []*bool{&config.Cloudtrail, &newTopic},
		},
		// If new bucket created, should this have encryption enabled
		{
			Prompt:   &survey.Confirm{Message: QuestionSnsEnableEncryption, Default: config.SnsEncryptionEnabled},
			Response: &config.SnsEncryptionEnabled,
			Checks:   []*bool{&config.Cloudtrail, &newTopic},
		},
		// Allow the user to set the SSE Key ARN if required
		{
			Prompt:   &survey.Input{Message: QuestionSnsEncryptionKeyArn, Default: config.SnsEncryptionKeyArn},
			Response: &config.SnsEncryptionKeyArn,
			Opts:     []survey.AskOpt{survey.WithValidator(validateOptionalAwsArnFormat)},
			Checks:   []*bool{&config.Cloudtrail, &newTopic, &config.SnsEncryptionEnabled},
		},
	}, config.Cloudtrail); err != nil {
		return err
	}
	// SQS Options
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		// New queue created, allow user to optionally name the queue
		{
			Prompt:   &survey.Input{Message: QuestionSqsQueueName, Default: config.SqsQueueName},
			Response: &config.SqsQueueName,
			Checks:   []*bool{&config.Cloudtrail},
		},
		// New queue created, should this have encryption enabled
		{
			Prompt:   &survey.Confirm{Message: QuestionSqsEnableEncryption, Default: config.SqsEncryptionEnabled},
			Response: &config.SqsEncryptionEnabled,
			Checks:   []*bool{&config.Cloudtrail},
		},
		// Allow the user to set the SSE Key ARN if required
		{
			Prompt:   &survey.Input{Message: QuestionSqsEncryptionKeyArn, Default: config.SqsEncryptionKeyArn},
			Response: &config.SqsEncryptionKeyArn,
			Opts:     []survey.AskOpt{survey.WithValidator(validateOptionalAwsArnFormat)},
			Checks:   []*bool{&config.Cloudtrail, &config.SqsEncryptionEnabled},
		},
	}, config.Cloudtrail); err != nil {
		return err
	}

	return nil
}

func existingCloudTrailDisabled(extraState *AwsGenerateCommandExtraState) *bool {
	existingCloudTrailDisabled := !extraState.UseExistingCloudtrail
	return &existingCloudTrailDisabled
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
			Prompt:   &survey.Input{Message: QuestionExistingIamRoleExtID, Default: config.ExistingIamRole.ExternalId},
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

func promptCustomizeAwsOutputLocation(extraState *AwsGenerateCommandExtraState) error {
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionAwsCustomizeOutputLocation, Default: extraState.Output},
		Response: &extraState.Output,
		Opts:     []survey.AskOpt{survey.WithValidator(validPathExists)},
		Required: true,
	}); err != nil {
		return err
	}

	return nil
}

func promptCustomizeConfigOptions(config *aws.GenerateAwsTfConfigurationArgs) error {
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionConfigName, Default: config.ConfigName},
		Checks:   []*bool{&config.Config},
		Response: &config.ConfigName,
	}); err != nil {
		return err
	}

	return nil
}

func askAdvancedAwsOptions(config *aws.GenerateAwsTfConfigurationArgs, extraState *AwsGenerateCommandExtraState) error {
	answer := ""

	// Prompt for options
	for answer != AwsAdvancedOptDone {
		// Construction of this slice is a bit strange at first look, but the reason for that is because we have to do string
		// validation to know which option was selected due to how survey works; and doing it by index (also supported) is
		// difficult when the options are dynamic (which they are)
		//
		// Only ask about more accounts if consolidated cloudtrail is setup (matching scenarios doc)
		// Scenario document suggests that this is no longer the case and that
		// we can have other accounts even if we only have Config integration (Scenario 7)
		var options []string

		// Determine if user specified name for Config is potentially required
		if config.Config {
			options = append(options, QuestionCustomizeConfigName)
		}

		// Only show Advanced CloudTrail options if CloudTrail integration is set to true
		if config.Cloudtrail {
			options = append(options, AdvancedOptCloudTrail)
		}

		// Scenario document suggests that this is no longer the case and that
		// we can have other accounts even if we only have Config integration (Scenario 7)
		//if config.ConsolidatedCloudtrail {
		options = append(options, AdvancedOptAwsAccounts)
		//}

		options = append(options, AdvancedOptIamRole, AwsAdvancedOptLocation, AwsAdvancedOptDone)
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Prompt: &survey.Select{
				Message: "Which options would you like to configure?",
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
		case AwsAdvancedOptLocation:
			if err := promptCustomizeAwsOutputLocation(extraState); err != nil {
				return err
			}
		case QuestionCustomizeConfigName:
			if err := promptCustomizeConfigOptions(config); err != nil {
				return err
			}
		}

		// Re-prompt if not done
		innerAskAgain := true
		if answer == AwsAdvancedOptDone {
			innerAskAgain = false
		}

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Checks:   []*bool{&innerAskAgain},
			Prompt:   &survey.Confirm{Message: QuestionAwsAnotherAdvancedOpt, Default: false},
			Response: &innerAskAgain,
		}); err != nil {
			return err
		}

		if !innerAskAgain {
			answer = AwsAdvancedOptDone
		}
	}

	return nil
}

func configOrCloudtrailEnabled(config *aws.GenerateAwsTfConfigurationArgs) *bool {
	cloudtrailOrConfigEnabled := config.Cloudtrail || config.Config
	return &cloudtrailOrConfigEnabled
}

func awsConfigIsEmpty(g *aws.GenerateAwsTfConfigurationArgs) bool {
	return !g.Cloudtrail &&
		!g.Config &&
		!g.ConsolidatedCloudtrail &&
		g.AwsProfile == "" &&
		g.AwsRegion == "" &&
		g.ExistingCloudtrailBucketArn == "" &&
		g.ExistingIamRole == nil &&
		g.ExistingSnsTopicArn == "" &&
		g.LaceworkProfile == "" &&
		!g.ForceDestroyS3Bucket &&
		g.SubAccounts == nil
}

func writeAwsGenerationArgsCache(a *aws.GenerateAwsTfConfigurationArgs) {
	if !awsConfigIsEmpty(a) {
		// If ExistingIamRole is partially set, don't write this to cache; the values won't work when loaded
		if a.ExistingIamRole.IsPartial() {
			a.ExistingIamRole = nil
		}
		cli.WriteAssetToCache(CachedAwsAssetIacParams, time.Now().Add(time.Hour*1), a)
	}
}

// entry point for launching a survey to build out the required generation parameters
func promptAwsGenerate(
	config *aws.GenerateAwsTfConfigurationArgs,
	existingIam *aws.ExistingIamRoleDetails,
	extraState *AwsGenerateCommandExtraState,
) error {
	// Cache for later use if generation is abandon and in interactive mode
	if cli.InteractiveMode() {
		defer writeAwsGenerationArgsCache(config)
		defer extraState.writeCache()
	}

	// Set ExistingIamRole details, if provided as cli flags; otherwise don't initialize
	if existingIam.Arn != "" ||
		existingIam.Name != "" ||
		existingIam.ExternalId != "" {
		config.ExistingIamRole = existingIam
	}

	// These are the core questions that should be asked.  Region required for provider block
	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Prompt:   &survey.Confirm{Message: QuestionAwsEnableConfig, Default: config.Config},
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
		Prompt:   &survey.Confirm{Message: QuestionAwsConfigAdvanced, Default: askAdvanced},
		Response: &askAdvanced,
	}); err != nil {
		return err
	}

	// Keep prompting for advanced options until the say done
	if askAdvanced {
		if err := askAdvancedAwsOptions(config, extraState); err != nil {
			return err
		}
	}

	return nil
}
