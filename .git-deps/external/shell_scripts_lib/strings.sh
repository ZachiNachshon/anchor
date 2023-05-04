#!/bin/bash

# Title         String utilities for common string manipulation actions
# Author        Zachi Nachshon <zachi.nachshon@gmail.com>
# Supported OS  Linux & macOS
# Description   Use this file instead of searching how to perform string
#               manipulation all over again
#==============================================================================

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/checks.sh"

to_upper() {
  local str=$1
  echo "${str}" | tr '[:lower:]' '[:upper:]'
}

is_comment() {
  local line=$1
  [[ "${line}" == \#* ]]
}

#######################################
# Split a string by a delimiter
# Globals:
#   None
# Arguments:
#   str        - string to manipulate
#   delimiter  - (optional) delimiter to use
# Usage:
#   split_newlines_by_delimiter "one two three"
#   split_newlines_by_delimiter "one;two;three" ";"
#######################################
split_newlines_by_delimiter() {
  local str=$1
  local delimiter=$2

  # By default use space as delimiter
  if [[ -z "${delimiter}" ]]; then
    delimiter=" "
  fi

  echo "${str}" | tr "${delimiter}" '\n'
}

#######################################
# Find difference between two strings
# Globals:
#   None
# Arguments:
#   str_a    - 1st string
#   str_b    - 2nd string
# Usage:
#   is_text_equal "hello" "hella"
#######################################
is_text_equal() {
  local str_a=$1
  local str_b=$2

  check_tool diff

  # Note: Those diff options won't apply on lines that start with a whitespace char 

  #  -E, --ignore-tab-expansion
  #         ignore changes due to tab expansion

  #  -b, --ignore-space-change
  #         ignore changes in the amount of white space

  #  -w, --ignore-all-space
  #         ignore all white space

  #  -B, --ignore-blank-lines
  #         ignore changes whose lines are all blank

  # Using sed to remove any color codes from text (example: \033[0;36m)
  # diff -wEbB -y --suppress-common-lines \
  diff -wbB -y --suppress-common-lines \
  <(echo -e "${str_a}" | sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g") \
  <(echo -e "${str_b}" | sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g")
}

# is_text_equal() {
#   local str_a=$1
#   local str_b=$2

#   check_tool comm
#   # -1     suppress column 1 (lines unique to FILE1)
#   # -2     suppress column 2 (lines unique to FILE2)
#   # -3     suppress column 3 (lines that appear in both files)
#   local formatted_str_a=$(echo -e "${str_a}" | xargs | sort)
#   local formatted_str_b=$(echo -e "${str_b}" | xargs | sort)
#   result=$(comm -3 <(echo ${formatted_str_a}) <(echo ${formatted_str_b}))
#   [[ -z "${result}" ]]
# }