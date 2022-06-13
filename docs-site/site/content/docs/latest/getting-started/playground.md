---
layout: docs
title: Playground
description: Take <code>anchor</code> for a spin by using it on a remote playground git repository.
group: getting-started
toc: true
---

## Instructions
Follow these steps to connect to a remote git playground repository and check the dynamic CLI live experience. **All actions are no-op**, you can safely run them as they only print to stdout.

1. Register to a remote git playground respository and set it as the default config context:

   ```bash
   anchor config set-context-entry playground \
      --repository.remote.url=https://github.com/ZachiNachshon/anchor-playground.git \
      --repository.remote.autoUpdate=false \
      --set-current-context
   ```

1. Type `anchor` to fetch the repository and print all available commands

1. Check which items are available under the `team-infra` command:

   ```bash
   anchor team-infra status
   ```
   
1. Select the `team-infra` command to start an interactive action selection, try running an action/workflow:

   ```txt
   anchor team-infra select
   ```

1. Use the `run` command to run an action non-interactively:

   ```bash
   anchor team-infra run backoffice --action=install-jenkins-master
   ```

1. Run an action-set (workflow) non-interactively:

   ```bash
   anchor team-infra run backoffice --workflow=provision-jenkins-server-agents
   ```

1. Use other playground commands and run different actions to check different use cases

<br>

{{< callout info >}}
This is a quick overview just to get a grasp of how simple it is to use `anchor`.<br>To add `anchor` support to an existing or new git repository, please see the [structure section](/docs/{{< param docs_version >}}/repository/structure).
{{< /callout >}}
