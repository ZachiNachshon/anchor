package config

import (
	"github.com/ZachiNachshon/anchor/pkg/utils/installer"
	"github.com/ZachiNachshon/anchor/pkg/utils/locator"
	"os"

	"github.com/ZachiNachshon/anchor/pkg/common"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/shell"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func CheckPrerequisites() error {
	var repoPath = ""
	if repoPath = os.Getenv("DOCKER_FILES"); len(repoPath) <= 0 {
		return errors.Errorf("DOCKER_FILES environment variable is missing, must contain path to 'dockerfiles' git repository.")
	}
	common.GlobalOptions.DockerRepositoryPath = repoPath

	// TODO: resolve shell type from configuration (https://github.com/spf13/viper ?)
	common.ShellExec = shell.NewShellExecutor(shell.BASH)

	setDefaultEnvVar()
	LoadEnvVars(common.GlobalOptions.DockerRepositoryPath)

	if err := installer.NewGoInstaller(common.ShellExec).Check(); err != nil {
		return err
	}

	return nil
}

func setDefaultEnvVar() {
	// Docker
	_ = os.Setenv("REGISTRY", common.GlobalOptions.DockerRegistryDnsWithIp)
	_ = os.Setenv("NAMESPACE", common.GlobalOptions.DockerImageNamespace)
	_ = os.Setenv("TAG", common.GlobalOptions.DockerImageTag)
}

func LoadEnvVars(identifier string) {
	var envFilePath = ""

	if identifier == common.GlobalOptions.DockerRepositoryPath {
		// dirname should be the repository root folder only at config CheckPrerequisites stage
		envFilePath = common.GlobalOptions.DockerRepositoryPath + "/.env"
	} else {
		ctxDir, _ := locator.DirLocator.DockerContext(identifier)
		envFilePath = ctxDir + "/.env"
	}
	loadEnvVarsInner(envFilePath)
}

func loadEnvVarsInner(envFilePath string) {
	if err := godotenv.Overload(envFilePath); err != nil {
		if common.GlobalOptions.Verbose {
			// TODO: Change to warn once implemented
			logger.Info(err.Error())
		}
	}

	if v := os.Getenv("NAMESPACE"); len(v) > 0 {
		common.GlobalOptions.DockerImageNamespace = v
	}

	if v := os.Getenv("TAG"); len(v) > 0 {
		common.GlobalOptions.DockerImageTag = v
	}
}

const KubernetesNamespaceManifest = `
apiVersion: v1
kind: Namespace
metadata:
  name: NAMESPACE-TO-REPLACE
`

const RegistryContainerdConfigTemplate = `disabled_plugins = ["aufs", "btrfs", "zfs"]
root = "/var/lib/containerd"
state = "/run/containerd"
oom_score = 0

[grpc]
  address = "/run/containerd/containerd.sock"
  uid = 0
  gid = 0
  max_recv_message_size = 16777216
  max_send_message_size = 16777216

[debug]
  address = ""
  uid = 0
  gid = 0
  level = ""

[metrics]
  address = ""
  grpc_histogram = false

[cgroup]
  path = ""

[plugins]
  [plugins.cgroups]
    no_prometheus = false
  [plugins.cri]
    stream_server_address = "127.0.0.1"
    stream_server_port = "0"
    enable_selinux = false
    sandbox_image = "k8s.gcr.io/pause:3.1"
    stats_collect_period = 10
    systemd_cgroup = false
    enable_tls_streaming = false
    max_container_log_line_size = 16384
    [plugins.cri.containerd]
      snapshotter = "overlayfs"
      no_pivot = false
      [plugins.cri.containerd.default_runtime]
        runtime_type = "io.containerd.runtime.v1.linux"
        runtime_engine = ""
        runtime_root = ""
      [plugins.cri.containerd.untrusted_workload_runtime]
        runtime_type = ""
        runtime_engine = ""
        runtime_root = ""
    [plugins.cri.cni]
      bin_dir = "/opt/cni/bin"
      conf_dir = "/etc/cni/net.d"
      conf_template = ""
    [plugins.cri.registry]
      [plugins.cri.registry.mirrors]
        [plugins.cri.registry.mirrors."local.insecure-registry.io"]
          endpoint = ["http://127.0.0.1:32001"]
        [plugins.cri.registry.mirrors."registry.anchor:32001"]
          endpoint = ["http://127.0.0.1:32001"]
        [plugins.cri.registry.mirrors."docker.io"]
          endpoint = ["https://registry-1.docker.io"]
    [plugins.cri.x509_key_pair_streaming]
      tls_cert_file = ""
      tls_key_file = ""
  [plugins.diff-service]
    default = ["walking"]
  [plugins.linux]
    shim = "containerd-shim"
    runtime = "runc"
    runtime_root = ""
    no_shim = false
    shim_debug = false
  [plugins.opt]
    path = "/opt/containerd"
  [plugins.restart]
    interval = "10s"
  [plugins.scheduler]
    pause_threshold = 0.02
    deletion_threshold = 0
    mutation_threshold = 100
    schedule_delay = "0s"
    startup_delay = "100ms"
`

