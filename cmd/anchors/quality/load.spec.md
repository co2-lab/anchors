<!-- @anchors
  code: LDQRS
  updated_at: 2026-11-19
  layer: comando
-->
# SuiteLoad — a suite's run waits for the load queue, and the commit hook jumps it

> **Code**: `LDQRS`

## Overview

`anchors test` and `anchors mutation` run a project's suites — whole, by the impact path, by a
selection or a budget. With `load:` declared (LDCFS), each run waits for the queue (LDQUS) at the
level its files' scope costs, and holds its lease through the ingestion after it, which weighs too.
A whole run names no file, and is placed by its workspace.

The commit hook jumps the queue: `check --index` and `verify --phase pre-commit` turn it off for
themselves and every command they run, so a commit never waits behind another worktree's heavy run.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the run's files | absolute or relative to the root; none for a whole run | — | the suite runner |
| the project's configuration | `anchors.yaml` with or without `load:` | one that does not load | this unit: the run starts at once |

## Effects

| Effect | Description |
| --- | --- |
| `LDQRS-B01` | A suite's run holds a lease at the level its files' scope costs — a whole run placed by its workspace, files outside every scope at the task's default — from before the run until after its ingestion; with no `load:`, with the queue off (the commit hook), or with a queue that cannot be read, the run starts at once. (`admitRun`, `skipLoadQueue`) |

## Errors

none — what the queue cannot answer never stops a run.
