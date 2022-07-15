# Release Notes
Another day, another release. These are the release notes for the version `v0.39.0`.

## Features
* feat(cli): colorize CDK commands (#859) (Salim Afiune)([41e0785](https://github.com/lacework/go-sdk/commit/41e07854448b5053312bc57791c3b20a7aa4de60))
* feat(cli): fetch existing policy for cli update (#841) (Darren)([8bc93ba](https://github.com/lacework/go-sdk/commit/8bc93baa67e3f8f063afc9b4f05199801cd8c170))
* feat(cli): new policy-exceptions command (#844) (Darren)([7e894c9](https://github.com/lacework/go-sdk/commit/7e894c9d5d2f3d550f3d117db684a6116e851ecd))
* feat(cli): include status in vuln ctr scan assesment output (#853) (Darren)([5ff586c](https://github.com/lacework/go-sdk/commit/5ff586cb3a12c6f60fa79164b6a3d3e74461d6f4))
* feat(cli): display CDK version (Salim Afiune Maya)([d21bbef](https://github.com/lacework/go-sdk/commit/d21bbefeb4da27cc3036d1981c0bec44ec88d97a))
* feat(cdk): life cycle event 'cdk-init' (Salim Afiune Maya)([12eee01](https://github.com/lacework/go-sdk/commit/12eee017c12145c6bc315050803e9f15f04a651e))
* feat: support inline scanner integrations (Salim Afiune Maya)([99aacd8](https://github.com/lacework/go-sdk/commit/99aacd80e2a451f2286f49d6cd2079adf3a87105))
## Bug Fixes
* fix: omit vulnerability_exception fixable parameter when unset (#857) (Darren)([e671857](https://github.com/lacework/go-sdk/commit/e6718576f977a79bdef231c164c60d0717df290d))
* fix(cli): use proper env var name LW_API_TOKEN (#856) (Salim Afiune)([22efef4](https://github.com/lacework/go-sdk/commit/22efef47783986824d3c404cee86484660ca7ff4))
* fix(cli): install.ps1 avoid Windows agent collision (#854) (Salim Afiune)([1d52d27](https://github.com/lacework/go-sdk/commit/1d52d27057cf3e4f99896509db26a6cf2704ecc6))
## Documentation Updates
* docs(api): fix path to examples in readme (#858) (Yasmine Hal)([a8a5feb](https://github.com/lacework/go-sdk/commit/a8a5febb8a8cdd373e4e07cde64041f4d0abbd4f))
## Other Changes
* chore(cli): update content-library component name  (#855) (hazedav)([2794070](https://github.com/lacework/go-sdk/commit/27940704c140ca229235f8fc1ac8b1c81f1ad269))
* ci: version bump to v0.38.1-dev (Lacework)([196406a](https://github.com/lacework/go-sdk/commit/196406aedce90674661551ac146f6cdec2d8d634))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
