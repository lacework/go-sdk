---
title: "lacework policy create"
slug: lacework_policy_create
hide_title: true
---

## lacework policy create

Create a policy

### Synopsis

Create a policy.

A policy is represented in either JSON or YAML format.

The following attributes are minimally required:

    ---
    policyType: Violation
    queryId: MyQuery
    title: My Policy
    enabled: false
    description: My Policy Description
    remediation: My Policy Remediation
    severity: high
    evalFrequency: Daily
    alertEnabled: false
    alertProfile: LW_CloudTrail_Alerts


```
lacework policy create [flags]
```

### Options

```
  -f, --file string   path to a policy to create
  -h, --help          help for create
  -u, --url string    url to a policy to create
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

* [lacework policy](lacework_policy.md)	 - Manage policies

