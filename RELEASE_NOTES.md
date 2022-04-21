# Release Notes
Another day, another release. These are the release notes for the version `v0.32.0`.

## Features
* feat(cli): output query & policy in YAML 📝 (#775) (Salim Afiune)([9173efe](https://github.com/lacework/go-sdk/commit/9173efef1603bb204ac998c7d8ed4a65a10bcdc8))
* feat(api): AlertProfiles v2 service (#765) (Darren)([683ceca](https://github.com/lacework/go-sdk/commit/683ceca0c99971c9cde7371e3fc5cbd24adb9d33))
* feat(cli): ubuntu ESM support (#768) (Salim Afiune)([7527e83](https://github.com/lacework/go-sdk/commit/7527e834cb15f246839fd984f73e4885949a9fea))
## Bug Fixes
* fix(cli): only display datasource relationships when avail (hazedav)([c9589ed](https://github.com/lacework/go-sdk/commit/c9589ed5abbac73511f829fb9b69571cc795ed60))
## Documentation Updates
* docs: enable policy enable/disable help text. (#766) (Darren)([768cfbf](https://github.com/lacework/go-sdk/commit/768cfbf4125f1f79b9801674dc7f9a4ca8fe0365))
## Other Changes
* chore(deps): bump github.com/BurntSushi/toml from 0.4.1 to 1.1.0 (#763) (dependabot[bot])([a3b6c5b](https://github.com/lacework/go-sdk/commit/a3b6c5b34bdd7272a2fc93267536865b37339330))
* chore(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1 (#740) (dependabot[bot])([db6541c](https://github.com/lacework/go-sdk/commit/db6541c0e32ae343629e01644971fdb2620d2404))
* chore(deps): bump github.com/abiosoft/colima from 0.3.2 to 0.3.4 (#723) (dependabot[bot])([b9bef2a](https://github.com/lacework/go-sdk/commit/b9bef2a4c57abfdc89ed5660e0f9250e755ef149))
* ci: version bump to v0.31.1-dev (Lacework)([7868fe7](https://github.com/lacework/go-sdk/commit/7868fe791cdf4a8cf2e574776798175f573a2939))
* test(lql): fix case issue with host tests (hazedav)([30d7d7b](https://github.com/lacework/go-sdk/commit/30d7d7b77fb477a81361b0eca4d18eaed637f596))
* test: rename LW_CLI_ -> CLI_ (#769) (Darren)([5cfedda](https://github.com/lacework/go-sdk/commit/5cfedda359d2cba7d17a25f1eecc55207989ab67))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
