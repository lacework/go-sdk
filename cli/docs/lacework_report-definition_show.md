---
title: "lacework report-definition show"
slug: lacework_report-definition_show
hide_title: true
---

## lacework report-definition show

Show a report definition by ID

### Synopsis

Show a single report definition by it's ID.
To show specific report definition version:

    lacework report-definition show <report_definition_id> --version <version>

To show all versions of a report definition:

    lacework report-definition show <report_definition_id> --version all



```
lacework report-definition show <report_definition_id> [flags]
```

### Options

```
  -h, --help             help for show
      --version string   show a version of a report definition
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

* [lacework report-definition](lacework_report-definition.md)	 - Manage report definitions

