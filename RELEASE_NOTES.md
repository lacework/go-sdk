# Release Notes
Another day, another release. These are the release notes for the version `v0.10.1`.

## Features
* feat(cli): apiv2 support for query validate (hazedav)([bcd81ce](https://github.com/lacework/go-sdk/commit/bcd81ce50c98b4b5204122301f4490a0a74c33c6))
## Bug Fixes
* fix(api): callbacks panic when response is nil (#486) (Salim Afiune)([8b7472a](https://github.com/lacework/go-sdk/commit/8b7472abe64cb3c837e34be7a2ec6703eacf0eb4))
* fix(cli): Output gcp comp list with --json flag (#482) (Darren)([9f95aa8](https://github.com/lacework/go-sdk/commit/9f95aa839dc16b6a64c0bbce9c66185b8fb3eef0))
## Other Changes
* test(cli): create VMs of supported platforms (#480) (Salim Afiune)([7f7c06e](https://github.com/lacework/go-sdk/commit/7f7c06e4013a0b379a685d21e78e03de3422ce56))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
