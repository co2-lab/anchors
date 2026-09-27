<!-- @anchors
  code: CHLGC
  updated_at: 2026-09-26
  layer: infra
-->
# CheckLog — the check's output mirrored to a file, so it can be reread without re-running

> **Code**: `CHLGC`

## Overview

`anchors check --all` takes minutes. Without a copy of its output, every question about the result
("what failed?", "which gates passed?") means running everything again just to reread what was already
printed, and a terminal with limited scroll, or an agent session, loses it. This unit duplicates the
check's standard output into a file under the project's `.anchors/` state folder, while the output
still goes to the screen.

The file is chosen by SCOPE: the full check and the check of changed files write to different files,
because the pre-commit runs the changed-files check on every commit and would otherwise erase the
expensive full snapshot at the moment it is most useful.

The file opens with a header that records the command, when it ran, the commit (HEAD) and the state of
the working tree. Without it the file lies by omission: yesterday's reading looks like today's. The
tree state is never claimed clean without having been counted.

The mirror is a convenience: if the file cannot be opened, the check still runs and writes only to the
screen.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a folder where `.anchors/` exists or can be created | a root where `.anchors` cannot be created | this unit: the mirror is not opened and the check continues |
| the dirty-file count | a count of modified files, or a negative value meaning "could not count" | — | the caller passes what the git reader returned; this unit renders every value |
| the HEAD and subject | a short hash and a subject, or empty when there is no commit to name | — | the caller; an empty HEAD is left out |

## Effects

| Effect | Description |
| --- | --- |
| `CHLGC-B01` | The full check is mirrored to `.anchors/check-all.txt` and the changed-files check to `.anchors/check-changed.txt`. (`ScopeName`) |
| `CHLGC-B02` | The header is written first, then everything the command writes to standard output is copied to the file. |
| `CHLGC-B03` | An output longer than the pipe's buffer is copied whole, and closing the mirror does not hang. |
| `CHLGC-B04` | When the mirror file cannot be opened, no mirror is returned, closing it and asking its path are safe (the path is empty), and standard output stays usable. |
| `CHLGC-B05` | The header records the command, the moment it ran (date, time and offset) and, when a HEAD is known, the HEAD with its subject; with no HEAD, no HEAD line is written. (`Header`) |
| `CHLGC-B06` | The tree line reads `clean` for zero modified files, `1 modified file` for one, and `N modified file(s)` for more. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CHLGC-I01` | The two scopes never overwrite each other: a changed-files check leaves the full snapshot intact. | runs a full mirror, then a changed-files mirror, and checks the full file still holds its output and the changed file does not |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CHLGC-X01` | The header never claims a clean tree it could not count: a negative count reads `unknown`. | A snapshot stamped "clean" without anyone having looked describes a state that never existed, and the reader has no way to suspect it. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CHLGC-E01` | REF[CHLGC-B04]: the state folder or the mirror file cannot be created, which B04 answers by running the check without a mirror | — | — <!-- @resilient: the mirror is a convenience copy of what standard output already shows, so a check without it loses nothing the person sees --> |

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
