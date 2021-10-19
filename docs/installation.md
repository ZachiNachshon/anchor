<h1 id="installation" align="center">⎈ Installation<br><br></h1>

- [Pre-Built Release](#pre-built-release)
- [Build from Source](#build-from-source)

<br>

<h2 id="pre-built-release">Pre-Built Release</h2>

Download and install `anchor` binary (copy & paste into a terminal):

```bash
bash <<'EOF'

# Change Version, OS and Architecture accordingly
VERSION=0.1.0

OS_ARCH=Darwin_x86_64
# Options: 
#   - Darwin_arm64
#   - Linux_arm64
#   - Linux_armv6
#   - Linux_x86_64

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

<br>


<h2 id="pre-built-release">Build from Source</h2>

Clone `anchor` repository into a directory of your choice:

```bash
git clone git@github.com:ZachiNachshon/anchor.git; cd anchor
```

<br>

<h4>Build to Custom Path</h4>

1. Build a binary to current directory

   ```bash
   make build
   ```

1. Copy the binary to a bin folder in use

   ```bash
   # Make sure '${HOME}/.local/bin' exists in PATH or sourced on every new session
   cp anchor ${HOME}/.local/bin
   
   # Alternatively, copy directly to /usr/local/bin
   cp anchor /usr/local/bin
   ```

   | :bulb: Note |
   | :--------------------------------------- |
   | Alternatively, save the binary to a unified location and create a symlink to it from `/usr/local/bin`. |

<br>

<h4>Build to <code>GOPATH</code></h4>

1. Run the following to build and place the binary in `${GOPATH}/bin`

   ```bash
   make build-to-gopath
   ```

<br>

