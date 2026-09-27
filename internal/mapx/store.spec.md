<!-- @anchors
  code: GRPRG
  updated_at: 2026-09-26
  layer: mapa
-->
# GraphPersistence — saving and loading the map file without churn and without partial reads

> **Code**: `GRPRG`

## Overview

The map is a versioned file at the root of the project, written by the binary and committed by the team.
This unit writes the graph to that file and reads it back.

Two decisions shape it. First, the format belongs to whoever WRITES: a save always stamps the current
format and the release of the binary that wrote it, whatever the in-memory graph carried, so a freshly
assembled graph can never be born asking for a migration. Second, a save that would change nothing but
the name of the binary that wrote it is not a change: a local build writes "dev" and continuous
integration writes the published release, and rewriting the line each time would make the map flip
between the two on every run — the very oscillation the line exists to reveal.

On load, the format is checked BEFORE the graph is handed to anyone. A file this binary cannot read is
refused whole, because interpreting part of it and saving later would silently drop what was not
understood.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the graph to save | any graph, whatever format number it carries | — | this unit: the format number is overwritten on save |
| the path | a writable file path | a path whose directory does not exist | this unit: the write error is returned |
| the file to load | a map in a readable format | a missing file, text that is not the map's structure, a format outside the readable range | this unit: each is returned as an error and no graph is given back |

## Effects

| Effect | Description |
| --- | --- |
| `GRPRG-B01` | Saving (`Save`) always writes the current format and the release of the running binary, even when the graph carried another format number. |
| `GRPRG-B02` | The saved file starts with a fixed comment header that names the file and how it is generated. |
| `GRPRG-B03` | When the only difference between the new content and the file on disk is the line naming the binary that wrote it, the file is left untouched. |
| `GRPRG-B04` | When anything else changed, the file is rewritten, carrying the release of the binary that is running. |
| `GRPRG-B05` | Loading a map whose format this binary cannot read returns the format refusal and no graph. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GRPRG-I01` | A graph saved and loaded back has the same nodes, edges, stamps, judgments and signals it was saved with. | closed cycle: save a populated graph, load it, compare |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GRPRG-X01` | Loading never hands back a partially read graph: an unreadable format yields no graph at all. | A partial graph saved later would drop what was not understood, with nothing reported. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `GRPRG-E01` | The file to load does not exist or cannot be read. | The read error is returned, with no graph. | There is no map to reason about; a caller that wants an empty graph decides that itself. |
| `GRPRG-E02` | The file's text is not the map's structure. | The parse error is returned, with no graph. | A half-parsed map would be the partial read the unit exists to prevent. |
| `GRPRG-E03` | REF[GRPRG-B05]: a format outside the readable range is the failure B05 answers with the format refusal | — | — |
| `GRPRG-E04` | The file cannot be written (its directory does not exist). | The write error is returned. | Swallowing it would let the caller believe the map was saved. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/format.go` | `ConfereFormato` | mapa — the readable-format check run on load |
| DEP2 | `internal/mapx/model.go` | `Graph` | mapa — the structure written and read |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
