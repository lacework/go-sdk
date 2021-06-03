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
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/lacework/go-sdk/lwconfig"
)

// The name of the directory we will store backups of configuration files before migrating them
const ConfigBackupDir = "cfg_backups"

// Migrations executes automatic configuration migrations,
// if a configuration file does not exist, it will only update
// the CLI state to the appropriate parameters
func (c *cliState) Migrations() (err error) {
	if !needMigration() {
		return nil
	}

	c.Log.Debugw("executing v2 migration")
	c.Event.Feature = featMigrateConfigV2
	defer func() {
		if err == nil {
			c.SendHoneyvent()
		} else {
			err = errors.Wrap(err, "during v2 migration")
		}

		// update global honeyvent with updated state
		c.Event.Account = c.Account
		c.Event.Subaccount = c.Subaccount
		c.Event.CfgVersion = c.CfgVersion
	}()

	err = c.VerifySettings()
	if err != nil {
		return err
	}

	orgInfo, err := c.LwApi.Account.GetOrganizationInfo()
	if err != nil {
		return err
	}

	// set new v2 config version and notify our feature event
	c.CfgVersion = 2
	c.Event.AddFeatureField("config_version", c.CfgVersion)
	c.Event.AddFeatureField("org_account", orgInfo.OrgAccount)
	// NOTE: @afiune this will be a constant pattern below where
	// we will update settings and notify the feature event

	if orgInfo.OrgAccount {
		// we only need to update the account/sub-account
		// if the user has an organizational account
		c.Log.Debugw("organizational account detected")
		c.Event.AddFeatureField("org_account_url", orgInfo.OrgAccountURL)

		primaryAccount := strings.ToLower(orgInfo.AccountName())

		// if the user is accessing a sub-account, that is, if the current
		// account is different from the primary account name, set it as
		// a what it is, the sub-account
		if primaryAccount != c.Account {
			c.Log.Debugw("updating account settings for APIv2",
				"old_account", c.Account,
				"new_account", primaryAccount,
			)
			c.Subaccount = c.Account
			c.Account = primaryAccount

			c.Event.AddFeatureField("account", c.Account)
			c.Event.AddFeatureField("subaccount", c.Subaccount)

			c.Log.Debugw("generating new API client")
			err = c.NewClient()
			if err != nil {
				return err
			}
		}
	}

	// if the configuration file does not exist, most likely the user
	// is executing the CLI via env variables or flags, update feature
	// field and exit migration
	if !fileExists(viper.ConfigFileUsed()) {
		c.Log.Debugw("config file not found, skipping profile migration")
		c.Event.AddFeatureField("config_file", "not_found")
		return nil
	}

	c.Log.Debugw("config found, migrating profile", "profile", c.Profile)
	migratedProfile := lwconfig.Profile{
		Account:    c.Account,
		Subaccount: c.Subaccount,
		ApiKey:     c.KeyID,
		ApiSecret:  c.Secret,
		Version:    c.CfgVersion,
	}

	// create a backup before modifying the user's configuration
	bkpPath, err := createConfigurationBackup()
	if err != nil {
		return err
	}
	c.Log.Debugw("configuration backup", "path", bkpPath)
	c.Event.AddFeatureField("backup_file", path.Base(bkpPath))

	// store the migrated profile
	err = lwconfig.StoreProfileAt(viper.ConfigFileUsed(), c.Profile, migratedProfile)
	if err != nil {
		return errors.Wrap(err, "unable to store migrated profile")
	}

	c.Log.Debugw("configuration migrated successfully")
	return nil
}

func createConfigurationBackup() (string, error) {
	profiles, err := lwconfig.LoadProfilesFrom(viper.ConfigFileUsed())
	if err != nil {
		return "", err
	}

	cacheDir, err := versionCacheDir()
	if err != nil {
		return "", err
	}

	backupDir := path.Join(cacheDir, ConfigBackupDir)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	backupCfgPath := path.Join(backupDir,
		fmt.Sprintf(".lacework.toml.%s.%s.bkp",
			time.Now().Format("20060102150405"), newID()),
	)
	return backupCfgPath, lwconfig.StoreAt(backupCfgPath, profiles)
}

func needMigration() bool {
	return cli.CfgVersion != 2 // &&
}
