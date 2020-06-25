# v0.1.20

## Features
* feat(cli): add time range flags to events list cmd (#154) (Salim Afiune)([e055bc0](https://github.com/lacework/go-sdk/commit/e055bc045509620239600d4f35087817ee5d7fdc))
---
# v0.1.19

## Features
* feat(cli): set User-Agent header (backend metrics) (Salim Afiune Maya)([bb4cfc8](https://github.com/lacework/go-sdk/commit/bb4cfc81d0176bda39bb67e4bcdb3ebb422f8110))
* feat: inject client version into User-Agent header (Salim Afiune Maya)([87261d2](https://github.com/lacework/go-sdk/commit/87261d2a356b3e92dc0979c6ae6070d6558d1bf4))
* feat(api): set User-Agent header (backend metrics) (Salim Afiune Maya)([5c5001b](https://github.com/lacework/go-sdk/commit/5c5001b340f3c8e19a9ff131dab939d36f263bdd))
---
# v0.1.18

## Features
* feat(cli): add --fixable flag to vulnerability cmd (#148) (Salim Afiune)([d649e2a](https://github.com/lacework/go-sdk/commit/d649e2a754be958e8504347c68ea1286dc16a58e))
---
# v0.1.17

## Other Changes
* ci: fix vuln scan cli matrix (#143) (Salim Afiune)([646faac](https://github.com/lacework/go-sdk/commit/646faacc762b1f361de3bc61d2e543db9b674c3c))
* ci: fix release commit message (#144) (Salim Afiune)([6c6f357](https://github.com/lacework/go-sdk/commit/6c6f357d4cd1e6dae08cf55e637cea4ca56aebaa))
---
# v0.1.16

## Other Changes
* ci: dogfooding lacework vulnerability scans (orb) (Salim Afiune Maya)([e74a188](https://github.com/lacework/go-sdk/commit/e74a18814127127395f496de908ec8bb4cb22072))
* ci: build/release docker containers automatically (Salim Afiune Maya)([897b05a](https://github.com/lacework/go-sdk/commit/897b05ae9dba9eb12e44d9a09bf48092f2af3764))
---
# v0.1.15

## Bug Fixes
* fix: vulnerability scans of unsupported images (Salim Afiune Maya)([3d33a78](https://github.com/lacework/go-sdk/commit/3d33a78baa23cd024b4e9afcd2bbaa3652274967))
## Other Changes
* chore(cli): remove deprecated --digest flag (Salim Afiune Maya)([aaecce1](https://github.com/lacework/go-sdk/commit/aaecce1e815ae89c761c12842bf227156432a889))
* ci: update release process to be automated v.1 (#134) (Salim Afiune)([374b4b0](https://github.com/lacework/go-sdk/commit/374b4b01180985fb721a250efe463eed36474286))
* ci: create release from git tag (Salim Afiune Maya)([ec95742](https://github.com/lacework/go-sdk/commit/ec95742ca8f0ef117d96a0c4d2d18e96fd5304c6))
* ci: upload artifacts to release (#140) (Salim Afiune)([7e8e03f](https://github.com/lacework/go-sdk/commit/7e8e03f5635a4ceddd45ed4caf2a133f646b4803))
* ci: add slack notifications (Salim Afiune Maya)([d7523b8](https://github.com/lacework/go-sdk/commit/d7523b8a593c2ca78ef46bcf84aa0c6400bc8d10))
* ci: enable integration tests in CircleCI (Salim Afiune Maya)([a17c238](https://github.com/lacework/go-sdk/commit/a17c238bf397b6cad2036d299971c672fd116b09))
---
# v0.1.14

## Features
* feat: understand vuln reports with 0 vulnerabilities (#124) (Salim Afiune)([6af13b0](https://github.com/lacework/go-sdk/commit/6af13b06ac04ff8b2efb156248a70fbb50908dde))
* feat: auto-populate account with --profile flag (#121) (Salim Afiune)([3539ec4](https://github.com/lacework/go-sdk/commit/3539ec409285a7d3f0335e6bfc2676f03c5fbb4c))
## Bug Fixes
* fix(spelling) Fixes event header misspelling (Scott Ford)([e55a6c1](https://github.com/lacework/go-sdk/commit/e55a6c16f93059d93c8ce0985a16d5bf4a7ad020))
* fix(release): update release link and version message (#117) (Salim Afiune Maya)([2969722](https://github.com/lacework/go-sdk/commit/2969722f94745fe348cc9c58d1c08ae22b81cf23))
## Documentation Updates
* doc: update cli documentation cli/README.md (#125) (Salim Afiune)([e31c4fc](https://github.com/lacework/go-sdk/commit/e31c4fc7bacaa22afa734fb35885b1eff056b98d))
## Other Changes
* chore: fix typos in AWS events (#129) (Salim Afiune)([46d1bb6](https://github.com/lacework/go-sdk/commit/46d1bb69203344b784976f1fb00537a65374ab69))
* chore: bump version to v0.1.14-dev (Salim Afiune Maya)([8e7ac41](https://github.com/lacework/go-sdk/commit/8e7ac41badd51ffc1287088ca525419d6bfb5ba2))
* ci: switch Shippable in favor of CircleCI (#120) (Salim Afiune Maya)([630e8bf](https://github.com/lacework/go-sdk/commit/630e8bf308d5c944ccccd8311a566d859891a927))
---
# v0.1.13

## Features
* feat(cli): avoid displaying API key secret (#115) (Salim Afiune Maya)([3305b09](https://github.com/lacework/go-sdk/commit/3305b095fb43a3352255e472f38ba8f19b6d7c4b))
* feat(release): add version bump after release (Salim Afiune Maya)([4c67b3f](https://github.com/lacework/go-sdk/commit/4c67b3fbb74fa9a05db1a712c73d1570246ffc89))
## Bug Fixes
* fix(release): purge the docker manifest to udate (Salim Afiune Maya)([ed58109](https://github.com/lacework/go-sdk/commit/ed58109a5ea45b7e7b7f4d9fde86f81e183f726b))
---
# v0.1.12

## Features
* feat(cli): manage compliance reports (GCP Azure AWS) (Salim Afiune Maya)([1d0155f](https://github.com/lacework/go-sdk/commit/1d0155f48ca4dee6a4f9381870645f3c07597dff))
* feat(api): add compliance service (Salim Afiune Maya)([862812c](https://github.com/lacework/go-sdk/commit/862812c4635ded3647f3e7b76e2807de06c652ba))
* feat(cli): list integrations of a specific type (Salim Afiune Maya)([e1d3674](https://github.com/lacework/go-sdk/commit/e1d36740f7d7fe496f7746624519c81a670d054a))
## Documentation Updates
* docs(cli): remove the need to install using sudo (Salim Afiune Maya)([4534c57](https://github.com/lacework/go-sdk/commit/4534c576779ca769d053c7c19e85a6029741810e))
## Other Changes
* ci: fix typo in release.sh script (Salim Afiune Maya)([cf6a836](https://github.com/lacework/go-sdk/commit/cf6a8369e2a6b906fb604afc6213cf7c04df8095))
* ci: add docker images to release notes (Salim Afiune Maya)([4f8f945](https://github.com/lacework/go-sdk/commit/4f8f945f49d2af51856617d994cd031b02ba6678))
* test(integration): add compliance tests (Salim Afiune Maya)([d41fb49](https://github.com/lacework/go-sdk/commit/d41fb49838a7c7990acd4b7f4fd40f0a98f2452a))
---
# v0.1.11

## Features
* feat: incident analysis, visualize event details (Salim Afiune Maya)([532f11d](https://github.com/lacework/go-sdk/commit/532f11d461759c9214730a1ec5b92d9ad39afbaf))
## Bug Fixes
* fix(api): use correct types on events response (Salim Afiune Maya)([86d8b7b](https://github.com/lacework/go-sdk/commit/86d8b7b533ef77f4b9bcf63fc839ae88be12000b))
## Other Changes
* style(cli): show help without errors (Salim Afiune Maya)([a72ba55](https://github.com/lacework/go-sdk/commit/a72ba55a1a35e9c0e9626d8af4c9e1ea102c6e7c))
* ci: add badge to README and encrypted keys (Salim Afiune Maya)([c03a416](https://github.com/lacework/go-sdk/commit/c03a41664771d6a0fcfc858223e99a347b506a20))
* test(integration): adds end-to-end tests (Salim Afiune Maya)([e2eb449](https://github.com/lacework/go-sdk/commit/e2eb4493bfaf73f575a3e0c1297ba4186ace34ec))
* test(integration): new framework to write CLI tests (Salim Afiune Maya)([402b2a2](https://github.com/lacework/go-sdk/commit/402b2a28d05a5f5bf8bfd198145d091feb2461fe))
---
# v0.1.10

## Features
* feat(cli): add aliases to integration and event cmds (Salim Afiune Maya)([9e8cd5c](https://github.com/lacework/go-sdk/commit/9e8cd5c4d2eb0d9cbed715a89985978e62eab9c0))
* feat(cli): preconfigure using key JSON file (WebUI) (Salim Afiune Maya)([80c48e7](https://github.com/lacework/go-sdk/commit/80c48e7bbaf95c888b9422249c8e09818c0a83b2))
* feat(cli): new 'integration show' cmd  (#91) (Salim Afiune Maya)([5bedf53](https://github.com/lacework/go-sdk/commit/5bedf5348c9fcc1748bc66534d8ac2e6475e6c64))
## Bug Fixes
* fix(docker): fix build/release of CLI containers (Salim Afiune Maya)([2146ecb](https://github.com/lacework/go-sdk/commit/2146ecbd6c0d4d0a9f8f608a902aeffebdce3cf9))
* fix(api): parsing event details 'cpu_percentage' (Salim Afiune Maya)([5f978ea](https://github.com/lacework/go-sdk/commit/5f978ead44bd6700f520ccb0742d5355464cfece))
## Other Changes
* chore: consistency with ID fields in Go structs (Salim Afiune Maya)([79b874e](https://github.com/lacework/go-sdk/commit/79b874ed3410b033b52a59c4fa98acb719aacfcf))
---
# v0.1.9

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
---
# v0.1.8

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
---
# v0.1.7

## Bug Fixes
* fix(cli): access integration state securely (Salim Afiune Maya)([543562b](https://github.com/lacework/go-sdk/commit/543562b282f214ce877fabe12d1dde27017ddf30))
---
# v0.1.6

## Features
* feat(api): add EventsService to inspect events (Salim Afiune Maya)([533a271](https://github.com/lacework/go-sdk/commit/533a2713f5c179e50c90c63318991643f005a750))
* feat(api): add Details func to EventsService (Salim Afiune Maya)([56b95ca](https://github.com/lacework/go-sdk/commit/56b95ca2c02c2f8af24dd351e7fe6247b4da7eba))
* feat(cli): new event list command (Salim Afiune Maya)([d7c9f9e](https://github.com/lacework/go-sdk/commit/d7c9f9e2c41a1bd1411b92b6b8632aa2e32845dd))
* feat(cli): new event show command (Salim Afiune Maya)([8f75c78](https://github.com/lacework/go-sdk/commit/8f75c78d222f3fa00d24cf41e8d1e712f6600122))
* feat(cli): `--noninteractive` mode flag (Salim Afiune Maya)([10536af](https://github.com/lacework/go-sdk/commit/10536afe1d6ce76ba3391c145f3533b4d6725484))
## Bug Fixes
* fix(api): omitempty integration responses fields (Salim Afiune Maya)([44e2314](https://github.com/lacework/go-sdk/commit/44e2314f4ca02f1c0e6a134bedbc82161a81473c))
## Documentation Updates
* docs(cli): document environment variables (Salim Afiune Maya)([0012ec1](https://github.com/lacework/go-sdk/commit/0012ec14f574f6e4c1dc2a5d774e17ef038f1308))
## Other Changes
* chore(cli): update usage of commands (Salim Afiune Maya)([5dd3057](https://github.com/lacework/go-sdk/commit/5dd3057371fe87434f2da54f68cdcc3dc5fd754a))
* chore(cli): style updates to release scripts (Salim Afiune Maya)([f4355bf](https://github.com/lacework/go-sdk/commit/f4355bf481a3349fe9bac3700a1eaa8e80227238))
---
# v0.1.5

## Features
* feat(cli): implement JSON format for all commands (Salim Afiune Maya)([c7d4fee](https://github.com/lacework/go-sdk/commit/c7d4fee38f7b02a678e2fd30ecdb8f5bff82c11b))
* feat(cli): `vul scan run` command can poll for status (Salim Afiune Maya)([e2c8c8d](https://github.com/lacework/go-sdk/commit/e2c8c8dd1ee98d800902e6b671477e955e30a7fe))
## Bug Fixes
* fix(install.ps1) copy lacework.exe to ProgramData (Salim Afiune Maya)([53e685f](https://github.com/lacework/go-sdk/commit/53e685f583b70ff56b40ffdd8b2500d258891580))
---
# v0.1.4

## Features
* feat(install.ps1): support to install cli on windows (Salim Afiune Maya)([ae53d6f](https://github.com/lacework/go-sdk/commit/ae53d6ffe0f79b7d1f0cdbbd4feab560f727f29c))
## Bug Fixes
* fix(cli): use correct variable inside install.sh (Salim Afiune Maya)([eda17d5](https://github.com/lacework/go-sdk/commit/eda17d5ef4b2986bcbe99e4db0213ca7ebd18876))
* fix(cli): support colors for windows (Salim Afiune Maya)([ea48379](https://github.com/lacework/go-sdk/commit/ea4837955064c32575311722b9b137fd08161417))
## Documentation Updates
* doc(api): adds `_examples` or token-generation (Salim Afiune Maya)([c9dbc02](https://github.com/lacework/go-sdk/commit/c9dbc022831ed6b79d5216138c439662c7ad39e4))
* docs(cli): update install.sh to use CLI (Salim Afiune Maya)([d0fda04](https://github.com/lacework/go-sdk/commit/d0fda0460f9e247b53f92f4f72e32fa248de9b8c))
* docs: added lacework-cli profiles (Salim Afiune Maya)([d2292b0](https://github.com/lacework/go-sdk/commit/d2292b0425cde9334afc7b92e9fdc76ec63649c4))
## Other Changes
* chore(typo): misspelled word inside install.ps1 (Salim Afiune Maya)([ec93a96](https://github.com/lacework/go-sdk/commit/ec93a964407a84523708abaaee8569390216ec21))
* chore(timeout): increase api timeout to 60s (Salim Afiune Maya)([1f17bad](https://github.com/lacework/go-sdk/commit/1f17badb40690c04389d87453f556686fac17fb0))
---
# v0.1.3

## Features
* feat(vul): show number of fixable vulnerabilities (Salim Afiune Maya)([6403029](https://github.com/lacework/go-sdk/commit/6403029284085579a10cd4f74b41c3fd53ca765f))
* feat(cli): new vulnerability command (Salim Afiune Maya)([494d8d8](https://github.com/lacework/go-sdk/commit/494d8d8e317cad38bb6012567c4ff8ba9a2d3aa4))
* feat(api): add vulnerabilities service (Salim Afiune Maya)([d0b2c3b](https://github.com/lacework/go-sdk/commit/d0b2c3b5f2ad197a5ee9ba0d276c2d16d618891d))
* feat: introducing named profiles (Salim Afiune Maya)([6fb64fd](https://github.com/lacework/go-sdk/commit/6fb64fd2e0f953bf10103a44a821f840292991c4))
* feat: disallow extra arguments on sub-commands (#48) (Salim Afiune Maya)([f67ca9a](https://github.com/lacework/go-sdk/commit/f67ca9af38863e3d4a44a16d1130f919f7c2592e))
* feat: add configure command (#47) (Salim Afiune Maya)([f334fda](https://github.com/lacework/go-sdk/commit/f334fda595e000cc6ab830314c88f950eabc6761))
## Other Changes
* chore: adds a couple new go package dependencies (Salim Afiune Maya)([1842700](https://github.com/lacework/go-sdk/commit/1842700a8e7dfc91e91b11d5b21a62247a755b34))
---
# v0.1.2

## Features
* feat(lwloggder): go package for logging messages (Salim Afiune Maya)([cb5feee](https://github.com/lacework/go-sdk/commit/cb5feeeb6c1ddad54c6163c1e2b2c4dfdb6381fa))
## Refactor
* refactor(cli): rename cli binary to lacework (Salim Afiune Maya)([51ce22f](https://github.com/lacework/go-sdk/commit/51ce22f579f2984dfcfd54faaee844b621d7a617))
---
# v0.1.1

## Features
* feat(api): debug logs for all requests & responses (Salim Afiune Maya)([209f7ee](https://github.com/lacework/go-sdk/commit/209f7ee6240ccacabaac136b258409a48f673e8f))
* feat(api): add api client IDs for multi-client req (Salim Afiune Maya)([82c209f](https://github.com/lacework/go-sdk/commit/82c209fd04690a98109cc4a197431de29490c419))
* feat(api): implement a logging mechanism using zap (Salim Afiune Maya)([c078a70](https://github.com/lacework/go-sdk/commit/c078a70c002cb1646f9ba21e9d6d0b18aee85fdc))
## Bug Fixes
* fix(cli): error when account is empty (Salim Afiune Maya)([7dc59aa](https://github.com/lacework/go-sdk/commit/7dc59aa5958016dc2f3524495296cd27366fe660))
* fix(cli): load debug state correctly (Salim Afiune Maya)([8f7343c](https://github.com/lacework/go-sdk/commit/8f7343c93c11b757b64e2fed163684634a95ff36))
* fix(cli) Update environment variable prefix (Scott Ford)([484ca39](https://github.com/lacework/go-sdk/commit/484ca3951ecc934fd12ca88691b20a9781d245c7))
* docs(README) Update cli README to add documentation for ENV VARS (Scott Ford)([484ca39](https://github.com/lacework/go-sdk/commit/484ca3951ecc934fd12ca88691b20a9781d245c7))
## Documentation Updates
* doc(cli): fix single quote typo (Salim Afiune Maya)([3770b89](https://github.com/lacework/go-sdk/commit/3770b89d3d90aed18822fd2fa521cef3f23351bd))
* doc(logo): add logo to main README (Salim Afiune Maya)([620b992](https://github.com/lacework/go-sdk/commit/620b9921fa9d04147f4dfce3f4a24dfdec9a4238))
## Other Changes
* chore(cli): hide integration sub-commands (Salim Afiune Maya)([791ef7d](https://github.com/lacework/go-sdk/commit/791ef7d4df542c382607a12de08b262b95059311))
* chore(typo): fix RestfulAPI typo (Salim Afiune Maya)([39a7298](https://github.com/lacework/go-sdk/commit/39a72989366c85288b91d9ad033ed431a7796395))
* build: fix release checks (Salim Afiune Maya)([08bdb7d](https://github.com/lacework/go-sdk/commit/08bdb7d67fa22cde8b16a8ab392646b3288e5aa9))
* build(release): generate changelog and release notes (Salim Afiune Maya)([3aa0a91](https://github.com/lacework/go-sdk/commit/3aa0a91f31984fb0935b079c573b2f33ab5e7831))
---
# v0.1.0

## Features
* feat(cli): Installation scripts and documentation 🎉 (Salim Afiune Maya)([bb96b3b](https://github.com/lacework/go-sdk/commit/bb96b3bf26f105137afc50011b2c88c67e4ed0c7))
* feat(cli): the new lacework-cli MVP 🔥🔥 (Salim Afiune Maya)([34a73b6](https://github.com/lacework/go-sdk/commit/34a73b6d8df6e58225186831ae62a86a1724d747))
* feat(integrations): add AZURE_CFG and polish the rest (Salim Afiune Maya)([abd5bee](https://github.com/lacework/go-sdk/commit/abd5bee7d21141116f630faef44bdff707385e1d))
* feat(api): List integrations by type (Salim Afiune Maya)([f96a15b](https://github.com/lacework/go-sdk/commit/f96a15bc492d1acfc427b1851e3dfc12dc83a48b))
* feat: implement service model (Salim Afiune Maya)([d0cbf9f](https://github.com/lacework/go-sdk/commit/d0cbf9f20d7605a0f8a782909c09e78b3f1a6d8e))
* feat(api): new GetIntegrationSchema() (Salim Afiune Maya)([1aaec6c](https://github.com/lacework/go-sdk/commit/1aaec6cc8f3784d93f01d9819547dcab8cefddd4))
* feat(integrations): CRUD azure config integrations (Salim Afiune Maya)([0f83504](https://github.com/lacework/go-sdk/commit/0f83504636753c5e3d8cede843c1e18c73826c84))
* feat(integrations): CRUD aws config integrations (Salim Afiune Maya)([93475b0](https://github.com/lacework/go-sdk/commit/93475b0b229100dbb44ca4f779eb19120e93d298))
* feat(request): trigger token generation if missing (Salim Afiune Maya)([8cd82d6](https://github.com/lacework/go-sdk/commit/8cd82d677f65b7606111f91685e22e0618e6404d))
* feat(fakeAPI): New LaceworkServer to mock API req (Salim Afiune Maya)([c8211c1](https://github.com/lacework/go-sdk/commit/c8211c1e91a2ac116b65067ba890ab5dcae38e98))
* feat(client): Option to trigger a new token gen (Salim Afiune Maya)([96c8c6b](https://github.com/lacework/go-sdk/commit/96c8c6b78c8e73a17b8a71f3d2f25a78ce598b57))
## Refactor
* refactor: leverage integration structs for all gcp (Salim Afiune Maya)([922d117](https://github.com/lacework/go-sdk/commit/922d11755d4200cfbecdf56cb7e5f96775b9f136))
* refactor: leverage integration structs for all azure (Salim Afiune Maya)([1037d1b](https://github.com/lacework/go-sdk/commit/1037d1b8ee0028e4c688dbd9468499e7865eea6f))
* refactor: leverage integration structs for all aws (Salim Afiune Maya)([1146348](https://github.com/lacework/go-sdk/commit/11463481c419b1b9bb691ef1c15a43550b268095))
* refactor(integration): make space for New() funcs (Salim Afiune Maya)([1da9746](https://github.com/lacework/go-sdk/commit/1da9746332bb60dbb06c8345164b362fb0d06ce3))
* refactor(integration): move CRUD gcp config code (Salim Afiune Maya)([962191b](https://github.com/lacework/go-sdk/commit/962191b6c9383d6b52c9e6d78bde2507872541ba))
## Bug Fixes
* fix(install): configurable installation_dir (Salim Afiune Maya)([9d17b1f](https://github.com/lacework/go-sdk/commit/9d17b1f2d5e424acbee873d5b9934ff85e4f67f0))
* fix(release): tar linux binaries (Salim Afiune Maya)([9311b8f](https://github.com/lacework/go-sdk/commit/9311b8f2f5bd6e65381f7d0e8900da87925f8f78))
* fix(update): GCP CFG api path (mjunglw)([3508e78](https://github.com/lacework/go-sdk/commit/3508e78d483dd53b46e5e8d9c2b26ded257b1ef4))
* fix(gcp): update missing fields in structs (#14) (lwmobeent)([ce9745f](https://github.com/lacework/go-sdk/commit/ce9745f6202dac7e1c27b3f0eef582df29bc21a6))
* fix(enums): integrationType and gcpResourceLevel from array to map (mjunglw)([37c4d77](https://github.com/lacework/go-sdk/commit/37c4d771b27475eec0a05a55d00375e41d8ba53b))
* fix(client): expose Client struct for provider to use (mjunglw)([5e97951](https://github.com/lacework/go-sdk/commit/5e9795177ce2e76e029f66a60bbfb1497514ba08))
* fix(lint): various lint fixes (Salim Afiune Maya)([92efbab](https://github.com/lacework/go-sdk/commit/92efbab9b952e240045687717388f6badff072de))
## Documentation Updates
* docs(api): update README's and code comments (Salim Afiune Maya)([9d3e739](https://github.com/lacework/go-sdk/commit/9d3e7399275becf01ac23aaf702fd288c4518c2e))
* docs(README): Add usage and descriptions (Salim Afiune Maya)([9ed08dc](https://github.com/lacework/go-sdk/commit/9ed08dc604ef2482c8566f2d431a476920376239))
## Other Changes
* chore(deps): Add a few Go dependencies 🙌 (Salim Afiune Maya)([4ae8b8e](https://github.com/lacework/go-sdk/commit/4ae8b8e7d2671b40a98ae4b863192e6fe89f7e5b))
* chore(api): make response structs public (Salim Afiune Maya)([6b84e2c](https://github.com/lacework/go-sdk/commit/6b84e2c5786012d0d672df9e46c1e1166a38d086))
* ci(prepare): prepare the pipeline (Salim Afiune Maya)([94a8b0f](https://github.com/lacework/go-sdk/commit/94a8b0fddd0951734bd39137a77acb79cd3e9224))
* ci(tests): run tests in shippable ci (Salim Afiune Maya)([cb73c4b](https://github.com/lacework/go-sdk/commit/cb73c4bca8d9689ee505ccbe4848b9e8cfe72fa9))
* test: fix integration delete requests (Salim Afiune Maya)([4181580](https://github.com/lacework/go-sdk/commit/41815805e7039b7bb066432ced636c6d59e1d81b))
* test(integrations): generic Delete() func (Salim Afiune Maya)([c491f1a](https://github.com/lacework/go-sdk/commit/c491f1af2574330b7eda8c7ed285f05abb17a099))
* test(integrations): generic Get() func (Salim Afiune Maya)([0d0c8b0](https://github.com/lacework/go-sdk/commit/0d0c8b0e9d014b868a4e59004fd5407d63090ab0))
* test(unit): verify integrationType is well map (Salim Afiune Maya)([320640f](https://github.com/lacework/go-sdk/commit/320640f35068046ebbefccd9bfef95234e3dad4c))
