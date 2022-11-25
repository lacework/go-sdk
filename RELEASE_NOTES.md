# Release Notes

Another day, another release. These are the release notes for the version `v1.0.0`.

## Features

* feat: Add new cloud account type GCP Sidekick Integration (#1028) (
  ammarekbote)([fa428fb3](https://github.com/lacework/go-sdk/commit/fa428fb37add07373131930117c74def4b8129e8))
* feat(lwseverity): add valid severity use cases (#1024) (
  hazedav)([6b0b1330](https://github.com/lacework/go-sdk/commit/6b0b133005e387a32b6a705dd187cf233769853d))
* feat(cdk): bidirectional communication via gRPC (#1023) (Salim
  Afiune)([ff5654bd](https://github.com/lacework/go-sdk/commit/ff5654bdebdf92c0e95a4b9f8fde73a90cf8fae6))
* feat: lacework alert command (
  hazedav)([1059cc56](https://github.com/lacework/go-sdk/commit/1059cc566afbeccc3ca3df35d118b555385df0d6))

## Refactor

* refactor(v1): migrate cli to use v2 api (#1025) (
  Darren)([9f12eb22](https://github.com/lacework/go-sdk/commit/9f12eb226c47f03d5bafe53ece9a8558d77f9454))

## Bug Fixes

* fix: remove event command and api service (#1030) (
  hazedav)([87fd557d](https://github.com/lacework/go-sdk/commit/87fd557dfa158e141e005b9b44045a038d3b32e8))
* fix(api): container registry dockerhub 🐛  (#1032) (Salim
  Afiune)([0703eb93](https://github.com/lacework/go-sdk/commit/0703eb93b543684ecb631aa98813880427fb67ca))
* fix: compliance evaluation search cmd, searches greater than last 7 days  (#977) (
  Darren)([cf921309](https://github.com/lacework/go-sdk/commit/cf921309d6060621dc18231b43b585cc5e0f7bb4))

## Other Changes

* ci: version bump to v0.46.1-dev (
  Lacework)([a4a42bcb](https://github.com/lacework/go-sdk/commit/a4a42bcbc4240a3b63845fc8a69acc4d687f64ba))

## Docker Images

* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
