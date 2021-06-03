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
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// OutputJSON will print out the JSON representation of the provided data
func (c *cliState) OutputJSON(v interface{}) error {
	pretty, err := c.JsonF.Marshal(v)
	if err != nil {
		c.Log.Debugw("unable to pretty print JSON object", "raw", v)
		return err
	}
	fmt.Fprintln(color.Output, string(pretty))
	return nil
}

// OutputHumanRead will print out the provided message if the cli state is
// configured to talk to humans, to switch to json format use --json
func (c *cliState) OutputHuman(format string, a ...interface{}) {
	if c.HumanOutput() {
		fmt.Fprintf(os.Stdout, format, a...)
	}
}

// OutputJSONString is just like OutputJSON but from a JSON string
func (c *cliState) OutputJSONString(s string) error {
	pretty, err := c.FormatJSONString(s)
	if err != nil {
		return err
	}
	fmt.Fprintln(color.Output, string(pretty))
	return nil
}

// FormatJSONString formats a JSON string into a pretty JSON format
func (c *cliState) FormatJSONString(s string) (string, error) {
	pretty, err := c.JsonF.Format([]byte(strings.Trim(s, "'")))
	if err != nil {
		c.Log.Debugw("unable to pretty print JSON string", "raw", s)
		return "", err
	}
	return string(pretty), nil
}

// Used to clean CSV inputs prior to rendering
func csvCleanData(input []string) []string {
	var data []string
	for _, h := range input {
		data = append(data, strings.Replace(h, "\n", "", -1))
	}
	return data
}

// Used to produce CSV output
func renderAsCSV(headers []string, data [][]string) (error, string) {
	csvOut := &strings.Builder{}
	csv := csv.NewWriter(csvOut)

	if len(headers) > 0 {
		if err := csv.Write(csvCleanData(headers)); err != nil {
			return errors.Wrap(err, "Failed to build csv output"), ""
		}
	}

	for _, record := range data {
		if err := csv.Write(csvCleanData(record)); err != nil {
			return errors.Wrap(err, "Failed to build csv output"), ""
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	csv.Flush()
	return nil, csvOut.String()
}

// OutputCSV will print out the provided headers/data in CSV format if the cli state is
// configured to output CSV
func (c *cliState) OutputCSV(headers []string, data [][]string) error {
	err, csv := renderAsCSV(headers, data)
	if err != nil {
		return err
	}

	fmt.Print(csv)
	return nil
}
