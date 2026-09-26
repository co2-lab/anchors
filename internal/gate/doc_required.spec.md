<!-- @anchors
  code: DCRQD
  updated_at: 2026-09-26
  layer: gate
-->
# DocRequired — the aggregated document the unit must feed

> **Code**: `DCRQD`

## Overview

Confronts a unit against the documents the project declared as a CONTRACT: **the unit was
delivered, but did it feed what it was supposed to feed?**

The Structure declares which documents are mandatory and WHEN each one must be touched —
a layer triggers a duty. That resolution already existed, and two places consulted it: the
one that informs whoever picks up the card, and the one that lists the duties. **Neither
confronted.**

What that cost, measured: a delivery shipped a new unit — with spec, code, feature and
test — without touching either the interface contract or the data schema. Both were
declared mandatory for that layer, both were untouched, and every check went green. It
only surfaced because two agents collided on the same card and there was something to
compare: one branch carried 267 lines the other did not. **Without the collision nobody
would have noticed** — there is nothing to compare when only one piece of work exists.

**Two questions, and the second is the one that matters:** does the document EXIST, and
does it MENTION this unit? The first alone would approve an empty file created to silence
the gate.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the declared duties | what the Structure declares, or nothing | — (no declaration is a case, not an error) | the Structure: what is not declared is not charged |
| the unit's layer | the layer of the UNIT the spec describes | the layer of the node itself, which for a spec is always `spec` | this unit: reading the node's layer would charge every spec the same duty |
| the map | a built graph, or none | — | this unit: with no map the aggregated verdict is skipped, never approved |

## Effects

| Effect | Description |
| --- | --- |
| `DCRQD-B01` | A mandatory document that DOES NOT EXIST fails. |
| `DCRQD-B02` | A document that exists and does not MENTION the unit fails too — existence alone would approve an empty file created to silence the gate. |
| `DCRQD-B03` | A mention by the unit's identity CODE counts as documented. |
| `DCRQD-B04` | A mention by the FILE NAME also counts: the document speaks of the unit either way. |
| `DCRQD-B05` | Satisfying one of two declared duties is not enough — each document is charged on its own. |
| `DCRQD-B06` | Without a declaration in the Structure nothing is charged: the ruler is what the project committed to, not what one supposes it owes. |
| `DCRQD-B07` | A layer with no trigger declared is not charged, even when the document exists. |
| `DCRQD-B08` | Aggregated, the verdict is ONE PER DOCUMENT, not one per unit. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCRQD-I01` | The duty starts from the SPEC, not from the code. The spec is what declares the unit; starting from the code would charge the duty of a file that merely realises it. | declares the duty and verifies the charge lands on the spec |
| `DCRQD-I02` | The layer used is the UNIT's, not the node's. A spec node always has layer `spec`; reading it would charge every spec of the project the same duty, or none. | confronts a spec whose unit belongs to a triggering layer and verifies the charge |
| `DCRQD-I03` | Without a map the aggregated verdict is skipped, never approved. Approving without being able to look would stamp what was not measured. | runs the aggregate with no graph and verifies it does not approve |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCRQD-X01` | Does not understand the document's CONTENT. | The search is coarse on purpose: it looks for the identity code and the file name. A contract that cites the route and describes the wrong shape passes here — and that is acceptable. The gate separates "not documented" from "documented"; judging the quality of the documentation is another ruler. |
| `DCRQD-X02` | Does not decide WHICH documents are mandatory. | That is a decision of the project, declared in the Structure. A gate that invented duties would charge what nobody committed to. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `RequiredFor` | core — the duties and their triggers are declared in the Structure |
| DEP2 | `internal/mapx/model.go` | `KindSpec` | core — the duty starts from the spec, and the layer comes from the unit |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
