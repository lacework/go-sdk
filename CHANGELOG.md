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
