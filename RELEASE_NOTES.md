# Release Notes
Another day, another release. These are the release notes for the version `v0.2.2`.

## Features
* feat(lql): --file flag to load LQL query from disk (Salim Afiune Maya)([4804319](https://github.com/lacework/go-sdk/commit/4804319a0c26211119c10eb3dc4d889b3da7e227))
* feat(cli): --file to pass a package manifest file (Salim Afiune Maya)([75680d8](https://github.com/lacework/go-sdk/commit/75680d8d9469d8679b17c46979439340d8869da9))
* feat: human-readable output for scan-pkg-manifest (Salim Afiune Maya)([783f550](https://github.com/lacework/go-sdk/commit/783f55015c1e6a1071927e19810266376ecbe082))
* feat(lql): improve running queries (Salim Afiune Maya)([61c5ee5](https://github.com/lacework/go-sdk/commit/61c5ee51aac65626aff4f81ebceb96633865d2f7))
## Bug Fixes
* fix(ci): remove slack notification for windows (#214) (Salim Afiune)([a2c5124](https://github.com/lacework/go-sdk/commit/a2c51242c08c1683cfb9c80c832be2559058f957))
## Other Changes
* ci(slack): notify pipeline failures (#213) (Salim Afiune)([85ad396](https://github.com/lacework/go-sdk/commit/85ad396f6cb049ab246ff36fa2f29d46fab6459d))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
