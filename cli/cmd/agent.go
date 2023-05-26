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
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	agentCmdState = struct {
		TokenUpdateEnable     bool
		TokenUpdateDisable    bool
		TokenUpdateName       string
		TokenUpdateDesc       string
		InstallForce          bool
		InstallSshUser        string
		InstallSshPort        int
		InstallAgentToken     string
		InstallTrustHostKey   bool
		InstallPassword       string
		InstallIdentityFile   string
		InstallTagKey         string
		InstallTag            []string
		InstallIncludeRegions []string
		InstallDryRun         bool
		InstallProjectId      string
		InstallMaxParallelism int
		InstallBYORole        string
		InstallSkipCreatInfra bool
		InstallForceReinstall bool
		InstallServerURL      string
		InstallAWSProfile     string
	}{}

	defaultSshIdentityKey = "~/.ssh/id_rsa"

	agentCmd = &cobra.Command{
		Use:   "agent",
		Short: "Manage Lacework agents",
		Long: `Manage agents and agent access tokens in your account.

To analyze application, host, and user behavior, Lacework uses a lightweight agent,
which securely forwards collected metadata to the Lacework cloud for analysis. The
agent requires minimal system resources and runs on most 64-bit Linux distributions.

For a complete list of supported operating systems, visit:

  https://docs.lacework.com/supported-operating-systems`,
	}

	agentTokenCmd = &cobra.Command{
		Use:     "token",
		Aliases: []string{"tokens"},
		Short:   "Manage agent access tokens",
		Long: `Manage agent access tokens in your account.

Agent tokens should be treated as secret and not published. A token uniquely identifies
a Lacework customer. If you suspect your token has been publicly exposed or compromised,
generate a new token, update the new token on all machines using the old token. When
complete, the old token can safely be disabled without interrupting Lacework services.`,
	}

	agentTokenListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all agent access tokens",
		Args:    cobra.NoArgs,
		RunE:    listAgentTokens,
	}

	agentTokenCreateCmd = &cobra.Command{
		Use:   "create <name> [description]",
		Short: "Create a new agent access token",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  createAgentToken,
	}

	agentTokenShowCmd = &cobra.Command{
		Use:   "show <token>",
		Short: "Show details about an agent access token",
		Args:  cobra.ExactArgs(1),
		RunE:  showAgentToken,
	}

	agentTokenUpdateCmd = &cobra.Command{
		Use:   "update <token>",
		Short: "Update an agent access token",
		Long: `Update an agent access token.

To update the token name and description:

    lacework agent token update <token> --name dev --description "k8s deployment for dev"

To disable a token:

    lacework agent token update <token> --disable

To enable a token:

    lacework agent token update <token> --enable`,
		Args: cobra.ExactArgs(1),
		RunE: updateAgentToken,
	}

	// TODO hidden for now
	agentGenerateCmd = &cobra.Command{
		Use:    "generate",
		Short:  "Generate agent deployment scripts",
		Long:   `TBA`,
		Hidden: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}

	agentInstallCmd = &cobra.Command{
		Use:   "install <[user@]host[:port]>",
		Short: "Install the datacollector agent on a remote host",
		Args:  cobra.ExactArgs(1),
		Long: `For single host installation of the Lacework agent via Secure Shell (SSH).

When this command is executed without any additional flag, an interactive prompt will be
launched to help gather the necessary authentication information to access the remote host.

To authenticate to the remote host with a username and password.

    lacework agent install <host> --ssh_username <your-user> --ssh_password <secret>

To authenticate to the remote host with an identity file instead.

    lacework agent install <user@host> -i /path/to/your/key

To provide an agent access token of your choice, use the command 'lacework agent token list',
select a token and pass it to the '--token' flag.

    lacework agent install <user@host> -i /path/to/your/key --token <token>

To authenticate to the remote host on a non-standard SSH port use the '--ssh_port' flag or
pass it directly via the argument.

    lacework agent install <user@host:port>

To explicitly specify the server URL that the agent will connect to:

    lacework agent install --server_url https://your.server.url.lacework.net

To list all active agents in your environment. 

    lacework agent list

NOTE: New agents could take up to an hour to report back to the platform.`,
		RunE: installRemoteAgent,
	}

	agentAWSInstallCmd = &cobra.Command{
		Use:   "aws-install",
		Args:  cobra.NoArgs,
		Short: "Install the datacollector agent on all remote AWS hosts",
	}

	agentGCPInstallCmd = &cobra.Command{
		Use:   "gcp-install",
		Args:  cobra.NoArgs,
		Short: "Install the datacollector agent on all remote GCE hosts",
	}
)

