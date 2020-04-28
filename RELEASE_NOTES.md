# Release Notes
Another day, another release. These are the release notes for the version `v0.1.8`.

# Keep Users Up-To-Date
We have introduced the new **Lacework Updater** mechanism into the Lacework
CLI, we are aiming to keep our users up-to-date with our release cadence, this new
feature will check for any available update and will suggest executing the update
command.

This feature is only enabled inside the `lacework version` command:
```
$ lacework version
lacework v0.1.7 (sha:861ce9e227a3f1cd95c78011d95f0f54e6b72ec2) (time:20200427020330)

A newer version of the Lacework CLI is available! The latest version is v0.1.8,
to update execute the following command:

  $ curl https://raw.githubusercontent.com/lacework/go-sdk/master/cli/install.sh | sudo bash
```

# 📦 Happy Updates!

## Features
* feat(cli/vul): show layer content instead of hash (Salim Afiune Maya)([a15e767](https://github.com/lacework/go-sdk/commit/a15e767ca4c5d05ff3fa888ce9f1333de7b545ac))
* feat(cli): add --details flag to vulnerability cmd (Salim Afiune Maya)([227a7b2](https://github.com/lacework/go-sdk/commit/227a7b2dbceab31c8cd982f802373a83a08c7e1b))
* feat(cli): check for available updates 👓 ✨ (Salim Afiune Maya)([9318952](https://github.com/lacework/go-sdk/commit/93189526b2cc518b75c0de524e748d08b27247b5))
* feat: new go library lwupdater 🆕 ⭐ (Salim Afiune Maya)([0f7637e](https://github.com/lacework/go-sdk/commit/0f7637e31b01ee05b0b2ced2f740c3836bcddbe2))
## Refactor
* refactor(cli): consistency between image ID & Digest (Salim Afiune Maya)([4f59376](https://github.com/lacework/go-sdk/commit/4f5937672218427128501bc20637485a40692f81))
* refactor(api): request and response log messages (Salim Afiune Maya)([e4a3b3c](https://github.com/lacework/go-sdk/commit/e4a3b3ca8e4a2db9b6ef841ff33e8feab4d9bb6e))
## Bug Fixes
* fix(cli): sort vulnerabilities by severity (Salim Afiune Maya)([1e0de4c](https://github.com/lacework/go-sdk/commit/1e0de4c3a6c2107d4bffd47a9f188e07c9e3ca79))
## Documentation Updates
* docs(cli): func comments and cmd style updates (Salim Afiune Maya)([b50f987](https://github.com/lacework/go-sdk/commit/b50f987c37946dcee4579f5c6ce67e496734c12f))
* docs(lwlogger): add basic usage example (Salim Afiune Maya)([c994534](https://github.com/lacework/go-sdk/commit/c99453485c61f179c3495df03291a50ecc689947))
## Other Changes
* style(cli): align vulnerability summary report (Salim Afiune Maya)([0b37cf6](https://github.com/lacework/go-sdk/commit/0b37cf6dd480800966a271ac86f486f1d790675a))
* style(cli): remove dup vul report summary footer (Salim Afiune Maya)([6e36455](https://github.com/lacework/go-sdk/commit/6e36455dfd0f7755d3727bd9ab22dd17b25c1b86))
* style: avoid mixing duties between api/ and cli/ (Salim Afiune Maya)([fb7b7c2](https://github.com/lacework/go-sdk/commit/fb7b7c254498bb504b16c7b42580ca0c274efeda))

