<!-- @anchors
  code: MGLTR
  updated_at: 2026-09-28
  layer: apoio
-->
# CodeLetters — rewriting the letter of the codes of one kind of unit

> **Code**: `MGLTR`

## Overview

A migration step may rename a code LETTER for the units of one kind: a plan's phase `F` became
`W`, a flow's step `P` became `T`, a result `R` — of an action, or of a flow fitted as a piece —
became `O` (format 5). Unlike a renamed
key, these codes live in the project's own files — plans, flows, actions, and everything that cites
them. This unit rewrites them in one file's text, given which unit codes are of which kind; the
command finds the units and the files.

The kind scopes the rename. The same letter means something else elsewhere: a spec's `-R01` is a
permission and stays, and a plan's revision block `-R0001` has four digits and is not an item.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the text | any file's text | — | this unit |
| the kinds | unit code → kind, for the units found | a code of no known kind | this unit: left as it is |
| the renames | the letter renames of the steps crossed | — | the steps (`LetterRenamesBetween`) |

## Effects

| Effect | Description |
| --- | --- |
| `MGLTR-B01` | The letter renames of the steps between two formats are gathered in order; outside that interval none are (`LetterRenamesBetween`). |
| `MGLTR-B02` | A code whose unit is of a renamed kind and whose letter is that kind's old one is rewritten to the new letter, keeping its unit and number (`RewriteCodeLetters`). |
| `MGLTR-B03` | Every other code is left as it is: a unit of another kind or of no known kind, another letter of the same unit, and a code with more than two digits. |
| `MGLTR-B04` | The rewrite reports, per code, how many times it was replaced, and the report is listed in a stable order (`SortedCounts`). |
| `MGLTR-B05` | Rewriting a text already rewritten changes nothing. |

## Errors

none — the rewrite reads only the text it is given; a code it does not recognize is a code it
leaves alone, not a failure

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/migra/step.go` | `LetterRename` | infra — the steps that declare the renames |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
