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
readonly org_name=lacework
readonly package_name=lacework-cli
readonly binary_name=lacework
readonly docker_org=lacework
readonly git_user="Lacework Inc."
readonly git_email="tech-ally@lacework.net"

VERSION=$(cat VERSION)
TARGETS=(
  ${package_name}-darwin-amd64
  ${package_name}-darwin-arm64
  ${package_name}-windows-386.exe
  ${package_name}-windows-amd64.exe
  ${package_name}-linux-386
  ${package_name}-linux-amd64
  ${package_name}-linux-arm
  ${package_name}-linux-arm64
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
    verify     Check if the release is ready to be applied
    trigger    Trigger a release by creating a git tag
    publish    Download binaries, generate shasums and creates a Github tag like 'v0.1.0'
    build      Builds binaries and upload them to s3
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
    build)
     build_artifacts
      ;;
    verify)
      verify_release
      ;;
    trigger)
      trigger_release
      ;;
    *)
      usage
      ;;
  esac
}

trigger_release() {
  if [[ "$VERSION" =~ "-dev" ]]; then
      log "No release needed. (VERSION=${VERSION})"
      log ""
      log "Read more about the release process at:"
      log "  - https://github.com/${org_name}/${project_name}/wiki/Release-Process"
    else
      log "VERSION ready to be released to 'x.y.z' tag. Triggering a release!"
      log ""
      tag_release
      bump_version
  fi
}

verify_release() {
  log "verifying new release"
  _changed_file=$(git whatchanged --name-only --pretty="" origin..HEAD)
  _required_files_for_release=(
    RELEASE_NOTES.md
    CHANGELOG.md
    VERSION
    api/version.go
  )
  for f in "${_required_files_for_release[@]}"; do
    if [[ "$_changed_file" =~ "$f" ]]; then
      log "(required) '$f' has been modified. Great!"
    else
      warn "$f needs to be updated"
      warn ""
      warn "Read more about the release process at:"
      warn "  - https://github.com/${org_name}/${project_name}/wiki/Release-Process"
      exit 123
    fi
  done

  if [[ "$VERSION" =~ "-dev" ]]; then
      warn "the 'VERSION' needs to be cleaned up to be only 'x.y.z' tag"
      warn ""
      warn "Read more about the release process at:"
      warn "  - https://github.com/${org_name}/${project_name}/wiki/Release-Process"
      exit 123
    else
      log "(required) VERSION has been cleaned up to 'x.y.z' tag. Great!"
  fi
}

prepare_release() {
  log "preparing new release"
  prerequisites
  remove_tag_version
  check_for_minor_version_bump
  cli_generate_files
  generate_release_notes
  update_changelog
  push_release
  open_pull_request
}

build_artifacts() {
  log "building artifacts for v$VERSION"
  clean_cache
  build_cli_cross_platform
  upload_artifacts
}

publish_release() {
  log "releasing v$VERSION"
  download_artifacts
  compress_targets
  generate_shasums
  create_release
}

cli_generate_files() {
  make generate-docs
  make generate-databox
}

download_artifacts() {
  log "downloading signed artifacts for v$VERSION"
  aws s3 sync "s3://lacework-cli/builds/v${VERSION}/signed" bin/signed

  while [ ! -n "$(ls -1 bin/signed/.completed 2>/dev/null)"  ]; do
    log "waiting for signed artifacts..."
    sleep 5
    aws s3 sync "s3://lacework-cli/builds/v${VERSION}/signed" bin/signed
  done

  # sync everything
  aws s3 sync "s3://lacework-cli/builds/v${VERSION}" bin/

  log "moving signed/ artifacts to bin/ directory"
  mv bin/signed/lacework* bin/.
}

