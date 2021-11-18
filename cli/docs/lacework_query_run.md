## lacework query run

run a query

### Synopsis

Run an LQL query via editor:

    lacework query run --range today

Run a query via ID (uses active profile):

    lacework query run MyQuery --start "-1w@w" --end "@w"

Start and End times are required to run a query:

1.  Start and End times must be specified in one of the following formats:

    A. A relative time specifier
    B. RFC3339 Date and Time
    C. Epoch time in milliseconds

2. Start and End times must be specified in one of the following ways:

    A.  As StartTimeRange and EndTimeRange in the ParamInfo block within the query
    B.  As start_time_range and end_time_range if specifying JSON
    C.  As --start and --end CLI flags

3. Start and End time precedence:

    A.  CLI flags take precedence over JSON specifications
    B.  JSON specifications take precedence over ParamInfo specifications

```
lacework query run [query_id] [flags]
```

### Options

```
      --end string      end time for query (default "now")
  -f, --file string     path to a query to run
  -h, --help            help for run
      --range string    natural time range for query
      --start string    start time for query (default "@d")
  -u, --url string      url to a query to run
      --validate_only   validate query only (do not run)
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

* [lacework query](lacework_query.md)	 - run and manage queries

