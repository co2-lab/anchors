<!-- @anchors
  code: SNGTU
  updated_at: 2026-10-08
  layer: gate
-->
# SingleTestPerUnit — a unit has one test file per test layer

> **Code**: `SNGTU`

## Overview

Two files testing one unit in one layer split what the unit is proven by, and each looks complete
on its own: in the reference app a hook had `useIapNative.test.ts` and `useIapNative.test.tsx`, and
five models had tests both under `__tests__/unit/models` and under `__tests__/unit/lambdas/`. Which
one to extend, which one a gate reads, which one a mutation run counts — each answer was a guess.

The unit a test tests is the project's own derivation read backwards (`TestedUnits`), so nothing
assumes a language or a layout. A split the project means is declared in one of the files with
`@split-test: <why>`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a code file | any other kind | this unit: Skip |
| the map | the scan's map | no map | this unit: Pending, naming `map build` |

## Effects

| Effect | Description |
| --- | --- |
| `SNGTU-B01` | A unit tested by more than one file in the same test layer fails, naming the layer and the files; files in different layers are not a split. |
| `SNGTU-B02` | A `@split-test: <why>` in any of the files of a layer declares the split, and the unit passes for that layer. |
| `SNGTU-B03` | A node that is not code, and a unit no test file tests, are skipped; without a map the gate is pending. |

## Errors

none — the gate reads the map and the test files it names; a file it cannot read declares no split

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
