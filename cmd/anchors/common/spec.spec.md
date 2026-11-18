<!-- @anchors
  code: UNCDN
  updated_at: 2026-10-08
  layer: comando
-->
# UnitCodes — the identity code of a unit, read from a header, from the map, or from the codes a file names

> **Code**: `UNCDN`

## Overview

Several commands need the identity code of the unit they are working on: to date a plan's progress, to
pick the documents a unit change requires, to record which rules a file carries when the map is
ingested. This unit answers that question from three sources.

From a header: the `code:` line of an artifact's header, whatever comment style the file uses. Only a
code of the project's configured length, in capitals and digits, is read, so prose that happens to say
"code:" is not taken for an identity.

From the map: the exact node first, and when that node carries no code, any node of the same unit stem.
This is how the code file of a unit answers with the code its spec declared. Without a map there is no
answer.

From a file's text: the rule codes the file names, kept to the unit's own prefix, each once, in the order
they first appear. A code of another unit that the file merely mentions is left out.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the header text | any text; the code line may sit in a line, hash, block or markup comment | a code of another length, lower case | this unit: no code is read |
| the unit path | a path relative to the project root, either slash style | a path of a unit absent from the map | this unit: answers the empty string |
| the map | the project's map file | a missing or unreadable map | this unit: answers the empty string |
| the unit code filter | a unit code, or empty for every code | — | this unit |
| the file to scan | a readable file | a missing or unreadable file | this unit: returns the read error |

## Effects

| Effect | Description |
| --- | --- |
| `UNCDN-B01` | `CodeDoHeaderSpec`, through `SpecHeaderCodeRE`: the header code is read from a `code:` line with no comment mark or behind a line, hash, block or markup comment mark. |
| `UNCDN-B02` | A code shorter or longer than the configured code length, or in lower case, is not read as the header code. |
| `UNCDN-B03` | `CodeOfUnit`: the unit's code is the code of the map node whose identifier is exactly the unit path. |
| `UNCDN-B04` | When the exact node carries no code, or there is no exact node, the code of a node of the same unit stem answers; with neither, the answer is empty. |
| `UNCDN-B05` | With no map, the unit has no code. |
| `UNCDN-B06` | `CodesInFileOfUnit`: the codes of a file are the rule codes it names that carry the unit's prefix, each once, in the order they first appear; with no unit every rule code counts. |
| `UNCDN-B07` | A data state the file writes bare (`DS-seen-no`, as a data-state table names it) is listed as the unit's own code (`ARSCA-DS-seen-no`), after the prefixed codes; a data state prefixed with another unit's code is not the unit's. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `UNCDN-I01` | No code listed for a unit belongs to another unit, and none is listed twice. | scans a file that names its own codes twice and another unit's code once, and compares the exact list |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `UNCDN-E01` | The file whose codes are asked for cannot be read. | The read error is returned, with no codes. | An empty list would say "this file names no rule", which is a different fact from "the file was never read". |
| `UNCDN-E02` | REF[UNCDN-B05]: a missing or unreadable map is the failure B05 answers with no code | — | — |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
