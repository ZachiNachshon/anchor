---
layout: docs
title: Menu Command
description: Create a dynamic CLI menu command.
toc: true
group: repository
---

## The `command.yaml` file

This file represents a dynamically created CLI command, a folder containing the `command.yaml` file is treated as a CLI menu command of its own. 

The content of the file dictates how the command will appear on the CLI menu and how to interact with it.

```yaml
name: <command-name>
description: <command-long-description>
command:
  use: <cli-usage>
  short: <cli-short-description>
```

{{< bs-table >}}
| Task | Description |
| --- | --- |
| `name` | Command name |
| `description` | Command description |
| `command.use` | CLI menu command usage keyword |
| `command.short` | CLI menu short description |
{{< /bs-table >}}

Example:

```yaml
name: "Infrastructure team"
description: "Provisioning actions for: backoffice, developement and production"
command:
  use: "team-infra"
  short: "Infrastructure team provisioning actions"
```

<div class="col-lg-6">
   <img style="vertical-align: top;" src="/docs/latest/assets/img/anchor-command.svg" width="800" >
</div>

## Sub-commands preset

Every menu command derived from the `command.yaml` file will contain a preset of sub-commands. Those sub-commands interact with the [actions](/docs/{{< param docs_version >}}/repository/action) or [workflows](/docs/{{< param docs_version >}}/repository/action-set-workflow) defined in the `instructions.yaml` file of the respective folder.

{{< bs-table >}}
| Task | Description |
| --- | --- |
| `run` | Run an action non-interactively |
| `select` | Select an action using an interactive menu |
| `status` | Check status validity of available actions |
{{< /bs-table >}}

<div class="col-lg-6">
   <img style="vertical-align: top;" src="/docs/latest/assets/img/anchor-sub-commands.svg" width="800" >
</div>


## Interactive command

The `select` sub-command should be used when you wish to start an interactive menu selection for action or workflow.

Example:

```text
anchor team-infra select
```

## Non-interactive command

The `run` sub-command should be used when you wish to run an action or workflow directly, skipping the selection menu.

**Action:**

```bash
anchor team-infra run backoffice --action=install-jenkins-master
```

**Workflow:**

```bash
anchor team-infra run backoffice --workflow=provision-jenkins-server-agents
```

{{< callout info >}}
The `run` command appears on every action or workflow information panel when using an interactive selection menu.
{{< /callout >}}

<div class="col-lg-6">
   <img style="vertical-align: top;" src="/docs/latest/assets/img/anchor-select-menu-jenkins.svg" width="800" >
</div>