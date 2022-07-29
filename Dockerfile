FROM golang:1.18-alpine3.16 as alpine
RUN apk add -U --no-cache ca-certificates

COPY . /opt/go-sdk

WORKDIR /opt/go-sdk/cli

ENV CGO_ENABLED=0
RUN go build -ldflags="-X github.com/lacework/go-sdk/cli/cmd.Version=$(cat /opt/go-sdk/VERSION)" -o /opt/lacework 

FROM  scratch
LABEL maintainer="tech-ally@lacework.net" \
      description="The Lacework CLI helps you manage the Lacework cloud security platform"

COPY LICENSE /
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine /opt/lacework /usr/local/bin/lacework
#ADD bin/lacework-cli-linux-amd64 /usr/local/bin/lacework
#ENTRYPOINT ["/usr/local/bin/lacework"]

ENTRYPOINT ["/usr/local/bin/lacework"]
