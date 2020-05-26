#!/bin/bash
#
# Name::        release.sh
# Description:: Use this script to prepare a new release on Github,
#               the automation will build cross-platform binaries,
#               compress all generated targets, generate shasum
#               hashes, and create a GH tag like 'v0.1.0'
#               (using the VERSION file)
# Author::      Salim Afiune Maya (<afiune@lacework.net>)
#
set -eou pipefail

readonly project_name=go-sdk
readonly package_name=lacework-cli
readonly binary_name=lacework
readonly docker_org=techallylw
readonly docker_tags=(
  latest
  scratch
  ubi-8
  centos-8
  debian-10
  ubuntu-1804
  amazonlinux-2
#  windows-nanoserver
)

VERSION=$(cat VERSION)
TARGETS=(
  ${package_name}-darwin-386
  ${package_name}-darwin-amd64
  ${package_name}-windows-386.exe
  ${package_name}-windows-amd64.exe
  ${package_name}-linux-386
  ${package_name}-linux-amd64
)

usage() {
  local _cmd
  _cmd="$(basename "${0}")"
  cat <<USAGE
${_cmd}: A tool that helps you do releases!

Use this script to prepare a new Github release, the automation will build
cross-platform binaries, compress all generated targets, generate shasum hashes,
generate release notes, update the changelog and create a Github tag like 'v0.1.0'.

USAGE:
    ${_cmd} [command]

COMMANDS:
    prepare    Generates release notes, updates version and CHANGELOG.md
    publish    Builds binaries, shasums and creates a Github tag like 'v0.1.0'

Update version after release:
    version [kind] Prepare the version after release, it adds the '-dev' tag
                   Kinds of version bumps: [patch, minor, major]
USAGE
}

main() {
  case "${1:-}" in
    prepare)
      prepare_release
      ;;
    publish)
      publish_release
      ;;
    version)
      bump_version $2
      ;;
    *)
      usage
      ;;
  esac
}

prepare_release() {
  log "preparing new release"
  prerequisites
  remove_dev_version
  generate_release_notes
  update_changelog
  push_release
}

publish_release() {
  log "releasing v$VERSION"
  prerequisites
  release_check
  clean_cache
  build_cli_cross_platform
  compress_targets
  generate_shasums
  tag_release
}

update_changelog() {
  log "updating CHANGELOG.md"
  _changelog=$(cat CHANGELOG.md)
  echo "# v$VERSION" > CHANGELOG.md
  echo "" >> CHANGELOG.md
  echo "$(cat CHANGES.md)" >> CHANGELOG.md
  echo "---" >> CHANGELOG.md
  echo "$_changelog" >> CHANGELOG.md
  # clean changes file since we don't need it anymore
  rm CHANGES.md
}

load_list_of_changes() {
  latest_version=$(find_latest_version)
  local _list_of_changes=$(git log --no-merges --pretty="* %s (%an)([%h](https://github.com/lacework/${project_name}/commit/%H))" ${latest_version}..master)
  echo "## Features" > CHANGES.md
  echo "$_list_of_changes" | grep "\* feat[:(]" >> CHANGES.md
  echo "## Refactor" >> CHANGES.md
  echo "$_list_of_changes" | grep "\* refactor[:(]" >> CHANGES.md
  echo "## Performance Improvements" >> CHANGES.md
  echo "$_list_of_changes" | grep "\* perf[:(]" >> CHANGES.md
  echo "## Bug Fixes" >> CHANGES.md
  echo "$_list_of_changes" | grep "\* fix[:(]" >> CHANGES.md
  echo "## Documentation Updates" >> CHANGES.md
  echo "$_list_of_changes" | grep "\* doc[:(]" >> CHANGES.md
  echo "$_list_of_changes" | grep "\* docs[:(]" >> CHANGES.md
  echo "## Other Changes" >> CHANGES.md
  echo "$_list_of_changes" | grep "\* style[:(]" >> CHANGES.md
  echo "$_list_of_changes" | grep "\* chore[:(]" >> CHANGES.md
  echo "$_list_of_changes" | grep "\* build[:(]" >> CHANGES.md
  echo "$_list_of_changes" | grep "\* ci[:(]" >> CHANGES.md
  echo "$_list_of_changes" | grep "\* test[:(]" >> CHANGES.md
}

generate_release_notes() {
  log "generating release notes at RELEASE_NOTES.md"
  load_list_of_changes
  echo "# Release Notes" > RELEASE_NOTES.md
  echo "Another day, another release. These are the release notes for the version \`v$VERSION\`." >> RELEASE_NOTES.md
  echo "" >> RELEASE_NOTES.md
  echo "$(cat CHANGES.md)" >> RELEASE_NOTES.md

  # Add Docker Images Footer
  echo "" >> RELEASE_NOTES.md
  echo "## Docker Images" >> RELEASE_NOTES.md
  for tag in "${docker_tags[@]}"; do
    echo "* \`docker pull ${docker_org}/${package_name}:${tag}\`" >> RELEASE_NOTES.md
  done
}

