---
title: "lacework policy"
slug: lacework_policy
hide_title: true
---

## lacework policy

Manage policies

### Synopsis

Manage policies in your Lacework account.

A policy is a mechanism used to add annotated metadata to a Lacework query for improving
the context of alerts, reports, and information displayed in the Lacework Console.

A policy also facilitates the scheduled execution of a Lacework query

A query is a mechanism used to interactively request information from a specific
curated dataset. A query has a defined structure for authoring detections.

Lacework ships a set of default LQL policies that are available in your account.

Limitations:
  * The maximum number of records that each policy will return is 1000
  * The maximum number of API calls is 120 per hour for ad-hoc LQL query executions

To view all the policies in your Lacework account.

    lacework policy ls

To view more details about a single policy.

    lacework policy show <policy_id>

To view the LQL query associated with the policy, use the query id shown.

    lacework query show <query_id>

**NOTE: LQL syntax may change.**


### Options

```
  -h, --help   help for policy
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

* [lacework](lacework.md)	 - A tool to manage the Lacework cloud security platform.
* [lacework policy create](lacework_policy_create.md)	 - Create a policy
* [lacework policy delete](lacework_policy_delete.md)	 - Delete a policy
* [lacework policy disable](lacework_policy_disable.md)	 - Disable Policies
* [lacework policy enable](lacework_policy_enable.md)	 - Enable Policies
* [lacework policy list](lacework_policy_list.md)	 - List policies
* [lacework policy list-tags](lacework_policy_list-tags.md)	 - List policy tags
* [lacework policy show](lacework_policy_show.md)	 - Show policy
* [lacework policy update](lacework_policy_update.md)	 - Update a policy

