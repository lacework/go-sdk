# Release Notes
Another day, another release. These are the release notes for the version `v0.2.16`.

## Features
* feat(cli): New Service Now alert channel (#303) (Darren)([512f2d9](https://github.com/lacework/go-sdk/commit/512f2d9c38d8124904dbdd661d9ab3b8441fc86d))
## Refactor
* refactor: Change input method for private_key field (#305) (Darren)([e56cdc6](https://github.com/lacework/go-sdk/commit/e56cdc68d74074f5ee904712aa56779dc9a0e1ed))
## Bug Fixes
* fix: Use select for issue grouping (#304) (Darren)([799d9c3](https://github.com/lacework/go-sdk/commit/799d9c34723d4af03f9ed811028880181f1757f5))
* fix: Add issue_grouping field to gcp pub sub (#301) (Darren)([1a66d2c](https://github.com/lacework/go-sdk/commit/1a66d2cc20a33674348705fbd3e552bf9222c787))
* fix(cli): install.sh should try curl and wget (Salim Afiune Maya)([f6b0bd7](https://github.com/lacework/go-sdk/commit/f6b0bd725992961e8c8b02ba3090164f31056388))
* fix(cli): install.sh should respect target override (Salim Afiune Maya)([4164f58](https://github.com/lacework/go-sdk/commit/4164f5872dd381471ad8608abd9c2fa821dac524))
## Other Changes
* chore(cli): install.sh print exitcodes for debugging (Salim Afiune Maya)([5e66c11](https://github.com/lacework/go-sdk/commit/5e66c11264f0b9ba012be20370b0b65cc114fe4c))
* ci: trigger homebrew update script (#299) (Darren)([9247cb1](https://github.com/lacework/go-sdk/commit/9247cb12b62582596939be55861f4d6c08bf8690))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
