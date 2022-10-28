# Release Notes
Another day, another release. These are the release notes for the version `v0.44.0`.

## Features
* feat: new v2 reports service (#943) (Darren)([fa766162](https://github.com/lacework/go-sdk/commit/fa766162f2562fc52f30d4424f8a6f6e4f1773a3))
## Refactor
* refactor(cli): migrate cmd compliance aws|gcp|azure get-report (#959) (Darren)([b613dfec](https://github.com/lacework/go-sdk/commit/b613dfec83be6cff23703dd3f54dd32c2444b421))
## Bug Fixes
* fix(cli): add query parameters to paging path (#980) (hazedav)([d4d1641a](https://github.com/lacework/go-sdk/commit/d4d1641ac483cd7c8b13172566f7260d360494a5))
* fix: fix gcp/azure --type flag  (#974) (Darren)([d67ea021](https://github.com/lacework/go-sdk/commit/d67ea0213f6ae315a36d2978a7ea5adb3c208613))
## Other Changes
* style(cli): rename `aws-install` files (#985) (Nick Schmeller)([c221cdab](https://github.com/lacework/go-sdk/commit/c221cdaba7cf6feb51b32ee30556f613badafa40))
* chore: Add .go-version file for goenv (#973) (Ross)([64ffaaea](https://github.com/lacework/go-sdk/commit/64ffaaeaa252e328128c4130695a2e1bcc3cb2ff))
* chore: bump golangci-lint to 1.50.0 (#969) (Ross)([524d7232](https://github.com/lacework/go-sdk/commit/524d7232ee7fa027a915388315e04cf6aa4748fd))
* chore: bump go jUnit to v2.0.0 (#972) (Ross)([b2cce81f](https://github.com/lacework/go-sdk/commit/b2cce81f3ccc34dba775f85918f24efaed09b12f))
* chore: bump goimports to 0.1.12 (#970) (Ross)([b297b0b1](https://github.com/lacework/go-sdk/commit/b297b0b181ed1b1c568966d3ecc7d5f71dde4410))
* chore: bump gotestsum 1.8.2 (#971) (Ross)([bed04947](https://github.com/lacework/go-sdk/commit/bed0494759a43ee1a37be1ade6e69181d22f4ca5))
* build: fix minor typo in pull request template (#968) (hazedav)([caeb5494](https://github.com/lacework/go-sdk/commit/caeb54942a4874af3f7c5c3fa4f51b7dd8a22f90))
* ci: version bump to v0.43.1-dev (Lacework)([5e9c0066](https://github.com/lacework/go-sdk/commit/5e9c0066010dc1f327fb8b69f1f5359a92f1b877))
* test(cli): enable compliance tests (#965) (Salim Afiune)([ee5389b6](https://github.com/lacework/go-sdk/commit/ee5389b6d757315e83bc701482f120df50ce4a18))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
