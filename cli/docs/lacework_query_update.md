---
title: "lacework query update"
slug: lacework_query_update
---

## lacework query update

update a query

### Synopsis

Update a single LQL query.

There are multiple ways you can update a query:

  * Typing the query into your default editor (via $EDITOR)
  * From a local file on disk using the flag '--file'
  * From a URL using the flag '--url'

There are also multiple formats you can use to define a query:

  * Javascript Object Notation (JSON)
  * YAML Ain't Markup Language (YAML)

To launch your default editor and update a query.

    lacework query update


```
lacework query update [flags]
```

### Options

```
  -f, --file string   path to a query to update
  -h, --help          help for update
  -u, --url string    url to a query to update
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

* [lacework query](lacework_query.md)	 - run and manage queries

