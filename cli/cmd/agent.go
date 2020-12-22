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
	"path"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwrunner"
)

var (
	agentCmdState = struct {
		TokenUpdateEnable   bool
		TokenUpdateDisable  bool
		TokenUpdateName     string
		TokenUpdateDesc     string
		InstallSshUser      string
		InstallAgentToken   string
		InstallPassword     string
		InstallIdentityFile string
	}{}

	defaultSshIdentityKey = "~/.ssh/id_rsa"

	agentCmd = &cobra.Command{
		Use:   "agent",
		Short: "manage Lacework agents",
		Long: `Manage agents and agent access tokens in your account.

To analyze application, host, and user behavior, Lacework uses a lightweight agent,
which securely forwards collected metadata to the Lacework cloud for analysis. The
agent requires minimal system resources and runs on most 64-bit Linux distributions.

For a complete list of supported operating systems, visit:

    https://support.lacework.com/hc/en-us/articles/360005230014-Supported-Operating-Systems`,
	}

	agentTokenCmd = &cobra.Command{
		Use:     "token",
		Aliases: []string{"tokens"},
		Short:   "manage agent access tokens",
		Long: `Manage agent access tokens in your account.

Agent tokens should be treated as secret and not published. A token uniquely identifies
a Lacework customer. If you suspect your token has been publicly exposed or compromised,
generate a new token, update the new token on all machines using the old token. When
complete, the old token can safely be disabled without interrupting Lacework services.`,
	}

	agentTokenListCmd = &cobra.Command{
		Use:   "list",
		Short: "list all agent access tokens",
		Long:  `List all agent access tokens.`,
		Args:  cobra.NoArgs,
		RunE:  listAgentTokens,
	}

	agentTokenCreateCmd = &cobra.Command{
		Use:   "create <name> [description]",
		Short: "create a new agent access token",
		Long:  `Create a new agent access token.`,
		Args:  cobra.RangeArgs(1, 2),
		RunE:  createAgentToken,
	}

	agentTokenShowCmd = &cobra.Command{
		Use:   "show <token>",
		Short: "show details about an agent access token",
		Long:  `Show details about an agent access token.`,
		Args:  cobra.ExactArgs(1),
		RunE:  showAgentToken,
	}

	agentTokenUpdateCmd = &cobra.Command{
		Use:   "update <token>",
		Short: "update an agent access token",
		Long: `Update an agent access token.

To update the token name and description:

    $ lacework agent token update <token> --name dev --description "k8s deployment for dev"

To disable a token:

    $ lacework agent token update <token> --disable

To enable a token:

    $ lacework agent token update <token> --enable`,
		Args: cobra.ExactArgs(1),
		RunE: updateAgentToken,
	}

	// TODO hidden for now
	agentListCmd = &cobra.Command{
		Use:    "list",
		Short:  "list all hosts with a running agent",
		Long:   `List all hosts in your environment that has a running agent`,
		Hidden: true,
		RunE:   listAgents,
	}

	// TODO hidden for now
	agentGenerateCmd = &cobra.Command{
		Use:    "generate",
		Short:  "generate agent deployment scripts",
		Long:   `TBA`,
		Hidden: true,
		RunE:   listAgentTokens,
	}

	// TODO hidden for now
	agentInstallCmd = &cobra.Command{
		Use:   "install <host>",
		Short: "install the datacollector agent on a remote host",
		Args:  cobra.ExactArgs(1),
		Long: `For single host installation of the Lacework agent via Secure Shell (SSH).

When this command is executed without any additional flag, an interactive prompt will be
launched to help gather the necessary authentication information to access the remote host.

To authenticate to the remote host with a username and password.

    $ lacework agent install <host> --ssh_username <your-user> --ssh_password <secret>

To authenticate to the remote host with an identity file instead.

    $ lacework agent install <user@host> -i /path/to/your/key

To provide an agent access token of your choice, use the command 'lacework agent token list',
select a token and pass it to the '--token' flag.

    $ lacework agent install <user@host> -i /path/to/your/key --token <token>
    `,
		RunE: installRemoteAgent,
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

	// add the list sub-command to the 'agent token' cmd
	agentTokenCmd.AddCommand(agentTokenListCmd)
	agentTokenCmd.AddCommand(agentTokenCreateCmd)
	agentTokenCmd.AddCommand(agentTokenShowCmd)
	agentTokenCmd.AddCommand(agentTokenUpdateCmd)

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
		"ssh_username", "", "user to log in as on the remote host",
	)
	agentInstallCmd.Flags().StringVar(&agentCmdState.InstallAgentToken,
		"token", "", "agent access token",
	)
}

