//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/lwcomponent"
)

const (
	componentTypeAnnotation string = "component"
	componentsCacheKey      string = "components"
)

var (
	// componentsCmd represents the components command
	componentsCmd = &cobra.Command{
		Use:     "component",
		Hidden:  true,
		Aliases: []string{"components"},
		Short:   "Manage components",
		Long:    `Manage components to extend your experience with the Lacework platform`,
	}

	// componentsListCmd represents the list sub-command inside the components command
	componentsListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all components",
		Long:    `List all available components and their current state`,
		RunE:    runComponentsList,
	}

	// componentsShowCmd represents the show sub-command inside the components command
	componentsShowCmd = &cobra.Command{
		Use:   "show <component>",
		Short: "Show details about a component",
		Long:  "Show details about a component",
		Args:  cobra.ExactArgs(1),
		RunE:  runComponentsShow,
	}

	// componentsInstallCmd represents the install sub-command inside the components command
	componentsInstallCmd = &cobra.Command{
		Use:   "install <component>",
		Short: "Install a new component",
		Long:  `Install a new component`,
		Args:  cobra.ExactArgs(1),
		RunE:  runComponentsInstall,
	}

	// componentsUpdateCmd represents the update sub-command inside the components command
	componentsUpdateCmd = &cobra.Command{
		Use:   "update <component>",
		Short: "Update an existing component",
		Long:  `Update an existing component`,
		Args:  cobra.ExactArgs(1),
		RunE:  runComponentsUpdate,
	}

	// componentsUninstallCmd represents the uninstall sub-command inside the components command
	componentsUninstallCmd = &cobra.Command{
		Use:     "uninstall <component>",
		Aliases: []string{"delete", "remove", "rm"},
		Short:   "Uninstall an existing component",
		Long:    `Uninstall an existing component`,
		Args:    cobra.ExactArgs(1),
		RunE:    runComponentsDelete,
	}

	// componentsDevModeCmd represents the dev sub-command inside the components command
	componentsDevModeCmd = &cobra.Command{
		Use:    "dev <component>",
		Hidden: true,
		Short:  "Enter development mode of a new or existing component",
		Args:   cobra.ExactArgs(1),
		RunE:   runComponentsDevMode,
	}

	versionArg string
)

func init() {
	// add the components command
	rootCmd.AddCommand(componentsCmd)

	componentsInstallCmd.PersistentFlags().StringVar(&versionArg, "version", "",
		"require a specific version to be installed (default is latest)")
	componentsUpdateCmd.PersistentFlags().StringVar(&versionArg, "version", "",
		"update to a specific version (default is latest)")

	// add sub-commands to the components command
	componentsCmd.AddCommand(componentsListCmd)
	componentsCmd.AddCommand(componentsShowCmd)
	componentsCmd.AddCommand(componentsInstallCmd)
	componentsCmd.AddCommand(componentsUpdateCmd)
	componentsCmd.AddCommand(componentsUninstallCmd)
	componentsCmd.AddCommand(componentsDevModeCmd)

	// load components dynamically
	cli.PrototypeLoadComponents()

	// v1 components
	cli.LoadComponents()
}

// hasInstalledCommands is used inside the cobra template for generating the usage
// of commands, it returns true if there are installed commands via the CDK
func hasInstalledCommands() bool {
	return cli.installedCmd
}

// isComponent is used inside the cobra template for generating the usage of
// commands, it needs the annotations of the command and it will return true
// if the command was installed from the CDK
func isComponent(annotations map[string]string) bool {
	t, found := annotations["type"]
	if found && t == componentTypeAnnotation {
		return true
	}
	return false
}

// IsComponentInstalled returns true if component is
// valid and installed
func (c *cliState) IsComponentInstalled(name string) bool {
	var err error
	c.LwComponents, err = lwcomponent.LocalState()
	if err != nil || c.LwComponents == nil {
		return false
	}

	component, found := c.LwComponents.GetComponent(name)
	if found && component.IsInstalled() {
		return true
	}
	return false
}

