---
title: "lacework agent"
slug: lacework_agent
hide_title: true
---

## lacework agent

Manage Lacework agents

### Synopsis

Manage agents and agent access tokens in your account.

To analyze application, host, and user behavior, Lacework uses a lightweight agent,
which securely forwards collected metadata to the Lacework cloud for analysis. The
agent requires minimal system resources and runs on most 64-bit Linux distributions.

For a complete list of supported operating systems, visit:

  https://docs.lacework.com/supported-operating-systems

### Options

```
  -h, --help   help for agent
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
* [lacework agent install](lacework_agent_install.md)	 - Install the datacollector agent on a remote host
* [lacework agent list](lacework_agent_list.md)	 - List all hosts with a running agent
* [lacework agent token](lacework_agent_token.md)	 - Manage agent access tokens

