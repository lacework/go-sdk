FROM  ubuntu:18.04
LABEL maintainer="tech-ally@lacework.net" \
      description="The Lacework CLI helps you manage the Lacework cloud security platform"

RUN apt-get update
RUN apt-get install curl -y
COPY ./LICENSE /
RUN curl https://raw.githubusercontent.com/lacework/go-sdk/master/cli/install.sh | bash
