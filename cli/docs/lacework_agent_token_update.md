## lacework agent token update

update an agent access token

### Synopsis

Update an agent access token.

To update the token name and description:

    $ lacework agent token update <token> --name dev --description "k8s deployment for dev"

To disable a token:

    $ lacework agent token update <token> --disable

To enable a token:

    $ lacework agent token update <token> --enable

```
lacework agent token update <token> [flags]
```

### Options

```
      --description string   new agent access token description
      --disable              disable agent access token
      --enable               enable agent access token
  -h, --help                 help for update
      --name string          new agent access token name
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

* [lacework agent token](lacework_agent_token.md)	 - manage agent access tokens

