---
title: "lacework compliance google list-projects"
slug: lacework_compliance_google_list-projects
hide_title: true
---

## lacework compliance google list-projects

List projects from an organization

### Synopsis

List all GCP projects from the provided organization ID.

Use the following command to list all GCP integrations in your account:

    lacework cloud-account list --type GcpCfg

Then, select one GUID from an integration and visualize its details using the command:

    lacework cloud-account show <int_guid>


```
lacework compliance google list-projects <organization_id> [flags]
```

### Options

```
  -h, --help   help for list-projects
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

* [lacework compliance google](lacework_compliance_google.md)	 - Compliance for Google Cloud

