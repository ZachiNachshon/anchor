#!/bin/bash

#######################################
# Get the active shell name
# Globals:
#   None
# Arguments:
#   None
# Usage:
#   shell_get_name
#######################################
shell_get_name() {
  # Get the shell suffix i.e. zsh from /bin/zsh
  # echo "${SHELL##*/}"
  echo "${SHELL}"
}

shell_is_zsh() {
  [[ ${SHELL} == *"zsh"* ]]
}

shell_is_bash() {
  [[ ${SHELL} == *"bash"* ]]
}

get_rc_file_path() {
  echo ""
}
