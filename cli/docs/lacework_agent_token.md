## lacework agent token

manage agent access tokens

### Synopsis

Manage agent access tokens in your account.

Agent tokens should be treated as secret and not published. A token uniquely identifies
a Lacework customer. If you suspect your token has been publicly exposed or compromised,
generate a new token, update the new token on all machines using the old token. When
complete, the old token can safely be disabled without interrupting Lacework services.

### Options

```
  -h, --help   help for token
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

* [lacework agent](lacework_agent.md)	 - manage Lacework agents
* [lacework agent token create](lacework_agent_token_create.md)	 - create a new agent access token
* [lacework agent token list](lacework_agent_token_list.md)	 - list all agent access tokens
* [lacework agent token show](lacework_agent_token_show.md)	 - show details about an agent access token
* [lacework agent token update](lacework_agent_token_update.md)	 - update an agent access token

