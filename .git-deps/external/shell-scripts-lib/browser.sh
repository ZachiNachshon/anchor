#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"
source "${CURRENT_FOLDER_ABS_PATH}/os.sh"

#######################################
# Open a chrome/default browser under supplied address
# Globals:
#   None
# Arguments:
#   address - url address to open the browser with
# Usage:
#   open_browser "http://www.google.com"
#######################################
open_browser() {
  local address=$1

  local os_type=$(read_os_type)
  log_info "Opening browser. os: ${os_type}, address: ${address}"
  if [[ "${os_type}" == "linux" ]]; then
    xdg-open "${address}"
  elif [[ "${os_type}" == "darwin" ]]; then
    open -a "Google Chrome" "${address}"
  fi
}
