<!-- @anchors
  code: MPSYN
  updated_at: 2026-09-28
  layer: comando
-->
# MapSyncForCommit — the commit carries the map a build of the commit makes

> **Code**: `MPSYN`

## Overview

The map committed with a change was the map of the moment before: the hook had just dated the
staged files, which changed their revision, and a node's date comes from the last commit that
touched the file — the one not made yet. The CI's `map build` of the commit differed from the
committed map on nearly every commit, and the next rebuild dropped the proofs of the dated files
over a date. The pre-commit hook now writes the map as it will be once the commit exists, and
stages it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project | a git work tree whose map git tracks | a map git does not track, or no map | this unit: left alone, nothing compares it |
| the dated files | the files the hook just dated, with their content before and after | — | `anchors touch` |

## Effects

| Effect | Description |
| --- | --- |
| `MPSYN-B01` | After the commit, a full `map build` of the committed files makes exactly the map the hook committed. |
| `MPSYN-B02` | A file the hook only dated keeps what was measured at its revision — its proofs survive. |
| `MPSYN-B03` | The map is built from the index: an unstaged edit and an untracked file do not enter it. |
| `MPSYN-B04` | A map git does not track, or no map, is left alone and nothing is staged. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MPSYN-E01` | The index cannot be read or the map cannot be written or staged. | The error comes back; the hook warns and does not block. | The map-freshness checks still hold the commit to the map; a sync that failed must not be what blocks it. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/scan/scan.go` | `WalkStaged`, `ShortHash` | scan — the files as the index has them |
| DEP2 | `internal/mapx/stamp.go` | `RebaseRev` | mapa — a date is not a new revision to measure |
| DEP3 | `internal/mapx/lock.go` | `WithLock` | mapa — the map's lock |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
