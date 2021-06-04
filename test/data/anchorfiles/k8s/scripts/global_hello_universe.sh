#!/bin/bash

source ./utils/logger.sh

function say_global_hello() {
  log_info "Global Hello Universe ${NAME}"
}

main() {
  decorate_greeting $(say_global_hello)
}

main "$@"