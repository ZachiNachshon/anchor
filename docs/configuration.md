<h1 id="configuration" align="center">Configuration<br><br></h1>

- [Overview](#overview)
- [Use Context](#use-context)
- [Set Context Entry](#set-context-entry)
  - [Usage](#set-context-usage)
  - [Flags](#set-context-flags)
- [View Configuration](#view-config)
- [Edit Configuration](#edit-config)

<br>

<h3 id="overview">Overview</h3>

The `config` command is in charge of setting the remote marketplace repository context for `anchor` command.

<br>

<h3 id="use-context">Use Context</h3>

Anchor configuration is comprised of a single/multiple config contexts. Every context points to a remote/local repository that acts as a marketplace. When running anchor command, it must know which repository to interact with for that specific run.

* If no context is defined - the user is prompted to select a context

   <details><summary>Show</summary>
      <img style="vertical-align: top;" src="../assets/images/configuration/anchor-select-context.png" width="500" >
   </details>

* If a context is already defined - a repository scan takes place and the user is prompted with an interactive select menu afterwards

   <details><summary>Show</summary>
      <img style="vertical-align: top;" src="../assets/images/configuration/anchor-select-context.png" width="500" >
   </details>

<br>

<h3 id="set-context-entry">Set Context Entry</h3>

Update an existing registered remote markeplace repository values and/or add a new one. It is possible to control each remote repository with default values and/or custom ones. Every remote marketplace repository is identified with a ***config context*** name.

<h4 id="set-context-usage">Usage</h4>

Use the following command to set up a remote marketpace with custom attributes on the ***ops-team-context*** config context:

```bash
anchor config set-context-entry ops-team-context \
    --repository.remote.url=git@github.com:Organization/OpsMarketplace.git \
    --repository.remote.branch=my-branch \
    --repository.remote.revision=123456789abcdefg \
    --repository.remote.clonePath=/custom/path/to/config \
    --repository.remote.autoUpdate=true
```

<h4 id="set-context-flags">Flags</h4>

| **Name**                                                     | **Type** | **Default value** |
| :----------------------------------------------------------- | :------- | :---------------- |
| `repository.remote.url`                                      | `string` | **mandatory**     |
| The URI of the remote Git repository                         |          |                   |
| `repository.remote.branch`                                   | `string` | master            |
| Specify a branch to fetch HEAD references                    |          |                   |
| `repository.remote.revision`                                 | `string` |                   |
| Specific commit to get checked out                           |          |                   |
| `repository.remote.clonePath`                                | `string` |                   |
| Local path to clone the repository into.<br/>Default location is used to clone into `$HOME/.config/anchor/repositories/<context-name>` |          |                   |
| `repository.remote.autoUpdate`                               | `bool`   | False             |
| Allow checking remote marketplace for available changes and auto update if there are any.<br>Performs a remote check on every `anchor` CLI action |          |                   |

<br>

<h3 id="view-config">View Configuration</h3>

View the existing local configuration by printing it to the terminal as plain text:

```bash
anchor config view
```

<br>

<h3 id="edit-config">Edit Configuration</h3>

Edit the exiting local configuration by entering into a text edit mode:

```bash
anchor config edit
```

<br>
