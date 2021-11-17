---
title: "lacework agent token update"
slug: lacework_agent_token_update
---

## lacework agent token update

update an agent access token

### Synopsis

Update an agent access token.

To update the token name and description:

    $ lacework agent token update <token> --name dev --description "k8s deployment for dev"

To disable a token:

    $ lacework agent token update <token> --disable

To enable a token:

    $ lacework agent token update <token> --enable

```
lacework agent token update <token> [flags]
```

### Options

```
      --description string   new agent access token description
      --disable              disable agent access token
      --enable               enable agent access token
  -h, --help                 help for update
      --name string          new agent access token name
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

* [lacework agent token](lacework_agent_token.md)	 - manage agent access tokens

