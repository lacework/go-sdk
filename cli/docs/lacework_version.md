## lacework version

print the Lacework CLI version

### Synopsis


Prints out the installed version of the Lacework CLI and checks for newer
versions available for update.

Set the environment variable 'LW_UPDATES_DISABLE=1' to avoid checking for updates.

```
lacework version [flags]
```

### Options

```
  -h, --help   help for version
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

* [lacework](lacework.md)	 - A tool to manage the Lacework cloud security platform.

