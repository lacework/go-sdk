# Release Notes
Another day, another release. These are the release notes for the version `v0.1.22`.

## Features
* feat(cli): Create PagerDuty Alert Channels 🚨 (#174) (Salim Afiune)([5cc424e](https://github.com/lacework/go-sdk/commit/5cc424e21598482f817288037c8f8e54397c13bd))
* feat(api): add PagerDuty alert channel integrations (#173) (Salim Afiune)([f46316c](https://github.com/lacework/go-sdk/commit/f46316c7f4150ccf99646640a12d801cb407134b))
* feat(cli): Create AWS CloudWatch Alert Channels 🚨 (Salim Afiune Maya)([201b59b](https://github.com/lacework/go-sdk/commit/201b59be0a97d661916ff401da0be903fee06f2f))
* feat(api): add AWS CloudWatch Alert Channels Int (Salim Afiune Maya)([d9a11ec](https://github.com/lacework/go-sdk/commit/d9a11ec5c242b09e19338c6b8a5a39ddf6ad368d))
* feat(api): enum AlertLevel for alert severity levels (Salim Afiune Maya)([d3bf436](https://github.com/lacework/go-sdk/commit/d3bf436933a794b6bbcc733da724159a9dc79a95))
* feat(api): get/update container registry integrations (#168) (Salim Afiune)([a072c46](https://github.com/lacework/go-sdk/commit/a072c46aff03e619fbef03488ba5b65730264b91))
## Refactor
* refactor(api): AlertChannel prefix in funcs/structs (Salim Afiune Maya)([b0429ef](https://github.com/lacework/go-sdk/commit/b0429efd0efa56ec9ccbe338a37a6e6ae2dc3bc5))
* refactor(api): use AlertLevel enum for Slack Alerts (Salim Afiune Maya)([4b5acf9](https://github.com/lacework/go-sdk/commit/4b5acf989fda4c052c3dc6b0206db866aa57f243))
## Bug Fixes
* fix(cli): missing fields for Slack integrations (#170) (Salim Afiune)([a8ce9a9](https://github.com/lacework/go-sdk/commit/a8ce9a90f52dd81281fca78b077435229bdbafaf))
## Other Changes
* chore(api): adds alert channel \_examples/ (Salim Afiune Maya)([f967206](https://github.com/lacework/go-sdk/commit/f967206db3dd209f94e694b5f4db98dd8b11f113))

## Docker Images
* `docker pull techallylw/lacework-cli:latest`
* `docker pull techallylw/lacework-cli:scratch`
* `docker pull techallylw/lacework-cli:ubi-8`
* `docker pull techallylw/lacework-cli:centos-8`
* `docker pull techallylw/lacework-cli:debian-10`
* `docker pull techallylw/lacework-cli:ubuntu-1804`
* `docker pull techallylw/lacework-cli:amazonlinux-2`
