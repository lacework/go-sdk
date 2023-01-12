//go:build query && !windows

// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

package integration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
)

func _TestQueryUpdateFromIDEditor(t *testing.T) { // TODO reenable
	// get temp file
	file, err := createTemporaryFile(
		"TestQueryUpdateFromIDEditor",
		fmt.Sprintf(queryJSONTemplate, queryID, queryUpdateText),
	)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	dir := createTOMLConfigFromCIvars()
	defer os.RemoveAll(dir)

	// setup
	LaceworkCLIWithHome(dir, "query", "create", "-u", queryURL)
	// teardown
	defer LaceworkCLIWithHome(dir, "query", "delete", queryID)

	_ = runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			expectString(t, c, "Update query")
			c.SendLine("")
			time.Sleep(time.Millisecond)
			// Move to line number 4 and add comment "--- Updated from CLI Editor"
			c.Send("4Go--- Updated from CLI Editor\x1b")
			c.SendLine(":wq!") // save and close
			time.Sleep(time.Millisecond)
			expectString(t, c,
				fmt.Sprintf("The query %s was updated.", queryID))
		},
		"query", "update", queryID,
	)

	t.Run("verify query editions", func(t *testing.T) {
		stdout, stderr, exitcode := LaceworkCLIWithHome(dir, "query", "show", queryID)
		assert.Empty(t,
			stderr.String(),
			"STDERR should be empty")
		assert.Contains(t,
			stdout.String(),
			"--- Updated from CLI Editor",
			"the query was not editted correctly")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
	})
}
