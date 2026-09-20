<!-- @anchors
  code: RTDCL
  updated_at: 2026-09-20
  layer: gate
-->
# RouteDeclared — a screen declares how one arrives, and names its neighbours

> **Code**: `RTDCL`

## Overview

Confronts a SCREEN spec against the navigation graph it belongs to: **is there a named
route that reaches this screen, and do its edges point at concrete screens?**

It is the navigation traceability ruler. Without the route the screen is a loose node —
something exists that nobody can reach. With generic terms in the navigation tables
("Next screen", "Main menu") the edge points nowhere: it reads like a link and connects to
no node at all, which is worse than an absent edge, because an absent edge is visible and
a vague one passes for a declared one.

**Jurisdiction is the whole design.** Only `layer: screen` is charged. Hooks, business
logic, stores and DAOs have no route, and charging them was the vice of the legacy
validator this gate replaces — it knew only screen and component, so it treated every
non-component as a screen and produced a false positive for every unit that legitimately
has no route. A gate that cries over what cannot be fixed teaches the team to ignore it.

That is also why the Skip carries a stated REASON. A bare indeterminate count leaves the
reader wondering whether the silence is their problem; saying "this is not a screen"
closes the question in one line.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the layer | `screen` is charged; every other layer is left alone | — (an absent layer is a case, not an error) | this unit, reading the header first and the node's tags second |
| the route line | a named route written in the header, in either declared language | a route the project never wrote | the spec's author |
| the navigation sections | the entry and exit tables, in either declared language | prose outside those sections | this unit: only table rows are confronted |

## Effects

| Effect | Description |
| --- | --- |
| `RTDCL-B01` | An artifact whose layer is not `screen` leaves without a verdict, and the verdict SAYS why — a bare indeterminate count would leave the reader wondering whether the silence is their problem. |
| `RTDCL-B02` | A screen with no named route FAILS: without it the screen is a node nobody can reach. |
| `RTDCL-B03` | A screen whose navigation table carries a generic term FAILS — an edge that names no concrete screen points nowhere. |
| `RTDCL-B04` | A screen with a named route and concrete neighbours passes. |
| `RTDCL-B05` | Route and navigation are recognised in either declared language, so a project writing in its own language is measured and not silently approved. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RTDCL-I01` | The header's declared layer is the source of truth of identity, and the node's tags are only the fallback when the header is silent. The header is what the author wrote on purpose; a tag can come from inference. | confronts an artifact whose header and tags disagree and verifies the header wins |
| `RTDCL-I02` | Only TABLE ROWS of the navigation sections are confronted, never surrounding prose. An explanatory sentence mentioning "the next screen" is the author writing, not an edge being declared. | confronts a screen whose prose carries a generic term outside any table and verifies it is not accused |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RTDCL-X01` | Does not charge a route from any layer other than `screen`. | Hooks, business logic, stores and DAOs have no route. This was the legacy validator's vice — it knew only screen and component, treated every non-component as a screen, and produced a false positive for every unit that legitimately has none. |
| `RTDCL-X02` | Does not verify that the declared route EXISTS in the router. | The ruler here is the spec's internal coherence: a route is named and the neighbours are concrete. Confronting the name against the real routing table is a different confrontation, over a different artifact. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Node` | core — the node's tags are the fallback for identity |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
