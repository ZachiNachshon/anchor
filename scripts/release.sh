#!/usr/bin/env bash

source ./scripts/logger.sh
source ./scripts/tag.sh

main() {
  if [[ -z "${GITHUB_TOKEN}" ]]; then
    log_fatal "missing env var GITHUB_TOKEN"
  fi

  tag=$(create_tag_based_on_version_file)

  if [[ -z "${tag}" ]]; then
    log_fatal "Cannot release due to invalid tag"
  fi

  new_line
  printf "Release tag ${tag} (y/n): " yn
  read yn
  new_line

  if [[ "${yn}" == "y" ]]; then
    log_info "Releasing binaries... (tag: ${tag})"
    goreleaser release --rm-dist
  else
    log_info "Nothing was released."
  fi
}

main "$@"