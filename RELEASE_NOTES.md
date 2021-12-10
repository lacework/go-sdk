# Release Notes
Another day, another release. These are the release notes for the version `v0.23.0`.

## Features
* feat(api): CRUD operations for org team member (#628) (vatasha)([80015b9](https://github.com/lacework/go-sdk/commit/80015b9807e3dc541a6e3e3c1b8c726a84995209))
* feat(cli): generate & execute aws terraform code (#626) (Matt Cadorette)([992bfa6](https://github.com/lacework/go-sdk/commit/992bfa6c3fbd08c302d53bb02a27831278f0a493))
* feat(cli): Manage Report Rules (#625) (Darren)([890350b](https://github.com/lacework/go-sdk/commit/890350b2b34af6380de15df8d665828bb313776a))
* feat(cli): add breadcrums to compliance commands (Salim Afiune Maya)([f133971](https://github.com/lacework/go-sdk/commit/f1339716c91f7cdb7d2a4fd64bac9f457cedb6ab))
* feat(cli): cache compliance reports for 30 minutes (Salim Afiune Maya)([77f82aa](https://github.com/lacework/go-sdk/commit/77f82aad8e39c208afbef6a35ad62e4206fe56ae))
* feat(api): implement crud operations for team members (#614) (vatasha)([4aa40a2](https://github.com/lacework/go-sdk/commit/4aa40a23ed2c53beeea0bc60a851bf47488b4d52))
## Bug Fixes
* fix(cli): make sure progress bar starts with space (#635) (Salim Afiune)([4153f64](https://github.com/lacework/go-sdk/commit/4153f64f071bb85c0718967be4f71a9bd7071e65))
* fix(cli): fix cache directory for azure comp reports (Salim Afiune Maya)([19be146](https://github.com/lacework/go-sdk/commit/19be1465a3f73f6359c32454863312a407318e6b))
* fix: report rule notifications (#621) (Darren)([f37b913](https://github.com/lacework/go-sdk/commit/f37b91349c031cf0d92c79a7f4dfc181446e69e3))
## Other Changes
* ci: fix test resource Terraform pipeline (#634) (Salim Afiune)([c14756e](https://github.com/lacework/go-sdk/commit/c14756e12c0f051e871e0f7a7aeccd1089cdd875))
* ci: version bump to v0.22.1-dev (Lacework)([e12f022](https://github.com/lacework/go-sdk/commit/e12f022356992ae477dbbf211f9439dd4b080dce))
* test(policy): update TestPolicyCreateURL (hazedav)([a89177c](https://github.com/lacework/go-sdk/commit/a89177ce3d82660c8700b555c0596f7a98e53d47))
* test(policy): adds resource for dynamic accnt (hazedav)([b28d290](https://github.com/lacework/go-sdk/commit/b28d2902aa6d1a557af913c60a0be7a68c3b9895))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
