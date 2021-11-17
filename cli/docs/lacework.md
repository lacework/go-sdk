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

* [lacework access-token](lacework_access-token.md)	 - generate temporary API access tokens
* [lacework account](lacework_account.md)	 - manage accounts in an organization (org admins only)
* [lacework agent](lacework_agent.md)	 - manage Lacework agents
* [lacework alert-rule](lacework_alert-rule.md)	 - manage alert rules
* [lacework api](lacework_api.md)	 - helper to call Lacework's API
* [lacework compliance](lacework_compliance.md)	 - manage compliance reports
* [lacework configure](lacework_configure.md)	 - configure the Lacework CLI
* [lacework event](lacework_event.md)	 - inspect Lacework events
* [lacework integration](lacework_integration.md)	 - manage external integrations
* [lacework policy](lacework_policy.md)	 - manage policies
* [lacework query](lacework_query.md)	 - run and manage queries
* [lacework resource-group](lacework_resource-group.md)	 - manage resource groups
* [lacework version](lacework_version.md)	 - print the Lacework CLI version
* [lacework vulnerability](lacework_vulnerability.md)	 - container and host vulnerability assessments

