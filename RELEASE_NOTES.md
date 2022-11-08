# Release Notes
Another day, another release. These are the release notes for the version `v0.45.0`.

## Features
* feat: add JIRA bidirectional and grouping support (#940) (hazedav)([441fc00d](https://github.com/lacework/go-sdk/commit/441fc00d317f2593c70f11b52621ed87043da536))
* feat(generate): Add Terraform generation for AWS EKS Audit Log (#960) (Ross)([67dfb49b](https://github.com/lacework/go-sdk/commit/67dfb49bdfd1b16377c1feb7488eeb0ffd0e14ee))
## Bug Fixes
* fix(cdk): Avoid double component error messages (#1000) (Salim Afiune)([363c1d40](https://github.com/lacework/go-sdk/commit/363c1d40ebf7e0d2e675745e037c273dc646e317))
* fix(cdk): Pass all global flags via env variables (#993) (Salim Afiune)([1dc9fadd](https://github.com/lacework/go-sdk/commit/1dc9fadd9ef62376ccf5941f1ad632269845f104))
* fix(all) LimitNumScans was wrong type (#978) (Christien Rioux)([eab66dab](https://github.com/lacework/go-sdk/commit/eab66dab16e9930295a942a521b000148761470a))
## Other Changes
* style(cdk): detect lwcomponent error inside Execute() (#1001) (Salim Afiune)([5f88e520](https://github.com/lacework/go-sdk/commit/5f88e52067f2f4db6bbc3c7fb95da83e3fc253de))
* chore(deps): bump github.com/abiosoft/colima from 0.4.5 to 0.4.6 (#937) (dependabot[bot])([d8f67d1c](https://github.com/lacework/go-sdk/commit/d8f67d1cb66c9ad65cd567d9acf4e88388f35b9d))
* ci: version bump to v0.44.2-dev (Lacework)([10f4e321](https://github.com/lacework/go-sdk/commit/10f4e321abd57dc00c3e2b6b756ed98a017696b8))
* test(cdk): fix tests due to IaC component release (#998) (Salim Afiune)([4af2e7f1](https://github.com/lacework/go-sdk/commit/4af2e7f17e4c854c58b77f52765c60d97769664a))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
