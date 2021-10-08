#!/bin/bash
#
set -eou pipefail

# If the variable $LW_DEBUG is set, print all shell commands executed
if [ -n "${LW_DEBUG:-}" ]; then set -x; fi

readonly github_releases="https://github.com/lacework/go-sdk/releases"
readonly default_install_dir=/usr/local/bin
readonly package_name=lacework-cli
readonly binary_name=lacework

usage() {
  local _cmd
  _cmd="$(basename "${0}")"
  cat <<USAGE
${_cmd}: Installs the Lacework Command Line Interface.

Authors: Technology Alliances <tech-ally@lacework.net>

USAGE:
    ${_cmd} [FLAGS]

FLAGS:
    -h    Prints help information
    -v    Specifies a version (ex: v0.1.0)
    -d    The installation directory (default: $default_install_dir)
    -t    Specifies the target of the program to download (default: linux-amd64)
USAGE
}

main() {
  version=""

  # Parse command line flags and options
  while getopts "h:v:t:d:" opt; do
    case "${opt}" in
      h)
        usage
        exit 0
        ;;
      v)
        version="${OPTARG}"
        ;;
      t)
        target="${OPTARG}"
        check_target
        ;;
      d)
        installation_dir="${OPTARG}"
        ;;
      \?)
        echo "" >&2
        usage >&2
        exit_with "Invalid option" 1
        ;;
    esac
  done

  if [ -z "${installation_dir:-}" ]; then
    installation_dir=$default_install_dir
  fi

  if [ -z "${target:-}" ]; then
    detect_platform
  fi

  log "Installing the Lacework CLI"
  create_workdir
  download_archive "$version" "$target"
  verify_archive
  extract_archive
  install_cli
  print_cli_version
  log "The Lacework CLI has been successfully installed."
}

create_workdir() {
  if [ -d /var/tmp ]; then
    local _tmp=/var/tmp
  else
    local _tmp=/tmp
  fi

  workdir="$(mktemp -d -p "$_tmp" 2> /dev/null || mktemp -d "${_tmp}/lacework.XXXX")"
  # add a trap to clean up work directory
  trap 'code=$?; rm -rf $workdir; exit $code' INT TERM EXIT
  cd "${workdir}"
}

check_target() {
  if [[ ! $target =~ "-" ]]; then
    exit_with "malformed target '${target}' (format: system-arch)" 5
  fi

  sys=$(echo "$target" | cut -d- -f1)
  if [ -z "${sys}" ]; then
    exit_with "malformed target '${target}' (format: system-arch)" 5
  fi

  arch=$(echo "$target" | cut -d- -f2)
  if [ -z "${arch}" ]; then
    exit_with "malformed target '${target}' (format: system-arch)" 5
  fi

  verify_platform
}

detect_platform() {
  local _ostype
  _ostype="$(uname -s)"

  case "${_ostype}" in
    Darwin|Linux)
      sys="$(uname -s | tr '[:upper:]' '[:lower:]')"
      arch="$(uname -m | tr '[:upper:]' '[:lower:]')"
      ;;
    *)
      exit_with "unable to determine OS platform type: ${_ostype}" 2
      ;;
  esac

  verify_platform

  target="${sys}-${arch}"
}

verify_platform() {
  case "${sys}" in
    darwin)
      ext=zip
      shasum_cmd="shasum -a 256"
      ;;
    linux)
      ext=tar.gz
      shasum_cmd="sha256sum"
      ;;
    *)
      exit_with "unable to determine system type, perhaps is not supported: ${sys}" 3
      ;;
  esac

  # The following architectures match our cross-platform build process
  # https://golang.org/doc/install/source#environment
  case "${arch}" in
    x86_64 | amd64)
      arch=amd64
      ;;
   i686 | 386)
      arch=386
      ;;
   aarch64* | armv8* | arm64)
      arch=arm64
      ;;
   armv7* | armv6* | arm)
      arch=arm
      ;;
    *)
      exit_with "architecture not supported: ${arch}" 3
      ;;
  esac
}

download_archive() {
  local _version="${1:-latest}"
  local -r _target="${2:?}"
  local url

  if [ "$_version" == "latest" ]; then
    url="${github_releases}/latest/download/${package_name}-${_target}.${ext}"
  else
    url="${github_releases}/download/${_version}/${package_name}-${_target}.${ext}"
  fi

  download_file "${url}" "${workdir}/${package_name}-${_version}.${ext}"
  download_file "${url}.sha256sum" "${workdir}/${package_name}-${_version}.${ext}.sha256sum"

  archive="${package_name}-${_target}.${ext}"
  sha_file="${package_name}-${_target}.${ext}.sha256sum"

  mv -v "${workdir}/${package_name}-${_version}.${ext}" "${archive}"
  mv -v "${workdir}/${package_name}-${_version}.${ext}.sha256sum" "${sha_file}"
}

verify_archive() {
  log "Verifying the shasum digest matches the downloaded archive"
  ${shasum_cmd} -c "${sha_file}"
}

extract_archive() {
  log "Extracting ${archive}"
  case "${ext}" in
    tar.gz)
      archive_dir="${archive%.tar.gz}"
      mkdir "${archive_dir}"
      zcat "${archive}" | tar --extract --directory "${archive_dir}" --strip-components=1

      ;;
    zip)
      archive_dir="${archive%.zip}"
      unzip -j "${archive}" -d "${archive_dir}"
      ;;
    *)
      exit_with "[extract] Unknown file extension: ${ext}" 4
      ;;
  esac
}

install_cli() {
  log "Installing Lacework CLI into $installation_dir"
  mkdir -pv "$installation_dir"
  install -v "${archive_dir}/${binary_name}" "${installation_dir}/${binary_name}"
}

print_cli_version() {
  log "Verifying installed Lacework CLI version"
  LW_NOCACHE=1 LW_TELEMETRY_DISABLE=1 LW_UPDATES_DISABLE=1 "${installation_dir}/${binary_name}" version
}

download_file() {
  local _url="${1}"
  local _dst="${2}"
  local _code
  local _wget_extra_args=""
  local _curl_extra_args=""

  # try to download with wget
  if command -v wget > /dev/null; then
    log "Downloading via wget: ${_url}"

    set +e
    wget -q -O "${_dst}" "${_url}"
    _code="$?"
    set -eou pipefail

    if [ $_code -eq 0 ]; then
      return 0
    else
      warn "wget failed to download file, trying to download with curl (exitcode: ${_code})"
    fi
  fi

  # try to download with curl
  if command -v curl > /dev/null; then
    log "Downloading via curl: ${_url}"

    set +e
    curl -sSfL "${_url}" -o "${_dst}"
    _code="$?"
    set -eou pipefail

    if [ $_code -eq 0 ]; then
      return 0
    else
      warn "curl failed to download file (exitcode: ${_code})"
    fi
  fi

  # wget and curl have failed, inform the user
  exit_with "Required: SSL-enabled 'curl' or 'wget' on PATH" 6
}

log() {
  echo "--> install: $1"
}

warn() {
  echo "xxx install: $1" >&2
}

exit_with() {
  warn "$1"
  exit "${2:-10}"
}

main "$@"
