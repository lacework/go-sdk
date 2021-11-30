# Release Notes
Another day, another release. These are the release notes for the version `v0.22.0`.

## Features
* feat(cli): cache host vuln assessments for 1 hour (Salim Afiune Maya)([7f50bd2](https://github.com/lacework/go-sdk/commit/7f50bd23654a5566c35ba13e88ad2f4c9c8dafa9))
* feat(cli): logic to store assets with expiration time (Salim Afiune Maya)([ad920ca](https://github.com/lacework/go-sdk/commit/ad920ca5d638adb4f6331a80550381d8d1b3021d))
* feat(lwgenerate): Add initial hcl generation support (#611) (Matt Cadorette)([f9e9e78](https://github.com/lacework/go-sdk/commit/f9e9e781b73c9a92415702e052d10fed19a72ac2))
## Bug Fixes
* fix(cli): PKG typo in host vuln table (Salim Afiune Maya)([21ff81a](https://github.com/lacework/go-sdk/commit/21ff81a762d8c1dd3335bf88f8598ecebacc5686))
* fix(cli): avoid using "Fixed" host vulnerabilities (Salim Afiune Maya)([49ce7a7](https://github.com/lacework/go-sdk/commit/49ce7a70bca8eb3125c53a22e99c219fa61d3d33))
## Other Changes
* ci: version bump to v0.21.1-dev (Lacework)([30637bc](https://github.com/lacework/go-sdk/commit/30637bcce5b7aefb164c1ff3587753dbd4b5844d))
* test(cli): verify host vuln report rendering (Salim Afiune Maya)([b7d72e2](https://github.com/lacework/go-sdk/commit/b7d72e29b14b0987d021ecb03f618bbf3ac779a0))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
