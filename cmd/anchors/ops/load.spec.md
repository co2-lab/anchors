<!-- @anchors
  code: LDCMS
  updated_at: 2026-11-19
  layer: comando
-->
# LoadCommand — what the machine lends its runs, and who holds it

> **Code**: `LDCMS`

## Overview

A run that waits looks like a hang unless something says what it waits for. `anchors load` reads
the machine's leases (LDQUS) and shows the capacity, what is in use, the runs holding the machine
and the runs waiting, from every worktree and project on the machine.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the leases file | the queue's JSON, or none | a file that is not one | this unit: the read error |

## Effects

| Effect | Description |
| --- | --- |
| `LDCMS-B01` | The command shows the capacity and what is in use, each run holding the machine — project, task, level, cost, for how long, pid — and each run waiting, oldest first, marking those that waited past their limit; an idle machine shows a dash for each. (`printLoad`) |

## Errors

none — a leases file that cannot be read is the command's error.
