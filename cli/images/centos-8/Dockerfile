FROM  centos:8
LABEL maintainer="tech-ally@lacework.net" \
      description="The Lacework CLI helps you manage the Lacework cloud security platform"

# From December 2021 CentOS 8 Life was ended
# (@afiune) Workaround: use a mirror
RUN sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-Linux-*
RUN sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-Linux-*

RUN yum update -y
COPY ./LICENSE /
RUN curl https://raw.githubusercontent.com/lacework/go-sdk/main/cli/install.sh | bash
