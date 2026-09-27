<!-- @anchors
  code: MPCMM
  updated_at: 2026-09-26
  layer: comando
-->
# MapCommand — builds the dependency map from the project and answers questions about it

> **Code**: `MPCMM`

## Overview

The `map` command owns the project's dependency map. Its `build` walks the project reading text, never
parsing code, lets the map builder infer the edges, and writes the map file. Its `show` reads that file
and answers the questions a person or an agent asks of it: what surrounds one file, which files are
islands, how big the map is, and in which order to work through it.

A rebuild must not erase what the map remembers. The workflow runs `map build` before every `check`, so
a rebuild that started from nothing would wipe the judgment stamps of the previous step, and with them
the evidence that someone already judged the work; the flow graph, which only `flow build` fills from
files `map build` does not walk, would vanish too. So the build carries both from the previous map. And
because a stamp can still be lost for good reasons — a node was removed — or bad ones — a merge resolved
the wrong way — the build compares the stamps per gate and warns when a gate has fewer than before.

The build also says where it had to guess. When two layer patterns match the same file and neither
declares a priority, the pattern length decides, and a wrong decision takes the file out of reach of
every gate that measures the right layer. The build names each pair of layers that was decided this way,
grouped, because it is where the decision is made. Both warnings inform and never fail the build.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project for `build` | a root with a configuration file | a root with no configuration | this unit: refuses, pointing at `anchors init` |
| the map for `show` | an existing map file | no map file | this unit: refuses, pointing at `anchors map build` |
| the selector of `show` | a file of the map, the orphans switch, the statistics switch or the worklist switch | none of them, or a file that is not a node | this unit: refuses naming what to give |
| the pending switch | the worklist with a configuration to run the gates | the worklist with no configuration | this unit: refuses, naming the configuration |

## Effects

### Build

| Effect | Description |
| --- | --- |
| `MPCMM-B01` | A rebuild keeps the judgment stamps the previous map carried on the edges that still exist. |
| `MPCMM-B02` | A rebuild keeps the flow graph of the previous map, which the build itself does not produce. |
| `MPCMM-B03` | When a gate has fewer judgment stamps than in the previous map, the build warns naming the gate and the count before and after; with no loss, or with a gain, it says nothing. |
| `MPCMM-B04` | The edge summary lists the triad's edge types first, in a fixed order, then every other type alphabetically, so no type is hidden; with no edges there is no summary. |
| `MPCMM-B05` | When the layer of a file was decided by pattern length, the build warns once per pair of winning and losing layers, with the number of files and one example, and says what the wrong choice costs; with no such file it says nothing. |

### Show

| Effect | Description |
| --- | --- |
| `MPCMM-B06` | Showing a file lists the edges that reach it and the edges that leave it, marking a file with nothing above it as a top and one with nothing below as a leaf. |
| `MPCMM-B07` | The orphans switch lists the nodes with no edge at all, with their kind and count. |
| `MPCMM-B08` | The statistics switch gives the node and edge totals, the nodes by kind and the edges by type, each in a fixed order. |
| `MPCMM-B09` | The worklist lists every node in topological order, the ruler before what it governs and the spec before its code, with the node count. |
| `MPCMM-B10` | The worklist with the pending switch runs the gates and lists only the nodes with a failing gate; a judgment that is only pending does not count. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MPCMM-I01` | Recording a judgment and then rebuilding the map leaves the judgment answered: the rebuild never undoes a judgment on an edge that still exists. | judges a target, rebuilds the map, and finds the gate still answered on the target |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MPCMM-X01` | The stamp-loss and layer-ambiguity warnings never fail the build. | Removing a node legitimately removes its stamps, and the length rule is usually right; failing would reject projects that are correct. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MPCMM-E01` | `build` finds no configuration file. | It fails pointing at `anchors init`. | Without the layers nothing can be classified. |
| `MPCMM-E02` | `show` finds no map file. | It fails pointing at `anchors map build`. | There is nothing to query yet. |
| `MPCMM-E03` | `show` is given a file that is not a node of the map. | It fails naming the file. | An empty neighbourhood would read as an isolated file, which is a different fact. |
| `MPCMM-E04` | `show` is given no file and no switch. | It fails naming the accepted selectors. | There is no default question to answer. |
| `MPCMM-E05` | The worklist is asked for pending nodes with no configuration file. | It fails naming the configuration. | Pending means a failing gate, and the gates live in the configuration. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/scan/scan.go` | `Walk`, `Ambiguities` | scan — the files and the layers decided by length |
| DEP2 | `internal/mapx/build.go` | `Build`, `PreserveStamps` | mapa — the graph and what a rebuild carries over |
| DEP3 | `internal/mapx/store.go` | `Load`, `Save` | mapa — the map file |
| DEP4 | `internal/gate/gate.go` | `RunWithConfig` | gate — the failing gates of the pending worklist |
| DEP5 | `internal/config/config.go` | `Load` | config — the layers and the gates |
| DEP6 | `internal/mapx/query.go` | `Neighbors`, `Orphans`, `Statistics`, `TopoOrder` | mapa — the questions `show` answers |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
