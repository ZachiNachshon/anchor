---
layout: docs
title: Contribute
description: Contribute to the development of `anchor` using the documentation, build scripts and tests.
group: getting-started
toc: true
aliases: "/docs/latest/getting-started/contribute/"
---

## Tooling setup

- [Go lang](https://go.dev/dl/) `v1.20` is required for compilation and tests
- [Node.js](https://nodejs.org/en/download/) is optional for managing the documentation site

{{< callout info >}}
Docs site is using npm scripts to build the documentation and compile source files. The `package.json` houses these scripts which used for various docs development actions.
{{< /callout >}}

## Guidelines

- PRs need to have a clear description of the problem they are solving
- PRs should be small
- Code without tests is not accepted, PRs must not reduce tests coverage
- Contributions must not add additional dependencies
- Before creating a PR, make sure your code is well formatted, abstractions are named properly and design is simple
- In case your contribution can't comply with any of the above please start a GitHub issue for discussion

## How to Contribute?

1. Fork this repository
1. Create a PR on the forked repository
1. Send a pull request to the upstream repository

## Development Scripts

The `makefile` within this repository contains numerous tasks used for project development. Run `make` to see all the available scripts in your terminal.

{{< bs-table >}}
| Task | Description |
| --- | --- |
| `update-externals` | Update external source dependents |
| `deps` | Tidy, verify and vendor go dependencies |
| `fmtcheck` | Validate Go code format and imports |
| `fmt` | Format Go code using gofmt style and sort imports |
| `test` | Run tests suite on host machine |
| `test-containerized` | Run tests suite within a Docker container |
| `test-with-coverage` | Run tests suite on host machine with coverage report |
| `docs-site` | Run a local documentation site |
| `docs-site-lan` | Run a local documentation site (LAN available) |
| `build` | Build a binary for system OS/Arch |
| `build-main-package` | Build main package for system OS/Arch |
| `install` | Build and Install a Go binary locally |
| `delete` | Delete a locally installed Go binary |
| `github-release-create` | Build and publish Go binary(ies) as GitHub release |
| `github-release-delete` | Prompt for a GitHub release tag to delete |
{{< /bs-table >}}

## Testing Locally

This repository had been developed using the TDD methodology (Test Driven Development). Tests allow you to make sure your changes work as expected, don't break existing code and keeping code coverage high.

Running tests locally allows you to have short validation cycles instead of waiting for the PR status to complete.

**How to run a test suite?**

1. Clone the `anchor` repository
2. Run `make tests` to use the locally installed `go` runtime
3. Alternatively, run `make test-containerized` to use the same `go` runtime which is supported by this repository

{{< callout info >}}
`tparse` should be installed for performing `go test` output analysis ([instructions in here](https://github.com/mfridman/tparse))
{{< /callout >}}

## Documentation Scripts

The `/docs-site/package.json` includes numerous tasks for developing the documentation site. Run `npm run` to see all the available npm scripts in your terminal. Primary tasks include:

{{< bs-table >}}
| Task | Description |
| --- | --- |
| `npm run docs-build` | Cleans the Hugo destination directory for a fresh serve |
| `npm run docs-serve` | Builds and runs the documentation locally |
{{< /bs-table >}}

## Local documentation 

Running our documentation locally requires the use of Hugo, which gets installed via the `hugo-bin` npm package. Hugo is a blazingly fast and quite extensible static site generator. Here’s how to get it started:

- Run through the [tooling setup](#tooling-setup) above to install all dependencies
- Navigate to `/docs-site` directory and run `npm install` to install local dependencies listed in `package.json`
- From `/docs-site` directory, run `npm run docs-serve` in the command line
- Open [http://localhost:9001/](http://localhost:9001/) in your browser

Learn more about using Hugo by reading its [documentation](https://gohugo.io/documentation/).

## Troubleshooting

In case you encounter problems with missing dependencies, run `make align-deps`.
