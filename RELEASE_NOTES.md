# Release Notes
Another day, another release. These are the release notes for the version `v0.34.0`.

## Features
* feat: Add support for APIv2 CloudAccounts List (#784) (Ross)([563d1da](https://github.com/lacework/go-sdk/commit/563d1da63b7eae5f7e7691117787570d2d0e8cfc))
* feat(cli): add gcp compliance report disable warning (#787) (Ross)([0cc3734](https://github.com/lacework/go-sdk/commit/0cc373442eff0bc970b10d5eabfde0296ba699a1))
* feat(cli): add azure compliance report disable warning (#786) (Ross)([a028ecb](https://github.com/lacework/go-sdk/commit/a028ecb3b6a5a5e61f82464c9ef83226dd393227))
## Bug Fixes
* fix(test): Fix TestCloudAccountsListByType unit test failure (#788) (Ross)([48fb4c7](https://github.com/lacework/go-sdk/commit/48fb4c78d6f2d47c353341c0d47978fcefa25651))
* fix(api): add newResponse for function SearchAllPages (#770) (aircraft-cerier)([474a163](https://github.com/lacework/go-sdk/commit/474a16396b30d70b934a07839891436c3cef6643))
* fix(cli): Remove Additional Cloudtrail/AuditLog options when Cloudtrail/AuditLog Integration is set to No (#782) (Ross)([c9d7f48](https://github.com/lacework/go-sdk/commit/c9d7f48572e028be691e7b984b52208681cfad5a))
## Other Changes
* ci: version bump to v0.33.1-dev (Lacework)([3b3aacf](https://github.com/lacework/go-sdk/commit/3b3aacf9e8075d39532980b063e2083c9e0dc88a))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
