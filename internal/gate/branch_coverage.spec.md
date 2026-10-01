<!-- @anchors
  code: BRCOV
  updated_at: 2026-09-30
  layer: gate
-->
# BranchCoverage — the tests take the branches the code has

> **Code**: `BRCOV`

## Overview

Line coverage says a line ran. A line with a condition runs whichever way the condition goes,
and the other way can stay untested for good: in the reference app a "red" state was unreachable
(the slider's maximum was below the base it had to cross), a percentage was always zero, an
empty-state skeleton never rendered — every line covered, every branch not.

The lcov format carries branches (`BRDA`) whatever the language that wrote it, and the coverage
ingestion keeps them on the node: how many, and the ones no suite took. A branch never taken on a
line where the mutation run also found mutants no test ran is reported apart as likely dead —
nothing the tests do reaches it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the node | a code file | any other kind | this unit: Skip |
| the coverage signal | an ingested lcov at the node's revision | none, or measured at another revision | this unit: Pending |
| the branches | the lcov's `BRDA` entries | a report with none (Go's profile, for one) | this unit: Skip, nothing to measure |
| the floor | the gate entry's `min_percent`, 0 to 100 | — | the configuration load |

## Effects

| Effect | Description |
| --- | --- |
| `BRCOV-B01` | A file whose share of branches taken is below the floor fails, naming the share, the floor, the branches missed out of the total and their lines in order; at or above the floor it passes. With no `min_percent` the floor is 100: every branch. |
| `BRCOV-B02` | A branch on a line carrying `@no-branch: <why>`, or on the line after it, is left out of the missed; a waiver with no reason waives nothing. |
| `BRCOV-B03` | A missed branch on a line where the mutation, measured at the node's revision, found a mutant no test ran fails as likely dead, naming the lines — whatever the floor. A mutation measured at another revision says nothing about these lines. |
| `BRCOV-B04` | A node that is not code is skipped; a node with no coverage, or coverage of another revision, is pending; a coverage with no branch is skipped. |
| `BRCOV-B05` | A file a coverage report listed with no instrumentable line is skipped, and one a whole run of its suite left out of the report is a divergence, as `line-coverage` reads them. (`coverageAbsence`) |

## Errors

none — the gate reads the signal the ingestion wrote and the file's text

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/ingest.go` | `NoCoverageLines` | mapx — the lines of the mutants no test ran |
| DEP2 | `internal/gate/rule_uses.go` | `gateEntry` | gate — the gate's own settings |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
