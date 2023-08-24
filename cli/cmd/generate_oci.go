package cmd

import (
	"time"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/lwgenerate/oci"
	"github.com/pkg/errors"
)

type OciGenerateCommandExtraState struct {
	AskAdvanced    bool
	Output         string
	TerraformApply bool
}

var (
	// questions
	QuestionOciEnableConfig            = "Enable configuration integration?"
	QuestionOciTenantOcid              = "Specify the OCID of the tenant to be integrated"
	QuestionOciUserEmail               = "Specify the email address to associate with the integration OCI user"
	QuestionOciConfigAdvanced          = "Configure advanced integration options?"
	QuestionCustomizeOciConfigName     = "Customize Config integration name?"
	QuestionOciConfigName              = "Specify name of config integration (optional)"
	QuestionOciCustomizeOutputLocation = "Provide the location for the output to be written:"
	QuestionOciAnotherAdvancedOpt      = "Configure another advanced integration option"

	// options
	OciAdvancedOptDone            = "Done"
	OciAdvancedOptLocation        = "Customize output location"
	OciAdvancedOptIntegrationName = "Customize integration name"

	// state
	GenerateOciCommandState      = &oci.GenerateOciTfConfigurationArgs{}
	GenerateOciCommandExtraState = &OciGenerateCommandExtraState{}

	// cache keys
	CachedOciAssetIacParams  = "iac-oci-generate-params"
	CachedAssetOciExtraState = "iac-oci-extra-state"

	// oci command is used to generate TF code for OCI
	generateOciTfCommand = &cobra.Command{
		Use:   "oci",
		Short: "Generate and/or execute Terraform code for OCI integration",
		Long: `Use this command to generate Terraform code for deploying Lacework into an OCI tenant.

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
		RunE:    runGenerateOci,
		PreRunE: preRunGenerateOci,
	}
)

func runGenerateOci(cmd *cobra.Command, args []string) error {
	// Generate TF Code
	cli.StartProgress("Generating Terraform Code...")

	// Explicitly set Lacework profile if it was passed in main args
	if cli.Profile != "default" {
		GenerateOciCommandState.LaceworkProfile = cli.Profile
	}

	// Setup modifiers for NewTerraform constructor
	mods := []oci.OciTerraformModifier{
		oci.WithLaceworkProfile(GenerateOciCommandState.LaceworkProfile),
		oci.WithConfigName(GenerateOciCommandState.ConfigName),
		oci.WithTenantOcid(GenerateOciCommandState.TenantOcid),
		oci.WithUserEmail(GenerateOciCommandState.OciUserEmail),
	}

	// Create new struct
	data := oci.NewTerraform(
		GenerateOciCommandState.Config,
		mods...)

	// Generate
	hcl, err := data.Generate()
	cli.StopProgress()

	if err != nil {
		return errors.Wrap(err, "failed to generate terraform code")
	}

	// Write-out generated code to location specified
	dirname, _, err := writeGeneratedCodeToLocation(cmd, hcl, "oci")
	if err != nil {
		return err
	}

	// Prompt to execute
	err = SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Default: GenerateOciCommandExtraState.TerraformApply, Message: QuestionRunTfPlan},
		Response: &GenerateOciCommandExtraState.TerraformApply,
	})

	if err != nil {
		return errors.Wrap(err, "failed to prompt for terraform execution")
	}

	locationDir, _ := determineOutputDirPath(dirname, "oci")
	if GenerateOciCommandExtraState.TerraformApply {
		// Execution pre-run check
		err := executionPreRunChecks(dirname, locationDir, "oci")
		if err != nil {
			return err
		}
	}

	// Output where code was generated
	if !GenerateOciCommandExtraState.TerraformApply {
		cli.OutputHuman(provideGuidanceAfterExit(false, false, locationDir, "terraform"))
	}

	return nil
}

func preRunGenerateOci(cmd *cobra.Command, _ []string) error {
	// Validate output location is OK if supplied
	dirname, err := cmd.Flags().GetString("output")
	if err != nil {
		return errors.Wrap(err, "failed to load command flags")
	}
	if err := validateOutputLocation(dirname); err != nil {
		return err
	}

	// Validate tenant OCID
	tenantOcid, err := cmd.Flags().GetString("tenant_ocid")
	if err != nil {
		return errors.Wrap(err, "failed to load command flags")
	}
	if err := validateOciTenantOcid(tenantOcid); tenantOcid != "" && err != nil {
		return err
	}

	// Validate user email
	ociUserEmail, err := cmd.Flags().GetString("oci_user_email")
	if err != nil {
		return errors.Wrap(err, "failed to load command flags")
	}
	if err := validateOciUserEmail(ociUserEmail); ociUserEmail != "" && err != nil {
		return err
	}

	// Load any cached inputs if interactive
	if cli.InteractiveMode() {
		cachedOptions := &oci.GenerateOciTfConfigurationArgs{}
		iacParamsExpired := cli.ReadCachedAsset(CachedOciAssetIacParams, &cachedOptions)
		if iacParamsExpired {
			cli.Log.Debug("loaded previously set values for OCI iac generation")
		}

		extraState := &OciGenerateCommandExtraState{}
		extraStateParamsExpired := cli.ReadCachedAsset(CachedAssetOciExtraState, &extraState)
		if extraStateParamsExpired {
			cli.Log.Debug("loaded previously set values for OCI iac generation (extra state)")
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
			if err := mergo.Merge(GenerateOciCommandState, cachedOptions); err != nil {
				return errors.Wrap(err, "failed to load saved options")
			}
			if err := mergo.Merge(GenerateOciCommandExtraState, extraState); err != nil {
				return errors.Wrap(err, "failed to load saved options")
			}
		}
	}

	// Collect and/or confirm parameters
	err = promptOciGenerate(GenerateOciCommandState, GenerateOciCommandExtraState)
	if err != nil {
		return errors.Wrap(err, "collecting/confirming parameters")
	}

	return nil
}

func initGenerateOciTfCommandFlags() {
	// add flags to sub commands
	generateOciTfCommand.PersistentFlags().BoolVar(
		&GenerateOciCommandState.Config,
		"config",
		false,
		"enable config integration")
	generateOciTfCommand.PersistentFlags().StringVar(
		&GenerateOciCommandState.ConfigName,
		"config_name",
		"",
		"specify name of config integration")
	generateOciTfCommand.PersistentFlags().StringVar(
		&GenerateOciCommandState.TenantOcid,
		"tenant_ocid",
		"",
		"specify the OCID of the tenant to integrate")
	generateOciTfCommand.PersistentFlags().StringVar(
		&GenerateOciCommandState.OciUserEmail,
		"oci_user_email",
		"",
		"specify the email of the OCI user created for integration")
	generateOciTfCommand.PersistentFlags().BoolVar(
		&GenerateOciCommandExtraState.TerraformApply,
		"apply",
		false,
		"run terraform apply without executing plan or prompting",
	)
	generateOciTfCommand.PersistentFlags().StringVar(
		&GenerateOciCommandExtraState.Output,
		"output",
		"",
		"location to write generated content (default is ~/lacework/oci)",
	)
}

// basic validation of Tenant OCID format
func validateOciTenantOcid(val interface{}) error {
	return validateStringWithRegex(
		val,
		// https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm
		`ocid1\.tenancy\.[^\.\s]*\.[^\.\s]*(\.[^\.\s]+)?\.[^\.\s]+`,
		"invalid tenant OCID supplied",
	)
}

// basic validation of email
func validateOciUserEmail(val interface{}) error {
	return validateEmailAddress(val)
}

func promptCustomizeOciOutputLocation(extraState *OciGenerateCommandExtraState) error {
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionOciCustomizeOutputLocation, Default: extraState.Output},
		Response: &extraState.Output,
		Opts:     []survey.AskOpt{survey.WithValidator(validPathExists)},
		Required: true,
	}); err != nil {
		return err
	}

	return nil
}

func promptCustomizeOciConfigOptions(config *oci.GenerateOciTfConfigurationArgs) error {
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionOciConfigName, Default: config.ConfigName},
		Checks:   []*bool{&config.Config},
		Response: &config.ConfigName,
	}); err != nil {
		return err
	}

	return nil
}

func askAdvancedOciOptions(config *oci.GenerateOciTfConfigurationArgs, extraState *OciGenerateCommandExtraState) error {
	answer := ""

	// Prompt for options
	for answer != OciAdvancedOptDone {
		var options []string

		// Determine if user specified name for Config is potentially required
		if config.Config {
			options = append(options, OciAdvancedOptIntegrationName)
		}

		options = append(options, OciAdvancedOptLocation)

		options = append(options, OciAdvancedOptDone)
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
		case OciAdvancedOptLocation:
			if err := promptCustomizeOciOutputLocation(extraState); err != nil {
				return err
			}
		case OciAdvancedOptIntegrationName:
			if err := promptCustomizeOciConfigOptions(config); err != nil {
				return err
			}
		}

		// Re-prompt if not done
		innerAskAgain := true
		if answer == OciAdvancedOptDone {
			innerAskAgain = false
		}

		if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
			Checks:   []*bool{&innerAskAgain},
			Prompt:   &survey.Confirm{Message: QuestionOciAnotherAdvancedOpt, Default: false},
			Response: &innerAskAgain,
		}); err != nil {
			return err
		}

		if !innerAskAgain {
			answer = OciAdvancedOptDone
		}
	}

	return nil
}

func configEnabled(config *oci.GenerateOciTfConfigurationArgs) *bool {
	return &config.Config
}

func (a *OciGenerateCommandExtraState) isEmpty() bool {
	return a.Output == "" &&
		!a.TerraformApply &&
		!a.AskAdvanced
}

// Flush current state of the struct to disk, provided it's not empty
func (a *OciGenerateCommandExtraState) writeCache() {
	if !a.isEmpty() {
		cli.WriteAssetToCache(CachedAssetOciExtraState, time.Now().Add(time.Hour*1), a)
	}
}

func ociConfigIsEmpty(g *oci.GenerateOciTfConfigurationArgs) bool {
	return !g.Config &&
		g.ConfigName == "" &&
		g.LaceworkProfile == "" &&
		g.TenantOcid == "" &&
		g.OciUserEmail == ""
}

func writeOciGenerationArgsCache(a *oci.GenerateOciTfConfigurationArgs) {
	if !ociConfigIsEmpty(a) {
		cli.WriteAssetToCache(CachedOciAssetIacParams, time.Now().Add(time.Hour*1), a)
	}
}

// entry point for launching a survey to build out the required generation parameters
func promptOciGenerate(
	config *oci.GenerateOciTfConfigurationArgs,
	extraState *OciGenerateCommandExtraState,
) error {
	// Cache for later use if generation is abandon and in interactive mode
	if cli.InteractiveMode() {
		defer writeOciGenerationArgsCache(config)
		defer extraState.writeCache()
	}

	// These are the core questions that should be asked.
	if err := SurveyMultipleQuestionWithValidation(
		[]SurveyQuestionWithValidationArgs{
			{
				Prompt:   &survey.Confirm{Message: QuestionOciEnableConfig, Default: config.Config},
				Response: &config.Config,
			},
		}); err != nil {
		return err
	}

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionOciTenantOcid, Default: config.TenantOcid},
		Response: &config.TenantOcid,
		Opts:     []survey.AskOpt{survey.WithValidator(survey.Required), survey.WithValidator(validateOciTenantOcid)},
		Checks:   []*bool{configEnabled(config)},
	}); err != nil {
		return err
	}

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionOciUserEmail, Default: config.OciUserEmail},
		Response: &config.OciUserEmail,
		Opts:     []survey.AskOpt{survey.WithValidator(survey.Required), survey.WithValidator(validateOciUserEmail)},
		Checks:   []*bool{configEnabled(config)},
	}); err != nil {
		return err
	}

	// Validate that config was enabled. Otherwise throw error.
	if !config.Config {
		return errors.New("must enable config to continue")
	}

	// Find out if the customer wants to specify more advanced features
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Message: QuestionOciConfigAdvanced, Default: extraState.AskAdvanced},
		Response: &extraState.AskAdvanced,
	}); err != nil {
		return err
	}

	// Keep prompting for advanced options until the say done
	if extraState.AskAdvanced {
		if err := askAdvancedOciOptions(config, extraState); err != nil {
			return err
		}
	}

	return nil
}
