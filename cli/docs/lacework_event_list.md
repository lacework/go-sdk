## lacework event list

list all events (default last 7 days)

### Synopsis

List all events for the last 7 days by default, or pass --start and --end to
specify a custom time period. You can also pass --serverity to filter by a
severity threshold.

Additionally, pass --days to list events for a specified number of days.

For example, to list all events from the last day with severity medium and above
(Critical, High and Medium) run:

    $ lacework events list --severity medium --days 1

```
lacework event list [flags]
```

### Options

```
      --days int          list events for specified number of days (max: 7 days)
      --end string        end of the time range in UTC (format: yyyy-MM-ddTHH:mm:ssZ)
  -h, --help              help for list
      --severity string   filter events by severity threshold (critical, high, medium, low, info)
      --start string      start of the time range in UTC (format: yyyy-MM-ddTHH:mm:ssZ)
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
  -u, --subaccount string   sub-account name inside your organization (org admins only)
```

### SEE ALSO

* [lacework event](lacework_event.md)	 - inspect Lacework events

