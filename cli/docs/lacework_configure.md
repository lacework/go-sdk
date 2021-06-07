## lacework configure

configure the Lacework CLI

### Synopsis

Configure settings that the Lacework CLI uses to interact with the Lacework
platform. These include your Lacework account, API access key and secret.

To create a set of API keys, log in to your Lacework account via WebUI and
navigate to Settings > API Keys and click + Create New. Enter a name for
the key and an optional description, then click Save. To get the secret key,
download the generated API key file.

Use the flag --json_file to preload the downloaded API key file.

If this command is run with no flags, the Lacework CLI will store all
settings under the default profile. The information in the default profile
is used any time you run a Lacework CLI command that doesn't explicitly
specify a profile to use.

You can configure multiple profiles by using the --profile flag. If a
config file does not exist (the default location is ~/.lacework.toml),
the Lacework CLI will create it for you.

```
lacework configure [flags]
```

### Options

```
  -h, --help               help for configure
  -j, --json_file string   loads the API key JSON file downloaded from the WebUI
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

* [lacework](lacework.md)	 - A tool to manage the Lacework cloud security platform.
* [lacework configure list](lacework_configure_list.md)	 - list all configured profiles at ~/.lacework.toml
* [lacework configure show](lacework_configure_show.md)	 - show current configuration data

