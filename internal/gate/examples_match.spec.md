<!-- @anchors
  code: EXMCH
  updated_at: 2026-10-08
  layer: gate
-->
# ExamplesMatch — every Examples row is a case its tests run

> **Code**: `EXMCH`

## Overview

A scenario outline says "for each of these rows"; the parameterised test that proves it runs its
own table, and nothing kept the two together. In the reference app they drifted apart in eight
features: the Examples said `mesada` and the test `allowance`, "Crédito" against "Cartão de
crédito", and a row the code no longer had at all.

The gate reads no test library. Every value of a row must appear in the body of a test that cites
the scenario's code, as a whole token and with its case. The body is read as `test-has-assertion`
reads it. Examples often show what the screen says while the test uses a key: a column whose
header carries `(label)`, in any supported language, is display text and is not looked for.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a feature | any other kind | this unit: Skip |
| the outlines | coded scenarios with an examples table, in any Gherkin language | none | this unit: Skip |
| the tests source | the dialect's `tests`, pattern or script | none declared | this unit: Skip, nothing to measure |
| the gate entry | `labels: true` or nothing | — | the configuration load |

## Effects

| Effect | Description |
| --- | --- |
| `EXMCH-B01` | A row whose values do not all appear in the tests citing its scenario's code fails, naming the code, the row's line and the missing values; a feature whose rows all appear passes. The tests of every citing test join: a row may be in any of them. |
| `EXMCH-B02` | A value appears only as a whole token with its case: not glued to a letter, a digit or an underscore on either side. Surrounding quotes or backticks in the cell are not part of the value, and an empty cell asks for nothing. (`tokenIn`) |
| `EXMCH-B03` | A column whose header carries `(label)` in any supported language (`(rótulo)`, `(etiqueta)`) is display text and is not looked for. (`isLabelColumn`) |
| `EXMCH-B04` | The rows are those under each examples table's header, of the scenario whose coded tag precedes its title, in any Gherkin language; a table ends at the first line that is neither a row, blank nor a comment, a cell past the header's columns is looked for, and a scenario with no code has no rows. (`exampleTables`) |
| `EXMCH-B05` | A test's body is read as `test-has-assertion` reads it: with `labels: true` on this gate an empty test stands for the block around it, and a script's `end` bounds it. |
| `EXMCH-B06` | A node that is not a feature, a feature with no coded outline with a row, a project with no tests source, and a feature no test of which cites an outline's code are skipped; an outline no test cites is left to `scenario-coverage`. |

## Errors

| Error | When | What the user sees |
| --- | --- | --- |
| `EXMCH-E01` | The tests cannot be listed (a script that fails, output outside the contract) | the gate fails with the reason |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
