# Release Notes
Another day, another release. These are the release notes for the version `v0.12.0`.

## Features
* feat(cli): Output v2 integration state details (#505) (Darren)([e2cc6f1](https://github.com/lacework/go-sdk/commit/e2cc6f1cbd3da43df6e014ccd9075ec59af6ae06))
* feat: Add apiV2AlertChannelTest (#502) (Darren)([f968b89](https://github.com/lacework/go-sdk/commit/f968b89bb0fa55d8eb84cf439626e2610507b883))
* feat(cli): create Github Registry integrations (Salim Afiune Maya)([d3756e4](https://github.com/lacework/go-sdk/commit/d3756e4b666de9039e7ccc01b706f5ed4e750f1e))
* feat(cli): create Google Artifact Registry integrations (Salim Afiune Maya)([0f1bb68](https://github.com/lacework/go-sdk/commit/0f1bb6866bc27e5e9e39fa79cdda350513c3f789))
* feat(api): add Github Registry support (Salim Afiune Maya)([697ca98](https://github.com/lacework/go-sdk/commit/697ca98ad540c6a0e79eeda26d0e2b0cb7f9271a))
* feat(api): add Google Artifact Registry support (Salim Afiune Maya)([503ae99](https://github.com/lacework/go-sdk/commit/503ae9953f510688974dafd35287d3d3fe4e12b0))
* feat(api): new v2 Container Registries service (Salim Afiune Maya)([46212c6](https://github.com/lacework/go-sdk/commit/46212c6b72fe964eb56f3618afec67dc2a641673))
## Refactor
* refactor(cli): policy update api args (hazedav)([a9add4b](https://github.com/lacework/go-sdk/commit/a9add4b7144b9461c4a398801e2ae6d10f6936cf))
* refactor(cli): split survey limits into diff funcs (Salim Afiune Maya)([1e1dca4](https://github.com/lacework/go-sdk/commit/1e1dca49a7a485ba81ce21a7dfa90c3ea7cb6dd2))
## Bug Fixes
* fix: comp [aws | gcp | az] list json output (#497) (Darren)([c3fb601](https://github.com/lacework/go-sdk/commit/c3fb60168e186611868943e3eaea3d92973c7aec))
## Other Changes
* chore(deps): bump github.com/honeycombio/libhoney-go (#494) (dependabot[bot])([9c1985c](https://github.com/lacework/go-sdk/commit/9c1985c0dbd3341ad2c47b797ebbce0a1941b0c1))
* chore: fixed minor typos (Salim Afiune Maya)([a13e58b](https://github.com/lacework/go-sdk/commit/a13e58bda7b1294936538481609e3421fdc4c0df))
* chore(deps): add github.com/fatih/structs package (Salim Afiune Maya)([89d6e0d](https://github.com/lacework/go-sdk/commit/89d6e0d30cb5405540f9493874eebd7652f5f5ef))
* ci: sign lacework-releng commits (#500) (Salim Afiune)([523097c](https://github.com/lacework/go-sdk/commit/523097ce8b4c6ede3317d1143446cf41df85a27f))
* test(cli): updating policy integration resource (hazedav)([acceada](https://github.com/lacework/go-sdk/commit/acceadae2c7bca2364ebe5ddda51e39b6158ced4))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
