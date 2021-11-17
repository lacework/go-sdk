---
title: "lacework resource-group"
slug: lacework_resource-group
---

## lacework resource-group

manage resource groups

### Synopsis

Manage Lacework-identifiable assets via the use of resource groups.

### Options

```
  -h, --help   help for resource-group
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

* [lacework](/cli/commands/lacework/)	 - A tool to manage the Lacework cloud security platform.
* [lacework resource-group create](/cli/commands/lacework_resource-group_create/)	 - create a new resource group
* [lacework resource-group delete](/cli/commands/lacework_resource-group_delete/)	 - delete a resource group
* [lacework resource-group list](/cli/commands/lacework_resource-group_list/)	 - list all resource groups
* [lacework resource-group show](/cli/commands/lacework_resource-group_show/)	 - get resource group by id

