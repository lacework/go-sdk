<img src="https://techally-content.s3-us-west-1.amazonaws.com/public-content/lacework_logo_full.png" width="600">

# Lacework cli

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

## Configuration File

The Lacework cli looks for a file named `.lacework.toml` inside your home
directory (`$HOME/.lacework.toml`) to access the following settings:
* `account`: Account subdomain of URL (i.e. `<ACCOUNT>.lacework.net`)
* `api_key`: API Access Key
* `api_secret`: API Access Secret


An example of a Lacework configuration file:
```toml
account = "example"
api_key = "EXAMPLE_1234567890ABC"
api_secret = "_super_secret_key"
```

You can provide a different configuration file with the option `--config`.

### Environment Variables
Default configuration parameters found in the `.lacework.toml` may also be 
overriden by setting environment variables prefixed with `LW_`. 

#### Example
To override the `account`, `api_key`, and `api_secret`  configurations:
```
$ export LW_ACCOUNT='<MY_ACCOUNT>'
$ export LW_API_KEY='<MY_API_KEY>'
$ export LW_API_SECRET='<MY_API_SECRET>'
```

## Basic Usage
Once you have created your configuration file `$HOME/.lacework.toml`,
you are ready to use the Lacework cli, a few basic commands are:

1) List all integration in your account:
```bash
$ lacework integration list
```
2) Use the `api` command to access Lacework's RestfulAPI, for example,
to get details about a specific event:
```bash
$ lacework api get '/external/events/GetEventDetails?EVENT_ID=16700'
```

## Development
To build and install the CLI from source, use the `make install-cli` directive
defined at the top level of this repository:
```
$ make install-cli
$ lacework help
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
