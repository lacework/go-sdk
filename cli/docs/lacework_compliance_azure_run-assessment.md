---
title: "lacework compliance azure run-assessment"
slug: lacework_compliance_azure_run-assessment
hide_title: true
---

## lacework compliance azure run-assessment

Run a new Azure compliance assessment

### Synopsis

Run a compliance assessment of the provided Azure tenant.

To list all Azure tenants and subscriptions configured in your account:

    lacework compliance azure list

```
lacework compliance azure run-assessment <tenant_id> [flags]
```

### Options

```
  -h, --help   help for run-assessment
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

