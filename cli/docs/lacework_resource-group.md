---
title: "lacework resource-group"
slug: lacework_resource-group
hide_title: true
---

## lacework resource-group

Manage resource groups

### Synopsis

Manage Lacework-identifiable assets via the use of resource groups.

### Options

```
  -h, --help   help for resource-group
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

* [lacework](lacework.md)	 - A tool to manage the Lacework cloud security platform.
* [lacework resource-group create](lacework_resource-group_create.md)	 - Create a new resource group
* [lacework resource-group delete](lacework_resource-group_delete.md)	 - Delete a resource group
* [lacework resource-group list](lacework_resource-group_list.md)	 - List all resource groups
* [lacework resource-group show](lacework_resource-group_show.md)	 - Get resource group by ID

