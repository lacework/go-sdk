# Release Notes
Another day, another release. These are the release notes for the version `v0.35.0`.

## Features
* feat: show query start time in progress spinner (#798) (Darren)([2179616](https://github.com/lacework/go-sdk/commit/21796160c18f1c7acbf6ece7f2fbc8c3d5526ec0))
* feat(api): Add S3BucketArn field to allow BYOB for AWS EKS Audit Log (#793) (Ross)([93dbb2d](https://github.com/lacework/go-sdk/commit/93dbb2d74f700fc7c536566c7130d8096ea3febc))
* feat(cli): alert profiles command (#792) (Darren)([c0bc9a3](https://github.com/lacework/go-sdk/commit/c0bc9a3a694a7ae28f65c10f264dfff6e0042102))
## Refactor
* refactor(cli): query run progress message (#804) (hazedav)([f64d461](https://github.com/lacework/go-sdk/commit/f64d461b0dc674473122ff29035323b9d5b29dbe))
* refactor(cli): updated terminology for fail_on (#803) (hazedav)([f5b85f5](https://github.com/lacework/go-sdk/commit/f5b85f587f19ad656aef1095096d613cd5a328d8))
## Bug Fixes
* fix: allow resource scope to be optional (#807) (Salim Afiune)([6e8e55c](https://github.com/lacework/go-sdk/commit/6e8e55cdd0d3e54b3a8b9d450b36f84e9ceff614))
* fix(cli): Add MachineDetail tags (hazedav)([8cd1ea0](https://github.com/lacework/go-sdk/commit/8cd1ea0aebdee72d87d7e8faf2a7bae2cd87fe4c))
* fix(cli): panic while detecting ESM on Ubuntu sys (#799) (Salim Afiune)([dce2d18](https://github.com/lacework/go-sdk/commit/dce2d18ab0e88e1f106d69b352ac041ee7810221))
## Other Changes
* chore(deps): bump github.com/hashicorp/go-version from 1.4.0 to 1.5.0 (#801) (dependabot[bot])([7665edd](https://github.com/lacework/go-sdk/commit/7665edd61c4365d0fc12ce363ba4ff94421f41f8))
* chore(deps): bump github.com/hashicorp/terraform-exec (#773) (dependabot[bot])([a7915cc](https://github.com/lacework/go-sdk/commit/a7915cc74e2883fa048dff47e7b100a1e976358d))
* chore(deps): bump github.com/abiosoft/colima from 0.4.1 to 0.4.2 (#795) (dependabot[bot])([6ea8ef3](https://github.com/lacework/go-sdk/commit/6ea8ef30b61474daba5b44fadd0efd06b6a6ccfe))
* chore(deps): bump github.com/abiosoft/colima from 0.3.4 to 0.4.1 (#791) (dependabot[bot])([cf19e21](https://github.com/lacework/go-sdk/commit/cf19e21eebe9aa1d88dc179a221ee5e427b58e1b))
* ci: version bump to v0.34.1-dev (Lacework)([8811c84](https://github.com/lacework/go-sdk/commit/8811c84d050aad00f250b35b47ee8422e1907156))
* test: fix pipeline to check for low vulns (#805) (Salim Afiune)([233c24f](https://github.com/lacework/go-sdk/commit/233c24f372a913f61d64c509183c77323e376982))
* test: disable file test since CI runs as root (#800) (Salim Afiune)([5b271c6](https://github.com/lacework/go-sdk/commit/5b271c671d058ebe9fc898ec785c5c750e167d65))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
