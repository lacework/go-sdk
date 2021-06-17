# Release Notes
Another day, another release. These are the release notes for the version `v0.9.0`.

## Features
* feat(cli): new account list command (#444) (Salim Afiune)([735d34e](https://github.com/lacework/go-sdk/commit/735d34e346c82485b3c830f9dabcb0fa8aa9935a))
* feat(cli): LQLv2 support (#441) (hazedav)([8ad4abf](https://github.com/lacework/go-sdk/commit/8ad4abffdaadcf9cd8a3da0ffac8e012f4736672))
## Bug Fixes
* fix(cli): migration of Frankfurt organizational accounts (#447) (Salim Afiune)([73ed79a](https://github.com/lacework/go-sdk/commit/73ed79ad08280c47d9e54b179971a35dd79bfa19))
* fix(cli): adding more policy keys for json output (#442) (hazedav)([7a0d77f](https://github.com/lacework/go-sdk/commit/7a0d77fd028abad45cfce17d5c68f043f07c83a0))
* fix(cli): better error messaging for policy create (#443) (hazedav)([06174ff](https://github.com/lacework/go-sdk/commit/06174ff3de64f1ca25f7c098ec34a052dfbddbea))
## Other Changes
* chore(deps): bump github.com/fatih/color from 1.10.0 to 1.12.0 (#429) (dependabot[bot])([1b6bfe6](https://github.com/lacework/go-sdk/commit/1b6bfe692c3887725b9fc9104b6c216e7310397b))
* test(cli): adding lqlv2 test resource (#439) (hazedav)([df22443](https://github.com/lacework/go-sdk/commit/df22443556c1c709450d39c0b35f7a64228677ad))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
