# Release Notes
Another day, another release. These are the release notes for the version `v1.50.0`.

## Features
* feat: expose raw machine tags for external automation (#1618) (Matt Cadorette)([5e2e2c3c](https://github.com/lacework/go-sdk/commit/5e2e2c3c6247186e46750720f2cf91521d494e29))
* feat(GROW-2876): support terraform output blocks in lwgenerate (#1609) (Matt Cadorette)([5440c57d](https://github.com/lacework/go-sdk/commit/5440c57d49156c11cbf1b87f9474d3e5b7134f50))
## Bug Fixes
* fix: use buffered channel for downloadComplete to prevent blocking (#1617) (Matt Cadorette)([7ddbd389](https://github.com/lacework/go-sdk/commit/7ddbd389bc8134a50df4c23c75727aea57232db4))
* fix: policy severity test correct test criteria (#1614) (Matt Cadorette)([ef808de2](https://github.com/lacework/go-sdk/commit/ef808de2738e75f7d4b4bc9876f69922b67e3809))
* fix: correct integration tests for gcp (#1612) (Matt Cadorette)([8271a9f4](https://github.com/lacework/go-sdk/commit/8271a9f4ddc56eb623777cb77d93dac289bcbfc8))
* fix: failing tests due to transient data states (#1608) (Matt Cadorette)([259b911b](https://github.com/lacework/go-sdk/commit/259b911b4d3a94cf73d6240293729d9f63249df7))
* fix(cli): remove option to create AzureAD intg from CLI (manual) (#1604) (Manan Bhatia)([79711080](https://github.com/lacework/go-sdk/commit/79711080fcd947668c61846cc1ef3f7a1b330ad3))
* fix(api): correct handling for resourceTags in policy exceptions (#1605) (Matt Cadorette)([fbe45aff](https://github.com/lacework/go-sdk/commit/fbe45aff6f3b84b5e4e738bdac432bdbf726c3af))
## Other Changes
* ci: version bump to v1.49.1-dev (Lacework)([db33cada](https://github.com/lacework/go-sdk/commit/db33cada5778db7c70eea9f74865ff1b6d70e49c))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
