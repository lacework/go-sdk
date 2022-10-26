# Release Notes
Another day, another release. These are the release notes for the version `v0.43.0`.

## Features
* feat(cli): parallelize installation with goroutines (#948) (Nick Schmeller)([6ee9d6d1](https://github.com/lacework/go-sdk/commit/6ee9d6d1c3f5eaeca36e2a412b544c23a221bb52))
* feat(cdk): dev-mode command (#957) (Salim Afiune)([faff0ae7](https://github.com/lacework/go-sdk/commit/faff0ae724206eff05d173ae3516fc5095330cc7))
* feat(api): New AwsSidekickOrg Cloud Account type (#951) (Teddy Reed)([4ab28fa0](https://github.com/lacework/go-sdk/commit/4ab28fa09cd0c2d480b0ab1e86237be3cae5cf36))
* feat(cli): find and fix AWS instances w/o agents (#934) (Nick Schmeller)([4e507c0b](https://github.com/lacework/go-sdk/commit/4e507c0b89edde2a8f04e03581a92aae9789bdc9))
* feat: GKE audit-log (#927) (jonathan stewart)([76e6cc88](https://github.com/lacework/go-sdk/commit/76e6cc88ed923dbabe29b5c2f68d7e382c93f0ed))
* feat: generate command rework (#928) (jonathan stewart)([56cbdcd1](https://github.com/lacework/go-sdk/commit/56cbdcd192a79148c24e735c2ae814fd1b70bb66))
* feat(cdk): component dev mode (#911) (Salim Afiune)([d477e688](https://github.com/lacework/go-sdk/commit/d477e688f3a360ea6e6cae6b08d8ac0a92df4a1c))
* feat: add new parameters to gcp iac generation (#918) (jonathan stewart)([e353487b](https://github.com/lacework/go-sdk/commit/e353487b3d58394fefcdf6c90aa05b22c919bafd))
## Bug Fixes
* fix(cli): token from flag into arg for `aws-install` (#952) (Nick Schmeller)([0633e5cb](https://github.com/lacework/go-sdk/commit/0633e5cb75a149214226c75634ffd9c8c1eaf44c))
* fix: gke integration test flags (#949) (jonathan stewart)([7cdd53e5](https://github.com/lacework/go-sdk/commit/7cdd53e5ff29fb10d962ba4f8c09cda3cbfa8599))
* fix(cli): update `aws-install` command in help msg (#947) (Nick Schmeller)([c1396e19](https://github.com/lacework/go-sdk/commit/c1396e197ef22c9c1c4c467a0706284563647e78))
* fix(lwcomponent): better JSON handling components (#946) (hazedav)([8edb92de](https://github.com/lacework/go-sdk/commit/8edb92dee96c769ca40246e0397d229de2e775a2))
* fix: aws terraform sns variables (#944) (jonathan stewart)([c8ead2c2](https://github.com/lacework/go-sdk/commit/c8ead2c2b7a6f74b88aa220a7fceaee92d6df1fb))
* fix: azure Lacework provider & profile flag (#938) (jonathan stewart)([75ed6ba4](https://github.com/lacework/go-sdk/commit/75ed6ba4c804104278a6819e13e5248a7c9f5e9c))
* fix: azure interactive caching (#939) (jonathan stewart)([cff15ed2](https://github.com/lacework/go-sdk/commit/cff15ed268990c53b5f292bacc995a7b62f08dd1))
* fix: generate TF overwrite (#935) (jonathan stewart)([b33ad901](https://github.com/lacework/go-sdk/commit/b33ad90181e69aad1e3a92da49eda7a2936a117f))
* fix(lwrunner): create ~/.ssh directory if not exist (#933) (Salim Afiune)([bfc90990](https://github.com/lacework/go-sdk/commit/bfc909902c92a4dc5490f03d69357fa3ab2e6c84))
* fix: gcr limitByRep json field name (#924) (Darren)([0167eabe](https://github.com/lacework/go-sdk/commit/0167eabe68c13172db82d4d7aa33cbe5535c32b2))
* fix: gcr limit by repo json field name (#923) (Darren)([6f13c4c0](https://github.com/lacework/go-sdk/commit/6f13c4c07e8419dcfe97309cfac42126ecdcdddc))
* fix: limit_by_labels json field name (#922) (Darren)([857ea019](https://github.com/lacework/go-sdk/commit/857ea019c8ef55ff92626da3fb90a7b8ab18eb62))
* fix: omitempty fields ctr registry v2 (#921) (Darren)([14a8a225](https://github.com/lacework/go-sdk/commit/14a8a2250e18b3f9f679ad6ec63057d140656a52))
* fix(lql): fix handling of large numbers in json (#919) (hazedav)([6c934ad3](https://github.com/lacework/go-sdk/commit/6c934ad353119d31c30e5792a1db017673aac9d9))
* fix: DockerhubV2Data fix fields (#920) (Darren)([cbf3e395](https://github.com/lacework/go-sdk/commit/cbf3e3959b56799a832daaf21373c7b7a7dc4c9f))
* fix(ALLY-1188): Enable the disable of enable_ebla (#916) (jonathan stewart)([ba3c3562](https://github.com/lacework/go-sdk/commit/ba3c356216a3970743ab91d3fb9bb4fed5a1461a))
## Other Changes
* chore(deps): bump github.com/honeycombio/libhoney-go (#931) (dependabot[bot])([e5624291](https://github.com/lacework/go-sdk/commit/e56242913e0b7973957b2de2f11668eddbe06720))
* chore(deps): bump github.com/abiosoft/colima from 0.4.4 to 0.4.5 (#930) (dependabot[bot])([3aa5ff93](https://github.com/lacework/go-sdk/commit/3aa5ff9344493cbe936a49ea7277a59c61e67e36))
* chore(deps): bump github.com/hashicorp/hcl/v2 from 2.13.0 to 2.14.1 (#925) (dependabot[bot])([e872b5a4](https://github.com/lacework/go-sdk/commit/e872b5a4058380c61be8b10a47c11a11ad2621b2))
* chore(deps): bump go.uber.org/zap from 1.21.0 to 1.23.0 (#901) (dependabot[bot])([593cdbaa](https://github.com/lacework/go-sdk/commit/593cdbaaef8cef5dbe21bfc77c2ead7ad0e7cea8))
* chore(deps): bump github.com/hashicorp/terraform-exec (#908) (dependabot[bot])([976c486c](https://github.com/lacework/go-sdk/commit/976c486c318a6df7ab14014ef9c2880062687da1))
* chore(deps): bump github.com/spf13/viper from 1.12.0 to 1.13.0 (#910) (dependabot[bot])([658b600a](https://github.com/lacework/go-sdk/commit/658b600ae729cf039eb6de1d0b74bd07b8c25cea))
* build: fix install-cli Makefile target to account for ARM builds (#950) (0xdeadbeefJERKY)([4fb46f23](https://github.com/lacework/go-sdk/commit/4fb46f23944bb176428116b297cf1ae56dd3e32e))
* ci: version bump to v0.42.1-dev (Lacework)([3044a87f](https://github.com/lacework/go-sdk/commit/3044a87f12b54618fef9c3b6308dc2c82fccba93))
* test(cli): disable compliance tests (#964) (Salim Afiune)([6989e67e](https://github.com/lacework/go-sdk/commit/6989e67e077dede960e002d923b11ee4001a1d19))
* test: update error msg test (#942) (Darren)([ad04af52](https://github.com/lacework/go-sdk/commit/ad04af525202fc883adee52ea706ccff8e38a36b))
* test(win): disable failing integration test (#932) (Salim Afiune)([7e2adfc0](https://github.com/lacework/go-sdk/commit/7e2adfc02e9cb29a41824cd6dc9ca31c7d8a49ff))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
