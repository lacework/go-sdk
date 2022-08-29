# Release Notes
Another day, another release. These are the release notes for the version `v0.41.0`.

## Features
* feat(cli): aws compliance search command (#893) (Darren)([5c665cb](https://github.com/lacework/go-sdk/commit/5c665cba2bc0d7f89c669d6096d3f60932eb8d2d))
* feat(cli): Addition of extra AWS terraform options (#890) (djmctavish)([de34097](https://github.com/lacework/go-sdk/commit/de34097c80f239e51335a33655ef89f35f5b7c06))
* feat(api): add v2 AwsCfg (#889) (Darren)([53bd1ca](https://github.com/lacework/go-sdk/commit/53bd1ca21bd5d6e4893397ffad14403154f184d9))
* feat: DataExportRule service (#887) (Darren)([9c6ca08](https://github.com/lacework/go-sdk/commit/9c6ca082f9ae0131bb2c3b0d65aec9fb51b4e6b6))
* feat: add credentials to sidekick data (#886) (Darren)([3311753](https://github.com/lacework/go-sdk/commit/331175320a594ef4a4646cf2ebdd5b8c06eee407))
## Refactor
* refactor(cli): migrate fetch accountIds in policy exception to API v2 (#896) (Darren)([510a41f](https://github.com/lacework/go-sdk/commit/510a41f0bdb6ba9e8d3df24db55c691ce0349fa0))
## Bug Fixes
* fix: omitempty profileVersions field (#897) (Darren)([2170edb](https://github.com/lacework/go-sdk/commit/2170edb47e5353597149c6208fedbe7604fd7a4f))
* fix(cli): account flag domain fix (#895) (jonathan stewart)([ed892dc](https://github.com/lacework/go-sdk/commit/ed892dc23bbef80eda97985ad1fc9c385ace4260))
* fix: add description field to dataexportrules (#891) (Darren)([0531907](https://github.com/lacework/go-sdk/commit/05319072e362cfc258b9555c52953138a62d9953))
* fix(lwupdater): avoid comparing dev versions (#881) (Salim Afiune)([a8f9ee4](https://github.com/lacework/go-sdk/commit/a8f9ee444c9899768b57fe510308878e2d6aa13c))
* fix: alert profile filter (#883) (Darren)([91e979a](https://github.com/lacework/go-sdk/commit/91e979a12156c5fced9fa9cf7490d6499ff87bb9))
## Documentation Updates
* docs(cli): display IaC default output location (#898) (sthetland)([139da68](https://github.com/lacework/go-sdk/commit/139da68fb6ec1033b1ecefe8a8295e08e955404e))
## Other Changes
* build(cli): sign Lacework CLI artifacts (#899) (Salim Afiune)([dcc4605](https://github.com/lacework/go-sdk/commit/dcc460567947f316d7eeb7dcda2406813cf83dc3))
* ci: fix policy test (#888) (Salim Afiune)([294f99b](https://github.com/lacework/go-sdk/commit/294f99b9e9767aea33450b561a46e78231f1bf4a))
* ci: version bump to v0.40.1-dev (Lacework)([599d7ff](https://github.com/lacework/go-sdk/commit/599d7ff61c0a21c961862d633d11f9c450cfc471))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