upload_artifacts() {
  log "uploading artifacts for v$VERSION"
  aws s3 sync bin/ "s3://lacework-cli/builds/v$VERSION"
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

release_contains_features() {
  latest_version=$(find_latest_version)
  log "found latest version: $latest_version"
  git log --no-merges --pretty="%s" ${latest_version}..main | grep "feat[:(]" >/dev/null
  return $?
}

load_list_of_changes() {
  latest_version=$(find_latest_version)
  local _list_of_changes=$(git log --no-merges --pretty="* %s (%an)([%h](https://github.com/${org_name}/${project_name}/commit/%H))" ${latest_version}..main)

  # init changes file
  true > CHANGES.md

  _feat=$(echo "$_list_of_changes" | grep "\* feat[:(]")
  _refactor=$(echo "$_list_of_changes" | grep "\* refactor[:(]")
  _perf=$(echo "$_list_of_changes" | grep "\* perf[:(]")
  _fix=$(echo "$_list_of_changes" | grep "\* fix[:(]")
  _doc=$(echo "$_list_of_changes" | grep "\* doc[:(]")
  _docs=$(echo "$_list_of_changes" | grep "\* docs[:(]")
  _metric=$(echo "$_list_of_changes" | grep "\* metric[:(]")
  _style=$(echo "$_list_of_changes" | grep "\* style[:(]")
  _chore=$(echo "$_list_of_changes" | grep "\* chore[:(]")
  _build=$(echo "$_list_of_changes" | grep "\* build[:(]")
  _ci=$(echo "$_list_of_changes" | grep "\* ci[:(]")
  _test=$(echo "$_list_of_changes" | grep "\* test[:(]")

  if [ "$_feat" != "" ]; then
    echo "## Features" >> CHANGES.md
    echo "$_feat" >> CHANGES.md
  fi

  if [ "$_refactor" != "" ]; then
    echo "## Refactor" >> CHANGES.md
    echo "$_refactor" >> CHANGES.md
  fi

  if [ "$_perf" != "" ]; then
    echo "## Performance Improvements" >> CHANGES.md
    echo "$_perf" >> CHANGES.md
  fi

  if [ "$_fix" != "" ]; then
    echo "## Bug Fixes" >> CHANGES.md
    echo "$_fix" >> CHANGES.md
  fi

  if [ "${_docs}${_doc}" != "" ]; then
    echo "## Documentation Updates" >> CHANGES.md
    if [ "$_doc" != "" ]; then echo "$_doc" >> CHANGES.md; fi
    if [ "$_docs" != "" ]; then echo "$_docs" >> CHANGES.md; fi
  fi

  if [ "${_style}${_chore}${_build}${_ci}${_test}" != "" ]; then
    echo "## Other Changes" >> CHANGES.md
    if [ "$_style" != "" ]; then echo "$_style" >> CHANGES.md; fi
    if [ "$_chore" != "" ]; then echo "$_chore" >> CHANGES.md; fi
    if [ "$_build" != "" ]; then echo "$_build" >> CHANGES.md; fi
    if [ "$_ci" != "" ]; then echo "$_ci" >> CHANGES.md; fi
    if [ "$_metric" != "" ]; then echo "$_metric" >> CHANGES.md; fi
    if [ "$_test" != "" ]; then echo "$_test" >> CHANGES.md; fi
  fi
}

generate_release_notes() {
  log "generating release notes at RELEASE_NOTES.md"
  load_list_of_changes
  echo "# Release Notes" > RELEASE_NOTES.md
  echo "Another day, another release. These are the release notes for the version \`v$VERSION\`." >> RELEASE_NOTES.md
  echo "" >> RELEASE_NOTES.md
  echo "$(cat CHANGES.md)" >> RELEASE_NOTES.md

  # Add Docker Image Footer
  echo "" >> RELEASE_NOTES.md
  echo "## :whale: [Docker Image](https://hub.docker.com/r/lacework/lacework-cli)" >> RELEASE_NOTES.md
  echo '```' >> RELEASE_NOTES.md
  echo "docker pull ${docker_org}/${package_name}" >> RELEASE_NOTES.md
  echo '```' >> RELEASE_NOTES.md
}

push_release() {
  log "commiting and pushing the release to github"
  _version_no_tag=$(echo $VERSION | awk -F. '{printf("%d.%d.%d", $1, $2, $3)}')
  if [ "$CI" != "" ]; then
    log "configuring git user email, user name and signingkey"
    git config --global user.email $git_email
    git config --global user.name $git_user
    git config --global user.signingkey $GPG_SIGNING_KEY
  fi
  git checkout -B release
  git commit -sS -am "release: v$_version_no_tag"
  git push origin release -f
}

open_pull_request() {
  local _body="/tmp/pr.json"
  local _pr="/tmp/pr.out"

  log "opening GH pull request"
  generate_pr_body "$_body"
  curl -XPOST -H "Authorization: token $GITHUB_TOKEN" --data  "@$_body" \
        https://api.github.com/repos/${org_name}/${project_name}/pulls > $_pr

  # @afiune just to debug the issue where the field `html_url` comes as `null`
  echo "$_pr" | jq .

  _pr_url=$(jq .html_url $_pr)
  log ""
  log "It is time to review the release!"
  log "    $_pr_url"
}

tag_release() {
  local _tag="v$VERSION"
  log "creating github tag: $_tag"
  git tag "$_tag"
  git push origin "$_tag"
}

prerequisites() {
  # gox is the tool we use to generate cross-platform binaries
  if ! command -v "gox" > /dev/null 2>&1; then
    warn "Required command 'gox' not found on PATH"
    warn "Try running 'make prepare'"
    exit 127
  fi

  local _branch=$(git rev-parse --abbrev-ref HEAD)
  if [ "$_branch" != "main" ]; then
    warn "Releases must be generated from the 'main' branch. (current $_branch)"
    warn "Switch to the main branch and try again."
    exit 127
  fi

  local _unsaved_changes=$(git status -s)
  if [ "$_unsaved_changes" != "" ]; then
    warn "You have unsaved changes in the main branch. Are you resuming a release?"
    warn "To resume a release you have to start over, to remove all unsaved changes run the command:"
    warn "  git reset --hard origin/main"
    exit 127
  fi
}

find_latest_version() {
  local _pattern="v[0-9]\+.[0-9]\+.[0-9]\+"
  local _versions
  _versions=$(git ls-remote --tags --quiet | grep $_pattern | tr '/' ' ' | awk '{print $NF}')
  echo "$_versions" | tr '.' ' ' | sort -r -V | tr ' ' '.' | head -1
}

add_tag_version() {
  _tag=${1:-dev}
  echo $VERSION | awk -F. '{printf("%d.%d.%d-'$_tag'", $1, $2, $3)}' > VERSION
  VERSION=$(cat VERSION)
  scripts/version_updater.sh
  log "updated version to v$VERSION"
}

check_for_minor_version_bump() {
  if release_contains_features; then
    log "new feature detected, minor version bump"
    echo $VERSION | awk -F. '{printf("%d.%d.0", $1, $2+1)}' > VERSION
    VERSION=$(cat VERSION)
    scripts/version_updater.sh
    log "updated version to v$VERSION"
  fi
}

remove_tag_version() {
  echo $VERSION | awk -F. '{printf("%d.%d.%d", $1, $2, $3)}' > VERSION
  VERSION=$(cat VERSION)
  scripts/version_updater.sh
  log "updated version to v$VERSION"
}

bump_version() {
  log "updating version after tagging release"
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
    scripts/version_updater.sh
    log "version bumped from $latest_version to v$VERSION"
  else
    log "skipping version bump. Already bumped to v$VERSION"
    return
  fi

  log "commiting and pushing the vertion bump to github"
  if [ "$CI" != "" ]; then
    log "configuring git user email, user name and signingkey"
    git config --global user.email $git_email
    git config --global user.name $git_user
    git config --global user.signingkey $GPG_SIGNING_KEY
  fi
  git add VERSION
  git add api/version.go # file genereted by scripts/version_updater.sh
  git commit -sS -m "ci: version bump to v$VERSION"
  git push origin main
}

clean_cache() {
  log "cleaning cache bin/ directory"
  rm -rf bin/*
}

build_cli_cross_platform() {
  log "building cross-platform binaries"
  make build-cli-cross-platform
  log "creating signed/ folder"
  mkdir -p bin/signed
  touch bin/signed/.keeper
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

    cp "bin/${target}" "$_cli_name"

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

create_release() {
  local _tag
  _tag=$(git describe --tags)
  local _body="/tmp/release.json"
  local _release="/tmp/release.out"

  log "generating GH release $_tag"
  generate_release_body "$_body"
  curl -XPOST -H "Authorization: token $GITHUB_TOKEN" --data  "@$_body" \
        https://api.github.com/repos/${org_name}/${project_name}/releases > $_release

  local _content_type
  local _artifact
  local _upload_url
  _upload_url=$(jq .upload_url $_release | sed 's/"//g' | cut -d{ -f1)
  log "uploading artifacts to GH release at ($_upload_url)"
  for target in ${TARGETS[*]}; do

    if [[ "$target" =~ linux ]]; then
      _artifact="$target.tar.gz"
      _content_type="application/gzip"
    else
      _artifact="$target.zip"
      _content_type="application/zip"
    fi

    log "uploading bin/$_artifact.sha256sum"
    curl -s -H "Authorization: token $GITHUB_TOKEN"  \
        -H "Content-Type: $_content_type"            \
        --data-binary "@bin/${_artifact}.sha256sum"  \
        "${_upload_url}?name=${_artifact}.sha256sum"

    log "uploading bin/$_artifact"
    curl -s -H "Authorization: token $GITHUB_TOKEN"  \
        -H "Content-Type: $_content_type"            \
        --data-binary "@bin/$_artifact"              \
        "${_upload_url}?name=${_artifact}"

  done

  log "the release has been completed!"
  log ""
  log " -> https://github.com/${org_name}/${project_name}/releases/tag/${_tag}"
}

generate_pr_body() {
  _file=${1:-pr.json}
  _version_no_tag=$(echo $VERSION | awk -F. '{printf("%d.%d.%d", $1, $2, $3)}')
  _release_notes=$(jq -aRs .  <<< cat RELEASE_NOTES.md)
  cat <<EOF > $_file
{
  "base": "main",
  "head": "release",
  "title": "Release v$_version_no_tag",
  "body": $_release_notes
}
EOF
}

generate_release_body() {
  _file=${1:-release.json}
  _tag=$(git describe --tags)
  _release_notes=$(jq -aRs .  <<< cat RELEASE_NOTES.md)
  cat <<EOF > $_file
{
  "tag_name": "$_tag",
  "name": "$_tag",
  "draft": false,
  "prerelease": false,
  "body": $_release_notes
}
EOF
}

log() {
  echo "--> ${project_name}: $1"
}

warn() {
  echo "xxx ${project_name}: $1" >&2
}

main "$@" || exit 99
