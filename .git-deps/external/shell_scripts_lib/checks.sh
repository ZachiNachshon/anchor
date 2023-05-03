#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"

#######################################
# Checks if local utility exists
# Globals:
#   None
# Arguments:
#   name - utility CLI name
# Usage:
#   is_tool_exist "kubectl"
#######################################
is_tool_exist() {
  local name=$1
  [[ $(command -v "${name}") ]]
}

#######################################
# Verify if local utility exists, fail otherwise
# Globals:
#   None
# Arguments:
#   name - utility CLI name
# Usage:
#   check_tool "kubectl"
#######################################
check_tool() {
  local name=$1
  local exists=$(command -v "${name}")
  if [[ "${exists}" != *${name}* && "${exists}" != 0 ]]; then
    log_fatal "missing utility. name: ${name}"
  fi
}

#######################################
# Checks local Docker image exists
# Globals:
#   None
# Arguments:
#   name - Docker image name
# Usage:
#   is_image_exists "busybox"
#######################################
is_image_exists() {
  local name=$1
  local result=$(docker images | grep "${name}" | awk {'print $3'})
  [[ -n ${result} ]]
}

#######################################
# Verify if local Docker image exists, fail otherwise
# Globals:
#   None
# Arguments:
#   image_name - Docker image name
#   message    - (optional) failure message
# Usage:
#   check_image "busybox"
#######################################
check_image() {
  local image_name=$1
  local message=$2

  local exists=$(docker images -a | grep "${image_name}" | awk {'print $3'})
  if [[ -z "${exists}" ]]; then
    if [[ -z "${message}" ]]; then
      log_fatal "No docker image could be found. name: ${image_name}"
    else
      log_fatal "No docker image could be found. name: ${image_name}, message: ${message}"
    fi
  fi
}
