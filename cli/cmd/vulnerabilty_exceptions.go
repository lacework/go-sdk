//
// Author:: Darren Murray(<darren.murray@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/lacework/go-sdk/api"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// vulnerability-exceptions command is used to manage lacework vulnerability exceptions
	vulnerabilityExceptionCommand = &cobra.Command{
		Use:     "vulnerability-exception",
		Aliases: []string{"vulnerability-exceptions", "ve", "vuln-exception", "vuln-exceptions"},
		Short:   "Manage vulnerability exceptions",
		Long:    "Manage vulnerability exceptions to control and customize your alert profile for hosts and containers.",
	}

	// list command is used to list all lacework vulnerability exceptions
	vulnerabilityExceptionListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all vulnerability exceptions",
		Long:    "List all vulnerability exceptions configured in your Lacework account.",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			vulnerabilityExceptions, err := cli.LwApi.V2.VulnerabilityExceptions.List()
			if err != nil {
				return errors.Wrap(err, "unable to get vulnerability exceptions")
			}
			if len(vulnerabilityExceptions.Data) == 0 {
				msg := `There are no vulnerability exceptions configured in your account.

Get started by integrating your vulnerability exceptions to manage alerting using the command:

    lacework vulnerability-exception create

If you prefer to configure vulnerability exceptions via the WebUI, log in to your account at:

    https://%s.lacework.net

Then navigate to Vulnerabilities > Exceptions.
`
				cli.OutputHuman(fmt.Sprintf(msg, cli.Account))
				return nil
			}

			if cli.JSONOutput() {
				return cli.OutputJSON(vulnerabilityExceptions)
			}

			var rows [][]string
			for _, vuln := range vulnerabilityExceptions.Data {
				rows = append(rows, []string{vuln.Guid, vuln.ExceptionName, vuln.ExceptionType, vuln.Status()})
			}

			cli.OutputHuman(renderCustomTable([]string{"GUID", "NAME", "TYPE", "STATE"}, rows,
				tableFunc(func(t *tablewriter.Table) {
					t.SetBorder(false)
					t.SetAutoWrapText(false)
				})))
			return nil
		},
	}
	// show command is used to retrieve a lacework vulnerability exception by id
	vulnerabilityExceptionShowCommand = &cobra.Command{
		Use:   "show",
		Short: "Get vulnerability exception by ID",
		Long:  "Get a single vulnerability exception by it's vulnerability exception ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var response api.VulnerabilityExceptionResponse
			err := cli.LwApi.V2.VulnerabilityExceptions.Get(args[0], &response)
			if err != nil {
				return errors.Wrap(err, "unable to get vulnerability exception")
			}
			vuln := response.Data

			if cli.JSONOutput() {
				return cli.OutputJSON(vuln)
			}

			var groupCommon [][]string
			groupCommon = append(groupCommon, []string{vuln.Guid, vuln.ExceptionName, vuln.ExceptionType, vuln.Status()})

			cli.OutputHuman(renderSimpleTable([]string{"GUID", "NAME", "TYPE", "STATUS"}, groupCommon))
			cli.OutputHuman("\n")
			cli.OutputHuman(buildVulnerabilityExceptionsPropsTable(vuln))
			return nil
		},
	}

	// delete command is used to remove a lacework vulnerability exception by id
	vulnerabilityExceptionDeleteCommand = &cobra.Command{
		Use:   "delete",
		Short: "Delete a vulnerability exception",
		Long:  "Delete a single vulnerability exception by it's vulnerability exception ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			err := cli.LwApi.V2.VulnerabilityExceptions.Delete(args[0])
			if err != nil {
				return errors.Wrap(err, "unable to delete vulnerability exception")
			}
			return nil
		},
	}

	// create command is used to create a new lacework vulnerability exception
	vulnerabilityExceptionCreateCommand = &cobra.Command{
		Use:   "create",
		Short: "Create a new vulnerability exception",
		Long:  "Creates a new single vulnerability exception.",
		RunE: func(_ *cobra.Command, args []string) error {
			if !cli.InteractiveMode() {
				return errors.New("interactive mode is disabled")
			}

			vulnID, err := promptCreateVulnerabilityException()
			if err != nil {
				return errors.Wrap(err, "unable to create vulnerability exception")
			}

			cli.OutputHuman("The vulnerability exception created with GUID %s.\n", vulnID)
			return nil
		},
	}
)

