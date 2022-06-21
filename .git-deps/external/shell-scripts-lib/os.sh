#!/bin/bash

#######################################
# Return OS type as plain string
# Globals:
#   OSTYPE
# Arguments:
#   None
# Usage:
#   read_os_type
#######################################
read_os_type() {
  if [[ "${OSTYPE}" == "linux"* ]]; then
    echo "linux"
  elif [[ "${OSTYPE}" == "darwin"* ]]; then
    echo "darwin"
  else
    echo "OS type is not supported. os: ${OSTYPE}"
  fi
}

#######################################
# Return OS_Arch tuple as plain string
# Allow overriding arch with custom name
# Globals:
#   None
# Arguments:
#    string - (optional) custom mapping for arch e.g "x86_64:amd64"
# Usage:
#   identify_os_arch
#   identify_os_arch "x86_64:amd64" "armv:arm"
#######################################
identify_os_arch() {
  local amd64="amd64"
  local arm="arm"
  local arm64="arm64"
  local i386="386"
  local override_arch=$(if [[ "$#" -gt 0 ]]; then echo "true"; else echo "false"; fi)

  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      x86_64*)
        amd64=$(cut -d : -f 2- <<<"${1}")
        shift
        ;;
      i386*)
        i386=$(cut -d : -f 2- <<<"${1}")
        shift
        ;;
      armv*)
        arm=$(cut -d : -f 2- <<<"${1}")
        shift
        ;;
      arm64*)
        arm64=$(cut -d : -f 2- <<<"${1}")
        shift
        ;;
    esac
  done

  local os=$(uname | tr '[:upper:]' '[:lower:]')
  local arch=$(uname -m | tr '[:upper:]' '[:lower:]')
  local result="${os}_${arch}"

  # Replace arch with custom mapping, if supplied
  if [[ "${override_arch}" == "true" ]]; then
    case "${arch}" in
      x86_64*)
        result="${os}_${amd64}"
        ;;
      386*)
        result="${os}_${i386}"
        ;;
      armv*)
        result="${os}_${arm}"
        ;;
      arm64*)
        result="${os}_${arm64}"
        ;;
    esac
  fi

  echo "${result}"
}
