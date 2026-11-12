<!-- @anchors
  code: SYMXS
  updated_at: 2026-11-12
  layer: gate
-->
# SymbolExists — a name marked as a symbol is in the project's code

> **Code**: `SYMXS`

## Overview

A comment or a spec that names a symbol makes a claim about the code: this option, this
function, this field exists. Nothing confronted the claim. Measured in MIF (2026-11-12): a
comment of the import handler explained a rule by an option, `isEvalSupported`, that the PDF
library had removed and the reader no longer passed; the rule's code was right at all four
ends, every gate passed, and only a review caught the stale name.

Prose cannot be told from a name by its shape, so the project marks the names it cites:
`@@name`, `#name` or `$name`, with dotted parts (`#Reader.open`). The gate reads the marks of a
code file's comments and of a spec's text, and confronts each against the identifiers of the
map's code and test files, read with their comments removed. Only the project is searched: a
name of a dependency the project does not use is cited on purpose, and says so with
`@no-symbol-check: <reason>` on its line.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | a code file whose comments the project declares, and a spec | any other kind | this unit: it skips the others |
| a mark | `@@`, `#` or `$` before a name, the character before the marker not a word character, `@`, `#`, `$`, `/`, `.` or `&`, nor the `](` of a link, and the name not followed by `-` | a `#` that opens a comment, `#region` and `#endregion`, a `#` color, a fenced block of a spec | this unit |
| the project's names | the identifiers of the map's code and test files outside their comments, `$` a character of a name | names held only by comments, specs or features | this unit, read once per map |

## Effects

| Effect | Description |
| --- | --- |
| `SYMXS-B01` | A marked name — every dotted part of it — that is an identifier of the project's code passes; one that is not fails, naming the file, the line and the name. (`checkSymbolExists`, `projectSymbols`) |
| `SYMXS-B02` | A code file is read in its comments only, a spec outside its fenced blocks; a mark glued to a word, a path or a link, a name followed by a dash, an editor's folding directive (`#region`, `#endregion`), a color (`#` and 3, 4, 6 or 8 hex digits) and the `#` that opens a comment are not marks. (`citedSymbols`, `markedAt`, `commentText`, `withoutFences`) |
| `SYMXS-B03` | A line carrying `@no-symbol-check: <reason>` is not confronted. |
| `SYMXS-B04` | An artifact with no mark, or of another kind, is skipped. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `SYMXS-E01` | The project declares no comments for the code file. | Pending, naming the file and `dialect.comments`. | The core reads no language's comments; with none declared the marks cannot be found, and Pass would claim a confrontation nobody made. |
| `SYMXS-E02` | No map is built. | Pending, saying there is no map. | The project's names are read from the map's files; without them every mark would fail for a reason that is not the mark's. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
