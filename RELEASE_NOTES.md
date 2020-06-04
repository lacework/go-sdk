# Release Notes
Another day, another release. These are the release notes for the version `v0.1.14`.

## Display vulnerability reports with 0 vulnerabilities
With this new release, when you scan a container image that has 0 vulnerabilities, you
will now get the following output as a result:
```
$ lacework vul report sha256:32c5a42735581453d1e9ea664604fa81f4b635459dc88799ee1f8c373dca4819
Great news! This container image has no vulnerabilities.
```

And for the JSON format:
```
$ lacework vul report sha256:32c5a42735581453d1e9ea664604fa81f4b635459dc88799ee1f8c373dca4819 --json
{
  "critical_vulnerabilities": 0,
  "fixable_vulnerabilities": 0,
  "high_vulnerabilities": 0,
  "info_vulnerabilities": 0,
  "low_vulnerabilities": 0,
  "medium_vulnerabilities": 0,
  "scan_status": "Success",
  "total_vulnerabilities": 0
}
```

## Features
* feat: understand vuln reports with 0 vulnerabilities (#124) (Salim Afiune)([6af13b0](https://github.com/lacework/go-sdk/commit/6af13b06ac04ff8b2efb156248a70fbb50908dde))
* feat: auto-populate account with --profile flag (#121) (Salim Afiune)([3539ec4](https://github.com/lacework/go-sdk/commit/3539ec409285a7d3f0335e6bfc2676f03c5fbb4c))
## Bug Fixes
* fix(spelling) Fixes event header misspelling (Scott Ford)([e55a6c1](https://github.com/lacework/go-sdk/commit/e55a6c16f93059d93c8ce0985a16d5bf4a7ad020))
* fix(release): update release link and version message (#117) (Salim Afiune Maya)([2969722](https://github.com/lacework/go-sdk/commit/2969722f94745fe348cc9c58d1c08ae22b81cf23))
## Documentation Updates
* doc: update cli documentation cli/README.md (#125) (Salim Afiune)([e31c4fc](https://github.com/lacework/go-sdk/commit/e31c4fc7bacaa22afa734fb35885b1eff056b98d))
## Other Changes
* chore: fix typos in AWS events (#129) (Salim Afiune)([46d1bb6](https://github.com/lacework/go-sdk/commit/46d1bb69203344b784976f1fb00537a65374ab69))
* chore: bump version to v0.1.14-dev (Salim Afiune Maya)([8e7ac41](https://github.com/lacework/go-sdk/commit/8e7ac41badd51ffc1287088ca525419d6bfb5ba2))
* ci: switch Shippable in favor of CircleCI (#120) (Salim Afiune Maya)([630e8bf](https://github.com/lacework/go-sdk/commit/630e8bf308d5c944ccccd8311a566d859891a927))

## Docker Images
* `docker pull techallylw/lacework-cli:latest`
* `docker pull techallylw/lacework-cli:scratch`
* `docker pull techallylw/lacework-cli:ubi-8`
* `docker pull techallylw/lacework-cli:centos-8`
* `docker pull techallylw/lacework-cli:debian-10`
* `docker pull techallylw/lacework-cli:ubuntu-1804`
* `docker pull techallylw/lacework-cli:amazonlinux-2`
