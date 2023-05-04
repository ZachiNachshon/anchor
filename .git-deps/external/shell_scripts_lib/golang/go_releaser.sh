#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")
ROOT_FOLDER_ABS_PATH=$(dirname "${CURRENT_FOLDER_ABS_PATH}")

source "${ROOT_FOLDER_ABS_PATH}/logger.sh"
source "${ROOT_FOLDER_ABS_PATH}/io.sh"
source "${ROOT_FOLDER_ABS_PATH}/os.sh"
source "${ROOT_FOLDER_ABS_PATH}/cmd.sh"
source "${ROOT_FOLDER_ABS_PATH}/strings.sh"
source "${ROOT_FOLDER_ABS_PATH}/checks.sh"
source "${ROOT_FOLDER_ABS_PATH}/github.sh"

SCRIPT_MENU_TITLE="Go Releaser"

DEFAULT_BUILD_PATH="./..."
DEFAULT_BUILD_OUTPUT_PATH="dist"
DEFAULT_GO_RELEASER_CONFIG_FILEPATH=".goreleaser.yml"

RESOLVED_GOPATH="${GOPATH}"
RESOLVED_GOBIN="${GOBIN:-${RESOLVED_GOPATH}/bin}"

CLI_ARGUMENT_BUILD=""
CLI_ARGUMENT_INSTALL=""
CLI_ARGUMENT_PUBLISH=""
CLI_ARGUMENT_DELETE=""

CLI_FLAG_BUILD_MAIN_PACKAGE_PATH=""   
CLI_FLAG_BUILD_OUTPUT_PATH=""         
CLI_FLAG_BUILD_OS_ARCH_TUPLE=""       
CLI_FLAG_BUILD_ARCHIVE=""       

CLI_FLAG_RELEASE_TYPE="" # options: github
CLI_FLAG_RELEASE_TAG=""            
CLI_FLAG_GO_RELEASER_CONFIG_FILEPATH=""            

CLI_FLAG_DELETE_ORIGIN="" # options: gobin-local/github
CLI_FLAG_DELETE_TAG=""

CLI_VALUE_BUILD_MAIN_PACKAGE_PATH=""
CLI_VALUE_BUILD_OUTPUT_PATH=""
CLI_VALUE_BUILD_OS_ARCH_TUPLES_ARRAY=()

CLI_VALUE_RELEASE_TYPE=""
CLI_VALUE_RELEASE_TAG=""
CLI_VALUE_GO_RELEASER_CONFIG_FILEPATH=""

CLI_VALUE_DELETE_ORIGIN=""         
CLI_VALUE_DELETE_TAG=""

is_build() {
  [[ -n "${CLI_ARGUMENT_BUILD}" ]]
}

is_main_package_build_path() {
  [[ -n "${CLI_FLAG_BUILD_MAIN_PACKAGE_PATH}" && -n "${CLI_VALUE_BUILD_MAIN_PACKAGE_PATH}" ]]
}