// Load v1 components
func (c *cliState) LoadComponents() {
	components, err := lwcomponent.LoadLocalComponents()
	if err != nil {
		c.Log.Debugw("unable to load components", "error", err)
		return
	}

	for _, component := range components {
		exists := false

		for _, cmd := range rootCmd.Commands() {
			if cmd.Use == component.Name {
				exists = true
				break
			}
		}

		// Skip components that were added by the prototype code
		if exists {
			continue
		}

		version := component.InstalledVersion()

		if version != nil {
			componentCmd := &cobra.Command{
				Use:                   component.Name,
				Short:                 component.Description,
				Annotations:           map[string]string{"type": componentTypeAnnotation},
				Version:               version.String(),
				SilenceUsage:          true,
				DisableFlagParsing:    true,
				DisableFlagsInUseLine: true,
				RunE: func(cmd *cobra.Command, args []string) error {
					return v1ComponentCommand(c, cmd)
				},
			}

			rootCmd.AddCommand(componentCmd)
		}
	}
}

// Grpc server used for components to communicate back to the CLI
func startGrpcServer(c *cliState) {
	if err := c.Serve(); err != nil {
		c.Log.Errorw("couldn't serve gRPC server", "error", err)
	}
}

func v1ComponentCommand(c *cliState, cmd *cobra.Command) error {
	// Parse component -v/--version flag
	versionVal, _ := cmd.Flags().GetBool("version")
	if versionVal {
		cmd.Printf("%s version %s\n", cmd.Use, cmd.Version)
		return nil
	}

	go startGrpcServer(c)

	catalog, err := LoadCatalog(cmd.Use, false)
	if err != nil {
		return errors.Wrap(err, "unable to load component Catalog")
	}

	component, err := catalog.GetComponent(cmd.Use)
	if err != nil {
		return err
	}

	if !component.Exec.Executable() {
		return errors.New("component is not executable")
	}

	c.Log.Debugw("running component", "component", cmd.Use,
		"args", c.componentParser.componentArgs,
		"cli_flags", c.componentParser.cliArgs)

	envs := []string{
		fmt.Sprintf("LW_COMPONENT_NAME=%s", cmd.Use),
	}

	envs = append(envs, c.envs()...)

	err = component.Exec.ExecuteInline(c.componentParser.componentArgs, envs...)
	if err != nil {
		return err
	}

	shouldPrint, err := dailyComponentUpdateAvailable(component.Name)
	if err != nil {
		cli.Log.Debugw("unable to load components last check cache", "error", err)
	}
	if shouldPrint && component.ApiInfo != nil && component.InstalledVersion().LessThan(component.ApiInfo.Version) {
		format := "\n%s v%s available: to update, run `lacework component update %s`\n"
		cli.OutputHuman(format, cmd.Use, component.ApiInfo.Version, cmd.Use)
	}

	return nil
}

// LoadComponents reads the local components state and loads all installed components
// of type `CLI_COMMAND` dynamically into the root command of the CLI (`rootCmd`)
func (c *cliState) PrototypeLoadComponents() {
	c.Log.Debugw("loading local components")
	state, err := lwcomponent.LocalState()
	if err != nil || state == nil {
		c.Log.Debugw("unable to load components", "error", err)
		return
	}

	c.LwComponents = state

	// @dhazekamp how do we ensure component command names don't overlap with other commands?

	for _, component := range c.LwComponents.Components {
		if component.IsInstalled() && component.IsCommandType() {
			c.installedCmd = true

			ver, err := component.CurrentVersion()
			if err != nil {
				c.Log.Errorw("unable to load dynamic cli command",
					"component", component.Name, "error", err,
				)
				continue
			}

			c.Log.Debugw("loading dynamic cli command",
				"component", component.Name, "version", ver,
			)
			componentCmd :=
				&cobra.Command{
					Use:                   component.Name,
					Short:                 component.Description,
					Annotations:           map[string]string{"type": componentTypeAnnotation},
					Version:               ver.String(),
					SilenceUsage:          true,
					DisableFlagParsing:    true,
					DisableFlagsInUseLine: true,
					RunE: func(cmd *cobra.Command, args []string) error {
						// cobra will automatically add a -v/--version flag to
						// the command, but because for components we're not
						// parsing the args at the usual point in time, we have
						// to repeat the check for -v here
						versionVal, _ := cmd.Flags().GetBool("version")
						if versionVal {
							cmd.Printf("%s version %s\n", cmd.Use, cmd.Version)
							return nil
						}
						go func() {
							// Start the gRPC server for components to communicate back
							if err := c.Serve(); err != nil {
								c.Log.Errorw("couldn't serve gRPC server", "error", err)
							}
						}()

						c.Log.Debugw("running component", "component", cmd.Use,
							"args", c.componentParser.componentArgs,
							"cli_flags", c.componentParser.cliArgs)
						f, ok := c.LwComponents.GetComponent(cmd.Use)
						if ok {
							shouldPrint, compVerErr := dailyComponentUpdateAvailable(f.Name)
							if compVerErr != nil {
								// Log an error but do not fail
								cli.Log.Debugw("unable to run daily component version check", "error", err)
							}
							if shouldPrint && f.Status() == lwcomponent.UpdateAvailable {
								format := "%s v%s available: to update, run `lacework component update %s`\n"
								cli.OutputHuman(fmt.Sprintf(format, cmd.Use, f.LatestVersion.String(), cmd.Use))
							}
							envs := []string{
								fmt.Sprintf("LW_COMPONENT_NAME=%s", cmd.Use),
							}
							envs = append(envs, c.envs()...)
							return f.RunAndOutput(c.componentParser.componentArgs, envs...)
						}

						// We will land here only if we couldn't run the component, which is not
						// possible since we are adding the components dynamically, still if it
						// happens, let the user know that we would love to hear their feedback
						return errors.New("something went pretty wrong here, contact support@lacework.net")
					},
				}
			rootCmd.AddCommand(componentCmd)
		}
	}
}

