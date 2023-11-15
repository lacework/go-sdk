package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
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
	// Core questions
	QuestionEnableAwsOrganization = "Enable integrations for AWS organization?"
	QuestionMainAwsProfile        = "Main AWS account profile:"
	QuestionMainAwsRegion         = "Main AWS account region:"

	// Agentless questions
	QuestionEnableAgentless                  = "Enable Agentless integration?"
	QuestionAgentlessManagementAccountID     = "AWS management account ID:"
	QuestionAgentlessManagementAccountRegion = "AWS management account region:"

	QuestionAgentlessScanningAccountProfile = "Scanning AWS account profile:"
	QuestionAgentlessScanningAccountRegion  = "Scanning AWS account region:"
	QuestionAgentlessScanningAccountAddMore = "Add another scanning AWS account?"

	QuestionAgentlessScanningAccountsReplace = "Currently configured scanning accounts: %s, replace?"
	QuestionAgentlessMonitoredAccountIDs     = "Monitored AWS account ID list:"
	QuestionAgentlessMonitoredAccountIDsHelp = "Please provide a comma seprated list that may " +
		"contain account IDs, OUs, or the organization root (e.g. 123456789000,ou-abcd-12345678,r-abcd)."

	QuestionAgentlessMonitoredAccountProfile  = "Monitored AWS account profile:"
	QuestionAgentlessMonitoredAccountRegion   = "Monitored AWS account region:"
	QuestionAgentlessMonitoredAccountAddMore  = "Add another monitored AWS account?"
	QuestionAgentlessMonitoredAccountsReplace = "Currently configured monitored accounts: %s, replace?"

	// Config questions
	QuestionEnableConfig                    = "Enable configuration integration?"
	QuestionConfigAdditionalAccountProfile  = "Addtional AWS account profile:"
	QuestionConfigAdditionalAccountRegion   = "Addtional AWS account region:"
	QuestionConfigAdditionalAccountsReplace = "Currently configured additional accounts: %s, replace?"
	QuestionConfigAdditionalAccountAddMore  = "Add another AWS account?"

	// Config Org questions
	QuestionConfigOrgLWAccount        = "Lacework account:"
	QuestionConfigOrgLWSubaccount     = "Lacework subaccount (optional):"
	QuestionConfigOrgLWAccessKeyId    = "Lacework access key ID:"
	QuestionConfigOrgLWSecretKey      = "Lacework secret key:"
	QuestionConfigOrgId               = "AWS organization ID:"
	QuestionConfigOrgUnits            = "AWS organization units (multiple can be supplied comma separated):"
	QuestionConfigOrgCfResourcePrefix = "Cloudformation resource prefix:"

	// CloudTrail questions
	QuestionEnableCloudtrail   = "Enable CloudTrail integration?"
	QuestionCloudtrailName     = "Name of cloudtrail integration (optional):"
	QuestionCloudtrailAdvanced = "Configure advanced options?"

	// CloudTrail advanced options
	OptCloudtrailMessage = "Which options would you like to configure?"

	OptCloudtrailOrg  = "Configure org account mappings"
	OptCloudtrailS3   = "Configure S3 bucket"
	OptCloudtrailSNS  = "Configure SNS topic"
	OptCloudtrailSQS  = "Configure SQS queue"
	OptCloudtrailIAM  = "Configure an existing IAM role"
	OptCloudtrailDone = "Done"

	// CloudTrail Org questions
	QuestionCloudtrailOrgAccountMappingsDefaultLWAccount = "Org account mappings default Lacework account:"
	QuestionCloudtrailOrgAccountMappingsAnotherAddMore   = "Add another org account mapping?"
	QuestionCloudtrailOrgAccountMappingsLWAccount        = "Lacework account:"
	QuestionCloudtrailOrgAccountMappingsAwsAccounts      = "AWS accounts:"

	// CloudTrail S3 Bucket Questions
	QuestionCloudtrailUseConsolidated          = "Use consolidated CloudTrail?"
	QuestionCloudtrailUseExistingS3            = "Use an existing CloudTrail?"
	QuestionCloudtrailS3ExistingBucketArn      = "Existing S3 bucket ARN used for CloudTrail logs:"
	QuestionCloudtrailS3BucketEnableEncryption = "Enable S3 bucket encryption"

	QuestionCloudtrailS3BucketSseKeyArn    = "Existing KMS encryption key arn for S3 bucket (optional):"
	QuestionCloudtrailS3BucketName         = "New S3 bucket name (optional):"
	QuestionCloudtrailS3BucketNotification = "Enable S3 bucket notifications"

	// CloudTrail SNS Topic Questions
	QuestionCloudtrailUseExistingSNSTopic = "Use an existing SNS topic?"
	QuestionCloudtrailSnsExistingTopicArn = "Existing SNS topic arn:"
	QuestionCloudtrailSnsEnableEncryption = "Enable encryption on SNS topic?"
	QuestionCloudtrailSnsEncryptionKeyArn = "Existing KMS encryption key arn for SNS topic (optional):"
	QuestionCloudtrailSnsTopicName        = "New SNS topic name (optional):"

	// CloudTrail SQS Queue Questions
	QuestionCloudtrailSqsEnableEncryption = "Enable encryption on SQS queue:"
	QuestionCloudtrailSqsEncryptionKeyArn = "Existing KMS encryption key arn for SQS queue (optional):"
	QuestionCloudtrailSqsQueueName        = "New SQS queue name (optional):"

	// CloudTrail IAM Role Questions
	QuestionCloudtrailExistingIamRoleName  = "Existing IAM role name for CloudTrail access:"
	QuestionCloudtrailExistingIamRoleArn   = "Existing IAM role ARN for CloudTrail access:"
	QuestionCloudtrailExistingIamRoleExtID = "External ID for the existing IAM role:"

	// Custom location Question
	QuestionAwsOutputLocation = "Custom output location (optional):"

	// Other options
	AwsAdvancedOptDone = "Done" // Used in aws controltower and eks_audit

	// Question labels
	IconAgentless  = "[Agentless]"
	IconConfig     = "[Configuration]"
	IconCloudTrail = "[CloudTrail]"

	// AwsArnRegex original source: https://regex101.com/r/pOfxYN/1
	AwsArnRegex = `^arn:(?P<Partition>[^:\n]*):(?P<Service>[^:\n]*):(?P<Region>[^:\n]*):(?P<AccountID>[^:\n]*):(?P<Ignore>(?P<ResourceType>[^:\/\n]*)[:\/])?(?P<Resource>.*)$` //nolint
	// AwsRegionRegex regex used for validating region input; note intentionally does not match gov cloud
	AwsRegionRegex              = `(af|ap|ca|eu|me|sa|us)-(central|(north|south)?(east|west)?)-\d`
	AwsProfileRegex             = `([A-Za-z_0-9-]+)`
	AwsAccountIDRegex           = `^\d{12}$`
	AwsOUIDRegex                = `^ou-[0-9a-z]{4,32}-[a-z0-9]{8,32}$`
	AWSRootIDRegex              = `^r-[0-9a-z]{4,32}$`
	AwsAssumeRoleRegex          = `^arn:aws:iam::\d{12}:role\/.*$`
	ValidateSubAccountFlagRegex = fmt.Sprintf(`%s:%s`, AwsProfileRegex, AwsRegionRegex)

	GenerateAwsCommandState = &aws.GenerateAwsTfConfigurationArgs{
		ExistingIamRole: &aws.ExistingIamRoleDetails{},
	}
	GenerateAwsCommandExtraState = &aws.AwsGenerateCommandExtraState{}

	CachedAwsArgsKey       = "iac-aws-generate-args"
	CachedAwsExtraStateKey = "iac-aws-extra-state"

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
	* If Terraform is not installed, or the version installed is not compatible, a new
    version will be installed into a temporary location
	* Once Terraform is detected or installed, Terraform plan will be executed
	* The command will prompt with the outcome of the plan and allow to view more details
    or continue with Terraform apply
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
				aws.WithAwsRegion(GenerateAwsCommandState.AwsRegion),
				aws.WithAwsAssumeRole(GenerateAwsCommandState.AwsAssumeRole),
				aws.WithLaceworkProfile(GenerateAwsCommandState.LaceworkProfile),
				aws.WithLaceworkAccountID(GenerateAwsCommandState.LaceworkAccountID),
				aws.WithAgentlessManagementAccountID(GenerateAwsCommandState.AgentlessManagementAccountID),
				aws.WithAgentlessMonitoredAccountIDs(GenerateAwsCommandState.AgentlessMonitoredAccountIDs),
				aws.WithAgentlessMonitoredAccounts(GenerateAwsCommandState.AgentlessMonitoredAccounts...),
				aws.WithAgentlessScanningAccounts(GenerateAwsCommandState.AgentlessScanningAccounts...),
				aws.WithConfigAdditionalAccounts(GenerateAwsCommandState.ConfigAdditionalAccounts...),
				aws.WithConfigOrgLWAccount(GenerateAwsCommandState.ConfigOrgLWAccount),
				aws.WithConfigOrgLWSubaccount(GenerateAwsCommandState.ConfigOrgLWSubaccount),
				aws.WithConfigOrgLWAccessKeyId(GenerateAwsCommandState.ConfigOrgLWAccessKeyId),
				aws.WithConfigOrgLWSecretKey(GenerateAwsCommandState.ConfigOrgLWSecretKey),
				aws.WithConfigOrgId(GenerateAwsCommandState.ConfigOrgId),
				aws.WithConfigOrgUnits(GenerateAwsCommandState.ConfigOrgUnits),
				aws.WithConfigOrgCfResourcePrefix(GenerateAwsCommandState.ConfigOrgCfResourcePrefix),
				aws.WithConsolidatedCloudtrail(GenerateAwsCommandState.ConsolidatedCloudtrail),
				aws.WithCloudtrailUseExistingS3(GenerateAwsCommandState.CloudtrailUseExistingS3),
				aws.WithCloudtrailUseExistingSNSTopic(GenerateAwsCommandState.CloudtrailUseExistingSNSTopic),
				aws.WithExistingCloudtrailBucketArn(GenerateAwsCommandState.ExistingCloudtrailBucketArn),
				aws.WithExistingSnsTopicArn(GenerateAwsCommandState.ExistingSnsTopicArn),
				aws.WithSubaccounts(GenerateAwsCommandState.SubAccounts...),
				aws.WithExistingIamRole(GenerateAwsCommandState.ExistingIamRole),
				aws.WithCloudtrailName(GenerateAwsCommandState.CloudtrailName),
				aws.WithOrgAccountMappings(GenerateAwsCommandState.OrgAccountMappings),
				aws.WithBucketName(GenerateAwsCommandState.BucketName),
				aws.WithBucketEncryptionEnabled(GenerateAwsCommandState.BucketEncryptionEnabled),
				aws.WithBucketSSEKeyArn(GenerateAwsCommandState.BucketSseKeyArn),
				aws.WithSnsTopicName(GenerateAwsCommandState.SnsTopicName),
				aws.WithSnsTopicEncryptionEnabled(GenerateAwsCommandState.SnsTopicEncryptionEnabled),
				aws.WithSnsTopicEncryptionKeyArn(GenerateAwsCommandState.SnsTopicEncryptionKeyArn),
				aws.WithSqsQueueName(GenerateAwsCommandState.SqsQueueName),
				aws.WithSqsEncryptionEnabled(GenerateAwsCommandState.SqsEncryptionEnabled),
				aws.WithSqsEncryptionKeyArn(GenerateAwsCommandState.SqsEncryptionKeyArn),
				aws.WithS3BucketNotification(GenerateAwsCommandState.S3BucketNotification),
			}

			// Create new struct
			data := aws.NewTerraform(
				GenerateAwsCommandState.AwsOrganization,
				GenerateAwsCommandState.Agentless,
				GenerateAwsCommandState.Config,
				GenerateAwsCommandState.Cloudtrail,
				mods...,
			)

			// Generate
			hcl, err := data.Generate()
			cli.StopProgress()

			if err != nil {
				return errors.Wrap(err, "failed to generate terraform code")
			}

			// Write-out generated code to location specified
			dirname, _, err := writeGeneratedCodeToLocation(cmd, hcl, "aws")
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

			locationDir, _ := determineOutputDirPath(dirname, "aws")
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

			// Validate aws assume role, if passed
			assumeRole, err := cmd.Flags().GetString("aws_assume_role")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsAssumeRole(assumeRole); assumeRole != "" && err != nil {
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

			// Parse cloudtrail org_account_mapping json, if passed
			if cmd.Flags().Changed("cloudtrail_org_account_mapping") {
				if err := parseCloudtrailOrgAccountMappingsFlag(GenerateAwsCommandState); err != nil {
					return err
				}
			}

			// Validate cloudtrail bucket arn, if passed
			arn, err := cmd.Flags().GetString("existing_bucket_arn")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsArnFormat(arn); arn != "" && err != nil {
				return err
			}
			if arn != "" {
				GenerateAwsCommandState.CloudtrailUseExistingS3 = true
			}

			// Validate SNS Topic Arn if passed
			arn, err = cmd.Flags().GetString("existing_sns_topic_arn")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsArnFormat(arn); arn != "" && err != nil {
				return err
			}
			if arn != "" {
				GenerateAwsCommandState.CloudtrailUseExistingSNSTopic = true
			}

			// Load any cached inputs if interactive
			if cli.InteractiveMode() {
				cachedOptions := &aws.GenerateAwsTfConfigurationArgs{}
				awsArgsExpired := cli.ReadCachedAsset(CachedAwsArgsKey, &cachedOptions)
				if awsArgsExpired {
					cli.Log.Debug("loaded previously set values for AWS iac generation")
				}

				extraState := &aws.AwsGenerateCommandExtraState{}
				extraStateExpired := cli.ReadCachedAsset(CachedAwsExtraStateKey, &extraState)
				if extraStateExpired {
					cli.Log.Debug("loaded previously set values for AWS iac generation (extra state)")
				}

				// Determine if previously cached options exists; prompt user if they'd like to continue
				answer := false
				if !awsArgsExpired || !extraStateExpired {
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

				// Collect and/or confirm parameters
				err = promptAwsGenerate(GenerateAwsCommandState, GenerateAwsCommandExtraState)
				if err != nil {
					return errors.Wrap(err, "collecting/confirming parameters")
				}
			}

			// Parse passed in AWS accounts
			if len(GenerateAwsCommandExtraState.AwsSubAccounts) > 0 {
				accounts, err := parseAwsAccountsFromCommandFlag(GenerateAwsCommandExtraState.AwsSubAccounts)
				if err != nil {
					return err
				}
				GenerateAwsCommandState.SubAccounts = accounts
				GenerateAwsCommandState.ConfigAdditionalAccounts = accounts
			}

			// Parse passed in Agentless monirtoed AWS accounts
			if len(GenerateAwsCommandExtraState.AgentlessMonitoredAccounts) > 0 {
				accounts, err := parseAwsAccountsFromCommandFlag(GenerateAwsCommandExtraState.AgentlessMonitoredAccounts)
				if err != nil {
					return err
				}
				GenerateAwsCommandState.AgentlessMonitoredAccounts = accounts
			}

			// Parse passed in Agentless scanning AWS accounts
			if len(GenerateAwsCommandExtraState.AgentlessScanningAccounts) > 0 {
				accounts, err := parseAwsAccountsFromCommandFlag(GenerateAwsCommandExtraState.AgentlessScanningAccounts)
				if err != nil {
					return err
				}
				GenerateAwsCommandState.AgentlessScanningAccounts = accounts
			}

			return nil
		},
	}
)

