<!-- @anchors
  code: EVFRA
  updated_at: 2026-10-03
  layer: mapa
-->
# EvidenceFreshness — a test's evidence expires when anything it exercises changes, not only its own file

> **Code**: `EVFRA`

## Overview

Two kinds of freshness live in the map and must not be confused. An edge stamp records a CONFRONTATION —
text against text. A node's signal records an EXECUTION — this test ran, against the revision recorded at
ingestion. Treating a scoreboard as still valid after the code it ran against changed is the opposite of
having proof.

Judging only the test file itself is too little. A shared script (`utils/login.yaml`) is composed by
hundreds of test scripts; when it changes, every one of them should lose its evidence, and none of their
own files changed. So at ingestion the map records the CLOSURE of the test — every node reached by
descending its outgoing edges, with the revision each one had — and later compares those recorded
revisions with the current ones.

The closure descends (what the test composes and depends on), never climbs to the spec and feature it
proves. It stops at `@noPropagation` nodes, like the impact analysis: a node that declared it does not
propagate is included but not walked through, so a deliberately volatile file cannot expire half the
suite.

A test that was never ingested has no verdict here: missing proof is not stale proof. And a signal
recorded before closures existed is judged by its own file only — declaring every old signal stale would
turn a precision improvement into a flood of false expirations on the day it shipped.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the test node | any node id | an id with no node, or a node never ingested | this unit: no verdict is returned |
| the recorded closure | the revisions captured at ingestion, possibly empty | — | the ingestion records it; an empty one falls back to the node alone |

## Effects

| Effect | Description |
| --- | --- |
| `EVFRA-B01` | A node with no signal, or with a signal carrying no ingestion revision, gets no verdict. |
| `EVFRA-B02` | When the test's own revision moved since ingestion, the verdict marks the evidence stale by its own change. |
| `EVFRA-B03` | When a node of the recorded closure now has a different revision, the verdict names it as a culprit (culprits sorted), without marking the test's own change. |
| `EVFRA-B04` | When neither the test nor any node of its closure moved, there is no verdict. |
| `EVFRA-B05` | A signal recorded with no closure is judged by the test's own file only: a change in what it composes does not expire it. |
| `EVFRA-B06` | The closure of a test is every node reached by descending its outgoing edges transitively, each with its current revision. |
| `EVFRA-B07` | A `@noPropagation` node enters the closure but the walk does not descend through it. |
| `EVFRA-B08` | A node recorded in the closure that no longer exists in the graph is not a culprit. |
| `EVFRA-B09` | A `captures` target — the unit a visual-regression test captures, and its images — enters the closure but the walk does not descend through it. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `EVFRA-I01` | The closure never contains the test node itself. | computes the closure of a test inside a composition chain and checks the test is absent |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `EVFRA-X01` | The closure only descends; the spec and feature above the test are never part of it. | They are the requirement the test proves, not the inputs that make it prove anything. |

## Errors

none — a node with no signal is normal input answered by B01 (no verdict), and the closure walks an in-memory graph that cannot fail.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Graph`, `Node` | mapa — the graph walked and the node whose signal's recorded closure is compared |
| DEP2 | `internal/mapx/impact.go` | `adjacency`, `noPropSet` | mapa — the shared edge index and the `@noPropagation` set |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
