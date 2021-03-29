#!/bin/bash
#
# Name::        git_env.sh
# Description:: Use this script to configure local git env,
#               including commit message validation
# Author::      Darren Murray (<dmurray-lacework@lacework.net>)
#

readonly commit_hook=.git/hooks/prepare-commit-msg
readonly new_commit_hook=scripts/githooks/prepare-commit-msg
readonly project_name=go-sdk

create_prepare_commit_msg(){
  if [ ! -f $commit_hook ]; then
    log "Adding git commit hooks"
    cp scripts/githooks/prepare-commit-msg $commit_hook
    chmod +x $commit_hook
  else 
    log "Git commit hook already exists"
    update
  fi
}

update() {
  currentVersion=$(cat $commit_hook | grep 'readonly version=v'| sed 's/.*readonly version=\(.*\)/\1/')
  newVersion=$(cat $new_commit_hook | grep 'readonly version=v'| sed 's/.*readonly version=\(.*\)/\1/')
    if [ $newVersion != $currentVersion ]; then
        log "Updating git commit hooks version $currentVersion -> $newVersion"
        cp scripts/githooks/prepare-commit-msg $commit_hook
        chmod +x $commit_hook
    else 
        log "No changes made to version"
  fi
}

log() {
  echo "--> ${project_name}: $1"
}

create_prepare_commit_msg