func parseAwsAccountsFromCommandFlag(accountsInput []string) ([]aws.AwsSubAccount, error) {
	// Validate the format of supplied values is correct
	if err := validateAwsSubAccounts(accountsInput); err != nil {
		return nil, err
	}
	accounts := []aws.AwsSubAccount{}
	for _, account := range accountsInput {
		accountDetails := strings.Split(account, ":")
		profile := accountDetails[0]
		region := accountDetails[1]
		alias := fmt.Sprintf("%s-%s", profile, region)
		accounts = append(accounts, aws.NewAwsSubAccount(profile, region, alias))
	}
	return accounts, nil
}

func parseCloudtrailOrgAccountMappingsFlag(args *aws.GenerateAwsTfConfigurationArgs) error {
	if err := json.Unmarshal([]byte(args.OrgAccountMappingsJson), &args.OrgAccountMappings); err != nil {
		return errors.Wrap(err, "failed to parse 'cloudtrail_org_account_mapping'")
	}
	return nil
}

func initGenerateAwsTfCommandFlags() {
	// add flags to sub commands
	// TODO Share the help with the interactive generation
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.AwsOrganization,
		"aws_organization",
		false,
		"enable organization integration")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.Agentless,
		"agentless",
		false,
		"enable agentless integration")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.AgentlessManagementAccountID,
		"agentless_management_account_id",
		"",
		"AWS management account ID for Agentless integration")
	generateAwsTfCommand.PersistentFlags().StringSliceVar(
		&GenerateAwsCommandState.AgentlessMonitoredAccountIDs,
		"agentless_monitored_account_ids",
		[]string{},
		"AWS monitored account IDs for Agentless integrations")
	generateAwsTfCommand.PersistentFlags().StringSliceVar(
		&GenerateAwsCommandExtraState.AgentlessMonitoredAccounts,
		"agentless_monitored_accounts",
		[]string{},
		"AWS monitored accounts for Agentless integrations; value format must be <aws profile>:<region>")
	generateAwsTfCommand.PersistentFlags().StringSliceVar(
		&GenerateAwsCommandExtraState.AgentlessScanningAccounts,
		"agentless_scanning_accounts",
		[]string{},
		"AWS scanning accounts for Agentless integrations; value format must be <aws profile>:<region>")
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
		&GenerateAwsCommandState.ConfigOrgLWAccount,
		"config_lacework_account",
		"",
		"specify lacework account for Config organization integration")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ConfigOrgLWSubaccount,
		"config_lacework_sub_account",
		"",
		"specify lacework sub-account for Config organization integration")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ConfigOrgLWAccessKeyId,
		"config_lacework_access_key_id",
		"",
		"specify AWS access key ID for Config organization integration")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ConfigOrgLWSecretKey,
		"config_lacework_secret_key",
		"",
		"specify AWS secret key for Config organization integration")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ConfigOrgId,
		"config_organization_id",
		"",
		"specify AWS organization ID for Config organization integration")
	generateAwsTfCommand.PersistentFlags().StringSliceVar(
		&GenerateAwsCommandState.ConfigOrgUnits,
		"config_organization_units",
		nil,
		"specify AWS organization unit for Config organization integration")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ConfigOrgCfResourcePrefix,
		"config_cf_resource_prefix",
		"",
		"specify Cloudformation resource prefix for Config organization integration")
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
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.AwsAssumeRole,
		"aws_assume_role",
		"",
		"specify aws assume role")
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
		&GenerateAwsCommandState.ExistingIamRole.Arn,
		"existing_iam_role_arn",
		"",
		"specify existing iam role arn to use")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ExistingIamRole.Name,
		"existing_iam_role_name",
		"",
		"specify existing iam role name to use")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ExistingIamRole.ExternalId,
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
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.OrgAccountMappingsJson,
		"cloudtrail_org_account_mapping", "", "Org account mapping json string. Example: "+
			"'{\"default_lacework_account\":\"main\", \"mapping\": [{ \"aws_accounts\": [\"123456789011\"], "+
			"\"lacework_account\": \"sub-account-1\"}]}'")

	// DEPRECATED
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.ForceDestroyS3Bucket,
		"force_destroy_s3",
		true,
		"enable force destroy S3 bucket")
	errcheckWARN(generateAwsTfCommand.PersistentFlags().MarkDeprecated(
		"force_destroy_s3", "by default, force destroy is enabled.",
	))
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.ConfigName,
		"config_name",
		"",
		"specify name of config integration")
	errcheckWARN(generateAwsTfCommand.PersistentFlags().MarkDeprecated(
		"config_name", "default config is used.",
	))
	// ---

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
		&GenerateAwsCommandState.SnsTopicEncryptionEnabled,
		"sns_topic_encryption_enabled",
		true,
		"enable encryption on SNS topic when creating one")
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.SnsTopicEncryptionKeyArn,
		"sns_topic_encryption_key_arn",
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
	generateAwsTfCommand.PersistentFlags().StringVar(
		&GenerateAwsCommandState.LaceworkAccountID,
		"lacework_aws_account_id",
		"",
		"the Lacework AWS root account id")
	generateAwsTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsCommandState.S3BucketNotification,
		"use_s3_bucket_notification",
		false,
		"enable S3 bucket notifications")
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

