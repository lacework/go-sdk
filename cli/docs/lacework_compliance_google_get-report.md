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
      --details           increase details about the compliance report
  -h, --help              help for get-report
      --pdf               download report in PDF format
      --pdf-file string   (DEPRECATED) use --pdf
      --type string       report type to display, supported types: CIS, SOC, or PCI (default "CIS")
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
```

### SEE ALSO

* [lacework compliance google](lacework_compliance_google.md)	 - compliance for Google Cloud

