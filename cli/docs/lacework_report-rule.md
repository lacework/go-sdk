---
title: "lacework report-rule"
slug: lacework_report-rule
hide_title: true
---

## lacework report-rule

Manage report rules

### Synopsis

Manage report rules to route reports to one or more email alert channels.		

A report rule has four parts:

  1. Email alert channel(s) that should receive the report
  2. One or more severities to include
  3. Resource group(s) containing the subset of your environment to consider
  4. Notification types containing which report information to send


### Options

```
  -h, --help   help for report-rule
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
* [lacework report-rule create](lacework_report-rule_create.md)	 - Create a new report rule
* [lacework report-rule delete](lacework_report-rule_delete.md)	 - Delete a report rule
* [lacework report-rule list](lacework_report-rule_list.md)	 - List all report rules
* [lacework report-rule show](lacework_report-rule_show.md)	 - Show a report rule by ID

