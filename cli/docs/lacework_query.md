---
title: "lacework query"
slug: lacework_query
---

## lacework query

run and manage queries

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

   https://support.lacework.com/hc/en-us/articles/4402301824403-LQL-Overview

To view all LQL queries in your Lacework account.

    lacework query ls

To show a query.

    lacework query show <query_id>

To execute a query.

    lacework query run <query_id>

** NOTE: LQL syntax may change. **


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

* [lacework](/cli/commands/lacework/)	 - A tool to manage the Lacework cloud security platform.
* [lacework query create](/cli/commands/lacework_query_create/)	 - create a query
* [lacework query delete](/cli/commands/lacework_query_delete/)	 - delete a query
* [lacework query list](/cli/commands/lacework_query_list/)	 - list queries
* [lacework query list-sources](/cli/commands/lacework_query_list-sources/)	 - list Lacework query data sources
* [lacework query run](/cli/commands/lacework_query_run/)	 - run a query
* [lacework query show](/cli/commands/lacework_query_show/)	 - show a query
* [lacework query show-source](/cli/commands/lacework_query_show-source/)	 - show Lacework query data source
* [lacework query update](/cli/commands/lacework_query_update/)	 - update a query
* [lacework query validate](/cli/commands/lacework_query_validate/)	 - validate a query

