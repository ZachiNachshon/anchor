---
layout: docs
title: Configuration
description: Set a remote git repository as the dynamic context for <code>anchor</code>.
group: content
toc: true
---

## Use Context

`anchor` configuration is based on config contexts, each is pointing to a remote/local git repository. When running an `anchor` command, it must have an active context in order to reflect its dynamic CLI menu based on the git repository structure.

{{< callout info >}}
When no context is defined as active i.e. `currentContext`, the user will be prompoted by an interactive menu to select which of the available config contexts he would like to use for a specific command run.
{{< /callout >}}

* If no context is defined - the user is prompted to select a context
<div class="col-lg-6">
   <img style="vertical-align: top;" src="/docs/latest/assets/img/anchor-context-selection.svg" width="600" >
</div>

* If a context is already defined - a remote repository scan takes place for possible changes if `autoUpdate` was enabled and the command is executed afterwards

<div class="col-lg-6">
   <img style="vertical-align: top;" src="/docs/latest/assets/img/anchor-auto-update-flow.svg" width="800" >
</div>

## Set Context Entry

Register a new remote git repository or update an existing one as an `anchor` config context. Use the flags to control each remote repository with default values and/or custom ones. Every remote marketplace repository is identified under a distinct ***config context*** name.

### Flags

{{< bs-table >}}
| Task | Description | Type | Default value | 
| --- | --- | --- | --- |
| `--repository.remote.url` | URL of the remote Git repository | `string` | **Required** |
| `--repository.remote.branch` | Specify a branch to fetch HEAD references | `string` | master
| `--repository.remote.revision` | Specific commit to check out | `string` | - |
| `--repository.remote.clonePath` | Local path to clone the repository into.<br/>(**default path**: `$HOME/.config/anchor/repositories/<context-name>`) | `string` | - |
| `--repository.remote.autoUpdate` | Allow checking remote marketplace for available changes and auto update if there are any.<br>Performs a remote check on every `anchor` CLI action | `bool` | False
| `--repository.local.path` | Local path of an `anchor` structured project / repository / any folder | `string` | - |
| `--set-current-context` | Set the newly registered config context as the active one | `string` | - |
{{< /bs-table >}}

### Usage

An example of how to register a remote git repository with custom attributes under ***ops-team*** as its config context identifier:

```bash
anchor config set-context-entry ops-team \
    --repository.remote.url=git@github.com:Organization/OpsMarketplace.git \
    --repository.remote.branch=my-branch \
    --repository.remote.revision=123456789abcdefg \
    --repository.remote.clonePath=/custom/path/to/config \
    --repository.remote.autoUpdate=true
```

## View Configuration 

View the local configuration by printing it to the terminal as plain text:

```bash
anchor config view
```

## Edit Configuration

Edit the local configuration by entering into a text edit mode:

```bash
anchor config edit
```

