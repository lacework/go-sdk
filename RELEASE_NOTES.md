# Release Notes
Another day, another release. These are the release notes for the version `v0.2.12`.

## Install Lacework Agent On Remote Hosts
Users can now install agents with the command:
```
$ lacework agent install
```

When this command is executed without any additional flag, an interactive prompt will be
launched to help gather the necessary authentication information to access the remote host.

To authenticate to the remote host with a username and password.
```
$ lacework agent install <host> --ssh_username <your-user> --ssh_password <secret>
```

To authenticate to the remote host with an identity file.

```
$ lacework agent install <user@host> -i /path/to/your/key
```

To provide an agent access token of your choice, use the command `lacework agent token list`,
select a token and pass it to the agent install command using the `--token` flag.

```
$ lacework agent install <user@host> -i /path/to/your/key --token <token>
```

### Agent Install Example

![lacework-cli-agent-install-release](https://user-images.githubusercontent.com/5712253/103636405-26d83100-4f07-11eb-8bcf-0784a5abd8af.gif)

## Features
* feat(cli): add telemetry (#278) (Salim Afiune)([5aeec3c](https://github.com/lacework/go-sdk/commit/5aeec3c51184fc7e43e1c9dc413d256c98b8c516))
* feat(cli): pull latest agent version from S3 (Salim Afiune Maya)([63cf1ab](https://github.com/lacework/go-sdk/commit/63cf1ab82933600189904abe0b25958769a42ec9))
* feat: add --force to agent install (Salim Afiune Maya)([6de4775](https://github.com/lacework/go-sdk/commit/6de47756973f3396b9d3f5d6e044db3308e1700a))
* feat: verify if agent is installed on remote host (Salim Afiune Maya)([252b9a6](https://github.com/lacework/go-sdk/commit/252b9a602781a68ee88d1d0c9e14ee290c310a79))
* feat(cli): check for known hosts and allow custom callbacks (Salim Afiune Maya)([ebedf22](https://github.com/lacework/go-sdk/commit/ebedf221f4a1569080aeaf8de1441661845d22b2))
* feat: add AWS S3 alert channel integration (#273) (Darren)([383de18](https://github.com/lacework/go-sdk/commit/383de18bedfa1d85eb140f5b82ecb2c69ba231be))
* feat(cli): enable agent install command (Salim Afiune Maya)([f13d58a](https://github.com/lacework/go-sdk/commit/f13d58a2bbedf7772ddd63330a4cb813f926f541))
## Refactor
* refactor: verify host connectivity before select token (Salim Afiune Maya)([829cf82](https://github.com/lacework/go-sdk/commit/829cf821d457e5178c13e3d98bd9f31c60be3ded))
* refactor(api): remove automatic report trigger (#271) (Salim Afiune)([18e624f](https://github.com/lacework/go-sdk/commit/18e624f74e68fddc2f180e5e608353a824bac9b7))
## Bug Fixes
* fix(cli): propagate errors from install.sh (#277) (Salim Afiune)([296be65](https://github.com/lacework/go-sdk/commit/296be658d106ad84cf9a4a3ced1d4f6122ce4db8))
* fix(cli): avoid showing unnamed tokens (Salim Afiune Maya)([7545444](https://github.com/lacework/go-sdk/commit/754544441972f73a55181a4255453f6f911f81d0))
## Documentation Updates
* docs: update agent install use (Salim Afiune Maya)([62195c1](https://github.com/lacework/go-sdk/commit/62195c1a2b429b02120a8d797e0debaa448016e8))
## Other Changes
* chore: update long desc of agent list cmd (Salim Afiune Maya)([8a24914](https://github.com/lacework/go-sdk/commit/8a2491456d361d22de4760a79abfbbb0dcc51559))
* build: stop publishing containers to old docker repo (Salim Afiune Maya)([ea23a30](https://github.com/lacework/go-sdk/commit/ea23a3085e4c8ef35acc4fe06d3ba972be4d932a))
* ci: send slack notifications to team alias ⭐ (Salim Afiune Maya)([5e4c0e6](https://github.com/lacework/go-sdk/commit/5e4c0e69824ef00289e4d86adecf48209709bb59))
* ci: fix mv on non exisitent directory (#272) (Darren)([4f101cf](https://github.com/lacework/go-sdk/commit/4f101cfe8c8aeff5981264e99bdb411b548e02e9))
* test(cli): increase agent install test coverage (#276) (Salim Afiune)([da5b4ae](https://github.com/lacework/go-sdk/commit/da5b4aea9730c55c10d541c976dcb7ccf16aca28))
* test: fix lwrunner tests (Salim Afiune Maya)([23587cd](https://github.com/lacework/go-sdk/commit/23587cdd98c694e65a8f0791c269817ce7252d4c))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
