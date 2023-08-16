---
title: "lacework compliance aws get-report"
slug: lacework_compliance_aws_get-report
hide_title: true
---

## lacework compliance aws get-report

Get the latest AWS compliance report

### Synopsis

Get the latest compliance assessment report from the provided AWS account, these
reports run on a regular schedule, typically once a day. The available report formats
are human-readable (default), json and pdf.

To list all AWS accounts configured in your account:

    lacework compliance aws list-accounts

To show recommendation details and affected resources for a recommendation id:

    lacework compliance aws get-report <account_id> [recommendation_id]

To retrieve a specific report by its report name:

    lacework compliance aws get-report <account_id> --report_name 'AWS CSA CCM 4.0.5'


```
lacework compliance aws get-report <account_id> [recommendation_id] [flags]
```

### Options

```
      --category strings     filter report details by category (identity-and-access-management, s3, logging...)
      --csv                  output report in CSV format
      --details              increase details about the compliance report
  -h, --help                 help for get-report
      --pdf                  download report in PDF format
      --report_name string   report name to display, run 'lacework report-definitions list' for more information. (default "CIS Amazon Web Services Foundations Benchmark v1.4.0")
      --service strings      filter report details by service (aws:s3, aws:iam, aws:cloudtrail, ...)
      --severity string      filter report details by severity threshold (critical, high, medium, low, info)
      --status string        filter report details by status (non-compliant, requires-manual-assessment, suppressed, compliant, could-not-assess)
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

* [lacework compliance aws](lacework_compliance_aws.md)	 - Compliance for AWS

