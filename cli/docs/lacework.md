---
title: "lacework"
slug: lacework
---

## lacework

A tool to manage the Lacework cloud security platform.

### Synopsis

The Lacework Command Line Interface is a tool that helps you manage the
Lacework cloud security platform. Use it to manage compliance reports,
external integrations, vulnerability scans, and other operations.

Start by configuring the Lacework CLI with the command:

    lacework configure

This will prompt you for your Lacework account and a set of API access keys.

### Options

```
  -a, --account string      account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
  -k, --api_key string      access key id
  -s, --api_secret string   secret access key
      --api_token string    access token (replaces the use of api_key and api_secret)
      --debug               turn on debug logging
  -h, --help                help for lacework
      --json                switch commands output from human-readable to json format
      --nocache             turn off caching
      --nocolor             turn off colors
      --noninteractive      turn off interactive mode (disable spinners, prompts, etc.)
      --organization        access organization level data sets (org admins only)
  -p, --profile string      switch between profiles configured at ~/.lacework.toml
      --subaccount string   sub-account name inside your organization (org admins only)
```

### SEE ALSO

* [lacework access-token](/cli/commands/lacework_access-token/)	 - generate temporary API access tokens
* [lacework account](/cli/commands/lacework_account/)	 - manage accounts in an organization (org admins only)
* [lacework agent](/cli/commands/lacework_agent/)	 - manage Lacework agents
* [lacework alert-rule](/cli/commands/lacework_alert-rule/)	 - manage alert rules
* [lacework api](/cli/commands/lacework_api/)	 - helper to call Lacework's API
* [lacework compliance](/cli/commands/lacework_compliance/)	 - manage compliance reports
* [lacework configure](/cli/commands/lacework_configure/)	 - configure the Lacework CLI
* [lacework event](/cli/commands/lacework_event/)	 - inspect Lacework events
* [lacework integration](/cli/commands/lacework_integration/)	 - manage external integrations
* [lacework policy](/cli/commands/lacework_policy/)	 - manage policies
* [lacework query](/cli/commands/lacework_query/)	 - run and manage queries
* [lacework resource-group](/cli/commands/lacework_resource-group/)	 - manage resource groups
* [lacework version](/cli/commands/lacework_version/)	 - print the Lacework CLI version
* [lacework vulnerability](/cli/commands/lacework_vulnerability/)	 - container and host vulnerability assessments

