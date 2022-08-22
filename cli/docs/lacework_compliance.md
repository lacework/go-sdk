---
title: "lacework compliance"
slug: lacework_compliance
hide_title: true
---

## lacework compliance

Manage compliance reports

### Synopsis

Manage compliance reports for Google, Azure, or AWS cloud providers.

Lacework cloud security platform provides continuous Compliance monitoring against
cloud security best practices and compliance standards as CIS, PCI DSS, SoC II and
HIPAA benchmark standards.

Get started by integrating one or more cloud accounts using the command:

    lacework integration create

If you prefer to configure the integration via the WebUI, log in to your account at:

    https://<ACCOUNT>.lacework.net

Then navigate to Settings > Integrations > Cloud Accounts.

Use the following command to list all available integrations in your account:

    lacework integrations list


### Options

```
  -h, --help   help for compliance
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
* [lacework compliance aws](lacework_compliance_aws.md)	 - Compliance for AWS
* [lacework compliance azure](lacework_compliance_azure.md)	 - Compliance for Azure Cloud
* [lacework compliance google](lacework_compliance_google.md)	 - Compliance for Google Cloud

