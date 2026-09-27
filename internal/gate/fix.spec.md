<!-- @anchors
  code: FXIXX
  updated_at: 2026-09-26
  layer: gate
-->
# Fix — the self-healer that applies the mechanical, safe repairs of `check --fix`

> **Code**: `FXIXX`

## Overview

Some findings have a repair that is mechanical and safe: writing the right date into the
`updated_at` field of a header needs no judgement, only git. `anchors check --fix` applies those
repairs, and this unit is what it calls. It holds a registry of fixers, one per check; a check with
no registered fixer is only reported, never repaired.

Today the registry holds one fixer, for the `updated-at-atual` check. It rewrites only the date
inside an existing field: the right date is today when the file has an uncommitted edit, and the
date of its last commit otherwise. When there is no way to know the right date (no git, no commit
and no edit) the file is left alone, because a guessed date would be worse than a stale one.

The repair is driven by the gate's scope, not by its verdict: every node the gate applies to is
handed to the fixer, and a file whose date is already right comes back unchanged, so it is not
reported. Each repair written, or attempted and failed, is returned for the caller to print.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the gates | the project's configured gates, any check | — | this unit: gates whose check has no fixer are ignored |
| the nodes | the map's nodes, with paths relative to the root | a node the gate does not apply to, or whose file is not on disk | this unit: out-of-scope nodes and unreadable files are skipped |
| the repository | a git working tree containing the files | a directory outside any repository | this unit: without git nothing is rewritten |

## Effects

| Effect | Description |
| --- | --- |
| `FXIXX-B01` | Only a check with a registered fixer is fixable (`Fixable`); today that is `updated-at-atual` alone. |
| `FXIXX-B02` | A stale `updated_at` on a committed file with no pending edit is rewritten to the date of the file's last commit, and the repair is reported as fixed, naming the gate and the file. |
| `FXIXX-B03` | A file with an uncommitted edit takes today's date. |
| `FXIXX-B04` | A date that already matches is left alone, and nothing is reported for that file. |
| `FXIXX-B05` | Only gates with a fixer are run, on every node the gate applies to whatever its verdict (the fixer decides whether there is anything to correct), and only for files present on disk. |
| `FXIXX-B06` | Outside a git repository the file is left untouched: there is no right date to find. |
| `FXIXX-B07` | A file that was never committed and has no pending edit is left untouched: there is nothing to compare against. |
| `FXIXX-B08` | The detail of each repair (fixed, or a write that failed) is written in the project's language, through i18n. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FXIXX-I01` | A repair replaces only the date inside the field; every other byte of the file stays as it was. | rewrites a stale date and compares the whole file with the expected content |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `FXIXX-X01` | Does not create a missing `updated_at` field; it only corrects the value of one that exists. | Adding the field is the author's work, charged by the header gate; a fixer that invented headers would write structure nobody reviewed. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FXIXX-E01` | Writing the repaired file fails. | The file is reported as not fixed, with the cause in the detail. | A repair that silently did not happen would let the caller claim the finding was fixed. <!-- @resilient: the write failure is not hidden: it becomes the file's result, not fixed and with the cause, which the caller prints; the other files are still repaired --> |
| `FXIXX-E02` | REF[FXIXX-B05]: a node whose file cannot be read is one of the files B05 leaves out, and nothing is reported for it | — | — <!-- @resilient: an unreadable file has no content to repair, and a node missing from disk is already reported by the map's own checks --> |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Gate` | core — the configured gates and their checks |
| DEP2 | `internal/gitmeta/gitmeta.go` | `UncommittedChanges`, `LastCommitDate`, `Today` | infra — the dates git knows |
| DEP3 | `internal/mapx/model.go` | `Node` | core — the nodes the gates apply to |
| DEP4 | `internal/i18n/i18n.go` | `T` | core — the localized detail of each repair |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
