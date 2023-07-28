<img src="https://techally-content.s3-us-west-1.amazonaws.com/public-content/lacework_logo_full.png" width="600">

# Lacework CLI

The Lacework Command Line Interface is a tool that helps you manage the
Lacework cloud security platform. You can use it to manage compliance
reports, external integrations, vulnerability scans, and other operations.

🐳 [CLI Docker Containers](https://hub.docker.com/r/lacework/lacework-cli)

## Installation

### Bash:

```
curl https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.sh | bash
```

### Powershell:

```
Set-ExecutionPolicy Bypass -Scope Process -Force;
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.ps1'))
```

### Homebrew:
```
brew install lacework/tap/lacework-cli
```

### Chocolatey:
```
choco install lacework-cli
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
lacework cloud-account list --profile prod
```

If there is no `--profile` option, the CLI will default to the `default` profile.

### Environment Variables
Default configuration parameters found in the `.lacework.toml` may also be 
overriden by setting environment variables prefixed with `LW_`. 

To override the `account`, `api_key`, and `api_secret`  configurations:
```
export LW_ACCOUNT="<YOUR_ACCOUNT>"
export LW_API_KEY="<YOUR_API_KEY>"
export LW_API_SECRET="<YOUR_API_SECRET>"
```

To override the profile to use:
```
export LW_PROFILE=prod
```

This is a list of all environment variables that can be used to modify the
operation of the Lacework CLI.

| Environment Variable | Description |
|----------------------|-------------|
|`LW_NOCOLOR=1`|turn off colors|
|`LW_NOCACHE=1`|turn off caching|
|`LW_DEBUG=1`|turn on debug logging|
|`LW_JSON=1`|switch commands output from human-readable to JSON format|
|`LW_NONINTERACTIVE=1`|disable interactive progress bars (i.e. spinners)|
|`LW_UPDATES_DISABLE=1`|disable daily version checks|
|`LW_TELEMETRY_DISABLE=1`|disable sending telemetry data|
|`LW_PROFILE="<name>"`|switch between profiles configured at `~/.lacework.toml`|
|`LW_ACCOUNT="<account>"`|account subdomain of URL (i.e. `<ACCOUNT>.lacework.net`)|
|`LW_SUBACCOUNT="<subaccount>"`|sub-account name inside your organization|
|`LW_API_KEY="<key>"`|API access key id|
|`LW_API_SECRET="<secret>"`|API secret access key|
|`LW_CDK_TARGET="<localhost:port>"`|address to dial the CDK server|

## Basic Usage
A few basic commands are:

1) List all cloud account integrations in your account:
```bash
lacework cloud-account list
```
2) List all events from the last 7 days in your account:
```bash
lacework events list
```
3) Request an on-demand container vulnerability scan:
```bash
lacework vulnerability container scan index.docker.io lacework/lacework-cli latest
```
4) Use the `api` command to access Lacework API, for example,
to list all available SCHEMAS in API v2:
```bash
lacework api get /schemas
```

## CLI Documentation
For more CLI documentation, see http://docs.lacework.com/cli

## Development

To build and install the CLI from source, use the `make install-cli` directive
defined at the top level of this repository, the automation will install the
tool at `/usr/local/bin/lacework`:
```
$ make prepare
$ make install-cli
$ lacework version
lacework 0.1.1-dev (sha:ca9f95d17f4f2092f89dba7b64eaed6db7493a5a) (time:20200406091143)
```

### Test Supported Platforms

The Lacework CLI runs on almost any operating system out there, it runs on Darwin,
Windows, and many Linux distributions. After setting up your development environment
you can test the generated binary by standing up a virtual machine of any supported
platform, to do that, you will need to install [Vagrant](https://www.vagrantup.com/) and
[VirtualBox](https://www.virtualbox.org/wiki/Downloads) on your workstation.

To start a supported host, run:
```
$ make vagrant-windows-up  # Stand up a Windows 10 VM
$ make vagrant-macos-up    # Stand up a Macos Sierra VM
$ make vagrant-linux-up    # Stand up a Ubuntu 18.04 VM
```

To access the VM, run:
```
$ make vagrant-windows-login  # Access the Windows 10 VM via WinRM/Powershell
$ make vagrant-macos-login    # Access the Macos Sierra VM via SSH
$ make vagrant-linux-login    # Access the Ubuntu 18.04 VM via SSH
```
__NOTE: When accessing a Windows VM from a Linux or MacOS system, you will need
to use the VirtualBox GUI rather than your terminal.__

To destroy the VM, run:
```
$ make vagrant-windows-destroy  # Destroy the Windows 10 VM
$ make vagrant-macos-destroy    # Destroy the Macos Sierra VM
$ make vagrant-linux-destroy    # Destroy the Ubuntu 18.04 VM
```

### Unit Tests

Running unit tests should be as simple as executing the `make test` directive.

### Integration Tests

The integration tests are end-to-end tests that run against a real Lacework API
Server, for that reason, it requires a set of Lacework API keys. To run these tests
locally you need to setup the following environment variables and use the directive
`make integration`, an example of the command you can use is:
```
CI_ACCOUNT="<YOUR_ACCOUNT>" \
  CI_SUBACCOUNT="<YOUR_SUBACCOUNT_IF_ANY>" \
  CI_API_KEY="<YOUR_API_KEY>" \
  CI_API_SECRET="<YOUR_API_SECRET>" \
  LW_INT_TEST_AWS_ACC="<YOUR_AWS_ACCOUNT>" make integration
```
This is a list of all environment variables used in the running the integration tests.

| Environment Variable | Description |
|----------------------|-------------|
|`CI_ACCOUNT="<YOUR_ACCOUNT>"` | account subdomain of URL (i.e. `<ACCOUNT>.lacework.net`)|
|`CI_SUBACCOUNT="<YOUR_ACCOUNT>"` | (orgs only) a sub-account|
|`CI_API_KEY="<YOUR_ACCOUNT>"` | API access key id|
|`CI_API_SECRET="<YOUR_ACCOUNT>"` | API secret access key|
|`LW_INT_TEST_AWS_ACC="<secret>"`|AWS Account for integration tests|
|`CI_STANDALONE_ACCOUNT=<bool>`|set to `true` if the Lacework account is not an organization|

#### Running Specific Integration Tests (RegEx)

When working on new tests or existing tests, you can use a regex to run
only specific integration tests. For example, to run only the tests related
to the command `lacework query update` use the command:

```
make integration regex=TestQueryUpdate
```

**Note that it is a best practice to follow a naming convention where we name
test functions after their actual commands so that we can use these patterns.**

This command will match the regex `TestQueryUpdate*` and will execute any integration
test that matches that pattern. For more information about what regular expressions you
can use, visit https://pkg.go.dev/cmd/go#hdr-Testing_flags.

**TIP:** If you are NOT modifying the CLI code and instead, you are only updating
the integration tests, you can use the directive `make integration-only` instead to
avoid rebuilding the CLI binary.

### Telemetry via Honeycomb

We use [Honeycomb](https://www.honeycomb.io/) for observability, to enable sending
tracing data to our development dataset, you must configure the environment variable
`HONEYAPIKEY`. This variable as well as the above CI environment variables can be
configured inside your bash profile (or any other shell profile you prefer).

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
