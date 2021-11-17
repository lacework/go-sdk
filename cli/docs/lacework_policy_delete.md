---
title: "lacework policy delete"
slug: lacework_policy_delete
---

## lacework policy delete

delete a policy

### Synopsis

Delete a policy by providing the policy id.

Use the command 'lacework policy list' to list the registered policies in
your Lacework account.

```
lacework policy delete <policy_id> [flags]
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

* [lacework policy](/cli/commands/lacework_policy/)	 - manage policies

