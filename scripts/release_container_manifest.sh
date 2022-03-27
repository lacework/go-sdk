#!/bin/bash
#

set -eou pipefail

# The repository where we are hosting the lacework-cli containers
readonly repository="lacework/lacework-cli"
readonly project_name=lacework-cli

log() {
  echo "--> ${project_name}: $1"
}

# Enable docker experimental mode
log "enabling experimental mode to use/upload docker manifest"
_docker_config=$HOME/.docker/config.json
mkdir -p ~/.docker
[[ ! -f "$_docker_config" ]] && echo '{"experimental": "enabled"}' > "$_docker_config"

# Authenticate to dockerhub if needed
[[ "unset" != "${DOCKERHUB_PASS:-unset}" ]] && echo "$DOCKERHUB_PASS" | docker \
  login -u "$DOCKERHUB_USERNAME" \
  --password-stdin

# when updating the distributions below, please make sure to update
# the script 'release.sh' inside the 'script/' folder
distros=(
  scratch
  ubi-8
  debian-10
  ubuntu-1804
  amazonlinux-2
#  windows-nanoserver
)

for dist in "${distros[@]}"; do
  log "pulling container for ${dist}"
  docker pull "${repository}:${dist}"
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

log "Docker manifest released! (https://hub.docker.com/repository/docker/${repository})"
