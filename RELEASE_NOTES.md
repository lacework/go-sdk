# Release Notes
Another day, another release. These are the release notes for the version `v1.45.0`.

## Features
* feat: add cdk component client (#1538) (Matt Cadorette)([4fca2248](https://github.com/lacework/go-sdk/commit/4fca2248f2dfff90237f167f320e40c5682a42d9))
## Bug Fixes
* fix: component install windows (#1543) (jonathan stewart)([6c906b5b](https://github.com/lacework/go-sdk/commit/6c906b5ba5019ce9f81c456260d72796294e8b2f))
* fix: reduce component API usage (#1542) (jonathan stewart)([5455e173](https://github.com/lacework/go-sdk/commit/5455e173edac65652c57a56df307c9d82c1b8a5d))
## Other Changes
* chore: limit the concurrency of test-build workflow runs for a branch (#1540) (Pengyuan Zhao)([78fe2e17](https://github.com/lacework/go-sdk/commit/78fe2e1760ba7e0a5791940175995590b5519f7c))
* chore(cli): print new component version available message (#1530) (Pengyuan Zhao)([d12d9ccd](https://github.com/lacework/go-sdk/commit/d12d9ccdd515592d8d2debec7bfef49bf0273326))
* chore(CLI): cache CDK API data in disk (#1529) (Pengyuan Zhao)([e9693715](https://github.com/lacework/go-sdk/commit/e969371562aa830945e7b4ace6a78423b98d57d5))
* chore(deps): bump golang.org/x/crypto from 0.14.0 to 0.17.0 (#1492) (dependabot[bot])([ff59352b](https://github.com/lacework/go-sdk/commit/ff59352b4bebf910d5d5014d679b4ff899efea09))
* chore(deps): bump github.com/BurntSushi/toml from 1.2.0 to 1.3.2 (#1490) (dependabot[bot])([4a5f655b](https://github.com/lacework/go-sdk/commit/4a5f655b2af644694b742467ad29d51171839f7c))
* chore: use CDK v1 API for lacework component dev (#1525) (Pengyuan Zhao)([14058d5c](https://github.com/lacework/go-sdk/commit/14058d5ca73834941e4d3452c7f3c0c2dc94043e))
* chore(cli): only fetch component versions when needed (#1526) (Pengyuan Zhao)([183e16ed](https://github.com/lacework/go-sdk/commit/183e16ed37544a3ab731f3abfe921ae13157e343))
* ci: Create workflow release.yml (#1539) (Pengyuan Zhao)([78307db3](https://github.com/lacework/go-sdk/commit/78307db3fb6144c8da7035ea3cdb78b7a48066d5))
* ci: Create workflow prepare-release.yml (#1536) (Pengyuan Zhao)([d662565f](https://github.com/lacework/go-sdk/commit/d662565fb398ced41c75f2196d794f8133882ebe))
* ci: Create workflow verify-release.yml (#1537) (Pengyuan Zhao)([b7a2af65](https://github.com/lacework/go-sdk/commit/b7a2af650d0b0da19bf8bc235cf996d2e09ebea3))
* ci: Create workflow update-homebrew-formula.yml (#1532) (Pengyuan Zhao)([ec268e36](https://github.com/lacework/go-sdk/commit/ec268e360acedb4ae75234850168e50e424bc88a))
* ci: Create workflow update-cli-docs.yml (#1531) (Pengyuan Zhao)([b5842a4f](https://github.com/lacework/go-sdk/commit/b5842a4fcfb1b67c21dc7d32fe951d49d1333257))
* ci: Create workflow test-build.yml (#1534) (Pengyuan Zhao)([8e089d78](https://github.com/lacework/go-sdk/commit/8e089d78cf36e99cdeef609e7490c2848b949604))
* ci: Create workflow nightly-build.yml (#1535) (Pengyuan Zhao)([6fa444b3](https://github.com/lacework/go-sdk/commit/6fa444b3ddaa6e3eba9443a7ec4558cb408a8e95))
* ci: version bump to v1.44.5-dev (Lacework)([05aa1e9b](https://github.com/lacework/go-sdk/commit/05aa1e9be199783dce0362bfee26b8cc50b3c9f3))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
