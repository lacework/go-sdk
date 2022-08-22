---
title: "lacework access-token"
slug: lacework_access-token
hide_title: true
---

## lacework access-token

Generate temporary API access tokens

### Synopsis

Generates a temporary API access token that can be used to access the
Lacework API. The token will be valid for the duration that you specify.

```
lacework access-token [flags]
```

### Options

```
  -d, --duration_seconds int   duration in seconds that the access token should remain valid (default 3600)
  -h, --help                   help for access-token
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

