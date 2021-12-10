#!/usr/bin/env bash

source ./scripts/logger.sh

_read_from_file() {
  echo $(cat ./resources/version.txt)
}

_push_remote_tag() {
  tag=$1
  if [[ ! -z "${tag}" ]]; then
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
  tag="v$(_read_from_file)"
  log_info "Successfully read tag from version file"
  _push_remote_tag "${tag}"
}


#main() {
#  create_tag
#}
#
#main "$@"