func runComponentsList(_ *cobra.Command, _ []string) (err error) {
	if !lwcomponent.CatalogV1Enabled(cli.LwApi) {
		return prototypeRunComponentsList()
	}

	return listComponents()
}

func listComponents() error {
	catalog, err := LoadCatalog("", false)
	if err != nil {
		return errors.Wrap(err, "unable to load component Catalog")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(catalog)
	}

	if catalog.ComponentCount() == 0 {
		msg := "There are no components available, " +
			"come back later or contact support. (version: %s)\n"
		cli.OutputHuman(fmt.Sprintf(msg, cli.LwComponents.Version))

		return nil
	}

	printComponents(catalog.PrintComponents())

	return nil
}

func printComponent(data []string) {
	printComponents([][]string{data})
}

func printComponents(data [][]string) {
	cli.OutputHuman(
		renderCustomTable(
			[]string{"Status", "Name", "Version", "Description"},
			data,
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetColumnSeparator(" ")
				t.SetAutoWrapText(false)
				t.SetAlignment(tablewriter.ALIGN_LEFT)
			}),
		),
	)
}

func runComponentsInstall(cmd *cobra.Command, args []string) (err error) {
	if !lwcomponent.CatalogV1Enabled(cli.LwApi) {
		return prototypeRunComponentsInstall(cmd, args)
	}

	return installComponent(args)
}

func installComponent(args []string) (err error) {
	var (
		componentName    string                 = args[0]
		downloadComplete                        = make(chan int8)
		params           map[string]interface{} = make(map[string]interface{})
		start            time.Time
	)

	cli.Event.Component = componentName
	cli.Event.Feature = "install_component"
	defer cli.SendHoneyvent()

	catalog, err := LoadCatalog(componentName, false)
	if err != nil {
		err = errors.Wrap(err, "unable to load component Catalog")
		return
	}

	component, err := catalog.GetComponent(componentName)
	if err != nil {
		return
	}

	cli.OutputChecklist(successIcon, fmt.Sprintf("Component %s found\n", component.Name))

	cli.StartProgress(fmt.Sprintf("Staging component %s...", componentName))

	start = time.Now()

	progressClosure := func(path string, sizeB int64) {
		downloadProgress(downloadComplete, path, sizeB)
	}

	stageClose, err := catalog.Stage(component, versionArg, progressClosure)
	defer stageClose()
	if err != nil {
		cli.StopProgress()
		return
	}

	downloadComplete <- 0

	params["stage_duration_ms"] = time.Since(start).Milliseconds()
	cli.Event.FeatureData = params

	cli.StopProgress()
	if err != nil {
		return
	}
	cli.OutputChecklist(successIcon, "Component %s staged\n", color.HiYellowString(componentName))

	cli.StartProgress("Verifing component signature...")

	err = catalog.Verify(component)

	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "verification of component signature failed")
		return
	}
	cli.OutputChecklist(successIcon, "Component signature verified\n")

	cli.StartProgress("Installing component...")

	err = catalog.Install(component)

	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "Install of component failed")
		return
	}
	cli.OutputChecklist(successIcon, "Component version %s installed\n", component.InstalledVersion())

	cli.StartProgress("Configuring component...")

	stdout, stderr, errCmd := component.Exec.Execute([]string{"cdk-init"}, cli.envs()...)
	if errCmd != nil {
		if errCmd != lwcomponent.ErrNonExecutable {
			cli.Log.Warnw("component life cycle",
				"error", errCmd.Error(), "stdout", stdout, "stderr", stderr)
		}
	} else {
		cli.Log.Infow("component life cycle", "stdout", stdout, "stderr", stderr)
	}
	cli.StopProgress()

	cli.OutputChecklist(successIcon, "Component configured\n")
	cli.OutputHuman("\nInstallation completed.\n")

	if component.InstallMessage != "" {
		cli.OutputHuman(fmt.Sprintf("\n%s\n", component.InstallMessage))
	}

	return
}

