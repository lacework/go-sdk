# Release Notes
Another day, another release. These are the release notes for the version `v0.14.0`.

## Features
* feat(cli): More detailed CSV compliance reports (#522) (Matt Cadorette)([bd02a44](https://github.com/lacework/go-sdk/commit/bd02a44009e5cf00901ced32385add97e04be0aa))
* feat(api): AwsS3 Alert Channel (v2) (Salim Afiune Maya)([507cac1](https://github.com/lacework/go-sdk/commit/507cac14781b986662a2a78b7d436cc2ba1250ee))
* feat(api): SlackChannel Alert Channel (v2) (Salim Afiune Maya)([e6ac8cb](https://github.com/lacework/go-sdk/commit/e6ac8cb8eed31646020aeb3bbf5d2ef539258389))
## Refactor
* refactor: rename alert channel enums (Salim Afiune Maya)([f727b9b](https://github.com/lacework/go-sdk/commit/f727b9b78a79047b5df5665324d54d0cb4302b7b))
## Bug Fixes
* fix(cli): removing policyUi from policy schema (#523) (hazedav)([5586aec](https://github.com/lacework/go-sdk/commit/5586aec415034a431490c0b3cfb4d51385bfdd23))
* fix(cli): handle scan status 'Partial' (#525) (Darren)([a8ba4b7](https://github.com/lacework/go-sdk/commit/a8ba4b722cf211f527c685519504359525bd0ebe))
* fix(cli): query run via editor (hazedav)([ca3f2a3](https://github.com/lacework/go-sdk/commit/ca3f2a34bb76b97f42b6a5cafb664839a1fc3689))
* fix(cli): Detect if running in cloudshell to output correct update command (#521) (Darren)([477fadd](https://github.com/lacework/go-sdk/commit/477fadd9a96398198f397eb57c29231761173a9b))
## Other Changes
* chore(api): update v2 alert channel tests (Salim Afiune Maya)([3581d0a](https://github.com/lacework/go-sdk/commit/3581d0a2237d176c6e1d1d0dfa749b1a3bce8480))
* ci: version bump to v0.13.1-dev (Lacework)([39cc865](https://github.com/lacework/go-sdk/commit/39cc865741096ae4f7a29c09722bd38d7637752d))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
