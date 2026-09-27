<!-- @anchors
  code: GRQRG
  updated_at: 2026-09-26
  layer: mapa
-->
# GraphQueries — read-only questions over the loaded map: who governs what, neighbours, orphans, counts and a parents-first order

> **Code**: `GRQRG`

## Overview

The map is walked by several commands that only need to ask it questions: which files a guide governs
directly (the size of a judgment audit for that guide), how many each guide governs, what is one step
away from a node, which nodes touch nothing at all (the structural orphan), how the graph breaks down by
kind and type, and in which order to process the nodes so that no child is fixed before the parent that
affects it.

The order is the one with a real decision in it. Only the vertical edges impose precedence — governs,
specifies, covered-by and tested-by, where the source rules or precedes the target. Informative edges
(references, reuse, seeds and the rest) do not. Nodes with no pending parent come first, and whenever
several are ready they are taken in alphabetical order, so the order is the same on every run. A cycle
must never stall the walk: the nodes caught in it are appended at the end, in alphabetical order.

Every answer here is computed from the graph in memory. Nothing is written and nothing is read from disk.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the graph | any graph, including cycles and edges repeated | — | this unit: repeated edges and cycles are tolerated |
| a node or guide id | any id | an id absent from the graph | this unit: the answer is simply empty |

## Effects

| Effect | Description |
| --- | --- |
| `GRQRG-B01` | Asking what a guide governs returns the direct targets of its `governs` edges only, sorted. |
| `GRQRG-B02` | The governance summary counts, per source, its direct `governs` edges; edges of other types are not counted. |
| `GRQRG-B03` | The neighbourhood of a node holds its incoming and its outgoing edges one step away, each list sorted by source, then target, then type. |
| `GRQRG-B04` | The orphans are the nodes touched by no edge at all, sorted by id. |
| `GRQRG-B05` | The statistics count the nodes, the edges, the nodes per kind and the edges per type. |
| `GRQRG-B06` | In the parents-first order, the source of every governs, specifies, covered-by or tested-by edge comes before its target. |
| `GRQRG-B07` | Edges of any other type impose no precedence: nodes linked only by them keep the alphabetical order. |
| `GRQRG-B08` | Nodes caught in a cycle are not dropped: they are appended after the rest, in alphabetical order. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GRQRG-I01` | The parents-first order does not depend on the order the nodes are stored in: ties are always broken alphabetically. | orders the same graph with its nodes stored in two different orders and compares |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GRQRG-X01` | The queries never change the graph. | They feed read-only views; a query that edited the map would change what the next command sees. |

## Errors

none — every query is a walk over an in-memory graph; an unknown id and a cycle are normal input answered by an empty list and by B08.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Graph`, `Edge`, `Node`, `Kind`, `EdgeType` | mapa — the graph queried |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
