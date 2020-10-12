## lacework

A tool to manage the Lacework cloud security platform.

### Synopsis

The Lacework Command Line Interface is a tool that helps you manage the
Lacework cloud security platform. Use it to manage compliance reports,
external integrations, vulnerability scans, and other operations.

Start by configuring the Lacework CLI with the command:

    $ lacework configure

This will prompt you for your Lacework account and a set of API access keys.

### Options

```
  -a, --account string      account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
  -k, --api_key string      access key id
  -s, --api_secret string   secret access key
      --debug               turn on debug logging
  -h, --help                help for lacework
      --json                switch commands output from human-readable to json format
      --nocolor             turn off colors
      --noninteractive      turn off interactive mode (disable spinners, prompts, etc.)
  -p, --profile string      switch between profiles configured at ~/.lacework.toml
```

### SEE ALSO

* [lacework access-token](lacework_access-token.md)	 - generate temporary access tokens
* [lacework api](lacework_api.md)	 - helper to call Lacework's RestfulAPI
* [lacework compliance](lacework_compliance.md)	 - manage compliance reports
* [lacework configure](lacework_configure.md)	 - configure the Lacework CLI
* [lacework event](lacework_event.md)	 - inspect Lacework events
* [lacework integration](lacework_integration.md)	 - manage external integrations
* [lacework version](lacework_version.md)	 - print the Lacework CLI version
* [lacework vulnerability](lacework_vulnerability.md)	 - container and host vulnerability assessments

