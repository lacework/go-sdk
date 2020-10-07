# Release Notes
Another day, another release. These are the release notes for the version `v0.2.4`.

### Allow custom installation directory
Now you can install the Lacework CLI on any directory you wish to install it!
```
$ curl https://raw.githubusercontent.com/lacework/go-sdk/master/cli/install.sh | bash -s -- -d $HOME/bin
```

The `install.sh` script has other flags like `-v` that allows you to rollback
to a previous version of the Lacework CLI:
```
bash: Installs the Lacework Command Line Interface.

Authors: Technology Alliances <tech-ally@lacework.net>

USAGE:
    bash [FLAGS]

FLAGS:
    -h    Prints help information
    -v    Specifies a version (ex: v0.1.0)
    -d    The installation directory (default: /usr/local/bin)
    -t    Specifies the target of the program to download (default: linux-amd64)
```

### Programatic access to profile data
We have seen a number of Lacework CLI users with a high number of configured
profiles, to help them manage their profiles we have introduced a new `list`
command that lists all profiles configured into the config file located at
`~/.lacework.toml`:
```
$ lacework configure list
```

Additionally, there is a new `show` command that allows you to programmatically
access the current computed configuration data. This makes it easy to export
the necessary environment variables for different workflows:
```
$ export LW_API_KEY="$(lacework configure show api_key --profile dev)"
```

## Features
* feat(cli): programatic access to profile data (#225) (Salim Afiune)([ab7ce7c](https://github.com/lacework/go-sdk/commit/ab7ce7cfe8e94053ca6bf8d32d929c5e748496e4))
* feat(cli): allow custom installation directory -d 📁 (#223) (Salim Afiune)([ee9e686](https://github.com/lacework/go-sdk/commit/ee9e686c46029b32e711f9534ecd7755926ec22b))
## Documentation Updates
* docs: automatically generate cli docs (#224) (Salim Afiune)([5b91e1e](https://github.com/lacework/go-sdk/commit/5b91e1e788128dd3cddf457bce565749c73eddae))
## Other Changes
* chore: add badges to README.md (#222) (Salim Afiune)([db7235d](https://github.com/lacework/go-sdk/commit/db7235d20e7af012cb8e8f3041a02728d4f28719))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
