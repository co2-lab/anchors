<!-- @anchors
  code: FLRGF
  updated_at: 2026-10-12
  layer: comando
-->
# FlowRegister — attach the flow domain's commands to the root command, once each

> **Code**: `FLRGF`

## Overview

The flow domain groups the commands an agent uses to move work along: take a task (`queue`, `next`,
`done`, `drop`, `reclaim`), compose the work prompt (`work`), record and confront a delivery
(`deliver`, `merge-progress`), escalate and resolve what blocks (`escalate`, `decided`, `unblock`,
`backfill-labels`, `discard`), report where a task stands (`task-status`, `pr-body`) and run the
watcher (`watch`), and track the project's progress (`progress`). This unit is the single place where
those commands are attached to the root command of the CLI.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the root command | the CLI's root command, before it parses arguments | a command that already holds the flow commands | the CLI's root: it registers each domain once |

## Effects

| Effect | Description |
| --- | --- |
| `FLRGF-B01` | After registration the root holds exactly the nineteen flow commands: backfill-labels, decided, deliver, discard, done, drop, escalate, merge-progress, monitor, next, pr-body, progress, queue, reclaim, report-bug, task-status, unblock, watch and work. |
| `FLRGF-B02` | The watcher's control commands (start, run, status, stop, pause, resume, logs) arrive under `watch`, not at the root. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FLRGF-I01` | Each flow command is attached to the root exactly once. | the names of the root's children after registration are counted, and none appears twice |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |

## Errors

none — registration only attaches commands; it reads no input and handles no failure.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
