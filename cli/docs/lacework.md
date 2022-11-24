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
* [lacework alert](lacework_alert.md)	 - Inspect and manage alerts
* [lacework alert-channel](lacework_alert-channel.md)	 - Manage alert channels
* [lacework alert-profile](lacework_alert-profile.md)	 - Manage alert profiles
* [lacework alert-rule](lacework_alert-rule.md)	 - Manage alert rules
* [lacework api](lacework_api.md)	 - Helper to call Lacework's API
* [lacework cloud-account](lacework_cloud-account.md)	 - Manage cloud accounts
* [lacework compliance](lacework_compliance.md)	 - Manage compliance reports
* [lacework configure](lacework_configure.md)	 - Configure the Lacework CLI
* [lacework container-registry](lacework_container-registry.md)	 - Manage container registries
* [lacework generate](lacework_generate.md)	 - Generate code to onboard your account
* [lacework policy](lacework_policy.md)	 - Manage policies
* [lacework policy-exception](lacework_policy-exception.md)	 - Manage policy exceptions
* [lacework query](lacework_query.md)	 - Run and manage queries
* [lacework report-definition](lacework_report-definition.md)	 - Manage report definitions
* [lacework report-rule](lacework_report-rule.md)	 - Manage report rules
* [lacework resource-group](lacework_resource-group.md)	 - Manage resource groups
* [lacework team-member](lacework_team-member.md)	 - Manage team members
* [lacework version](lacework_version.md)	 - Print the Lacework CLI version
* [lacework vulnerability](lacework_vulnerability.md)	 - Container and host vulnerability assessments
* [lacework vulnerability-exception](lacework_vulnerability-exception.md)	 - Manage vulnerability exceptions

