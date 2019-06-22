#!/bin/bash

set -e
set -o pipefail

DOCKER_REPO_PREFIX=${DOCKER_REPO_PREFIX}
DOCKER_REGISTRY=${DOCKER_REGISTRY}
ERRORS="$(pwd)/errors"

echo ----------------------- Listing all Docker images ------------------------

list_dirs() {

  # Dockerfile relative path i.e alpine/Dockerfile
  docker_rel_path=$(find -L . -iname '*Dockerfile' | sed 's|./||' | sort)

  if [[ ! -z ${docker_rel_path} ]]; then

    for name in ${docker_rel_path}; do
      build_dir=$(dirname "$name")

      echo "${build_dir}"
    done

  fi
}

main() {
  list_dirs
}

main "$@"


echo -e "\n    Done.
  "