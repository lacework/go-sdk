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
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/lacework/go-sdk/lwcomponent"
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
)

func init() {
	// add the components command
	rootCmd.AddCommand(componentsCmd)

	// add sub-commands to the components command
	componentsCmd.AddCommand(componentsListCmd)
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
	if found && t == "component" {
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
				c.Log.Warnw("unable to load dynamic cli command",
					"component", component.Name,
					"error", err.Error(),
				)
				continue
			}

			c.Log.Debugw("loading dynamic cli command",
				"component", component.Name, "version", ver,
			)
			rootCmd.AddCommand(
				&cobra.Command{
					// @afiune strip `lw-` from component?
					Use:                   component.Name,
					Short:                 component.Description,
					Annotations:           map[string]string{"type": "component"},
					Version:               ver.String(),
					SilenceUsage:          true,
					DisableFlagParsing:    true,
					DisableFlagsInUseLine: true,
					RunE: func(cmd *cobra.Command, args []string) error {
						globalFlags := []*pflag.Flag{}
						cmd.Flags().VisitAll(func(f *pflag.Flag) {
							// At runtime, we visit all global flags defined at the root command
							// and we essentially make a slice of the flag names with their shorthand
							// so that we can remove them all from the raw arguments (`args`)
							globalFlags = append(globalFlags, f)
						})

						// Use the list of global CLI flags to filter the provided arguments and return
						// the filteres arguments to pass to the underlying component and the real list
						// of CLI flags. The later is used to parse and then run the global CLI init func
						filteredArgs, filteredCLIFlags := filterCLIFlagsFromComponentArgs(args, globalFlags)

						// Parse all global CLI flags provided (filtered) by the user, then run the global
						// CLI init function to initialize our logger, api client, and other global config
						err := cmd.Flags().Parse(filteredCLIFlags)
						initConfig() // @afiune NOTE we purposely run this func first and then check the err
						if err != nil {
							cli.Log.Debugw("unable to parse global flags",
								"provided_flags", filteredCLIFlags, "error", err)
						}

						cli.Log.Debugw("running component", "component", cmd.Use,
							"args", filteredArgs, "cli_flags", filteredCLIFlags)
						f, ok := cli.LwComponents.GetComponent(cmd.Use)
						if ok {
							// @afiune what if the component needs other env variables
							envs := []string{fmt.Sprintf("LW_COMPONENT_NAME=%s", cmd.Use)}
							envs = append(envs, c.envs()...)
							return f.RunAndOutput(filteredArgs, envs...)
						}

						// We will land here only if we couldn't run the component, which is not
						// possible since we are adding the components dynamically, still if it
						// happens, let the user know that we would love to hear their feedback
						return errors.New("something went pretty wrong here, contact support@lacework.net")
					},
				},
			)
		}
	}
}

// filterCLIFlagsFromComponentArgs uses the arguments provided by the user and
// the list of global CLI flags to return the real list of component arguments
// and the real list of CLI flags provided as arguments
func filterCLIFlagsFromComponentArgs(args []string, globalFlags []*pflag.Flag) (
	componentArgs []string, cliFlags []string,
) {

	// this variable is used to store a flag of type `string` so that we can check
	// the next argument and pass it as the value of the flag
	stringFlag := ""

	for _, arg := range args {

		// if the stringFlag variable is not empty it means that the current argument
		// is the value of the provided string flag, add it, empty it and move on
		if stringFlag != "" {
			cliFlags = append(cliFlags, stringFlag, arg) // add the flag and value
			stringFlag = ""                              // empty the flag
			continue                                     // move to the next argument
		}

		// assume the argument is an argument unless it is a flag
		isArg := true

		// flags must have a prefix of `--` or `-`, if the argument has that prefix it
		// means that it could be a flag, strip it and compare it with all global flags
		argFlag := strings.TrimPrefix(strings.TrimPrefix(arg, "--"), "-")

		for _, flag := range globalFlags {
			if flag == nil { // avoid panics trying to access pointer
				continue
			}

			if flag.Name == argFlag || flag.Shorthand == argFlag {
				// the argument is indeed a flag

				if flag.Value == nil { // avoid panics trying to access interface
					continue
				}

				switch flag.Value.Type() { // check the type
				case "bool":
					isArg = false
					cliFlags = append(cliFlags, arg)
				case "string":
					isArg = false
					stringFlag = arg
				}
			}
		}

		if isArg {
			// the argument is actually an argument
			componentArgs = append(componentArgs, arg)
		}
	}

	return
}

func runComponentsList(_ *cobra.Command, _ []string) (err error) {
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

func runComponentsInstall(_ *cobra.Command, args []string) (err error) {
	cli.StartProgress("Loading components state...")
	// @afiune maybe move the state to the cache and fetch if it if has expired
	cli.LwComponents, err = lwcomponent.LoadState(cli.LwApi)
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

	cli.StartProgress(fmt.Sprintf("Installing component %s...", component.Name))
	err = cli.LwComponents.Install(args[0])
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to install component")
		return
	}
	cli.OutputChecklist(successIcon, "Component %s installed\n", color.HiYellowString(component.Name))
	cli.OutputChecklist(successIcon, "Signature verified\n")

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

	if component.Breadcrumbs.InstallationMessage != "" {
		cli.OutputHuman("\n")
		cli.OutputHuman(component.Breadcrumbs.InstallationMessage)
	}
	return
}