func runComponentsShow(_ *cobra.Command, args []string) (err error) {
	if !lwcomponent.CatalogV1Enabled(cli.LwApi) {
		return prototypeRunComponentsShow(args)
	}

	return showComponent(args)
}

func showComponent(args []string) error {
	var (
		componentName string = args[0]
	)

	catalog, err := LoadCatalog(componentName, true)
	if err != nil {
		return errors.Wrap(err, "unable to load component Catalog")
	}

	component, err := catalog.GetComponent(componentName)
	if err != nil {
		return err
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(component)
	}

	printComponent(component.PrintSummary())

	allVersions, err := catalog.ListComponentVersions(component)
	if err != nil {
		return err
	}

	printAvailableVersions(component.InstalledVersion(), allVersions)

	return nil
}

func printAvailableVersions(installedVersion *semver.Version, availableVersions []*semver.Version) {
	cli.OutputHuman("\n")

	result := "The following versions of this component are available to install:"
	foundInstalled := false

	for _, version := range availableVersions {
		result += "\n"
		result += " - " + version.String()
		if installedVersion != nil && version.Equal(installedVersion) {
			result += " (installed)"
			foundInstalled = true
		}
	}

	if installedVersion != nil && !foundInstalled {
		result += fmt.Sprintf(
			"\n\nThe currently installed version %s is no longer available to install.",
			installedVersion.String(),
		)
	}

	cli.OutputHuman(result)
	cli.OutputHuman("\n")
}

func runComponentsUpdate(_ *cobra.Command, args []string) (err error) {
	if !lwcomponent.CatalogV1Enabled(cli.LwApi) {
		return prototypeRunComponentsUpdate(args)
	}

	return updateComponent(args)
}

func updateComponent(args []string) (err error) {
	var (
		componentName    string                 = args[0]
		downloadComplete                        = make(chan int8)
		params           map[string]interface{} = make(map[string]interface{})
		start            time.Time
		targetVersion    *semver.Version
	)

	catalog, err := LoadCatalog(componentName, false)
	if err != nil {
		return errors.Wrap(err, "unable to load component Catalog")
	}

	component, err := catalog.GetComponent(componentName)
	if err != nil {
		return err
	}

	cli.OutputChecklist(successIcon, fmt.Sprintf("Component %s found\n", component.Name))

	installedVersion := component.InstalledVersion()
	if installedVersion == nil {
		return errors.Errorf("component %s not installed", color.HiYellowString(componentName))
	}

	latestVersion := component.LatestVersion()
	if latestVersion == nil {
		return errors.Errorf("component %s not available in API", color.HiYellowString(componentName))
	}

	if versionArg == "" {
		targetVersion = latestVersion
	} else {
		targetVersion, err = semver.NewVersion(versionArg)
		if err != nil {
			return errors.Errorf("invalid semantic version %s", versionArg)
		}
	}

	if installedVersion.Equal(targetVersion) {
		return errors.Errorf("You are already running version %s of this component",
			color.HiYellowString(installedVersion.String()))
	}

	cli.StartProgress(fmt.Sprintf("Staging component %s...", color.HiYellowString(componentName)))

	start = time.Now()

	progressClosure := func(path string, sizeB int64) {
		downloadProgress(downloadComplete, path, sizeB)
	}

	stageClose, err := catalog.Stage(component, versionArg, progressClosure)
	defer stageClose()
	if err != nil {
		cli.StopProgress()
		return
	}

	downloadComplete <- 0

	params["stage_duration_ms"] = time.Since(start).Milliseconds()
	cli.Event.FeatureData = params

	cli.StopProgress()
	if err != nil {
		return
	}
	cli.OutputChecklist(successIcon, "Component %s staged\n", color.HiYellowString(componentName))

	cli.StartProgress("Verifing component signature...")

	err = catalog.Verify(component)

	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "verification of component signature failed")
		return
	}
	cli.OutputChecklist(successIcon, "Component signature verified\n")

	cli.StartProgress(fmt.Sprintf("Updating component %s to version %s...", component.Name, targetVersion.String()))

	err = catalog.Install(component)

	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "Update of component failed")
		return
	}

	cli.OutputChecklist(successIcon, "Component %s updated to %s\n",
		color.HiYellowString(component.Name),
		color.HiCyanString(targetVersion.String()))

	cli.StartProgress("Configuring component...")

	stdout, stderr, errCmd := component.Exec.Execute([]string{"cdk-reconfigure"}, cli.envs()...)
	if errCmd != nil {
		if errCmd != lwcomponent.ErrNonExecutable {
			cli.Log.Warnw("component life cycle",
				"error", errCmd.Error(), "stdout", stdout, "stderr", stderr)
		}
	} else {
		cli.Log.Infow("component life cycle", "stdout", stdout, "stderr", stderr)
	}
	cli.StopProgress()

	cli.OutputChecklist(successIcon, "Component reconfigured\n")

	if component.UpdateMessage != "" {
		cli.OutputHuman(fmt.Sprintf("\n%s\n", component.UpdateMessage))
	}

	return
}