func validateAwsAccountID(val interface{}) error {
	return validateStringWithRegex(val, AwsAccountIDRegex, "invalid account ID supplied")
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

func validateAgentlessMonitoredAccountIDList(val interface{}) error {
	switch value := val.(type) {
	case string:
		regex := fmt.Sprintf(`%s|%s|%s`, AwsAccountIDRegex, AwsOUIDRegex, AWSRootIDRegex)
		ids := strings.Split(value, ",")
		for _, id := range ids {
			if err := validateStringWithRegex(
				id,
				regex,
				fmt.Sprintf("invalid account ID, OU ID or root ID supplied: %s", id),
			); err != nil {
				return err
			}
		}
	default:
		// if the value passed is not a string
		return errors.New("value must be a string")
	}
	return nil
}

// survey.Validator for aws region
func validateAwsRegion(val interface{}) error {
	return validateStringWithRegex(val, AwsRegionRegex, "invalid region supplied")
}

// survey.Validator for aws profile
func validateAwsProfile(val interface{}) error {
	return validateStringWithRegex(val, fmt.Sprintf(`^%s$`, AwsProfileRegex), "invalid profile name supplied")
}

// survey.Validator for aws profile
func validateAwsAssumeRole(val interface{}) error {
	return validateStringWithRegex(val, AwsAssumeRoleRegex, "invalid assume name supplied")
}

func promptAgentlessQuestions(config *aws.GenerateAwsTfConfigurationArgs) error {
	if !config.Agentless {
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Icon:     IconAgentless,
			Prompt:   &survey.Confirm{Message: QuestionEnableAgentless, Default: config.Agentless},
			Response: &config.Agentless,
		}); err != nil {
			return err
		}
	}

	if !config.Agentless {
		return nil
	}

	if config.AwsOrganization {
		monitoredAccountIDListInput := ""

		if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
			{
				Icon: IconAgentless,
				Prompt: &survey.Input{
					Message: QuestionAgentlessManagementAccountID,
					Default: config.AgentlessManagementAccountID,
				},
				Opts:     []survey.AskOpt{survey.WithValidator(validateAwsAccountID)},
				Response: &config.AgentlessManagementAccountID,
				Required: true,
			},
			{
				Icon: IconAgentless,
				Prompt: &survey.Input{
					Message: QuestionAgentlessMonitoredAccountIDs,
					Default: strings.Join(config.AgentlessMonitoredAccountIDs, ","),
					Help:    QuestionAgentlessMonitoredAccountIDsHelp,
				},
				Opts:     []survey.AskOpt{survey.WithValidator(validateAgentlessMonitoredAccountIDList)},
				Response: &monitoredAccountIDListInput,
				Required: true,
			},
		}, config.AwsOrganization); err != nil {
			return err
		}

		config.AgentlessMonitoredAccountIDs = strings.Split(monitoredAccountIDListInput, ",")

		if err := promptAwsAccountsQuestions(
			&config.AgentlessMonitoredAccounts,
			IconAgentless,
			QuestionAgentlessMonitoredAccountProfile,
			QuestionAgentlessMonitoredAccountRegion,
			QuestionAgentlessMonitoredAccountAddMore,
			QuestionAgentlessMonitoredAccountsReplace,
			false,
		); err != nil {
			return err
		}
	}

	if err := promptAwsAccountsQuestions(
		&config.AgentlessScanningAccounts,
		IconAgentless,
		QuestionAgentlessScanningAccountProfile,
		QuestionAgentlessScanningAccountRegion,
		QuestionAgentlessScanningAccountAddMore,
		QuestionAgentlessScanningAccountsReplace,
		!config.AwsOrganization,
	); err != nil {
		return err
	}

	return nil
}