func runComponentsUpdate(_ *cobra.Command, args []string) (err error) {
	cli.StartProgress("Loading components state...")
	// @afiune maybe move the state to the cache and fetch if it if has expired
	cli.LwComponents, err = lwcomponent.LoadState(cli.LwApi)
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
	// @afiune end boilerplate load components
	update, err := component.UpdateAvailable()
	if err != nil {
		return err
	}

	if !update {
		cli.OutputHuman(
			"You are running the latest version of the component %s.\n", args[0],
		)
		return nil
	}

	currentVersion, err := component.CurrentVersion()
	if err != nil {
		return err
	}

	cli.StartProgress(fmt.Sprintf("Updating component %s...", component.Name))
	err = cli.LwComponents.Install(args[0])
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to update component")
		return
	}
	cli.OutputChecklist(successIcon, "Component %s updated to %s\n",
		color.HiYellowString(component.Name),
		color.HiCyanString(fmt.Sprintf("v%s", component.LatestVersion.String())))
	cli.OutputChecklist(successIcon, "Signature verified\n")

	cli.StartProgress(fmt.Sprintf("Reconfiguring %s component...", component.Name))
	// component life cycle: reconfigure
	stdout, stderr, errCmd := component.RunAndReturn(
		[]string{"cdk-reconfigure", currentVersion.String(), component.LatestVersion.String()},
		nil, cli.envs()...)
	if errCmd != nil {
		cli.Log.Warnw("component life cycle",
			"error", errCmd.Error(), "stdout", stdout, "stderr", stderr)
	} else {
		cli.Log.Infow("component life cycle", "stdout", stdout, "stderr", stderr)
	}
	cli.StopProgress()

	cli.OutputChecklist(successIcon, "Component reconfigured\n")
	cli.OutputHuman("\nUpdate completed.\n")

	if component.Breadcrumbs.UpdateMessage != "" {
		cli.OutputHuman("\n")
		cli.OutputHuman(component.Breadcrumbs.UpdateMessage)
	}
	return
}

func runComponentsDelete(_ *cobra.Command, args []string) (err error) {
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
	// component life cycle: remove
	stdout, stderr, errCmd := component.RunAndReturn([]string{"cdk-remove"}, nil, cli.envs()...)
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

func runComponentsDevMode(_ *cobra.Command, args []string) error {
	cli.StartProgress("Loading components state...")
	var err error
	cli.LwComponents, err = lwcomponent.LoadState(cli.LwApi)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to load components")
	}

	component, found := cli.LwComponents.GetComponent(args[0])
	if !found {
		component = &lwcomponent.Component{
			Name: args[0],
		}

		if component.UnderDevelopment() {
			return errors.New("component already under development.")
		}

		cli.OutputHuman("Component '%s' not found. Defining a new component.\n",
			color.HiYellowString(component.Name))

		var (
			cType   string
			helpMsg = fmt.Sprintf("What are these component types ?\n"+
				"\n'%s' - A regular standalone-binary (this component type is not accessible via the CLI)"+
				"\n'%s' - A binary accessible via the Lacework CLI (Users will run 'lacework <COMPONENT_NAME>')"+
				"\n'%s' - A library that only provides content for the CLI or other components\n",
				lwcomponent.BinaryType, lwcomponent.CommandType, lwcomponent.LibraryType)
		)
		if err := survey.AskOne(&survey.Select{
			Message: "Select the type of component you are developing:",
			Help:    helpMsg,
			Options: []string{
				lwcomponent.BinaryType,
				lwcomponent.CommandType,
				lwcomponent.LibraryType,
			},
		}, &cType); err != nil {
			return err
		}

		component.Type = lwcomponent.Type(cType)

		if err := survey.AskOne(&survey.Input{
			Message: "What is this component about? (component description):",
		}, &component.Description); err != nil {
			return err
		}
	}

	if err := component.EnterDevelopmentMode(); err != nil {
		return errors.Wrap(err, "unable to enter development mode")
	}

	rPath, err := component.RootPath()
	if err != nil {
		return errors.New("unable to detect RootPath")
	}

	cli.OutputHuman("Component '%s' in now in development mode.\n\n",
		color.HiYellowString(component.Name))
	cli.OutputHuman("Root path: %s\n", rPath)
	cli.OutputHuman("Dev specs: %s\n", filepath.Join(rPath, ".dev"))
	if component.Type == lwcomponent.CommandType {
		cli.OutputHuman("\nDeploy your dev component at: %s\n",
			color.HiYellowString(filepath.Join(rPath, component.Name)))
	}
	return nil
}
