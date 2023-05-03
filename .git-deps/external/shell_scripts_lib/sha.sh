#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"
source "${CURRENT_FOLDER_ABS_PATH}/io.sh"

shasum_calculate() {
  local url=$1
  local filename=$(basename "${url}")

  local download_path=$(mktemp -d "${TMPDIR:-/tmp}"/shell-scripts-lib-shasum.XXXXXX)
  cwd=$(pwd)
  cd "${download_path}" || exit
  curl -LJO --fail --silent --show-error "${url}" --output "${filename}"

  local shasum=""
  if is_file_exist "${download_path}/${filename}"; then
    shasum=$(shasum -a 256 "${download_path}/${filename}")
  else
    log_error "Invalid url. path: ${url}"
  fi

  cd "${cwd}" || exit
  echo "${shasum}"
}