func init() {
	// add the vulnerability-exception command
	rootCmd.AddCommand(vulnerabilityExceptionCommand)

	// add sub-commands to the vulnerability-exception command
	vulnerabilityExceptionCommand.AddCommand(vulnerabilityExceptionListCommand)
	vulnerabilityExceptionCommand.AddCommand(vulnerabilityExceptionShowCommand)
	vulnerabilityExceptionCommand.AddCommand(vulnerabilityExceptionCreateCommand)
	vulnerabilityExceptionCommand.AddCommand(vulnerabilityExceptionDeleteCommand)
}

func buildVulnerabilityExceptionsPropsTable(vuln api.VulnerabilityException) string {
	props := setProps(vuln)

	return renderOneLineCustomTable("VULNERABILITY EXCEPTION PROPS",
		renderCustomTable([]string{}, props,
			tableFunc(func(t *tablewriter.Table) {
				t.SetBorder(false)
				t.SetColumnSeparator(" ")
				t.SetAutoWrapText(false)
				t.SetAlignment(tablewriter.ALIGN_LEFT)
			}),
		),
		tableFunc(func(t *tablewriter.Table) {
			t.SetBorder(false)
			t.SetAutoWrapText(false)
		}),
	)
}

func setProps(vuln api.VulnerabilityException) [][]string {
	var (
		details [][]string
	)
	details = append(details, []string{"DESCRIPTION", vuln.Props.Description})
	details = append(details, []string{"UPDATED BY", vuln.Props.UpdatedBy})
	details = append(details, []string{"LAST UPDATED", vuln.UpdatedTime})
	details = append(details, []string{"CREATED", vuln.CreatedTime})
	details = append(details, []string{"REASON", vuln.ExceptionReason})

	if vuln.ExceptionType == api.VulnerabilityExceptionTypeHost.String() {
		details = append(details, []string{"NAMESPACES", strings.Join(vuln.ResourceScope.Namespace, ", ")})
		details = append(details, []string{"HOSTNAMES", strings.Join(vuln.ResourceScope.Hostname, ", ")})
		details = append(details, []string{"EXTERNAL IPS", strings.Join(vuln.ResourceScope.ExternalIP, ",")})
		details = append(details, []string{"CLUSTER NAMES", strings.Join(vuln.ResourceScope.ClusterName, ", ")})
	} else if vuln.ExceptionType == api.VulnerabilityExceptionTypeContainer.String() {
		details = append(details, []string{"NAMESPACES", strings.Join(vuln.ResourceScope.Namespace, ", ")})
		details = append(details, []string{"IMAGE IDS", strings.Join(vuln.ResourceScope.ImageID, ", ")})
		details = append(details, []string{"IMAGE TAGS", strings.Join(vuln.ResourceScope.ImageTag, ", ")})
		details = append(details, []string{"REGISTRIES", strings.Join(vuln.ResourceScope.Registry, ", ")})
		details = append(details, []string{"REPOSITORIES", strings.Join(vuln.ResourceScope.Repository, ", ")})
	}

	details = append(details, []string{"FIXABLE", vulnerabilityExceptionFixableEnabled(vuln.VulnerabilityCriteria.Fixable)})
	details = append(details, []string{"CVES", strings.Join(vuln.VulnerabilityCriteria.Cve, ", ")})
	details = append(details, []string{"SEVERITIES", strings.Join(vuln.VulnerabilityCriteria.Severity, ", ")})

	if vuln.VulnerabilityCriteria.Package != nil {
		packages, err := json.Marshal(vuln.VulnerabilityCriteria.Package)
		if err != nil {
			packages = []byte{}
		}
		details = append(details, []string{"PACKAGES", string(packages)})
	} else {
		details = append(details, []string{"PACKAGES", ""})
	}
	return details
}

