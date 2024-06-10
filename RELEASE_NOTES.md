# Release Notes
Another day, another release. These are the release notes for the version `v1.51.0`.

## Features
* feat: prosast component data upload (#1639) (jonathan stewart)([635f2d30](https://github.com/lacework/go-sdk/commit/635f2d30643ce70c3ebd6b0c49542954a5e7a828))
* feat(GROW-2928): Upgrade the gcp terraform version (#1637) (Lei Jin)([02bcccab](https://github.com/lacework/go-sdk/commit/02bcccabd229620b4d08755b95615d4f4c17a71a))
* feat(GROW-2931): Add support for GCP default labels (#1633) (Lei Jin)([6d765866](https://github.com/lacework/go-sdk/commit/6d765866a99a07e3cb0bf0ccc9eaa73612b2f3c4))
* feat(GROW-2929): update lwgenerate to support GCP generation (#1631) (Lei Jin)([2c64710e](https://github.com/lacework/go-sdk/commit/2c64710ee725d77d7341274e940dc5cec27d3c3d))
* feat(GROW-2819): enable custom blocks with azure generate (#1630) (Matt Cadorette)([67c7bbc9](https://github.com/lacework/go-sdk/commit/67c7bbc9b88a7492fa3d387d44d5e1f290eef79e))
* feat(GROW-2908): lwgenerate enable custom root terraform blocks and provider arguments (#1626) (Matt Cadorette)([8d50d317](https://github.com/lacework/go-sdk/commit/8d50d317e44bf379523c0138fcb12ec34b4ed302))
* feat(GROW-2883): lwgenerate enable adding aws provider default tags (#1623) (Matt Cadorette)([5d6ef3df](https://github.com/lacework/go-sdk/commit/5d6ef3dfbbc1349abaf5e1c0fca325f0ac65960a))
## Bug Fixes
* fix(GROW-2931): GCP default label should be attributes (#1636) (Lei Jin)([e575eac8](https://github.com/lacework/go-sdk/commit/e575eac82b34e8a9ce840928d6f73e13675291e0))
* fix(GROW-2819): split azuread/azurerm gen args (#1632) (Matt Cadorette)([6ae460a7](https://github.com/lacework/go-sdk/commit/6ae460a7212365947c80889ec0ff02499a05bdd5))
## Other Changes
* chore: update cloud-account migrate cmd desc (#1625) (Manan Bhatia)([8c76d482](https://github.com/lacework/go-sdk/commit/8c76d482af778977c8eeb35c150fea9e1fed7be1))
* ci: version bump to v1.50.1-dev (Lacework)([f469da52](https://github.com/lacework/go-sdk/commit/f469da5287671e0ef893df381b6e64ae4530c50d))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
