//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2024, Lacework Inc.
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
	"strings"

	"github.com/spf13/cobra"
)

var (
	// grpcCmd is a hidden command that developers use to debug CDK components
	grpcCmd = &cobra.Command{
		Use:    "grpc",
		Hidden: true,
		Short:  "Starts a CDK gRPC server (developer mode)",
		Args:   cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cli.OutputHuman("\nDevelopment mode for CDK components")
			cli.OutputHuman("\n===================================\n\n")
			cli.OutputHuman("When debugging a component, it might expect some environment variables and a\n")
			cli.OutputHuman("running gRPC server, this command starts the CDK server and shows the variables\n")
			cli.OutputHuman("that your component might need:\n\n")
			vars := cli.envs()
			cli.OutputHuman("export %s", strings.Join(vars, " \\\n  "))
			cli.OutputHuman("\n\n'Ctrl+c' to stop the server. ")
			return cli.Serve()
		},
	}
)

func init() {
	rootCmd.AddCommand(grpcCmd)
}