func runComponentsDelete(_ *cobra.Command, args []string) (err error) {
	if !lwcomponent.CatalogV1Enabled(cli.LwApi) {
		return prototypeRunComponentsDelete(args)
	}

	return deleteComponent(args)
}

func deleteComponent(args []string) (err error) {
	var (
		componentName string = args[0]
	)

	catalog, err := LoadCatalog(componentName, false)
	if err != nil {
		return errors.Wrap(err, "unable to load component Catalog")
	}

	component, err := catalog.GetComponent(componentName)
	if err != nil {
		return err
	}

	cli.OutputChecklist(successIcon, fmt.Sprintf("Component %s found\n", component.Name))

	cli.StartProgress("Cleaning component data...")

	stdout, stderr, errCmd := component.Exec.Execute([]string{"cdk-cleanup"}, cli.envs()...)
	if errCmd != nil {
		if errCmd != lwcomponent.ErrNonExecutable {
			cli.Log.Warnw("component life cycle",
				"error", errCmd.Error(), "stdout", stdout, "stderr", stderr)
		}
	} else {
		cli.Log.Infow("component life cycle", "stdout", stdout, "stderr", stderr)
	}
	cli.StopProgress()

	cli.OutputChecklist(successIcon, "Component data removed\n")

	cli.StartProgress("Deleting component...")
	defer cli.StopProgress()

	err = catalog.Delete(component)
	if err != nil {
		return
	}

	cli.StopProgress()

	cli.OutputChecklist(successIcon, "Component %s deleted\n", color.HiYellowString(component.Name))

	msg := fmt.Sprintf(`\n- We will do better next time.\n\nDo you want to provide feedback?\nReach out to us at %s\n`,
		color.HiCyanString("support@lacework.net"))

	cli.OutputHuman(msg)

	return
}

func prototypeRunComponentsList() (err error) {
	cli.StartProgress("Loading components state...")
	cli.LwComponents, err = lwcomponent.LoadState(cli.LwApi)
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to list components")
		return
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(cli.LwComponents)
	}

	if len(cli.LwComponents.Components) == 0 {
		msg := "There are no components available, " +
			"come back later or contact support. (version: %s)\n"
		cli.OutputHuman(fmt.Sprintf(msg, cli.LwComponents.Version))
		return
	}

	cli.OutputHuman(
		renderCustomTable(
			[]string{"Status", "Name", "Version", "Description"},
			componentsToTable(),
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetColumnSeparator(" ")
				t.SetAutoWrapText(false)
				t.SetAlignment(tablewriter.ALIGN_LEFT)
			}),
		),
	)

	cli.OutputHuman("\nComponents version: %s\n", cli.LwComponents.Version)
	return
}

