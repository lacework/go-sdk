---
title: "lacework cloud-account delete"
slug: lacework_cloud-account_delete
hide_title: true
---

## lacework cloud-account delete

Delete a cloud account integration

### Synopsis

Delete a single cloud account integration by its integration GUID:

    lacework cloud-account delete <intg_guid>

Or delete many at once from a backup file produced by 'cloud-account backup':

    lacework cloud-account delete --bulk --file backup.json

```
lacework cloud-account delete [flags]
```

### Options

```
      --bulk          delete every integration listed in --file
      --dry-run       show what would be deleted without deleting
  -f, --file string   backup file listing the integrations to delete (with --bulk)
  -h, --help          help for delete
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

* [lacework cloud-account](lacework_cloud-account.md)	 - Manage cloud accounts

