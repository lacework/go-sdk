package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/imdario/mergo"
	"github.com/lacework/go-sdk/lwgenerate/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type AwsGenerateCommandExtraState struct {
	Output                string
	UseExistingCloudtrail bool
	AwsSubAccounts        []string
}

var (
	QuestionRunTfPlan            = "Run Terraform plan now?"
	QuestionUsePreviousCache     = "Previous IaC generation detected, load cached values?"
	GenerateAwsCommandState      = &aws.GenerateAwsTfConfigurationArgs{}
	GenerateAwsExistingRoleState = &aws.ExistingIamRoleDetails{}
	GenerateAwsCommandExtraState = &AwsGenerateCommandExtraState{}
	ValidateSubAccountFlagRegex  = fmt.Sprintf(`%s:%s`, AwsProfileRegex, AwsRegionRegex)
	CachedAssetIacParams         = "iac-aws-generate-params"
	CachedAssetAwsExtraState     = "iac-aws-extra-state"

	// iac-generate command is used to create IaC code for various environments
	generateTfCommand = &cobra.Command{
		Use:     "iac-generate",
		Aliases: []string{"iac"},
		Short:   "Create IaC code",
		Long:    "Create IaC content for various different cloud environments and configurations",
	}

	// aws command is used to generate TF code for aws
	generateAwsTfCommand = &cobra.Command{
		Use:   "aws",
		Short: "Genereate terraform code for deploying into a new AWS environment",
		Long: `Use this command to generate Terraform code for deploying Lacework into an AWS environment.

By default, this command will function interactively, prompting for the required information to setup the new cloud account. In interactive mode, this command will:

* Prompt for the required information to setup the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version will be confirmed suitable for use
	* If Terraform is not installed, or the version installed is not suitable, a new version will be installed into a temporary location
	* Once Terraform is detected or installed, Terraform plan will be executed
	* The command will prompt with the outcome of the plan and allow to view more details or continue with Terraform apply
	* If confirmed, Terraform apply will be run, completing the setup of the cloud account

This command can also be run in noninteractive mode however, only generation is supported at this time.  See help output for more details on supplying required values for generation.
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
			dirname, err := cmd.Flags().GetString("output")
			if err != nil {
				return errors.Wrap(err, "failed to parse output location")
			}

			ok, err := writeHclOutputPrecheck(dirname)
			if err != nil {
				return errors.Wrap(err, "failed to validate output location")
			}

			if !ok {
				return errors.Wrap(err, "aborting to avoid overwriting existing terraform code")
			}

			location, err := writeHclOutput(hcl, dirname)
			if err != nil {
				return errors.Wrap(err, "failed to write terrraform code to disk")
			}

			// Prompt to execute
			execute := false
			err = SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
				Prompt:   &survey.Confirm{Default: execute, Message: QuestionRunTfPlan},
				Response: &execute,
			})

			if err != nil {
				return errors.Wrap(err, "failed to run terraform execution")
			}

			// Execution pre-run check
			ok, err = TerraformExecutePreRunCheck(dirname)
			if err != nil {
				return errors.Wrap(err, "failed to check for existing terraform state")
			}

			if !ok {
				return errors.Wrap(err, "aborting to avoid overwriting existing terraform state")
			}

			// Execute
			locationDir := filepath.Dir(location)
			if execute {
				if err := TerraformPlanAndExecute(locationDir); err != nil {
					return errors.Wrap(err, "failed to run terraform apply")
				}
			}

			// Output where code was generated
			if !execute {
				provideGuidanceAfterExit(false, false, locationDir, "terraform")
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

			// Load any cached inputs if interactive
			if cli.InteractiveMode() {
				cachedOptions := &aws.GenerateAwsTfConfigurationArgs{}
				iacParamsExpired := cli.ReadCachedAsset(CachedAssetIacParams, &cachedOptions)
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
						Prompt:   &survey.Confirm{Message: QuestionUsePreviousCache, Default: answer},
						Response: &answer,
					}); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
				}

				// If the user decides NOT to use the previous values; we won't load them.  However, every time the command runs
				// we are going to write out new cached values, so if they run it - bail out - and run it again they'll get
				// reprompted.
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
				// validate consolidated_cloudtrail is enabled - otherwise this flag cannot be used
				if ok, _ := cmd.Flags().GetBool("consolidated_cloudtrail"); !ok {
					return errors.New("aws subaccounts can only be supplied with consolidated cloudtrail enabled")
				}

				// validate the format of supplied values is correct
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

// Prompt use for question, only if the CLI is in interactive mode
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

// Prompt for many values at once
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

// Prompt for confirmation if main.tf already exists; return true to continue
func writeHclOutputPrecheck(outputLocation string) (bool, error) {
	// If noninteractive, continue
	if !cli.InteractiveMode() {
		return true, nil
	}

	outputDir, err := determineOutputDirPath(outputLocation)
	if err != nil {
		return false, err
	}

	hclPath := filepath.FromSlash(fmt.Sprintf("%s/main.tf", outputDir))

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

// Write HCL output
func writeHclOutput(hcl string, location string) (string, error) {
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

	// Create HCL file
	outputLocation := filepath.FromSlash(fmt.Sprintf("%s/main.tf", dirname))
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

// This function used to validate provided output location exists and is a directory
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
