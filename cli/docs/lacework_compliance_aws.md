---
title: "lacework compliance aws"
slug: lacework_compliance_aws
hide_title: true
---

## lacework compliance aws

Compliance for AWS

### Synopsis

Manage compliance reports for Amazon Web Services (AWS).

To list all AWS accounts configured in your account:

    lacework compliance aws list-accounts

To get the latest AWS compliance assessment report:

    lacework compliance aws get-report <account_id>

These reports run on a regular schedule, typically once a day.


### Options

```
  -h, --help   help for aws
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

* [lacework compliance](lacework_compliance.md)	 - Manage compliance reports
* [lacework compliance aws get-report](lacework_compliance_aws_get-report.md)	 - Get the latest AWS compliance report
* [lacework compliance aws list-accounts](lacework_compliance_aws_list-accounts.md)	 - List all AWS accounts configured
* [lacework compliance aws search](lacework_compliance_aws_search.md)	 - Search for all known violations of a given resource arn

