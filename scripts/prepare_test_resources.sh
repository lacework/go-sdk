#!/bin/bash
#
# Name::        prepare_test_resources.sh
# Description:: Make ready resources required by integration tests
# Author::      Darren Murray (<darren.murray@lacework.net>)
#

main() {
  case "${1:-}" in
  clean)
    build_clean
    ;;
  dirty)
    build_dirty
    ;;
  go_component)
    go_component
    ;;
  all)
    build_clean
    build_dirty
    go_component
    ;;
  *)
    echo "invalid argument"
    ;;
esac
}

go_component() {
  echo "building integration/test_resources/cdk/go-component"
	gox -output="integration/test_resources/cdk/go-component/bin/go-component-{{.OS}}-{{.Arch}}" \
            -os="linux windows" \
            -arch="amd64 386" \
            -osarch="darwin/amd64 darwin/arm64 linux/arm linux/arm64" \
            github.com/lacework/go-sdk/integration/test_resources/cdk/go-component
}


build_clean() {
  echo "$DOCKERHUB_PASS" | docker login -u "$DOCKERHUB_USERNAME" --password-stdin
  echo "building clean container"
  docker build --no-cache -f "integration/test_resources/clean.Dockerfile" -t techallylw/test-cli-clean .
  docker push techallylw/test-cli-clean
}

build_dirty() {
  echo "$DOCKERHUB_PASS" | docker login -u "$DOCKERHUB_USERNAME" --password-stdin
  echo "building dirty container"
  docker build -f "integration/test_resources/vuln_scan/dirty.Dockerfile" -t techallylw/test-cli-dirty .
  docker push techallylw/test-cli-dirty
}

main "$@" || exit 99
