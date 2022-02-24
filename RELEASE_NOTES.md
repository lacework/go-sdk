# Release Notes
Another day, another release. These are the release notes for the version `v0.26.0`.

## Features
* feat(cli): show resources affected by a compliance recommendation (#697) (Darren)([86f86ee](https://github.com/lacework/go-sdk/commit/86f86eef38840a28be310347f67261ba588a850c))
* feat(cli): improve user feedback running terraform commands (#706) (Matt Cadorette)([8526d2c](https://github.com/lacework/go-sdk/commit/8526d2ccfb8a48dbc8ac1315d84006aea058822a))
* feat(cli): add lacework query preview-source (hazedav)([40e8bbb](https://github.com/lacework/go-sdk/commit/40e8bbb03fa3fca6eecac48c24ffca1468821815))
* feat(cli): switch-profile command (#705) (Salim Afiune)([75237cd](https://github.com/lacework/go-sdk/commit/75237cdb15508c87cad3cc0694e304524a0da3f6))
## Other Changes
* chore(deps): bump github.com/spf13/cobra from 1.2.1 to 1.3.0 (#647) (dependabot[bot])([2cb15fb](https://github.com/lacework/go-sdk/commit/2cb15fb48186a45064aac0ecce340efcfb2e8de1))
* chore(deps): bump github.com/spf13/viper from 1.9.0 to 1.10.1 (#644) (dependabot[bot])([5dac81c](https://github.com/lacework/go-sdk/commit/5dac81cfb57b29dc8addcbb32b16f11abeeca836))
* chore(deps): bump github.com/honeycombio/libhoney-go (#652) (dependabot[bot])([6e5c7bd](https://github.com/lacework/go-sdk/commit/6e5c7bd61c01f1a3d84c9285d0d08d7a2dc6b3ef))
* chore(deps): bump github.com/hashicorp/terraform-exec (#704) (dependabot[bot])([9767bc8](https://github.com/lacework/go-sdk/commit/9767bc8a1c2531f936e685f0143400017198d6a7))
* chore(deps): bump go.uber.org/zap from 1.19.1 to 1.21.0 (#700) (dependabot[bot])([b8be0b0](https://github.com/lacework/go-sdk/commit/b8be0b0fef410a5a6179495c7e2a4a44d5e9b3c6))
* chore(deps): bump github.com/hashicorp/go-version from 1.3.0 to 1.4.0 (#702) (dependabot[bot])([1729f4c](https://github.com/lacework/go-sdk/commit/1729f4cb513d8b32d6a7133f8df72bc3afb1f3a9))
* chore(deps): bump github.com/briandowns/spinner from 1.13.0 to 1.17.0 (#701) (dependabot[bot])([9d18563](https://github.com/lacework/go-sdk/commit/9d185636fa0fa7c8d46dedae01d349d18d0a60c1))
* ci(cli): release CentOS end-of-life problem (#699) (Salim Afiune)([4bc4a7c](https://github.com/lacework/go-sdk/commit/4bc4a7ce8c37f1f4f1a21b058fe315ba855687b6))
* ci: version bump to v0.25.1-dev (Lacework)([20551d1](https://github.com/lacework/go-sdk/commit/20551d1b3b748bedc696d190acb7c13cfffc54e4))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
