<h3 align="center" id="anchor-logo"><img src="assets/anchor-logo.png" height="300"></h3>

<p align="center">
  <a href="https://img.shields.io/github/go-mod/go-version/ZachiNachshon/anchor/pivot">
    <img src="https://img.shields.io/github/go-mod/go-version/ZachiNachshon/anchor/pivot" alt="Go Version"/>
  </a>
  <a href="https://github.com/ZachiNachshon/anchor/actions/workflows/ci.yaml/badge.svg?branch=pivot">
    <img src="https://github.com/ZachiNachshon/anchor/actions/workflows/ci.yaml/badge.svg?branch=pivot" alt="GitHub CI status"/>
  </a>
  <a href="https://goreportcard.com/badge/ZachiNachshon/anchor">
    <img src="https://goreportcard.com/badge/ZachiNachshon/anchor" alt="Go Report Card"/>
  </a>
  <a href="https://coveralls.io/repos/github/ZachiNachshon/anchor/badge.svg?branch=pivot">
    <img src="https://coveralls.io/repos/github/ZachiNachshon/anchor/badge.svg?branch=pivot" alt="Go Coverage"/>
  </a>
  <a href="https://github.com/ZachiNachshon/anchor/releases">
    <img src="https://img.shields.io/github/v/release/ZachiNachshon/anchor?include_prereleases&style=flat-square" alt="Go Releases"/>
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"/>
  </a>
  <a href="https://www.paypal.me/ZachiNachshon">
    <img src="https://img.shields.io/badge/$-donate-ff69b4.svg?maxAge=2592000&amp;style=flat">
  </a>
</p>

<p align="center">
  <a href="#requirements">Requirements</a> •
  <a href="#quickstart">QuickStart</a> •
  <a href="#overview">Overview</a> •
  <a href="#support">Support</a> •
  <a href="#license">License</a>
</p>

<br>

**Anchor** is a lightweight CLI tool that grants the **dynamic marketplace** experience for local environment needs by connecting to single/multiple remote repositories, each represent a different marketplace of executable actions.

Every such marketplace repository allows Anchor to **centralize and organize** executable single action and/or multiple actions (workflows) in a coherent, visible and easy-to-use approach. 

Action refers to anything that can be executed within a shell session, Anchor connects to remote git repositories, each containing a specific structure that allows anchor to understand what is available, exposing the repository content with **actions / workflows (actions-sets)** using an **interactive** selection prompter enriched with **documentation**.

<br>

<details><summary>Create a repository with the expected structure</summary>
<img style="vertical-align: top;" src="assets/images/anchorfiles-structure.png" height="400" >
</details>

<details><summary>Define actions/workflows within the <code>instructions.yaml</code> file</summary>
<img style="vertical-align: top;" src="assets/images/anchorfiles-structure.png" height="400" >
</details>

<details><summary>Add the remote marketplace to <code>Anchor</code> configuration file</summary>
<img style="vertical-align: top;" src="assets/images/anchor-config.png" width="500" >
</details>

<details><summary>Use <code>anchor select app</code> to interact with the marketplace</summary>
<img style="vertical-align: top;" src="assets/images/anchor-select-app.png" width="500" >
</details>

<br>

<h2 id="requirements">🏴‍☠️ Requirements</h2>

- A Unix-like operating system: macOS, Linux
- `git` (recommended v2.30.0 or higher)

<br>

<h2 id="quickstart">⚡️ QuickStart</h2>

1. Download and install `anchor` binary

   ```bash
   curl -sfL https://get.anchor/install.sh | sh -
   ```

2. Setup config with a dynamic remote marketplace

   ```bash
   curl -sfL https://get.anchor/setup-config.sh | sh -
   ```

<br>

<h2 id="overview">⚓️ Overview</h2>

- [Why creating `Anchor`?](#why-creating-anchor)
- [Installation](docs/installation.md)
- [Configuration](docs/configuration.md)
- [Available Features](docs/available-features.md)
- [Create a Marketplace Repository](docs/create-anchorfiles.md)
- [Common Questions](docs/common-questions.md)
- [Troubleshooting](docs/troubleshooting.md)

**Maintainers / Contributors:**

- [Contribute Guides](docs/contribute-guides.md)

<br>

<h3 id="incentive">⛵ Why Creating <code>Anchor</code>?</h3>

1. I beleive that local environment management should be a smooth sailing, well documented process with minimum context switches for *running scripts / installing applications / orchestrate installations / do whatever you require* on it
1. Allowing to compose different actions from multiple channels (shell scripts, CLI utilities etc..) into a coherent well documented workflow with rollback procedure
1. Having an action / workflow execution plan explained in plain english and managed via a central versioned controlled remote repository that can be shared with others to use
1. Using an agnostic client that doesn't change, rather, changes are reflected based on remote marketplace state

<br>

<h2 id="support">Support</h2>

<a href="https://www.buymeacoffee.com/ZachiNachshon" target="_blank"><img src="assets/images/bmc-orig.svg" height="57" width="200" alt="Buy Me A Coffee"></a>

<br>

<h2 id="licence">Licence</h2>

MIT

<br>
