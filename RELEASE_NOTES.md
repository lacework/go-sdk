# Release Notes
Another day, another release. These are the release notes for the version `v0.2.7`.

## Container Vulnerability Assessments in HTML format 
You have now the ability to output container vulnerability assessments in HTML format, this will help teams have clear, actionable insights to understand vulnerabilities so they can remediate.

Use the new flag `--html` in the following commands:
* `lacework vulnerability container scan`
* `lacework vulnerability container scan-status`
* `lacework vulnerability container show-assessment`

The flag will generate a vulnerability assessment in HTML format in addition to the regular human-readable output.

```
$ lacework vuln ctr show sha256:4861f2bca4e9087558cf9594249da90f42fb781f9b8a857db21e25abc2c57486 --html
                                  CONTAINER IMAGE DETAILS                                 |        VULNERABILITIES
------------------------------------------------------------------------------------------+---------------------------------
    ID          sha256:948cb0da5637b7d7662168a030930e9ab115c302b7d6833c21cb20dac9deb2e4   |   SEVERITY   COUNT   FIXABLE
    Digest      sha256:4861f2bca4e9087558cf9594249da90f42fb781f9b8a857db21e25abc2c57486   | -----------+-------+----------
    Registry    index.docker.io                                                           |   Critical      33        28
    Repository  dhubtesting/test                                                          |   High         220       192
    Size        163.2 MB                                                                  |   Medium       432       310
    Created At  2020-10-06T07:58:58+0000                                                  |   Low          246       124
    Tags        sentry_8.6-onbuild_06102020075836                                         |   Info         219        38
                                                                                          |
The container vulnerability assessment was stored at 'dhubtesting-test-sha256:4861f2bca4e9087558cf9594249da90f42fb781f9b8a857db21e25abc2c57486.html'
```

The HTML is a standalone file that can be shared without additional artifacts!

<img width="1904" alt="Screen Shot 2020-11-12 at 10 31 15 PM" src="https://user-images.githubusercontent.com/5712253/98998667-c2df5080-24f3-11eb-8c99-6b9aee637324.png">

## Features
* feat(cli): enable html copy to clipboard icons (Salim Afiune Maya)([ec2d1fa](https://github.com/lacework/go-sdk/commit/ec2d1fa5b796e2b51e49a850deff01d2f64ded18))
* feat(cli): HTML format for vulnerability assessments (Salim Afiune Maya)([00c2f43](https://github.com/lacework/go-sdk/commit/00c2f43613e554afd8ed283cbc12eb0b8eed0179))
* feat(cli): add ARM support (#236) (Salim Afiune)([821b8e6](https://github.com/lacework/go-sdk/commit/821b8e699e61eefda7d287a71b08ef26382a4ad7))
## Bug Fixes
* fix(cli): remove html column sort icons (Salim Afiune Maya)([dc4c0f6](https://github.com/lacework/go-sdk/commit/dc4c0f64055bfaada503b8f4f21ceda707bc5e55))
## Other Changes
* ci(fix) Update CI test node (#233) (Scott Ford)([ddbf86e](https://github.com/lacework/go-sdk/commit/ddbf86e8fbf9053af43bab9d57c04645383e529e))

## Docker Images
* `docker pull lacework/lacework-cli:latest`
* `docker pull lacework/lacework-cli:scratch`
* `docker pull lacework/lacework-cli:ubi-8`
* `docker pull lacework/lacework-cli:centos-8`
* `docker pull lacework/lacework-cli:debian-10`
* `docker pull lacework/lacework-cli:ubuntu-1804`
* `docker pull lacework/lacework-cli:amazonlinux-2`
