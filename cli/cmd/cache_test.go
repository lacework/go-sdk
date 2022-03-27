//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

package cmd

import (
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheGlobal(t *testing.T) {
	// create a temporal directory for our global cache
	dir, err := ioutil.TempDir("", "lacework-cli-cache")
	if err != nil {
		panic(err)
	}
	cli.InitCache(dir)

	defer func() {
		cli.Cache = nil
	}()

	key := "global/file"
	expected := []byte("data")

	err = cli.Cache.Write(key, expected)
	assert.Nil(t, err)
	cache, err := cli.Cache.Read(key)
	assert.Nil(t, err)
	assert.Equal(t, expected, cache)

	// verify global cache location
	assert.FileExists(t, path.Join(dir, "cache", "file"))
}

func TestCacheGlobalWithLongPath(t *testing.T) {
	// create a temporal directory for our global cache
	dir, err := ioutil.TempDir("", "lacework-cli-cache")
	if err != nil {
		panic(err)
	}
	cli.InitCache(dir)

	defer func() {
		cli.Cache = nil
	}()

	key := "global/path/to/file"
	expected := []byte("data")

	err = cli.Cache.Write(key, expected)
	assert.Nil(t, err)
	cache, err := cli.Cache.Read(key)
	assert.Nil(t, err)
	assert.Equal(t, expected, cache)

	// verify global cache location
	assert.FileExists(t, path.Join(dir, "cache", "path", "to", "file"))
}

func TestCacheScopedStandalone(t *testing.T) {
	// create a temporal directory for our global cache
	dir, err := ioutil.TempDir("", "lacework-cli-cache")
	if err != nil {
		panic(err)
	}
	cli.InitCache(dir)

	// setting up required Lacework CLI config
	cli.Account = "account"
	cli.Subaccount = "subaccount"
	cli.KeyID = "key"

	defer func() {
		cli.Cache = nil
		cli.Account = ""
		cli.KeyID = ""
		cli.Subaccount = ""
	}()

	key := "compliance_report"
	expected := []byte("data")

	err = cli.Cache.Write(key, expected)
	assert.Nil(t, err)
	cache, err := cli.Cache.Read(key)
	assert.Nil(t, err)
	assert.Equal(t, expected, cache)

	// verify global cache location
	assert.FileExists(t, path.Join(dir, "cache", "account", "subaccount", "key", "compliance_report"))
}

func TestCacheScopedOrgAccounts(t *testing.T) {
	// create a temporal directory for our global cache
	dir, err := ioutil.TempDir("", "lacework-cli-cache")
	if err != nil {
		panic(err)
	}
	cli.InitCache(dir)

	// setting up required Lacework CLI config
	cli.Account = "account"
	cli.KeyID = "key"

	defer func() {
		cli.Cache = nil
		cli.Account = ""
		cli.KeyID = ""
	}()

	key := "vuln_assessment"
	expected := []byte("data")

	err = cli.Cache.Write(key, expected)
	assert.Nil(t, err)
	cache, err := cli.Cache.Read(key)
	assert.Nil(t, err)
	assert.Equal(t, expected, cache)

	// verify global cache location
	assert.FileExists(t, path.Join(dir, "cache", "account", "standalone", "key", "vuln_assessment"))
}

func TestCacheEndToEnd(t *testing.T) {
	// create a temporal directory for our global cache
	dir, err := ioutil.TempDir("", "lacework-cli-cache")
	if err != nil {
		panic(err)
	}
	cli.InitCache(dir)

	defer func() {
		cli.Cache = nil
	}()

	key := "test"
	expected := []byte("data")

	// Create
	err = cli.Cache.Write(key, expected)
	assert.Nil(t, err)

	// Read
	cached, err := cli.Cache.Read(key)
	assert.Nil(t, err)
	assert.Equal(t, expected, cached)

	// Update
	expectedUpdate := []byte("better data")
	err = cli.Cache.Write(key, expectedUpdate)
	assert.Nil(t, err)
	newCache, err := cli.Cache.Read(key)
	assert.Nil(t, err)
	assert.Equal(t, expectedUpdate, newCache)
	assert.NotEqual(t, cached, newCache)

	// Delete
	err = cli.Cache.Erase(key)
	assert.Nil(t, err)
	err = cli.Cache.Erase(key)
	assert.NotNil(t, err)
}

func TestWriteReadAssetToCache(t *testing.T) {
	// create a temporal directory for our global cache
	dir, err := ioutil.TempDir("", "lacework-cli-cache")
	if err != nil {
		panic(err)
	}
	cli.InitCache(dir)

	t.Run("primitive type: string", func(t *testing.T) {
		key := "mocks/string_value"
		expected := "super cool value (string)"

		// write mock asset with an expiration time of 1 hour
		cli.WriteAssetToCache(key, time.Now().Add(time.Hour*1), expected)
		assert.FileExists(t, path.Join(dir, "cache", "standalone", key))

		// read asset
		var value string
		expired := cli.ReadCachedAsset(key, &value)
		if assert.False(t, expired) {
			assert.Equal(t, expected, value)
		}
	})

	t.Run("primitive type: bool", func(t *testing.T) {
		key := "mocks/bool_value"
		expected := true

		// write mock asset with an expiration time of 5 seconds
		cli.WriteAssetToCache(key, time.Now().Add(time.Second*5), expected)
		assert.FileExists(t, path.Join(dir, "cache", "standalone", key))

		// read asset
		var value bool
		expired := cli.ReadCachedAsset(key, &value)
		if assert.False(t, expired) {
			assert.Equal(t, expected, value)
		}
	})

	t.Run("JSON object", func(t *testing.T) {
		type myMockAsset struct {
			Name      string `json:"name"`
			Important bool   `json:"important"`
		}
		key := "mock_asset"
		expiresAt := time.Now().Add(time.Minute * 1)
		expected := myMockAsset{"foo", true}

		// write mock asset with an expiration time of 1 minute
		cli.WriteAssetToCache(key, expiresAt, expected)
		assert.FileExists(t, path.Join(dir, "cache", "standalone", key))

		// read asset
		var asset myMockAsset
		expired := cli.ReadCachedAsset(key, &asset)
		if assert.False(t, expired) {
			assert.Equal(t, expected, asset)
		}
	})

	t.Run("expired cache", func(t *testing.T) {
		key := "mocked_expired_should_not_exist"
		notExpected := 123

		// write mock asset with an expiration time of NOW!
		cli.WriteAssetToCache(key, time.Now(), notExpected)
		assert.NoFileExists(t, path.Join(dir, "cache", "standalone", key))

		// read asset
		var value int
		expired := cli.ReadCachedAsset(key, &value)
		assert.True(t, expired)
		assert.NotEqual(t, notExpected, value)
	})

	t.Run("time before NOW should NOT write asset", func(t *testing.T) {
		key := "mocked_data"
		notExpected := "foo"

		// try to write mock asset with an expiration time before NOW!
		cli.WriteAssetToCache(key, time.Now().Add(time.Duration(-1)*time.Minute), notExpected)
		assert.NoFileExists(t, path.Join(dir, "cache", "standalone", key))

		// read asset
		var value string
		expired := cli.ReadCachedAsset(key, &value)
		assert.True(t, expired)
		assert.NotEqual(t, notExpected, value)
	})
}
