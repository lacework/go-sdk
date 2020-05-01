FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM  scratch
LABEL maintainer="tech-ally@lacework.net" \
      description="The Lacework CLI helps you manage the Lacework cloud security platform"

COPY LICENSE /
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ADD bin/lacework-cli-linux-amd64 /usr/local/bin/lacework
ENTRYPOINT ["/usr/local/bin/lacework"]
