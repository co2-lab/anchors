<!-- @anchors
  code: MGCMM
  updated_at: 2026-10-03
  layer: comando
-->
# MigrateCommand — the command the format error promises, bringing the map and the config up to this binary's format

> **Code**: `MGCMM`

## Overview

When a binary meets a map or a config written in an older format, its error tells the user
to run `anchors migrate`. This unit is that command: an error that promises a command which
does not exist is worse than no error, because whoever reads it tries, fails, and stops
trusting the next message.

The command rewrites the two versioned files of a project, the map and the config, into the
format the running binary reads, through the migration steps of the migration package. It
reports each file on its own line: already current, migrated (with the format it went from
and to), or skipped with the reason it could not be read. For every file that changed it
lists the renamed keys, one per line with the number of occurrences, so whoever reviews the
diff knows what to expect before opening it.

A dry run says what would change and writes nothing. A real run ends by reminding that the
files are versioned and the migration must be committed, since a migration left on one
machine makes the next agent migrate again and produce the same diff.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a directory; either file may be missing | — | this unit: a missing file is reported and skipped |
| each file's format | a format at or below the binary's | a file whose format cannot be determined | the migration package reports it; this unit prints it and moves on |
| the project's config, crossing format 7 | an `anchors.yaml` this binary reads | a missing or unreadable config | this unit: no file is given a code, the reason is printed on a `·` line, and the rest of the migration goes on |

## Effects

| Effect | Description |
| --- | --- |
| `MGCMM-B01` | One run brings both the map and the config to the binary's current format, each reported as migrated from its old format. |
| `MGCMM-B02` | A file that is missing or cannot be read is reported on a `·` line with the reason, and the other file still migrates. |
| `MGCMM-B03` | A file already in the current format is reported as such and left untouched. |
| `MGCMM-B04` | For each migrated file, the renamed keys are listed in alphabetical order, each with its number of occurrences. |
| `MGCMM-B05` | With `--dry-run` the command reports what would be migrated, writes nothing, and says that nothing was written. |
| `MGCMM-B07` | When the project crosses a step that renames code letters, the codes of its plans (files of a `kind: plan` layer), flows (`.flow.md`) and actions (`.action.md`), found by the code their header declares, are rewritten in every versioned text file, each file listed with its rewrites; with `--dry-run` they are listed and not written; codes of other units are not touched, and a second run rewrites nothing. The installed pipelines are not rewritten: the command points to `anchors doctor --fix`, which updates the ones nobody edited. |
| `MGCMM-B08` | When the project crosses format 7, each four-character code a file owns is widened to five — the old code kept as the prefix — wherever it is written: the project's files, the names of the files that carry it (with or without a `recode:` block), the config and the map; the config's `code_lengths` becomes `[5]`. |
| `MGCMM-B09` | When the project crosses format 7, every governed file that can carry a comment and has no `code:` gets one in its `@anchors` header, generated from its name and its type (layer or kind), unique in the project; the `ref:` it had stays. |
| `MGCMM-B10` | When the project crosses format 7, each file the migration rewrote and the map had measured at the content it found has its measurements carried to its new revision. |
| `MGCMM-B11` | When the project crosses format 7, each renamed code is appended, old → new under the date, to `anchors.renames.yaml` (`RenamesFile`). |
| `MGCMM-B12` | When the project crosses format 7, a file other than the spec that carries, as its own `code:`, the code its unit's spec owns has that line turned into `ref: <code>` with a code of its own above it. |
| `MGCMM-B06` | After a real migration the command tells the user to commit it. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MGCMM-I01` | The migration is idempotent: a second run reports the files as current, asks for no commit, and leaves them byte for byte as the first run wrote them. | migrates, runs again, and compares the output and the file |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MGCMM-X01` | The command does not know which keys each format renamed: the steps belong to the migration package, and this unit only applies them and reports. | A format change is declared once, beside its reason; a second copy here would drift. |
| `MGCMM-X02` | The command never commits: it only reminds. | The commit belongs to whoever ran it, who reviews the diff first. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MGCMM-E01` | REF[MGCMM-B02]: a file that cannot be read is the failure this command handles, by reporting it and migrating the other file | — | — |
| `MGCMM-E02` | Crossing format 7, a project file, the config or the map cannot be written | the command fails with the write error; the map and the config go back to the format they had, and what was already written stays | run again after fixing the cause, it crosses the step again and codes only what is left |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/migra/migrate.go` | `MigrateFile` | apoio — the migration steps |
| DEP2 | `internal/mapx/model.go` | `DefaultPath`, `FormatoAtual` | mapa — the map's path and the current format |
| DEP3 | `internal/config/config.go` | `AbsRoot`, `DefaultFile` | config — the config's path |
| DEP4 | `internal/recode/plan.go` | `BuildPlan` | apoio — widening a code wherever it is written, file names included |
| DEP7 | `internal/recode/rewrite.go` | `Rewrite` | apoio — the code rewritten in the config and the map |
| DEP5 | `internal/migra/filecodes.go` | `WidenedCode`, `FileCode`, `CanCarryCode`, `WithHeaderCode` | apoio — what a code becomes and where the header line goes |
| DEP6 | `internal/scan/scan.go` | `Walk`, `ScanPaths` | scan — the governed files and their revisions |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
