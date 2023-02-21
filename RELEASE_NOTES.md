# Release Notes
Another day, another release. These are the release notes for the version `v1.11.0`.

## Features
* feat(cli): support bulk policy updates for policy enable/disable and … (#1132) (Darren)([19f1f3ca](https://github.com/lacework/go-sdk/commit/19f1f3ca3d4e8d989ac4fae25830874e23a72731))
## Bug Fixes
* fix: update VulnerabilityHost MachineTags to support GCE host types (#1151) (Darren)([f148dd48](https://github.com/lacework/go-sdk/commit/f148dd48eb11fe41a865a9708bd3932f6b073adb))
* fix: parsing of account names containing `lacework.net` (#1152) (Edoardo Pirovano)([94da6505](https://github.com/lacework/go-sdk/commit/94da65056125046da3f5a0cb44fc68bb5ebddd78))
* fix(ec2ssm): safely dereference SSM command output (#1147) (Nick Schmeller)([0ffd665a](https://github.com/lacework/go-sdk/commit/0ffd665a1b418ba7adae81ff87278203a5ea1f67))
* fix(cli): gRPC port collision (#1146) (Kolbeinn)([a9d0ef5d](https://github.com/lacework/go-sdk/commit/a9d0ef5d7c221a4827111659fc92f2830e13b703))
## Other Changes
* chore: fix go.mod (#1156) (Darren)([59e488e9](https://github.com/lacework/go-sdk/commit/59e488e9a2f115d523e4c5aa567d6ddd6fed94bc))
* chore: add growth-team to CODEOWNERS file (#1153) (Darren)([0d587535](https://github.com/lacework/go-sdk/commit/0d58753567a5c559a1067909a4a8b5dabb3cd0e3))
* ci: version bump to v1.10.1-dev (Lacework)([95ec0a8b](https://github.com/lacework/go-sdk/commit/95ec0a8bd1f0e20fa8a8334ed2b76d2287efa260))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
