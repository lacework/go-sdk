# Release Notes
Another day, another release. These are the release notes for the version `v1.17.0`.

## Features
* feat(cli): New 'compliance scan' command (#1238) (Darren)([6bb952fe](https://github.com/lacework/go-sdk/commit/6bb952fe721733d97938356fcab412f96f657c1d))
* feat(api): Addition of S3 Bucket ARN for AWS EKS Audit log (David McTavish)([e3b80e67](https://github.com/lacework/go-sdk/commit/e3b80e6780d375b4eccb1611f2a5312942f3dbc4))
## Refactor
* refactor: do not filter csv output for vuln ctr show (#1220) (Darren)([dfecd161](https://github.com/lacework/go-sdk/commit/dfecd16139f63041c5878826388eca4f4c1117dc))
## Bug Fixes
* fix: vuln host scan-pkg-manifest fail_on_severity (#1242) (Darren)([e05db169](https://github.com/lacework/go-sdk/commit/e05db16978cc08e7b32d3998718ea46b109ed391))
* fix: increase default TLS Handshake Timeout to 63s (#1237) (Darren)([0d3aa066](https://github.com/lacework/go-sdk/commit/0d3aa0666a4b43e83e5772696e526dfc0643d559))
* fix: unsupported aws region (#1231) (jonathan stewart)([be742c0c](https://github.com/lacework/go-sdk/commit/be742c0c97b7f9a61e31f611125af1e7bcbf62fa))
* fix(api): Omit s3 bucket arn from Aws EKS Audit if not set (#1234) (djmctavish)([4b264e9b](https://github.com/lacework/go-sdk/commit/4b264e9b5d3d0e2c8dacda0af5286e7be56fae5c))
* fix(cli): allow multiple recommendation ID types (#1230) (Salim Afiune)([f0e8a8d9](https://github.com/lacework/go-sdk/commit/f0e8a8d9b56bf291c9d9d2282389085765d1cc54))
## Other Changes
* style(fmt): fix format issues (Salim Afiune Maya)([0986f72e](https://github.com/lacework/go-sdk/commit/0986f72ef5b14bf689e33dbdfbec320246e5ada0))
* chore(deps): bump github.com/mattn/go-isatty from 0.0.14 to 0.0.18 (#1205) (dependabot[bot])([e9a4bd49](https://github.com/lacework/go-sdk/commit/e9a4bd49af8c300e050b0cc394dce2f40bb49b7f))
* chore(deps): bump golang.org/x/net from 0.2.0 to 0.7.0 (#1154) (dependabot[bot])([f88e6f13](https://github.com/lacework/go-sdk/commit/f88e6f13053fd6d0c34078185b333e5fcdf92607))
* chore(deps): bump github.com/abiosoft/colima from 0.4.6 to 0.5.4 (#1172) (dependabot[bot])([9f69228a](https://github.com/lacework/go-sdk/commit/9f69228a9b0994d1f0d99fdec2dd8444c95e3e6c))
* chore(cli): use an example URL for CTF docs (#1228) (Nick Schmeller)([e18af0db](https://github.com/lacework/go-sdk/commit/e18af0dbd5c1421b4a44c959f7a378daf8cd2de6))
* ci: fix intermittency issues with virtual terminals (#1235) (Salim Afiune)([e281108d](https://github.com/lacework/go-sdk/commit/e281108dfe4b0e3d137dcad6a9d27a4ef07879eb))
* ci(fmt): avoid formatting 'databox' package (Salim Afiune Maya)([93cedfef](https://github.com/lacework/go-sdk/commit/93cedfef5d2ce88c4d0c8d0225a7ae5d5b199c81))
* ci: version bump to v1.16.1-dev (Lacework)([916a63f5](https://github.com/lacework/go-sdk/commit/916a63f5afbbd9b5fbe54d2ff73f996b7dff48c3))
* test: fix compliance scan test (#1243) (Darren)([1a497bab](https://github.com/lacework/go-sdk/commit/1a497bab6ee12d506b4503a155ecf248f6ab413e))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
