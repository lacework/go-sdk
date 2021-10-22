# Release Notes
Another day, another release. These are the release notes for the version `v0.18.0`.

## Features
* feat(cli): unhiding query "sources" cmd (hazedav)([4e2513f](https://github.com/lacework/go-sdk/commit/4e2513fc3ccd4f4b9818378c98f771df9f3a4546))
* feat(api): Alert Rules v2 service (#583) (Darren)([a549a8c](https://github.com/lacework/go-sdk/commit/a549a8c501e00a42be879d38197fc13a4984c914))
* feat(cli): Migrate New Relic alert channel to API v2 (#586) (Darren)([6f32767](https://github.com/lacework/go-sdk/commit/6f327671b19e221859bff761e495ba696279bb61))
* feat(cli): Migrate Pager Duty alert channel to API v2 (#585) (Darren)([68be1ec](https://github.com/lacework/go-sdk/commit/68be1eccbce846810ec35ece9819f99b1c73d5e9))
* feat(api): Migrate GCP Pub Sub alert channel to API v2 (#581) (vatasha)([7847d3b](https://github.com/lacework/go-sdk/commit/7847d3bddcaf7e2f8174684f25588559db4cdb37))
* feat(cli): Migrate Ibm QRadar alert channel to API v2 (#584) (Darren)([b2a0cde](https://github.com/lacework/go-sdk/commit/b2a0cde29e8f6400934d9ac4dc40fab4fe1ef513))
* feat(api): Migrate Service Now alert channel to API v2 (#582) (Darren)([692c8e5](https://github.com/lacework/go-sdk/commit/692c8e58650ed98891cd9202ea273400f51744b3))
* feat: Migrate Splunk alert channel to V2 (#580) (Darren)([c0c422b](https://github.com/lacework/go-sdk/commit/c0c422b71b7893958c0641c7d75b2401546bfee5))
* feat(cli): Add support for new GCP and AWS compliance reports (#579) (iancrichardson)([1de409b](https://github.com/lacework/go-sdk/commit/1de409b0aaaa9d20b8b882f04db5dae679157f07))
* feat(cli): adding lql query data sources (hazedav)([8828823](https://github.com/lacework/go-sdk/commit/882882383471db9b32076e047d9d00ab3d0a37c7))
* feat(api): Microsoft Teams alert channel APIv2 (Salim Afiune Maya)([57b318e](https://github.com/lacework/go-sdk/commit/57b318e102749c5bd34901f595424671b61c0b1e))
* feat(api): Migrate Cisco Spark Webhook alert channel to API v2 (#577) (vatasha)([9481999](https://github.com/lacework/go-sdk/commit/9481999f2faeabde3a9085601f5fbfdfee17531c))
## Refactor
* refactor(cli): create Microsoft Teams alert channels via APIv2 (Salim Afiune Maya)([74372b8](https://github.com/lacework/go-sdk/commit/74372b8ba1303eeb259799b3d080c0329fb34253))
## Bug Fixes
* fix: Splunk ssl (#587) (Darren)([eab94cf](https://github.com/lacework/go-sdk/commit/eab94cf33b16b48cb7fb3146c2aa903a7ae37d7e))
## Other Changes
* style: rename TeamsUrl in favor of TeamsURL (Salim Afiune Maya)([bd46e0d](https://github.com/lacework/go-sdk/commit/bd46e0d46c37017834fcd6521ade92d28cb0f9f9))
* ci: version bump to v0.17.1-dev (Lacework)([9df310c](https://github.com/lacework/go-sdk/commit/9df310c0a20d97b79ecfdf4c6ffe520b87456d93))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
