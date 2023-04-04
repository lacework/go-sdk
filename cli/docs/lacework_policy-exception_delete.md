---
title: "lacework policy-exception delete"
slug: lacework_policy-exception_delete
hide_title: true
---

## lacework policy-exception delete

Delete a policy exception

### Synopsis

Delete a policy exception. 

To remove a policy exception, run the delete command with policy ID and exception ID arguments:

    lacework policy-exception delete <policy_id> <exception_id>

```
lacework policy-exception delete <policy_id> <exception_id> [flags]
```

### Options

```
  -h, --help   help for delete
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

* [lacework policy-exception](lacework_policy-exception.md)	 - Manage policy exceptions

