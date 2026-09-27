<!-- @anchors
  code: DLCND
  updated_at: 2026-09-26
  layer: comando
-->
# DeliveryConfront — confront what a delivery declares against the disk, at the moment it is declared

> **Code**: `DLCND`

## Overview

When a delivery is recorded, the author asserts things that can be checked on the spot. Two
divergences were measured in the first round where the deliver-then-review flow worked: a record
claimed to update the counts of a spec and the file had not been touched; another declared the
functions paginated, citing the pagination gate by name, while that gate was red on that very file.
Neither is solved by asking for more care — the first is memory, the second is the cost of rereading an
informative gate's output — and the record is the right place to catch them, because it is where the
author asserts.

So after a delivery is recorded, the declared files are compared with what git reports as changed; the
unit is checked for a test without a mutation signal; and the project's gates are run on the delivered
files, listing the failing ones, informative gates included, since those are exactly the ones that go
unnoticed. When git cannot be read, the output says the confrontation did not happen — silence would
read as "the declared files check out", a claim nobody verified.

The confrontation never blocks the record: it only prints. Blocking would push authors to declare less,
and the record's value is in declaring more.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the declared files | paths relative to the project root, or absolute | — | this unit: each is made relative to the root before comparing |
| the unit | the delivered unit's path, or empty | — | this unit: an empty unit gets no mutation check |
| the working tree | a git repository at the project root, or no repository | — | this unit: without git it says the confrontation did not happen |

## Effects

| Effect | Description |
| --- | --- |
| `DLCND-B01` | The confrontation only prints: a delivery with untouched files and failing gates is still recorded and succeeds. |
| `DLCND-B02` | Without a readable git repository, the output says the declared files could not be confronted and that nobody verified them, and no file is accused. |
| `DLCND-B03` | A declared file that git reports neither modified nor new is listed under "declared files that do NOT appear in the diff", with the warning that the record asserts something the disk does not confirm. |
| `DLCND-B04` | A declared file that git reports modified is not listed, and neither is a file inside a new directory that git reports as a whole. |
| `DLCND-B05` | A unit that has a test file beside it and no mutation signal in the map is warned that no mutation signal was ingested; a unit without a test, a unit already measured, or no unit, gets no warning. |
| `DLCND-B06` | The project's gates run on the map's nodes of the declared files and the unit, and each failing gate is listed as gate, target and the first sentence of its detail — informative gates included. |
| `DLCND-B07` | Without configured gates, without a map, or without a delivered file in the map, no gate is listed. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DLCND-I01` | "Could not look" and "looked and found nothing" are never the same answer. | outside git the comparison reports it did not happen; inside git, with every file touched, it reports it happened and accuses none |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DLCND-X01` | The confrontation writes nothing: not the record, not the map, not the working tree. | It reports on the declaration; changing what it confronts would make the report describe a state it created. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DLCND-E01` | REF[DLCND-B02]: git cannot be read, and B02 says the confrontation did not happen instead of accusing or clearing any file | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate` | `RunWithConfig` | gate — the project's gates |
| DEP2 | `internal/mapx` | `Load` | mapa — the delivered nodes and the mutation signal |
| DEP3 | `internal/gitmeta` | `Check`, `Explain` | why git could not be read |
| DEP4 | `cmd/anchors/flow/deliver.go` | the deliver command | comando — the only caller |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
