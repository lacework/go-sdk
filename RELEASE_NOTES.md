# Release Notes
Another day, another release. These are the release notes for the version `v0.40.0`.

## Features
* feat(cli): alert-profile update cmd (#879) (Darren)([d8024e5](https://github.com/lacework/go-sdk/commit/d8024e5aad8ea7e240d2be2a65680351678034ce))
* feat(cdk): leave breadcrumbs after installation and update (#876) (Salim Afiune)([d16cc15](https://github.com/lacework/go-sdk/commit/d16cc153cc6483dd20e05bfd3035f5a0445ef17a))
* feat(lql): add limit option (hazedav)([d598def](https://github.com/lacework/go-sdk/commit/d598deff17ae4b644cc655c084e96fa6b5d70ad4))
* feat: add cloud account integration aws agentless scanning (#870) (Darren)([8807a97](https://github.com/lacework/go-sdk/commit/8807a97ac04a9d2637947a6b74b5f6bfd148e91e))
## Refactor
* refactor: use v2/AgentInfo for agents list command (#866) (Darren)([06531ed](https://github.com/lacework/go-sdk/commit/06531ed7ebcc8fba41d95dd998aa2e8d0125aa7c))
## Bug Fixes
* fix(lql): add line break when no preview found (#880) (hazedav)([83a1216](https://github.com/lacework/go-sdk/commit/83a1216a20162af5fdc4c30c5f0572d65a773e64))
* fix(cdk): JSON output for 'component list' command (#878) (Salim Afiune)([a1dc162](https://github.com/lacework/go-sdk/commit/a1dc162441b5ac849010dc87b5a09fbce3f13d16))
* fix(api): detect FQDN as account in api.NewClient() (#872) (Salim Afiune)([6e14d57](https://github.com/lacework/go-sdk/commit/6e14d577fc28f6e5b4eabbda0fe5f08d080e4473))
* fix: vulnerabilty_exception fixable parameter omit when unset (#863) (Darren)([f517e0e](https://github.com/lacework/go-sdk/commit/f517e0e48fe6cf0ac8f05d59d65c25dc4dbb7ddd))
* fix(install.ps1): modify User $Path correctly (#869) (Salim Afiune)([73c69f7](https://github.com/lacework/go-sdk/commit/73c69f74a4215799f9826d14badaa5214150ba97))
## Documentation Updates
* docs(cli): update help message from LPP commands (#873) (lacework-aaronscheiber)([810d005](https://github.com/lacework/go-sdk/commit/810d0056decf864ab8cd522c2f571475b825a76b))
* docs(cli): updating help for query run (#864) (hazedav)([8953d19](https://github.com/lacework/go-sdk/commit/8953d19f8538b6247f0018d66777ed891eb50c51))
## Other Changes
* chore(deps): bump github.com/honeycombio/libhoney-go (#867) (dependabot[bot])([e77d133](https://github.com/lacework/go-sdk/commit/e77d133e3d7564de0fae979b08d8fbfe05221cf4))
* chore(deps): bump github.com/BurntSushi/toml from 1.1.0 to 1.2.0 (#871) (dependabot[bot])([fa69c24](https://github.com/lacework/go-sdk/commit/fa69c24b079e09ac4012a8beca22947cc01ce686))
* chore(lql): fully remove evaluatorId (#865) (hazedav)([1112f55](https://github.com/lacework/go-sdk/commit/1112f55c47577591d94cd7719ab66c07f51f21e4))
* ci: version bump to v0.39.1-dev (Lacework)([23827ee](https://github.com/lacework/go-sdk/commit/23827ee45a92c9c23784af1e61967eff6ef9dae9))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
