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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// queryCreateCmd represents the lql create command
	queryCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create a query",
		Long: `Create a query.

There are multiple ways you can create a query:

  * Typing the query into your default editor (via $EDITOR)
  * From a local file on disk using the flag '--file'
  * From a URL using the flag '--url'

There are also multiple formats you can use to define a query:

  * Javascript Object Notation (JSON)
  * YAML Ain't Markup Language (YAML)

To launch your default editor and create a new query.

    lacework lql create

The following example comes from Lacework's implementation of a policy query:

---
evaluatorId: Cloudtrail
queryId: LW_Global_AWS_CTA_AccessKeyDeleted
queryText: |-
  LW_Global_AWS_CTA_AccessKeyDeleted {
      source {
          CloudTrailRawEvents
      }
      filter {
          EVENT_SOURCE = 'iam.amazonaws.com'
          and EVENT_NAME = 'DeleteAccessKey'
          and ERROR_CODE is null
      }
      return distinct {
          INSERT_ID,
          INSERT_TIME,
          EVENT_TIME,
          EVENT
      }
  }

Identifier of the query that executes while running the policy


This query specifies an identifier named 'LW_Global_AWS_CTA_AccessKeyDeleted'.
Policy evaluation uses this dataset (along with the filters) to identify AWS
CloudTrail events that signify that an IAM access key was deleted. The query
is delimited by '{ }' and contains three sections:

  * Source data is specified in the 'source' clause. The source of data is the
  'CloudTrailRawEvents' dataset. LQL queries generally refer to other datasets,
  and customizable policies always target a suitable dataset.

  * Records of interest are specified by the 'filter' clause. In the example, the
  records available in 'CloudTrailRawEvents' are filtered for those whose source
  is 'iam.amazonaws.com', whose event name is 'DeleteAccessKey', and that do not
  have any error code. The syntax for this filtering expression strongly resembles SQL.

  * The fields this query exposes are listed in the 'return' clause. Because there
  may be unwanted duplicates among result records when Lacework composes them from
  just these four columns, the distinct modifier is added. This behaves like a SQL
  'SELECT DISTINCT'. Each returned column in this case is just a field that is present
  in 'CloudTrailRawEvents', but we can compose results by manipulating strings, dates,
  JSON and numbers as well.

The resulting dataset is shaped like a table. The table's columns are named with the
names of the columns selected. If desired, you could alias them to other names as well.

For more information about LQL, visit:

    https://support.lacework.com/hc/en-us/articles/4402301824403-LQL-Overview
`,
		Args: cobra.NoArgs,
		RunE: createQuery,
	}
)

func init() {
	// add sub-commands to the lql command
	queryCmd.AddCommand(queryCreateCmd)

	setQuerySourceFlags(queryCreateCmd)
}

func createQuery(cmd *cobra.Command, args []string) error {
	msg := "unable to create query"

	// input query
	queryString, err := inputQuery(cmd)
	if err != nil {
		return errors.Wrap(err, msg)
	}
	// parse query
	newQuery, err := parseQuery(queryString)
	if err != nil {
		return errors.Wrap(err, msg)
	}

	cli.Log.Debugw("creating query", "query", queryString)
	cli.StartProgress(" Creating query...")
	create, err := cli.LwApi.V2.Query.Create(newQuery)
	cli.StopProgress()
	if err != nil {
		return errors.Wrap(err, msg)
	}

	if cli.JSONOutput() {
		return cli.OutputJSON(create.Data)
	}

	cli.OutputHuman("The query %s was created.\n", create.Data.QueryID)
	return nil
}
