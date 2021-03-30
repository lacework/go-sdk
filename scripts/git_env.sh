#!/bin/bash
#
# Name::        git_env.sh
# Description:: Use this script to configure local git env
# Author::      Darren Murray (<dmurray-lacework@lacework.net>)
#

readonly commit_hook=scripts/githooks/prepare-commit-msg
readonly hooks_path=scripts/githooks/
readonly project_name=go-sdk

prepare_git_env(){
  chmod +x $commit_hook
  log "Setting git-hooks path to $hooks_path"
  git config core.hooksPath $hooks_path
}

log() {
  echo "--> ${project_name}: $1"
}

prepare_git_env