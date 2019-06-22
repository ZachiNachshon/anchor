#!/bin/bash
#
# This script allows you to clean container
# images from previous docker build and
# additional unknown container layers.
#
set -e
set -o pipefail

DOCKER_REPO_PREFIX=${DOCKER_REPO_PREFIX}
DOCKER_REGISTRY=${DOCKER_REGISTRY}
DOT_FILES=${DOT_FILES}

echo ------------------------- Cleanup: Previous Images ---------------------------

#-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
#         CLEANUP: IMAGES
#-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

clean_unknown_images() {
  # Remove docker images under same name if exist and remove partially created images, if any
  unknownImages=$(docker images | grep "<none>" | awk {'print $3'})
  if [[ ! -z "${unknownImages}" ]]; then
      echo "Removing all unknown <none> images..."
      docker rmi -f ${unknownImages}
  fi
}

clean_container_images() {
  img_name=$1

  previousImages=$(docker images | grep "${DOCKER_REPO_PREFIX}/${img_name}" | awk {'print $3'})
    if [[ ! -z "${previousImages}" ]]; then
        echo "Removing previous docker images..."
        docker rmi -f ${previousImages}
    fi
}

main() {

  if [[ $# -eq 0 ]]; then
    echo "Usage: $0 image1 image2 ..."
    exit 1
  fi

  for name in "$@"; do
    if [[ ! -d "$name" ]]; then
      echo "Unable to find container configuration directory with name: $name"
      exit 1
    else
      echo "Cleaning: ${name}"
      clean_unknown_images || true
      clean_container_images ${name} || true
    fi

    shift
  done
}

main "$@"

echo -e "\n    Done.
  "