func listAgents(_ *cobra.Command, _ []string) error {
	// @afiune POC - This depends on LQL
	time.Sleep(500 * time.Millisecond)
	response, err := loadAgents()
	if err != nil {
		return errors.Wrap(err, "beta feature not yet supported")
	}
	if cli.JSONOutput() {
		return cli.OutputJSON(response.Data)
	}

	cli.OutputHuman(
		renderSimpleTable(
			[]string{"Hostname", "Name", "IP Address", "External IP", "Status", "OS Arch", "Version"},
			agentsToTable(response.Data),
		),
	)
	return nil
}

func selectAgentAccessToken() (string, error) {
	cli.StartProgress(" Searching for agent access tokens...")
	response, err := cli.LwApi.Agents.ListTokens()
	cli.StopProgress()
	if err != nil {
		return "", errors.Wrap(err, "unable to list agent access token")
	}

	var (
		tokenNames = make([]string, len(response.Data))
		tokenName  = ""
	)
	for i, aTkn := range response.Data {
		tokenNames[i] = aTkn.TokenAlias
	}

	err = survey.AskOne(&survey.Select{
		Message: "Choose an agent access token: ",
		Options: tokenNames,
	}, &tokenName, survey.WithValidator(survey.Required))
	if err != nil {
		return "", errors.Wrap(err, "unable to ask for agent access token")
	}
	for _, aTkn := range response.Data {
		if tokenName == aTkn.TokenAlias {
			return aTkn.AccessToken, nil
		}
	}

	// @afiune this should never happen
	return "", errors.New("something went pretty wrong here, contact support@lacework.net")
}

