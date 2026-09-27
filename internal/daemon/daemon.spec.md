<!-- @anchors
  code: DMSTD
  updated_at: 2026-09-26
  layer: infra
-->
# DaemonState — the background watcher's state files: PID, pause flag, log and meta

> **Code**: `DMSTD`

## Overview

The watcher runs in the background so the terminal stays free. Everything the other commands need to
know about it lives in plain files in the project's `.anchors/` state folder: which process is the
watcher (the PID file), whether it is paused (a flag file), where it logs, and since when and over which
root it runs (the meta file). Files are simple and inspectable: a person can see and fix the state by
hand.

This unit answers "is the watcher running?" honestly: a PID file that points at a process that has
exited is stale, and it is removed rather than reported as a live watcher. It stops the watcher,
refusing when none is running, and removes the PID file itself instead of trusting the watcher's loop
to do it. How a process is probed and terminated is platform-specific and belongs to `DMRND`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a folder where `.anchors/` exists or can be created | — | the caller (`anchors watch`) |
| the PID file content | a positive integer written by this unit | text, zero or a negative number (a hand-edited or corrupted file) | this unit: such a file reads as "not running" |

## Effects

| Effect | Description |
| --- | --- |
| `DMSTD-B01` | The state files are `.anchors/watch.pid`, `watch.log`, `watch.paused` and `watch.meta` under the project root. (`PathsFor`) |
| `DMSTD-B02` | "Running" answers the PID written in the PID file when that process is alive, and 0 when the file is missing, not a positive integer, or points at no live process. (`WritePID`) |
| `DMSTD-B03` | A PID file pointing at a process that has exited is removed when "running" is asked. |
| `DMSTD-B04` | Stopping when no watcher is running is refused with "watcher is not running". |
| `DMSTD-B05` | Stopping a running watcher terminates its process and removes the PID file. |
| `DMSTD-B06` | The watcher is paused exactly while the pause flag file exists: pausing creates it and resuming removes it. (`Pause`, `Resume`, `IsPaused`) |
| `DMSTD-B07` | Cleanup removes the PID file and the pause flag. |
| `DMSTD-B08` | The meta file holds `started=` with the start moment in RFC 3339 and `root=` with the project root, one per line; reading a missing meta gives empty text. (`WriteMeta`, `ReadMeta`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DMSTD-I01` | A PID file never outlives the answer "not running" for an exited process: after the check, the file is gone. | writes the PID of a process that has exited, asks whether the watcher runs, and checks the answer is 0 and the file is gone |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DMSTD-E01` | REF[DMSTD-B04]: stopping with no watcher is the failure B04 answers with "watcher is not running" | — | — |
| `DMSTD-E02` | The state folder cannot be created when the PID is written. | The error is returned and no PID file exists. | A watcher whose PID was not recorded is invisible to `stop` and `status`; the starter must know. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/daemon/daemon_unix.go` | `alive`, `terminate` | infra — the platform probe and termination (`DMRND`) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
