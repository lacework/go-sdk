//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lacework/go-sdk/lwlogger"
)

var (
	// the global cli state with defaults
	cli = NewDefaultState()

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:           "lacework",
		Short:         "A tool to manage the Lacework cloud security platform.",
		SilenceErrors: true,
		Long: `
The Lacework Command Line Interface is a tool that helps you manage the
Lacework cloud security platform. Use it to manage compliance reports,
external integrations, vulnerability scans, and other operations.

Start by configuring the Lacework CLI with the command:

    $ lacework configure

This will prompt you for your Lacework account and a set of API access keys.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			switch cmd.Use {
			case "help [command]", "configure", "version":
				return nil
			default:
				return cli.NewClient()
			}
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// first, verify if the user provided a command to execute,
	// if no command was provided, only print out the usage message
	if noCommandProvided() {
		errcheckWARN(rootCmd.Help())
		os.Exit(127)
	}

	errcheckEXIT(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().Bool("debug", false,
		"turn on debug logging",
	)
	rootCmd.PersistentFlags().Bool("nocolor", false,
		"turn off colors",
	)
	rootCmd.PersistentFlags().Bool("noninteractive", false,
		"disable interactive progress bars (i.e. 'spinners')",
	)
	rootCmd.PersistentFlags().Bool("json", false,
		"switch commands output from human-readable to json format",
	)
	rootCmd.PersistentFlags().StringP("profile", "p", "",
		"switch between profiles configured at ~/.lacework.toml",
	)
	rootCmd.PersistentFlags().StringP("api_key", "k", "",
		"access key id",
	)
	rootCmd.PersistentFlags().StringP("api_secret", "s", "",
		"secret access key",
	)
	rootCmd.PersistentFlags().StringP("account", "a", "",
		"account subdomain of URL (i.e. <ACCOUNT>.lacework.net)",
	)

	errcheckWARN(viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")))
	errcheckWARN(viper.BindPFlag("nocolor", rootCmd.PersistentFlags().Lookup("nocolor")))
	errcheckWARN(viper.BindPFlag("noninteractive", rootCmd.PersistentFlags().Lookup("noninteractive")))
	errcheckWARN(viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json")))
	errcheckWARN(viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile")))
	errcheckWARN(viper.BindPFlag("account", rootCmd.PersistentFlags().Lookup("account")))
	errcheckWARN(viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api_key")))
	errcheckWARN(viper.BindPFlag("api_secret", rootCmd.PersistentFlags().Lookup("api_secret")))
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	// Find home directory
	home, err := homedir.Dir()
	errcheckEXIT(err)

	// Search config in home directory with name ".lacework" (without extension)
	viper.AddConfigPath(home)
	viper.SetConfigName(".lacework")

	viper.SetConfigType("toml") // set TOML as the config format
	viper.SetEnvPrefix("LW")    // set prefix for all env variables LW_ABC
	viper.AutomaticEnv()        // read in environment variables that match

	if viper.GetBool("debug") {
		cli.LogLevel = "DEBUG"
	}

	// initialize a Lacework logger
	cli.Log = lwlogger.New(cli.LogLevel).Sugar()

	if viper.GetBool("nocolor") {
		cli.Log.Info("turning off colors")
		cli.JsonF.DisabledColor = true
	}

	if viper.GetBool("noninteractive") {
		cli.NonInteractive()
	}

	if viper.GetBool("json") {
		cli.EnableJSONOutput()
	}

	// by default the cli logs are going to be visualized in
	// a console format unless the user wants the opposite
	if os.Getenv("LW_LOG_FORMAT") == "" {
		if cli.JSONOutput() {
			os.Setenv("LW_LOG_FORMAT", "JSON")
		} else {
			os.Setenv("LW_LOG_FORMAT", "CONSOLE")
		}
	}

	// try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// the config file was not found; ignore error
			cli.Log.Debugw("configuration file not found")
		} else {
			// the config file was found but another error was produced
			exitwith(errors.Wrap(err, "Error: unable to read in config"))
		}
	} else {
		cli.Log.Debugw("using configuration file",
			"path", viper.ConfigFileUsed(),
		)
	}

	// get the profile passed as a parameter or environment variable
	// if any, set it into the CLI state, that will trigger to load the
	// state, if no profile was specified just load the default state
	if p := viper.GetString("profile"); len(p) != 0 {
		err = cli.SetProfile(p)
	} else {
		err = cli.LoadState()
	}

	if err != nil {
		if isCommand("configure") {
			cli.Log.Debugw(
				"error ignored",
				"reason", "running configure cmd",
				"error", err,
			)
		} else {
			exitwith(err)
		}
	}
}

// isCommand checks the overall arguments passed to the lacework cli
// and returns true if the provided command name is the one running
func isCommand(cmd string) bool {
	if len(os.Args) <= 1 {
		return false
	}

	if os.Args[1] == cmd {
		return true
	}

	return false
}

// noCommandProvided checks if a command or argument was provided
func noCommandProvided() bool {
	return len(os.Args) <= 1
}

// errcheckEXIT is a simple macro to check Golang errors, if the provided
// error is nil, it doesn't do anything, but if the error has something,
// it exits the program
func errcheckEXIT(err error) {
	if err != nil {
		exitwith(err)
	}
}

// errcheckWARN is similar to errcheckEXIT but it doesn't exit the program,
// it only prints a WARNING message to the user, useful for those cases where
// we know there won't be aproblem but the linter still asks to check all errors
func errcheckWARN(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN %s\n", err)
	}
}

// exitwith prints out an error message and exits the program with exit code 1
func exitwith(err error) {
	exitwithCode(err, 1)
}

// exitwithCode prints out an error message and exits the program with
// the provided exit code
func exitwithCode(err error, code int) {
	fmt.Fprintf(os.Stderr, "ERROR %s\n", err)
	os.Exit(code)
}
