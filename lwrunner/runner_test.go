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

package lwrunner_test

import (
	"testing"

	"github.com/lacework/go-sdk/lwrunner"

	"github.com/stretchr/testify/assert"
)

func TestLwRunnerNew(t *testing.T) {
	subject := lwrunner.New("ubuntu", "my-test-host")
	assert.Equal(t, subject.Port, 22)
	assert.Equal(t, subject.User, "ubuntu")
}