const AutoCompletionFuncBash = `# bash completion for anchor                               -*- shell-script -*-

__anchor_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__anchor_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__anchor_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__anchor_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__anchor_handle_reply()
{
    __anchor_debug "${FUNCNAME[0]}"
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            COMPREPLY=( $(compgen -W "${allflags[*]}" -- "$cur") )
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%=*}"
                __anchor_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __anchor_index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions=("${must_have_one_noun[@]}")
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    COMPREPLY=( $(compgen -W "${completions[*]}" -- "$cur") )

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        COMPREPLY=( $(compgen -W "${noun_aliases[*]}" -- "$cur") )
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
		if declare -F __anchor_custom_func >/dev/null; then
			# try command name qualified custom func
			__anchor_custom_func
		else
			# otherwise fall back to unqualified for compatibility
			declare -F __custom_func >/dev/null && __custom_func
		fi
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__anchor_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__anchor_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1
}

__anchor_handle_flag()
{
    __anchor_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __anchor_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __anchor_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __anchor_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if [[ ${words[c]} != *"="* ]] && __anchor_contains_word "${words[c]}" "${two_word_flags[@]}"; then
			  __anchor_debug "${FUNCNAME[0]}: found a flag ${words[c]}, skip the next argument"
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__anchor_handle_noun()
{
    __anchor_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __anchor_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __anchor_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__anchor_handle_command()
{
    __anchor_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_anchor_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __anchor_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__anchor_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __anchor_handle_reply
        return
    fi
    __anchor_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __anchor_handle_flag
    elif __anchor_contains_word "${words[c]}" "${commands[@]}"; then
        __anchor_handle_command
    elif [[ $c -eq 0 ]]; then
        __anchor_handle_command
    elif __anchor_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __anchor_handle_command
        else
            __anchor_handle_noun
        fi
    else
        __anchor_handle_noun
    fi
    __anchor_handle_word
}

#compdef _anchor anchor


function _anchor {
  local -a commands

  _arguments -C \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "docker:Docker commands"
      "help:Help about any command"
      "kubernetes:Kubernetes commands"
      "list:List all supported directories from DOCKER_FILES repository"
      "version:Print anchor version"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  docker)
    _anchor_docker
    ;;
  help)
    _anchor_help
    ;;
  kubernetes)
    _anchor_kubernetes
    ;;
  list)
    _anchor_list
    ;;
  version)
    _anchor_version
    ;;
  esac
}


function _anchor_docker {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for docker]' \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "build:Builds a docker image"
      "clean:Clean docker containers and images"
      "purge:Purge all docker images and containers"
      "push:Push a docker image to repository [registry.anchor:32001]"
      "run:Run a docker container"
      "stop:Stop a docker container"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  build)
    _anchor_docker_build
    ;;
  clean)
    _anchor_docker_clean
    ;;
  purge)
    _anchor_docker_purge
    ;;
  push)
    _anchor_docker_push
    ;;
  run)
    _anchor_docker_run
    ;;
  stop)
    _anchor_docker_stop
    ;;
  esac
}

function _anchor_docker_build {
  _arguments \
    '(-t --Docker image tag)'{-t,--Docker image tag}'[anchor docker build <name> -t my_tag]:' \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_docker_clean {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_docker_purge {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_docker_push {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_docker_run {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_docker_stop {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_help {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}


function _anchor_kubernetes {
  local -a commands

  _arguments -C \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "connect:Connect to a kubernetes pod by name"
      "create:Create a local Kubernetes cluster"
      "dashboard:Deploy a Kubernetes dashboard"
      "deploy:Deploy a container Kubernetes manifest"
      "destroy:Destroy local Kubernetes cluster"
      "expose:Expose a container port to the host instance"
      "registry:Create a private docker registry [registry.anchor:32001]"
      "remove:Removed a previously deployed container manifest"
      "status:Print cluster [registry.anchor:32001] status"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  connect)
    _anchor_kubernetes_connect
    ;;
  create)
    _anchor_kubernetes_create
    ;;
  dashboard)
    _anchor_kubernetes_dashboard
    ;;
  deploy)
    _anchor_kubernetes_deploy
    ;;
  destroy)
    _anchor_kubernetes_destroy
    ;;
  expose)
    _anchor_kubernetes_expose
    ;;
  registry)
    _anchor_kubernetes_registry
    ;;
  remove)
    _anchor_kubernetes_remove
    ;;
  status)
    _anchor_kubernetes_status
    ;;
  esac
}

function _anchor_kubernetes_connect {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_create {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_dashboard {
  _arguments \
    '(-d --Delete Kubernetes dashboard)'{-d,--Delete Kubernetes dashboard}'[anchor cluster dashboard -d]' \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_deploy {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_destroy {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_expose {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_registry {
  _arguments \
    '(-d --Delete Kubernetes docker registry as a pod)'{-d,--Delete Kubernetes docker registry as a pod}'[anchor cluster registry -d]' \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_remove {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_status {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_list {
  _arguments \
    '(-a --filter by affinity)'{-a,--filter by affinity}'[anchor list -a affinity-name]:' \
    '(-k --filter kubernetes manifests only)'{-k,--filter kubernetes manifests only}'[anchor list -k]' \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_version {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}


_anchor_docker_build()
{
    last_command="anchor_docker_build"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--Docker image tag=")
    two_word_flags+=("--Docker image tag")
    two_word_flags+=("-t")
    local_nonpersistent_flags+=("--Docker image tag=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_docker_clean()
{
    last_command="anchor_docker_clean"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_docker_purge()
{
    last_command="anchor_docker_purge"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_docker_push()
{
    last_command="anchor_docker_push"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_docker_run()
{
    last_command="anchor_docker_run"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_docker_stop()
{
    last_command="anchor_docker_stop"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_docker()
{
    last_command="anchor_docker"

    command_aliases=()

    commands=()
    commands+=("build")
    commands+=("clean")
    commands+=("purge")
    commands+=("push")
    commands+=("run")
    commands+=("stop")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("build")
    must_have_one_noun+=("clean")
    must_have_one_noun+=("list")
    must_have_one_noun+=("purge")
    must_have_one_noun+=("push")
    must_have_one_noun+=("run")
    must_have_one_noun+=("stop")
    noun_aliases=()
}

_anchor_kubernetes_connect()
{
    last_command="anchor_kubernetes_connect"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_kubernetes_create()
{
    last_command="anchor_kubernetes_create"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_kubernetes_dashboard()
{
    last_command="anchor_kubernetes_dashboard"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--Delete Kubernetes dashboard")
    flags+=("-d")
    local_nonpersistent_flags+=("--Delete Kubernetes dashboard")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_kubernetes_deploy()
{
    last_command="anchor_kubernetes_deploy"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_kubernetes_destroy()
{
    last_command="anchor_kubernetes_destroy"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_kubernetes_expose()
{
    last_command="anchor_kubernetes_expose"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_kubernetes_registry()
{
    last_command="anchor_kubernetes_registry"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--Delete Kubernetes docker registry as a pod")
    flags+=("-d")
    local_nonpersistent_flags+=("--Delete Kubernetes docker registry as a pod")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_kubernetes_remove()
{
    last_command="anchor_kubernetes_remove"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_kubernetes_status()
{
    last_command="anchor_kubernetes_status"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_kubernetes()
{
    last_command="anchor_kubernetes"

    command_aliases=()

    commands=()
    commands+=("connect")
    commands+=("create")
    commands+=("dashboard")
    commands+=("deploy")
    commands+=("destroy")
    commands+=("expose")
    commands+=("registry")
    commands+=("remove")
    commands+=("status")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_list()
{
    last_command="anchor_list"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--filter by affinity=")
    two_word_flags+=("--filter by affinity")
    two_word_flags+=("-a")
    flags+=("--filter kubernetes manifests only")
    flags+=("-k")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_version()
{
    last_command="anchor_version"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_anchor_root_command()
{
    last_command="anchor"

    command_aliases=()

    commands=()
    commands+=("docker")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("d")
        aliashash["d"]="docker"
    fi
    commands+=("kubernetes")
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        command_aliases+=("k")
        aliashash["k"]="kubernetes"
    fi
    commands+=("list")
    commands+=("version")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("docker")
    must_have_one_noun+=("kubernetes")
    must_have_one_noun+=("list")
    noun_aliases=()
}

__start_anchor()
{
    local cur prev words cword
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __anchor_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("anchor")
    local must_have_one_flag=()
    local must_have_one_noun=()
    local last_command
    local nouns=()

    __anchor_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_anchor anchor
else
    complete -o default -o nospace -F __start_anchor anchor
fi

# ex: ts=4 sw=4 et filetype=sh
`

