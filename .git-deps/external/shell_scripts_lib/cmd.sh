#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"

#######################################
# Run a command from string
# Globals:
#   is_verbose - based on env var LOGGER_VERBOSE
#   is_dry_run - based on env var LOGGER_DRY_RUN
# Arguments:
#   cmd_string - shell command in string format
# Usage:
#   cmd_run "echo 'hello world'"
#######################################
cmd_run() {
  local cmd_string=$1
  if is_verbose; then
    echo """
  ${cmd_string}
""" >&1
  fi
  if ! is_dry_run; then
    eval "${cmd_string}"
  fi
}
