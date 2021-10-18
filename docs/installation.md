<h1 id="installation" align="center">⎈ Installation<br><br></h1>

- [Pre-Built Release](#pre-built-release)
- [Build from Source](#build-from-source)

<br>

<h2 id="pre-built-release">Pre-Built Release</h2>

1. Download and install `anchor` binary

   ```bash
   curl -sfL https://get.anchor/install.sh | sh -
   ```

1. Setup config with a dynamic remote marketplace

   ```bash
   curl -sfL https://get.anchor/setup-config.sh | sh -
   ```

<br>


<h2 id="pre-built-release">Build from Source</h2>

1. Clone `anchor` repository

   ```bash
   git clone git@github.com:ZachiNachshon/anchor.git; cd anchor
   ```
   
1. Build a binary

   ```bash
   make build
   ```
   
1. Copy the binary to the bin folder in use

   ```bash
   # Make sure '${HOME}/.local/bin' exists in PATH or sourced on every new session
   cp anchor ${HOME}/.local/bin
   
   # Alternatively, copy directly to /usr/local/bin
   cp anchor /usr/local/bin
   ```

   | :bulb: Note |
   | :--------------------------------------- |
   | Another option is to save the binary in a unified location and create a symlink within the bin folder. |

<br>

