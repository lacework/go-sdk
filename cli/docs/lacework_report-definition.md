---
title: "lacework report-definition"
slug: lacework_report-definition
hide_title: true
---

## lacework report-definition

Manage report definitions

### Synopsis

Manage report definitions to configure the data retrieval and layout information for a report.


### Options

```
  -h, --help   help for report-definition
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
* [lacework report-definition create](lacework_report-definition_create.md)	 - Create a report definition
* [lacework report-definition delete](lacework_report-definition_delete.md)	 - Delete a report definition
* [lacework report-definition diff](lacework_report-definition_diff.md)	 - Compare two versions of a report definition
* [lacework report-definition list](lacework_report-definition_list.md)	 - List all report definitions
* [lacework report-definition revert](lacework_report-definition_revert.md)	 - Update a report definition
* [lacework report-definition show](lacework_report-definition_show.md)	 - Show a report definition by ID
* [lacework report-definition update](lacework_report-definition_update.md)	 - Update a report definition

