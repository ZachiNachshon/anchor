---
layout: docs
title: Action
description: Create executable actions for a menu command.
group: repository
toc: false
---

Actions are defined under the `actions` YAML attribute of the `instructions.yaml` file.

```yaml
instructions:
  actions:
    - id: "action-id"
      displayName: "Selection name"
      title: "Shorter information"
      description: "Longer information"
      showOutput: true
      script: |
        echo "Hello World !"
```

{{< callout info >}}
Environment variables substitution is supported, use the env vars within `script` or `scriptFile` attributes.
{{< /callout >}}

{{< bs-table >}}
| Task | Description | Type | Default value | 
| --- | --- | --- | --- |
| `instructions.actions[].id` | Identifier, by default this is the text used as the interactive selector | `string` | **Required** |
| `instructions.actions[].displayName` | Alternative name to use for displaying the action on the interactive selector (Optional) | `string` | - |
| `instructions.actions[].title` | Short title about what this action is all about, appended to the interactive selector | `string` | - |
| `instructions.actions[].description` | Longer info about the action responsibilities, shown on the information section - off the interactive selector | `string` | - |
| `instructions.actions[].showOutput` | Print action execution output to stdout rather than using a spinner indicator for success / failure | `bool` | False |
| `instructions.actions[].script` | Free text script content to get executed upon selection<br>(mutual exclusive with `scriptFile`) | `string` | **Required** |
| `instructions.actions[].scriptFile` | Scripted file to get executed upon selection, path is relative to repository root folder<br>(`SHEBANG` is required)<br>(mutual exclusive with `script`) | `string` | **Required** |
| `instructions.actions[].context` | Context scope of the action, affects the pre-action execution flow with custom prompt messages.<br>Available values:<br>&nbsp;&nbsp;• **application** - prints basic prompt message<br>&nbsp;&nbsp;• **kubernetes** - prints k8s current context and check if `kubectl` exists| `string` | application |
| `globals.context` | Same as action context, but applied to all actions in the file | `string` | - |
{{< /bs-table >}}

Example:

```yaml
instructions:
  actions:
    - id: install-jenkins-agents
      displayName: "Jenkins: Install Agents"
      title: "No-Op action, used as an example"
      showOutput: true
      scriptFile: infra-team/backoffice/scripts/jenkins-install-agents.sh ${JENKINS-MASTER-TOKEN}
      description: "Install a Jenkins agent on multiple machines and connect them to a running server.

   GCP Cloud Provider:
     • Region: us-west1

   Ansible playbook vars:
     • Server IP Address
     • Number of agents (default x3)
     • Secrets fetched from Vault"

    - id: install-jenkins-master
      displayName: "Jenkins: Install Server"
      title: "No-Op action, used as an example"
      showOutput: true
      script: |
        echo "SSH into a remote machine and install a Jenkins server..."
        sleep 1
      description: "Install a Jenkins server on a remote machine.

   GCP Cloud Provider:
     • Region: us-west1
     
   Ansible playbook vars:
     • IP Address
     • Secrets fetched from Vault"
```

{{< callout info >}}
In the following example, `backoffice` is the name of the folder containing the `instructions.yaml` file.
{{< /callout >}}

<div class="col-lg-6">
   <img style="vertical-align: top;" src="/docs/latest/assets/img/anchor-select-menu-jenkins.svg" width="800" >
</div>
