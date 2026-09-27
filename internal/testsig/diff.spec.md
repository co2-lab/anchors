<!-- @anchors
  code: DCLDF
  updated_at: 2026-09-26
  layer: infra
-->
# DiffChangedLines — which lines of which files a change added, read from a unified diff

> **Code**: `DCLDF`

## Overview

Diff coverage asks one question of the tests: do the lines this change added run under a test? This unit
answers the first half — which lines were added — by reading a unified diff. The other half, which lines are
covered, comes from the coverage report (see the coverage reading).

Only the NEW side matters: an added line is one that needs a test, while a removed line no longer exists.
The unit follows the file named by each new-file header and the line cursor set by each hunk header, and
records each added line under that file. A deleted file, and any file left with no added line, produce no
entry.

The diff comes either from the project's git working copy — against the current commit, or against a given
reference — or from a diff file, which is the fallback for a project without git: any tool that emits a
unified diff serves.

## Effects

| Effect | Description |
| --- | --- |
| `DCLDF-B01` | Each added line is recorded under the path of the file the preceding new-file header names; removed lines are not recorded. |
| `DCLDF-B02` | A hunk header sets the new-side line number the next added line gets, in both the `+start,count` and the bare `+start` form; a header without a new side gives line 0. |
| `DCLDF-B03` | The `a/` and `b/` prefixes and a tab-separated timestamp after the path are stripped, so the path is the file's path relative to the root. |
| `DCLDF-B04` | A deleted file (new side `/dev/null`) records nothing, and a file left with no added line is dropped from the result. |
| `DCLDF-B05` | Unchanged context lines advance the new-side line number without being recorded. |
| `DCLDF-B06` | `GitDiff`: From git, the working copy is compared against the current commit when no reference is given, and against the reference when one is. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the diff text | a unified diff, from git or from any `diff -u` source | a diff in another format, which yields no lines | the caller, choosing the diff source |
| the root | a git working copy, for the git path | a directory that is not a repository | this unit: git's refusal is returned as an error |
| the reference | any revision git understands, or empty | — | git, which rejects an unknown revision |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCLDF-I01` | Removed lines never shift the new-side numbering: a line added after removals keeps the number its hunk header gives it. | reads a hunk that removes one line and adds one, and verifies the added line's number is the hunk's new start |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCLDF-X01` | `ParseDiffFile`: Reading a diff file needs no git: the file is read as text, wherever it came from. | A project without git still gets diff coverage from any tool that writes a unified diff. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DCLDF-E01` | The root is not a git repository, or git fails. | The git error is returned and no lines are reported. | An empty result would read as "the change added nothing", and diff coverage would pass a change it never saw. |
| `DCLDF-E02` | The diff file cannot be read. | The read error is returned. | Same reason: a missing diff must not look like an empty change. |

## Dependencies

none — the unit reads text and runs git; it depends on nothing of the project.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
