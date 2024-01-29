//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2023, Lacework Inc.
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
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/abiosoft/colima/util/terminal"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/internal/databox"
	"github.com/lacework/go-sdk/lwcomponent"
)

var cdkDevState = struct {
	Type        string
	Scaffolding string
	Description string
}{}

var cdkGolangScaffoldingRequirements = map[string]string{
	"go": "https://go.dev/dl/",
}

var cdkPythonScaffoldingRequirements = map[string]string{
	"python3": "https://www.python.org/downloads/",
	"poetry":  "https://python-poetry.org/docs/",
}

func init() {
	componentsDevModeCmd.Flags().StringVar(
		&cdkDevState.Type,
		"type", "",
		fmt.Sprintf("component type (%s, %s, %s)",
			lwcomponent.BinaryType,
			lwcomponent.CommandType,
			lwcomponent.LibraryType,
		),
	)

	componentsDevModeCmd.Flags().StringVar(
		&cdkDevState.Description,
		"description", "",
		"component description",
	)

	componentsDevModeCmd.Flags().StringVar(
		&cdkDevState.Scaffolding, "scaffolding", "",
		"autogenerate code for a new component (available: Golang, Python)",
	)
}

func runComponentsDevMode(_ *cobra.Command, args []string) error {
	if !lwcomponent.CatalogV1Enabled(cli.LwApi) {
		return prototypeRunComponentsDevMode(args)
	}

	return devModeComponent(args)
}

func devModeComponent(args []string) error {
	var componentName string = args[0]

	cli.StartProgress("Loading components Catalog...")

	catalog, err := lwcomponent.NewCatalog(cli.LwApi, lwcomponent.NewStageTarGz, false)
	defer catalog.Persist()

	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to load component Catalog")
	}

	component, _ := catalog.GetComponent(componentName)
	if component == nil {
		newComponent := lwcomponent.NewCDKComponent(componentName, "", lwcomponent.EmptyType, nil, nil)
		component = &newComponent

		cli.OutputHuman("Component '%s' not found. Defining a new component.\n",
			color.HiYellowString(component.Name))

		var (
			helpMsg = fmt.Sprintf("What are these component types ?\n"+
				"\n'%s' - A binary accessible via the Lacework CLI (Users will run 'lacework <COMPONENT_NAME>')"+
				"\n'%s' - A regular standalone-binary (this component type is not accessible via the CLI)"+
				"\n'%s' - A library that only provides content for the Lacework CLI or other components\n",
				lwcomponent.CommandType, lwcomponent.BinaryType, lwcomponent.LibraryType)
		)
		if cdkDevState.Type == "" {
			if err := survey.AskOne(&survey.Select{
				Message: "Select the type of component you are developing:",
				Help:    helpMsg,
				Options: []string{
					lwcomponent.CommandType,
					lwcomponent.BinaryType,
					lwcomponent.LibraryType,
				},
			}, &cdkDevState.Type); err != nil {
				return err
			}
		}

		component.Type = lwcomponent.Type(cdkDevState.Type)

		if cdkDevState.Description == "" {
			if err := survey.AskOne(&survey.Input{
				Message: "What is this component about? (component description):",
			}, &component.Description); err != nil {
				return err
			}
		} else {
			component.Description = cdkDevState.Description
		}
	}

	if err := component.EnterDevMode(); err != nil {
		return errors.Wrap(err, "unable to enter development mode")
	}

	componentDir, err := component.Dir()
	if err != nil {
		return errors.New("unable to detect RootPath")
	}

	cli.OutputHuman("Component '%s' in now in development mode.\n",
		color.HiYellowString(component.Name))

	if component.Type == lwcomponent.CommandType {
		// Offer the creation of a component scaffolding
		if cdkDevState.Scaffolding == "" && cli.InteractiveMode() {
			if err := survey.AskOne(&survey.Select{
				Message: "Would you like to initialize your component with scaffolding? ",
				Options: []string{"No. Start from scratch", "Golang", "Python"},
			}, &cdkDevState.Scaffolding); err != nil {
				return err
			}
		}

		switch cdkDevState.Scaffolding {
		case "Golang":
			if err := cdkGolangScaffolding(component.Name, componentDir); err != nil {
				return err
			}

		case "Python":
			if err := cdkPythonScaffolding(component.Name, componentDir); err != nil {
				return err
			}

		default:
			cli.OutputHuman("\nDeploy your dev component at: %s\n",
				color.HiYellowString(filepath.Join(componentDir, component.Name)))
		}
	}

	cli.OutputHuman("\nComponent directory: %s\n", color.HiCyanString(componentDir))
	cli.OutputHuman("Dev specs: %s\n", color.HiCyanString(filepath.Join(componentDir, ".dev")))
	return nil
}

