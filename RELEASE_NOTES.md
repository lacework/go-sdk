# Release Notes
Another day, another release. These are the release notes for the version `v0.2.11`.

## Features
* feat(cli): daily version check (#269) (Salim Afiune)([5c15eef](https://github.com/lacework/go-sdk/commit/5c15eef84f428ec0534954babb28a3db92d5a7c5))
* feat(api): add Webhook integration (#267) (Darren)([f32572e](https://github.com/lacework/go-sdk/commit/f32572ecdadd5c179227cd228bf1fdd7cf618763))
## Refactor
* refactor(cli): abstract rendering tables (human-readable) (#263) (Salim Afiune)([8a10b4c](https://github.com/lacework/go-sdk/commit/8a10b4cf10de03d9b4c0409e495fdd7118974b92))
## Bug Fixes
* fix(cli): render account mapping file correctly (#266) (Salim Afiune)([4c327d7](https://github.com/lacework/go-sdk/commit/4c327d7e6081d0f7726a1bc007b1b736a106933f))
* fix(api): new request body for lql service (#260) (Salim Afiune)([4e2b439](https://github.com/lacework/go-sdk/commit/4e2b439ff394d632cd6ebf214da376050da46812))
* fix(api): avoid updating AgentTokenResponse.Props (#259) (Salim Afiune)([c3fe8bc](https://github.com/lacework/go-sdk/commit/c3fe8bcc41efd995f756f97a5ffca8bb961e89e4))
## Documentation Updates
* docs: update READMEs and _examples/ (#268) (Salim Afiune)([3791da0](https://github.com/lacework/go-sdk/commit/3791da01005335c34852446c57eb99e51a6d3ce1))
## Other Changes
* build: upgrade Go version to 1.15 (#265) (Salim Afiune)([06d41f5](https://github.com/lacework/go-sdk/commit/06d41f56add71f8369ffae68ea7ba5d738eb4d5b))
* ci: update hostname from our test machine (#262) (Salim Afiune)([beb289e](https://github.com/lacework/go-sdk/commit/beb289e732c177e2f3d062d61dee5dd9f1593ce9))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
