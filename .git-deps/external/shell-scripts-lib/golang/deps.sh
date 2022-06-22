#!/bin/bash

# Title         Align go dependencies
# Author        Zachi Nachshon <zachi.nachshon@gmail.com>
# Supported OS  Linux & macOS
# Description   Add/remove/clean go dependencies
#==============================================================================
CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")
ROOT_FOLDER_ABS_PATH=$(dirname "${CURRENT_FOLDER_ABS_PATH}")

# shellcheck source=../../logger.sh
source "${ROOT_FOLDER_ABS_PATH}/logger.sh"

prerequisites() {
  check_tool go
}

#######################################
# Tidy, verify and vendor go dependencies for existing project
# Globals:
#   None
# Arguments:
#   None
# Usage:
#   /shell-scripts-lib/golang/deps.sh
#######################################
main() {
  log_info "Tidying dependencies..."
  go mod tidy
  log_info "Verifying dependencies..."
  go mod verify
  log_info "Vendoring dependencies..."
  go mod vendor
  log_info "Done."
}

main "$@"
