---
layout: docs
title: Structure
description: Learn how to connect an existing or new git repository to `anchor`, reflecting the repository content via command-line-interface utility.
group: repository
toc: false
aliases: "/anchor/docs/latest/repository/"
sections:
  - title: Menu Command
    description: Create a dynamic CLI menu command
  - title: Instructions
    description: Define actions and action-sets for a menu command
---

## Overview

`anchor` connects to any git repository and expose executable commands as dynamic command-line-interface utility to use from any environment, CI and local. 

In order for it to properly scan the respository and extract the commands with their respective actions or workflows, a basic structure should be followed introducing a set of YAML files.

Example of such structure:

```text
├── ...
├── <cli-command-1>                   
│   └── <command-1-actions-1>               
│       ├── instructions.yaml
│       ├── <additional-files-and-folders>
│       └── ...       
│   ├── command.yaml
│   ├── <additional-files-and-folders>
│   └── ...  
├── <cli-command-2>                   
│   └── <command-2-actions-1>               
│       ├── instructions.yaml
│       └── ...       
│   └── <command-2-actions-2>               
│       ├── instructions.yaml
│       └── ...       
│   ├── command.yaml
│   └── ...  
```