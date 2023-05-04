#!/bin/bash

CURRENT_FOLDER_ABS_PATH=$(dirname "${BASH_SOURCE[0]}")
ROOT_FOLDER_ABS_PATH=$(dirname "${CURRENT_FOLDER_ABS_PATH}")

source "${ROOT_FOLDER_ABS_PATH}/logger.sh"
source "${ROOT_FOLDER_ABS_PATH}/io.sh"
source "${ROOT_FOLDER_ABS_PATH}/cmd.sh"
source "${ROOT_FOLDER_ABS_PATH}/strings.sh"
source "${ROOT_FOLDER_ABS_PATH}/checks.sh"
source "${ROOT_FOLDER_ABS_PATH}/github.sh"

SCRIPT_MENU_TITLE="Go Development Commands"

DEFAULT_TESTS_PATH="./..."
DEFAULT_TEST_COVERAGE_REPORT_FILE="coverage.out"
DEFAULT_CONTAINERIZED_TESTS_GO_VERSION="1.20"
DEFAULT_CONTAINERIZED_TESTS_TPARSE_VERSION="v0.10.0"

CLI_ARGUMENT_DEPS=""
CLI_ARGUMENT_FMT=""
CLI_ARGUMENT_TEST=""
CLI_ARGUMENT_DOCS=""

CLI_FLAG_FMT_CHECK_ONLY=""                 # true/false if missing
CLI_FLAG_TESTS_PATH=""                     
CLI_FLAG_TEST_CONTAINERIZED=""              # true/false if missing
CLI_FLAG_TEST_CONTAINERIZED_GO_VERSION=""   
CLI_FLAG_TEST_DENSE_MODE=""                # true/false if missing
CLI_FLAG_TEST_COVERAGE=""                  # true/false if missing
CLI_FLAG_DOCS_LAN=""                       # true/false if missing

CLI_VALUE_TESTS_PATH=""
CLI_VALUE_TEST_CONTAINERIZE_GO_VERSION=""

is_deps() {
  [[ -n "${CLI_ARGUMENT_DEPS}" ]]
}

is_fmt() {
  [[ -n "${CLI_ARGUMENT_FMT}" ]]
}

is_fmt_check_only() {
  [[ -n "${CLI_FLAG_FMT_CHECK_ONLY}" ]]
}

is_test() {
  [[ -n "${CLI_ARGUMENT_TEST}" ]]
}

is_test_generate_coverage() {
  [[ -n "${CLI_FLAG_TEST_COVERAGE}" ]]
}

is_test_with_custom_test_path() {
  [[ -n "${CLI_FLAG_TESTS_PATH}" ]]
}

get_test_path() {
  echo "${CLI_VALUE_TESTS_PATH:-${DEFAULT_TESTS_PATH}}"
}

is_test_containerize() {
  [[ -n "${CLI_FLAG_TEST_CONTAINERIZED}" ]]
}

get_test_containerize_go_version() {
  echo "${CLI_VALUE_TEST_CONTAINERIZE_GO_VERSION:-${DEFAULT_CONTAINERIZED_TESTS_GO_VERSION}}"
}

is_test_dense_mode() {
  [[ -n "${CLI_FLAG_TEST_DENSE_MODE}" ]]
}

is_docs() {
  [[ -n "${CLI_ARGUMENT_DOCS}" ]]
}

is_docs_lan() {
  [[ -n "${CLI_FLAG_DOCS_LAN}" ]]
}

vendor_go_deps() {
  log_info "Tidying dependencies"
  cmd_run "go mod tidy"
  
  log_info "Verifying dependencies"
  new_line
  cmd_run "go mod verify"
  new_line

  log_info "Vendoring dependencies"
  cmd_run "go mod vendor"

  log_info "All cleaned up !"
}

report_on_format_errors() {
  log_info "Checking that Go code complies with gofmt requirements..."
  local gofmt_files="<List of go files with format errors>"
  if ! is_dry_run; then
    gofmt_files=$(gofmt -l $(find . -name '*.go' | grep -v vendor))
  fi
  if is_verbose; then
    new_line
    echo "  gofmt -l \$(find . -name '*.go' | grep -v vendor)"
    new_line
  fi
  if [[ -n ${gofmt_files} ]]; then
    log_info 'Found go files with formatting errors:'
    new_line
    echo "${gofmt_files}"
    new_line
    log_fatal "Please run fmt to reformat the code."
  else
    log_info "Go code is well formatted !"
  fi
}

