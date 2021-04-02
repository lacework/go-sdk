# Release Notes
Another day, another release. These are the release notes for the version `v0.4.0`.

## Features
* feat: add new flags to fail host and container vuln scans (#352) (Darren)([1209d59](https://github.com/lacework/go-sdk/commit/1209d59138365a2aa86e9e0f8af66586542201db))
## Refactor
* refactor: Remove deprecated ListRange function (Darren Murray)([2a0019f](https://github.com/lacework/go-sdk/commit/2a0019f28dba6fcc6d470c2f6d54136581aac817))
* refactor: Remove deprecated vulnerability report command (Darren Murray)([d0de656](https://github.com/lacework/go-sdk/commit/d0de6564d7f6f87f682db902ed51e7aba0efbe87))
* refactor: Remove deprecated pdfName flag (Darren Murray)([94b8f67](https://github.com/lacework/go-sdk/commit/94b8f67aca806265de8ca8e6128fc3f137381efb))
## Bug Fixes
* fix(cli): mask ServiceNow integration password (#359) (Salim Afiune)([8e37697](https://github.com/lacework/go-sdk/commit/8e376976df0d429a403319c43e46ff1ad87dbc97))
* fix(api): update time format in API response (#360) (Salim Afiune)([a0045b4](https://github.com/lacework/go-sdk/commit/a0045b440a4164ff5ffd5c42d5a049471893b18f))
## Other Changes
* chore(deps): bump github.com/spf13/viper from 1.6.2 to 1.7.1 (#364) (dependabot[bot])([4f8f838](https://github.com/lacework/go-sdk/commit/4f8f838d3ed4352d5d660d75ce99a121ffc73cf0))
* chore(deps): bump github.com/stretchr/testify from 1.6.1 to 1.7.0 (#363) (dependabot[bot])([eed32e8](https://github.com/lacework/go-sdk/commit/eed32e8de351736c670594affda7d28002b7b778))
* chore(deps): bump github.com/fatih/color from 1.9.0 to 1.10.0 (#365) (dependabot[bot])([008253b](https://github.com/lacework/go-sdk/commit/008253bc78d27e2cb31434082284d6b64e0e51a4))
* chore(deps): bump github.com/spf13/cobra from 1.0.0 to 1.1.3 (#362) (dependabot[bot])([f3427ea](https://github.com/lacework/go-sdk/commit/f3427ea014da326218cb639f24f1ca58ae3e0515))
* chore: Validation script for git commit message conventions (#355) (Darren)([7fe9678](https://github.com/lacework/go-sdk/commit/7fe9678e854890a6a393df722007be38290630ff))
* ci: fix githooks and release scripts (#366) (Salim Afiune)([9036013](https://github.com/lacework/go-sdk/commit/90360137f823a4efa9651ade65ffaedfb4f61861))
* ci: switch over our CI to CodeFresh (#357) (Salim Afiune)([734acb7](https://github.com/lacework/go-sdk/commit/734acb7971c35ae79d6435ecc89d5a7fb4d493e7))
* ci: Fix prepare_test_resources script (Darren Murray)([e89beb6](https://github.com/lacework/go-sdk/commit/e89beb63d11941c67a97d91044fcbd6fe7ee2ad1))
* ci: Add setup_remote_docker to prepare-test-resources job (Darren Murray)([60628cd](https://github.com/lacework/go-sdk/commit/60628cd7731d0c15693a9b274f087a0b3eba590c))
* test(cli): lacework query command (#358) (hazedav)([d5473df](https://github.com/lacework/go-sdk/commit/d5473df7369b23d4e04dccb287de9b509b646342))
* test: removed Azure tenants from test account (Salim Afiune Maya)([e181ed1](https://github.com/lacework/go-sdk/commit/e181ed174fa9096e7c7a4a01ee95cc22951ba55d))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
