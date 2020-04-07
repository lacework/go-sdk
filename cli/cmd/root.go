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

	prettyjson "github.com/hokaccha/go-prettyjson"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/lacework/go-sdk/lwlogger"
)

type cliState struct {
	Account  string
	KeyID    string
	Secret   string
	Token    string
	LogLevel string

	JsonF *prettyjson.Formatter
	Log   *zap.SugaredLogger
}

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile string
	cli     = cliState{
		JsonF: prettyjson.NewFormatter(),
	}
	rootCmd = &cobra.Command{
		Use:              "lacework",
		Short:            "A tool to manage the Lacework cloud security platform.",
		PersistentPreRun: loadStateFromViper,
		SilenceErrors:    true,
		Long: `
The Lacework Command Line Interface is a tool that helps you manage the
Lacework cloud security platform. You can use it to manage compliance
reports, external integrations, vulnerability scans, and other operations.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile,
		"config", "c", "",
		"config file (default is $HOME/.lacework.toml)",
	)
	rootCmd.PersistentFlags().Bool("debug", false,
		"turn on debug logging",
	)
	rootCmd.PersistentFlags().StringP("api_key", "k", "",
		"access key id",
	)
	rootCmd.PersistentFlags().StringP("api_secret", "s", "",
		"access secret key",
	)
	rootCmd.PersistentFlags().StringP("account", "a", "",
		"account subdomain of URL (i.e. <ACCOUNT>.lacework.net)",
	)

	checkBindError(viper.BindPFlag("account", rootCmd.PersistentFlags().Lookup("account")))
	checkBindError(viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api_key")))
	checkBindError(viper.BindPFlag("api_secret", rootCmd.PersistentFlags().Lookup("api_secret")))
	checkBindError(viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")))
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if cfgFile != "" {
		// Use config file from flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".lacework" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(".lacework")
	}

	viper.SetConfigType("toml") // set TOML as the config format
	viper.SetEnvPrefix("LW")    // set prefix for all env variables LW_ABC
	viper.AutomaticEnv()        // read in environment variables that match

	if viper.GetBool("debug") {
		cli.LogLevel = "DEBUG"
	}

	// by default the cli logs are going to be visualized in
	// a console format unless the user wants the opposite
	if os.Getenv("LW_LOG_FORMAT") == "" {
		os.Setenv("LW_LOG_FORMAT", "CONSOLE")
	}

	// initialize a Lacework logger
	cli.Log = lwlogger.New(cli.LogLevel).Sugar()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		cli.Log.Debugw("using config file",
			"path", viper.ConfigFileUsed(),
		)
	}
}

func checkBindError(err error) {
	if err != nil {
		// this check happens before we have initialized the logger,
		// so we need to use native fmt prints
		fmt.Printf("WARN unable to bind parameter: %v\n", err)
	}
}

func loadStateFromViper(_ *cobra.Command, _ []string) {
	cli.KeyID = viper.GetString("api_key")
	cli.Secret = viper.GetString("api_secret")
	cli.Account = viper.GetString("account")

	cli.Log.Debugw("state loaded",
		"account", cli.Account,
		"api_key", cli.KeyID,
		"api_secret", cli.Secret,
	)
}
