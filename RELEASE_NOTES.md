# Release Notes
Another day, another release. These are the release notes for the version `v0.17.0`.

## Features
* feat(api): Migrate VictorOps alert channel to V2 (#573) (Darren)([6bc3b88](https://github.com/lacework/go-sdk/commit/6bc3b881ab72d5d767ae54fe5b9c14dc9dc0a17a))
* feat(api): Webhook alert channel(v2) (#565) (Darren)([06915a0](https://github.com/lacework/go-sdk/commit/06915a0271330abfedce23606f04c5a1cf16a077))
## Refactor
* refactor(cli): use APIv2 to create Datadog integrations + update api/examples (#571) (vatasha)([9d6ea0e](https://github.com/lacework/go-sdk/commit/9d6ea0eba0fc34ef594fac141d55cded3b0d7ea6))
## Bug Fixes
* fix(cli): disable cache mechanism in install.sh (#563) (Salim Afiune)([0ccd32e](https://github.com/lacework/go-sdk/commit/0ccd32e682e99eda369121c97eb7d6c695ddc9e8))
* fix:Add issue grouping field to cloudwatch alert channel data (#566) (vatasha)([1e93e91](https://github.com/lacework/go-sdk/commit/1e93e91beb7a0cd95c49b211486b1ae032a0b2ae))
* fix(cli): install.ps1 avoid fail for False boolean (#564) (Larry Eichenbaum)([32604df](https://github.com/lacework/go-sdk/commit/32604df46c095bf8efa28dcc1b32c0fef14b1e16))
## Documentation Updates
* docs: improve developer guidelines (Salim Afiune Maya)([80f8e9e](https://github.com/lacework/go-sdk/commit/80f8e9e2f62ce0a567883824eb2105866ddb274a))
* docs: add signed commits documentation (Salim Afiune Maya)([b8b6325](https://github.com/lacework/go-sdk/commit/b8b63257de8a8ae2c7b0fc6b4fa7711eaa87ff70))
* docs: fix broken link in CONTRIBUTING.md (Salim Afiune Maya)([92a4d0c](https://github.com/lacework/go-sdk/commit/92a4d0cfcdd208305f1c4b3c321cd6639fd5dda1))
* docs: improve pull request template (Salim Afiune Maya)([967f9d6](https://github.com/lacework/go-sdk/commit/967f9d6d642b563b043ed0b6207ed09374d86523))
## Other Changes
* chore(deps): bump go.uber.org/zap from 1.19.0 to 1.19.1 (#568) (dependabot[bot])([3f4a1f2](https://github.com/lacework/go-sdk/commit/3f4a1f2e07e9fd4c3009e2dcc6d0a8309f53fe23))
* chore(deps): bump github.com/fatih/color from 1.12.0 to 1.13.0 (#556) (dependabot[bot])([bf79562](https://github.com/lacework/go-sdk/commit/bf7956236a1e9dcdc1a7cdb891da57918d510471))
* chore(deps): bump github.com/honeycombio/libhoney-go (#557) (dependabot[bot])([39fa065](https://github.com/lacework/go-sdk/commit/39fa065c3f46449c2033311b89eebe036f994769))
* chore(deps): bump github.com/spf13/viper from 1.8.1 to 1.9.0 (#553) (dependabot[bot])([bf4a39a](https://github.com/lacework/go-sdk/commit/bf4a39a802daa8e017d141f7a101afb03427ebb1))
* chore(deps): bump github.com/mitchellh/mapstructure from 1.4.1 to 1.4.2 (#552) (dependabot[bot])([649c729](https://github.com/lacework/go-sdk/commit/649c729f37f5067cda2ecac55deac1a32c0ec6a6))
* ci: fix downgrading of go packages (#575) (vatasha)([1a816c9](https://github.com/lacework/go-sdk/commit/1a816c9f8b7d47831326b3344d94072ff8350c6c))
* ci: version bump to v0.16.1-dev (Lacework)([bb5caa0](https://github.com/lacework/go-sdk/commit/bb5caa057d2ae5515aebbaf2a6d9ec2247f1e4ea))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
