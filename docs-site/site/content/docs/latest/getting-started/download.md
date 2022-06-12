---
layout: docs
title: Download
description: Download `anchor` binary / source-files to any environment, local or CI.
group: getting-started
toc: true
---

## Package Managers

Pull in `anchor`'s binary using popular package managers.

### Homebrew

The fastest way (for `macOS` and `Linux`) to install `anchor` is using [Homebrew](https://brew.sh/):

```bash
brew install ZachiNachshon/tap/anchor
```

Alternatively, tap into the formula to have brew search capabilities on that tap formulas:

1. Tap into `ZachiNachshon` formula

    ```bash
    brew tap ZachiNachshon/tap
    ```

1. Install the latest `anhcor` binary

    ```bash
    brew install anchor
    ```

## Pre-Built Release

1. Update the download script with the following parameters:

     - **VERSION:** binary released version
     - **OS_ARCH:** operating system &amp; architecture tuple

1. Download and install `anchor` binary (copy & paste into a terminal):

```bash
bash <<'EOF'

# Change Version, OS and Architecture accordingly
VERSION=0.1.0

OS_ARCH=darwin_amd64
# Options: 
#   - darwin_arm64
#   - linux_arm64
#   - linux_armv6
#   - linux_amd64

# Create a temporary folder
repo_temp_path=$(mktemp -d ${TMPDIR:-/tmp}/anchor-repo.XXXXXX)
cwd=$(pwd)
cd ${repo_temp_path}

# Download & extract
echo -e "\nDownloading anchor to temp directory...\n"
curl -SL "https://github.com/ZachiNachshon/anchor/releases/download/v${VERSION}/anchor_${VERSION}_${OS_ARCH}.tar.gz" | tar -xz

# Create a dest directory and move the binary
echo -e "\nMoving binary to ~/.local/bin"
mkdir -p ${HOME}/.local/bin; mv anchor ${HOME}/.local/bin

# Add this line to your *rc file (zshrc, bashrc etc..) to make `anchor` available on new sessions
echo "Exporting ~/.local/bin (make sure to have it available on PATH)"
export PATH="${PATH}:${HOME}/.local/bin"

cd ${cwd}

# Cleanup
if [[ ! -z ${repo_temp_path} && -d ${repo_temp_path} && ${repo_temp_path} == *"anchor-repo"* ]]; then
	echo "Deleting temp directory"
	rm -rf ${repo_temp_path}
fi

echo -e "\nDone (type 'anchor' for help)\n"

EOF
```

Alternatively, you can download a release directy from GitHub

<a href="{{< param "download.dist" >}}" class="btn btn-bd-primary" onclick="ga('send', 'event', 'Getting started', 'Download', 'Download Anchor');" target="_blank">Download Specific Release</a>

{{< callout warning >}}
## `PATH` awareness

Make sure `${HOME}/.local/bin` exists on the `PATH` or sourced on every new shell session.
{{< /callout >}}

## Build from Source

Clone `anchor` repository into a directory of your choice:

```bash
git clone https://github.com/ZachiNachshon/anchor.git; cd anchor
```

{{< callout info >}}
{{< js.inline >}}
***Note:** Go `{{ $.Site.Params.go_version }}` is required to build from source.*
{{< /js.inline >}}
{{< /callout >}}

<br>

#### Build to Custom Path

1. Change directory to the local `anchor` cloned repository
1. Build the binary (destination: `PWD`)

   ```bash
   make build
   ```

1. Copy the binary to a bin folder that exists on `PATH`

   ```bash
   cp anchor ${HOME}/.local/bin
   ```

1. **(Optional)** Alternatively, copy directly to `/usr/local/bin`

    ```bash
   cp anchor /usr/local/bin
   ```

{{< callout warning >}}
## `PATH` awareness

Make sure `${HOME}/.local/bin` exists on the `PATH` or sourced on every new shell session.
{{< /callout >}}

<br>

#### Build to GOPATH

1. Change directory to the local `anchor` cloned repository
1. Build the binary (destination: `GOPATH/bin`)

   ```bash
   make build-to-gopath
   ```
