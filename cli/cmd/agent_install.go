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
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/masterzen/winrm"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"

	"github.com/lacework/go-sdk/lwrunner"
)

// Official download url for installing Lacework agents
const (
	linuxAgentInstallDownloadURL   = "https://packages.lacework.net/install.sh"
	windowsAgentInstallDownloadURL = "https://updates.lacework.net/windows/preGA-0.5.0.2151/LWDataCollector.msi"
)

func installRemoteAgent(_ *cobra.Command, args []string) error {
	var (
		user    = agentCmdState.InstallSshUser
		port    = agentCmdState.InstallSshPort
		host    = args[0]
		authSet = false
	)

	// verify if the user specified the username via "user@host"
	if strings.Contains(host, "@") {
		userHost := strings.Split(host, "@")
		user = userHost[0]
		host = userHost[1]
	}

	// verify if the user specified the port via "host:port"
	if strings.Contains(host, ":") {
		userHost := strings.Split(host, ":")
		host = userHost[0]
		p, err := strconv.Atoi(userHost[1])
		if err != nil {
			return errors.Wrap(err, "invalid port")
		}
		port = p
	}

	token := agentCmdState.InstallAgentToken
	if token == "" {
		// user didn't provide an agent token
		cli.Log.Debugw("agent token not provided")
		var err error
		token, err = selectAgentAccessToken()
		if err != nil {
			return err
		}
	}

	if agentCmdState.WinRM {
		cli.Log.Debugw("installing windows agent", "user", user, "host", host)
		return installWindowsAgent(user, host, token, port)
	}

	cli.Log.Debugw("creating runner", "user", user, "host", host)
	runner := lwrunner.New(user, host, verifyHostCallback)

	if runner.Port != port {
		cli.Log.Debugw("ssh settings", "port", port)
		runner.Port = port
	}

	if runner.User == "" {
		cli.Log.Debugw("ssh username not set")
		user, err := askForUsername()
		if err != nil {
			return err
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
		identityFile, err := lwrunner.DefaultIdentityFilePath()
		if err != nil {
			return err
		}

		err = runner.UseIdentityFile(identityFile)
		if err != nil {
			cli.Log.Debugw("unable to use default identity file", "error", err)

			// if the default identity file didn't work, ask the user for auth details
			cli.Log.Debugw("ssh auth settings not configured")
			if err := askForAuthenticationDetails(runner); err != nil {
				return err
			}
		}
	}

	if err := verifyAccessToRemoteHost(runner); err != nil {
		return err
	}

	if err := isAgentInstalledOnRemoteHost(runner); err != nil {
		return err
	}

	cmd := fmt.Sprintf("sudo sh -c \"curl -sSL %s | sh -s -- %s\"", linuxAgentInstallDownloadURL, token)
	return runInstallCommandOnRemoteHost(runner, cmd)
}

func runInstallCommandOnRemoteHost(runner *lwrunner.Runner, cmd string) error {
	cli.StartProgress(" Installing agent on the remote host...")
	cli.Log.Debugw("exec remote command", "cmd", cmd)
	stdout, stderr, err := runner.Exec(cmd)
	cli.StopProgress()
	cli.Log.Debugw("remote command results",
		"cmd", cmd,
		"stdout", stdout.String(),
		"stderr", stderr.String(),
		"error", err,
	)
	if err != nil {
		return errors.Wrap(formatRunnerError(stdout, stderr, err), "unable to install agent on the remote host")
	}

	cli.OutputHuman("Lacework agent installed successfully on host %s\n\n", runner.Hostname)
	cli.OutputHuman(renderOneLineCustomTable("Installation Details", stdout.String(),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetColumnSeparator(" ")
			t.SetAutoWrapText(false)
		})))
	return nil
}

func isAgentInstalledOnRemoteHost(runner *lwrunner.Runner) error {
	agentVersionCmd := "sudo sh -c \"/var/lib/lacework/datacollector -v\""

	cli.StartProgress("Verifying previous agent installations...")
	cli.Log.Debugw("exec remote command", "cmd", agentVersionCmd)
	stdout, stderr, err := runner.Exec(agentVersionCmd)
	cli.StopProgress()
	cli.Log.Debugw("remote command results", "cmd", agentVersionCmd,
		"stdout", stdout.String(),
		"stderr", stderr.String(),
		"error", err,
	)

	if err != nil {
		// if we couldn't run the agent version command it means that
		// the agent is not yet installed, so we return nil to continue
		// with the agent installation process
		return nil
	}

	if agentCmdState.InstallForce {
		cli.Log.Debugw("forcing previous agent installation on remote host")
		return nil
	}

	return errors.Errorf("agent already installed on the remote host. %s", stderr.String())
}

