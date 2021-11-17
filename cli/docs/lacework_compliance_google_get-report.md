---
title: "lacework compliance google get-report"
slug: lacework_compliance_google_get-report
---

## lacework compliance google get-report

get the latest GCP compliance report

### Synopsis

Get the latest compliance assessment report, these reports run on a regular schedule,
typically once a day. The available report formats are human-readable (default), json and pdf.

To run an ad-hoc compliance assessment use the command:

    $ lacework compliance gcp run-assessment <project_id>


```
lacework compliance google get-report <organization_id> <project_id> [flags]
```

### Options

```
      --category strings   filter report details by category (storage, networking, identity-and-access-management, ...)
      --csv                output report in CSV format
      --details            increase details about the compliance report
  -h, --help               help for get-report
      --pdf                download report in PDF format
      --service strings    filter report details by service (gcp:storage:bucket, gcp:kms:cryptoKey, gcp:project, ...)
      --severity string    filter report details by severity threshold (critical, high, medium, low, info)
      --status string      filter report details by status (non-compliant, requires-manual-assessment, suppressed, compliant, could-not-assess)
      --type string        report type to display, supported types: CIS, CIS12, K8S, HIPAA, SOC, or PCI (default "CIS")
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

