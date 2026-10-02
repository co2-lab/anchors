<!-- @anchors
  code: TSRCH
  updated_at: 2026-10-02
  layer: gate
-->
# TestReach — a test reaches the unit it says it tests

> **Code**: `TSRCH`

## Overview

A test can carry every code, match every title and assert — over a copy. In the reference app a
test pasted the unit's function into its own file and exercised the paste; another imported a
neighbour; an end-to-end test's `ref:` named one handler and invoked another. The gates that read
titles saw nothing wrong, and the unit could change freely.

Two gates ask it on one routine, apart so the project decides per gate what blocks:

- `test-exercises-unit` — the unit the project's derivation pairs with the test (`TestedUnits`);
- `test-ref-matches-unit` — the unit the test's `ref:` names.

Reaching is read without knowing the language. A test reaches a unit when an import line names
the unit's module, when it names something the unit defines (the dialect's `definition` says what
a definition looks like), or — for the `ref:` — when a call the project declares on the gate
(`invocations`) names it. A test that defines a name its unit defines is exercising a copy.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a test file that is not a support file | any other kind, a support file | this unit: Skip |
| the definition | the dialect's `definition` | none declared | this unit: `test-exercises-unit` skips; `test-ref-matches-unit` reads imports and invocations only |
| the invocations | the gate entry's patterns, each with a capture group | a pattern with no group | the configuration load |
| the unit's file | a readable file | one that cannot be read | this unit: it defines nothing, and only an import or an invocation reaches it |

## Effects

| Effect | Description |
| --- | --- |
| `TSRCH-B01` | `test-exercises-unit`: a test that names, as a whole word, something its unit defines passes; one that neither imports the unit nor names anything it defines fails naming the unit. Every unit the derivation pairs with the test is confronted. |
| `TSRCH-B02` | An import line of the test — by the dialect's `import_pattern` when declared, or the forms most languages share — that names the unit's file name without extension, or its directory as a path, reaches the unit; the name elsewhere in the test, or a unit at the root with no directory, does not count as an import. (`importsModule`) |
| `TSRCH-B03` | A test that defines a name its unit defines fails naming the unit and the names, sorted — even when it reaches the unit otherwise: it is exercising its own copy. A name it copies does not count as reaching. |
| `TSRCH-B04` | A defined name shorter than three characters is anyone's, and neither reaches nor copies; a definition's first non-empty capture group is its name. (`definedNames`) |
| `TSRCH-B05` | `test-ref-matches-unit`: a test that reaches any of the code files the `ref:`'s spec governs passes; one that reaches none fails naming the `ref:` and the files, sorted. |
| `TSRCH-B06` | A declared invocation whose first non-empty capture is the unit's file name without extension, or one of its directories, reaches the unit; another capture does not. (`invokes`) |
| `TSRCH-B07` | A node that is not a test or is a support file is skipped by both gates; `test-exercises-unit` skips a test no unit pairs with and a project with no definition; `test-ref-matches-unit` skips a test with no `ref:` and a `ref:` whose spec governs no code. |
| `TSRCH-B09` | A definition is identified with its owner when the dialect's `definition` captures one (a Go method by its receiver type, `GormPinger.Ping`), and a definition indented under a type with no owner captured is a member; a test that defines a member of its own type — a fake that implements the unit's interface — copies nothing, and only the same name of the same owner defined again is a copy. A member's bare name still reaches the unit. (`definedNames`) |
| `TSRCH-B08` | A test that carries `@no-unit-import: <why>` passes `test-exercises-unit`, copies included; a waiver with no reason waives nothing. |

## Errors

none — an unreadable unit defines nothing, and the gates read what the map and the files say

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/gate/reverse_match.go` | `testedUnits` | gate — the derivation read backwards, cached per map |
| DEP2 | `internal/gate/proof_crosses_boundary.go` | `importLines` | gate — the lines that import something |
| DEP3 | `internal/gate/proof_crosses_boundary.go` | `unitFiles` | gate — a code's spec to the code files it governs |
| DEP4 | `internal/gate/ref_resolves.go` | `refHeaderRE` | gate — the `ref:` header |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
