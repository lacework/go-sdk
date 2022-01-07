package cmd

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/internal/file"
	"github.com/lacework/go-sdk/lwgenerate/gcp"
)

var (
	// Define question text here to be reused in testing
	QuestionGcpEnableConfig            = "Enable Config Integration?"
	QuestionGcpEnableAuditLog          = "Enable AuditLog Integration?"
	QuestionGcpOrganizationIntegration = "Organization Integration?"
	QuestionGcpOrganizationID          = "Specify the GCP Organization ID:"
	QuestionGcpProjectID               = "Specify the Project ID to be used to provision Lacework resources: (optional)"
	QuestionGcpServiceAccountCredsPath = "Specify Service Account credentials JSON path: (optional)"

	QuestionGcpConfigureAdvanced             = "Configure advanced integration options?"
	GcpAdvancedOptExistingServiceAccount     = "Configure & use existing Service Account"
	QuestionExistingServiceAccountName       = "Specify an existing Service Account name:"
	QuestionExistingServiceAccountPrivateKey = "Specify an existing Service Account Private key (base64 encoded):"

	GcpAdvancedOptAuditLog              = "Configure additional AuditLog options"
	QuestionGcpUseExistingBucket        = "Use an existing Bucket?"
	QuestionGcpExistingBucketName       = "Specify an existing Bucket name:"
	QuestionGcpConfigureNewBucket       = "Configure settings for new Bucket?"
	QuestionGcpBucketName               = "Specify new Bucket name: (optional)"
	QuestionGcpBucketRegion             = "Specify the Bucket Region: (optional)"
	QuestionGcpBucketLocation           = "Specify the Bucket Location:  (optional)"
	QuestionGcpBucketRetention          = "Specify the Bucket Retention Days:  (optional)"
	QuestionGcpBucketLifecycle          = "Specify the Bucket Lifecycle Rule Age:  (optional)"
	QuestionGcpEnableUBLA               = "Enable Uniform Bucket Level Access(UBLA)?"
	QuestionGcpEnableBucketForceDestroy = "Enable Bucket Force Destroy?"
	QuestionGcpUseExistingSink          = "Use an existing Sink?"
	QuestionGcpExistingSinkName         = "Specify the existing Sink name"

	GcpAdvancedOptIntegrationName      = "Customize Integration name(s)"
	QuestionGcpConfigIntegrationName   = "Specify the custom Config integration name:"
	QuestionGcpAuditLogIntegrationName = "Specify the custom AuditLog integration name:"

	QuestionGcpAnotherAdvancedOpt      = "Configure another advanced integration option"
	GcpAdvancedOptLocation             = "Customize output location"
	QuestionGcpCustomizeOutputLocation = "Provide the location for the output to be written:"
	GcpAdvancedOptDone                 = "Done"

	// GcpRegionRegex regex used for validating region input
	GcpRegionRegex = `(asia|australia|europe|northamerica|southamerica|us)-(central|(north|south)?(east|west)?)\d`

	GenerateGcpCommandState                  = &gcp.GenerateGcpTfConfigurationArgs{}
	GenerateGcpExistingServiceAccountDetails = &gcp.ExistingServiceAccountDetails{}
	GenerateGcpCommandExtraState             = &GcpGenerateCommandExtraState{}
	CachedGcpAssetIacParams                  = "iac-gcp-generate-params"
	CachedAssetGcpExtraState                 = "iac-gcp-extra-state"

	// gcp command is used to generate TF code for gcp
	generateGcpTfCommand = &cobra.Command{
		Use:   "gcp",
		Short: "Generate and/or execute Terraform code for GCP integration",
		Long: `Use this command to generate Terraform code for deploying Lacework into an GCP environment.

By default, this command will function interactively, prompting for the required information to setup the new cloud account. In interactive mode, this command will:

* Prompt for the required information to setup the integration
* Generate new Terraform code using the inputs
* Optionally, run the generated Terraform code:
  * If Terraform is already installed, the version will be confirmed suitable for use
	* If Terraform is not installed, or the version installed is not suitable, a new version will be installed into a temporary location
	* Once Terraform is detected or installed, Terraform plan will be executed
	* The command will prompt with the outcome of the plan and allow to view more details or continue with Terraform apply
	* If confirmed, Terraform apply will be run, completing the setup of the cloud account

This command can also be run in noninteractive mode. See help output for more details on supplying required values for generation.
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
				gcp.WithConfigIntegrationName(GenerateGcpCommandState.ConfigIntegrationName),
				gcp.WithAuditLogLabels(GenerateGcpCommandState.AuditLogLabels),
				gcp.WithBucketLabels(GenerateGcpCommandState.BucketLabels),
				gcp.WithPubSubSubscriptionLabels(GenerateGcpCommandState.PubSubSubscriptionLabels),
				gcp.WithPubSubTopicLabels(GenerateGcpCommandState.PubSubTopicLabels),
				gcp.WithBucketRegion(GenerateGcpCommandState.BucketRegion),
				gcp.WithBucketLocation(GenerateGcpCommandState.BucketLocation),
				gcp.WithBucketName(GenerateGcpCommandState.BucketName),
				gcp.WithExistingLogBucketName(GenerateGcpCommandState.ExistingLogBucketName),
				gcp.WithExistingLogSinkName(GenerateGcpCommandState.ExistingLogSinkName),
				gcp.WithAuditLogIntegrationName(GenerateGcpCommandState.AuditLogIntegrationName),
				gcp.WithLaceworkProfile(GenerateGcpCommandState.LaceworkProfile),
			}

			if GenerateGcpCommandState.OrganizationIntegration {
				mods = append(mods, gcp.WithOrganizationIntegration(GenerateGcpCommandState.OrganizationIntegration))
			}

			if GenerateGcpCommandState.EnableForceDestroyBucket {
				mods = append(mods, gcp.WithEnableForceDestroyBucket())
			}

			if GenerateGcpCommandState.EnableUBLA {
				mods = append(mods, gcp.WithEnableUBLA())
			}

			if GenerateGcpCommandState.LogBucketLifecycleRuleAge != nil {
				mods = append(mods, gcp.WithLogBucketLifecycleRuleAge(*GenerateGcpCommandState.LogBucketLifecycleRuleAge))
			}

			if GenerateGcpCommandState.LogBucketRetentionDays != 0 {
				mods = append(mods, gcp.WithLogBucketRetentionDays(GenerateGcpCommandState.LogBucketRetentionDays))
			}

			// Create new struct
			data := gcp.NewTerraform(
				GenerateGcpCommandState.Config,
				GenerateGcpCommandState.AuditLog,
				mods...)

			// Generate
			hcl, err := data.Generate()
			cli.StopProgress()

			if err != nil {
				return errors.Wrap(err, "failed to generate terraform code")
			}

			// Write-out generated code to location specified
			dirname, location, err := writeGeneratedCodeToLocation(cmd, hcl)
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

			// Execute
			locationDir := filepath.Dir(location)
			if GenerateGcpCommandExtraState.TerraformApply {
				// Execution pre-run check
				err := executionPreRunChecks(dirname, locationDir)
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
			if err := validateGcpServiceAccountCredentials(gcpSaCredentials); gcpSaCredentials != "" && err != nil {
				return err
			}

			// Validate gcp region, if passed
			region, err := cmd.Flags().GetString("bucket_region")
			if err != nil {
				return errors.Wrap(err, "failed to load command flags")
			}
			if err := validateGcpRegion(region); region != "" && err != nil {
				return err
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
			err = promptGcpGenerate(GenerateGcpCommandState, GenerateGcpExistingServiceAccountDetails, GenerateGcpCommandExtraState)
			if err != nil {
				return errors.Wrap(err, "collecting/confirming parameters")
			}

			return nil
		},
	}
)

func validateServiceAccountCredentialsFile(credFile string) error {
	if file.FileExists(credFile) {
		jsonFile, err := os.Open(credFile)
		// if we os.Open returns an error then handle it
		if err != nil {
			return errors.Wrap(err, "Issue opening GCP credentials file")
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var result map[string]interface{}
		err = json.Unmarshal(byteValue, &result)
		if err != nil {
			return errors.Wrap(err, "Unable to parse credentials file.")
		}
		valid := validateSaCredFileContent(result)
		if !valid {
			return errors.New("Invalid GCP Service Account credentials file. " +
				"The private_key and client_email fields MUST be present. private_key must be base64 encoded")
		}

	}
	return nil
}

func validateSaCredFileContent(credFileContent map[string]interface{}) bool {
	if credFileContent["private_key"] != nil && credFileContent["client_email"] != nil {
		err := validateStringIsBase64(credFileContent["private_key"].(string))
		if err == nil {
			return true
		}
	}
	return false
}

// create survey.Validator for string is base64
func validateStringIsBase64(val interface{}) error {
	switch value := val.(type) {
	case string:
		// if value isn't base64, return error
		_, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return errors.Wrap(err, "provided private key is not base64 encoded")
		}
	default:
		// if the value passed is not a string
		return errors.New("value must be a string")
	}

	return nil
}

// survey.Validator for gcp region
func validateGcpRegion(val interface{}) error {
	return validateStringWithRegex(val, GcpRegionRegex, "invalid region name supplied")
}

// survey.Validator for gcp profile
func validateGcpServiceAccountCredentials(credentials string) error {
	return validateServiceAccountCredentialsFile(credentials)
}

func promptGcpAuditLogQuestions(config *gcp.GenerateGcpTfConfigurationArgs, extraState *GcpGenerateCommandExtraState) error {
	// Only ask these questions if configure audit log is true
	err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionGcpUseExistingBucket, Default: extraState.UseExistingBucket},
			Checks:   []*bool{&config.AuditLog},
			Response: &extraState.UseExistingBucket,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpExistingBucketName, Default: config.ExistingLogBucketName},
			Checks:   []*bool{&config.AuditLog, &extraState.UseExistingBucket},
			Required: true,
			Response: &config.ExistingLogBucketName,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpConfigureNewBucket, Default: config.ExistingLogBucketName},
			Checks:   []*bool{&config.AuditLog, &extraState.UseExistingBucket},
			Required: true,
			Response: &config.ExistingLogBucketName,
		},
	}, config.AuditLog)

	newBucket := !extraState.UseExistingBucket
	err = SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Confirm{Message: QuestionGcpConfigureNewBucket, Default: extraState.ConfigureNewBucketSettings},
			Checks:   []*bool{&config.AuditLog, &newBucket},
			Required: true,
			Response: &extraState.ConfigureNewBucketSettings,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpBucketName, Default: config.BucketName},
			Checks:   []*bool{&config.AuditLog, &newBucket, &extraState.ConfigureNewBucketSettings},
			Required: true,
			Response: &config.BucketName,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpBucketRegion, Default: config.BucketName},
			Checks:   []*bool{&config.AuditLog, &newBucket, &extraState.ConfigureNewBucketSettings},
			Required: true,
			Response: &config.BucketRegion,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpBucketLocation, Default: config.BucketName},
			Checks:   []*bool{&config.AuditLog, &newBucket, &extraState.ConfigureNewBucketSettings},
			Required: true,
			Response: &config.BucketLocation,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpBucketRetention, Default: string(rune(config.LogBucketRetentionDays))},
			Checks:   []*bool{&config.AuditLog, &newBucket, &extraState.ConfigureNewBucketSettings},
			Required: true,
			Response: &config.LogBucketRetentionDays,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpBucketLifecycle, Default: string(rune(*config.LogBucketLifecycleRuleAge))},
			Checks:   []*bool{&config.AuditLog, &newBucket, &extraState.ConfigureNewBucketSettings},
			Required: true,
			Response: &config.LogBucketLifecycleRuleAge,
		},
		{
			Prompt:   &survey.Confirm{Message: QuestionGcpEnableUBLA, Default: config.EnableUBLA},
			Checks:   []*bool{&config.AuditLog, &newBucket, &extraState.ConfigureNewBucketSettings},
			Required: true,
			Response: &config.EnableUBLA,
		},
		{
			Prompt:   &survey.Confirm{Message: QuestionGcpEnableBucketForceDestroy, Default: config.EnableForceDestroyBucket},
			Checks:   []*bool{&config.AuditLog, &newBucket, &extraState.UseExistingBucket},
			Required: true,
			Response: &config.EnableForceDestroyBucket,
		},
		{
			Prompt:   &survey.Confirm{Message: QuestionGcpUseExistingSink, Default: extraState.UseExistingSink},
			Checks:   []*bool{&config.AuditLog, &newBucket, &extraState.ConfigureNewBucketSettings},
			Required: true,
			Response: &extraState.UseExistingSink,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpExistingSinkName, Default: config.ExistingLogSinkName},
			Checks:   []*bool{&config.AuditLog, &newBucket, &extraState.ConfigureNewBucketSettings, &extraState.UseExistingSink},
			Required: true,
			Response: &config.ExistingLogSinkName,
		},
	}, config.AuditLog)

	if err != nil {
		return err
	}

	return nil
}

func promptGcpExistingServiceAccountQuestions(config *gcp.GenerateGcpTfConfigurationArgs) error {
	// ensure struct is initialized
	if config.ExistingServiceAccount == nil {
		config.ExistingServiceAccount = &gcp.ExistingServiceAccountDetails{}
	}

	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionExistingServiceAccountName, Default: config.ExistingServiceAccount.Name},
			Response: &config.ExistingServiceAccount.Name,
			Opts:     []survey.AskOpt{survey.WithValidator(survey.Required)},
		},
		{
			Prompt:   &survey.Input{Message: QuestionExistingServiceAccountPrivateKey, Default: config.ExistingServiceAccount.PrivateKey},
			Response: &config.ExistingServiceAccount.PrivateKey,
			Opts:     []survey.AskOpt{survey.WithValidator(survey.Required), survey.WithValidator(validateStringIsBase64)},
		}}); err != nil {
		return err
	}

	return nil
}

func promptGcpIntegrationNameQuestions(config *gcp.GenerateGcpTfConfigurationArgs) error {
	if err := SurveyMultipleQuestionWithValidation([]SurveyQuestionWithValidationArgs{
		{
			Prompt:   &survey.Input{Message: QuestionGcpConfigIntegrationName, Default: config.ConfigIntegrationName},
			Checks:   []*bool{&config.Config},
			Required: true,
			Response: &config.ConfigIntegrationName,
		},
		{
			Prompt:   &survey.Input{Message: QuestionGcpAuditLogIntegrationName, Default: config.AuditLogIntegrationName},
			Checks:   []*bool{&config.AuditLog},
			Required: true,
			Response: &config.AuditLogIntegrationName,
		}}); err != nil {
		return err
	}

	return nil
}

func promptCustomizeGcpOutputLocation(extraState *GcpGenerateCommandExtraState) error {
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionGcpCustomizeOutputLocation, Default: extraState.Output},
		Response: &extraState.Output,
		Opts:     []survey.AskOpt{survey.WithValidator(validPathExists)},
		Required: true,
	}); err != nil {
		return err
	}

	return nil
}

func askAdvancedOptions(config *gcp.GenerateGcpTfConfigurationArgs, extraState *GcpGenerateCommandExtraState) error {
	answer := ""

	// Prompt for options
	for answer != GcpAdvancedOptDone {
		// Construction of this slice is a bit strange at first look, but the reason for that is because we have to do string
		// validation to know which option was selected due to how survey works; and doing it by index (also supported) is
		// difficult when the options are dynamic (which they are)
		options := []string{GcpAdvancedOptAuditLog, GcpAdvancedOptExistingServiceAccount, GcpAdvancedOptIntegrationName}

		options = append(options, GcpAdvancedOptLocation, GcpAdvancedOptDone)
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

func configOrAuditLogEnabled(config *gcp.GenerateGcpTfConfigurationArgs) *bool {
	auditLogOrConfigEnabled := config.AuditLog || config.Config
	return &auditLogOrConfigEnabled
}

func gcpConfigIsEmpty(g *gcp.GenerateGcpTfConfigurationArgs) bool {
	return !g.AuditLog &&
		!g.Config &&
		g.ServiceAccountCredentials == "" &&
		g.GcpOrganizationId == "" &&
		g.LaceworkProfile == ""
}

func writeGcpGenerationArgsCache(a *gcp.GenerateGcpTfConfigurationArgs) {
	if !gcpConfigIsEmpty(a) {
		// If ExistingIamRole is partially set, don't write this to cache; the values won't work when loaded
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
				Prompt:   &survey.Confirm{Message: QuestionGcpEnableConfig, Default: config.Config},
				Response: &config.Config,
			},
			{
				Prompt:   &survey.Confirm{Message: QuestionGcpEnableAuditLog, Default: config.AuditLog},
				Response: &config.AuditLog,
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
				Prompt:   &survey.Input{Message: QuestionGcpProjectID, Default: config.GcpProjectId},
				Response: &config.GcpProjectId,
			},
			{
				Prompt:   &survey.Input{Message: QuestionGcpServiceAccountCredsPath, Default: config.ServiceAccountCredentials},
				Response: &config.ServiceAccountCredentials,
			},
		}); err != nil {
		return err
	}

	// Validate one of config or audit log was enabled; otherwise error out
	if !config.Config && !config.AuditLog {
		return errors.New("must enable audit log or config")
	}

	// Find out if the customer wants to specify more advanced features
	askAdvanced := false
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Message: QuestionGcpConfigureAdvanced, Default: askAdvanced},
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
