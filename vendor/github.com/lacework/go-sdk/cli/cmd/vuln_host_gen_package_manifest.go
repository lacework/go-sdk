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

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// vulHostGenPkgManifestCmd represents the 'lacework vuln host generate-pkg-manifest' command
	vulHostGenPkgManifestCmd = &cobra.Command{
		Use:   "generate-pkg-manifest",
		Args:  cobra.NoArgs,
		Short: "Generates a package-manifest from the local host",
		Long: `Generates a package-manifest formatted for usage with the Lacework
scan package-manifest API.

Additionally, you can automatically generate a package-manifest from
the local host and send it directly to the Lacework API with the command:

    lacework vulnerability host scan-pkg-manifest --local`,
		RunE: func(_ *cobra.Command, _ []string) error {
			manifest, err := cli.GeneratePackageManifest()
			if err != nil {
				return errors.Wrap(err, "unable to generate package manifest")
			}

			return cli.OutputJSON(manifest)
		},
	}
)
