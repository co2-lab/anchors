<!-- @anchors
  code: GRINC
  updated_at: 2026-10-08
  layer: mapa
-->
# IncrementalMap — new files enter the map without the tree being read

> **Code**: `GRINC`

## Overview

A file created after the last `map build` had no node, and whatever reached the map for it was
dropped: `anchors test` proved a new spec's rules and had nowhere to write the proof, so
`scenario-coverage` blocked a spec whose tests passed. Rebuilding the whole map fixes it and reads
every file for the sake of one.

The map already knows most of what a new file relates to: every path, kind, layer and declared
code. So the new files are read, every other file stands in as what its node keeps, and the same
`Build` runs over that set in memory; only what involves the new files is merged. The files whose
content decides the new files' links — the whole unit the derivation ties them to, and the anchors
beside them — are read too, in a second pass. Three ways in use it: `anchors new`, the watcher, and
the ingestion that follows a run.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the paths | any relative paths | — | this unit: known ones are left alone; the reader drops ignored and unclassified ones |
| the reader | a function that reads files as the scan does | — | the caller (`scan.ScanPaths` on disk) |
| the map on disk | a map, or none yet | — | this unit: with no map, nothing is written — the first `map build` makes it whole |

## Effects

| Effect | Description |
| --- | --- |
| `GRINC-B01` | New files enter the map with the same nodes and edges the full build gives them, reading only them and their unit — never a file outside it. (`AddFiles`) |
| `GRINC-B02` | The derivation links inside the re-read unit are replaced by what the derivation says now, removals included: a new feature takes the place of the spec→test link of a unit that had none. |
| `GRINC-B03` | A new anchor gives its declared code to its derived siblings that declare none. |
| `GRINC-B04` | What the new file declares — its seeds, needs, dependencies, realizes — reaches the existing files it names. |
| `GRINC-B05` | A path the map already has is not read; a path the reader does not give (ignored, in no layer, missing) adds nothing. |
| `GRINC-B07` | The anchors that may own a new file — beside it, or whose name the new file carries — are read, so an override their header layer chooses decides the new file's links as in the full build. |
| `GRINC-B06` | On disk, the addition runs under the map's lock and the map is re-read first; with no map, nothing is written; the tree's missing governed files can be added by listing its paths, without reading the others. (`AddFilesAt`, `AddMissingAt`, `Reader`) |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GRINC-X01` | A relation an EXISTING file declares toward the new one — a dependency row, a plan's seed, an `@realizes`, a scenario code cited in a test of another directory — waits for the next `map build`. | The map does not keep those declarations; finding them means reading the files that make them, which is the tree. The commit hook and the CI still demand the full map. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `GRINC-E01` | The reader fails. | The error comes back and the map is not changed. | A half-read unit merged would be a map the full build contradicts. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