func prototypeRunComponentsDevMode(args []string) error {
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
			helpMsg = fmt.Sprintf("What are these component types ?\n"+
				"\n'%s' - A binary accessible via the Lacework CLI (Users will run 'lacework <COMPONENT_NAME>')"+
				"\n'%s' - A regular standalone-binary (this component type is not accessible via the CLI)"+
				"\n'%s' - A library that only provides content for the Lacework CLI or other components\n",
				lwcomponent.CommandType, lwcomponent.BinaryType, lwcomponent.LibraryType)
		)
		if cdkDevState.Type == "" {
			if err := survey.AskOne(&survey.Select{
				Message: "Select the type of component you are developing:",
				Help:    helpMsg,
				Options: []string{
					lwcomponent.CommandType,
					lwcomponent.BinaryType,
					lwcomponent.LibraryType,
				},
			}, &cdkDevState.Type); err != nil {
				return err
			}
		}

		component.Type = lwcomponent.Type(cdkDevState.Type)

		if cdkDevState.Description == "" {
			if err := survey.AskOne(&survey.Input{
				Message: "What is this component about? (component description):",
			}, &component.Description); err != nil {
				return err
			}
		} else {
			component.Description = cdkDevState.Description
		}
	}

	if err := component.EnterDevelopmentMode(); err != nil {
		return errors.Wrap(err, "unable to enter development mode")
	}

	rPath, err := component.RootPath()
	if err != nil {
		return errors.New("unable to detect RootPath")
	}

	cli.OutputHuman("Component '%s' in now in development mode.\n",
		color.HiYellowString(component.Name))

	if component.Type == lwcomponent.CommandType {
		// Offer the creation of a component scaffolding
		if cdkDevState.Scaffolding == "" && cli.InteractiveMode() {
			if err := survey.AskOne(&survey.Select{
				Message: "Would you like to initialize your component with scaffolding? ",
				Options: []string{"No. Start from scratch", "Golang", "Python"},
			}, &cdkDevState.Scaffolding); err != nil {
				return err
			}
		}

		switch cdkDevState.Scaffolding {
		case "Golang":
			if err := cdkGolangScaffolding(component.Name, rPath); err != nil {
				return err
			}

		case "Python":
			if err := cdkPythonScaffolding(component.Name, rPath); err != nil {
				return err
			}

		default:
			cli.OutputHuman("\nDeploy your dev component at: %s\n",
				color.HiYellowString(filepath.Join(rPath, component.Name)))
		}
	}

	cli.OutputHuman("\nRoot path: %s\n", color.HiCyanString(rPath))
	cli.OutputHuman("Dev specs: %s\n", color.HiCyanString(filepath.Join(rPath, ".dev")))
	return nil
}

