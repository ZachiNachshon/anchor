#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"

github_release_create_tag() {
  local tag=$1
  check_tool "gh"
  log_info "Creating a new GitHub release. tag: ${tag}"
  gh release create "${tag}"
}

github_release_upload_file() {
  local tag=$1
  local filepath=$2
  check_tool "gh"
  log_info "Uploading file. tag: ${tag}, path: ${filepath}"
  gh release upload "${tag}" "${filepath}"
}

github_release_delete_tag() {
  local tag=$1
  check_tool "gh"
  log_info "Deleting release tag from local. name: ${tag}"
  git tag -d "${tag}"
  new_line
  log_info "Deleting release tag from remote. name: ${tag}"
  git push origin ":refs/tags/${tag}"
}
