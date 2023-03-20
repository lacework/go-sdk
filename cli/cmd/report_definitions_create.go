//
// Author:: Darren Murray(<darren.murray@lacework.net>)
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

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// create command is used to create a new lacework report definition
var reportDefinitionsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a report definition",
	Long: `Create a new report definition to view the evaluation of a set of policies in a report.

To create a new report definition:

    lacework report-definition create

To create a new report definition from an existing file:

    lacework report-definition create --file custom-report.json
`,
	Args: cobra.NoArgs,
	RunE: createReportDefinition,
}

func createReportDefinition(_ *cobra.Command, args []string) error {
	var (
		reportDefinition api.ReportDefinition
		err              error
	)

	if reportDefinitionsCmdState.File != "" {
		fileInput, err := inputReportDefinitionFromFile(reportDefinitionsCmdState.File)
		if err != nil {
			return err
		}

		cfg, err := parseNewReportDefinition(fileInput)
		if err != nil {
			return err
		}
		reportDefinition = api.NewReportDefinition(cfg)
	} else {
		reportDefinition, err = promptCreateReportDefinition()
		if err != nil {
			return err
		}
	}

	cli.StartProgress("Creating report definition...")
	resp, err := cli.LwApi.V2.ReportDefinitions.Create(reportDefinition)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to create report definition")
	}

	cli.OutputHuman("New report definition created. To view the report run:\n\n"+
		"lacework report-definition show %s \n", resp.Data.ReportDefinitionGuid)
	return nil
}

func inputReportDefinitionFromFile(filePath string) (string, error) {
	fileData, err := ioutil.ReadFile(filePath)

	if err != nil {
		return "", errors.Wrap(err, "unable to read file")
	}

	return string(fileData), nil
}

func promptCreateReportDefinition() (api.ReportDefinition, error) {
	var useExisting bool

	if err := survey.AskOne(&survey.Confirm{Message: CreateReportDefinitionQuestion}, &useExisting); err != nil {
		return api.ReportDefinition{}, err
	}

	if useExisting {
		return promptCreateReportDefinitionFromExisting()
	}

	return promptCreateReportDefinitionFromNew()
}

func promptCreateReportDefinitionFromNew() (reportDefinition api.ReportDefinition, err error) {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: CreateReportDefinitionReportNameQuestion},
			Validate: survey.Required,
		},
		{
			Name:     "display",
			Prompt:   &survey.Input{Message: CreateReportDefinitionDisplayNameQuestion},
			Validate: survey.Required,
		},
		{
			Name:     "subType",
			Prompt:   &survey.Select{Message: CreateReportDefinitionReportSubTypeQuestion, Options: api.ReportDefinitionSubTypes()},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name        string `survey:"name"`
		DisplayName string `survey:"display"`
		SubType     string `survey:"subType"`
	}{}

	if err = survey.Ask(questions, &answers, survey.WithIcons(promptIconsFunc)); err != nil {
		return
	}

	reportDefinition = api.ReportDefinition{
		ReportName:    answers.Name,
		DisplayName:   answers.DisplayName,
		SubReportType: answers.SubType,
		ReportType:    api.ReportDefinitionTypeCompliance.String(),
	}

	var sections []api.ReportDefinitionSection

	// add sections
	if err = promptAddReportDefinitionSection(&sections); err != nil {
		return
	}

	addSection := false
	for {
		if err = survey.AskOne(&survey.Confirm{
			Message: CreateReportDefinitionAddSectionQuestion,
		}, &addSection); err != nil {
			return
		}

		if addSection {
			if err = promptAddReportDefinitionSection(&sections); err != nil {
				return
			}
		} else {
			break
		}
	}

	reportDefinition.ReportDefinitionDetails = api.ReportDefinitionDetails{Sections: sections}

	return
}

