## lacework compliance aws get-report

get the latest AWS compliance report

### Synopsis

Get the latest compliance assessment report from the provided AWS account, these
reports run on a regular schedule, typically once a day. The available report formats
are human-readable (default), json and pdf.

To list all AWS accounts configured in your account:

    $ lacework compliance aws list-accounts

To run an ad-hoc compliance assessment of an AWS account:

    $ lacework compliance aws run-assessment <account_id>


```
lacework compliance aws get-report <account_id> [flags]
```

### Options

```
      --details           increase details about the compliance report
  -h, --help              help for get-report
      --pdf               download report in PDF format
      --pdf-file string   (DEPRECATED) use --pdf
      --type string       report type to display, supported types: CIS, NIST_800-53_Rev4, ISO_2700, HIPAA, SOC, or PCI (default "CIS")
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

* [lacework compliance aws](lacework_compliance_aws.md)	 - compliance for AWS