func cdkGolangScaffolding(componentName string, componentDir string) error {
	if err := cdkScaffoldingPreflightCheck("Golang", cdkGolangScaffoldingRequirements); err != nil {
		return err
	}

	cli.OutputHuman("\nDeploying %s scaffolding:\n", color.HiMagentaString("Golang"))

	for _, file := range databox.ListFilesFromDir("/scaffoldings/golang") {
		content, found := databox.Get(file)
		if found {
			// Create directory, if needed
			subDir := filepath.Dir(file)
			subDir = strings.TrimPrefix(subDir, "/scaffoldings/golang")
			fileDir := filepath.Join(componentDir, subDir)
			if subDir != "" {
				if err := os.MkdirAll(fileDir, 0755); err != nil {
					return errors.Wrap(err, "unable to create subdirectory from scaffolding")
				}
			}

			var (
				buff     = &bytes.Buffer{}
				fileName = filepath.Base(file)
				filePath = filepath.Join(fileDir, fileName)
				tmpl     = template.Must(template.New(fileName).Delims("[[", "]]").Parse(string(content)))
				cData    = struct{ Component string }{
					Component: componentName,
				}
			)
			if err := tmpl.Execute(buff, cData); err != nil {
				return errors.Wrap(err, "unable to generate files from go scaffolding")
			}
			if err := os.WriteFile(filePath, buff.Bytes(), os.ModePerm); err != nil {
				cli.OutputChecklist(failureIcon, "Unable to write file %s\n", color.HiRedString(filePath))
				cli.Log.Debugw("unable to write file", "error", err)
			} else {
				cli.OutputChecklist(successIcon, "File %s deployed\n", color.HiYellowString(filePath))
			}
		}
	}

	// Missing tasks we can do for the developer
	//
	// 1) Change directory to Root path
	//    > Command: 'cd ...'
	// 2) Initialize git repository
	//    > Command: 'git init'
	// 3) Create your initial commit
	//    > Command: 'git add .; git commit -m "feat: init component"'
	// 4) Dowload Go dependencies
	//    > Command: 'make go-vendor'
	// 5) Build the component
	//    > Command: 'make build'
	// 6) Run the component via the Lacework CLI
	//    > Command: 'lacework <component_name> placeholder'
	//
	cli.StartProgress("Initializing Git repository...")
	err := cdkInitGitRepo(componentDir)
	cli.StopProgress()
	if err != nil {
		cli.OutputChecklist(failureIcon, "Unable to initialize Git repository\n")
		cli.Log.Debugw("unable to initialize Git repository", "error", err)
	} else {
		cli.OutputChecklist(successIcon, "Git repository initialized\n")
	}

	cli.StartProgress("Downloading Go dependencies...")
	err = cdkGoVendor(componentDir)
	cli.StopProgress()
	if err != nil {
		cli.OutputChecklist(failureIcon, "Unable to download Go dependencies\n")
		cli.Log.Debugw("unable to download Go dependencies", "error", err)
	} else {
		cli.OutputChecklist(successIcon, "Go dependencies downloaded\n")
	}

	cli.StartProgress("Building your component...")
	err = cdkGoBuild(componentDir)
	cli.StopProgress()
	if err != nil {
		cli.OutputChecklist(failureIcon, "Unable to build your Go component\n")
		cli.Log.Debugw("unable to build your Go component", "error", err)
	} else {
		cli.OutputChecklist(successIcon, "Dev component built at %s\n",
			color.HiYellowString(filepath.Join(componentDir, componentName)))
	}

	cli.StartProgress("Verifying component...")
	err = cdkGoRunVerify(componentName)
	cli.StopProgress()
	if err != nil {
		// this is not on the developer, it's on this codebase, notify to fix it
		cli.OutputChecklist(failureIcon, "Unable run scaffolding component\n")
		cli.Log.Debugw("unable to run scaffolding component", "error", err)
	} else {
		cli.OutputChecklist(successIcon, "Component verified\n")
	}

	cli.OutputHuman("\nDeployment completed! Time for %s\n", randomEmoji())
	return nil
}

