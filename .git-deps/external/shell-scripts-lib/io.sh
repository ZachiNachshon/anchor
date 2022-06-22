#!/bin/bash

is_file_exist() {
  local path=$1
  [[ -f "${path}" || $(is_symlink "${path}") ]]
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

is_directory_exist() {
  local path=$1
  [[ -d "${path}" ]]
}

is_file_contain() {
  local filepath=$1
  local text=$2
  grep -q -w "${text}" "${filepath}"
}
