#!/usr/bin/env bash

source ./scripts/logger.sh

_prompt_for_tag() {
  read -p "Enter tag (v0.0.0): v" input
  echo -e ${input}
}

_push_remote_tag() {
  tag=$1
  if [[ ! -z "${tag}" ]]; then
    log_info "Creating GitHub tag: ${tag}"
    if git tag ${tag}; then
      git push origin tag ${tag}
    else
      log_error "Tag already exist locally, please remove (git tag -d ${tag})"
      tag=""
    fi
  else
    log_error "Tag cannot be empty, aborting"
  fi

  echo ${tag}
}

# Delete local tag: git tag -d <tag-name>
# Delete remote tag: git push --delete origin <tag-name>
create_tag() {
  tag="v$(_prompt_for_tag)"
  _push_remote_tag "${tag}"
}


#main() {
#  create_tag
#}
#
#main "$@"