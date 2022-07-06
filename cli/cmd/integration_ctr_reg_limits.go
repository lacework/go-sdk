//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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
	"bytes"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func castMapStringSliceToString(m []map[string]string) string {
	b := new(bytes.Buffer)
	for _, s := range m {
		for key, value := range s {
			fmt.Fprintf(b, "%s=%s, ", key, value)
		}
	}
	return b.String()
}

func castStringToLimitByLabel(labels string) []map[string]string {
	out := make([]map[string]string, 0)

	for _, label := range strings.Split(labels, "\n") {
		kv := strings.Split(label, ":")
		if len(kv) != 2 {
			cli.Log.Warnw("malformed limit_by_label entry, ignoring",
				"label", label,
				"expected_format", "key:value",
			)
			continue
		}
		out = append(out, map[string]string{kv[0]: kv[1]})
	}

	return out
}

func askV2LimitByTags(answers interface{}) error {
	custom := false
	if err := survey.AskOne(&survey.Confirm{
		Message: "Configure limit of scans by tags?",
	}, &custom); err != nil {
		return err
	}

	if custom {
		questions := []*survey.Question{{
			Name:   "limit_tags",
			Prompt: &survey.Multiline{Message: "List of tags to scan:"},
		}}

		if err := survey.Ask(questions, answers,
			survey.WithIcons(promptIconsFunc),
		); err != nil {
			return err
		}
	}

	return nil
}

func askV2LimitByLabels(answers interface{}) error {
	custom := false
	if err := survey.AskOne(&survey.Confirm{
		Message: "Configure limit of scans by labels?",
	}, &custom); err != nil {
		return err
	}

	if custom {
		questions := []*survey.Question{{
			Name:   "limit_labels",
			Prompt: &survey.Multiline{Message: "List of 'key:value' labels to scan:"},
		}}

		if err := survey.Ask(questions, answers,
			survey.WithIcons(promptIconsFunc),
		); err != nil {
			return err
		}
	}

	return nil
}

func askV2LimitByRepositories(answers interface{}) error {
	custom := false
	if err := survey.AskOne(&survey.Confirm{
		Message: "Configure limit of scans by repositories?",
	}, &custom); err != nil {
		return err
	}

	if custom {
		questions := []*survey.Question{{
			Name:   "limit_repos",
			Prompt: &survey.Multiline{Message: "List of repositories to scan:"},
		}}

		if err := survey.Ask(questions, answers,
			survey.WithIcons(promptIconsFunc),
		); err != nil {
			return err
		}
	}

	return nil
}

func askV2Limits(answers interface{}) error {
	if err := askV2LimitByTags(answers); err != nil {
		return err
	}
	if err := askV2LimitByLabels(answers); err != nil {
		return err
	}
	return askV2LimitByRepositories(answers)
}
