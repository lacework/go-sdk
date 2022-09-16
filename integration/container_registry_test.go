//go:build container_registry

// Author:: Darren Murray (<darren.murray@lacework.net>)
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerRegistryShow(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("container-registry", "show", "TECHALLY_7891701909B297518DB6E1C7E993E706B4C1BE1641A4EE3")
	// Summary Table
	assert.Contains(t, out.String(), "CONTAINER REGISTRY GUID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "NAME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "TYPE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATE",
		"STDOUT table headers changed, please check")

	// Details Table
	assert.Contains(t, out.String(), "REGISTRY TYPE",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "REGISTRY TYPE",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "PRIVATE KEY ID",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "PRIVATE KEY",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "CLIENT EMAIL",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "REGISTRY DOMAIN",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "REGISTRY DOMAIN",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "UPDATED AT",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "UPDATED BY",
		"STDOUT details headers changed, please check")
	assert.Contains(t, out.String(), "LAST SUCCESSFUL STATE",
		"STDOUT details headers changed, please check")

	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}

func TestContainerRegistryList(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("cr", "list")
	assert.Contains(t, out.String(), "CONTAINER REGISTRY GUID",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "NAME",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "TYPE",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATUS",
		"STDOUT table headers changed, please check")
	assert.Contains(t, out.String(), "STATE",
		"STDOUT table headers changed, please check")
	assert.Empty(t,
		err.String(),
		"STDERR should be empty")
	assert.Equal(t, 0, exitcode,
		"EXITCODE is not the expected one")
}
