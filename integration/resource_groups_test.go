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

package integration

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func createResourceGroup(typ string) (string, error) {
	iType, _ := api.FindResourceGroupType(typ)

	var testQuery api.RGQuery
	err := json.Unmarshal([]byte(iType.QueryTemplate()), &testQuery)
	if err != nil {
		return "", errors.Wrap(err, "error serializing query template")
	}

	var testResourceGroup api.ResourceGroupDataWithQuery = api.ResourceGroupDataWithQuery{
		Name:        fmt.Sprintf("CLI_TestCreateResourceGroup_%s", iType.String()),
		Type:        iType.String(),
		Query:       &testQuery,
		Description: "Resource Group Created By CLI Integration Testing",
		Enabled:     1,
	}

	testResourceGroupBytes, err := json.Marshal(testResourceGroup)
	if err != nil {
		return "", errors.Wrap(err, "error marshaling test resource group")
	}

	out, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"api", "post", "v2/ResourceGroups", "-d", string(testResourceGroupBytes), "--json",
	)
	if stderr.String() != "" {
		return "", errors.New(stderr.String())
	}
	if exitcode != 0 {
		return "", errors.New("non-zero exit code")
	}

	var resourceGroupV2Response api.ResourceGroupV2Response
	err = json.Unmarshal(out.Bytes(), &resourceGroupV2Response)
	if err != nil {
		return "", err
	}
	return resourceGroupV2Response.Data.ID(), nil
}

func popResourceGroup() (string, error) {
	type resourceGroup struct {
		Id string `json:"resource_guid"`
	}
	type listResourceGroupsResponse struct {
		ResourceGroups []resourceGroup `json:"resource_groups"`
	}

	var resourceGroups listResourceGroupsResponse

	out, stderr, exitcode := LaceworkCLIWithTOMLConfig(
		"resource-group", "list", "--json", "--nocache",
	)
	if stderr.String() != "" {
		return "", errors.New(stderr.String())
	}
	if exitcode != 0 {
		return "", errors.New("non-zero exit code")
	}

	err := json.Unmarshal(out.Bytes(), &resourceGroups)
	if err != nil {
		return "", err
	}

	for _, rg := range resourceGroups.ResourceGroups {
		return rg.Id, nil
	}
	return "", errors.New("no resource groups found")
}

func TestResourceGroupCreateEditor(t *testing.T) {
	// create
	out, err, exitcode := LaceworkCLIWithTOMLConfig("resource-group", "create")
	assert.Contains(t, out.String(), "Choose a resource group type to create")
	assert.Contains(t, out.String(), "[Use arrows to move, type to filter]")
	assert.Contains(t, err.String(), "ERROR unable to create resource group:")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestResourceGroupList(t *testing.T) {
	// list (output human)
	out, err, exitcode := LaceworkCLIWithTOMLConfig("resource-group", "list")
	assert.Contains(t, out.String(), "RESOURCE GROUP ID")
	assert.Contains(t, out.String(), "All Aws Accounts")
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")

	// list (output json)
	out, err, exitcode = LaceworkCLIWithTOMLConfig("resource-group", "list", "--json")
	assert.Contains(t, out.String(), `"resource_guid"`)
	assert.Contains(t, out.String(), `"type": "AWS"`)
	assert.Empty(t, err.String(), "STDERR should be empty")
	assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
}

func TestPolicyShowNoInput(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("resource-group", "show")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR accepts 1 arg(s), received 0")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestResourceGroupShow(t *testing.T) {
	resourceGroupShowID, err := popResourceGroup()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	t.Run("Human Output", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithTOMLConfig("resource-group", "show", resourceGroupShowID)
		assert.Contains(t, out.String(), "RESOURCE GROUP ID")
		assert.Contains(t, out.String(), resourceGroupShowID)
		assert.Empty(t, err.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	})

	t.Run("JSON Output", func(t *testing.T) {
		out, err, exitcode := LaceworkCLIWithTOMLConfig("resource-group", "show", resourceGroupShowID, "--json")
		assert.Contains(t, out.String(), `"resource_guid"`)
		assert.Contains(t, out.String(), `"`+resourceGroupShowID+`"`)
		assert.Empty(t, err.String(), "STDERR should be empty")
		assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
	})
}

func TestResourceGroupDeleteNoInput(t *testing.T) {
	out, err, exitcode := LaceworkCLIWithTOMLConfig("resource-group", "delete")
	assert.Empty(t, out.String(), "STDOUT should be empty")
	assert.Contains(t, err.String(), "ERROR accepts 1 arg(s), received 0")
	assert.Equal(t, 1, exitcode, "EXITCODE is not the expected one")
}

func TestResourceGroupDelete(t *testing.T) {

	// test each RGv2 type against its default template
	for i := range api.ResourceGroupTypes {
		switch i {
		case api.NoneResourceGroup, api.LwAccountResourceGroup:
			// these resource groups are not applicable
			continue
		default:
			// skip lw_account
			t.Run(i.String(), func(t *testing.T) {
				// setup resource group
				resourceGroupID, err := createResourceGroup(i.String())
				if err != nil && !strings.Contains(err.Error(), "already exists in the account") {
					assert.FailNow(t, err.Error())
				}

				// delete resource group
				out, stderr, exitcode := LaceworkCLIWithTOMLConfig("resource-group", "delete", resourceGroupID)
				assert.Contains(t, out.String(), "The resource group was deleted.")
				assert.Empty(t, stderr.String(), "STDERR should be empty")
				assert.Equal(t, 0, exitcode, "EXITCODE is not the expected one")
			})
		}
	}

}
