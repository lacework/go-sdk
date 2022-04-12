# Release Notes
Another day, another release. These are the release notes for the version `v0.31.0`.

## Features
* feat(cli): disable/enable policies by tag (#758) (Darren)([3a45a01](https://github.com/lacework/go-sdk/commit/3a45a0186ca908920f2383c2eec2c7c7db6e17bf))
* feat(cli): default template when creating LQL queries (#761) (Salim Afiune)([4bb9c19](https://github.com/lacework/go-sdk/commit/4bb9c1930b64e3e080ee44fc71028276bb44d8d8))
* feat(cli): disable/enable all <report-type> compliance command (#755) (Darren)([e694304](https://github.com/lacework/go-sdk/commit/e69430484c93245919b2e48e264dc3c24d6a7297))
* feat(cli): adds lql datasource relationships (#759) (hazedav)([4b1fc27](https://github.com/lacework/go-sdk/commit/4b1fc271db940d5272cd652f5b65d06e061af874))
* feat(api): Recommendations v1 Service (#751) (Darren)([253cc38](https://github.com/lacework/go-sdk/commit/253cc38420ec9a6423aec2dc57958255eb1e5cf9))
## Refactor
* refactor(cli): inefficient policy sort (hazedav)([9a35ff6](https://github.com/lacework/go-sdk/commit/9a35ff665acd795f7a713b470d6603771f060beb))
## Bug Fixes
* fix(cli): command arguments (#760) (Salim Afiune)([574dfc1](https://github.com/lacework/go-sdk/commit/574dfc1dbb8eb066ca3e8668ed72be58731b0ff3))
* fix: add missing gcp compliance report types to get-report command (#756) (Darren)([c05a07f](https://github.com/lacework/go-sdk/commit/c05a07ff0664164d8e4d53009c4e2af6e2451c27))
## Other Changes
* chore(deps): bump github.com/kyokomi/emoji/v2 from 2.2.8 to 2.2.9 (#724) (dependabot[bot])([717a725](https://github.com/lacework/go-sdk/commit/717a725edfdb205701456bd6f960b1b0a3788988))
* chore(deps): bump github.com/spf13/cobra from 1.3.0 to 1.4.0 (#733) (dependabot[bot])([d1531ee](https://github.com/lacework/go-sdk/commit/d1531eef51e37966bc39e002d0589d400c548990))
* ci: add make cmd for output go tests in junit format (#754) (Darren)([f3d0324](https://github.com/lacework/go-sdk/commit/f3d032431bea1770ab8b731a5aefef23ab1c4594))
* ci: version bump to v0.30.1-dev (Lacework)([5cfe47a](https://github.com/lacework/go-sdk/commit/5cfe47a62c3c650cb35e879528614aa8155f9ac8))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
