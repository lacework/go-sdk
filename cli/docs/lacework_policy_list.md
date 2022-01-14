---
title: "lacework policy list"
slug: lacework_policy_list
hide_title: true
---

## lacework policy list

List policies

### Synopsis

List all the registered policies in your Lacework account.

```
lacework policy list [flags]
```

### Options

```
      --alert_enabled     only show alert_enabled policies
      --enabled           only show enabled policies
  -h, --help              help for list
      --severity string   filter policies by severity threshold (critical, high, medium, low, info)
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