func cdkPythonScaffolding(componentName string, componentDir string) error {
	if err := cdkScaffoldingPreflightCheck("Python", cdkPythonScaffoldingRequirements); err != nil {
		return err
	}

	cli.OutputHuman("\nDeploying %s scaffolding:\n", color.HiMagentaString("Python"))

	for _, file := range databox.ListFilesFromDir("/scaffoldings/python") {
		content, found := databox.Get(file)
		if found {
			// Create directory, if needed
			subDir := filepath.Dir(file)
			subDir = strings.TrimPrefix(subDir, "/scaffoldings/python")
			fileDir := filepath.Join(componentDir, subDir)
			if subDir != "" {
				if err := os.MkdirAll(fileDir, 0755); err != nil {
					return errors.Wrap(err, "unable to create subdirectory from scaffolding")
				}
			}

			var (
				fileName = filepath.Base(file)
				filePath = filepath.Join(fileDir, fileName)
			)
			if err := os.WriteFile(filePath, content, os.ModePerm); err != nil {
				cli.OutputChecklist(failureIcon, "Unable to write file %s\n", color.HiRedString(filePath))
				cli.Log.Debugw("unable to write file", "error", err)
			} else {
				cli.OutputChecklist(successIcon, "File %s deployed\n", color.HiYellowString(filePath))
			}
		}
	}

	cli.StartProgress("Initializing Git repository...")
	err := cdkInitGitRepo(componentDir)
	cli.StopProgress()
	if err != nil {
		cli.OutputChecklist(failureIcon, "Unable to initialize Git repository\n")
		cli.Log.Debugw("unable to initialize Git repository", "error", err)
	} else {
		cli.OutputChecklist(successIcon, "Git repository initialized\n")
	}

	// Poetry repository structure `project-name/src/project-name/__init__.py`
	err = os.Rename(filepath.Join(componentDir, "src/package"), filepath.Join(componentDir, "src", componentName))
	if err != nil {
		cli.OutputChecklist(failureIcon, "Unable to rename package directory\n")
		cli.Log.Debugw("unable to rename package directory", "error", err)
		return err
	}

	cli.StartProgress("Poetry init...")
	err = cdkExec(componentDir,
		"poetry",
		"init",
		"--no-interaction",
		// Because of https://github.com/python-poetry/poetry/issues/5975
		fmt.Sprintf("--name=%s", strings.ReplaceAll(componentName, "-", "")),
		"--python=^3.11,<3.12",
		"--dev-dependency=pyinstaller",
		"--dev-dependency=poethepoet")
	cli.StopProgress()
	if err != nil {
		cli.OutputChecklist(failureIcon, "Unable to initialize Poetry\n")
		cli.Log.Debugw("unable to initialize Poetry", "error", err)
		return err
	} else {
		cli.OutputChecklist(successIcon, "Poetry init\n")
	}

	f, err := os.OpenFile(filepath.Join(componentDir, "pyproject.toml"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.WriteString("[tool.poe.tasks]\n")
	if err != nil {
		return err
	}

	_, err = f.WriteString("build.shell = \"poetry run pyinstaller src/")
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf(
		"%s/__main__.py --collect-submodules application -D --name %s --distpath dist;",
		componentName, componentName,
	))
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf(
		" mv dist/%s/* .\"\n",
		componentName,
	))
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf(
		"clean = \"rm -r build/ %s %s.spec\"\n",
		componentName, componentName,
	))
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	cli.StartProgress("Poetry install...")
	err = cdkExec(componentDir, "poetry", "install")
	cli.StopProgress()
	if err != nil {
		cli.OutputChecklist(failureIcon, "Unable to Poetry install\n")
		cli.Log.Debugw("unable to Poetry install", "error", err)
		return err
	} else {
		cli.OutputChecklist(successIcon, "Poetry install\n")
	}

	cli.StartProgress("Building your component...")
	err = cdkExec(componentDir, "poetry", "run", "poe", "build")
	cli.StopProgress()
	if err != nil {
		cli.OutputChecklist(failureIcon, "Unable to build your Python component\n")
		cli.Log.Debugw("unable to build your Python component", "error", err)
	} else {
		cli.OutputChecklist(successIcon, "Dev component built at %s\n",
			color.HiYellowString(filepath.Join(componentDir, componentName)))
	}

	cli.OutputHuman("\nDeployment completed! Time for %s\n", randomEmoji())
	return nil
}

