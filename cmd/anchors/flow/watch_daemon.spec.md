<!-- @anchors
  code: WTDMW
  updated_at: 2026-09-26
  layer: comando
-->
# WatchDaemon — the watcher started in the background survives the terminal that started it

> **Code**: `WTDMW`

## Overview

`anchors watch start` re-runs the CLI as a background watcher and returns the terminal at once.
That child must outlive the terminal: closing the window or pressing Ctrl+C in it must not take the
watcher down with it. How a process is detached from its terminal is a property of the operating
system, so the detachment is implemented once per platform, and exactly one implementation is
built for any target.

On unix-like systems the child starts in a new session, so it leads its own process group and the
terminal's hang-up and interrupt signals, which go to the terminal's group, do not reach it. On
Windows, which has no sessions, the child is created in a new process group, the closest portable
equivalent: the console's Ctrl+C sent to the parent's group does not reach it. That half is not stated as a rule here: its test only runs on Windows, and no run has proven it yet.

The detachment only prepares how the child will be started; starting it, writing its pid and its
metadata belong to the watch command.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the child process | a command not yet started | a process already running, whose attributes can no longer change | the watch command: it detaches before starting |

## Effects

| Effect | Description |
| --- | --- |
| `WTDMW-B01` | On a unix-like system, a detached child leads its own process group, distinct from the parent's. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `WTDMW-I01` | Exactly one detachment implementation is built for any target platform. | the package's file list is read for a unix-like and for a Windows target, and each holds exactly one of the two files |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `WTDMW-X01` | Detaching does not start the child. | The caller starts it and records its pid; a detachment that started the process would take that duty from the caller. |

## Errors

none — detaching only sets the attributes the child will be started with; it cannot fail, and the failure to start is the caller's.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `cmd/anchors/flow/watch.go` | `watch start` | comando — the only caller |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
