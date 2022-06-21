#!/bin/bash

COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW="\033[0;33m"
COLOR_WHITE='\033[1;37m'
COLOR_LIGHT_CYAN='\033[0;36m'
COLOR_NONE='\033[0m'

CLI_OPTION_SILENT=""

exit_on_error() {
  exit_code=$1
  message=$2
  if [ $exit_code -ne 0 ]; then
    #        >&2 echo "\"${message}\" command failed with exit code ${exit_code}."
    # >&2 echo "\"${message}\""
    exit $exit_code
  fi
}

is_silent() {
  [[ ! -z ${CLI_OPTION_SILENT} ]]
}

_log_base() {
  prefix=$1
  shift
  echo -e "${prefix}$*" >&2
}

log_info() {
  if ! is_silent; then
    _log_base "${COLOR_GREEN}INFO${COLOR_NONE}: " "$@"
  fi
}

log_warning() {
  if ! is_silent; then
    _log_base "${COLOR_YELLOW}WARNING${COLOR_NONE}: " "$@"
  fi
}

log_error() {
  _log_base "${COLOR_RED}ERROR${COLOR_NONE}: " "$@"
}

log_fatal() {
  _log_base "${COLOR_RED}ERROR${COLOR_NONE}: " "$@"
  message="$@"
  exit_on_error 1 "${message}"
}

new_line() {
  echo -e "" >&2
}