push_release() {
  log "commiting and pushing the release to github"
  git checkout -B release
  git commit -am "Release v$VERSION"
  git push origin release
  log ""
  log "Follow the above url and open a pull request"
}

tag_release() {
  local _tag="v$VERSION"
  log "creating github tag: $_tag"
  git tag "$_tag"
  git push origin "$_tag"
  log "go to https://github.com/lacework/${project_name}/releases/tag/${_tag} and upload all files from 'bin/'"
}

release_check() {
  if git ls-remote --tags 2>/dev/null | grep "tags/v$VERSION" >/dev/null; then
    warn "The git tag 'v$VERSION' already exists at github.com"
    warn "This is usually the case where a release wasn't finished properly"
    warn "Remove the tag from the remote and try to release again"
    warn ""
    warn "To remove the tag run the following commands:"
    warn "  git tag -d v$VERSION"
    warn "  git push --delete origin v$VERSION"
    exit 127
  fi
}

prerequisites() {
  # gox is the tool we use to generate cross-platform binaries
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

  local _unsaved_changes=$(git status -s)
  if [ "$_unsaved_changes" != "" ]; then
    warn "You have unsaved changes in the master branch. Are you resuming a release?"
    warn "To resume a release you have to start over, to remove all unsaved changes run the command:"
    warn "  git reset --hard origin/master"
    exit 127
  fi
}

find_latest_version() {
  local _pattern="v[0-9]\+.[0-9]\+.[0-9]\+"
  local _versions
  _versions=$(git ls-remote --tags --quiet | grep $_pattern | tr '/' ' ' | awk '{print $NF}')
  echo "$_versions" | tr '.' ' ' | sort -nr -k 1 -k 2 -k 3 | tr ' ' '.' | head -1
}

remove_dev_version() {
  if [[ "$VERSION" =~ "dev" ]]; then
    echo $VERSION | cut -d- -f1 > VERSION
    VERSION=$(cat VERSION)
    log "updated version for release v$VERSION"
  fi
}

bump_version() {
  log "updating version after release"
  prerequisites
  latest_version=$(find_latest_version)

  if [[ "v$VERSION" == "$latest_version" ]]; then
    case "${1:-}" in
      major)
        echo $VERSION | awk -F. '{printf("%d.%d.%d-dev", $1+1, $2, $3)}' > VERSION
        ;;
      minor)
        echo $VERSION | awk -F. '{printf("%d.%d.%d-dev", $1, $2+1, $3)}' > VERSION
        ;;
      *)
        echo $VERSION | awk -F. '{printf("%d.%d.%d-dev", $1, $2, $3+1)}' > VERSION
        ;;
    esac
    VERSION=$(cat VERSION)
    log "version bumped from $latest_version to v$VERSION"
  else
    log "skipping version bump. Already bumped to v$VERSION"
    return
  fi

  log "pushing version bump to 'master'. [Press Enter to continue]"
  read
  git add VERSION
  git commit -m "version bump to v$VERSION"
  git push origin master
}

clean_cache() {
  log "cleaning cache bin/ directory"
  rm -rf bin/*
}

build_cli_cross_platform() {
  log "building cross-platform binaries"
  make build-cli-cross-platform
}

generate_shasums() {
  ( cd bin/
    local _compressed
    log "generating sha256sum Hashes"
    for target in ${TARGETS[*]}; do

      if [[ "$target" =~ linux ]]; then
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
  log "compressing target binaries"
  local _target_with_ext
  local _cli_name

  for target in ${TARGETS[*]}; do
    if [[ "$target" =~ exe ]]; then
      _cli_name="bin/${binary_name}.exe"
    else
      _cli_name="bin/${binary_name}"
    fi

    mv "bin/${target}" "$_cli_name"

    if [[ "$target" =~ linux ]]; then
      _target_with_ext="bin/${target}.tar.gz"
      tar -czvf "$_target_with_ext" "$_cli_name" 2>/dev/null
    else
      _target_with_ext="bin/${target}.zip"
      zip "$_target_with_ext" "$_cli_name" >/dev/null
    fi

    log $_target_with_ext
    rm -f "$_cli_name"
  done
}

log() {
  echo "--> ${project_name}: $1"
}

warn() {
  echo "xxx ${project_name}: $1" >&2
}

main "$@" || exit 99
