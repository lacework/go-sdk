#!/bin/bash
#
# Name::        lacework_cli_release.sh
# Description:: Use this script to prepare a new release on Github,
#               the automation will build cross-platform binaries,
#               compress all generated targets, generate shasum
#               hashes, and create a GH tag like v0.1.0 (using the
#               VERSION file inside the cli/ directory)
# Author::      Salim Afiune Maya (<afiune@lacework.net>)
#
set -eou pipefail

CLINAME=lacework-cli
VERSION=$(cat cli/VERSION)
TARGETS=(
  ${CLINAME}-darwin-386
  ${CLINAME}-darwin-amd64
  ${CLINAME}-windows-386.exe
  ${CLINAME}-windows-amd64.exe
  ${CLINAME}-linux-386
  ${CLINAME}-linux-amd64
)

main() {
  log "Preparing release v$VERSION"
  prerequisites
  build_cli_cross_platform
  compress_targets
  generate_shasums
  create_git_tag
}

create_git_tag() {
  local _tag="v$VERSION"
  log "Creating github tag: $_tag"
  git tag "$_tag"
  git push origin "$_tag"
  log "Go to https://github.com/lacework/go-sdk/releases and upload all files from 'bin/'"
}

prerequisites() {
  if ! command -v "gox" > /dev/null 2>&1; then
    warn "Required command 'gox' not found on PATH"
    warn "Try running 'make prepare'"
    exit 127
  fi

  local _branch=$(git rev-parse --abbrev-ref HEAD)
  if [ "$_branch" != "master" ]; then
    warn "Releases must be generated from the 'master' branch. (current $_branch)"
    warn "Switch to the master branch and try again."
    exit 127
  fi
}

clean_cache() {
  rm -rf bin/*
}

build_cli_cross_platform() {
  clean_cache
  make build-cli-cross-platform
}

generate_shasums() {
  ( cd bin/
    local _compressed
    log "Generating sha256sum Hashes"
    for target in ${TARGETS[*]}; do

      if [[ "$target" =~ /linux/ ]]; then
	_compressed="$target.tar.gz"
      else
	_compressed="$target.zip"
      fi

      log "bin/$_compressed.sha256sum"
      shasum -a 256 $_compressed > $_compressed.sha256sum 

    done
  )
}

# compress_targets will compress all targets and remove the raw
# binaries (already compressed), this is a release so we don't
# need the raw binaries anymore.
compress_targets() {
  log "Compressing target binaries"
  for target in ${TARGETS[*]}; do
    if [[ "$target" =~ /linux/ ]]; then
      tar -czvf "bin/${target}.tar.gz" "bin/${target}"
    else
      zip "bin/${target}.zip" "bin/${target}"
    fi
    rm -f "bin/${target}"
  done
}

log() {
  echo "--> ${CLINAME}: $1"
}

warn() {
  echo "xxx ${CLINAME}: $1" >&2
}

main || exit 99