func promptConfigQuestions(config *aws.GenerateAwsTfConfigurationArgs) error {
	if !config.Config {
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Icon:     IconConfig,
			Prompt:   &survey.Confirm{Message: QuestionEnableConfig, Default: config.Config},
			Response: &config.Config,
		}); err != nil {
			return err
		}
	}

	if !config.Config {
		return nil
	}

	if config.AwsOrganization {
		if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
			{
				Icon:     IconConfig,
				Prompt:   &survey.Input{Message: QuestionConfigOrgLWAccount, Default: config.ConfigOrgLWAccount},
				Response: &config.ConfigOrgLWAccount,
				Required: true,
			},
			{
				Icon:     IconConfig,
				Prompt:   &survey.Input{Message: QuestionConfigOrgLWSubaccount, Default: config.ConfigOrgLWSubaccount},
				Response: &config.ConfigOrgLWSubaccount,
			},
			{
				Icon:     IconConfig,
				Prompt:   &survey.Input{Message: QuestionConfigOrgLWAccessKeyId, Default: config.ConfigOrgLWAccessKeyId},
				Response: &config.ConfigOrgLWAccessKeyId,
				Required: true,
			},
			{
				Icon:     IconConfig,
				Prompt:   &survey.Input{Message: QuestionConfigOrgLWSecretKey, Default: config.ConfigOrgLWSecretKey},
				Response: &config.ConfigOrgLWSecretKey,
				Required: true,
			},
			{
				Icon:     IconConfig,
				Prompt:   &survey.Input{Message: QuestionConfigOrgId, Default: config.ConfigOrgId},
				Response: &config.ConfigOrgId,
				Required: true,
			},
		}); err != nil {
			return err
		}

		var orgUnitsInput string
		if err := survey.AskOne(
			&survey.Input{Message: QuestionConfigOrgUnits, Default: strings.Join(config.ConfigOrgUnits, ",")}, &orgUnitsInput,
			survey.WithValidator(survey.Required), survey.WithIcons(customPromptIconsFunc(IconConfig)),
		); err != nil {
			return err
		}
		config.ConfigOrgUnits = strings.Split(orgUnitsInput, ",")

		if err := survey.AskOne(
			&survey.Input{
				Message: QuestionConfigOrgCfResourcePrefix, Default: config.ConfigOrgCfResourcePrefix,
			}, &config.ConfigOrgCfResourcePrefix,
			survey.WithValidator(survey.Required), survey.WithIcons(customPromptIconsFunc(IconConfig)),
		); err != nil {
			return err
		}

		return nil
	}

	if err := promptAwsAccountsQuestions(
		&config.ConfigAdditionalAccounts,
		IconConfig,
		QuestionConfigAdditionalAccountProfile,
		QuestionConfigAdditionalAccountRegion,
		QuestionConfigAdditionalAccountAddMore,
		QuestionConfigAdditionalAccountsReplace,
		true,
	); err != nil {
		return err
	}

	return nil
}