func init() {
	// add the agent command
	rootCmd.AddCommand(agentCmd)

	// add the token sub-command to the agent cmd
	agentCmd.AddCommand(agentTokenCmd)
	agentCmd.AddCommand(agentInstallCmd)
	agentCmd.AddCommand(agentGenerateCmd)
	agentCmd.AddCommand(agentListCmd)
	agentCmd.AddCommand(agentAWSInstallCmd)
	agentCmd.AddCommand(agentGCPInstallCmd)

	// add the list sub-command to the 'agent token' cmd
	agentTokenCmd.AddCommand(agentTokenListCmd)
	agentTokenCmd.AddCommand(agentTokenCreateCmd)
	agentTokenCmd.AddCommand(agentTokenShowCmd)
	agentTokenCmd.AddCommand(agentTokenUpdateCmd)

	// add sub-commands to the 'agent aws-install' command for different install methods
	agentAWSInstallCmd.AddCommand(agentInstallAWSEC2ICCmd)
	agentAWSInstallCmd.AddCommand(agentInstallAWSSSHCmd)
	agentAWSInstallCmd.AddCommand(agentInstallAWSSSMCmd)

	// add sub-commands to the 'agent gcp-install' command for different install methods
	agentGCPInstallCmd.AddCommand(agentInstallGCPOSLCmd)

	// 'agent token update' flags
	agentTokenUpdateCmd.Flags().BoolVar(&agentCmdState.TokenUpdateEnable,
		"enable", false, "enable agent access token",
	)
	agentTokenUpdateCmd.Flags().BoolVar(&agentCmdState.TokenUpdateDisable,
		"disable", false, "disable agent access token",
	)
	agentTokenUpdateCmd.Flags().StringVar(&agentCmdState.TokenUpdateName,
		"name", "", "new agent access token name",
	)
	agentTokenUpdateCmd.Flags().StringVar(&agentCmdState.TokenUpdateDesc,
		"description", "", "new agent access token description",
	)

	// 'agent install' flags
	agentInstallCmd.Flags().StringVarP(&agentCmdState.InstallIdentityFile,
		"identity_file", "i", defaultSshIdentityKey,
		"identity (private key) for public key authentication",
	)
	agentInstallCmd.Flags().StringVar(&agentCmdState.InstallPassword,
		"ssh_password", "", "password for authentication",
	)
	agentInstallCmd.Flags().StringVar(&agentCmdState.InstallSshUser,
		"ssh_username", "", "username to login with",
	)
	agentInstallCmd.Flags().IntVar(&agentCmdState.InstallSshPort,
		"ssh_port", 22, "port to connect to on the remote host",
	)
	agentInstallCmd.Flags().BoolVar(&agentCmdState.InstallForce,
		"force", false, "override any pre-installed agent",
	)
	agentInstallCmd.Flags().StringVar(&agentCmdState.InstallAgentToken,
		"token", "", "agent access token",
	)
	agentInstallCmd.Flags().BoolVar(&agentCmdState.InstallTrustHostKey,
		"trust_host_key", false, "automatically add host keys to the ~/.ssh/known_hosts file",
	)
	agentInstallCmd.Flags().StringVar(&agentCmdState.InstallServerURL,
		"server_url", "https://agent.lacework.net", "server URL that agents will talk to, prefixed with `https://`",
	)
}

