---
title: "lacework alert show"
slug: lacework_alert_show
hide_title: true
---

## lacework alert show

Show details about a specific alert

### Synopsis

Show details about a specific alert.

There are different types of alert details that can be shown to assist
with alert investigation. These types are referred to as alert detail scopes.

The following alert detail scopes are available:

  * Details (default)
  * Investigation
  * Events
  * RelatedAlerts
  * Integrations
  * Timeline

View an alert's timeline details:

  lacework alert show <alert_id> --scope Timeline


```
lacework alert show <alert_id> [flags]
```

### Options

```
  -h, --help           help for show
      --scope string   type of alert details to show (default "Details")
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