format_go_sources() {
  log_info "Formatting Go source code to comply with gofmt requirements..."
  local gofmt_files="<List of go files with format errors>"
  if ! is_dry_run; then
    gofmt_files=$(gofmt -l $(find . -name '*.go' | grep -v vendor))
  fi
  if is_verbose; then
    new_line
    echo "  gofmt -l \$(find . -name '*.go' | grep -v vendor)"
  fi
  if [[ -n ${gofmt_files} ]]; then
    cmd_run "gofmt -w ${gofmt_files}"
    log_info "List of formatted go source files:"
    new_line
    echo "${gofmt_files}"
  else
    log_info "Go code is well formatted !"
  fi
}

maybe_format_go_sources() {
  check_tool "gofmt"
  if is_fmt_check_only; then
    report_on_format_errors
  else
    format_go_sources
  fi
}

run_tests_on_host() {
  local tests_path=$(get_test_path)
  log_info "Runing tests suite on: ${COLOR_YELLOW}HOST MACHINE${COLOR_NONE}"

  local maybe_pipe=""
  if is_tool_exist tparse; then
    if is_test_dense_mode; then
      maybe_pipe=" | tparse -all -notests -smallscreen"
    else
      maybe_pipe=" | tparse -all"
    fi
  fi

  local maybe_coverage=""
  if is_test_generate_coverage; then
    # -coverprofile=coverage.out was added for GitHub workflow integration with 
    # jandelgado/gcov2lcov-action
    maybe_coverage="-cover -covermode=count -coverprofile=coverage.out.temp"
  fi

  new_line
  cmd_run "go test -v ${tests_path} -json ${maybe_coverage} ${maybe_pipe}"
  new_line

  if is_test_generate_coverage && is_file_exist "coverage.out.temp"; then
    log_info "Preparing tests coverage report"
    cmd_run "cat coverage.out.temp | grep -v '_testkit\|_fakes' >${DEFAULT_TEST_COVERAGE_REPORT_FILE}"
    if is_file_exist "${DEFAULT_TEST_COVERAGE_REPORT_FILE}"; then
      local cov=$(cmd_run "go tool cover -func ${DEFAULT_TEST_COVERAGE_REPORT_FILE} | grep total | awk '{print \$3}'")
      new_line
      log_info "TESTS COVERAGE: ${cov}"
    fi
  fi
}

run_tests_containerized() {
  local go_version=$(get_test_containerize_go_version)
  local cwd=$(pwd)
  local working_dir_name="/home/go_tests_suite_dir"

  log_info "Runing tests suite on: ${COLOR_YELLOW}Containerized Environment${COLOR_NONE}. Go: v${go_version}, tparse: ${DEFAULT_CONTAINERIZED_TESTS_TPARSE_VERSION}"

  cmd_run """
docker run -it \
    -v ${cwd}:${working_dir_name} \
    --entrypoint /bin/sh golang:${go_version} \
    -c '
go install github.com/mfridman/tparse@${DEFAULT_CONTAINERIZED_TESTS_TPARSE_VERSION}
cd ${working_dir_name} || exit
make test '
"""
}

run_tests_suite() {
  if is_test_containerize; then
    run_tests_containerized
  else
    run_tests_on_host
  fi
}

start_local_docs_site() {
  check_tool "npm"
  check_tool "hugo"
  if ! is_dry_run; then
    cd docs-site || exit
  fi
  if is_docs_lan; then
    log_info "Running a local docs site opened for LAN access (http://192.168.x.xx:9001)"
    new_line
    cmd_run "npm run docs-serve-lan"
  else
    log_info "Running a local docs site (http://localhost:9001/anchor/)"
    new_line
    cmd_run "npm run docs-serve"
  fi
}