func promptCloudtrailQuestions(
	config *aws.GenerateAwsTfConfigurationArgs,
	extraState *aws.AwsGenerateCommandExtraState,
) error {
	if !config.Cloudtrail {
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Confirm{Message: QuestionEnableCloudtrail, Default: config.Cloudtrail},
			Response: &config.Cloudtrail,
		}); err != nil {
			return err
		}
	}

	if !config.Cloudtrail {
		return nil
	}

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Icon:     IconCloudTrail,
		Prompt:   &survey.Confirm{Message: QuestionCloudtrailUseConsolidated, Default: config.ConsolidatedCloudtrail},
		Response: &config.ConsolidatedCloudtrail,
	}); err != nil {
		return err
	}

	// Find out if the customer wants to specify more advanced features
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Icon:     IconCloudTrail,
		Prompt:   &survey.Confirm{Message: QuestionCloudtrailAdvanced, Default: extraState.CloudtrailAdvanced},
		Response: &extraState.CloudtrailAdvanced,
	}); err != nil {
		return err
	}

	// Keep prompting for advanced options until the say done
	if extraState.CloudtrailAdvanced {
		answer := ""
		options := []string{
			OptCloudtrailS3,
			OptCloudtrailSNS,
			OptCloudtrailSQS,
			OptCloudtrailIAM,
			OptCloudtrailDone,
		}
		if config.AwsOrganization {
			options = append([]string{OptCloudtrailOrg}, options...)
		}
		for answer != OptCloudtrailDone {
			if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
				Icon: IconCloudTrail,
				Prompt: &survey.Select{
					Message: OptCloudtrailMessage,
					Options: options,
				},
				Response: &answer,
			}); err != nil {
				return err
			}
			switch answer {
			case OptCloudtrailOrg:
				if err := promptCloudtrailOrgQuestions(config); err != nil {
					return err
				}
			case OptCloudtrailS3:
				if err := promptCloudtrailS3Questions(config); err != nil {
					return err
				}
			case OptCloudtrailSNS:
				if err := promptCloudtrailSNSQuestions(config); err != nil {
					return err
				}
			case OptCloudtrailSQS:
				if err := promptCloudtrailSQSQuestions(config); err != nil {
					return err
				}
			case OptCloudtrailIAM:
				if err := promptCloudtrailIAMQuestions(config); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func promptCloudtrailOrgAccountMappingQuestions(config *aws.GenerateAwsTfConfigurationArgs) error {
	mapping := aws.OrgAccountMap{}
	var accountsAnswer string
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailOrgAccountMappingsLWAccount},
			Response: &mapping.LaceworkAccount,
		},
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Multiline{Message: QuestionCloudtrailOrgAccountMappingsAwsAccounts},
			Response: &accountsAnswer,
		},
	}); err != nil {
		return err
	}
	mapping.AwsAccounts = strings.Split(accountsAnswer, "\n")
	config.OrgAccountMappings.Mapping = append(config.OrgAccountMappings.Mapping, mapping)
	return nil
}

