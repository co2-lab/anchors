<!-- @anchors
  code: SGINA
  updated_at: 2026-09-27
  layer: mapa
-->
# SignalIngestion — hanging the runner's results on the map's nodes: executions, proven rules, coverage and mutation

> **Code**: `SGINA`

## Overview

Anchors does not run a project's tests. It consumes the artefacts the project's runner already writes —
execution reports, line coverage, mutation reports — and hangs each measurement on the node it measured.
This unit is that hanging: it matches the paths in a report to the nodes of the map and records the result
together with the revision each node had at the time, so a later change of the file makes the
measurement stale instead of silently still valid.

Reports rarely use the map's paths. A runner writes absolute paths or paths relative to its own workspace,
so matching is by path suffix, and only at a path boundary. In a monorepo the same relative path can exist
in two workspaces; the report paths are then resolved before ingestion, preferring the node that shares
the most leading directories with the report file itself, and leaving a real tie with no owner rather
than asserting a proof nobody measured.

An ingestion is a MEASUREMENT, not an addition. A full run that no longer proves a rule removes it: "no
rule proven" is information, and keeping the last proof forever would let the map claim what is no longer
true. But a measurement speaks only for what it measured. When the suite is named, each suite's proof and
coverage are kept apart and the node carries their union, so ingesting the mobile suite does not erase the
backend's; a partial run speaks only for the rule codes its cases named. The union is only as fresh as its
oldest contributor, and a contribution with no recorded revision is never fresh. A report ingested from
outside the repository is an ad-hoc run that never runs again; it can be dropped when a real run lands, so
it does not keep nodes stale forever.

Coverage from several suites is a union of LINES — a line covered by any suite is covered, once. Only
suites measured against the current text take part: a suite recorded before the file was edited numbers
other lines. A report that gave only totals cannot join a union of lines, and the best single suite is used
instead — an under-estimate, never a double count.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| execution per report path | passed, failed and skipped counts per test file path, absolute or relative | — | the report reader; unmatched paths are simply ignored |
| proven and seen rule codes | the codes of passing cases, and for a partial run the codes its cases named | — | the report reader; no seen set means a full run |
| the declared codes per node | the rule codes each spec declares, read by the caller | — | the caller: this unit does not touch the disk |
| the layer and the suite | any name, or nothing | — | this unit: no layer means `unit`, no suite means the report is the whole measurement |
| coverage per report path | covered and total lines, optionally line by line, optionally marked as older than the file | — | the report reader |
| mutation per report path | killed, survived, no coverage, ignored, score, and the tool's thresholds | — | the report reader |

## Effects

### Execution

| Effect | Description |
| --- | --- |
| `SGINA-B01` | A test node matched by a report path receives the counts under the ingestion's layer — `unit` when none is named — its totals become the sum over all its layers, and the node's revision, the revisions of its closure and the ingestion date are recorded. |
| `SGINA-B02` | Ingesting a layer again replaces that layer's counts and leaves the other layers as they were. |

### Proven rules

| Effect | Description |
| --- | --- |
| `SGINA-B03` | A node's proven rules are the ones it declares that a passing case named; a full run with no suite writes that result even when empty, so a proof that stopped existing is erased. |
| `SGINA-B04` | With a suite named, each suite's proof is kept apart and the node's proven rules are their union: another suite's ingestion never erases it, and ingesting the same suite again replaces its own entry. |
| `SGINA-B05` | A suite that no longer proves anything on a node leaves the union and leaves no empty entry behind. |
| `SGINA-B06` | The union is only as fresh as its oldest contributor: the node reads stale while any contributing suite measured an earlier revision, and an entry with no recorded revision is never fresh. |
| `SGINA-B07` | A partial run changes only what it saw: a node none of whose declared rules the run named keeps its proof, a rule it named takes this run's result, and a rule it did not name keeps its earlier proof. |
| `SGINA-B08` | Dropping the external suites removes the proof and coverage recorded by reports from outside the repository, recomputes the node's union, freshness and coverage from what remains, and names each suite dropped, sorted. |

### Coverage

| Effect | Description |
| --- | --- |
| `SGINA-B09` | With no suite named, a code node matched by a report path receives its covered and total lines, its percentage and its revision, and the percentage it replaces is kept as the baseline for the delta. |
| `SGINA-B10` | With a suite named, the node's coverage is the union of the lines covered by any suite, a line covered by two suites counting once, and ingesting a suite again replaces its own lines. |
| `SGINA-B11` | When a suite reported only totals, the union falls back to the suite with the most covered lines. |
| `SGINA-B12` | Only suites measured against the node's current revision join the union; a suite measured against an earlier one stays stored and keeps the node stale. |
| `SGINA-B13` | A report older than the file's current text is kept with no revision: it neither joins the union nor passes for fresh. |
| `SGINA-B14` | The lines of a suite are stored as compact ranges and read back as the same lines. |
| `SGINA-B19` | A report that instruments no line of the file sets the node's percentage to zero with its zero total, with or without a suite: the previous ingestion's percentage is never left beside a total of zero. |

### Mutation

| Effect | Description |
| --- | --- |
| `SGINA-B15` | A code node matched by a mutation report receives the killed, survived, uncovered and ignored mutants, the score and the tool's thresholds as the latest totals; with a scope named, the same numbers are also kept under that scope with the revision it was measured at. |

### Matching and freshness

| Effect | Description |
| --- | --- |
| `SGINA-B16` | A node's signal is stale when it recorded a revision and the node has moved since; a signal with no recorded revision is not stale. |
| `SGINA-B17` | A report path matches a node when the two are equal or one ends with the other at a path boundary, in either direction. |
| `SGINA-B18` | Resolving report paths gives each path its only matching node, keeps a path no node matches as it came, breaks a tie by the node sharing the most leading directories with the report file, and returns a path still tied as ambiguous, with no owner. |
| `SGINA-B20` | When several paths of one report match a node, the node receives exactly one of them, always the same: the exact path, else the one closest in length to the node's, else the first in lexical order. |
| `SGINA-B21` | A mutation ingestion records, beside the killed, how many of them were killed by the time limit. |
| `SGINA-B22` | An execution ingestion records each test file's run time — the sum of its cases' times — under the suite, or under the layer when the ingestion names no suite; a file whose cases carry no time records none. |
| `SGINA-B24` | The revs the tree has now replace the map's, by path; a node the tree does not give keeps its rev, and the count of changed nodes is returned. (`RefreshRevs`) |
| `SGINA-B23` | A run time Anchors measured itself is recorded on the node under the suite; a node the map does not have is ignored. (`RecordRunSeconds`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SGINA-I01` | A node's proven rules are always exactly the union of its per-suite proofs, without repetition and in sorted order. | ingests two suites proving overlapping rules and compares the node's proven rules with their sorted union |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SGINA-X01` | Execution counts land only on test nodes, and coverage and mutation only on code nodes, even when another node's path matches the report. | A count hung on the wrong kind of node would be read by a gate that asks another question. |

## Errors

none — a report path that matches no node, a node that declares nothing and a suite that proves nothing are normal input answered by B03–B05 and B18; ingestion works on the graph in memory and cannot fail.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Graph`, `TestSignal`, `SuiteCoverage`, `MutationScope` | mapa — the signals recorded on nodes |
| DEP2 | `internal/mapx/evidence.go` | `EvidenceClosure` | mapa — the closure revisions recorded with an execution |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
