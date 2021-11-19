# Release Notes
Another day, another release. These are the release notes for the version `v0.21.0`.

## Features
* feat: Auto docs generation for the Lacework CLI (#616) (Salim Afiune)([81b7fdb](https://github.com/lacework/go-sdk/commit/81b7fdbf50c1e560277c44be38e9fd2473ca15c7))
* feat: Filter vuln ctr list-assessments by registry (#609) (Darren)([264b450](https://github.com/lacework/go-sdk/commit/264b4504633267cd6f1929eb39109f6cafb74692))
* feat(api): Report Rules v2 Service (#607) (Darren)([26e96cc](https://github.com/lacework/go-sdk/commit/26e96cc5c95e5188b6d0b82977371f578eba55e9))
## Bug Fixes
* fix(cli): remove sub-account during reconfiguration (#602) (Salim Afiune)([d623afe](https://github.com/lacework/go-sdk/commit/d623afedc2152cbeb7abcd90f87e6d8a475e95a0))
## Documentation Updates
* docs: Run make generate-docs (Darren Murray)([b6438a1](https://github.com/lacework/go-sdk/commit/b6438a1595363faa331a1e79c89dab492e8c5e44))
* docs: Vuln ctr list-assessments registry flag (Darren Murray)([cf32b2f](https://github.com/lacework/go-sdk/commit/cf32b2f41b306fcd18103ec4fae992058de4d273))
## Other Changes
* chore: ran make go-vendor with go 1.17 (#615) (Salim Afiune)([432104f](https://github.com/lacework/go-sdk/commit/432104faf0522c58ee404232600714360aa64bea))
* chore: Update minimum go version -> 1.17 (#612) (Darren)([c6e183c](https://github.com/lacework/go-sdk/commit/c6e183c22d576b11698989a34443e4c4c53d235a))
* chore(deps): bump github.com/honeycombio/libhoney-go (#605) (dependabot[bot])([7226624](https://github.com/lacework/go-sdk/commit/7226624e163015d12cd5b58042ec552df4ad64ff))
* chore(cli): Update s3, slack, and cloudwatch alert channels cli commands to point to API v2 (#606) (vatasha)([65887be](https://github.com/lacework/go-sdk/commit/65887be28d97aa2d16b5f528ca2cc227ea7e94fd))
* ci: add Azure integration as a test resource (#608) (Salim Afiune)([5cab0ef](https://github.com/lacework/go-sdk/commit/5cab0ef1d1435c64fd3eb52c086ee63e756c41ad))
* ci: Update go conditinal build comment to 1.17 preferred format (#604) (Darren)([a8f90f6](https://github.com/lacework/go-sdk/commit/a8f90f668ac4c9ac7b6eaec7b3720de390d02a5f))
* ci: version bump to v0.20.1-dev (Lacework)([15761a7](https://github.com/lacework/go-sdk/commit/15761a7f5e6f1cc750f3c1b0911f2dabe7eb3c00))
* test: Fix failing alert rule test (#603) (Darren)([c1b8ca5](https://github.com/lacework/go-sdk/commit/c1b8ca59b088bda0715562f846ddb05bfee7bfe1))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
