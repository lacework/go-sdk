## lacework compliance google list-projects

list projects from an organization

### Synopsis

List all GCP projects from the provided organization ID.

Use the following command to list all GCP integrations in your account:

    $ lacework integrations list --type GCP_CFG

Then, select one GUID from an integration and visualize its details using the command:

    $ lacework integration show <int_guid>


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
      --debug               turn on debug logging
      --json                switch commands output from human-readable to json format
      --nocolor             turn off colors
      --noninteractive      turn off interactive mode (disable spinners, prompts, etc.)
  -p, --profile string      switch between profiles configured at ~/.lacework.toml
```

### SEE ALSO

* [lacework compliance google](lacework_compliance_google.md)	 - compliance for Google Cloud

