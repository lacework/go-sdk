## lacework integration delete

delete an external integrations

### Synopsis

Delete an external integration by providing its integration GUID. Integration
GUIDs can be found by using the 'lacework integration list' command.

```
lacework integration delete <int_guid> [flags]
```

### Options

```
  -h, --help   help for delete
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

* [lacework integration](lacework_integration.md)	 - manage external integrations

