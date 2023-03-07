//
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
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestComponentArgs(t *testing.T) {
	assert := assert.New(t)
	flags := &pflag.FlagSet{}
	rootCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		flags.AddFlag(f)
	})
	// cobra adds a -v flag very late in the execution flow
	flags.BoolP("version", "v", false, "version of command")
	for _, k := range []struct{ args, expectedComponent, expectedCLI []string }{
		{
			[]string{"iac", "-v", "-d", "foo"},
			[]string{"iac", "-d", "foo"},
			[]string{"-v"},
		},
		{
			[]string{"--profile", "none", "--debug"},
			[]string{},
			[]string{"--profile", "none", "--debug"}},
		{
			[]string{
				"iac", "profile", "--iac-organization", "1234", "--profile", "none", "--json",
				"--upload=false",
			},
			[]string{"iac", "profile", "--iac-organization", "1234", "--upload=false"},
			[]string{"--profile", "none", "--json"},
		},
		{
			[]string{
				"iac", "tf-scan", "-a", "qan", "--profile", "foo", "-d", "src", "--var-file", "src/prod.tfvars,src/base.tfvars",
			},
			[]string{"iac", "tf-scan", "-d", "src", "--var-file", "src/prod.tfvars,src/base.tfvars"},
			[]string{"-a", "qan", "--profile", "foo"},
		},
		{
			[]string{
				"iac", "org", "--help",
			},
			[]string{"iac", "org", "--help"},
			[]string{},
		},
		{
			[]string{
				"iac", "configure", "set-profile", "--debug", "--help",
			},
			[]string{
				"iac", "configure", "set-profile", "--help",
			},
			[]string{
				"--debug",
			},
		},
	} {
		p := componentArgParser{}
		p.parseArgs(flags, k.args)
		assert.Equal(k.expectedComponent, p.componentArgs, "parsing %v expecting component args", k.args)
		assert.Equal(k.expectedCLI, p.cliArgs, k.args, "parsing %v expecting cli args", k.args)
	}
}
