# Release Notes
Another day, another release. These are the release notes for the version `v0.2.13`.

## Features
* feat(cli): support manifest bigger than 1k packages (Salim Afiune Maya)([eebddb9](https://github.com/lacework/go-sdk/commit/eebddb9325ede76ffa1853d00508da54cb5b9678))
* feat(cli): gen-pkg-manifest detect running kernel (Salim Afiune Maya)([9151be1](https://github.com/lacework/go-sdk/commit/9151be15a05b48f3d7456571cd75411f2ba7ddb9))
## Refactor
* refactor: simplify removeEpochFromPkgVersion func (Salim Afiune)([04aba5b](https://github.com/lacework/go-sdk/commit/04aba5bda340283f86d93496f01e0089a500468d))
## Bug Fixes
* fix(cli): ensure api client has valid auth token (Salim Afiune Maya)([056eda5](https://github.com/lacework/go-sdk/commit/056eda5cb7bde11e2334b6f38bd338afe111ade9))
## Other Changes
* ci: generate code coverage in HTML format (Salim Afiune Maya)([a58b58a](https://github.com/lacework/go-sdk/commit/a58b58a6477ec8d12c06bff3672093aef826c1f1))
* ci: add 'metric' as a valid commit message (Salim Afiune Maya)([dd7b601](https://github.com/lacework/go-sdk/commit/dd7b6010969d1f99055b7dbc9442498fa9f002cf))
* ci: fix slack notifications team alias ⭐ (Salim Afiune Maya)([ca51f92](https://github.com/lacework/go-sdk/commit/ca51f92693a48f113dd7661d9ef03eef7c26a17a))
* metric(cli): detect feature split_pkg_manifest (Salim Afiune Maya)([fdb9f4a](https://github.com/lacework/go-sdk/commit/fdb9f4a1c1eae2b9a44ea846fae413a93f073ca9))
* metric(cli): detect feature gen_pkg_manifest (Salim Afiune Maya)([78905bb](https://github.com/lacework/go-sdk/commit/78905bb73f398bf26a6e297e3929e5993e4965dc))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
