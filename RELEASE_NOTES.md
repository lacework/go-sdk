# Release Notes
Another day, another release. These are the release notes for the version `v0.2.10`.

## Features
* feat(cli): new agent access token command (#256) (Salim Afiune)([7f8ba11](https://github.com/lacework/go-sdk/commit/7f8ba113b38ecd768f61e54ba712badf6596a587))
* feat(compliance): new aws list-accounts command (Salim Afiune Maya)([705f2eb](https://github.com/lacework/go-sdk/commit/705f2ebf9f1b9b5af2eb745c86498fe31c01e174))
## Refactor
* refactor: account mapping file for consolidated CT (#252) (Salim Afiune)([402a363](https://github.com/lacework/go-sdk/commit/402a3634765ef8c6f1f65d1be13da2ad34cf2960))
## Bug Fixes
* fix(install.sh): avoid logging with 'info' cmd (#254) (Salim Afiune)([df5f8cf](https://github.com/lacework/go-sdk/commit/df5f8cfbc7228ff9bff25e6e22a2ab68acd47fa4))
* fix: false positive results in pkg manifest scan (#255) (Salim Afiune)([a6d6cda](https://github.com/lacework/go-sdk/commit/a6d6cda9f36b38f8b653bd01ef258bd431611908))
* fix(databox): remove hardcoded LW account (Salim Afiune Maya)([c806157](https://github.com/lacework/go-sdk/commit/c80615749827c12dbfef5e1c76bf5857cd3dae7a))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
