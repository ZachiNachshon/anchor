---
layout: docs
title: Overview
description: Learn about <code>anchor</code>, why it was created and the pain it comes to solve.
group: getting-started
toc: true
---

## Why creating `anchor`?

1. Allow a better experience on repositories containing lots of scripts managed by multiple teams, make them approachable and safe to use by having a documented and controlled process with minimum context switches for *running scripts / installing applications / orchestrate installations / do whatever you require*

1. Allowing to compose different actions from multiple channels (shell scripts, CLI utilities etc..) into a coherent well documented workflow with rollback procedure

1. Having an action / workflow execution plan explained in plain english and managed via a central versioned controlled remote repository in a GitOps way that can be shared with others to use easily

1. Remove the fear of running an arbitrary undocumeted script that relies on ENV vars to control its execution

1. Using an agnostic client that doesn't change, rather, changes are reflected based on the remote git repository(ies) it relies on

1. Reduce the amount of CLI utilities created in a variety of languages in an organization

## In a nutshell

`anchor` is a lightweight CLI tool written in `go` that allows to create dynamic CLI's based on remote git repositories containing a simple opinionated structure. It exposes any executable commands it finds as a dynamic command-line-interface utility.

Think of it as a dynamic scripts / commands marketplace experience for any environment, local and CI. By&nbsp;connecting to any git remote repositories, each represents a different domain (examples: `guild-onboarding` | `team-infra` | `team-build`), each with its own executable actions.

Every repository allows `anchor` to centralize and organize a set of domain actions into their own categories, every domain containing a list of executable single action and/or grouped actions-sets (workflows) in a coherent, visible and easy-to-use approach via intractive shell experience.

`anchor` connects to remote git repositories containing an opinionated structure that allows it to understand what is available, exposing pre-defined categories as a dynamicly created CLI commands with their underlying domain items as actions / actions-sets (workflows) as follows:

   - Interactive selector enriched with documentation
   - Non-interactive mode i.e. direct CLI command
   