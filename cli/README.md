<img src="https://techally-content.s3-us-west-1.amazonaws.com/public-content/lacework_logo_full.png" width="600">

# Lacework CLI

The Lacework Command Line Interface is a tool that helps you manage the
Lacework cloud security platform. You can use it to manage compliance
reports, external integrations, vulnerability scans, and other operations.

## Install

### Bash:

```
$ curl https://raw.githubusercontent.com/lacework/go-sdk/master/cli/install.sh | sudo bash
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

>To create a set of API keys, log in to your Lacework account and navigate to
>Settings -> API Keys, then click + Create New. Enter a name for the key and
>an optional description and click Save. To get the secret key, download the
>generated API key file and open it in an editor.

The following example shows sample values. Replace them with your own.

```bash
$ lacework configure
Account: example
Access Key ID: EXAMPLE_1234567890ABCDE1EXAMPLE1EXAMPLE123456789EXAMPLE
Secret Access Key: _d12345dcbde000d1232bbbe51234a609

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
Secret Access Key: _12345prode11111232bbbe51234a609

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

This is a list of all environment variables that can be user to modify the
operation of the Lacework CLI.

| Environment Variable | Description |
|----------------------|-------------|
|`LW_NOCOLOR=1`|turn off colors|
|`LW_DEBUG=1`|turn on debug logging|
|`LW_JSON=1`|switch commands output from human-readable to JSON format|
|`LW_NONINTERACTIVE=1`|disable interactive progress bars (i.e. spinners)|
|`LW_PROFILE="<name>"`|switch between profiles configured at `~/.lacework.toml`|
|`LW_ACCOUNT="<account>"`|account subdomain of URL (i.e. `<ACCOUNT>.lacework.net`)|
|`LW_API_KEY="<key>"`|access key id|
|`LW_API_SECRET="<secret>"`|secret access key|

## Basic Usage
A few basic commands are:

1) List all integration in your account:
```bash
$ lacework integration list
```
2) Request an on-demand vulnerability scan:
```bash
$ lacework vulnerability scan run index.docker.io techallylw/lacework-cli latest
```
3) Use the `api` command to access Lacework's RestfulAPI, for example,
to get details about a specific event:
```bash
$ lacework api get '/external/events/GetEventDetails?EVENT_ID=16700'
```

## Development

To build and install the CLI from source, use the `make install-cli` directive
defined at the top level of this repository, the automation will install the
tool at `/usr/local/bin/lacework`:
```
$ make install-cli
$ lacework version
lacework v.dev (sha:ca9f95d17f4f2092f89dba7b64eaed6db7493a5a) (time:20200406091143)
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
