# Release Notes
Introducing the new Lacework CLI `v0.2.0`

## What's new?

### Top-level vulnerability command
The new `lacework vulnerability` command have two new sub-commands:
* `container`: for container vulnerability assessments
* `host`: for host vulnerability assessments

### Container Vulnerability
The new `lacework vulnerability container` command has the following sub-commands:
* `list-assessments`: list container vulnerability assessments (default last 7 days)
* `show-assessment`: show results of a container vulnerability assessment
* `scan`: request an on-demand container vulnerability assessment
* `scan-status`: check the status of an on-demand container vulnerability assessment

### Host Vulnerability
The new `lacework vulnerability host` command has the following sub-commands:
* `list-cves`: list the CVEs found in the hosts in your environment
* `list-hosts`: list the hosts that contain a specified CVE id in your environment
* `show-assessment`: show results of a host vulnerability assessment
* `scan-pkg-manifest`: request an on-demand host vulnerability assessment from a package-manifest

### Lacework Events
* The new `lacework event open <event_id>` allows you to quickly open events in the Lacework Console to do further investigation
* The `lacework event show <event_id>` also provides a direct link to the event in the Lacework Console
* The `lacework events list` command can now filter by event severity, as well as filter by a specific number of days (max 7 days)

## What's deprecated?
The following commands have been deprecated:

* `lacework vulnerability report` use now `lacework vulnerability container show-assessment`
* `lacework vulnerability scan run` use now `lacework vulnerability container scan`
* `lacework vulnerability scan show` use now `lacework vulnerability container scan-status`

These commands will continue to work for a number of releases showing a deprecation message similar to:
```
$ lacework vulnerability report
(DEPRECATED) This command has been moved, use now the following command:

  $ lacework vulnerability container show-assessment <sha256:hash>
```

## Features
* feat(cli): new event open command (#197) (Salim Afiune)([42e0309](https://github.com/lacework/go-sdk/commit/42e03096cf387a55329275c22a787ccf239c1baa))
* feat(cli): filter events by severity (Salim Afiune Maya)([2d8fdf4](https://github.com/lacework/go-sdk/commit/2d8fdf46b391562205d036a8f866b4e940377f9c))
* feat(cli): list events from a number of days (Salim Afiune Maya)([0474765](https://github.com/lacework/go-sdk/commit/047476548e6b86dcd249c8f37b0cfb65a49a401d))
* feat(cli): allow users to pass only --start flag (Salim Afiune Maya)([547dc1d](https://github.com/lacework/go-sdk/commit/547dc1d3a8db23e9d9b411e045b6bbce6b99e161))
* feat(cli): filter assessments for specific repos (Salim Afiune Maya)([6482d8e](https://github.com/lacework/go-sdk/commit/6482d8ea6ad712077fc595011cbdfee0715c04bc))
* feat(cli): --active & --fixable flags to container vuln (Salim Afiune Maya)([9f027b9](https://github.com/lacework/go-sdk/commit/9f027b9b56c2b4c110281246971988881f8f1164))
* feat(cli): --active & --fixable flags to host vuln (Salim Afiune Maya)([27f5197](https://github.com/lacework/go-sdk/commit/27f5197c17488a9575a8ba47f17293590a8cdbbf))
* feat(cli): add emoji support for windows (Salim Afiune Maya)([0762814](https://github.com/lacework/go-sdk/commit/07628145c9e034bc8492d9e833bf9cef962996da))
* feat(cli): add an emoji Go package for 🍺 🍕 🌮 (Salim Afiune Maya)([cafb8d8](https://github.com/lacework/go-sdk/commit/cafb8d8cf721e7d3259f7de5f06613d3136c28f0))
* feat(cli): order vulnerabilities by total of hosts (Salim Afiune Maya)([5cfe695](https://github.com/lacework/go-sdk/commit/5cfe69538cb1c869909e4b4f321eeab7c3ac1b19))
* feat(cli): new vulnerability list-assessments command (Salim Afiune Maya)([7e7191a](https://github.com/lacework/go-sdk/commit/7e7191ab1aa4b765081c91573df307d5c9113f9c))
## Refactor
* refactor(cli): container and host vulnerability cmds (Salim Afiune Maya)([c5c0117](https://github.com/lacework/go-sdk/commit/c5c0117492eec958159b13df36b738af48f5a5e0))
* refactor: host vulnerability feature (Salim Afiune Maya)([5e9f770](https://github.com/lacework/go-sdk/commit/5e9f7700acd422f5bf0b79d3faf58ffc6ed0034b))
* refactor: container vulnerability feature (Salim Afiune Maya)([bdaf126](https://github.com/lacework/go-sdk/commit/bdaf12641851b3a3bb514617ca3ae61e062bbb07))
## Performance Improvements
* perf(cli): retry polling on-demand container scan statuses (Salim Afiune Maya)([d14ea35](https://github.com/lacework/go-sdk/commit/d14ea3598c2f5d4ea795f3930c0e6b48698e9777))
## Other Changes
* chore(cli): update help messages (Salim Afiune Maya)([f1c164c](https://github.com/lacework/go-sdk/commit/f1c164c14703e6dc1faecbd566ff7be3aae822ae))
* chore(cli): consistent help message for vuln cmds (Salim Afiune Maya)([f796c58](https://github.com/lacework/go-sdk/commit/f796c5835f91c5224701e60f8236fc55e663b83e))
* chore(cli): leave breadcrumbs for host vuln cmds (Salim Afiune Maya)([45d8427](https://github.com/lacework/go-sdk/commit/45d8427554a9a74f40f3e97c2e0f8c0251a8450f))
* ci(integration): run full tests on windows (#190) (Salim Afiune)([c5c8cf4](https://github.com/lacework/go-sdk/commit/c5c8cf4c80a2fcb40e84dcefbec4f733c5d8bc52))
* test(integration): add host vulnerability tests (Salim Afiune Maya)([a5cb795](https://github.com/lacework/go-sdk/commit/a5cb7951832c4c95c64b24c80f73e06293920283))
* test(integration): add container vulnerability tests (Salim Afiune Maya)([9b2c49d](https://github.com/lacework/go-sdk/commit/9b2c49d88ca962274e145028eaebb58f88ff417b))

## Docker Images
* `docker pull techallylw/lacework-cli:latest`
* `docker pull techallylw/lacework-cli:scratch`
* `docker pull techallylw/lacework-cli:ubi-8`
* `docker pull techallylw/lacework-cli:centos-8`
* `docker pull techallylw/lacework-cli:debian-10`
* `docker pull techallylw/lacework-cli:ubuntu-1804`
* `docker pull techallylw/lacework-cli:amazonlinux-2`
