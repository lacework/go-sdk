//
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

package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestEntities_Images_Search(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/Images/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockImagesResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response := api.ImagesEntityResponse{}
	err = c.V2.Entities.Search(&response, api.SearchFilter{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 4, len(response.Data)) {
		assert.Equal(t, "gcr.io/techally-hipstershop-275821/shippingservice", response.Data[0].Repo)
	}
}

func TestEntities_Images_List(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/Images/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockImagesResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Entities.ListImages()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	// only one since all pages would be two
	if assert.Equal(t, 4, len(response.Data)) {
		assert.Equal(t, 120, response.Data[0].Mid)
		assert.Equal(t,
			"sha256:7c9d94b8d689bf7a7b9b669bbeabc90d9f40ed517e62bd3b0fc3ffcdb6151961",
			response.Data[0].ImageID)
		assert.Empty(t, response.Paging.Urls.NextPage)
	}
}

func TestEntities_Images_List_All(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/Images/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockImagesResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Entities.ListAllImages()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	if assert.Equal(t, 4, len(response.Data)) {
		assert.Equal(t, 120, response.Data[0].Mid)
		assert.Equal(t,
			"sha256:7c9d94b8d689bf7a7b9b669bbeabc90d9f40ed517e62bd3b0fc3ffcdb6151961",
			response.Data[0].ImageID)
		assert.Equal(t, 36, response.Data[1].Mid)
		assert.Equal(t,
			"sha256:1efcf523c90b0268ffdf05b7a73ab0007332067faaffd816ff9f8e733063d889",
			response.Data[1].ImageID)
		assert.Empty(t, response.Paging, "paging should be empty")
	}
}

func TestEntities_Images_ListAll_EmptyData(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockToken("TOKEN")
	fakeServer.MockAPI("Entities/Images/search",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method, "Search() should be a POST method")
			fmt.Fprintf(w, mockPaginationEmptyResponse())
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithApiV2(),
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.NoError(t, err)

	response, err := c.V2.Entities.ListAllImages()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, len(response.Data))
}

func mockImagesResponse() string {
	return `
{
  "data": [
    {
      "containerType": "DOCKER",
      "createdTime": "2022-02-09T18:37:34.959Z",
      "imageId": "sha256:7c9d94b8d689bf7a7b9b669bbeabc90d9f40ed517e62bd3b0fc3ffcdb6151961",
      "mid": 120,
      "repo": "gcr.io/techally-hipstershop-275821/shippingservice",
      "size": 35018063,
      "tag": "sha256:940db804030f211f6f780a941e4ccaecab8f650b971f6e45d48ea0590fad8b2b"
    },
    {
      "containerType": "DOCKER",
      "createdTime": "2022-02-09T18:38:45.953Z",
      "imageId": "sha256:1efcf523c90b0268ffdf05b7a73ab0007332067faaffd816ff9f8e733063d889",
      "mid": 36,
      "repo": "123456789012.dkr.ecr.us-east-2.amazonaws.com/amazon-k8s-test",
      "size": 312076970,
      "tag": "v1.7.5"
    },
    {
      "containerType": "DOCKER",
      "createdTime": "2022-02-09T18:52:34.974Z",
      "imageId": "sha256:10d28bedfe5dec59da9ebf8e6260224ac9008ab5c11dbbe16ee3ba3e4439ac2c",
      "mid": 120,
      "repo": "k8s.gcr.io/awesome",
      "size": 682696,
      "tag": "3.2"
    },
    {
      "containerType": "DOCKER",
      "createdTime": "2022-02-09T19:19:46.176Z",
      "imageId": "sha256:190162514693718242e49d7ed7e90cf1597554c05d16b8407600902a7bd33831",
      "mid": 122,
      "repo": "gcr.io/techally-hipstershop-275821/currencyservice",
      "size": 160124851,
      "tag": "sha256:1e57cf3949fb02025105b51ea8b36e6c7b805208652117a03da6b2bbc78a765a"
    }
  ],
  "paging": {
    "rows": 4,
    "totalRows": 4,
    "urls": {
      "nextPage": null
    }
  }
}
	`
}
