#!/bin/bash

source ./app/first-app/.env
source ./utils/logger.sh

main() {
  name=${NAME}

  if [[ -z ${name} ]]; then
    log_fatal "Cannot resolve NAME, aborting"
  fi

  log_info "Goodbye World from ${NAME} !"
}

main "$@"

