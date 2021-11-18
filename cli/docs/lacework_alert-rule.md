## lacework alert-rule

manage alert rules

### Synopsis

Manage alert rules to route events to the appropriate people or tools.		
An alert rule has three parts:
  1. Alert channel(s) that should receive the event notification
  2. Event severity and categories to include
  3. Resource group(s) containing the subset of your environment to consider


### Options

```
  -h, --help   help for alert-rule
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
* [lacework alert-rule create](lacework_alert-rule_create.md)	 - create a new alert rule
* [lacework alert-rule delete](lacework_alert-rule_delete.md)	 - delete a alert rule
* [lacework alert-rule list](lacework_alert-rule_list.md)	 - list all alert rules
* [lacework alert-rule show](lacework_alert-rule_show.md)	 - show an alert rule by id

