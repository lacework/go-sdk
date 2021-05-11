# Release Notes
Another day, another release. These are the release notes for the version `v0.7.0`.

## Features
* feat: new lwconfig Go package to handle Lacework config (#402) (Salim Afiune)([18d4864](https://github.com/lacework/go-sdk/commit/18d4864b0f5fd904ebc31e5f66ded697f068c209))
* feat(cli): lacework query natural/relative times (hazedav)([624235d](https://github.com/lacework/go-sdk/commit/624235dffb40ca6bfd853a736ed71f155441b5ec))
* feat(cli): lacework query crumbs (hazedav)([333bec3](https://github.com/lacework/go-sdk/commit/333bec3dcf9295a30967cb06d7daf8a1560154dc))
## Refactor
* refactor: negligible -> info (#393) (Darren)([be9d203](https://github.com/lacework/go-sdk/commit/be9d2034cfdcbd207ff220ce11f223dac54bd811))
## Bug Fixes
* fix(cli): Fix filters with json output for vuln host command (#388) (Darren)([282404f](https://github.com/lacework/go-sdk/commit/282404f01dbace5c88d30fc5268addf839c5b40e))
* fix(cli): Json output for vuln ctr adheres to filter flags (#383) (Darren)([cf7bb6d](https://github.com/lacework/go-sdk/commit/cf7bb6d29652fe75cbb7564793d76524e1ecd320))
## Other Changes
* style: Add gotestsum to Makefile (#397) (Ross)([6b98f91](https://github.com/lacework/go-sdk/commit/6b98f9136f545fd576452ac47b507406b98c2864))
* chore(deps): bump github.com/AlecAivazis/survey/v2 from 2.0.7 to 2.2.12 (#395) (dependabot[bot])([3b3fd4e](https://github.com/lacework/go-sdk/commit/3b3fd4e4b37753666e448aa5ba0ce2da824d122d))
* ci: fix release script for automatic releases (#405) (Salim Afiune)([3b657ee](https://github.com/lacework/go-sdk/commit/3b657ee1f61aa2148c5a4c88c040c6cfa0552329))
* ci: Add Codefresh badge (#404) (Darren)([f94d4ab](https://github.com/lacework/go-sdk/commit/f94d4abff5e17674a48777f958036d00612c0c35))
* ci: Remove Circleci config (#398) (Darren)([6b95e84](https://github.com/lacework/go-sdk/commit/6b95e849efb78def42a156eb026202049171ca46))
* ci: update ci tools (#403) (Salim Afiune)([580947b](https://github.com/lacework/go-sdk/commit/580947b066c28217733954fc00cadbc543b482bc))
* ci: disable on-deman ctr vulnerability tests (#400) (Salim Afiune)([bcaaaa7](https://github.com/lacework/go-sdk/commit/bcaaaa7116145e0c1eb8898fec512aa988442257))
* metric: normalize account inside Honeyvent (#407) (Salim Afiune)([3eb3b56](https://github.com/lacework/go-sdk/commit/3eb3b565a211905524f17ca81a91f327237c1d97))
* test(policy): adding policy test resource (#391) (hazedav)([74d8cfe](https://github.com/lacework/go-sdk/commit/74d8cfe6510c02332f223b1ca3bf447905561761))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
