#!/bin/bash

# Title         Format or check go code for gofmt standards
# Author        Zachi Nachshon <zachi.nachshon@gmail.com>
# Supported OS  Linux & macOS
# Description   Run gofmt to format/check to comply with gofmt standards
#==============================================================================
CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")
ROOT_FOLDER_ABS_PATH=$(dirname "${CURRENT_FOLDER_ABS_PATH}")

# shellcheck source=../../logger.sh
source "${ROOT_FOLDER_ABS_PATH}/logger.sh"
# shellcheck source=../../checks.sh
source "${ROOT_FOLDER_ABS_PATH}/checks.sh"

prerequisites() {
  check_tool go
  check_tool zip
}

format() {
  log_info "Formatting Go source code to comply with gofmt requirements..."
  gofmt_files=$(gofmt -l $(find . -name '*.go' | grep -v vendor))
  if [[ -n ${gofmt_files} ]]; then
    log_info "About to format the following files:"
    echo "${gofmt_files}"
    gofmt -w $(find . -name '*.go' | grep -v vendor)
    new_line
  fi
  log_info "Done."
}

check() {
  log_info "Checking that Go code complies with gofmt requirements..."
  gofmt_files=$(gofmt -l $(find . -name '*.go' | grep -v vendor))
  if [[ -n ${gofmt_files} ]]; then
    log_warning 'gofmt needs to run on the following files:'
    log_warning "${gofmt_files}"
    log_warning "Please run \`make fmt\` to reformat the code."
    exit 1
  fi

  log_info "Go code is well formatted !"
  exit 0
}

#######################################
# Format or check go source code to gofmt standards
# Scan path is working directory (PWD)
# Globals:
#   None
# Arguments:
#   action - --check or --format
# Usage:
#   /shell-scripts-lib/golang/fmt.sh --check
#######################################
main() {
  local action=$1
  prerequisites

  if [[ "${action}" == "--check" ]]; then
    check
  elif [[ "${action}" == "--format" ]]; then
    format
  else
    log_fatal "Invalid action flag, supported flags: --check, --format."
  fi
}

main "$@"
