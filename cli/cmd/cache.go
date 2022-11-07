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
	"time"

	"github.com/lacework/go-sdk/internal/cache"
	"github.com/lacework/go-sdk/internal/format"
	"github.com/peterbourgon/diskv/v3"
)

const MaxCacheSize = 1024 * 1024 * 1024

// InitCache initializes the Lacework CLI cache to store data on disk,
// this functions accepts an specific path to store the cache or, by
// default, it will use the default location.
//
// Simple CRUD example:
//
// ```go
// cli.Cache.WriteString("data", "something useful")  // Create
// myData := cli.Cache.Read("data")                   // Read
// cli.Cache.Write("data", []byte("cool update"))     // Update
// cli.Cache.Erase("data")                            // Delete
// ```
func (c *cliState) InitCache(d ...string) {
	if len(d) == 0 {
		dir, err := cache.CacheDir()
		if err == nil {
			d = []string{dir}
		}
	}

	cache := strings.Join(d, "/")
	c.Cache = diskv.New(diskv.Options{
		BasePath:          path.Join(cache, "cache"),
		AdvancedTransform: CacheTransform,
		InverseTransform:  InverseCacheTransform,
		CacheSizeMax:      MaxCacheSize,
	})
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
		pathToKey := []string{}
		if len(keys) > 2 {
			pathToKey = keys[1 : len(keys)-1]
		}
		return &diskv.PathKey{
			Path:     pathToKey,
			FileName: keys[len(keys)-1],
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
	pathToKey := []string{cli.Account, subaccount, cli.KeyID}

	// if the key contains "/" we need to split the path
	if strings.Contains(key, "/") {
		keys := strings.Split(key, "/")
		pathToKey = append(pathToKey, keys[0:len(keys)-1]...)
		key = keys[len(keys)-1]
	}

	return &diskv.PathKey{
		Path:     pathToKey,
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

	if c.Token == "" || c.tokenCache.ExpiresAt.Before(time.Now().Add(-10*time.Second)) {
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

// cliAsset is a simple private struct that acks as an envelope for storing
// assets that has a expiration time into the local cache. The Data field
// is an interface on purpose to allow developers store any kind of asset,
// from primitives such as strings, ints, bools, to JSON objects
type cliAsset struct {
	Data      interface{} `json:"data"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// WriteAssetToCache stores an asset with an expiration time
//
// Simple Example: Having a struct named vulnReport
//
// ```go
// cli.WriteAssetToCache("my-report", time.Now().Add(time.Hour * 1), vulnReport{Foo: "bar"})
// ```
func (c *cliState) WriteAssetToCache(key string, expiresAt time.Time, data interface{}) {
	if c.noCache {
		return
	}

	if expiresAt.Before(time.Now()) {
		return // avoid writing assets that are already expired
	}

	c.Log.Debugw("saving asset",
		"feature", "cache",
		"path", key,
		"data", data,
		"expires_at", expiresAt,
	)
	err := c.Cache.Write(key, structToString(cliAsset{data, expiresAt}))
	if err != nil {
		c.Log.Warnw("unable to write asset in cache",
			"feature", "cache",
			"error", err.Error(),
		)
	}
}

// ReadCachedAsset tries to reads an asset with an expiration time, if the
// asset has expired, it returns "true", otherwise it returns "false"
//
// Simple Example: Having a struct named vulnReport
//
// ```go
// var report vulnReport
//
//	if expired := cli.ReadCachedAsset("my-report", &report); !expired {
//	    fmt.Printf("My report: %v\n", report)
//	}
//
// ```
func (c *cliState) ReadCachedAsset(key string, data interface{}) bool {
	if c.noCache {
		return true // if the cache is disabled, all assets are treated like expired
	}

	if dataJSON, err := c.Cache.Read(key); err == nil {
		var asset cliAsset
		if err := json.Unmarshal(dataJSON, &asset); err == nil {
			c.Log.Debugw("asset loaded from cache",
				"feature", "cache",
				"path", key,
				"expires_at", asset.ExpiresAt,
			)

			// check if the cache expired
			if time.Now().After(asset.ExpiresAt) {
				c.Log.Debugw("asset expired, removing from cache",
					"feature", "cache",
					"path", key,
					"time_now", time.Now(),
					"expires_at", asset.ExpiresAt,
				)
				if err := c.Cache.Erase(key); err != nil {
					c.Log.Warnw("unable to erase asset from cache",
						"feature", "cache", "path", key,
					)
				}
				return true
			}

			if err := json.Unmarshal(structToString(asset.Data), &data); err == nil {
				// we successfully retrieved the asset, which has not expired,
				// and we cast it to the proper type
				return false
			}
			c.Log.Warnw("unable to cast asset data",
				"feature", "cache",
				"error", err.Error(),
			)
		}
	}

	return true
}
