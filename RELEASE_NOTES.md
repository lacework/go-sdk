# Release Notes
Another day, another release. These are the release notes for the version `v0.10.0`.

## Features
* feat(cli): New compliance google list command  (#466) (Darren)([f3aec8d](https://github.com/lacework/go-sdk/commit/f3aec8d5a08801fa517bd9aec405f1fb2cdafed1))
* feat(cli): caching mechanism (Salim Afiune Maya)([426c3eb](https://github.com/lacework/go-sdk/commit/426c3eb6ce569dc06295e30d2ca51f2c05213a4d))
* feat(api): LifecycleCallbacks (Salim Afiune Maya)([7b0c03e](https://github.com/lacework/go-sdk/commit/7b0c03e8a5bbb12af3dfbe7494fda5ea2030b9d5))
* feat: new internal/format go package (Salim Afiune Maya)([aa49474](https://github.com/lacework/go-sdk/commit/aa49474d6f9cdfa28d5e829e43e5d49c0fbaafe4))
* feat: install diskv go package (Salim Afiune Maya)([b0f7c43](https://github.com/lacework/go-sdk/commit/b0f7c43ac521a0830a109f6af44bfaf8fc3c00c8))
* feat(cli): apiv2 support for lacework policy (hazedav)([efaa50e](https://github.com/lacework/go-sdk/commit/efaa50e4a5023e65f10c9145355140092d4d7806))
* feat(cli): apiv2 support for Queries (hazedav)([e3b2cfa](https://github.com/lacework/go-sdk/commit/e3b2cfaffbf601bf9baa583e783cff699aabbe21))
* feat(cli): add --organization global flag (Salim Afiune Maya)([408a340](https://github.com/lacework/go-sdk/commit/408a34003d54c4e5098e8148836cc6bdeb0e3bb0))
* feat(api): Add AwsCtSqs Cloud Account Get/Update funcs (Salim Afiune Maya)([435686e](https://github.com/lacework/go-sdk/commit/435686ed61f9d3bf2691358d691bdef9594d3f93))
* feat(api): Add V2 CloudAccountsService (Salim Afiune Maya)([656b257](https://github.com/lacework/go-sdk/commit/656b2577ead5f84a8445de55f0aaee160ba79ddd))
* feat(api): WithOrgAccess Option to access org level (Salim Afiune Maya)([15f672e](https://github.com/lacework/go-sdk/commit/15f672e8765fa9102bc1139e30831b0ee9392bc9))
* feat(api): CopyClient() to duplicate an API client (Salim Afiune Maya)([084fe25](https://github.com/lacework/go-sdk/commit/084fe254d3edeed8e1229b4285b378a241662664))
* feat(api): propagate APIv2 error messages (#453) (Salim Afiune)([9a4c697](https://github.com/lacework/go-sdk/commit/9a4c6971f10e0b8945e780e09dccf82528c83391))
## Refactor
* refactor: migrate Agent Access Tokens to API v2  (#461) (Salim Afiune)([e08d947](https://github.com/lacework/go-sdk/commit/e08d94788004211f5b7d982539a9bad60d0c8653))
* refactor(cli): create AwsCtSqs CloudAccounts using APIv2 (Salim Afiune Maya)([b8b1605](https://github.com/lacework/go-sdk/commit/b8b1605eaf3e09c8213996430301fcc7022e7cc6))
* refactor(cli): improved LQL service func names (#438) (hazedav)([08f7ce4](https://github.com/lacework/go-sdk/commit/08f7ce4dbcd25e364c48eff92c00e07e814f69d4))
## Bug Fixes
* fix(api): add 'description' field for Host Vulnerability (#481) (Salim Afiune)([dba83cb](https://github.com/lacework/go-sdk/commit/dba83cbfbedbcb2542e91f6903899fe862ecad02))
* fix(cli): remove subaccount during configure (#474) (Salim Afiune)([f228f08](https://github.com/lacework/go-sdk/commit/f228f089fe8923065d559b177cd994fb802f4156))
* fix: check http status code before erase cache (Darren Murray)([7f1a29d](https://github.com/lacework/go-sdk/commit/7f1a29d26ce49582fa5b70a00ded3ece38c9a096))
* fix(cli): non-interactive should not run daily version check (#462) (Salim Afiune)([4bb7e9f](https://github.com/lacework/go-sdk/commit/4bb7e9feaf8baff84c2f94408f633abb69f98d16))
## Other Changes
* chore(deps): revert bump AlecAivazis/survey from 2.2.14 to 2.2.12 (#476) (Salim Afiune)([2d2176d](https://github.com/lacework/go-sdk/commit/2d2176d3a6a455c2de78c3c520f3a820561f7a8a))
* chore(deps): revert bump briandowns/spinner from 1.16.0 to 1.12.0 (#475) (Salim Afiune)([f2e19fb](https://github.com/lacework/go-sdk/commit/f2e19fbe4b64769c9dfb3fe50eb1dcfaa173d3d3))
* chore(deps): bump github.com/spf13/cobra from 1.1.3 to 1.2.1 (#469) (dependabot[bot])([6725a81](https://github.com/lacework/go-sdk/commit/6725a811d7815fc1047ead5490002069a54bc660))
* chore(deps): bump github.com/honeycombio/libhoney-go (#433) (dependabot[bot])([6a833a1](https://github.com/lacework/go-sdk/commit/6a833a16863f434c232b2638a6b90130976d3e22))
* chore(deps): bump github.com/spf13/viper from 1.7.1 to 1.8.1 (#459) (dependabot[bot])([7296b4a](https://github.com/lacework/go-sdk/commit/7296b4a4ca603baa0871c36aa2a079c078d568a0))
* chore(deps): bump github.com/briandowns/spinner from 1.12.0 to 1.16.0 (#452) (dependabot[bot])([06050a2](https://github.com/lacework/go-sdk/commit/06050a27fc8ad5a5554dfa2f041cc071d7d43253))
* chore(deps): bump go.uber.org/zap from 1.14.1 to 1.18.1 (#460) (dependabot[bot])([d311493](https://github.com/lacework/go-sdk/commit/d31149382f623616885f200d349c300ced8aad86))
* chore(deps): bump github.com/AlecAivazis/survey/v2 from 2.2.12 to 2.2.14 (#468) (dependabot[bot])([86b6a0f](https://github.com/lacework/go-sdk/commit/86b6a0f7e373d95d572243516183508a67897d4c))
* chore(api): rename cloud_account_aws.go -> cloud_accounts_aws_ct_sqs.go (Salim Afiune Maya)([dfd9ce2](https://github.com/lacework/go-sdk/commit/dfd9ce2e8569432e62c5538104016de2ab5b1059))
* ci: decouple integrations tests (#463) (Darren)([da7d429](https://github.com/lacework/go-sdk/commit/da7d429f53acaeaf2d1f30d8d0ae53d0f898ad65))
* ci: prepare-integrations pipeline (#457) (Darren)([d0e2436](https://github.com/lacework/go-sdk/commit/d0e2436e22d939d05dbe0cf7bf4fdc745c91dae1))
* test(cli): cache mechanism (Salim Afiune Maya)([6a0e9ca](https://github.com/lacework/go-sdk/commit/6a0e9ca6b59f50d015f99cbc6c0f2b9804f71e20))
* test(cli): fixing policy integration url (hazedav)([191a7ed](https://github.com/lacework/go-sdk/commit/191a7ed607632ef4b57d833415ea1a9e7cd76922))
* test(cli): apiv2 support for policy integration (hazedav)([68d324c](https://github.com/lacework/go-sdk/commit/68d324c6a48ca22952a580ee9fdcda305714f2b1))
* test(cli): add apiv2 query (#465) (hazedav)([bbd5b94](https://github.com/lacework/go-sdk/commit/bbd5b946d38213db2763be3e9a43bc643dc07e95))
* test(api): CloudAccount interface and CRUD tests (Salim Afiune Maya)([b0b103e](https://github.com/lacework/go-sdk/commit/b0b103e0116172b98be47f594b6ee7933555197c))
* test(cli): updated naming for test resources (hazedav)([4c3a389](https://github.com/lacework/go-sdk/commit/4c3a389cf2a684cda2c1fa1c61929bd29e4fb9e5))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
