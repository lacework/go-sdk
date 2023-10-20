# Release Notes
Another day, another release. These are the release notes for the version `v1.35.0`.

## Features
* feat(resource-groups): add per-group-type templates (#1414) (hazedav)([2ccd309f](https://github.com/lacework/go-sdk/commit/2ccd309f186297672cd9091ba2bde4b877ffc0d3))
* feat: check if cdk component is tainted (#1412) (Pengyuan Zhao)([1cea71bb](https://github.com/lacework/go-sdk/commit/1cea71bb645a8a9789f53b42338c8472a5be2d3a))
* feat(resource-groups): add support for OCI (#1410) (hazedav)([76033703](https://github.com/lacework/go-sdk/commit/7603370345edbf037d6a819b4d91cb073b55c5a9))
* feat(lql): add option to generate empty lql file (#1335) (Daniel Thorpe)([f0ad17a5](https://github.com/lacework/go-sdk/commit/f0ad17a532b357ba1e39203e6ecb13cc22cd1a53))
* feat: Components daily version check (#1407) (Linpeng Zhang)([242635d5](https://github.com/lacework/go-sdk/commit/242635d51e4fce2984e09f725140f376282a99e4))
* feat(SPM-1183): support OCI souce for alert_rules (#1405) (Pengyuan Zhao)([6fedf904](https://github.com/lacework/go-sdk/commit/6fedf9045d19792539ce4591efe1afd889429e59))
## Refactor
* refactor(state.Install): add ability to install Python oneFolder archives (#1409) (Timothy MacDonald)([50ef8812](https://github.com/lacework/go-sdk/commit/50ef88122305d2cd3a86cc15bc43367b261e6253))
## Bug Fixes
* fix(cdk): place Python component binary in component root (#1419) (Timothy MacDonald)([6809238b](https://github.com/lacework/go-sdk/commit/6809238b2efbeb2f687ce71c9862f4add6dce569))
* fix(scaffolding): Python components can't have '-' (#1391) (Salim Afiune)([4150c4d8](https://github.com/lacework/go-sdk/commit/4150c4d8a536a89a486f9fce70a08c40fe59962a))
* fix: CDK bugs (#1411) (jonathan stewart)([93923cdd](https://github.com/lacework/go-sdk/commit/93923cdd81ea401768d316a7d39156abec72652c))
* fix: CDK v1 signature verify (#1408) (jonathan stewart)([2c5afe27](https://github.com/lacework/go-sdk/commit/2c5afe27680174808ab66ccdfb8e56bd7c2d93d4))
## Other Changes
* chore(deps): bump golang.org/x/net from 0.7.0 to 0.17.0 (#1406) (dependabot[bot])([72b67b2a](https://github.com/lacework/go-sdk/commit/72b67b2a05edbe6e128dc7c2ec06e1633beece2a))
* chore(deps): bump github.com/honeycombio/libhoney-go (#1328) (dependabot[bot])([fa9a4c9c](https://github.com/lacework/go-sdk/commit/fa9a4c9cb65f7497dfb627cee5c1677298f54da2))
* chore(cli): add info about false-positive virus on Windows (#1374) (Salim Afiune)([3072d076](https://github.com/lacework/go-sdk/commit/3072d07611a711cf608ae3e30a4270659d185a68))
* ci: migrate windows pipelines to GH actions (#1417) (Salim Afiune)([8999a523](https://github.com/lacework/go-sdk/commit/8999a523bf9232148ba2a4d664417857b70934b8))
* ci: version bump to v1.34.1-dev (Lacework)([3437a181](https://github.com/lacework/go-sdk/commit/3437a181d88f9dd9b768e999ab27151c71473779))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
