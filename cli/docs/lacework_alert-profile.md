---
title: "lacework alert-profile"
slug: lacework_alert-profile
hide_title: true
---

## lacework alert-profile

Manage alert profiles

### Synopsis

Manage alert profiles to define how your LQL queries get consumed into alerts.

An alert profile consists of the ID of the new profile, the ID of an existing profile that
the new profile extends, and a list of alert templates.

### Options

```
  -h, --help   help for alert-profile
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
* [lacework alert-profile create](lacework_alert-profile_create.md)	 - Create a new alert profile
* [lacework alert-profile delete](lacework_alert-profile_delete.md)	 - Delete an alert profile
* [lacework alert-profile list](lacework_alert-profile_list.md)	 - List all alert profiles
* [lacework alert-profile show](lacework_alert-profile_show.md)	 - Show an alert profile by ID
* [lacework alert-profile update](lacework_alert-profile_update.md)	 - Update alert templates from an existing alert profile

