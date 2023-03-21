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
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

// update command is used to update a new lacework report definition
var reportDefinitionsUpdateCommand = &cobra.Command{
	Use:   "update <report_definition_id>",
	Short: "Update a report definition",
	Long: `Update an existing custom report definition.

To update an existing report definition:

    lacework report-definition update <report_definition_id>

To update a new report definition from an existing file:

    lacework report-definition update <report_definition_id> --file custom-report.json
`,
	Args: cobra.ExactArgs(1),
	RunE: updateReportDefinition,
}

func updateReportDefinition(_ *cobra.Command, args []string) error {
	var (
		reportDefinition api.ReportDefinitionUpdate
		err              error
	)

	cli.StartProgress("Fetching report definition...")
	existingReport, err := cli.LwApi.V2.ReportDefinitions.Get(args[0])
	cli.StopProgress()

	if err != nil {
		return err
	}

	if existingReport.Data.CreatedBy == "SYSTEM" {
		return errors.New("only user created report definitions can be modified")
	}

	if reportDefinitionsCmdState.File != "" {
		fileInput, err := inputReportDefinitionFromFile(reportDefinitionsCmdState.File)
		if err != nil {
			return err
		}

		cfg, err := parseNewReportDefinition(fileInput)
		if err != nil {
			return err
		}
		reportDefinition = api.NewReportDefinitionUpdate(cfg)
	} else {
		reportDefinition, err = promptUpdateReportDefinition(existingReport.Data)
		if err != nil {
			return err
		}
	}

	cli.StartProgress("Updating report definition...")
	resp, err := cli.LwApi.V2.ReportDefinitions.Update(args[0], reportDefinition)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, "unable to update report definition")
	}

	cli.OutputHuman("The report definition %s was updated. \n", resp.Data.ReportDefinitionGuid)
	return nil
}

func promptUpdateReportDefinition(existingReport api.ReportDefinition) (api.ReportDefinitionUpdate, error) {
	var useEditor bool

	if err := survey.AskOne(&survey.Confirm{Message: UpdateReportDefinitionQuestion}, &useEditor); err != nil {
		return api.ReportDefinitionUpdate{}, err
	}

	if useEditor {
		return promptUpdateReportDefinitionFromEditor(existingReport)
	}

	return promptUpdateReportDefinitionQuestions(existingReport)
}

func promptUpdateReportDefinitionQuestions(existingReport api.ReportDefinition) (reportDefinition api.ReportDefinitionUpdate, err error) {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: UpdateReportDefinitionReportNameQuestion, Default: existingReport.ReportName},
			Validate: survey.Required,
		},
		{
			Name:     "display",
			Prompt:   &survey.Input{Message: UpdateReportDefinitionDisplayNameQuestion, Default: existingReport.DisplayName},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name        string `survey:"name"`
		DisplayName string `survey:"display"`
	}{}

	if err = survey.Ask(questions, &answers, survey.WithIcons(promptIconsFunc)); err != nil {
		return
	}

	reportDefinition = api.ReportDefinitionUpdate{
		ReportName:  answers.Name,
		DisplayName: answers.DisplayName,
	}

	sections, err := promptUpdateReportDefinitionSections(existingReport)
	if err != nil {
		return
	}

	reportDefinition.ReportDefinitionDetails = &api.ReportDefinitionDetails{Sections: sections}

	return
}

