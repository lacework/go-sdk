---
title: "lacework configure show"
slug: lacework_configure_show
hide_title: true
---

## lacework configure show

Show current configuration data

### Synopsis

Prints the current computed configuration data from the specified configuration
key. The order of precedence to compute the configuration is flags, environment
variables, and the configuration file ~/.lacework.toml. 

The available configuration keys are:

* profile
* account
* subaccount
* api_secret
* api_key

To show the configuration from a different profile, use the flag --profile.

    lacework configure show account --profile my-profile

```
lacework configure show <config_key> [flags]
```

### Options

```
  -h, --help   help for show
```

### Options inherited from parent commands

```
  -a, --account string      account URL (i.e. <ACCOUNT>[.CUSTER][.corp].lacework.net)
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

* [lacework configure](lacework_configure.md)	 - Configure the Lacework CLI

