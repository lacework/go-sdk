package cmd

import (
	"regexp"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/lwgenerate/gcp"
)

var (
	// Define question text here to be reused in testing
	QuestionGcpEnableAgentless         = "Enable Agentless integration?"
	QuestionGcpEnableConfiguration     = "Enable Configuration integration?"
	QuestionGcpEnableAuditLog          = "Enable Audit Log integration?"
	QuestionUsePubSubAudit             = "Use Pub Sub Audit Log?"
	QuestionGcpOrganizationIntegration = "Organization integration?"
	QuestionGcpOrganizationID          = "Specify the GCP organization ID:"
	QuestionGcpProjectID               = "Specify the project ID to be used to provision Lacework resources:"
	QuestionGcpServiceAccountCredsPath = "Specify service account credentials JSON path: (optional)"

	QuestionGcpConfigureAdvanced             = "Configure advanced integration options?"
	GcpAdvancedOptExistingServiceAccount     = "Configure & use existing service account"
	QuestionExistingServiceAccountName       = "Specify an existing service account name:"
	QuestionExistingServiceAccountPrivateKey = "Specify an existing service account private key (base64 encoded):"

	GcpAdvancedOptAgentless      = "Configure additional Agentless options"
	QuestionGcpProjectFilterList = "Specify a comma separated list of Google Cloud projects that " +
		"you want to monitor: (optional)"
	QuestionGcpRegions = "Specify a comma separated list of regions to deploy Agentless:"

	GcpAdvancedOptAuditLog        = "Configure additional Audit Log options"
	QuestionGcpUseExistingBucket  = "Use an existing bucket?"
	QuestionGcpExistingBucketName = "Specify an existing bucket name:"
	QuestionGcpConfigureNewBucket = "Configure settings for new bucket?"
	QuestionGcpBucketRegion       = "Specify the bucket region: (optional)"
	QuestionGcpCustomBucketName   = "Specify a custom bucket name: (optional)"
	QuestionGcpBucketLifecycle    = "Specify the bucket lifecycle rule age: (optional)"
	QuestionGcpEnableUBLA         = "Enable uniform bucket level access(UBLA)?"
	QuestionGcpUseExistingSink    = "Use an existing sink?"
	QuestionGcpExistingSinkName   = "Specify the existing sink name"

	GcpAdvancedOptIntegrationName           = "Customize integration name(s)"
	QuestionGcpConfigurationIntegrationName = "Specify a custom configuration integration name: (optional)"
	QuestionGcpAuditLogIntegrationName      = "Specify a custom Audit Log integration name: (optional)"

	QuestionGcpAnotherAdvancedOpt      = "Configure another advanced integration option"
	GcpAdvancedOptLocation             = "Customize output location"
	GcpAdvancedOptProjects             = "Configure multiple projects"
	QuestionGcpCustomizeOutputLocation = "Provide the location for the output to be written:"
	QuestionGcpCustomizeProjects       = "Provide comma separated list of project ID"
	QuestionGcpCustomFilter            = "Specify a custom Audit Log filter which supersedes all other filter options"
	GcpAdvancedOptDone                 = "Done"

	// GcpRegionRegex regex used for validating region input
	GcpRegionRegex = `(asia|australia|europe|northamerica|southamerica|us)-(central|(north|south)?(east|west)?)\d`

	GenerateGcpCommandState                  = &gcp.GenerateGcpTfConfigurationArgs{}
	GenerateGcpExistingServiceAccountDetails = &gcp.ExistingServiceAccountDetails{}
	GenerateGcpCommandExtraState             = &GcpGenerateCommandExtraState{}
	CachedGcpAssetIacParams                  = "iac-gcp-generate-params"
	CachedAssetGcpExtraState                 = "iac-gcp-extra-state"

	InvalidProjectIDMessage = "invalid GCP project ID. " +
		"It must be 6 to 30 lowercase ASCII letters, digits, or hyphens. " +
		"It must start with a letter. Trailing hyphens are prohibited. Example: tokyo-rain-123"

	// gcp command is used to generate TF code for gcp
	generateGcpTfCommand = &cobra.Command{
		Use:   "gcp",
		Short: "Generate and/or execute Terraform code for GCP integration",
		Long: `Use this command to generate Terraform code for deploying Lacework into an GCP environment.

By default, this command interactively prompts for the required information to setup the new cloud account.
In interactive mode, this command will:

* Prompt for the required information to setup the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version is verified as compatible for use
  * If Terraform is not installed, or the version installed is not compatible, a new version will be
    installed into a temporary location
  * Once Terraform is detected or installed, Terraform plan will be executed
  * The command will prompt with the outcome of the plan and allow to view more details or continue with
    Terraform apply
  * If confirmed, Terraform apply will be run, completing the setup of the cloud account

This command can also be run in noninteractive mode.
See help output for more details on the parameter value(s) required for Terraform code generation.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Generate TF Code
			cli.StartProgress("Generating Terraform Code...")

			// Explicitly set Lacework profile if it was passed in main args
			if cli.Profile != "default" {
				GenerateGcpCommandState.LaceworkProfile = cli.Profile
			}

			// Setup modifiers for NewTerraform constructor
			mods := []gcp.GcpTerraformModifier{
				gcp.WithGcpServiceAccountCredentials(GenerateGcpCommandState.ServiceAccountCredentials),
				gcp.WithOrganizationId(GenerateGcpCommandState.GcpOrganizationId),
				gcp.WithProjectId(GenerateGcpCommandState.GcpProjectId),
				gcp.WithExistingServiceAccount(GenerateGcpCommandState.ExistingServiceAccount),
				gcp.WithConfigurationIntegrationName(GenerateGcpCommandState.ConfigurationIntegrationName),
				gcp.WithAuditLogLabels(GenerateGcpCommandState.AuditLogLabels),
				gcp.WithBucketLabels(GenerateGcpCommandState.BucketLabels),
				gcp.WithPubSubSubscriptionLabels(GenerateGcpCommandState.PubSubSubscriptionLabels),
				gcp.WithPubSubTopicLabels(GenerateGcpCommandState.PubSubTopicLabels),
				gcp.WithCustomBucketName(GenerateGcpCommandState.CustomBucketName),
				gcp.WithBucketRegion(GenerateGcpCommandState.BucketRegion),
				gcp.WithExistingLogBucketName(GenerateGcpCommandState.ExistingLogBucketName),
				gcp.WithExistingLogSinkName(GenerateGcpCommandState.ExistingLogSinkName),
				gcp.WithAuditLogIntegrationName(GenerateGcpCommandState.AuditLogIntegrationName),
				gcp.WithLaceworkProfile(GenerateGcpCommandState.LaceworkProfile),
				gcp.WithLogBucketLifecycleRuleAge(GenerateGcpCommandState.LogBucketLifecycleRuleAge),
				gcp.WithFoldersToInclude(GenerateGcpCommandState.FoldersToInclude),
				gcp.WithFoldersToExclude(GenerateGcpCommandState.FoldersToExclude),
				gcp.WithCustomFilter(GenerateGcpCommandState.CustomFilter),
				gcp.WithGoogleWorkspaceFilter(GenerateGcpCommandState.GoogleWorkspaceFilter),
				gcp.WithK8sFilter(GenerateGcpCommandState.K8sFilter),
				gcp.WithPrefix(GenerateGcpCommandState.Prefix),
				gcp.WithWaitTime(GenerateGcpCommandState.WaitTime),
				gcp.WithEnableUBLA(GenerateGcpCommandState.EnableUBLA),
				gcp.WithMultipleProject(GenerateGcpCommandState.Projects),
				gcp.WithProjectFilterList(GenerateGcpCommandState.ProjectFilterList),
				gcp.WithRegions(GenerateGcpCommandState.Regions),
			}

			if GenerateGcpCommandState.OrganizationIntegration {
				mods = append(mods, gcp.WithOrganizationIntegration(GenerateGcpCommandState.OrganizationIntegration))
			}

			if len(GenerateGcpCommandState.FoldersToExclude) > 0 {
				mods = append(mods, gcp.WithIncludeRootProjects(GenerateGcpCommandState.IncludeRootProjects))
			}

			// Create new struct
			data := gcp.NewTerraform(
				GenerateGcpCommandState.Agentless,
				GenerateGcpCommandState.Configuration,
				GenerateGcpCommandState.AuditLog,
				GenerateGcpCommandState.UsePubSubAudit,
				mods...)

			// Generate
			hcl, err := data.Generate()
			cli.StopProgress()

			if err != nil {
				return errors.Wrap(err, "failed to generate terraform code")
			}

			// Write-out generated code to location specified
			dirname, _, err := writeGeneratedCodeToLocation(cmd, hcl, "gcp")
			if err != nil {
				return err
			}

			// Prompt to execute
			err = SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
				Prompt:   &survey.Confirm{Default: GenerateGcpCommandExtraState.TerraformApply, Message: QuestionRunTfPlan},
				Response: &GenerateGcpCommandExtraState.TerraformApply,
			})

			if err != nil {
				return errors.Wrap(err, "failed to prompt for terraform execution")
			}

			locationDir, _ := determineOutputDirPath(dirname, "gcp")
			if GenerateGcpCommandExtraState.TerraformApply {
				// Execution pre-run check
				err := executionPreRunChecks(dirname, locationDir, "gcp")
				if err != nil {
					return err
				}
			}

			// Output where code was generated
			if !GenerateGcpCommandExtraState.TerraformApply {
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

			// Validate gcp sa credentials file, if passed
			gcpSaCredentials, err := cmd.Flags().GetString("service_account_credentials")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}

			if gcpSaCredentials != "" {
				if err := gcp.ValidateServiceAccountCredentialsFile(gcpSaCredentials); err != nil {
					return err
				}
			}

			// Validate gcp region, if passed
			region, err := cmd.Flags().GetString("bucket_region")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateGcpRegion(region); err != nil {
				return err
			}

			projectId, err := cmd.Flags().GetString("project_id")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if projectId == "" && !cli.InteractiveMode() {
				return errors.New("project_id must be provided")
			}

			// Load any cached inputs if interactive
			if cli.InteractiveMode() {
				cachedOptions := &gcp.GenerateGcpTfConfigurationArgs{}
				iacParamsExpired := cli.ReadCachedAsset(CachedGcpAssetIacParams, &cachedOptions)
				if iacParamsExpired {
					cli.Log.Debug("loaded previously set values for GCP iac generation")
				}

				extraState := &GcpGenerateCommandExtraState{}
				extraStateParamsExpired := cli.ReadCachedAsset(CachedAssetGcpExtraState, &extraState)
				if extraStateParamsExpired {
					cli.Log.Debug("loaded previously set values for GCP iac generation (extra state)")
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
					if err := mergo.Merge(GenerateGcpCommandState, cachedOptions); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
					if err := mergo.Merge(GenerateGcpCommandExtraState, extraState); err != nil {
						return errors.Wrap(err, "failed to load saved options")
					}
				}
			}

			// Collect and/or confirm parameters
			if err := promptGcpGenerate(
				GenerateGcpCommandState,
				GenerateGcpExistingServiceAccountDetails,
				GenerateGcpCommandExtraState,
			); err != nil {
				return errors.Wrap(err, "collecting/confirming parameters")
			}

			return nil
		},
	}
)

type GcpGenerateCommandExtraState struct {
	AskAdvanced                bool
	Output                     string
	ConfigureNewBucketSettings bool
	UseExistingServiceAccount  bool
	UseExistingBucket          bool
	UseExistingSink            bool
	TerraformApply             bool
}

func (gcp *GcpGenerateCommandExtraState) isEmpty() bool {
	return gcp.Output == "" &&
		!gcp.AskAdvanced &&
		!gcp.UseExistingServiceAccount &&
		!gcp.UseExistingBucket &&
		!gcp.UseExistingSink &&
		!gcp.TerraformApply
}

// Flush current state of the struct to disk, provided it's not empty
func (gcp *GcpGenerateCommandExtraState) writeCache() {
	if !gcp.isEmpty() {
		cli.WriteAssetToCache(CachedAssetGcpExtraState, time.Now().Add(time.Hour*1), gcp)
	}
}

func initGenerateGcpTfCommandFlags() {
	// add flags to sub commands
	// TODO Share the help with the interactive generation
	generateGcpTfCommand.PersistentFlags().BoolVar(
		&GenerateGcpCommandState.Agentless,
		"agentless",
		false,
		"enable agentless integration")
	generateGcpTfCommand.PersistentFlags().BoolVar(
		&GenerateGcpCommandState.Agentless,
		"audit_log",
		false,
		"enable audit log integration")
	generateGcpTfCommand.PersistentFlags().BoolVar(
		&GenerateGcpCommandState.Configuration,
		"configuration",
		false,
		"enable configuration integration")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.ServiceAccountCredentials,
		"service_account_credentials",
		"",
		"specify service account credentials JSON file path (leave blank to make use of google credential ENV vars)")
	generateGcpTfCommand.PersistentFlags().BoolVar(
		&GenerateGcpCommandState.OrganizationIntegration,
		"organization_integration",
		false,
		"enable organization integration")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.GcpOrganizationId,
		"organization_id",
		"",
		"specify the organization id (only set if organization_integration is set)")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.GcpProjectId,
		"project_id",
		"",
		"specify the project id to be used to provision lacework resources (required)")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpExistingServiceAccountDetails.Name,
		"existing_service_account_name",
		"",
		"specify existing service account name")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpExistingServiceAccountDetails.PrivateKey,
		"existing_service_account_private_key",
		"",
		"specify existing service account private key (base64 encoded)")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.ConfigurationIntegrationName,
		"configuration_integration_name",
		"",
		"specify a custom configuration integration name")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.CustomBucketName,
		"custom_bucket_name",
		"",
		"override prefix based storage bucket name generation with a custom name")
	// TODO: Implement AuditLogLabels, BucketLabels, PubSubSubscriptionLabels & PubSubTopicLabels
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.BucketRegion,
		"bucket_region",
		"",
		"specify bucket region")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.ExistingLogBucketName,
		"existing_bucket_name",
		"",
		"specify existing bucket name")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.ExistingLogSinkName,
		"existing_sink_name",
		"",
		"specify existing sink name")
	generateGcpTfCommand.PersistentFlags().StringSliceVar(
		&GenerateGcpCommandState.ProjectFilterList,
		"project_filter_list",
		[]string{},
		"List of GCP project IDs to monitor for Agentless integration")
	generateGcpTfCommand.PersistentFlags().StringSliceVar(
		&GenerateGcpCommandState.Regions,
		"regions",
		[]string{},
		"List of GCP regions to deploy for Agentless integration")

	// DEPRECATED
	generateGcpTfCommand.PersistentFlags().BoolVar(
		&GenerateGcpCommandState.EnableForceDestroyBucket,
		"enable_force_destroy_bucket",
		true,
		"enable force bucket destroy")
	errcheckWARN(generateGcpTfCommand.PersistentFlags().MarkDeprecated(
		"enable_force_destroy_bucket", "by default, force destroy is enabled.",
	))
	// ---

	generateGcpTfCommand.PersistentFlags().BoolVar(
		&GenerateGcpCommandState.EnableUBLA,
		"enable_ubla",
		true,
		"enable universal bucket level access(ubla)")
	generateGcpTfCommand.PersistentFlags().IntVar(
		&GenerateGcpCommandState.LogBucketLifecycleRuleAge,
		"bucket_lifecycle_rule_age",
		-1,
		"specify the lifecycle rule age")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.CustomFilter,
		"custom_filter",
		"",
		"Audit Log filter which supersedes all other filter options when defined")
	generateGcpTfCommand.PersistentFlags().BoolVar(
		&GenerateGcpCommandState.GoogleWorkspaceFilter,
		"google_workspace_filter",
		true,
		"filter out Google Workspace login logs from GCP Audit Log sinks")
	generateGcpTfCommand.PersistentFlags().BoolVar(
		&GenerateGcpCommandState.K8sFilter,
		"k8s_filter",
		true,
		"filter out GKE logs from GCP Audit Log sinks")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.Prefix,
		"prefix",
		"",
		"prefix that will be used at the beginning of every generated resource")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.WaitTime,
		"wait_time",
		"",
		"amount of time to wait before the next resource is provisioned")
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandState.AuditLogIntegrationName,
		"audit_log_integration_name",
		"",
		"specify a custom audit log integration name")
	generateGcpTfCommand.PersistentFlags().StringArrayVarP(
		&GenerateGcpCommandState.FoldersToExclude,
		"folders_to_exclude",
		"e",
		[]string{},
		"List of root folders to exclude for an organization-level integration")
	generateGcpTfCommand.PersistentFlags().BoolVar(
		&GenerateGcpCommandState.IncludeRootProjects,
		"include_root_projects",
		true,
		"Disables logic that includes root-level projects if excluding folders")
	generateGcpTfCommand.PersistentFlags().StringArrayVarP(
		&GenerateGcpCommandState.FoldersToInclude,
		"folders_to_include",
		"i",
		[]string{},
		"list of root folders to include for an organization-level integration")
	generateGcpTfCommand.PersistentFlags().BoolVar(
		&GenerateGcpCommandExtraState.TerraformApply,
		"apply",
		false,
		"run terraform apply without executing plan or prompting",
	)
	generateGcpTfCommand.PersistentFlags().StringVar(
		&GenerateGcpCommandExtraState.Output,
		"output",
		"",
		"location to write generated content (default is ~/lacework/gcp)",
	)
	generateGcpTfCommand.PersistentFlags().BoolVar(
		&GenerateGcpCommandState.UsePubSubAudit,
		"use_pub_sub",
		false,
		"use pub/sub for the audit log data rather than bucket")
	generateGcpTfCommand.PersistentFlags().StringSliceVar(
		&GenerateGcpCommandState.Projects,
		"projects",
		[]string{},
		"list of project IDs to integrate with (project-level integrations)")
}

// survey.Validator for gcp region
func validateGcpRegion(val interface{}) error {
	switch value := val.(type) {
	case string:
		// as this field is optional, it is valid for this field to be empty
		if value != "" {
			// if value doesn't match regex, return invalid arn
			ok, err := regexp.MatchString(GcpRegionRegex, value)
			if err != nil {
				return errors.Wrap(err, "failed to validate input")
			}

			if !ok {
				return errors.New("invalid region name supplied")
			}
		}
	default:
		// if the value passed is not a string
		return errors.New("value must be a string")
	}

	return nil
}

func promptGcpAgentlessQuestions(
	config *gcp.GenerateGcpTfConfigurationArgs,
	extraState *GcpGenerateCommandExtraState,
) error {
	projectFilterListInput := ""

	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionGcpProjectFilterList, Default: strings.Join(config.ProjectFilterList, ",")},
			Response: &projectFilterListInput,
		},
	}, config.Agentless)

	if projectFilterListInput != "" {
		config.ProjectFilterList = strings.Split(projectFilterListInput, ",")
	}

	return err
}

func promptGcpAuditLogQuestions(
	config *gcp.GenerateGcpTfConfigurationArgs,
	extraState *GcpGenerateCommandExtraState,
) error {

	// Only ask these questions if configure audit log is true
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionUsePubSubAudit, Default: config.UsePubSubAudit},
			Checks:   []*bool{&config.AuditLog},
			Response: &config.UsePubSubAudit,
		},
	}, config.AuditLog); err != nil {
		return err
	}
	// Present the user with Bucket Configuration options, if required
	if err := promptGcpBucketConfiguration(config, extraState); err != nil {
		return err
	}
	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionGcpUseExistingSink, Default: extraState.UseExistingSink},
			Checks:   []*bool{&config.AuditLog},
			Required: true,
			Response: &extraState.UseExistingSink,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpExistingSinkName, Default: config.ExistingLogSinkName},
			Checks:   []*bool{&config.AuditLog, &extraState.UseExistingSink},
			Required: true,
			Response: &config.ExistingLogSinkName,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpCustomFilter, Default: config.CustomFilter},
			Checks:   []*bool{&config.AuditLog},
			Response: &config.CustomFilter,
		},
	}, config.AuditLog)

	return err
}

func promptGcpBucketConfiguration(
	config *gcp.GenerateGcpTfConfigurationArgs, extraState *GcpGenerateCommandExtraState,
) error {
	// Prompt to configure bucket information (not required when using the Pub Sub Audit Log)
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionGcpUseExistingBucket, Default: extraState.UseExistingBucket},
			Checks:   []*bool{&config.AuditLog, usePubSubActivityDisabled(config)},
			Response: &extraState.UseExistingBucket,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpExistingBucketName, Default: config.ExistingLogBucketName},
			Checks:   []*bool{&config.AuditLog, &extraState.UseExistingBucket, usePubSubActivityDisabled(config)},
			Required: true,
			Response: &config.ExistingLogBucketName,
		},
	}, config.AuditLog); err != nil {
		return err
	}

	newBucket := !extraState.UseExistingBucket
	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionGcpConfigureNewBucket, Default: extraState.ConfigureNewBucketSettings},
			Checks:   []*bool{&config.AuditLog, &newBucket, usePubSubActivityDisabled(config)},
			Required: true,
			Response: &extraState.ConfigureNewBucketSettings,
		},
		{
			Prompt: &survey.Input{Message: QuestionGcpBucketRegion, Default: config.BucketRegion},
			Checks: []*bool{&config.AuditLog,
				&newBucket,
				&extraState.ConfigureNewBucketSettings,
				usePubSubActivityDisabled(config)},
			Opts:     []survey.AskOpt{survey.WithValidator(validateGcpRegion)},
			Response: &config.BucketRegion,
		},
		{
			Prompt: &survey.Input{Message: QuestionGcpCustomBucketName, Default: config.CustomBucketName},
			Checks: []*bool{&config.AuditLog,
				&newBucket,
				&extraState.ConfigureNewBucketSettings,
				usePubSubActivityDisabled(config)},
			Response: &config.CustomBucketName,
		},
		{
			Prompt: &survey.Input{Message: QuestionGcpBucketLifecycle, Default: "-1"},
			Checks: []*bool{&config.AuditLog,
				&newBucket,
				&extraState.ConfigureNewBucketSettings,
				usePubSubActivityDisabled(config)},
			Response: &config.LogBucketLifecycleRuleAge,
		},
		{
			Prompt: &survey.Confirm{Message: QuestionGcpEnableUBLA, Default: config.EnableUBLA},
			Checks: []*bool{&config.AuditLog,
				&newBucket,
				&extraState.ConfigureNewBucketSettings,
				usePubSubActivityDisabled(config)},
			Required: true,
			Response: &config.EnableUBLA,
		},
	}, config.AuditLog)

	return err
}

func usePubSubActivityDisabled(config *gcp.GenerateGcpTfConfigurationArgs) *bool {
	usePubSubActivityDisabled := !config.UsePubSubAudit
	return &usePubSubActivityDisabled
}
func promptGcpExistingServiceAccountQuestions(config *gcp.GenerateGcpTfConfigurationArgs) error {
	// ensure struct is initialized
	if config.ExistingServiceAccount == nil {
		config.ExistingServiceAccount = &gcp.ExistingServiceAccountDetails{}
	}

	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionExistingServiceAccountName, Default: config.ExistingServiceAccount.Name},
			Response: &config.ExistingServiceAccount.Name,
			Opts:     []survey.AskOpt{survey.WithValidator(survey.Required)},
		},
		{
			Prompt: &survey.Input{
				Message: QuestionExistingServiceAccountPrivateKey,
				Default: config.ExistingServiceAccount.PrivateKey,
			},
			Response: &config.ExistingServiceAccount.PrivateKey,
			Opts: []survey.AskOpt{
				survey.WithValidator(survey.Required),
				survey.WithValidator(gcp.ValidateStringIsBase64),
			},
		}})

	return err
}

func promptGcpIntegrationNameQuestions(config *gcp.GenerateGcpTfConfigurationArgs) error {
	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt: &survey.Input{
				Message: QuestionGcpConfigurationIntegrationName,
				Default: config.ConfigurationIntegrationName,
			},
			Checks:   []*bool{&config.Configuration},
			Response: &config.ConfigurationIntegrationName,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpAuditLogIntegrationName, Default: config.AuditLogIntegrationName},
			Checks:   []*bool{&config.AuditLog},
			Response: &config.AuditLogIntegrationName,
		}})

	return err
}

func promptCustomizeGcpOutputLocation(extraState *GcpGenerateCommandExtraState) error {
	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionGcpCustomizeOutputLocation, Default: extraState.Output},
		Response: &extraState.Output,
		Opts:     []survey.AskOpt{survey.WithValidator(validPathExists)},
		Required: true,
	})

	return err
}

func promptCustomizeGcpProjects(config *gcp.GenerateGcpTfConfigurationArgs) error {

	validation := func(val interface{}) error {
		switch value := val.(type) {
		case string:
			for _, id := range strings.Split(value, ",") {
				err := validateGcpProjectId(strings.TrimSpace(id))
				if err != nil {
					return err
				}
			}
		default:
			return errors.New("value must be a string")
		}

		return nil
	}

	var projects string

	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionGcpCustomizeProjects},
		Response: &projects,
		Opts:     []survey.AskOpt{survey.WithValidator(validation)},
		Required: true,
	})

	if err != nil {
		return err
	}

	for _, id := range strings.Split(projects, ",") {
		config.Projects = append(config.Projects, strings.TrimSpace(id))
	}

	return nil
}

func askAdvancedOptions(config *gcp.GenerateGcpTfConfigurationArgs, extraState *GcpGenerateCommandExtraState) error {
	answer := ""

	// Prompt for options
	for answer != GcpAdvancedOptDone {
		// Construction of this slice is a bit strange at first look, but the reason for that is because we have to do
		// string validation to know which option was selected due to how survey works; and doing it by index (also
		// supported) is difficult when the options are dynamic (which they are)
		var options []string

		// Only show Advanced Agentless options if Agentless integration is set to true
		if config.Agentless {
			options = append(options, GcpAdvancedOptAgentless)
		}

		// Only show Advanced AuditLog options if AuditLog integration is set to true
		if config.AuditLog {
			options = append(options, GcpAdvancedOptAuditLog)
		}

		options = append(options,
			GcpAdvancedOptExistingServiceAccount,
			GcpAdvancedOptIntegrationName,
			GcpAdvancedOptLocation,
			GcpAdvancedOptProjects,
			GcpAdvancedOptDone)
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
		case GcpAdvancedOptAgentless:
			if err := promptGcpAgentlessQuestions(config, extraState); err != nil {
				return err
			}
		case GcpAdvancedOptAuditLog:
			if err := promptGcpAuditLogQuestions(config, extraState); err != nil {
				return err
			}
		case GcpAdvancedOptExistingServiceAccount:
			if err := promptGcpExistingServiceAccountQuestions(config); err != nil {
				return err
			}
		case GcpAdvancedOptIntegrationName:
			if err := promptGcpIntegrationNameQuestions(config); err != nil {
				return err
			}
		case GcpAdvancedOptLocation:
			if err := promptCustomizeGcpOutputLocation(extraState); err != nil {
				return err
			}
		case GcpAdvancedOptProjects:
			if err := promptCustomizeGcpProjects(config); err != nil {
				return err
			}
		}

		// Re-prompt if not done
		innerAskAgain := true
		if answer == GcpAdvancedOptDone {
			innerAskAgain = false
		}

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Checks:   []*bool{&innerAskAgain},
			Prompt:   &survey.Confirm{Message: QuestionGcpAnotherAdvancedOpt, Default: false},
			Response: &innerAskAgain,
		}); err != nil {
			return err
		}

		if !innerAskAgain {
			answer = GcpAdvancedOptDone
		}
	}

	return nil
}

func gcpConfigIsEmpty(g *gcp.GenerateGcpTfConfigurationArgs) bool {
	return !g.Agentless &&
		!g.AuditLog &&
		!g.Configuration &&
		g.ServiceAccountCredentials == "" &&
		g.GcpOrganizationId == "" &&
		g.LaceworkProfile == ""
}

func writeGcpGenerationArgsCache(a *gcp.GenerateGcpTfConfigurationArgs) {
	if !gcpConfigIsEmpty(a) {
		// If ExistingServiceAccount is partially set, don't write this to cache; the values won't work when loaded
		if a.ExistingServiceAccount.IsPartial() {
			a.ExistingServiceAccount = nil
		}
		cli.WriteAssetToCache(CachedGcpAssetIacParams, time.Now().Add(time.Hour*1), a)
	}
}

// entry point for launching a survey to build out the required generation parameters
func promptGcpGenerate(
	config *gcp.GenerateGcpTfConfigurationArgs,
	existingServiceAccount *gcp.ExistingServiceAccountDetails,
	extraState *GcpGenerateCommandExtraState,
) error {
	// Cache for later use if generation is abandon and in interactive mode
	if cli.InteractiveMode() {
		defer writeGcpGenerationArgsCache(config)
		defer extraState.writeCache()
	}

	// Set ExistingIamRole details, if provided as cli flags; otherwise don't initialize
	if existingServiceAccount.Name != "" ||
		existingServiceAccount.PrivateKey != "" {
		config.ExistingServiceAccount = existingServiceAccount
	}

	// These are the core questions that should be asked.
	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Prompt:   &survey.Confirm{Message: QuestionGcpEnableAgentless, Default: config.Agentless},
				Response: &config.Agentless,
			},
			{
				Prompt:   &survey.Confirm{Message: QuestionGcpEnableConfiguration, Default: config.Configuration},
				Response: &config.Configuration,
			},
			{
				Prompt:   &survey.Confirm{Message: QuestionGcpEnableAuditLog, Default: config.AuditLog},
				Response: &config.AuditLog,
			},
		}); err != nil {
		return err
	}

	// Validate one of configuration or audit log was enabled; otherwise error out
	if !config.Agentless && !config.Configuration && !config.AuditLog {
		return errors.New("must enable agentless, audit log or configuration")
	}

	configOrAuditLogEnabled := config.Configuration || config.AuditLog
	regionsInput := ""

	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Prompt:   &survey.Input{Message: QuestionGcpProjectID, Default: config.GcpProjectId},
				Opts:     []survey.AskOpt{survey.WithValidator(validateGcpProjectId)},
				Required: true,
				Response: &config.GcpProjectId,
			},
			{
				Prompt:   &survey.Input{Message: QuestionGcpRegions, Default: strings.Join(config.Regions, ",")},
				Checks:   []*bool{&config.Agentless},
				Response: &regionsInput,
				Required: true,
			},
			{
				Prompt:   &survey.Confirm{Message: QuestionGcpOrganizationIntegration, Default: config.OrganizationIntegration},
				Response: &config.OrganizationIntegration,
			},
			{
				Prompt:   &survey.Input{Message: QuestionGcpOrganizationID, Default: config.GcpOrganizationId},
				Checks:   []*bool{&config.OrganizationIntegration},
				Required: true,
				Response: &config.GcpOrganizationId,
			},
			{
				Prompt:   &survey.Input{Message: QuestionGcpServiceAccountCredsPath, Default: config.ServiceAccountCredentials},
				Opts:     []survey.AskOpt{survey.WithValidator(gcp.ValidateServiceAccountCredentials)},
				Checks:   []*bool{&configOrAuditLogEnabled},
				Response: &config.ServiceAccountCredentials,
			},
		}); err != nil {
		return err
	}

	if regionsInput != "" {
		config.Regions = strings.Split(regionsInput, ",")
	}

	// Find out if the customer wants to specify more advanced features
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Message: QuestionGcpConfigureAdvanced, Default: extraState.AskAdvanced},
		Response: &extraState.AskAdvanced,
	}); err != nil {
		return err
	}

	// Keep prompting for advanced options until the say done
	if extraState.AskAdvanced {
		if err := askAdvancedOptions(config, extraState); err != nil {
			return err
		}
	}

	return nil
}

func validateGcpProjectId(val interface{}) error {
	switch value := val.(type) {

	case string:
		match, err := regexp.MatchString("(^[a-z][a-z0-9-]{4,28}[a-z0-9]$|^$)", value)
		if err != nil {
			return err
		}

		if !match {
			return errors.New(InvalidProjectIDMessage)
		}
	default:
		return errors.New("value must be a string")
	}

	return nil
}
