#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"

#######################################
# Return a config value for specific key from a properties file
# Globals:
#   None
# Arguments:
#   dir_path    - config.properties file absolute path
#   key         - key to return its configuration
#   do_not_fail - (optional) do not fail if config cannot be found, default: fail
# Usage:
#   property "absolute/path/to/config/dir" "some.test.key"
#   property "absolute/path/to/config/dir" "some.test.key" "do_not_fail"
#######################################
property() {
  local dir_path=$1
  local key=$2
  local do_not_fail=$3
  local value=""

  # Consider piping ' | envsubst' at the end of the grep command
  value=$(grep "${key}" "${dir_path}/config.properties" | cut -d '=' -f2)
  if [[ -z "${value}" && "${do_not_fail}" != "do_not_fail" ]]; then
    log_fatal "missing property. key: ${key}"
  fi

  echo "${value}"
}

#######################################
# Return a formatted string pattern from a properties file
# Globals:
#   None
# Arguments:
#   dir_path - pattern.properties file absolute path
#   pattern  - string pattern identifier
#   ...      - string arguments that correspond to the %s within the pattern
# Usage:
#   pattern "absolute/path/to/pattern/dir" "some.test.pattern" "first" "second"
#######################################
pattern() {
  local dir_path=$1
  local pattern=$2

  local format=$(grep "${pattern}" "${dir_path}/pattern.properties" | cut -d '=' -f2- | envsubst)
  if [[ -z "${format}" ]]; then
    log_fatal "missing pattern. pattern: ${pattern}"
  fi

  shift
  shift
  printf ${format} "$@"
}
