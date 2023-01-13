# v1.5.2

## Other Changes
* chore: re-add compliance enable/disable commands (#1106) (Ross)([f25c72c9](https://github.com/lacework/go-sdk/commit/f25c72c9902bbdcdb22ccf98d28f9e39c57c9dc5))
* ci: version bump to v1.5.2-dev (Lacework)([9571b913](https://github.com/lacework/go-sdk/commit/9571b91369161598dcb88e8dad22835ea723976e))
---
# v1.5.1

## Bug Fixes
* fix: policy-exception json keys (#1107) (jonathan stewart)([6ea84168](https://github.com/lacework/go-sdk/commit/6ea84168ef4f4dbf01c07f6c619686bd2c27b9c3))
* fix: compliance search (#1105) (Darren)([5faea9b8](https://github.com/lacework/go-sdk/commit/5faea9b859f76eacd9a97b8ce3ffdf512d62345a))
## Other Changes
* chore(deps): bump github.com/aws/aws-sdk-go-v2/config (#1103) (dependabot[bot])([daaee23b](https://github.com/lacework/go-sdk/commit/daaee23b9c06a1b66741389a4a065c3a78cda5ab))
* chore(deps): bump cloud.google.com/go/oslogin from 1.7.0 to 1.9.0 (#1102) (dependabot[bot])([f38bf58e](https://github.com/lacework/go-sdk/commit/f38bf58e2b87d332774a6615c532b43d39f7db73))
* ci: version bump to v1.5.1-dev (Lacework)([4bc4663f](https://github.com/lacework/go-sdk/commit/4bc4663fe44dbaefa57ab2c4ee1f765c0e14f791))
---
# v1.5.0

## Features
* feat(RAIN-45215): Add Hidden Azure & GCP suppressions list commands (#1100) (Ross)([796bd92c](https://github.com/lacework/go-sdk/commit/796bd92c4f870d8fb002df5f5aec9c0d8107541d))
## Other Changes
* ci: version bump to v1.4.2-dev (Lacework)([8ecf94b3](https://github.com/lacework/go-sdk/commit/8ecf94b30bade1e0d0670809576246283255751f))
---
# v1.4.1

## Refactor
* refactor: Use alpine as our base Docker image (#1097) (Salim Afiune)([d89d748a](https://github.com/lacework/go-sdk/commit/d89d748a6ace003b5f23d732a4099303954fd62e))
## Other Changes
* style(cdk): rename lifecycle event cdk-cleanup (#1092) (Salim Afiune)([4969c91a](https://github.com/lacework/go-sdk/commit/4969c91a08cf680ec07bcb66e49979eca885c1af))
* ci: update chocolatey scripts (#1093) (Darren)([606fc228](https://github.com/lacework/go-sdk/commit/606fc228cf31ea8645b2c5da2b39f5aa9333e6ff))
* ci: version bump to v1.4.1-dev (Lacework)([d7df6551](https://github.com/lacework/go-sdk/commit/d7df6551435587b4e7b529ff6627ac419b09e9a9))
---
# v1.4.0

## Features
* feat(GROW-1058): Improve management of component versions (#1088) (Edoardo Pirovano)([952f2c30](https://github.com/lacework/go-sdk/commit/952f2c30cac5cacbf02494eab031d81f6b95dee5))
* feat(cdk): pass LW_CLI_VERSION env variable (#1085) (Salim Afiune)([f0382ae8](https://github.com/lacework/go-sdk/commit/f0382ae850c95d9956d902138d80465d5e7fad05))
* feat(lwlogger): add merge logger function (#1087) (Kolbeinn)([1e569be8](https://github.com/lacework/go-sdk/commit/1e569be8c5d53335605e2540b02f0b94ec21925e))
* feat(cli): pass LW_API_SERVER_URL env variable to components (#1077) (Kolbeinn)([c85d1541](https://github.com/lacework/go-sdk/commit/c85d154189d8ce830bf1b28bbc505fd0b48f692d))
## Bug Fixes
* fix(lwlogger): do not override explicit level argument (#1084) (Kolbeinn)([452dff83](https://github.com/lacework/go-sdk/commit/452dff8325497eae61a2ba76f59cd52eaf535678))
## Other Changes
* chore(deps): Run `make go-vendor` (#1094) (Edoardo Pirovano)([8beaf31f](https://github.com/lacework/go-sdk/commit/8beaf31fa3ba37de5d5b2dffd5bd167157580131))
* chore(deps): bump github.com/aws/aws-sdk-go-v2/service/ec2 (#1083) (dependabot[bot])([c5a04eab](https://github.com/lacework/go-sdk/commit/c5a04eab6aafee96ab7c2bda4c66fdf986ed75e4))
* chore(deps): bump github.com/honeycombio/libhoney-go (#1081) (dependabot[bot])([b9217256](https://github.com/lacework/go-sdk/commit/b921725609140ced50d5e0bfa9f91395129de34b))
* chore(deps): bump github.com/aws/aws-sdk-go-v2/config (#1075) (dependabot[bot])([58711b27](https://github.com/lacework/go-sdk/commit/58711b27dd08cdd4572b79b14c2792615ca438fd))
* chore(deps): bump github.com/hashicorp/hcl/v2 from 2.14.1 to 2.15.0 (#1044) (dependabot[bot])([44586503](https://github.com/lacework/go-sdk/commit/44586503459658e6be01f567861edec0de664fad))
* chore(deps): bump github.com/zclconf/go-cty from 1.11.0 to 1.12.1 (#1020) (dependabot[bot])([16a79870](https://github.com/lacework/go-sdk/commit/16a798704cb9dd638118d28321efb7c999185867))
* ci: version bump to v1.3.2-dev (Lacework)([9a254aee](https://github.com/lacework/go-sdk/commit/9a254aee3c0096f35eb498369cd1a739fe3db078))
---
# v1.3.1

## Refactor
* refactor: Release of docker images (#1074) (Salim Afiune)([3e198981](https://github.com/lacework/go-sdk/commit/3e1989815e31f3a403eb17afb68263df813222d4))
## Bug Fixes
* fix: avoid injecting empty alertProfile (#1078) (Salim Afiune)([485cc4cf](https://github.com/lacework/go-sdk/commit/485cc4cfb1b69b2f5bf958a0a5c5d1d0b8010641))
## Other Changes
* ci: version bump to v1.3.1-dev (Lacework)([60ff20ce](https://github.com/lacework/go-sdk/commit/60ff20ce3cf9eae4dea666ee004ac2648797a6da))
---
# v1.3.0

## Features
* feat(cli): add filter for fixable alerts (#1073) (hazedav)([98f7684c](https://github.com/lacework/go-sdk/commit/98f7684c9d511707015360f4a25ee55fb1fab172))
* feat: Add support from hidden suppressions migration command (#1069) (Ross)([2441d452](https://github.com/lacework/go-sdk/commit/2441d45216f276623fad92d9dee23c89a71138e7))
* feat: Display supported policy exception constraints (#1068) (Ross)([c6ebd686](https://github.com/lacework/go-sdk/commit/c6ebd686ee8d42d333e5e6f501fe6d57ff0122e3))
* feat: Add support from hidden suppressions api (#1066) (Ross)([69a2247a](https://github.com/lacework/go-sdk/commit/69a2247ad596faa7fff8586c850acac885c1dc04))
## Bug Fixes
* fix: revert compliance csv report_type field format (#1057) (Darren)([123d46b2](https://github.com/lacework/go-sdk/commit/123d46b2b3e1fefc67d78cd333a0ea4b03395e10))
* fix(cli): Ensure log level is correctly passed to components (#1071) (Kolbeinn)([0134e0a2](https://github.com/lacework/go-sdk/commit/0134e0a29663d97560280e9a07e83ccc13bb339a))
## Other Changes
* ci: version bump to v1.2.1-dev (Lacework)([24e59cb3](https://github.com/lacework/go-sdk/commit/24e59cb378b477e92c7d6ce4fd4fee500371473a))
---
# v1.2.0

## Features
* feat: add proxy scanner token management (#1018) (Credibleforce)([06b93071](https://github.com/lacework/go-sdk/commit/06b930715167c526ce9f8301d9a0cae23c8e9bdc))
## Bug Fixes
* fix: vuln ctr scan to use tag/digest argument to filter (#1070) (djmctavish)([1ec57f67](https://github.com/lacework/go-sdk/commit/1ec57f6799020755679bd09d49023bf26de6f29f))
## Documentation Updates
* docs: remove obsolete 'event' generated docs (#1067) (sthetland)([391a1626](https://github.com/lacework/go-sdk/commit/391a162679feabf6ad612b878e9566ce3d15bca7))
* docs(cli): list supported hosts for `ec2ic` (#1053) (Nick Schmeller)([c45f1dc1](https://github.com/lacework/go-sdk/commit/c45f1dc118be7de5b07554d9c15d669e263c9e9c))
## Other Changes
* build(cli): Support aarch64 architecture (#1063) (Kolbeinn)([99d2eb27](https://github.com/lacework/go-sdk/commit/99d2eb27c641c895b195072fcb2d2e959d7641ea))
* ci: version bump to v1.1.3-dev (Lacework)([19462b11](https://github.com/lacework/go-sdk/commit/19462b116cdd9f0c88b4e49378f6d549fec6db11))
---
# v1.1.2

## Bug Fixes
* fix: fix 1.1.2 issue (#1061) (Darren)([20d90d4e](https://github.com/lacework/go-sdk/commit/20d90d4ed5b6da6d9860e278ccb5ff23dcad998c))
* fix: fix 1.1.2 release (#1059) (Darren)([07c7a7c3](https://github.com/lacework/go-sdk/commit/07c7a7c369fc7c4ac31a424b414bf8f2763bc5a5))
* fix(api): Gcp sidekick integration server token (#1054) (ammarekbote)([6a16fdf0](https://github.com/lacework/go-sdk/commit/6a16fdf0dc2268981a52bffa99d6693ffb60a81a))
* fix: generate success message (#1056) (jonathan stewart)([57c282c5](https://github.com/lacework/go-sdk/commit/57c282c5b30d2110e52bf9ec3e25865f2df8b37a))
* fix: generate EKS non-interactive (#1055) (jonathan stewart)([7fbb78fc](https://github.com/lacework/go-sdk/commit/7fbb78fc0a9cdbb11c045d449c0d77ad0099ce50))
## Documentation Updates
* docs: fix documentation autogeneration (#1050) (Darren)([80106749](https://github.com/lacework/go-sdk/commit/80106749c054a05bc3a3bf4fde00131ad4897a1d))
## Other Changes
* ci: version bump to v1.1.3-dev (Lacework)([14914582](https://github.com/lacework/go-sdk/commit/14914582905d4b0490144ec908c821e112353a7f))
* ci: version bump to v1.1.3-dev (Lacework)([5dda8bca](https://github.com/lacework/go-sdk/commit/5dda8bca17f3f83ac48245c72352a75b367a3b03))
* ci: version bump to v1.1.2-dev (Lacework)([75c91670](https://github.com/lacework/go-sdk/commit/75c9167046643f17e9636a68deb14f5dcd44f343))
---
# v1.1.2

## Bug Fixes
* fix: fix 1.1.2 release (#1059) (Darren)([07c7a7c3](https://github.com/lacework/go-sdk/commit/07c7a7c369fc7c4ac31a424b414bf8f2763bc5a5))
## Other Changes
* ci: version bump to v1.1.3-dev (Lacework)([5dda8bca](https://github.com/lacework/go-sdk/commit/5dda8bca17f3f83ac48245c72352a75b367a3b03))
---
# v1.1.2

## Bug Fixes
* fix: generate success message (#1056) (jonathan stewart)([57c282c5](https://github.com/lacework/go-sdk/commit/57c282c5b30d2110e52bf9ec3e25865f2df8b37a))
* fix: generate EKS non-interactive (#1055) (jonathan stewart)([7fbb78fc](https://github.com/lacework/go-sdk/commit/7fbb78fc0a9cdbb11c045d449c0d77ad0099ce50))
## Documentation Updates
* docs: fix documentation autogeneration (#1050) (Darren)([80106749](https://github.com/lacework/go-sdk/commit/80106749c054a05bc3a3bf4fde00131ad4897a1d))
## Other Changes
* ci: version bump to v1.1.2-dev (Lacework)([75c91670](https://github.com/lacework/go-sdk/commit/75c9167046643f17e9636a68deb14f5dcd44f343))
---
# v1.1.1

## Bug Fixes
* fix: fix vuln ctr fail_on_severity (#1048) (Darren)([dcf79d3f](https://github.com/lacework/go-sdk/commit/dcf79d3fabe9465deff1127f5730c9ae4d67b6a9))
## Other Changes
* ci: version bump to v1.1.1-dev (Lacework)([b1e7b28c](https://github.com/lacework/go-sdk/commit/b1e7b28ce432b053a093bb0951f538f393c5b3a4))
---
# v1.1.0

## Features
* feat(cli): find and fix GCP instances w/o agents (#1022) (Nick Schmeller)([08edf28a](https://github.com/lacework/go-sdk/commit/08edf28a69693926b5b01be54750ca61b671942a))
## Refactor
* refactor: remove all unused functions, types and tests (#1037) (Darren)([7170808b](https://github.com/lacework/go-sdk/commit/7170808bd3c7ced956d4a4bd836a647321cfea43))
## Bug Fixes
* fix(cli): set correct SSH user for Amazon Linux 2 (#1046) (Nick Schmeller)([a92b37d1](https://github.com/lacework/go-sdk/commit/a92b37d13c9a90d65817d7a8d1ab1e29e693be3f))
## Other Changes
* chore(deps): bump github.com/aws/aws-sdk-go-v2/service/ec2instanceconnect (#1006) (dependabot[bot])([8cd69ed1](https://github.com/lacework/go-sdk/commit/8cd69ed190c908eff374ee39c31cf09e4dcdbea2))
* chore(deps): bump github.com/aws/aws-sdk-go-v2/service/ec2 (#1043) (dependabot[bot])([357ce2db](https://github.com/lacework/go-sdk/commit/357ce2dbbf910da1c3bf8bf19feb31d15cc70599))
* ci: fix release script; sort versions in find_latest_version (#1042) (Darren)([330a6a16](https://github.com/lacework/go-sdk/commit/330a6a16ad7d78718932e02f58f3bb4160f1f46a))
* ci: add logging line to output latest version found (#1040) (Darren)([38ec6ca1](https://github.com/lacework/go-sdk/commit/38ec6ca1a73016c5bd4244bce30742bd3a9bf396))
---
# v1.0.0

## Features

* feat: Add new cloud account type GCP Sidekick Integration (#1028) (
  ammarekbote)([fa428fb3](https://github.com/lacework/go-sdk/commit/fa428fb37add07373131930117c74def4b8129e8))
* feat(lwseverity): add valid severity use cases (#1024) (
  hazedav)([6b0b1330](https://github.com/lacework/go-sdk/commit/6b0b133005e387a32b6a705dd187cf233769853d))
* feat(cdk): bidirectional communication via gRPC (#1023) (Salim
  Afiune)([ff5654bd](https://github.com/lacework/go-sdk/commit/ff5654bdebdf92c0e95a4b9f8fde73a90cf8fae6))
* feat: lacework alert command (
  hazedav)([1059cc56](https://github.com/lacework/go-sdk/commit/1059cc566afbeccc3ca3df35d118b555385df0d6))

## Refactor

* refactor(v1): migrate cli to use v2 api (#1025) (
  Darren)([9f12eb22](https://github.com/lacework/go-sdk/commit/9f12eb226c47f03d5bafe53ece9a8558d77f9454))

## Bug Fixes

* fix: remove event command and api service (#1030) (
  hazedav)([87fd557d](https://github.com/lacework/go-sdk/commit/87fd557dfa158e141e005b9b44045a038d3b32e8))
* fix(api): container registry dockerhub 🐛  (#1032) (Salim
  Afiune)([0703eb93](https://github.com/lacework/go-sdk/commit/0703eb93b543684ecb631aa98813880427fb67ca))
* fix: compliance evaluation search cmd, searches greater than last 7 days  (#977) (
  Darren)([cf921309](https://github.com/lacework/go-sdk/commit/cf921309d6060621dc18231b43b585cc5e0f7bb4))

## Other Changes

* ci: version bump to v0.46.1-dev (
  Lacework)([a4a42bcb](https://github.com/lacework/go-sdk/commit/a4a42bcbc4240a3b63845fc8a69acc4d687f64ba))

---

# v0.46.0

## Features

* feat(cdk): improve honeyvent data for components (#1016) (Salim
  Afiune)([8204f095](https://github.com/lacework/go-sdk/commit/8204f095f6527143d115fc424b4b542f8f12c8c7))
* feat(cli)!: pass agent token as flag for `aws-install` (#1010) (Nick
  Schmeller)([6dbcb233](https://github.com/lacework/go-sdk/commit/6dbcb233fc4167ed9470fd886730e4728eb4e4b5))

## Refactor

* refactor(cli): fixup + test `aws-install` SSH user (#1009) (Nick
  Schmeller)([d2167710](https://github.com/lacework/go-sdk/commit/d21677104c2e3769418c3089908ed9e863e4b4f9))

## Bug Fixes

* fix(cdk): improve component arg parsing (#1014) (Sam
  Shen)([daa9c271](https://github.com/lacework/go-sdk/commit/daa9c27104c47d5bda34c88d6352eec15936d223))
* fix(cli): rootCmd to intelligently PreRun components (#1015) (
  hazedav)([dbfd4063](https://github.com/lacework/go-sdk/commit/dbfd40635e72fb771823f9fad79060337871e87f))
* fix(cli): minor fix related to refreshing tokens (#1013) (
  hazedav)([9f4c6193](https://github.com/lacework/go-sdk/commit/9f4c6193fd726d116643e89bad7cf7ef5a5ffd95))

## Other Changes

* ci: version bump to v0.45.2-dev (
  Lacework)([605a8573](https://github.com/lacework/go-sdk/commit/605a8573efb3610a8590daf252ac8a3ae8155163))

---

# v0.45.1

## Bug Fixes

* fix(cdk) pass valid LW_API_TOKEN to components (#1005) (Sam
  Shen)([136f123a](https://github.com/lacework/go-sdk/commit/136f123aa6a77553b7f96df33bb5e35372a12e38))
* fix(cli): do not `aws-install` on nil runners (#994) (Nick
  Schmeller)([9014bf23](https://github.com/lacework/go-sdk/commit/9014bf234c5df5370fd1e9d8287c20c7d43ce39b))

## Other Changes

* chore: bump lacework provider version (#1011) (jonathan
  stewart)([1448e08d](https://github.com/lacework/go-sdk/commit/1448e08d6a52ae546e778d7add066301774d3539))
* ci: version bump to v0.45.1-dev (
  Lacework)([e3e0b52c](https://github.com/lacework/go-sdk/commit/e3e0b52cbe8918db2a69907f3f5b5ae5cb7f56fa))

---

# v0.45.0

## Features

* feat: add JIRA bidirectional and grouping support (#940) (
  hazedav)([441fc00d](https://github.com/lacework/go-sdk/commit/441fc00d317f2593c70f11b52621ed87043da536))
* feat(generate): Add Terraform generation for AWS EKS Audit Log (#960) (
  Ross)([67dfb49b](https://github.com/lacework/go-sdk/commit/67dfb49bdfd1b16377c1feb7488eeb0ffd0e14ee))

## Bug Fixes

* fix(cdk): Avoid double component error messages (#1000) (Salim
  Afiune)([363c1d40](https://github.com/lacework/go-sdk/commit/363c1d40ebf7e0d2e675745e037c273dc646e317))
* fix(cdk): Pass all global flags via env variables (#993) (Salim
  Afiune)([1dc9fadd](https://github.com/lacework/go-sdk/commit/1dc9fadd9ef62376ccf5941f1ad632269845f104))
* fix(all) LimitNumScans was wrong type (#978) (Christien
  Rioux)([eab66dab](https://github.com/lacework/go-sdk/commit/eab66dab16e9930295a942a521b000148761470a))

## Other Changes

* style(cdk): detect lwcomponent error inside Execute() (#1001) (Salim
  Afiune)([5f88e520](https://github.com/lacework/go-sdk/commit/5f88e52067f2f4db6bbc3c7fb95da83e3fc253de))
* chore(deps): bump github.com/abiosoft/colima from 0.4.5 to 0.4.6 (#937) (
  dependabot[bot])([d8f67d1c](https://github.com/lacework/go-sdk/commit/d8f67d1cb66c9ad65cd567d9acf4e88388f35b9d))
* ci: version bump to v0.44.2-dev (
  Lacework)([10f4e321](https://github.com/lacework/go-sdk/commit/10f4e321abd57dc00c3e2b6b756ed98a017696b8))
* test(cdk): fix tests due to IaC component release (#998) (Salim
  Afiune)([4af2e7f1](https://github.com/lacework/go-sdk/commit/4af2e7f17e4c854c58b77f52765c60d97769664a))

---

# v0.44.1

## Bug Fixes

* fix(lint): redundant break statement (#991) (Salim
  Afiune)([aac27888](https://github.com/lacework/go-sdk/commit/aac278887b3b4c7c51071ec8372de6caa04c8082))
* fix(cdk): parse global flags when running components (#984) (Salim
  Afiune)([29e80a22](https://github.com/lacework/go-sdk/commit/29e80a22dc791db8f5cb07fda187d0e9c5f12ef2))

## Other Changes

* ci: version bump to v0.44.1-dev (
  Lacework)([e647232c](https://github.com/lacework/go-sdk/commit/e647232cdeea2e7b3db2fb887cfbf9409d2f75ae))

---

# v0.44.0

## Features

* feat: new v2 reports service (#943) (
  Darren)([fa766162](https://github.com/lacework/go-sdk/commit/fa766162f2562fc52f30d4424f8a6f6e4f1773a3))

## Refactor

* refactor(cli): migrate cmd compliance aws|gcp|azure get-report (#959) (
  Darren)([b613dfec](https://github.com/lacework/go-sdk/commit/b613dfec83be6cff23703dd3f54dd32c2444b421))

## Bug Fixes

* fix(cli): add query parameters to paging path (#980) (
  hazedav)([d4d1641a](https://github.com/lacework/go-sdk/commit/d4d1641ac483cd7c8b13172566f7260d360494a5))
* fix: fix gcp/azure --type flag  (#974) (
  Darren)([d67ea021](https://github.com/lacework/go-sdk/commit/d67ea0213f6ae315a36d2978a7ea5adb3c208613))

## Other Changes

* style(cli): rename `aws-install` files (#985) (Nick
  Schmeller)([c221cdab](https://github.com/lacework/go-sdk/commit/c221cdaba7cf6feb51b32ee30556f613badafa40))
* chore: Add .go-version file for goenv (#973) (
  Ross)([64ffaaea](https://github.com/lacework/go-sdk/commit/64ffaaeaa252e328128c4130695a2e1bcc3cb2ff))
* chore: bump golangci-lint to 1.50.0 (#969) (
  Ross)([524d7232](https://github.com/lacework/go-sdk/commit/524d7232ee7fa027a915388315e04cf6aa4748fd))
* chore: bump go jUnit to v2.0.0 (#972) (
  Ross)([b2cce81f](https://github.com/lacework/go-sdk/commit/b2cce81f3ccc34dba775f85918f24efaed09b12f))
* chore: bump goimports to 0.1.12 (#970) (
  Ross)([b297b0b1](https://github.com/lacework/go-sdk/commit/b297b0b181ed1b1c568966d3ecc7d5f71dde4410))
* chore: bump gotestsum 1.8.2 (#971) (
  Ross)([bed04947](https://github.com/lacework/go-sdk/commit/bed0494759a43ee1a37be1ade6e69181d22f4ca5))
* build: fix minor typo in pull request template (#968) (
  hazedav)([caeb5494](https://github.com/lacework/go-sdk/commit/caeb54942a4874af3f7c5c3fa4f51b7dd8a22f90))
* ci: version bump to v0.43.1-dev (
  Lacework)([5e9c0066](https://github.com/lacework/go-sdk/commit/5e9c0066010dc1f327fb8b69f1f5359a92f1b877))
* test(cli): enable compliance tests (#965) (Salim
  Afiune)([ee5389b6](https://github.com/lacework/go-sdk/commit/ee5389b6d757315e83bc701482f120df50ce4a18))

---

# v0.43.0

## Features

* feat(cli): parallelize installation with goroutines (#948) (Nick
  Schmeller)([6ee9d6d1](https://github.com/lacework/go-sdk/commit/6ee9d6d1c3f5eaeca36e2a412b544c23a221bb52))
* feat(cdk): dev-mode command (#957) (Salim
  Afiune)([faff0ae7](https://github.com/lacework/go-sdk/commit/faff0ae724206eff05d173ae3516fc5095330cc7))
* feat(api): New AwsSidekickOrg Cloud Account type (#951) (Teddy
  Reed)([4ab28fa0](https://github.com/lacework/go-sdk/commit/4ab28fa09cd0c2d480b0ab1e86237be3cae5cf36))
* feat(cli): find and fix AWS instances w/o agents (#934) (Nick
  Schmeller)([4e507c0b](https://github.com/lacework/go-sdk/commit/4e507c0b89edde2a8f04e03581a92aae9789bdc9))
* feat: GKE audit-log (#927) (jonathan
  stewart)([76e6cc88](https://github.com/lacework/go-sdk/commit/76e6cc88ed923dbabe29b5c2f68d7e382c93f0ed))
* feat: generate command rework (#928) (jonathan
  stewart)([56cbdcd1](https://github.com/lacework/go-sdk/commit/56cbdcd192a79148c24e735c2ae814fd1b70bb66))
* feat(cdk): component dev mode (#911) (Salim
  Afiune)([d477e688](https://github.com/lacework/go-sdk/commit/d477e688f3a360ea6e6cae6b08d8ac0a92df4a1c))
* feat: add new parameters to gcp iac generation (#918) (jonathan
  stewart)([e353487b](https://github.com/lacework/go-sdk/commit/e353487b3d58394fefcdf6c90aa05b22c919bafd))

## Bug Fixes

* fix(cli): token from flag into arg for `aws-install` (#952) (Nick
  Schmeller)([0633e5cb](https://github.com/lacework/go-sdk/commit/0633e5cb75a149214226c75634ffd9c8c1eaf44c))
* fix: gke integration test flags (#949) (jonathan
  stewart)([7cdd53e5](https://github.com/lacework/go-sdk/commit/7cdd53e5ff29fb10d962ba4f8c09cda3cbfa8599))
* fix(cli): update `aws-install` command in help msg (#947) (Nick
  Schmeller)([c1396e19](https://github.com/lacework/go-sdk/commit/c1396e197ef22c9c1c4c467a0706284563647e78))
* fix(lwcomponent): better JSON handling components (#946) (
  hazedav)([8edb92de](https://github.com/lacework/go-sdk/commit/8edb92dee96c769ca40246e0397d229de2e775a2))
* fix: aws terraform sns variables (#944) (jonathan
  stewart)([c8ead2c2](https://github.com/lacework/go-sdk/commit/c8ead2c2b7a6f74b88aa220a7fceaee92d6df1fb))
* fix: azure Lacework provider & profile flag (#938) (jonathan
  stewart)([75ed6ba4](https://github.com/lacework/go-sdk/commit/75ed6ba4c804104278a6819e13e5248a7c9f5e9c))
* fix: azure interactive caching (#939) (jonathan
  stewart)([cff15ed2](https://github.com/lacework/go-sdk/commit/cff15ed268990c53b5f292bacc995a7b62f08dd1))
* fix: generate TF overwrite (#935) (jonathan
  stewart)([b33ad901](https://github.com/lacework/go-sdk/commit/b33ad90181e69aad1e3a92da49eda7a2936a117f))
* fix(lwrunner): create ~/.ssh directory if not exist (#933) (Salim
  Afiune)([bfc90990](https://github.com/lacework/go-sdk/commit/bfc909902c92a4dc5490f03d69357fa3ab2e6c84))
* fix: gcr limitByRep json field name (#924) (
  Darren)([0167eabe](https://github.com/lacework/go-sdk/commit/0167eabe68c13172db82d4d7aa33cbe5535c32b2))
* fix: gcr limit by repo json field name (#923) (
  Darren)([6f13c4c0](https://github.com/lacework/go-sdk/commit/6f13c4c07e8419dcfe97309cfac42126ecdcdddc))
* fix: limit_by_labels json field name (#922) (
  Darren)([857ea019](https://github.com/lacework/go-sdk/commit/857ea019c8ef55ff92626da3fb90a7b8ab18eb62))
* fix: omitempty fields ctr registry v2 (#921) (
  Darren)([14a8a225](https://github.com/lacework/go-sdk/commit/14a8a2250e18b3f9f679ad6ec63057d140656a52))
* fix(lql): fix handling of large numbers in json (#919) (
  hazedav)([6c934ad3](https://github.com/lacework/go-sdk/commit/6c934ad353119d31c30e5792a1db017673aac9d9))
* fix: DockerhubV2Data fix fields (#920) (
  Darren)([cbf3e395](https://github.com/lacework/go-sdk/commit/cbf3e3959b56799a832daaf21373c7b7a7dc4c9f))
* fix(ALLY-1188): Enable the disable of enable_ebla (#916) (jonathan
  stewart)([ba3c3562](https://github.com/lacework/go-sdk/commit/ba3c356216a3970743ab91d3fb9bb4fed5a1461a))

## Other Changes

* chore(deps): bump github.com/honeycombio/libhoney-go (#931) (
  dependabot[bot])([e5624291](https://github.com/lacework/go-sdk/commit/e56242913e0b7973957b2de2f11668eddbe06720))
* chore(deps): bump github.com/abiosoft/colima from 0.4.4 to 0.4.5 (#930) (
  dependabot[bot])([3aa5ff93](https://github.com/lacework/go-sdk/commit/3aa5ff9344493cbe936a49ea7277a59c61e67e36))
* chore(deps): bump github.com/hashicorp/hcl/v2 from 2.13.0 to 2.14.1 (#925) (
  dependabot[bot])([e872b5a4](https://github.com/lacework/go-sdk/commit/e872b5a4058380c61be8b10a47c11a11ad2621b2))
* chore(deps): bump go.uber.org/zap from 1.21.0 to 1.23.0 (#901) (
  dependabot[bot])([593cdbaa](https://github.com/lacework/go-sdk/commit/593cdbaaef8cef5dbe21bfc77c2ead7ad0e7cea8))
* chore(deps): bump github.com/hashicorp/terraform-exec (#908) (
  dependabot[bot])([976c486c](https://github.com/lacework/go-sdk/commit/976c486c318a6df7ab14014ef9c2880062687da1))
* chore(deps): bump github.com/spf13/viper from 1.12.0 to 1.13.0 (#910) (
  dependabot[bot])([658b600a](https://github.com/lacework/go-sdk/commit/658b600ae729cf039eb6de1d0b74bd07b8c25cea))
* build: fix install-cli Makefile target to account for ARM builds (#950) (
  0xdeadbeefJERKY)([4fb46f23](https://github.com/lacework/go-sdk/commit/4fb46f23944bb176428116b297cf1ae56dd3e32e))
* ci: version bump to v0.42.1-dev (
  Lacework)([3044a87f](https://github.com/lacework/go-sdk/commit/3044a87f12b54618fef9c3b6308dc2c82fccba93))
* test(cli): disable compliance tests (#964) (Salim
  Afiune)([6989e67e](https://github.com/lacework/go-sdk/commit/6989e67e077dede960e002d923b11ee4001a1d19))
* test: update error msg test (#942) (
  Darren)([ad04af52](https://github.com/lacework/go-sdk/commit/ad04af525202fc883adee52ea706ccff8e38a36b))
* test(win): disable failing integration test (#932) (Salim
  Afiune)([7e2adfc0](https://github.com/lacework/go-sdk/commit/7e2adfc02e9cb29a41824cd6dc9ca31c7d8a49ff))

---

# v0.42.0

## Features

* feat(cli): alert-channel cmd & container-registry cmd (#904) (
  Darren)([bdb6317](https://github.com/lacework/go-sdk/commit/bdb63178d5b421bea431c12839974b358d254fc2))

## Refactor

* refactor(chocolatey): address moderation requested changes (#906) (
  Darren)([e33ecfc](https://github.com/lacework/go-sdk/commit/e33ecfc5d721960a2bf9d9ad3e01cc2fbcfa84e3))

## Other Changes

* chore(deps): bump github.com/kyokomi/emoji/v2 from 2.2.9 to 2.2.10 (#884) (
  dependabot[bot])([bee54a0](https://github.com/lacework/go-sdk/commit/bee54a0b0301b836f04210133c9f8ec8f83887c9))
* chore(deps): bump github.com/zclconf/go-cty from 1.10.0 to 1.11.0 (#894) (
  dependabot[bot])([0436d26](https://github.com/lacework/go-sdk/commit/0436d26123b40a5c4615188bec2fbc024b400b92))
* ci: version bump to v0.41.2-dev (
  Lacework)([0e6cd1e](https://github.com/lacework/go-sdk/commit/0e6cd1e8adb1b5da2be63785e7166e2d8596b583))

---

# v0.41.1

## Other Changes

* build(cli): detect when signing process has completed (#902) (Salim
  Afiune)([05bca71](https://github.com/lacework/go-sdk/commit/05bca713ac5eaddbc080b5eac2ef3090c65f48e5))
* ci: version bump to v0.41.1-dev (
  Lacework)([d176ac5](https://github.com/lacework/go-sdk/commit/d176ac5c30a0963752e993816aa2be5c8418f405))

---

# v0.41.0

## Features

* feat(cli): aws compliance search command (#893) (
  Darren)([5c665cb](https://github.com/lacework/go-sdk/commit/5c665cba2bc0d7f89c669d6096d3f60932eb8d2d))
* feat(cli): Addition of extra AWS terraform options (#890) (
  djmctavish)([de34097](https://github.com/lacework/go-sdk/commit/de34097c80f239e51335a33655ef89f35f5b7c06))
* feat(api): add v2 AwsCfg (#889) (
  Darren)([53bd1ca](https://github.com/lacework/go-sdk/commit/53bd1ca21bd5d6e4893397ffad14403154f184d9))
* feat: DataExportRule service (#887) (
  Darren)([9c6ca08](https://github.com/lacework/go-sdk/commit/9c6ca082f9ae0131bb2c3b0d65aec9fb51b4e6b6))
* feat: add credentials to sidekick data (#886) (
  Darren)([3311753](https://github.com/lacework/go-sdk/commit/331175320a594ef4a4646cf2ebdd5b8c06eee407))

## Refactor

* refactor(cli): migrate fetch accountIds in policy exception to API v2 (#896) (
  Darren)([510a41f](https://github.com/lacework/go-sdk/commit/510a41f0bdb6ba9e8d3df24db55c691ce0349fa0))

## Bug Fixes

* fix: omitempty profileVersions field (#897) (
  Darren)([2170edb](https://github.com/lacework/go-sdk/commit/2170edb47e5353597149c6208fedbe7604fd7a4f))
* fix(cli): account flag domain fix (#895) (jonathan
  stewart)([ed892dc](https://github.com/lacework/go-sdk/commit/ed892dc23bbef80eda97985ad1fc9c385ace4260))
* fix: add description field to dataexportrules (#891) (
  Darren)([0531907](https://github.com/lacework/go-sdk/commit/05319072e362cfc258b9555c52953138a62d9953))
* fix(lwupdater): avoid comparing dev versions (#881) (Salim
  Afiune)([a8f9ee4](https://github.com/lacework/go-sdk/commit/a8f9ee444c9899768b57fe510308878e2d6aa13c))
* fix: alert profile filter (#883) (
  Darren)([91e979a](https://github.com/lacework/go-sdk/commit/91e979a12156c5fced9fa9cf7490d6499ff87bb9))

## Documentation Updates

* docs(cli): display IaC default output location (#898) (
  sthetland)([139da68](https://github.com/lacework/go-sdk/commit/139da68fb6ec1033b1ecefe8a8295e08e955404e))

## Other Changes

* build(cli): sign Lacework CLI artifacts (#899) (Salim
  Afiune)([dcc4605](https://github.com/lacework/go-sdk/commit/dcc460567947f316d7eeb7dcda2406813cf83dc3))
* ci: fix policy test (#888) (Salim
  Afiune)([294f99b](https://github.com/lacework/go-sdk/commit/294f99b9e9767aea33450b561a46e78231f1bf4a))
* ci: version bump to v0.40.1-dev (
  Lacework)([599d7ff](https://github.com/lacework/go-sdk/commit/599d7ff61c0a21c961862d633d11f9c450cfc471))

---

# v0.40.0

## Features

* feat(cli): alert-profile update cmd (#879) (
  Darren)([d8024e5](https://github.com/lacework/go-sdk/commit/d8024e5aad8ea7e240d2be2a65680351678034ce))
* feat(cdk): leave breadcrumbs after installation and update (#876) (Salim
  Afiune)([d16cc15](https://github.com/lacework/go-sdk/commit/d16cc153cc6483dd20e05bfd3035f5a0445ef17a))
* feat(lql): add limit option (
  hazedav)([d598def](https://github.com/lacework/go-sdk/commit/d598deff17ae4b644cc655c084e96fa6b5d70ad4))
* feat: add cloud account integration aws agentless scanning (#870) (
  Darren)([8807a97](https://github.com/lacework/go-sdk/commit/8807a97ac04a9d2637947a6b74b5f6bfd148e91e))

## Refactor

* refactor: use v2/AgentInfo for agents list command (#866) (
  Darren)([06531ed](https://github.com/lacework/go-sdk/commit/06531ed7ebcc8fba41d95dd998aa2e8d0125aa7c))

## Bug Fixes

* fix(lql): add line break when no preview found (#880) (
  hazedav)([83a1216](https://github.com/lacework/go-sdk/commit/83a1216a20162af5fdc4c30c5f0572d65a773e64))
* fix(cdk): JSON output for 'component list' command (#878) (Salim
  Afiune)([a1dc162](https://github.com/lacework/go-sdk/commit/a1dc162441b5ac849010dc87b5a09fbce3f13d16))
* fix(api): detect FQDN as account in api.NewClient() (#872) (Salim
  Afiune)([6e14d57](https://github.com/lacework/go-sdk/commit/6e14d577fc28f6e5b4eabbda0fe5f08d080e4473))
* fix: vulnerabilty_exception fixable parameter omit when unset (#863) (
  Darren)([f517e0e](https://github.com/lacework/go-sdk/commit/f517e0e48fe6cf0ac8f05d59d65c25dc4dbb7ddd))
* fix(install.ps1): modify User $Path correctly (#869) (Salim
  Afiune)([73c69f7](https://github.com/lacework/go-sdk/commit/73c69f74a4215799f9826d14badaa5214150ba97))

## Documentation Updates

* docs(cli): update help message from LPP commands (#873) (
  lacework-aaronscheiber)([810d005](https://github.com/lacework/go-sdk/commit/810d0056decf864ab8cd522c2f571475b825a76b))
* docs(cli): updating help for query run (#864) (
  hazedav)([8953d19](https://github.com/lacework/go-sdk/commit/8953d19f8538b6247f0018d66777ed891eb50c51))

## Other Changes

* chore(deps): bump github.com/honeycombio/libhoney-go (#867) (
  dependabot[bot])([e77d133](https://github.com/lacework/go-sdk/commit/e77d133e3d7564de0fae979b08d8fbfe05221cf4))
* chore(deps): bump github.com/BurntSushi/toml from 1.1.0 to 1.2.0 (#871) (
  dependabot[bot])([fa69c24](https://github.com/lacework/go-sdk/commit/fa69c24b079e09ac4012a8beca22947cc01ce686))
* chore(lql): fully remove evaluatorId (#865) (
  hazedav)([1112f55](https://github.com/lacework/go-sdk/commit/1112f55c47577591d94cd7719ab66c07f51f21e4))
* ci: version bump to v0.39.1-dev (
  Lacework)([23827ee](https://github.com/lacework/go-sdk/commit/23827ee45a92c9c23784af1e61967eff6ef9dae9))

---

# v0.39.0

## Features

* feat(cli): colorize CDK commands (#859) (Salim
  Afiune)([41e0785](https://github.com/lacework/go-sdk/commit/41e07854448b5053312bc57791c3b20a7aa4de60))
* feat(cli): fetch existing policy for cli update (#841) (
  Darren)([8bc93ba](https://github.com/lacework/go-sdk/commit/8bc93baa67e3f8f063afc9b4f05199801cd8c170))
* feat(cli): new policy-exceptions command (#844) (
  Darren)([7e894c9](https://github.com/lacework/go-sdk/commit/7e894c9d5d2f3d550f3d117db684a6116e851ecd))
* feat(cli): include status in vuln ctr scan assesment output (#853) (
  Darren)([5ff586c](https://github.com/lacework/go-sdk/commit/5ff586cb3a12c6f60fa79164b6a3d3e74461d6f4))
* feat(cli): display CDK version (Salim Afiune
  Maya)([d21bbef](https://github.com/lacework/go-sdk/commit/d21bbefeb4da27cc3036d1981c0bec44ec88d97a))
* feat(cdk): life cycle event 'cdk-init' (Salim Afiune
  Maya)([12eee01](https://github.com/lacework/go-sdk/commit/12eee017c12145c6bc315050803e9f15f04a651e))
* feat: support inline scanner integrations (Salim Afiune
  Maya)([99aacd8](https://github.com/lacework/go-sdk/commit/99aacd80e2a451f2286f49d6cd2079adf3a87105))

## Bug Fixes

* fix: omit vulnerability_exception fixable parameter when unset (#857) (
  Darren)([e671857](https://github.com/lacework/go-sdk/commit/e6718576f977a79bdef231c164c60d0717df290d))
* fix(cli): use proper env var name LW_API_TOKEN (#856) (Salim
  Afiune)([22efef4](https://github.com/lacework/go-sdk/commit/22efef47783986824d3c404cee86484660ca7ff4))
* fix(cli): install.ps1 avoid Windows agent collision (#854) (Salim
  Afiune)([1d52d27](https://github.com/lacework/go-sdk/commit/1d52d27057cf3e4f99896509db26a6cf2704ecc6))

## Documentation Updates

* docs(api): fix path to examples in readme (#858) (Yasmine
  Hal)([a8a5feb](https://github.com/lacework/go-sdk/commit/a8a5febb8a8cdd373e4e07cde64041f4d0abbd4f))

## Other Changes

* chore(cli): update content-library component name  (#855) (
  hazedav)([2794070](https://github.com/lacework/go-sdk/commit/27940704c140ca229235f8fc1ac8b1c81f1ad269))
* ci: version bump to v0.38.1-dev (
  Lacework)([196406a](https://github.com/lacework/go-sdk/commit/196406aedce90674661551ac146f6cdec2d8d634))

---

# v0.38.0

## Features

* feat: chocolatey package scripts (#832) (
  Darren)([22c1005](https://github.com/lacework/go-sdk/commit/22c1005f64c6a7cb81957eeb06961b5de2856575))

## Refactor

* refactor: policy exception constraint type (#842) (
  Darren)([88bfdb5](https://github.com/lacework/go-sdk/commit/88bfdb58eb73abf00bacc797075aa4350cda0683))

## Bug Fixes

* fix(cli): subtract 1 minute from machine details end time (#848) (
  Darren)([c53388f](https://github.com/lacework/go-sdk/commit/c53388f794ae9b3b9690728fed76c75c9749862a))
* fix(cli): Fix icon, desc and summary of .nuspec file (#843) (
  Darren)([0c5c456](https://github.com/lacework/go-sdk/commit/0c5c45689742903c47632b167d44388ebff9af9d))

## Other Changes

* chore(deps): bump github.com/hashicorp/hcl/v2 from 2.12.0 to 2.13.0 (#840) (
  dependabot[bot])([77d78f0](https://github.com/lacework/go-sdk/commit/77d78f0c996aedaf595d2f9bd7e90cc398ad0256))
* chore(deps): bump github.com/hashicorp/terraform-exec (#845) (
  dependabot[bot])([9db51aa](https://github.com/lacework/go-sdk/commit/9db51aa8aff39b7f2090d638760c4791780d0cc0))
* chore(deps): bump github.com/abiosoft/colima from 0.4.2 to 0.4.4 (#849) (
  dependabot[bot])([0b1b60e](https://github.com/lacework/go-sdk/commit/0b1b60e18d9164a2ba3b575e72610a434b0380c1))
* chore(deps): bump github.com/stretchr/testify from 1.7.5 to 1.8.0 (#846) (
  dependabot[bot])([ba2cdbf](https://github.com/lacework/go-sdk/commit/ba2cdbf0a4316bf33c192e5fd7981f389c8ce0b9))
* chore(deps): bump github.com/hashicorp/go-version from 1.5.0 to 1.6.0 (#847) (
  dependabot[bot])([b8b232e](https://github.com/lacework/go-sdk/commit/b8b232e999086b23b79b9267e9d74099200d5a37))
* chore(deps): bump github.com/hashicorp/terraform-exec (#837) (
  dependabot[bot])([c7b06b6](https://github.com/lacework/go-sdk/commit/c7b06b6f42546e4c9455886779e08e750dabe292))
* chore(deps): bump github.com/hashicorp/hc-install from 0.3.2 to 0.4.0 (#836) (
  dependabot[bot])([2e388bd](https://github.com/lacework/go-sdk/commit/2e388bda6f5666c752db824335143df0ec8bd4ac))
* chore(deps): bump github.com/stretchr/testify from 1.7.3 to 1.7.5 (#838) (
  dependabot[bot])([d3d0ef7](https://github.com/lacework/go-sdk/commit/d3d0ef772df59d52af1034ab83f70749a4fa8410))
* ci: version bump to v0.37.1-dev (
  Lacework)([841d251](https://github.com/lacework/go-sdk/commit/841d251b3def5e8c2cb4bd0f51bc49b4602b5df6))

---

# v0.37.0

## Features

* feat(api): add policy exceptions service (#826) (
  Darren)([fe5b9bf](https://github.com/lacework/go-sdk/commit/fe5b9bf508800d82377e379805958925d10e8a08))
* feat(cdk): pass env variable LW_COMPONENT_NAME (#829) (Salim
  Afiune)([4b7451a](https://github.com/lacework/go-sdk/commit/4b7451a3823daab4fa90a00bda0e6ea6666b5ebd))

## Refactor

* refactor: move vulnerabilityExceptionFixableEnabled to go-sdk (#820) (
  Darren)([941570e](https://github.com/lacework/go-sdk/commit/941570e9edc44277da081a5b6db9c5e36a3036d2))

## Bug Fixes

* fix(GcpGkeAudit): Casing issue with PrivateKeyId & PrivateKeyId (#834) (
  Ross)([b7b328b](https://github.com/lacework/go-sdk/commit/b7b328b4ca0557ec569460d37895875288ffc83a))
* fix: misspelled namespace field (#831) (Salim
  Afiune)([6312843](https://github.com/lacework/go-sdk/commit/6312843f7a1a2dfbe68af3eba552abf8338bc481))
* fix(cli): Make the html table scrollable (#825) (John
  John)([85d5214](https://github.com/lacework/go-sdk/commit/85d5214a6d09ded8ace5056145b4f45083e97789))
* fix(cdk): run component dynamically (#823) (Salim
  Afiune)([bb1ce0d](https://github.com/lacework/go-sdk/commit/bb1ce0dc0b9c74bff80385735028e71941e5aa66))

## Other Changes

* chore(deps): bump github.com/hashicorp/terraform-json (#808) (
  dependabot[bot])([ac2b094](https://github.com/lacework/go-sdk/commit/ac2b09446833a6cc1abb0bf74df547889f7107f1))
* chore(deps): bump github.com/spf13/viper from 1.10.1 to 1.12.0 (#802) (
  dependabot[bot])([000024a](https://github.com/lacework/go-sdk/commit/000024a90f37963dce2d5a80139f2ca5e8ab977e))
* chore(deps): bump github.com/stretchr/testify from 1.7.2 to 1.7.3 (#828) (
  dependabot[bot])([fb93f55](https://github.com/lacework/go-sdk/commit/fb93f558c6c9b07be8001d2f8335b1bddc46ea3a))
* ci: version bump to v0.36.1-dev (
  Lacework)([cdd6d7b](https://github.com/lacework/go-sdk/commit/cdd6d7b0b1596a50261449742c9e134916b9a535))

---

# v0.36.0

## Features

* feat(api): Implement GcpGkeAudit CloudAccount interface (#821) (
  Ross)([fd01c7f](https://github.com/lacework/go-sdk/commit/fd01c7f8c92eeebe7465561e7a978b39f26d7854))
* feat(cli): lacework query fail_on_count (#415) (
  hazedav)([0b1de6a](https://github.com/lacework/go-sdk/commit/0b1de6acd3873c0b754302292db396eee12d6b08))
* feat: Lacework CDK (Cloud Development Kit) (#544) (Salim
  Afiune)([b0ba78c](https://github.com/lacework/go-sdk/commit/b0ba78c742b5e09381a81fe7446c1b5fdb3ef655))

## Other Changes

* style(cli): use switch instead of if statements (#822) (Salim
  Afiune)([e40322c](https://github.com/lacework/go-sdk/commit/e40322c37dd9175de1d796582604b554ecdcc69c))
* chore(deps): bump github.com/mitchellh/mapstructure from 1.4.3 to 1.5.0 (#779) (
  dependabot[bot])([af66b8e](https://github.com/lacework/go-sdk/commit/af66b8ec01a31cee0b69190978bf51285d966dd4))
* chore(deps): bump github.com/imdario/mergo from 0.3.12 to 0.3.13 (#819) (
  dependabot[bot])([47007f0](https://github.com/lacework/go-sdk/commit/47007f0582835ddaf7291d71773dd0e9439f6b41))
* ci: version bump to v0.35.2-dev (
  Lacework)([1ff58d2](https://github.com/lacework/go-sdk/commit/1ff58d2803d7e4cd509310667dd5766d2766241c))

---

# v0.35.1

## Bug Fixes

* fix: expiration time format (#811) (
  Darren)([521b7ac](https://github.com/lacework/go-sdk/commit/521b7ac7e3836669db704566b450e86800818e9e))
* fix(cli): compliance gcp enable command (#813) (
  Darren)([0a931a8](https://github.com/lacework/go-sdk/commit/0a931a8a70eff44aa228afb8495777a31f984b9b))

## Other Changes

* style(cli): add missing carrier return (#816) (Salim
  Afiune)([ca37520](https://github.com/lacework/go-sdk/commit/ca37520f3f05c21f06ea49221f53895ecf3029ed))
* chore(deps): bump github.com/hashicorp/hcl/v2 from 2.11.1 to 2.12.0 (#780) (
  dependabot[bot])([9471ffd](https://github.com/lacework/go-sdk/commit/9471ffd3657e82ceb4958143220668ddb2e3a30d))
* chore(deps): bump github.com/stretchr/testify from 1.7.1 to 1.7.2 (#806) (
  dependabot[bot])([d936875](https://github.com/lacework/go-sdk/commit/d9368757d23cc2eb83dda67a81b5260852d16da9))
* chore(lwgenerate): bump GCP config module to ~> 2.0 (#814) (Salim
  Afiune)([913740b](https://github.com/lacework/go-sdk/commit/913740b86e5fd8dcd3c69eceb944c7b86273535f))
* ci: version bump to v0.35.1-dev (
  Lacework)([7efab5f](https://github.com/lacework/go-sdk/commit/7efab5f1d7a7c2595c7408a5b39d59ecd8fc7462))

---

# v0.35.0

## Features

* feat: show query start time in progress spinner (#798) (
  Darren)([2179616](https://github.com/lacework/go-sdk/commit/21796160c18f1c7acbf6ece7f2fbc8c3d5526ec0))
* feat(api): Add S3BucketArn field to allow BYOB for AWS EKS Audit Log (#793) (
  Ross)([93dbb2d](https://github.com/lacework/go-sdk/commit/93dbb2d74f700fc7c536566c7130d8096ea3febc))
* feat(cli): alert profiles command (#792) (
  Darren)([c0bc9a3](https://github.com/lacework/go-sdk/commit/c0bc9a3a694a7ae28f65c10f264dfff6e0042102))

## Refactor

* refactor(cli): query run progress message (#804) (
  hazedav)([f64d461](https://github.com/lacework/go-sdk/commit/f64d461b0dc674473122ff29035323b9d5b29dbe))
* refactor(cli): updated terminology for fail_on (#803) (
  hazedav)([f5b85f5](https://github.com/lacework/go-sdk/commit/f5b85f587f19ad656aef1095096d613cd5a328d8))

## Bug Fixes

* fix: allow resource scope to be optional (#807) (Salim
  Afiune)([6e8e55c](https://github.com/lacework/go-sdk/commit/6e8e55cdd0d3e54b3a8b9d450b36f84e9ceff614))
* fix(cli): Add MachineDetail tags (
  hazedav)([8cd1ea0](https://github.com/lacework/go-sdk/commit/8cd1ea0aebdee72d87d7e8faf2a7bae2cd87fe4c))
* fix(cli): panic while detecting ESM on Ubuntu sys (#799) (Salim
  Afiune)([dce2d18](https://github.com/lacework/go-sdk/commit/dce2d18ab0e88e1f106d69b352ac041ee7810221))

## Other Changes

* chore(deps): bump github.com/hashicorp/go-version from 1.4.0 to 1.5.0 (#801) (
  dependabot[bot])([7665edd](https://github.com/lacework/go-sdk/commit/7665edd61c4365d0fc12ce363ba4ff94421f41f8))
* chore(deps): bump github.com/hashicorp/terraform-exec (#773) (
  dependabot[bot])([a7915cc](https://github.com/lacework/go-sdk/commit/a7915cc74e2883fa048dff47e7b100a1e976358d))
* chore(deps): bump github.com/abiosoft/colima from 0.4.1 to 0.4.2 (#795) (
  dependabot[bot])([6ea8ef3](https://github.com/lacework/go-sdk/commit/6ea8ef30b61474daba5b44fadd0efd06b6a6ccfe))
* chore(deps): bump github.com/abiosoft/colima from 0.3.4 to 0.4.1 (#791) (
  dependabot[bot])([cf19e21](https://github.com/lacework/go-sdk/commit/cf19e21eebe9aa1d88dc179a221ee5e427b58e1b))
* ci: version bump to v0.34.1-dev (
  Lacework)([8811c84](https://github.com/lacework/go-sdk/commit/8811c84d050aad00f250b35b47ee8422e1907156))
* test: fix pipeline to check for low vulns (#805) (Salim
  Afiune)([233c24f](https://github.com/lacework/go-sdk/commit/233c24f372a913f61d64c509183c77323e376982))
* test: disable file test since CI runs as root (#800) (Salim
  Afiune)([5b271c6](https://github.com/lacework/go-sdk/commit/5b271c671d058ebe9fc898ec785c5c750e167d65))

---

# v0.34.0

## Features

* feat: Add support for APIv2 CloudAccounts List (#784) (
  Ross)([563d1da](https://github.com/lacework/go-sdk/commit/563d1da63b7eae5f7e7691117787570d2d0e8cfc))
* feat(cli): add gcp compliance report disable warning (#787) (
  Ross)([0cc3734](https://github.com/lacework/go-sdk/commit/0cc373442eff0bc970b10d5eabfde0296ba699a1))
* feat(cli): add azure compliance report disable warning (#786) (
  Ross)([a028ecb](https://github.com/lacework/go-sdk/commit/a028ecb3b6a5a5e61f82464c9ef83226dd393227))

## Bug Fixes

* fix(test): Fix TestCloudAccountsListByType unit test failure (#788) (
  Ross)([48fb4c7](https://github.com/lacework/go-sdk/commit/48fb4c78d6f2d47c353341c0d47978fcefa25651))
* fix(api): add newResponse for function SearchAllPages (#770) (
  aircraft-cerier)([474a163](https://github.com/lacework/go-sdk/commit/474a16396b30d70b934a07839891436c3cef6643))
* fix(cli): Remove Additional Cloudtrail/AuditLog options when Cloudtrail/AuditLog Integration is set to No (#782) (
  Ross)([c9d7f48](https://github.com/lacework/go-sdk/commit/c9d7f48572e028be691e7b984b52208681cfad5a))

## Other Changes

* ci: version bump to v0.33.1-dev (
  Lacework)([3b3aacf](https://github.com/lacework/go-sdk/commit/3b3aacf9e8075d39532980b063e2083c9e0dc88a))

---

# v0.33.0

## Features

* feat(cli): add compliance report disable warning (#778) (
  hazedav)([ae2ffb7](https://github.com/lacework/go-sdk/commit/ae2ffb7cb065b9be2c9dee33985deed0ac7f7e20))

## Bug Fixes

* fix(cli): Remove aws profile "default" (#777) (
  Ross)([e47ea55](https://github.com/lacework/go-sdk/commit/e47ea553fce260597d8262288c18d7d913422473))

## Other Changes

* ci: version bump to v0.32.1-dev (
  Lacework)([209ec0a](https://github.com/lacework/go-sdk/commit/209ec0ac86cac33b8ecab6998da9b09e10b96b12))

---

# v0.32.0

## Features

* feat(cli): output query & policy in YAML 📝 (#775) (Salim
  Afiune)([9173efe](https://github.com/lacework/go-sdk/commit/9173efef1603bb204ac998c7d8ed4a65a10bcdc8))
* feat(api): AlertProfiles v2 service (#765) (
  Darren)([683ceca](https://github.com/lacework/go-sdk/commit/683ceca0c99971c9cde7371e3fc5cbd24adb9d33))
* feat(cli): ubuntu ESM support (#768) (Salim
  Afiune)([7527e83](https://github.com/lacework/go-sdk/commit/7527e834cb15f246839fd984f73e4885949a9fea))

## Bug Fixes

* fix(cli): only display datasource relationships when avail (
  hazedav)([c9589ed](https://github.com/lacework/go-sdk/commit/c9589ed5abbac73511f829fb9b69571cc795ed60))

## Documentation Updates

* docs: enable policy enable/disable help text. (#766) (
  Darren)([768cfbf](https://github.com/lacework/go-sdk/commit/768cfbf4125f1f79b9801674dc7f9a4ca8fe0365))

## Other Changes

* chore(deps): bump github.com/BurntSushi/toml from 0.4.1 to 1.1.0 (#763) (
  dependabot[bot])([a3b6c5b](https://github.com/lacework/go-sdk/commit/a3b6c5b34bdd7272a2fc93267536865b37339330))
* chore(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1 (#740) (
  dependabot[bot])([db6541c](https://github.com/lacework/go-sdk/commit/db6541c0e32ae343629e01644971fdb2620d2404))
* chore(deps): bump github.com/abiosoft/colima from 0.3.2 to 0.3.4 (#723) (
  dependabot[bot])([b9bef2a](https://github.com/lacework/go-sdk/commit/b9bef2a4c57abfdc89ed5660e0f9250e755ef149))
* ci: version bump to v0.31.1-dev (
  Lacework)([7868fe7](https://github.com/lacework/go-sdk/commit/7868fe791cdf4a8cf2e574776798175f573a2939))
* test(lql): fix case issue with host tests (
  hazedav)([30d7d7b](https://github.com/lacework/go-sdk/commit/30d7d7b77fb477a81361b0eca4d18eaed637f596))
* test: rename LW_CLI_ -> CLI_ (#769) (
  Darren)([5cfedda](https://github.com/lacework/go-sdk/commit/5cfedda359d2cba7d17a25f1eecc55207989ab67))

---

# v0.31.0

## Features

* feat(cli): disable/enable policies by tag (#758) (
  Darren)([3a45a01](https://github.com/lacework/go-sdk/commit/3a45a0186ca908920f2383c2eec2c7c7db6e17bf))
* feat(cli): default template when creating LQL queries (#761) (Salim
  Afiune)([4bb9c19](https://github.com/lacework/go-sdk/commit/4bb9c1930b64e3e080ee44fc71028276bb44d8d8))
* feat(cli): disable/enable all <report-type> compliance command (#755) (
  Darren)([e694304](https://github.com/lacework/go-sdk/commit/e69430484c93245919b2e48e264dc3c24d6a7297))
* feat(cli): adds lql datasource relationships (#759) (
  hazedav)([4b1fc27](https://github.com/lacework/go-sdk/commit/4b1fc271db940d5272cd652f5b65d06e061af874))
* feat(api): Recommendations v1 Service (#751) (
  Darren)([253cc38](https://github.com/lacework/go-sdk/commit/253cc38420ec9a6423aec2dc57958255eb1e5cf9))

## Refactor

* refactor(cli): inefficient policy sort (
  hazedav)([9a35ff6](https://github.com/lacework/go-sdk/commit/9a35ff665acd795f7a713b470d6603771f060beb))

## Bug Fixes

* fix(cli): command arguments (#760) (Salim
  Afiune)([574dfc1](https://github.com/lacework/go-sdk/commit/574dfc1dbb8eb066ca3e8668ed72be58731b0ff3))
* fix: add missing gcp compliance report types to get-report command (#756) (
  Darren)([c05a07f](https://github.com/lacework/go-sdk/commit/c05a07ff0664164d8e4d53009c4e2af6e2451c27))

## Other Changes

* chore(deps): bump github.com/kyokomi/emoji/v2 from 2.2.8 to 2.2.9 (#724) (
  dependabot[bot])([717a725](https://github.com/lacework/go-sdk/commit/717a725edfdb205701456bd6f960b1b0a3788988))
* chore(deps): bump github.com/spf13/cobra from 1.3.0 to 1.4.0 (#733) (
  dependabot[bot])([d1531ee](https://github.com/lacework/go-sdk/commit/d1531eef51e37966bc39e002d0589d400c548990))
* ci: add make cmd for output go tests in junit format (#754) (
  Darren)([f3d0324](https://github.com/lacework/go-sdk/commit/f3d032431bea1770ab8b731a5aefef23ab1c4594))
* ci: version bump to v0.30.1-dev (
  Lacework)([5cfe47a](https://github.com/lacework/go-sdk/commit/5cfe47a62c3c650cb35e879528614aa8155f9ac8))

---

# v0.30.0

## Features

* feat(lql): more suitable default for lql run (
  hazedav)([73a7132](https://github.com/lacework/go-sdk/commit/73a7132765125f0d2370aef36fa4c8039c1c0cdd))
* feat(cli/cmd): CLI survey code for azure - ALLY-394 (#720) (
  djmctavish)([73f8b0c](https://github.com/lacework/go-sdk/commit/73f8b0c83eaf6c67728f871cc3f4e4544c851ba7))

## Bug Fixes

* fix: fix fail_on_severity when no vulns are found (#750) (
  Darren)([479869f](https://github.com/lacework/go-sdk/commit/479869f95ae1c6d15b9108a59fcf36d77d448c4b))

## Other Changes

* chore: deprecate CLI Wiki (#739) (Salim
  Afiune)([eeb6478](https://github.com/lacework/go-sdk/commit/eeb64783f2224c397e61d5d3fe4fd93523e1fc1f))
* ci: version bump to v0.29.1-dev (
  Lacework)([fc74684](https://github.com/lacework/go-sdk/commit/fc74684e5143e7a01158f876d28e8a9773c8d2e6))

---

# v0.29.0

## Features

* feat(api): Implement AwsEksAudit CloudAccount interface (#742) (
  Ross)([879048c](https://github.com/lacework/go-sdk/commit/879048c11430e4606b270be61329a7e71d3138fb))

## Other Changes

* chore: update go version to 1.18 (#735) (
  Darren)([ec936cb](https://github.com/lacework/go-sdk/commit/ec936cb983f4f8dd69a24479252608468cac6a88))
* ci: version bump to v0.28.1-dev (
  Lacework)([6d7faae](https://github.com/lacework/go-sdk/commit/6d7faae6b398fd14b7da3b6fdabbb4e7cd772556))

---

# v0.28.0

## Features

* feat(cli): generate & execute gcp terraform code (#716) (
  Ross)([36e0c8c](https://github.com/lacework/go-sdk/commit/36e0c8c2cd700d7d2a8ae31e205178876ecf349c))

## Refactor

* refactor(cli): leverage anonymous lql (
  hazedav)([c607d81](https://github.com/lacework/go-sdk/commit/c607d81f199a5baaf522eaedae4d5d5b83296c77))

## Bug Fixes

* fix: use latest versions of TF provider and modules (#737) (Salim
  Afiune)([fc14f79](https://github.com/lacework/go-sdk/commit/fc14f799327a0a8348614b0742452d414c307673))

## Other Changes

* ci: fix qa tests (#736) (Salim
  Afiune)([bb6afa5](https://github.com/lacework/go-sdk/commit/bb6afa5261b91c36e84ff7987b59d9f658555e9f))
* ci: version bump to v0.27.1-dev (
  Lacework)([cb9055d](https://github.com/lacework/go-sdk/commit/cb9055d3e0ef95e5df9c57d3a14a53ea35da9748))
* test: fix policy severity filter test (#734) (
  Darren)([0aa05b1](https://github.com/lacework/go-sdk/commit/0aa05b11770263bb4bb447c75f5ea213d283e1c6))

---

# v0.27.0

## Features

* feat(cli): update query by loading it from the platform (#721) (Salim
  Afiune)([2f68adb](https://github.com/lacework/go-sdk/commit/2f68adbac8e63946f3f8e6bf6a92a5837ca8fed4))
* feat(cli): sort query and policy list output (
  hazedav)([68de649](https://github.com/lacework/go-sdk/commit/68de649549139b292d90ed41da87e0e902b12048))
* feat(cli): Add breadcrumb for "query show-source" (
  hazedav)([062421c](https://github.com/lacework/go-sdk/commit/062421ccfbbd25e06ffffa930338f5fe4afddafa))
* feat(lwgenerate): Add initial hcl generation support Azure - ALLY-394 (#682) (
  djmctavish)([66845d2](https://github.com/lacework/go-sdk/commit/66845d264ee2570ed7f7592e5a7e5af1cb517e3b))

## Refactor

* refactor(cli): further deprecate evaluatorId (
  hazedav)([7a63eca](https://github.com/lacework/go-sdk/commit/7a63ecae49c82583f2dc1cd7f8c4d38f381ce571))
* refactor(cli): Create cloud specific TF file (#712) (
  Ross)([3f9e0d8](https://github.com/lacework/go-sdk/commit/3f9e0d847404026b753d676a8f2f75f5096361e7))
* refactor(cli): vulnerability_expections packages human readable output (#711) (
  Darren)([c37afbf](https://github.com/lacework/go-sdk/commit/c37afbfd70d9d7f8cfcabbdf6e78b3a6ce89c3d7))

## Bug Fixes

* fix: use online machine for testing (#722) (
  chi-dev-2021)([e960b64](https://github.com/lacework/go-sdk/commit/e960b6437df33495287c32799d408839de8c3dae))
* fix(cli): delete listed twice on help (
  hazedav)([d606f9d](https://github.com/lacework/go-sdk/commit/d606f9d9d9531ccb94897086a93197d81c1e2d67))

## Other Changes

* chore(cli): run regex of integration test (#729) (Salim
  Afiune)([03d55cf](https://github.com/lacework/go-sdk/commit/03d55cf3025cdac9f3a0ed047bba2db3ec465ba3))
* build: deprecate centos-8 images (#709) (Salim
  Afiune)([5636bdd](https://github.com/lacework/go-sdk/commit/5636bdd699e10870da75dbd21be6d4a3ea64155c))
* ci: fix release pipeline (#730) (
  Darren)([6470a4a](https://github.com/lacework/go-sdk/commit/6470a4a12588681149d0c59a1cd27fb4a7f88a1b))
* ci: version bump to v0.26.1-dev (
  Lacework)([dccd3f5](https://github.com/lacework/go-sdk/commit/dccd3f5cad014dfdfa26040cf1f4af936ca15ce3))
* test(windows): fix TestQueryUpdateFromIDEditor (#728) (Salim
  Afiune)([424c939](https://github.com/lacework/go-sdk/commit/424c939f7aedd0b623b652965798903298ddf50a))
* test(integration): fix finding host id (#725) (Salim
  Afiune)([4e246d9](https://github.com/lacework/go-sdk/commit/4e246d972a1e0683b0f14a95340189668a368ba6))
* test(lql): disable policies tests (
  hazedav)([2647565](https://github.com/lacework/go-sdk/commit/2647565a354c33b46329b01e5db2172f836c0d04))
* test(cli): add comprehensive help text testing (#715) (
  hazedav)([b38c696](https://github.com/lacework/go-sdk/commit/b38c69625edc6561f2b0da0c455d5d30d04f7177))

---

# v0.26.0

## Features

* feat(cli): show resources affected by a compliance recommendation (#697) (
  Darren)([86f86ee](https://github.com/lacework/go-sdk/commit/86f86eef38840a28be310347f67261ba588a850c))
* feat(cli): improve user feedback running terraform commands (#706) (Matt
  Cadorette)([8526d2c](https://github.com/lacework/go-sdk/commit/8526d2ccfb8a48dbc8ac1315d84006aea058822a))
* feat(cli): add lacework query preview-source (
  hazedav)([40e8bbb](https://github.com/lacework/go-sdk/commit/40e8bbb03fa3fca6eecac48c24ffca1468821815))
* feat(cli): switch-profile command (#705) (Salim
  Afiune)([75237cd](https://github.com/lacework/go-sdk/commit/75237cdb15508c87cad3cc0694e304524a0da3f6))

## Other Changes

* chore(deps): bump github.com/spf13/cobra from 1.2.1 to 1.3.0 (#647) (
  dependabot[bot])([2cb15fb](https://github.com/lacework/go-sdk/commit/2cb15fb48186a45064aac0ecce340efcfb2e8de1))
* chore(deps): bump github.com/spf13/viper from 1.9.0 to 1.10.1 (#644) (
  dependabot[bot])([5dac81c](https://github.com/lacework/go-sdk/commit/5dac81cfb57b29dc8addcbb32b16f11abeeca836))
* chore(deps): bump github.com/honeycombio/libhoney-go (#652) (
  dependabot[bot])([6e5c7bd](https://github.com/lacework/go-sdk/commit/6e5c7bd61c01f1a3d84c9285d0d08d7a2dc6b3ef))
* chore(deps): bump github.com/hashicorp/terraform-exec (#704) (
  dependabot[bot])([9767bc8](https://github.com/lacework/go-sdk/commit/9767bc8a1c2531f936e685f0143400017198d6a7))
* chore(deps): bump go.uber.org/zap from 1.19.1 to 1.21.0 (#700) (
  dependabot[bot])([b8be0b0](https://github.com/lacework/go-sdk/commit/b8be0b0fef410a5a6179495c7e2a4a44d5e9b3c6))
* chore(deps): bump github.com/hashicorp/go-version from 1.3.0 to 1.4.0 (#702) (
  dependabot[bot])([1729f4c](https://github.com/lacework/go-sdk/commit/1729f4cb513d8b32d6a7133f8df72bc3afb1f3a9))
* chore(deps): bump github.com/briandowns/spinner from 1.13.0 to 1.17.0 (#701) (
  dependabot[bot])([9d18563](https://github.com/lacework/go-sdk/commit/9d185636fa0fa7c8d46dedae01d349d18d0a60c1))
* ci(cli): release CentOS end-of-life problem (#699) (Salim
  Afiune)([4bc4a7c](https://github.com/lacework/go-sdk/commit/4bc4a7ce8c37f1f4f1a21b058fe315ba855687b6))
* ci: version bump to v0.25.1-dev (
  Lacework)([20551d1](https://github.com/lacework/go-sdk/commit/20551d1b3b748bedc696d190acb7c13cfffc54e4))

---

# v0.25.0

## Features

* feat(api): add registry notifications field to docker v2 (#696) (
  Darren)([aac6801](https://github.com/lacework/go-sdk/commit/aac68016b548cc9c2799bb5b12c1f6e0f16aa8bd))
* feat(api): abstract primitives to interact with /v2/Vulnerabilities (#694) (Salim
  Afiune)([cd2622a](https://github.com/lacework/go-sdk/commit/cd2622a6d642c14e542e9997621b06f284a21245))
* feat(policies): add tags support (
  hazedav)([2a0ad67](https://github.com/lacework/go-sdk/commit/2a0ad670b86ec304b1acbbf72cb3c127348e7409))
* feat(cli): be able to change global LW_LQL_OPERATOR (Salim Afiune
  Maya)([a4c7746](https://github.com/lacework/go-sdk/commit/a4c77462c93d298adaeb7a80a6945a5d96e986ad))
* feat(cli): filter list of agents (Salim Afiune
  Maya)([8e200a0](https://github.com/lacework/go-sdk/commit/8e200a0c618c95a5206112565237a8cd7caac6eb))
* feat(api): implement Entities/Images/serch endpoint (Salim Afiune
  Maya)([33b58ef](https://github.com/lacework/go-sdk/commit/33b58ef0a2433b802baad3bae2d0586c422a4998))
* feat(api): implement Entities/Users/serch endpoint (Salim Afiune
  Maya)([b86f893](https://github.com/lacework/go-sdk/commit/b86f893be1d8b6fc9958611599a82bdef64b33d7))
* feat(api): support v2/Entities/Users (Salim Afiune
  Maya)([40c71d7](https://github.com/lacework/go-sdk/commit/40c71d7e2eba75bf5696c475265226bacc7a0943))
* feat(api): v2 pagination (Salim Afiune
  Maya)([255079d](https://github.com/lacework/go-sdk/commit/255079dd7c09873570e8412eab26ab01d03f96f5))
* feat(cli): new agent list command (Salim Afiune
  Maya)([8b056f5](https://github.com/lacework/go-sdk/commit/8b056f5100764ec3cca5413df38130014cad416d))
* feat(api): search data-export entities (Salim Afiune
  Maya)([5b87d33](https://github.com/lacework/go-sdk/commit/5b87d337bb585d18216715d89467203c1e94f42f))
* feat(cli): Manage Team Members (#671) (
  Darren)([6930449](https://github.com/lacework/go-sdk/commit/6930449b156bcf33a92fd4ef885c7ffe166ef369))
* feat(cli): Manage Vulnerability Exceptions (#663) (
  Darren)([0d62f11](https://github.com/lacework/go-sdk/commit/0d62f11b2a5264314489aa9bd23d7ca2a19d0b28))
* feat(cli): add --cascade flag to delete policy and query (#677) (
  hazedav)([6842a6b](https://github.com/lacework/go-sdk/commit/6842a6b162113471391118d34a03b1a022c5694d))
* feat(cli): display status of Azure compliance subscriptions (Salim Afiune
  Maya)([c479863](https://github.com/lacework/go-sdk/commit/c479863e71d63f301ea2d926db085f2eb8af2a4a))
* feat(cli): display status of GCP compliance projects (Salim Afiune
  Maya)([30e8448](https://github.com/lacework/go-sdk/commit/30e844823c26deee4ec01c3d40e2fdbc05111d6c))
* feat(cli): display status of AWS compliance accounts (Salim Afiune
  Maya)([e43978c](https://github.com/lacework/go-sdk/commit/e43978c9d40fd444e238deaa61059d6184df544d))

## Refactor

* refactor(cli): Refactor generate & generate_aws (#691) (
  Ross)([3fabe7b](https://github.com/lacework/go-sdk/commit/3fabe7b330284b0fbe044ea6af3a9b05c7e35b72))
* refactor(api): rename Pagination interface to Pageable (Salim Afiune
  Maya)([3b1e881](https://github.com/lacework/go-sdk/commit/3b1e881f52f0ee399084c072eaf53020a3510630))
* refactor(api): add logging to NextPage() func (Salim Afiune
  Maya)([810b723](https://github.com/lacework/go-sdk/commit/810b7231e0c682e6725f86d0f993c98ff7330c6e))
* refactor(cli): Refactor generate & generate_aws (#689) (
  Ross)([c2d4241](https://github.com/lacework/go-sdk/commit/c2d42418c7c1e56848048e475ad40b27694a7a94))
* refactor(cli): abstract shared funcs to internal package (
  hazedav)([5aa880d](https://github.com/lacework/go-sdk/commit/5aa880d5ccb09c75a9039cccf8313db30a020dff))
* refactor(cli): abstract shared funcs to internal package (
  hazedav)([5617a2d](https://github.com/lacework/go-sdk/commit/5617a2db0e1ccccee8e4d6df3b94c1a20d00bfc6))
* refactor(cli): rename 'aws' s3/cloudwatch to 'amazon' in cli integration prompt (#657) (
  Darren)([2c8757d](https://github.com/lacework/go-sdk/commit/2c8757d7b021c46e7c43c5a4eac353eafed5e269))

## Bug Fixes

* fix(cli): remove text wrapping on list-sources output (#692) (
  Darren)([49ca2ef](https://github.com/lacework/go-sdk/commit/49ca2ef0fb1826d38ce8e4d445ef14fc1e36d185))
* fix(cli): validate user filters via --filters (Salim Afiune
  Maya)([6ba02b0](https://github.com/lacework/go-sdk/commit/6ba02b086e2d650f36f4bf74fb482589244f9aec))
* fix(api): omit empty fields in v2 SearchFilter (Salim Afiune
  Maya)([5f4d2e7](https://github.com/lacework/go-sdk/commit/5f4d2e7065af5018c2768f898a333e7067640f69))
* fix(cli): check the account has agents to display (Salim Afiune
  Maya)([7e4ca27](https://github.com/lacework/go-sdk/commit/7e4ca27d43f9a6a9118a696e857884da4666b933))
* fix(lwgenerate): support lists for module attributes (#687) (Matt
  Cadorette)([6adfeeb](https://github.com/lacework/go-sdk/commit/6adfeeb1945f694491346e6f004251d900cecf5a))

## Other Changes

* style(cli): add breadcrumb to query list-sources cmd (#685) (
  Darren)([66ee797](https://github.com/lacework/go-sdk/commit/66ee797ea33ab813c2c5795cf374e4647889c360))
* chore(cli): update spinner message for compliance list cmds (Salim Afiune
  Maya)([c3d4aee](https://github.com/lacework/go-sdk/commit/c3d4aee0f6e5eb7066ab16b6b78b5657c55db8b3))
* ci: fix v1/v2 environment variables (#688) (Salim
  Afiune)([de4699a](https://github.com/lacework/go-sdk/commit/de4699a26a249cd325405f49602c6560d43b1822))
* ci: version bump to v0.24.1-dev (
  Lacework)([2df7446](https://github.com/lacework/go-sdk/commit/2df744647538fd2371b34cc0c1cc3dcc42cd71a7))
* test(cli): for compliance az list command (Salim Afiune
  Maya)([a98599c](https://github.com/lacework/go-sdk/commit/a98599c529efee480031f19f60042148a8bd31f2))

---

# v0.24.0

## Features

* feat(cli): make policyID optional for policy create (
  hazedav)([38006bf](https://github.com/lacework/go-sdk/commit/38006bf39a3cae293595e13c6b6bfb9d05936f1d))
* feat(lwgenerate): Add initial hcl generation support GCP - ALLY-393 (#654) (
  Ross)([fb34a64](https://github.com/lacework/go-sdk/commit/fb34a64b0ddc71bc2140b42c66fa0f38160e762a))
* feat(cli): add csv output support to vuln reports (#649) (Matt
  Cadorette)([c76095e](https://github.com/lacework/go-sdk/commit/c76095e88c60f02c02b449c09ffbeaf2dd5bf1b1))
* feat(api): Vulnerability Exceptions v2 Service (#627) (
  Darren)([ef65e23](https://github.com/lacework/go-sdk/commit/ef65e232b7c7c650562db5cded0b3ee31a461a55))
* feat(cli): add apply flag for aws iac generation & fix cache issues (#642) (Matt
  Cadorette)([c5014bb](https://github.com/lacework/go-sdk/commit/c5014bb73f1db4977a9500dc4f92e686fa9d1c74))

## Refactor

* refactor(test_resources): remove Terraform code dep (#636) (Salim
  Afiune)([aca0101](https://github.com/lacework/go-sdk/commit/aca01017734598d64526b9a7008fe9171d74fe5c))

## Bug Fixes

* fix: expiration time (#656) (
  Darren)([4726f68](https://github.com/lacework/go-sdk/commit/4726f68ca5731208e140e006c98716aba3d626a3))
* fix: non_os_package_support omitted from json request when false (#655) (
  Darren)([c8ab03a](https://github.com/lacework/go-sdk/commit/c8ab03ae4ad1616ef172db3230fd3f1baff6e389))
* fix: handle unsupported integration types (#645) (
  Darren)([6c1a5d0](https://github.com/lacework/go-sdk/commit/6c1a5d05f937eaac89893e5cb483e7d6c045e17a))
* fix(cli): clarify container error message (#643) (Salim
  Afiune)([74592ff](https://github.com/lacework/go-sdk/commit/74592ff8fadf8f78606c940d347d9024dfbe8a8c))

## Documentation Updates

* docs(iac-generation): update aws iac tf command help (#653) (Matt
  Cadorette)([5c09072](https://github.com/lacework/go-sdk/commit/5c09072b7079f72be40fd154da471755d2eeccc7))
* docs(cli): Improve iac-generate/aws help message (#638) (Matt
  Cadorette)([127e221](https://github.com/lacework/go-sdk/commit/127e2212018fb2124b3d15ec2fcdd155e1326533))

## Other Changes

* chore(deps): bump github.com/hashicorp/hcl/v2 from 2.10.1 to 2.11.1 (#631) (
  dependabot[bot])([a82d22b](https://github.com/lacework/go-sdk/commit/a82d22b3a3d3ad6eec69a7d03a45a19e298be893))
* chore(deps): bump github.com/mitchellh/mapstructure from 1.4.2 to 1.4.3 (#630) (
  dependabot[bot])([4f08643](https://github.com/lacework/go-sdk/commit/4f0864393b611ce5d6527393ff2e8d5df64cfd24))
* chore(deps): bump github.com/peterbourgon/diskv/v3 from 3.0.0 to 3.0.1 (#610) (
  dependabot[bot])([c1792d5](https://github.com/lacework/go-sdk/commit/c1792d51bd563f852cf154a7fdfbf5d7560192e0))
* chore(deps): Bump GoImports version to v0.1.8 (#633) (
  Ross)([22f5283](https://github.com/lacework/go-sdk/commit/22f5283a93d0e3552c7343ad7c6104c0f3389ab8))
* chore(deps): Bump gotestsum version and use instead of go test (#632) (
  Ross)([7f05f44](https://github.com/lacework/go-sdk/commit/7f05f441ce579a9654e7809d75a7a81845e9b380))
* ci: test on windows systems  (#648) (Salim
  Afiune)([70fe066](https://github.com/lacework/go-sdk/commit/70fe0665af92afba20e28af7d1c3319043ba78b1))
* ci: version bump to v0.23.1-dev (
  Lacework)([cfc7bc0](https://github.com/lacework/go-sdk/commit/cfc7bc00f6ff515a0f8b64b8d6d6bb6ceb684fc0))
* test: fix host vuln test (#659) (
  Darren)([d28ea4c](https://github.com/lacework/go-sdk/commit/d28ea4c366eeca098547efc0c7230f30d9e7a141))
* test: avoid running v1/v2 migration on every test (#651) (Salim
  Afiune)([63735c2](https://github.com/lacework/go-sdk/commit/63735c2db0e518ff381eaa8949ff5fc33cfca64b))

---

# v0.23.0

## Features

* feat(api): CRUD operations for org team member (#628) (
  vatasha)([80015b9](https://github.com/lacework/go-sdk/commit/80015b9807e3dc541a6e3e3c1b8c726a84995209))
* feat(cli): generate & execute aws terraform code (#626) (Matt
  Cadorette)([992bfa6](https://github.com/lacework/go-sdk/commit/992bfa6c3fbd08c302d53bb02a27831278f0a493))
* feat(cli): Manage Report Rules (#625) (
  Darren)([890350b](https://github.com/lacework/go-sdk/commit/890350b2b34af6380de15df8d665828bb313776a))
* feat(cli): add breadcrums to compliance commands (Salim Afiune
  Maya)([f133971](https://github.com/lacework/go-sdk/commit/f1339716c91f7cdb7d2a4fd64bac9f457cedb6ab))
* feat(cli): cache compliance reports for 30 minutes (Salim Afiune
  Maya)([77f82aa](https://github.com/lacework/go-sdk/commit/77f82aad8e39c208afbef6a35ad62e4206fe56ae))
* feat(api): implement crud operations for team members (#614) (
  vatasha)([4aa40a2](https://github.com/lacework/go-sdk/commit/4aa40a23ed2c53beeea0bc60a851bf47488b4d52))

## Bug Fixes

* fix(cli): make sure progress bar starts with space (#635) (Salim
  Afiune)([4153f64](https://github.com/lacework/go-sdk/commit/4153f64f071bb85c0718967be4f71a9bd7071e65))
* fix(cli): fix cache directory for azure comp reports (Salim Afiune
  Maya)([19be146](https://github.com/lacework/go-sdk/commit/19be1465a3f73f6359c32454863312a407318e6b))
* fix: report rule notifications (#621) (
  Darren)([f37b913](https://github.com/lacework/go-sdk/commit/f37b91349c031cf0d92c79a7f4dfc181446e69e3))

## Other Changes

* ci: fix test resource Terraform pipeline (#634) (Salim
  Afiune)([c14756e](https://github.com/lacework/go-sdk/commit/c14756e12c0f051e871e0f7a7aeccd1089cdd875))
* ci: version bump to v0.22.1-dev (
  Lacework)([e12f022](https://github.com/lacework/go-sdk/commit/e12f022356992ae477dbbf211f9439dd4b080dce))
* test(policy): update TestPolicyCreateURL (
  hazedav)([a89177c](https://github.com/lacework/go-sdk/commit/a89177ce3d82660c8700b555c0596f7a98e53d47))
* test(policy): adds resource for dynamic accnt (
  hazedav)([b28d290](https://github.com/lacework/go-sdk/commit/b28d2902aa6d1a557af913c60a0be7a68c3b9895))

---

# v0.22.0

## Features

* feat(cli): cache host vuln assessments for 1 hour (Salim Afiune
  Maya)([7f50bd2](https://github.com/lacework/go-sdk/commit/7f50bd23654a5566c35ba13e88ad2f4c9c8dafa9))
* feat(cli): logic to store assets with expiration time (Salim Afiune
  Maya)([ad920ca](https://github.com/lacework/go-sdk/commit/ad920ca5d638adb4f6331a80550381d8d1b3021d))
* feat(lwgenerate): Add initial hcl generation support (#611) (Matt
  Cadorette)([f9e9e78](https://github.com/lacework/go-sdk/commit/f9e9e781b73c9a92415702e052d10fed19a72ac2))

## Bug Fixes

* fix(cli): PKG typo in host vuln table (Salim Afiune
  Maya)([21ff81a](https://github.com/lacework/go-sdk/commit/21ff81a762d8c1dd3335bf88f8598ecebacc5686))
* fix(cli): avoid using "Fixed" host vulnerabilities (Salim Afiune
  Maya)([49ce7a7](https://github.com/lacework/go-sdk/commit/49ce7a70bca8eb3125c53a22e99c219fa61d3d33))

## Other Changes

* ci: version bump to v0.21.1-dev (
  Lacework)([30637bc](https://github.com/lacework/go-sdk/commit/30637bcce5b7aefb164c1ff3587753dbd4b5844d))
* test(cli): verify host vuln report rendering (Salim Afiune
  Maya)([b7d72e2](https://github.com/lacework/go-sdk/commit/b7d72e29b14b0987d021ecb03f618bbf3ac779a0))

---

# v0.21.0

## Features

* feat: Auto docs generation for the Lacework CLI (#616) (Salim
  Afiune)([81b7fdb](https://github.com/lacework/go-sdk/commit/81b7fdbf50c1e560277c44be38e9fd2473ca15c7))
* feat: Filter vuln ctr list-assessments by registry (#609) (
  Darren)([264b450](https://github.com/lacework/go-sdk/commit/264b4504633267cd6f1929eb39109f6cafb74692))
* feat(api): Report Rules v2 Service (#607) (
  Darren)([26e96cc](https://github.com/lacework/go-sdk/commit/26e96cc5c95e5188b6d0b82977371f578eba55e9))

## Bug Fixes

* fix(cli): remove sub-account during reconfiguration (#602) (Salim
  Afiune)([d623afe](https://github.com/lacework/go-sdk/commit/d623afedc2152cbeb7abcd90f87e6d8a475e95a0))

## Documentation Updates

* docs: Run make generate-docs (Darren
  Murray)([b6438a1](https://github.com/lacework/go-sdk/commit/b6438a1595363faa331a1e79c89dab492e8c5e44))
* docs: Vuln ctr list-assessments registry flag (Darren
  Murray)([cf32b2f](https://github.com/lacework/go-sdk/commit/cf32b2f41b306fcd18103ec4fae992058de4d273))

## Other Changes

* chore: ran make go-vendor with go 1.17 (#615) (Salim
  Afiune)([432104f](https://github.com/lacework/go-sdk/commit/432104faf0522c58ee404232600714360aa64bea))
* chore: Update minimum go version -> 1.17 (#612) (
  Darren)([c6e183c](https://github.com/lacework/go-sdk/commit/c6e183c22d576b11698989a34443e4c4c53d235a))
* chore(deps): bump github.com/honeycombio/libhoney-go (#605) (
  dependabot[bot])([7226624](https://github.com/lacework/go-sdk/commit/7226624e163015d12cd5b58042ec552df4ad64ff))
* chore(cli): Update s3, slack, and cloudwatch alert channels cli commands to point to API v2 (#606) (
  vatasha)([65887be](https://github.com/lacework/go-sdk/commit/65887be28d97aa2d16b5f528ca2cc227ea7e94fd))
* ci: add Azure integration as a test resource (#608) (Salim
  Afiune)([5cab0ef](https://github.com/lacework/go-sdk/commit/5cab0ef1d1435c64fd3eb52c086ee63e756c41ad))
* ci: Update go conditinal build comment to 1.17 preferred format (#604) (
  Darren)([a8f90f6](https://github.com/lacework/go-sdk/commit/a8f90f668ac4c9ac7b6eaec7b3720de390d02a5f))
* ci: version bump to v0.20.1-dev (
  Lacework)([15761a7](https://github.com/lacework/go-sdk/commit/15761a7f5e6f1cc750f3c1b0911f2dabe7eb3c00))
* test: Fix failing alert rule test (#603) (
  Darren)([c1b8ca5](https://github.com/lacework/go-sdk/commit/c1b8ca59b088bda0715562f846ddb05bfee7bfe1))

---

# v0.20.0

## Features

* feat(cli): Manage Alert Rules in Lacework CLI (#597) (
  Darren)([4d5648c](https://github.com/lacework/go-sdk/commit/4d5648cce2430b670167d3d66c3d8aee860b5f39))

## Bug Fixes

* fix(cli): avoid panic when assessment has NO vulns (#600) (Salim
  Afiune)([4638c7e](https://github.com/lacework/go-sdk/commit/4638c7e7d6b08e8f9bb123ed897321897a33ab4e))
* fix: ResourceGroups lastUpdated time format (#599) (
  Darren)([506954f](https://github.com/lacework/go-sdk/commit/506954fd1451d192a43d2ae1a38fb1e2b6261b4e))
* fix(cli): making evaluatorId optional (
  hazedav)([a7da3c0](https://github.com/lacework/go-sdk/commit/a7da3c0bc346ac1db957a3854c04c2c5a6a537c6))

## Other Changes

* ci: version bump to v0.19.1-dev (
  Lacework)([e8e36dd](https://github.com/lacework/go-sdk/commit/e8e36ddf314463540288dea6d577a12b5d41456a))
* test(cli): adds host query integration test file (
  hazedav)([7be96a7](https://github.com/lacework/go-sdk/commit/7be96a79d2bb3f9df84c259615fe0027110342bf))

---

# v0.19.0

## Features

* feat: AlertRuleSeverities constructor from int slice (#595) (
  Darren)([637bcab](https://github.com/lacework/go-sdk/commit/637bcabd9f59a3d73cf0c740b58d7cfecb5b0c90))
* feat(api): Migrate Jira Cloud and Jira Server to API v2 (#590) (
  vatasha)([15369f3](https://github.com/lacework/go-sdk/commit/15369f35ef472fe8d2bf162e4e0e4d99304413aa))

## Refactor

* refactor: Add function to convert string slice to AlertRuleSeverities (#591) (
  Darren)([b3ae293](https://github.com/lacework/go-sdk/commit/b3ae2931f2a201059909f72b91b79f7c7c20ffc8))

## Bug Fixes

* fix(cli): affected resource count in compliance reports (#592) (Salim
  Afiune)([7dc18d9](https://github.com/lacework/go-sdk/commit/7dc18d9288285ad3196ff3f186641891cda8ec7d))

## Other Changes

* ci: version bump to v0.18.1-dev (
  Lacework)([21b392a](https://github.com/lacework/go-sdk/commit/21b392ad93c6a2f83d003d80125bc3da16feadc4))

---

# v0.18.0

## Features

* feat(cli): unhiding query "sources" cmd (
  hazedav)([4e2513f](https://github.com/lacework/go-sdk/commit/4e2513fc3ccd4f4b9818378c98f771df9f3a4546))
* feat(api): Alert Rules v2 service (#583) (
  Darren)([a549a8c](https://github.com/lacework/go-sdk/commit/a549a8c501e00a42be879d38197fc13a4984c914))
* feat(cli): Migrate New Relic alert channel to API v2 (#586) (
  Darren)([6f32767](https://github.com/lacework/go-sdk/commit/6f327671b19e221859bff761e495ba696279bb61))
* feat(cli): Migrate Pager Duty alert channel to API v2 (#585) (
  Darren)([68be1ec](https://github.com/lacework/go-sdk/commit/68be1eccbce846810ec35ece9819f99b1c73d5e9))
* feat(api): Migrate GCP Pub Sub alert channel to API v2 (#581) (
  vatasha)([7847d3b](https://github.com/lacework/go-sdk/commit/7847d3bddcaf7e2f8174684f25588559db4cdb37))
* feat(cli): Migrate Ibm QRadar alert channel to API v2 (#584) (
  Darren)([b2a0cde](https://github.com/lacework/go-sdk/commit/b2a0cde29e8f6400934d9ac4dc40fab4fe1ef513))
* feat(api): Migrate Service Now alert channel to API v2 (#582) (
  Darren)([692c8e5](https://github.com/lacework/go-sdk/commit/692c8e58650ed98891cd9202ea273400f51744b3))
* feat: Migrate Splunk alert channel to V2 (#580) (
  Darren)([c0c422b](https://github.com/lacework/go-sdk/commit/c0c422b71b7893958c0641c7d75b2401546bfee5))
* feat(cli): Add support for new GCP and AWS compliance reports (#579) (
  iancrichardson)([1de409b](https://github.com/lacework/go-sdk/commit/1de409b0aaaa9d20b8b882f04db5dae679157f07))
* feat(cli): adding lql query data sources (
  hazedav)([8828823](https://github.com/lacework/go-sdk/commit/882882383471db9b32076e047d9d00ab3d0a37c7))
* feat(api): Microsoft Teams alert channel APIv2 (Salim Afiune
  Maya)([57b318e](https://github.com/lacework/go-sdk/commit/57b318e102749c5bd34901f595424671b61c0b1e))
* feat(api): Migrate Cisco Spark Webhook alert channel to API v2 (#577) (
  vatasha)([9481999](https://github.com/lacework/go-sdk/commit/9481999f2faeabde3a9085601f5fbfdfee17531c))

## Refactor

* refactor(cli): create Microsoft Teams alert channels via APIv2 (Salim Afiune
  Maya)([74372b8](https://github.com/lacework/go-sdk/commit/74372b8ba1303eeb259799b3d080c0329fb34253))

## Bug Fixes

* fix: Splunk ssl (#587) (
  Darren)([eab94cf](https://github.com/lacework/go-sdk/commit/eab94cf33b16b48cb7fb3146c2aa903a7ae37d7e))

## Other Changes

* style: rename TeamsUrl in favor of TeamsURL (Salim Afiune
  Maya)([bd46e0d](https://github.com/lacework/go-sdk/commit/bd46e0d46c37017834fcd6521ade92d28cb0f9f9))
* ci: version bump to v0.17.1-dev (
  Lacework)([9df310c](https://github.com/lacework/go-sdk/commit/9df310c0a20d97b79ecfdf4c6ffe520b87456d93))

---

# v0.17.0

## Features

* feat(api): Migrate VictorOps alert channel to V2 (#573) (
  Darren)([6bc3b88](https://github.com/lacework/go-sdk/commit/6bc3b881ab72d5d767ae54fe5b9c14dc9dc0a17a))
* feat(api): Webhook alert channel(v2) (#565) (
  Darren)([06915a0](https://github.com/lacework/go-sdk/commit/06915a0271330abfedce23606f04c5a1cf16a077))
* feat(api): Cloudwatch alert channel(v2) (#562) (
  Vatasha)([d0c9be7](https://github.com/lacework/go-sdk/commit/d0c9be752ec6a2552cc7abd18acfa39fbca2f5c6))
* feat(api): Datadog alert channel(v2) (#570) (
  Vatasha)([3135b95](https://github.com/lacework/go-sdk/commit/3135b955ebbfe5d1a68d1328c01b900f413a1db9))

## Refactor

* refactor(cli): use APIv2 to create Datadog integrations + update api/examples (#571) (
  vatasha)([9d6ea0e](https://github.com/lacework/go-sdk/commit/9d6ea0eba0fc34ef594fac141d55cded3b0d7ea6))

## Bug Fixes

* fix(cli): disable cache mechanism in install.sh (#563) (Salim
  Afiune)([0ccd32e](https://github.com/lacework/go-sdk/commit/0ccd32e682e99eda369121c97eb7d6c695ddc9e8))
* fix:Add issue grouping field to cloudwatch alert channel data (#566) (
  vatasha)([1e93e91](https://github.com/lacework/go-sdk/commit/1e93e91beb7a0cd95c49b211486b1ae032a0b2ae))
* fix(cli): install.ps1 avoid fail for False boolean (#564) (Larry
  Eichenbaum)([32604df](https://github.com/lacework/go-sdk/commit/32604df46c095bf8efa28dcc1b32c0fef14b1e16))

## Documentation Updates

* docs: improve developer guidelines (Salim Afiune
  Maya)([80f8e9e](https://github.com/lacework/go-sdk/commit/80f8e9e2f62ce0a567883824eb2105866ddb274a))
* docs: add signed commits documentation (Salim Afiune
  Maya)([b8b6325](https://github.com/lacework/go-sdk/commit/b8b63257de8a8ae2c7b0fc6b4fa7711eaa87ff70))
* docs: fix broken link in CONTRIBUTING.md (Salim Afiune
  Maya)([92a4d0c](https://github.com/lacework/go-sdk/commit/92a4d0cfcdd208305f1c4b3c321cd6639fd5dda1))
* docs: improve pull request template (Salim Afiune
  Maya)([967f9d6](https://github.com/lacework/go-sdk/commit/967f9d6d642b563b043ed0b6207ed09374d86523))

## Other Changes

* chore(deps): bump go.uber.org/zap from 1.19.0 to 1.19.1 (#568) (
  dependabot[bot])([3f4a1f2](https://github.com/lacework/go-sdk/commit/3f4a1f2e07e9fd4c3009e2dcc6d0a8309f53fe23))
* chore(deps): bump github.com/fatih/color from 1.12.0 to 1.13.0 (#556) (
  dependabot[bot])([bf79562](https://github.com/lacework/go-sdk/commit/bf7956236a1e9dcdc1a7cdb891da57918d510471))
* chore(deps): bump github.com/honeycombio/libhoney-go (#557) (
  dependabot[bot])([39fa065](https://github.com/lacework/go-sdk/commit/39fa065c3f46449c2033311b89eebe036f994769))
* chore(deps): bump github.com/spf13/viper from 1.8.1 to 1.9.0 (#553) (
  dependabot[bot])([bf4a39a](https://github.com/lacework/go-sdk/commit/bf4a39a802daa8e017d141f7a101afb03427ebb1))
* chore(deps): bump github.com/mitchellh/mapstructure from 1.4.1 to 1.4.2 (#552) (
  dependabot[bot])([649c729](https://github.com/lacework/go-sdk/commit/649c729f37f5067cda2ecac55deac1a32c0ec6a6))
* ci: fix downgrading of go packages (#575) (
  vatasha)([1a816c9](https://github.com/lacework/go-sdk/commit/1a816c9f8b7d47831326b3344d94072ff8350c6c))
* ci: version bump to v0.16.1-dev (
  Lacework)([bb5caa0](https://github.com/lacework/go-sdk/commit/bb5caa057d2ae5515aebbaf2a6d9ec2247f1e4ea))

---

# v0.16.0

## Features

* feat: Add NonOSPackageSupport container registries (#560) (
  Darren)([c61e4d3](https://github.com/lacework/go-sdk/commit/c61e4d3a4fbf26b821c831574a2c5556b90e424b))
* feat(cli): pass-through used profile in recommendations (#554) (Salim
  Afiune)([d5340c1](https://github.com/lacework/go-sdk/commit/d5340c15e2a196897f0cd4339b45d01236d0d5e9))
* feat(telemetry): send information about flags (#551) (Salim
  Afiune)([8f626d6](https://github.com/lacework/go-sdk/commit/8f626d63905bd43e017fad9f83476972a7aa5ae0))

## Documentation Updates

* docs: Contributor Documentation (#550) (
  Darren)([a278374](https://github.com/lacework/go-sdk/commit/a27837457423d79aedeaf3d1e66bba0cea4069df))

## Other Changes

* chore(deps): bump github.com/AlecAivazis/survey/v2 from 2.2.12 to 2.3.2 (#549) (
  dependabot[bot])([00deaa0](https://github.com/lacework/go-sdk/commit/00deaa034ae2b02508381756bb900c47e799e9f9))
* chore(deps): bump go.uber.org/zap from 1.18.1 to 1.19.1 (#548) (
  dependabot[bot])([f98cdc3](https://github.com/lacework/go-sdk/commit/f98cdc300bc1acf0027c6128c149a18260ab8256))
* chore(deps): bump github.com/BurntSushi/toml from 0.3.1 to 0.4.1 (#509) (
  dependabot[bot])([e1184a1](https://github.com/lacework/go-sdk/commit/e1184a1af99435b8348a4a15ad87bc2e96ac2b16))
* ci: ran make go-vendor to fix pipeline (#558) (Salim
  Afiune)([576b127](https://github.com/lacework/go-sdk/commit/576b1272b58534a6b7726e10caec751cbab9f509))
* ci: version bump to v0.15.1-dev (
  Lacework)([aa0751f](https://github.com/lacework/go-sdk/commit/aa0751fd8940ddd63759ab0a9cdc8ad6f9660dfc))

---

# v0.15.0

## Features

* feat(lwdomain): export go package (#545) (Salim
  Afiune)([13fb167](https://github.com/lacework/go-sdk/commit/13fb1675fdc30446ef903404a191104cf71195ee))
* feat(cli): stdin support for query/policy (
  hazedav)([770f0ee](https://github.com/lacework/go-sdk/commit/770f0ee26c613a59a166673f9d4a372a586f7da1))
* feat(cli): Manage Resource Groups in the lacework cli (#538) (
  Darren)([5e27cc8](https://github.com/lacework/go-sdk/commit/5e27cc8f0fba66d62b9113be1af0f97957529063))
* feat: Resource Groups (#537) (
  Darren)([754e8f4](https://github.com/lacework/go-sdk/commit/754e8f47ba2fbf1a68ca8b71fe3a761a36f1eb6a))
* feat(api): Add Non-OS Package Support for APIv1 (#539) (Andre
  Elizondo)([92e20e6](https://github.com/lacework/go-sdk/commit/92e20e6fa48be3ce192b84cfec832d9261f892e2))
* feat(cli): switch over to API v2 by default (#519) (Salim
  Afiune)([d605b45](https://github.com/lacework/go-sdk/commit/d605b45ad95920990ce7b0e8e662b0c12369afff))

## Refactor

* refactor: Add resource group specific create functions(#542) (
  Darren)([8a78438](https://github.com/lacework/go-sdk/commit/8a7843885afdfcde4ac0d98381748b0f220dbfcc))

## Bug Fixes

* fix(cli): Output correct resource group json (#543) (
  Darren)([3311ef2](https://github.com/lacework/go-sdk/commit/3311ef262edc058f29f0f151d938c3e76ab15f05))
* fix(install.ps1): Force Set-Variable as Constant (#541) (Salim
  Afiune)([6087ac2](https://github.com/lacework/go-sdk/commit/6087ac2b977132ba13350b248e2551c15c75881c))

## Documentation Updates

* docs(cli): correct typo v1 to v2 (#536) (
  lhasadreams)([cff6064](https://github.com/lacework/go-sdk/commit/cff6064a25d711d784dd0623c41739af78f0ff92))

## Other Changes

* ci: version bump to v0.14.1-dev (
  Lacework)([9417890](https://github.com/lacework/go-sdk/commit/941789044a3cdba6565209570149cf15475ad5fe))
* test(cli): re-enable vuln ctr scan tests (#532) (Salim
  Afiune)([7697c23](https://github.com/lacework/go-sdk/commit/7697c23c76f581f1610e16c83f0c787bac74695d))
* test(cli): disable vuln ctr scan tests (#530) (Salim
  Afiune)([0787b1b](https://github.com/lacework/go-sdk/commit/0787b1b6f7f74a79a528e0dde5157b172732344d))
* test(cli): re-enable query/policy integration (#526) (
  hazedav)([d54b51a](https://github.com/lacework/go-sdk/commit/d54b51a46a2154d8912859c64eb761d6bb72256f))

---

# v0.14.0

## Features

* feat(cli): More detailed CSV compliance reports (#522) (Matt
  Cadorette)([bd02a44](https://github.com/lacework/go-sdk/commit/bd02a44009e5cf00901ced32385add97e04be0aa))
* feat(api): AwsS3 Alert Channel (v2) (Salim Afiune
  Maya)([507cac1](https://github.com/lacework/go-sdk/commit/507cac14781b986662a2a78b7d436cc2ba1250ee))
* feat(api): SlackChannel Alert Channel (v2) (Salim Afiune
  Maya)([e6ac8cb](https://github.com/lacework/go-sdk/commit/e6ac8cb8eed31646020aeb3bbf5d2ef539258389))

## Refactor

* refactor: rename alert channel enums (Salim Afiune
  Maya)([f727b9b](https://github.com/lacework/go-sdk/commit/f727b9b78a79047b5df5665324d54d0cb4302b7b))

## Bug Fixes

* fix(cli): removing policyUi from policy schema (#523) (
  hazedav)([5586aec](https://github.com/lacework/go-sdk/commit/5586aec415034a431490c0b3cfb4d51385bfdd23))
* fix(cli): handle scan status 'Partial' (#525) (
  Darren)([a8ba4b7](https://github.com/lacework/go-sdk/commit/a8ba4b722cf211f527c685519504359525bd0ebe))
* fix(cli): query run via editor (
  hazedav)([ca3f2a3](https://github.com/lacework/go-sdk/commit/ca3f2a34bb76b97f42b6a5cafb664839a1fc3689))
* fix(cli): Detect if running in cloudshell to output correct update command (#521) (
  Darren)([477fadd](https://github.com/lacework/go-sdk/commit/477fadd9a96398198f397eb57c29231761173a9b))

## Other Changes

* chore(api): update v2 alert channel tests (Salim Afiune
  Maya)([3581d0a](https://github.com/lacework/go-sdk/commit/3581d0a2237d176c6e1d1d0dfa749b1a3bce8480))
* ci: version bump to v0.13.1-dev (
  Lacework)([39cc865](https://github.com/lacework/go-sdk/commit/39cc865741096ae4f7a29c09722bd38d7637752d))

---

# v0.13.0

## Features

* feat(cli): install agents via SSH on a custom port (#516) (Salim
  Afiune)([89ff730](https://github.com/lacework/go-sdk/commit/89ff730d03ef79a48ff505de9e1ff39ae48d5661))
* feat(cli): new agent install --trust_host_key flag (#515) (Salim
  Afiune)([fbd9934](https://github.com/lacework/go-sdk/commit/fbd99344d44f72e60f784695cf79fa3b56a243b5))
* feat(cli): new --api_token global flag (#513) (Salim
  Afiune)([9cf1271](https://github.com/lacework/go-sdk/commit/9cf12719981516e83afe4390996d42925ccde3c9))
* feat(cli): unhide query and policy commands (
  hazedav)([d1990d2](https://github.com/lacework/go-sdk/commit/d1990d2c065b0ea6efe0de247629b27afa9ee828))

## Bug Fixes

* fix(cli): vuln host scan-pkg-manifest --local centos 6.10 (#514) (
  Darren)([83f8884](https://github.com/lacework/go-sdk/commit/83f888455cb022087ad5ac28354a2a9fe4989bfe))
* fix(install.ps1): verify if running as administrator (#510) (Salim
  Afiune)([d9c2284](https://github.com/lacework/go-sdk/commit/d9c2284b2ba44ec9c15133c0e382348acbff98b3))

## Documentation Updates

* docs(cli): add inline help docs for query/policy cmds (#517) (Salim
  Afiune)([e471393](https://github.com/lacework/go-sdk/commit/e47139331a0e3482c34f39ee56b2917b69b70b33))

## Other Changes

* ci: version bump to v0.12.1-dev (
  Lacework)([af1cb50](https://github.com/lacework/go-sdk/commit/af1cb501c8b75034c46c1e089cc45c1c62a73f19))
* test(cli): go test modularity via build tags (#511) (
  hazedav)([b83a4af](https://github.com/lacework/go-sdk/commit/b83a4afb799a0f7fa1297e7ed027cf1c4c1b829a))

---

# v0.12.0

## Features

* feat(cli): Output v2 integration state details (#505) (
  Darren)([e2cc6f1](https://github.com/lacework/go-sdk/commit/e2cc6f1cbd3da43df6e014ccd9075ec59af6ae06))
* feat: Add apiV2AlertChannelTest (#502) (
  Darren)([f968b89](https://github.com/lacework/go-sdk/commit/f968b89bb0fa55d8eb84cf439626e2610507b883))
* feat(cli): create Github Registry integrations (Salim Afiune
  Maya)([d3756e4](https://github.com/lacework/go-sdk/commit/d3756e4b666de9039e7ccc01b706f5ed4e750f1e))
* feat(cli): create Google Artifact Registry integrations (Salim Afiune
  Maya)([0f1bb68](https://github.com/lacework/go-sdk/commit/0f1bb6866bc27e5e9e39fa79cdda350513c3f789))
* feat(api): add Github Registry support (Salim Afiune
  Maya)([697ca98](https://github.com/lacework/go-sdk/commit/697ca98ad540c6a0e79eeda26d0e2b0cb7f9271a))
* feat(api): add Google Artifact Registry support (Salim Afiune
  Maya)([503ae99](https://github.com/lacework/go-sdk/commit/503ae9953f510688974dafd35287d3d3fe4e12b0))
* feat(api): new v2 Container Registries service (Salim Afiune
  Maya)([46212c6](https://github.com/lacework/go-sdk/commit/46212c6b72fe964eb56f3618afec67dc2a641673))

## Refactor

* refactor(cli): policy update api args (
  hazedav)([a9add4b](https://github.com/lacework/go-sdk/commit/a9add4b7144b9461c4a398801e2ae6d10f6936cf))
* refactor(cli): split survey limits into diff funcs (Salim Afiune
  Maya)([1e1dca4](https://github.com/lacework/go-sdk/commit/1e1dca49a7a485ba81ce21a7dfa90c3ea7cb6dd2))

## Bug Fixes

* fix: comp [aws | gcp | az] list json output (#497) (
  Darren)([c3fb601](https://github.com/lacework/go-sdk/commit/c3fb60168e186611868943e3eaea3d92973c7aec))

## Other Changes

* chore(deps): bump github.com/honeycombio/libhoney-go (#494) (
  dependabot[bot])([9c1985c](https://github.com/lacework/go-sdk/commit/9c1985c0dbd3341ad2c47b797ebbce0a1941b0c1))
* chore: fixed minor typos (Salim Afiune
  Maya)([a13e58b](https://github.com/lacework/go-sdk/commit/a13e58bda7b1294936538481609e3421fdc4c0df))
* chore(deps): add github.com/fatih/structs package (Salim Afiune
  Maya)([89d6e0d](https://github.com/lacework/go-sdk/commit/89d6e0d30cb5405540f9493874eebd7652f5f5ef))
* ci: sign lacework-releng commits (#500) (Salim
  Afiune)([523097c](https://github.com/lacework/go-sdk/commit/523097ce8b4c6ede3317d1143446cf41df85a27f))
* test(cli): updating policy integration resource (
  hazedav)([acceada](https://github.com/lacework/go-sdk/commit/acceadae2c7bca2364ebe5ddda51e39b6158ced4))

---

# v0.11.1

## Bug Fixes

* fix(cli): scan-pkg-manifest returns valid JSON when no vulns are found (#498) (
  Darren)([0cdcc7e](https://github.com/lacework/go-sdk/commit/0cdcc7ed90e446b61db4bb4f55f47e1db35ff14f))

---

# v0.11.0

## Features

* feat(cli): manage EmailUser alert channels (Salim Afiune
  Maya)([4624dfb](https://github.com/lacework/go-sdk/commit/4624dfb6524f653f98eae7c87357d99c247e6f6c))
* feat(api): EmailUser Alert Channel (v2) (#492) (Salim
  Afiune)([45cc8f4](https://github.com/lacework/go-sdk/commit/45cc8f4eef3e0eacac66f4cb30be47b20fec9a6b))
* feat(cli): apiv2 support for query run (
  hazedav)([6d8c766](https://github.com/lacework/go-sdk/commit/6d8c766b3838cbfe45ce874c3741820aeafdd5bc))

## Bug Fixes

* fix(cli): integration show align left (Salim Afiune
  Maya)([df17625](https://github.com/lacework/go-sdk/commit/df176256fe0f8cd49bc95298debe3001c539933d))
* fix(cli): integration show single object (Salim Afiune
  Maya)([b26a3d0](https://github.com/lacework/go-sdk/commit/b26a3d07f062ee39be6ba745b052de49ba4ccea5))
* fix(cli): evalFrequency optional for policyCreate (
  hazedav)([03d527e](https://github.com/lacework/go-sdk/commit/03d527edb8e02229625176ed30548c50a43c0e92))

## Other Changes

* chore(cli): apply lint recommendations (Salim Afiune
  Maya)([931b2db](https://github.com/lacework/go-sdk/commit/931b2db86daad2cfa11f97aa6dd3766ee06ff162))
* test(cli): re-enable on-demand ctr scans (#490) (Salim
  Afiune)([ddddbc6](https://github.com/lacework/go-sdk/commit/ddddbc66aad78c4f69de7bb0104c4c5e80a058b9))

---

# v0.10.1

## Features

* feat(cli): apiv2 support for query validate (
  hazedav)([bcd81ce](https://github.com/lacework/go-sdk/commit/bcd81ce50c98b4b5204122301f4490a0a74c33c6))

## Bug Fixes

* fix(api): callbacks panic when response is nil (#486) (Salim
  Afiune)([8b7472a](https://github.com/lacework/go-sdk/commit/8b7472abe64cb3c837e34be7a2ec6703eacf0eb4))
* fix(cli): Output gcp comp list with --json flag (#482) (
  Darren)([9f95aa8](https://github.com/lacework/go-sdk/commit/9f95aa839dc16b6a64c0bbce9c66185b8fb3eef0))

## Other Changes

* test(cli): create VMs of supported platforms (#480) (Salim
  Afiune)([7f7c06e](https://github.com/lacework/go-sdk/commit/7f7c06e4013a0b379a685d21e78e03de3422ce56))

---

# v0.10.0

## Features

* feat(cli): New compliance google list command  (#466) (
  Darren)([f3aec8d](https://github.com/lacework/go-sdk/commit/f3aec8d5a08801fa517bd9aec405f1fb2cdafed1))
* feat(cli): caching mechanism (Salim Afiune
  Maya)([426c3eb](https://github.com/lacework/go-sdk/commit/426c3eb6ce569dc06295e30d2ca51f2c05213a4d))
* feat(api): LifecycleCallbacks (Salim Afiune
  Maya)([7b0c03e](https://github.com/lacework/go-sdk/commit/7b0c03e8a5bbb12af3dfbe7494fda5ea2030b9d5))
* feat: new internal/format go package (Salim Afiune
  Maya)([aa49474](https://github.com/lacework/go-sdk/commit/aa49474d6f9cdfa28d5e829e43e5d49c0fbaafe4))
* feat: install diskv go package (Salim Afiune
  Maya)([b0f7c43](https://github.com/lacework/go-sdk/commit/b0f7c43ac521a0830a109f6af44bfaf8fc3c00c8))
* feat(cli): apiv2 support for lacework policy (
  hazedav)([efaa50e](https://github.com/lacework/go-sdk/commit/efaa50e4a5023e65f10c9145355140092d4d7806))
* feat(cli): apiv2 support for Queries (
  hazedav)([e3b2cfa](https://github.com/lacework/go-sdk/commit/e3b2cfaffbf601bf9baa583e783cff699aabbe21))
* feat(cli): add --organization global flag (Salim Afiune
  Maya)([408a340](https://github.com/lacework/go-sdk/commit/408a34003d54c4e5098e8148836cc6bdeb0e3bb0))
* feat(api): Add AwsCtSqs Cloud Account Get/Update funcs (Salim Afiune
  Maya)([435686e](https://github.com/lacework/go-sdk/commit/435686ed61f9d3bf2691358d691bdef9594d3f93))
* feat(api): Add V2 CloudAccountsService (Salim Afiune
  Maya)([656b257](https://github.com/lacework/go-sdk/commit/656b2577ead5f84a8445de55f0aaee160ba79ddd))
* feat(api): WithOrgAccess Option to access org level (Salim Afiune
  Maya)([15f672e](https://github.com/lacework/go-sdk/commit/15f672e8765fa9102bc1139e30831b0ee9392bc9))
* feat(api): CopyClient() to duplicate an API client (Salim Afiune
  Maya)([084fe25](https://github.com/lacework/go-sdk/commit/084fe254d3edeed8e1229b4285b378a241662664))
* feat(api): propagate APIv2 error messages (#453) (Salim
  Afiune)([9a4c697](https://github.com/lacework/go-sdk/commit/9a4c6971f10e0b8945e780e09dccf82528c83391))

## Refactor

* refactor: migrate Agent Access Tokens to API v2  (#461) (Salim
  Afiune)([e08d947](https://github.com/lacework/go-sdk/commit/e08d94788004211f5b7d982539a9bad60d0c8653))
* refactor(cli): create AwsCtSqs CloudAccounts using APIv2 (Salim Afiune
  Maya)([b8b1605](https://github.com/lacework/go-sdk/commit/b8b1605eaf3e09c8213996430301fcc7022e7cc6))
* refactor(cli): improved LQL service func names (#438) (
  hazedav)([08f7ce4](https://github.com/lacework/go-sdk/commit/08f7ce4dbcd25e364c48eff92c00e07e814f69d4))

## Bug Fixes

* fix(api): add 'description' field for Host Vulnerability (#481) (Salim
  Afiune)([dba83cb](https://github.com/lacework/go-sdk/commit/dba83cbfbedbcb2542e91f6903899fe862ecad02))
* fix(cli): remove subaccount during configure (#474) (Salim
  Afiune)([f228f08](https://github.com/lacework/go-sdk/commit/f228f089fe8923065d559b177cd994fb802f4156))
* fix: check http status code before erase cache (Darren
  Murray)([7f1a29d](https://github.com/lacework/go-sdk/commit/7f1a29d26ce49582fa5b70a00ded3ece38c9a096))
* fix(cli): non-interactive should not run daily version check (#462) (Salim
  Afiune)([4bb7e9f](https://github.com/lacework/go-sdk/commit/4bb7e9feaf8baff84c2f94408f633abb69f98d16))

## Other Changes

* chore(deps): revert bump AlecAivazis/survey from 2.2.14 to 2.2.12 (#476) (Salim
  Afiune)([2d2176d](https://github.com/lacework/go-sdk/commit/2d2176d3a6a455c2de78c3c520f3a820561f7a8a))
* chore(deps): revert bump briandowns/spinner from 1.16.0 to 1.12.0 (#475) (Salim
  Afiune)([f2e19fb](https://github.com/lacework/go-sdk/commit/f2e19fbe4b64769c9dfb3fe50eb1dcfaa173d3d3))
* chore(deps): bump github.com/spf13/cobra from 1.1.3 to 1.2.1 (#469) (
  dependabot[bot])([6725a81](https://github.com/lacework/go-sdk/commit/6725a811d7815fc1047ead5490002069a54bc660))
* chore(deps): bump github.com/honeycombio/libhoney-go (#433) (
  dependabot[bot])([6a833a1](https://github.com/lacework/go-sdk/commit/6a833a16863f434c232b2638a6b90130976d3e22))
* chore(deps): bump github.com/spf13/viper from 1.7.1 to 1.8.1 (#459) (
  dependabot[bot])([7296b4a](https://github.com/lacework/go-sdk/commit/7296b4a4ca603baa0871c36aa2a079c078d568a0))
* chore(deps): bump github.com/briandowns/spinner from 1.12.0 to 1.16.0 (#452) (
  dependabot[bot])([06050a2](https://github.com/lacework/go-sdk/commit/06050a27fc8ad5a5554dfa2f041cc071d7d43253))
* chore(deps): bump go.uber.org/zap from 1.14.1 to 1.18.1 (#460) (
  dependabot[bot])([d311493](https://github.com/lacework/go-sdk/commit/d31149382f623616885f200d349c300ced8aad86))
* chore(deps): bump github.com/AlecAivazis/survey/v2 from 2.2.12 to 2.2.14 (#468) (
  dependabot[bot])([86b6a0f](https://github.com/lacework/go-sdk/commit/86b6a0f7e373d95d572243516183508a67897d4c))
* chore(api): rename cloud_account_aws.go -> cloud_accounts_aws_ct_sqs.go (Salim Afiune
  Maya)([dfd9ce2](https://github.com/lacework/go-sdk/commit/dfd9ce2e8569432e62c5538104016de2ab5b1059))
* ci: decouple integrations tests (#463) (
  Darren)([da7d429](https://github.com/lacework/go-sdk/commit/da7d429f53acaeaf2d1f30d8d0ae53d0f898ad65))
* ci: prepare-integrations pipeline (#457) (
  Darren)([d0e2436](https://github.com/lacework/go-sdk/commit/d0e2436e22d939d05dbe0cf7bf4fdc745c91dae1))
* test(cli): cache mechanism (Salim Afiune
  Maya)([6a0e9ca](https://github.com/lacework/go-sdk/commit/6a0e9ca6b59f50d015f99cbc6c0f2b9804f71e20))
* test(cli): fixing policy integration url (
  hazedav)([191a7ed](https://github.com/lacework/go-sdk/commit/191a7ed607632ef4b57d833415ea1a9e7cd76922))
* test(cli): apiv2 support for policy integration (
  hazedav)([68d324c](https://github.com/lacework/go-sdk/commit/68d324c6a48ca22952a580ee9fdcda305714f2b1))
* test(cli): add apiv2 query (#465) (
  hazedav)([bbd5b94](https://github.com/lacework/go-sdk/commit/bbd5b946d38213db2763be3e9a43bc643dc07e95))
* test(api): CloudAccount interface and CRUD tests (Salim Afiune
  Maya)([b0b103e](https://github.com/lacework/go-sdk/commit/b0b103e0116172b98be47f594b6ee7933555197c))
* test(cli): updated naming for test resources (
  hazedav)([4c3a389](https://github.com/lacework/go-sdk/commit/4c3a389cf2a684cda2c1fa1c61929bd29e4fb9e5))

---

# v0.9.1

## Bug Fixes

* fix(api): Parse LastEvaluationTime correctly in hostVulnListHostsResponse (#446) (
  Darren)([7918c4c](https://github.com/lacework/go-sdk/commit/7918c4c2cf816cad193f07a35150ce8de91481e0))
* fix(api): Enable FirstSeenTime field in hostVulnListHostsResponse (#445) (
  Darren)([195d44a](https://github.com/lacework/go-sdk/commit/195d44a98d019b293edaf8643c8aa876b996b945))

---

# v0.9.0

## Features

* feat(cli): new account list command (#444) (Salim
  Afiune)([735d34e](https://github.com/lacework/go-sdk/commit/735d34e346c82485b3c830f9dabcb0fa8aa9935a))
* feat(cli): LQLv2 support (#441) (
  hazedav)([8ad4abf](https://github.com/lacework/go-sdk/commit/8ad4abffdaadcf9cd8a3da0ffac8e012f4736672))

## Bug Fixes

* fix(cli): migration of Frankfurt organizational accounts (#447) (Salim
  Afiune)([73ed79a](https://github.com/lacework/go-sdk/commit/73ed79ad08280c47d9e54b179971a35dd79bfa19))
* fix(cli): adding more policy keys for json output (#442) (
  hazedav)([7a0d77f](https://github.com/lacework/go-sdk/commit/7a0d77fd028abad45cfce17d5c68f043f07c83a0))
* fix(cli): better error messaging for policy create (#443) (
  hazedav)([06174ff](https://github.com/lacework/go-sdk/commit/06174ff3de64f1ca25f7c098ec34a052dfbddbea))

## Other Changes

* chore(deps): bump github.com/fatih/color from 1.10.0 to 1.12.0 (#429) (
  dependabot[bot])([1b6bfe6](https://github.com/lacework/go-sdk/commit/1b6bfe692c3887725b9fc9104b6c216e7310397b))
* test(cli): adding lqlv2 test resource (#439) (
  hazedav)([df22443](https://github.com/lacework/go-sdk/commit/df22443556c1c709450d39c0b35f7a64228677ad))

---

# v0.8.0

## Features

* feat(cli): Add CSV rendering output for select commands (#424) (Matt
  Cadorette)([1bbe07c](https://github.com/lacework/go-sdk/commit/1bbe07c2055e61d6d34d0f4acb8bd9194358fcaf))
* feat(cli): Automate migration of v2 configuration (#431) (Salim
  Afiune)([2ac806a](https://github.com/lacework/go-sdk/commit/2ac806a91e8fea03eb6dec96e2ce8e8b3944b5e7))
* feat(pkg): to disseminate a domain URL (Salim Afiune
  Maya)([0770102](https://github.com/lacework/go-sdk/commit/07701029b5a057fd1b4da97d10683b27a381f365))
* feat(cli): v2 configure cmd undestands subaccounts (Salim Afiune
  Maya)([5e092f7](https://github.com/lacework/go-sdk/commit/5e092f70e601e081a02bb7ec18204c6373cdb6ec))
* feat: AWS Gov Cloud CloudTrail integration (#425) (
  Darren)([63097d5](https://github.com/lacework/go-sdk/commit/63097d55ae37616a61022f8cb0cba3cc5dcb9134))
* feat(cli): AWS GovCloud Config integration (#421) (
  Darren)([68d7087](https://github.com/lacework/go-sdk/commit/68d708713289c84606613a7261e015462c413299))
* feat(cli): lacework policy command (
  hazedav)([8c53e8e](https://github.com/lacework/go-sdk/commit/8c53e8eac2a06916dab4cf0062ea8ffc3781895e))
* feat: introducing first APIv2 endpoint UserProfile (#411) (Salim
  Afiune)([c183cce](https://github.com/lacework/go-sdk/commit/c183ccec205c473e461eae19d116f863a613ab9b))

## Refactor

* refactor(cli): remove shorthand of subaccount "u" (#436) (Salim
  Afiune)([2decc86](https://github.com/lacework/go-sdk/commit/2decc8610b77244c0f1b04370b34d368ea6b9206))
* refactor(cli): use new internal domain go package (Salim Afiune
  Maya)([1695993](https://github.com/lacework/go-sdk/commit/16959935dc6f876f82c970c9f1844ecd4068ccb7))
* refactor(cli): improve configure command (Salim Afiune
  Maya)([b14938b](https://github.com/lacework/go-sdk/commit/b14938be7ac4f22dc83d9659d9b062071f981cf9))

## Bug Fixes

* fix: Fix vuln container list-assessments json output (#430) (
  Darren)([02ff56d](https://github.com/lacework/go-sdk/commit/02ff56d37ab0fe0468aaeb1dd8afe28a53abbf0b))

## Other Changes

* ci: fix a bunch of things with our CI pipeline (#432) (Salim
  Afiune)([af603e1](https://github.com/lacework/go-sdk/commit/af603e145fe138a79fd21445da75087b2de93a49))
* test(integration): configure command for OrgAdmins (Salim Afiune
  Maya)([6504ad9](https://github.com/lacework/go-sdk/commit/6504ad991058eac4c0db2cd238d42580f9193e43))
* test: new DEBUG env variable for interactive shell (Salim Afiune
  Maya)([bd58489](https://github.com/lacework/go-sdk/commit/bd584894848e27593349c4ee9a779f4ee7fff9b2))
* test(integration): fix missed help message update (#419) (Salim
  Afiune)([f1fcbf3](https://github.com/lacework/go-sdk/commit/f1fcbf32b3dc571978b42293509f4def1d7c56cb))

---

# v0.7.1

## Bug Fixes

* fix(api): graceful degradation for HTTP 204 (#409) (
  hazedav)([44eb159](https://github.com/lacework/go-sdk/commit/44eb159af877e359be7a840cb72d4fd05b0cc718))

## Other Changes

* build(cli): Upgrade Go version 1.15 -> 1.16, Add support for darwin/arm64 (#408) (
  Darren)([c9d0a6f](https://github.com/lacework/go-sdk/commit/c9d0a6f8807eb12d3ebe86985e087450f095558e))
* ci: Disable failing test (#418) (
  Darren)([c5d7733](https://github.com/lacework/go-sdk/commit/c5d7733d589cacf2b44e58f1d55c7ddc6ef9d5c5))
* ci(cli): fix integration test (#417) (Salim
  Afiune)([b9ad070](https://github.com/lacework/go-sdk/commit/b9ad0706341d2583cd16cca70860695d55e9401b))
* ci: Update go deps (#413) (
  Darren)([8ee2393](https://github.com/lacework/go-sdk/commit/8ee239338741d5e624336347521514e86b9fa31e))
* test(cli): add container vuln scan NO polling (#414) (Salim
  Afiune)([793103f](https://github.com/lacework/go-sdk/commit/793103fb77170952e7e7f8498459f68d38c875c6))

---

# v0.7.0

## Features

* feat: new lwconfig Go package to handle Lacework config (#402) (Salim
  Afiune)([18d4864](https://github.com/lacework/go-sdk/commit/18d4864b0f5fd904ebc31e5f66ded697f068c209))
* feat(cli): lacework query natural/relative times (
  hazedav)([624235d](https://github.com/lacework/go-sdk/commit/624235dffb40ca6bfd853a736ed71f155441b5ec))
* feat(cli): lacework query crumbs (
  hazedav)([333bec3](https://github.com/lacework/go-sdk/commit/333bec3dcf9295a30967cb06d7daf8a1560154dc))

## Refactor

* refactor: negligible -> info (#393) (
  Darren)([be9d203](https://github.com/lacework/go-sdk/commit/be9d2034cfdcbd207ff220ce11f223dac54bd811))

## Bug Fixes

* fix(cli): Fix filters with json output for vuln host command (#388) (
  Darren)([282404f](https://github.com/lacework/go-sdk/commit/282404f01dbace5c88d30fc5268addf839c5b40e))
* fix(cli): Json output for vuln ctr adheres to filter flags (#383) (
  Darren)([cf7bb6d](https://github.com/lacework/go-sdk/commit/cf7bb6d29652fe75cbb7564793d76524e1ecd320))

## Other Changes

* style: Add gotestsum to Makefile (#397) (
  Ross)([6b98f91](https://github.com/lacework/go-sdk/commit/6b98f9136f545fd576452ac47b507406b98c2864))
* chore(deps): bump github.com/AlecAivazis/survey/v2 from 2.0.7 to 2.2.12 (#395) (
  dependabot[bot])([3b3fd4e](https://github.com/lacework/go-sdk/commit/3b3fd4e4b37753666e448aa5ba0ce2da824d122d))
* ci: fix release script for automatic releases (#405) (Salim
  Afiune)([3b657ee](https://github.com/lacework/go-sdk/commit/3b657ee1f61aa2148c5a4c88c040c6cfa0552329))
* ci: Add Codefresh badge (#404) (
  Darren)([f94d4ab](https://github.com/lacework/go-sdk/commit/f94d4abff5e17674a48777f958036d00612c0c35))
* ci: Remove Circleci config (#398) (
  Darren)([6b95e84](https://github.com/lacework/go-sdk/commit/6b95e849efb78def42a156eb026202049171ca46))
* ci: update ci tools (#403) (Salim
  Afiune)([580947b](https://github.com/lacework/go-sdk/commit/580947b066c28217733954fc00cadbc543b482bc))
* ci: disable on-deman ctr vulnerability tests (#400) (Salim
  Afiune)([bcaaaa7](https://github.com/lacework/go-sdk/commit/bcaaaa7116145e0c1eb8898fec512aa988442257))
* metric: normalize account inside Honeyvent (#407) (Salim
  Afiune)([3eb3b56](https://github.com/lacework/go-sdk/commit/3eb3b565a211905524f17ca81a91f327237c1d97))
* test(policy): adding policy test resource (#391) (
  hazedav)([74d8cfe](https://github.com/lacework/go-sdk/commit/74d8cfe6510c02332f223b1ca3bf447905561761))

---

# v0.6.0

## Features

* feat(cli): "lacework query" (hidden) command (
  hazedav)([4ae9ea8](https://github.com/lacework/go-sdk/commit/4ae9ea8bd077a5571cb66be4f75cef7f6c4be5af))

## Bug Fixes

* fix(cli): Hard coding polling mechanism to start after 64 seconds (#389) (Scott
  Ford)([f9644a6](https://github.com/lacework/go-sdk/commit/f9644a69b1ce9d4dde40b29fdc27df458c4d289f))

---

# v0.5.0

## Features

* feat: Filter vuln host list cves and show assessment by severity (#375) (
  Darren)([7e9313e](https://github.com/lacework/go-sdk/commit/7e9313e82f6a9af30560aac441f67bd6bee0a619))
* feat(cli): split list of Google project id and alias (#377) (Salim
  Afiune)([3f8dd94](https://github.com/lacework/go-sdk/commit/3f8dd9419350d60bd3bdd5b416a2e5d7e23af89a))

## Refactor

* refactor(cli): use official url for agent installs (#369) (Salim
  Afiune)([ca2d3f7](https://github.com/lacework/go-sdk/commit/ca2d3f7063dce713f6743ace9a4a07d5cc6c2119))

## Bug Fixes

* fix(fmt): undetected format problem (#385) (Salim
  Afiune)([8161d00](https://github.com/lacework/go-sdk/commit/8161d001edcff2e0d1bd05c2a500653d9c403904))
* fix(cli): split Azure subscriptions with parentheses (Salim Afiune
  Maya)([790759e](https://github.com/lacework/go-sdk/commit/790759e4c765f7dc68afebaf02be9dc7853eb3ad))
* fix(cli): get compliance reports from AWS accounts w/aliases (Salim Afiune
  Maya)([2dddcb1](https://github.com/lacework/go-sdk/commit/2dddcb10094703fbd9e59cfc2f03697b50f32dc0))
* fix(cli): Json result adheres to filters for cmd get-report (aws,gcp,azure) (#379) (
  Darren)([0da3db3](https://github.com/lacework/go-sdk/commit/0da3db3220e6c25b5a581a424ae9180d4e35e690))
* fix(cli): split GCP projects with parentheses (#380) (Salim
  Afiune)([339d13e](https://github.com/lacework/go-sdk/commit/339d13ef3af47cc4d52e6833ae308e0d96e5349f))
* fix(cli): match err message changed in server resp (#376) (Salim
  Afiune)([de7b137](https://github.com/lacework/go-sdk/commit/de7b13796ce7a5edfb173f2a868d281dce35ddc1))

## Other Changes

* chore: run goimports-reviser to order imports (#378) (Salim
  Afiune)([44b8e58](https://github.com/lacework/go-sdk/commit/44b8e589e3861b849e6092db2a248f98882f5918))
* chore(deps): bump github.com/briandowns/spinner from 1.10.0 to 1.12.0 (#374) (
  dependabot[bot])([e6058eb](https://github.com/lacework/go-sdk/commit/e6058ebd1ec3e82c6c7c8763ac6b1a9ef9a552db))
* chore(deps): bump github.com/mitchellh/mapstructure from 1.3.0 to 1.4.1 (#373) (
  dependabot[bot])([2f40036](https://github.com/lacework/go-sdk/commit/2f40036fb827181ac9e7fc37d71a66ff916bc775))
* chore(deps): bump github.com/kyokomi/emoji/v2 from 2.2.5 to 2.2.8 (#372) (
  dependabot[bot])([ec6ec28](https://github.com/lacework/go-sdk/commit/ec6ec281c41fd047521c8c6715b2b9a9e78e0a56))
* chore(deps): bump github.com/olekukonko/tablewriter from 0.0.4 to 0.0.5 (#370) (
  dependabot[bot])([bdfa365](https://github.com/lacework/go-sdk/commit/bdfa365f936f775704f7c74320e162f000066ebc))
* chore(deps): bump github.com/honeycombio/libhoney-go (#371) (
  dependabot[bot])([0494850](https://github.com/lacework/go-sdk/commit/0494850865c95d1f07348c17a069ea116e6fe121))
* test(integration): fix failing test from API changes (#368) (Salim
  Afiune)([57b33e3](https://github.com/lacework/go-sdk/commit/57b33e360cb207b431c077937b99255f7bef4352))

---

# v0.4.0

## Features

* feat: add new flags to fail host and container vuln scans (#352) (
  Darren)([1209d59](https://github.com/lacework/go-sdk/commit/1209d59138365a2aa86e9e0f8af66586542201db))

## Refactor

* refactor: Remove deprecated ListRange function (Darren
  Murray)([2a0019f](https://github.com/lacework/go-sdk/commit/2a0019f28dba6fcc6d470c2f6d54136581aac817))
* refactor: Remove deprecated vulnerability report command (Darren
  Murray)([d0de656](https://github.com/lacework/go-sdk/commit/d0de6564d7f6f87f682db902ed51e7aba0efbe87))
* refactor: Remove deprecated pdfName flag (Darren
  Murray)([94b8f67](https://github.com/lacework/go-sdk/commit/94b8f67aca806265de8ca8e6128fc3f137381efb))

## Bug Fixes

* fix(cli): mask ServiceNow integration password (#359) (Salim
  Afiune)([8e37697](https://github.com/lacework/go-sdk/commit/8e376976df0d429a403319c43e46ff1ad87dbc97))
* fix(api): update time format in API response (#360) (Salim
  Afiune)([a0045b4](https://github.com/lacework/go-sdk/commit/a0045b440a4164ff5ffd5c42d5a049471893b18f))

## Other Changes

* chore(deps): bump github.com/spf13/viper from 1.6.2 to 1.7.1 (#364) (
  dependabot[bot])([4f8f838](https://github.com/lacework/go-sdk/commit/4f8f838d3ed4352d5d660d75ce99a121ffc73cf0))
* chore(deps): bump github.com/stretchr/testify from 1.6.1 to 1.7.0 (#363) (
  dependabot[bot])([eed32e8](https://github.com/lacework/go-sdk/commit/eed32e8de351736c670594affda7d28002b7b778))
* chore(deps): bump github.com/fatih/color from 1.9.0 to 1.10.0 (#365) (
  dependabot[bot])([008253b](https://github.com/lacework/go-sdk/commit/008253bc78d27e2cb31434082284d6b64e0e51a4))
* chore(deps): bump github.com/spf13/cobra from 1.0.0 to 1.1.3 (#362) (
  dependabot[bot])([f3427ea](https://github.com/lacework/go-sdk/commit/f3427ea014da326218cb639f24f1ca58ae3e0515))
* chore: Validation script for git commit message conventions (#355) (
  Darren)([7fe9678](https://github.com/lacework/go-sdk/commit/7fe9678e854890a6a393df722007be38290630ff))
* ci: fix githooks and release scripts (#366) (Salim
  Afiune)([9036013](https://github.com/lacework/go-sdk/commit/90360137f823a4efa9651ade65ffaedfb4f61861))
* ci: switch over our CI to CodeFresh (#357) (Salim
  Afiune)([734acb7](https://github.com/lacework/go-sdk/commit/734acb7971c35ae79d6435ecc89d5a7fb4d493e7))
* ci: Fix prepare_test_resources script (Darren
  Murray)([e89beb6](https://github.com/lacework/go-sdk/commit/e89beb63d11941c67a97d91044fcbd6fe7ee2ad1))
* ci: Add setup_remote_docker to prepare-test-resources job (Darren
  Murray)([60628cd](https://github.com/lacework/go-sdk/commit/60628cd7731d0c15693a9b274f087a0b3eba590c))
* test(cli): lacework query command (#358) (
  hazedav)([d5473df](https://github.com/lacework/go-sdk/commit/d5473df7369b23d4e04dccb287de9b509b646342))
* test: removed Azure tenants from test account (Salim Afiune
  Maya)([e181ed1](https://github.com/lacework/go-sdk/commit/e181ed174fa9096e7c7a4a01ee95cc22951ba55d))

---

# v0.3.0

## Features

* feat(ux): improve container scan error messages (Salim Afiune
  Maya)([34e6ee7](https://github.com/lacework/go-sdk/commit/34e6ee7df622e47fd29d0845f0abe6158bb6e97d))
* feat(api): implement account information endpoint (#349) (Salim
  Afiune)([1392ecb](https://github.com/lacework/go-sdk/commit/1392ecb98c2908ab143ee630f0df4e0453a07767))
* feat(cli): new list-registries sub-command (Salim Afiune
  Maya)([7a828fe](https://github.com/lacework/go-sdk/commit/7a828fef4fe728873eb19bf689c892160ecb6a68))

## Other Changes

* chore: ran make prepare (Salim Afiune
  Maya)([662b220](https://github.com/lacework/go-sdk/commit/662b220a67f900fe003a0d275462589c1df0eb1b))
* ci: fix prepare-test-resources job (#348) (
  Darren)([a071c04](https://github.com/lacework/go-sdk/commit/a071c040a4c14eaa4c5b154361d94b5ffcb7d458))
* ci: open release pull request automatically (Salim Afiune
  Maya)([f227182](https://github.com/lacework/go-sdk/commit/f2271824616649f827fcffc5a15c99fdc809dd04))
* ci: automatic minor version bump (Salim Afiune
  Maya)([73e2cb9](https://github.com/lacework/go-sdk/commit/73e2cb9cb11f22909f7d698e10aa28a6a961061f))
* test: change target for container vuln scanning (#343) (
  Darren)([c348b01](https://github.com/lacework/go-sdk/commit/c348b0156f6846fb999202a29380639f0899c490))

---

# v0.2.23

## Features

* feat(cli): add Azure list-tenants sub-command (#341) (
  Darren)([960a8b7](https://github.com/lacework/go-sdk/commit/960a8b7f53a52febb66e89dce20f28211569a3fb))
* feat(cli): filter vulnerability assessments by severity (#338) (
  Darren)([07af9b1](https://github.com/lacework/go-sdk/commit/07af9b1091c8387d5d9e094a0fb5e9024749e4ef))

---

# v0.2.22

## Bug Fixes

* fix: implement both ECR auth methods (Salim Afiune
  Maya)([7af312c](https://github.com/lacework/go-sdk/commit/7af312c0308d3dad52f4e0264ee354800ed4d616))
* fix(api): type-o in host vulnerability status (#336) (Scott
  Ford)([85b271b](https://github.com/lacework/go-sdk/commit/85b271bdb35c33eec12df5591c530800688b6d90))

## Other Changes

* ci: run nightly integration tests on ARM-linux (#317) (Salim
  Afiune)([525b51d](https://github.com/lacework/go-sdk/commit/525b51d6866550a26ae7faefb2d8a8de99dd9323))
* test(cli): use a cli tag with vuln (Salim Afiune
  Maya)([eae52b8](https://github.com/lacework/go-sdk/commit/eae52b8d8350fd1e6166c036fa8b86e3462517de))

---

# v0.2.21

## Features

* feat(cli): load account from new UI API JSON file (#331) (Salim
  Afiune)([e841210](https://github.com/lacework/go-sdk/commit/e841210501eaa25620e102fad1e1ecf7a8c3bd3b))
* feat(cli): filtering flags for compliance report recommendations (#330) (
  Darren)([d04f09f](https://github.com/lacework/go-sdk/commit/d04f09f7cebaf278aaf22085cdb0f0a9f54b0b85))
* feat: Add support for ServiceNow Alert custom JSON template (#327) (
  Darren)([129bc28](https://github.com/lacework/go-sdk/commit/129bc2861f7ba7a807cb35c9d001b7ae575d97ab))

## Bug Fixes

* fix: Generate a new token upon a request with an expired token (#332) (
  Darren)([2bbc8b3](https://github.com/lacework/go-sdk/commit/2bbc8b38d7f5332174eb2f7266e300763054bdce))

## Documentation Updates

* docs: update go version batch in README (#329) (Salim
  Afiune)([161753e](https://github.com/lacework/go-sdk/commit/161753e35eac1990f69cc58d24c2b8afc5a5364d))
* docs: Add Homebrew installation to README (#328) (
  Darren)([e4ad780](https://github.com/lacework/go-sdk/commit/e4ad7803de81ed330b1a6eb9e5f61c0e3353c5e2))

## Other Changes

* ci: delete 'master' branch (#334) (Salim
  Afiune)([792e458](https://github.com/lacework/go-sdk/commit/792e458aed01012d2c665998a3e1e09f4a7facbb))

---

# v0.2.20

## Features

* feat(cli): New IBM QRadar alert channel (#325) (
  Darren)([0e9f6f5](https://github.com/lacework/go-sdk/commit/0e9f6f57e16e96ccfb84b9bd618e595378cf0fb3))
* feat(cli): New Relic Insights alert channel (#323) (
  Darren)([d7242b8](https://github.com/lacework/go-sdk/commit/d7242b84a525a4b9484324f2dfb49b172528302b))

## Documentation Updates

* doc(cli): update scan-pkg-manifest help to 10k pkgs (#324) (Salim
  Afiune)([0348800](https://github.com/lacework/go-sdk/commit/03488001d7c0aab308adf00ec6d8ac4a73aec78c))

---

# v0.2.19

## Features

* feat(cli): New VictorOps alert channel (#318) (
  Darren)([dfcd34a](https://github.com/lacework/go-sdk/commit/dfcd34adae5b84d697bd5f76831f749765f885fa))
* feat(cli): New CiscoWebex alert channel (#316) (
  Darren)([8e0071a](https://github.com/lacework/go-sdk/commit/8e0071a7d8a879d59c4545c8528e01287b90041e))
* feat: New Microsoft Teams alert channel (#315) (
  Darren)([e414226](https://github.com/lacework/go-sdk/commit/e41422624b5274a5d2b571fca4399ee2c9df3af3))

## Refactor

* refactor(cli): exponential retries polling scans (Salim Afiune
  Maya)([2bb881d](https://github.com/lacework/go-sdk/commit/2bb881d90fa6497003a89045acb41978398d5502))

## Bug Fixes

* fix(cli): match API client timeout with NGINX (#321) (Salim
  Afiune)([10b7a28](https://github.com/lacework/go-sdk/commit/10b7a28677af003e176de26b0dc1558c4837d1c4))

## Other Changes

* ci: increase integration test timeout to 30m (Salim Afiune
  Maya)([3081e3c](https://github.com/lacework/go-sdk/commit/3081e3cc6264e85b7fd80019a36d89dda1d9e5c9))
* test: change ctr vuln scan tag (Salim Afiune
  Maya)([c15bd1b](https://github.com/lacework/go-sdk/commit/c15bd1b2e0454ff61f6a7549f62c6bc791190079))
* test(cli): disable failing tests (RAIN-15300) (#320) (Salim
  Afiune)([e2afb31](https://github.com/lacework/go-sdk/commit/e2afb31277e8a2781c5f8e55ebf352992cd6d2b0))

---

# v0.2.18

## Features

* feat(cli): new Datadog alert channel  (#313) (
  Darren)([8298022](https://github.com/lacework/go-sdk/commit/8298022367d99ccc0237e456c8ccd6b45885e1d7))

## Bug Fixes

* fix(cli): avoid daily update check during install (Salim Afiune
  Maya)([2bc94c0](https://github.com/lacework/go-sdk/commit/2bc94c002392cb900ce6ca47574eafa3f4e55959))

## Other Changes

* test: fix intermittent events test (#312) (Salim
  Afiune)([d69983b](https://github.com/lacework/go-sdk/commit/d69983b63c9c771f533d72325bc58cd63e5d3a66))
* test: fix intermittent events test (Salim Afiune
  Maya)([15c371c](https://github.com/lacework/go-sdk/commit/15c371cc59ae9b2cbbe25588ff20574c1272efbe))

---

# v0.2.17

## Bug Fixes

* fix: Json mapping for Snow Username (#307) (
  Darren)([1ef8d99](https://github.com/lacework/go-sdk/commit/1ef8d9912dd66522712bfa4a15d4ab728d409e65))

---

# v0.2.16

## Features

* feat(cli): New Service Now alert channel (#303) (
  Darren)([512f2d9](https://github.com/lacework/go-sdk/commit/512f2d9c38d8124904dbdd661d9ab3b8441fc86d))

## Refactor

* refactor: Change input method for private_key field (#305) (
  Darren)([e56cdc6](https://github.com/lacework/go-sdk/commit/e56cdc68d74074f5ee904712aa56779dc9a0e1ed))

## Bug Fixes

* fix: Use select for issue grouping (#304) (
  Darren)([799d9c3](https://github.com/lacework/go-sdk/commit/799d9c34723d4af03f9ed811028880181f1757f5))
* fix: Add issue_grouping field to gcp pub sub (#301) (
  Darren)([1a66d2c](https://github.com/lacework/go-sdk/commit/1a66d2cc20a33674348705fbd3e552bf9222c787))
* fix(cli): install.sh should try curl and wget (Salim Afiune
  Maya)([f6b0bd7](https://github.com/lacework/go-sdk/commit/f6b0bd725992961e8c8b02ba3090164f31056388))
* fix(cli): install.sh should respect target override (Salim Afiune
  Maya)([4164f58](https://github.com/lacework/go-sdk/commit/4164f5872dd381471ad8608abd9c2fa821dac524))

## Other Changes

* chore(cli): install.sh print exitcodes for debugging (Salim Afiune
  Maya)([5e66c11](https://github.com/lacework/go-sdk/commit/5e66c11264f0b9ba012be20370b0b65cc114fe4c))
* ci: trigger homebrew update script (#299) (
  Darren)([9247cb1](https://github.com/lacework/go-sdk/commit/9247cb12b62582596939be55861f4d6c08bf8690))

---

# v0.2.15

## Features

* feat: add telemetry to detect Homebrew installations (#297) (
  Darren)([fa81abc](https://github.com/lacework/go-sdk/commit/fa81abc1c044cb362cd29608c65fc820d0c8a706))
* feat(cli): New Gcp PubSub alert channel (#294) (
  Darren)([08a3e61](https://github.com/lacework/go-sdk/commit/08a3e61469a1bb35675d63070474bad3e7988ad4))

---

# v0.2.14

## Features

* feat(cli): support Homebrew upgrade command (#291) (
  Darren)([bedfa5d](https://github.com/lacework/go-sdk/commit/bedfa5d2895fa7aabe1de6ef2902eb260133c3fd))
* feat(cli): Add Splunk alert channel (#289) (
  Darren)([04679a5](https://github.com/lacework/go-sdk/commit/04679a51bf35469f4986e693ed19a75cc36dbbb2))
* feat(cli): add account check to catch http(s):// (#288) (Salim
  Afiune)([3d770a1](https://github.com/lacework/go-sdk/commit/3d770a1171814e33de7f6420d1637ad8c03f30c8))

## Bug Fixes

* fix(cli): skip daily version check for version cmd (#290) (Salim
  Afiune)([5c9f4ca](https://github.com/lacework/go-sdk/commit/5c9f4cab61620636e0962e0fbe4edce97c41e8dc))

---

# v0.2.13

## Features

* feat(cli): support manifest bigger than 1k packages (Salim Afiune
  Maya)([eebddb9](https://github.com/lacework/go-sdk/commit/eebddb9325ede76ffa1853d00508da54cb5b9678))
* feat(cli): gen-pkg-manifest detect running kernel (Salim Afiune
  Maya)([9151be1](https://github.com/lacework/go-sdk/commit/9151be15a05b48f3d7456571cd75411f2ba7ddb9))

## Refactor

* refactor: simplify removeEpochFromPkgVersion func (Salim
  Afiune)([04aba5b](https://github.com/lacework/go-sdk/commit/04aba5bda340283f86d93496f01e0089a500468d))

## Bug Fixes

* fix(cli): ensure api client has valid auth token (Salim Afiune
  Maya)([056eda5](https://github.com/lacework/go-sdk/commit/056eda5cb7bde11e2334b6f38bd338afe111ade9))

## Other Changes

* ci: generate code coverage in HTML format (Salim Afiune
  Maya)([a58b58a](https://github.com/lacework/go-sdk/commit/a58b58a6477ec8d12c06bff3672093aef826c1f1))
* ci: add 'metric' as a valid commit message (Salim Afiune
  Maya)([dd7b601](https://github.com/lacework/go-sdk/commit/dd7b6010969d1f99055b7dbc9442498fa9f002cf))
* ci: fix slack notifications team alias ⭐ (Salim Afiune
  Maya)([ca51f92](https://github.com/lacework/go-sdk/commit/ca51f92693a48f113dd7661d9ef03eef7c26a17a))
* metric(cli): detect feature split_pkg_manifest (Salim Afiune
  Maya)([fdb9f4a](https://github.com/lacework/go-sdk/commit/fdb9f4a1c1eae2b9a44ea846fae413a93f073ca9))
* metric(cli): detect feature gen_pkg_manifest (Salim Afiune
  Maya)([78905bb](https://github.com/lacework/go-sdk/commit/78905bb73f398bf26a6e297e3929e5993e4965dc))

---

# v0.2.12

## Features

* feat(cli): add telemetry (#278) (Salim
  Afiune)([5aeec3c](https://github.com/lacework/go-sdk/commit/5aeec3c51184fc7e43e1c9dc413d256c98b8c516))
* feat(cli): pull latest agent version from S3 (Salim Afiune
  Maya)([63cf1ab](https://github.com/lacework/go-sdk/commit/63cf1ab82933600189904abe0b25958769a42ec9))
* feat: add --force to agent install (Salim Afiune
  Maya)([6de4775](https://github.com/lacework/go-sdk/commit/6de47756973f3396b9d3f5d6e044db3308e1700a))
* feat: verify if agent is installed on remote host (Salim Afiune
  Maya)([252b9a6](https://github.com/lacework/go-sdk/commit/252b9a602781a68ee88d1d0c9e14ee290c310a79))
* feat(cli): check for known hosts and allow custom callbacks (Salim Afiune
  Maya)([ebedf22](https://github.com/lacework/go-sdk/commit/ebedf221f4a1569080aeaf8de1441661845d22b2))
* feat: add AWS S3 alert channel integration (#273) (
  Darren)([383de18](https://github.com/lacework/go-sdk/commit/383de18bedfa1d85eb140f5b82ecb2c69ba231be))
* feat(cli): enable agent install command (Salim Afiune
  Maya)([f13d58a](https://github.com/lacework/go-sdk/commit/f13d58a2bbedf7772ddd63330a4cb813f926f541))

## Refactor

* refactor: verify host connectivity before select token (Salim Afiune
  Maya)([829cf82](https://github.com/lacework/go-sdk/commit/829cf821d457e5178c13e3d98bd9f31c60be3ded))
* refactor(api): remove automatic report trigger (#271) (Salim
  Afiune)([18e624f](https://github.com/lacework/go-sdk/commit/18e624f74e68fddc2f180e5e608353a824bac9b7))

## Bug Fixes

* fix(cli): propagate errors from install.sh (#277) (Salim
  Afiune)([296be65](https://github.com/lacework/go-sdk/commit/296be658d106ad84cf9a4a3ced1d4f6122ce4db8))
* fix(cli): avoid showing unnamed tokens (Salim Afiune
  Maya)([7545444](https://github.com/lacework/go-sdk/commit/754544441972f73a55181a4255453f6f911f81d0))

## Documentation Updates

* docs: update agent install use (Salim Afiune
  Maya)([62195c1](https://github.com/lacework/go-sdk/commit/62195c1a2b429b02120a8d797e0debaa448016e8))

## Other Changes

* chore: update long desc of agent list cmd (Salim Afiune
  Maya)([8a24914](https://github.com/lacework/go-sdk/commit/8a2491456d361d22de4760a79abfbbb0dcc51559))
* build: stop publishing containers to old docker repo (Salim Afiune
  Maya)([ea23a30](https://github.com/lacework/go-sdk/commit/ea23a3085e4c8ef35acc4fe06d3ba972be4d932a))
* ci: send slack notifications to team alias ⭐ (Salim Afiune
  Maya)([5e4c0e6](https://github.com/lacework/go-sdk/commit/5e4c0e69824ef00289e4d86adecf48209709bb59))
* ci: fix mv on non exisitent directory (#272) (
  Darren)([4f101cf](https://github.com/lacework/go-sdk/commit/4f101cfe8c8aeff5981264e99bdb411b548e02e9))
* test(cli): increase agent install test coverage (#276) (Salim
  Afiune)([da5b4ae](https://github.com/lacework/go-sdk/commit/da5b4aea9730c55c10d541c976dcb7ccf16aca28))
* test: fix lwrunner tests (Salim Afiune
  Maya)([23587cd](https://github.com/lacework/go-sdk/commit/23587cdd98c694e65a8f0791c269817ce7252d4c))

---

# v0.2.11

## Features

* feat(cli): daily version check (#269) (Salim
  Afiune)([5c15eef](https://github.com/lacework/go-sdk/commit/5c15eef84f428ec0534954babb28a3db92d5a7c5))
* feat(api): add Webhook integration (#267) (
  Darren)([f32572e](https://github.com/lacework/go-sdk/commit/f32572ecdadd5c179227cd228bf1fdd7cf618763))

## Refactor

* refactor(cli): abstract rendering tables (human-readable) (#263) (Salim
  Afiune)([8a10b4c](https://github.com/lacework/go-sdk/commit/8a10b4cf10de03d9b4c0409e495fdd7118974b92))

## Bug Fixes

* fix(cli): render account mapping file correctly (#266) (Salim
  Afiune)([4c327d7](https://github.com/lacework/go-sdk/commit/4c327d7e6081d0f7726a1bc007b1b736a106933f))
* fix(api): new request body for lql service (#260) (Salim
  Afiune)([4e2b439](https://github.com/lacework/go-sdk/commit/4e2b439ff394d632cd6ebf214da376050da46812))
* fix(api): avoid updating AgentTokenResponse.Props (#259) (Salim
  Afiune)([c3fe8bc](https://github.com/lacework/go-sdk/commit/c3fe8bcc41efd995f756f97a5ffca8bb961e89e4))

## Documentation Updates

* docs: update READMEs and _examples/ (#268) (Salim
  Afiune)([3791da0](https://github.com/lacework/go-sdk/commit/3791da01005335c34852446c57eb99e51a6d3ce1))

## Other Changes

* build: upgrade Go version to 1.15 (#265) (Salim
  Afiune)([06d41f5](https://github.com/lacework/go-sdk/commit/06d41f56add71f8369ffae68ea7ba5d738eb4d5b))
* ci: update hostname from our test machine (#262) (Salim
  Afiune)([beb289e](https://github.com/lacework/go-sdk/commit/beb289e732c177e2f3d062d61dee5dd9f1593ce9))

---

# v0.2.10

## Features

* feat(cli): new agent access token command (#256) (Salim
  Afiune)([7f8ba11](https://github.com/lacework/go-sdk/commit/7f8ba113b38ecd768f61e54ba712badf6596a587))
* feat(compliance): new aws list-accounts command (Salim Afiune
  Maya)([705f2eb](https://github.com/lacework/go-sdk/commit/705f2ebf9f1b9b5af2eb745c86498fe31c01e174))

## Refactor

* refactor: account mapping file for consolidated CT (#252) (Salim
  Afiune)([402a363](https://github.com/lacework/go-sdk/commit/402a3634765ef8c6f1f65d1be13da2ad34cf2960))

## Bug Fixes

* fix(install.sh): avoid logging with 'info' cmd (#254) (Salim
  Afiune)([df5f8cf](https://github.com/lacework/go-sdk/commit/df5f8cfbc7228ff9bff25e6e22a2ab68acd47fa4))
* fix: false positive results in pkg manifest scan (#255) (Salim
  Afiune)([a6d6cda](https://github.com/lacework/go-sdk/commit/a6d6cda9f36b38f8b653bd01ef258bd431611908))
* fix(databox): remove hardcoded LW account (Salim Afiune
  Maya)([c806157](https://github.com/lacework/go-sdk/commit/c80615749827c12dbfef5e1c76bf5857cd3dae7a))

---

# v0.2.9

## Features

* feat(api): enable account mapping file for CT int (#250) (Salim
  Afiune)([cb99f61](https://github.com/lacework/go-sdk/commit/cb99f61f5da717911d72ff2379a98ce4b7f6dd61))

## Refactor

* refactor(api): better error check handler (#247) (Salim
  Afiune)([b363347](https://github.com/lacework/go-sdk/commit/b363347409f242e6ad1a46e1104884a559877ada))

## Other Changes

* ci: set container tag to debian-10 that has vulns (#248) (Salim
  Afiune)([323b91e](https://github.com/lacework/go-sdk/commit/323b91e38c6616a064501c602d58f2cb0a2572f2))
* ci: dogfooding Lacework Orb html parameter (Salim Afiune
  Maya)([464d34d](https://github.com/lacework/go-sdk/commit/464d34db0d53fc7e0e9f4b6cbd7663eae963ad46))
* ci: remove slack alert for win systems (Salim Afiune
  Maya)([b6b5b45](https://github.com/lacework/go-sdk/commit/b6b5b458de4e69fd0d176522a47a168382a03397))

---

# v0.2.8

## Bug Fixes

* fix(cli): generate html for scan commands (Salim Afiune
  Maya)([6846ffd](https://github.com/lacework/go-sdk/commit/6846ffd4a6f5d76ca700652e0a5cd7adda6c9cdc))

## Other Changes

* ci: improve release notes and changelog generation (Salim Afiune
  Maya)([af22a7a](https://github.com/lacework/go-sdk/commit/af22a7a5518adc4b5f1c6fc21c23934230dd0705))
* ci: avoid release.sh to update version multiple times (Salim Afiune
  Maya)([d72149b](https://github.com/lacework/go-sdk/commit/d72149b50b197a9ddc7528e18bc88c1836f4847f))
* test(cli): HTML for container vulnerability (Salim Afiune
  Maya)([fee8505](https://github.com/lacework/go-sdk/commit/fee85056a22b460b925f60083cbb2610c8552c55))

---

# v0.2.7

## Features

* feat(cli): enable html copy to clipboard icons (Salim Afiune
  Maya)([ec2d1fa](https://github.com/lacework/go-sdk/commit/ec2d1fa5b796e2b51e49a850deff01d2f64ded18))
* feat(cli): HTML format for vulnerability assessments (Salim Afiune
  Maya)([00c2f43](https://github.com/lacework/go-sdk/commit/00c2f43613e554afd8ed283cbc12eb0b8eed0179))
* feat(cli): add ARM support (#236) (Salim
  Afiune)([821b8e6](https://github.com/lacework/go-sdk/commit/821b8e699e61eefda7d287a71b08ef26382a4ad7))

## Bug Fixes

* fix(cli): remove html column sort icons (Salim Afiune
  Maya)([dc4c0f6](https://github.com/lacework/go-sdk/commit/dc4c0f64055bfaada503b8f4f21ceda707bc5e55))

## Other Changes

* ci(fix) Update CI test node (#233) (Scott
  Ford)([ddbf86e](https://github.com/lacework/go-sdk/commit/ddbf86e8fbf9053af43bab9d57c04645383e529e))

---

# v0.2.6

## Features

* feat(api): trigger initial report automatically (#230) (Salim
  Afiune)([1e24a22](https://github.com/lacework/go-sdk/commit/1e24a229d2f2c54b81809ccb156a9d4283962c32))

## Documentation Updates

* docs(cli): disable timestamp for automatic docs (#229) (Salim
  Afiune)([f4d7841](https://github.com/lacework/go-sdk/commit/f4d78417c307c38507995892999cb85e7be74cf2))

---

# v0.2.5

## Bug Fixes

* fix(cli): add epoch to package manifest (Salim Afiune
  Maya)([17da487](https://github.com/lacework/go-sdk/commit/17da48755062265245d98ba6f4a330ae65fcdb6b))

## Other Changes

* chore(ci): make GH org a readonly parameter (Salim Afiune
  Maya)([b4f5f6d](https://github.com/lacework/go-sdk/commit/b4f5f6d5ba5a644a6198445bd820d68bf243907d))
* chore(cli): update pkg-manifest message for 0 vuln (Salim Afiune
  Maya)([5029dc8](https://github.com/lacework/go-sdk/commit/5029dc82aa51f260e84cd476acd6c64cab7f063a))

---

# v0.2.4

## Features

* feat(cli): programatic access to profile data (#225) (Salim
  Afiune)([ab7ce7c](https://github.com/lacework/go-sdk/commit/ab7ce7cfe8e94053ca6bf8d32d929c5e748496e4))
* feat(cli): allow custom installation directory -d 📁 (#223) (Salim
  Afiune)([ee9e686](https://github.com/lacework/go-sdk/commit/ee9e686c46029b32e711f9534ecd7755926ec22b))

## Documentation Updates

* docs: automatically generate cli docs (#224) (Salim
  Afiune)([5b91e1e](https://github.com/lacework/go-sdk/commit/5b91e1e788128dd3cddf457bce565749c73eddae))

## Other Changes

* chore: add badges to README.md (#222) (Salim
  Afiune)([db7235d](https://github.com/lacework/go-sdk/commit/db7235d20e7af012cb8e8f3041a02728d4f28719))

---

# v0.2.3

## Features

* feat(cli): add scan-pkg-manifest summary 📈 (#220) (Salim
  Afiune)([9b009c3](https://github.com/lacework/go-sdk/commit/9b009c3e98a69d294d424c2b912b1aadb675ee98))
* feat(ux): generate package-manifest command (#217) (Salim
  Afiune)([0c842ab](https://github.com/lacework/go-sdk/commit/0c842ab15c30b3f754a379ecd2aea014c367bae7))

## Refactor

* refactor: remove 'apk' as supported pkg manager (Salim Afiune
  Maya)([4165783](https://github.com/lacework/go-sdk/commit/41657839f06ea9b8eae85119451c77e632ec99bb))

## Other Changes

* chore(ci): update lacework circleci orb (Salim Afiune
  Maya)([3952c66](https://github.com/lacework/go-sdk/commit/3952c66f47dbb0024b3fef35f3f39087fa76844e))

---

# v0.2.2

## Features

* feat(lql): --file flag to load LQL query from disk (Salim Afiune
  Maya)([4804319](https://github.com/lacework/go-sdk/commit/4804319a0c26211119c10eb3dc4d889b3da7e227))
* feat(cli): --file to pass a package manifest file (Salim Afiune
  Maya)([75680d8](https://github.com/lacework/go-sdk/commit/75680d8d9469d8679b17c46979439340d8869da9))
* feat: human-readable output for scan-pkg-manifest (Salim Afiune
  Maya)([783f550](https://github.com/lacework/go-sdk/commit/783f55015c1e6a1071927e19810266376ecbe082))
* feat(lql): improve running queries (Salim Afiune
  Maya)([61c5ee5](https://github.com/lacework/go-sdk/commit/61c5ee51aac65626aff4f81ebceb96633865d2f7))

## Bug Fixes

* fix(ci): remove slack notification for windows (#214) (Salim
  Afiune)([a2c5124](https://github.com/lacework/go-sdk/commit/a2c51242c08c1683cfb9c80c832be2559058f957))

## Other Changes

* ci(slack): notify pipeline failures (#213) (Salim
  Afiune)([85ad396](https://github.com/lacework/go-sdk/commit/85ad396f6cb049ab246ff36fa2f29d46fab6459d))

---

# v0.2.1

## Features

* feat(ctr): use new lacework/lacework-cli repository (#206) (Salim
  Afiune)([fa1e268](https://github.com/lacework/go-sdk/commit/fa1e2682422f03288c53350d4fc6691bea6869c5))
* feat: add DockerV2, ECR and GCR container registries (#205) (Salim
  Afiune)([18a8c8b](https://github.com/lacework/go-sdk/commit/18a8c8b60ef6c869bcb3a72b870d4bcfd66ee794))
* feat: add decoder for jira custom_template_file (#201) (Salim
  Afiune)([2630ab5](https://github.com/lacework/go-sdk/commit/2630ab5fb746a8a3b4995734a79489554ff4f682))
* feat(cli): ask for JIRA Custom Template file 🚨 (Salim Afiune
  Maya)([5a4eb17](https://github.com/lacework/go-sdk/commit/5a4eb173b26eaff1845a33b586f9b87bbb59f449))
* feat(api): encode custom_template_file for Jira int (Salim Afiune
  Maya)([887ca15](https://github.com/lacework/go-sdk/commit/887ca157f8484cfc971b83f1d4e65f0fa2f10382))

## Documentation Updates

* docs(typo) fix spelling of visualize for compliance help command (#204) (Scott
  Ford)([75e0348](https://github.com/lacework/go-sdk/commit/75e03488a01704837498f630c4c0323d6a3ee6ef))

## Other Changes

* chore(api): remove MinAlertSeverity field from examples/ (Salim Afiune
  Maya)([274b8e9](https://github.com/lacework/go-sdk/commit/274b8e927b62d395074ba79628ba4ad8abdd5905))
* ci(cli): fix event time range test (Salim Afiune
  Maya)([9c2336b](https://github.com/lacework/go-sdk/commit/9c2336b9d7b2f96e3b8f8e22d57ba2d4fa77583b))

---

# v0.2.0

## Features

* feat(cli): new event open command (#197) (Salim
  Afiune)([42e0309](https://github.com/lacework/go-sdk/commit/42e03096cf387a55329275c22a787ccf239c1baa))
* feat(cli): filter events by severity (Salim Afiune
  Maya)([2d8fdf4](https://github.com/lacework/go-sdk/commit/2d8fdf46b391562205d036a8f866b4e940377f9c))
* feat(cli): list events from a number of days (Salim Afiune
  Maya)([0474765](https://github.com/lacework/go-sdk/commit/047476548e6b86dcd249c8f37b0cfb65a49a401d))
* feat(cli): allow users to pass only --start flag (Salim Afiune
  Maya)([547dc1d](https://github.com/lacework/go-sdk/commit/547dc1d3a8db23e9d9b411e045b6bbce6b99e161))
* feat(cli): filter assessments for specific repos (Salim Afiune
  Maya)([6482d8e](https://github.com/lacework/go-sdk/commit/6482d8ea6ad712077fc595011cbdfee0715c04bc))
* feat(cli): --active & --fixable flags to container vuln (Salim Afiune
  Maya)([9f027b9](https://github.com/lacework/go-sdk/commit/9f027b9b56c2b4c110281246971988881f8f1164))
* feat(cli): --active & --fixable flags to host vuln (Salim Afiune
  Maya)([27f5197](https://github.com/lacework/go-sdk/commit/27f5197c17488a9575a8ba47f17293590a8cdbbf))
* feat(cli): add emoji support for windows (Salim Afiune
  Maya)([0762814](https://github.com/lacework/go-sdk/commit/07628145c9e034bc8492d9e833bf9cef962996da))
* feat(cli): add an emoji Go package for 🍺 🍕 🌮 (Salim Afiune
  Maya)([cafb8d8](https://github.com/lacework/go-sdk/commit/cafb8d8cf721e7d3259f7de5f06613d3136c28f0))
* feat(cli): order vulnerabilities by total of hosts (Salim Afiune
  Maya)([5cfe695](https://github.com/lacework/go-sdk/commit/5cfe69538cb1c869909e4b4f321eeab7c3ac1b19))
* feat(cli): new vulnerability list-assessments command (Salim Afiune
  Maya)([7e7191a](https://github.com/lacework/go-sdk/commit/7e7191ab1aa4b765081c91573df307d5c9113f9c))

## Refactor

* refactor(cli): container and host vulnerability cmds (Salim Afiune
  Maya)([c5c0117](https://github.com/lacework/go-sdk/commit/c5c0117492eec958159b13df36b738af48f5a5e0))
* refactor: host vulnerability feature (Salim Afiune
  Maya)([5e9f770](https://github.com/lacework/go-sdk/commit/5e9f7700acd422f5bf0b79d3faf58ffc6ed0034b))
* refactor: container vulnerability feature (Salim Afiune
  Maya)([bdaf126](https://github.com/lacework/go-sdk/commit/bdaf12641851b3a3bb514617ca3ae61e062bbb07))

## Performance Improvements

* perf(cli): retry polling on-demand container scan statuses (Salim Afiune
  Maya)([d14ea35](https://github.com/lacework/go-sdk/commit/d14ea3598c2f5d4ea795f3930c0e6b48698e9777))

## Other Changes

* chore(cli): update help messages (Salim Afiune
  Maya)([f1c164c](https://github.com/lacework/go-sdk/commit/f1c164c14703e6dc1faecbd566ff7be3aae822ae))
* chore(cli): consistent help message for vuln cmds (Salim Afiune
  Maya)([f796c58](https://github.com/lacework/go-sdk/commit/f796c5835f91c5224701e60f8236fc55e663b83e))
* chore(cli): leave breadcrumbs for host vuln cmds (Salim Afiune
  Maya)([45d8427](https://github.com/lacework/go-sdk/commit/45d8427554a9a74f40f3e97c2e0f8c0251a8450f))
* ci(integration): run full tests on windows (#190) (Salim
  Afiune)([c5c8cf4](https://github.com/lacework/go-sdk/commit/c5c8cf4c80a2fcb40e84dcefbec4f733c5d8bc52))
* test(integration): add host vulnerability tests (Salim Afiune
  Maya)([a5cb795](https://github.com/lacework/go-sdk/commit/a5cb7951832c4c95c64b24c80f73e06293920283))
* test(integration): add container vulnerability tests (Salim Afiune
  Maya)([9b2c49d](https://github.com/lacework/go-sdk/commit/9b2c49d88ca962274e145028eaebb58f88ff417b))

---

# v0.1.24

## Features

* feat(cli): better ux in account validation (#187) (Salim
  Afiune)([cdd045a](https://github.com/lacework/go-sdk/commit/cdd045a830dcdc788daf77d9ea558ba4d296e003))
* feat(cli): new access-tokens command (#184) (Salim
  Afiune)([ee338c4](https://github.com/lacework/go-sdk/commit/ee338c4afb057bf4ea578a8d0ddb48b2d39b34d3))
* feat(cli): Create Jira Alert Channels 🚨 (Salim Afiune
  Maya)([6ca8cef](https://github.com/lacework/go-sdk/commit/6ca8ceffce1c17f3f84634da5514c059da952ca1))
* feat(api): add Jira alert channel integrations (Salim Afiune
  Maya)([0cdb2a4](https://github.com/lacework/go-sdk/commit/0cdb2a46d820f249c0fe918320303b1061e0f5ed))

## Refactor

* refactor: remove legacy field min_alert_severity (#186) (Salim
  Afiune)([54ca38c](https://github.com/lacework/go-sdk/commit/54ca38c8c509d800e2bddca5435529f5d0b60643))

## Bug Fixes

* fix(cli): display integration update by/update time (Salim Afiune
  Maya)([7060078](https://github.com/lacework/go-sdk/commit/7060078d3f8a09a82a4efaf98c4cb15f4856f753))

---

# v0.1.23

## Refactor

* refactor(cli): replace '--pdf-file' for '--pdf' (#180) (Salim
  Afiune)([80bbce6](https://github.com/lacework/go-sdk/commit/80bbce636cac49fe315118add45252bd8ee4bf6a))

## Bug Fixes

* fix(cli): missing integration details (#181) (Salim
  Afiune)([40355d3](https://github.com/lacework/go-sdk/commit/40355d3877c2674268c38bb5cc81a698dd115166))
* fix(cli): error showing non-existing integration (#178) (Salim
  Afiune)([252072f](https://github.com/lacework/go-sdk/commit/252072faa60aaac06fb7bbf2dd7ca82fa71d2b09))

## Other Changes

* ci: build statically linked binaries (Salim Afiune
  Maya)([43f6f80](https://github.com/lacework/go-sdk/commit/43f6f804ffac3f8e326dc31f4196808f39bc035d))
* ci(integration): add windows support (Salim Afiune
  Maya)([46632e7](https://github.com/lacework/go-sdk/commit/46632e72e0ab9ee45d690605e4c52efb1a8cf391))

---

# v0.1.22

## Features

* feat(cli): Create PagerDuty Alert Channels 🚨 (#174) (Salim
  Afiune)([5cc424e](https://github.com/lacework/go-sdk/commit/5cc424e21598482f817288037c8f8e54397c13bd))
* feat(api): add PagerDuty alert channel integrations (#173) (Salim
  Afiune)([f46316c](https://github.com/lacework/go-sdk/commit/f46316c7f4150ccf99646640a12d801cb407134b))
* feat(cli): Create AWS CloudWatch Alert Channels 🚨 (Salim Afiune
  Maya)([201b59b](https://github.com/lacework/go-sdk/commit/201b59be0a97d661916ff401da0be903fee06f2f))
* feat(api): add AWS CloudWatch Alert Channels Int (Salim Afiune
  Maya)([d9a11ec](https://github.com/lacework/go-sdk/commit/d9a11ec5c242b09e19338c6b8a5a39ddf6ad368d))
* feat(api): enum AlertLevel for alert severity levels (Salim Afiune
  Maya)([d3bf436](https://github.com/lacework/go-sdk/commit/d3bf436933a794b6bbcc733da724159a9dc79a95))
* feat(api): get/update container registry integrations (#168) (Salim
  Afiune)([a072c46](https://github.com/lacework/go-sdk/commit/a072c46aff03e619fbef03488ba5b65730264b91))

## Refactor

* refactor(api): AlertChannel prefix in funcs/structs (Salim Afiune
  Maya)([b0429ef](https://github.com/lacework/go-sdk/commit/b0429efd0efa56ec9ccbe338a37a6e6ae2dc3bc5))
* refactor(api): use AlertLevel enum for Slack Alerts (Salim Afiune
  Maya)([4b5acf9](https://github.com/lacework/go-sdk/commit/4b5acf989fda4c052c3dc6b0206db866aa57f243))

## Bug Fixes

* fix(cli): missing fields for Slack integrations (#170) (Salim
  Afiune)([a8ce9a9](https://github.com/lacework/go-sdk/commit/a8ce9a90f52dd81281fca78b077435229bdbafaf))

## Other Changes

* chore(api): adds alert channel \_examples/ (Salim Afiune
  Maya)([f967206](https://github.com/lacework/go-sdk/commit/f967206db3dd209f94e694b5f4db98dd8b11f113))

---

# v0.1.21

## Features

* feat(cli): Create Slack Channel Alerts 🚨 (#165) (Salim
  Afiune)([0d1f8c7](https://github.com/lacework/go-sdk/commit/0d1f8c74656c4e2043323b38cadde4e0456d6cfd))
* feat(api): add Slack Channel integrations (#164) (Salim
  Afiune)([fb81416](https://github.com/lacework/go-sdk/commit/fb81416b541882ef697d9be6dc0685772183e336))
* feat(api): new Vulnerabilities.ListEvaluations() func (#160) (Salim
  Afiune)([0060799](https://github.com/lacework/go-sdk/commit/0060799f47742091f9dc16eb987ba8cc5b5cee25))
* feat(cli): configure in non-interactive mode (#158) (Salim
  Afiune)([781f65b](https://github.com/lacework/go-sdk/commit/781f65b7f3449cbb2bb04831aa5443e7981a30e4))
* feat(cli): add --packages flag to vulnerability cmd (#149) (Salim
  Afiune)([3c34eaf](https://github.com/lacework/go-sdk/commit/3c34eaf8de21a1e5f23707034ca69c02cabf5e25))

## Other Changes

* chore(cli): remove deprecated old config loading (#159) (Salim
  Afiune)([1661939](https://github.com/lacework/go-sdk/commit/1661939b94c42c080f051039c70f4c82a56f2ad3))

---

# v0.1.20

## Features

* feat(cli): add time range flags to events list cmd (#154) (Salim
  Afiune)([e055bc0](https://github.com/lacework/go-sdk/commit/e055bc045509620239600d4f35087817ee5d7fdc))

---

# v0.1.19

## Features

* feat(cli): set User-Agent header (backend metrics) (Salim Afiune
  Maya)([bb4cfc8](https://github.com/lacework/go-sdk/commit/bb4cfc81d0176bda39bb67e4bcdb3ebb422f8110))
* feat: inject client version into User-Agent header (Salim Afiune
  Maya)([87261d2](https://github.com/lacework/go-sdk/commit/87261d2a356b3e92dc0979c6ae6070d6558d1bf4))
* feat(api): set User-Agent header (backend metrics) (Salim Afiune
  Maya)([5c5001b](https://github.com/lacework/go-sdk/commit/5c5001b340f3c8e19a9ff131dab939d36f263bdd))

---

# v0.1.18

## Features

* feat(cli): add --fixable flag to vulnerability cmd (#148) (Salim
  Afiune)([d649e2a](https://github.com/lacework/go-sdk/commit/d649e2a754be958e8504347c68ea1286dc16a58e))

---

# v0.1.17

## Other Changes

* ci: fix vuln scan cli matrix (#143) (Salim
  Afiune)([646faac](https://github.com/lacework/go-sdk/commit/646faacc762b1f361de3bc61d2e543db9b674c3c))
* ci: fix release commit message (#144) (Salim
  Afiune)([6c6f357](https://github.com/lacework/go-sdk/commit/6c6f357d4cd1e6dae08cf55e637cea4ca56aebaa))

---

# v0.1.16

## Other Changes

* ci: dogfooding lacework vulnerability scans (orb) (Salim Afiune
  Maya)([e74a188](https://github.com/lacework/go-sdk/commit/e74a18814127127395f496de908ec8bb4cb22072))
* ci: build/release docker containers automatically (Salim Afiune
  Maya)([897b05a](https://github.com/lacework/go-sdk/commit/897b05ae9dba9eb12e44d9a09bf48092f2af3764))

---

# v0.1.15

## Bug Fixes

* fix: vulnerability scans of unsupported images (Salim Afiune
  Maya)([3d33a78](https://github.com/lacework/go-sdk/commit/3d33a78baa23cd024b4e9afcd2bbaa3652274967))

## Other Changes

* chore(cli): remove deprecated --digest flag (Salim Afiune
  Maya)([aaecce1](https://github.com/lacework/go-sdk/commit/aaecce1e815ae89c761c12842bf227156432a889))
* ci: update release process to be automated v.1 (#134) (Salim
  Afiune)([374b4b0](https://github.com/lacework/go-sdk/commit/374b4b01180985fb721a250efe463eed36474286))
* ci: create release from git tag (Salim Afiune
  Maya)([ec95742](https://github.com/lacework/go-sdk/commit/ec95742ca8f0ef117d96a0c4d2d18e96fd5304c6))
* ci: upload artifacts to release (#140) (Salim
  Afiune)([7e8e03f](https://github.com/lacework/go-sdk/commit/7e8e03f5635a4ceddd45ed4caf2a133f646b4803))
* ci: add slack notifications (Salim Afiune
  Maya)([d7523b8](https://github.com/lacework/go-sdk/commit/d7523b8a593c2ca78ef46bcf84aa0c6400bc8d10))
* ci: enable integration tests in CircleCI (Salim Afiune
  Maya)([a17c238](https://github.com/lacework/go-sdk/commit/a17c238bf397b6cad2036d299971c672fd116b09))

---

# v0.1.14

## Features

* feat: understand vuln reports with 0 vulnerabilities (#124) (Salim
  Afiune)([6af13b0](https://github.com/lacework/go-sdk/commit/6af13b06ac04ff8b2efb156248a70fbb50908dde))
* feat: auto-populate account with --profile flag (#121) (Salim
  Afiune)([3539ec4](https://github.com/lacework/go-sdk/commit/3539ec409285a7d3f0335e6bfc2676f03c5fbb4c))

## Bug Fixes

* fix(spelling) Fixes event header misspelling (Scott
  Ford)([e55a6c1](https://github.com/lacework/go-sdk/commit/e55a6c16f93059d93c8ce0985a16d5bf4a7ad020))
* fix(release): update release link and version message (#117) (Salim Afiune
  Maya)([2969722](https://github.com/lacework/go-sdk/commit/2969722f94745fe348cc9c58d1c08ae22b81cf23))

## Documentation Updates

* doc: update cli documentation cli/README.md (#125) (Salim
  Afiune)([e31c4fc](https://github.com/lacework/go-sdk/commit/e31c4fc7bacaa22afa734fb35885b1eff056b98d))

## Other Changes

* chore: fix typos in AWS events (#129) (Salim
  Afiune)([46d1bb6](https://github.com/lacework/go-sdk/commit/46d1bb69203344b784976f1fb00537a65374ab69))
* chore: bump version to v0.1.14-dev (Salim Afiune
  Maya)([8e7ac41](https://github.com/lacework/go-sdk/commit/8e7ac41badd51ffc1287088ca525419d6bfb5ba2))
* ci: switch Shippable in favor of CircleCI (#120) (Salim Afiune
  Maya)([630e8bf](https://github.com/lacework/go-sdk/commit/630e8bf308d5c944ccccd8311a566d859891a927))

---

# v0.1.13

## Features

* feat(cli): avoid displaying API key secret (#115) (Salim Afiune
  Maya)([3305b09](https://github.com/lacework/go-sdk/commit/3305b095fb43a3352255e472f38ba8f19b6d7c4b))
* feat(release): add version bump after release (Salim Afiune
  Maya)([4c67b3f](https://github.com/lacework/go-sdk/commit/4c67b3fbb74fa9a05db1a712c73d1570246ffc89))

## Bug Fixes

* fix(release): purge the docker manifest to udate (Salim Afiune
  Maya)([ed58109](https://github.com/lacework/go-sdk/commit/ed58109a5ea45b7e7b7f4d9fde86f81e183f726b))

---

# v0.1.12

## Features

* feat(cli): manage compliance reports (GCP Azure AWS) (Salim Afiune
  Maya)([1d0155f](https://github.com/lacework/go-sdk/commit/1d0155f48ca4dee6a4f9381870645f3c07597dff))
* feat(api): add compliance service (Salim Afiune
  Maya)([862812c](https://github.com/lacework/go-sdk/commit/862812c4635ded3647f3e7b76e2807de06c652ba))
* feat(cli): list integrations of a specific type (Salim Afiune
  Maya)([e1d3674](https://github.com/lacework/go-sdk/commit/e1d36740f7d7fe496f7746624519c81a670d054a))

## Documentation Updates

* docs(cli): remove the need to install using sudo (Salim Afiune
  Maya)([4534c57](https://github.com/lacework/go-sdk/commit/4534c576779ca769d053c7c19e85a6029741810e))

## Other Changes

* ci: fix typo in release.sh script (Salim Afiune
  Maya)([cf6a836](https://github.com/lacework/go-sdk/commit/cf6a8369e2a6b906fb604afc6213cf7c04df8095))
* ci: add docker images to release notes (Salim Afiune
  Maya)([4f8f945](https://github.com/lacework/go-sdk/commit/4f8f945f49d2af51856617d994cd031b02ba6678))
* test(integration): add compliance tests (Salim Afiune
  Maya)([d41fb49](https://github.com/lacework/go-sdk/commit/d41fb49838a7c7990acd4b7f4fd40f0a98f2452a))

---

# v0.1.11

## Features

* feat: incident analysis, visualize event details (Salim Afiune
  Maya)([532f11d](https://github.com/lacework/go-sdk/commit/532f11d461759c9214730a1ec5b92d9ad39afbaf))

## Bug Fixes

* fix(api): use correct types on events response (Salim Afiune
  Maya)([86d8b7b](https://github.com/lacework/go-sdk/commit/86d8b7b533ef77f4b9bcf63fc839ae88be12000b))

## Other Changes

* style(cli): show help without errors (Salim Afiune
  Maya)([a72ba55](https://github.com/lacework/go-sdk/commit/a72ba55a1a35e9c0e9626d8af4c9e1ea102c6e7c))
* ci: add badge to README and encrypted keys (Salim Afiune
  Maya)([c03a416](https://github.com/lacework/go-sdk/commit/c03a41664771d6a0fcfc858223e99a347b506a20))
* test(integration): adds end-to-end tests (Salim Afiune
  Maya)([e2eb449](https://github.com/lacework/go-sdk/commit/e2eb4493bfaf73f575a3e0c1297ba4186ace34ec))
* test(integration): new framework to write CLI tests (Salim Afiune
  Maya)([402b2a2](https://github.com/lacework/go-sdk/commit/402b2a28d05a5f5bf8bfd198145d091feb2461fe))

---

# v0.1.10

## Features

* feat(cli): add aliases to integration and event cmds (Salim Afiune
  Maya)([9e8cd5c](https://github.com/lacework/go-sdk/commit/9e8cd5c4d2eb0d9cbed715a89985978e62eab9c0))
* feat(cli): preconfigure using key JSON file (WebUI) (Salim Afiune
  Maya)([80c48e7](https://github.com/lacework/go-sdk/commit/80c48e7bbaf95c888b9422249c8e09818c0a83b2))
* feat(cli): new 'integration show' cmd  (#91) (Salim Afiune
  Maya)([5bedf53](https://github.com/lacework/go-sdk/commit/5bedf5348c9fcc1748bc66534d8ac2e6475e6c64))

## Bug Fixes

* fix(docker): fix build/release of CLI containers (Salim Afiune
  Maya)([2146ecb](https://github.com/lacework/go-sdk/commit/2146ecbd6c0d4d0a9f8f608a902aeffebdce3cf9))
* fix(api): parsing event details 'cpu_percentage' (Salim Afiune
  Maya)([5f978ea](https://github.com/lacework/go-sdk/commit/5f978ead44bd6700f520ccb0742d5355464cfece))

## Other Changes

* chore: consistency with ID fields in Go structs (Salim Afiune
  Maya)([79b874e](https://github.com/lacework/go-sdk/commit/79b874ed3410b033b52a59c4fa98acb719aacfcf))

---

# v0.1.9

## Features

* feat: Add lacework-cli containers (Salim Afiune
  Maya)([73cdda0](https://github.com/lacework/go-sdk/commit/73cdda0413c56401e349162c04da261fe4e32bc7))
* feat(cli): create Azure integrations (Salim Afiune
  Maya)([29105e7](https://github.com/lacework/go-sdk/commit/29105e7fc85315b8c718906454af74245889f2a9))
* feat(cli): create GCP integrations (Salim Afiune
  Maya)([b2154a1](https://github.com/lacework/go-sdk/commit/b2154a16aa6d647514353c2a2d67c14cef9b608f))
* feat(cli): create AWS CloudTrail integrations (Salim Afiune
  Maya)([7e80795](https://github.com/lacework/go-sdk/commit/7e8079589f3f0d36c90f3e33c08ae7f168e13774))
* feat(cli): create integration sub-command (Salim Afiune
  Maya)([9842a0d](https://github.com/lacework/go-sdk/commit/9842a0db14cc059de9dd950408d2efc97de4b02a))
* feat(api): create container registry integrations (Salim Afiune
  Maya)([e33613d](https://github.com/lacework/go-sdk/commit/e33613ddcd10176464dfbcc02f09e986a5c5de01))
* feat(cli): delete external integrations (Salim Afiune
  Maya)([fe802b4](https://github.com/lacework/go-sdk/commit/fe802b45a05b70034d28bce8949362ba592aec2b))

## Refactor

* refactor(cli): new configure command using survey (Salim Afiune
  Maya)([d311ed4](https://github.com/lacework/go-sdk/commit/d311ed48ad758a48fc687db96b1ad5b2815cfeb6))

## Other Changes

* style: avoid mixing duties between api and cli (Salim Afiune
  Maya)([b245d9f](https://github.com/lacework/go-sdk/commit/b245d9f63765fdf7fb131bf933a762f9220969c8))
* style(cli): use appropriate icons per platform (Salim Afiune
  Maya)([c3e051e](https://github.com/lacework/go-sdk/commit/c3e051ed0124386796bf49d6addbad31c4d26ba4))
* chore(cli): update int create usage message (Salim Afiune
  Maya)([0959618](https://github.com/lacework/go-sdk/commit/095961838afce65a43ebf34b3405bb5b0fa09f80))
* chore(deps): remove promptui in favor of survey (Salim Afiune
  Maya)([0c663aa](https://github.com/lacework/go-sdk/commit/0c663aa23e1773aeec4162d8bf78aaadcf8f19b8))

---

# v0.1.8

## Features

* feat(cli/vul): show layer content instead of hash (Salim Afiune
  Maya)([a15e767](https://github.com/lacework/go-sdk/commit/a15e767ca4c5d05ff3fa888ce9f1333de7b545ac))
* feat(cli): add --details flag to vulnerability cmd (Salim Afiune
  Maya)([227a7b2](https://github.com/lacework/go-sdk/commit/227a7b2dbceab31c8cd982f802373a83a08c7e1b))
* feat(cli): check for available updates 👓 ✨ (Salim Afiune
  Maya)([9318952](https://github.com/lacework/go-sdk/commit/93189526b2cc518b75c0de524e748d08b27247b5))
* feat: new go library lwupdater 🆕 ⭐ (Salim Afiune
  Maya)([0f7637e](https://github.com/lacework/go-sdk/commit/0f7637e31b01ee05b0b2ced2f740c3836bcddbe2))

## Refactor

* refactor(cli): consistency between image ID & Digest (Salim Afiune
  Maya)([4f59376](https://github.com/lacework/go-sdk/commit/4f5937672218427128501bc20637485a40692f81))
* refactor(api): request and response log messages (Salim Afiune
  Maya)([e4a3b3c](https://github.com/lacework/go-sdk/commit/e4a3b3ca8e4a2db9b6ef841ff33e8feab4d9bb6e))

## Bug Fixes

* fix(cli): sort vulnerabilities by severity (Salim Afiune
  Maya)([1e0de4c](https://github.com/lacework/go-sdk/commit/1e0de4c3a6c2107d4bffd47a9f188e07c9e3ca79))

## Documentation Updates

* docs(cli): func comments and cmd style updates (Salim Afiune
  Maya)([b50f987](https://github.com/lacework/go-sdk/commit/b50f987c37946dcee4579f5c6ce67e496734c12f))
* docs(lwlogger): add basic usage example (Salim Afiune
  Maya)([c994534](https://github.com/lacework/go-sdk/commit/c99453485c61f179c3495df03291a50ecc689947))

## Other Changes

* style(cli): align vulnerability summary report (Salim Afiune
  Maya)([0b37cf6](https://github.com/lacework/go-sdk/commit/0b37cf6dd480800966a271ac86f486f1d790675a))
* style(cli): remove dup vul report summary footer (Salim Afiune
  Maya)([6e36455](https://github.com/lacework/go-sdk/commit/6e36455dfd0f7755d3727bd9ab22dd17b25c1b86))
* style: avoid mixing duties between api/ and cli/ (Salim Afiune
  Maya)([fb7b7c2](https://github.com/lacework/go-sdk/commit/fb7b7c254498bb504b16c7b42580ca0c274efeda))

---

# v0.1.7

## Bug Fixes

* fix(cli): access integration state securely (Salim Afiune
  Maya)([543562b](https://github.com/lacework/go-sdk/commit/543562b282f214ce877fabe12d1dde27017ddf30))

---

# v0.1.6

## Features

* feat(api): add EventsService to inspect events (Salim Afiune
  Maya)([533a271](https://github.com/lacework/go-sdk/commit/533a2713f5c179e50c90c63318991643f005a750))
* feat(api): add Details func to EventsService (Salim Afiune
  Maya)([56b95ca](https://github.com/lacework/go-sdk/commit/56b95ca2c02c2f8af24dd351e7fe6247b4da7eba))
* feat(cli): new event list command (Salim Afiune
  Maya)([d7c9f9e](https://github.com/lacework/go-sdk/commit/d7c9f9e2c41a1bd1411b92b6b8632aa2e32845dd))
* feat(cli): new event show command (Salim Afiune
  Maya)([8f75c78](https://github.com/lacework/go-sdk/commit/8f75c78d222f3fa00d24cf41e8d1e712f6600122))
* feat(cli): `--noninteractive` mode flag (Salim Afiune
  Maya)([10536af](https://github.com/lacework/go-sdk/commit/10536afe1d6ce76ba3391c145f3533b4d6725484))

## Bug Fixes

* fix(api): omitempty integration responses fields (Salim Afiune
  Maya)([44e2314](https://github.com/lacework/go-sdk/commit/44e2314f4ca02f1c0e6a134bedbc82161a81473c))

## Documentation Updates

* docs(cli): document environment variables (Salim Afiune
  Maya)([0012ec1](https://github.com/lacework/go-sdk/commit/0012ec14f574f6e4c1dc2a5d774e17ef038f1308))

## Other Changes

* chore(cli): update usage of commands (Salim Afiune
  Maya)([5dd3057](https://github.com/lacework/go-sdk/commit/5dd3057371fe87434f2da54f68cdcc3dc5fd754a))
* chore(cli): style updates to release scripts (Salim Afiune
  Maya)([f4355bf](https://github.com/lacework/go-sdk/commit/f4355bf481a3349fe9bac3700a1eaa8e80227238))

---

# v0.1.5

## Features

* feat(cli): implement JSON format for all commands (Salim Afiune
  Maya)([c7d4fee](https://github.com/lacework/go-sdk/commit/c7d4fee38f7b02a678e2fd30ecdb8f5bff82c11b))
* feat(cli): `vul scan run` command can poll for status (Salim Afiune
  Maya)([e2c8c8d](https://github.com/lacework/go-sdk/commit/e2c8c8dd1ee98d800902e6b671477e955e30a7fe))

## Bug Fixes

* fix(install.ps1) copy lacework.exe to ProgramData (Salim Afiune
  Maya)([53e685f](https://github.com/lacework/go-sdk/commit/53e685f583b70ff56b40ffdd8b2500d258891580))

---

# v0.1.4

## Features

* feat(install.ps1): support to install cli on windows (Salim Afiune
  Maya)([ae53d6f](https://github.com/lacework/go-sdk/commit/ae53d6ffe0f79b7d1f0cdbbd4feab560f727f29c))

## Bug Fixes

* fix(cli): use correct variable inside install.sh (Salim Afiune
  Maya)([eda17d5](https://github.com/lacework/go-sdk/commit/eda17d5ef4b2986bcbe99e4db0213ca7ebd18876))
* fix(cli): support colors for windows (Salim Afiune
  Maya)([ea48379](https://github.com/lacework/go-sdk/commit/ea4837955064c32575311722b9b137fd08161417))

## Documentation Updates

* doc(api): adds `_examples` or token-generation (Salim Afiune
  Maya)([c9dbc02](https://github.com/lacework/go-sdk/commit/c9dbc022831ed6b79d5216138c439662c7ad39e4))
* docs(cli): update install.sh to use CLI (Salim Afiune
  Maya)([d0fda04](https://github.com/lacework/go-sdk/commit/d0fda0460f9e247b53f92f4f72e32fa248de9b8c))
* docs: added lacework-cli profiles (Salim Afiune
  Maya)([d2292b0](https://github.com/lacework/go-sdk/commit/d2292b0425cde9334afc7b92e9fdc76ec63649c4))

## Other Changes

* chore(typo): misspelled word inside install.ps1 (Salim Afiune
  Maya)([ec93a96](https://github.com/lacework/go-sdk/commit/ec93a964407a84523708abaaee8569390216ec21))
* chore(timeout): increase api timeout to 60s (Salim Afiune
  Maya)([1f17bad](https://github.com/lacework/go-sdk/commit/1f17badb40690c04389d87453f556686fac17fb0))

---

# v0.1.3

## Features

* feat(vul): show number of fixable vulnerabilities (Salim Afiune
  Maya)([6403029](https://github.com/lacework/go-sdk/commit/6403029284085579a10cd4f74b41c3fd53ca765f))
* feat(cli): new vulnerability command (Salim Afiune
  Maya)([494d8d8](https://github.com/lacework/go-sdk/commit/494d8d8e317cad38bb6012567c4ff8ba9a2d3aa4))
* feat(api): add vulnerabilities service (Salim Afiune
  Maya)([d0b2c3b](https://github.com/lacework/go-sdk/commit/d0b2c3b5f2ad197a5ee9ba0d276c2d16d618891d))
* feat: introducing named profiles (Salim Afiune
  Maya)([6fb64fd](https://github.com/lacework/go-sdk/commit/6fb64fd2e0f953bf10103a44a821f840292991c4))
* feat: disallow extra arguments on sub-commands (#48) (Salim Afiune
  Maya)([f67ca9a](https://github.com/lacework/go-sdk/commit/f67ca9af38863e3d4a44a16d1130f919f7c2592e))
* feat: add configure command (#47) (Salim Afiune
  Maya)([f334fda](https://github.com/lacework/go-sdk/commit/f334fda595e000cc6ab830314c88f950eabc6761))

## Other Changes

* chore: adds a couple new go package dependencies (Salim Afiune
  Maya)([1842700](https://github.com/lacework/go-sdk/commit/1842700a8e7dfc91e91b11d5b21a62247a755b34))

---

# v0.1.2

## Features

* feat(lwloggder): go package for logging messages (Salim Afiune
  Maya)([cb5feee](https://github.com/lacework/go-sdk/commit/cb5feeeb6c1ddad54c6163c1e2b2c4dfdb6381fa))

## Refactor

* refactor(cli): rename cli binary to lacework (Salim Afiune
  Maya)([51ce22f](https://github.com/lacework/go-sdk/commit/51ce22f579f2984dfcfd54faaee844b621d7a617))

---

# v0.1.1

## Features

* feat(api): debug logs for all requests & responses (Salim Afiune
  Maya)([209f7ee](https://github.com/lacework/go-sdk/commit/209f7ee6240ccacabaac136b258409a48f673e8f))
* feat(api): add api client IDs for multi-client req (Salim Afiune
  Maya)([82c209f](https://github.com/lacework/go-sdk/commit/82c209fd04690a98109cc4a197431de29490c419))
* feat(api): implement a logging mechanism using zap (Salim Afiune
  Maya)([c078a70](https://github.com/lacework/go-sdk/commit/c078a70c002cb1646f9ba21e9d6d0b18aee85fdc))

## Bug Fixes

* fix(cli): error when account is empty (Salim Afiune
  Maya)([7dc59aa](https://github.com/lacework/go-sdk/commit/7dc59aa5958016dc2f3524495296cd27366fe660))
* fix(cli): load debug state correctly (Salim Afiune
  Maya)([8f7343c](https://github.com/lacework/go-sdk/commit/8f7343c93c11b757b64e2fed163684634a95ff36))
* fix(cli) Update environment variable prefix (Scott
  Ford)([484ca39](https://github.com/lacework/go-sdk/commit/484ca3951ecc934fd12ca88691b20a9781d245c7))
* docs(README) Update cli README to add documentation for ENV VARS (Scott
  Ford)([484ca39](https://github.com/lacework/go-sdk/commit/484ca3951ecc934fd12ca88691b20a9781d245c7))

## Documentation Updates

* doc(cli): fix single quote typo (Salim Afiune
  Maya)([3770b89](https://github.com/lacework/go-sdk/commit/3770b89d3d90aed18822fd2fa521cef3f23351bd))
* doc(logo): add logo to main README (Salim Afiune
  Maya)([620b992](https://github.com/lacework/go-sdk/commit/620b9921fa9d04147f4dfce3f4a24dfdec9a4238))

## Other Changes

* chore(cli): hide integration sub-commands (Salim Afiune
  Maya)([791ef7d](https://github.com/lacework/go-sdk/commit/791ef7d4df542c382607a12de08b262b95059311))
* chore(typo): fix RestfulAPI typo (Salim Afiune
  Maya)([39a7298](https://github.com/lacework/go-sdk/commit/39a72989366c85288b91d9ad033ed431a7796395))
* build: fix release checks (Salim Afiune
  Maya)([08bdb7d](https://github.com/lacework/go-sdk/commit/08bdb7d67fa22cde8b16a8ab392646b3288e5aa9))
* build(release): generate changelog and release notes (Salim Afiune
  Maya)([3aa0a91](https://github.com/lacework/go-sdk/commit/3aa0a91f31984fb0935b079c573b2f33ab5e7831))

---

# v0.1.0

## Features

* feat(cli): Installation scripts and documentation 🎉 (Salim Afiune
  Maya)([bb96b3b](https://github.com/lacework/go-sdk/commit/bb96b3bf26f105137afc50011b2c88c67e4ed0c7))
* feat(cli): the new lacework-cli MVP 🔥🔥 (Salim Afiune
  Maya)([34a73b6](https://github.com/lacework/go-sdk/commit/34a73b6d8df6e58225186831ae62a86a1724d747))
* feat(integrations): add AZURE_CFG and polish the rest (Salim Afiune
  Maya)([abd5bee](https://github.com/lacework/go-sdk/commit/abd5bee7d21141116f630faef44bdff707385e1d))
* feat(api): List integrations by type (Salim Afiune
  Maya)([f96a15b](https://github.com/lacework/go-sdk/commit/f96a15bc492d1acfc427b1851e3dfc12dc83a48b))
* feat: implement service model (Salim Afiune
  Maya)([d0cbf9f](https://github.com/lacework/go-sdk/commit/d0cbf9f20d7605a0f8a782909c09e78b3f1a6d8e))
* feat(api): new GetIntegrationSchema() (Salim Afiune
  Maya)([1aaec6c](https://github.com/lacework/go-sdk/commit/1aaec6cc8f3784d93f01d9819547dcab8cefddd4))
* feat(integrations): CRUD azure config integrations (Salim Afiune
  Maya)([0f83504](https://github.com/lacework/go-sdk/commit/0f83504636753c5e3d8cede843c1e18c73826c84))
* feat(integrations): CRUD aws config integrations (Salim Afiune
  Maya)([93475b0](https://github.com/lacework/go-sdk/commit/93475b0b229100dbb44ca4f779eb19120e93d298))
* feat(request): trigger token generation if missing (Salim Afiune
  Maya)([8cd82d6](https://github.com/lacework/go-sdk/commit/8cd82d677f65b7606111f91685e22e0618e6404d))
* feat(fakeAPI): New LaceworkServer to mock API req (Salim Afiune
  Maya)([c8211c1](https://github.com/lacework/go-sdk/commit/c8211c1e91a2ac116b65067ba890ab5dcae38e98))
* feat(client): Option to trigger a new token gen (Salim Afiune
  Maya)([96c8c6b](https://github.com/lacework/go-sdk/commit/96c8c6b78c8e73a17b8a71f3d2f25a78ce598b57))

## Refactor

* refactor: leverage integration structs for all gcp (Salim Afiune
  Maya)([922d117](https://github.com/lacework/go-sdk/commit/922d11755d4200cfbecdf56cb7e5f96775b9f136))
* refactor: leverage integration structs for all azure (Salim Afiune
  Maya)([1037d1b](https://github.com/lacework/go-sdk/commit/1037d1b8ee0028e4c688dbd9468499e7865eea6f))
* refactor: leverage integration structs for all aws (Salim Afiune
  Maya)([1146348](https://github.com/lacework/go-sdk/commit/11463481c419b1b9bb691ef1c15a43550b268095))
* refactor(integration): make space for New() funcs (Salim Afiune
  Maya)([1da9746](https://github.com/lacework/go-sdk/commit/1da9746332bb60dbb06c8345164b362fb0d06ce3))
* refactor(integration): move CRUD gcp config code (Salim Afiune
  Maya)([962191b](https://github.com/lacework/go-sdk/commit/962191b6c9383d6b52c9e6d78bde2507872541ba))

## Bug Fixes

* fix(install): configurable installation_dir (Salim Afiune
  Maya)([9d17b1f](https://github.com/lacework/go-sdk/commit/9d17b1f2d5e424acbee873d5b9934ff85e4f67f0))
* fix(release): tar linux binaries (Salim Afiune
  Maya)([9311b8f](https://github.com/lacework/go-sdk/commit/9311b8f2f5bd6e65381f7d0e8900da87925f8f78))
* fix(update): GCP CFG api path (
  mjunglw)([3508e78](https://github.com/lacework/go-sdk/commit/3508e78d483dd53b46e5e8d9c2b26ded257b1ef4))
* fix(gcp): update missing fields in structs (#14) (
  lwmobeent)([ce9745f](https://github.com/lacework/go-sdk/commit/ce9745f6202dac7e1c27b3f0eef582df29bc21a6))
* fix(enums): integrationType and gcpResourceLevel from array to map (
  mjunglw)([37c4d77](https://github.com/lacework/go-sdk/commit/37c4d771b27475eec0a05a55d00375e41d8ba53b))
* fix(client): expose Client struct for provider to use (
  mjunglw)([5e97951](https://github.com/lacework/go-sdk/commit/5e9795177ce2e76e029f66a60bbfb1497514ba08))
* fix(lint): various lint fixes (Salim Afiune
  Maya)([92efbab](https://github.com/lacework/go-sdk/commit/92efbab9b952e240045687717388f6badff072de))

## Documentation Updates

* docs(api): update README's and code comments (Salim Afiune
  Maya)([9d3e739](https://github.com/lacework/go-sdk/commit/9d3e7399275becf01ac23aaf702fd288c4518c2e))
* docs(README): Add usage and descriptions (Salim Afiune
  Maya)([9ed08dc](https://github.com/lacework/go-sdk/commit/9ed08dc604ef2482c8566f2d431a476920376239))

## Other Changes

* chore(deps): Add a few Go dependencies 🙌 (Salim Afiune
  Maya)([4ae8b8e](https://github.com/lacework/go-sdk/commit/4ae8b8e7d2671b40a98ae4b863192e6fe89f7e5b))
* chore(api): make response structs public (Salim Afiune
  Maya)([6b84e2c](https://github.com/lacework/go-sdk/commit/6b84e2c5786012d0d672df9e46c1e1166a38d086))
* ci(prepare): prepare the pipeline (Salim Afiune
  Maya)([94a8b0f](https://github.com/lacework/go-sdk/commit/94a8b0fddd0951734bd39137a77acb79cd3e9224))
* ci(tests): run tests in shippable ci (Salim Afiune
  Maya)([cb73c4b](https://github.com/lacework/go-sdk/commit/cb73c4bca8d9689ee505ccbe4848b9e8cfe72fa9))
* test: fix integration delete requests (Salim Afiune
  Maya)([4181580](https://github.com/lacework/go-sdk/commit/41815805e7039b7bb066432ced636c6d59e1d81b))
* test(integrations): generic Delete() func (Salim Afiune
  Maya)([c491f1a](https://github.com/lacework/go-sdk/commit/c491f1af2574330b7eda8c7ed285f05abb17a099))
* test(integrations): generic Get() func (Salim Afiune
  Maya)([0d0c8b0](https://github.com/lacework/go-sdk/commit/0d0c8b0e9d014b868a4e59004fd5407d63090ab0))
* test(unit): verify integrationType is well map (Salim Afiune
  Maya)([320640f](https://github.com/lacework/go-sdk/commit/320640f35068046ebbefccd9bfef95234e3dad4c))
