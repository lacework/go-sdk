# Release Notes
Another day, another release. These are the release notes for the version `v0.1.9`.

## 🐳 Ship it Small & Fast!
From this release on, we are shipping the Lacework CLI in a **7MB** Docker container!

Snip it:
```
$ docker run techallylw/lacework-cli

The Lacework Command Line Interface is a tool that helps you manage the
Lacework cloud security platform. Use it to manage compliance reports,
external integrations, vulnerability scans, and other operations.

Start by configuring the Lacework CLI with the command:

    $ lacework configure

This will prompt you for your Lacework account and a set of API access keys.

Usage:
  lacework [command]

Available Commands:
  api           helper to call Lacework's RestfulAPI
  configure     configure the Lacework CLI
  event         inspect Lacework events
  integration   manage external integrations
  version       print the Lacework CLI version
  vulnerability view vulnerability reports and run on-demand scans

Flags:
  -a, --account string      account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
  -k, --api_key string      access key id
  -s, --api_secret string   secret access key
      --debug               turn on debug logging
      --json                switch commands output from human-readable to json format
      --nocolor             turn off colors
      --noninteractive      disable interactive progress bars (i.e. 'spinners')
  -p, --profile string      switch between profiles configured at ~/.lacework.toml

Use "lacework [command] --help" for more information about a command.

```

## Features
* feat: Add lacework-cli containers (Salim Afiune Maya)([73cdda0](https://github.com/lacework/go-sdk/commit/73cdda0413c56401e349162c04da261fe4e32bc7))
* feat(cli): create Azure integrations (Salim Afiune Maya)([29105e7](https://github.com/lacework/go-sdk/commit/29105e7fc85315b8c718906454af74245889f2a9))
* feat(cli): create GCP integrations (Salim Afiune Maya)([b2154a1](https://github.com/lacework/go-sdk/commit/b2154a16aa6d647514353c2a2d67c14cef9b608f))
* feat(cli): create AWS CloudTrail integrations (Salim Afiune Maya)([7e80795](https://github.com/lacework/go-sdk/commit/7e8079589f3f0d36c90f3e33c08ae7f168e13774))
* feat(cli): create integration sub-command (Salim Afiune Maya)([9842a0d](https://github.com/lacework/go-sdk/commit/9842a0db14cc059de9dd950408d2efc97de4b02a))
* feat(api): create container registry integrations (Salim Afiune Maya)([e33613d](https://github.com/lacework/go-sdk/commit/e33613ddcd10176464dfbcc02f09e986a5c5de01))
* feat(cli): delete external integrations (Salim Afiune Maya)([fe802b4](https://github.com/lacework/go-sdk/commit/fe802b45a05b70034d28bce8949362ba592aec2b))
## Refactor
* refactor(cli): new configure command using survey (Salim Afiune Maya)([d311ed4](https://github.com/lacework/go-sdk/commit/d311ed48ad758a48fc687db96b1ad5b2815cfeb6))
## Other Changes
* style: avoid mixing duties between api and cli (Salim Afiune Maya)([b245d9f](https://github.com/lacework/go-sdk/commit/b245d9f63765fdf7fb131bf933a762f9220969c8))
* style(cli): use appropriate icons per platform (Salim Afiune Maya)([c3e051e](https://github.com/lacework/go-sdk/commit/c3e051ed0124386796bf49d6addbad31c4d26ba4))
* chore(cli): update int create usage message (Salim Afiune Maya)([0959618](https://github.com/lacework/go-sdk/commit/095961838afce65a43ebf34b3405bb5b0fa09f80))
* chore(deps): remove promptui in favor of survey (Salim Afiune Maya)([0c663aa](https://github.com/lacework/go-sdk/commit/0c663aa23e1773aeec4162d8bf78aaadcf8f19b8))
