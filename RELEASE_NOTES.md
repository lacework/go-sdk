# Release Notes
Another day, another release. These are the release notes for the version `v1.28.0`.

## Features
* feat: lwgenerate for aws controltower module (#1327) (Darren)([4c1c21c1](https://github.com/lacework/go-sdk/commit/4c1c21c1be37f1ebb2966f48e2318b2332602dad))
* feat(api): support new fields in alert_rules api (#1331) (Darren)([c23859c7](https://github.com/lacework/go-sdk/commit/c23859c77317fc2d8941dffde11bf9be20fb3e15))
* feat(lwgenerate): support complex/nested module attribute values (#1332) (Matt Cadorette)([7cbbb849](https://github.com/lacework/go-sdk/commit/7cbbb84946049bf9ad6c21ad2240d1db6cc6684e))
## Bug Fixes
* fix(COD-1559): Add some retrying in `updater.go` (#1350) (Edoardo Pirovano)([72d720cd](https://github.com/lacework/go-sdk/commit/72d720cd722efb9b4c6ea9a1ec99e0ced78184b0))
* fix: alert rules empty arrays (#1347) (jonathan stewart)([2910edf2](https://github.com/lacework/go-sdk/commit/2910edf2568642cddde8ffbcd19f9757e75c8a90))
* fix: remove HostVulnerability & ContainerVulnerability values (#1334) (Darren)([26c627dc](https://github.com/lacework/go-sdk/commit/26c627dccad19e15d1dd8821ba6eb8dc8dfbf54c))
## Documentation Updates
* docs: https://pkg.go.dev/github.com/lacework/go-sdk (#1323) (Salim Afiune)([4f753a68](https://github.com/lacework/go-sdk/commit/4f753a680e825749451e950677a00b68cd938147))
## Other Changes
* style(fmt): fixed imports-check directive (#1349) (Salim Afiune)([b16e7ef5](https://github.com/lacework/go-sdk/commit/b16e7ef5009bab1174a2f1d9b1ae176fa5083b62))
* style(lql.go): update default LQL query to match policy eng requirements (#1333) (Daniel Thorpe)([77001409](https://github.com/lacework/go-sdk/commit/770014098d11c7bbc63c08b7e42351186b62bbe4))
* ci: version bump to v1.27.3-dev (Lacework)([2557e79e](https://github.com/lacework/go-sdk/commit/2557e79e91672c5523773cc9db1c0117d6fb9116))
* test: fix VulnerabilityAssessment test (#1343) (Darren)([182e9ad7](https://github.com/lacework/go-sdk/commit/182e9ad750300772ec4279e18b08d5fc3c73a6c7))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
