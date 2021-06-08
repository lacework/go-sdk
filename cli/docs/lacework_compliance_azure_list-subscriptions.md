## lacework compliance azure list-subscriptions

list subscriptions from tenant

### Synopsis

List all Azure subscriptions from the provided Tenant ID.

Use the following command to list all Azure Tenants configured in your account:

    $ lacework compliance az list

```
lacework compliance azure list-subscriptions <tenant_id> [flags]
```

### Options

```
  -h, --help   help for list-subscriptions
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
      --subaccount string   sub-account name inside your organization (org admins only)
```

### SEE ALSO

* [lacework compliance azure](lacework_compliance_azure.md)	 - compliance for Azure Cloud