func verifyAccessToRemoteHost(runner *lwrunner.Runner) error {
	accessCmd := "echo we-are-in"

	cli.StartProgress("Verifying access to the remote host...")
	cli.Log.Debugw("exec remote command", "cmd", accessCmd)
	stdout, stderr, err := runner.Exec(accessCmd)
	cli.StopProgress()
	cli.Log.Debugw("remote command results", "cmd", accessCmd,
		"stdout", stdout.String(),
		"stderr", stderr.String(),
		"error", err,
	)

	if err != nil || !strings.Contains(stdout.String(), "we-are-in") {
		return errors.Wrap(formatRunnerError(stdout, stderr, err), "unable to connect to the remote host")
	}

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
		tokenNames = make([]string, 0)
		tokenName  = ""
	)
	for _, aTkn := range response.Data {
		// only display tokens that have a name (a.k.a Alias)
		if strings.TrimSpace(aTkn.TokenAlias) != "" {
			tokenNames = append(tokenNames, aTkn.TokenAlias)
		}
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

// ask for the ssh username
func askForUsername() (string, error) {
	var user string

	err := survey.AskOne(&survey.Input{
		Message: "SSH username:",
	}, &user, survey.WithValidator(survey.Required))
	if err != nil {
		return "", errors.Wrap(err, "unable to ask for username")
	}

	return user, nil
}

func verifyHostCallback(host string, remote net.Addr, key ssh.PublicKey) error {
	// error if key does not exist inside the default known_hosts file,
	// or if host in known_hosts file but key changed!
	hostFound, err := lwrunner.CheckKnownHost(host, remote, key, "")
	if hostFound && err != nil {
		// the host in known_hosts file was found but key mismatch
		return err
	}

	// handshake because public key already exists
	if hostFound && err == nil {
		return nil
	}

	if agentCmdState.InstallTrustHostKey {
		// the user wants to add the new host to known hosts file automatically
		return lwrunner.AddKnownHost(host, remote, key, "")
	}

	// ask user to check if he/she trust the host public key
	if askIsHostTrusted(host, key) {
		// add the new host to known hosts file.
		return lwrunner.AddKnownHost(host, remote, key, "")
	}

	// non trusted key
	return errors.New("you typed no, the agent installation was aborted!")
}

// ask user to check if he/she trust the host public key
func askIsHostTrusted(host string, key ssh.PublicKey) bool {
	// about to ask a question to the user
	cli.StopProgress()

	var (
		trust    = false
		question = fmt.Sprintf(
			"Unknown Host: %s\nFingerprint: %s\nWould you like to continue connecting?",
			host, ssh.FingerprintSHA256(key),
		)
		err = survey.AskOne(&survey.Confirm{
			Message: question,
			Help:    "By typing 'yes', the host will be added to the $HOME/.ssh/known_hosts file.",
		}, &trust)
	)
	if err != nil {
		cli.Log.Debugw("unable to ask if host is trusted", "error", err)
		return false
	}
	return trust
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

func formatRunnerError(stdout, stderr bytes.Buffer, err error) error {
	formatted := ""

	if stdout.String() != "" {
		formatted = fmt.Sprintf("%s\n\nSTDOUT:\n%s", formatted, stdout.String())
	}

	if stderr.String() != "" {
		formatted = fmt.Sprintf("%s\n\nSTDERR:\n%s", formatted, stderr.String())
	}

	if formatted == "" {
		return err
	}

	if err == nil {
		return errors.New(formatted)
	}

	return errors.Wrap(err, formatted)
}

func installWindowsAgent(user, host, token string, port int) error {
	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 5*time.Minute)
	client, err := winrm.NewClient(endpoint, user, agentCmdState.InstallPassword)
	if err != nil {
		return err
	}

	cli.StartProgress("Verifying access to the remote host...")
	defer cli.StopProgress()
	_, stderr, code, err := client.RunPSWithString("echo we-are-in", "")
	if err != nil || code != 0 {
		return errors.Wrap(err,
			fmt.Sprintf("unable to access the remote host via WinRM transport.\n\nexitcode: %d\nstderr: %s",
				code, stderr))
	}

	cli.StartProgress(fmt.Sprintf("Downloading Windows agent installer (%s) ...", windowsAgentInstallDownloadURL))
	// @afiune create a temporal directory and store the MSI there
	downloadMSI := `$msi_dest="C:/Users/vagrant/LWDataCollector.msi"
	$wc = New-Object System.Net.WebClient
	$wc.DownloadFile("` + windowsAgentInstallDownloadURL + `", $msi_dest)`
	_, stderr, code, err = client.RunPSWithString(downloadMSI, "")
	if err != nil || code != 0 {
		return errors.Wrap(err,
			fmt.Sprintf("unable to access the remote host via WinRM transport.\n\nexitcode: %d\nstderr: %s",
				code, stderr))
	}

	cli.StartProgress("Installing Lacework Windows agent...")
	installAgent := `Start-Process msiexec.exe -Wait -ArgumentList '/i C:\Users\vagrant\LWDataCollector.msi ACCESSTOKEN=` +
		token + ` SERVERURL=https://api.lacework.net /qn'`
	_, stderr, code, err = client.RunPSWithString(installAgent, "")
	if err != nil || code != 0 {
		return errors.Wrap(err,
			fmt.Sprintf("unable to access the remote host via WinRM transport.\n\nexitcode: %d\nstderr: %s",
				code, stderr))
	}

	// stdout, stderr, code, err = client.RunPSWithString("ls C:\\Windows\\System32\\drivers\\lwdcs.sys", "")
	cli.StartProgress("Verifying installation ...")
	var (
		stderrBuf    bytes.Buffer
		stdoutBuf    bytes.Buffer
		stderrWriter = bufio.NewWriter(&stderrBuf)
		stdoutWriter = bufio.NewWriter(&stdoutBuf)
	)
	code, err = client.Run("sc query lwdatacollector", stdoutWriter, stderrWriter)
	if err != nil || code != 0 {
		return errors.Wrap(err,
			fmt.Sprintf("unable to access the remote host via WinRM transport.\n\nexitcode: %d\nstdout: %s\nstderr: %s",
				code, stdoutBuf.String(), stderrBuf.String()))
	}
	cli.StopProgress()
	cli.OutputHuman(stdoutBuf.String())
	cli.OutputHuman("\nLacework agent installed successfully on host %s\n", host)
	return nil
}
