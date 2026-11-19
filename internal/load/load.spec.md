<!-- @anchors
  code: LDQUS
  updated_at: 2026-11-19
  layer: infra
-->
# LoadQueue — the machine's runs admitted by what they cost

> **Code**: `LDQUS`

## Overview

A run states what it costs in CPUs and memory (LDCFS), and starts when that fits in what the
machine has free. The leases of every run on the machine — every worktree, every project — live in
one file (`~/.anchors/load/leases.json`, or `$ANCHORS_LOAD_DIR`), under a lock, so the processes see
each other. A light run passes while a heavy one holds most of the machine. A run that waited past
its limit stops the admission of every run that came after it until it fits, so a stream of light
runs never starves a heavy one. A lease whose process died is reclaimed.

The capacity is the machine's CPUs and memory minus a reserve left to whoever uses it: one CPU, and
2 GB or a quarter of the memory, whichever is less. `$ANCHORS_LOAD_CAPACITY` (`cpu=8,memory=16GB`)
sets it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a request | a project, a task, a level, its cost and its wait limit | — | the caller |
| the leases file | the JSON this unit writes; none | a file that is not one | this unit: an error, and the caller runs at once |
| `$ANCHORS_LOAD` | `off`, `0`, `false` turn the queue off for the process and its children | — | the commit hook sets it |
| `$ANCHORS_LOAD_CAPACITY` | `cpu=<n>,memory=<size>`, either part | another key | this unit: the machine's detected capacity |

## Effects

| Effect | Description |
| --- | --- |
| `LDQUS-B01` | A run starts when its cost fits in CPU and in memory beside the leases; with nothing running it starts whatever its cost, so a run costlier than the machine runs alone; an unread memory does not bound the queue. (`Admit`, `Used`) |
| `LDQUS-B02` | A run that waited past its limit stops the admission of every run that came after it, and starts once it fits; within its limit it holds no one. |
| `LDQUS-B03` | The queue admits a project that declares `load:`, unless `$ANCHORS_LOAD` turns it off or the run is in CI and `load.ci` is not set. (`Enabled`) |
| `LDQUS-B04` | The capacity is the machine's CPUs and physical memory minus the reserve, or what `$ANCHORS_LOAD_CAPACITY` says; sizes are written in the largest unit they reach. (`MachineCapacity`, `ParseCapacity`, `Size`, `totalMemory`) |
| `LDQUS-B05` | Runs share the machine's leases through one file under a lock: a run that does not fit is recorded as waiting and says, then every 30 seconds, what holds the machine; a released lease admits it. (`Acquire`, `Read`, `Sorted`, `Dir`) |
| `LDQUS-B06` | A lease or a waiter whose process died is dropped at the next read under the lock. |

## Errors

none — a queue that cannot be read lets the run start at once, saying so (LDQRS): the queue spaces
runs out, it must not stop them.
