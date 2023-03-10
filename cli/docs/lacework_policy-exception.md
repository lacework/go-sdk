---
title: "lacework policy-exception"
slug: lacework_policy-exception
hide_title: true
---

## lacework policy-exception

Manage policy exceptions

### Synopsis

Manage policy exceptions in your Lacework account.

To view all the policies in your Lacework account.

    lacework policy list


### Options

```
  -h, --help   help for policy-exception
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
* [lacework policy-exception create](lacework_policy-exception_create.md)	 - Create a policy exception
* [lacework policy-exception delete](lacework_policy-exception_delete.md)	 - Delete a policy exception
* [lacework policy-exception list](lacework_policy-exception_list.md)	 - List all exceptions from a single policy
* [lacework policy-exception show](lacework_policy-exception_show.md)	 - Show details about a policy exception

