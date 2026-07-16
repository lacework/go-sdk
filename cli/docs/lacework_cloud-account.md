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
* [lacework cloud-account backup](lacework_cloud-account_backup.md)	 - Back up cloud account integrations of a given type to a file
* [lacework cloud-account cleanup](lacework_cloud-account_cleanup.md)	 - Delete leftover cloud resources for the integrations in a backup file
* [lacework cloud-account create](lacework_cloud-account_create.md)	 - Create a new cloud account integration
* [lacework cloud-account delete](lacework_cloud-account_delete.md)	 - Delete a cloud account integration
* [lacework cloud-account list](lacework_cloud-account_list.md)	 - List all available cloud account integrations
* [lacework cloud-account migrate](lacework_cloud-account_migrate.md)	 - Mark a GCPv1 (storage-based) cloud account integration for migration
* [lacework cloud-account repair](lacework_cloud-account_repair.md)	 - Re-register a missing cloud-account integration from its onboarding template
* [lacework cloud-account restore](lacework_cloud-account_restore.md)	 - Re-create cloud account integrations from a backup file
* [lacework cloud-account show](lacework_cloud-account_show.md)	 - Show a single cloud account integration

