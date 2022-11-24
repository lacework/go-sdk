---
title: "lacework compliance google"
slug: lacework_compliance_google
hide_title: true
---

## lacework compliance google

Compliance for Google Cloud

### Synopsis

Manage compliance reports for Google Cloud.

To list all GCP organizations and projects configured in your account:

    lacework compliance gcp list

To list all GCP projects from an organization, use the command:

    lacework compliance gcp list-projects <organization_id>

To get the latest GCP compliance assessment report, use the command:

    lacework compliance gcp get-report <organization_id> <project_id>

These reports run on a regular schedule, typically once a day.

To run an ad-hoc compliance assessment use the command:

    lacework compliance gcp run-assessment <org_or_project_id>


### Options

```
  -h, --help   help for google
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

* [lacework compliance](lacework_compliance.md)	 - Manage compliance reports
* [lacework compliance google get-report](lacework_compliance_google_get-report.md)	 - Get the latest GCP compliance report
* [lacework compliance google list](lacework_compliance_google_list.md)	 - List gcp projects and organizations
* [lacework compliance google list-projects](lacework_compliance_google_list-projects.md)	 - List projects from an organization