func showAgentToken(_ *cobra.Command, args []string) error {
	response, err := cli.LwApi.V2.AgentAccessTokens.Get(args[0])
	if err != nil {
		return errors.Wrap(err, "unable to get agent access token")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(response.Data)
	}

	cli.OutputHuman(buildAgentTokenDetailsTable(response.Data))
	return nil
}

func updateAgentToken(_ *cobra.Command, args []string) error {
	if agentCmdState.TokenUpdateEnable && agentCmdState.TokenUpdateDisable {
		return errors.New("specify only one --enable or --disable")
	}

	// read the current state
	response, err := cli.LwApi.V2.AgentAccessTokens.Get(args[0])
	if err != nil {
		return errors.Wrap(err, "unable to get agent access token")
	}
	actual := response.Data
	updated := api.AgentAccessTokenRequest{
		TokenAlias: actual.TokenAlias,
		Enabled:    actual.Enabled,
		Props: &api.AgentAccessTokenProps{
			CreatedTime: actual.Props.CreatedTime,
		},
	}

	if agentCmdState.TokenUpdateEnable {
		updated.Enabled = 1
	}

	if agentCmdState.TokenUpdateDisable {
		updated.Enabled = 0
	}

	if agentCmdState.TokenUpdateName != "" {
		updated.TokenAlias = agentCmdState.TokenUpdateName
	}

	if agentCmdState.TokenUpdateDesc != "" {
		updated.Props.Description = agentCmdState.TokenUpdateDesc
	}

	response, err = cli.LwApi.V2.AgentAccessTokens.Update(args[0], updated)
	if err != nil {
		return errors.Wrap(err, "unable to update the agent access token")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(response.Data)
	}

	cli.OutputHuman(buildAgentTokenDetailsTable(response.Data))
	return nil
}

func createAgentToken(_ *cobra.Command, args []string) error {
	var desc string
	if len(args) == 2 {
		desc = args[1]
	}

	response, err := cli.LwApi.V2.AgentAccessTokens.Create(args[0], desc)
	if err != nil {
		return errors.Wrap(err, "unable to create agent access token")
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(response.Data)
	}

	cli.OutputHuman(buildAgentTokenDetailsTable(response.Data))
	return nil
}

func listAgentTokens(_ *cobra.Command, _ []string) error {
	response, err := cli.LwApi.V2.AgentAccessTokens.List()
	if err != nil {
		return errors.Wrap(err, "unable to list agent access token")
	}

	if len(response.Data) == 0 {
		cli.OutputHuman(
			"There are no agent access tokens. Try creating one with 'lacework agent token create%s'\n",
			cli.OutputNonDefaultProfileFlag(),
		)
		return nil
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(response.Data)
	}

	cli.OutputHuman(
		renderSimpleTable(
			[]string{"Token", "Name", "State"},
			agentTokensToTable(response.Data),
		),
	)
	return nil
}

func agentTokensToTable(tokens []api.AgentAccessToken) [][]string {
	out := [][]string{}
	for _, token := range tokens {
		out = append(out, []string{
			token.AccessToken,
			token.TokenAlias,
			token.PrettyState(),
		})
	}
	return out
}

func buildAgentTokenDetailsTable(token api.AgentAccessToken) string {
	return renderOneLineCustomTable("Agent Access Token Details",
		renderSimpleTable([]string{},
			[][]string{
				{"TOKEN", token.AccessToken},
				{"NAME", token.TokenAlias},
				{"DESCRIPTION", token.Props.Description},
				{"VERSION", token.Version},
				{"STATE", token.PrettyState()},
				{"CREATED AT", token.Props.CreatedTime.Format(time.RFC3339)},
			},
		),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetAutoWrapText(false)
		}),
	)
}
