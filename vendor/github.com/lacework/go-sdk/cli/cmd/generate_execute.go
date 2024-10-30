package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Masterminds/semver"
	"github.com/abiosoft/colima/util/terminal"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/pkg/errors"
)

var (
	requiredTerraformVersion = ">= 0.15.1"
	installTerraformVersion  = "1.0.11"
)

type terraformVersion struct {
	Version string `json:"terraform_version"`
}

// helper function to create new *tfexec.Terraform object and return useful error if not found
func newTf(workingDir string, execPath string) (*tfexec.Terraform, error) {
	// Create new tf object
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not locate terraform binary")
	}

	return tf, nil
}

// LocateOrInstallTerraform Determine if terraform is installed, if that version is new enough,
// and if not install a new ephemeral binary of the correct version into tmp location
//
// forceInstall: if set always install ephemeral binary
func LocateOrInstallTerraform(forceInstall bool, workingDir string) (*tfexec.Terraform, error) {
	// find existing binary if not force install
	execPath, _ := exec.LookPath("terraform")
	if execPath != "" {
		cli.Log.Debugf("existing terraform path %s", execPath)
	}

	existingVersionOk := false
	if !forceInstall && execPath != "" {
		// Test if it's an OK version
		requiredVersion := requiredTerraformVersion
		constraint, _ := semver.NewConstraint(requiredVersion)

		// Extract tf version && check for unsupportedExistingVersion
		out, err := exec.Command("terraform", "--version", "--json").Output()
		if err != nil {
			return nil,
				errors.Wrap(
					err,
					fmt.Sprintf("failed to collect version from existing terraform install (%s)", execPath))
		}

		// If this version supports checking the version via --version --json, check if we can use it
		var data terraformVersion
		unsupportedVersionCheck := false
		err = json.Unmarshal(out, &data)
		if err != nil {
			// If this version does not support checking version via  --version --json, report and install new
			cli.OutputHuman(
				"Existing Terraform version cannot be used, version doesn't meet requirement %s, "+
					"installing short lived version\n",
				requiredVersion)
			unsupportedVersionCheck = true
		}
		cli.Log.Debugf("existing terraform version %s", data.Version)

		// Parse into new semver
		if !unsupportedVersionCheck {
			tfVersion, err := semver.NewVersion(data.Version)
			if err != nil {
				return nil,
					errors.Wrap(
						err,
						fmt.Sprintf("version from existing terraform install is invalid (%s)", data.Version))
			}

			// Test if it matches
			existingVersionOk, _ = constraint.Validate(tfVersion)
			if !existingVersionOk {
				cli.OutputHuman(
					"Existing Terraform version cannot be used, version %s doesn't meet requirement %s, "+
						"installing short lived version\n",
					data.Version,
					requiredVersion)
			}
			cli.Log.Debug("using existing terraform install")
		}
	}

	if !existingVersionOk {
		// If forceInstall was true or the existing version couldn't be used, install it
		installer := &releases.ExactVersion{
			Product: product.Terraform,
			Version: version.Must(version.NewVersion(installTerraformVersion)),
		}

		cli.StartProgress("Installing Terraform")
		installPath, err := installer.Install(context.Background())
		if err != nil {
			return nil, errors.Wrap(err, "error installing terraform")
		}
		cli.StopProgress()
		execPath = installPath
	}

	// Return *tfexec.Terraform object
	return newTf(workingDir, execPath)
}

// used to create suitable response messages when showing actions tf will take as a result of a plan execution
func createOrDestroy(create bool,
	destroy bool,
	update bool,
	read bool,
	noop bool,
	replace bool,
	createBeforeDestroy bool,
	destroyBeforeCreate bool,
) string {
	switch {
	case create:
		return "created"
	case destroy:
		return "destroyed"
	case update:
		return "update"
	case read:
		return "read"
	case noop:
		return "(noop)"
	case replace:
		return "replaced"
	case createBeforeDestroy:
		return "created before destroy"
	case destroyBeforeCreate:
		return "destroyed before create"
	default:
		return "unchanged"
	}
}

type TfPlanChangesSummary struct {
	plan    *tfjson.Plan
	create  int
	deleted int
	update  int
	replace int
}

