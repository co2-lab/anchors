<!-- @anchors
  code: HDTHD
  updated_at: 2026-09-30
  layer: comando
-->
# HeaderDateTouch — bumps the header date of the files that changed, and only of those

> **Code**: `HDTHD`

## Overview

The header-date gate charges every changed file whose header date is not the day of the change. Keeping
up with it meant editing the date by hand, or a one-off script after every bulk change. The touch command
writes the day's date into the date field of the header of every file that changed.

It does not lie about what changed. The date says when the content changed, so a file whose only difference
from the last commit is its own date is not bumped again, a file with no difference is not touched, and a
file already at the date is left alone. Only the date inside the header counts: the header must sit at the
top of the file, as a run of comment lines in any dialect or a closed HTML comment block. Text that merely
looks like a header further down — a test fixture holding a header in a string, for instance — is never
rewritten.

By default the candidates are the files changed in the working tree against the last commit, plus the
untracked ones. With the staged flag they are the files in the index; each bumped file is staged again, and
a file that also has changes outside the index is skipped, because staging it would put changes nobody
chose into the commit. During a partial commit git prepares a second, real index, and that one is dated
too. Exclusion globs from the flag and from the project's configuration add up, and a dry run says what it
would bump. The installed pre-commit runs this bump by default; the project can turn it off.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | the root of a git repository with at least one commit | a directory outside any repository | this unit: it fails (`HDTHD-E01`) |
| the changed files | any file git reports as changed, added or untracked | a file that cannot be read | this unit: it is skipped and named (`HDTHD-E02`) |
| the date | a day in year-month-day form; today when not given | — | the caller: the flag is written as given |
| the exclusion globs | globs from the flag and from the configuration | — | this unit: a matching file is skipped and named |

## Effects

| Effect | Description |
| --- | --- |
| `HDTHD-B01` | Only a date inside the header at the top of the file is bumped: a comment block containing the header marker within the first ten lines, in any comment dialect, or an HTML comment block that closes; header-like text elsewhere is not a header. |
| `HDTHD-B02` | A real change and a new file are bumped to the date. |
| `HDTHD-B03` | An unchanged file, a date-only change and a file already at the date are not bumped. |
| `HDTHD-B04` | Without the staged flag the candidates are the worktree changes and the untracked files. |
| `HDTHD-B05` | With the staged flag the index is dated and re-staged. |
| `HDTHD-B06` | A staged file with changes outside the index is skipped, and its worktree change does not reach the index. |
| `HDTHD-B07` | In a partial commit the real index is dated too. |
| `HDTHD-B08` | The exclude globs of the flag and of the configuration add up. |
| `HDTHD-B09` | The dry run says what it would bump and writes nothing. |
| `HDTHD-B10` | Each bump is listed with its old and new date and each skip with its reason, then the total. |
| `HDTHD-B11` | Without a date the day of the run is written. |
| `HDTHD-B12` | The pre-commit bump is on unless the project turns it off. |
| `HDTHD-B13` | A project root below the repository's top considers only its own files, named from the project root, in both the worktree and the staged modes. |
| `HDTHD-B14` | Files or folders named on the command line narrow the touch to the changed files among them — a path relative to where the command runs, or absolute; one outside the project is an error. With none named, every changed file is a candidate. (`touchPaths`, `within`) |
| `HDTHD-B15` | Every file named on the command line gets a verdict: one with no change from HEAD is reported as skipped for it, one that does not exist as unreadable, one with no `@anchors` header as such — a named file never leaves the list silently. A named folder brings only its changed files. (`named`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `HDTHD-I01` | Touching twice bumps nothing the second time: once a file carries the date, it is already dated. | runs touch twice with the same date and reads the second run's output |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `HDTHD-X01` | A changed file with no dated header is neither touched nor listed. | The command dates headers; a file without one is not its business, and listing every such file would bury the ones it did date. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `HDTHD-E01` | Outside a git repository touch fails. | Error naming the git command that could not list the changes. | "Changed" is measured against the last commit; without one there is no change to date. |
| `HDTHD-E02` | A changed file cannot be read. | It is skipped and named as unreadable, and the other files are dated. | One unreadable file must not stop the dates of the rest. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Load` | config — the exclusion globs and the pre-commit switch |
| DEP2 | `internal/gitmeta/gitmeta.go` | `Today` | infra — the day of the run |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
