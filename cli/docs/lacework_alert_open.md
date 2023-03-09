---
title: "lacework alert open"
slug: lacework_alert_open
hide_title: true
---

## lacework alert open

Open a specified alert in a web browser

### Synopsis

Open a specified alert in a web browser.

```
lacework alert open <alert_id> [flags]
```

### Options

```
  -h, --help   help for open
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

* [lacework alert](lacework_alert.md)	 - Inspect and manage alerts