// buildHumanReadablePlannedActions creates a summarized listing of expected changes from Terraform
func buildHumanReadablePlannedActions(workingDir string, execPath string, data []*tfjson.ResourceChange) string {
	outputString := strings.Builder{}
	outputString.WriteString("Resource details: \n")

	for _, c := range data {
		outputString.WriteString(fmt.Sprintf("  %s.%s will be %s\n", c.Type, c.Name,
			createOrDestroy(
				c.Change.Actions.Create(),
				c.Change.Actions.Delete(),
				c.Change.Actions.Update(),
				c.Change.Actions.Read(),
				c.Change.Actions.NoOp(),
				c.Change.Actions.Replace(),
				c.Change.Actions.CreateBeforeDestroy(),
				c.Change.Actions.DestroyBeforeCreate(),
			),
		),
		)
	}
	outputString.WriteString("\n")
	outputString.WriteString(fmt.Sprintf(
		"More details can be viewed by running:\n\n  cd %s\n  %s show tfplan.json\n",
		workingDir, execPath,
	))
	outputString.WriteString("\n")
	return outputString.String()
}

// DisplayTerraformPlanChanges used to display the results of a plan
//
// returns true if apply should run, false to exit
func DisplayTerraformPlanChanges(tf *tfexec.Terraform, data TfPlanChangesSummary) (bool, error) {
	// Prompt for next steps
	prompt := true
	previewShown := false
	var answer int

	// Displaying resources
	for prompt {
		id, err := promptForTerraformNextSteps(&previewShown, data)
		if err != nil {
			return false, errors.Wrap(err, "failed to run terraform")
		}

		switch {
		case id == 1 && !previewShown:
			cli.OutputHuman(buildHumanReadablePlannedActions(tf.WorkingDir(), tf.ExecPath(), data.plan.ResourceChanges))
		default:
			answer = id
			prompt = false
		}

		if id == 1 && !previewShown {
			previewShown = true
		}
	}

	// Run apply
	if answer == 0 {
		return true, nil
	}

	// Quit
	return false, nil
}

func parseTfPlanOutput(plan *tfjson.Plan) *TfPlanChangesSummary {
	// Build output of changes
	resourceCreate := 0
	resourceDelete := 0
	resourceUpdate := 0
	resourceReplace := 0

	for _, c := range plan.ResourceChanges {
		switch {
		case c.Change.Actions.Create():
			resourceCreate++
		case c.Change.Actions.Delete():
			resourceDelete++
		case c.Change.Actions.Update():
			resourceUpdate++
		case c.Change.Actions.Replace():
			resourceReplace++
		}
	}

	return &TfPlanChangesSummary{
		plan:    plan,
		create:  resourceCreate,
		deleted: resourceDelete,
		update:  resourceDelete,
		replace: resourceReplace,
	}
}

func processTfPlanChangesSummary(tf *tfexec.Terraform) (*TfPlanChangesSummary, error) {
	// Extract changes from tf plan
	cli.StartProgress("Getting terraform plan details")
	plan, err := tf.ShowPlanFile(context.Background(), "tfplan.json")
	if err != nil {
		return nil, errors.Wrap(err, "failed to inspect terraform plan")
	}
	cli.StopProgress()

	return parseTfPlanOutput(plan), nil
}

func TerraformInit(tf *tfexec.Terraform) error {
	cli.StartProgress("Running terraform init")
	err := tf.Init(context.Background(), tfexec.Upgrade(true))
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "failed to execute terraform init")
	}

	return nil
}

// TerraformExecPlan Run terraform plan using the workingDir from *tfexec.Terraform
//
// - Run plan
// - Get plan file details (returned)
func TerraformExecPlan(tf *tfexec.Terraform) (*TfPlanChangesSummary, error) {
	// Plan
	cli.StartProgress("Running terraform plan")
	_, err := tf.Plan(context.Background(), tfexec.Out("tfplan.json"))
	cli.StopProgress()
	if err != nil {
		return nil, err
	}

	// Gather changes from plan
	return processTfPlanChangesSummary(tf)
}

// TerraformExecApply Run terraform apply using the workingDir from *tfexec.Terraform
//
// - Run plan
// - Get plan file details (returned)
func TerraformExecApply(tf *tfexec.Terraform) error {
	cli.StartProgress("Running terraform apply")
	err := tf.Apply(context.Background())
	cli.StopProgress()
	if err != nil {
		return err
	}

	return nil
}

// Simple helper to prompt for next steps after TF plan
func promptForTerraformNextSteps(previewShown *bool, data TfPlanChangesSummary) (int, error) {
	options := []string{
		"Continue with Terraform Apply",
	}

	// Omit option to show details if we already have
	if !*previewShown {
		options = append(options, "Show details")
	}
	options = append(options, "Quit")

	var answer int
	err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt: &survey.Select{
			Message: fmt.Sprintf(
				"Terraform will create %d resources, delete %d resources, update %d resources, and replace %d resources.",
				data.create,
				data.deleted,
				data.update,
				data.replace),
			Options: options,
		},
		Response: &answer,
	})

	return answer, err
}

