package cmd

import (
	"regexp"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/v2/lwgenerate/gcp"
)

// Question labels
const (
	IconGcpAgentless     = "[Agentless]"
	IconGcpConfiguration = "[Configuration]"
	IconGcpAuditLog      = "[Audit Log]"
)

var (
	// Define question text here to be reused in testing
	// Core questions
	QuestionGcpEnableAgentless              = "Enable Agentless integration?"
	QuestionGcpEnableConfiguration          = "Enable Configuration integration?"
	QuestionGcpEnableAuditLog               = "Enable Audit Log integration?"
	QuestionGcpOrganizationIntegration      = "Organization integration?"
	QuestionGcpOrganizationID               = "Organization ID:"
	QuestionGcpProjectID                    = "Project ID to be used to provision Lacework resources:"
	QuestionGcpServiceAccountCredsPath      = "Service account credentials JSON path: (optional)"
	QuestionGcpRegions                      = "Comma separated list of regions to deploy Agentless:"
	QuestionGcpProjectFilterList            = "Comma separated list of project IDs to monitor: (optional)"
	QuestionGcpUseExistingSink              = "Use an existing sink?"
	QuestionGcpExistingSinkName             = "Existing sink name:"
	QuestionGcpConfigurationIntegrationName = "Custom Configuration integration name: (optional)"
	QuestionGcpAuditLogIntegrationName      = "Custom Audit Log integration name: (optional)"
	QuestionGcpCustomFilter                 = "Custom Audit Log filter which supersedes other filter options: (optional)"
	QuestionGcpCustomizeOutputLocation      = "Provide the location for the output to be written: (optional)"
	QuestionGcpCustomizeProjects            = "Provide comma separated list of project IDs to deploy: (optional)"

	// Service account questions
	QuestionUseExistingServiceAccount        = "Use existing service account details?"
	QuestionExistingServiceAccountName       = "Existing service account name:"
	QuestionExistingServiceAccountPrivateKey = "Existing service account private key (base64 encoded):"

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

	// List of valid GCP regions
	validGcpRegions = map[string]bool{
		"africa-south1":           true,
		"asia-east1":              true,
		"asia-east2":              true,
		"asia-northeast1":         true,
		"asia-northeast2":         true,
		"asia-northeast3":         true,
		"asia-south1":             true,
		"asia-south2":             true,
		"asia-southeast1":         true,
		"asia-southeast2":         true,
		"australia-southeast1":    true,
		"australia-southeast2":    true,
		"europe-central2":         true,
		"europe-north1":           true,
		"europe-north2":           true,
		"europe-southwest1":       true,
		"europe-west1":            true,
		"europe-west2":            true,
		"europe-west3":            true,
		"europe-west4":            true,
		"europe-west5":            true,
		"europe-west6":            true,
		"europe-west8":            true,
		"europe-west9":            true,
		"europe-west10":           true,
		"europe-west12":           true,
		"me-central1":             true,
		"me-central2":             true,
		"me-west1":                true,
		"northamerica-northeast1": true,
		"northamerica-northeast2": true,
		"northamerica-south1":     true,
		"southamerica-east1":      true,
		"southamerica-west1":      true,
		"us-central1":             true,
		"us-central2":             true,
		"us-east1":                true,
		"us-east4":                true,
		"us-east5":                true,
		"us-south1":               true,
		"us-west1":                true,
		"us-west2":                true,
		"us-west3":                true,
		"us-west4":                true,
	}

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
				gcp.WithPubSubSubscriptionLabels(GenerateGcpCommandState.PubSubSubscriptionLabels),
				gcp.WithPubSubTopicLabels(GenerateGcpCommandState.PubSubTopicLabels),
				gcp.WithExistingLogSinkName(GenerateGcpCommandState.ExistingLogSinkName),
				gcp.WithAuditLogIntegrationName(GenerateGcpCommandState.AuditLogIntegrationName),
				gcp.WithLaceworkProfile(GenerateGcpCommandState.LaceworkProfile),
				gcp.WithFoldersToInclude(GenerateGcpCommandState.FoldersToInclude),
				gcp.WithFoldersToExclude(GenerateGcpCommandState.FoldersToExclude),
				gcp.WithCustomFilter(GenerateGcpCommandState.CustomFilter),
				gcp.WithGoogleWorkspaceFilter(GenerateGcpCommandState.GoogleWorkspaceFilter),
				gcp.WithK8sFilter(GenerateGcpCommandState.K8sFilter),
				gcp.WithPrefix(GenerateGcpCommandState.Prefix),
				gcp.WithWaitTime(GenerateGcpCommandState.WaitTime),
				gcp.WithMultipleProject(GenerateGcpCommandState.Projects),
				gcp.WithProjectFilterList(GenerateGcpCommandState.ProjectFilterList),
				gcp.WithRegions(GenerateGcpCommandState.Regions),
				gcp.WithUsePubSubAudit(true), // always set to true, storage based integration deprecated
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
	Output                    string
	UseExistingServiceAccount bool
	UseExistingSink           bool
	TerraformApply            bool
}

