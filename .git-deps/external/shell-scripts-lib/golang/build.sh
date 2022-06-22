#!/bin/bash

# Title         Build go binaries for multiple OS/Architecture
# Author        Zachi Nachshon <zachi.nachshon@gmail.com>
# Supported OS  Linux & macOS
# Description   Issue a go build for single or multiple OS/Architecture combinations,
#               defaults: linux-amd64 / darwin-amd64
#==============================================================================
CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")
ROOT_FOLDER_ABS_PATH=$(dirname "${CURRENT_FOLDER_ABS_PATH}")

# shellcheck source=../../logger.sh
source "${ROOT_FOLDER_ABS_PATH}/logger.sh"
# shellcheck source=../../checks.sh
source "${ROOT_FOLDER_ABS_PATH}/checks.sh"
# shellcheck source=../../io.sh
source "${ROOT_FOLDER_ABS_PATH}/io.sh"
# shellcheck source=../../os.sh
source "${ROOT_FOLDER_ABS_PATH}/os.sh"

# Defaults for OS/Arch
XC_ARCH=${XC_ARCH:-"amd64"}
XC_OS=${XC_OS:-linux darwin}
#XC_OS=${XC_OS:-linux darwin windows}

pre_build_info() {
  go env
  new_line
  log_info "Go version:"
  go version

  new_line
  log_info "Go files path: ${go_files_path}"
  log_info "Artifacts distribution path: ${dist_path}"
}

prepare_binary_path() {
  local binary_name=$1
  local dist_os_arch_path=$2
  local binary_file_path="${dist_os_arch_path}/${binary_name}"
  if [[ "${os}" == "windows" ]]; then
    binary_file_path="${binary_file_path}.exe"
  fi
  echo "${binary_file_path}"
}

prepare_zip_dest_path() {
  local binary_name=$1
  local dist_os_arch_path=$2
  local os=$3
  local arch=$4
  echo "${dist_os_arch_path}/${binary_name}_${os}_${arch}.zip"
}

build_os_arch() {
  pre_build_info

  for os in ${XC_OS[@]}; do
    for arch in ${XC_ARCH[@]}; do

      new_line
      log_info "Building binary for ${os}/${arch}..."

      dist_os_arch_path="${dist_path}/${os}_${arch}"

      if ! is_directory_exist "${dist_os_arch_path}"; then
        log_info "Creating output directory. path: ${dist_os_arch_path}"
        mkdir -p "${dist_os_arch_path}"
      fi

      binary_file_path=$(prepare_binary_path "${binary_name}" "${dist_os_arch_path}")
      zip_file_path=$(prepare_zip_dest_path "${binary_name}" "${dist_os_arch_path}" "${os}" "${arch}")

      if is_debug; then
        echo """
      CGO_ENABLED=0 GOARCH=${arch} GOOS=${os}
      go build -o ${binary_file_path} -a -ldflags '-extldflags "-static"' -mod=readonly "${go_files_path}"
      """
      fi

      CGO_ENABLED=0 GOARCH=${arch} GOOS=${os} go build -o ${binary_file_path} -a -ldflags '-extldflags "-static"' -mod=readonly "${go_files_path}"

      if is_file_exist "${binary_file_path}"; then
        zip "${zip_file_path}" "${binary_file_path}"
      else
        log_warning "Failed to locate binary file, cannot zip. path: ${binary_file_path}"
      fi

    done
  done
}

build() {
  local dist_path_resolved="${dist_path}"

  # GOPATH folders cannot have sub-folders with /GOPATH/bin/OS_Arch/binary_name
  if [[ "${dist_path}" != *${GOPATH}* ]]; then 
    local os_arch=$(identify_os_arch "x86_64:amd64")
    dist_path_resolved="${dist_path}/${os_arch}"
  fi

  local binary_file_path=$(prepare_binary_path "${binary_name}" "${dist_path_resolved}")

  if is_debug; then
    echo """
    go build -o "${binary_file_path}" "${go_files_path}"
"""
  fi

  go build -o "${binary_file_path}" "${go_files_path}"
  log_info "Binary created at path: ${binary_file_path}"
}

parse_program_arguments() {
  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      action*)
        action=$(cut -d : -f 2- <<<"${1}" | xargs)
        shift
        ;;
      binary_name*)
        binary_name=$(cut -d : -f 2- <<<"${1}" | xargs)
        shift
        ;;
      dist_path*)
        dist_path=$(cut -d : -f 2- <<<"${1}" | xargs)
        shift
        ;;
      go_files_path*)
        go_files_path=$(cut -d : -f 2- <<<"${1}" | xargs)
        shift
        ;;
      debug*)
        debug="verbose"
        shift
        ;;
      *)
        break
        ;;
    esac
  done

  # Set defaults
  debug=${debug=''}
}

verify_program_arguments() {
  if [[ -z "${action}" ]]; then
    log_fatal "Missing mandatory param. name: action"
  elif [[ -z "${binary_name}" ]]; then
    log_fatal "Missing mandatory param. name: binary_name"
  elif [[ -z "${dist_path}" ]]; then
    log_fatal "Missing mandatory param. name: dist_path"
  elif [[ -z "${go_files_path}" ]]; then
    log_fatal "Missing mandatory param. name: go_files_path"
  fi
}

is_debug() {
  [[ -n "${debug}" ]]
}

prerequisites() {
  check_tool go
  check_tool zip
}

#######################################
# Build go binaries from sources to a custom destination
# Globals:
#   None
# Arguments:
#   action        - build/build-os-arch
#   binary_name   - name of the binary output
#   dist_path     - binary(ies) distribution path
#   go_files_path - Go source code folder path
#   debug         - (Optional) add logs verbosity
# Usage:
#   /shell-scripts-lib/golang/build.sh \
#		  action: build \
#		  binary_name: my-project \
#		  dist_path: $(GOPATH)/bin \
#		  go_files_path: /path/to/my-project/cmd/my-project/*.go
#######################################
main() {
  parse_program_arguments "$@"
  verify_program_arguments

  prerequisites

  if [[ "${action}" == "build" ]]; then
    build
  elif [[ "${action}" == "build-os-arch" ]]; then
    build_os_arch
  fi
}

main "$@"