// this helper function is called when terraform flow has been completely executed through apply
func provideGuidanceAfterSuccess(workingDir string, laceworkProfile string) string {
	out := new(strings.Builder)
	fmt.Fprintf(out, "Lacework integration was successful! Terraform code saved in %s\n\n", workingDir)
	fmt.Fprintln(out, "To view integration status:")

	laceworkCmd := "    lacework cloud-account list\n\n"
	if laceworkProfile != "" {
		laceworkCmd = fmt.Sprintf("    lacework -p %s cloud-account list\n\n", laceworkProfile)
	}
	fmt.Fprint(out, laceworkCmd)

	return out.String()
}

func provideGuidanceAfterFailure(err error, workingDir string, binaryLocation string) string {
	out := new(strings.Builder)
	fmt.Fprintf(out, "\n\n%s\n\n", err.Error())
	fmt.Fprintln(out, strings.Repeat("-", 80))
	fmt.Fprint(out, "Terraform encountered an error (see above)\n\n")
	fmt.Fprintf(out, "The Terraform code, state, and plan output have been saved in %s.\n\n", workingDir)
	fmt.Fprintln(out, "Once the issues have been resolved, the integration can be continued using the following commands:")
	fmt.Fprintf(out, "  cd %s\n", workingDir)
	fmt.Fprintf(out, "  %s apply\n\n", binaryLocation)
	fmt.Fprintln(out, "Should you simply want to clean up the failed deployment, use the following commands:")
	fmt.Fprintf(out, "  cd %s\n", workingDir)
	fmt.Fprintf(out, "  %s destroy\n\n", binaryLocation)

	return out.String()
}

// this helper function is called when the entire generation/apply flow is not completed; it provides
// guidance on how to proceed from the last point of execution
func provideGuidanceAfterExit(initRun bool, planRun bool, workingDir string, binaryLocation string) string {
	planNote := " and plan output"
	if !planRun {
		planNote = ""
	}

	out := new(strings.Builder)
	fmt.Fprintf(out, "Terraform code%s saved in %s\n\n", planNote, workingDir)
	fmt.Fprintln(out, "The generated code can be executed at any point in the future using the following commands:")
	fmt.Fprintf(out, "  cd %s\n", workingDir)

	if !initRun {
		fmt.Fprintf(out, "  %s init\n", binaryLocation)
	}

	fmt.Fprintf(out, "  %s plan\n", binaryLocation)
	fmt.Fprintf(out, "  %s apply\n\n", binaryLocation)
	return out.String()
}

// Execute a terraform plan & execute
func TerraformPlanAndExecute(workingDir string) error {
	// Ensure Terraform is installed
	tf, err := LocateOrInstallTerraform(false, workingDir)
	if err != nil {
		return err
	}

	vw := terminal.NewVerboseWriter(10)
	tf.SetStdout(vw)
	tf.SetStderr(vw)

	// Initialize tf project
	if err := TerraformInit(tf); err != nil {
		return err
	}

	// Write plan
	changes, err := TerraformExecPlan(tf)
	if err != nil {
		return err
	}

	vw.Close()

	// Display changes and determine if apply should proceed
	proceed, err := DisplayTerraformPlanChanges(tf, *changes)
	if err != nil {
		return err
	}

	// If not proceed; display guidance on how to continue outside of this session
	if !proceed {
		cli.OutputHuman(provideGuidanceAfterExit(true, true, tf.WorkingDir(), tf.ExecPath()))
		return nil
	}

	vw = terminal.NewVerboseWriter(10)
	tf.SetStdout(vw)
	tf.SetStderr(vw)

	// Apply plan
	if err := TerraformExecApply(tf); err != nil {
		return errors.New(provideGuidanceAfterFailure(err, tf.WorkingDir(), tf.ExecPath()))
	}
	vw.Close()

	cli.OutputHuman(provideGuidanceAfterSuccess(tf.WorkingDir(), GenerateAwsCommandState.LaceworkProfile))
	return nil
}

func TerraformExecutePreRunCheck(outputLocation string, cloud string) (bool, error) {
	// If noninteractive, continue
	if !cli.InteractiveMode() {
		return true, nil
	}

	dirname, err := determineOutputDirPath(outputLocation, cloud)
	if err != nil {
		return false, err
	}
	stateFile := filepath.FromSlash(fmt.Sprintf("%s/terraform.tfstate", dirname))

	// If the file doesn't exist, carry on
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return true, nil
	}

	// If it does exist; confirm overwrite
	answer := false
	if err := SurveyQuestionInteractiveOnly(SurveyQuestionWithValidationArgs{
		Prompt:   &survey.Confirm{Message: fmt.Sprintf("Terraform state file %s already exists, continue?", stateFile)},
		Response: &answer,
	}); err != nil {
		return false, err
	}

	return answer, nil
}
