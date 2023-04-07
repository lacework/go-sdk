# Release Notes
Another day, another release. These are the release notes for the version `v1.16.0`.

## Features
* feat(cli): enable default query for 'lql run' cmd (#1218) (Salim Afiune)([f91c5463](https://github.com/lacework/go-sdk/commit/f91c5463eab777624e7cf9ab18c50fb006f42a54))
* feat(cli): add `server_url` to `agent install` cmd (#1224) (Nick Schmeller)([475af398](https://github.com/lacework/go-sdk/commit/475af398f180f62736089f26919496c8d49527c3))
* feat(api): Implement feature flagging endpoint (#1212) (Linpeng Zhang)([8a3e4bd6](https://github.com/lacework/go-sdk/commit/8a3e4bd6c37c198b45a90384e8fdd851e3b58375))
* feat(cli): add `server_url` flag for CTF commands (#1219) (Nick Schmeller)([070ce61f](https://github.com/lacework/go-sdk/commit/070ce61fb929ef406b551905ef19f97a7c731c38))
* feat(cli): new report-definitions revert command (#1209) (Darren)([33961d48](https://github.com/lacework/go-sdk/commit/33961d484a66605c6e9fcb74027b4561a0b37624))
## Refactor
* refactor: show valid versions when invalid report definition version is given (#1216) (Darren)([52a0ae4c](https://github.com/lacework/go-sdk/commit/52a0ae4c207e4b0e880cd02b53ed5217ad61c523))
* refactor: move version check output to stderr (#1207) (Darren)([1793885c](https://github.com/lacework/go-sdk/commit/1793885c4af1ee09388d5432d32683bba63813cb))
## Bug Fixes
* fix(cli): custom clock offset, defaults now to -2s (#1222) (Darren)([3b63c959](https://github.com/lacework/go-sdk/commit/3b63c959629b9b0b1a2a28496940a0c9ff15fcc0))
## Other Changes
* chore(cli): run `make go-vendor` to fix release (#1221) (Nick Schmeller)([7553e4bd](https://github.com/lacework/go-sdk/commit/7553e4bde715528a86032206bc5c239570d0f6ab))
* chore(lwgenerate): bump azure modules to 2.0 (#1210) (Salim Afiune)([c4857515](https://github.com/lacework/go-sdk/commit/c485751561196c8c187256d695a8d55d5cbc1811))
* ci: version bump to v1.15.2-dev (Lacework)([1a969222](https://github.com/lacework/go-sdk/commit/1a9692220c8b61074d295340a5045c570aa71149))
* test(cli): fix daily version integration test (#1225) (Salim Afiune)([694a15b3](https://github.com/lacework/go-sdk/commit/694a15b3c8495bda6f3e364d46926acc397da9a2))
* test: update version test to check stderr for version update notifica… (#1211) (Darren)([5f82668b](https://github.com/lacework/go-sdk/commit/5f82668b3950e32374a95b8702dac69f00185f24))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
