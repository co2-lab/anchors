<!-- @anchors
  code: RCPLR
  updated_at: 2026-10-04
  layer: infra
-->
# RecodePlan — planning and applying the rename of a code across the whole project

> **Code**: `RCPLR`

## Overview

`anchors recode OLD NEW` renames a unit's identity code everywhere it appears. This unit plans the
rename over the whole project, without writing anything, so the command can show a dry run; and it
applies a plan when asked.

Planning refuses what cannot be a rename: a malformed code, a code renamed to itself, a new code another
unit already owns (the rename would merge two identities), or an old code that appears nowhere (nothing
to rename). For every project file that holds the old code, the plan carries the
rewritten content (through `RCRWR`) and the classified occurrences. When the project declares its own
recode dialect, the plan also rewrites the testID prefix derived from the code and lists the files whose
names carry the code, to be renamed (through `RCDLR`). A testID prefix that does not follow the
convention is never guessed: the plan warns that divergent testIDs exist and leaves them alone.

Applying writes each file with its original mode and performs the renames. Inside a repository a rename
goes through git, so the index records it and the history is kept; outside a repository, or for a file
git does not track yet, a plain rename is the right behaviour. A real refusal from git is never bypassed:
moving a file behind git's back, in a command that renames in mass, would leave the index diverging from
the disk in silence.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the old and new codes | upper-case codes of the project's declared lengths, distinct, the new one owned by no unit | malformed or equal codes, a new code already in use | this unit: planning is refused |
| the project | the files the project's layers declare, read from disk | files outside the layers | the scanner; file renames also look outside the layers (snapshots, flows) |
| the recode dialect | the project's `recode:` block, or none | — | the configuration |

## Effects

| Effect | Description |
| --- | --- |
| `RCPLR-B01` | A malformed old or new code is refused before anything is read. |
| `RCPLR-B02` | Renaming a code to itself is refused. |
| `RCPLR-B03` | The plan holds, sorted by path, every project file that holds the old code, with its rewritten content; files without the code are left out, and the total counts the content replacements. (`BuildPlan`) |
| `RCPLR-B04` | With the testID convention declared, the testID prefix derived from the old code is rewritten to the new one, counted apart from the content replacements; without the dialect no testID is touched. |
| `RCPLR-B05` | With file patterns declared, every file whose name matches them for the old code is planned for renaming to the new code; without the dialect no file is renamed. |
| `RCPLR-B06` | When the expected testID prefix appears nowhere but files holding the old code carry testIDs, the plan warns about a divergent prefix and does not rewrite it. |
| `RCPLR-B07` | An old code that appears in no file content, no testID and no file name is refused. |
| `RCPLR-B08` | Applying writes every planned file with its original mode and performs every rename, and answers how many files were written and renamed. |
| `RCPLR-B09` | Inside a repository a tracked file is renamed through git, so the index records the rename; outside a repository it is a plain rename. |
| `RCPLR-B10` | Inside a repository a file git does not track is moved by a plain rename, and a reported success means the file really moved. |
| `RCPLR-B11` | A target code another unit already owns — declared in its header, or the prefix of one of its scenario codes — is refused naming that file; a code that merely contains the target is not a collision. |
| `RCPLR-B12` | The refusal of a malformed code names the lengths the project declares in `code_lengths`, the same ones the validation applies. |
| `RCPLR-B13` | `BuildBatchCited` plans the rename of many codes at once, as a migration makes it: one walk of the governed files, each rewritten once for every code where it is cited as a code (`NewCitedSet`), the files whose names carry a code renamed with it (the given name patterns), the counts per code, and the bare mentions left, per code and file; `Batch.Apply` writes the files and makes the renames as `Plan.Apply` does. |
| `RCPLR-B14` | A batch never rewrites a binary as text — a file that is not valid UTF-8 or holds a NUL byte, such as a VR baseline —: it is renamed with its code when its name carries one, byte for byte. |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RCPLR-X01` | A refusal from git is returned as an error naming the move, and the file is not moved by other means. | Moving anyway leaves the index diverging from the disk, in silence, in a command that renames in mass — and the trace that would allow undoing is exactly what is not created. |
| `RCPLR-X02` | Planning writes nothing to the project. | The dry run must be safe to run at any time; only applying writes. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RCPLR-E01` | REF[RCPLR-B01]: a malformed code is refused by B01 | — | — |
| `RCPLR-E02` | REF[RCPLR-B02]: the same code on both sides is refused by B02 | — | — |
| `RCPLR-E03` | REF[RCPLR-B07]: an old code with nothing to rename is refused by B07 | — | — |
| `RCPLR-E04` | REF[RCPLR-X01]: git refusing a move is surfaced by X01 | — | — |
| `RCPLR-E05` | A planned file cannot be written while applying. | Applying stops with an error naming the file, and answers how many files were written before it. | The caller must know the rename is partial, and where it stopped. |
| `RCPLR-E06` | REF[RCPLR-B11]: a target code owned by another unit is refused by B11 | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/recode/rewrite.go` | `Rewrite`, `Find`, `ValidCode` | infra — the text rewrite (`RCRWR`) |
| DEP2 | `internal/recode/dialect.go` | `RewriteTestIDs`, `FileMatchesCode`, `RenameFilePath` | infra — the project's dialect (`RCDLR`) |
| DEP3 | `internal/scan/scan.go` | `Walk` | scan — the project's files |
| DEP4 | `internal/gitmeta/availability.go` | `Check` | infra — whether a move goes through git (`GTAVG`) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
