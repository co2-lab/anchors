<!-- @anchors
  code: FRGTF
  updated_at: 2026-10-20
  layer: gate
-->
# Frozen gates — a frozen part is quiet, and told when touched

> **Code**: `FRGTF`

## Overview

A freeze keeps a part of the project as it is by decision, until someone thaws it (DESIGN-frozen-units.md).
The map says which freeze covers what (`FRSCF`); this unit makes the gates honor it. A gate on a frozen part
answers `frozen`, with the freeze's reason, counted apart — neither failure nor pending. A frozen rule
leaves every question asked of a rule, the way a rule marked `@no-scenario` or `@retired` does. Two gates
still speak about frozen parts: `freeze-consistent` confronts the marks with the doctrine's records, and
`frozen-touched` tells a change to a frozen part — it informs, and never bars.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the map | the graph the run confronts, with its freezes | no map | this unit: nothing is frozen, and the freeze gates are pending |
| a target | a node of the run | — | the gate engine |

## Effects

| Effect | Description |
| --- | --- |
| `FRGTF-B01` | A gate on a node a freeze covers whole answers `frozen` with the freeze as it reads — its code with its recorded reason, or the reason —, and does not confront it; the freeze's own gates still run on it; a gate of the set (batch or project) leaves frozen targets out, and does not run when every target is frozen. (`Frozen`, `setFreezes`, `frozenResult`, `withoutFrozen`) |
| `FRGTF-B02` | A rule a freeze covers leaves every question asked of a rule — its scenario, its test, what it uses, its realizer —, and a doctrine's freeze entries ask for no realizer and are no rule another doctrine could repeat. (`ruleIsFrozen`) |
| `FRGTF-B03` | `freeze-consistent` fails naming each mark whose freeze code no doctrine records, with its place, and each recorded freeze no mark carries, with its document; it passes when every code is recorded and every record marked, skips when there is no freeze, and does not confront a mark that carries a reason. (`checkFreezeConsistent`) |
| `FRGTF-B04` | `frozen-touched` tells a frozen file with an uncommitted change — or a file with frozen rules, naming each rule and its freeze — as a declared divergence, which informs and never bars; an unchanged frozen part passes, and a part with no freeze is skipped. (`checkFrozenTouched`) |
| `FRGTF-B05` | The run's profile counts the `frozen` verdicts of each gate apart, and a frozen verdict neither fails nor bars. (`GateSummary`) |
| `FRGTF-B06` | `freeze-consistent` also names each glob of anchors.yaml's `frozen_paths` that matches no file of the map, with its freeze: a folder that freezes nothing is a glob mistyped, or a folder that moved; a project with frozen folders and no other freeze is confronted, not skipped. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FRGTF-I01` | Thawing a part gives back the verdicts it had: the freeze only changes what is asked while it holds. | a failing gate on a node, frozen then thawed: frozen, then the same failure |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FRGTF-E01` | `freeze-consistent` runs with no map. | Pending, naming the missing map. | With no map there are no marks and no records to confront, and passing would say they agree. |
