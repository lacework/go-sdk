---
title: "lacework alert-rule delete"
slug: lacework_alert-rule_delete
hide_title: true
---

## lacework alert-rule delete

Delete a alert rule

### Synopsis

Delete a single alert rule by it's ID.

```
lacework alert-rule delete <alert_rule_id> [flags]
```

### Options

```
  -h, --help   help for delete
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

* [lacework alert-rule](lacework_alert-rule.md)	 - Manage alert rules

