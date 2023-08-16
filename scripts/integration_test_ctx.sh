#!/bin/bash
#
# Name::        integration_test_ctx.sh
# Description:: Use this script to run integration tests based on changed files
# Author::      Darren Murray (<dmurray-lacework@lacework.net>)
#
readonly project_name=go-sdk

run_integration_tests(){
  BRANCH=$(git branch --show-current)
  CHANGES=$(git --no-pager diff --name-only $BRANCH $(git merge-base $BRANCH main))

  log "Changes -> ${CHANGES}"

  # Fetch relevant build tags
  TAGS=$(go run integration/context/ctx_cfg.go -- $CHANGES)

  # if no tags then exit
  if [ "$TAGS" = "" ]; then
      log "No tags for changes: ${CHANGES}"
      exit 0
  fi

  # Run integration tests matching build tags
  log "Running tests with tags ${TAGS}"
  PATH="${PWD}/bin:${PATH}" gotestsum -- -v github.com/lacework/go-sdk/integration -timeout 30m -tags="${TAGS}"
}

log() {
  echo "--> ${project_name}: $1"
}

run_integration_tests