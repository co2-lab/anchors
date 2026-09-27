<!-- @anchors
  code: VPFVR
  updated_at: 2026-09-26
  layer: comando
-->
# VerifyPhaseFacade — one invocation per phase, delegated to the check pipeline

> **Code**: `VPFVR`

## Overview

The git hooks used to call the check once per file, a process each, while the other tools of the phase lived
elsewhere with their own idea of when to run. Verify is the single command a phase calls: the phase is the
argument, and the configuration is the only source of which gates it charges.

It is a facade and holds no verdict of its own. It collects the scope — the files staged in the git index,
the files named, or the whole project — and hands it to the check command in a child process, so there is
one pipeline and one truth about what passing means. What it adds is the phase: an automatic phase (any
phase other than manual) asks the check for the computable gates only and for the findings only, because a
hook can neither wait for an AI judgment nor dump a table of forty gates on every commit; the manual phase,
or no phase, asks for the full report the person typed the command to see.

The pre-commit over the index also dates the staged files before checking them: the date in a header is a
mechanical fact (the file changed today) and the gate that checks it should see the right date. A project
can turn that dating off; a staged file with changes outside the index is named and left undated; and a
dating failure only warns, since the gate still checks the dates.

The child's exit code crosses the process boundary: the check's "not governed" exit stays "not governed",
so a commit that touches only files the project does not govern is let through, and every other failure
stays a failure.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the scope | the staged files, one or more named files, or the whole project | none of the three | this unit: it refuses (`VPFVR-E03`) |
| the phase | pre-commit, pre-push, ci, manual, or none | — | the check: an undeclared phase selects no gate there |
| the git index | a repository with an index | a directory outside any git repository | this unit: the staged scope is refused (`VPFVR-E04`) |

## Effects

| Effect | Description |
| --- | --- |
| `VPFVR-B01` | The staged scope is the files the index added, copied, modified or renamed; deletions and untracked files are not in it. |
| `VPFVR-B02` | When nothing is staged, it says there is nothing to verify and succeeds. |
| `VPFVR-B03` | Verify hands the files to the check command in a child process, one changed-file argument per file (or the full sweep). |
| `VPFVR-B04` | An automatic phase asks the check for the computable gates only and for the findings only. |
| `VPFVR-B05` | The manual phase, or no phase, asks the check for the full report, judgment gates included. |
| `VPFVR-B06` | The pre-commit phase over the index dates the staged files first and says which files it dated. |
| `VPFVR-B07` | A project that turned pre-commit dating off gets no dating: the staged files keep their dates. |
| `VPFVR-B08` | A staged file that also has changes outside the index is named as not dated. |
| `VPFVR-B09` | The phase, the commit message file, the category, the skip-slow and the no-record choices reach the check unchanged. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `VPFVR-I01` | The check is never asked for both the full sweep and a list of files. | builds the check invocation from a file list and from the full sweep and inspects both |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `VPFVR-X01` | Verify holds no verdict of its own: what passes is what the check says passes. | Two implementations of "passing" would drift apart; the facade adds the phase and the scope, not a second ruler. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `VPFVR-E01` | The child check exits with the not-governed code. | Verify returns the not-governed error, which the program turns back into the same exit code. | A generic failure there made the pre-commit bar a configuration-only commit, the very case the code exists to let through. |
| `VPFVR-E02` | The child check fails with any other code. | Verify fails. | Only not-governed has its own handling; any other code swallowed would hide a real failure. |
| `VPFVR-E03` | Neither the staged scope, nor a file, nor the full sweep is given. | Error asking for one of the three. | A verify with no scope would pass having confronted nothing. |
| `VPFVR-E04` | The staged scope is asked for outside a git repository. | Error explaining which of git or the repository is missing. | There is no index to list, and the raw git error does not say why. |
| `VPFVR-E05` | Dating the staged files fails. | A warning is printed and the verify goes on. | The updated-at gate still checks the dates; a failed convenience must not block the commit. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `cmd/anchors/quality/check.go` | `ExitNotGoverned`, `errNotGoverned` | comando — the check the facade delegates to |
| DEP2 | `cmd/anchors/quality/touch.go` | `touchRun`, `touchOnPreCommit` | comando — dating the staged files |
| DEP3 | `internal/gitmeta/availability.go` | `Check`, `Explain` | infra — why the index cannot be read |
| DEP4 | `internal/config/config.go` | `Load` | config — the touch setting |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
