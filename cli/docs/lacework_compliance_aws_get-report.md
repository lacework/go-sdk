---
title: "lacework compliance aws get-report"
slug: lacework_compliance_aws_get-report
hide_title: true
---

## lacework compliance aws get-report

Get the latest AWS compliance report

### Synopsis

Get the latest compliance assessment report from the provided AWS account, these
reports run on a regular schedule, typically once a day. The available report formats
are human-readable (default), json and pdf.

To list all AWS accounts configured in your account:

    lacework compliance aws list-accounts

To show recommendation details and affected resources for a recommendation id:

    lacework compliance aws get-report <account_id> [recommendation_id]

To retrieve a specific report by its report name:

    lacework compliance aws get-report <account_id> --report_name 'AWS CSA CCM 4.0.5'


```
lacework compliance aws get-report <account_id> [recommendation_id] [flags]
```

### Options

```
      --category strings     filter report details by category (identity-and-access-management, s3, logging...)
      --csv                  output report in CSV format
      --details              increase details about the compliance report
  -h, --help                 help for get-report
      --pdf                  download report in PDF format
      --report_name string   report name to display, run 'lacework report-definitions list' for more information.
      --service strings      filter report details by service (aws:s3, aws:iam, aws:cloudtrail, ...)
      --severity string      filter report details by severity threshold (critical, high, medium, low, info)
      --status string        filter report details by status (non-compliant, requires-manual-assessment, suppressed, compliant, could-not-assess)
      --type string          report type to display, run 'lacework report-definitions list' for more information.
                             valid types:
                             'AWS_CIS_14','AWS_CIS_1_4_ISO_IEC_27002_2022','AWS_CIS_S3','AWS_CMMC_1.02','AWS_CSA_CCM_4_0_5',
                             'AWS_Cyber_Essentials_2_2','AWS_HIPAA','AWS_ISO_27001:2013','AWS_NIST_800-171_rev2','AWS_NIST_800-53_rev5',
                             'AWS_NIST_CSF','AWS_PCI_DSS_3.2.1','AWS_SOC_2','AWS_SOC_Rev2','HIPAA',
                             'ISO_2700','LW_AWS_SEC_ADD_1_0','NIST_800-171_Rev2','NIST_800-53_Rev4','PCI',
                             'SOC', (default "AWS_CIS_14")
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

* [lacework compliance aws](lacework_compliance_aws.md)	 - Compliance for AWS

