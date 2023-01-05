FROM alpine:latest
LABEL maintainer="tech-ally@lacework.net" \
      description="The Lacework CLI helps you manage the Lacework cloud security platform"
RUN apk add -U --no-cache ca-certificates
COPY LICENSE /
ADD bin/lacework-cli-linux-amd64 /usr/local/bin/lacework
ENTRYPOINT ["/usr/local/bin/lacework"]
