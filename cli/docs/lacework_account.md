---
title: "lacework account"
slug: lacework_account
---

## lacework account

manage accounts in an organization (org admins only)

### Synopsis

Manage accounts inside your Lacework organization.

An organization can contain multiple accounts so you can also manage components
such as alerts, resource groups, team members, and audit logs at a more granular
level inside an organization. A team member may have access to multiple accounts
and can easily switch between them.

To enroll your Lacework account in an organization follow the documentation:

  https://support.lacework.com/hc/en-us/articles/360041727394-Organization-Overview
    

### Options

```
  -h, --help   help for account
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
* [lacework account list](lacework_account_list.md)	 - list all accounts

