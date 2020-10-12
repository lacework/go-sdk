## lacework access-token

generate temporary access tokens

### Synopsis

Generates a temporary access token that can be used to access the
Lacework API. The token will be valid for the duration that you specify.

```
lacework access-token [flags]
```

### Options

```
  -d, --duration_seconds int   duration in seconds that the access token should remain valid (default 3600)
  -h, --help                   help for access-token
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

