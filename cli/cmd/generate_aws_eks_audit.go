package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lacework/go-sdk/lwgenerate/aws_eks_audit"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

var (
	// Define question text here, so they can be reused in testing
	QuestionEksAuditMultiRegion          = "Integrate clusters in more than one region?"
	QuestionEksAuditRegionClusterCurrent = "Currently configured regions and clusters: %s. " +
		"Configure additional?"
	QuestionEksAuditRegion         = "Specify AWS region:"
	QuestionEksAuditRegionClusters = "Specify a comma-seperated list of clusters in region" +
		" to ingest EKS Audit Logs:"
	QuestionEksAuditAdditionalRegion = "Configure another AWS region?"

	QuestionEksAuditConfigureAdvanced = "Configure advanced integration options?"

	// S3 Bucket Questions
	EksAuditConfigureBucket              = "Configure bucket settings"
	QuestionEksAuditBucketVersioning     = "Enable access versioning on the new bucket?"
	QuestionEksAuditMfaDeleteS3Bucket    = "Should MFA object deletion be required for the new bucket?"
	QuestionEksAuditForceDestroyS3Bucket = "Should force destroy be enabled for the new bucket?"
	QuestionEksAuditBucketLifecycle      = "Specify the bucket lifecycle expiration days: (optional)"
	QuestionEksAuditBucketEncryption     = "Enable encryption for the new bucket?"
	QuestionEksAuditBucketSseAlgorithm   = "Specify the bucket SSE Algorithm: (optional)"
	QuestionEksAuditBucketExistingKey    = "Use existing KMS key?"
	QuestionEksAuditBucketKeyArn         = "Specify the bucket existing SSE KMS key ARN:"
	QuestionEksAuditKmsKeyRotation       = "Should the KMS key have rotation enabled?"
	QuestionEksAuditKmsKeyDeletionDays   = "Specify the KMS key deletion days: (optional)"

	// SNS Topic Questions
	EksAuditConfigureSns                = "Configure SNS settings"
	QuestionEksAuditSnsEncryption       = "Enable encryption on SNS topic when creating?"
	QuestionEksAuditSnsEncryptionKeyArn = "Specify existing KMS encryption key ARN for SNS topic (optional)"

	// Cloudwatch IAM Questions
	EksAuditExistingCwIamRole        = "Configure and use existing Cloudwatch IAM role"
	QuestionEksAuditExistingCwIamArn = "Specify an existing Cloudwatch IAM role ARN:"

	// Firehose Questions
	EksAuditConfigureFh                = "Configure Firehose settings"
	QuestionEksAuditExistingFhIamRole  = "Use existing Firehose IAM role?"
	QuestionEksAuditExistingFhIamArn   = "Specify an existing Firehose IAM role ARN:"
	QuestionEksAuditFhEncryption       = "Enable encryption on Firehose when creating?"
	QuestionEksAuditFhEncryptionKeyArn = "Specify existing KMS encryption key ARN for Firehose (optional)"

	// Cross Account IAM Questions
	EksAuditExistingCaIamRole          = "Configure and use existing Cross Account IAM role"
	QuestionEksAuditExistingCaIamArn   = "Specify an existing Cross Account IAM role ARN:"
	QuestionEksAuditExistingCaIamExtID = "Specify the external ID to be used with the existing IAM role:"

	// Customize integration name
	EksAuditIntegrationNameOpt            = "Customize integration name"
	QuestionEksAuditCustomIntegrationName = "Specify a custom integration name: (optional)"

	// Customize output location
	EksAuditAdvancedOptLocation             = "Customize output location"
	QuestionEksAuditCustomizeOutputLocation = "Provide the location for the output to be written:"

	QuestionEksAuditAnotherAdvancedOpt = "Configure another advanced integration option"
	EksAuditAdvancedOptDone            = "Done"

	// AwsEksAuditRegionRegex regex used for validating region input; note intentionally does not match gov cloud
	AwsEksAuditRegionRegex = `(us|ap|ca|eu|sa)-(central|(north|south)?(east|west)?)-\d`

	GenerateAwsEksAuditCommandState      = &aws_eks_audit.GenerateAwsEksAuditTfConfigurationArgs{}
	GenerateAwsEksAuditCommandExtraState = &AwsEksAuditGenerateCommandExtraState{}
	GenerateAwsEksAuditExistingRoleState = &aws_eks_audit.ExistingCrossAccountIamRoleDetails{}
	CachedAssetAwsEksAuditIacParams      = "iac-aws-eks-audit-generate-params"
	CachedAssetAwsEksAuditExtraState     = "iac-aws-eks-audit-extra-state"

	// aws-eks-audit-log command is used to generate TF code for aws-eks-audit-log
	generateAwsEksAuditTfCommand = &cobra.Command{
		Use:   "eks",
		Short: "Generate and/or execute Terraform code for EKS integration",
		Long: `Use this command to generate Terraform code for deploying Lacework into an EKS
environment.

By default, this command interactively prompts for the required information to set up the new cloud account.
In interactive mode, this command will:

* Prompt for the required information to set up the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version is verified as compatible for use
  * If Terraform is not installed, or the version installed is not compatible, a new version will be installed into a temporary location
  * Once Terraform is detected or installed, the Terraform plan is executed
  * The command prompts you with the outcome of the plan and allows you to view more details or continue with Terraform apply
  * If confirmed, Terraform apply runs, completing the setup of the cloud account

This command can also be run in noninteractive mode.
See help output for more details on the parameter values required for Terraform code generation.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generate TF Code
			cli.StartProgress("Generating Terraform Code...")

			// Explicitly set Lacework profile if it was passed in main args
			if cli.Profile != "default" {
				GenerateAwsEksAuditCommandState.LaceworkProfile = cli.Profile
			}

			// Setup modifiers for NewTerraform constructor
			mods := []aws_eks_audit.AwsEksAuditTerraformModifier{
				aws_eks_audit.WithAwsProfile(GenerateAwsEksAuditCommandState.AwsProfile),
				aws_eks_audit.WithLaceworkAccountID(GenerateAwsEksAuditCommandState.LaceworkAccountID),
				aws_eks_audit.WithBucketLifecycleExpirationDays(GenerateAwsEksAuditCommandState.BucketLifecycleExpirationDays),
				aws_eks_audit.WithBucketSseAlgorithm(GenerateAwsEksAuditCommandState.BucketSseAlgorithm),
				aws_eks_audit.WithBucketSseKeyArn(GenerateAwsEksAuditCommandState.BucketSseKeyArn),
				aws_eks_audit.WithEksAuditIntegrationName(GenerateAwsEksAuditCommandState.EksAuditIntegrationName),
				aws_eks_audit.WithExistingCloudWatchIamRoleArn(GenerateAwsEksAuditCommandState.ExistingCloudWatchIamRoleArn),
				aws_eks_audit.WithExistingCrossAccountIamRole(GenerateAwsEksAuditCommandState.ExistingCrossAccountIamRole),
				aws_eks_audit.WithExistingFirehoseIamRoleArn(GenerateAwsEksAuditCommandState.ExistingFirehoseIamRoleArn),
				aws_eks_audit.WithFilterPattern(GenerateAwsEksAuditCommandState.FilterPattern),
				aws_eks_audit.WithFirehoseEncryptionKeyArn(GenerateAwsEksAuditCommandState.FirehoseEncryptionKeyArn),
				aws_eks_audit.WithKmsKeyDeletionDays(GenerateAwsEksAuditCommandState.KmsKeyDeletionDays),
				aws_eks_audit.WithPrefix(GenerateAwsEksAuditCommandState.Prefix),
				aws_eks_audit.WithParsedRegionClusterMap(GenerateAwsEksAuditCommandState.ParsedRegionClusterMap),
				aws_eks_audit.WithSnsTopicEncryptionKeyArn(GenerateAwsEksAuditCommandState.SnsTopicEncryptionKeyArn),
				aws_eks_audit.WithLaceworkProfile(GenerateAwsEksAuditCommandState.LaceworkProfile),
				aws_eks_audit.EnableBucketEncryption(GenerateAwsEksAuditCommandState.BucketEnableEncryption),
				aws_eks_audit.EnableBucketVersioning(GenerateAwsEksAuditCommandState.BucketVersioning),
				aws_eks_audit.EnableFirehoseEncryption(GenerateAwsEksAuditCommandState.FirehoseEncryptionEnabled),
				aws_eks_audit.EnableSnsTopicEncryption(GenerateAwsEksAuditCommandState.SnsTopicEncryptionEnabled),
				aws_eks_audit.EnableBucketVersioning(GenerateAwsEksAuditCommandState.BucketVersioning),
				aws_eks_audit.EnableKmsKeyRotation(GenerateAwsEksAuditCommandState.KmsKeyRotation),
			}

			if GenerateAwsEksAuditCommandState.BucketEnableMfaDelete {
				mods = append(mods, aws_eks_audit.EnableBucketMfaDelete())
			}

			if GenerateAwsEksAuditCommandState.BucketForceDestroy {
				mods = append(mods, aws_eks_audit.EnableBucketForceDestroy())
			}

			// Create new struct
			data := aws_eks_audit.NewTerraform(mods...)

			// Generate
			hcl, err := data.Generate()
			cli.StopProgress()

			if err != nil {
				return errors.Wrap(err, "failed to generate terraform code")
			}

			// Write-out generated code to location specified
			dirname, _, err := writeGeneratedCodeToLocation(cmd, hcl, "aws_eks_audit")
			if err != nil {
				return err
			}

			// Prompt to execute
			err = SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
				Prompt:   &survey.Confirm{Default: GenerateAwsEksAuditCommandExtraState.TerraformApply, Message: QuestionRunTfPlan},
				Response: &GenerateAwsEksAuditCommandExtraState.TerraformApply,
			})

			if err != nil {
				return errors.Wrap(err, "failed to prompt for terraform execution")
			}

			// Execute
			locationDir, _ := determineOutputDirPath(dirname, "aws_eks_audit")
			if GenerateAwsEksAuditCommandExtraState.TerraformApply {
				// Execution pre-run check
				err := executionPreRunChecks(dirname, locationDir, "aws_eks_audit")
				if err != nil {
					return err
				}
			}

			// Output where code was generated
			if !GenerateAwsEksAuditCommandExtraState.TerraformApply {
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

			// Validate bucket sse key ARN, if passed
			bucketSseKeyArn, err := cmd.Flags().GetString("bucket_sse_key_arn")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsArnFormat(bucketSseKeyArn); bucketSseKeyArn != "" && err != nil {
				return err
			}

			// Validate firehose key ARN, if passed
			firehoseKeyArn, err := cmd.Flags().GetString("firehose_encryption_key_arn")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsArnFormat(firehoseKeyArn); firehoseKeyArn != "" && err != nil {
				return err
			}

			// Validate sns topic key ARN, if passed
			snsKeyArn, err := cmd.Flags().GetString("sns_topic_encryption_key_arn")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateAwsArnFormat(snsKeyArn); snsKeyArn != "" && err != nil {
				return err
			}

			// Load any cached inputs if interactive
			if cli.InteractiveMode() {
				cachedOptions := &aws_eks_audit.GenerateAwsEksAuditTfConfigurationArgs{}
				iacParamsExpired := cli.ReadCachedAsset(CachedAssetAwsEksAuditIacParams, &cachedOptions)
				if iacParamsExpired {
					cli.Log.Debug("loaded previously set values for AWS EKS Audit IAC generation")
				}

				extraState := &AwsEksAuditGenerateCommandExtraState{}
				extraStateParamsExpired := cli.ReadCachedAsset(CachedAssetAwsEksAuditExtraState, &extraState)
				if extraStateParamsExpired {
					cli.Log.Debug("loaded previously set values for AWS EKS Audit IAC generation (extra state)")
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
					if err := mergo.Merge(GenerateAwsEksAuditCommandState, cachedOptions); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
					if err := mergo.Merge(GenerateAwsEksAuditCommandExtraState, extraState); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
				}
			}

			// Parse regions passed as part of the region cluster map
			if len(GenerateAwsEksAuditCommandState.RegionClusterMap) > 0 {
				// validate the format of supplied values is correct

				awsParsedRegionClusterMap := make(map[string][]string)
				for region, clusters := range GenerateAwsEksAuditCommandState.RegionClusterMap {
					// verify each region is a valid aws region
					if err := validateStringWithRegex(region, AwsEksAuditRegionRegex,
						"invalid region name supplied"); err != nil {
						return err
					}
					// parse the cluster comma-seperated string into a list of clusters
					parsedClusters := strings.Split(clusters, ",")
					awsParsedRegionClusterMap[region] = append(awsParsedRegionClusterMap[region], parsedClusters...)
				}
				GenerateAwsEksAuditCommandState.ParsedRegionClusterMap = awsParsedRegionClusterMap
			}

			// Collect and/or confirm parameters
			err = promptAwsEksAuditGenerate(GenerateAwsEksAuditCommandState, GenerateAwsEksAuditExistingRoleState,
				GenerateAwsEksAuditCommandExtraState)
			if err != nil {
				return errors.Wrap(err, "collecting/confirming parameters")
			}

			return nil
		},
	}
)

type AwsEksAuditGenerateCommandExtraState struct {
	AskAdvanced             bool
	Output                  string
	ConfigureBucketSettings bool
	UseExistingKmsKey       bool
	MultiRegion             bool
	TerraformApply          bool
}

func (eks *AwsEksAuditGenerateCommandExtraState) isEmpty() bool {
	return eks.Output == "" &&
		!eks.AskAdvanced &&
		!eks.ConfigureBucketSettings &&
		!eks.UseExistingKmsKey &&
		!eks.TerraformApply
}

// Flush current state of the struct to disk, provided it's not empty
func (eks *AwsEksAuditGenerateCommandExtraState) writeCache() {
	if !eks.isEmpty() {
		cli.WriteAssetToCache(CachedAssetAwsEksAuditExtraState, time.Now().Add(time.Hour*1), eks)
	}
}

func initGenerateAwsEksAuditTfCommandFlags() {
	// add flags to sub commands
	// TODO Share the help with the interactive generation
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandState.AwsProfile, "aws_profile", "", "specify aws profile")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandState.LaceworkAccountID, "lacework_aws_account_id", "", "the Lacework AWS root account id")
	generateAwsEksAuditTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsEksAuditCommandState.BucketEnableMfaDelete,
		"enable_mfa_delete_s3",
		false,
		"enable mfa delete on s3 bucket. Requires bucket versioning.")
	generateAwsEksAuditTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsEksAuditCommandState.BucketEnableEncryption,
		"enable_encryption_s3",
		true,
		"enable encryption on s3 bucket")
	generateAwsEksAuditTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsEksAuditCommandState.BucketForceDestroy,
		"enable_force_destroy",
		false,
		"enable force destroy s3 bucket")
	generateAwsEksAuditTfCommand.PersistentFlags().IntVar(
		&GenerateAwsEksAuditCommandState.BucketLifecycleExpirationDays,
		"bucket_lifecycle_exp_days",
		0,
		"specify the s3 bucket lifecycle expiration days")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandState.BucketSseAlgorithm,
		"bucket_sse_algorithm",
		"",
		"specify the encryption algorithm to use for S3 bucket server-side encryption")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandState.BucketSseKeyArn,
		"bucket_sse_key_arn",
		"",
		"specify the kms key arn to be used for s3. (required when bucket_sse_algorithm is aws:kms & using an existing kms key)")
	generateAwsEksAuditTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsEksAuditCommandState.BucketVersioning,
		"enable_bucket_versioning",
		true,
		"enable s3 bucket versioning")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandState.EksAuditIntegrationName,
		"integration_name",
		"",
		"specify the name of the eks audit integration")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandState.ExistingCloudWatchIamRoleArn,
		"existing_cw_iam_role_arn",
		"",
		"specify existing cloudwatch iam role arn to use")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditExistingRoleState.Arn,
		"existing_ca_iam_role_arn",
		"",
		"specify existing cross account iam role arn to use")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditExistingRoleState.ExternalId,
		"existing_ca_iam_role_external_id",
		"",
		"specify existing cross account iam role external_id to use")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandState.ExistingFirehoseIamRoleArn,
		"existing_firehose_iam_role_arn",
		"",
		"specify existing firehose iam role arn to use")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandState.FilterPattern,
		"custom_filter_pattern",
		"",
		"specify a custom cloudwatch log filter pattern")
	generateAwsEksAuditTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsEksAuditCommandState.FirehoseEncryptionEnabled,
		"enable_firehose_encryption",
		true,
		"enable firehose encryption")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandState.FirehoseEncryptionKeyArn,
		"firehose_encryption_key_arn",
		"",
		"specify the kms key arn to be used with the Firehose")
	generateAwsEksAuditTfCommand.PersistentFlags().IntVar(
		&GenerateAwsEksAuditCommandState.KmsKeyDeletionDays,
		"kms_key_deletion_days",
		0,
		"specify the kms waiting period before deletion, in number of days")
	generateAwsEksAuditTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsEksAuditCommandState.KmsKeyRotation,
		"enable_kms_key_rotation",
		true,
		"enable automatic kms key rotation")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandState.Prefix,
		"prefix",
		"",
		"specify the prefix that will be used at the beginning of every generated resource")
	generateAwsEksAuditTfCommand.PersistentFlags().StringToStringVar(
		&GenerateAwsEksAuditCommandState.RegionClusterMap,
		"region_clusters",
		map[string]string{},
		"configure eks clusters per aws region. To configure multiple regions pass the flag"+
			" multiple times. Example format:  --region_clusters <region>=\"cluster,list\"")
	generateAwsEksAuditTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsEksAuditCommandState.SnsTopicEncryptionEnabled,
		"enable_sns_topic_encryption",
		true,
		"enable encryption on the sns topic")
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandState.SnsTopicEncryptionKeyArn,
		"sns_topic_encryption_key_arn",
		"",
		"specify the kms key arn to be used with the sns topic")
	generateAwsEksAuditTfCommand.PersistentFlags().BoolVar(
		&GenerateAwsEksAuditCommandExtraState.TerraformApply,
		"apply",
		false,
		"run terraform apply without executing plan or prompting",
	)
	generateAwsEksAuditTfCommand.PersistentFlags().StringVar(
		&GenerateAwsEksAuditCommandExtraState.Output,
		"output",
		"",
		"location to write generated content",
	)
}

// Validate the response is of type int
func validateResponseTypeInt(val interface{}) error {
	switch value := val.(type) {
	case string:
		if _, err := strconv.Atoi(value); err != nil {
			// if the value passed is not of type int
			return errors.New("value must be a number")
		}
	default:
		// if the value passed is not a string
		return errors.New("value must be a string")
	}
	return nil
}

func promptAwsEksAuditBucketQuestions(config *aws_eks_audit.GenerateAwsEksAuditTfConfigurationArgs) error {
	// Only ask these questions if configure bucket is true
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionEksAuditBucketVersioning, Default: config.BucketVersioning},
			Response: &config.BucketVersioning,
		},
		{
			Prompt:   &survey.Confirm{Message: QuestionEksAuditMfaDeleteS3Bucket, Default: config.BucketEnableMfaDelete},
			Response: &config.BucketEnableMfaDelete,
		},
		{
			Prompt:   &survey.Confirm{Message: QuestionEksAuditForceDestroyS3Bucket, Default: config.BucketForceDestroy},
			Response: &config.BucketForceDestroy,
		},
		{
			Prompt: &survey.Confirm{Message: QuestionEksAuditBucketEncryption,
				Default: config.BucketEnableEncryption},
			Response: &config.BucketEnableEncryption,
		},
	}); err != nil {
		return err
	}

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Message: QuestionEksAuditBucketExistingKey},
		Checks:   []*bool{&config.BucketEnableEncryption},
		Opts:     []survey.AskOpt{},
		Response: &config.ExistingBucketKmsKey,
	}); err != nil {
		return err
	}

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionEksAuditBucketSseAlgorithm},
		Checks:   []*bool{&config.BucketEnableEncryption},
		Opts:     []survey.AskOpt{},
		Response: &config.BucketSseAlgorithm,
	}); err != nil {
		return err
	}

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionEksAuditBucketKeyArn},
		Checks:   []*bool{&config.BucketEnableEncryption, &config.ExistingBucketKmsKey},
		Opts:     []survey.AskOpt{survey.WithValidator(validateAwsArnFormat)},
		Response: &config.BucketSseKeyArn,
	}); err != nil {
		return err
	}

	newKmsKey := config.BucketEnableEncryption && !config.ExistingBucketKmsKey
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionEksAuditKmsKeyRotation},
			Checks:   []*bool{&config.BucketEnableEncryption, &newKmsKey},
			Required: true,
			Opts:     []survey.AskOpt{},
			Response: &config.KmsKeyRotation,
		},
		{
			Prompt:   &survey.Input{Message: QuestionEksAuditKmsKeyDeletionDays, Default: strconv.Itoa(config.KmsKeyDeletionDays)},
			Checks:   []*bool{&config.BucketEnableEncryption, &newKmsKey},
			Opts:     []survey.AskOpt{survey.WithValidator(validateResponseTypeInt)},
			Response: &config.KmsKeyDeletionDays,
		},
		{
			Prompt:   &survey.Input{Message: QuestionEksAuditBucketLifecycle, Default: strconv.Itoa(config.BucketLifecycleExpirationDays)},
			Opts:     []survey.AskOpt{survey.WithValidator(validateResponseTypeInt)},
			Response: &config.BucketLifecycleExpirationDays,
		},
	}); err != nil {
		return err
	}

	return nil
}

func promptAwsEksAuditExistingCrossAccountIamQuestions(input *aws_eks_audit.
	GenerateAwsEksAuditTfConfigurationArgs) error {
	// ensure struct is initialized
	if input.ExistingCrossAccountIamRole == nil {
		input.ExistingCrossAccountIamRole = &aws_eks_audit.ExistingCrossAccountIamRoleDetails{}
	}

	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionEksAuditExistingCaIamArn, Default: input.ExistingCrossAccountIamRole.Arn},
			Response: &input.ExistingCrossAccountIamRole.Arn,
			Opts:     []survey.AskOpt{survey.WithValidator(survey.Required), survey.WithValidator(validateAwsArnFormat)},
		},
		{
			Prompt:   &survey.Input{Message: QuestionEksAuditExistingCaIamExtID, Default: input.ExistingCrossAccountIamRole.ExternalId},
			Response: &input.ExistingCrossAccountIamRole.ExternalId,
			Opts:     []survey.AskOpt{survey.WithValidator(survey.Required)},
		}})
	return err
}

func promptAwsEksAuditFirehoseQuestions(input *aws_eks_audit.
	GenerateAwsEksAuditTfConfigurationArgs) error {

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Confirm{
			Message: QuestionEksAuditExistingFhIamRole,
			Default: input.ExistingFirehoseIam,
		},
		Opts:     []survey.AskOpt{},
		Response: &input.ExistingFirehoseIam,
		Required: false,
	}); err != nil {
		return err
	}

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt: &survey.Input{
				Message: QuestionEksAuditExistingFhIamArn,
				Default: input.ExistingFirehoseIamRoleArn,
			},
			Checks:   []*bool{&input.ExistingFirehoseIam},
			Opts:     []survey.AskOpt{survey.WithValidator(validateAwsArnFormat)},
			Response: &input.ExistingFirehoseIamRoleArn,
			Required: true,
		},
		{
			Prompt: &survey.Confirm{
				Message: QuestionEksAuditFhEncryption,
				Default: input.FirehoseEncryptionEnabled,
			},
			Opts:     []survey.AskOpt{},
			Response: &input.FirehoseEncryptionEnabled,
			Required: true,
		},
	}); err != nil {
		return err
	}

	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{
			Message: QuestionEksAuditFhEncryptionKeyArn,
			Default: input.FirehoseEncryptionKeyArn,
		},
		Checks:   []*bool{&input.FirehoseEncryptionEnabled},
		Opts:     []survey.AskOpt{},
		Response: &input.FirehoseEncryptionKeyArn,
	})
	return err
}

func promptAwsEksAuditExistingCloudwatchIamQuestions(input *aws_eks_audit.
	GenerateAwsEksAuditTfConfigurationArgs) error {

	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{
			Message: QuestionEksAuditExistingCwIamArn,
			Default: input.ExistingCloudWatchIamRoleArn,
		},
		Opts:     []survey.AskOpt{survey.WithValidator(validateAwsArnFormat)},
		Response: &input.ExistingCloudWatchIamRoleArn,
		Required: true,
	})
	return err
}

func promptAwsEksAuditSnsQuestions(input *aws_eks_audit.
	GenerateAwsEksAuditTfConfigurationArgs) error {

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Confirm{
			Message: QuestionEksAuditSnsEncryption,
			Default: input.SnsTopicEncryptionEnabled,
		},
		Opts:     []survey.AskOpt{},
		Response: &input.SnsTopicEncryptionEnabled,
		Required: true,
	}); err != nil {
		return err
	}

	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{
			Message: QuestionEksAuditSnsEncryptionKeyArn,
			Default: input.SnsTopicEncryptionKeyArn,
		},
		Checks:   []*bool{&input.SnsTopicEncryptionEnabled},
		Opts:     []survey.AskOpt{},
		Response: &input.SnsTopicEncryptionKeyArn,
	})
	return err
}

func promptAwsEksAuditCustomIntegrationName(input *aws_eks_audit.
	GenerateAwsEksAuditTfConfigurationArgs) error {

	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{
			Message: QuestionEksAuditCustomIntegrationName,
			Default: input.EksAuditIntegrationName,
		},
		Opts:     []survey.AskOpt{},
		Response: &input.EksAuditIntegrationName,
	})
	return err
}

func promptAwsEksAuditAdditionalClusterRegionQuestions(
	config *aws_eks_audit.GenerateAwsEksAuditTfConfigurationArgs,
	extraState *AwsEksAuditGenerateCommandExtraState,
) error {
	// For each region, collect which clusters to integrate with
	askAgain := false
	if cli.InteractiveMode() {
		askAgain = true
	}

	if config.ParsedRegionClusterMap == nil {
		config.ParsedRegionClusterMap = make(map[string][]string)
	}

	// If there are existing region clusters configured (i.e., from the CLI) display them and ask if they want to add more
	if len(config.ParsedRegionClusterMap) > 0 {
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Prompt: &survey.Confirm{
				Message: fmt.Sprintf(
					QuestionEksAuditRegionClusterCurrent,
					config.ParsedRegionClusterMap,
				),
			},
			Response: &askAgain}); err != nil {
			return err
		}
	}

	// If we already have more than 1 region, don't bother asking the user if it's
	// multi region and instead just set MultiRegion to true
	if len(config.ParsedRegionClusterMap) > 1 {
		extraState.MultiRegion = true
	}

	// If only 1 region has been configured and the user wishes to add more clusters,
	// ask if they want this be to multi region
	if len(config.ParsedRegionClusterMap) <= 1 && askAgain && !extraState.MultiRegion {
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Prompt: &survey.Confirm{
				Message: QuestionEksAuditMultiRegion,
				Default: extraState.MultiRegion,
			},
			Checks:   []*bool{&askAgain},
			Opts:     []survey.AskOpt{},
			Required: true,
			Response: &extraState.MultiRegion,
		}); err != nil {
			return err
		}
	}

	// For each region to add, collect the list of clusters to integrate with
	for askAgain {
		var awsEksAuditRegion string
		var awsEksAuditRegionClusters string
		regionClustersQuestions := []SurveyQuestionWithValidationArgs{
			{
				Prompt:   &survey.Input{Message: QuestionEksAuditRegion},
				Opts:     []survey.AskOpt{survey.WithValidator(validateAwsRegion)},
				Required: true,
				Response: &awsEksAuditRegion,
			},
			{
				Prompt:   &survey.Input{Message: QuestionEksAuditRegionClusters},
				Opts:     []survey.AskOpt{},
				Required: true,
				Response: &awsEksAuditRegionClusters,
			},
		}

		if err := SurveyMultipleQuestionWithValidation(regionClustersQuestions); err != nil {
			return err
		}

		// append region clusters in case the user has input a region more than once
		config.ParsedRegionClusterMap[awsEksAuditRegion] = append(config.ParsedRegionClusterMap[awsEksAuditRegion], strings.Split(awsEksAuditRegionClusters, ",")...)

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Prompt:   &survey.Confirm{Message: QuestionEksAuditAdditionalRegion},
			Checks:   []*bool{&extraState.MultiRegion},
			Response: &askAgain}); err != nil {
			return err
		}

		if !extraState.MultiRegion {
			askAgain = false
		}
	}

	return nil
}

func promptCustomizeEksAuditOutputLocation(extraState *AwsEksAuditGenerateCommandExtraState) error {
	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Input{Message: QuestionEksAuditCustomizeOutputLocation,
			Default: extraState.Output},
		Response: &extraState.Output,
		Opts:     []survey.AskOpt{survey.WithValidator(validPathExists)},
		Required: true,
	})

	return err
}

func askAdvancedEksAuditOptions(config *aws_eks_audit.GenerateAwsEksAuditTfConfigurationArgs,
	extraState *AwsEksAuditGenerateCommandExtraState) error {
	answer := ""

	// Prompt for options
	for answer != AwsAdvancedOptDone {
		// Construction of this slice is a bit strange at first look, but the reason for that is because we have to do string
		// validation to know which option was selected due to how survey works; and doing it by index (also supported) is
		// difficult when the options are dynamic (which they are)
		//
		// Only ask about more accounts if consolidated cloudtrail is set up (matching scenario's doc)
		var options []string

		options = append(options,
			EksAuditConfigureBucket,
			EksAuditExistingCaIamRole,
			EksAuditConfigureFh,
			EksAuditExistingCwIamRole,
			EksAuditConfigureSns,
			EksAuditIntegrationNameOpt,
			EksAuditAdvancedOptLocation,
			EksAuditAdvancedOptDone)
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
		case EksAuditConfigureBucket:
			if err := promptAwsEksAuditBucketQuestions(config); err != nil {
				return err
			}
		case EksAuditExistingCaIamRole:
			if err := promptAwsEksAuditExistingCrossAccountIamQuestions(config); err != nil {
				return err
			}
		case EksAuditConfigureFh:
			if err := promptAwsEksAuditFirehoseQuestions(config); err != nil {
				return err
			}
		case EksAuditExistingCwIamRole:
			if err := promptAwsEksAuditExistingCloudwatchIamQuestions(config); err != nil {
				return err
			}
		case EksAuditConfigureSns:
			if err := promptAwsEksAuditSnsQuestions(config); err != nil {
				return err
			}
		case EksAuditIntegrationNameOpt:
			if err := promptAwsEksAuditCustomIntegrationName(config); err != nil {
				return err
			}
		case EksAuditAdvancedOptLocation:
			if err := promptCustomizeEksAuditOutputLocation(extraState); err != nil {
				return err
			}
		}

		// Re-prompt if not done
		innerAskAgain := true
		if answer == EksAuditAdvancedOptDone {
			innerAskAgain = false
		}

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Checks:   []*bool{&innerAskAgain},
			Prompt:   &survey.Confirm{Message: QuestionEksAuditAnotherAdvancedOpt, Default: false},
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

func eksAuditConfigIsEmpty(g *aws_eks_audit.GenerateAwsEksAuditTfConfigurationArgs) bool {
	return g.AwsProfile == "" &&
		len(g.ParsedRegionClusterMap) == 0 &&
		g.ExistingCrossAccountIamRole == nil &&
		g.LaceworkProfile == ""
}

func writeAwsEksAuditGenerationArgsCache(a *aws_eks_audit.GenerateAwsEksAuditTfConfigurationArgs) {
	if !eksAuditConfigIsEmpty(a) {
		// If ExistingIamRole is partially set, don't write this to cache; the values won't work when loaded
		if a.ExistingCrossAccountIamRole.IsPartial() {
			a.ExistingCrossAccountIamRole = nil
		}
		cli.WriteAssetToCache(CachedAssetAwsEksAuditIacParams, time.Now().Add(time.Hour*1), a)
	}
}

// entry point for launching a survey to build out the required generation parameters
func promptAwsEksAuditGenerate(config *aws_eks_audit.GenerateAwsEksAuditTfConfigurationArgs, existingIam *aws_eks_audit.ExistingCrossAccountIamRoleDetails, extraState *AwsEksAuditGenerateCommandExtraState) error {
	// Cache for later use if generation is abandoned and in interactive mode
	if cli.InteractiveMode() {
		defer writeAwsEksAuditGenerationArgsCache(config)
		defer extraState.writeCache()
	}

	// Set ExistingCrossAccountIamRole details, if provided as cli flags; otherwise don't initialize
	if existingIam.Arn != "" ||
		existingIam.ExternalId != "" {
		config.ExistingCrossAccountIamRole = existingIam
	}

	// These are the core questions that should be asked.
	if err := promptAwsEksAuditAdditionalClusterRegionQuestions(config, extraState); err != nil {
		return err
	}

	// Find out if the customer wants to specify more advanced features
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Confirm{Message: QuestionEksAuditConfigureAdvanced,
			Default: extraState.AskAdvanced},
		Response: &extraState.AskAdvanced,
	}); err != nil {
		return err
	}

	// Keep prompting for advanced options until the say done
	if extraState.AskAdvanced {
		if err := askAdvancedEksAuditOptions(config, extraState); err != nil {
			return err
		}
	}

	return nil
}
