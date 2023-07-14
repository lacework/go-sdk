package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lacework/go-sdk/lwgenerate/aws_controltower"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

var (
	QuestionAwsControlTowerCoreS3Bucket            = "Provide the Arn of the S3 Bucket for consolidated CloudTrail:"
	QuestionAwsControlTowerCoreSnsTopic            = "Provide the Arn of the SNS Topic:"
	QuestionAwsControlTowerCoreLogProfile          = "Provide the aws profile of the 'log_archive' account:"
	QuestionAwsControlTowerCoreLogRegion           = "Provide the aws region of the 'log_archive' account:"
	QuestionAwsControlTowerCoreAuditProfile        = "Provide the aws profile of the 'audit' account:"
	QuestionAwsControlTowerCoreAuditRegion         = "Provide the aws region of the 'audit' account:"
	QuestionAwsControlTowerConfigureAdvanced       = "Configure advanced integration options?"
	QuestionAwsControlTowerCustomizeOutputLocation = "Provide the location for the output to be written:"

	ControlTowerConfigureExistingIamRoleOpt                 = "Configure existing Iam Role?"
	QuestionAwsControlTowerCoreIamRoleName                  = "Specify Existing Iam Role name:"
	QuestionAwsControlTowerCoreIamRoleArn                   = "Specify Existing Iam Arn:"
	QuestionAwsControlTowerCoreIamRoleExternalID            = "Specify Existing Iam Role external ID:"
	ControlTowerIntegrationNameOpt                          = "Customize integration name?"
	QuestionControlTowerIntegrationName                     = "Specify a custom integration name:"
	ControlTowerIntegrationPrefixOpt                        = "Customize resource prefix name?"
	QuestionControlTowerPrefix                              = "Specify a prefix name for resources:"
	ControlTowerIntegrationSqsOpt                           = "Customize sqs queue name?"
	QuestionControlTowerSqsQueueName                        = "Specify a name for sqs queue:"
	QuestionControlTowerOrgAccountMappingsLWDefaultAccount  = "Specify org account mappings default Lacework account:"
	QuestionControlTowerOrgAccountMappingAnotherAdvancedOpt = "Configure another org account mapping?"
	QuestionControlTowerOrgAccountMappingsLWAccount         = "Specify lacework account: "
	QuestionControlTowerOrgAccountMappingsAwsAccounts       = "Specify aws accounts:"
	ControlTowerAdvancedOptLocation                         = "Customize output location"
	ControlTowerAdvancedOptMappings                         = "Configure Org Account Mappings"
	QuestionControlTowerAnotherAdvancedOpt                  = "Configure another advanced integration option?"
	ControlTowerAdvancedOptDone                             = "Done"

	GenerateAwsControlTowerCommandState      = &aws_controltower.GenerateAwsControlTowerTfConfigurationArgs{}
	GenerateAwsControlTowerCommandExtraState = &AwsControlTowerGenerateCommandExtraState{}
	CachedAssetAwsControlTowerIacParams      = "iac-aws-controltower-generate-params"
	CachedAssetAwsControlTowerExtraState     = "iac-aws-controltower-extra-state"

	generateAwsControlTowerTfCommand = &cobra.Command{
		Use:   "controltower",
		Short: "Generate and/or execute Terraform code for ControlTower integration",
		Long: `Use this command to generate Terraform code for deploying Lacework with Aws Cloudtrail and
ControlTower.

By default, this command interactively prompts for the required information to set up the new cloud account.
In interactive mode, this command will:

* Prompt for the required information to set up the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version is verified as compatible for use
  * If Terraform is not installed, or the version installed is not compatible, a new
    version will be installed into a temporary location
  * Once Terraform is detected or installed, the Terraform plan is executed
  * The command prompts you with the outcome of the plan and allows you to view more
    details or continue with Terraform apply
  * If confirmed, Terraform apply runs, completing the setup of the cloud account

This command can also be run in noninteractive mode.
See help output for more details on the parameter values required for Terraform code generation.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generate TF Code
			cli.StartProgress("Generating Terraform Code...")

			// Explicitly set Lacework profile if it was passed in main args
			if cli.Profile != "default" {
				GenerateAwsControlTowerCommandState.LaceworkProfile = cli.Profile
			}

			// Setup modifiers for NewTerraform constructor
			mods := []aws_controltower.AwsControlTowerTerraformModifier{
				aws_controltower.WithLaceworkAccountID(GenerateAwsControlTowerCommandState.LaceworkAccountID),
				aws_controltower.WithSubaccounts(GenerateAwsControlTowerCommandState.SubAccounts...),
				aws_controltower.WithSqsQueueName(GenerateAwsControlTowerCommandState.SqsQueueName),
				aws_controltower.WithPrefix(GenerateAwsControlTowerCommandState.Prefix),
				aws_controltower.WithCrossAccountPolicyName(GenerateAwsControlTowerCommandState.CrossAccountPolicyName),
				aws_controltower.WithSubaccounts(GenerateAwsControlTowerCommandState.SubAccounts...),
			}

			if useExistingIamRole(GenerateAwsControlTowerCommandState) {
				mods = append(mods, aws_controltower.WithExisitingIamRole(
					GenerateAwsControlTowerCommandState.IamRoleArn,
					GenerateAwsControlTowerCommandState.IamRoleName,
					GenerateAwsControlTowerCommandState.IamRoleExternalID,
				))
			}

			if !GenerateAwsControlTowerCommandState.OrgAccountMappings.IsEmpty() {
				mods = append(mods, aws_controltower.WithOrgAccountMappings(GenerateAwsControlTowerCommandState.OrgAccountMappings))
			}

			data := aws_controltower.NewTerraform(
				GenerateAwsControlTowerCommandState.S3BucketArn,
				GenerateAwsControlTowerCommandState.SNSTopicArn,
				mods...)

			// Generate
			hcl, err := data.Generate()
			cli.StopProgress()

			if err != nil {
				return errors.Wrap(err, "failed to generate terraform code")
			}

			// Write-out generated code to location specified
			dirname, _, err := writeGeneratedCodeToLocation(cmd, hcl, "aws_controltower")
			if err != nil {
				return err
			}

			// Prompt to execute
			err = SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
				Prompt: &survey.Confirm{
					Default: GenerateAwsControlTowerCommandExtraState.TerraformApply,
					Message: QuestionRunTfPlan,
				},
				Response: &GenerateAwsControlTowerCommandExtraState.TerraformApply,
			})

			if err != nil {
				return errors.Wrap(err, "failed to prompt for terraform execution")
			}

			// Execute
			locationDir, _ := determineOutputDirPath(dirname, "aws_controltower")
			if GenerateAwsControlTowerCommandExtraState.TerraformApply {
				// Execution pre-run check
				err := executionPreRunChecks(dirname, locationDir, "aws_controltower")
				if err != nil {
					return err
				}
			}

			// Output where code was generated
			if !GenerateAwsControlTowerCommandExtraState.TerraformApply {
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

			// Validate s3_bucket_arn, if passed
			s3BucketArn, err := cmd.Flags().GetString("s3_bucket_arn")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsArnFormat(s3BucketArn); s3BucketArn != "" && err != nil {
				return err
			}

			// Validate sns_topic_arn, if passed
			snsTopicArn, err := cmd.Flags().GetString("sns_topic_arn")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsArnFormat(snsTopicArn); snsTopicArn != "" && err != nil {
				return err
			}

			// Parse audit_account, if passed
			if cmd.Flags().Changed("audit_account") {
				if err := parseAuditAccountFlag(GenerateAwsControlTowerCommandState); err != nil {
					return err
				}
			}

			GenerateAwsControlTowerCommandState.SubAccounts = append(GenerateAwsControlTowerCommandState.SubAccounts,
				aws_controltower.AwsSubAccount{
					AwsProfile: GenerateAwsControlTowerCommandState.AuditProfile,
					AwsRegion:  GenerateAwsControlTowerCommandState.AuditRegion})

			// Parse log_archive_account, if passed
			if cmd.Flags().Changed("log_archive_account") {
				if err := parseLogArchiveAccountFlag(GenerateAwsControlTowerCommandState); err != nil {
					return err
				}
			}

			GenerateAwsControlTowerCommandState.SubAccounts = append(GenerateAwsControlTowerCommandState.SubAccounts,
				aws_controltower.AwsSubAccount{
					AwsProfile: GenerateAwsControlTowerCommandState.LogArchiveProfile,
					AwsRegion:  GenerateAwsControlTowerCommandState.LogArchiveRegion})

			// Parse org_account_mapping json, if passed
			if cmd.Flags().Changed("org_account_mapping") {
				if err := parseOrgAccountMappingsFlag(GenerateAwsControlTowerCommandState); err != nil {
					return err
				}
			}

			// Load any cached inputs if interactive
			if cli.InteractiveMode() {
				cachedOptions := &aws_controltower.GenerateAwsControlTowerTfConfigurationArgs{}
				iacParamsExpired := cli.ReadCachedAsset(CachedAssetAwsControlTowerIacParams, &cachedOptions)
				if iacParamsExpired {
					cli.Log.Debug("loaded previously set values for AWS ControlTower IAC generation")
				}

				extraState := &AwsControlTowerGenerateCommandExtraState{}
				extraStateParamsExpired := cli.ReadCachedAsset(CachedAssetAwsControlTowerExtraState, &extraState)
				if extraStateParamsExpired {
					cli.Log.Debug("loaded previously set values for AWS ControlTower IAC generation (extra state)")
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
					if err := mergo.Merge(GenerateAwsControlTowerCommandState, cachedOptions); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
					if err := mergo.Merge(GenerateAwsControlTowerCommandExtraState, extraState); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
				}
			}

			// Collect and/or confirm parameters
			err = promptAwsControlTowerGenerate(GenerateAwsControlTowerCommandState,
				GenerateAwsControlTowerCommandExtraState)
			if err != nil {
				return errors.Wrap(err, "collecting/confirming parameters")
			}

			return nil
		},
	}
)

func useExistingIamRole(args *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) bool {
	return args.IamRoleArn != "" && args.IamRoleExternalID != "" && args.IamRoleName != ""
}

func parseAuditAccountFlag(args *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) error {
	parsedAccount := strings.Split(args.AuditAccount, ":")
	if len(parsedAccount) != 2 {
		return errors.New("invalid audit_account. Format must be 'profile:region'")
	}

	args.AuditProfile = parsedAccount[0]
	args.AuditRegion = parsedAccount[1]
	return nil
}

func parseLogArchiveAccountFlag(args *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) error {
	parsedAccount := strings.Split(args.LogArchiveAccount, ":")

	if len(parsedAccount) != 2 {
		fmt.Println(parsedAccount)
		return errors.New("invalid log_archive_account. Format must be 'profile:region'")
	}

	args.LogArchiveProfile = parsedAccount[0]
	args.LogArchiveRegion = parsedAccount[1]
	return nil
}

func parseOrgAccountMappingsFlag(args *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) error {
	if err := json.Unmarshal([]byte(args.OrgAccountMappingsJson), &args.OrgAccountMappings); err != nil {
		return errors.Wrap(err, "failed to parse 'org_account_mapping'")
	}

	return nil
}

type AwsControlTowerGenerateCommandExtraState struct {
	AskAdvanced             bool
	Output                  string
	ConfigureBucketSettings bool
	UseExistingKmsKey       bool
	MultiRegion             bool
	TerraformApply          bool
}

func (controltower *AwsControlTowerGenerateCommandExtraState) isEmpty() bool {
	return controltower.Output == "" &&
		!controltower.AskAdvanced &&
		!controltower.ConfigureBucketSettings &&
		!controltower.UseExistingKmsKey &&
		!controltower.TerraformApply
}

// Flush current state of the struct to disk, provided it's not empty
func (controltower *AwsControlTowerGenerateCommandExtraState) writeCache() {
	if !controltower.isEmpty() {
		cli.WriteAssetToCache(CachedAssetAwsControlTowerExtraState, time.Now().Add(time.Hour*1), controltower)
	}
}

func initGenerateAwsControlTowerTfCommandFlags() {
	// add flags to sub commands

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandState.LaceworkAccountID,
		"lacework_aws_account_id", "", "the Lacework AWS root account id")

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandState.S3BucketArn,
		"s3_bucket_arn", "", "the S3 Bucket for consolidated CloudTrail")

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandState.SNSTopicArn,
		"sns_topic_arn", "", "the SNS Topic")

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandState.LogArchiveAccount,
		"log_archive_account", "", "The log archive account flag input in the format profile:region")

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandState.AuditAccount,
		"audit_account", "", "The audit account flag input in the format profile:region")

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandState.OrgAccountMappingsJson,
		"org_account_mapping", "", "Org account mapping json string. Example: "+
			"'{\"default_lacework_account\":\"main\", \"mapping\": [{ \"aws_accounts\": [\"123456789011\"], "+
			"\"lacework_account\": \"sub-account-1\"}]}'")

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandState.IamRoleExternalID,
		"iam_role_external_id",
		"",
		"specify the external id of the existing iam role")

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandState.IamRoleName,
		"iam_role_name",
		"",
		"specify the name of the existing iam role")

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandState.IamRoleArn,
		"iam_role_arn",
		"",
		"specify the arn of the existing iam role")

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandState.SqsQueueName,
		"sqs_queue_name",
		"",
		"specify the name of the sqs queue")

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandState.Prefix,
		"prefix",
		"",
		"specify the prefix that will be used at the beginning of every generated resource")

	generateAwsControlTowerTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsControlTowerCommandExtraState.TerraformApply,
		"apply",
		false,
		"run terraform apply without executing plan or prompting")

	generateAwsControlTowerTfCommand.PersistentFlags().StringVar(
		&GenerateAwsControlTowerCommandExtraState.Output,
		"output",
		"",
		"location to write generated content")
}

func promptCustomizeControlTowerOutputLocation(extraState *AwsControlTowerGenerateCommandExtraState) error {
	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{Message: QuestionAwsControlTowerCustomizeOutputLocation,
			Default: extraState.Output},
		Response: &extraState.Output,
		Opts:     []survey.AskOpt{survey.WithValidator(validPathExists)},
		Required: true,
	})

	return err
}

func promptAwsIamRoleQuestions(input *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) error {
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{
			Message: QuestionAwsControlTowerCoreIamRoleName,
			Default: input.IamRoleName,
		},
		Opts:     []survey.AskOpt{},
		Response: &input.IamRoleName,
	}); err != nil {
		return err
	}

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{
			Message: QuestionAwsControlTowerCoreIamRoleArn,
			Default: input.IamRoleArn,
		},
		Opts:     []survey.AskOpt{survey.WithValidator(validateAwsArnFormat)},
		Response: &input.IamRoleArn,
	}); err != nil {
		return err
	}

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{
			Message: QuestionAwsControlTowerCoreIamRoleExternalID,
			Default: input.IamRoleExternalID,
		},
		Response: &input.IamRoleExternalID,
	}); err != nil {
		return err
	}

	return nil
}

func promptCustomIntegrationName(input *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) error {
	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{
			Message: QuestionControlTowerIntegrationName,
			Default: input.LaceworkIntegrationName,
		},
		Opts:     []survey.AskOpt{},
		Response: &input.LaceworkIntegrationName,
	})

	return err
}

func promptCustomPrefix(input *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) error {
	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{
			Message: QuestionControlTowerPrefix,
			Default: input.Prefix,
		},
		Opts:     []survey.AskOpt{},
		Response: &input.Prefix,
	})

	return err
}

func promptCustomSqsQueueName(input *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) error {
	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{
			Message: QuestionControlTowerSqsQueueName,
			Default: input.SqsQueueName,
		},
		Opts:     []survey.AskOpt{},
		Response: &input.SqsQueueName,
	})

	return err
}

func promptControlTowerAddOrgAccountMappings(input *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) error {
	mapping := aws_controltower.OrgAccountMap{}
	var accountsAnswer string
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionControlTowerOrgAccountMappingsLWAccount},
			Response: &mapping.LaceworkAccount,
		},
		{
			Prompt:   &survey.Multiline{Message: QuestionControlTowerOrgAccountMappingsAwsAccounts},
			Response: &accountsAnswer,
		},
	}); err != nil {
		return err
	}
	mapping.AwsAccounts = strings.Split(accountsAnswer, "\n")
	input.OrgAccountMappings.Mapping = append(input.OrgAccountMappings.Mapping, mapping)
	return nil
}

func promptControlTowerOrgAccountMappings(input *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) error {
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt: &survey.Input{
				Message: QuestionControlTowerOrgAccountMappingsLWDefaultAccount,
				Default: input.OrgAccountMappings.DefaultLaceworkAccount},
			Response: &input.OrgAccountMappings.DefaultLaceworkAccount,
		},
	}); err != nil {
		return err
	}

	if err := promptControlTowerAddOrgAccountMappings(input); err != nil {
		return err
	}

	var askAgain bool
	for {
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Prompt:   &survey.Confirm{Message: QuestionControlTowerOrgAccountMappingAnotherAdvancedOpt},
			Response: &askAgain}); err != nil {
			return err
		}

		if !askAgain {
			break
		}

		if err := promptControlTowerAddOrgAccountMappings(input); err != nil {
			return err
		}
	}

	return nil
}

func askAdvancedControlTowerOptions(config *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs,
	extraState *AwsControlTowerGenerateCommandExtraState) error {
	answer := ""

	//Prompt for options
	for answer != AwsAdvancedOptDone {
		var options []string

		options = append(options,
			ControlTowerConfigureExistingIamRoleOpt,
			ControlTowerIntegrationNameOpt,
			ControlTowerAdvancedOptLocation,
			ControlTowerIntegrationPrefixOpt,
			ControlTowerIntegrationSqsOpt,
			ControlTowerAdvancedOptMappings,
			ControlTowerAdvancedOptDone)
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
		case ControlTowerConfigureExistingIamRoleOpt:
			if err := promptAwsIamRoleQuestions(config); err != nil {
				return err
			}
			config.UseExistingIamRole = true
		case ControlTowerIntegrationNameOpt:
			if err := promptCustomIntegrationName(config); err != nil {
				return err
			}
		case ControlTowerIntegrationPrefixOpt:
			if err := promptCustomPrefix(config); err != nil {
				return err
			}
		case ControlTowerIntegrationSqsOpt:
			if err := promptCustomSqsQueueName(config); err != nil {
				return err
			}
		case ControlTowerAdvancedOptLocation:
			if err := promptCustomizeControlTowerOutputLocation(extraState); err != nil {
				return err
			}
		case ControlTowerAdvancedOptMappings:
			if err := promptControlTowerOrgAccountMappings(config); err != nil {
				return err
			}
		}

		// Re-prompt if not done
		innerAskAgain := true
		if answer == ControlTowerAdvancedOptDone {
			innerAskAgain = false
		}

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Checks:   []*bool{&innerAskAgain},
			Prompt:   &survey.Confirm{Message: QuestionControlTowerAnotherAdvancedOpt, Default: false},
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

func controltowerConfigIsEmpty(g *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) bool {
	return g.SNSTopicArn == "" &&
		g.S3BucketArn == "" &&
		g.LaceworkProfile == ""
}

func writeAwsControlTowerGenerationArgsCache(a *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs) {
	if !controltowerConfigIsEmpty(a) {
		cli.WriteAssetToCache(CachedAssetAwsControlTowerIacParams, time.Now().Add(time.Hour*1), a)
	}
}

// entry point for launching a survey to build out the required generation parameters
func promptAwsControlTowerGenerate(
	config *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs,
	extraState *AwsControlTowerGenerateCommandExtraState,
) error {

	// Cache for later use if generation is abandoned and in interactive mode
	if cli.InteractiveMode() {
		defer writeAwsControlTowerGenerationArgsCache(config)
		defer extraState.writeCache()
	}

	// Set Flags if set

	// prompt ControlTower core questions
	if err := promptAwsControlTowerCoreQuestions(config, extraState); err != nil {
		return err
	}

	// Find out if the customer wants to specify more advanced features
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Confirm{Message: QuestionAwsControlTowerConfigureAdvanced,
			Default: extraState.AskAdvanced},
		Response: &extraState.AskAdvanced,
	}); err != nil {
		return err
	}

	// Keep prompting for advanced options until the say done
	if extraState.AskAdvanced {
		if err := askAdvancedControlTowerOptions(config, extraState); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	initGenerateAwsControlTowerTfCommandFlags()
}

func promptAwsControlTowerCoreQuestions(config *aws_controltower.GenerateAwsControlTowerTfConfigurationArgs,
	state *AwsControlTowerGenerateCommandExtraState) error {
	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Prompt:   &survey.Input{Message: QuestionAwsControlTowerCoreS3Bucket, Default: config.S3BucketArn},
				Response: &config.S3BucketArn,
				Opts:     []survey.AskOpt{survey.WithValidator(validateAwsArnFormat)},
			},
			{
				Prompt:   &survey.Input{Message: QuestionAwsControlTowerCoreSnsTopic, Default: config.SNSTopicArn},
				Response: &config.SNSTopicArn,
				Opts:     []survey.AskOpt{survey.WithValidator(validateAwsArnFormat)},
			},
			{
				Prompt:   &survey.Input{Message: QuestionAwsControlTowerCoreLogProfile, Default: config.LogArchiveProfile},
				Response: &config.LogArchiveProfile,
			},
			{
				Prompt:   &survey.Input{Message: QuestionAwsControlTowerCoreLogRegion, Default: config.LogArchiveRegion},
				Response: &config.LogArchiveRegion,
				Opts:     []survey.AskOpt{survey.WithValidator(validateAwsRegion)},
			},
			{
				Prompt:   &survey.Input{Message: QuestionAwsControlTowerCoreAuditProfile, Default: config.AuditProfile},
				Response: &config.AuditProfile,
			},
			{
				Prompt:   &survey.Input{Message: QuestionAwsControlTowerCoreAuditRegion, Default: config.AuditRegion},
				Response: &config.AuditRegion,
				Opts:     []survey.AskOpt{survey.WithValidator(validateAwsRegion)},
			},
		}); err != nil {
		return err
	}
	config.SubAccounts = []aws_controltower.AwsSubAccount{
		aws_controltower.NewAwsSubAccount(config.LogArchiveProfile, config.LogArchiveRegion, "log_archive"),
		aws_controltower.NewAwsSubAccount(config.AuditProfile, config.AuditRegion, "audit")}

	return nil
}
