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

// LoadComponents reads the local components state and loads all installed components
// of type `CLI_COMMAND` dynamically into the root command of the CLI (`rootCmd`)
func (c *cliState) LoadComponents() {
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
							if f.Status() == lwcomponent.UpdateAvailable {
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
	cli.StartProgress("Loading component Catalog...")

	catalog, err := lwcomponent.NewCatalog(cli.LwApi, lwcomponent.NewStageTarGz)
	defer catalog.Cache()

	cli.StopProgress()
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

	cli.OutputHuman("\nComponents version: %s\n", cli.LwComponents.Version)

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

	return installComponent(cmd, args)
}

func installComponent(cmd *cobra.Command, args []string) (err error) {
	var (
		componentName string                 = args[0]
		params        map[string]interface{} = make(map[string]interface{})
		start         time.Time
	)

	cli.Event.Component = componentName
	cli.Event.Feature = "install_component"
	defer cli.SendHoneyvent()

	cli.StartProgress("Loading component Catalog...")

	catalog, err := lwcomponent.NewCatalog(cli.LwApi, lwcomponent.NewStageTarGz)
	defer catalog.Cache()

	cli.StopProgress()
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

	stageClose, err := catalog.Stage(component, versionArg)
	if err != nil {
		return
	}
	defer stageClose()

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
	cli.OutputChecklist(successIcon, "Component installed\n")

	// @jon-stewart: TODO: Component lifecycle `cdk-init` command

	// @jon-stewart: TODO: print install message

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

	cli.StartProgress("Loading components Catalog...")

	catalog, err := lwcomponent.NewCatalog(cli.LwApi, lwcomponent.NewStageTarGz)
	defer catalog.Cache()

	cli.StopProgress()
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

	version := component.InstalledVersion()

	availableVersions, err := catalog.ListComponentVersions(component)
	if err != nil {
		return err
	}

	printAvailableVersions(version, availableVersions)

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
		componentName string                 = args[0]
		params        map[string]interface{} = make(map[string]interface{})
		start         time.Time
		targetVersion *semver.Version
	)

	cli.StartProgress("Loading components Catalog...")

	catalog, err := lwcomponent.NewCatalog(cli.LwApi, lwcomponent.NewStageTarGz)
	defer catalog.Cache()

	cli.StopProgress()
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
		cli.OutputHuman("component %s not installed", component.Name)
		return nil
	}

	latestVersion := component.LatestVersion()
	if latestVersion == nil {
		cli.OutputHuman("component %s not available in API", component.Name)
		return nil
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
		cli.OutputHuman("You are already running version %s of this component", installedVersion.String())
		return nil
	}

	cli.StartProgress(fmt.Sprintf("Staging component %s...", color.HiYellowString(componentName)))

	start = time.Now()

	stageClose, err := catalog.Stage(component, versionArg)
	if err != nil {
		return
	}
	defer stageClose()

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

	cli.OutputChecklist(successIcon, "Component %s updated to %s\n", color.HiYellowString(component.Name), color.HiCyanString(targetVersion.String()))

	// @jon-stewart: TODO: component lifecycle event

	// @jon-stewart: TODO: component update message

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

	cli.StartProgress("Loading components Catalog...")

	catalog, err := lwcomponent.NewCatalog(cli.LwApi, lwcomponent.NewStageTarGz)
	defer catalog.Cache()

	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to load component Catalog")
	}

	component, err := catalog.GetComponent(componentName)
	if err != nil {
		return err
	}

	cli.OutputChecklist(successIcon, fmt.Sprintf("Component %s found\n", component.Name))

	// @jon-stewart: TODO: component life cycle: cleanup

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

		out = append(out, []string{
			colorize.Sprintf(cdata.Status().String()),
			cdata.Name,
			cdata.LatestVersion.String(),
			cdata.Description,
		})
	}
	return out
}

func prototypeRunComponentsInstall(cmd *cobra.Command, args []string) (err error) {
	var (
		componentName string                 = args[0]
		version       string                 = versionArg
		params        map[string]interface{} = make(map[string]interface{})
		start         time.Time
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

	cli.StartProgress(fmt.Sprintf("Installing component %s...", component.Name))
	err = cli.LwComponents.Install(component, version)
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to install component")
		return
	}
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

	cli.StartProgress(fmt.Sprintf("Updating component %s to version %s...", component.Name, &updateTo))
	err = cli.LwComponents.Install(component, updateTo.String())
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to update component")
		return
	}
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
