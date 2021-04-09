# Release Notes
Another day, another release. These are the release notes for the version `v0.5.0`.

## Features
* feat: Filter vuln host list cves and show assessment by severity (#375) (Darren)([7e9313e](https://github.com/lacework/go-sdk/commit/7e9313e82f6a9af30560aac441f67bd6bee0a619))
* feat(cli): split list of Google project id and alias (#377) (Salim Afiune)([3f8dd94](https://github.com/lacework/go-sdk/commit/3f8dd9419350d60bd3bdd5b416a2e5d7e23af89a))
## Refactor
* refactor(cli): use official url for agent installs (#369) (Salim Afiune)([ca2d3f7](https://github.com/lacework/go-sdk/commit/ca2d3f7063dce713f6743ace9a4a07d5cc6c2119))
## Bug Fixes
* fix(cli): split GCP projects with parentheses (#380) (Salim Afiune)([339d13e](https://github.com/lacework/go-sdk/commit/339d13ef3af47cc4d52e6833ae308e0d96e5349f))
* fix(cli): match err message changed in server resp (#376) (Salim Afiune)([de7b137](https://github.com/lacework/go-sdk/commit/de7b13796ce7a5edfb173f2a868d281dce35ddc1))
## Other Changes
* chore: run goimports-reviser to order imports (#378) (Salim Afiune)([44b8e58](https://github.com/lacework/go-sdk/commit/44b8e589e3861b849e6092db2a248f98882f5918))
* chore(deps): bump github.com/briandowns/spinner from 1.10.0 to 1.12.0 (#374) (dependabot[bot])([e6058eb](https://github.com/lacework/go-sdk/commit/e6058ebd1ec3e82c6c7c8763ac6b1a9ef9a552db))
* chore(deps): bump github.com/mitchellh/mapstructure from 1.3.0 to 1.4.1 (#373) (dependabot[bot])([2f40036](https://github.com/lacework/go-sdk/commit/2f40036fb827181ac9e7fc37d71a66ff916bc775))
* chore(deps): bump github.com/kyokomi/emoji/v2 from 2.2.5 to 2.2.8 (#372) (dependabot[bot])([ec6ec28](https://github.com/lacework/go-sdk/commit/ec6ec281c41fd047521c8c6715b2b9a9e78e0a56))
* chore(deps): bump github.com/olekukonko/tablewriter from 0.0.4 to 0.0.5 (#370) (dependabot[bot])([bdfa365](https://github.com/lacework/go-sdk/commit/bdfa365f936f775704f7c74320e162f000066ebc))
* chore(deps): bump github.com/honeycombio/libhoney-go (#371) (dependabot[bot])([0494850](https://github.com/lacework/go-sdk/commit/0494850865c95d1f07348c17a069ea116e6fe121))
* test(integration): fix failing test from API changes (#368) (Salim Afiune)([57b33e3](https://github.com/lacework/go-sdk/commit/57b33e360cb207b431c077937b99255f7bef4352))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
