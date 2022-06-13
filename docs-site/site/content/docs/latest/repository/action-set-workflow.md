---
layout: docs
title: Action Sets (Workflow)
description: Create executable action-sets (a.k.a workflows) for a menu command.
group: repository
toc: false
---

Workflows are defined under the `workflows` YAML attribute of the `instructions.yaml` file.

```yaml
instructions:
  workflows:
    - id: "workflow-id"
      displayName: "Selection name"
      title: "Shorter information"
      description: "Longer information"
      tolerateFailures: true
      actionIds:
        - action-id-1
        - action-id-2
```

{{< bs-table >}}
| Task | Description | Type | Default value | 
| --- | --- | --- | --- |
| `instructions.workflows[].id` | Identifier, by default this is the text used as the interactive selector | `string` | **Required** |
| `instructions.workflows[].displayName` | Alternative name to use for displaying the workflow on the interactive selector (Optional) | `string` | - |
| `instructions.workflows[].title` | Short title about what this workflow is all about, appended to the interactive selector | `string` | - |
| `instructions.workflows[].description` | Longer info about the workflow responsibilities, shown on the information section - off the interactive selector| `string` | - |
| `instructions.workflows[].tolerateFailures` | Continue to the next action upon failure or fail the entire workflow | `bool` | False |
| `instructions.workflows[].context` | Context scope of the workflow<br>Available values:<br>&nbsp;&nbsp;• **application** - prints basic prompt message<br>&nbsp;&nbsp;• **kubernetes** - prints k8s current context and check if `kubectl` exists| `string` | application |
| `instructions.workflows[].actionIds` | List of ordered action ids to get executed upon workflow run | `string[]` | - |
| `globals.context` | Same as action context, but applied to all actions in the file | `string` | - |
{{< /bs-table >}}

Example:

```yaml
instructions:
  actions:
    - id: install-jenkins-master
        ...

    - id: install-jenkins-agents
        ...

  workflows:
    - id: provision-jenkins-server-agents
      displayName: "Jenkins: Provision Server/Agents"
      title: "No-Op action, used as an example"
      description: "Provision a Jenkins server and x3 agents"
      tolerateFailures: true
      actionIds:
        - install-jenkins-master
        - install-jenkins-agents
```

{{< callout info >}}
In the following example, `backoffice` is the name of the folder containing the `instructions.yaml` file.
{{< /callout >}}

<div class="col-lg-6">
   <img style="vertical-align: top;" src="/docs/latest/assets/img/anchor-workflow-selection.svg" width="800" >
</div>

