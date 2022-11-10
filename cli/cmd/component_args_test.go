package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentArgs(t *testing.T) {
	assert := assert.New(t)
	flags := rootCmd.PersistentFlags()
	for _, k := range []struct{ args, expectedComponent, expectedCLI []string }{
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