is_build_specific_os_arch() {
  [[ -n "${CLI_FLAG_BUILD_OS_ARCH_TUPLE}" ]] && [[ ${#CLI_VALUE_BUILD_OS_ARCH_TUPLES_ARRAY[@]} -gt 0 ]]
}

get_build_path() {
  echo "${CLI_VALUE_BUILD_MAIN_PACKAGE_PATH:-${DEFAULT_BUILD_PATH}}"
}

get_build_output_path() {
  echo "${CLI_VALUE_BUILD_OUTPUT_PATH:-${DEFAULT_BUILD_OUTPUT_PATH}}"
}

is_build_archive() {
  [[ -n "${CLI_FLAG_BUILD_ARCHIVE}" ]]
}

is_install() {
  [[ -n "${CLI_ARGUMENT_INSTALL}" ]]
}

is_publish() {
  [[ -n "${CLI_ARGUMENT_PUBLISH}" ]]
}

get_publish_goreleaser_config_path() {
  echo "${CLI_VALUE_GO_RELEASER_CONFIG_FILEPATH:-${DEFAULT_GO_RELEASER_CONFIG_FILEPATH}}"
}

is_github_release_type() {
  [[ "${CLI_VALUE_RELEASE_TYPE}" == "github" ]]
}

is_release_tag_exist() {
  [[ -n "${CLI_FLAG_RELEASE_TAG}" ]]
}

get_release_tag() {
  echo "${CLI_VALUE_RELEASE_TAG}"
}

is_delete() {
  [[ -n "${CLI_ARGUMENT_DELETE}" ]]
}

is_delete_exist() {
  [[ -n "${CLI_FLAG_DELETE_TAG}" && -n "${CLI_VALUE_DELETE_TAG}" ]]
}

get_delete_tag() {
  echo "${CLI_VALUE_DELETE_TAG}"
}

is_gobin_local_delete_origin() {
  [[ "${CLI_VALUE_DELETE_ORIGIN}" == "gobin-local" ]]
}

is_github_delete_origin() {
  [[ "${CLI_VALUE_DELETE_ORIGIN}" == "github" ]]
}

publish_binaries_to_github_using_goreleaser() {
  local tag=$(get_release_tag)
  local goreleaser_config_path=$(get_publish_goreleaser_config_path)
  new_line
  if github_prompt_for_approval_before_release "${tag}"; then
    if github_is_release_tag_exist "${tag}"; then
      log_fatal "GitHub release tag already exist, cannot override. tag: ${tag}"
    else
      github_create_release_tag "${tag}"
      log_info "Publishing Go binaries to GitHub using goreleaser. config: ${goreleaser_config_path}"
      cmd_run "GORELEASER_CURRENT_TAG=${tag} goreleaser release --clean --config=${goreleaser_config_path}"
    fi
  else
    log_warning "Nothing was uploaded."
  fi
}

publish_binaries() {
  if is_github_release_type; then
    check_tool "gh"
    check_tool "goreleaser"
    publish_binaries_to_github_using_goreleaser
  else
    log_fatal "Invalid publish release type. value: ${CLI_FLAG_RELEASE_TYPE}"
  fi
}

build_binaries() {
  if is_build_specific_os_arch; then
    for tuple in "${CLI_VALUE_BUILD_OS_ARCH_TUPLES_ARRAY[@]}"; do
      os_arch_split=(${tuple//_/ })
      build_binary_os_arch "${os_arch_split[0]}" "${os_arch_split[1]}"
    done
  else
    build_binary_os_arch
  fi
}

build_binary_os_arch() {
  local os=$1
  local arch=$2
  local build_path=$(get_build_path)
  local output_folder=$(get_build_output_path)
  local binary_name=$(get_project_name) # Only when the --main-pacakge flag is not used
  local build_vars=""
  local build_flags=""

  # When building for specific OS and ARCH, build duration takes 400% longer
  # Thus, when building locally we can avoid that penalty
  if [[ -n "${os}" && -n "${arch}" ]]; then
    output_folder="${output_folder}/${os}_${arch}"
    build_vars="CGO_ENABLED=0 GOARCH=${arch} GOOS=${os}"
    build_flags="-a -ldflags '-extldflags "-static"' -mod=readonly"
    log_info "Building Go binary. os: ${os}, arch: ${arch}"
  else
    local os_system=$(read_os_type)
    local arch_system=$(read_arch "x86_64:amd64")
    log_info "Building Go binary. os: ${os_system}, arch: ${arch_system}"
  fi

  if ! is_directory_exist "${output_folder}"; then
    cmd_run "mkdir -p ${output_folder}"
  fi

  # When running go build with a pattern use a directory path with -o flag
  # When running go build on a specific package, use a filepath with -o flag
  # Error when misused:
  #   go: cannot write multiple packages to non-directory
  if is_main_package_build_path; then
    output_folder="${output_folder}/${binary_name}"
  fi

  cmd_run "${build_vars} go build -o ${output_folder} ${build_flags} ${build_path}"
  log_info "Binary built successfully. path: ${output_folder}"

  # if is_build_archive; then
  #   local cwd=$(pwd)
  #   local archive_path="${output_folder}"
  #   local file_name="${binary_name}"

  #   # If we're using --main-package flag, we need to change dir into a parent dir
  #   if is_file_exist "${output_folder}"; then
  #     file_name=$(basename "${output_folder}")
  #     archive_path=$(dirname "${output_folder}")
  #   fi

  #   if ! is_dry_run; then
  #     cmd_run "cd ${archive_path} || exit"
  #   fi

  #   log_info "Archiving go binary into a tarball. file: ${file_name}.zip, path: ${archive_path}"

  #   local tarball_name="${file_name}.zip"
  #   if [[ -n "${os}" && -n "${arch}" ]]; then
  #     tarball_name="${file_name}_${os}_${arch}.zip"
  #   fi
  #   cmd_run "tar -zcvf ${tarball_name} ."
  #   if ! is_dry_run; then
  #     cmd_run "cd ${cwd} || exit"
  #   fi
  # fi
}

install_binary() {
  local build_path=$(get_build_path)
  log_info "Installing Go binary. build_path: ${build_path}"
  cmd_run "go install ${build_path}"
  log_info "Go binary installed. path: ${RESOLVED_GOBIN}"
}

delete_binary_from_gobin() {
  local binary_name=$(get_project_name)
  local delete_binary_path="${RESOLVED_GOBIN}/${binary_name}"
  log_info "Deleting Go binary. path: ${delete_binary_path}"
  if is_file_exist "${delete_binary_path}"; then
    cmd_run "rm -rf ${delete_binary_path}"
  else
    log_warning "No binary can be found for deletion"
  fi
}

delete_release_from_github() {
  local tag="no-tag"
  if is_delete_exist; then
    tag=$(get_delete_tag)
  else
    tag=$(prompt_user_input "Insert tag to delete")
  fi
  if [[ -n "${tag}" ]]; then
    if github_is_release_tag_exist "${tag}"; then
        log_info "GitHub release tag was found. tag: ${tag}"
        if github_prompt_for_approval_before_delete "${tag}"; then
          github_delete_release_tag "${tag}"
        else
          log_info "No GitHub release tag was deleted."
        fi
      else
        log_warning "No GitHub release tag was found. tag: ${tag}"
      fi
  fi
}

delete_binary_or_release() {
  if is_gobin_local_delete_origin; then
    delete_binary_from_gobin
  elif is_github_delete_origin; then
    check_tool "gh"
    delete_release_from_github
  else
    log_fatal "Flag --origin has invalid value or missing a value. value: ${CLI_VALUE_DELETE_ORIGIN}"
  fi
}

print_help_menu_and_exit() {
  local exec_filename=$1
  local file_name=$(basename "${exec_filename}")
  echo -e ""
  echo -e "${SCRIPT_MENU_TITLE} - Build, Install and Release a Go package"
  echo -e " "
  echo -e "${COLOR_WHITE}USAGE${COLOR_NONE}"
  echo -e "  "${file_name}" [command] [option] [flag]"
  echo -e " "
  echo -e "${COLOR_WHITE}ARGUMENTS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}build${COLOR_NONE}                     Build local Go binary(ies)"
  echo -e "  ${COLOR_LIGHT_CYAN}install${COLOR_NONE}                   Build and Install a Go binary locally"
  echo -e "  ${COLOR_LIGHT_CYAN}publish${COLOR_NONE}                   Build and publish Go binary(ies) as GitHub release"
  echo -e "  ${COLOR_LIGHT_CYAN}delete${COLOR_NONE}                    Delete a locally installed Go binary or a remote release"
  echo -e " "
  echo -e "${COLOR_WHITE}BUILD FLAGS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}--main-package${COLOR_NONE}            Optional package of the ${COLOR_GREEN}main.go${COLOR_NONE} file (default: ${COLOR_GREEN}${DEFAULT_BUILD_PATH}${COLOR_NONE})"
  echo -e "  ${COLOR_LIGHT_CYAN}--output-path${COLOR_NONE}             Optional output path for the generated binaries (default: ${COLOR_GREEN}${DEFAULT_BUILD_OUTPUT_PATH}${COLOR_NONE})"
  echo -e "  ${COLOR_LIGHT_CYAN}--os-arch${COLOR_NONE}                 Optional repeatable OS_ARCH flag i.e. linux_amd64 (default: ${COLOR_GREEN}System${COLOR_NONE} or ${COLOR_GREEN}\$GOOS_\$GOARCH${COLOR_NONE})"
  echo -e "  ${COLOR_LIGHT_CYAN}--archive${COLOR_NONE}                 Archive the build result into tar.gz"
  echo -e " "
  echo -e "${COLOR_WHITE}PUBLISH FLAGS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}--release-type${COLOR_NONE} <option>   Publish release destination [${COLOR_GREEN}options: github${COLOR_NONE}]"
  echo -e "  ${COLOR_LIGHT_CYAN}--release-tag${COLOR_NONE} <value>     Tag of the release"
  echo -e "  ${COLOR_LIGHT_CYAN}--config${COLOR_NONE} <value>          Path for goreleaser config file (default: ${COLOR_GREEN}${DEFAULT_GO_RELEASER_CONFIG_FILEPATH}${COLOR_NONE})"
  echo -e " "
  echo -e "${COLOR_WHITE}DELETE FLAGS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}--origin${COLOR_NONE} <option>         Origin of the Go binary/release to delete [${COLOR_GREEN}options: gobin-local/github${COLOR_NONE}]"
  echo -e "  ${COLOR_LIGHT_CYAN}--delete-tag${COLOR_NONE} <value>      Remote Go binary release tag to delete (github only)"
  echo -e " "  
  echo -e "${COLOR_WHITE}GENERAL FLAGS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}-y${COLOR_NONE} (--auto-prompt)        Do not prompt for approval and accept everything"
  echo -e "  ${COLOR_LIGHT_CYAN}-d${COLOR_NONE} (--dry-run)            Run all commands in dry-run mode without file system changes"
  echo -e "  ${COLOR_LIGHT_CYAN}-v${COLOR_NONE} (--verbose)            Output debug logs for commands executions"
  echo -e "  ${COLOR_LIGHT_CYAN}-s${COLOR_NONE} (--silent)             Do not output logs for commands executions"
  echo -e "  ${COLOR_LIGHT_CYAN}-h${COLOR_NONE} (--help)               Show available actions and their description"
  echo -e " "
  echo -e "${COLOR_WHITE}GLOBALS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}GITHUB_TOKEN${COLOR_NONE}              Valid GitHub token with write access for publishing releases"
  echo -e " "
  exit 0
}

parse_program_arguments() {
  if [ $# = 0 ]; then
    print_help_menu_and_exit "$0"
  fi

  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      build)
        CLI_ARGUMENT_BUILD="build"
        shift
        ;;
      install)
        CLI_ARGUMENT_INSTALL="install"
        shift
        ;;
      publish)
        CLI_ARGUMENT_PUBLISH="publish"
        shift
        ;;
      delete)
        CLI_ARGUMENT_DELETE="delete"
        shift
        ;;
      --main-package)
        CLI_FLAG_BUILD_MAIN_PACKAGE_PATH="main-package"
        shift
        CLI_VALUE_BUILD_MAIN_PACKAGE_PATH=$(cut -d ' ' -f 2- <<<"${1}" | xargs)
        shift
        ;;
      --output-path)
        CLI_FLAG_BUILD_OUTPUT_PATH="output-path"
        shift
        CLI_VALUE_BUILD_OUTPUT_PATH=$(cut -d ' ' -f 2- <<<"${1}" | xargs)
        shift
        ;;
      --os-arch)
        CLI_FLAG_BUILD_OS_ARCH_TUPLE="os-arch"
        shift
        os_arch_value=$(cut -d ' ' -f 2- <<<"${1}" | xargs)
        CLI_VALUE_BUILD_OS_ARCH_TUPLES_ARRAY+=("${os_arch_value}")
        shift
        ;;
      --archive)
        CLI_FLAG_BUILD_ARCHIVE="build-archive"
        shift
        ;;
      --release-type)
        CLI_FLAG_RELEASE_TYPE="release-type"
        shift
        CLI_VALUE_RELEASE_TYPE=$(cut -d ' ' -f 2- <<<"${1}" | xargs)
        shift
        ;;
      --release-tag)
        CLI_FLAG_RELEASE_TAG="release-tag"
        shift
        CLI_VALUE_RELEASE_TAG=$(cut -d ' ' -f 2- <<<"${1}" | xargs)
        shift
        ;;
      --config)
        CLI_FLAG_GO_RELEASER_CONFIG_FILEPATH="goreleaser-config"
        shift
        CLI_VALUE_GO_RELEASER_CONFIG_FILEPATH=$(cut -d ' ' -f 2- <<<"${1}" | xargs)
        shift
        ;;
      --origin)
        CLI_FLAG_DELETE_ORIGIN="origin"
        shift
        CLI_VALUE_DELETE_ORIGIN=$(cut -d ' ' -f 2- <<<"${1}" | xargs)
        shift
        ;;
      --delete-tag)
        CLI_FLAG_DELETE_TAG="delete-tag"
        shift
        CLI_VALUE_DELETE_TAG=$(cut -d ' ' -f 2- <<<"${1}" | xargs)
        shift
        ;;
      -d | --dry-run)
        # Used by logger.sh
        export LOGGER_DRY_RUN="true"
        shift
        ;;
      -y | --auto-prompt)
        # Used by prompter.sh
        export PROMPTER_SKIP_PROMPT="y"
        shift
        ;;
      -v | --verbose)
        # Used by logger.sh
        export LOGGER_VERBOSE="true"
        shift
        ;;
      -s | --silent)
        # Used by logger.sh
        export LOGGER_SILENT="true"
        shift
        ;;
      -h | --help)
        print_help_menu_and_exit "$0"
        ;;
      *)
        log_fatal "Unknown option $1 (did you mean =$1 ?)"
        ;;
    esac
  done
}

