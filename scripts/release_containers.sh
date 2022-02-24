#!/bin/bash
#

set -eou pipefail

# The repository where we are hosting the lacework-cli containers
readonly repository="lacework/lacework-cli"
readonly project_name=lacework-cli

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

log "releasing container from SCRATCH"
docker build -t "${repository}:scratch" --no-cache .
docker push "${repository}:scratch"

# when updating the distributions below, please make sure to update
# the script 'release.sh' inside the 'script/' folder
distros=(
  ubi-8
  centos-8
  debian-10
  ubuntu-1804
  amazonlinux-2
#  windows-nanoserver
)

for dist in "${distros[@]}"; do
  log "releasing container for ${dist}"
  docker build -f "cli/images/${dist}/Dockerfile" --no-cache -t "${repository}:${dist}" .
  docker push "${repository}:${dist}"
done

scripts/release_container_manifest.sh

log "All docker containers have been released! (https://hub.docker.com/repository/docker/${repository})"
