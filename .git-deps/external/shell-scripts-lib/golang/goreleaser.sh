#!/bin/bash

# Title         Create a GitHub for a specific project version
# Author        Zachi Nachshon <zachi.nachshon@gmail.com>
# Supported OS  Linux & macOS
# Description   Run goreleaser to release a new version for the project
#==============================================================================
CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")
ROOT_FOLDER_ABS_PATH=$(dirname "${CURRENT_FOLDER_ABS_PATH}")

# shellcheck source=../../logger.sh
source "${ROOT_FOLDER_ABS_PATH}/logger.sh"
# shellcheck source=../../checks.sh
source "${ROOT_FOLDER_ABS_PATH}/checks.sh"
# shellcheck source=../../prompter.sh
source "${ROOT_FOLDER_ABS_PATH}/prompter.sh"
# shellcheck source=../../io.sh
source "${ROOT_FOLDER_ABS_PATH}/io.sh"
# shellcheck source=../../envs.sh
source "${ROOT_FOLDER_ABS_PATH}/envs.sh"

read_version_from_file() {
  echo $(cat ${version_file_path})
}

push_remote_tag() {
  local tag=$1

  if [[ -n "${tag}" ]]; then
    log_info "Creating GitHub tag: ${tag}"

    if git tag ${tag}; then
      git push origin tag ${tag}
    else
      log_error """Tag already exists locally, please remove local/remote tags:

  • Local tag: git tag -d ${tag}
  • Remote tag: git push origin :refs/tags/${tag}
"""
      tag=""
    fi
  else
    log_error "Tag cannot be empty, aborting"
  fi

  echo ${tag}
}

# Delete local tag: git tag -d <tag-name>
# Delete remote tag: git push --delete origin <tag-name>
create_tag_based_on_version_file() {
  local tag="v$(read_version_from_file)"
  log_info "Successfully read tag from version file. path: ${version_file_path}"
  push_remote_tag "${tag}"
}

release_version() {
  local tag=$1

  log_info "Releasing binaries... (tag: ${tag})"
  if is_debug; then
    echo """
  goreleaser release --rm-dist --config=${config_file_path}
  """
  fi

  #    goreleaser release --rm-dist --config=${config_file_path}
}

parse_program_arguments() {
  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      config_file_path*)
        config_file_path=$(cut -d : -f 2- <<<"${1}" | xargs)
        shift
        ;;
      version_file_path*)
        version_file_path=$(cut -d : -f 2- <<<"${1}" | xargs)
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
  if [[ -z "${config_file_path}" ]]; then
    log_fatal "Missing mandatory param. name: config_file_path"
  elif [[ -z "${version_file_path}" ]]; then
    log_fatal "Missing mandatory param. name: version_file_path"
  fi

  if ! is_file_exist "${config_file_path}"; then
    log_fatal "Failed to locate go-releaser config file. path: ${config_file_path}"
  fi

  if ! is_file_exist "${version_file_path}"; then
    log_fatal "Failed to locate version file. path: ${version_file_path}"
  fi
}

is_debug() {
  [[ -n "${debug}" ]]
}

prerequisites() {
  check_tool goreleaser
  # Verify token exists, do not print to stdout
  # goreleaser is using this token as an authentication method
  env_var GITHUB_TOKEN >>/dev/null
}

#######################################
# Run tests with coverage report on multiple environments
# Globals:
#   GITHUB_TOKEN       - valid GitHub token with write access for tags creation
# Arguments:
#   config_file_path   - path to goreleaser YAML config file
#   version_file_path  - path to read the project version from
# Usage:
#   /shell-scripts-lib/golang/test.sh \
#		  config_file_path: ./.goreleaser.yml \
#		  version_file_path: ./resources/version.txt
#######################################
main() {
  parse_program_arguments "$@"
  verify_program_arguments

  prerequisites

  local tag=$(create_tag_based_on_version_file)

  if [[ -z "${tag}" ]]; then
    log_fatal "Cannot release due to invalid tag. tag: ${tag}"
  fi

  new_line
  local yn=$(prompt_yes_no "Release version ${tag}")
  if [[ "${yn}" == "y" ]]; then
    release_version ${github_token} ${tag}
  else
    log_info "Nothing was released."
}

main "$@"
