#!/bin/bash

COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW="\033[0;33m"
COLOR_BLUE="\033[0;34m"
COLOR_PURPLE="\033[0;35m"
COLOR_LIGHT_CYAN='\033[0;36m'
COLOR_WHITE='\033[1;37m'
COLOR_NONE='\033[0m'

LOGGER_VERBOSE=""
LOGGER_SILENT=""
LOGGER_DRY_RUN=""

ICON_GOOD="${COLOR_GREEN}✔${COLOR_NONE}"
ICON_WARN="${COLOR_YELLOW}⚠${COLOR_NONE}"
ICON_BAD="${COLOR_RED}✗${COLOR_NONE}"

exit_on_error() {
  exit_code=$1
  message=$2
  if [ $exit_code -ne 0 ]; then
    #        >&1 echo "\"${message}\" command failed with exit code ${exit_code}."
    # >&1 echo "\"${message}\""
    exit $exit_code
  fi
}

is_verbose() {
  [[ -n "${LOGGER_VERBOSE}" ]]
}

is_silent() {
  [[ -n ${LOGGER_SILENT} ]]
}

is_dry_run() {
  [[ -n ${LOGGER_DRY_RUN} ]]
}

evaluate_dry_run_mode() {
  if is_dry_run; then
    echo -e "${COLOR_YELLOW}Running in DRY RUN mode${COLOR_NONE}" >&0
    new_line
  fi
}

_log_base() {
  prefix=$1
  shift
  echo -e "${prefix}$*" >&0
}

log_debug() {
  local debug_level_txt="DEBUG"
  if is_dry_run; then
    debug_level_txt+=" (Dry Run)"
  fi

  if ! is_silent && is_verbose; then
    _log_base "${COLOR_WHITE}${debug_level_txt}${COLOR_NONE}: " "$@"
  fi
}

log_info() {
  local info_level_txt="INFO"
  if is_dry_run; then
    info_level_txt+=" (Dry Run)"
  fi

  if ! is_silent; then
    _log_base "${COLOR_GREEN}${info_level_txt}${COLOR_NONE}: " "$@"
  fi
}

log_warning() {
  local warn_level_txt="WARNING"
  if is_dry_run; then
    warn_level_txt+=" (Dry Run)"
  fi

  if ! is_silent; then
    _log_base "${COLOR_YELLOW}${warn_level_txt}${COLOR_NONE}: " "$@"
  fi
}

log_error() {
  local error_level_txt="ERROR"
  if is_dry_run; then
    error_level_txt+=" (Dry Run)"
  fi
  _log_base "${COLOR_RED}${error_level_txt}${COLOR_NONE}: " "$@"
}

log_fatal() {
  local fatal_level_txt="ERROR"
  if is_dry_run; then
    fatal_level_txt+=" (Dry Run)"
  fi
  _log_base "${COLOR_RED}${fatal_level_txt}${COLOR_NONE}: " "$@"
  message="$@"
  exit_on_error 1 "${message}"
}

new_line() {
  echo -e "" >&1
}

log_indicator_good() {
  local error_level_txt=""
  if is_dry_run; then
    error_level_txt+=" (Dry Run)"
  fi
  if ! is_silent; then
    _log_base "${ICON_GOOD}${error_level_txt} " "$@"
  fi
}

log_indicator_warning() {
  local error_level_txt=""
  if is_dry_run; then
    error_level_txt+=" (Dry Run)"
  fi
  if ! is_silent; then
    _log_base "${ICON_WARN}${error_level_txt} " "$@"
  fi
}

log_indicator_bad() {
  local error_level_txt=""
  if is_dry_run; then
    error_level_txt+=" (Dry Run)"
  fi
  if ! is_silent; then
    _log_base "${ICON_BAD}${error_level_txt} " "$@"
  fi
}
