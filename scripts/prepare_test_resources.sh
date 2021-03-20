#!/bin/bash
#
# Name::        prepare_test_resources.sh
# Description:: Make ready resources required by integration tests
# Author::      Darren Murray (<darren.murray@lacework.net>)
#

main() {
  if [[ -z $DOCKERHUB_PASS ]]; then
    echo "$DOCKERHUB_PASS" | docker login -u "$DOCKERHUB_USERNAME" --password-stdin
  fi

  case "${1:-}" in
  clean)
    build_clean
    ;;
  dirty)
    build_dirty
    ;;
  all)
    build_clean
    build_dirty
    ;;
  *)
    echo "invalid argument"
    ;;
esac
}

build_clean() {
  echo "building clean container"
  docker build --no-cache -f "integration/test_resources/clean.Dockerfile" -t techallylw/test-cli-clean .
  docker push techallylw/test-cli-clean
}

build_dirty() {
  echo "building dirty container"
  docker build -f "integration/test_resources/vuln_scan/dirty.Dockerfile" -t techallylw/test-cli-dirty .
  docker push techallylw/test-cli-dirty
}

main "$@" || exit 99