func promptCloudtrailOrgQuestions(config *aws.GenerateAwsTfConfigurationArgs) error {
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon: IconCloudTrail,
			Prompt: &survey.Input{
				Message: QuestionCloudtrailOrgAccountMappingsDefaultLWAccount,
				Default: config.OrgAccountMappings.DefaultLaceworkAccount},
			Response: &config.OrgAccountMappings.DefaultLaceworkAccount,
		},
	}); err != nil {
		return err
	}

	askAgain := true
	for askAgain {
		if err := promptCloudtrailOrgAccountMappingQuestions(config); err != nil {
			return err
		}
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Confirm{Message: QuestionCloudtrailOrgAccountMappingsAnotherAddMore},
			Response: &askAgain}); err != nil {
			return err
		}
	}

	return nil
}

func promptCloudtrailS3Questions(config *aws.GenerateAwsTfConfigurationArgs) error {
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Confirm{Message: QuestionCloudtrailUseExistingS3, Default: config.CloudtrailUseExistingS3},
			Response: &config.CloudtrailUseExistingS3,
		},
	}); err != nil {
		return err
	}

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailName, Default: config.CloudtrailName},
			Response: &config.CloudtrailName,
		},
		{
			Icon: IconCloudTrail,
			Prompt: &survey.Input{
				Message: QuestionCloudtrailS3ExistingBucketArn,
				Default: config.ExistingCloudtrailBucketArn,
			},
			Required: true,
			Opts:     []survey.AskOpt{survey.WithValidator(validateAwsArnFormat)},
			Response: &config.ExistingCloudtrailBucketArn,
		},
	}, config.CloudtrailUseExistingS3); err != nil {
		return err
	}

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailS3BucketName, Default: config.BucketName},
			Response: &config.BucketName,
		},
		{
			Icon: IconCloudTrail,
			Prompt: &survey.Confirm{
				Message: QuestionCloudtrailS3BucketEnableEncryption,
				Default: config.BucketEncryptionEnabled,
			},
			Response: &config.BucketEncryptionEnabled,
		},
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailS3BucketSseKeyArn, Default: config.BucketSseKeyArn},
			Response: &config.BucketSseKeyArn,
			Opts:     []survey.AskOpt{survey.WithValidator(validateOptionalAwsArnFormat)},
			Checks:   []*bool{&config.BucketEncryptionEnabled},
		},
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Confirm{Message: QuestionCloudtrailS3BucketNotification, Default: config.S3BucketNotification},
			Response: &config.S3BucketNotification,
		},
	}, !config.CloudtrailUseExistingS3); err != nil {
		return err
	}

	return nil
}