func prototypeRunComponentsShow(args []string) (err error) {
	cli.StartProgress("Loading components state...")
	cli.LwComponents, err = lwcomponent.LoadState(cli.LwApi)
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to load state of components")
		return
	}
	component, found := cli.LwComponents.GetComponent(args[0])
	if !found {
		return errors.New("component not found")
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(component)
	}

	var colorize *color.Color
	switch component.Status() {
	case lwcomponent.NotInstalled:
		colorize = color.New(color.FgWhite, color.Bold)
	case lwcomponent.Installed:
		colorize = color.New(color.FgGreen, color.Bold)
	case lwcomponent.UpdateAvailable:
		colorize = color.New(color.FgYellow, color.Bold)
	}

	cli.OutputHuman(
		renderCustomTable(
			[]string{"Name", "Status", "Description"},
			[][]string{{
				component.Name,
				colorize.Sprintf(component.Status().String()),
				component.Description,
			}},
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetColumnSeparator(" ")
				t.SetAutoWrapText(false)
				t.SetAlignment(tablewriter.ALIGN_LEFT)
			}),
		),
	)
	var currentVersion *semver.Version = nil
	if component.Status() == lwcomponent.Installed || component.Status() == lwcomponent.UpdateAvailable {
		installed, err := component.CurrentVersion()
		if err != nil {
			return err
		}
		currentVersion = installed
	}
	cli.OutputHuman("\n")
	cli.OutputHuman(component.ListVersions(currentVersion))
	cli.OutputHuman("\n")
	return
}

func componentsToTable() [][]string {
	out := [][]string{}
	for _, cdata := range cli.LwComponents.Components {
		var colorize *color.Color
		switch cdata.Status() {
		case lwcomponent.NotInstalled:
			colorize = color.New(color.FgWhite, color.Bold)
		case lwcomponent.Installed:
			colorize = color.New(color.FgGreen, color.Bold)
		case lwcomponent.UpdateAvailable:
			colorize = color.New(color.FgYellow, color.Bold)
		}

		// by default, we display the latest version
		version := cdata.LatestVersion.String()

		// but if the component is installed,
		// we display the current version instead
		if currentVersion, err := cdata.CurrentVersion(); err == nil {
			version = currentVersion.String()
		}

		out = append(out, []string{
			colorize.Sprintf(cdata.Status().String()),
			cdata.Name,
			version,
			cdata.Description,
		})
	}
	return out
}

func prototypeRunComponentsInstall(_ *cobra.Command, args []string) (err error) {
	var (
		componentName    string                 = args[0]
		downloadComplete                        = make(chan int8)
		version          string                 = versionArg
		params           map[string]interface{} = make(map[string]interface{})
		start            time.Time
	)

	cli.Event.Component = componentName
	cli.Event.Feature = "install_component"
	defer cli.SendHoneyvent()

	cli.StartProgress("Loading components state...")
	// @afiune maybe move the state to the cache and fetch if it if has expired
	cli.LwComponents, err = lwcomponent.LoadState(cli.LwApi)
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to load components")
		return
	}

	component, found := cli.LwComponents.GetComponent(componentName)
	if !found {
		err = errors.New(fmt.Sprintf("component %s not found. Try running 'lacework component list'", componentName))
		return
	}

	cli.OutputChecklist(successIcon, fmt.Sprintf("Component %s found\n", componentName))

	if version == "" {
		version = component.LatestVersion.String()
	}

	start = time.Now()

	progressClosure := func(path string, sizeB int64) {
		downloadProgress(downloadComplete, path, sizeB)
	}

	cli.StartProgress(fmt.Sprintf("Installing component %s...", component.Name))
	err = cli.LwComponents.Install(component, version, progressClosure)
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to install component")
		return
	}

	downloadComplete <- 0

	cli.OutputChecklist(successIcon, "Component %s installed\n", color.HiYellowString(component.Name))

	params["install_duration_ms"] = time.Since(start).Milliseconds()

	start = time.Now()

	cli.StartProgress("Verifing component signature...")
	err = cli.LwComponents.Verify(component, version)
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "verification of component signature failed")
		return
	}
	cli.OutputChecklist(successIcon, "Component signature verified\n")

	params["verify_duration_ms"] = time.Since(start).Milliseconds()

	start = time.Now()

	cli.StartProgress(fmt.Sprintf("Configuring component %s...", component.Name))
	// component life cycle: initialize
	stdout, stderr, errCmd := component.RunAndReturn([]string{"cdk-init"}, nil, cli.envs()...)
	if errCmd != nil {
		cli.Log.Warnw("component life cycle",
			"error", errCmd.Error(), "stdout", stdout, "stderr", stderr)
	} else {
		cli.Log.Infow("component life cycle", "stdout", stdout, "stderr", stderr)
	}
	cli.StopProgress()

	cli.OutputChecklist(successIcon, "Component configured\n")
	cli.OutputHuman("\nInstallation completed.\n")

	params["configure_duration_ms"] = time.Since(start).Milliseconds()

	cli.Event.FeatureData = params

	if component.Breadcrumbs.InstallationMessage != "" {
		cli.OutputHuman("\n")
		cli.OutputHuman(component.Breadcrumbs.InstallationMessage)
		cli.OutputHuman("\n")
	}
	return
}

