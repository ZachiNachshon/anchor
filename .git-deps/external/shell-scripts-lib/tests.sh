#!/bin/bash

print_test_title() {
  local name=$1
  local title="
---------------------------------
  Test: ${name}
---------------------------------\n"
  printf "${title}" >&1
}
