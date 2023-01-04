# Release Notes
Another day, another release. These are the release notes for the version `v1.4.0`.

## Features
* feat(GROW-1058): Improve management of component versions (#1088) (Edoardo Pirovano)([952f2c30](https://github.com/lacework/go-sdk/commit/952f2c30cac5cacbf02494eab031d81f6b95dee5))
* feat(cdk): pass LW_CLI_VERSION env variable (#1085) (Salim Afiune)([f0382ae8](https://github.com/lacework/go-sdk/commit/f0382ae850c95d9956d902138d80465d5e7fad05))
* feat(lwlogger): add merge logger function (#1087) (Kolbeinn)([1e569be8](https://github.com/lacework/go-sdk/commit/1e569be8c5d53335605e2540b02f0b94ec21925e))
* feat(cli): pass LW_API_SERVER_URL env variable to components (#1077) (Kolbeinn)([c85d1541](https://github.com/lacework/go-sdk/commit/c85d154189d8ce830bf1b28bbc505fd0b48f692d))
## Bug Fixes
* fix(lwlogger): do not override explicit level argument (#1084) (Kolbeinn)([452dff83](https://github.com/lacework/go-sdk/commit/452dff8325497eae61a2ba76f59cd52eaf535678))
## Other Changes
* chore(deps): Run `make go-vendor` (#1094) (Edoardo Pirovano)([8beaf31f](https://github.com/lacework/go-sdk/commit/8beaf31fa3ba37de5d5b2dffd5bd167157580131))
* chore(deps): bump github.com/aws/aws-sdk-go-v2/service/ec2 (#1083) (dependabot[bot])([c5a04eab](https://github.com/lacework/go-sdk/commit/c5a04eab6aafee96ab7c2bda4c66fdf986ed75e4))
* chore(deps): bump github.com/honeycombio/libhoney-go (#1081) (dependabot[bot])([b9217256](https://github.com/lacework/go-sdk/commit/b921725609140ced50d5e0bfa9f91395129de34b))
* chore(deps): bump github.com/aws/aws-sdk-go-v2/config (#1075) (dependabot[bot])([58711b27](https://github.com/lacework/go-sdk/commit/58711b27dd08cdd4572b79b14c2792615ca438fd))
* chore(deps): bump github.com/hashicorp/hcl/v2 from 2.14.1 to 2.15.0 (#1044) (dependabot[bot])([44586503](https://github.com/lacework/go-sdk/commit/44586503459658e6be01f567861edec0de664fad))
* chore(deps): bump github.com/zclconf/go-cty from 1.11.0 to 1.12.1 (#1020) (dependabot[bot])([16a79870](https://github.com/lacework/go-sdk/commit/16a798704cb9dd638118d28321efb7c999185867))
* ci: version bump to v1.3.2-dev (Lacework)([9a254aee](https://github.com/lacework/go-sdk/commit/9a254aee3c0096f35eb498369cd1a739fe3db078))

## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)
```
docker pull lacework/lacework-cli
```
