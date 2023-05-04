#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/io.sh"

#######################################
# Return an external git dependency folder path
# or just a folder path if it is not an external
# dependency.
# 
# Globals:
#   None
# 
# Arguments:
#   working_dir      - ABS of working directory path
#   dependency_name  - (Optional) external dependency name
#   inner_path       - Path from the working directory / external dependency folder
# 
# Usage:
#   get_external_folder_dependency_path \
#     "$PWD" \
#     "shell_scripts_lib" \
#     "runner/ansible/config"
#######################################
get_external_folder_dependency_path() {
  local working_dir=$1
  local dependency_name=$2
  local inner_path=$3

  if is_directory_exist "${working_dir}/external"; then
    echo "${working_dir}/external/${dependency_name}/${inner_path}"
  else
    echo "${working_dir}/${inner_path}"
  fi
}