func prototypeRunComponentsUpdate(args []string) (err error) {
	cli.StartProgress("Loading components state...")
	// @afiune maybe move the state to the cache and fetch if it if has expired
	cli.LwComponents, err = lwcomponent.LoadState(cli.LwApi)
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to load components")
		return
	}

	component_name := args[0]

	component, found := cli.LwComponents.GetComponent(args[0])
	if !found {
		err = errors.New(fmt.Sprintf("component %s not found. Try running 'lacework component list'", component_name))
		return
	}
	// @afiune end boilerplate load components

	cli.OutputChecklist(successIcon, fmt.Sprintf("Component %s found\n", component_name))

	updateTo := component.LatestVersion
	if versionArg != "" {
		parsedVersion, err := semver.NewVersion(versionArg)
		if err != nil {
			err = errors.Wrap(err, "invalid version specified")
			return err
		}
		updateTo = *parsedVersion
	}

	currentVersion, err := component.CurrentVersion()
	if err != nil {
		return err
	}

	if currentVersion.Equal(&updateTo) {
		cli.OutputHuman("You are already running version %s of this component", currentVersion.String())
		return nil
	}

	downloadComplete := make(chan int8)

	progressClosure := func(path string, sizeB int64) {
		downloadProgress(downloadComplete, path, sizeB)
	}

	cli.StartProgress(fmt.Sprintf("Updating component %s to version %s...", component.Name, &updateTo))
	err = cli.LwComponents.Install(component, updateTo.String(), progressClosure)
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to update component")
		return
	}

	downloadComplete <- 0

	cli.OutputChecklist(successIcon, "Component %s updated to %s\n",
		color.HiYellowString(component.Name),
		color.HiCyanString(fmt.Sprintf("v%s", updateTo.String())))

	cli.StartProgress("Verifing component signature...")
	err = cli.LwComponents.Verify(component, updateTo.String())
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "verification of component signature failed")
		return
	}
	cli.OutputChecklist(successIcon, "Component signature verified\n")

	cli.StartProgress(fmt.Sprintf("Reconfiguring %s component...", component.Name))
	// component life cycle: reconfigure
	stdout, stderr, errCmd := component.RunAndReturn(
		[]string{"cdk-reconfigure", currentVersion.String(), updateTo.String()},
		nil, cli.envs()...)
	if errCmd != nil {
		cli.Log.Warnw("component life cycle",
			"error", errCmd.Error(), "stdout", stdout, "stderr", stderr)
	} else {
		cli.Log.Infow("component life cycle", "stdout", stdout, "stderr", stderr)
	}
	cli.StopProgress()

	cli.OutputChecklist(successIcon, "Component reconfigured\n")
	cli.OutputHuman("\n")
	cli.OutputHuman(component.MakeUpdateMessage(*currentVersion, updateTo))
	cli.OutputHuman("\n")
	return
}

