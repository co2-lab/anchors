<!-- @anchors
  code: MPLCK
  updated_at: 2026-10-08
  layer: mapa
-->
# MapLock — the map changed by one writer at a time, each applying only what it changes

> **Code**: `MPLCK`

## Overview

The map is one file, and every command that changes it rewrote it whole: it read the map, changed
its part, and wrote everything back. Two processes doing that at once lost one of the changes —
reported from a project running three mutation ingests in parallel with an `anchors test` and the
pre-commit check: ten ingestions of one agent disappeared, and the map went back to the old scores.

Each command now says only what it changes and applies it to the map as it is on disk at that
moment, under an exclusive lock: re-read, apply, write (`Update`). The heavy work — running a
suite, parsing a report, running the gates — happens before, outside the lock, so the lock is held
for the time of a read and a write, not of a run.

The lock is a file beside the map (`anchors.graph.yaml.lock`), created only if it does not exist,
holding its owner's pid and host. A lock left by a process that died holding it is taken over.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the map path | a path whose directory can be written | a directory that cannot be written | this unit: the error names the lock file |
| the change | a function that changes the graph it is given, or refuses | — | the caller: it must change only its own part |

## Effects

| Effect | Description |
| --- | --- |
| `MPLCK-B01` | Taking the lock creates the lock file beside the map, holding the owner's pid and host; releasing it removes the file. (`Lock`, `LockPath`) |
| `MPLCK-B02` | While the lock is held, another writer waits, and takes it once it is released. |
| `MPLCK-B03` | A lock whose owner is a dead process of this host is taken over at once; a lock older than a minute is taken over whatever its owner; a lock of a live process, or of another host, is not. |
| `MPLCK-B04` | A change is applied to the map re-read under the lock and written: writers in parallel processes, each changing its own part, all reach the map. (`Update`) |
| `MPLCK-B05` | A function run under the lock releases it when it returns, with an error or without. (`WithLock`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MPLCK-E01` | Another writer holds the lock past the timeout. | Taking the lock fails saying another process is writing the map, naming the holder and the lock file; the holder's lock is left as it is. | Waiting forever would hang a hook; taking a live lock would bring the race back. |
| `MPLCK-E02` | The map cannot be read under the lock, or the change refuses. | `Update` returns the error, writes nothing, and releases the lock. | A change applied to no map, or a half-applied one, is not a map. |
| `MPLCK-E03` | The lock file cannot be created for a reason other than existing (the directory cannot be written). | The error names the lock file. | Retrying until the timeout would hide the real cause for two minutes. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