func installRemoteAgent(_ *cobra.Command, args []string) error {
	token := agentCmdState.InstallAgentToken
	if token == "" {
		// user didn't provide an agent token
		var err error
		token, err = selectAgentAccessToken()
		if err != nil {
			return err
		}
	}

	var (
		// TODO @afiune where can we get it?
		sha         = "3.3.5_2020-11-16_master_ac0e65055f11f4f59bab6ea4dfa61dcafaa9a3f1"
		downloadUrl = fmt.Sprintf("https://s3-us-west-2.amazonaws.com/www.lacework.net/download/%s/install.sh", sha)
		cmd         = fmt.Sprintf("sudo sh -c \"curl -sSL %s | sh -s -- %s\"", downloadUrl, token)
		user        = agentCmdState.InstallSshUser
		host        = args[0]
		authSet     = false
	)

	// verify if the user specified the username via user@host
	if strings.Contains(host, "@") {
		userHost := strings.Split(host, "@")
		user = userHost[0]
		host = userHost[1]
	}

	cli.Log.Debugw("creating runner", "user", user, "host", host)
	runner := lwrunner.New(user, host)

	if runner.User == "" {
		cli.Log.Debugw("ssh username not set")
		// ask for the username
		err := survey.AskOne(&survey.Input{
			Message: "SSH username:",
		}, &user, survey.WithValidator(survey.Required))
		if err != nil {
			return errors.Wrap(err, "unable to ask for username")
		}

		runner.User = user
		cli.Log.Debugw("ssh settings", "user", runner.User)
	}

	if agentCmdState.InstallIdentityFile != defaultSshIdentityKey {
		cli.Log.Debugw("ssh settings", "identity_file", agentCmdState.InstallIdentityFile)
		err := runner.UseIdentityFile(agentCmdState.InstallIdentityFile)
		if err != nil {
			return errors.Wrap(err, "unable to use provided identity file")
		}
		authSet = true
	}

	if agentCmdState.InstallPassword != "" {
		cli.Log.Debugw("ssh settings", "auth", "password_from_flag")
		runner.UsePassword(agentCmdState.InstallPassword)
		authSet = true
	}

	// if no authentication was set
	if !authSet {
		// try to use the default identity file
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		err = runner.UseIdentityFile(path.Join(home, ".ssh", "id_rsa"))
		if err != nil {
			cli.Log.Debugw("unable to use default identity file", "error", err)

			// if the default identity file didn't work, ask the user for auth details
			cli.Log.Debugw("ssh auth settings not configured")
			err := askForAuthenticationDetails(runner)
			if err != nil {
				return err
			}
		}
	}

	cli.StartProgress(" Installing agent on remote host...")
	cli.Log.Debugw("executing remote command", "cmd", cmd)
	stdout, stderr, err := runner.Exec(cmd)
	cli.StopProgress()
	if err != nil {
		cli.OutputHuman(stderr.String())
		return errors.Wrap(err, "unable to install agent")
	}

	cli.OutputHuman("Lacework agent installed successfully on host %s\n\n", host)
	cli.OutputHuman("INSTALLATION DETAILS\n")
	cli.OutputHuman("-----------------------------------------------------------------\n")
	cli.OutputHuman(stdout.String())
	return nil
}

func askForAuthenticationDetails(runner *lwrunner.Runner) error {
	authMethod := ""
	err := survey.AskOne(&survey.Select{
		Message: "Choose SSH authentication method: ",
		Options: []string{"Identity File", "Password"},
	}, &authMethod, survey.WithValidator(survey.Required))
	if err != nil {
		return errors.Wrap(err, "unable to ask for authentication method")
	}
	switch authMethod {
	case "Password":
		// ask for a password
		var password string
		err = survey.AskOne(&survey.Password{
			Message: "SSH password:",
		}, &password, survey.WithValidator(survey.Required))
		if err != nil {
			return errors.Wrap(err, "unable to ask for password")
		}

		runner.UsePassword(password)
		cli.Log.Debugw("ssh settings", "auth", "password_from_input")
	default:
		// ask for an identity file
		var identityFile string
		err = survey.AskOne(&survey.Input{
			Message: "SSH identity file:",
		}, &identityFile, survey.WithValidator(survey.Required))
		if err != nil {
			return errors.Wrap(err, "unable to ask for identity file")
		}

		err = runner.UseIdentityFile(identityFile)
		if err != nil {
			return errors.Wrap(err, "unable to use provided identity file")
		}
		cli.Log.Debugw("ssh settings", "identity_file", identityFile)
	}

	return nil
}

func showAgentToken(_ *cobra.Command, args []string) error {
	response, err := cli.LwApi.Agents.GetToken(args[0])
	if err != nil {
		return errors.Wrap(err, "unable to get agent access token")
	}

	if len(response.Data) == 0 {
		return errors.New(`unable to create agent access token

The platform did not return any token in the response body, this is very
unlikely to happen but, hey it happened. Please help us improve the
Lacework CLI by reporting this issue at:

  https://support.lacework.com/hc/en-us/requests/new
`)
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(response.Data[0])
	}

	cli.OutputHuman(buildAgentTokenDetailsTable(response.Data[0]))
	return nil
}

