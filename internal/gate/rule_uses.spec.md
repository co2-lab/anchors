<!-- @anchors
  code: RLUSG
  updated_at: 2026-10-08
  layer: gate
-->
# RuleUses — each rule says what it uses, and what it uses exists

> **Code**: `RLUSG`

## Overview

A spec has the two ends — the rules (behaviours, states, errors…) and the data (the data contract,
the domain, the props) — and nothing between them: what each rule reads lived, when it lived
anywhere, in prose. A field changed in the contract pointed at no rule, and a rule could check a
datum the contract never offered.

Three sections tie them: **Validations** (`V`: a condition on a datum → behaviour), **Presentation
validations** (`P`: a prop or state → appearance) and **Rule uses** (any other rule → what it
uses). Each row starts with the rule's code and goes on with what it uses; the row is read by
position, since the columns' names follow the project's language. Two gates read them:
`rule-uses-declared` asks every rule to say what it uses, and `rule-uses-resolve` asks that what
it uses exists in the spec.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the spec | any spec, with or without the three sections | another kind of node | this unit: Skip |
| the gate entry | `letters`, one letter A–Z each | anything else | the config: the load fails (`CNFGO-B51`) |

## Effects

| Effect | Description |
| --- | --- |
| `RLUSG-B01` | The three sections are found by their title in any supported language, or by the project's own title in `section_titles`; each row whose first cell carries a rule code gives that rule, and its second cell what it uses; a row whose uses are still a TODO says nothing yet. |
| `RLUSG-B02` | A uses cell lists its backticked items when it has any, and its comma-separated items otherwise. |
| `RLUSG-B03` | `rule-uses-declared` fails naming every rule of the letters asked about that has no row in the three sections; the gate entry's `letters` chooses the letters, and without it every letter but an open question's, a plan phase's and a flag scenario's is asked about; a rule whose line carries `@no-uses: <why>` is waived. |
| `RLUSG-B04` | `rule-uses-declared` skips a node that is not a spec, and a spec with no rule of the letters asked about. |
| `RLUSG-B05` | `rule-uses-resolve` accepts a field the spec declares — the first cell of a row of any of its other tables, or a backticked name in a heading —, a dotted or indexed name by its first segment, a `DEPn` that is a row of the dependencies table, and leaves codes to `code-reference-valid`; it fails naming each use that resolves to nothing, with its rule. |
| `RLUSG-B06` | `rule-uses-resolve` skips a node that is not a spec, and a spec whose rules do not say what they use yet. |
| `RLUSG-B07` | `rule-uses-implemented` fails naming each field a rule uses that no code file the spec governs mentions as a whole word — the name, or the first or last segment of a dotted one —; codes and `DEPn` are not fields; a spec with no rule uses, or governing no code, is skipped. |

## Errors

none — both gates read only the spec's text; what they cannot find is a verdict, not a failure

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
