<!-- @anchors
  code: FRSCF
  updated_at: 2026-10-20
  layer: mapx
-->
# Frozen scope — the map carries each freeze to everything it covers

> **Code**: `FRSCF`

## Overview

A freeze keeps a part of the project as it is until someone thaws it (DESIGN-frozen-units.md). The scan
reads its marks where they are written — `frozen:` in a header, `@frozen:` on a rule's or a heading's
line — and its record, an entry of a doctrine document coded with the freeze letter. This unit carries
each mark to everything it covers, because the map is the one that knows who is whose: a unit's spec
frozen covers the unit's code, feature and tests, and a frozen doctrine document or rule covers the spec
rules that realize it. The gates, the reports and the guides ask the node, and never read a mark.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the files | the scanned files, with their freeze marks and entries | — | the scan |
| a unit | the code its owner's header declares, and the files that `ref:` it or stand beside its anchor | a file that only cites the code in its text | this unit: such a file is not covered |
| a value | a freeze's code, or a reason | — | this unit: a code with no entry reads as the code alone |

## Effects

| Effect | Description |
| --- | --- |
| `FRSCF-B01` | A file's own `frozen:` freezes it; a unit's owner frozen — its spec, or the file that owns its code — freezes every file of the unit, the ones that `ref:` it and the ones derived beside its anchor; a file that only cites the unit's code in its text is not covered, and a file's own mark wins over the one it inherits. (`Frozen`) |
| `FRSCF-B02` | A rule's mark is carried on its node; a frozen doctrine document, or a frozen doctrine rule, freezes each spec rule that realizes it, with its value, unless that rule carries a mark of its own. (`FrozenRules`) |
| `FRSCF-B03` | A freeze's code reads with its recorded reason (`QLDTQ-F01 — the reason`), a reason reads as itself, and a code no doctrine records reads as the code. (`DescribeFreeze`, `FreezeRecords`, `FreezeRecord`) |
| `FRSCF-B04` | The frozen parts are every frozen file and every frozen rule of a file not frozen whole, by path and rule; with a value, only that freeze's; the project's frozen rules are every node's, and a path's freeze is its node's. (`FrozenParts`, `FrozenPart`, `FrozenAt`) |
| `FRSCF-B05` | A file a glob of anchors.yaml's `frozen_paths` matches is frozen by that entry's freeze, as its own header's mark would freeze it: a unit's owner matched freezes the unit's files wherever they are, and a file's own mark wins over the folder's. (`applyFreezes`, `Config.FrozenPathOf`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FRSCF-I01` | A map with no freeze mark has no frozen node and no frozen rule. | the build test's files, built, carry no freeze |

## Errors

none — the build reads files already scanned; an edge whose end is not a node of the map, or a file the build did not take, is normal flow, and carries no freeze.