func updateAgentToken(_ *cobra.Command, args []string) error {
	if agentCmdState.TokenUpdateEnable && agentCmdState.TokenUpdateDisable {
		return errors.New("specify only one --enable or --disable")
	}

	// read the current state
	response, err := cli.LwApi.Agents.GetToken(args[0])
	if err != nil {
		return errors.Wrap(err, "unable to get agent access token")
	}
	actual := response.Data[0]
	updated := api.AgentTokenRequest{
		TokenAlias: actual.TokenAlias,
		Enabled:    actual.EnabledInt(),
		Props: &api.AgentTokenProps{
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

	response, err = cli.LwApi.Agents.UpdateToken(args[0], updated)
	if err != nil {
		return errors.Wrap(err, "unable to update the agent access token")
	}

	if len(response.Data) == 0 {
		return errors.New(`unable to update the agent access token

The platform did not return any token in the response body, this is very
unlikely to happen but, hey it happened. Please help us improve the
Lacework CLI by reporting this issue at:

  https://support.lacework.com/hc/en-us/requests/new
`)
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(response.Data[0])
	}

	cli.OutputHuman(buildAgentTokenDetailsTable(response.Data[0]))
	return nil
}

func createAgentToken(_ *cobra.Command, args []string) error {
	var desc string
	if len(args) == 2 {
		desc = args[1]
	}

	response, err := cli.LwApi.Agents.CreateToken(args[0], desc)
	if err != nil {
		return errors.Wrap(err, "unable to create agent access token")
	}

	if len(response.Data) == 0 {
		return errors.New(`unable to create agent access token

The platform did not return any token in the response body, this is very
unlikely to happen but, hey it happened. Please help us improve the
Lacework CLI by reporting this issue at:

  https://support.lacework.com/hc/en-us/requests/new
`)
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(response.Data[0])
	}

	cli.OutputHuman(buildAgentTokenDetailsTable(response.Data[0]))
	return nil
}

func listAgentTokens(_ *cobra.Command, _ []string) error {
	response, err := cli.LwApi.Agents.ListTokens()
	if err != nil {
		return errors.Wrap(err, "unable to list agent access token")
	}

	if len(response.Data) == 0 {
		cli.OutputHuman("There are no agent access tokens. Try creating one with 'lacework agent token create'\n")
		return nil
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(response.Data)
	}

	cli.OutputHuman(
		renderSimpleTable(
			[]string{"Token", "Name", "Status"},
			agentTokensToTable(response.Data),
		),
	)
	return nil
}

func agentsToTable(agents []AgentHost) [][]string {
	out := [][]string{}
	for _, agent := range agents {
		out = append(out, []string{
			agent.MachineHostname,
			agent.Name,
			agent.MachineIP,
			agent.Tags.ExternalIP,
			agent.Status,
			fmt.Sprintf("%s/%s", agent.Tags.Os, agent.Tags.Arch),
			agent.AgentVersion,
		})
	}

	// order by severity
	sort.Slice(out, func(i, j int) bool {
		return out[i][1] < out[j][1]
	})

	return out
}

func agentTokensToTable(tokens []api.AgentToken) [][]string {
	out := [][]string{}
	for _, token := range tokens {
		out = append(out, []string{
			token.AccessToken,
			token.TokenAlias,
			token.PrettyStatus(),
		})
	}
	return out
}

func buildAgentTokenDetailsTable(token api.AgentToken) string {
	return renderOneLineCustomTable("Agent Token Details",
		renderSimpleTable([]string{},
			[][]string{
				[]string{"TOKEN", token.AccessToken},
				[]string{"NAME", token.TokenAlias},
				[]string{"DESCRIPTION", token.Props.Description},
				[]string{"ACCOUNT", token.Account},
				[]string{"VERSION", token.Version},
				[]string{"STATUS", token.PrettyStatus()},
				[]string{"CREATED AT", token.Props.CreatedTime.Format(time.RFC3339)},
				[]string{"UPDATED AT", token.LastUpdatedTime.Format(time.RFC3339)},
			},
		),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetAutoWrapText(false)
		}),
	)
}
