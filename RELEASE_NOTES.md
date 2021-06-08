# Release Notes
Another day, another release. These are the release notes for the version `v0.8.0`.

## Features
* feat(cli): Add CSV rendering output for select commands (#424) (Matt Cadorette)([1bbe07c](https://github.com/lacework/go-sdk/commit/1bbe07c2055e61d6d34d0f4acb8bd9194358fcaf))
* feat(cli): Automate migration of v2 configuration (#431) (Salim Afiune)([2ac806a](https://github.com/lacework/go-sdk/commit/2ac806a91e8fea03eb6dec96e2ce8e8b3944b5e7))
* feat(pkg): to disseminate a domain URL (Salim Afiune Maya)([0770102](https://github.com/lacework/go-sdk/commit/07701029b5a057fd1b4da97d10683b27a381f365))
* feat(cli): v2 configure cmd undestands subaccounts (Salim Afiune Maya)([5e092f7](https://github.com/lacework/go-sdk/commit/5e092f70e601e081a02bb7ec18204c6373cdb6ec))
* feat: AWS Gov Cloud CloudTrail integration (#425) (Darren)([63097d5](https://github.com/lacework/go-sdk/commit/63097d55ae37616a61022f8cb0cba3cc5dcb9134))
* feat(cli): AWS GovCloud Config integration (#421) (Darren)([68d7087](https://github.com/lacework/go-sdk/commit/68d708713289c84606613a7261e015462c413299))
* feat(cli): lacework policy command (hazedav)([8c53e8e](https://github.com/lacework/go-sdk/commit/8c53e8eac2a06916dab4cf0062ea8ffc3781895e))
* feat: introducing first APIv2 endpoint UserProfile (#411) (Salim Afiune)([c183cce](https://github.com/lacework/go-sdk/commit/c183ccec205c473e461eae19d116f863a613ab9b))
## Refactor
* refactor(cli): remove shorthand of subaccount "u" (#436) (Salim Afiune)([2decc86](https://github.com/lacework/go-sdk/commit/2decc8610b77244c0f1b04370b34d368ea6b9206))
* refactor(cli): use new internal domain go package (Salim Afiune Maya)([1695993](https://github.com/lacework/go-sdk/commit/16959935dc6f876f82c970c9f1844ecd4068ccb7))
* refactor(cli): improve configure command (Salim Afiune Maya)([b14938b](https://github.com/lacework/go-sdk/commit/b14938be7ac4f22dc83d9659d9b062071f981cf9))
## Bug Fixes
* fix: Fix vuln container list-assessments json output (#430) (Darren)([02ff56d](https://github.com/lacework/go-sdk/commit/02ff56d37ab0fe0468aaeb1dd8afe28a53abbf0b))
## Other Changes
* ci: fix a bunch of things with our CI pipeline (#432) (Salim Afiune)([af603e1](https://github.com/lacework/go-sdk/commit/af603e145fe138a79fd21445da75087b2de93a49))
* test(integration): configure command for OrgAdmins (Salim Afiune Maya)([6504ad9](https://github.com/lacework/go-sdk/commit/6504ad991058eac4c0db2cd238d42580f9193e43))
* test: new DEBUG env variable for interactive shell (Salim Afiune Maya)([bd58489](https://github.com/lacework/go-sdk/commit/bd584894848e27593349c4ee9a779f4ee7fff9b2))
* test(integration): fix missed help message update (#419) (Salim Afiune)([f1fcbf3](https://github.com/lacework/go-sdk/commit/f1fcbf32b3dc571978b42293509f4def1d7c56cb))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
