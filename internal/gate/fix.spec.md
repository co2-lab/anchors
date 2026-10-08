<!-- @anchors
  code: FXIXX
  updated_at: 2026-10-08
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
| `FXIXX-B09` | `--fix` writes the header a governed file lacks, at its top after any shebang: the `ref:` of the units the map ties it to (`UnitCodesOf`), or the `layer:` of a guide, a document or a test support file; to a header at the top with no identity it adds that line below `@anchors`, changing nothing written; a file with an identity, with no unit and no such layer, binary, or an executable script is left as it is. (`fixMissingHeader`) |
| `FXIXX-B10` | `--fix` gives a file whose header has an identity (or gets one) and no `code:` of its own a code generated from its name and its type, unique among the codes in the map, written below `@anchors` beside the identity; a header with its own code is left as it is. (`fixMissingHeader`) |
| `FXIXX-B11` | `FixWithConfig` runs the fixers with the project's config, so the ones that read the code by the dialect — the dependency chain's — have it; the config is theirs only for that run. |
| `FXIXX-B12` | Each repair says whether it added or removed lines (`LinesMoved`), so whoever carries the file's evidence knows whether the line-level signals still name the right lines. |
| `FXIXX-B13` | A run of the fixers that includes the date repair reads the repository's state once — the files with uncommitted changes and each file's last commit date —, and the date repair answers each file from it, with the same rule as one git question per file. |

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

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
