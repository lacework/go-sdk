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

package api_test

import (
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestAlertLevelValid(t *testing.T) {
	assert.True(t, api.AlertLevel(1).Valid())
	assert.True(t, api.AlertLevel(2).Valid())
	assert.True(t, api.AlertLevel(3).Valid())
	assert.True(t, api.AlertLevel(4).Valid())
	assert.True(t, api.AlertLevel(5).Valid())

	assert.False(t, api.AlertLevel(6).Valid())
	assert.False(t, api.AlertLevel(123).Valid())
}

func TestAlertLevelStrings(t *testing.T) {
	assert.Equal(t, "Critical", api.CriticalAlertLevel.String())
	assert.Equal(t, "High", api.HighAlertLevel.String())
	assert.Equal(t, "Medium", api.MediumAlertLevel.String())
	assert.Equal(t, "Low", api.LowAlertLevel.String())
	assert.Equal(t, "All", api.AllAlertLevel.String())
}

func TestAlertLevelInts(t *testing.T) {
	assert.Equal(t, 1, api.CriticalAlertLevel.Int())
	assert.Equal(t, 2, api.HighAlertLevel.Int())
	assert.Equal(t, 3, api.MediumAlertLevel.Int())
	assert.Equal(t, 4, api.LowAlertLevel.Int())
	assert.Equal(t, 5, api.AllAlertLevel.Int())
}
