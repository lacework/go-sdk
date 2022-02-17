---
title: "lacework query"
slug: lacework_query
hide_title: true
---

## lacework query

Run and manage queries

### Synopsis

Run and manage Lacework Query Language (LQL) queries.

To provide customizable specification of datasets, Lacework provides the Lacework
Query Language (LQL). LQL is a human-readable text syntax for specifying selection,
filtering, and manipulation of data.

Currently, Lacework has introduced LQL for configuration of AWS CloudTrail policies
and queries. This means you can use LQL to customize AWS CloudTrail policies only.
For all other policies, use the previous existing methods.

Lacework ships a set of default LQL queries that are available in your account.

For more information about LQL, visit:

  https://docs.lacework.com/lql-overview

To view all LQL queries in your Lacework account.

    lacework query ls

To show a query.

    lacework query show <query_id>

To execute a query.

    lacework query run <query_id>

**NOTE: LQL syntax may change.**


### Options

```
  -h, --help   help for query
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

* [lacework](lacework.md)	 - A tool to manage the Lacework cloud security platform.
* [lacework query create](lacework_query_create.md)	 - Create a query
* [lacework query delete](lacework_query_delete.md)	 - Delete a query
* [lacework query list](lacework_query_list.md)	 - List queries
* [lacework query list-sources](lacework_query_list-sources.md)	 - List Lacework query data sources
* [lacework query preview-source](lacework_query_preview-source.md)	 - Preview Lacework query data source
* [lacework query run](lacework_query_run.md)	 - Run a query
* [lacework query show](lacework_query_show.md)	 - Show a query
* [lacework query show-source](lacework_query_show-source.md)	 - Show Lacework query data source
* [lacework query update](lacework_query_update.md)	 - Update a query
* [lacework query validate](lacework_query_validate.md)	 - Validate a query

