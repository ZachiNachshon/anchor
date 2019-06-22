#!/bin/bash

set -e
set -o pipefail

DOCKER_REPO_PREFIX=${DOCKER_REPO_PREFIX}
DOCKER_REGISTRY=${DOCKER_REGISTRY}
ERRORS="$(pwd)/errors"

echo ------------------------- Building Docker Image --------------------------

build() {

  # Directory name i.e alpine
	dir_name=$1

  # Dockerfile relative path i.e alpine/Dockerfile
  docker_rel_path=$(find -L ./${dir_name} -iname '*Dockerfile' | sed 's|./||' | sort)

  build_dir=$(dirname "$docker_rel_path")
  image=${dir_name}

  # We could add an inner folder to indicate the version number
#	suite=${build_dir##*\/}
	suite="latest"

  echo "Building ${DOCKER_REPO_PREFIX}/${image}:${suite} for context ${build_dir}"
  echo "docker build -f ${docker_rel_path} -t "${DOCKER_REPO_PREFIX}/${image}:${suite}" "${build_dir}""
	docker build -f ${docker_rel_path} -t "${DOCKER_REPO_PREFIX}/${image}:${suite}" "${build_dir}" || return 1

  #
  # Note: --rm --force-rm removes the image before building anew
  #
#	echo "docker build --rm --force-rm -f ${docker_rel_path} -t "${DOCKER_REPO_PREFIX}/${image}:${suite}" "${build_dir}""
#	docker build --rm --force-rm -f ${docker_rel_path} -t "${DOCKER_REPO_PREFIX}/${image}:${suite}" "${build_dir}" || return 1
}

push() {
  echo "Pushing to ${DOCKER_REGISTRY}"
}

main() {

  if [[ $# -eq 0 ]]; then
    echo "Usage: $0 image1 image2 ..."
    exit 1
  fi

  # Directory name
	dir=$1

	# Flag --push used for pushing to registry
	push=$2

  if [[ ! -d "${dir}" ]]; then
    echo "Unable to find container configuration directory with name: ${dir}"
    exit 1
  else
    build ${dir}

    if [[ ${push} = "--push" ]]; then
      push
    fi
  fi
}

main "$@"


echo -e "\n    Done.
  "