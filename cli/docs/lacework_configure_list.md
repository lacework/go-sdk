## lacework configure list

list all configured profiles at ~/.lacework.toml

### Synopsis

List all profiles configured into the config file ~/.lacework.toml

To switch to a different profile permanently in your current terminal,
export the environment variable:

    export LW_PROFILE="my-profile"

```
lacework configure list [flags]
```

### Options

```
  -h, --help   help for list
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

* [lacework configure](lacework_configure.md)	 - configure the Lacework CLI

