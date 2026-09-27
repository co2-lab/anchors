<!-- @anchors
  code: FRZEX
  updated_at: 2026-09-26
  layer: comando
-->
# Freeze — the project is stopped in three layers with a written reason, and thawed by undoing exactly those layers

> **Code**: `FRZEX`

## Overview

Sometimes a whole project has to stop: a leaked credential to rotate, a release to protect,
a broken main to repair. `freeze` stops the work in three layers. It writes the freeze and
its reason into the project's configuration and pushes that change, so every agent and hook
that reads the configuration refuses to work; in github mode it also creates a repository
rule on the remote that blocks push and merge on every branch, and opens an issue with the
reason, so whoever hits the brake knows what happened. `thaw` undoes the three layers.

The reason is mandatory, because a freeze with no written reason is indistinguishable from
a broken configuration, and whoever hits it tries to work around it instead of reading. The
freeze is written as two lines at the top of the file, without reserializing it, so the
thaw can remove exactly those lines and give back the file as it was. A file that already
declares those keys has them replaced, because a duplicated key stops the configuration
from loading at all.

The commit and the push go out past the hooks, on purpose: the hooks the freeze has just
armed would refuse the freeze itself, and while the remote still says frozen they would
refuse the thaw. For the same reason, whoever fixes the problem passes: the admin bypasses
the remote rule. Each remote layer that fails is a warning that says which half of the brake
holds, never a failure that loses the local brake.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the reason | any non-blank text, including quotes and colons | a missing or blank reason | this unit: refuses before touching the file |
| the configuration | a loadable `anchors.yaml` | a missing or broken one | this unit: fails naming the load |
| the remote layers | github mode with a declared repository | any other workflow mode | this unit: only the file layer applies |

## Effects

| Effect | Description |
| --- | --- |
| `FRZEX-B01` | A missing or blank reason is refused and the file is not frozen. |
| `FRZEX-B02` | The freeze writes `enabled: false` and the quoted reason as the first two lines of the file, leaving the rest as it was; the file still loads, frozen, with the reason. |
| `FRZEX-B03` | The frozen file is committed and pushed past the hooks, with the subject `chore(anchors): freeze — <first line of the reason>`. |
| `FRZEX-B04` | With `--no-push`, nothing is committed or pushed and the frozen file is left as a local change. |
| `FRZEX-B05` | In github mode with a repository, the freeze creates the `anchors-freeze` rule on the remote and opens an issue with a fixed title, reporting its link; in any other mode it calls nothing on the remote. |
| `FRZEX-B06` | With `--no-ruleset`, the remote rule is not created and the issue is still opened. |
| `FRZEX-B07` | Freezing a project already frozen changes nothing and shows the existing reason. |
| `FRZEX-B08` | A remote layer that fails (push, rule or issue) is a warning naming it, and the command still succeeds with the local file frozen; the thaw warns the same way. |
| `FRZEX-B09` | The thaw removes the freeze lines, commits and pushes `chore(anchors): thaw`, deletes the rule found by its name and closes the issue found by its exact title. |
| `FRZEX-B10` | Thawing a project that is not frozen changes nothing and says so. |
| `FRZEX-B11` | When the remote has no freeze rule and no open freeze issue, the thaw deletes and closes nothing, without error. |
| `FRZEX-B12` | The deprecated Portuguese flag names (`--motivo`, `--sem-ruleset`, `--sem-push`) still drive the commands. |
| `FRZEX-B13` | A configuration that already declares a top-level `enabled:` or `freeze_reason:` (with any indented continuation lines) has those lines replaced, not stacked on: the frozen file holds exactly one of each key and still loads, frozen, with the new reason; the thaw then removes them and the project loads enabled. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FRZEX-I01` | A freeze followed by a thaw gives back the configuration byte for byte, when it declared no top-level `enabled:` or `freeze_reason:` of its own (otherwise see `FRZEX-B13`). | freezes a project, thaws it, and compares the file with the original |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `FRZEX-X01` | The configuration is never reserialized: the freeze only adds its two lines (after removing any top-level `enabled:`/`freeze_reason:` already there) and the thaw only removes top-level lines with those keys. | Reserializing would drop comments and reorder keys, turning a brake into an unreviewable diff. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FRZEX-E01` | `anchors.yaml` is missing or does not load, for freeze or thaw. | The command fails with "load anchors.yaml". | There is no project to freeze, and writing a new file would invent one. |
| `FRZEX-E02` | REF[FRZEX-B08]: a failing remote layer is the failure this command handles, by warning and keeping the local brake | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Load`, `Frozen`, `FreezeReasonText`, `GitHubMode` | config — the freeze state and the workflow mode |
| DEP2 | `cmd/anchors/common` | `AliasDeFlag`, `ResolveAliases`, `FirstLineOfReason` | comando — the deprecated flag names and the commit subject |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
