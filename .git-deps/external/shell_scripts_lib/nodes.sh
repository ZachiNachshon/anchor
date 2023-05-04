#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/props.sh"
source "${CURRENT_FOLDER_ABS_PATH}/math.sh"

read_master_address_as_ansible_str() {
  local master_name=$(property rpi.cluster.master.name)
  local master_address=$(property "rpi.cluster.address.${master_name}")
  echo "${master_name} ansible_host=${master_address}"
}

read_workers_addresses_as_ansible_str() {
  local workers_addresses=""
  local i=1
  local name=$(property "rpi.cluster.node.${i}.name" do_not_fail)
  local address=$(property "rpi.cluster.address.${name}" do_not_fail)
  while [[ -n "${name}" && -n "${address}" ]]; do
    workers_addresses="${workers_addresses}${name} ansible_host=${address}\n"
    ((i++))
    name=$(property "rpi.cluster.node.${i}.name" do_not_fail)
    address=$(property "rpi.cluster.address.${name}" do_not_fail)
  done

  # Use printf to enforce new lines
  printf "${workers_addresses}"
}

# arr=()	        Create an empty array
# arr=(1 2 3)	    Initialize array
# ${arr[2]}	      Retrieve third element
# ${arr[@]}	      Retrieve all elements
# ${!arr[@]}	    Retrieve array indices
# ${#arr[@]}	    Calculate array size
# arr[0]=3	      Overwrite 1st element
# arr+=(4)	      Append value(s)
# str=$(ls)	      Save ls output as a string
# arr=( $(ls) )	  Save ls output as an array of files
# ${arr[@]:s:n}	  Retrieve n elements starting at index s

select_node() {
  local prompt_master=$1
  local prompt_nodes=$2
  local skip_option_all=$3

  selection_array=()

  local menu_str="\nSelect a Raspberry Pi node:\n\n"
  menu_idx=1
  if [[ "${prompt_master}" == "prompt_master" ]]; then
    local master_name=$(property rpi.cluster.master.name)
    local master_address=$(property "rpi.cluster.address.${master_name}")
    local master_address_ansible="${master_name} ansible_host=${master_address}"
    menu_str+="  ${menu_idx}. ${master_address_ansible}\n"
    selection_array+=("${master_address_ansible}")
    menu_idx=$(increment ${menu_idx})
  fi

  if [[ "${prompt_nodes}" == "prompt_nodes" ]]; then
    local workers_addresses=()
    local node_idx=1
    local name=$(property "rpi.cluster.node.${node_idx}.name" do_not_fail)
    local address=$(property "rpi.cluster.address.${name}" do_not_fail)
    while [[ -n "${name}" && -n "${address}" ]]; do
      node_address_ansible="${name} ansible_host=${address}"
      menu_str+="  ${menu_idx}. ${node_address_ansible}\n"
      menu_idx=$(increment ${menu_idx})
      selection_array+=("${node_address_ansible}")
      node_idx=$(increment ${node_idx})
      name=$(property "rpi.cluster.node.${node_idx}.name" do_not_fail)
      address=$(property "rpi.cluster.address.${name}" do_not_fail)
    done
  fi

  if [[ -z "${prompt_master}" && -z "${prompt_nodes}" ]]; then
    echo -e "\n    Nothing to select." >&0
    exit 0
  fi

  if [[ -z "${skip_option_all}" ]]; then
    # Append a menu option to select all prompt options
    all_nodes_option="All the above"
    selection_array+=("${all_nodes_option}")
    menu_str+="  ${menu_idx}. ${all_nodes_option}\n"
  fi

  # Use printf to enforce new lines
  printf "${menu_str}\nPlease choose a node (enter to skip): " >&0
  read input
  if [[ -n "${input}" ]]; then
    selected_value="${selection_array[input - 1]}" >>/dev/null
    if [[ -z "${selected_value}" ]]; then
      echo -e "\n    Invalid selection." >&0
      exit 1
    else
      if [[ "${selected_value}" == "${all_nodes_option}" ]]; then
        selected_value=""
        y=0
        # Wrap up all nodes in an ansible format
        for _ in ${selection_array[@]}; do
          node=${selection_array[y]}
          if [[ "${node}" != "${all_nodes_option}" ]]; then
            selected_value+="${node}\n"
          fi
          y=$(increment ${y})
        done
        selected_value=$(printf "${selected_value}")
      fi
    fi
  else
    echo -e "\n    Nothing has changed." >&0
    exit 0
  fi

  echo "${selected_value}"
}
