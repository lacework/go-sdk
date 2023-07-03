package cmd

import (
	"time"

	"github.com/lacework/go-sdk/lwgenerate/gcp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type GkeGenerateCommandExtraState struct {
	AskAdvanced                bool
	Output                     string
	ConfigureNewBucketSettings bool
	UseExistingServiceAccount  bool
	UseExistingSink            bool
	TerraformApply             bool
}

func (g *GkeGenerateCommandExtraState) isEmpty() bool {
	return g.Output == "" &&
		!g.AskAdvanced &&
		!g.UseExistingServiceAccount &&
		!g.UseExistingSink &&
		!g.TerraformApply
}

func (g *GkeGenerateCommandExtraState) writeCache() {
	if !g.isEmpty() {
		cli.WriteAssetToCache(CachedGkeAssetExtraState, time.Now().Add(time.Hour*1), g)
	}
}

var (
	QuestionGkeOrganizationIntegration = "Organization integration?"
	QuestionGkeOrganizationID          = "Specify the GCP organization ID:"
	QuestionGkeProjectID               = "Specify the project ID to be used to provision Lacework resources:"
	QuestionGkeServiceAccountCredsPath = "Specify service account credentials JSON path: (optional)"

	QuestionGkeConfigureAdvanced  = "Configure advanced integration options?"
	GkeAdvancedOpt                = "Configure additional options"
	QuestionGkeUseExistingSink    = "Use an existing sink?"
	QuestionGkeExistingSinkName   = "Specify the existing sink name"
	GkeAdvancedOptIntegrationName = "Customize integration name(s)"
	QuestionGkeIntegrationName    = "Specify a custom integration name: (optional)"

	GkeAdvancedOptExistingServiceAccount        = "Configure & use existing service account"
	QuestionGkeExistingServiceAccountName       = "Specify an existing service account name:"
	QuestionGkeExistingServiceAccountPrivateKey = "Specify an existing service account private key" +
		" (base64 encoded):" // guardrails-disable-line

	GkeAdvancedOptLocation             = "Customize output location"
	QuestionGkeCustomizeOutputLocation = "Provide the location for the output to be written:"
	QuestionGkeAnotherAdvancedOpt      = "Configure another advanced integration option"
	GkeAdvancedOptDone                 = "Done"

	GenerateGkeCommandState           = &gcp.GenerateGkeTfConfigurationArgs{}
	GenerateGkeExistingServiceAccount = &gcp.ServiceAccount{}
	GenerateGkeCommandExtraState      = &GkeGenerateCommandExtraState{}
	CachedGkeAssetIacParams           = "iac-gke-generate-params"
	CachedGkeAssetExtraState          = "iac-gke-extra-state"

	generateGkeTfCommand = &cobra.Command{
		Use:   "gke",
		Short: "Generate and/or execute Terraform code for GKE integration",
		Long: `Use this command to generate Terraform code for deploying Lacework into a GKE environment.

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
			cli.StartProgress("Generating Terraform Code...")

			if cli.Profile != "default" {
				GenerateGkeCommandState.LaceworkProfile = cli.Profile
			}

			mods := []gcp.Modifier{
				gcp.WithGkeExistingServiceAccount(GenerateGkeCommandState.ExistingServiceAccount),
				gcp.WithGkeExistingSinkName(GenerateGkeCommandState.ExistingSinkName),
				gcp.WithGkeIntegrationName(GenerateGkeCommandState.IntegrationName),
				gcp.WithGkeLabels(GenerateGkeCommandState.Labels),
				gcp.WithGkeLaceworkProfile(GenerateGkeCommandState.LaceworkProfile),
				gcp.WithGkeOrganizationId(GenerateGkeCommandState.OrganizationId),
				gcp.WithGkeOrganizationIntegration(GenerateGkeCommandState.OrganizationIntegration),
				gcp.WithGkePrefix(GenerateGkeCommandState.Prefix),
				gcp.WithGkeProjectId(GenerateGkeCommandState.ProjectId),
				gcp.WithGkePubSubSubscriptionLabels(GenerateGkeCommandState.PubSubSubscriptionLabels),
				gcp.WithGkePubSubTopicLabels(GenerateGkeCommandState.PubSubTopicLabels),
				gcp.WithGkeServiceAccountCredentials(GenerateGkeCommandState.ServiceAccountCredentials),
				gcp.WithGkeWaitTime(GenerateGkeCommandState.WaitTime),
			}

			hcl, err := gcp.NewGkeTerraform(mods...).Generate()
			cli.StopProgress()
			if err != nil {
				return errors.Wrap(err, "failed to generate terraform code")
			}

			dirname, _, err := writeGeneratedCodeToLocation(cmd, hcl, "gke")
			if err != nil {
				return err
			}

			err = SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
				Prompt:   &survey.Confirm{Default: GenerateGkeCommandExtraState.TerraformApply, Message: QuestionRunTfPlan},
				Response: &GenerateGkeCommandExtraState.TerraformApply,
			})

			if err != nil {
				return errors.Wrap(err, "failed to prompt for terraform execution")
			}

			locationDir, _ := determineOutputDirPath(dirname, "gke")
			if GenerateGkeCommandExtraState.TerraformApply {
				err := executionPreRunChecks(dirname, locationDir, "gke")
				if err != nil {
					return err
				}
			}

			if !GenerateGkeCommandExtraState.TerraformApply {
				cli.OutputHuman(provideGuidanceAfterExit(false, false, locationDir, "terraform"))
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if err = gkeValidation(cmd); err != nil {
				return err
			}

			if cli.InteractiveMode() {
				if err = gkeCaching(); err != nil {
					return err
				}
			}

			err = promptGkeGenerate(GenerateGkeCommandState, GenerateGkeExistingServiceAccount, GenerateGkeCommandExtraState)
			if err != nil {
				return errors.Wrap(err, "collecting/confirming parameters")
			}

			return nil
		},
	}
)

func gkeValidation(cmd *cobra.Command) error {
	dirname, err := cmd.Flags().GetString("output")
	if err != nil {
		return errors.Wrap(err, "failed to load command flags")
	}

	if err := validateOutputLocation(dirname); err != nil {
		return err
	}

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

	return nil
}

func gkeCaching() error {
	cachedOptions := &gcp.GenerateGkeTfConfigurationArgs{}
	iacParamsExpired := cli.ReadCachedAsset(CachedGkeAssetIacParams, &cachedOptions)
	if iacParamsExpired {
		cli.Log.Debug("loaded previously set values for GCP iac generation")
	}

	extraState := &GkeGenerateCommandExtraState{}
	extraStateParamsExpired := cli.ReadCachedAsset(CachedGkeAssetExtraState, &extraState)
	if extraStateParamsExpired {
		cli.Log.Debug("loaded previously set values for GCP iac generation (extra state)")
	}

	answer := false
	if !iacParamsExpired || !extraStateParamsExpired {
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Prompt:   &survey.Confirm{Message: QuestionUsePreviousCache, Default: false},
			Response: &answer,
		}); err != nil {
			return errors.Wrap(err, "failed to load saved options")
		}
	}

	if answer {
		if err := mergo.Merge(GenerateGkeCommandState, cachedOptions); err != nil {
			return errors.Wrap(err, "failed to load saved options")
		}
		if err := mergo.Merge(GenerateGkeCommandExtraState, extraState); err != nil {
			return errors.Wrap(err, "failed to load saved options")
		}
	}

	return nil
}

func initGenerateGkeTfCommandFlags() {
	generateGkeTfCommand.PersistentFlags().StringVar(
		&GenerateGkeExistingServiceAccount.Name,
		"existing_service_account_name",
		"",
		"specify existing service account name")
	generateGkeTfCommand.PersistentFlags().StringVar(
		&GenerateGkeExistingServiceAccount.PrivateKey,
		"existing_service_account_private_key",
		"",
		"specify existing service account private key (base64 encoded)")
	generateGkeTfCommand.PersistentFlags().StringVar(
		&GenerateGkeCommandState.ExistingSinkName,
		"existing_sink_name",
		"",
		"specify existing sink name")
	generateGkeTfCommand.PersistentFlags().StringVar(
		&GenerateGkeCommandState.IntegrationName,
		"integration_name",
		"",
		"specify a custom integration name")
	generateGkeTfCommand.PersistentFlags().StringVar(
		&GenerateGkeCommandState.OrganizationId,
		"organization_id",
		"",
		"specify the organization id (only set if organization_integration is set)")
	generateGkeTfCommand.PersistentFlags().BoolVar(
		&GenerateGkeCommandState.OrganizationIntegration,
		"organization_integration",
		false,
		"enable organization integration")
	generateGkeTfCommand.PersistentFlags().StringVar(
		&GenerateGkeCommandState.Prefix,
		"prefix",
		"",
		"prefix that will be used at the beginning of every generated resource")
	generateGkeTfCommand.PersistentFlags().StringVar(
		&GenerateGkeCommandState.ProjectId,
		"project_id",
		"",
		"specify the project id to be used to provision lacework resources (required)")
	generateGkeTfCommand.PersistentFlags().StringVar(
		&GenerateGkeCommandState.ServiceAccountCredentials,
		"service_account_credentials",
		"",
		"specify service account credentials JSON file path (leave blank to make use of google credential ENV vars)")
	generateGkeTfCommand.PersistentFlags().StringVar(
		&GenerateGkeCommandState.WaitTime,
		"wait_time",
		"",
		"amount of time to wait before the next resource is provisioned")
	generateGkeTfCommand.PersistentFlags().BoolVar(
		&GenerateGkeCommandExtraState.TerraformApply,
		"apply",
		false,
		"run terraform apply without executing plan or prompting",
	)
	generateGkeTfCommand.PersistentFlags().StringVar(
		&GenerateGkeCommandExtraState.Output,
		"output",
		"",
		"location to write generated content (default is ~/lacework/gcp)",
	)
}

func gkeConfigIsEmpty(g *gcp.GenerateGkeTfConfigurationArgs) bool {
	return g.ServiceAccountCredentials == "" &&
		g.ProjectId == "" &&
		g.OrganizationId == "" &&
		g.LaceworkProfile == ""
}

func writeGkeGenerationArgsCache(a *gcp.GenerateGkeTfConfigurationArgs) {
	if !gkeConfigIsEmpty(a) {
		if a.ExistingServiceAccount.IsPartial() {
			a.ExistingServiceAccount = nil
		}
		cli.WriteAssetToCache(CachedGkeAssetIacParams, time.Now().Add(time.Hour*1), a)
	}
}

func promptGkeGenerate(
	config *gcp.GenerateGkeTfConfigurationArgs,
	existingServiceAccount *gcp.ServiceAccount,
	extraState *GkeGenerateCommandExtraState,
) error {
	if cli.InteractiveMode() {
		defer writeGkeGenerationArgsCache(config)
		defer extraState.writeCache()
	}

	if existingServiceAccount.Name != "" ||
		existingServiceAccount.PrivateKey != "" {
		config.ExistingServiceAccount = existingServiceAccount
	}

	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Prompt:   &survey.Input{Message: QuestionGkeProjectID, Default: config.ProjectId},
				Required: true,
				Opts:     []survey.AskOpt{survey.WithValidator(validateGcpProjectId)},
				Response: &config.ProjectId,
			},
			{
				Prompt:   &survey.Confirm{Message: QuestionGkeOrganizationIntegration, Default: config.OrganizationIntegration},
				Response: &config.OrganizationIntegration,
			},
			{
				Prompt:   &survey.Input{Message: QuestionGkeOrganizationID, Default: config.OrganizationId},
				Checks:   []*bool{&config.OrganizationIntegration},
				Required: true,
				Response: &config.OrganizationId,
			},
			{
				Prompt:   &survey.Input{Message: QuestionGkeServiceAccountCredsPath, Default: config.ServiceAccountCredentials},
				Opts:     []survey.AskOpt{survey.WithValidator(gcp.ValidateServiceAccountCredentials)},
				Response: &config.ServiceAccountCredentials,
			},
		}); err != nil {
		return err
	}

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Message: QuestionGkeConfigureAdvanced, Default: extraState.AskAdvanced},
		Response: &extraState.AskAdvanced,
	}); err != nil {
		return err
	}

	if extraState.AskAdvanced {
		if err := promptGkeAdvancedOptions(config, extraState); err != nil {
			return err
		}
	}

	return nil
}

func promptGkeAdvancedOptions(
	config *gcp.GenerateGkeTfConfigurationArgs, extraState *GkeGenerateCommandExtraState,
) error {
	answer := ""
	options := []string{
		GkeAdvancedOpt,
		GkeAdvancedOptExistingServiceAccount,
		GkeAdvancedOptIntegrationName,
		GkeAdvancedOptLocation,
		GkeAdvancedOptDone,
	}

	for answer != GkeAdvancedOptDone {
		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Prompt: &survey.Select{
				Message: "Which options would you like to configure?",
				Options: options,
			},
			Response: &answer,
		}); err != nil {
			return err
		}

		switch answer {
		case GkeAdvancedOpt:
			if err := promptGkeQuestions(config, extraState); err != nil {
				return err
			}
		case GkeAdvancedOptExistingServiceAccount:
			if err := promptGkeExistingServiceAccountQuestions(config); err != nil {
				return err
			}
		case GkeAdvancedOptIntegrationName:
			if err := promptGkeIntegrationNameQuestions(config); err != nil {
				return err
			}
		case GkeAdvancedOptLocation:
			if err := promptCustomizeGkeOutputLocation(extraState); err != nil {
				return err
			}
		}

		innerAskAgain := true
		if answer == GkeAdvancedOptDone {
			innerAskAgain = false
		}

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Checks:   []*bool{&innerAskAgain},
			Prompt:   &survey.Confirm{Message: QuestionGkeAnotherAdvancedOpt, Default: false},
			Response: &innerAskAgain,
		}); err != nil {
			return err
		}

		if !innerAskAgain {
			answer = GkeAdvancedOptDone
		}
	}

	return nil
}

func promptGkeQuestions(config *gcp.GenerateGkeTfConfigurationArgs, extraState *GkeGenerateCommandExtraState) error {
	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionGkeUseExistingSink, Default: extraState.UseExistingSink},
			Required: true,
			Response: &extraState.UseExistingSink,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGkeExistingSinkName, Default: config.ExistingSinkName},
			Checks:   []*bool{&extraState.UseExistingSink},
			Required: true,
			Response: &config.ExistingSinkName,
		},
	})

	return err
}

func promptGkeExistingServiceAccountQuestions(config *gcp.GenerateGkeTfConfigurationArgs) error {
	if config.ExistingServiceAccount == nil {
		config.ExistingServiceAccount = &gcp.ServiceAccount{}
	}

	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt: &survey.Input{
				Message: QuestionGkeExistingServiceAccountName,
				Default: config.ExistingServiceAccount.Name,
			},
			Response: &config.ExistingServiceAccount.Name,
			Opts:     []survey.AskOpt{survey.WithValidator(survey.Required)},
		},
		{
			Prompt: &survey.Input{
				Message: QuestionGkeExistingServiceAccountPrivateKey,
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

func promptGkeIntegrationNameQuestions(config *gcp.GenerateGkeTfConfigurationArgs) error {
	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionGkeIntegrationName, Default: config.IntegrationName},
			Response: &config.IntegrationName,
		},
	})

	return err
}

func promptCustomizeGkeOutputLocation(extraState *GkeGenerateCommandExtraState) error {
	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionGkeCustomizeOutputLocation, Default: extraState.Output},
		Response: &extraState.Output,
		Opts:     []survey.AskOpt{survey.WithValidator(validPathExists)},
		Required: true,
	})

	return err
}
