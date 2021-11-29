package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/imdario/mergo"
	"github.com/lacework/go-sdk/lwgenerate/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	GenerateAwsCommandState      = &aws.GenerateAwsTfConfigurationArgs{}
	GenerateAwsExistingRoleState = &aws.ExistingIamRoleDetails{}

	// iac-generate command is used to create IaC code for various environments
	generateTfCommand = &cobra.Command{
		Use:     "iac-generate",
		Aliases: []string{"iac-generate", "iac"},
		Short:   "create iac code",
		Long:    "Create IaC content for various different cloud environments and configurations",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			// If output location was supplied, validate it exists
			dirname, err := cmd.Flags().GetString("output")
			if err == nil {
				_, err := os.Stat(dirname)
				if err != nil {
					return errors.Wrap(err, "could not access specified output location!")
				}
			}

			return nil
		},
	}

	// aws command is used to generate TF code for aws
	generateAwsTfCommand = &cobra.Command{
		Use:   "aws",
		Short: "generate code for aws environment",
		Long:  "Genereate Terraform code for deploying into a new AWS environment.",
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
			location, err := writeHclOutput(hcl, cmd)
			if err != nil {
				return errors.Wrap(err, "failed to write terrraform code to disk")
			}

			// Prompt to execute
			execute := false
			err = SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
				Prompt:   &survey.Confirm{Default: execute, Message: "Run Terraform plan now?"},
				Response: &execute,
			})

			if err != nil {
				return errors.Wrap(err, "failed to run terraform execution")
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
			// Load any cached inputs
			cachedOptions := &aws.GenerateAwsTfConfigurationArgs{}
			if ok := cli.ReadCachedAsset("iac-aws-generate-params", &cachedOptions); ok {
				cli.Log.Debug("loaded previously arguments for AWS iac generation")
			}

			// Merge cached inputs to current options (current options win)
			if err := mergo.Merge(GenerateAwsCommandState, cachedOptions); err != nil {
				return errors.Wrap(err, "failed to load saved options!")
			}

			// Try to load cached location
			var outputLocation string
			if ok := cli.ReadCachedAsset("tfLocation", &outputLocation); !ok {
				cli.Log.Debug("couldn't load cached output location")
			}

			// Collect and/or confirm parameters
			err := promptAwsGenerate(GenerateAwsCommandState, GenerateAwsExistingRoleState, &outputLocation)
			if err != nil {
				return err
			}

			// Store location
			cli.WriteAssetToCache("tfLocation", time.Now().Add(1*time.Hour), outputLocation)

			// reset location if it was updated
			if outputLocation != "" {
				if err := cmd.Flags().Set("output", outputLocation); err != nil {
					return err
				}
			}

			return nil
		},
	}
)

func init() {
	// add the iac-generate command
	rootCmd.AddCommand(generateTfCommand)

	// Add global flags for iac generation
	generateTfCommand.PersistentFlags().String("output", "", "location to write generated content")

	// add flags to sub commands
	// TODO Share the help with the interactive generation
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.Cloudtrail, "cloudtrail", false, "enable cloudtrail integration")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.Config, "config", false, "enable config integration")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.AwsRegion, "aws_region", "", "specify aws region")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.AwsProfile, "aws_profile", "", "specify aws profile")
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

// Write HCL output
func writeHclOutput(hcl string, cmd *cobra.Command) (string, error) {
	// Write out
	var dirname string
	dirname, err := cmd.Flags().GetString("output")
	if err != nil {
		return "", err
	}

	if dirname == "" {
		dirname, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}

	directory := filepath.FromSlash(fmt.Sprintf("%s/%s", dirname, "lacework"))
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.Mkdir(directory, 0700)
		if err != nil {
			return "", err
		}
	}

	location := fmt.Sprintf("%s/%s/main.tf", dirname, "lacework")
	err = os.WriteFile(
		filepath.FromSlash(location),
		[]byte(hcl),
		0700,
	)
	if err != nil {
		return "", err
	}

	cli.StopProgress()
	return location, nil
}
