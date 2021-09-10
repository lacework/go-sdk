# Release Notes
Another day, another release. These are the release notes for the version `v0.15.0`.

## Features
* feat(lwdomain): export go package (#545) (Salim Afiune)([13fb167](https://github.com/lacework/go-sdk/commit/13fb1675fdc30446ef903404a191104cf71195ee))
* feat(cli): stdin support for query/policy (hazedav)([770f0ee](https://github.com/lacework/go-sdk/commit/770f0ee26c613a59a166673f9d4a372a586f7da1))
* feat(cli): Manage Resource Groups in the lacework cli (#538) (Darren)([5e27cc8](https://github.com/lacework/go-sdk/commit/5e27cc8f0fba66d62b9113be1af0f97957529063))
* feat: Resource Groups (#537) (Darren)([754e8f4](https://github.com/lacework/go-sdk/commit/754e8f47ba2fbf1a68ca8b71fe3a761a36f1eb6a))
* feat(api): Add Non-OS Package Support for APIv1 (#539) (Andre Elizondo)([92e20e6](https://github.com/lacework/go-sdk/commit/92e20e6fa48be3ce192b84cfec832d9261f892e2))
* feat(cli): switch over to API v2 by default (#519) (Salim Afiune)([d605b45](https://github.com/lacework/go-sdk/commit/d605b45ad95920990ce7b0e8e662b0c12369afff))
## Refactor
* refactor: Add resource group specific create functions(#542) (Darren)([8a78438](https://github.com/lacework/go-sdk/commit/8a7843885afdfcde4ac0d98381748b0f220dbfcc))
## Bug Fixes
* fix(cli): Output correct resource group json (#543) (Darren)([3311ef2](https://github.com/lacework/go-sdk/commit/3311ef262edc058f29f0f151d938c3e76ab15f05))
* fix(install.ps1): Force Set-Variable as Constant (#541) (Salim Afiune)([6087ac2](https://github.com/lacework/go-sdk/commit/6087ac2b977132ba13350b248e2551c15c75881c))
## Documentation Updates
* docs(cli): correct typo v1 to v2 (#536) (lhasadreams)([cff6064](https://github.com/lacework/go-sdk/commit/cff6064a25d711d784dd0623c41739af78f0ff92))
## Other Changes
* ci: version bump to v0.14.1-dev (Lacework)([9417890](https://github.com/lacework/go-sdk/commit/941789044a3cdba6565209570149cf15475ad5fe))
* test(cli): re-enable vuln ctr scan tests (#532) (Salim Afiune)([7697c23](https://github.com/lacework/go-sdk/commit/7697c23c76f581f1610e16c83f0c787bac74695d))
* test(cli): disable vuln ctr scan tests (#530) (Salim Afiune)([0787b1b](https://github.com/lacework/go-sdk/commit/0787b1b6f7f74a79a528e0dde5157b172732344d))
* test(cli): re-enable query/policy integration (#526) (hazedav)([d54b51a](https://github.com/lacework/go-sdk/commit/d54b51a46a2154d8912859c64eb761d6bb72256f))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
