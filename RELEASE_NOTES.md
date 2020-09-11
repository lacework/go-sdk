# Release Notes
Another day, another release. These are the release notes for the version `v0.2.1`.

## Features
* feat(ctr): use new lacework/lacework-cli repository (#206) (Salim Afiune)([fa1e268](https://github.com/lacework/go-sdk/commit/fa1e2682422f03288c53350d4fc6691bea6869c5))
* feat: add DockerV2, ECR and GCR container registries (#205) (Salim Afiune)([18a8c8b](https://github.com/lacework/go-sdk/commit/18a8c8b60ef6c869bcb3a72b870d4bcfd66ee794))
* feat: add decoder for jira custom_template_file (#201) (Salim Afiune)([2630ab5](https://github.com/lacework/go-sdk/commit/2630ab5fb746a8a3b4995734a79489554ff4f682))
* feat(cli): ask for JIRA Custom Template file 🚨 (Salim Afiune Maya)([5a4eb17](https://github.com/lacework/go-sdk/commit/5a4eb173b26eaff1845a33b586f9b87bbb59f449))
* feat(api): encode custom_template_file for Jira int (Salim Afiune Maya)([887ca15](https://github.com/lacework/go-sdk/commit/887ca157f8484cfc971b83f1d4e65f0fa2f10382))
## Documentation Updates
* docs(typo) fix spelling of visualize for compliance help command (#204) (Scott Ford)([75e0348](https://github.com/lacework/go-sdk/commit/75e03488a01704837498f630c4c0323d6a3ee6ef))
## Other Changes
* chore(api): remove MinAlertSeverity field from examples/ (Salim Afiune Maya)([274b8e9](https://github.com/lacework/go-sdk/commit/274b8e927b62d395074ba79628ba4ad8abdd5905))
* ci(cli): fix event time range test (Salim Afiune Maya)([9c2336b](https://github.com/lacework/go-sdk/commit/9c2336b9d7b2f96e3b8f8e22d57ba2d4fa77583b))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
