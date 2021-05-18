# Release Notes
Another day, another release. These are the release notes for the version `v0.7.1`.

## Bug Fixes
* fix(api): graceful degradation for HTTP 204 (#409) (hazedav)([44eb159](https://github.com/lacework/go-sdk/commit/44eb159af877e359be7a840cb72d4fd05b0cc718))
## Other Changes
* build(cli): Upgrade Go version 1.15 -> 1.16, Add support for darwin/arm64 (#408) (Darren)([c9d0a6f](https://github.com/lacework/go-sdk/commit/c9d0a6f8807eb12d3ebe86985e087450f095558e))
* ci: Disable failing test (#418) (Darren)([c5d7733](https://github.com/lacework/go-sdk/commit/c5d7733d589cacf2b44e58f1d55c7ddc6ef9d5c5))
* ci(cli): fix integration test (#417) (Salim Afiune)([b9ad070](https://github.com/lacework/go-sdk/commit/b9ad0706341d2583cd16cca70860695d55e9401b))
* ci: Update go deps (#413) (Darren)([8ee2393](https://github.com/lacework/go-sdk/commit/8ee239338741d5e624336347521514e86b9fa31e))
* test(cli): add container vuln scan NO polling (#414) (Salim Afiune)([793103f](https://github.com/lacework/go-sdk/commit/793103fb77170952e7e7f8498459f68d38c875c6))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
