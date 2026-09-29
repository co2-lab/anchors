<!-- @anchors
  code: EDSTD
  updated_at: 2026-09-29
  layer: mapa
-->
# EdgeStamping — recording on each relation that it was confronted, with what result, and since when

> **Code**: `EDSTD`

## Overview

Gates run per node; validation is recorded per relation. This unit is the glue between the two: given
the verdicts of the nodes a `check` confronted, it stamps the relations whose BOTH ends were confronted —
only then has the relation really been re-validated. The stamp records both ends' revisions and the
verdict, and that is what lets the relation go stale by itself the next time either end moves.

The stamp records CHANGE, not verification. Its date answers "since when is this relation as it is", so
re-confronting and finding the same revisions and the same verdict keeps the date; a new revision or a
new verdict moves it. A stamp with no date takes today's, or the gap would be preserved forever.

A waiver is a person's decision, and the mechanical check does not confirm it: both ends passing their
gates says nothing about why a rule was waived. The check leaves a waived stamp as it was; when an end
moves, the relation goes stale and a person decides again.

A judgment gate is answered by a person or an AI about ONE target, which no two-ended round can cover. So
a single relation, or every relation touching a node, can be stamped directly, and when the gate is named
its verdict is also recorded as a judgment of its own — a field the check never rewrites, because the
check rewrites the stamp on every round and would otherwise erase what was already answered. A judgment
holds while both ends keep the revisions it was given at; when either moves, it is a question again.

Judging keeps a waiver too, with one difference from the check, and it is chosen: a waiver answers ONE
gate's question. Another gate's verdict answers a different question, so it is recorded as that gate's
judgment and the waived stamp stays; the gate that waived, judging again, is a person answering the same
question again and replaces its own waiver; and a new waiver always lands. Keeping every waiver against
every judgment would leave a waiver its own gate could never undo; replacing it on any judgment would let
an unrelated gate erase a person's decision, which is what the check already refuses to do.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node verdicts of a round | any list of confronted nodes, each passed or failed | — | the caller: only nodes it confronted are listed |
| the date | the text the caller gives as today | — | the caller: the package never reads the clock |
| a relation or a node to judge | any pair of ids or any node id | a pair with no relation, a node with no relation | this unit: a missing relation answers false, a node with none stamps zero |
| the gate name | a judgment gate's name, or nothing | — | this unit: with no name, no judgment is recorded |

## Effects

| Effect | Description |
| --- | --- |
| `EDSTD-B01` | A round stamps only the relations whose two ends were both confronted in it, and reports how many it stamped. |
| `EDSTD-B02` | The round's verdict on a relation is issue when either end failed, and ok when both passed. |
| `EDSTD-B03` | A round never rewrites a waived stamp: its verdict, date and revisions stay as they were. |
| `EDSTD-B04` | A stamp records both ends' current revisions, so the relation is fresh right after it and stale again when either end moves. |
| `EDSTD-B05` | The stamp's date stays where it was when the revisions and the verdict are unchanged, and moves to the given date when either changes. |
| `EDSTD-B06` | A stamp with no date takes the given date even when nothing else changed. |
| `EDSTD-B07` | One relation can be stamped by its two ends; stamping a relation that does not exist answers false and stamps nothing. |
| `EDSTD-B08` | Stamping one relation with a gate named records the gate on the stamp and a judgment of that gate on the relation. |
| `EDSTD-B09` | Stamping a node stamps every relation that touches it, coming in or going out (a waiver aside, B14), reports how many, and with a gate named records that gate's judgment on each. |
| `EDSTD-B10` | A node counts as judged by a gate only while some relation touching it holds that gate's judgment at both ends' current revisions; another gate's judgment never answers for it. |
| `EDSTD-B11` | A judgment survives the next round of the check, which rewrites the stamp. |
| `EDSTD-B12` | A relation holds one judgment per gate: judging again replaces it, and keeps its date when the verdict and revisions are the same. |
| `EDSTD-B13` | The stale relations are listed: the never stamped and those with an end moved since the stamp. |
| `EDSTD-B15` | A snapshot of the stamps, taken before a round stamps its copy of the map, lets the round carry to the map as it is on disk only the stamps it changed; a stamp another process changed since the snapshot is kept as theirs and counted, and an edge the map on disk does not have is left out (`EdgeStamps`, `ApplyStampChanges`). |
| `EDSTD-B16` | Moving a file from one revision to another carries along everything measured at the first — its signal, proofs, coverage and mutation, the closures of the tests that reach it, and its edges' stamps and judgments —, and leaves what was measured at any other revision, and every other file, as they were. (`RebaseRev`) |
| `EDSTD-B17` | Keeping a file's evidence moves, from every earlier revision it holds to the current one, the scenario proofs and execution, the closures other tests recorded, and its edges' stamps and judgments; the node's revision becomes the current one. Coverage and mutation, which name lines, move only when lines are kept, and the node's shared revision stays while it stamps them. The declaration — the revisions carried, the current one, the reason, the day, and whether lines went along — is recorded on the node, the latest five kept; a file with nothing at an earlier revision carries nothing and records nothing; an unknown file is left alone. (`KeepEvidence`) |
| `EDSTD-B14` | Judging a node or a single relation keeps a waived stamp recorded by another gate, or when no gate is named, and still records the judging gate's judgment; the gate that recorded the waiver replaces it, and a new waiver always replaces the stamp. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `EDSTD-I01` | Two rounds over the same graph, with the same verdicts on the same date, produce the same stamps. | stamps two identical graphs on the same date and compares |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `EDSTD-X01` | The package never decides the date: every stamp carries the date the caller passed. | The map is versioned; a clock read here would make the file change on its own. |

## Errors

none — a relation or node that does not exist is normal input answered by B07 (false) and B09 (zero stamped); stamping works on the graph in memory and cannot fail.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Graph`, `Stamp`, `Judgment`, `Stale` | mapa — the relations stamped and the staleness rule the stamp feeds |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
