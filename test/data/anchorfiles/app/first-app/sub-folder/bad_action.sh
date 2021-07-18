#!/bin/bash

source ./utils/logger.sh

function say_hello() {
  log_info "Hello"
}

main() {
  fake_missing_method $(say_hello)
}

main "$@"