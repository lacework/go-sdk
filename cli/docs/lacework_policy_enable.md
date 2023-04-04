---
title: "lacework policy enable"
slug: lacework_policy_enable
hide_title: true
---

## lacework policy enable

Enable policies

### Synopsis

Enable policies by ID or all policies matching a tag.

To enter the policy enable prompt:

	lacework policy enable

To enable a single policy by its ID:

	lacework policy enable lacework-policy-id

To enable many policies by ID provide a list of policy ids:

	lacework policy enable lacework-policy-id-one lacework-policy-id-two

To enable all policies for AWS CIS 1.4.0:

	lacework policy enable --tag framework:cis-aws-1-4-0

To enable all policies for GCP CIS 1.3.0:

	lacework policy enable --tag framework:cis-gcp-1-3-0



```
lacework policy enable [policy_id...] [flags]
```

### Options

```
  -h, --help         help for enable
      --tag string   enable all policies with the specified tag
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

* [lacework policy](lacework_policy.md)	 - Manage policies

