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

package cmd

import (
	"encoding/json"
	"path"
	"strings"

	"github.com/lacework/go-sdk/internal/format"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/peterbourgon/diskv/v3"
)

const MaxCacheSize = 1024 * 1024 * 1024

// InitCache initializes the Lacework CLI cache to store data on disk
//
// Simple CRUD example:
//
// ```go
// cli.Cache.WriteString("data", "something useful")  // Create
// myData := cli.Cache.Read("data")                   // Read
// cli.Cache.Write("data", []byte("cool update"))     // Update
// cli.Cache.Erase("data")                            // Delete
// ```
//
func (c *cliState) InitCache() {
	// @afiune: add a way to disable the cache
	dir, err := cacheDir()
	if err == nil {
		c.Cache = diskv.New(diskv.Options{
			BasePath:          path.Join(dir, "cache"),
			AdvancedTransform: CacheTransform,
			InverseTransform:  InverseCacheTransform,
			CacheSizeMax:      MaxCacheSize,
		})
	}
}

func CacheTransform(key string) *diskv.PathKey {
	// Global cache
	//
	// The Lacework CLI will have times where we need to cache global things
	// such as the daily version checks. For those cases, we will use the
	// global cache that can be accessed like:
	//
	// cli.Cache.Read("global/version")
	//
	if strings.HasPrefix(key, "global/") {
		keys := strings.Split(key, "/")
		return &diskv.PathKey{
			Path:     []string{keys[1]},
			FileName: key,
		}
	}

	// Scoped cache
	//
	// If the Lacework CLI is not using the global cache, then we will land
	// in the scoped cache, this cache is individual per profile, that is,
	// a convination of /account/subaccount/key_id/{file}. This is the default
	// cache location when doing CRUD actions like:
	//
	// cli.Cache.WriteString("static_data", "{ ... some static data ... }")
	//
	subaccount := cli.Subaccount
	if subaccount == "" {
		subaccount = "standalone"
	}
	return &diskv.PathKey{
		Path:     []string{cli.Account, subaccount, cli.KeyID},
		FileName: key,
	}
}

func InverseCacheTransform(pathKey *diskv.PathKey) string {
	if strings.HasPrefix(pathKey.FileName, "global/") {
		keys := strings.Split(pathKey.FileName, "/")
		return strings.Join(pathKey.Path, "/") + keys[1]
	}
	return strings.Join(pathKey.Path, "/") + pathKey.FileName
}

func cacheDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".config", "lacework"), nil
}

func (c *cliState) EraseCachedToken() error {
	if c.noCache {
		return nil
	}

	c.Log.Debugw("token expired, removing from cache",
		"feature", "cache",
	)

	return c.Cache.Erase("token")
}

func (c *cliState) ReadCachedToken() {
	if c.noCache {
		return
	}

	if tokenJSON, err := c.Cache.Read("token"); err == nil {
		if err := json.Unmarshal(tokenJSON, &c.tokenCache); err == nil {
			c.Log.Debugw("token loaded from cache",
				"feature", "cache",
				"token", format.Secret(4, c.tokenCache.Token),
			)
			c.Token = c.tokenCache.Token
		}
	}
}

func (c *cliState) WriteCachedToken() error {
	if c.noCache {
		return nil
	}

	if c.tokenCache.Token == "" {
		response, err := c.LwApi.GenerateToken()
		if err != nil {
			return err
		}

		c.Log.Debugw("saving token",
			"feature", "cache",
			"token", format.Secret(4, response.Token),
			"expires_at", response.ExpiresAt,
		)
		err = c.Cache.Write("token", structToString(response))
		if err != nil {
			c.Log.Warnw("unable to write token in cache",
				"feature", "cache",
				"error", err.Error(),
			)
		}
		c.Token = response.Token
		c.tokenCache.Token = response.Token
		c.tokenCache.ExpiresAt = response.ExpiresAt
	}
	return nil
}

// structToString takes any arbitrary type and converts it into a string
func structToString(v interface{}) []byte {
	out, err := json.Marshal(v)
	if err != nil {
		return []byte{}
	}
	return out
}
