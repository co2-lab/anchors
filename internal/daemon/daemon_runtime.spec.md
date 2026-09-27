<!-- @anchors
  code: DMRND
  updated_at: 2026-09-26
  layer: infra
-->
# DaemonRuntime — how each platform probes and terminates the background watcher

> **Code**: `DMRND`

## Overview

The watcher's state rules (`DMSTD`) need two things only the operating system can answer: is the
process with this PID alive, and how is it asked to end. Each platform supplies both, in its own file,
and the two files together must give `DMSTD` the same contract: a process that has exited is not
alive, so its stale PID file is cleaned; and stopping ends the process, after which `DMSTD` removes the
PID file itself.

On Unix-like systems a process is probed with signal 0, which delivers nothing and only tests the
process, and it is ended with SIGTERM: a signal the watcher's loop catches, so it can exit cleanly.
On Windows there is no SIGTERM: the process is killed outright, and the loop gets no chance to clean
up. That asymmetry is why `DMSTD` removes the PID file on stop instead of relying on the loop. The
Windows file is not built on this project's development host, so none of the rules below claims its
behaviour as proven; they are the contract proven on the Unix side, which the Windows side is written
to honour.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the PID | a positive PID read from the watcher's PID file | zero or a negative PID | `DMSTD`, which never probes a PID that is not positive |
| the process to terminate | a process found for the running watcher's PID | — | `DMSTD`, which only stops what it found running |

## Effects

| Effect | Description |
| --- | --- |
| `DMRND-B01` | A process that exists is reported alive, and one that has exited is reported not alive. |
| `DMRND-B02` | On Unix-like systems, terminating sends SIGTERM, a signal the watcher can catch to exit cleanly, not an uncatchable kill. |
| `DMRND-B03` | On Unix-like systems, a process the probe may not signal (the probe answers "operation not permitted", EPERM) exists, so it is reported alive: a live watcher of another user is never taken for a stale one. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DMRND-I01` | Whatever the platform's termination does, stopping the watcher leaves no PID file behind: the removal does not depend on the watcher's loop running its cleanup. | starts a child that is not the watcher (it never cleans anything), stops it through the state unit, and checks the PID file is gone |

## Errors

none — the platform functions answer a yes/no probe or pass the operating system's own error back to `DMSTD`, which decides what to do with it.

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
