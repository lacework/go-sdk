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
	Output         string
	TerraformApply bool
}

var (
	// questions
	QuestionOciTenantOcid = "Specify the OCID of the tenant to be integrated"
	QuestionOciUserEmail  = "Specify the email address to associate with the integration OCI user"

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
	cli.StartProgress("Generating Terraform Code...")

	if cli.Profile != "default" {
		GenerateOciCommandState.LaceworkProfile = cli.Profile
	}

	// generate tf code
	hcl, err := oci.NewTerraform(
		GenerateOciCommandState.Config,
		oci.WithLaceworkProfile(GenerateOciCommandState.LaceworkProfile),
		oci.WithTenantOcid(GenerateOciCommandState.TenantOcid),
		oci.WithUserEmail(GenerateOciCommandState.OciUserEmail),
	).Generate()

	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "failed to generate terraform code")
	}

	dirname, _, err := writeGeneratedCodeToLocation(cmd, hcl, "oci")
	if err != nil {
		return err
	}

	// Prompt to execute
	err = SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Confirm{
			Default: GenerateOciCommandExtraState.TerraformApply,
			Message: QuestionRunTfPlan,
		},
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
	err := promptOciGenerate(GenerateOciCommandState, GenerateOciCommandExtraState)
	if err != nil {
		return errors.Wrap(err, "collecting/confirming parameters")
	}

	return nil
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

func (a *OciGenerateCommandExtraState) isEmpty() bool {
	return a.Output == "" &&
		!a.TerraformApply
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

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionOciTenantOcid, Default: config.TenantOcid},
		Response: &config.TenantOcid,
		Opts:     []survey.AskOpt{survey.WithValidator(survey.Required), survey.WithValidator(validateOciTenantOcid)},
		// Checks:   []*bool{configEnabled(config)},
	}); err != nil {
		return err
	}

	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Input{Message: QuestionOciUserEmail, Default: config.OciUserEmail},
		Response: &config.OciUserEmail,
		Opts:     []survey.AskOpt{survey.WithValidator(survey.Required), survey.WithValidator(validateOciUserEmail)},
		// Checks:   []*bool{configEnabled(config)},
	}); err != nil {
		return err
	}

	return nil
}
