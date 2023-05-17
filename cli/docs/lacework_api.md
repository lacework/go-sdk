---
title: "lacework api"
slug: lacework_api
hide_title: true
---

## lacework api

Helper to call Lacework's API

### Synopsis

Use this command as a helper to call any available Lacework API v2 endpoint.

### API v2

To list all available Lacework schema types:

    lacework api get /v2/schemas

To receive a json response of all machines within the given time window:

    lacework api post /api/v2/Entities/Machines/search -d "{}"

To receive a json response of all agents within the given time window:

    lacework api post /api/v2/AgentInfo/search -d "{}"

For a complete list of available API v2 endpoints visit:

    https://<ACCOUNT>.lacework.net/api/v2/docs


```
lacework api <method> <path> [flags]
```

### Options

```
  -d, --data string   data to send only for post and patch requests
  -h, --help          help for api
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