func (gcp *GcpGenerateCommandExtraState) isEmpty() bool {
	return gcp.Output == "" &&
		!gcp.UseExistingServiceAccount &&
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
		&GenerateGcpCommandState.AuditLog,
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
		"specify the organization id (only set if agentless integration or organization_integration is set)")
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

	// ---

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
		true,
		"deprecated: pub/sub audit log integration is always used and only supported type")
	generateGcpTfCommand.PersistentFlags().StringSliceVar(
		&GenerateGcpCommandState.Projects,
		"projects",
		[]string{},
		"list of project IDs to integrate with (project-level integrations)")
}

func validateGcpRegion(val interface{}) error {
	switch value := val.(type) {
	case string:
		regions := strings.Split(value, ",")
		for _, region := range regions {
			region = strings.TrimSpace(region)
			if !validGcpRegions[region] {
				return errors.New("invalid GCP region. Please provide a valid GCP region (e.g., 'us-central1', 'europe-west1')")
			}
		}
	default:
		return errors.New("value must be a string")
	}

	return nil
}

func promptGcpAgentlessQuestions(config *gcp.GenerateGcpTfConfigurationArgs) error {
	regionsInput := ""
	projectFilterListInput := ""

	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon: IconGcpAgentless,
			Prompt: &survey.Input{
				Message: QuestionGcpRegions,
				Default: strings.Join(config.Regions, ","),
				Help:    "Enter a valid GCP region (e.g., 'us-central1', 'europe-west1')",
			},
			Response: &regionsInput,
			Required: true,
			Opts:     []survey.AskOpt{survey.WithValidator(validateGcpRegion)},
		},
		{
			Icon: IconGcpAgentless,
			Prompt: &survey.Input{
				Message: QuestionGcpProjectFilterList,
				Default: strings.Join(config.ProjectFilterList, ","),
			},
			Response: &projectFilterListInput,
		},
	})

	if err != nil {
		return err
	}

	if regionsInput != "" {
		config.Regions = strings.Split(regionsInput, ",")
	}
	if projectFilterListInput != "" {
		config.ProjectFilterList = strings.Split(projectFilterListInput, ",")
	}

	return nil
}

func promptGcpConfigurationQuestions(config *gcp.GenerateGcpTfConfigurationArgs) error {
	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon: IconGcpConfiguration,
			Prompt: &survey.Input{
				Message: QuestionGcpConfigurationIntegrationName,
				Default: config.ConfigurationIntegrationName,
			},
			Response: &config.ConfigurationIntegrationName,
		},
	})

	return err
}

func promptGcpAuditLogQuestions(
	config *gcp.GenerateGcpTfConfigurationArgs,
	extraState *GcpGenerateCommandExtraState,
) error {
	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon: IconGcpAuditLog,
			Prompt: &survey.Confirm{
				Message: QuestionGcpUseExistingSink,
				Default: extraState.UseExistingSink,
			},
			Response: &extraState.UseExistingSink,
		},
		{
			Icon: IconGcpAuditLog,
			Prompt: &survey.Input{
				Message: QuestionGcpExistingSinkName,
				Default: config.ExistingLogSinkName,
			},
			Checks:   []*bool{&extraState.UseExistingSink},
			Response: &config.ExistingLogSinkName,
		},
		{
			Icon: IconGcpAuditLog,
			Prompt: &survey.Input{
				Message: QuestionGcpAuditLogIntegrationName,
				Default: config.AuditLogIntegrationName,
			},
			Response: &config.AuditLogIntegrationName,
		},
		{
			Icon: IconGcpAuditLog,
			Prompt: &survey.Input{
				Message: QuestionGcpCustomFilter,
				Default: config.CustomFilter,
			},
			Response: &config.CustomFilter,
		},
	})

	return err
}

func promptGcpOrganizationQuestions(config *gcp.GenerateGcpTfConfigurationArgs) error {
	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt: &survey.Input{
				Message: QuestionGcpOrganizationID,
				Default: config.GcpOrganizationId,
			},
			Required: true,
			Response: &config.GcpOrganizationId,
		},
	})

	return err
}

func promptGcpServiceAccountQuestions(config *gcp.GenerateGcpTfConfigurationArgs) error {
	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionGcpServiceAccountCredsPath, Default: config.ServiceAccountCredentials},
			Opts:     []survey.AskOpt{survey.WithValidator(gcp.ValidateServiceAccountCredentials)},
			Response: &config.ServiceAccountCredentials,
		},
	})

	return err
}

func promptGcpExistingServiceAccountQuestions(config *gcp.GenerateGcpTfConfigurationArgs) error {
	// ensure struct is initialized
	if config.ExistingServiceAccount == nil {
		config.ExistingServiceAccount = &gcp.ExistingServiceAccountDetails{}
	}

	// Prompt for service account name first (optional)
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt: &survey.Input{
				Message: QuestionExistingServiceAccountName,
				Default: config.ExistingServiceAccount.Name,
			},
			Response: &config.ExistingServiceAccount.Name,
		},
	}); err != nil {
		return err
	}

	// If service account name is set (either from flag or prompt), require private key and validate
	if strings.TrimSpace(config.ExistingServiceAccount.Name) != "" {
		if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
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
			},
		}); err != nil {
			return err
		}
	} else {
		// If name is not set, clear private key as well
		config.ExistingServiceAccount.PrivateKey = ""
	}

	return nil
}

