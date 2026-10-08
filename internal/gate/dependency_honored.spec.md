<!-- @anchors
  code: DEPHN
  updated_at: 2026-10-08
  layer: gate
-->
# DependencyHonored — a spec declares no dependency: the code does

> **Code**: `DEPHN`

> **DEPHN-R0001:** the gate stopped confronting a spec's Dependencies table with the code.
> A spec precedes the code, and the files a unit imports and the methods it calls are born
> with it: they are declared where the import is, by the `@dep:` flag, and confronted by
> `dep-declared` and `dep-honored` (DESIGN-dependencies-out-of-the-spec.md). What this gate
> asks now is that the table is gone.
> **Revises:** `DEPHN-B01`, `DEPHN-B02`, `DEPHN-B03`, `DEPHN-B04`, `DEPHN-B05`, `DEPHN-B06`, `DEPHN-B07`, `DEPHN-B08`, `DEPHN-I01`, `DEPHN-I02`, `DEPHN-I03`, `DEPHN-X01`, `DEPHN-X02`, `DEPHN-E01`.

## Overview

A spec that still carries a Dependencies table — rows opening with a `DEPn` — declares, in the
artifact that comes first, what only the code can know. The gate names it as a divergence to
migrate: `anchors migrate` removes the table and rewrites the `DEPn` a rule cites. A spec with
no table passes.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — | the gate engine, routing by `on: [spec]` |
| the spec's text | rows opening with a `DEPn` | a `DEPn` cited in another column (a rule's uses) | this unit, reading only the rows that open with it |

## Effects

| Effect | Description |
| --- | --- |
| `DEPHN-B01` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` |
| `DEPHN-B02` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` |
| `DEPHN-B03` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` |
| `DEPHN-B04` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` |
| `DEPHN-B05` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` |
| `DEPHN-B06` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` |
| `DEPHN-B07` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` |
| `DEPHN-B08` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` |
| `DEPHN-B09` | A spec whose text still has rows opening with a `DEPn` — a Dependencies table — diverges, naming how many rows and pointing at `anchors migrate`; a spec with none passes, and an artifact that is not a spec leaves without a verdict. (`checkDependencyHonored`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DEPHN-I01` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` | — |
| `DEPHN-I02` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` | — |
| `DEPHN-I03` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` | — |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DEPHN-X01` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` | — |
| `DEPHN-X02` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` | — |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DEPHN-E01` | @retired: DEPHN-R0001 — the spec declares no dependency any more; the code does, by `@dep:` | — | — |

## Open Decisions

none
