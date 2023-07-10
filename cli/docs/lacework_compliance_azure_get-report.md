---
title: "lacework compliance azure get-report"
slug: lacework_compliance_azure_get-report
hide_title: true
---

## lacework compliance azure get-report

Get the latest Azure compliance report

### Synopsis

Get the latest Azure compliance assessment report, these reports run on a regular schedule,
typically once a day. The available report formats are human-readable (default), json and pdf.

To list all Azure tenants and subscriptions configured in your account:

    lacework compliance azure list

To show recommendation details and affected resources for a recommendation id:

    lacework compliance azure get-report <tenant_id> <subscriptions_id> [recommendation_id]

To retrieve a specific report by its report name:

    lacework compliance azure get-report <tenant_id> <subscriptions_id> --report_name 'Azure CIS 1.3.1 Report'


```
lacework compliance azure get-report <tenant_id> <subscriptions_id> [flags]
```

### Options

```
      --category strings     filter report details by category (networking, storage, ...)
      --csv                  output report in CSV format
      --details              increase details about the compliance report
  -h, --help                 help for get-report
      --pdf                  download report in PDF format
      --report_name string   report name to display, run 'lacework report-definitions list' for more information. (default "CIS Microsoft Azure Foundations Benchmark v1.5.0")
      --service strings      filter report details by service (azure:ms:storage, azure:ms:sql, azure:ms:network, ...)
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

* [lacework compliance azure](lacework_compliance_azure.md)	 - Compliance for Azure Cloud

