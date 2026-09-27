<!-- @anchors
  code: MPMRM
  updated_at: 2026-09-26
  layer: comando
-->
# MapMerge — the git merge driver that unites two versions of the map instead of merging text

> **Code**: `MPMRM`

## Overview

The map file is derived, and git does not know it: it merges the file line by line. When both sides touch
the same lines it asks for a manual resolution nobody can do reliably on a graph of hundreds of edges;
when they do not, it resolves on its own, and that is where the damage happens. Measured in a real
project, a conflict-free merge of the map erased sixty-two judgment stamps, and a later one dropped every
node that only the incoming branch had created while git reported that the merge went well. The warning
of `map build` cannot see either loss, because it happens inside git, before Anchors runs.

This unit is the driver git calls instead, once the repository declares it for the map file. It receives
the three versions git passes — the common base, our side and the other side — and writes the result
onto our side, which is where git expects it. The result is our map with what only the other side had
added to it: the nodes it alone has, the edges it alone has (with the judgments they carry), the
judgments of edges both sides share — joined gate by gate, so two branches that judged the same edge with
different gates keep both —, the failures observed on shared nodes, and the flow graph. A judgment is state
of the work — the report behind it lives in the command that recorded it, not in the file — so losing it
sends the judgment back to the queue and the evidence has to be redone.

It is not meant to be run by hand, and it reports on the error stream how many judgments the result
keeps and how many came from the other side.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the arguments | exactly three paths: base, our side, the other side | any other number of arguments | this unit: refuses any other count |
| our side and the other side | readable map files | a missing or unparseable file | this unit: fails naming the side it could not read |
| the base | any path; it is never read | — | this unit: the base is ignored |

## Effects

| Effect | Description |
| --- | --- |
| `MPMRM-B01` | The result is written onto our side's file, which is what git expects of a merge driver. |
| `MPMRM-B02` | A node that exists only on the other side is added to the result, and a node only on our side stays. |
| `MPMRM-B03` | A node both sides have keeps our side's version: its revision is not reconciled here, because only a rebuild over the merged tree knows the file's content. What was ingested on it is joined instead (B08). |
| `MPMRM-B04` | An edge that exists only on the other side is added to the result with the judgments it carries; an edge is identified by its type and its two ends. |
| `MPMRM-B05` | When only one side has judged an edge both sides share, that side's judgments are in the result, whichever side it is: a side with no judgment on the edge never erases the other's. |
| `MPMRM-B07` | On an edge both sides share, judgments are joined by gate: a gate only one side judged is kept, and a gate both sides judged keeps the judgment with the latest change date, ours on a tie; the edge's check stamp follows the same rule. |
| `MPMRM-B08` | A node both sides have gets the other side's observed failures joined by rule (the latest last occurrence wins, ours on a tie), and the other side's test signal when both carry the same revision; the flow graph is the union of both sides' states, by code, and transitions. |
| `MPMRM-B06` | After writing, the error stream reports how many judgments the result keeps and, when there are more than our side had, how many came from the other side. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MPMRM-I01` | Every node and every edge of either side is in the result: the merge only adds to our side, never removes from it. | merges two sides with nodes and edges exclusive to each and finds all of them in the result |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MPMRM-X01` | Does not read the base version. | The union does not need to know what was there before: a judgment is not deleted by its absence on one side. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MPMRM-E01` | Our side or the other side cannot be read or parsed as a map. | The merge fails naming which side it could not read, and our side's file is left as it was. | A merge written from half the information would silently drop the other half, which is the loss the driver exists to prevent. |
| `MPMRM-E02` | The driver receives other than three paths. | The command is refused before reading anything. | Git always passes three; any other call is a misconfiguration, not a merge. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/store.go` | `Load`, `Save` | mapa — reading and writing the map files |
| DEP2 | `internal/mapx/model.go` | `Graph`, `Edge`, `Node`, `FlowGraph` | mapa — the judgments, failures and flow the driver joins |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
