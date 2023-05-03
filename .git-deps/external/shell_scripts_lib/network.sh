#!/bin/bash

# Title         Network utilities such as discoverability and information
# Author        Zachi Nachshon <zachi.nachshon@gmail.com>
# Supported OS  Linux & macOS
# Description   TODO
#==============================================================================

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/math.sh"
source "${CURRENT_FOLDER_ABS_PATH}/prompter.sh"

ARP_SCAN_DARWIN_AMD64="${CURRENT_FOLDER_ABS_PATH}/tooling/arp-scan/darwin_amd64/arp-scan"
ARP_SCAN_LINUX_AMD64="${CURRENT_FOLDER_ABS_PATH}/tooling/arp-scan/linux_amd64/arp-scan"
ARP_SCAN_LINUX_ARM="${CURRENT_FOLDER_ABS_PATH}/tooling/arp-scan/linux_arm/arp-scan"
ARP_SCAN_LINUX_ARM64="${CURRENT_FOLDER_ABS_PATH}/tooling/arp-scan/linux_arm64/arp-scan"

ask_for_rpi_node_username() {
  local default_user=$1
  local user=$(prompt_user_input "Enter RPi node user" ${default_user})
  echo "${user}"
}

ask_for_rpi_node_password() {
  local default_pass=$1
  local pass=$(prompt_for_password "Enter RPi node password" "${default_pass}")
  echo "${pass}"
}

read_network_address() {
  #  echo "0.0.0.0"
  #  echo "127.0.0.1"
  #  local public_ip=$(curl ifconfig.me)
  #  echo "${public_ip}"
  echo "192.168.1.102"
}

get_all_lan_network_devices() {
  local os_arch=$1
  local ip_range=$2
  local filter=$3
  local result=""

  if [[ "${os_arch}" == *"linux"* ]]; then

    local arpscan_bin="${ARP_SCAN_LINUX_AMD64}"
    if [[ "${os_arch}" == "linux_arm64" ]]; then
      arpscan_bin="${ARP_SCAN_LINUX_ARM64}"
    elif [[ "${os_arch}" == "linux_arm" ]]; then
      arpscan_bin="${ARP_SCAN_LINUX_ARM}"
    fi

    if [[ -n "${filter}" ]]; then
      result=$(sudo "${arpscan_bin}" "${ip_range}" | grep -v "WARNING" | grep '\t' | grep -i "${filter}")
    else
      result=$(sudo "${arpscan_bin}" "${ip_range}" | grep -v "WARNING" | grep '\t')
    fi

  elif [[ "${os_arch}" == *"darwin"* ]]; then

    local arpscan_bin="${ARP_SCAN_DARWIN_AMD64}"
    if [[ "${os_arch}" == *"arm"* ]]; then
      log_fatal "Apple silicon M1 is not supported yet."
    fi

    if [[ -n "${filter}" ]]; then
      result=$(sudo "${arpscan_bin}" "${ip_range}" | grep -v "WARNING" | grep '\t' | grep -i "${filter}")
    else
      result=$(sudo "${arpscan_bin}" "${ip_range}" | grep -v "WARNING" | grep '\t')
    fi

  else
    log_fatal "OS not supported. name: ${os_arch}"
  fi

  echo "${result}"
}

select_network_LAN_device() {
  local os_arch=$1
  local ip_range=$2
  local filter=$3
  local devices=$(get_all_lan_network_devices "${os_arch}" "${ip_range}" "${filter}")

  local delimiter="$"
  local devices_str=""
  while read -r line; do
    entry="${line}"
    devices_str+="${entry}${delimiter}"
  done <<<"${devices}"

  local selection=$(prompt_selection_allow_text "Select a RPi address ('r' to reload)" "${devices_str}" '$')
  if [[ -z "${selection}" ]]; then
    echo -e "\n    Invalid selection." >&0
    exit 1
  elif [[ "${selection}" == "r" ]]; then
    new_line
    selection=$(select_network_LAN_device "${os_arch}" "${ip_range}" "${filter}")
  fi

  echo "${selection}"
}

append_new_host_locally() {
  local hostname=$1
  local ip_address=$2
  local filepath="/etc/hosts"

  if is_file_contain "${filepath}" "${hostname}"; then
    log_info "Entry '${hostname}' already exists in ${filepath}"
  else
    echo -e "${ip_address}\t${hostname}" | sudo tee -a "${filepath}" >/dev/null
    log_info "Added entry to ${filepath} successfully"
  fi
}

print_instructions_network_scan() {
  echo -e """${COLOR_YELLOW}
  ================================================================================================
  ${COLOR_RED}Elevated user permissions are required for this step !${COLOR_YELLOW}

  This step scans all devices on the LAN network and lists the following:
    • IP Address
    • MAC Address
    • Device Name
  ================================================================================================${COLOR_NONE}
""" >&1
}

print_instructions_connect_via_ssh() {
  local ip_address=$1
  local default_user=$2
  local default_pass=$3

  echo -e """${COLOR_YELLOW}
  ================================================================================================
  About to run a script over SSH on address ${ip_address}.

  Requirements:
    • Ansible or Docker
      If Ansible is missing, a Docker image will be built and used instead.

  ${COLOR_RED}This step prompts for connection access information (press ENTER for defaults):
    • Raspberry Pi node user     (default: ${default_user})
    • Raspberry Pi node password (default: ${default_pass})${COLOR_YELLOW}
  ================================================================================================${COLOR_NONE}
""" >&1
}
