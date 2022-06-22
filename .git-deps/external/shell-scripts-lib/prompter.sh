#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")

source "${CURRENT_FOLDER_ABS_PATH}/logger.sh"
source "${CURRENT_FOLDER_ABS_PATH}/math.sh"

#######################################
# Prompt for enter key
# Globals:
#   None
# Arguments:
#   None
# Usage:
#   prompt_for_enter
#######################################
prompt_for_enter() {
  printf "${COLOR_GREEN}  Press ENTER to continue...${COLOR_NONE}" >&0
  read input
}

#######################################
# Prompt for a password, input is obfuscated
# Globals:
#   None
# Arguments:
#   message - prompt message
#   default - (optional) use a default password
# Usage:
#   prompt_for_password "Insert secret" "default-pass"
#######################################
prompt_for_password() {
  local message=$1
  local default=$2

  if [[ -n "${default}" ]]; then
    # Use printf to enforce new lines
    printf "${message} (default: ${default}): " >&0
  else
    printf "${message} (enter to skip): " >&0
  fi

  read -s password
  if [[ -z "${password}" ]]; then
    if [[ -n "${default}" ]]; then
      password="${default}"
    else
      echo -e "\n\n    Nothing has changed." >&0
      exit 0
    fi
  fi

  echo "${password}"
}

#######################################
# Prompt a yes/no question with severity levels
# Globals:
#   None
# Arguments:
#   message - prompt message
#   level   - (optional: critical/warning) highlight text colors
# Usage:
#   prompt_yes_no "Do you want to proceed" "warning"
#######################################
prompt_yes_no() {
  local message=$1
  local level=$2

  local prompt=""
  if [[ ${level} == "critical" ]]; then
    prompt="${COLOR_RED}${message}? (y/n):${COLOR_NONE} "
  elif [[ ${level} == "warning" ]]; then
    prompt="${COLOR_YELLOW}${message}? (y/n):${COLOR_NONE} "
  else
    prompt="${message}? (y/n): "
  fi

  printf "${prompt}" >&0
  read input
  if [[ "${input}" != "y" ]]; then
    input=""
  fi

  echo "${input}"
}

#######################################
# Prompt for user input and return
# Globals:
#   None
# Arguments:
#   message - prompt message
#   default - (optional) fallback option if no selection
# Usage:
#   prompt_user_input "Enter input" "default_value"
#######################################
prompt_user_input() {
  local message=$1
  local default=$2

  if [[ -n "${default}" ]]; then
    # Use printf to enforce new lines
    printf "${message} (default: ${default}): " >&0
  else
    printf "${message} (enter to abort): " >&0
  fi

  read input
  if [[ -z "${input}" ]]; then
    if [[ -n "${default}" ]]; then
      input="${default}"
    else
      echo -e "\n    Nothing has changed." >&0
      exit 0
    fi
  fi

  echo "${input}"
}

#######################################
# Prompt for selection from a delimited string
# Globals:
#   None
# Arguments:
#   title               - prompt message title
#   delimited_items_str - items as string separated by a delimiter
#   delimiter           - (optional) delimiter to use as items separator, defaults to ' '
# Usage:
#   prompt_selection "Please select" "one;two;three" ";"
#######################################
prompt_selection() {
  local title=$1
  local delimited_items_str=$2
  local delimiter=$3
  local numeric_menu_str="\n${title}:\n\n"

  # By default use space as delimiter
  if [[ -z "${delimiter}" ]]; then
    delimiter=" "
  fi

  if [[ -n "${delimited_items_str}" ]]; then
    local saveIFS=$IFS
    IFS="${delimiter}"
    read -r -a selection_array <<<"${delimited_items_str}"
    IFS=${saveIFS}
    local menu_idx=1

    for ((i = 0; i < ${#selection_array[@]}; i++)); do
      local value=${selection_array[i]}
      numeric_menu_str+="  ${menu_idx}. ${value}\n"
      menu_idx=$(increment ${menu_idx})
    done

    numeric_menu_str="${numeric_menu_str}\nPlease choose (enter to skip): "

    # Use printf to enforce new lines
    printf "${numeric_menu_str}" >&0
    read input
    if [[ -n "${input}" ]]; then

      local is_number=$(is_number "${input}")
      if [[ -n ${is_number} ]]; then
        local selected_value="${selection_array[input - 1]}" >>/dev/null
        if [[ -z "${selected_value}" ]]; then
          echo -e "\n    Invalid selection." >&0
          # exit 1
        fi
        result="${selected_value}"
      else
        echo -e "\n    Invalid input." >&0
        exit 1
      fi
    else
      echo -e "\n    Nothing has changed." >&0
      exit 0
    fi
  else
    log_warning "No values to prompt, skipping"
  fi

  echo "${result}"
}

#######################################
# Prompt for selection from a delimited string
# Allow custom text selection which is not in the avaialble options
# Globals:
#   None
# Arguments:
#   title               - prompt message title
#   delimited_items_str - items as string separated by a delimiter
#   delimiter           - (optional) delimiter to use as items separator, defaults to ' '
# Usage:
#   prompt_selection_allow_text "Select a value ('r' to reload)" "one$two$three" '$'
#######################################
prompt_selection_allow_text() {
  local title=$1
  local delimited_items_str=$2
  local delimiter=$3
  local numeric_menu_str="\n${title}:\n\n"

  # By default use space as delimiter
  if [[ -z "${delimiter}" ]]; then
    delimiter=" "
  fi

  if [[ -n "${delimited_items_str}" ]]; then
    local saveIFS=$IFS
    IFS="${delimiter}"
    read -r -a selection_array <<<"${delimited_items_str}"
    IFS=${saveIFS}
    local menu_idx=1

    for ((i = 0; i < ${#selection_array[@]}; i++)); do
      local value=${selection_array[i]}
      numeric_menu_str+="  ${menu_idx}. ${value}\n"
      menu_idx=$(increment ${menu_idx})
    done

    numeric_menu_str="${numeric_menu_str}\nPlease choose (enter to skip): "

    # Use printf to enforce new lines
    printf "${numeric_menu_str}" >&0
    read input
    if [[ -n "${input}" ]]; then

      local is_number=$(is_number "${input}")
      if [[ -n ${is_number} ]]; then
        local selected_value="${selection_array[input - 1]}" >>/dev/null
        if [[ -z "${selected_value}" ]]; then
          echo -e "\n    Invalid selection." >&0
          exit 1
        fi
        result="${selected_value}"
      else
        # Do not fail, allow non-numeric values, can be used for selection retry actions
        # or customized input
        result=${input}
      fi
    else
      echo -e "\n    Nothing has changed." >&0
      exit 0
    fi
  else
    log_warning "No values to prompt, skipping"
  fi

  echo "${result}"
}
