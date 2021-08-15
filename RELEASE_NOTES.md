# Release Notes
Another day, another release. These are the release notes for the version `v0.13.0`.

## Features
* feat(cli): install agents via SSH on a custom port (#516) (Salim Afiune)([89ff730](https://github.com/lacework/go-sdk/commit/89ff730d03ef79a48ff505de9e1ff39ae48d5661))
* feat(cli): new agent install --trust_host_key flag (#515) (Salim Afiune)([fbd9934](https://github.com/lacework/go-sdk/commit/fbd99344d44f72e60f784695cf79fa3b56a243b5))
* feat(cli): new --api_token global flag (#513) (Salim Afiune)([9cf1271](https://github.com/lacework/go-sdk/commit/9cf12719981516e83afe4390996d42925ccde3c9))
* feat(cli): unhide query and policy commands (hazedav)([d1990d2](https://github.com/lacework/go-sdk/commit/d1990d2c065b0ea6efe0de247629b27afa9ee828))
## Bug Fixes
* fix(cli): vuln host scan-pkg-manifest --local centos 6.10 (#514) (Darren)([83f8884](https://github.com/lacework/go-sdk/commit/83f888455cb022087ad5ac28354a2a9fe4989bfe))
* fix(install.ps1): verify if running as administrator (#510) (Salim Afiune)([d9c2284](https://github.com/lacework/go-sdk/commit/d9c2284b2ba44ec9c15133c0e382348acbff98b3))
## Documentation Updates
* docs(cli): add inline help docs for query/policy cmds (#517) (Salim Afiune)([e471393](https://github.com/lacework/go-sdk/commit/e47139331a0e3482c34f39ee56b2917b69b70b33))
## Other Changes
* ci: version bump to v0.12.1-dev (Lacework)([af1cb50](https://github.com/lacework/go-sdk/commit/af1cb501c8b75034c46c1e089cc45c1c62a73f19))
* test(cli): go test modularity via build tags (#511) (hazedav)([b83a4af](https://github.com/lacework/go-sdk/commit/b83a4afb799a0f7fa1297e7ed027cf1c4c1b829a))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
