---
title: "lacework compliance azure"
slug: lacework_compliance_azure
hide_title: true
---

## lacework compliance azure

Compliance for Azure Cloud

### Synopsis

Manage compliance reports for Azure Cloud.

To list all Azure tenants configured in your account:

    lacework compliance azure list-tenants

To list all Azure subscriptions from a tenant, use the command:

    lacework compliance azure list-subscriptions <tenant_id>

To get the latest Azure compliance assessment report, use the command:

    lacework compliance azure get-report <tenant_id> <subscription_id>

These reports run on a regular schedule, typically once a day.

To run an ad-hoc compliance assessment use the command:

    lacework compliance azure run-assessment <tenant_id>


### Options

```
  -h, --help   help for azure
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
* [lacework compliance azure get-report](lacework_compliance_azure_get-report.md)	 - Get the latest Azure compliance report
* [lacework compliance azure list](lacework_compliance_azure_list.md)	 - List Azure tenants and subscriptions
* [lacework compliance azure list-subscriptions](lacework_compliance_azure_list-subscriptions.md)	 - List subscriptions from tenant
* [lacework compliance azure run-assessment](lacework_compliance_azure_run-assessment.md)	 - Run a new Azure compliance assessment

