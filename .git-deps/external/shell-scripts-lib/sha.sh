#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"

shasum_calculate() {
  local url=$1
  local filename=$(basename "${url}")

  local download_path=$(mktemp -d "${TMPDIR:-/tmp}"/shell-scripts-lib-shasum.XXXXXX)
  cwd=$(pwd)
  cd "${download_path}" || exit
  curl -s "${url}" \
    -L -o "${filename}"

  log_info "SHA 256:"
  shasum -a 256 "${download_path}/${filename}"
  cd "${cwd}" || exit
}
