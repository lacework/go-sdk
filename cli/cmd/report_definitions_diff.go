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
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// reportDefinitionsDiffCommand command is used to compare 2 lacework report definition versions
var reportDefinitionsDiffCommand = &cobra.Command{
	Use:   "diff <report_definition_id> <version> <version>",
	Short: "Compare two versions of a report definition",
	Long: `Compare two versions of a report definition.

To see a diff of two report definition versions:

    lacework report-definition diff <report_definition_id> <current_version> <new_version>
`,
	Args: cobra.ExactArgs(3),
	RunE: diffReportDefinition,
}

func diffReportDefinition(_ *cobra.Command, args []string) error {
	var (
		err        error
		versionOne int
		versionTwo int
		reportOne  *diffCfg
		reportTwo  *diffCfg
	)

	if versionOne, err = strconv.Atoi(args[1]); err != nil {
		return errors.Wrap(err, "unable to parse version")
	}

	if versionTwo, err = strconv.Atoi(args[2]); err != nil {
		return errors.Wrap(err, "unable to parse version")
	}

	cli.StartProgress("Fetching all report definition versions...")
	response, err := cli.LwApi.V2.ReportDefinitions.GetVersions(args[0])
	cli.StopProgress()

	if err != nil {
		return err
	}

	for _, r := range response.Data {
		if r.Version == versionOne {
			reportOne = &diffCfg{
				name:   fmt.Sprintf("Version-%d", r.Version),
				object: r,
			}
		}
		if r.Version == versionTwo {
			reportTwo = &diffCfg{
				name:   fmt.Sprintf("Version-%d", r.Version),
				object: r,
			}
		}
	}

	if reportOne == nil || reportTwo == nil {
		return errors.New("unable to find report definition versions")
	}

	diff, err := diffAsYamlString(*reportOne, *reportTwo)
	if err != nil {
		return err
	}

	cli.OutputHuman(diff)
	return nil
}

type diffCfg struct {
	name   string
	object any
}

func diffAsYamlString(objectOne, objectTwo diffCfg) (string, error) {
	yamlBytesOne, err := yaml.Marshal(objectOne.object)
	if err != nil {
		return "", err
	}

	yamlBytesTwo, err := yaml.Marshal(objectTwo.object)
	if err != nil {
		return "", err
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(yamlBytesOne)),
		B:        difflib.SplitLines(string(yamlBytesTwo)),
		FromFile: objectOne.name,
		ToFile:   objectTwo.name,
		Context:  3,
	}

	diffText, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", err
	}

	output := prettyPrintDiff(diffText)
	return output, nil
}

func prettyPrintDiff(diff string) string {
	if diff == "" {
		return ""
	}

	var sb = &strings.Builder{}
	lines := strings.Split(diff, "\n")

	//colourize lines in diff
	for _, s := range lines {
		if strings.HasPrefix(s, "+") {
			addition := color.HiGreenString(fmt.Sprintf("%s\n", s))
			sb.WriteString(addition)
			continue
		}

		if strings.HasPrefix(s, "-") {
			subtraction := color.HiRedString(fmt.Sprintf("%s\n", s))
			sb.WriteString(subtraction)
			continue
		}

		if strings.HasPrefix(s, "@@") {
			lineDiff := color.HiBlueString(fmt.Sprintf("%s\n", s))
			sb.WriteString(lineDiff)
			continue
		}

		sb.WriteString(fmt.Sprintf("%s\n", s))

	}

	return sb.String()
}
