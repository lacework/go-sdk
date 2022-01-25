# Release Notes
Another day, another release. These are the release notes for the version `v0.24.0`.

## Features
* feat(cli): make policyID optional for policy create (hazedav)([38006bf](https://github.com/lacework/go-sdk/commit/38006bf39a3cae293595e13c6b6bfb9d05936f1d))
* feat(lwgenerate): Add initial hcl generation support GCP - ALLY-393 (#654) (Ross)([fb34a64](https://github.com/lacework/go-sdk/commit/fb34a64b0ddc71bc2140b42c66fa0f38160e762a))
* feat(cli): add csv output support to vuln reports (#649) (Matt Cadorette)([c76095e](https://github.com/lacework/go-sdk/commit/c76095e88c60f02c02b449c09ffbeaf2dd5bf1b1))
* feat(api): Vulnerability Exceptions v2 Service (#627) (Darren)([ef65e23](https://github.com/lacework/go-sdk/commit/ef65e232b7c7c650562db5cded0b3ee31a461a55))
* feat(cli): add apply flag for aws iac generation & fix cache issues (#642) (Matt Cadorette)([c5014bb](https://github.com/lacework/go-sdk/commit/c5014bb73f1db4977a9500dc4f92e686fa9d1c74))
## Refactor
* refactor(test_resources): remove Terraform code dep (#636) (Salim Afiune)([aca0101](https://github.com/lacework/go-sdk/commit/aca01017734598d64526b9a7008fe9171d74fe5c))
## Bug Fixes
* fix: expiration time (#656) (Darren)([4726f68](https://github.com/lacework/go-sdk/commit/4726f68ca5731208e140e006c98716aba3d626a3))
* fix: non_os_package_support omitted from json request when false (#655) (Darren)([c8ab03a](https://github.com/lacework/go-sdk/commit/c8ab03ae4ad1616ef172db3230fd3f1baff6e389))
* fix: handle unsupported integration types (#645) (Darren)([6c1a5d0](https://github.com/lacework/go-sdk/commit/6c1a5d05f937eaac89893e5cb483e7d6c045e17a))
* fix(cli): clarify container error message (#643) (Salim Afiune)([74592ff](https://github.com/lacework/go-sdk/commit/74592ff8fadf8f78606c940d347d9024dfbe8a8c))
## Documentation Updates
* docs(iac-generation): update aws iac tf command help (#653) (Matt Cadorette)([5c09072](https://github.com/lacework/go-sdk/commit/5c09072b7079f72be40fd154da471755d2eeccc7))
* docs(cli): Improve iac-generate/aws help message (#638) (Matt Cadorette)([127e221](https://github.com/lacework/go-sdk/commit/127e2212018fb2124b3d15ec2fcdd155e1326533))
## Other Changes
* chore(deps): bump github.com/hashicorp/hcl/v2 from 2.10.1 to 2.11.1 (#631) (dependabot[bot])([a82d22b](https://github.com/lacework/go-sdk/commit/a82d22b3a3d3ad6eec69a7d03a45a19e298be893))
* chore(deps): bump github.com/mitchellh/mapstructure from 1.4.2 to 1.4.3 (#630) (dependabot[bot])([4f08643](https://github.com/lacework/go-sdk/commit/4f0864393b611ce5d6527393ff2e8d5df64cfd24))
* chore(deps): bump github.com/peterbourgon/diskv/v3 from 3.0.0 to 3.0.1 (#610) (dependabot[bot])([c1792d5](https://github.com/lacework/go-sdk/commit/c1792d51bd563f852cf154a7fdfbf5d7560192e0))
* chore(deps): Bump GoImports version to v0.1.8 (#633) (Ross)([22f5283](https://github.com/lacework/go-sdk/commit/22f5283a93d0e3552c7343ad7c6104c0f3389ab8))
* chore(deps): Bump gotestsum version and use instead of go test (#632) (Ross)([7f05f44](https://github.com/lacework/go-sdk/commit/7f05f441ce579a9654e7809d75a7a81845e9b380))
* ci: test on windows systems  (#648) (Salim Afiune)([70fe066](https://github.com/lacework/go-sdk/commit/70fe0665af92afba20e28af7d1c3319043ba78b1))
* ci: version bump to v0.23.1-dev (Lacework)([cfc7bc0](https://github.com/lacework/go-sdk/commit/cfc7bc00f6ff515a0f8b64b8d6d6bb6ceb684fc0))
* test: fix host vuln test (#659) (Darren)([d28ea4c](https://github.com/lacework/go-sdk/commit/d28ea4c366eeca098547efc0c7230f30d9e7a141))
* test: avoid running v1/v2 migration on every test (#651) (Salim Afiune)([63735c2](https://github.com/lacework/go-sdk/commit/63735c2db0e518ff381eaa8949ff5fc33cfca64b))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
