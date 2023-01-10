#!/bin/bash
#

set -eou pipefail

# The repository where we are hosting the lacework-cli container
readonly repository="lacework/lacework-cli"
readonly project_name=lacework-cli
VERSION=$(cat VERSION)

log() {
  echo "--> ${project_name}: $1"
}

# Make sure we have the binary needed for the SCRATCH container
if [ ! -f "bin/lacework-cli-linux-amd64" ]; then
  log "building Lacework CLI cross-platform"
  make build-cli-cross-platform
fi

# Authenticate to dockerhub
echo "$DOCKERHUB_PASS" | docker login -u "$DOCKERHUB_USERNAME" --password-stdin

log "releasing docker container v${VERSION} from SCRATCH"

docker context create cf-environment
docker buildx create --name multiarch-builder \
  --driver docker-container \
  --platform linux/386,linux/amd64,linux/arm64,linux/armhf \
  --use cf-environment
docker buildx build . --push --no-cache \
  -t "${repository}:v${VERSION}" \
  -t "${repository}:latest" \
  --platform linux/386,linux/amd64,linux/arm64,linux/armhf

log "Docker container released! (https://hub.docker.com/repository/docker/${repository})"
