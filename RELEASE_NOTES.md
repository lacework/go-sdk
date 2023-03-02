# Release Notes
Another day, another release. These are the release notes for the version `v1.12.0`.

## Features
* feat(cli): display query inside 'policy show' command (Salim Afiune Maya)([11533129](https://github.com/lacework/go-sdk/commit/115331290d7c22138f7550bff58e91431cc71fb2))
* feat(cli): Add 'subtype' filter flag to 'lacework report-definitions … (#1163) (Darren)([d2563764](https://github.com/lacework/go-sdk/commit/d2563764f8f104c4e89c49e5540ed899c355d50d))
* feat(cli): show valid constraint keys for policy exception create prompt (#1159) (Darren)([30b51441](https://github.com/lacework/go-sdk/commit/30b51441250ee84b20219b7181355e4efb962657))
## Refactor
* refactor(cli): 'query list-sources' output (Salim Afiune Maya)([6fe50b95](https://github.com/lacework/go-sdk/commit/6fe50b9557f1f0e22b3e5387061e95ba45a72758))
* refactor(cli): improve 'policy show' command output (Salim Afiune Maya)([0f157af9](https://github.com/lacework/go-sdk/commit/0f157af96887998e03ddbf08d45e9363f8223f01))
* refactor(cli): improve policy list command output (Salim Afiune Maya)([88829437](https://github.com/lacework/go-sdk/commit/888294371c6243e85980efb42789308bebc14b71))
## Bug Fixes
* fix: resolve inaccurate total packages count on 'vuln ctr show --deta… (#1166) (Darren)([d46a4acd](https://github.com/lacework/go-sdk/commit/d46a4acda022a01a00c73497ac98e8f43b04ddb3))
## Other Changes
* ci: disable mock terminal test on windows ci (Darren Murray)([1b5b419f](https://github.com/lacework/go-sdk/commit/1b5b419f34fbbad6a83ebcc889e3a2f58eef007e))
* ci: version bump to v1.11.2-dev (Lacework)([0da4c899](https://github.com/lacework/go-sdk/commit/0da4c899911c5e234620eb1a0d7082cb9d47db2a))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
