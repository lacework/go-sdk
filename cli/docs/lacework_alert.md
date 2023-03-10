---
title: "lacework alert"
slug: lacework_alert
hide_title: true
---

## lacework alert

Inspect and manage alerts

### Synopsis

Inspect and manage alerts.

Lacework provides real-time alerts that are interactive and manageable.
Each alert contains various metadata information, such as severity level, type, status, alert category, and associated tags.

You can also post a comment to an alert's timeline; or change an alert status from Open to Closed.

For more information about alerts, visit:

https://docs.lacework.com/console/alerts-overview

To view all alerts in your Lacework account.

    lacework alert ls

To show an alert.

    lacework alert show <alert_id>

To close an alert.

    lacework alert close <alert_id>


### Options

```
  -h, --help   help for alert
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

* [lacework](lacework.md)	 - A tool to manage the Lacework cloud security platform.
* [lacework alert close](lacework_alert_close.md)	 - Close an alert
* [lacework alert comment](lacework_alert_comment.md)	 - Add a comment
* [lacework alert list](lacework_alert_list.md)	 - List all alerts
* [lacework alert open](lacework_alert_open.md)	 - Open a specified alert in a web browser
* [lacework alert show](lacework_alert_show.md)	 - Show details about a specific alert

