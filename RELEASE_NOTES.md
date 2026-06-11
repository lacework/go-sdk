# Release Notes
Another day, another release. These are the release notes for the version `v2.14.0`.

## Features
* feat(AWLS2-1027): tenant level support for Azure DSPM (#1835) (Joseph Wilder)([bdf7c23d](https://github.com/lacework/go-sdk/commit/bdf7c23dbcca66eb4b970ed0efd6059be9d50754))
## Bug Fixes
* fix(pre-flight): stop filtering out non-Recommended Azure regions in region discovery (CAD-2132) (#1833) (Lokesh Vadlamudi)([76b64083](https://github.com/lacework/go-sdk/commit/76b64083ee5798c9b25e035da326e9dc1fda5f44))
* fix: report missing secretsmanager:TagResource in agentless preflight (CAD-2093) (#1832) (Lokesh Vadlamudi)([bf89134d](https://github.com/lacework/go-sdk/commit/bf89134d8d6d378f452a5429ff10c247b0d23b44))
## Other Changes
* chore: run go mod tidy to mark golang.org/x/sync as a direct dependency (#1836) (Lokesh Vadlamudi)([1738055b](https://github.com/lacework/go-sdk/commit/1738055bfb7a8b9a0aaacba7e2598a9f4cbba686))
* chore(CAD-2106): clean up AWS Config and CloudTrail permissions (#1831) (Pengyuan Zhao)([65b649bd](https://github.com/lacework/go-sdk/commit/65b649bd03587f69281708280d34260105384a41))
* chore: skip AWS compliance report integration tests (no CI data) (#1834) (Lokesh Vadlamudi)([b59dff0f](https://github.com/lacework/go-sdk/commit/b59dff0f6d899e599e877352fcedd1773e95a03e))
* ci: version bump to v2.13.1-dev (Lacework)([d6091405](https://github.com/lacework/go-sdk/commit/d609140568b310c37ca734a8fa1717d853434d10))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
