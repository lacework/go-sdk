//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/internal/cache"
	"github.com/lacework/go-sdk/internal/file"
)

var (
	// componentsCmd represents the components command
	componentsCmd = &cobra.Command{
		Use:     "component",
		Aliases: []string{"components"},
		Short:   "manage components",
		Long:    `Manage components to extend your experience with the Lacework platform`,
	}

	// componentsListCmd represents the azure sub-command inside the components command
	componentsListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list all components",
		Long:    `List all components`,
		RunE:    runComponentsList,
	}

	// componentsInstallCmd represents the gcp sub-command inside the components command
	componentsInstallCmd = &cobra.Command{
		Use:   "install",
		Short: "install a new component",
		Long:  `Install a new component`,
		Args:  cobra.ExactArgs(1),
		RunE:  runComponentsInstall,
	}

	// componentsUpdateCmd represents the aws sub-command inside the components command
	componentsUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "update an existing component",
		Long:  `Update an existing component`,
		Args:  cobra.ExactArgs(1),
		RunE:  runComponentsUpdate,
	}

	// componentsDeleteCmd represents the aws sub-command inside the components command
	componentsDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete an existing component",
		Long:  `Delete an existing component`,
		Args:  cobra.ExactArgs(1),
		RunE:  runComponentsDelete,
	}
)

func init() {
	// add the components command
	rootCmd.AddCommand(componentsCmd)

	// add sub-commands to the components command
	componentsCmd.AddCommand(componentsListCmd)
	componentsCmd.AddCommand(componentsInstallCmd)
	componentsCmd.AddCommand(componentsUpdateCmd)
	componentsCmd.AddCommand(componentsDeleteCmd)

	// load components
	cli.LoadComponents()
}

// @afiune how do we pass arguments?
func runComponent(cmd string, args []string) error {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	return c.Run()
}

func (c *cliState) LoadComponents() {
	c.LwComponents = loadComponentsFromDirk()
	for _, component := range c.LwComponents.Components {
		if component.Status == "Installed" && component.CLICommand {
			var (
				cmd     = component.Name
				cmdName = component.CommandName
			)
			cli.Log.Debugw("loading cli command", "component", cmd, "command_name", cmdName)
			rootCmd.AddCommand(
				&cobra.Command{
					Use:   cmdName,
					Short: fmt.Sprintf("%s component", cmd),
					Run: func(_ *cobra.Command, args []string) {
						runComponent(cmd, args)
					},
				},
			)
		}
	}
}

type LwComponentState struct {
	Version    string      `json:"version"`
	Components []Component `json:"components"`
}

type Component struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Status      string `json:"status"`
	//Size ?

	// will this component be accessible via the CLI
	CLICommand  bool   `json:"cli_command"`
	CommandName string `json:"command_name"`

	// the component is a binary
	Binary bool `json:"binary"`

	// the component is a library, only provides content for the CLI or other components
	Library bool `json:"library"`
}

func loadComponentsFromDirk() *LwComponentState {
	state := new(LwComponentState)
	// @afiune log more information about loading components
	cli.Log.Debugw("loading components")
	cacheDir, err := cache.CacheDir()
	if err != nil {
		return state
	}

	componentsFile := path.Join(cacheDir, "components")
	if file.FileExists(componentsFile) {
		componentState, err := ioutil.ReadFile(componentsFile)
		if err != nil {
			return state
		}

		err = json.Unmarshal(componentState, state)
		if err != nil {
			cli.Log.Debugw("unable to load components",
				"file", componentsFile,
				"error", err,
			)
		}
	}

	return state
}

func runComponentsList(_ *cobra.Command, _ []string) error {
	cli.OutputHuman(
		renderCustomTable(
			[]string{"Status", "Name", "Description"},
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
	return nil
}

func componentsToTable() [][]string {
	out := [][]string{}
	for _, cdata := range cli.LwComponents.Components {
		out = append(out, []string{
			cdata.Status,
			cdata.Name,
			cdata.Description,
		})
	}
	return out
}

func runComponentsInstall(_ *cobra.Command, args []string) error {
	cacheDir, err := cache.CacheDir()
	if err != nil {
		return err
	}

	componentsFile := path.Join(cacheDir, "components")

	exists := false
	for i, component := range cli.LwComponents.Components {
		if component.Name == args[0] {
			cli.LwComponents.Components[i].Status = "Installed"
			exists = true
		}
	}

	if !exists {
		return errors.New("component not found. Try running 'lacework component list'")
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(cli.LwComponents); err != nil {
		return err
	}

	if err := ioutil.WriteFile(componentsFile, buf.Bytes(), 0644); err != nil {
		return err
	}

	cli.StartProgress(" Installing component...")
	time.Sleep(5 * time.Second)
	cli.StopProgress()

	cli.OutputHuman("The component %s was installed.\n", args[0])
	return nil
}
func runComponentsUpdate(_ *cobra.Command, _ []string) error {
	return nil
}
func runComponentsDelete(_ *cobra.Command, args []string) error {
	cacheDir, err := cache.CacheDir()
	if err != nil {
		return err
	}

	componentsFile := path.Join(cacheDir, "components")

	exists := false
	for i, component := range cli.LwComponents.Components {
		if component.Name == args[0] {
			cli.LwComponents.Components[i].Status = "Not Installed"
			exists = true
		}
	}

	if !exists {
		return errors.New("component not found. Try running 'lacework component list'")
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(cli.LwComponents); err != nil {
		return err
	}

	if err := ioutil.WriteFile(componentsFile, buf.Bytes(), 0644); err != nil {
		return err
	}

	cli.StartProgress(" Deleting component...")
	time.Sleep(5 * time.Second)
	cli.StopProgress()

	cli.OutputHuman("The component %s was deleted.\n", args[0])
	return nil
}