func promptCreateVulnerabilityException() (string, error) {
	var (
		group  = ""
		prompt = &survey.Select{
			Message: "Choose a vulnerability exception type to create: ",
			Options: []string{
				"Host",
				"Container",
			},
		}
		err = survey.AskOne(prompt, &group)
	)
	if err != nil {
		return "", err
	}

	switch group {
	case "Host":
		return createHostVulnerabilityException()
	case "Container":
		return createContainerVulnerabilityException()
	default:
		return "", errors.New("unknown vulnerability exception type")
	}
}

func vulnerabilityExceptionFixableEnabled(fixable []int) string {
	if len(fixable) == 0 {
		return "false"
	}
	return strconv.FormatBool(fixable[0] == 1)
}

func askVulnerabilityExceptionCriteria(answers interface{}, criteria []string) error {
	var questions []*survey.Question
	for _, c := range criteria {
		if c == "CVEs" {
			questions = append(questions,
				&survey.Question{
					Name:     "cves",
					Prompt:   &survey.Multiline{Message: "List of CVE IDs:"},
					Validate: validateCveFormat(),
				})
			continue
		}

		if c == "Severities" {
			questions = append(questions,
				&survey.Question{
					Name: "severities",
					Prompt: &survey.MultiSelect{
						Message: "Select severities:",
						Options: []string{"Critical", "High", "Medium", "Low", "Info"},
					},
					Validate: validateSeverities(),
				})
			continue
		}
		if c == "Packages" {
			questions = append(questions,
				&survey.Question{
					Name:   "packages",
					Prompt: &survey.Multiline{Message: "List of 'package:version' packages to include:"},
				})
			continue
		}

		if c == "Packages" {
			questions = append(questions,
				&survey.Question{
					Name:   "fixable",
					Prompt: &survey.Confirm{Message: "Include Fixable:"},
				})
			continue
		}
	}

	err := survey.Ask(questions, answers, survey.WithIcons(promptIconsFunc))
	if err != nil {
		return err
	}
	return nil
}

func transformVulnerabilityExceptionPackages(packages string) []api.VulnerabilityExceptionPackage {
	if packages == "" {
		return []api.VulnerabilityExceptionPackage{}
	}
	var vulnPackages []api.VulnerabilityExceptionPackage
	packageList := strings.Split(packages, "\n")
	for _, pack := range packageList {
		vulnPackage := strings.Split(pack, ":")
		vulnPackages = append(vulnPackages, api.VulnerabilityExceptionPackage{Name: vulnPackage[0], Version: vulnPackage[1]})
	}
	return vulnPackages
}

func validateCveFormat() survey.Validator {
	return func(val interface{}) error {
		cveRegEx, _ := regexp.Compile(`(?i)CVE-\d{4}-\d{4,7}`)
		if list, ok := val.([]core.OptionAnswer); ok {
			for _, i := range list {
				if !cveRegEx.MatchString(i.Value) {
					return fmt.Errorf("CVE format is invalid. Please format corretly eg: CVE-2014-0001")
				}
			}
		} else {
			value := val.(string)
			if !cveRegEx.MatchString(value) {
				return fmt.Errorf("CVE format is invalid. Please format corretly eg: CVE-2014-0001")
			}
		}
		return nil
	}
}

func validateSeverities() survey.Validator {
	return func(val interface{}) error {
		if list, ok := val.([]core.OptionAnswer); ok {
			for _, i := range list {
				match := strings.Contains(strings.Join(api.ValidEventSeverities, ", "), strings.ToLower(i.Value))
				if !match {
					return fmt.Errorf("severity '%s' is invalid. Must be one of 'Critical', 'High', 'Medium', 'Low', 'Info'", i.Value)
				}
			}
		} else {
			value := val.(core.OptionAnswer).Value
			match := strings.Contains(strings.Join(api.ValidEventSeverities, ", "), strings.ToLower(value))
			if !match {
				return fmt.Errorf("severity '%s' is invalid. Must be one of 'Critical', 'High', 'Medium', 'Low', 'Info'", value)
			}
		}
		return nil
	}
}
