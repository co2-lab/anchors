<!-- @anchors
  code: PRSNT
  updated_at: 2026-09-28
  layer: gate
-->
# PresentationGates — the presentation validations, confronted

> **Code**: `PRSNT`

## Overview

Each row of `Presentation validations` says: this prop or state, under this condition, makes the
unit look like this. Four questions follow from that shape, all answered from the spec's own text —
no language, no rendering: is every value of the prop decided, does one condition lead to one
appearance, is the text shown a message code rather than copy, and can a test point at what
changes. All four are informational by default.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the spec | any spec | another kind of node | this unit: Skip |
| the rows | rows of four cells, rule code first | a row still TODO, or shorter | this unit: left out |
| the value sets | backticked values of the props table, first column of a table under a backticked heading | a set of fewer than two values | this unit: not an enumeration, not confronted |

## Effects

| Effect | Description |
| --- | --- |
| `PRSNT-B01` | The rows are read from the Presentation validations section by position — rule, what it reads, condition, appearance —, and a spec with none is skipped by all four gates. |
| `PRSNT-B02` | `presentation-exhaustive` fails naming, per prop or state, the declared values no row's condition names; a row whose condition says "otherwise" (in any supported language) covers every value; a prop with no declared set of two values or more is not confronted. |
| `PRSNT-B03` | `presentation-conflict` fails naming the rules that give one prop and one condition two different appearances. |
| `PRSNT-B04` | `presentation-copy-single-source` fails naming the rule whose appearance carries text in double quotes (or “ ” « ») without citing a message code; single quotes are values, not copy. |
| `PRSNT-B05` | `presentation-observable` fails naming the rules whose appearance cites, in backticks, no identifier of the spec's Test Identifiers section. |

## Errors

none — the gates read only the spec's text

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate/rule_uses.go` | `splitSections`, `declaredNames`, `usedItems` | gate — the sections and cells of the spec |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
