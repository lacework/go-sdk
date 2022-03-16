# Release Notes
Another day, another release. These are the release notes for the version `v0.27.0`.

## Features
* feat(cli): update query by loading it from the platform (#721) (Salim Afiune)([2f68adb](https://github.com/lacework/go-sdk/commit/2f68adbac8e63946f3f8e6bf6a92a5837ca8fed4))
* feat(cli): sort query and policy list output (hazedav)([68de649](https://github.com/lacework/go-sdk/commit/68de649549139b292d90ed41da87e0e902b12048))
* feat(cli): Add breadcrumb for "query show-source" (hazedav)([062421c](https://github.com/lacework/go-sdk/commit/062421ccfbbd25e06ffffa930338f5fe4afddafa))
* feat(lwgenerate): Add initial hcl generation support Azure - ALLY-394 (#682) (djmctavish)([66845d2](https://github.com/lacework/go-sdk/commit/66845d264ee2570ed7f7592e5a7e5af1cb517e3b))
## Refactor
* refactor(cli): further deprecate evaluatorId (hazedav)([7a63eca](https://github.com/lacework/go-sdk/commit/7a63ecae49c82583f2dc1cd7f8c4d38f381ce571))
* refactor(cli): Create cloud specific TF file (#712) (Ross)([3f9e0d8](https://github.com/lacework/go-sdk/commit/3f9e0d847404026b753d676a8f2f75f5096361e7))
* refactor(cli): vulnerability_expections packages human readable output (#711) (Darren)([c37afbf](https://github.com/lacework/go-sdk/commit/c37afbfd70d9d7f8cfcabbdf6e78b3a6ce89c3d7))
## Bug Fixes
* fix: use online machine for testing (#722) (chi-dev-2021)([e960b64](https://github.com/lacework/go-sdk/commit/e960b6437df33495287c32799d408839de8c3dae))
* fix(cli): delete listed twice on help (hazedav)([d606f9d](https://github.com/lacework/go-sdk/commit/d606f9d9d9531ccb94897086a93197d81c1e2d67))
## Other Changes
* chore(cli): run regex of integration test (#729) (Salim Afiune)([03d55cf](https://github.com/lacework/go-sdk/commit/03d55cf3025cdac9f3a0ed047bba2db3ec465ba3))
* build: deprecate centos-8 images (#709) (Salim Afiune)([5636bdd](https://github.com/lacework/go-sdk/commit/5636bdd699e10870da75dbd21be6d4a3ea64155c))
* ci: fix release pipeline (#730) (Darren)([6470a4a](https://github.com/lacework/go-sdk/commit/6470a4a12588681149d0c59a1cd27fb4a7f88a1b))
* ci: version bump to v0.26.1-dev (Lacework)([dccd3f5](https://github.com/lacework/go-sdk/commit/dccd3f5cad014dfdfa26040cf1f4af936ca15ce3))
* test(windows): fix TestQueryUpdateFromIDEditor (#728) (Salim Afiune)([424c939](https://github.com/lacework/go-sdk/commit/424c939f7aedd0b623b652965798903298ddf50a))
* test(integration): fix finding host id (#725) (Salim Afiune)([4e246d9](https://github.com/lacework/go-sdk/commit/4e246d972a1e0683b0f14a95340189668a368ba6))
* test(lql): disable policies tests (hazedav)([2647565](https://github.com/lacework/go-sdk/commit/2647565a354c33b46329b01e5db2172f836c0d04))
* test(cli): add comprehensive help text testing (#715) (hazedav)([b38c696](https://github.com/lacework/go-sdk/commit/b38c69625edc6561f2b0da0c455d5d30d04f7177))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
