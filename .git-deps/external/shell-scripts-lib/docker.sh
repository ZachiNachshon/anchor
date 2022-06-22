#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"
source "${CURRENT_FOLDER_ABS_PATH}/prompter.sh"

prompt_container_registries() {
  # Create a list of all available registries
  local domain_name=$(property rpi.cluster.k8s.ingress.domain.name)
  local sub_domain_name=$(property app.docker.registry.k8s.ingress.sub.domain)
  local url="${sub_domain_name}.${domain_name}.com"

  local registries_list="${url}"
  local selection=$(prompt_selection "Select a Registry" "${registries_list}")
  echo "${selection}"
}

docker_login() {
  local registry=$1
  local user=$2
  local pass=$3
  local url="https://${registry}"

  echo "${pass}" | docker login --username "${user}" --password-stdin "${url}"
}

clean_container() {
  local name=$1
  stop_docker_container "${name}"
  remove_stopped_container "${name}"
}

remove_image() {
  local name=$1
  clean_container "${name}"
  docker_remove_image "${name}"
}

#######################################
# Log a running Docker container by name/prefix/suffix
# Globals:
#   None
# Arguments:
#   name - Docker running container name (or prefix/suffix)
# Usage:
#   docker_log_container "busybox"
#######################################
docker_log_container() {
  local name=$1
  local containersToLog=$(docker ps | grep "${name}" | awk {'print $1'})
  if [[ -n "${containersToLog}" ]]; then
    log_info "Logging container. name: ${name}, id: ${containersToLog}"
    docker logs -f "${containersToLog}"
  else
    log_info "No running container found. name: ${name}"
  fi
}

#######################################
# Stop a running Docker container by name/prefix/suffix
# Globals:
#   None
# Arguments:
#   name - Docker running container name (or prefix/suffix)
# Usage:
#   docker_stop_running_container "busybox"
#######################################
docker_stop_running_container() {
  local name=$1
  local containersToStop=$(docker ps | grep "${name}" | awk {'print $1'})
  if [[ -n "${containersToStop}" ]]; then
    log_info "Stopping container. name: ${name}, id: ${containersToStop}"
    docker stop "${containersToStop}"
  else
    log_info "No running container found. name: ${name}"
  fi
}

#######################################
# Remove Docker images by name/prefix/suffix
# Globals:
#   None
# Arguments:
#   name - Docker image name (or prefix/suffix)
# Usage:
#   docker_remove_image "busybox"
#######################################
docker_remove_image() {
  local name=$1
  local imagesToRemove=$(docker images | grep "${name}" | awk {'print $1'} | xargs)
  if [[ -n "${imagesToRemove}" ]]; then
    log_info "Removing container image. name: ${imagesToRemove}"
    docker rmi -f "${imagesToRemove}"
  else
    log_info "No container image can be found. name: ${name}"
  fi
}

#######################################
# Remove Docker dangling intermediate images (a.k.a <none>)
# Globals:
#   None
# Arguments:
#   None
# Usage:
#   docker_remove_dangling_intermediate_images
#######################################
docker_remove_dangling_intermediate_images() {
  local dangling_images="$(docker images -f "dangling=true" -q)"
  if [[ -n "${dangling_images}" ]]; then
    log_info "Removing dangling intermediate images..."
    local images_arr=($(echo $dangling_images | tr '\n' ' '))
    for ((i = 0; i < ${#images_arr[@]}; i++)); do
      local image_to_remove=${images_arr[i]}
      if [[ -n ${image_to_remove} ]]; then
        log_info "Removing image. id: ${image_to_remove}"
        docker rmi -f "${image_to_remove}"
      fi
    done
  else
    log_info "No dangling images found"
  fi
}

#######################################
# Remove Docker stopped container by name/prefix/suffix
# Globals:
#   None
# Arguments:
#   name - Docker stopped container name (or prefix/suffix)
# Usage:
#   docker_remove_stopped_container "busybox"
#######################################
docker_remove_stopped_container() {
  local name=$1
  local containerToRemove=$(docker ps -a | grep "${name}" | awk {'print $1'})
  if [[ -n "${containerToRemove}" ]]; then
    log_info "Removing stopped container. name: ${name}, id: ${containerToRemove}"
    docker rm -f ${containerToRemove}
  else
    log_info "No stopped container found. name: ${name}"
  fi
}

#######################################
# Remove Docker stopped containers (a.k.a docker ps -a)
# Globals:
#   None
# Arguments:
#   None
# Usage:
#   docker_remove_stopped_containers
#######################################
docker_remove_stopped_containers() {
  docker container prune -f
}
