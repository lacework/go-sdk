# Release Notes
Another day, another release. These are the release notes for the version `v0.2.19`.

## Features
* feat(cli): New VictorOps alert channel (#318) (Darren)([dfcd34a](https://github.com/lacework/go-sdk/commit/dfcd34adae5b84d697bd5f76831f749765f885fa))
* feat(cli): New CiscoWebex alert channel (#316) (Darren)([8e0071a](https://github.com/lacework/go-sdk/commit/8e0071a7d8a879d59c4545c8528e01287b90041e))
* feat: New Microsoft Teams alert channel (#315) (Darren)([e414226](https://github.com/lacework/go-sdk/commit/e41422624b5274a5d2b571fca4399ee2c9df3af3))
## Refactor
* refactor(cli): exponential retries polling scans (Salim Afiune Maya)([2bb881d](https://github.com/lacework/go-sdk/commit/2bb881d90fa6497003a89045acb41978398d5502))
## Bug Fixes
* fix(cli): match API client timeout with NGINX (#321) (Salim Afiune)([10b7a28](https://github.com/lacework/go-sdk/commit/10b7a28677af003e176de26b0dc1558c4837d1c4))
## Other Changes
* ci: increase integration test timeout to 30m (Salim Afiune Maya)([3081e3c](https://github.com/lacework/go-sdk/commit/3081e3cc6264e85b7fd80019a36d89dda1d9e5c9))
* test: change ctr vuln scan tag (Salim Afiune Maya)([c15bd1b](https://github.com/lacework/go-sdk/commit/c15bd1b2e0454ff61f6a7549f62c6bc791190079))
* test(cli): disable failing tests (RAIN-15300) (#320) (Salim Afiune)([e2afb31](https://github.com/lacework/go-sdk/commit/e2afb31277e8a2781c5f8e55ebf352992cd6d2b0))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
