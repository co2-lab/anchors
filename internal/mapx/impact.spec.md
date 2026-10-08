<!-- @anchors
  code: IMANM
  updated_at: 2026-10-08
  layer: mapa
-->
# ImpactAnalysis — what changing one file propagates to, and what it must be confronted against

> **Code**: `IMANM`

## Overview

Changing a file has two consequences, and they point in opposite directions. Going DOWN, from parent to
child, are the files that depend on the changed one and may need to be redone — the wave of
propagation. Going UP, from child to parent, are the files the changed one must be CONFRONTED against —
not propagation but validation: a divergence found there is reconciled either in the child or in the
parent.

The two walks differ on purpose. The downward walk follows every outgoing edge transitively, but a child
that declared `@noPropagation` stops it: the child is still reached (it is a direct dependent), and the
wave does not continue through it. The upward walk climbs transitively — the parent is itself validated
against the grandparent — but it never turns back down, so climbing to a shared guide does not drag in
the sibling units that guide also rules.

The analysis is a pure query: it opens no issue and changes nothing. Opening the issue belongs to the
gate, and executing the wave belongs to `propagate`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the graph | any graph of nodes and edges, cycles included | — | this unit: visited nodes are never walked twice |
| the starting node | any node id | an id with no node or no edges | this unit: it simply reaches nothing |

## Effects

| Effect | Description |
| --- | --- |
| `IMANM-B01` | The propagation list holds every node reached by descending the outgoing edges transitively from the changed node, sorted. |
| `IMANM-B02` | A node marked `@noPropagation` is listed when the wave reaches it, but the wave does not descend through it. |
| `IMANM-B03` | The validation list holds every node reached by climbing the incoming edges transitively — parent, grandparent and up — sorted. |
| `IMANM-B04` | Climbing never turns downward: a sibling that shares a parent with the changed node is in neither list. |
| `IMANM-B05` | A node with no edges has an empty propagation list and an empty validation list. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `IMANM-I01` | The changed node never appears in its own lists, even when a cycle leads back to it. | analyses a node inside a cycle and checks it is absent from both lists |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `IMANM-X01` | The analysis does not modify the graph. | Opening issues is the gate's job and running the wave is `propagate`'s; the analysis only shows the path. |

## Errors

none — the analysis walks an in-memory graph; an unknown or isolated starting node is normal input that reaches nothing (B05), not a failure.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
