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

package main

import (
	"fmt"
	"os"

	"github.com/lacework/go-sdk/cli/cmd"
	"github.com/lacework/go-sdk/lwcomponent"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if componentError, ok := err.(*lwcomponent.RunError); ok {
			// by default, all our components should display the error
			// to the end user, which is why we don't output it, but we
			// still exit the main program with the exit code from the component
			os.Exit(componentError.ExitCode)
		}

		fmt.Fprintf(os.Stderr, "ERROR %s\n", err)
		os.Exit(1)
	}
}
