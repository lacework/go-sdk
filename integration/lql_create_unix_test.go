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

func TestQueryCreateFromEditor(t *testing.T) {
	dir := createTOMLConfigFromCIvars()
	defer os.RemoveAll(dir)

	// teardown
	defer LaceworkCLIWithHome(dir, "query", "delete", queryID)

	_ = runFakeTerminalTestFromDir(t, dir,
		func(c *expect.Console) {
			time.Sleep(time.Millisecond)
			expectString(t, c, "Type a query to create")
			c.SendLine("")
			time.Sleep(time.Millisecond)
			// Replace first like with query id
			c.Send(fmt.Sprintf("c$queryId: %s\x1b", queryID))
			time.Sleep(time.Millisecond)
			// Move to line number 5 and add source
			c.Send("5Gc$CloudTrailRawEvents\x1b")
			time.Sleep(time.Millisecond)
			// Move to line number 8 and add filter
			c.Send("8Gc$EVENT_NAME like 'GetBucket%'\x1b")
			time.Sleep(time.Millisecond)
			// Move to line number 11 and add return
			c.Send("11Gc$INSERT_ID\x1b")
			time.Sleep(time.Millisecond)
			c.SendLine(":wq!") // save and close
			time.Sleep(time.Millisecond)
			expectString(t, c,
				fmt.Sprintf("The query %s was created.", queryID))
		},
		"query", "create",
	)

	t.Run("verify query editions", func(t *testing.T) {
		stdout, stderr, exitcode := LaceworkCLIWithHome(dir, "query", "show", queryID)
		assert.Empty(t,
			stderr.String(),
			"STDERR should be empty")
		assert.Contains(t, stdout.String(), "CloudTrailRawEvents", "STDOUT changed, please update")
		assert.Contains(t, stdout.String(), "EVENT_NAME like 'GetBucket%'", "STDOUT changed, please update")
		assert.Contains(t, stdout.String(), "INSERT_ID", "STDOUT changed, please update")
		assert.Equal(t, 0, exitcode,
			"EXITCODE is not the expected one")
	})
}
