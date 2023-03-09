---
title: "lacework alert close"
slug: lacework_alert_close
hide_title: true
---

## lacework alert close

Close an alert

### Synopsis

Use this command to change the status of an alert to closed.

The reason for closing the alert must be provided from the following options:

  * 0 - Other
  * 1 - False positive
  * 2 - Not enough information
  * 3 - Malicious and have resolution in place
  * 4 - Expected because of routine testing.

Reasons may be provided inline or via prompt.

If you choose Other, a comment is required and should contain a brief explanation of why the alert is closed.
Comments may be provided inline or via editor.

**Note: A closed alert cannot be reopened. You will be prompted to confirm closure of the alert.  
This prompt can be bypassed with the --noninteractive flag**


```
lacework alert close <alert_id> [flags]
```

### Options

```
  -c, --comment string   a comment to associate with the alert closure
  -h, --help             help for close
  -r, --reason int       the reason for closing the alert (default -1)
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

