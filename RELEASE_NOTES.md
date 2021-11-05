# Release Notes
Another day, another release. These are the release notes for the version `v0.20.0`.

## Features
* feat(cli): Manage Alert Rules in Lacework CLI (#597) (Darren)([4d5648c](https://github.com/lacework/go-sdk/commit/4d5648cce2430b670167d3d66c3d8aee860b5f39))
## Bug Fixes
* fix(cli): avoid panic when assessment has NO vulns (#600) (Salim Afiune)([4638c7e](https://github.com/lacework/go-sdk/commit/4638c7e7d6b08e8f9bb123ed897321897a33ab4e))
* fix: ResourceGroups lastUpdated time format (#599) (Darren)([506954f](https://github.com/lacework/go-sdk/commit/506954fd1451d192a43d2ae1a38fb1e2b6261b4e))
* fix(cli): making evaluatorId optional (hazedav)([a7da3c0](https://github.com/lacework/go-sdk/commit/a7da3c0bc346ac1db957a3854c04c2c5a6a537c6))
## Other Changes
* ci: version bump to v0.19.1-dev (Lacework)([e8e36dd](https://github.com/lacework/go-sdk/commit/e8e36ddf314463540288dea6d577a12b5d41456a))
* test(cli): adds host query integration test file (hazedav)([7be96a7](https://github.com/lacework/go-sdk/commit/7be96a79d2bb3f9df84c259615fe0027110342bf))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
