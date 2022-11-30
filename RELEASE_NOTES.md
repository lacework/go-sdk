# Release Notes
Another day, another release. These are the release notes for the version `v1.1.0`.

## Features
* feat(cli): find and fix GCP instances w/o agents (#1022) (Nick Schmeller)([08edf28a](https://github.com/lacework/go-sdk/commit/08edf28a69693926b5b01be54750ca61b671942a))
## Refactor
* refactor: remove all unused functions, types and tests (#1037) (Darren)([7170808b](https://github.com/lacework/go-sdk/commit/7170808bd3c7ced956d4a4bd836a647321cfea43))
## Bug Fixes
* fix(cli): set correct SSH user for Amazon Linux 2 (#1046) (Nick Schmeller)([a92b37d1](https://github.com/lacework/go-sdk/commit/a92b37d13c9a90d65817d7a8d1ab1e29e693be3f))
## Other Changes
* chore(deps): bump github.com/aws/aws-sdk-go-v2/service/ec2instanceconnect (#1006) (dependabot[bot])([8cd69ed1](https://github.com/lacework/go-sdk/commit/8cd69ed190c908eff374ee39c31cf09e4dcdbea2))
* chore(deps): bump github.com/aws/aws-sdk-go-v2/service/ec2 (#1043) (dependabot[bot])([357ce2db](https://github.com/lacework/go-sdk/commit/357ce2dbbf910da1c3bf8bf19feb31d15cc70599))
* ci: fix release script; sort versions in find_latest_version (#1042) (Darren)([330a6a16](https://github.com/lacework/go-sdk/commit/330a6a16ad7d78718932e02f58f3bb4160f1f46a))
* ci: add logging line to output latest version found (#1040) (Darren)([38ec6ca1](https://github.com/lacework/go-sdk/commit/38ec6ca1a73016c5bd4244bce30742bd3a9bf396))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