const AutoCompletionFuncZsh = `#compdef _anchor anchor


function _anchor {
  local -a commands

  _arguments -C \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "completion:Generate auto completion script for bash/zsh"
      "docker:Docker commands"
      "help:Help about any command"
      "kubernetes:Kubernetes commands"
      "list:List all supported directories from DOCKER_FILES repository"
      "version:Print anchor version"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  completion)
    _anchor_completion
    ;;
  docker)
    _anchor_docker
    ;;
  help)
    _anchor_help
    ;;
  kubernetes)
    _anchor_kubernetes
    ;;
  list)
    _anchor_list
    ;;
  version)
    _anchor_version
    ;;
  esac
}

function _anchor_completion {
  _arguments \
    '(-f --file)'{-f,--file}'[anchor completion -f]' \
    '(-h --help)'{-h,--help}'[help for completion]' \
    '(-s --shell)'{-s,--shell}'[anchor completion -s bash/zsh]:' \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}


function _anchor_docker {
  local -a commands

  _arguments -C \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "build:Builds a docker image"
      "purge:Purge all docker images and containers"
      "push:Push a docker image to repository [registry.anchor]"
      "remove:Remove docker containers and images"
      "run:Run a docker container"
      "stop:Stop a docker container"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  build)
    _anchor_docker_build
    ;;
  purge)
    _anchor_docker_purge
    ;;
  push)
    _anchor_docker_push
    ;;
  remove)
    _anchor_docker_remove
    ;;
  run)
    _anchor_docker_run
    ;;
  stop)
    _anchor_docker_stop
    ;;
  esac
}

function _anchor_docker_build {
  _arguments \
    '(-t --tag)'{-t,--tag}'[anchor docker build <name> -t my_tag]:' \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_docker_purge {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_docker_push {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_docker_remove {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_docker_run {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_docker_stop {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_help {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}


function _anchor_kubernetes {
  local -a commands

  _arguments -C \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "connect:Connect to a kubernetes pod by name"
      "create:Create a local Kubernetes cluster"
      "dashboard:Deploy a Kubernetes dashboard"
      "deploy:Deploy a container Kubernetes manifest"
      "destroy:Destroy local Kubernetes cluster"
      "expose:Expose a container port to the host instance"
      "registry:Create a private docker registry [registry.anchor]"
      "remove:Removed a previously deployed container manifest"
      "status:Print cluster [anchor] status"
      "token:Generate export KUBECONFIG command and load to clipboard"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  connect)
    _anchor_kubernetes_connect
    ;;
  create)
    _anchor_kubernetes_create
    ;;
  dashboard)
    _anchor_kubernetes_dashboard
    ;;
  deploy)
    _anchor_kubernetes_deploy
    ;;
  destroy)
    _anchor_kubernetes_destroy
    ;;
  expose)
    _anchor_kubernetes_expose
    ;;
  registry)
    _anchor_kubernetes_registry
    ;;
  remove)
    _anchor_kubernetes_remove
    ;;
  status)
    _anchor_kubernetes_status
    ;;
  token)
    _anchor_kubernetes_token
    ;;
  esac
}

function _anchor_kubernetes_connect {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_create {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_dashboard {
  _arguments \
    '(-d --Delete Kubernetes dashboard)'{-d,--Delete Kubernetes dashboard}'[anchor cluster dashboard -d]' \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_deploy {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_destroy {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_expose {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_registry {
  _arguments \
    '(-d --delete)'{-d,--delete}'[anchor cluster registry -d]' \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_remove {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_status {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_kubernetes_token {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_list {
  _arguments \
    '(-a --filter by affinity)'{-a,--filter by affinity}'[anchor list -a affinity-name]:' \
    '(-k --filter kubernetes manifests only)'{-k,--filter kubernetes manifests only}'[anchor list -k]' \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

function _anchor_version {
  _arguments \
    '(-v --verbose)'{-v,--verbose}'[anchor <command> -v]'
}

`

