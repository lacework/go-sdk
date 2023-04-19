//
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

package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestInventorySearch(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Inventory/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, `{"data":[]}`)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response := api.InventoryAwsResponse{}
	err = c.V2.Inventory.Search(&response, api.InventorySearch{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, len(response.Data))
}

func TestInventorySearchInvalid(t *testing.T) {
	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
	)
	assert.NoError(t, err)

	type MockInvalidInventoryResponse struct{}
	response := MockInvalidInventoryResponse{}
	err = c.V2.Inventory.Search(&response, api.InventorySearch{})
	assert.Error(t, err)
}
