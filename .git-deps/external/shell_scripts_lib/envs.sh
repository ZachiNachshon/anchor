#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"

#######################################
# Verify if environment variable exists, fail otherwise
# Globals:
#   None
# Arguments:
#   name - environment variable name
# Usage:
#   env_var "MY_SPECIAL_SECRET"
#######################################
env_var() {
  local name=$1
  eval "value=\${$name}"
  if [[ -z "${value}" ]]; then
    log_fatal "missing env var. name: ${name}"
  fi
  echo "${value}"
}
