# Release Notes
Another day, another release. These are the release notes for the version `v0.37.0`.

## Features
* feat(api): add policy exceptions service (#826) (Darren)([fe5b9bf](https://github.com/lacework/go-sdk/commit/fe5b9bf508800d82377e379805958925d10e8a08))
* feat(cdk): pass env variable LW_COMPONENT_NAME (#829) (Salim Afiune)([4b7451a](https://github.com/lacework/go-sdk/commit/4b7451a3823daab4fa90a00bda0e6ea6666b5ebd))
## Refactor
* refactor: move vulnerabilityExceptionFixableEnabled to go-sdk (#820) (Darren)([941570e](https://github.com/lacework/go-sdk/commit/941570e9edc44277da081a5b6db9c5e36a3036d2))
## Bug Fixes
* fix(GcpGkeAudit): Casing issue with PrivateKeyId & PrivateKeyId (#834) (Ross)([b7b328b](https://github.com/lacework/go-sdk/commit/b7b328b4ca0557ec569460d37895875288ffc83a))
* fix: misspelled namespace field (#831) (Salim Afiune)([6312843](https://github.com/lacework/go-sdk/commit/6312843f7a1a2dfbe68af3eba552abf8338bc481))
* fix(cli): Make the html table scrollable (#825) (John John)([85d5214](https://github.com/lacework/go-sdk/commit/85d5214a6d09ded8ace5056145b4f45083e97789))
* fix(cdk): run component dynamically (#823) (Salim Afiune)([bb1ce0d](https://github.com/lacework/go-sdk/commit/bb1ce0dc0b9c74bff80385735028e71941e5aa66))
## Other Changes
* chore(deps): bump github.com/hashicorp/terraform-json (#808) (dependabot[bot])([ac2b094](https://github.com/lacework/go-sdk/commit/ac2b09446833a6cc1abb0bf74df547889f7107f1))
* chore(deps): bump github.com/spf13/viper from 1.10.1 to 1.12.0 (#802) (dependabot[bot])([000024a](https://github.com/lacework/go-sdk/commit/000024a90f37963dce2d5a80139f2ca5e8ab977e))
* chore(deps): bump github.com/stretchr/testify from 1.7.2 to 1.7.3 (#828) (dependabot[bot])([fb93f55](https://github.com/lacework/go-sdk/commit/fb93f558c6c9b07be8001d2f8335b1bddc46ea3a))
* ci: version bump to v0.36.1-dev (Lacework)([cdd6d7b](https://github.com/lacework/go-sdk/commit/cdd6d7b0b1596a50261449742c9e134916b9a535))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
