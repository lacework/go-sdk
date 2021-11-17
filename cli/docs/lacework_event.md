---
title: "lacework event"
slug: lacework_event
---

## lacework event

inspect Lacework events

### Synopsis

Inspect events reported by the Lacework platform

### Options

```
  -h, --help   help for event
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
* [lacework event list](/cli/commands/lacework_event_list/)	 - list all events (default last 7 days)
* [lacework event open](/cli/commands/lacework_event_open/)	 - open a specified event in a web browser
* [lacework event show](/cli/commands/lacework_event_show/)	 - show details about a specific event

