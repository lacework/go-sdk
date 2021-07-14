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
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lacework/go-sdk/api"
)

var (
	// duration of the access token in seconds
	durationSeconds int

	// accessTokenCmd represents the access-token command
	accessTokenCmd = &cobra.Command{
		Use:   "access-token",
		Short: "generate temporary API access tokens",
		Long: `Generates a temporary API access token that can be used to access the
Lacework API. The token will be valid for the duration that you specify.`,
		Args: cobra.NoArgs,
		RunE: generateAccessToken,
	}
)

func init() {
	// add the access-token command
	rootCmd.AddCommand(accessTokenCmd)

	accessTokenCmd.Flags().IntVarP(&durationSeconds,
		"duration_seconds", "d", api.DefaultTokenExpiryTime,
		"duration in seconds that the access token should remain valid",
	)
}

func generateAccessToken(_ *cobra.Command, args []string) error {
	var (
		response *api.TokenData
		err      error
	)

	if durationSeconds == api.DefaultTokenExpiryTime {
		response, err = cli.LwApi.GenerateToken()
		if err != nil {
			return errors.Wrap(err, "unable to generate access token")
		}
	} else {
		// if the duration is different from the default,
		// regenerate the lacework api client
		client, err := api.NewClient(cli.Account,
			api.WithApiV2(),
			api.WithLogLevel(cli.LogLevel),
			api.WithExpirationTime(durationSeconds),
			api.WithHeader("User-Agent", fmt.Sprintf("Command-Line/%s", Version)),
		)
		if err != nil {
			return errors.Wrap(err, "unable to generate api client")
		}

		response, err = client.GenerateTokenWithKeys(cli.KeyID, cli.Secret)
		if err != nil {
			return errors.Wrap(err, "unable to generate access token")
		}
	}

	// cache new token
	err = cli.Cache.Write("token", structToString(response))
	if err != nil {
		cli.Log.Warnw("unable to write token in cache",
			"feature", "cache",
			"error", err.Error(),
		)
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(response)
	}

	cli.OutputHuman(response.Token)
	cli.OutputHuman("\n")
	return nil
}
