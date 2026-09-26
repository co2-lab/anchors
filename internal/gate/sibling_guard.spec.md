<!-- @anchors
  code: SBGRD
  updated_at: 2026-09-26
  layer: gate
-->
# SiblingGuard — sibling functions treat the same parameter consistently

> **Code**: `SBGRD`

## Overview

Confronts a module against its own internal asymmetry: **when two sibling functions guard
a parameter and a third does not, the one that does not is almost always forgetfulness,
not decision.**

The motivating case was real. A versioning module exported three functions over the same
history. The TWO that wrote filtered the history by key before deciding; the one that READ
did not filter — and it was exactly the one receiving the multi-key array straight from
the repository. Asking for one key's version returned another key's, in silence. It passed
11 green gates and 16 tests.

**The asymmetry is the signal**, and it is what makes this detectable without understanding
the domain. The gate does not know what the guard does, only that the siblings apply it and
one does not. That is also why it cannot be a lint rule: nothing in the syntax is wrong.

Deliberately CONSERVATIVE. It accuses only when three conditions hold at once: three or
more exported functions receive the same parameter name, the MAJORITY applies a
recognisable guard over it, and at least one applies none. Below that bar it stays silent,
because a false positive here teaches the team to ignore the gate — and a gate that is
ignored defends nothing.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the stack dialect | declared in the project's Structure | — (absent is a case, not an error) | this unit: with no dialect the verdict is undetermined, never approval |
| the exported functions | whatever the declared dialect recognises as exported | a construct the dialect does not reach | the dialect declared by the project, never a pattern invented by the gate |
| the parameter | any name shared by three or more siblings | a name fewer than three siblings share | this unit, by the conservatism bar |

## Effects

| Effect | Description |
| --- | --- |
| `SBGRD-B01` | An artifact that is not code leaves without a verdict: the gate reads function bodies. |
| `SBGRD-B02` | Without a declared dialect the verdict is UNDETERMINED, and says so — reading code it cannot recognise would stamp what it never checked. |
| `SBGRD-B03` | Fewer than three siblings on the same parameter is left alone: below the bar, asymmetry is not evidence. |
| `SBGRD-B04` | When the majority guards and at least one does not, the one that does not is ACCUSED. |
| `SBGRD-B05` | When every sibling guards, nothing is accused — there is no asymmetry to report. |
| `SBGRD-B06` | When no sibling guards, nothing is accused either: a module that never guards is a decision, not an oversight. |
| `SBGRD-B07` | The verdict NAMES the function that fails to guard and the parameter at stake, so the reader does not diff the module by hand. |
| `SBGRD-B08` | A waiver declared on the function WITH A WRITTEN REASON silences the accusation — the sibling that legitimately delegates the check says so where whoever reads the function will see it. A bare marker does not waive: a waiver with no why is the silence the gate exists to end. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SBGRD-I01` | The gate never judges WHAT the guard does. It reads that the siblings apply one and that one does not — understanding the guard would require understanding the domain, and that is what makes this detectable at all. | confronts a module whose guards differ in content and verifies only the absence is reported |
| `SBGRD-I02` | The three conservatism conditions hold TOGETHER. Any one of them alone would produce the false positive that costs the gate its credibility. | drops each condition in turn and verifies nothing is accused |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SBGRD-X01` | Does not invent what an exported function or a guard looks like. | Both come from the dialect the project declares. A pattern hardcoded here would recognise one ecosystem and report green over every other, which is the worst failure a measuring instrument can have. |
| `SBGRD-X02` | Does not accuse a single function in isolation. | With no siblings there is no asymmetry, and asymmetry is the whole evidence. Accusing a lone function would mean judging the domain — which this gate explicitly cannot do. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/dialect.go` | `KnownDialectFamilies` | core — recognising an exported function and a guard belongs to the project's dialect |
| DEP2 | `internal/mapx/model.go` | `KindCode` | core — the gate only has jurisdiction over code |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
