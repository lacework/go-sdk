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

package lwtime_test

import (
	"testing"
	"time"

	"github.com/lacework/go-sdk/lwtime"
	"github.com/stretchr/testify/assert"
)

func TestParseRelativeOK(t *testing.T) {
	rt, err := lwtime.ParseRelative("now")
	assert.Nil(t, err)

	dur := rt.Unix() - time.Now().Unix()
	assert.LessOrEqual(t, dur, int64(2))
}

func TestParseRelativeErr(t *testing.T) {
	_, err := lwtime.ParseRelative("jackie weaver")
	assert.NotNil(t, err)
}
