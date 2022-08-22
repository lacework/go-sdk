---
title: "lacework policy delete"
slug: lacework_policy_delete
hide_title: true
---

## lacework policy delete

Delete a policy

### Synopsis

Delete a policy by providing the policy ID.

Use the command 'lacework policy list' to list the registered policies in
your Lacework account.

```
lacework policy delete <policy_id> [flags]
```

### Options

```
      --cascade   delete policy and its associated query
  -h, --help      help for delete
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

* [lacework policy](lacework_policy.md)	 - Manage policies

