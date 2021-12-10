#!/bin/bash

COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW="\033[0;33m"
COLOR_WHITE='\033[1;37m'
COLOR_NONE='\033[0m'

_log_base() {
  prefix=$1
  shift
  echo -e "${prefix}$*" >&2
}

log_info() {
  _log_base "${COLOR_GREEN}INFO${COLOR_NONE}: " "$@"
}

log_warning() {
  _log_base "${COLOR_YELLOW}WARNING${COLOR_NONE}: " "$@"
}

log_error() {
  _log_base "${COLOR_RED}ERROR${COLOR_NONE}: " "$@"
}

log_fatal() {
  _log_base "${COLOR_RED}ERROR${COLOR_NONE}: " "$@"
  exit 1
}

new_line() {
  echo -e "" >&2
}