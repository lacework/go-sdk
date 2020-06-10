FROM  centos:8
LABEL maintainer="tech-ally@lacework.net" \
      description="The Lacework CLI helps you manage the Lacework cloud security platform"

COPY ./LICENSE /

# Fix (curl, libcurl-minimal)
# * CVE-2019-5436
# * CVE-2019-5482
# * CVE-2019-5481
ADD cli/images/centos-8/centos-8-stream.repo /etc/yum.repos.d/
RUN yum update -y

# Fix CVE-2019-6477
RUN yum remove bind-export-libs -y

RUN curl https://raw.githubusercontent.com/lacework/go-sdk/master/cli/install.sh | bash
