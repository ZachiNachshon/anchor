---
layout: docs
title: Instructions
description: Define actions and action-sets for a menu command.
group: repository
toc: true
sections:
  - title: Action
    description: Create executable actions for a menu command
  - title: Action-Set (Workflow)
    description: Create executable action-sets (a.k.a workflows)
---

## The `instructions.yaml` file

This file purpose is declaring which actions or workflows are avaiable for a specific menu command and how one should interact with it.

- It is being scanned by `anchor` at runtime upon menu command execution
- **The name of the folder** containing the `instructions.yaml` file will be the name that is used as its identifier

The basic structure of an instruction file is as follows:

```yaml
instructions:

  actions:
    ...

  workflows:
    ...

```

## Contextual instructions

Every `instructions.yaml` file has an option to define a unified global context for all of its comprising actions and workflows, these affects the pre-execution flow with custom prompt messages.

```yaml
globals:
  context: <context-in-here>
  
instructions:

  actions:
    ...

  workflows:
    ...
```

### Available contexts

- `application`: prints basic prompt message before action is executed

<div class="col-lg-6">
   <img style="vertical-align: top;" src="/docs/latest/assets/img/anchor-app-context.svg" width="800" >
</div>

- `kubernetes`: prints Kubernetes current context information before action is executed

<div class="col-lg-6">
   <img style="vertical-align: top;" src="/docs/latest/assets/img/anchor-k8s-context.svg" width="800" >
</div>

{{< callout info >}}
When an action or workspace overrides the `context` attribute and it differs from the global one then the overridden context will be the one in use.
{{< /callout >}}

