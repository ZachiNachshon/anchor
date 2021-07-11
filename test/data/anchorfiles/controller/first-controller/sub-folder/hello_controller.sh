#!/bin/bash

source ./controller/first-controller/.env
source ./utils/logger.sh

main() {
  name=${NAME}

  if [[ -z ${name} ]]; then
    log_fatal "Cannot resolve NAME, aborting"
  fi

  log_info "Hello Controller from ${NAME} !"
}

main "$@"
