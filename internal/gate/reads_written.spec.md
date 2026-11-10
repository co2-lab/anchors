<!-- @anchors
  code: RDWRS
  updated_at: 2026-11-10
  layer: gate
-->
# ReadsWritten — every data a header reads is written by some file's header

> **Code**: `RDWRS`

## Overview

A file that persists a record and a file that consumes it are tied by nothing an import says: the
writer and the reader of an audit item may never import each other. Measured in MIF: seven
defects where the writer and the reader of one persisted record disagreed on a field, a sign, a
vocabulary or a rename, each side green against its own fixture.

The headers say it: `writes: AuditItem` on the writer, `reads: AuditItem` on the reader. The map
ties the reader to every writer of the data (`reads-from`, GRBLG-B33), and the reader's evidence
holds the writer, so a change to the writer's shape stales the reader's proof (EVFRA-B16).

This gate closes the silent end of that tie: a data a header reads that no header writes makes no
edge, and the reader would look tied while tied to nothing — a typo, or a writer that never
declared itself. A file whose header reads nothing is not its subject.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the reads | the header's `reads:`, a comma-separated list of names | — | the scan (RPSCR-B60) |
| the ties | the map's `reads-from` edges out of the file | — | the map (GRBLG-B33) |

## Effects

| Effect | Description |
| --- | --- |
| `RDWRS-B01` | A data the header reads that ties the file to no writer fails, naming it; every read tied to a writer passes. (`checkReadsWritten`) |
| `RDWRS-B02` | A file whose header reads nothing is skipped. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RDWRS-E01` | No map is built. | Pending, saying there is no map. | The ties are the map's edges; without them, Pass would claim a tie nobody checked. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
