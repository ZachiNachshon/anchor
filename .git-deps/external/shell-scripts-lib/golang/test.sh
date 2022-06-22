#!/bin/bash

# Title         Run tests with coverage report on multiple environments
# Author        Zachi Nachshon <zachi.nachshon@gmail.com>
# Supported OS  Linux & macOS
# Description   Run tests suite with coverage on environments: local/containerized/CI
#==============================================================================
CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")
ROOT_FOLDER_ABS_PATH=$(dirname "${CURRENT_FOLDER_ABS_PATH}")

# shellcheck source=../../logger.sh
source "${ROOT_FOLDER_ABS_PATH}/logger.sh"
# shellcheck source=../../checks.sh
source "${ROOT_FOLDER_ABS_PATH}/checks.sh"

test_local() {
  local tests_path=$1
  go test -v ${tests_path} -json -cover -covermode=count -coverprofile=coverage.out.temp | tparse -all -notests -smallscreen
  cat coverage.out.temp | grep -v '_testkit\|_fakes' >coverage.out
  go tool cover -func coverage.out | grep total | awk '{print $3}'
}

test_containerized() {
  if [[ -z "${go_version}" ]]; then
    log_fatal "Missing mandatory param. name: go_version"
  elif [[ -z "${binary_name}" ]]; then
    log_fatal "Missing mandatory param. name: binary_name"
  elif [[ -z "${project_root_path}" ]]; then
    log_fatal "Missing mandatory param. name: project_root_path"
  fi

  docker run -it \
    -v ${project_root_path}:/home/${binary_name} \
    --entrypoint /bin/sh golang:${go_version} \
    -c """
go install github.com/mfridman/tparse@v0.10.0
cd /home/${binary_name}
go test -v ${tests_path} -json -cover -covermode=count -coverprofile=coverage.out.temp | tparse -all -notests -smallscreen
cat coverage.out.temp | grep -v '_testkit\|_fakes' > coverage.out
go tool cover -func coverage.out | grep total | awk '{print $3}'
"""
}

test_github_ci() {
  local tests_path=$1
  go test -v ${tests_path} -json -cover -covermode=count -coverprofile=coverage.out.temp | tparse -all
  cat coverage.out.temp | grep -v '_testkit\|_fakes' >coverage.out
  # -coverprofile=coverage.out was added for GitHub workflow integration with jandelgado/gcov2lcov-action
  # Error:
  #   /tmp/gcov2lcov-linux-amd64 -infile coverage.out -outfile coverage.lcov
  #   2021/08/01 07:21:57 error opening input file: open coverage.out: no such file or directory
}

prerequisites() {
  check_tool go
  check_tool zip
}

parse_program_arguments() {
  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      action*)
        action=$(cut -d : -f 2- <<<"${1}" | xargs)
        shift
        ;;
      tests_path*)
        tests_path=$(cut -d : -f 2- <<<"${1}" | xargs)
        shift
        ;;
      go_version*)
        go_version=$(cut -d : -f 2- <<<"${1}" | xargs)
        shift
        ;;
      binary_name*)
        binary_name=$(cut -d : -f 2- <<<"${1}" | xargs)
        shift
        ;;
      project_root_path*)
        project_root_path=$(cut -d : -f 2- <<<"${1}" | xargs)
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
  elif [[ -z "${tests_path}" ]]; then
    log_fatal "Missing mandatory param. name: tests_path"
  fi
}

#######################################
# Run tests with coverage report on multiple environments
# Globals:
#   None
# Arguments:
#   action            - local/containerized/github-ci
#   tests_path        - test source files path
#   go_version        - (containerized) go version to run tests with
#   binary_name       - (containerized) project binary name used for volume mapping
#   project_root_path - (containerized) project root folder used for volume mapping
#   debug             - (Optional) add logs verbosity
# Usage:
#   /shell-scripts-lib/golang/test.sh \
#		  action: local \
#		  tests_path: ./...

#   /shell-scripts-lib/golang/test.sh \
#		  action: containerized \
#		  tests_path: ./... \
#		  project_root_path: $PWD \
#		  binary_name: my-project \
#		  go_version: 1.18
#######################################
main() {
  parse_program_arguments "$@"
  verify_program_arguments

  prerequisites

  log_info "Running tests on environment: ${action}..."

  if [[ "${action}" == "local" ]]; then
    check_tool tparse
    test_local "${tests_path}"

  elif [[ "${action}" == "containerized" ]]; then
    check_tool docker
    test_containerized "${tests_path}"

  elif [[ "${action}" == "github-ci" ]]; then
    check_tool tparse
    test_github_ci "${tests_path}"

  else
    log_fatal "Invalid action flag, supported flags: --local, --containerized, --ci."
  fi
}

main "$@"
