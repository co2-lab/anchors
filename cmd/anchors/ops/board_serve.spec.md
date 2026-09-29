<!-- @anchors
  code: BRSRB
  updated_at: 2026-09-29
  layer: comando
-->
# BoardServe — the board page served locally with live state, read from the host only when something changed

> **Code**: `BRSRB`

## Overview

The published board is a snapshot: the pipeline runs, writes the board data, publishes it,
and between two runs what everyone sees is the past. In a project with many agents
delivering, minutes are enough for the snapshot to lie about who holds which card.

`board serve` brings the same board page up on a local server, next to board data read live
from the repository's issues. It is the same page and the same data contract as the
published board, so the two renderings never diverge; the live data only adds that it is
live, which makes the page count agent activity from now rather than from the time of a
cached read.

The reads are economical because the host's rate limit is shared by every agent of the
project. A floor between two reads returns the last read without asking anything. Past the
floor, the command asks one cheap question — the most recently updated card, sorted by
update time — and re-reads everything only when that card is newer than the newest one
already served. The full sweep reads every page of the issues, drops pull requests, and adds
the comments of the open cards (from which the board takes each card's owner). When a sweep
fails after a good read, the board keeps serving the good read; with nothing to serve, the
data route answers with the error.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the repository | `owner/name`, given or discovered from the current clone | a directory that is not a clone, with no `--repo` | this unit: fails pointing at `--repo` |
| the floor | any duration | — | this unit: zero means always ask |
| the port | a port the machine can listen on | an invalid or taken port | the listener fails and the command returns its error |
| the host's answers | the issues pages, the newest update time, the comments pages | refusals and malformed answers | this unit: refusals are reported or tolerated as each effect says |

## Effects

| Effect | Description |
| --- | --- |
| `BRSRB-B01` | The board data is the collected items wrapped with `live: true` and the read time in UTC as both `generated` and `takenAt`. |
| `BRSRB-B02` | Within the floor after a read, the last read is returned and the host is not called. |
| `BRSRB-B03` | Past the floor, the command asks for the most recently updated card and sweeps again only when its update time is newer than the newest already served. |
| `BRSRB-B04` | When a sweep fails after a good read, the good read is served again, without error. |
| `BRSRB-B05` | The full sweep reads every page of the repository's issues, stitches them, drops pull requests, and takes each card's owner from the comments of the open cards; the pipeline's own jq expression runs inside the binary, so no `jq` has to be installed. |
| `BRSRB-B06` | The comments of the open cards are read page by page following the cursor; a refused or malformed answer, or a repository name without owner and name, yields no comments instead of breaking the board. |
| `BRSRB-B07` | Without `--repo`, the repository comes from the current clone; the command prints the address, the repository and the floor before serving. |
| `BRSRB-B08` | The root route serves the published board page as HTML, and the data route serves the board data as JSON that the browser must not cache. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `BRSRB-I01` | The baseline of the incremental question is the newest update time among the items of the board being served. | reads a board and compares the baseline with the newest `updated` of its items |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `BRSRB-X01` | The page served is the pipeline's own board page; this unit renders nothing itself. | Two renderings of the board would diverge, and the local one would lie in another way. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `BRSRB-E01` | The host refuses the full sweep. | The collection fails with "collect from GitHub" and the host's own message. | A refusal without the host's reason (a rate limit, a permission) leaves nothing to act on. |
| `BRSRB-E02` | The sweep fails and there is no previous read. | The data route answers 502 with the error as a JSON `error` field. | There is nothing true to serve, and an empty board would read as "no cards". |
| `BRSRB-E03` | No `--repo` is given and the repository cannot be discovered from the clone. | The command fails saying to declare `--repo <owner/name>`. | Guessing a repository would serve another project's board. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/initx` | `BoardHTML`, `BoardCollectJQ` | apoio — the published board page and its data contract |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
