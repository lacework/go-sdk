## lacework compliance aws

compliance for AWS

### Synopsis

Manage compliance reports for Amazon Web Services (AWS).

To list all AWS accounts configured in your account:

    $ lacework compliance aws list-accounts

To get the latest AWS compliance assessment report:

    $ lacework compliance aws get-report <account_id>

These reports run on a regular schedule, typically once a day.

To run an ad-hoc compliance assessment:

    $ lacework compliance aws run-assessment <account_id>


### Options

```
  -h, --help   help for aws
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

* [lacework compliance](lacework_compliance.md)	 - manage compliance reports
* [lacework compliance aws get-report](lacework_compliance_aws_get-report.md)	 - get the latest AWS compliance report
* [lacework compliance aws list-accounts](lacework_compliance_aws_list-accounts.md)	 - list all AWS accounts configured
* [lacework compliance aws run-assessment](lacework_compliance_aws_run-assessment.md)	 - run a new AWS compliance report

