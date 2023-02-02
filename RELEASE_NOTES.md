# Release Notes
Another day, another release. These are the release notes for the version `v1.9.0`.

## Features
* feat(api): v2PageMetadata to log page information (Salim Afiune Maya)([fd5338b6](https://github.com/lacework/go-sdk/commit/fd5338b69ac3ea22a611cc57799e4a139c545e8a))
* feat(cli): improve caching using hash of filters (Salim Afiune Maya)([b8b5e9ea](https://github.com/lacework/go-sdk/commit/b8b5e9eaa93514a2bf7defa08159eb432b947bb4))
## Refactor
* refactor(cli): deprecate --image_id flag (Salim Afiune Maya)([82db2c29](https://github.com/lacework/go-sdk/commit/82db2c2990d4503616940831fe6007434aef14c9))
* refactor(cli): use evaluation GUID to retrieve scan results (Salim Afiune Maya)([7b0fb126](https://github.com/lacework/go-sdk/commit/7b0fb126dce5886f0bced4ba9ce475242ab1b1cf))
* refactor(cli): avoid unnecessary memory consumption (Salim Afiune Maya)([a0c0b543](https://github.com/lacework/go-sdk/commit/a0c0b54301e0159c7000659c287e54424be37311))
## Bug Fixes
* fix(cli): --fixable and other filters (Salim Afiune Maya)([d3cb69a4](https://github.com/lacework/go-sdk/commit/d3cb69a4e8ccdb6641977d017e47a1e67bb8af47))
* fix(cli): --registry flag fails faster for list-assessments command (Salim Afiune Maya)([cdd71461](https://github.com/lacework/go-sdk/commit/cdd7146133186711bd29357b8dc12b6687cf2931))
## Other Changes
* style: use cobra deprecation feature for --image_id (Salim Afiune Maya)([53f3b114](https://github.com/lacework/go-sdk/commit/53f3b11459a13a84f73139dbc6c411a4029ea31e))
* ci: version bump to v1.8.1-dev (Lacework)([9adc8a12](https://github.com/lacework/go-sdk/commit/9adc8a12498ca2976d4707ecdf759090070f47a0))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
