## lacework compliance azure

compliance for Azure Cloud

### Synopsis

Manage compliance reports for Azure Cloud.

To list all Azure Tenants configured in your account:

    $ lacework compliance azure list-tenants

To list all Azure Subscriptions from a Tenant, use the command:

    $ lacework compliance azure list-subscriptions <tenant_id>

To get the latest Azure compliance assessment report, use the command:

    $ lacework compliance azure get-report <tenant_id> <subscriptions_id>

These reports run on a regular schedule, typically once a day.

To run an ad-hoc compliance assessment use the command:

    $ lacework compliance azure run-assessment <tenant_id>


### Options

```
  -h, --help   help for azure
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

* [lacework compliance](lacework_compliance.md)	 - manage compliance reports
* [lacework compliance azure get-report](lacework_compliance_azure_get-report.md)	 - get the latest Azure compliance report
* [lacework compliance azure list-subscriptions](lacework_compliance_azure_list-subscriptions.md)	 - list subscriptions from tenant
* [lacework compliance azure list-tenants](lacework_compliance_azure_list-tenants.md)	 - list all Azure Tenants configured
* [lacework compliance azure run-assessment](lacework_compliance_azure_run-assessment.md)	 - run a new Azure compliance assessment

