# Release Notes
Another day, another release. These are the release notes for the version `v0.25.0`.

## Features
* feat(api): add registry notifications field to docker v2 (#696) (Darren)([aac6801](https://github.com/lacework/go-sdk/commit/aac68016b548cc9c2799bb5b12c1f6e0f16aa8bd))
* feat(api): abstract primitives to interact with /v2/Vulnerabilities (#694) (Salim Afiune)([cd2622a](https://github.com/lacework/go-sdk/commit/cd2622a6d642c14e542e9997621b06f284a21245))
* feat(policies): add tags support (hazedav)([2a0ad67](https://github.com/lacework/go-sdk/commit/2a0ad670b86ec304b1acbbf72cb3c127348e7409))
* feat(cli): be able to change global LW_LQL_OPERATOR (Salim Afiune Maya)([a4c7746](https://github.com/lacework/go-sdk/commit/a4c77462c93d298adaeb7a80a6945a5d96e986ad))
* feat(cli): filter list of agents (Salim Afiune Maya)([8e200a0](https://github.com/lacework/go-sdk/commit/8e200a0c618c95a5206112565237a8cd7caac6eb))
* feat(api): implement Entities/Images/serch endpoint (Salim Afiune Maya)([33b58ef](https://github.com/lacework/go-sdk/commit/33b58ef0a2433b802baad3bae2d0586c422a4998))
* feat(api): implement Entities/Users/serch endpoint (Salim Afiune Maya)([b86f893](https://github.com/lacework/go-sdk/commit/b86f893be1d8b6fc9958611599a82bdef64b33d7))
* feat(api): support v2/Entities/Users (Salim Afiune Maya)([40c71d7](https://github.com/lacework/go-sdk/commit/40c71d7e2eba75bf5696c475265226bacc7a0943))
* feat(api): v2 pagination (Salim Afiune Maya)([255079d](https://github.com/lacework/go-sdk/commit/255079dd7c09873570e8412eab26ab01d03f96f5))
* feat(cli): new agent list command (Salim Afiune Maya)([8b056f5](https://github.com/lacework/go-sdk/commit/8b056f5100764ec3cca5413df38130014cad416d))
* feat(api): search data-export entities (Salim Afiune Maya)([5b87d33](https://github.com/lacework/go-sdk/commit/5b87d337bb585d18216715d89467203c1e94f42f))
* feat(cli): Manage Team Members (#671) (Darren)([6930449](https://github.com/lacework/go-sdk/commit/6930449b156bcf33a92fd4ef885c7ffe166ef369))
* feat(cli): Manage Vulnerability Exceptions (#663) (Darren)([0d62f11](https://github.com/lacework/go-sdk/commit/0d62f11b2a5264314489aa9bd23d7ca2a19d0b28))
* feat(cli): add --cascade flag to delete policy and query (#677) (hazedav)([6842a6b](https://github.com/lacework/go-sdk/commit/6842a6b162113471391118d34a03b1a022c5694d))
* feat(cli): display status of Azure compliance subscriptions (Salim Afiune Maya)([c479863](https://github.com/lacework/go-sdk/commit/c479863e71d63f301ea2d926db085f2eb8af2a4a))
* feat(cli): display status of GCP compliance projects (Salim Afiune Maya)([30e8448](https://github.com/lacework/go-sdk/commit/30e844823c26deee4ec01c3d40e2fdbc05111d6c))
* feat(cli): display status of AWS compliance accounts (Salim Afiune Maya)([e43978c](https://github.com/lacework/go-sdk/commit/e43978c9d40fd444e238deaa61059d6184df544d))
## Refactor
* refactor(cli): Refactor generate & generate_aws (#691) (Ross)([3fabe7b](https://github.com/lacework/go-sdk/commit/3fabe7b330284b0fbe044ea6af3a9b05c7e35b72))
* refactor(api): rename Pagination interface to Pageable (Salim Afiune Maya)([3b1e881](https://github.com/lacework/go-sdk/commit/3b1e881f52f0ee399084c072eaf53020a3510630))
* refactor(api): add logging to NextPage() func (Salim Afiune Maya)([810b723](https://github.com/lacework/go-sdk/commit/810b7231e0c682e6725f86d0f993c98ff7330c6e))
* refactor(cli): Refactor generate & generate_aws (#689) (Ross)([c2d4241](https://github.com/lacework/go-sdk/commit/c2d42418c7c1e56848048e475ad40b27694a7a94))
* refactor(cli): abstract shared funcs to internal package (hazedav)([5aa880d](https://github.com/lacework/go-sdk/commit/5aa880d5ccb09c75a9039cccf8313db30a020dff))
* refactor(cli): abstract shared funcs to internal package (hazedav)([5617a2d](https://github.com/lacework/go-sdk/commit/5617a2db0e1ccccee8e4d6df3b94c1a20d00bfc6))
* refactor(cli): rename 'aws' s3/cloudwatch to 'amazon' in cli integration prompt (#657) (Darren)([2c8757d](https://github.com/lacework/go-sdk/commit/2c8757d7b021c46e7c43c5a4eac353eafed5e269))
## Bug Fixes
* fix(cli): remove text wrapping on list-sources output (#692) (Darren)([49ca2ef](https://github.com/lacework/go-sdk/commit/49ca2ef0fb1826d38ce8e4d445ef14fc1e36d185))
* fix(cli): validate user filters via --filters (Salim Afiune Maya)([6ba02b0](https://github.com/lacework/go-sdk/commit/6ba02b086e2d650f36f4bf74fb482589244f9aec))
* fix(api): omit empty fields in v2 SearchFilter (Salim Afiune Maya)([5f4d2e7](https://github.com/lacework/go-sdk/commit/5f4d2e7065af5018c2768f898a333e7067640f69))
* fix(cli): check the account has agents to display (Salim Afiune Maya)([7e4ca27](https://github.com/lacework/go-sdk/commit/7e4ca27d43f9a6a9118a696e857884da4666b933))
* fix(lwgenerate): support lists for module attributes (#687) (Matt Cadorette)([6adfeeb](https://github.com/lacework/go-sdk/commit/6adfeeb1945f694491346e6f004251d900cecf5a))
## Other Changes
* style(cli): add breadcrumb to query list-sources cmd (#685) (Darren)([66ee797](https://github.com/lacework/go-sdk/commit/66ee797ea33ab813c2c5795cf374e4647889c360))
* chore(cli): update spinner message for compliance list cmds (Salim Afiune Maya)([c3d4aee](https://github.com/lacework/go-sdk/commit/c3d4aee0f6e5eb7066ab16b6b78b5657c55db8b3))
* ci: fix v1/v2 environment variables (#688) (Salim Afiune)([de4699a](https://github.com/lacework/go-sdk/commit/de4699a26a249cd325405f49602c6560d43b1822))
* ci: version bump to v0.24.1-dev (Lacework)([2df7446](https://github.com/lacework/go-sdk/commit/2df744647538fd2371b34cc0c1cc3dcc42cd71a7))
* test(cli): for compliance az list command (Salim Afiune Maya)([a98599c](https://github.com/lacework/go-sdk/commit/a98599c529efee480031f19f60042148a8bd31f2))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
