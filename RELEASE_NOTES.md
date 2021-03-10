# Release Notes
Another day, another release. These are the release notes for the version `v0.2.22`.

## Bug Fixes
* fix: implement both ECR auth methods (Salim Afiune Maya)([7af312c](https://github.com/lacework/go-sdk/commit/7af312c0308d3dad52f4e0264ee354800ed4d616))
* fix(api): type-o in host vulnerability status (#336) (Scott Ford)([85b271b](https://github.com/lacework/go-sdk/commit/85b271bdb35c33eec12df5591c530800688b6d90))
## Other Changes
* ci: run nightly integration tests on ARM-linux (#317) (Salim Afiune)([525b51d](https://github.com/lacework/go-sdk/commit/525b51d6866550a26ae7faefb2d8a8de99dd9323))
* test(cli): use a cli tag with vuln (Salim Afiune Maya)([eae52b8](https://github.com/lacework/go-sdk/commit/eae52b8d8350fd1e6166c036fa8b86e3462517de))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
