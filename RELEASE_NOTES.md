# Release Notes
Another day, another release. These are the release notes for the version `v1.42.0`.

## Features
* feat(cli): Add time filter to vulnerability host list-cve command (#1481) (kendall-choy-lw)([19a61c8d](https://github.com/lacework/go-sdk/commit/19a61c8d939578e0a1b9809ec7825cad28915bad))
* feat(cli): show query language in policies when available (#1475) (Kun Chen)([b754b4ee](https://github.com/lacework/go-sdk/commit/b754b4ee59624af7d4b1f33eb44bc148a4d9f3d8))
* feat(RAIN-91461): support Rego in query validate (#1472) (Kun Chen)([6ed380ec](https://github.com/lacework/go-sdk/commit/6ed380ecb97c4f8771d893e23dcd00e23de406e3))
* feat(cli): Select query language only on query creation (#1467) (Kun Chen)([89986135](https://github.com/lacework/go-sdk/commit/899861358b76f4ced74c3518d5c0d71016d936da))
## Bug Fixes
* fix: improve error message when token generate fails in cli (#1482) (Matt Cadorette)([8814db38](https://github.com/lacework/go-sdk/commit/8814db38d054650fe9c76f28de83f1ac7e70978b))
* fix: temporarily install Semgrep to fix SAST (#1485) (Jeremy Dubreil)([38149934](https://github.com/lacework/go-sdk/commit/3814993407121d084834d592dc4596f5092adae0))
* fix: always set alias "main" for the main profile (#1483) (Pengyuan Zhao)([985c9e04](https://github.com/lacework/go-sdk/commit/985c9e040957acb3f9d84f4fb8035adaeaaa17a1))
## Other Changes
* chore(deps): bump google.golang.org/grpc from 1.51.0 to 1.56.3 (#1486) (dependabot[bot])([28b3debb](https://github.com/lacework/go-sdk/commit/28b3debbdb29739cafcdc2ce49f45bcd660714f8))
* ci: fix TestContainerVulnerabilityCommandsEndToEnd (#1484) (Salim Afiune)([0f1e35d4](https://github.com/lacework/go-sdk/commit/0f1e35d4ba62134b0c4582b07ddaea2863b0cb2d))
* ci: add GitHub workflow for SCA (#1478) (Pengyuan Zhao)([904f51cf](https://github.com/lacework/go-sdk/commit/904f51cfbce4905c418fbb3d9c4a851445b845c0))
* ci: version bump to v1.41.2-dev (Lacework)([f37368c9](https://github.com/lacework/go-sdk/commit/f37368c96a42f5830e62d9943a982a4fd556eedd))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
