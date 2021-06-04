#!/bin/bash

ANCHOR_VERSION="v0.4.0"

print_banner() {
  echo -e "
 █████╗ ███╗   ██╗ ██████╗██╗  ██╗ ██████╗ ██████╗
██╔══██╗████╗  ██║██╔════╝██║  ██║██╔═══██╗██╔══██╗
███████║██╔██╗ ██║██║     ███████║██║   ██║██████╔╝
██╔══██║██║╚██╗██║██║     ██╔══██║██║   ██║██╔══██╗
██║  ██║██║ ╚████║╚██████╗██║  ██║╚██████╔╝██║  ██║
╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝ (${ANCHOR_VERSION})
"
}

clone_anchor_dockerfiles_repo() {
  echo -e "
=======================================================================
                          Cloning Repositories
======================================================================="

  echo -e "
===========
Repository: anchor-dockerfiles
===========
"

  git clone "https://github.com/ZachiNachshon/anchor-dockerfiles.git" ${HOME}/.anchor/anchor-dockerfiles
  echo -e "\n    Done."
}

set_dockerfiles_env_var() {

  echo -e "
=======================================================================
                       Setting Environment Variables
======================================================================="

  # Copy to clipboard on macOS/Linux
  local exportCmd="export DOCKER_FILES=${HOME}/.anchor/anchor-dockerfiles"
  if [[ "$OSTYPE" == "linux"* ]]; then
    xclip -selection "${exportCmd}"
  elif [[ "$OSTYPE" == "darwin"* ]]; then
    echo "${exportCmd}" | pbcopy
  fi

  echo -e "\nRun the following command to link between anchor to DOCKER_FILES (paste from clipboard):

  ${exportCmd}
"
}

verify_pre_setup() {
  read -p "Do you want to download anchor-dockerfiles to ${HOME}/.anchor? (y/n): " input
  if [[ ${input} != "y" ]]; then
    echo -e "\n    Nothing has changed.\n"
    exit 0
  fi
}

main() {
  print_banner
  verify_pre_setup
  clone_anchor_dockerfiles_repo
  set_dockerfiles_env_var
}

main "$@"
