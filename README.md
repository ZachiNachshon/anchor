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
  <a href="https://img.shields.io/github/downloads/ZachiNachshon/anchor/total">
    <img src="https://img.shields.io/github/downloads/ZachiNachshon/anchor/total" alt="Downloads"/>
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

**Anchor** is a lightweight CLI tool that grants the **dynamic marketplace** experience for local / CI environment by connecting to single/multiple remote repositories, each represents a different marketplace of items, each with its own executable actions.

Every marketplace repository allows Anchor to **centralize and organize** a set of domain items into their own categories, every domain containing a list of executable single **action** and/or grouped actions (**workflows**) per item in a coherent, visible and easy-to-use approach. 

**Anchor** connects to remote git repositories, each containing an opinionated structure that allows it to understand what is available, exposing to the user every **Category** as a dynamicly created CLI command, underlying domain items with their defined **actions / workflows (actions-sets)** using an **interactive** selector enriched with **documentation**.

| :heavy_exclamation_mark: WARNING |
| :--------------------------------------- |
| Anchor is still in **alpha stage**, breaking changes might occur, use it with caution ! |

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
- [How does it work?](#how-does-it-work)
- [Installation](docs/installation.md)
- [Configuration](docs/configuration.md)
- [Available features](docs/available-features.md)
- [Create a marketplace repository](docs/create-anchorfiles.md)
- [Common questions](docs/common-questions.md)
- [Troubleshooting](docs/troubleshooting.md)

**Maintainers / Contributors:**

- [Contribute Guides](docs/contribute-guides.md)

<br>

<h3 id="why-creating-anchor">⛵ Why Creating <code>Anchor</code>?</h3>

1. I believe that local environment management should be a *smooth sailing*, well documented process with minimum context switches for *running scripts / installing applications / orchestrate installations / do whatever you require* on it
1. Allowing to compose different actions from multiple channels (shell scripts, CLI utilities etc..) into a coherent well documented workflow with rollback procedure
1. Having an action / workflow execution plan explained in plain english and managed via a central versioned controlled remote repository that can be shared with others to use
1. Using an agnostic client that doesn't change, rather, changes are reflected based on remote marketplace state

<br>

<h3 id="how-does-it-work">🗺 How Does It Work?</h3>

1. Create a structured repository (a.k.a anchorfiles) as the remote marketplace

   <details><summary>Show</summary>

   </details>
   
1. Define actions/workflows within the an ***instructions.yaml*** file

   <details><summary>Show</summary>
   <img style="vertical-align: top;" src="assets/images/anchorfiles-structure.png" height="400" >
   </details>

1. Set the remote marketplace in anchor ***config.yaml*** file

   <details><summary>Show</summary>
   <img style="vertical-align: top;" src="assets/images/anchor-config.png" width="500" >
   </details>

1. Use anchor CLI commands to interact with the marketplace

   <details><summary>Show</summary>
   <img style="vertical-align: top;" src="assets/images/anchor-select-app.png" width="500" >
   </details>

<br>

<h2 id="support">Support</h2>

Anchor is an open source project that is currently self maintained in addition to my day job, you are welcome to show your appreciation by sending me cups of coffee using the the following link as it is a known fact that it is the fuel that drives software engineering ☕

<a href="https://www.buymeacoffee.com/ZachiNachshon" target="_blank"><img src="assets/images/bmc-orig.svg" height="57" width="200" alt="Buy Me A Coffee"></a>

<br>

<h2 id="licence">Licence</h2>

MIT

<br>