func promptUpdateReportDefinitionSections(report api.ReportDefinition) ([]api.ReportDefinitionSection, error) {
	var (
		sections    = report.ReportDefinitionDetails.Sections
		newSections []api.ReportDefinitionSection
		err         error
	)

	cli.StartProgress("Fetching list of policy ids...")
	resp, err := cli.LwApi.V2.Policy.List()
	cli.StopProgress()

	if err != nil {
		return nil, err
	}

	//filter the policies not in current report's domain
	var policies []api.Policy
	for _, p := range resp.Data {
		domain := strings.ToUpper(report.SubReportType)
		if slices.Contains(p.Tags, fmt.Sprintf("domain:%s", domain)) {
			policies = append(policies, p)
		}
	}

	if len(policies) == 0 {
		return nil, errors.New("unable to find policies")
	}

	// edit existing sections
	editSection := false
	if err := survey.AskOne(&survey.Confirm{
		Message: UpdateReportDefinitionEditSectionQuestion,
	}, &editSection); err != nil {
		return nil, err
	}

	if editSection {
		if sections, err = promptUpdateReportDefinitionSection(&sections, policies); err != nil {
			return nil, err
		}

		editAnotherSection := false
		for {
			if err := survey.AskOne(&survey.Confirm{
				Message: UpdateReportDefinitionEditAnotherSectionQuestion,
			}, &editAnotherSection); err != nil {
				return nil, err
			}

			if editAnotherSection {
				if sections, err = promptUpdateReportDefinitionSection(&newSections, policies); err != nil {
					return nil, err
				}
			} else {
				break
			}
		}
	}

	// add new sections
	addSection := false
	if err := survey.AskOne(&survey.Confirm{
		Message: UpdateReportDefinitionAddSectionQuestion,
	}, &addSection); err != nil {
		return nil, err
	}

	if addSection {
		if err := promptAddReportDefinitionSection(&newSections, policies); err != nil {
			return nil, err
		}

		addAnotherSection := false
		for {
			if err := survey.AskOne(&survey.Confirm{
				Message: CreateReportDefinitionAddSectionQuestion,
			}, &addAnotherSection); err != nil {
				return nil, err
			}

			if addAnotherSection {
				if err := promptAddReportDefinitionSection(&newSections, policies); err != nil {
					return nil, err
				}
			} else {
				break
			}
		}
	}

	sections = append(sections, newSections...)

	return sections, nil

}

func promptUpdateReportDefinitionFromEditor(existingReport api.ReportDefinition) (reportDefinition api.ReportDefinitionUpdate, err error) {
	var reportDefinitionConfig api.ReportDefinitionConfig

	if err != nil {
		return
	}

	updateCfg := api.NewReportDefinitionUpdate(existingReport.Config())
	reportTemplateYaml, err := yaml.Marshal(updateCfg)
	if err != nil {
		return
	}

	// open editor with report yaml
	report, err := inputReportDefinitionFromEditor("update", string(reportTemplateYaml))
	if err != nil {
		return
	}
	err = yaml.Unmarshal([]byte(report), &reportDefinitionConfig)

	reportDefinition = api.NewReportDefinitionUpdate(reportDefinitionConfig)

	return
}

func promptUpdateReportDefinitionSection(currentSections *[]api.ReportDefinitionSection, policies []api.Policy) ([]api.ReportDefinitionSection, error) {
	type sectionMapping struct {
		section  api.ReportDefinitionSection
		position int
	}

	var (
		policyIDs        []string
		selectedSections = make(map[string]sectionMapping)
		sectionTitles    []string
		sections         = *currentSections
	)

	for _, policy := range policies {
		policyIDs = append(policyIDs, policy.PolicyID)
	}

	for i, section := range sections {
		sectionTitles = append(sectionTitles, section.Title)
		selectedSections[section.Title] = sectionMapping{section: section, position: i}
	}

	var selectedTitle string
	if err := survey.AskOne(&survey.Select{
		Message: UpdateReportDefinitionSelectSectionQuestion,
		Options: sectionTitles,
	}, &selectedTitle); err != nil {
		return nil, err
	}

	selectedSection := selectedSections[selectedTitle]

	questions := []*survey.Question{
		{
			Name:     "title",
			Prompt:   &survey.Input{Message: CreateReportDefinitionSectionTitleQuestion, Default: selectedSection.section.Title},
			Validate: survey.Required,
		},
		{
			Name:     "policies",
			Prompt:   &survey.MultiSelect{Message: CreateReportDefinitionPoliciesQuestion, Options: policyIDs, Default: selectedSection.section.Policies},
			Validate: survey.MinItems(1),
		},
	}

	answers := struct {
		Title    string   `survey:"title"`
		Policies []string `survey:"policies"`
	}{}

	if err := survey.Ask(questions, &answers, survey.WithIcons(promptIconsFunc)); err != nil {
		return nil, err
	}

	updatedSection := api.ReportDefinitionSection{
		Title:    answers.Title,
		Policies: answers.Policies,
	}

	sections[selectedSection.position] = updatedSection
	return sections, nil
}