const KubernetesDashboardManifest = `# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ------------------- Dashboard Secret ------------------- #

apiVersion: v1
kind: Secret
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard-certs
  namespace: kube-system
type: Opaque

---
# ------------------- Dashboard Service Account ------------------- #

apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kube-system

---
# ------------------- Dashboard Role & Role Binding ------------------- #

kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kubernetes-dashboard-minimal
  namespace: kube-system
rules:
  # Allow Dashboard to create 'kubernetes-dashboard-key-holder' secret.
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["create"]
    # Allow Dashboard to create 'kubernetes-dashboard-settings' config map.
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["create"]
    # Allow Dashboard to get, update and delete Dashboard exclusive secrets.
  - apiGroups: [""]
    resources: ["secrets"]
    resourceNames: ["kubernetes-dashboard-key-holder", "kubernetes-dashboard-certs"]
    verbs: ["get", "update", "delete"]
    # Allow Dashboard to get and update 'kubernetes-dashboard-settings' config map.
  - apiGroups: [""]
    resources: ["configmaps"]
    resourceNames: ["kubernetes-dashboard-settings"]
    verbs: ["get", "update"]
    # Allow Dashboard to get metrics from heapster.
  - apiGroups: [""]
    resources: ["services"]
    resourceNames: ["heapster"]
    verbs: ["proxy"]
  - apiGroups: [""]
    resources: ["services/proxy"]
    resourceNames: ["heapster", "http:heapster:", "https:heapster:"]
    verbs: ["get"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: kubernetes-dashboard-minimal
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: kubernetes-dashboard-minimal
subjects:
  - kind: ServiceAccount
    name: kubernetes-dashboard
    namespace: kube-system

---
# ------------------- Dashboard Deployment ------------------- #

kind: Deployment
apiVersion: apps/v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kube-system
spec:
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      k8s-app: kubernetes-dashboard
  template:
    metadata:
      labels:
        k8s-app: kubernetes-dashboard
    spec:
      containers:
        - name: kubernetes-dashboard
          image: k8s.gcr.io/kubernetes-dashboard-amd64:v1.10.1
          ports:
            - containerPort: 8443
              protocol: TCP
          args:
            - --auto-generate-certificates
            # Uncomment the following line to manually specify Kubernetes API server Host
            # If not specified, Dashboard will attempt to auto discover the API server and connect
            # to it. Uncomment only if the default does not work.
            # - --apiserver-host=http://my-address:port
          volumeMounts:
            - name: kubernetes-dashboard-certs
              mountPath: /certs
              # Create on-disk volume to store exec logs
            - mountPath: /tmp
              name: tmp-volume
          livenessProbe:
            httpGet:
              scheme: HTTPS
              path: /
              port: 8443
            initialDelaySeconds: 30
            timeoutSeconds: 30
      volumes:
        - name: kubernetes-dashboard-certs
          secret:
            secretName: kubernetes-dashboard-certs
        - name: tmp-volume
          emptyDir: {}
      serviceAccountName: kubernetes-dashboard
      # Comment the following tolerations if Dashboard must not be deployed on master
      tolerations:
        - key: node-role.kubernetes.io/master
          effect: NoSchedule

---
# ------------------- Dashboard Service ------------------- #

kind: Service
apiVersion: v1
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kube-system
spec:
  ports:
    - port: 443
      targetPort: 8443
  selector:
    k8s-app: kubernetes-dashboard`

const KubernetesRegistryManifest = `---
apiVersion: v1
kind: Namespace
metadata:
  name: container-registry
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: registry-claim
  namespace: container-registry
spec:
  accessModes:
    - ReadWriteMany
  volumeMode: Filesystem
  resources:
    requests:
      storage: 5Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: registry
  name: registry
  namespace: container-registry
spec:
  replicas: 1
  selector:
    matchLabels:
      app: registry
  template:
    metadata:
      labels:
        app: registry
    spec:
      containers:
        - name: registry
          image: cdkbot/registry-amd64:2.6
          env:
            - name: REGISTRY_HTTP_ADDR
              value: :5000
            - name: REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY
              value: /var/lib/registry
            - name: REGISTRY_STORAGE_DELETE_ENABLED
              value: "yes"
          ports:
            - containerPort: 5000
              name: registry
              protocol: TCP
          volumeMounts:
            - mountPath: /var/lib/registry
              name: registry-data
      volumes:
        - name: registry-data
          persistentVolumeClaim:
            claimName: registry-claim
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: registry
  name: registry
  namespace: container-registry
spec:
  type: NodePort
  selector:
    app: registry
  ports:
    - name: "registry"
      port: 5000
      targetPort: 5000
      nodePort: 32001`
