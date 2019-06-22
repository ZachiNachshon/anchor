#!/bin/bash
#
# This script allows you to run specific container image
# by the directory name the Dockerfile resides in.
#
set -e
set -o pipefail

DOCKER_REPO_PREFIX=${DOCKER_REPO_PREFIX}
DOCKER_REGISTRY=${DOCKER_REGISTRY}
DOT_FILES=${DOT_FILES}

echo ---------------- Deploying Docker Container on LOCAL Machine -----------------

run() {
  docker_func_name=$1
  source ${DOT_FILES}/.dockerfunc
  ${docker_func_name}

  sleep 1
}

main() {

  if [[ $# -eq 0 ]]; then
    echo "Usage: $0 image1 image2 ..."
    exit 1
  fi

  # directory name
	dir=$1

	# tail logs after docker run
	tail=$2

  if [[ ! -d "${dir}" ]]; then
    echo "Unable to find container configuration directory with name: ${dir}"
    exit 1
  else
    run ${dir}

    echo -e "\n    Done.
      "

    if [[ ${tail} = "--tail" ]]; then
      docker_id=$(docker ps -a | grep ${DOCKER_REPO_PREFIX}/${dir} | awk '{print $1}')
      echo "docker logs -f ${docker_id}"
      docker logs -f ${docker_id}
    fi
  fi
}

main "$@"



