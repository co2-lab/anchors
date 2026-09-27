<!-- @anchors
  code: OPRGP
  updated_at: 2026-09-27
  layer: comando
-->
# OpsRegister — the operation commands reach the CLI through one registration point

> **Code**: `OPRGP`

## Overview

The CLI is assembled by domain: each command package hands its commands to the root through
one registration function, and the root knows only the domains. This unit is the registration
of the operations domain — the commands that set a project up and operate it: init, new,
migrate, freeze and thaw, settings, install-hooks, commit-msg, board, docs, changelog, suggest,
synthesize, code and generated-paths.

Registering in one place is what makes a command exist for the user, for the root's freeze
allowlist, and for the guides and pipelines that cite commands by name; a command missing here
is an `unknown command` for whoever follows the instruction.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the root command | the CLI's root, or any command acting as one | — | the caller (the CLI's root assembly) |

## Effects

| Effect | Description |
| --- | --- |
| `OPRGP-B01` | The root receives the fifteen operation commands: board, changelog, code, commit-msg, docs, freeze, generated-paths, init, install-hooks, migrate, new, settings, suggest, synthesize and thaw. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `OPRGP-I01` | Each operation command is registered exactly once. | registers on a bare root and compares the sorted names with the list |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `OPRGP-X01` | The registration adds commands and nothing else: no flag, no hook, no behaviour of the root. | The root's cross-cutting behaviour (language, freeze) belongs to the root, once for every domain. |

## Errors

none — registration only attaches commands built by their own units; it has no input to reject and no failure path.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `cmd/anchors/root.go` | the root assembly that calls this registration | comando — CLRTC |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
