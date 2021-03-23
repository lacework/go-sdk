# v0.3.0

## Features
* feat(api): implement account information endpoint (#349) (Salim Afiune)([1392ecb](https://github.com/lacework/go-sdk/commit/1392ecb98c2908ab143ee630f0df4e0453a07767))
## Other Changes
* chore: ran make prepare (Salim Afiune Maya)([662b220](https://github.com/lacework/go-sdk/commit/662b220a67f900fe003a0d275462589c1df0eb1b))
* ci: fix prepare-test-resources job (#348) (Darren)([a071c04](https://github.com/lacework/go-sdk/commit/a071c040a4c14eaa4c5b154361d94b5ffcb7d458))
* ci: open release pull request automatically (Salim Afiune Maya)([f227182](https://github.com/lacework/go-sdk/commit/f2271824616649f827fcffc5a15c99fdc809dd04))
* ci: automatic minor version bump (Salim Afiune Maya)([73e2cb9](https://github.com/lacework/go-sdk/commit/73e2cb9cb11f22909f7d698e10aa28a6a961061f))
* test: change target for container vuln scanning (#343) (Darren)([c348b01](https://github.com/lacework/go-sdk/commit/c348b0156f6846fb999202a29380639f0899c490))
---
# v0.2.23

## Features
* feat(cli): add Azure list-tenants sub-command (#341) (Darren)([960a8b7](https://github.com/lacework/go-sdk/commit/960a8b7f53a52febb66e89dce20f28211569a3fb))
* feat(cli): filter vulnerability assessments by severity (#338) (Darren)([07af9b1](https://github.com/lacework/go-sdk/commit/07af9b1091c8387d5d9e094a0fb5e9024749e4ef))
---
# v0.2.22

## Bug Fixes
* fix: implement both ECR auth methods (Salim Afiune Maya)([7af312c](https://github.com/lacework/go-sdk/commit/7af312c0308d3dad52f4e0264ee354800ed4d616))
* fix(api): type-o in host vulnerability status (#336) (Scott Ford)([85b271b](https://github.com/lacework/go-sdk/commit/85b271bdb35c33eec12df5591c530800688b6d90))
## Other Changes
* ci: run nightly integration tests on ARM-linux (#317) (Salim Afiune)([525b51d](https://github.com/lacework/go-sdk/commit/525b51d6866550a26ae7faefb2d8a8de99dd9323))
* test(cli): use a cli tag with vuln (Salim Afiune Maya)([eae52b8](https://github.com/lacework/go-sdk/commit/eae52b8d8350fd1e6166c036fa8b86e3462517de))
---
# v0.2.21

## Features
* feat(cli): load account from new UI API JSON file (#331) (Salim Afiune)([e841210](https://github.com/lacework/go-sdk/commit/e841210501eaa25620e102fad1e1ecf7a8c3bd3b))
* feat(cli): filtering flags for compliance report recommendations (#330) (Darren)([d04f09f](https://github.com/lacework/go-sdk/commit/d04f09f7cebaf278aaf22085cdb0f0a9f54b0b85))
* feat: Add support for ServiceNow Alert custom JSON template (#327) (Darren)([129bc28](https://github.com/lacework/go-sdk/commit/129bc2861f7ba7a807cb35c9d001b7ae575d97ab))
## Bug Fixes
* fix: Generate a new token upon a request with an expired token (#332) (Darren)([2bbc8b3](https://github.com/lacework/go-sdk/commit/2bbc8b38d7f5332174eb2f7266e300763054bdce))
## Documentation Updates
* docs: update go version batch in README (#329) (Salim Afiune)([161753e](https://github.com/lacework/go-sdk/commit/161753e35eac1990f69cc58d24c2b8afc5a5364d))
* docs: Add Homebrew installation to README (#328) (Darren)([e4ad780](https://github.com/lacework/go-sdk/commit/e4ad7803de81ed330b1a6eb9e5f61c0e3353c5e2))
## Other Changes
* ci: delete 'master' branch (#334) (Salim Afiune)([792e458](https://github.com/lacework/go-sdk/commit/792e458aed01012d2c665998a3e1e09f4a7facbb))
---
# v0.2.20

## Features
* feat(cli): New IBM QRadar alert channel (#325) (Darren)([0e9f6f5](https://github.com/lacework/go-sdk/commit/0e9f6f57e16e96ccfb84b9bd618e595378cf0fb3))
* feat(cli): New Relic Insights alert channel (#323) (Darren)([d7242b8](https://github.com/lacework/go-sdk/commit/d7242b84a525a4b9484324f2dfb49b172528302b))
## Documentation Updates
* doc(cli): update scan-pkg-manifest help to 10k pkgs (#324) (Salim Afiune)([0348800](https://github.com/lacework/go-sdk/commit/03488001d7c0aab308adf00ec6d8ac4a73aec78c))
---
# v0.2.19

## Features
* feat(cli): New VictorOps alert channel (#318) (Darren)([dfcd34a](https://github.com/lacework/go-sdk/commit/dfcd34adae5b84d697bd5f76831f749765f885fa))
* feat(cli): New CiscoWebex alert channel (#316) (Darren)([8e0071a](https://github.com/lacework/go-sdk/commit/8e0071a7d8a879d59c4545c8528e01287b90041e))
* feat: New Microsoft Teams alert channel (#315) (Darren)([e414226](https://github.com/lacework/go-sdk/commit/e41422624b5274a5d2b571fca4399ee2c9df3af3))
## Refactor
* refactor(cli): exponential retries polling scans (Salim Afiune Maya)([2bb881d](https://github.com/lacework/go-sdk/commit/2bb881d90fa6497003a89045acb41978398d5502))
## Bug Fixes
* fix(cli): match API client timeout with NGINX (#321) (Salim Afiune)([10b7a28](https://github.com/lacework/go-sdk/commit/10b7a28677af003e176de26b0dc1558c4837d1c4))
## Other Changes
* ci: increase integration test timeout to 30m (Salim Afiune Maya)([3081e3c](https://github.com/lacework/go-sdk/commit/3081e3cc6264e85b7fd80019a36d89dda1d9e5c9))
* test: change ctr vuln scan tag (Salim Afiune Maya)([c15bd1b](https://github.com/lacework/go-sdk/commit/c15bd1b2e0454ff61f6a7549f62c6bc791190079))
* test(cli): disable failing tests (RAIN-15300) (#320) (Salim Afiune)([e2afb31](https://github.com/lacework/go-sdk/commit/e2afb31277e8a2781c5f8e55ebf352992cd6d2b0))
---
# v0.2.18

## Features
* feat(cli): new Datadog alert channel  (#313) (Darren)([8298022](https://github.com/lacework/go-sdk/commit/8298022367d99ccc0237e456c8ccd6b45885e1d7))
## Bug Fixes
* fix(cli): avoid daily update check during install (Salim Afiune Maya)([2bc94c0](https://github.com/lacework/go-sdk/commit/2bc94c002392cb900ce6ca47574eafa3f4e55959))
## Other Changes
* test: fix intermittent events test (#312) (Salim Afiune)([d69983b](https://github.com/lacework/go-sdk/commit/d69983b63c9c771f533d72325bc58cd63e5d3a66))
* test: fix intermittent events test (Salim Afiune Maya)([15c371c](https://github.com/lacework/go-sdk/commit/15c371cc59ae9b2cbbe25588ff20574c1272efbe))
---
# v0.2.17

## Bug Fixes
* fix: Json mapping for Snow Username (#307) (Darren)([1ef8d99](https://github.com/lacework/go-sdk/commit/1ef8d9912dd66522712bfa4a15d4ab728d409e65))
---
# v0.2.16

## Features
* feat(cli): New Service Now alert channel (#303) (Darren)([512f2d9](https://github.com/lacework/go-sdk/commit/512f2d9c38d8124904dbdd661d9ab3b8441fc86d))
## Refactor
* refactor: Change input method for private_key field (#305) (Darren)([e56cdc6](https://github.com/lacework/go-sdk/commit/e56cdc68d74074f5ee904712aa56779dc9a0e1ed))
## Bug Fixes
* fix: Use select for issue grouping (#304) (Darren)([799d9c3](https://github.com/lacework/go-sdk/commit/799d9c34723d4af03f9ed811028880181f1757f5))
* fix: Add issue_grouping field to gcp pub sub (#301) (Darren)([1a66d2c](https://github.com/lacework/go-sdk/commit/1a66d2cc20a33674348705fbd3e552bf9222c787))
* fix(cli): install.sh should try curl and wget (Salim Afiune Maya)([f6b0bd7](https://github.com/lacework/go-sdk/commit/f6b0bd725992961e8c8b02ba3090164f31056388))
* fix(cli): install.sh should respect target override (Salim Afiune Maya)([4164f58](https://github.com/lacework/go-sdk/commit/4164f5872dd381471ad8608abd9c2fa821dac524))
## Other Changes
* chore(cli): install.sh print exitcodes for debugging (Salim Afiune Maya)([5e66c11](https://github.com/lacework/go-sdk/commit/5e66c11264f0b9ba012be20370b0b65cc114fe4c))
* ci: trigger homebrew update script (#299) (Darren)([9247cb1](https://github.com/lacework/go-sdk/commit/9247cb12b62582596939be55861f4d6c08bf8690))
---
# v0.2.15

## Features
* feat: add telemetry to detect Homebrew installations (#297) (Darren)([fa81abc](https://github.com/lacework/go-sdk/commit/fa81abc1c044cb362cd29608c65fc820d0c8a706))
* feat(cli): New Gcp PubSub alert channel (#294) (Darren)([08a3e61](https://github.com/lacework/go-sdk/commit/08a3e61469a1bb35675d63070474bad3e7988ad4))
---
# v0.2.14

## Features
* feat(cli): support Homebrew upgrade command (#291) (Darren)([bedfa5d](https://github.com/lacework/go-sdk/commit/bedfa5d2895fa7aabe1de6ef2902eb260133c3fd))
* feat(cli): Add Splunk alert channel (#289) (Darren)([04679a5](https://github.com/lacework/go-sdk/commit/04679a51bf35469f4986e693ed19a75cc36dbbb2))
* feat(cli): add account check to catch http(s):// (#288) (Salim Afiune)([3d770a1](https://github.com/lacework/go-sdk/commit/3d770a1171814e33de7f6420d1637ad8c03f30c8))
## Bug Fixes
* fix(cli): skip daily version check for version cmd (#290) (Salim Afiune)([5c9f4ca](https://github.com/lacework/go-sdk/commit/5c9f4cab61620636e0962e0fbe4edce97c41e8dc))
---
# v0.2.13

## Features
* feat(cli): support manifest bigger than 1k packages (Salim Afiune Maya)([eebddb9](https://github.com/lacework/go-sdk/commit/eebddb9325ede76ffa1853d00508da54cb5b9678))
* feat(cli): gen-pkg-manifest detect running kernel (Salim Afiune Maya)([9151be1](https://github.com/lacework/go-sdk/commit/9151be15a05b48f3d7456571cd75411f2ba7ddb9))
## Refactor
* refactor: simplify removeEpochFromPkgVersion func (Salim Afiune)([04aba5b](https://github.com/lacework/go-sdk/commit/04aba5bda340283f86d93496f01e0089a500468d))
## Bug Fixes
* fix(cli): ensure api client has valid auth token (Salim Afiune Maya)([056eda5](https://github.com/lacework/go-sdk/commit/056eda5cb7bde11e2334b6f38bd338afe111ade9))
## Other Changes
* ci: generate code coverage in HTML format (Salim Afiune Maya)([a58b58a](https://github.com/lacework/go-sdk/commit/a58b58a6477ec8d12c06bff3672093aef826c1f1))
* ci: add 'metric' as a valid commit message (Salim Afiune Maya)([dd7b601](https://github.com/lacework/go-sdk/commit/dd7b6010969d1f99055b7dbc9442498fa9f002cf))
* ci: fix slack notifications team alias ⭐ (Salim Afiune Maya)([ca51f92](https://github.com/lacework/go-sdk/commit/ca51f92693a48f113dd7661d9ef03eef7c26a17a))
* metric(cli): detect feature split_pkg_manifest (Salim Afiune Maya)([fdb9f4a](https://github.com/lacework/go-sdk/commit/fdb9f4a1c1eae2b9a44ea846fae413a93f073ca9))
* metric(cli): detect feature gen_pkg_manifest (Salim Afiune Maya)([78905bb](https://github.com/lacework/go-sdk/commit/78905bb73f398bf26a6e297e3929e5993e4965dc))
---
# v0.2.12

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
---
# v0.2.11

## Features
* feat(cli): daily version check (#269) (Salim Afiune)([5c15eef](https://github.com/lacework/go-sdk/commit/5c15eef84f428ec0534954babb28a3db92d5a7c5))
* feat(api): add Webhook integration (#267) (Darren)([f32572e](https://github.com/lacework/go-sdk/commit/f32572ecdadd5c179227cd228bf1fdd7cf618763))
## Refactor
* refactor(cli): abstract rendering tables (human-readable) (#263) (Salim Afiune)([8a10b4c](https://github.com/lacework/go-sdk/commit/8a10b4cf10de03d9b4c0409e495fdd7118974b92))
## Bug Fixes
* fix(cli): render account mapping file correctly (#266) (Salim Afiune)([4c327d7](https://github.com/lacework/go-sdk/commit/4c327d7e6081d0f7726a1bc007b1b736a106933f))
* fix(api): new request body for lql service (#260) (Salim Afiune)([4e2b439](https://github.com/lacework/go-sdk/commit/4e2b439ff394d632cd6ebf214da376050da46812))
* fix(api): avoid updating AgentTokenResponse.Props (#259) (Salim Afiune)([c3fe8bc](https://github.com/lacework/go-sdk/commit/c3fe8bcc41efd995f756f97a5ffca8bb961e89e4))
## Documentation Updates
* docs: update READMEs and _examples/ (#268) (Salim Afiune)([3791da0](https://github.com/lacework/go-sdk/commit/3791da01005335c34852446c57eb99e51a6d3ce1))
## Other Changes
* build: upgrade Go version to 1.15 (#265) (Salim Afiune)([06d41f5](https://github.com/lacework/go-sdk/commit/06d41f56add71f8369ffae68ea7ba5d738eb4d5b))
* ci: update hostname from our test machine (#262) (Salim Afiune)([beb289e](https://github.com/lacework/go-sdk/commit/beb289e732c177e2f3d062d61dee5dd9f1593ce9))
---
# v0.2.10

## Features
* feat(cli): new agent access token command (#256) (Salim Afiune)([7f8ba11](https://github.com/lacework/go-sdk/commit/7f8ba113b38ecd768f61e54ba712badf6596a587))
* feat(compliance): new aws list-accounts command (Salim Afiune Maya)([705f2eb](https://github.com/lacework/go-sdk/commit/705f2ebf9f1b9b5af2eb745c86498fe31c01e174))
## Refactor
* refactor: account mapping file for consolidated CT (#252) (Salim Afiune)([402a363](https://github.com/lacework/go-sdk/commit/402a3634765ef8c6f1f65d1be13da2ad34cf2960))
## Bug Fixes
* fix(install.sh): avoid logging with 'info' cmd (#254) (Salim Afiune)([df5f8cf](https://github.com/lacework/go-sdk/commit/df5f8cfbc7228ff9bff25e6e22a2ab68acd47fa4))
* fix: false positive results in pkg manifest scan (#255) (Salim Afiune)([a6d6cda](https://github.com/lacework/go-sdk/commit/a6d6cda9f36b38f8b653bd01ef258bd431611908))
* fix(databox): remove hardcoded LW account (Salim Afiune Maya)([c806157](https://github.com/lacework/go-sdk/commit/c80615749827c12dbfef5e1c76bf5857cd3dae7a))
---
# v0.2.9

## Features
* feat(api): enable account mapping file for CT int (#250) (Salim Afiune)([cb99f61](https://github.com/lacework/go-sdk/commit/cb99f61f5da717911d72ff2379a98ce4b7f6dd61))
## Refactor
* refactor(api): better error check handler (#247) (Salim Afiune)([b363347](https://github.com/lacework/go-sdk/commit/b363347409f242e6ad1a46e1104884a559877ada))
## Other Changes
* ci: set container tag to debian-10 that has vulns (#248) (Salim Afiune)([323b91e](https://github.com/lacework/go-sdk/commit/323b91e38c6616a064501c602d58f2cb0a2572f2))
* ci: dogfooding Lacework Orb html parameter (Salim Afiune Maya)([464d34d](https://github.com/lacework/go-sdk/commit/464d34db0d53fc7e0e9f4b6cbd7663eae963ad46))
* ci: remove slack alert for win systems (Salim Afiune Maya)([b6b5b45](https://github.com/lacework/go-sdk/commit/b6b5b458de4e69fd0d176522a47a168382a03397))
---
# v0.2.8

## Bug Fixes
* fix(cli): generate html for scan commands (Salim Afiune Maya)([6846ffd](https://github.com/lacework/go-sdk/commit/6846ffd4a6f5d76ca700652e0a5cd7adda6c9cdc))
## Other Changes
* ci: improve release notes and changelog generation (Salim Afiune Maya)([af22a7a](https://github.com/lacework/go-sdk/commit/af22a7a5518adc4b5f1c6fc21c23934230dd0705))
* ci: avoid release.sh to update version multiple times (Salim Afiune Maya)([d72149b](https://github.com/lacework/go-sdk/commit/d72149b50b197a9ddc7528e18bc88c1836f4847f))
* test(cli): HTML for container vulnerability (Salim Afiune Maya)([fee8505](https://github.com/lacework/go-sdk/commit/fee85056a22b460b925f60083cbb2610c8552c55))
---
# v0.2.7

## Features
* feat(cli): enable html copy to clipboard icons (Salim Afiune Maya)([ec2d1fa](https://github.com/lacework/go-sdk/commit/ec2d1fa5b796e2b51e49a850deff01d2f64ded18))
* feat(cli): HTML format for vulnerability assessments (Salim Afiune Maya)([00c2f43](https://github.com/lacework/go-sdk/commit/00c2f43613e554afd8ed283cbc12eb0b8eed0179))
* feat(cli): add ARM support (#236) (Salim Afiune)([821b8e6](https://github.com/lacework/go-sdk/commit/821b8e699e61eefda7d287a71b08ef26382a4ad7))
## Bug Fixes
* fix(cli): remove html column sort icons (Salim Afiune Maya)([dc4c0f6](https://github.com/lacework/go-sdk/commit/dc4c0f64055bfaada503b8f4f21ceda707bc5e55))
## Other Changes
* ci(fix) Update CI test node (#233) (Scott Ford)([ddbf86e](https://github.com/lacework/go-sdk/commit/ddbf86e8fbf9053af43bab9d57c04645383e529e))
---
# v0.2.6

## Features
* feat(api): trigger initial report automatically (#230) (Salim Afiune)([1e24a22](https://github.com/lacework/go-sdk/commit/1e24a229d2f2c54b81809ccb156a9d4283962c32))
## Documentation Updates
* docs(cli): disable timestamp for automatic docs (#229) (Salim Afiune)([f4d7841](https://github.com/lacework/go-sdk/commit/f4d78417c307c38507995892999cb85e7be74cf2))
---
# v0.2.5

## Bug Fixes
* fix(cli): add epoch to package manifest (Salim Afiune Maya)([17da487](https://github.com/lacework/go-sdk/commit/17da48755062265245d98ba6f4a330ae65fcdb6b))
## Other Changes
* chore(ci): make GH org a readonly parameter (Salim Afiune Maya)([b4f5f6d](https://github.com/lacework/go-sdk/commit/b4f5f6d5ba5a644a6198445bd820d68bf243907d))
* chore(cli): update pkg-manifest message for 0 vuln (Salim Afiune Maya)([5029dc8](https://github.com/lacework/go-sdk/commit/5029dc82aa51f260e84cd476acd6c64cab7f063a))
---
# v0.2.4

## Features
* feat(cli): programatic access to profile data (#225) (Salim Afiune)([ab7ce7c](https://github.com/lacework/go-sdk/commit/ab7ce7cfe8e94053ca6bf8d32d929c5e748496e4))
* feat(cli): allow custom installation directory -d 📁 (#223) (Salim Afiune)([ee9e686](https://github.com/lacework/go-sdk/commit/ee9e686c46029b32e711f9534ecd7755926ec22b))
## Documentation Updates
* docs: automatically generate cli docs (#224) (Salim Afiune)([5b91e1e](https://github.com/lacework/go-sdk/commit/5b91e1e788128dd3cddf457bce565749c73eddae))
## Other Changes
* chore: add badges to README.md (#222) (Salim Afiune)([db7235d](https://github.com/lacework/go-sdk/commit/db7235d20e7af012cb8e8f3041a02728d4f28719))
---
# v0.2.3

## Features
* feat(cli): add scan-pkg-manifest summary 📈 (#220) (Salim Afiune)([9b009c3](https://github.com/lacework/go-sdk/commit/9b009c3e98a69d294d424c2b912b1aadb675ee98))
* feat(ux): generate package-manifest command (#217) (Salim Afiune)([0c842ab](https://github.com/lacework/go-sdk/commit/0c842ab15c30b3f754a379ecd2aea014c367bae7))
## Refactor
* refactor: remove 'apk' as supported pkg manager (Salim Afiune Maya)([4165783](https://github.com/lacework/go-sdk/commit/41657839f06ea9b8eae85119451c77e632ec99bb))
## Other Changes
* chore(ci): update lacework circleci orb (Salim Afiune Maya)([3952c66](https://github.com/lacework/go-sdk/commit/3952c66f47dbb0024b3fef35f3f39087fa76844e))
---
# v0.2.2

## Features
* feat(lql): --file flag to load LQL query from disk (Salim Afiune Maya)([4804319](https://github.com/lacework/go-sdk/commit/4804319a0c26211119c10eb3dc4d889b3da7e227))
* feat(cli): --file to pass a package manifest file (Salim Afiune Maya)([75680d8](https://github.com/lacework/go-sdk/commit/75680d8d9469d8679b17c46979439340d8869da9))
* feat: human-readable output for scan-pkg-manifest (Salim Afiune Maya)([783f550](https://github.com/lacework/go-sdk/commit/783f55015c1e6a1071927e19810266376ecbe082))
* feat(lql): improve running queries (Salim Afiune Maya)([61c5ee5](https://github.com/lacework/go-sdk/commit/61c5ee51aac65626aff4f81ebceb96633865d2f7))
## Bug Fixes
* fix(ci): remove slack notification for windows (#214) (Salim Afiune)([a2c5124](https://github.com/lacework/go-sdk/commit/a2c51242c08c1683cfb9c80c832be2559058f957))
## Other Changes
* ci(slack): notify pipeline failures (#213) (Salim Afiune)([85ad396](https://github.com/lacework/go-sdk/commit/85ad396f6cb049ab246ff36fa2f29d46fab6459d))
---
# v0.2.1

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
---
# v0.2.0

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
---
# v0.1.24

## Features
* feat(cli): better ux in account validation (#187) (Salim Afiune)([cdd045a](https://github.com/lacework/go-sdk/commit/cdd045a830dcdc788daf77d9ea558ba4d296e003))
* feat(cli): new access-tokens command (#184) (Salim Afiune)([ee338c4](https://github.com/lacework/go-sdk/commit/ee338c4afb057bf4ea578a8d0ddb48b2d39b34d3))
* feat(cli): Create Jira Alert Channels 🚨 (Salim Afiune Maya)([6ca8cef](https://github.com/lacework/go-sdk/commit/6ca8ceffce1c17f3f84634da5514c059da952ca1))
* feat(api): add Jira alert channel integrations (Salim Afiune Maya)([0cdb2a4](https://github.com/lacework/go-sdk/commit/0cdb2a46d820f249c0fe918320303b1061e0f5ed))
## Refactor
* refactor: remove legacy field min_alert_severity (#186) (Salim Afiune)([54ca38c](https://github.com/lacework/go-sdk/commit/54ca38c8c509d800e2bddca5435529f5d0b60643))
## Bug Fixes
* fix(cli): display integration update by/update time (Salim Afiune Maya)([7060078](https://github.com/lacework/go-sdk/commit/7060078d3f8a09a82a4efaf98c4cb15f4856f753))
---
# v0.1.23

## Refactor
* refactor(cli): replace '--pdf-file' for '--pdf' (#180) (Salim Afiune)([80bbce6](https://github.com/lacework/go-sdk/commit/80bbce636cac49fe315118add45252bd8ee4bf6a))
## Bug Fixes
* fix(cli): missing integration details (#181) (Salim Afiune)([40355d3](https://github.com/lacework/go-sdk/commit/40355d3877c2674268c38bb5cc81a698dd115166))
* fix(cli): error showing non-existing integration (#178) (Salim Afiune)([252072f](https://github.com/lacework/go-sdk/commit/252072faa60aaac06fb7bbf2dd7ca82fa71d2b09))
## Other Changes
* ci: build statically linked binaries (Salim Afiune Maya)([43f6f80](https://github.com/lacework/go-sdk/commit/43f6f804ffac3f8e326dc31f4196808f39bc035d))
* ci(integration): add windows support (Salim Afiune Maya)([46632e7](https://github.com/lacework/go-sdk/commit/46632e72e0ab9ee45d690605e4c52efb1a8cf391))
---
# v0.1.22

## Features
* feat(cli): Create PagerDuty Alert Channels 🚨 (#174) (Salim Afiune)([5cc424e](https://github.com/lacework/go-sdk/commit/5cc424e21598482f817288037c8f8e54397c13bd))
* feat(api): add PagerDuty alert channel integrations (#173) (Salim Afiune)([f46316c](https://github.com/lacework/go-sdk/commit/f46316c7f4150ccf99646640a12d801cb407134b))
* feat(cli): Create AWS CloudWatch Alert Channels 🚨 (Salim Afiune Maya)([201b59b](https://github.com/lacework/go-sdk/commit/201b59be0a97d661916ff401da0be903fee06f2f))
* feat(api): add AWS CloudWatch Alert Channels Int (Salim Afiune Maya)([d9a11ec](https://github.com/lacework/go-sdk/commit/d9a11ec5c242b09e19338c6b8a5a39ddf6ad368d))
* feat(api): enum AlertLevel for alert severity levels (Salim Afiune Maya)([d3bf436](https://github.com/lacework/go-sdk/commit/d3bf436933a794b6bbcc733da724159a9dc79a95))
* feat(api): get/update container registry integrations (#168) (Salim Afiune)([a072c46](https://github.com/lacework/go-sdk/commit/a072c46aff03e619fbef03488ba5b65730264b91))
## Refactor
* refactor(api): AlertChannel prefix in funcs/structs (Salim Afiune Maya)([b0429ef](https://github.com/lacework/go-sdk/commit/b0429efd0efa56ec9ccbe338a37a6e6ae2dc3bc5))
* refactor(api): use AlertLevel enum for Slack Alerts (Salim Afiune Maya)([4b5acf9](https://github.com/lacework/go-sdk/commit/4b5acf989fda4c052c3dc6b0206db866aa57f243))
## Bug Fixes
* fix(cli): missing fields for Slack integrations (#170) (Salim Afiune)([a8ce9a9](https://github.com/lacework/go-sdk/commit/a8ce9a90f52dd81281fca78b077435229bdbafaf))
## Other Changes
* chore(api): adds alert channel \_examples/ (Salim Afiune Maya)([f967206](https://github.com/lacework/go-sdk/commit/f967206db3dd209f94e694b5f4db98dd8b11f113))
---
# v0.1.21

## Features
* feat(cli): Create Slack Channel Alerts 🚨 (#165) (Salim Afiune)([0d1f8c7](https://github.com/lacework/go-sdk/commit/0d1f8c74656c4e2043323b38cadde4e0456d6cfd))
* feat(api): add Slack Channel integrations (#164) (Salim Afiune)([fb81416](https://github.com/lacework/go-sdk/commit/fb81416b541882ef697d9be6dc0685772183e336))
* feat(api): new Vulnerabilities.ListEvaluations() func (#160) (Salim Afiune)([0060799](https://github.com/lacework/go-sdk/commit/0060799f47742091f9dc16eb987ba8cc5b5cee25))
* feat(cli): configure in non-interactive mode (#158) (Salim Afiune)([781f65b](https://github.com/lacework/go-sdk/commit/781f65b7f3449cbb2bb04831aa5443e7981a30e4))
* feat(cli): add --packages flag to vulnerability cmd (#149) (Salim Afiune)([3c34eaf](https://github.com/lacework/go-sdk/commit/3c34eaf8de21a1e5f23707034ca69c02cabf5e25))
## Other Changes
* chore(cli): remove deprecated old config loading (#159) (Salim Afiune)([1661939](https://github.com/lacework/go-sdk/commit/1661939b94c42c080f051039c70f4c82a56f2ad3))
---
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
