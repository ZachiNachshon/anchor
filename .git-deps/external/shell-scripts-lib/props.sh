#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"

_get_directory_path() {
  local name=$1
  local dir_path=""

  case "${name}" in
    app.cert.manager*)
      dir_path="app/cert-manager/"
      ;;
    app.docker.registry.ui.*)
      dir_path="app/docker-registry-ui/"
      ;;
    app.docker.registry.*)
      dir_path="app/docker-registry/"
      ;;
    app.kubernetes.dashboard.*)
      dir_path="app/kubernetes-dashboard/"
      ;;
    app.traefik.dashboard.*)
      dir_path="app/traefik-dashboard/"
      ;;
    app.traefik.*)
      dir_path="app/traefik/"
      ;;
    app.whoami.*)
      dir_path="app/whoami/"
      ;;
    docker.ansible.*)
      dir_path="docker/ansible/"
      ;;
    docker.bazel.*)
      dir_path="docker/bazel/"
      ;;
    docker.frp.*)
      dir_path="docker/frp/"
      ;;
    docker.heroku.*)
      dir_path="docker/heroku/"
      ;;
    docker.debian_buster.*)
      dir_path="docker/debian_buster/"
      ;;
    docker.registries.*)
      dir_path="docker/registries/"
      ;;
    k8s.k3s.*)
      dir_path="k8s/k3s/"
      ;;
    k8s.users.*)
      dir_path="k8s/users/"
      ;;
    rpi.nodes.*)
      dir_path="rpi/nodes/"
      ;;
    rpi.tools.*)
      dir_path="rpi/tools/"
      ;;
    rpi.cluster*)
      dir_path=""
      ;;
    rpi.sd_card*)
      dir_path="rpi/sd_card/"
      ;;
    rpi.proxy*)
      dir_path="rpi/proxy/"
      ;;
    tools.info*)
      dir_path="tools/info/"
      ;;
    *)
      echo -n "unknown"
      ;;
  esac

  echo "${dir_path}"
}

#######################################
# Return a config value for specific key from a properties file
# Globals:
#   None
# Arguments:
#   dir_path    - config.properties file absolute path
#   key         - key to return its configuration
#   do_not_fail - (optional) do not fail if config cannot be found, default: fail
# Usage:
#   property "absolute/path/to/config/dir" "some.test.key"
#   property "absolute/path/to/config/dir" "some.test.key" "do_not_fail"
#######################################
property() {
  # local dir_path=$(_get_directory_path "${name}")
  local dir_path=$1
  local key=$2
  local do_not_fail=$3
  local value=""

  # Consider piping ' | envsubst' at the end of the grep command
  value=$(grep "${key}" "${dir_path}/config.properties" | cut -d '=' -f2)
  if [[ -z "${value}" && "${do_not_fail}" != "do_not_fail" ]]; then
    log_fatal "missing property. key: ${key}"
  fi

  echo "${value}"
}

#######################################
# Return a formatted string pattern from a properties file
# Globals:
#   None
# Arguments:
#   dir_path - pattern.properties file absolute path
#   pattern  - string pattern identifier
#   ...      - string arguments that correspond to the %s within the pattern
# Usage:
#   pattern "absolute/path/to/pattern/dir" "some.test.pattern" "first" "second"
#######################################
pattern() {
  local dir_path=$1
  local pattern=$2

  local format=$(grep "${pattern}" "${dir_path}/pattern.properties" | cut -d '=' -f2- | envsubst)
  if [[ -z "${format}" ]]; then
    log_fatal "missing pattern. pattern: ${pattern}"
  fi

  shift
  shift
  printf ${format} "$@"
}
