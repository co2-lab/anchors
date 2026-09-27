<!-- @anchors
  code: GTMTG
  updated_at: 2026-09-26
  layer: infra
-->
# GitMeta — what git knows about the files: last commit dates, pending changes, HEAD, dirty count

> **Code**: `GTMTG`

## Overview

Several parts of Anchors need facts only git has: the date a file last changed (the source of truth for
the `updated_at` stamp, more reliable than a date kept by hand), whether a file has uncommitted edits
(a file being edited right now changed today, not at its last commit), which commit a report describes,
and how many files are modified in the tree (a report run on a dirty tree describes a state that no
commit holds).

Dates are days only, never times: a file edited at one hour and committed at another on the same day is
the same day. The bulk date reader answers every file of the repository with a single log walk instead
of one git call per file.

The unit never claims a fact it could not read. Outside a repository, or without git, "no pending
changes" and "could not ask" are different answers, and "zero modified files" is never reported when
the files could not be counted: a report stamped "clean" without anyone having looked describes a
snapshot that never existed.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | any folder; a folder outside a repository gets the "unknown" answers | — | this unit: every git failure becomes an explicit "not known" |
| the file | a path relative to the root | an absolute path | the caller (the map and the touch commands pass relative paths) |

## Effects

| Effect | Description |
| --- | --- |
| `GTMTG-B01` | Today is the system date as year-month-day. |
| `GTMTG-B02` | The last commit date of a file is the day (year-month-day) of the most recent commit that touched it; a file never committed, or a root outside a repository, gives no date. (`LastCommitDate`) |
| `GTMTG-B03` | The bulk date reader gives, for every file of the repository, the day of the most recent commit that touched it, and an empty map outside a repository. (`AllCommitDates`) |
| `GTMTG-B04` | The pending-change question answers whether the file has uncommitted edits (a new or edited file has them, a just-committed file does not), together with whether the question could be asked at all. (`UncommittedChanges`) |
| `GTMTG-B05` | The yes/no pending-change shortcut answers yes for a new or edited file in a repository and no for a just-committed one. (`HasUncommittedChanges`) |
| `GTMTG-B06` | HEAD is the short hash and the subject of the last commit; it is not known outside a repository or in a repository with no commit. (`Head`) |
| `GTMTG-B07` | The dirty count is the number of files with uncommitted changes in the whole tree: 0 in a clean repository, 1 after one new file. (`DirtyCount`) |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GTMTG-X01` | The dirty count never answers 0 when it could not count: outside a repository or without git it answers a negative number. | Zero says "clean tree"; stamping a report with it without having looked is the reassuring silence, the worst kind. |
| `GTMTG-X02` | The pending-change question never answers "known" when it could not ask: outside a repository it answers not known and no change. | "The file is clean" and "there is no repository to ask" would otherwise reach the caller as the same no, and that difference decides whether there is a verdict to give. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `GTMTG-E01` | REF[GTMTG-X01]: git failing to list the tree's status is the failure X01 answers with a negative count | — | — <!-- @resilient: the failure becomes a negative count that no caller can read as a clean tree, and the report header says the tree is unknown --> |
| `GTMTG-E02` | REF[GTMTG-X02]: git failing to report a file's status is the failure X02 answers with "not known" | — | — <!-- @resilient: the failure becomes not known, which the caller distinguishes from no change and answers by giving no verdict --> |

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