check_legal_arguments() {
  ! is_build && ! is_install && ! is_publish && ! is_delete
}

check_delete_invalid_origin() {
  is_delete && ! is_gobin_local_delete_origin && ! is_github_delete_origin
}

check_publish_release_type() {
  if is_publish; then
    if is_github_release_type && [[ -z "${GITHUB_TOKEN}" ]]; then
      log_fatal "Publish command is missing an authentication token. name: GITHUB_TOKEN"
    fi
    if ! is_github_release_type; then
      log_fatal "Publish command has an invalid release type value. options: github"
    fi
  fi
}

check_publish_release() {
  if is_publish; then
    if [[ -z "${CLI_VALUE_RELEASE_TYPE}" ]]; then
      log_fatal "Publish command is missing a release tag. flag: --release-type"
    fi
    if [[ -z "${CLI_VALUE_RELEASE_TAG}" ]]; then
      log_fatal "Publish command is missing a release tag. flag: --release-tag"
    fi
  fi
}

# Currently redundant since a user will get prompted is a delete tag is missing
# check_delete_missing_tag() {
#   if is_delete; then
#     if is_github_delete_origin && [[ -z "${CLI_VALUE_DELETE_TAG}" ]]; then
#       log_fatal "Delete command from GitHub requires a delete tag. flag: --delete-tag"
#     fi
#   fi
# }

check_unsupported_actions() {
  if is_build_archive; then
    log_fatal "Archiving a build result is not supported yet"
  fi
}

verify_program_arguments() {
  if check_legal_arguments; then
    log_fatal "Missing mandatory command argument. Options: build/install/publish/delete"
  fi
  if check_delete_invalid_origin; then
    log_fatal "Command argument 'delete' is missing a mandatory flag value or has an invalid value. flag: --origin, options: gobin-local/github"
  fi
  check_publish_release
  check_publish_release_type
  # check_delete_missing_tag
  check_unsupported_actions
  evaluate_dry_run_mode
}

prerequisites() {
  check_tool "go"
}

get_project_name() {
  basename "$(pwd)"
}

main() {
  parse_program_arguments "$@"
  verify_program_arguments

  prerequisites

  if is_build; then
    build_binaries
  fi

  if is_install; then
    build_binaries
    install_binary
  fi

  if is_delete; then
    delete_binary_or_release
  fi

  if is_publish; then
    build_binaries
    publish_binaries
  fi
}

main "$@"
