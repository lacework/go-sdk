---
title: "lacework"
slug: lacework
hide_title: true
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

* [lacework access-token](lacework_access-token.md)	 - Generate temporary API access tokens
* [lacework account](lacework_account.md)	 - Manage accounts in an organization (org admins only)
* [lacework agent](lacework_agent.md)	 - Manage Lacework agents
* [lacework alert-rule](lacework_alert-rule.md)	 - Manage alert rules
* [lacework api](lacework_api.md)	 - Helper to call Lacework's API
* [lacework compliance](lacework_compliance.md)	 - Manage compliance reports
* [lacework configure](lacework_configure.md)	 - Configure the Lacework CLI
* [lacework event](lacework_event.md)	 - Inspect Lacework events
* [lacework integration](lacework_integration.md)	 - Manage external integrations
* [lacework policy](lacework_policy.md)	 - Manage policies
* [lacework query](lacework_query.md)	 - Run and manage queries
* [lacework resource-group](lacework_resource-group.md)	 - Manage resource groups
* [lacework version](lacework_version.md)	 - Print the Lacework CLI version
* [lacework vulnerability](lacework_vulnerability.md)	 - Container and host vulnerability assessments

