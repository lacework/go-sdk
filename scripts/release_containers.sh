#!/bin/bash
#

set -eou pipefail

# The repository where we are hosting the lacework-cli containers
readonly repository="lacework/lacework-cli"
# @afiune let us continue posting the Lacework CLI to the old reporitory for a few
readonly old_repo="techallylw/lacework-cli"
readonly project_name=lacework-cli

log() {
  echo "--> ${project_name}: $1"
}

# Make sure we have the binary needed for the SCRATCH container
if [ ! -f "bin/lacework-cli-linux-amd64" ]; then
  log "building Lacework CLI cross-platform"
  make build-cli-cross-platform
fi

# Enable docker experimental mode
log "enabling experimental mode to use/upload docker manifest"
mkdir -p ~/.docker
echo '{"experimental": "enabled"}' > ~/.docker/config.json

# Authenticate to dockerhub
echo "$DOCKERHUB_PASS" | docker login -u "$DOCKERHUB_USERNAME" --password-stdin

log "releasing container from SCRATCH"
docker build -t "${repository}:scratch" --no-cache .
docker push "${repository}:scratch"

# @afiune let us continue posting the Lacework CLI to the old reporitory for a few
log "releasing old container for SCRATCH"
docker image tag "${repository}:scratch" "${old_repo}:scratch"
docker push "${old_repo}:scratch"

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

  # @afiune let us continue posting the Lacework CLI to the old reporitory for a few
  log "releasing old container for ${dist}"
  docker image tag "${repository}:${dist}" "${old_repo}:${dist}"
  docker push "${old_repo}:${dist}"
done

log "creating docker manifest"
docker manifest create   "${repository}:latest"      \
                         "${repository}:scratch"     \
                         "${repository}:ubi-8"       \
                         "${repository}:centos-8"    \
                         "${repository}:debian-10"   \
                         "${repository}:ubuntu-1804" \
                         "${repository}:amazonlinux-2" --amend

log "pushing docker manifest"
docker manifest push "${repository}:latest" --purge

# @afiune let us continue posting the Lacework CLI to the old reporitory for a few
log "creating docker manifest for the old repository ${old_repo}"
docker manifest create   "${old_repo}:latest"      \
                         "${old_repo}:scratch"     \
                         "${old_repo}:ubi-8"       \
                         "${old_repo}:centos-8"    \
                         "${old_repo}:debian-10"   \
                         "${old_repo}:ubuntu-1804" \
                         "${old_repo}:amazonlinux-2" --amend

log "pushing docker manifest for the old repository ${old_repo}"
docker manifest push "${old_repo}:latest" --purge

log "All docker containers have been released! (https://hub.docker.com/repository/docker/${repository})"
log "All docker containers have been released! (https://hub.docker.com/repository/docker/${old_repo})"
