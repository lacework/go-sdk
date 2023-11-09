//
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

package cdk

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "[[.Component]]",
	Short: "A brief description of your component",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

[[.Component]] is a scaffolding for Lacework CDK components. It helps developers
create new components fast by autogenerating a skeleton with boilerplate code.

To quickly get going with this new component, start modifying the placeholder command:

    lacework [[.Component]] placeholder
`,

	// NOTE: If your component does NOT need multiple commands, uncomment
	//       this line and remove the file 'placeholder.go'
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().Bool("debug", false,
		"turn on debug logging",
	)
	RootCmd.PersistentFlags().Bool("nocolor", false,
		"turn off colors",
	)
	RootCmd.PersistentFlags().Bool("nocache", false,
		"turn off caching",
	)
	RootCmd.PersistentFlags().Bool("noninteractive", false,
		"turn off interactive mode (disable spinners, prompts, etc.)",
	)
	RootCmd.PersistentFlags().Bool("json", false,
		"switch commands output from human-readable to json format",
	)

	errcheckWARN(viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug")))
	errcheckWARN(viper.BindPFlag("nocolor", RootCmd.PersistentFlags().Lookup("nocolor")))
	errcheckWARN(viper.BindPFlag("nocache", RootCmd.PersistentFlags().Lookup("nocache")))
	errcheckWARN(viper.BindPFlag("noninteractive", RootCmd.PersistentFlags().Lookup("noninteractive")))
	errcheckWARN(viper.BindPFlag("json", RootCmd.PersistentFlags().Lookup("json")))

	// Here you will define your flags and configuration settings.
}

// errcheckWARN is a simple macro to check Golang errors, if the provided error
// is nil, it doesn't do anything, but if the error has something, it prints a
// WARNING message to the user, useful for those cases where we know there won't
// be a problem but the linter still asks to check all errors
func errcheckWARN(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN %s\n", err)
	}
}
