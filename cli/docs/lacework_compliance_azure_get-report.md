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

To run an ad-hoc compliance assessment use the command:

    lacework compliance azure run-assessment <tenant_id>

To show recommendation details and affected resources for a recommendation id:

    lacework compliance azure get-report <tenant_id> <subscriptions_id> [recommendation_id]


```
lacework compliance azure get-report <tenant_id> <subscriptions_id> [flags]
```

### Options

```
      --category strings   filter report details by category (networking, storage, ...)
      --csv                output report in CSV format
      --details            increase details about the compliance report
  -h, --help               help for get-report
      --pdf                download report in PDF format
      --service strings    filter report details by service (azure:ms:storage, azure:ms:sql, azure:ms:network, ...)
      --severity string    filter report details by severity threshold (critical, high, medium, low, info)
      --status string      filter report details by status (non-compliant, requires-manual-assessment, suppressed, compliant, could-not-assess)
      --type string        report type to display, supported types: CIS, SOC, or PCI (default "CIS")
```

### Options inherited from parent commands

```
  -a, --account string      account URL (i.e. <ACCOUNT>[.CUSTER][.corp].lacework.net)
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

