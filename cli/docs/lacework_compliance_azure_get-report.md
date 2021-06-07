## lacework compliance azure get-report

get the latest Azure compliance report

### Synopsis

Get the latest Azure compliance assessment report, these reports run on a regular schedule,
typically once a day. The available report formats are human-readable (default), json and pdf.

To run an ad-hoc compliance assessment use the command:

    $ lacework compliance azure run-assessment <tenant_id>


```
lacework compliance azure get-report <tenant_id> <subscriptions_id> [flags]
```

### Options

```
      --category strings   filter report details by category (networking, storage, ...)
      --csv                output report in CSV format
      --details            increase details about the compliance report
  -h, --help               help for get-report
      --pdf                download report in PDF format
      --service strings    filter report details by service (azure:ms:storage, azure:ms:sql, azure:ms:network, ...)
      --severity string    filter report details by severity threshold (critical, high, medium, low, info)
      --status string      filter report details by status (non-compliant, requires-manual-assessment, suppressed, compliant, could-not-assess)
      --type string        report type to display, supported types: CIS, SOC, or PCI (default "CIS")
```

### Options inherited from parent commands

```
  -a, --account string      account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
  -k, --api_key string      access key id
  -s, --api_secret string   secret access key
      --debug               turn on debug logging
      --json                switch commands output from human-readable to json format
      --nocolor             turn off colors
      --noninteractive      turn off interactive mode (disable spinners, prompts, etc.)
  -p, --profile string      switch between profiles configured at ~/.lacework.toml
  -u, --subaccount string   sub-account name inside your organization (org admins only)
```

### SEE ALSO

* [lacework compliance azure](lacework_compliance_azure.md)	 - compliance for Azure Cloud