func promptCustomizeGcpOutputLocation(extraState *GcpGenerateCommandExtraState) error {
	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionGcpCustomizeOutputLocation, Default: extraState.Output},
		Response: &extraState.Output,
		Opts:     []survey.AskOpt{survey.WithValidator(validPathExists)},
	})

	return err
}

func promptCustomizeGcpProjects(config *gcp.GenerateGcpTfConfigurationArgs) error {
	// Determine the correct label
	var label string
	if config.Configuration && config.AuditLog {
		label = "[Configuration & AuditLog]"
	} else if config.Configuration {
		label = IconGcpConfiguration
	} else {
		label = IconGcpAuditLog
	}

	validation := func(val interface{}) error {
		switch value := val.(type) {
		case string:
			// If empty string, that's valid (optional field)
			if value == "" {
				return nil
			}
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

	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon:     label,
			Prompt:   &survey.Input{Message: QuestionGcpCustomizeProjects},
			Response: &projects,
			Opts:     []survey.AskOpt{survey.WithValidator(validation)},
		},
	})

	if err != nil {
		return err
	}

	if projects != "" {
		for _, id := range strings.Split(projects, ",") {
			config.Projects = append(config.Projects, strings.TrimSpace(id))
		}
	}

	return nil
}

// writeGcpGenerationArgsCache writes the current state to the cache
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

	// Ask for project ID first as it's required for all integrations
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt: &survey.Input{
				Message: QuestionGcpProjectID,
				Default: config.GcpProjectId,
			},
			Opts:     []survey.AskOpt{survey.WithValidator(validateGcpProjectId)},
			Required: true,
			Response: &config.GcpProjectId,
		},
	}); err != nil {
		return err
	}

	// Ask organization integration
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt: &survey.Confirm{
				Message: QuestionGcpOrganizationIntegration,
				Default: config.OrganizationIntegration,
			},
			Response: &config.OrganizationIntegration,
		},
	}); err != nil {
		return err
	}

	// Ask for organization ID if organization integration is enabled
	if config.OrganizationIntegration {
		if err := promptGcpOrganizationQuestions(config); err != nil {
			return err
		}
	}

	// Ask about each integration type and immediately ask related questions
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon: IconGcpAgentless,
			Prompt: &survey.Confirm{
				Message: QuestionGcpEnableAgentless,
				Default: config.Agentless,
			},
			Response: &config.Agentless,
		},
	}); err != nil {
		return err
	}

	// If Agentless is enabled, ask Agentless-specific questions immediately
	if config.Agentless {
		if err := promptGcpAgentlessQuestions(config); err != nil {
			return err
		}
	}

	// Ask for organization ID if organization integration is not enabled and agentless is enabled
	if !config.OrganizationIntegration && config.Agentless {
		if err := promptGcpOrganizationQuestions(config); err != nil {
			return err
		}
	}

	// Ask about Configuration integration
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon: IconGcpConfiguration,
			Prompt: &survey.Confirm{
				Message: QuestionGcpEnableConfiguration,
				Default: config.Configuration,
			},
			Response: &config.Configuration,
		},
	}); err != nil {
		return err
	}

	// If Configuration is enabled, ask Configuration-specific questions immediately
	if config.Configuration {
		if err := promptGcpConfigurationQuestions(config); err != nil {
			return err
		}
	}

	// Ask about Audit Log integration
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Icon: IconGcpAuditLog,
			Prompt: &survey.Confirm{
				Message: QuestionGcpEnableAuditLog,
				Default: config.AuditLog,
			},
			Response: &config.AuditLog,
		},
	}); err != nil {
		return err
	}

	// If Audit Log is enabled, ask Audit Log-specific questions immediately
	if config.AuditLog {
		if err := promptGcpAuditLogQuestions(config, extraState); err != nil {
			return err
		}
	}

	// Validate one of configuration or audit log was enabled; otherwise error out
	if !config.Agentless && !config.Configuration && !config.AuditLog {
		return errors.New("must enable agentless, audit log or configuration")
	}

	// Ask for multiple projects & service account credentials if configuration or audit log is enabled
	if config.Configuration || config.AuditLog {
		if err := promptCustomizeGcpProjects(config); err != nil {
			return err
		}
		if err := promptGcpServiceAccountQuestions(config); err != nil {
			return err
		}
		// If no credentials path provided, ask for direct service account details
		if config.ServiceAccountCredentials == "" {
			// yes or no question for using existing service account
			useExisting := false
			if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
				{
					Prompt:   &survey.Confirm{Message: QuestionUseExistingServiceAccount, Default: false},
					Response: &useExisting,
				},
			}); err != nil {
				return err
			}
			if useExisting {
				if err := promptGcpExistingServiceAccountQuestions(config); err != nil {
					return err
				}
			}
		}
	}

	// Ask for output location
	if err := promptCustomizeGcpOutputLocation(extraState); err != nil {
		return err
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
