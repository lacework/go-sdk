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
	"os"

	"github.com/fatih/color"
)

// OutputJSON will print out the JSON representation of the provided data
func (c *cliState) OutputJSON(v interface{}) error {
	pretty, err := c.JsonF.Marshal(v)
	if err != nil {
		cli.Log.Debugw("unable to pretty print JSON object", "raw", v)
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
