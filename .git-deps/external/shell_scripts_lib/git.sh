#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"
source "${CURRENT_FOLDER_ABS_PATH}/cmd.sh"

git_get_current_commit_hash() {
  local clone_path=$1
  local branch=$2
  cmd_run "git -C ${clone_path} rev-parse ${branch}"
}

git_get_branch_name() {
  cmd_run "git rev-parse --abbrev-ref HEAD"
}