print_help_menu_and_exit() {
  local exec_filename=$1
  local file_name=$(basename "${exec_filename}")
  echo -e ""
  echo -e "${SCRIPT_MENU_TITLE} - Run common development commands on a Go project"
  echo -e " "
  echo -e "${COLOR_WHITE}USAGE${COLOR_NONE}"
  echo -e "  "${file_name}" [command] [option] [flag]"
  echo -e " "
  echo -e "${COLOR_WHITE}ARGUMENTS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}deps${COLOR_NONE}                      Tidy, verify and vendor go dependencies"
  echo -e "  ${COLOR_LIGHT_CYAN}fmt${COLOR_NONE}                       Format Go code using gofmt style and sort imports"
  echo -e "  ${COLOR_LIGHT_CYAN}test${COLOR_NONE}                      Run tests suite"
  echo -e "  ${COLOR_LIGHT_CYAN}docs${COLOR_NONE}                      Run a local documentation site (http://localhost:9001/<project-name>/)"
  echo -e " "
  echo -e "${COLOR_WHITE}FORMAT FLAGS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}--check-only${COLOR_NONE}              Only validate Go code format and imports"
  echo -e " "
  echo -e "${COLOR_WHITE}TEST FLAGS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}--tests-path${COLOR_NONE} <path>       Tests folder path (default: ${COLOR_GREEN}'./...'${COLOR_NONE})"
  echo -e "  ${COLOR_LIGHT_CYAN}--dense-mode${COLOR_NONE}              Output tests suite in dense mode for small screens"
  echo -e "  ${COLOR_LIGHT_CYAN}--coverage${COLOR_NONE}                Generate tests coverage report output (path: ./${COLOR_GREEN}${DEFAULT_TEST_COVERAGE_REPORT_FILE}${COLOR_NONE})"
  echo -e "  ${COLOR_LIGHT_CYAN}--containerized${COLOR_NONE}           Run tests suite within a Docker container"
  echo -e "  ${COLOR_LIGHT_CYAN}--go-version${COLOR_NONE}              Go version for the containerized tests suite run (default: ${COLOR_GREEN}v1.20${COLOR_NONE})"
  echo -e " "
  echo -e "${COLOR_WHITE}DOCS FLAGS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}--lan${COLOR_NONE}                     Make the documentation site avaialble within LAN (http://192.168.x.xx:9001/)"
  echo -e " "  
  echo -e "${COLOR_WHITE}GENERAL FLAGS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}-y${COLOR_NONE} (--auto-prompt)        Do not prompt for approval and accept everything"
  echo -e "  ${COLOR_LIGHT_CYAN}-d${COLOR_NONE} (--dry-run)            Run all commands in dry-run mode without file system changes"
  echo -e "  ${COLOR_LIGHT_CYAN}-v${COLOR_NONE} (--verbose)            Output debug logs for commands executions"
  echo -e "  ${COLOR_LIGHT_CYAN}-s${COLOR_NONE} (--silent)             Do not output logs for commands executions"
  echo -e "  ${COLOR_LIGHT_CYAN}-h${COLOR_NONE} (--help)               Show available actions and their description"
  echo -e " "
  echo -e "${COLOR_WHITE}GLOBALS${COLOR_NONE}"
  echo -e "  ${COLOR_LIGHT_CYAN}PYPI_TOKEN${COLOR_NONE}                Valid PyPI token with write access for publishing releases"
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
      deps)
        CLI_ARGUMENT_DEPS="deps"
        shift
        ;;
      fmt)
        CLI_ARGUMENT_FMT="fmt"
        shift
        ;;
      test)
        CLI_ARGUMENT_TEST="test"
        shift
        ;;
      docs)
        CLI_ARGUMENT_DOCS="docs"
        shift
        ;;
      --check-only)
        CLI_FLAG_FMT_CHECK_ONLY="check-only"
        shift
        ;;
      --tests-path)
        CLI_FLAG_TESTS_PATH="test"
        shift
        CLI_VALUE_TESTS_PATH=$(cut -d ' ' -f 2- <<<"${1}" | xargs)
        shift
        ;;
      --containerized)
        CLI_FLAG_TEST_CONTAINERIZED="true"
        shift
        ;;
      --coverage)
        CLI_FLAG_TEST_COVERAGE="true"
        shift
        ;;
      --go-version)
        CLI_FLAG_TEST_CONTAINERIZED_GO_VERSION="go-version"
        shift
        CLI_VALUE_TEST_CONTAINERIZE_GO_VERSION=$(cut -d ' ' -f 2- <<<"${1}" | xargs)
        shift
        ;;
      --dense-mode)
        CLI_FLAG_TEST_DENSE_MODE="true"
        shift
        ;;
      --lan)
        CLI_FLAG_DOCS_LAN="true"
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
  if ! is_deps && ! is_fmt && ! is_test && ! is_docs; then
    log_fatal "No legal command could be found"
  fi
}

check_tests_path_has_value() {
  if is_test_with_custom_test_path && [[ -z "${CLI_VALUE_TESTS_PATH}" ]]; then
    log_fatal "Tests path flag is missing a folder path. flag: --tests-path"
  fi
}

verify_program_arguments() {
  check_legal_arguments
  check_tests_path_has_value
  evaluate_dry_run_mode
}

prerequisites() {
  check_tool "go"
}

main() {
  parse_program_arguments "$@"
  verify_program_arguments

  prerequisites

  if is_deps; then
    vendor_go_deps
  fi

  if is_fmt; then
    maybe_format_go_sources
  fi

  if is_test; then
    run_tests_suite
  fi

  if is_docs; then
    start_local_docs_site
  fi
}

main "$@"
