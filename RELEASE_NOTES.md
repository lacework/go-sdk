# Release Notes
Another day, another release. These are the release notes for the version `v1.38.0`.

## Features
* feat(PSP-1220): add new cloud account type for Azure Agentless integration (#1443) (AoZhang)([1bea0128](https://github.com/lacework/go-sdk/commit/1bea012813e775ae22c76fa2dc0f8f42a0da16a6))
* feat: support Agentless origanization integration with auto snapshot (#1424) (Pengyuan Zhao)([6084c103](https://github.com/lacework/go-sdk/commit/6084c103c6c92d75b2a725c8e1a3961d5b01e51e))
* feat(Makefile): build CLI simpler & faster (Salim Afiune Maya)([d5969abb](https://github.com/lacework/go-sdk/commit/d5969abbff31168e7efb646bd1e58d86add46bf0))
## Bug Fixes
* fix(lwtime): fix issues related to daylight savings (#1444) (hazedav)([755a47af](https://github.com/lacework/go-sdk/commit/755a47afcaed1b50106a4ac46805ab9d3fbb5ed8))
* fix(component): List installed version instead of latest (Salim Afiune Maya)([d92b2af6](https://github.com/lacework/go-sdk/commit/d92b2af66009b6692f9bd3ddcb2f6f069d788eb6))
## Other Changes
* chore: Implement --aws_assume_role flag for CLI aws integration (#1434) (Pengyuan Zhao)([a6bcd7f9](https://github.com/lacework/go-sdk/commit/a6bcd7f98776725cc4a06af3a4933cb958801522))
* build(release): fix chocolately release action (#1440) (Timothy MacDonald)([554f9742](https://github.com/lacework/go-sdk/commit/554f97421b5c720317aae068628fa28f9e5d2613))
* ci: version bump to v1.37.1-dev (Lacework)([ad406eaf](https://github.com/lacework/go-sdk/commit/ad406eaf329eda76bec39da65adb97d6ceceec78))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