func promptAddReportDefinitionSection(sections *[]api.ReportDefinitionSection) error {
	var policyIDs []string

	cli.StartProgress("Fetching list of policy ids...")
	resp, err := cli.LwApi.V2.Policy.List()
	cli.StopProgress()

	if err != nil {
		return err
	}

	if len(resp.Data) == 0 {
		return errors.New("unable to find policies")
	}

	for _, policy := range resp.Data {
		policyIDs = append(policyIDs, policy.PolicyID)
	}

	questions := []*survey.Question{
		{
			Name:     "title",
			Prompt:   &survey.Input{Message: CreateReportDefinitionSectionTitleQuestion},
			Validate: survey.Required,
		},
		{
			Name:     "policies",
			Prompt:   &survey.MultiSelect{Message: CreateReportDefinitionPoliciesQuestion, Options: policyIDs},
			Validate: survey.MinItems(1),
		},
	}

	answers := struct {
		Title    string   `survey:"title"`
		Policies []string `survey:"policies"`
	}{}

	if err = survey.Ask(questions, &answers, survey.WithIcons(promptIconsFunc)); err != nil {
		return err
	}

	section := api.ReportDefinitionSection{
		Title:    answers.Title,
		Policies: answers.Policies,
	}

	*sections = append(*sections, section)
	return nil
}

func promptCreateReportDefinitionFromExisting() (reportDefinition api.ReportDefinition, err error) {
	var (
		reports        = make(map[string]api.ReportDefinition)
		reportNames    []string
		selectedReport string
	)

	cli.StartProgress("Fetching existing report definitions...")
	resp, err := cli.LwApi.V2.ReportDefinitions.List()
	cli.StopProgress()

	if err != nil {
		return
	}

	for _, report := range resp.Data {
		reports[report.ReportName] = report
		reportNames = append(reportNames, report.ReportName)
	}

	// Add option for blank template
	reports["BLANK TEMPLATE"] = api.ReportDefinition{ReportName: "TEMPLATE",
		ReportType:              api.ReportDefinitionTypeCompliance.String(),
		ReportDefinitionDetails: api.ReportDefinitionDetails{Sections: []api.ReportDefinitionSection{{Title: "CUSTOM SECTION TITLE", Policies: []string{"example-policy-1"}}}}}
	reportNames = append([]string{"BLANK TEMPLATE"}, reportNames...)

	if err = survey.AskOne(&survey.Select{Message: SelectReportDefinitionQuestion, Options: reportNames}, &selectedReport); err != nil {
		return
	}

	reportTemplateYaml, err := yaml.Marshal(reports[selectedReport].Config())
	if err != nil {
		return
	}

	// open editor with report yaml
	report, err := inputReportDefinitionFromEditor("create", string(reportTemplateYaml))
	if err != nil {
		return
	}
	err = yaml.Unmarshal([]byte(report), &reportDefinition)

	return
}

func inputReportDefinitionFromEditor(action string, reportYaml string) (report string, err error) {
	prompt := &survey.Editor{
		Message:       fmt.Sprintf("Use the editor to %s your policy", action),
		FileName:      "report-definition*.yaml",
		HideDefault:   true,
		AppendDefault: true,
		Default:       reportYaml,
	}

	err = survey.AskOne(prompt, &report)
	return
}

func parseNewReportDefinition(s string) (report api.ReportDefinitionConfig, err error) {
	var res api.ReportDefinitionResponse
	if err = json.Unmarshal([]byte(s), &res); err == nil && res.Data.ReportName != "" {
		report = api.ReportDefinitionConfig{
			ReportName:    res.Data.ReportName,
			ReportType:    res.Data.ReportType,
			SubReportType: res.Data.SubReportType,
			DisplayName:   res.Data.DisplayName,
			Sections:      res.Data.ReportDefinitionDetails.Sections,
		}
		return report, nil
	}

	var cfg api.ReportDefinition
	if err = json.Unmarshal([]byte(s), &cfg); err == nil && cfg.ReportName != "" {
		report = api.ReportDefinitionConfig{
			ReportName:    cfg.ReportName,
			ReportType:    cfg.ReportType,
			SubReportType: cfg.SubReportType,
			DisplayName:   cfg.DisplayName,
			Sections:      cfg.ReportDefinitionDetails.Sections,
		}
		return report, nil
	}

	var yamlCfg api.ReportDefinition
	if err = yaml.Unmarshal([]byte(s), &yamlCfg); err == nil && yamlCfg.ReportName != "" {
		report = api.ReportDefinitionConfig{
			ReportName:    yamlCfg.ReportName,
			ReportType:    yamlCfg.ReportType,
			SubReportType: yamlCfg.SubReportType,
			DisplayName:   yamlCfg.DisplayName,
			Sections:      yamlCfg.ReportDefinitionDetails.Sections,
		}
		return report, nil
	}

	return report, errors.New("unable to parse report definition file")
}