func cdkScaffoldingPreflightCheck(scaffolding string, requirements map[string]string) error {
	errMessage := ""
	for file, site := range requirements {
		if _, err := exec.LookPath(file); err != nil {
			errMessage += fmt.Sprintf(`
%s is required to create the %s scaffolding. Please install it before proceeding:
  %s: %s`, file, scaffolding, file, site)
		}
	}
	if errMessage != "" {
		return errors.New(errMessage)
	}
	return nil
}

func cdkExec(rootPath string, name string, args ...string) error {
	var (
		vw  = terminal.NewVerboseWriter(10)
		cmd = exec.Command(name, args...)
	)
	if _, err := vw.Write([]byte(fmt.Sprintf("Command: %s %v\n", name, args))); err != nil {
		cli.Log.Debugw("unable to write to virtual terminal", "error", err)
	}
	cmd.Env = os.Environ()
	cmd.Dir = rootPath
	cmd.Stdout = vw
	cmd.Stderr = vw

	defer func() {
		if _, err := vw.Write([]byte("\n")); err != nil {
			cli.Log.Debugw("unable to write to virtual terminal", "error", err)
		}
	}()

	return cmd.Run()
}

func cdkInitGitRepo(rootPath string) error {
	eMsg := "unable to initialize Git repo"

	repo, err := git.PlainInit(rootPath, false)
	if err != nil {
		return errors.Wrap(err, eMsg)
	}

	w, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, eMsg)
	}

	_, err = w.Add(".")
	if err != nil {
		return errors.Wrap(err, eMsg)
	}

	_, err = w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Component Scaffolding",
			Email: "support@lacework.net",
			When:  time.Now(),
		},
	})

	return err
}

func cdkGoVendor(rootPath string) error {
	var (
		vw  = terminal.NewVerboseWriter(10)
		cmd = exec.Command("make", "go-vendor")
	)
	if _, err := vw.Write([]byte("Command: make go-vendor\n")); err != nil {
		cli.Log.Debugw("unable to write to virtual terminal", "error", err)
	}
	cmd.Env = os.Environ()
	cmd.Dir = rootPath
	cmd.Stdout = vw
	cmd.Stderr = vw

	// @afiune silly workaround to clean the spinner output
	defer func() {
		if _, err := vw.Write([]byte("\n")); err != nil {
			cli.Log.Debugw("unable to write to virtual terminal", "error", err)
		}
	}()
	return cmd.Run()
}

func cdkGoBuild(rootPath string) error {
	var (
		vw  = terminal.NewVerboseWriter(10)
		cmd = exec.Command("make", "build")
	)
	if _, err := vw.Write([]byte("Command: make build\n")); err != nil {
		cli.Log.Debugw("unable to write to virtual terminal", "error", err)
	}
	cmd.Env = os.Environ()
	cmd.Dir = rootPath
	cmd.Stdout = vw
	cmd.Stderr = vw
	return cmd.Run()
}

func cdkGoRunVerify(componentName string) error {
	var (
		vw  = terminal.NewVerboseWriter(10)
		cmd = exec.Command(laceworkCLIBinary(), componentName, "placeholder")
	)
	_, err := vw.Write([]byte(fmt.Sprintf("Command: %s\n", strings.Join(cmd.Args, " "))))
	if err != nil {
		cli.Log.Debugw("unable to write to virtual terminal", "error", err)
	}
	cmd.Env = os.Environ()
	cmd.Stdout = vw
	cmd.Stderr = vw
	return cmd.Run()
}

func laceworkCLIBinary() string {
	if os.Getenv("LW_CLI_INTEGRATION_MODE") != "" {
		return fmt.Sprintf(
			"lacework-cli-%s-%s",
			runtime.GOOS, runtime.GOARCH,
		)
	}

	return "lacework"
}
