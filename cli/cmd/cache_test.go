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
