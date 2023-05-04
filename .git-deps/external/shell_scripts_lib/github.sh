#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"
source "${CURRENT_FOLDER_ABS_PATH}/prompter.sh"
source "${CURRENT_FOLDER_ABS_PATH}/io.sh"
source "${CURRENT_FOLDER_ABS_PATH}/cmd.sh"

github_prompt_for_approval_before_release() {
  local tag=$1
  log_warning "Make sure to update all version releated files/variables before you continue !"
  new_line
  [[ $(prompt_yes_no "Release tag ${tag}" "warning") == "y" ]];
}

github_prompt_for_approval_before_delete() {
  local tag=$1
  log_warning "Make sure that you are deleting the correct tag before you continue !"
  new_line
  [[ $(prompt_yes_no "Delete tag ${tag}" "warning") == "y" ]];
}

github_is_release_tag_exist() {
  local tag=$1
  log_info "Checking if release tag exist. tag: ${tag}"
  cmd_run "gh release view ${tag} >/dev/null 2>&1"
}

github_is_draft_release_tag() {
  local tag=$1
  log_info "Checking if release tag is of type draft. tag: ${tag}"
  cmd_run "gh release view ${tag} --json isDraft | grep true"
}

github_create_release_tag() {
  local tag=$1
  log_info "Creating a new GitHub release. tag: ${tag}"
  cmd_run "gh release create ${tag}"
}

github_upload_release_asset() {
  local tag=$1
  local filepath=$2
  log_info "Uploading file. tag: ${tag}, path: ${filepath}"
  cmd_run "gh release upload ${tag} ${filepath}"
}

github_delete_release_tag() {
  local tag=$1
  log_info "Deleting release tag from local. name: ${tag}"
  cmd_run "git tag -d ${tag}"
  new_line
  log_info "Deleting release tag from remote. name: ${tag}"
  cmd_run "git push origin :refs/tags/${tag}"
}

github_download_release_asset() {
  local tag=$1
  local asset_name=$2
  log_info "Downloading GitHub release asset. tag: ${tag}, name: ${asset_name}"
  cmd_run "gh release download ${tag} --pattern ${asset_name}"
}

github_delete_release_tag_asset_from_any_repo() {
  local owner=$1  
  local repo=$2
  local tag_name=$3
  local asset_name=$4
  local token=$5

  local header=""
  if [[ -n "${token}" ]]; then
    header="-H \"Authorization: Bearer ${token}\""
  fi

  local curl_flags="-LJO"
  if is_verbose; then
    curl_flags="-LJOv"
  fi

  # Get the release information
  release_info=$(cmd_run "curl ${curl_flags} ${header} https://api.github.com/repos/${owner}/${repo}/releases/tags/${tag_name}")

  # Get the asset ID
  asset_id=$(echo "${release_info}" | jq ".assets[] | select(.name == \"${asset_name}\") | .id")

  if ! is_dry_run && [[ -z "${asset_id}" ]]; then
    log_fatal "Failed to retrieve asset id from GitHub release. tag: ${tag_name}, asset_name: ${asset_name}"
  fi

  # Delete the asset
  cmd_run "curl ${curl_flags} -X DELETE ${header} https://api.github.com/repos/${owner}/${repo}/releases/assets/${asset_id}"
}

github_download_release_asset_from_any_repo() {
  local owner=$1  
  local repo=$2
  local tag_name=$3
  local asset_name=$4
  local dl_path=$5
  local token=$6

  local header=""
  if [[ -n "${token}" ]]; then
    header="-H \"Authorization: Bearer ${token}\""
  fi

  cwd=$(pwd)
  if [[ -n "${dl_path}" ]] && ! is_directory_exist "${dl_path}"; then
    cmd_run "mkdir -p ${dl_path}"
  fi

  if [[ -n "${dl_path}" ]]; then
    cmd_run "cd ${dl_path} || exit"
  fi

  local curl_flags="-LJO"
  if is_verbose; then
    curl_flags="-LJOv"
  fi

  # Get the release information
  release_info=$(cmd_run "curl ${curl_flags} ${header} https://api.github.com/repos/${owner}/${repo}/releases/tags/${tag_name}")

  # Get the asset ID
  asset_id=$(echo "${release_info}" | jq ".assets[] | select(.name == \"${asset_name}\") | .id")

  if ! is_dry_run && [[ -z "${asset_id}" ]]; then
    log_fatal "Failed to retrieve asset id from GitHub release. tag: ${tag_name}, asset_name: ${asset_name}"
  fi

  # Download the asset
  cmd_run "curl ${curl_flags} ${header} -H \"Accept: application/octet-stream\" https://api.github.com/repos/${repo}/releases/assets/${asset_id}"

  if [[ -n "${dl_path}" ]]; then
    cmd_run "cd ${cwd} || exit"
  fi
}