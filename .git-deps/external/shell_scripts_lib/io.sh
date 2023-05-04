#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/cmd.sh"

is_file_exist() {
  local path=$1
  [[ -f "${path}" || $(is_symlink "${path}") ]]
}

is_file_extension() {
  local filename=$1
  local extension=$2
  [[ "${filename}" == *".${extension}"* ]]
}

is_file_size_bigger_than_zero() {
  local path=$1
  [[ -e "${path}" ]]
}

is_file_has_name() {
  local path=$1
  local name=$2
  [[ "${path}" == *"${name}"* ]]
}

is_symlink() {
  local abs_path=$1
  [[ -L "${abs_path}" ]]
}

is_symlink_target() {
  local symlink=$1
  local target=$2
  local link_dest=$(readlink "${symlink}")
  local result="${target}"
  [[ "${link_dest}" != "${target}" ]]
}

is_directory_exist() {
  local path=$1
  [[ -d "${path}" ]]
}

is_directory_empty() {
  local path=$1
  local result=""
  if is_directory_exist "${path}"; then
    result=$(ls -A "${path}")
  fi

  [[ -z "${result}" ]]
}

is_file_contain() {
  local filepath=$1
  local text=$2
  grep -q -w "${text}" "${filepath}"
}

create_symlink() {
  local link_path=$1
  local destination=$2
  cmd_run "ln -sfn "${destination}" "${link_path}""
}

remove_symlink() {
  local link_path=$1
  cmd_run "unlink "${link_path}" 2>/dev/null"
}