func prototypeRunComponentsDelete(args []string) (err error) {
	cli.StartProgress("Loading components state...")
	// @afiune maybe move the state to the cache and fetch if it if has expired
	// @afiune DO WE NEED THIS? It should already be loaded
	cli.LwComponents, err = lwcomponent.LocalState()
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to load components")
		return
	}

	component, found := cli.LwComponents.GetComponent(args[0])
	if !found {
		err = errors.New("component not found. Try running 'lacework component list'")
		return
	}

	if component.UnderDevelopment() {
		cli.OutputHuman("Component '%s' in under development. Bypassing checks.\n\n",
			color.HiYellowString(component.Name))
	} else if !component.IsInstalled() {
		err = errors.Errorf(
			"component not installed. Try running 'lacework component install %s'",
			args[0],
		)
		return
	}

	cli.StartProgress("Cleaning component data...")
	// component life cycle: cleanup
	stdout, stderr, errCmd := component.RunAndReturn([]string{"cdk-cleanup"}, nil, cli.envs()...)
	if errCmd != nil {
		cli.Log.Warnw("component life cycle",
			"error", errCmd.Error(), "stdout", stdout, "stderr", stderr)
	} else {
		cli.Log.Infow("component life cycle", "stdout", stdout, "stderr", stderr)
	}
	cli.StopProgress()

	cli.OutputChecklist(successIcon, "Component data removed\n")

	cli.StartProgress("Deleting component...")
	defer cli.StopProgress()

	cPath, err := component.RootPath()
	if err != nil {
		err = errors.Wrap(err, "unable to delete component")
		return
	}

	err = os.RemoveAll(cPath)
	if err != nil {
		err = errors.Wrap(err, "unable to delete component")
		return
	}

	cli.StopProgress()

	cli.OutputChecklist(successIcon, "Component %s deleted\n", color.HiYellowString(component.Name))
	cli.OutputHuman("\n- We will do better next time.\n")
	cli.OutputHuman("\nDo you want to provide feedback?\n")
	cli.OutputHuman("Reach out to us at %s\n", color.HiCyanString("support@lacework.net"))
	return
}

func downloadProgress(complete chan int8, path string, sizeB int64) {
	file, err := os.Open(path)
	if err != nil {
		cli.Log.Errorf("Failed to open component file: %s", err.Error())
		return
	}
	defer file.Close()

	var (
		previous      float64 = 0
		stop          bool    = false
		spinnerSuffix string  = ""
	)

	if !cli.nonInteractive {
		spinnerSuffix = cli.spinner.Suffix
	}

	for !stop {
		select {
		case <-complete:
			stop = true
		default:
			info, err := file.Stat()
			if err != nil {
				cli.Log.Errorf("Failed to stat component file: %s", err.Error())
				return
			}

			size := info.Size()
			if size == 0 {
				size = 1
			}

			if sizeB == 0 {
				mb := float64(size) / (1 << 20)

				if mb > previous {
					if !cli.nonInteractive {
						cli.spinner.Suffix = fmt.Sprintf("%s Downloaded: %.0fmb", spinnerSuffix, mb)
					} else {
						cli.OutputHuman("..Downloaded: %.0fmb\n", mb)
					}

					previous = mb
				}
			} else {
				percent := float64(size) / float64(sizeB) * 100

				if percent > previous {
					if !cli.nonInteractive {
						cli.spinner.Suffix = fmt.Sprintf("%s Downloaded: %.0f%s", spinnerSuffix, percent, "%")
					} else {
						cli.OutputHuman("..Downloaded: %.0f%s\n", percent, "%")
					}

					previous = percent
				}
			}
		}

		time.Sleep(time.Second)
	}
}

func LoadCatalog(componentName string, getAllVersions bool) (*lwcomponent.Catalog, error) {
	cli.StartProgress("Loading component catalog...")
	defer cli.StopProgress()

	var componentsApiInfo map[string]*lwcomponent.ApiInfo

	// try to load components Catalog from cache
	if !cli.noCache {
		expired := cli.ReadCachedAsset(componentsCacheKey, &componentsApiInfo)
		if !expired && !getAllVersions {
			cli.Log.Infow("loaded components from cache", "components", componentsApiInfo)
			return lwcomponent.NewCachedCatalog(cli.LwApi, lwcomponent.NewStageTarGz, componentsApiInfo)
		}
	}

	// load components Catalog from API
	catalog, err := lwcomponent.NewCatalog(cli.LwApi, lwcomponent.NewStageTarGz)
	if err != nil {
		return nil, err
	}

	// Retrieve the list of all available versions for a single component
	if getAllVersions {
		component, err := catalog.GetComponent(componentName)
		if err != nil {
			return nil, err
		}

		vers, err := catalog.ListComponentVersions(component)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to fetch component '%s' versions", componentName))
		}

		component.ApiInfo.AllVersions = vers
		catalog.Components[componentName] = lwcomponent.NewCDKComponent(
			component.Name,
			component.Description,
			component.Type,
			component.ApiInfo,
			component.HostInfo)
	}

	componentsApiInfo = make(map[string]*lwcomponent.ApiInfo, len(catalog.Components))

	for _, c := range catalog.Components {
		if c.ApiInfo != nil {
			componentsApiInfo[c.Name] = c.ApiInfo
		}
	}

	cli.WriteAssetToCache(componentsCacheKey, time.Now().Add(time.Hour*12), componentsApiInfo)

	return catalog, nil
}