func promptCloudtrailSNSQuestions(config *aws.GenerateAwsTfConfigurationArgs) error {
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon: IconCloudTrail,
			Prompt: &survey.Confirm{
				Message: QuestionCloudtrailUseExistingSNSTopic,
				Default: config.CloudtrailUseExistingSNSTopic,
			},
			Response: &config.CloudtrailUseExistingSNSTopic,
		},
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailSnsExistingTopicArn, Default: config.ExistingSnsTopicArn},
			Checks:   []*bool{&config.CloudtrailUseExistingSNSTopic},
			Required: true,
			Opts:     []survey.AskOpt{survey.WithValidator(validateAwsArnFormat)},
			Response: &config.ExistingSnsTopicArn,
		},
	}); err != nil {
		return err
	}

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailSnsTopicName, Default: config.SnsTopicName},
			Response: &config.SnsTopicName,
		},
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Confirm{Message: QuestionCloudtrailSnsEnableEncryption, Default: config.SnsTopicEncryptionEnabled},
			Response: &config.SnsTopicEncryptionEnabled,
		},
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailSnsEncryptionKeyArn, Default: config.SnsTopicEncryptionKeyArn},
			Response: &config.SnsTopicEncryptionKeyArn,
			Opts:     []survey.AskOpt{survey.WithValidator(validateOptionalAwsArnFormat)},
			Checks:   []*bool{&config.SnsTopicEncryptionEnabled},
		},
	}, !config.CloudtrailUseExistingSNSTopic); err != nil {
		return err
	}

	return nil
}

func promptCloudtrailSQSQuestions(config *aws.GenerateAwsTfConfigurationArgs) error {
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailSqsQueueName, Default: config.SqsQueueName},
			Response: &config.SqsQueueName,
		},
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Confirm{Message: QuestionCloudtrailSqsEnableEncryption, Default: config.SqsEncryptionEnabled},
			Response: &config.SqsEncryptionEnabled,
		},
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailSqsEncryptionKeyArn, Default: config.SqsEncryptionKeyArn},
			Response: &config.SqsEncryptionKeyArn,
			Opts:     []survey.AskOpt{survey.WithValidator(validateOptionalAwsArnFormat)},
			Checks:   []*bool{&config.SqsEncryptionEnabled},
		},
	}); err != nil {
		return err
	}
	return nil
}

