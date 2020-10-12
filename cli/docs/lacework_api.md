## lacework api

helper to call Lacework's RestfulAPI

### Synopsis

Use this command as a helper to call any available Lacework API endpoint.

For example, to list all integrations configured in your account run:

    lacework api get /external/integrations

For a complete list of available API endpoints visit:

    https://<ACCOUNT>.lacework.net/api/v1/external/docs


```
lacework api <method> <path> [flags]
```

### Options

```
  -d, --data string   data to send only for post and patch requests
  -h, --help          help for api
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

