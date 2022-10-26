---
title: "lacework cloud-account"
slug: lacework_cloud-account
hide_title: true
---

## lacework cloud-account

Manage cloud accounts

### Synopsis

Manage cloud account integrations with Lacework

### Options

```
  -h, --help   help for cloud-account
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
* [lacework cloud-account create](lacework_cloud-account_create.md)	 - Create a new cloud account integration
* [lacework cloud-account delete](lacework_cloud-account_delete.md)	 - Delete a cloud account integration
* [lacework cloud-account list](lacework_cloud-account_list.md)	 - List all available cloud account integrations
* [lacework cloud-account show](lacework_cloud-account_show.md)	 - Show a single cloud account integration

