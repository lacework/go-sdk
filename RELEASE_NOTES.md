# Release Notes
Another day, another release. These are the release notes for the version `v0.46.0`.

## Features
* feat(cdk): improve honeyvent data for components (#1016) (Salim Afiune)([8204f095](https://github.com/lacework/go-sdk/commit/8204f095f6527143d115fc424b4b542f8f12c8c7))
* feat(cli)!: pass agent token as flag for `aws-install` (#1010) (Nick Schmeller)([6dbcb233](https://github.com/lacework/go-sdk/commit/6dbcb233fc4167ed9470fd886730e4728eb4e4b5))
## Refactor
* refactor(cli): fixup + test `aws-install` SSH user (#1009) (Nick Schmeller)([d2167710](https://github.com/lacework/go-sdk/commit/d21677104c2e3769418c3089908ed9e863e4b4f9))
## Bug Fixes
* fix(cdk): improve component arg parsing (#1014) (Sam Shen)([daa9c271](https://github.com/lacework/go-sdk/commit/daa9c27104c47d5bda34c88d6352eec15936d223))
* fix(cli): rootCmd to intelligently PreRun components (#1015) (hazedav)([dbfd4063](https://github.com/lacework/go-sdk/commit/dbfd40635e72fb771823f9fad79060337871e87f))
* fix(cli): minor fix related to refreshing tokens (#1013) (hazedav)([9f4c6193](https://github.com/lacework/go-sdk/commit/9f4c6193fd726d116643e89bad7cf7ef5a5ffd95))
## Other Changes
* ci: version bump to v0.45.2-dev (Lacework)([605a8573](https://github.com/lacework/go-sdk/commit/605a8573efb3610a8590daf252ac8a3ae8155163))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
