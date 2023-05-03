#!/bin/bash

#######################################
# Archive using tar with exclusion option
# Globals:
#   None
# Arguments:
#   source_directory - directory of the source files
#   archive_filename - the archive filename
#   dest_folder      - destination folder
#   exclusions       - exclusion list to exclude from the archive
# Usage:
#   open_browser "http://www.google.com"
#######################################
tar_archive() {
  local all_args=("$@")
  local source_directory=$1
  local archive_filename=$2
  local dest_folder=$3
  local exclusions=$("${all_args[@]:3}")

  tar -czvf
  #   tar --exclude='file1.txt' --exclude='folder1' -zcvf backup.tar.gz .
}
