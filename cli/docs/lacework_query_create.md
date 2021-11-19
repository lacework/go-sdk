---
title: "lacework query create"
slug: lacework_query_create
hide_title: true
---

## lacework query create

Create a query

### Synopsis


There are multiple ways you can create a query:

  * Typing the query into your default editor (via $EDITOR)
  * Piping a query to the Lacework CLI command (via $STDIN)
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

  https://docs.lacework.com/lql-overview


```
lacework query create [flags]
```

### Options

```
  -f, --file string   path to a query to create
  -h, --help          help for create
  -u, --url string    url to a query to create
```

### Options inherited from parent commands

```
  -a, --account string      account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
  -k, --api_key string      access key id
  -s, --api_secret string   secret access key
      --api_token string    access token (replaces the use of api_key and api_secret)
      --debug               turn on debug logging
      --json                switch commands output from human-readable to json format
      --nocache             turn off caching
      --nocolor             turn off colors
      --noninteractive      turn off interactive mode (disable spinners, prompts, etc.)
      --organization        access organization level data sets (org admins only)
  -p, --profile string      switch between profiles configured at ~/.lacework.toml
      --subaccount string   sub-account name inside your organization (org admins only)
```

### SEE ALSO

* [lacework query](lacework_query.md)	 - Run and manage queries

