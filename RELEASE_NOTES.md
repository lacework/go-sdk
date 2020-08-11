# Release Notes
Another day, another release. These are the release notes for the version `v0.1.23`.

## Refactor
* refactor(cli): replace '--pdf-file' for '--pdf' (#180) (Salim Afiune)([80bbce6](https://github.com/lacework/go-sdk/commit/80bbce636cac49fe315118add45252bd8ee4bf6a))
## Bug Fixes
* fix(cli): missing integration details (#181) (Salim Afiune)([40355d3](https://github.com/lacework/go-sdk/commit/40355d3877c2674268c38bb5cc81a698dd115166))
* fix(cli): error showing non-existing integration (#178) (Salim Afiune)([252072f](https://github.com/lacework/go-sdk/commit/252072faa60aaac06fb7bbf2dd7ca82fa71d2b09))
## Other Changes
* ci: build statically linked binaries (Salim Afiune Maya)([43f6f80](https://github.com/lacework/go-sdk/commit/43f6f804ffac3f8e326dc31f4196808f39bc035d))
* ci(integration): add windows support (Salim Afiune Maya)([46632e7](https://github.com/lacework/go-sdk/commit/46632e72e0ab9ee45d690605e4c52efb1a8cf391))

## Docker Images
* `docker pull techallylw/lacework-cli:latest`
* `docker pull techallylw/lacework-cli:scratch`
* `docker pull techallylw/lacework-cli:ubi-8`
* `docker pull techallylw/lacework-cli:centos-8`
* `docker pull techallylw/lacework-cli:debian-10`
* `docker pull techallylw/lacework-cli:ubuntu-1804`
* `docker pull techallylw/lacework-cli:amazonlinux-2`
