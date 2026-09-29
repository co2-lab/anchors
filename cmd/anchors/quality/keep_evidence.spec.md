<!-- @anchors
  code: KPEVD
  updated_at: 2026-09-29
  layer: comando
-->
# KeepEvidence — a change that proves nothing new keeps the files' evidence

> **Code**: `KPEVD`

## Overview

A file's evidence is tied to the revision it was measured at, and any edit ages it. An edit may
touch only what Anchors reads — a `@realizes` on a rule line, `#region` markers, a comment citing a
code — and rerunning every suite for it is waste. In the reference app a doctrine commit added
`@realizes` to 134 specs and every one lost its proofs. Anchors cannot tell such an edit from a
behaviour change in every language; the one who made the change can, and this command is how
they say it, with the reason recorded.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the files | paths of files in the map | a file the map does not have, or that cannot be read | this unit: the command fails naming it |
| the reason | any non-blank text | none, or blank | this unit: the command fails |
| the map | the project's map, and the one at HEAD when there is one | no map | this unit: the command fails pointing at `map build` |

## Effects

| Effect | Description |
| --- | --- |
| `KPEVD-B01` | Each file's evidence moves to its current content, with the reason recorded on its node, and the command says what it carried; with `--lines` the coverage and mutation go along. (`KeepEvidence` does the moving.) |
| `KPEVD-B02` | When the map was rebuilt after the change and no longer holds the file's evidence, it is taken from the map at HEAD; with no commit or no committed map there is nothing to recover. |
| `KPEVD-B03` | A file whose evidence is already at its content, or that has none, is said to have nothing to keep, and is left as it was. |
| `KPEVD-B04` | The `@contract` stamps of the doubles pointing at the files kept are refreshed and listed, under the same declaration. |

## Errors

| Error | When | What the user sees |
| --- | --- | --- |
| `KPEVD-E01` | The reason is missing or blank, a file is not in the map or cannot be read, or there is no map | the command fails saying which, and the map is not written |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/stamp.go` | `KeepEvidence` | mapx — the evidence moved, the declaration recorded |
| DEP2 | `cmd/anchors/quality/stamp.go` | `refreshStamps` | comando — the `@contract` stamps refreshed and listed |
| DEP3 | `internal/mapx/store.go` | `LoadBytes` | mapx — the map at HEAD |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
