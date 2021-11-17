---
title: "lacework compliance google run-assessment"
slug: lacework_compliance_google_run-assessment
---

## lacework compliance google run-assessment

run a new GCP compliance assessment

### Synopsis

Run a compliance assessment for the provided GCP organization or project.

```
lacework compliance google run-assessment <org_or_project_id> [flags]
```

### Options

```
  -h, --help   help for run-assessment
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

* [lacework compliance google](lacework_compliance_google.md)	 - compliance for Google Cloud

