---
title: "lacework alert list"
slug: lacework_alert_list
hide_title: true
---

## lacework alert list

List all alerts

### Synopsis

List all alerts.

By default, alerts are shown for the last 24 hours.
Use a custom time range by suppling a range flag...

    lacework alert ls --range "last 7 days"

Or by specifying start and end flags.

    lacework alert ls --start "-7d@d" --end "now"

Start and end times may be specified in one of the following formats:
    A. A relative time specifier
    B. RFC3339 date and time
    C. Epoch time in milliseconds

To list open alerts of type "NewViolations" with high or critical severity.

    lacework alert ls --status Open --severity high --type NewViolations


```
lacework alert list [flags]
```

### Options

```
      --end string        end time for alerts (default "now")
  -h, --help              help for list
      --range string      natural time range for alerts
      --severity string   filter alerts by severity threshold (critical, high, medium, low, info)
      --start string      start time for alerts (default "-24h")
      --status string     filter alerts by status (Open, Closed)
      --type string       filter alerts by type
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

