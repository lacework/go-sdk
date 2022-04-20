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

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

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
)

func init() {
	// add the components command
	rootCmd.AddCommand(componentsCmd)

	// add sub-commands to the components command
	componentsCmd.AddCommand(componentsListCmd)
	componentsCmd.AddCommand(componentsInstallCmd)
	componentsCmd.AddCommand(componentsUpdateCmd)
	componentsCmd.AddCommand(componentsUninstallCmd)

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
					Use: component.Name,
					// @afiune should we build this or add it as component specification
					Short:                 fmt.Sprintf("component %s", component.Name),
					Long:                  component.Description,
					Annotations:           map[string]string{"type": "component"},
					Version:               ver.String(),
					DisableFlagParsing:    true,
					DisableFlagsInUseLine: true,
					RunE: func(cmd *cobra.Command, args []string) error {
						cli.Log.Debugw("running component", "component", cmd.Use, "args", args)
						// @afiune what if the component needs other env variables
						return component.RunAndOutput(args, c.envs()...)
					},
				},
			)
		}
	}
}
func runComponentsList(_ *cobra.Command, _ []string) (err error) {
	cli.LwComponents, err = lwcomponent.LoadState(cli.LwApi)
	if err != nil {
		err = errors.Wrap(err, "unable to list components")
		return
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
		out = append(out, []string{
			cdata.Status().String(),
			cdata.Name,
			cdata.LatestVersion.String(),
			cdata.Description,
		})
	}
	return out
}

func runComponentsInstall(_ *cobra.Command, args []string) (err error) {
	// @afiune maybe move the state to the cache and fetch if it if has expired
	cli.LwComponents, err = lwcomponent.LoadState(cli.LwApi)
	if err != nil {
		err = errors.Wrap(err, "unable to load components")
		return
	}

	component := cli.LwComponents.GetComponent(args[0])
	if component == nil {
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

	cli.OutputHuman("The component %s was installed.\n", args[0])
	return
}

func runComponentsUpdate(_ *cobra.Command, args []string) (err error) {
	// @afiune maybe move the state to the cache and fetch if it if has expired
	cli.LwComponents, err = lwcomponent.LoadState(cli.LwApi)
	if err != nil {
		err = errors.Wrap(err, "unable to load components")
		return
	}

	component := cli.LwComponents.GetComponent(args[0])
	if component == nil {
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

	cli.StartProgress(fmt.Sprintf("Updating component %s...", component.Name))
	err = cli.LwComponents.Install(args[0])
	cli.StopProgress()
	if err != nil {
		err = errors.Wrap(err, "unable to update component")
		return
	}

	cli.OutputHuman("The component %s was updated.\n", args[0])
	return
}

func runComponentsDelete(_ *cobra.Command, args []string) (err error) {
	// @afiune maybe move the state to the cache and fetch if it if has expired
	// @afiune DO WE NEED THIS? It should already be loaded
	cli.LwComponents, err = lwcomponent.LocalState()
	if err != nil {
		err = errors.Wrap(err, "unable to load components")
		return
	}

	component := cli.LwComponents.GetComponent(args[0])
	if component == nil {
		err = errors.New("component not found. Try running 'lacework component list'")
		return
	}

	if component.Status() != lwcomponent.Installed {
		err = errors.Errorf(
			"component not installed. Try running 'lacework component install %s'",
			args[0],
		)
		return
	}

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
	cli.OutputHuman("The component %s was deleted.\n", args[0])
	return
}
