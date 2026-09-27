<!-- @anchors
  code: INHKN
  updated_at: 2026-09-26
  layer: comando
-->
# InstallHooks — the git hooks that confront every commit and push with the gates and the freeze, installed without taking a hook the user wrote

> **Code**: `INHKN`

## Overview

The gates only protect a project if they run where the work leaves a machine. `install-hooks`
writes three git hooks and registers two merge drivers.

The pre-commit first checks whether the project is frozen on the remote, because learning it
at push time costs the whole afternoon's work; the remote is consulted at most once every ten
minutes, and only a "not frozen" answer is trusted from that cache, so a thawed project never
keeps someone blocked. It then runs the pre-commit phase of the gates once over the staged
files. A "not governed" verdict (nothing staged matches a layer, such as a lockfile change)
passes. A failure is not final there: the commit message does not exist yet, and it may
declare a waiver, so when the commit-msg hook is present the pre-commit reports and defers.
The commit-msg hook checks the message format first, then re-runs the gates with the message
in hand and blocks what no waiver covers. The pre-push refuses a push while the remote is
frozen and refuses a binary older than the minimum version the remote declares, so a fix can
reach everyone before work proceeds.

The merge drivers let git merge the map's judgment stamps and the plans' progress items
instead of merging them as text. Installing is safe to repeat: the hooks anchors wrote are
recognized and updated, and the attribute lines are written once. A hook the user wrote is
theirs and is not replaced without `--force`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a git repository or worktree holding `anchors.yaml` | a directory with no `anchors.yaml`, or outside git | this unit: refuses with the reason |
| the hooks directory | `core.hooksPath` (absolute or relative to the root), or the common git directory's hooks | — | this unit: asks git |
| existing hooks | hooks anchors wrote (current or legacy marker) and hooks the user wrote | — | this unit: tells them apart by the marker |

## Effects

| Effect | Description |
| --- | --- |
| `INHKN-B01` | The hooks go to `core.hooksPath` when set, resolved against the root when relative, and otherwise to the `hooks` directory of the common git directory. |
| `INHKN-B02` | A fresh install writes the pre-commit, commit-msg and pre-push hooks, executable, each holding the managed script. |
| `INHKN-B03` | A pre-commit the user wrote refuses the whole install unless `--force`, and nothing is written; a user's commit-msg or pre-push is left in place with a warning while the rest installs; `--force` replaces them all. |
| `INHKN-B04` | A hook anchors wrote, recognized by the current marker or the legacy commit-msg marker, is replaced on reinstall without `--force`, and is not reported as foreign. |
| `INHKN-B05` | The installed pre-commit refuses a commit while the remote's configuration is frozen, showing the freeze reason. |
| `INHKN-B06` | When the gates answer that nothing staged is governed, the pre-commit says there is nothing to confront and lets the commit proceed. |
| `INHKN-B07` | When the gates fail and the commit-msg hook is installed, the pre-commit reports the failure and how to declare a waiver, and the commit-msg re-runs the gates with the message and blocks the commit. |
| `INHKN-B08` | The installed pre-push refuses a push while the remote's configuration is frozen. |
| `INHKN-B09` | The installed commit-msg refuses the commit when the message check refuses the subject. |
| `INHKN-B10` | The installed pre-push refuses a push when the local binary is older than the minimum version the remote's configuration declares, and lets that exact version through. |
| `INHKN-B11` | The install registers the map and progress merge drivers in the repository's git config and adds their attribute lines to `.gitattributes`. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `INHKN-I01` | However many times the install runs, each merge attribute line appears once in `.gitattributes` and the user's own lines are kept intact. | installs twice over a `.gitattributes` without a final newline and counts the lines |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `INHKN-X01` | A hook the user wrote is never replaced without `--force`. | Their hook is their work; silently replacing it would switch off whatever it guarded. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `INHKN-E01` | The root has no `anchors.yaml`. | The install fails pointing at `anchors init`, and no hook is written. | A hook in an ungoverned project would run gates over nothing and train people to skip it. |
| `INHKN-E02` | The root is not inside a git repository. | The install fails explaining that installing the pre-commit needs git. | Git's raw error about a missing repository names neither the hook nor the fix. |
| `INHKN-E03` | REF[INHKN-B03]: a pre-commit the user wrote is the failure this command handles, by refusing and pointing at `--force` | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gitmeta` | `Check`, `Explain` | apoio — explains a directory outside git |
| DEP2 | `internal/mapx/model.go` | `DefaultPath` | mapa — the map's attribute line |
| DEP3 | `cmd/anchors/ops/commit_msg.go` | the message check the commit-msg hook runs | comando — CMMSC |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