func promptCloudtrailIAMQuestions(config *aws.GenerateAwsTfConfigurationArgs) error {
	// ensure struct is initialized
	if config.ExistingIamRole == nil {
		config.ExistingIamRole = &aws.ExistingIamRoleDetails{}
	}

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailExistingIamRoleName, Default: config.ExistingIamRole.Name},
			Response: &config.ExistingIamRole.Name,
			// Opts:     []survey.AskOpt{survey.WithValidator(survey.Required)},
			Required: true,
		},
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailExistingIamRoleArn, Default: config.ExistingIamRole.Arn},
			Response: &config.ExistingIamRole.Arn,
			Opts:     []survey.AskOpt{survey.WithValidator(validateAwsArnFormat)},
			Required: true,
		},
		{
			Icon:     IconCloudTrail,
			Prompt:   &survey.Input{Message: QuestionCloudtrailExistingIamRoleExtID, Default: config.ExistingIamRole.ExternalId},
			Response: &config.ExistingIamRole.ExternalId,
			// Opts:     []survey.AskOpt{survey.WithValidator(survey.Required)},
			Required: true,
		}}); err != nil {
		return err
	}

	return nil
}

func promptAwsAccountsQuestions(
	accounts *[]aws.AwsSubAccount,
	questionIcon string,
	questionProfile string,
	questionRegion string,
	questionAddMore string,
	questionReplace string,
	askFirst bool,
) error {
	if !cli.InteractiveMode() {
		return nil
	}

	askAgain := true
	newAccounts := []aws.AwsSubAccount{}

	// Ask if replacing existing accounts
	if len(*accounts) > 0 {
		accountListing := []string{}
		for _, account := range *accounts {
			accountListing = append(
				accountListing,
				fmt.Sprintf("%s:%s", account.AwsProfile, account.AwsRegion),
			)
		}
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Icon: questionIcon,
			Prompt: &survey.Confirm{
				Message: fmt.Sprintf(
					questionReplace,
					strings.Trim(strings.Join(strings.Fields(fmt.Sprint(accountListing)), ", "), "[]"),
				),
			},
			Response: &askAgain,
		}); err != nil {
			return err
		}
	}

	if askFirst {
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Icon:     questionIcon,
			Prompt:   &survey.Confirm{Message: questionAddMore},
			Response: &askAgain,
		}); err != nil {
			return err
		}
	}

	for askAgain {
		var profile, region string
		if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
			{
				Icon:     questionIcon,
				Prompt:   &survey.Input{Message: questionProfile},
				Opts:     []survey.AskOpt{survey.WithValidator(validateAwsProfile)},
				Required: true,
				Response: &profile,
			},
			{
				Icon:     questionIcon,
				Prompt:   &survey.Input{Message: questionRegion},
				Opts:     []survey.AskOpt{survey.WithValidator(validateAwsRegion)},
				Required: true,
				Response: &region,
			},
		}); err != nil {
			return err
		}
		alias := fmt.Sprintf("%s-%s", profile, region)
		newAccounts = append(
			newAccounts,
			aws.AwsSubAccount{AwsProfile: profile, AwsRegion: region, Alias: alias},
		)
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Icon:     questionIcon,
			Prompt:   &survey.Confirm{Message: questionAddMore},
			Response: &askAgain,
		}); err != nil {
			return err
		}
	}

	if len(newAccounts) > 0 {
		*accounts = newAccounts
	}

	return nil
}

func writeArgsCache(a *aws.GenerateAwsTfConfigurationArgs) {
	if !a.IsEmpty() {
		// If ExistingIamRole is partially set, don't write this to cache; the values won't work when loaded
		if a.ExistingIamRole.IsPartial() {
			a.ExistingIamRole = nil
		}
		cli.WriteAssetToCache(CachedAwsArgsKey, time.Now().Add(time.Hour*1), a)
	}
}

func writeExtraStateCache(a *aws.AwsGenerateCommandExtraState) {
	if !a.IsEmpty() {
		cli.WriteAssetToCache(CachedAwsExtraStateKey, time.Now().Add(time.Hour*1), a)
	}
}

// Entry point for launching a survey to build out the required generation parameters
func promptAwsGenerate(
	config *aws.GenerateAwsTfConfigurationArgs,
	extraState *aws.AwsGenerateCommandExtraState,
) error {
	// Cache for later use if generation is abandon and in interactive mode
	if cli.InteractiveMode() {
		defer writeArgsCache(config)
		defer writeExtraStateCache(extraState)
	}

	// Core questions
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt: &survey.Confirm{
				Message: QuestionEnableAwsOrganization,
				Default: config.AwsOrganization,
			},
			Response: &config.AwsOrganization,
		},
		{
			Prompt:   &survey.Input{Message: QuestionMainAwsProfile, Default: config.AwsProfile},
			Opts:     []survey.AskOpt{survey.WithValidator(validateAwsProfile)},
			Response: &config.AwsProfile,
			Required: true,
		},
		{
			Prompt:   &survey.Input{Message: QuestionMainAwsRegion, Default: config.AwsRegion},
			Opts:     []survey.AskOpt{survey.WithValidator(validateAwsRegion)},
			Response: &config.AwsRegion,
			Required: true,
		},
	}); err != nil {
		return err
	}

	if err := promptAgentlessQuestions(config); err != nil {
		return err
	}
	if err := promptConfigQuestions(config); err != nil {
		return err
	}
	if err := promptCloudtrailQuestions(config, extraState); err != nil {
		return err
	}

	// Custom ouput location
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionAwsOutputLocation, Default: extraState.Output},
		Response: &extraState.Output,
		Opts:     []survey.AskOpt{survey.WithValidator(validPathExists)},
	}); err != nil {
		return err
	}

	return nil
}
