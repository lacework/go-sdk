---
title: "lacework compliance azure list-subscriptions"
slug: lacework_compliance_azure_list-subscriptions
hide_title: true
---

## lacework compliance azure list-subscriptions

List subscriptions from tenant

### Synopsis

List all Azure subscriptions from the provided Tenant ID.

Use the following command to list all Azure Tenants configured in your account:

    lacework compliance az list

```
lacework compliance azure list-subscriptions <tenant_id> [flags]
```

### Options

```
  -h, --help   help for list-subscriptions
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

