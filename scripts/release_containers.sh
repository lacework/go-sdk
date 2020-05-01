#!/bin/bash
#

set -eou pipefail

# The repository where we are hosting the lacework-cli containers
# TODO @afiune switch it to "lacework/lacework-cli" repository
readonly repository="techallylw/lacework-cli"
readonly project_name=lacework-cli

log() {
  echo "--> ${project_name}: $1"
}

log "releasing container from SCRATCH"
docker build -t "${repository}:scratch" .
docker push "${repository}:scratch"

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
  docker build -f "cli/images/${dist}/Dockerfile" -t "${repository}:ubi-8" .
  docker push "${repository}:${dist}"
done

if ! docker manifest inspect "$repository"; then
  log "creating docker manifest"
  docker manifest create "${repository}:latest"      \
                         "${repository}:scratch"     \
                         "${repository}:ubi-8"       \
                         "${repository}:centos-8"    \
                         "${repository}:debian-10"   \
                         "${repository}:ubuntu-1804" \
                         "${repository}:amazonlinux-2"
fi

log "pushing docker manifest"
docker manifest push "${repository}:latest"

log "All docker containers have been released! (https://hub.docker.com/repository/docker/${repository})"
