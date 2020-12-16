<img src="https://techally-content.s3-us-west-1.amazonaws.com/public-content/lacework_logo_full.png" width="600">

# Lacework CLI

The Lacework Command Line Interface is a tool that helps you manage the
Lacework cloud security platform. You can use it to manage compliance
reports, external integrations, vulnerability scans, and other operations.

🐳 [CLI Docker Containers](https://hub.docker.com/r/lacework/lacework-cli)

## Installation

### Bash:

```
$ curl https://raw.githubusercontent.com/lacework/go-sdk/master/cli/install.sh | bash
```

### Powershell:

```
C:\> Set-ExecutionPolicy Bypass -Scope Process -Force
C:\> iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/lacework/go-sdk/master/cli/install.ps1'))
```

## Quick Configuration

The `lacework configure` command is the fastest way to set up your Lacework
CLI installation. The command prompts you for three things:

* `account`: Account subdomain of URL (i.e. `<ACCOUNT>.lacework.net`)
* `api_key`: API Access Key
* `api_secret`: API Access Secret

>To create a set of API keys, log in to your Lacework account via WebUI and
>navigate to Settings > API Keys and click + Create New. Enter a name for
>the key and an optional description, then click Save. To get the secret key,
>download the generated API key file.

_**NOTE:** Use the argument `--json_file` to preload the downloaded API key file._

The following example shows sample values. Replace them with your own.

```bash
$ lacework configure
Account: example
Access Key ID: EXAMPLE_1234567890ABCDE1EXAMPLE1EXAMPLE123456789EXAMPLE
Secret Access Key: **********************************

You are all set!
```

The result of this command is the generation of a file named `.lacework.toml`
inside your home directory (`$HOME/.lacework.toml`) with a single profile
named `default`.

### Multiple Profiles
You can add additional profiles that you can refer to with a name by specifying
the `--profile` option. The following example creates a profile named `prod`.

```bash
$ lacework configure --profile prod
Account: prod.example
Access Key ID: PROD_1234567890ABCDE1EXAMPLE1EXAMPLE123456789EXAMPLE
Secret Access Key: **********************************

You are all set!
```

Then, when you run a command, you can specify a `--profile prod` and use the
credentials and settings stored under that name.

```bash
$ lacework integration list --profile prod
```

If there is no `--profile` option, the CLI will default to the `default` profile.

### Environment Variables
Default configuration parameters found in the `.lacework.toml` may also be 
overriden by setting environment variables prefixed with `LW_`. 

To override the `account`, `api_key`, and `api_secret`  configurations:
```
$ export LW_ACCOUNT="<YOUR_ACCOUNT>"
$ export LW_API_KEY="<YOUR_API_KEY>"
$ export LW_API_SECRET="<YOUR_API_SECRET>"
```

To override the profile to use:
```
$ export LW_PROFILE=prod
```

This is a list of all environment variables that can be used to modify the
operation of the Lacework CLI.

| Environment Variable | Description |
|----------------------|-------------|
|`LW_NOCOLOR=1`|turn off colors|
|`LW_DEBUG=1`|turn on debug logging|
|`LW_JSON=1`|switch commands output from human-readable to JSON format|
|`LW_UPDATES_DISABLE=1`|disable daily version checks|
|`LW_NONINTERACTIVE=1`|disable interactive progress bars (i.e. spinners)|
|`LW_PROFILE="<name>"`|switch between profiles configured at `~/.lacework.toml`|
|`LW_ACCOUNT="<account>"`|account subdomain of URL (i.e. `<ACCOUNT>.lacework.net`)|
|`LW_API_KEY="<key>"`|access key id|
|`LW_API_SECRET="<secret>"`|secret access key|

## Basic Usage
A few basic commands are:

1) List all integration in your account:
```bash
$ lacework integrations list
```
2) List all events from the last 7 days in your account:
```bash
$ lacework events list
```
3) Request an on-demand container vulnerability scan:
```bash
$ lacework vulnerability container scan index.docker.io lacework/lacework-cli latest
```
4) Use the `api` command to access Lacework's RestfulAPI, for example,
to look at the SCHEMA of the `WEBHOOK` integration:
```bash
$ lacework api get /external/integrations/schema/WEBHOOK
```

## CLI Documentation
For more CLI documentation, see https://github.com/lacework/go-sdk/wiki/CLI-Documentation.

## Development

To build and install the CLI from source, use the `make install-cli` directive
defined at the top level of this repository, the automation will install the
tool at `/usr/local/bin/lacework`:
```
$ make install-cli
$ lacework version
lacework 0.1.1-dev (sha:ca9f95d17f4f2092f89dba7b64eaed6db7493a5a) (time:20200406091143)
```

### Unit Tests

Running unit tests should be as simple as executing the `make test` directive.

### Integration Tests

The integration tests are end-to-end tests that are run against a real Lacework API
Server, for that reason it requires a set of Lacework API keys, to run these tests
locally you need to setup the following environment variables and use the directive
`make integration`, an example of the command you can use is:
```
$ CI_ACCOUNT="<YOUR_ACCOUNT>" CI_API_KEY="<YOUR_API_KEY>" CI_API_SECRET="<YOUR_API_SECRET>" make integration
```

## License and Copyright
Copyright 2020, Lacework Inc.
```
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
